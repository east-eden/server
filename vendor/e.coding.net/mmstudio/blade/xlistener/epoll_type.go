package xlistener

import (
	"net"
	"os"
)

type epollEvent uint32

type fdInfo struct {
	file      *os.File
	fd        int
	conn      net.Conn
	action    func(info *fdInfo, event epollEvent)
	tsTimeout int64
}

func (evt epollEvent) String() (str string) {
	name := func(event epollEvent, name string) {
		if evt&event == 0 {
			return
		}
		if str != "" {
			str += "|"
		}
		str += name
	}

	name(EpollIn, "EPOLLIN")
	name(EpollRdHup, "EPOLLRDHUP")
	name(EpollClosed, "_EPOLLCLOSED")
	return
}
