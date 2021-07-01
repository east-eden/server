package game

import (
	"context"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"strings"
	"sync"

	"e.coding.net/mmstudio/blade/server/transport"
	"e.coding.net/mmstudio/blade/server/transport/codec"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/valyala/bytebufferpool"
	"go.uber.org/atomic"
)

var (
	GNetRecvMaxSize        uint32 = 1024 * 1024
	ErrGNetReadFail               = errors.New("gnet read failed")
	ErrGNetReadLengthLimit        = errors.New("gnet read length limit")
)

type GNetCodec struct{}

// Encode encodes frames upon server responses into TCP stream.
func (codec *GNetCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return buf, nil
}

// Decode decodes frames from TCP stream via specific implementation.
func (codec *GNetCodec) Decode(c gnet.Conn) ([]byte, error) {
	bufLen := c.BufferLength()
	log.Info().Int("buffer_Length", bufLen).Msg("decode...")
	if bufLen <= 0 {
		return nil, nil
	}

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

	if c.ShiftN(4) != 4 {
		c.ResetBuffer()
		return sizeBuf, ErrGNetReadFail
	}

	bodySize, bodyBuf := c.ReadN(int(msgLen))
	if bodySize != int(msgLen) {
		c.ResetBuffer()
		return bodyBuf, ErrGNetReadFail
	}

	if c.ShiftN(bodySize) != bodySize {
		c.ResetBuffer()
		return bodyBuf, ErrGNetReadFail
	}

	return bodyBuf, nil
}

type GNetSocket struct {
	gnet.Conn
	transport.Socket
	closed atomic.Bool
}

func (c *GNetSocket) Recv(transport.Register) (*transport.Message, *transport.MessageHandler, error) {
	return nil, nil, nil
}

func (c *GNetSocket) Send(m *transport.Message) error {
	body, err := c.PbMarshaler().Marshal(m.Body)
	if err != nil {
		return err
	}

	// Message Header:
	// 4 bytes message size, size = 4 bytes name crc + proto binary size
	// 4 bytes message name crc32 id,
	// Message Body:
	var bodySize uint32 = uint32(4 + len(body))
	var nameCrc uint32 = crc32.ChecksumIEEE([]byte(m.Name))
	buffer := bytebufferpool.Get()
	defer bytebufferpool.Put(buffer)

	_ = binary.Write(buffer, binary.LittleEndian, bodySize)
	_ = binary.Write(buffer, binary.LittleEndian, uint32(nameCrc))
	_, _ = buffer.Write(body)

	// todo add a writer buffer, cache bytes which didn't sended, then try resend
	return c.AsyncWrite(buffer.Bytes())
}

func (c *GNetSocket) Close() {
	if c.closed.Load() {
		return
	}

	c.closed.Store(true)
	err := c.Conn.Close()
	utils.ErrPrint(err, "gnet.Conn Close failed", c.Local(), c.Remote())
}

func (c *GNetSocket) IsClosed() bool {
	return c.closed.Load()
}

func (c *GNetSocket) Local() string {
	return c.Conn.LocalAddr().String()
}

func (c *GNetSocket) Remote() string {
	return c.Conn.RemoteAddr().String()
}

func (c *GNetSocket) PbMarshaler() codec.Marshaler {
	return &codec.ProtoBufMarshaler{}
}

type GNetServer struct {
	reg   transport.Register
	g     *Game
	conns map[gnet.Conn]transport.Socket
	sync.RWMutex

	*gnet.EventServer
	pool *goroutine.Pool
}

func NewGNetServer(ctx *cli.Context, g *Game) *GNetServer {
	maxConns := ctx.Int("account_connect_max")
	s := &GNetServer{
		g:     g,
		reg:   g.msgRegister.r,
		conns: make(map[gnet.Conn]transport.Socket, maxConns),
		pool:  goroutine.Default(),
	}

	s.serve(ctx)
	return s
}

func (s *GNetServer) serve(ctx *cli.Context) {
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
	sock := &GNetSocket{
		Conn:   c,
		closed: *atomic.NewBool(false),
	}

	s.Lock()
	s.conns[c] = sock
	s.Unlock()

	log.Info().Msg("gnet OnOpened")
	return
}

// OnClosed fires when a connection has been closed.
// The parameter:err is the last known connection error.
func (s *GNetServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	s.Lock()
	delete(s.conns, c)
	s.Unlock()

	log.Info().Err(err).Msg("gnet OnClosed")
	return
}

// React fires when a connection sends the server data.
// Call c.Read() or c.ReadN(n) within the parameter:c to read incoming data from client.
// Parameter:out is the return value which is going to be sent back to the client.
func (s *GNetServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if len(frame) < 4 {
		return
	}

	s.RLock()
	sock, ok := s.conns[c]
	s.RUnlock()

	if !ok {
		return nil, gnet.Close
	}

	nameCrc := binary.LittleEndian.Uint32(frame[:4])
	h, err := s.reg.GetHandler(nameCrc)
	if err != nil {
		return nil, gnet.Close
	}

	var message transport.Message
	message.Name = h.Name
	message.Body, err = sock.PbMarshaler().Unmarshal(frame[4:], h.RType)
	if err != nil {
		return nil, gnet.Close
	}

	s.pool.Submit(func() {
		err = h.Fn(context.Background(), sock, &message)
	})

	log.Info().Bytes("frame", frame).Interface("message", message).Msg("gnet React")
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
