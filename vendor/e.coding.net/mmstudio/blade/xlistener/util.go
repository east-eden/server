package xlistener

import (
	"fmt"
	"syscall"
)

var (
	ErrClosed        = fmt.Errorf("epoll instance is closed")
	ErrRegistered    = fmt.Errorf("file descriptor is already registered in epoll instance")
	ErrNotRegistered = fmt.Errorf("file descriptor was not registered before in epoll instance")
)

func temporaryErr(err error) bool {
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	return errno.Temporary()
}
