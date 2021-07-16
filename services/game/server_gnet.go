package game

import (
	"bytes"
	"context"
	"encoding/binary"
	"hash/crc32"
	"os"
	"strings"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/transport/codec"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.uber.org/atomic"
)

type GNetServer struct {
	tr    transport.Transport
	reg   transport.Register
	codec codec.Marshaler
	g     *Game
	conns map[gnet.Conn]*gnetTransportSocket
	sync.RWMutex

	*gnet.EventServer
	pool *ants.Pool
}

func NewGNetServer(ctx *cli.Context, g *Game) *GNetServer {
	maxAccount := ctx.Int("account_connect_max")

	s := &GNetServer{
		g:     g,
		reg:   g.msgRegister.r,
		codec: &codec.ProtoBufMarshaler{},
		conns: make(map[gnet.Conn]*gnetTransportSocket, maxAccount),
	}

	var err error
	s.pool, err = ants.NewPool(maxAccount, ants.WithExpiryDuration(10*time.Second))
	if !utils.ErrCheck(err, "NewPool failed when NewGNetServer", maxAccount) {
		return nil
	}

	err = s.serve(ctx)
	if !utils.ErrCheck(err, "serve failed when NewGNetServer", maxAccount) {
		return nil
	}
	return s
}

func (s *GNetServer) serve(ctx *cli.Context) error {
	tcpAddr := strings.Join([]string{"tcp", ctx.String("tcp_listen_addr")}, "://")
	s.tr = transport.NewTransport("gnet")

	err := s.tr.Init(
		transport.Timeout(transport.DefaultServeTimeout),
		transport.Codec(&codec.ProtoBufMarshaler{}),
	)

	if err != nil {
		return err
	}

	go func() {
		defer utils.CaptureException()

		err := s.tr.ListenAndServe(ctx.Context, tcpAddr, s)
		if err != nil {
			log.Warn().Err(err).Msg("gnet server ListenAndServe return with error")
			os.Exit(1)
		}
	}()

	log.Info().
		Str("addr", tcpAddr).
		Msg("gnet server serve at address")

	return nil
}

func (s *GNetServer) HandleSocket(context.Context, transport.Socket) {
}

// OnOpened fires when a new connection has been opened.
// The parameter:c has information about the connection such as it's local and remote address.
// Parameter:out is the return value which is going to be sent back to the client.
// It is generally not recommended to send large amounts of data back to the client in OnOpened.
//
// Note that the bytes returned by OnOpened will be sent back to client without being encoded.
func (s *GNetServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	sock := &gnetTransportSocket{
		Conn:   c,
		closed: *atomic.NewBool(false),
		codec:  &codec.ProtoBufMarshaler{},
	}

	sock.ctx, sock.cancel = context.WithCancel(context.Background())

	s.Lock()
	s.conns[c] = sock
	s.Unlock()

	return
}

// OnClosed fires when a connection has been closed.
// The parameter:err is the last known connection error.
func (s *GNetServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	s.Lock()
	sock, ok := s.conns[c]
	delete(s.conns, c)
	s.Unlock()

	if ok {
		sock.closed.Store(true)
		sock.cancel()
	}

	// log.Info().Err(err).Msg("gnet OnClosed")
	return
}

// Tick fires immediately after the server starts and will fire again
// following the duration specified by the delay return value.
func (s *GNetServer) Tick() (delay time.Duration, action gnet.Action) {
	delay = time.Second
	return
}

// React fires when a connection sends the server data.
// Call c.Read() or c.ReadN(n) within the parameter:c to read incoming data from client.
// Parameter:out is the return value which is going to be sent back to the client.
func (s *GNetServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if len(frame) < 4 {
		return
	}

	nameCrc := binary.LittleEndian.Uint32(frame[:4])
	h, err := s.reg.GetHandler(nameCrc)
	if err != nil {
		return nil, gnet.Close
	}

	var message transport.Message
	message.Name = h.Name
	message.Body, err = s.codec.Unmarshal(frame[4:], h.RType)
	if err != nil {
		return nil, gnet.Close
	}

	s.RLock()
	defer s.RUnlock()
	sock, ok := s.conns[c]
	if !ok {
		return nil, gnet.Close
	}
	// err = s.pool.Submit(func() {
	err = h.Fn(sock.ctx, sock, &message)
	// })
	utils.ErrPrint(err, "gnet Submit failed when GNetServer.React")

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

type gnetTransportSocket struct {
	gnet.Conn
	codec  codec.Marshaler
	closed atomic.Bool
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *gnetTransportSocket) Recv(transport.Register) (*transport.Message, *transport.MessageHandler, error) {
	return nil, nil, nil
}

func (s *gnetTransportSocket) Send(m *transport.Message) error {
	body, err := s.codec.Marshal(m.Body)
	if err != nil {
		return err
	}

	// Message Header:
	// 4 bytes message size, size = 4 bytes name crc + proto binary size
	// 4 bytes message name crc32 id,
	// Message Body:
	var bodySize uint32 = uint32(4 + len(body))
	var nameCrc uint32 = crc32.ChecksumIEEE([]byte(m.Name))
	buffer := new(bytes.Buffer)

	_ = binary.Write(buffer, binary.LittleEndian, bodySize)
	_ = binary.Write(buffer, binary.LittleEndian, uint32(nameCrc))
	_, _ = buffer.Write(body)

	return s.AsyncWrite(buffer.Bytes())
}

func (s *gnetTransportSocket) Close() {
	if s.closed.Load() {
		return
	}

	s.closed.Store(true)
	err := s.Conn.Close()
	utils.ErrPrint(err, "gnet.Conn Close failed")
}

func (s *gnetTransportSocket) IsClosed() bool {
	return s.closed.Load()
}

func (s *gnetTransportSocket) Local() string {
	if s.closed.Load() {
		return ""
	}
	return s.Conn.LocalAddr().String()
}

func (s *gnetTransportSocket) Remote() string {
	if s.closed.Load() {
		return ""
	}
	return s.Conn.RemoteAddr().String()
}
