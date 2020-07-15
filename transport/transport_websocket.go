// Package ws provides a websocket transport
package transport

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yokaiio/yokai_server/transport/codec"
)

var wsReadBufMax = 1024 * 1024 * 2
var upgrader = websocket.Upgrader{}

func newWsTransportSocket() interface{} {
	return &wsTransportSocket{
		codecs:        []codec.Marshaler{&codec.ProtoBufMarshaler{}, &codec.JsonMarshaler{}},
		evictedHandle: []func(Socket){},
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
		conn:          conn,
		codecs:        []codec.Marshaler{&codec.ProtoBufMarshaler{}, &codec.JsonMarshaler{}},
		timeout:       t.opts.Timeout,
		evictedHandle: []func(Socket){},
	}, nil
}

func (t *wsTransport) ListenAndServe(ctx context.Context, addr string, handler TransportHandler, opts ...ListenOption) error {
	wsHandler := &wsServeHandler{
		ctx:     ctx,
		fn:      handler,
		timeout: t.opts.Timeout,
	}

	wsHandler.sockPool.New = newWsTransportSocket

	server := &http.Server{
		Addr:      addr,
		Handler:   wsHandler,
		TLSConfig: t.opts.TLSConfig,
	}

	return server.ListenAndServeTLS("", "")
}

type wsServeHandler struct {
	ctx      context.Context
	fn       TransportHandler
	timeout  time.Duration
	sockPool sync.Pool
}

func (h *wsServeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, error := upgrader.Upgrade(w, r, nil)
	if error != nil {
		http.NotFound(w, r)
		return
	}

	sock := h.sockPool.Get().(*wsTransportSocket)
	sock.timeout = h.timeout
	sock.conn = conn
	sock.closed = false
	sock.evictedHandle = []func(Socket){}

	// handle in workerpool
	subCtx, cancel := context.WithCancel(h.ctx)
	h.fn(subCtx, sock, func() {
		cancel()
		h.sockPool.Put(sock)
	})
}

type wsTransportSocket struct {
	conn          *websocket.Conn
	codecs        []codec.Marshaler
	timeout       time.Duration
	evictedHandle []func(Socket)
	closed        bool
}

func (t *wsTransportSocket) AddEvictedHandle(f func(Socket)) {
	t.evictedHandle = append(t.evictedHandle, f)
}

func (t *wsTransportSocket) Close() error {
	for _, handle := range t.evictedHandle {
		handle(t)
	}

	t.closed = true
	return t.conn.Close()
}

func (t *wsTransportSocket) IsClosed() bool {
	return t.closed
}

func (t *wsTransportSocket) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *wsTransportSocket) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *wsTransportSocket) Recv(r Register) (*Message, *MessageHandler, error) {
	if t.IsClosed() {
		return nil, nil, errors.New("wsTransportSocket.Recv failed: socket closed")
	}

	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetReadDeadline(time.Now().Add(t.timeout))
	}

	// Message Header:
	// 4 bytes message size, size = all_size - Header(10 bytes)
	// 2 bytes message type,
	// 4 bytes message name crc32 id,
	// Message Body:

	_, data, err := t.conn.ReadMessage()
	if err != nil {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv read message error:%v", err)
	}

	var msgLen uint32
	var msgType uint16
	var nameCrc uint32
	msgLen = binary.LittleEndian.Uint32(data[:4])
	msgType = binary.LittleEndian.Uint16(data[4:6])
	nameCrc = binary.LittleEndian.Uint32(data[6:10])

	// check len
	if msgLen > uint32(wsReadBufMax) || msgLen < 0 {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv failed: message length<%d> too long", msgLen)
	}

	// check msg type
	if msgType < BodyBegin || msgType >= BodyEnd {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv failed: marshal type<%d> error", msgType)
	}

	// get register handler
	h, err := r.GetHandler(nameCrc)
	if err != nil {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv failed: %w", err)
	}

	bodyData := data[10:]
	var message Message
	message.Type = codec.CodecType(msgType)
	message.Name = h.Name
	message.Body, err = t.codecs[message.Type].Unmarshal(bodyData, h.RType)
	if err != nil {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv unmarshal message body failed: %w", err)
	}

	return &message, h, err
}

func (t *wsTransportSocket) Send(m *Message) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetWriteDeadline(time.Now().Add(t.timeout))
	}

	if m.Type < BodyBegin || m.Type >= BodyEnd {
		return fmt.Errorf("wsTransportSocket.Send marshal type<%d> error", m.Type)
	}

	out, err := t.codecs[m.Type].Marshal(m.Body)
	if err != nil {
		return err
	}

	// Message Header:
	// 4 bytes message size, size = all_size - Header(10 bytes)
	// 2 bytes message type,
	// 4 bytes message name crc32 id,
	// Message Body:
	var bodySize uint32 = uint32(len(out))
	items := strings.Split(m.Name, ".")
	protoName := items[len(items)-1]
	var nameCrc uint32 = crc32.ChecksumIEEE([]byte(protoName))
	var data []byte = make([]byte, 10+bodySize)

	binary.LittleEndian.PutUint32(data[:4], bodySize)
	binary.LittleEndian.PutUint16(data[4:6], uint16(m.Type))
	binary.LittleEndian.PutUint32(data[6:10], uint32(nameCrc))
	copy(data[10:], out)

	if err := t.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return err
	}

	return nil
}
