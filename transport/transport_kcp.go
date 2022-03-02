// Package kcp provides a KCP transport
package transport

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"time"

	maddr "github.com/asim/go-micro/v3/util/addr"
	mnet "github.com/asim/go-micro/v3/util/net"
	mls "github.com/asim/go-micro/v3/util/tls"
	"github.com/valyala/bytebufferpool"
	"github.com/xtaci/kcp-go"
	"go.uber.org/atomic"
	"google.golang.org/protobuf/proto"

	"github.com/east-eden/server/transport/codec"
	"github.com/east-eden/server/utils/writer"
)

func newKcpTransportSocket() *kcpTransportSocket {
	return &kcpTransportSocket{
		codecs: []codec.Marshaler{&codec.ProtoBufMarshaler{}, &codec.JsonMarshaler{}},
	}
}

type kcpTransport struct {
	opts *Options
}

func (t *kcpTransport) Init(opts ...Option) {
	t.opts = DefaultTransportOptions()

	for _, o := range opts {
		o(t.opts)
	}
}

func (t *kcpTransport) Options() *Options {
	return t.opts
}

func (t *kcpTransport) Protocol() string {
	return "kcp"
}

func (t *kcpTransport) Dial(addr string, opts ...DialOption) (Socket, error) {
	dopts := DialOptions{
		Timeout: DefaultDialTimeout,
	}

	for _, opt := range opts {
		opt(&dopts)
	}

	var conn net.Conn
	var err error

	// TODO: support dial option here rather than using internal config
	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig
		if config == nil {
			config = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		conn, err = kcp.DialWithOptions(addr, nil, 10, 3)
	} else {
		conn, err = kcp.DialWithOptions(addr, nil, 10, 3)
	}

	if err != nil {
		return nil, err
	}

	return &kcpTransportSocket{
		conn:    conn,
		writer:  writer.NewBinaryWriter(bufio.NewWriterSize(conn, writer.DefaultBinaryWriterSize), -1),
		reader:  bufio.NewReader(conn),
		codecs:  []codec.Marshaler{&codec.ProtoBufMarshaler{}, &codec.JsonMarshaler{}},
		timeout: t.opts.Timeout,
		closed:  *atomic.NewBool(false),
	}, nil
}

func (t *kcpTransport) ListenAndServe(ctx context.Context, addr string, server TransportServer, opts ...ListenOption) error {
	l, err := t.Listen(addr, opts...)
	if err != nil {
		return err
	}

	defer l.Close()
	return l.Accept(ctx, server)
}

func (t *kcpTransport) Listen(addr string, opts ...ListenOption) (Listener, error) {
	options := DefaultListenOptions()
	for _, o := range opts {
		o(options)
	}

	var l net.Listener
	var err error

	// TODO: support use of listen options
	if t.opts.Secure || t.opts.TLSConfig != nil {
		config := t.opts.TLSConfig

		fn := func(addr string) (net.Listener, error) {
			if config == nil {
				hosts := []string{addr}

				// check if its a valid host:port
				if host, _, err := net.SplitHostPort(addr); err == nil {
					if len(host) == 0 {
						hosts = maddr.IPs()
					} else {
						hosts = []string{host}
					}
				}

				// generate a certificate
				cert, err := mls.Certificate(hosts...)
				if err != nil {
					return nil, err
				}
				config = &tls.Config{Certificates: []tls.Certificate{cert}}
			}
			return kcp.ListenWithOptions(addr, nil, 10, 3)
		}

		l, err = mnet.Listen(addr, fn)
	} else {
		fn := func(addr string) (net.Listener, error) {
			return kcp.ListenWithOptions(addr, nil, 10, 3)
		}

		l, err = mnet.Listen(addr, fn)
	}

	if err != nil {
		return nil, err
	}

	ls := &kcpTransportListener{
		opts:     options,
		timeout:  t.opts.Timeout,
		listener: l,
	}

	// ls.sockPool.New = newTcpTransportSocket

	return ls, nil
}

type kcpTransportListener struct {
	opts     *ListenOptions
	listener net.Listener
	timeout  time.Duration
	// sockPool sync.Pool
}

func (t *kcpTransportListener) Addr() string {
	return t.listener.Addr().String()
}

func (t *kcpTransportListener) Close() error {
	return t.listener.Close()
}

func (t *kcpTransportListener) Accept(ctx context.Context, server TransportServer) error {
	var tempDelay time.Duration

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		c, err := t.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				fmt.Printf("tcp: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		// sock := t.sockPool.Get().(*tcpTransportSocket)
		sock := newKcpTransportSocket()
		sock.conn = c
		sock.reader = bufio.NewReader(sock.conn)
		sock.writer = writer.NewBinaryWriter(
			bufio.NewWriterSize(sock.conn, writer.DefaultBinaryWriterSize),
			t.opts.WriterLatency,
		)
		sock.timeout = t.timeout
		sock.closed.Store(false)

		server.HandleSocket(ctx, sock)
	}
}

type kcpTransportSocket struct {
	conn    net.Conn
	writer  writer.BinaryWriter
	reader  *bufio.Reader
	codecs  []codec.Marshaler
	timeout time.Duration
	closed  atomic.Bool
}

func (t *kcpTransportSocket) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *kcpTransportSocket) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *kcpTransportSocket) Close() {
	if t.closed.Load() {
		return
	}

	t.writer.Stop()
	t.closed.Store(true)
	_ = t.conn.Close()
}

func (t *kcpTransportSocket) IsClosed() bool {
	return t.closed.Load()
}

func (t *kcpTransportSocket) PbMarshaler() codec.Marshaler {
	return t.codecs[0]
}

func (t *kcpTransportSocket) Recv(r Register) (proto.Message, *MessageHandler, error) {
	if t.IsClosed() {
		return nil, nil, errors.New("tcpTransportSocket.Recv failed: socket closed")
	}

	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		if err := t.conn.SetDeadline(time.Now().Add(t.timeout)); err != nil {
			return nil, nil, err
		}
	}

	// Message Header:
	// 4 bytes message size, size = 4 bytes name crc + proto binary size
	// 4 bytes message name crc32 id,
	// Message Body:
	var header [8]byte
	if _, err := io.ReadFull(t.reader, header[:]); err != nil {
		return nil, nil, fmt.Errorf("tcpTransportSocket.Recv header failed: %w", err)
	}

	var msgLen uint32
	var nameCrc uint32
	msgLen = binary.LittleEndian.Uint32(header[:4])
	nameCrc = binary.LittleEndian.Uint32(header[4:8])

	if msgLen > TcpPacketMaxSize {
		return nil, nil, ErrTransportReadSizeTooLong
	}

	// read body bytes
	bodyData := make([]byte, msgLen-4)
	if _, err := io.ReadFull(t.reader, bodyData); err != nil {
		return nil, nil, fmt.Errorf("tcpTransportSocket.Recv body failed: %w", err)
	}

	// get register handler
	h, err := r.GetHandler(nameCrc)
	if err != nil {
		return nil, nil, fmt.Errorf("tcpTransportSocket.Recv failed: %w", err)
	}

	message, err := t.codecs[0].Unmarshal(bodyData, h.RType)
	if err != nil {
		return nil, nil, fmt.Errorf("tcpTransportSocket.Recv unmarshal message body failed: %w", err)
	}

	return message.(proto.Message), h, err
}

func (t *kcpTransportSocket) Send(m proto.Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		if err := t.conn.SetDeadline(time.Now().Add(t.timeout)); err != nil {
			return err
		}
	}

	body, err := t.codecs[0].Marshal(m)
	if err != nil {
		return err
	}

	// Message Header:
	// 4 bytes message size, size = 4 bytes name crc + proto binary size
	// 4 bytes message name crc32 id,
	// Message Body:
	var bodySize uint32 = uint32(4 + len(body))
	name := m.ProtoReflect().Descriptor().Name()
	var nameCrc uint32 = crc32.ChecksumIEEE([]byte(name))
	buffer := bytebufferpool.Get()
	defer bytebufferpool.Put(buffer)

	_ = binary.Write(buffer, binary.LittleEndian, bodySize)
	_ = binary.Write(buffer, binary.LittleEndian, uint32(nameCrc))
	_, _ = buffer.Write(body)

	// todo add a writer buffer, cache bytes which didn't sended, then try resend
	_, err = t.writer.Write(buffer.Bytes())
	return err
}

// io.Writer
func (t *kcpTransportSocket) Write(body []byte) (int, error) {
	_ = t.conn.SetWriteDeadline(time.Now().Add(t.timeout))
	return t.writer.Write(body)
}

// io.Reader
func (t *kcpTransportSocket) Read(body []byte) (int, error) {
	_ = t.conn.SetReadDeadline(time.Now().Add(t.timeout))
	return t.reader.Read(body)
}
