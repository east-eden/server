package link

import (
	"io"
)

type Protocol interface {
	NewTransporter(rw io.ReadWriter) (Transporter, error)
}

type BytesPoolProvider interface {
	SetBytesPool(BytesPool)
	GetBytesPool() BytesPool
}

type ProtocolFunc func(rw io.ReadWriter) (Transporter, error)

func (pf ProtocolFunc) NewTransporter(rw io.ReadWriter) (Transporter, error) {
	return pf(rw)
}
