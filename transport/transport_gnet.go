// Package tcp provides a TCP transport
package transport

import (
	"context"
	"encoding/binary"
	"errors"

	"github.com/panjf2000/gnet"
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

func (t *gnetTransport) ListenAndServe(ctx context.Context, addr string, server TransportServer, opts ...ListenOption) error {
	eventServer, ok := server.(gnet.EventHandler)
	if !ok {
		return ErrGNetEventHandlerAssert
	}

	return gnet.Serve(
		eventServer,
		addr,
		gnet.WithCodec(&gnetCodec{}),
		gnet.WithTicker(true),
	)
}

type gnetCodec struct{}

// Encode encodes frames upon server responses into TCP stream.
func (codec *gnetCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return buf, nil
}

// Decode decodes frames from TCP stream via specific implementation.
func (codec *gnetCodec) Decode(c gnet.Conn) ([]byte, error) {
	bufLen := c.BufferLength()
	if bufLen <= 0 {
		return nil, nil
	}

	size, sizeBuf := c.ReadN(4)
	if size != 4 {
		c.ResetBuffer()
		return sizeBuf, ErrGNetReadFail
	}

	msgLen := binary.LittleEndian.Uint32(sizeBuf)
	if msgLen > GNetRecvMaxSize {
		c.ResetBuffer()
		return sizeBuf, ErrGNetReadLengthLimit
	}

	if c.ShiftN(4) != 4 {
		c.ResetBuffer()
		return sizeBuf, ErrGNetReadFail
	}

	bodySize, bodyBuf := c.ReadN(int(msgLen))
	if bodySize != int(msgLen) {
		c.ResetBuffer()
		return bodyBuf, ErrGNetReadFail
	}

	if c.ShiftN(bodySize) != bodySize {
		c.ResetBuffer()
		return bodyBuf, ErrGNetReadFail
	}

	return bodyBuf, nil
}
