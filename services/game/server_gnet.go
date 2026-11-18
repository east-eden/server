package game

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/east-eden/server/transport"
	"github.com/east-eden/server/transport/codec"
	"github.com/east-eden/server/utils"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.uber.org/atomic"
	"google.golang.org/protobuf/proto"
)

// gnet frame layout (little endian):
//
//	| 4 bytes payload size | 4 bytes message name crc32 | payload |
//
// payload size = 4 + len(proto body)
const gnetFrameHeaderSize = 8

type GNetServer struct {
	tr    transport.Transport
	reg   transport.Register
	codec codec.Marshaler
	g     *Game
	conns map[gnet.Conn]*gnetTransportSocket
	sync.RWMutex

	*gnet.BuiltinEventEngine
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

	s.tr.Init(
		transport.Timeout(transport.DefaultServeTimeout),
	)

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

// OnOpen fires when a new connection has been opened.
func (s *GNetServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
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

// OnClose fires when a connection has been closed.
func (s *GNetServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	s.Lock()
	sock, ok := s.conns[c]
	delete(s.conns, c)
	s.Unlock()

	if ok {
		sock.closed.Store(true)
		sock.cancel()
	}

	return
}

// OnTick fires immediately after the server starts and will fire again
// following the duration specified by the delay return value.
func (s *GNetServer) OnTick() (delay time.Duration, action gnet.Action) {
	delay = time.Second
	return
}

// OnTraffic fires when a connection sends the server data.
//
// gnet/v2 has no codec hook, so the length-prefixed framing is decoded here.
func (s *GNetServer) OnTraffic(c gnet.Conn) gnet.Action {
	for {
		head, err := c.Peek(4)
		if err != nil {
			// not enough bytes for the payload size, wait for the next OnTraffic
			return gnet.None
		}

		payloadSize := binary.LittleEndian.Uint32(head)
		if payloadSize < 4 || payloadSize > transport.GNetRecvMaxSize {
			log.Warn().
				Uint32("payload_size", payloadSize).
				Str("remote", c.RemoteAddr().String()).
				Msg("gnet frame size invalid, close connection")
			return gnet.Close
		}

		total := 4 + int(payloadSize)
		if c.InboundBuffered() < total {
			return gnet.None
		}

		frame, err := c.Next(total)
		if err != nil {
			return gnet.Close
		}

		if action := s.handleFrame(c, frame[4:]); action != gnet.None {
			return action
		}
	}
}

func (s *GNetServer) handleFrame(c gnet.Conn, payload []byte) gnet.Action {
	if len(payload) < 4 {
		return gnet.Close
	}

	nameCrc := binary.LittleEndian.Uint32(payload[:4])
	h, err := s.reg.GetHandler(nameCrc)
	if err != nil {
		log.Warn().Err(err).Msg("gnet get handler failed")
		return gnet.Close
	}

	message, err := s.codec.Unmarshal(payload[4:], h.RType)
	if err != nil {
		log.Warn().Err(err).Msg("gnet unmarshal message failed")
		return gnet.Close
	}

	s.RLock()
	sock, ok := s.conns[c]
	s.RUnlock()
	if !ok {
		return gnet.Close
	}

	if err := h.Fn(sock.ctx, sock, message.(proto.Message)); err != nil && !errors.Is(err, io.EOF) {
		utils.ErrPrint(err, "gnet handle message failed")
	}

	return gnet.None
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

func (s *gnetTransportSocket) Recv(transport.Register) (proto.Message, *transport.MessageHandler, error) {
	return nil, nil, nil
}

func (s *gnetTransportSocket) Send(m proto.Message) error {
	body, err := s.codec.Marshal(m)
	if err != nil {
		return err
	}

	// Message Header:
	// 4 bytes message size, size = 4 bytes name crc + proto binary size
	// 4 bytes message name crc32 id,
	// Message Body:
	bodySize := uint32(4 + len(body))
	name := m.ProtoReflect().Descriptor().Name()
	nameCrc := crc32.ChecksumIEEE([]byte(name))
	buffer := new(bytes.Buffer)

	_ = binary.Write(buffer, binary.LittleEndian, bodySize)
	_ = binary.Write(buffer, binary.LittleEndian, nameCrc)
	_, _ = buffer.Write(body)

	return s.Conn.AsyncWrite(buffer.Bytes(), nil)
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

func (s *gnetTransportSocket) Write(body []byte) (int, error) {
	err := s.Conn.AsyncWrite(body, nil)
	return len(body), err
}

// gnet does not expose a blocking Read, inbound data is delivered by OnTraffic.
func (s *gnetTransportSocket) Read(body []byte) (int, error) {
	return 0, nil
}
