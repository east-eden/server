package xlistener

import (
	"syscall"
	"time"

	"e.coding.net/mmstudio/blade/xlistener/internal/sync2"
	"golang.org/x/sys/unix"
)

// http://man7.org/linux/man-pages/man2/epoll_ctl.2.html
const (
	EpollIn = unix.EPOLLIN // available for read
	//https://github.com/golang/go/issues/10940
	EpollRdHup   = unix.EPOLLRDHUP // since Linux 2.6.17, stream socket peer closed connection, or shut down writing half of connection.
	EpollOneShot = unix.EPOLLONESHOT
	EpollClosed  = 0x20 // the epoll instance is closed.
)

type epoll struct {
	epollFd        int
	eventFd        int
	waitDone       chan struct{}
	closeFlag      sync2.AtomicInt32
	callbacks      ConcurrentMapInt64FdInfo
	timeoutSigRead time.Duration
}

func epollCreate(timeoutCanRead time.Duration, defaultOnWaitError func(error)) (ep *epoll, err error) {
	var fd int
	var r uintptr
	var errno syscall.Errno
	// same as go runtime flag with EpollCreate1
	if fd, err = unix.EpollCreate1(0x80000); err != nil {
		return
	}
	if r, _, errno = unix.Syscall(unix.SYS_EVENTFD2, 0, 0, 0); errno != 0 {
		err = errno
		return
	}
	eventFd := int(r)

	// Set finalizer for write end of socket pair to avoid data races when
	// closing epoll instance and EBADF errors on writing ctl bytes from callers.
	err = unix.EpollCtl(fd, unix.EPOLL_CTL_ADD, eventFd, &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(eventFd),
	})

	if err != nil {
		_ = unix.Close(fd)
		_ = unix.Close(eventFd)
		return
	}

	ep = &epoll{
		epollFd:        fd,
		eventFd:        eventFd,
		callbacks:      NewConcurrentMapInt64FdInfo(),
		waitDone:       make(chan struct{}),
		timeoutSigRead: timeoutCanRead,
	}

	go ep.wait(defaultOnWaitError)
	return
}

// Close stops wait loop and closes all underlying resources.
func (ep *epoll) Close() (err error) {
	if ep.closeFlag.CompareAndSwap(0, 1) {
		if _, err = unix.Write(ep.eventFd, []byte{1, 0, 0, 0, 0, 0, 0, 0}); err != nil {
			return
		}
		<-ep.waitDone
		if err = unix.Close(ep.eventFd); err != nil {
			return
		}
		for info := range ep.callbacks.Iter() {
			if info.Val != nil {
				info.Val.action(info.Val, EpollClosed)
			}
		}
		return nil
	}
	return ErrClosed
}

func (ep *epoll) isClosed() bool {
	return ep.closeFlag.Get() == 1
}

//Add adds fd to epoll set with given events.
func (ep *epoll) Add(fd int, events epollEvent, fa *fdInfo) (err error) {
	if ep.isClosed() {
		return ErrClosed
	}

	if ep.callbacks.Has(fd) {
		return ErrRegistered
	}

	ep.callbacks.Set(fd, fa)
	if err = unix.EpollCtl(ep.epollFd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: uint32(events), Fd: int32(fd)}); err != nil {
		ep.callbacks.Remove(fd)
	}
	return
}

// Del removes epollFd from epoll set.
func (ep *epoll) Del(fd int) (err error) {
	if ep.isClosed() {
		return ErrClosed
	}
	if !ep.callbacks.Has(fd) {
		return ErrNotRegistered
	}
	ep.callbacks.Remove(fd)
	return unix.EpollCtl(ep.epollFd, unix.EPOLL_CTL_DEL, fd, nil)
}

// Mod sets to listen events on epollFd.
func (ep *epoll) Mod(fd int, events epollEvent) (err error) {
	if ep.isClosed() {
		return ErrClosed
	}
	if !ep.callbacks.Has(fd) {
		return ErrNotRegistered
	}
	return unix.EpollCtl(ep.epollFd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{Events: uint32(events), Fd: int32(fd)})
}

const (
	maxWaitEventsBegin = 1024
	maxWaitEventsStop  = 32768
)

func (ep *epoll) wait(onError func(error)) {
	defer func() {
		if err := unix.Close(ep.epollFd); err != nil {
			onError(err)
		}
		close(ep.waitDone)
	}()

	var n int
	var err error
	var now time.Time
	var nowUnix int64

	timeLastCheckTimeout := time.Now()
	timeoutCanReadMicroseconds := int(ep.timeoutSigRead.Microseconds())
	events := make([]unix.EpollEvent, maxWaitEventsBegin)

	for {
		n, err = unix.EpollWait(ep.epollFd, events, timeoutCanReadMicroseconds)
		if err != nil {
			if temporaryErr(err) {
				continue
			}
			onError(err)
			return
		}

		for i := 0; i < n; i++ {
			fd := int(events[i].Fd)
			if fd == ep.eventFd { // signal to close
				return
			}
			if fdInfo, ok := ep.callbacks.Get(fd); ok {
				fdInfo.action(fdInfo, epollEvent(events[i].Events))
			}
		}

		// check timeout
		now = time.Now()
		nowUnix = now.Unix()
		if timeLastCheckTimeout.Add(ep.timeoutSigRead).Before(now) {
			for info := range ep.callbacks.Iter() {
				if info.Val != nil && info.Val.tsTimeout < nowUnix {
					_ = ep.Del(info.Key)
				}
			}
			timeLastCheckTimeout = now
		}

		if n == len(events) && n*2 <= maxWaitEventsStop {
			events = make([]unix.EpollEvent, n*2)
		}
	}
}
