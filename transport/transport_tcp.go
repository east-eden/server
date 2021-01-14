// Package tcp provides a TCP transport
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
	"sync"
	"time"

	maddr "github.com/micro/go-micro/v2/util/addr"
	mnet "github.com/micro/go-micro/v2/util/net"
	mls "github.com/micro/go-micro/v2/util/tls"

	"e.coding.net/mmstudio/blade/server/transport/codec"
	"e.coding.net/mmstudio/blade/server/transport/writer"
)

var (
	tcpRecvBufMax = 1024 * 1024 * 2 // tcp recv buf
)

func newTcpTransportSocket() interface{} {
	return &tcpTransportSocket{
		codecs: []codec.Marshaler{&codec.ProtoBufMarshaler{}, &codec.JsonMarshaler{}},
	}
}

type tcpTransport struct {
	opts Options
}

func (t *tcpTransport) Init(opts ...Option) error {
	for _, o := range opts {
		o(&t.opts)
	}

	return nil
}

func (t *tcpTransport) Options() Options {
	return t.opts
}

func (t *tcpTransport) Protocol() string {
	return "tcp"
}

func (t *tcpTransport) Dial(addr string, opts ...DialOption) (Socket, error) {
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
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: dopts.Timeout}, "tcp", addr, config)
	} else {
		conn, err = net.DialTimeout("tcp", addr, dopts.Timeout)
	}

	if err != nil {
		return nil, err
	}

	return &tcpTransportSocket{
		conn:          conn,
		writer:        writer.NewWriter(bufio.NewWriterSize(conn, writer.DefaultWriterSize), -1),
		reader:        bufio.NewReader(conn),
		codecs:        []codec.Marshaler{&codec.ProtoBufMarshaler{}, &codec.JsonMarshaler{}},
		timeout:       t.opts.Timeout,
		evictedHandle: []func(Socket){},
		closed:        false,
	}, nil
}

func (t *tcpTransport) ListenAndServe(ctx context.Context, addr string, handler TransportHandler, opts ...ListenOption) error {
	l, err := t.Listen(addr, opts...)
	if err != nil {
		return err
	}

	defer l.Close()
	return l.Accept(ctx, handler)
}

func (t *tcpTransport) Listen(addr string, opts ...ListenOption) (Listener, error) {
	var options ListenOptions
	for _, o := range opts {
		o(&options)
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
			return tls.Listen("tcp", addr, config)
		}

		l, err = mnet.Listen(addr, fn)
	} else {
		fn := func(addr string) (net.Listener, error) {
			return net.Listen("tcp", addr)
		}

		l, err = mnet.Listen(addr, fn)
	}

	if err != nil {
		return nil, err
	}

	ls := &tcpTransportListener{
		timeout:  t.opts.Timeout,
		listener: l,
	}

	ls.sockPool.New = newTcpTransportSocket

	return ls, nil
}

type tcpTransportListener struct {
	listener net.Listener
	timeout  time.Duration
	sockPool sync.Pool
}

func (t *tcpTransportListener) Addr() string {
	return t.listener.Addr().String()
}

func (t *tcpTransportListener) Close() error {
	return t.listener.Close()
}

func (t *tcpTransportListener) Accept(ctx context.Context, fn TransportHandler) error {
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

		sock := t.sockPool.Get().(*tcpTransportSocket)
		sock.conn = c
		sock.reader = bufio.NewReader(sock.conn)
		sock.writer = writer.NewWriter(bufio.NewWriterSize(sock.conn, writer.DefaultWriterSize), writer.DefaultWriterLatency)
		sock.timeout = t.timeout
		sock.evictedHandle = []func(Socket){}
		sock.closed = false

		// callback with exit func
		subCtx, cancel := context.WithCancel(ctx)
		fn(subCtx, sock, func() {
			cancel()
			sock.Close()
			t.sockPool.Put(sock)
		})
	}
}

type tcpTransportSocket struct {
	conn          net.Conn
	writer        writer.Writer
	reader        *bufio.Reader
	codecs        []codec.Marshaler
	timeout       time.Duration
	evictedHandle []func(Socket)
	closed        bool
}

func (t *tcpTransportSocket) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *tcpTransportSocket) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *tcpTransportSocket) AddEvictedHandle(f func(Socket)) {
	t.evictedHandle = append(t.evictedHandle, f)
}

func (t *tcpTransportSocket) Close() error {
	for _, handle := range t.evictedHandle {
		handle(t)
	}

	t.writer.Stop()
	t.closed = true
	return t.conn.Close()
}

func (t *tcpTransportSocket) IsClosed() bool {
	return t.closed
}

func (t *tcpTransportSocket) Recv(r Register) (*Message, *MessageHandler, error) {
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
	// 2 bytes message size, size = all_size - Header(6 bytes)
	// 4 bytes message name crc32 id,
	// Message Body:
	var header [10]byte
	if _, err := io.ReadFull(t.reader, header[:]); err != nil {
		return nil, nil, fmt.Errorf("tcpTransportSocket.Recv header failed: %w", err)
	}

	var msgLen uint16
	// var msgType uint16
	var nameCrc uint32
	msgLen = binary.LittleEndian.Uint16(header[:2])
	// msgType = binary.LittleEndian.Uint16(header[4:6])
	nameCrc = binary.LittleEndian.Uint32(header[2:6])

	// check len
	// if msgLen > math.MaxUint16 || msgLen < 0 {
	// 	return nil, nil, fmt.Errorf("tcpTransportSocket.Recv failed: message length<%d> too long", msgLen)
	// }

	// check msg type
	// if msgType < BodyBegin || msgType >= BodyEnd {
	// 	return nil, nil, fmt.Errorf("tcpTransportSocket.Recv failed: marshal type<%d> error", msgType)
	// }

	// read body bytes
	bodyData := make([]byte, msgLen)
	if _, err := io.ReadFull(t.reader, bodyData); err != nil {
		return nil, nil, fmt.Errorf("tcpTransportSocket.Recv body failed: %w", err)
	}

	// get register handler
	h, err := r.GetHandler(nameCrc)
	if err != nil {
		return nil, nil, fmt.Errorf("tcpTransportSocket.Recv failed: %w", err)
	}

	var message Message
	// message.Type = codec.CodecType(msgType)
	message.Name = h.Name
	// message.Body, err = t.codecs[message.Type].Unmarshal(bodyData, h.RType)
	message.Body, err = t.codecs[0].Unmarshal(bodyData, h.RType)
	if err != nil {
		return nil, nil, fmt.Errorf("tcpTransportSocket.Recv unmarshal message body failed: %w", err)
	}

	return &message, h, err
}

func (t *tcpTransportSocket) Send(m *Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		if err := t.conn.SetDeadline(time.Now().Add(t.timeout)); err != nil {
			return err
		}
	}

	// if m.Type < BodyBegin || m.Type >= BodyEnd {
	// 	return fmt.Errorf("tcpTransportSocket.Send marshal type<%d> error", m.Type)
	// }

	// body, err := t.codecs[m.Type].Marshal(m.Body)
	body, err := t.codecs[0].Marshal(m.Body)
	if err != nil {
		return err
	}

	// Message Header:
	// 2 bytes message size, size = all_size - Header(6 bytes)
	// 4 bytes message name crc32 id,
	// Message Body:
	var bodySize uint16 = uint16(len(body))
	// items := strings.Split(m.Name, ".")
	// protoName := items[len(items)-1]
	var nameCrc uint32 = crc32.ChecksumIEEE([]byte(m.Name))
	var header []byte = make([]byte, 10)

	binary.LittleEndian.PutUint16(header[:2], bodySize)
	// binary.LittleEndian.PutUint16(header[4:6], uint16(m.Type))
	binary.LittleEndian.PutUint32(header[2:6], uint32(nameCrc))

	if _, err := t.writer.Write(header); err != nil {
		return err
	}

	if _, err := t.writer.Write(body); err != nil {
		return err
	}

	return nil
}
