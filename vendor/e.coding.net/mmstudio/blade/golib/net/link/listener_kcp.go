// +build kcp

package link

import (
	"net"

	kcp "github.com/xtaci/kcp-go"
)

func init() {
	RegisterMakeListener(KCP, kcpMakeListener)
}

func kcpMakeListener(address string) (ln net.Listener, err error) {
	return kcp.ListenWithOptions(address, nil, 10, 3)
}
