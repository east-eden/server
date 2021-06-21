package link

import (
	"github.com/gorilla/websocket"
	"net"
	"time"
)

type wsConn struct {
	// Connections support one concurrent reader and one concurrent writer.
	c          *websocket.Conn
	remain     []byte
}

// NewWebSocketConn return net.Conn with websocket protocol
// Connections support one concurrent reader and one concurrent writer.
func NewWebSocketConn(c *websocket.Conn) net.Conn {
	return &wsConn{
		c:      c,
		remain: make([]byte, 0),
	}
}

// Read reads data from the connection.
// support one concurrent reader
func (w *wsConn) Read(b []byte) (int, error) {
	if len(w.remain) > 0 {
		n := copy(b, w.remain)
		if n == len(w.remain) {
			w.remain = w.remain[:0]
		} else if n < len(w.remain) {
			copy(w.remain[0:], w.remain[n:])
		}
		return n, nil
	}
	t, payload, err := w.c.ReadMessage()
	if err != nil {
		return 0, err
	}
	if t == websocket.TextMessage || t == websocket.BinaryMessage {
		n := copy(b, payload)
		if n < len(payload) {
			w.remain = w.remain[:0]
			w.remain = append(w.remain, payload[n:]...)
		}
		return n, nil
	}
	return 0, nil
}

// Write writes data to the connection.
// support one concurrent writer
func (w *wsConn) Write(b []byte) (int, error) {
	if e := w.c.WriteMessage(websocket.BinaryMessage, b); e != nil {
		return 0, e
	}
	return len(b), nil
}

func (w *wsConn) Close() error {
	return w.c.Close()
}

func (w *wsConn) LocalAddr() net.Addr {
	return w.c.LocalAddr()
}

func (w *wsConn) RemoteAddr() net.Addr {
	return w.c.RemoteAddr()
}

func (w *wsConn) SetDeadline(t time.Time) error {
	e := w.c.SetReadDeadline(t)
	if e != nil {
		return e
	}
	return w.c.SetWriteDeadline(t)
}

func (w *wsConn) SetReadDeadline(t time.Time) error {
	return w.c.SetReadDeadline(t)
}

func (w *wsConn) SetWriteDeadline(t time.Time) error {
	return w.c.SetWriteDeadline(t)
}
