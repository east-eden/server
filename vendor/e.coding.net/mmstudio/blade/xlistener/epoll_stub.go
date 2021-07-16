// +build !linux

package xlistener

import (
	"time"
)

var (
	EpollIn      epollEvent = 0x1    // EPOLLIN = 0x1
	EpollRdHup   epollEvent = 0x2000 // since Linux 2.6.17, stream socket peer closed connection, or shut down writing half of connection.
	EpollOneShot epollEvent = 0x40000000
	EpollClosed  epollEvent = 0x20 // the epoll instance is closed.
)

type epoll struct {
}

func epollCreate(time.Duration, func(error)) (ep *epoll, err error) {
	return &epoll{}, nil
}

func (ep *epoll) Close() (err error) {
	return nil
}

func (ep *epoll) Add(fd int, events epollEvent, fa *fdInfo) (err error) {
	fa.action(fa, EpollIn)
	return nil
}

func (ep *epoll) Del(fd int) (err error) {
	return nil
}

func (ep *epoll) Mod(fd int, events epollEvent) (err error) {
	return nil
}
