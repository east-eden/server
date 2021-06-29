package game

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"e.coding.net/mmstudio/blade/server/services/game/player"
	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/transport/codec"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	GNetRecvMaxSize        uint32 = 1024 * 1024
	ErrGNetReadFail               = errors.New("gnet read failed")
	ErrGNetReadLengthLimit        = errors.New("gnet read length limit")
)

type GNetCodec struct{}

// Encode encodes frames upon server responses into TCP stream.
func (codec *GNetCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	log.Info().Bytes("buf", buf).Msg("gnet encoding...")
	return buf, nil
}

// Decode decodes frames from TCP stream via specific implementation.
func (codec *GNetCodec) Decode(c gnet.Conn) ([]byte, error) {
	log.Info().Msg("gnet decoding...")
	size, sizeBuf := c.ReadN(4)
	if size != 4 {
		c.ResetBuffer()
		return sizeBuf, ErrGNetReadFail
	}

	msgLen := binary.LittleEndian.Uint32(sizeBuf)
	if msgLen > GNetRecvMaxSize {
		c.ResetBuffer()
		return sizeBuf, ErrGNetReadLengthLimit
	}

	shiftN := c.ShiftN(4)
	if shiftN != 4 {
		c.ResetBuffer()
		return sizeBuf, ErrGNetReadFail
	}

	sz, bodyBuf := c.ReadN(int(msgLen))
	if sz != int(msgLen) {
		c.ResetBuffer()
		return bodyBuf, ErrGNetReadFail
	}

	return bodyBuf, nil
}

type GNetServer struct {
	reg transport.Register
	g   *Game

	*gnet.EventServer
	pool *goroutine.Pool
}

func NewGNetServer(ctx *cli.Context, g *Game) *GNetServer {
	s := &GNetServer{
		g:   g,
		reg: g.msgRegister.r,

		pool: goroutine.Default(),
	}

	// s.wp = workerpool.New(ctx.Int("account_connect_max"))
	err := s.serve(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("tcpserver serve return error")
	}

	return s
}

func (s *GNetServer) serve(ctx *cli.Context) error {
	tcpAddr := strings.Join([]string{"tcp", ctx.String("tcp_listen_addr")}, "://")
	go func() {
		err := gnet.Serve(
			s,
			tcpAddr,
			gnet.WithCodec(&GNetCodec{}),
			gnet.WithTicker(true),
		)
		_ = utils.ErrCheck(err, "gnet.Serve failed", tcpAddr)
	}()

	log.Info().Str("addr", tcpAddr).Msg("tcp server serve at address")
	return nil
}

// OnInitComplete fires when the server is ready for accepting connections.
// The parameter:server has information and various utilities.
func (s *GNetServer) OnInitComplete(server gnet.Server) (action gnet.Action) {
	log.Info().Msg("gnet OnInitComplete")
	return
}

// OnShutdown fires when the server is being shut down, it is called right after
// all event-loops and connections are closed.
func (s *GNetServer) OnShutdown(server gnet.Server) {
	log.Info().Msg("gnet OnShutdonw")
}

// OnOpened fires when a new connection has been opened.
// The parameter:c has information about the connection such as it's local and remote address.
// Parameter:out is the return value which is going to be sent back to the client.
// It is generally not recommended to send large amounts of data back to the client in OnOpened.
//
// Note that the bytes returned by OnOpened will be sent back to client without being encoded.
func (s *GNetServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Info().Msg("gnet OnOpened")
	return
}

// OnClosed fires when a connection has been closed.
// The parameter:err is the last known connection error.
func (s *GNetServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Info().Msg("gnet OnClosed")
	return
}

// PreWrite fires just before any data is written to any client socket, this event function is usually used to
// put some code of logging/counting/reporting or any prepositive operations before writing data to client.
func (s *GNetServer) PreWrite() {

}

// React fires when a connection sends the server data.
// Call c.Read() or c.ReadN(n) within the parameter:c to read incoming data from client.
// Parameter:out is the return value which is going to be sent back to the client.
func (s *GNetServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	nameCrc := binary.LittleEndian.Uint32(frame[:4])
	h, err := s.reg.GetHandler(nameCrc)
	if err != nil {
		return nil, gnet.Close
	}

	var message transport.Message
	message.Name = h.Name
	codec := &codec.ProtoBufMarshaler{}
	message.Body, err = codec.Unmarshal(frame[4:], h.RType)
	if err != nil {
		return nil, gnet.Close
	}

	err := h.Fn(context.Background(), sock, &message)

	log.Info().Bytes("frame", frame).Interface("message", message).Msg("gnet React")
	return
}

// Tick fires immediately after the server starts and will fire again
// following the duration specified by the delay return value.
func (s *GNetServer) Tick() (delay time.Duration, action gnet.Action) {
	log.Info().Msg("gnet Tick")
	delay = time.Second * 5
	return
}

func (s *GNetServer) Run(ctx context.Context) error {
	<-ctx.Done()
	log.Info().Msg("tcp server context done...")
	return nil
}

func (s *GNetServer) Exit() {
	log.Info().Msg("tcp server exit...")
}

func (s *GNetServer) handleSocket(ctx context.Context, sock transport.Socket) {
	subCtx, cancel := context.WithCancel(ctx)
	_ = s.pool.Submit(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			sock.Close()
			cancel()
		}()

		for {
			ct := time.Now()

			select {
			case <-subCtx.Done():
				return
			default:
			}

			msg, h, err := sock.Recv(s.reg)
			if err != nil {
				if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
					log.Warn().Err(err).Msg("TcpServer.handleSocket error")
				}
				return
			}

			if err := h.Fn(subCtx, sock, msg); err != nil {
				// account need disconnect
				if errors.Is(err, player.ErrAccountDisconnect) {
					log.Info().Msg("TcpServer.handleSocket account disconnect initiativly")
					return
				}

				log.Warn().Caller().Err(err).Str("msg", msg.Name).Msg("TcpServer.handleSocket callback error")
			}

			time.Sleep(tpcRecvInterval - time.Since(ct))
		}
	})
}
