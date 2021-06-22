// Package transport is an interface for synchronous connection based communication
package transport

import (
	"context"
	"errors"
	"reflect"
	"time"

	"e.coding.net/mmstudio/blade/server/transport/codec"
)

var (
	ErrTransportTcpPacketTooLong        = errors.New("transport tcp send packet too long")
	TcpPacketMaxSize             uint32 = 1024 * 1024 // 单个tcp包数据上限
)

const (
	BodyBegin    = 0
	BodyProtobuf = 0
	BodyJson     = 1
	BodyEnd      = 2
)

// Transport is an interface which is used for communication between
// services. It uses connection based socket send/recv semantics and
// has various implementations; http, grpc, quic.
type TransportHandler func(context.Context, Socket)
type Transport interface {
	Init(...Option) error
	Options() Options
	Dial(addr string, opts ...DialOption) (Socket, error)
	ListenAndServe(ctx context.Context, addr string, handler TransportHandler, opts ...ListenOption) error
	Protocol() string
}

type Listener interface {
	Addr() string
	Close() error
	Accept(context.Context, TransportHandler) error
}

type Message struct {
	// Type codec.CodecType
	Name string
	Body interface{}
}

type MessageFunc func(context.Context, Socket, *Message) error
type MessageHandler struct {
	Name  string
	RType reflect.Type
	Fn    MessageFunc
}

type Socket interface {
	Recv(Register) (*Message, *MessageHandler, error)
	Send(*Message) error
	Close()
	IsClosed() bool
	Local() string
	Remote() string
	PbMarshaler() codec.Marshaler
}

type Option func(*Options)

type DialOption func(*DialOptions)

type ListenOption func(*ListenOptions)

var (
	DefaultDialTimeout  = time.Second * 10
	DefaultServeTimeout = time.Second * 20
	DefaultRegister     = NewTransportRegister()
)

func NewTransport(proto string) Transport {
	switch proto {
	case "tcp":
		return &tcpTransport{}
	case "ws":
		return &wsTransport{}
	default:
		return nil
	}
}

func NewTransportRegister() Register {
	return &defaultTransportRegister{msgHandler: make(map[uint32]*MessageHandler)}
}
