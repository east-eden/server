// Package transport is an interface for synchronous connection based communication
package transport

import (
	"time"
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
type Transport interface {
	Init(...Option) error
	Options() Options
	Dial(addr string, opts ...DialOption) (Socket, error)
	Listen(addr string, opts ...ListenOption) (Listener, error)
	String() string
}

type Message struct {
	Type int
	Name string
	Body interface{}
}

type MessageFunc func(Socket, *Message)
type MessageHandler struct {
	Name string
	Fn   MessageFunc
}

type Register interface {
	RegisterMessage(string, MessageFunc) error
	GetHandler(uint32) (*MessageHandler, error)
}

type Socket interface {
	Recv() (*Message, *MessageHandler, error)
	Send(*Message) error
	Close() error
	Local() string
	Remote() string
}

type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error
}

type Option func(*Options)

type DialOption func(*DialOptions)

type ListenOption func(*ListenOptions)

var (
	DefaultDialTimeout = time.Second * 10
	DefaultRegister    = NewTransportRegister()
)

func NewTransport(opts ...Option) Transport {
	var options Options

	for _, o := range opts {
		o(&options)
	}
	return &tcpTransport{
		opts: options,
	}
}

func NewTransportRegister() Register {
	return &tcpTransportRegister{
		msgHandler: make(map[uint32]*MessageHandler, 0),
	}
}
