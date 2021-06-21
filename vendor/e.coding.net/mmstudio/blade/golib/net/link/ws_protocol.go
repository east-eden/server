package link

import (
	"io"
)

// transport builder
type WsProtocol struct {
	opts *WsProtocolOptions
}

func (f *WsProtocol) NewTransporter(rw io.ReadWriter) (t Transporter, err error) {
	conn, ok := rw.(*wsConn)
	if !ok {
		panic("now we need wsConn")
	}
	return NewWsTransporter(conn.c, f.opts), nil
}

func NewWsProtocol(opts ...WsProtocolOption) Protocol {
	return &WsProtocol{opts: NewWsProtocolOptions(opts...)}
}
