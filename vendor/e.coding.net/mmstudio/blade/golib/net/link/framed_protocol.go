package link

import (
	"encoding/binary"
	"io"
	"net"
)

// transport builder
type FramedProtocol struct {
	opts *ProtocolOptions
}

func (f *FramedProtocol) SetBytesPool(bp BytesPool) { f.opts.BytesPool = bp }
func (f *FramedProtocol) GetBytesPool() BytesPool   { return f.opts.BytesPool }
func (f *FramedProtocol) NewTransporter(rw io.ReadWriter) (t Transporter, err error) {
	conn, ok := rw.(net.Conn)
	if !ok {
		panic("now we need net.Conn")
	}
	return NewFramedTransporter(conn, f.opts), nil
}
func NewProtocol(opts ...ProtocolOption) *FramedProtocol {
	return &FramedProtocol{opts: NewProtocolOptions(opts...)}
}

// Deprecated:  use NewProtocol instead
func NewFramedProtocol(sizeofLen, maxRecv, maxSend int, byteOrder binary.ByteOrder) *FramedProtocol {
	if sizeofLen != SizeOfLen {
		panic("size of len must be 4")
	}
	return NewProtocol(WithProtocolOptionMaxSend(maxSend), WithProtocolOptionMaxRecv(maxRecv), WithProtocolOptionByteOrder(byteOrder))
}
