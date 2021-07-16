package link

import (
	"errors"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/gorilla/websocket"
)

func NewWsTransporter(c *websocket.Conn, opts *WsProtocolOptions) Transporter {
	if opts == nil {
		panic("opts should not ne nil")
	}
	return &WsTransporter{
		conn: c,
		opts: opts,
	}
}

type WsTransporter struct {
	conn             *websocket.Conn
	opts             *WsProtocolOptions
	readDeadlineSet  bool
	writeDeadlineSet bool
}

func (w *WsTransporter) Receive() (interface{}, error) {
	var err error
	retryCount := w.opts.RetryCountWhenTempError
	retry := false
	var buff []byte
	var t int
	for {
		if retry {
			log.Info().Int("retry_left", retryCount).Bool("deadline_set", w.readDeadlineSet).Msg("Receive with retry")
			if w.readDeadlineSet {
				// 读取超时，允许进行buffer wait内的超时设置
				if errSetDeadline := w.SetReadDeadline(time.Now().Add(w.opts.RecvBufferWait)); errSetDeadline != nil {
					return nil, err
				}
			} else {
				// 系统层未设置超时
				return nil, err
			}
		}
		t, buff, err = w.conn.ReadMessage()
		if err == nil && (t == websocket.TextMessage || t == websocket.BinaryMessage) {
			return buff, nil
		}
		if err != nil && !isTemporaryErr(err) {
			// 连接关闭或非临时性错误
			return nil, err
		}
		if err == nil && t != websocket.TextMessage && t != websocket.BinaryMessage {
			// 返回temporary err 上次继续重试
			err = &wsTemporaryError{"websocket message type error"}
		}
		if retryCount == 0 {
			return 0, err
		}
		retryCount--
		time.Sleep(w.opts.ReadRetryIntervalWhenTempError)
		retry = true
	}

	return nil, err
}

func (w *WsTransporter) Send(p interface{}) (err error) {
	switch buff := p.(type) {
	case []byte:
		// TODO retry
		err = w.conn.WriteMessage(websocket.BinaryMessage, buff)
	default:
		return errors.New("type not supported")
	}
	return
}

func (w *WsTransporter) Close() error {
	return w.conn.Close()
}

func (w *WsTransporter) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *WsTransporter) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *WsTransporter) SetReadDeadline(t time.Time) error {
	e := w.conn.SetReadDeadline(t)
	if e == nil {
		w.readDeadlineSet = true
	}
	return e
}

func (w *WsTransporter) SetWriteDeadline(t time.Time) error {
	e := w.conn.SetWriteDeadline(t)
	if e == nil {
		w.writeDeadlineSet = true
	}
	return e
}

type wsTemporaryError struct {
	c string
}

func (w *wsTemporaryError) Error() string {
	return w.c
}
func (w *wsTemporaryError) Timeout() bool {
	return true
}
func (w *wsTemporaryError) Temporary() bool {
	return true
}
