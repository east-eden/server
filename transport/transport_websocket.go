// Package ws provides a websocket transport
package transport

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"net/http"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/transport/codec"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func newWsTransportSocket() interface{} {
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

	// handle in workerpool
	subCtx, cancel := context.WithCancel(h.ctx)
	h.fn(subCtx, sock, func() {
		cancel()
		h.sockPool.Put(sock)
	})
}

type wsTransportSocket struct {
	conn    *websocket.Conn
	codecs  []codec.Marshaler
	timeout time.Duration
	closed  bool
}

func (t *wsTransportSocket) Close() error {
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
	// 2 bytes message size, size = all_size - Header(6 bytes)
	// 4 bytes message name crc32 id,
	// Message Body:

	_, data, err := t.conn.ReadMessage()
	if err != nil {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv read message error:%v", err)
	}

	// var msgLen uint16
	// var msgType uint16
	var nameCrc uint32
	// msgLen = binary.LittleEndian.Uint16(data[:2])
	// msgType = binary.LittleEndian.Uint16(data[4:6])
	nameCrc = binary.LittleEndian.Uint32(data[2:6])

	// check len
	// if msgLen > uint32(wsReadBufMax) || msgLen < 0 {
	// 	return nil, nil, fmt.Errorf("wsTransportSocket.Recv failed: message length<%d> too long", msgLen)
	// }

	// check msg type
	// if msgType < BodyBegin || msgType >= BodyEnd {
	// 	return nil, nil, fmt.Errorf("wsTransportSocket.Recv failed: marshal type<%d> error", msgType)
	// }

	// get register handler
	h, err := r.GetHandler(nameCrc)
	if err != nil {
		return nil, nil, fmt.Errorf("wsTransportSocket.Recv failed: %w", err)
	}

	bodyData := data[6:]
	var message Message
	// message.Type = codec.CodecType(msgType)
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
		t.conn.SetWriteDeadline(time.Now().Add(t.timeout))
	}

	// if m.Type < BodyBegin || m.Type >= BodyEnd {
	// 	return fmt.Errorf("wsTransportSocket.Send marshal type<%d> error", m.Type)
	// }

	out, err := t.codecs[0].Marshal(m.Body)
	if err != nil {
		return err
	}

	// Message Header:
	// 2 bytes message size, size = all_size - Header(6 bytes)
	// 4 bytes message name crc32 id,
	// Message Body:
	var bodySize uint16 = uint16(len(out))
	// items := strings.Split(m.Name, ".")
	// protoName := items[len(items)-1]
	var nameCrc uint32 = crc32.ChecksumIEEE([]byte(m.Name))
	var data []byte = make([]byte, 10+bodySize)

	binary.LittleEndian.PutUint16(data[:2], bodySize)
	// binary.LittleEndian.PutUint16(data[4:6], uint16(m.Type))
	binary.LittleEndian.PutUint32(data[2:6], uint32(nameCrc))
	copy(data[6:], out)

	if err := t.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return err
	}

	return nil
}
