// Package transport is an interface for synchronous connection based communication
package transport

import (
	"context"
	"errors"
	"reflect"
	"time"

	"google.golang.org/protobuf/proto"
)

var (
	ErrTransportTcpPacketTooLong        = errors.New("transport tcp send packet too long")
	ErrTransportReadSizeTooLong         = errors.New("transport recv msg length > 1MB")
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
type TransportServer interface {
	HandleSocket(context.Context, Socket)
}

type Transport interface {
	Init(...Option) error
	Options() Options
	Dial(addr string, opts ...DialOption) (Socket, error)
	ListenAndServe(ctx context.Context, addr string, server TransportServer, opts ...ListenOption) error
	Protocol() string
}

type Listener interface {
	Addr() string
	Close() error
	Accept(context.Context, TransportServer) error
}

type MessageFunc func(context.Context, Socket, proto.Message) error
type MessageHandler struct {
	Name  string
	RType reflect.Type
	Fn    MessageFunc
}

type Socket interface {
	Recv(Register) (proto.Message, *MessageHandler, error)
	Send(proto.Message) error
	Close()
	IsClosed() bool
	Local() string
	Remote() string
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
	case "gnet":
		return &gnetTransport{}
	default:
		return nil
	}
}

func NewTransportRegister() Register {
	return &defaultTransportRegister{msgHandler: make(map[uint32]*MessageHandler)}
}
