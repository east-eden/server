// Package ws provides a websocket transport
package transport

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"net/http"
	"time"

	"github.com/east-eden/server/transport/codec"
	"github.com/gorilla/websocket"
	"github.com/valyala/bytebufferpool"
	"go.uber.org/atomic"
)

var upgrader = websocket.Upgrader{}

func newWsTransportSocket() *wsTransportSocket {
	return &wsTransportSocket{
		codecs: []codec.Marshaler{&codec.ProtoBufMarshaler{}, &codec.JsonMarshaler{}},
	}
}

type wsTransport struct {
	opts Options
}

func (t *wsTransport) Init(opts ...Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *wsTransport) Options() Options {
	return t.opts
}

func (t *wsTransport) Protocol() string {
	return "ws"
}

func (t *wsTransport) Dial(addr string, opts ...DialOption) (Socket, error) {
	dopts := DialOptions{
		Timeout: DefaultDialTimeout,
	}

	for _, opt := range opts {
		opt(&dopts)
	}

	if t.opts.TLSConfig != nil {
		websocket.DefaultDialer.TLSClientConfig = t.opts.TLSConfig
	}

	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	return &wsTransportSocket{
		conn:    conn,
		codecs:  []codec.Marshaler{&codec.ProtoBufMarshaler{}, &codec.JsonMarshaler{}},
		timeout: t.opts.Timeout,
	}, nil
}

func (t *wsTransport) ListenAndServe(ctx context.Context, addr string, handler TransportServer, opts ...ListenOption) error {
	wsHandler := &wsServeHandler{
		ctx:     ctx,
		srv:     handler,
		timeout: t.opts.Timeout,
	}

	// wsHandler.sockPool.New = newWsTransportSocket

	server := &http.Server{
		Addr:      addr,
		Handler:   wsHandler,
		TLSConfig: t.opts.TLSConfig,
	}

	return server.ListenAndServeTLS("", "")
}

type wsServeHandler struct {
	ctx     context.Context
	srv     TransportServer
	timeout time.Duration
	// sockPool sync.Pool
}

func (h *wsServeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, error := upgrader.Upgrade(w, r, nil)
	if error != nil {
		http.NotFound(w, r)
		return
	}

	// sock := h.sockPool.Get().(*wsTransportSocket)
	sock := newWsTransportSocket()
	sock.timeout = h.timeout
	sock.conn = conn
	sock.closed.Store(false)

	// handle in workerpool
	h.srv.HandleSocket(h.ctx, sock)
}

type wsTransportSocket struct {
	conn    *websocket.Conn
	codecs  []codec.Marshaler
	timeout time.Duration
	closed  atomic.Bool
}

func (t *wsTransportSocket) Close() {
	t.closed.Store(true)
	_ = t.conn.Close()
}

func (t *wsTransportSocket) IsClosed() bool {
	return t.closed.Load()
}

func (t *wsTransportSocket) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *wsTransportSocket) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *wsTransportSocket) PbMarshaler() codec.Marshaler {
	return t.codecs[0]
}

func (t *wsTransportSocket) Recv(r Register) (*Message, *MessageHandler, error) {
	if t.IsClosed() {
		return nil, nil, errors.New("wsTransportSocket.Recv failed: socket closed")
	}

	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		_ = t.conn.SetReadDeadline(time.Now().Add(t.timeout))
	}

	// Message Header:
	// 4 bytes message size, size = 4 bytes name crc + proto binary size
	// 4 bytes message name crc32 id,
	// Message Body:

	_, data, err := t.conn.ReadMessage()
	if err != nil {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv read message error:%v", err)
	}

	msgLen := binary.LittleEndian.Uint32(data[:4])
	nameCrc := binary.LittleEndian.Uint32(data[4:8])

	if msgLen > TcpPacketMaxSize {
		return nil, nil, ErrTransportReadSizeTooLong
	}

	// get register handler
	h, err := r.GetHandler(nameCrc)
	if err != nil {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv failed: %w", err)
	}

	bodyData := data[8:]
	var message Message
	message.Name = h.Name
	message.Body, err = t.codecs[0].Unmarshal(bodyData, h.RType)
	if err != nil {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv unmarshal message body failed: %w", err)
	}

	return &message, h, err
}

func (t *wsTransportSocket) Send(m *Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		_ = t.conn.SetWriteDeadline(time.Now().Add(t.timeout))
	}

	// if m.Type < BodyBegin || m.Type >= BodyEnd {
	// 	return fmt.Errorf("wsTransportSocket.Send marshal type<%d> error", m.Type)
	// }

	body, err := t.codecs[0].Marshal(m.Body)
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

	if err := t.conn.WriteMessage(websocket.BinaryMessage, buffer.Bytes()); err != nil {
		return err
	}

	return nil
}
