// +build windows

package reuseport

import "net"

// NewReusablePortListener returns net.FileListener that created from a file discriptor for a socket with SO_REUSEPORT option.
func NewReusablePortListener(proto, addr string) (l net.Listener, err error) {
	return net.Listen(proto, addr)
}
