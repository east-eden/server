package link

import (
	"io"
	"net"
	"strings"
	"time"

	"e.coding.net/mmstudio/blade/golib/net/tcpkeepalive"
)

func TCPConnSetup(conn *net.TCPConn, spec *Spec) *net.TCPConn {
	_ = conn.SetNoDelay(spec.TCPNoDelay)
	_ = conn.SetWriteBuffer(spec.TCPWriteBuffer)
	_ = conn.SetReadBuffer(spec.TCPReadBuffer)
	_ = conn.SetKeepAlive(true)
	_ = conn.SetLinger(spec.TCPLingerSecond)
	if spec.KeepAlivePeriod != 0 {
		_ = tcpkeepalive.SetKeepAlive(conn, spec.KeepAlivePeriod, 6, 5*time.Second)
	}
	return conn
}

func isTemporaryErr(err error) bool {
	netErr, ok := err.(net.Error)
	if ok && netErr.Timeout() && netErr.Temporary() {
		return true
	}
	return false
}

func Dial(spec *Spec, protocol Protocol) (Session, error) {
	conn, err := net.Dial(spec.Proto, spec.Address)
	if err != nil {
		return nil, err
	}
	trans, err := protocol.NewTransporter(conn)
	if err != nil {
		return nil, err
	}
	return NewSession(trans, spec), nil
}

func DialTimeout(spec *Spec, protocol Protocol, timeout time.Duration) (Session, error) {
	conn, err := net.DialTimeout(spec.Proto, spec.Address, timeout)
	if err != nil {
		return nil, err
	}
	trans, err := protocol.NewTransporter(conn)
	if err != nil {
		return nil, err
	}
	return NewSession(trans, spec), nil
}

func Accept(listener net.Listener) (net.Conn, error) {
	var tempDelay time.Duration
	for {
		conn, err := listener.Accept()
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
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil, io.EOF
			}
			return nil, err
		}
		tempDelay = 0
		return conn, nil
	}
}
