// Package tcp provides a TCP transport
package transport

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	maddr "github.com/micro/go-micro/util/addr"
	"github.com/micro/go-micro/util/log"
	mnet "github.com/micro/go-micro/util/net"
	mls "github.com/micro/go-micro/util/tls"
	"github.com/yokaiio/yokai_server/internal/codec"
)

var tcpReadBufMax = 1024 * 1024 * 2

type tcpTransport struct {
	opts Options
}

type tcpTransportSocket struct {
	conn    net.Conn
	codecs  []codec.Marshaler
	timeout time.Duration
}

type tcpTransportListener struct {
	listener net.Listener
	timeout  time.Duration
}

func (t *tcpTransportSocket) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *tcpTransportSocket) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *tcpTransportSocket) Recv() (*Message, error) {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}

	// Header:
	// 4 bytes message size, size = all_size - Header(8 bytes)
	// 2 bytes message type,
	// 2 bytes message name length,
	// Message Name:
	// len(message length) bytes message name,
	// Message Body:
	var header [8]byte
	if _, err := io.ReadFull(t.conn, header[:]); err != nil {
		return nil, fmt.Errorf("connection read message header error:%v", err)
	}

	var msgLen uint32
	var msgType uint16
	var nameLen uint16
	msgLen = binary.LittleEndian.Uint32(header[:4])
	msgType = binary.LittleEndian.Uint16(header[4:6])
	nameLen = binary.LittleEndian.Uint16(header[6:8])

	// check len
	if msgLen > uint32(tcpReadBufMax) {
		return nil, fmt.Errorf("connection read failed with too long message:%v", msgLen)
	} else if msgLen < 8 {
		return nil, fmt.Errorf("connection read failed with too short message:%v", msgLen)
	}

	// check msg type
	if msgType < BodyBegin || msgType >= BodyEnd {
		return nil, fmt.Errorf("marshal type error:%v", msgType)
	}

	// read other bytes
	otherData := make([]byte, msgLen)
	if _, err := io.ReadFull(t.conn, otherData); err != nil {
		return nil, fmt.Errorf("connection read message body failed:%v", err)
	}

	var err error
	var message Message
	message.Type = int(msgType)
	message.Name = string(otherData[:nameLen])
	message.Body, err = t.codecs[message.Type].Unmarshal(otherData[nameLen:], message.Name)
	if err != nil {
		return nil, fmt.Errorf("connection unmarshal message body failed:%v", err)
	}

	return &message, err
}

func (t *tcpTransportSocket) Send(m *Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}

	if m.Type < BodyBegin || m.Type >= BodyEnd {
		return fmt.Errorf("marshal type error:%v", m.Type)
	}

	out, err := t.codecs[m.Type].Marshal(m.Body)
	if err != nil {
		return err
	}

	// Header:
	// 4 bytes message size, size = all_size - Header(8 bytes)
	// 2 bytes message type,
	// 2 bytes message name length,
	// Message Name:
	// len(message length) bytes message name,
	// Message Body:
	var bodySize uint32 = uint32(len(out))
	var nameLen uint32 = uint32(len(m.Name))
	var data []byte = make([]byte, 8+nameLen+bodySize)
	binary.LittleEndian.PutUint32(data[:4], nameLen+bodySize)
	binary.LittleEndian.PutUint16(data[4:6], uint16(m.Type))
	binary.LittleEndian.PutUint16(data[6:8], uint16(nameLen))
	copy(data[8:8+nameLen], []byte(m.Name))
	copy(data[8+nameLen:], out)

	if _, err := t.conn.Write(data); err != nil {
		return err
	}

	return nil
}

func (t *tcpTransportSocket) Close() error {
	return t.conn.Close()
}

func (t *tcpTransportListener) Addr() string {
	return t.listener.Addr().String()
}

func (t *tcpTransportListener) Close() error {
	return t.listener.Close()
}

func (t *tcpTransportListener) Accept(fn func(Socket)) error {
	var tempDelay time.Duration

	for {
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
				log.Logf("http: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		c.(*net.TCPConn).SetKeepAlive(true)
		c.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)

		sock := &tcpTransportSocket{
			timeout: t.timeout,
			conn:    c,
			codecs:  []codec.Marshaler{codec.NewProtobufCodec(), codec.NewJsonCodec()},
		}

		go func() {
			// TODO: think of a better error response strategy
			defer func() {
				if r := recover(); r != nil {
					sock.Close()
				}
			}()

			fn(sock)
		}()
	}
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
		conn:    conn,
		codecs:  []codec.Marshaler{codec.NewProtobufCodec(), codec.NewJsonCodec()},
		timeout: t.opts.Timeout,
	}, nil
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

	return &tcpTransportListener{
		timeout:  t.opts.Timeout,
		listener: l,
	}, nil
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

func (t *tcpTransport) String() string {
	return "tcp"
}
