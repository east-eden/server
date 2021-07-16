// +build !windows

package gate

import (
	"e.coding.net/mmstudio/blade/golib/osutil"
	"e.coding.net/mmstudio/blade/xlistener"
	"net"
)

func (g *Gate) XListener() (net.Listener, error) {
	return xlistener.Listen(
		func() (listener net.Listener, e error) {
			return g.listener, nil
		}, xlistener.WithBacklogAccept(g.spec.XListenerBacklogAccept),
		xlistener.WithTimeoutCanRead(g.spec.XListenerTimeoutCanRead),
		xlistener.WithEnableHandshake(g.spec.XListenerHandshake))
}

func setMaxOpenFile() {
	osutil.Setrlimit(maxOpenfile)
}
