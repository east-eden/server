// Package tcp provides a TCP transport
package transport

import (
	"context"
	"errors"

	"github.com/panjf2000/gnet/v2"
)

var (
	GNetRecvMaxSize           uint32 = 1024 * 1024
	ErrGNetReadFail                  = errors.New("gnet read failed")
	ErrGNetReadLengthLimit           = errors.New("gnet read length limit")
	ErrGNetEventHandlerAssert        = errors.New("gnet event handler assert")
)

type gnetTransport struct {
	opts *Options
}

func (t *gnetTransport) Init(opts ...Option) {
	t.opts = DefaultTransportOptions()

	for _, o := range opts {
		o(t.opts)
	}
}

func (t *gnetTransport) Options() *Options {
	return t.opts
}

func (t *gnetTransport) Protocol() string {
	return "gnet"
}

func (t *gnetTransport) Dial(addr string, opts ...DialOption) (Socket, error) {
	return nil, nil
}

// ListenAndServe starts a gnet engine with the given TransportServer.
//
// gnet/v2 dropped the Codec abstraction, so framing must be implemented by the
// TransportServer itself inside its OnTraffic callback.
func (t *gnetTransport) ListenAndServe(ctx context.Context, addr string, server TransportServer, opts ...ListenOption) error {
	eventServer, ok := server.(gnet.EventHandler)
	if !ok {
		return ErrGNetEventHandlerAssert
	}

	return gnet.Run(
		eventServer,
		addr,
		gnet.WithTicker(true),
		gnet.WithReuseAddr(true),
		gnet.WithMulticore(false),
	)
}
