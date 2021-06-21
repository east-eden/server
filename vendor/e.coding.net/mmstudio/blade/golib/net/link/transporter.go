package link

import (
	"net"
	"time"
)

// just define how messages are communicated between client and server.
type Transporter interface {
	Receive() (interface{}, error)
	Send(interface{}) (err error)
	Close() error
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

type ClearSendChan interface {
	ClearSendChan(Session, <-chan interface{})
}

type RemoteAddr interface {
	RemoteAddr() net.Addr
}
