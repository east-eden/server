package link

import (
	"net"
)

type NetworkType string

const (
	TCP       NetworkType = "tcp"
	TCP4      NetworkType = "tcp4"
	TCP6      NetworkType = "tcp6"
	HTTP      NetworkType = "http"
	KCP       NetworkType = "kcp"
	ReusePort NetworkType = "reuseport"
)

var newListeners = make(map[NetworkType]NewListener)

func init() {
	newListeners[TCP] = tcpNewListener(TCP)
	newListeners[TCP4] = tcpNewListener(TCP)
	newListeners[TCP6] = tcpNewListener(TCP)
	newListeners[HTTP] = tcpNewListener(TCP)
}

// RegisterMakeListener registers a NewListener for network.
func RegisterMakeListener(network NetworkType, ml NewListener) {
	newListeners[network] = ml
}

type NewListener func(s *Server, address string) (ln net.Listener, err error)

func tcpNewListener(network NetworkType) func(s *Server, address string) (ln net.Listener, err error) {
	return func(s *Server, address string) (ln net.Listener, err error) {
		return net.Listen(string(network), address)
	}
}
