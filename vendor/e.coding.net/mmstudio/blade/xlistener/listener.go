package xlistener

import (
	"net"
	"os"
	"sync"
	"syscall"
	"time"

	"e.coding.net/mmstudio/blade/xlistener/internal/ip"
	"e.coding.net/mmstudio/blade/xlistener/internal/sync2"
)

var _ net.Listener = &Listener{}

type Listener struct {
	base       net.Listener
	closed     sync2.AtomicBool
	acceptChan chan net.Conn
	closeOnce  sync.Once
	closeChan  chan struct{}
	epoll      *epoll
	conf       *Conf
	addr       []byte
	connId     sync2.AtomicUint64
}

func (l *Listener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.acceptChan:
		return conn, nil
	case <-l.closeChan:
	}
	return nil, os.ErrClosed
}

func (l *Listener) Addr() net.Addr {
	return l.base.Addr()
}

func (l *Listener) Close() error {
	l.closeOnce.Do(func() {
		l.closed.Set(true)
		close(l.closeChan)
	})
	return l.base.Close()
}

func Listen(baseListen func() (net.Listener, error), opts ...ConfOption) (net.Listener, error) {
	base, err := baseListen()
	if err != nil {
		return nil, err
	}

	l := &Listener{base: base, closeChan: make(chan struct{}), conf: NewConf(opts...)}
	l.acceptChan = make(chan net.Conn, l.conf.BacklogAccept)
	if l.epoll, err = epollCreate(l.conf.TimeoutCanRead, l.onWaitError); err != nil {
		return nil, err
	}
	_, port, err := net.SplitHostPort(base.Addr().String())
	if err != nil {
		return nil, err
	}
	l.addr = []byte(net.JoinHostPort(ip.GetLocalIP(), port))

	go l.acceptLoop()

	return l, nil
}

func (l *Listener) onWaitError(err error) {
	l.conf.Warningf("epoll: wait loop error: %s", err)
}

func (l *Listener) acceptLoop() {
	for {
		conn, err := l.base.Accept()
		if err != nil {
			if !l.closed.Get() {
				l.conf.Warningf("accept failed: %v", err)
			}
			break
		}
		if err := l.handlerConn(conn); err != nil {
			l.conf.Warningf("epoll handler conn with error:%s,fallback to go classic", err.Error())
			l.toBacklogAccept(conn)
		}
	}
}

func (l *Listener) toBacklogAccept(conn net.Conn) {
	if l.conf.EnableHandshake {
		xc, err := newXConn(conn, l.connId.Add(1))
		if err != nil {
			l.conf.Warningf("new xconn error:%s addr:%s", err.Error(), conn.RemoteAddr())
			_ = conn.Close()
			return
		}
		go func(xcIn *xConn) {
			if err := xcIn.handshake(l.addr, l.conf.HandshakeTimeout); err != nil {
				l.conf.Warningf("handshake error:%s addr:%s", err.Error(), xcIn.RemoteAddr())
				_ = xcIn.Close()
				return
			} else {
				select {
				case l.acceptChan <- xcIn:
				case <-l.closeChan:
				default:
					_ = xcIn.Close()
					l.conf.Warningf("the accept queue blocked, close the connection:%s", xcIn.RemoteAddr())
				}
			}
		}(xc)
	} else {
		select {
		case l.acceptChan <- conn:
		case <-l.closeChan:
		default:
			_ = conn.Close()
			l.conf.Warningf("the accept queue blocked, close the connection:%s", conn.RemoteAddr())
		}
	}
}

func (l *Listener) onEpollEvent(fi *fdInfo, evt epollEvent) {
	_ = l.epoll.Del(fi.fd)
	_ = fi.file.Close()

	// we must check EpollRdHup first
	if evt&EpollRdHup != 0 || evt&EpollClosed != 0 {
		_ = fi.conn.Close()
	} else if evt&EpollIn != 0 {
		l.toBacklogAccept(fi.conn)
	} else {
		// fixme should just pass it to acceptChan?
		_ = fi.conn.Close()
		l.conf.Warningf("onEpollEvent got unexpected event:%v", evt)
	}
}

// filer describes an object that has ability to return os.File.
type filer interface {
	// File returns a copy of object's file descriptor.
	File() (*os.File, error)
}

func (l *Listener) handlerConn(conn net.Conn) error {
	f, err := conn.(filer).File()
	if err != nil {
		return err
	}

	fd := int(f.Fd())

	// Set the file back to non blocking mode since calling Fd() sets underlying
	// os.File to blocking mode. This is useful to get conn.Set{Read}Deadline
	// methods still working on source Conn.
	//
	// See https://golang.org/pkg/net/#TCPConn.File
	// See /usr/local/go/src/net/net.go: conn.File()
	if err = syscall.SetNonblock(fd, true); err != nil {
		_ = f.Close()
		return os.NewSyscallError("setnonblock", err)
	}

	// we must hold this os.File and close it when finished.
	return l.epoll.Add(fd,
		EpollRdHup|EpollIn|EpollOneShot, // use one shot
		&fdInfo{
			file:      f,
			fd:        fd,
			conn:      conn,
			tsTimeout: time.Now().Unix() + int64(l.conf.TimeoutCanRead.Seconds()),
			action:    l.onEpollEvent,
		})
}
