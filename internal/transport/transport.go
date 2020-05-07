// Package transport is an interface for synchronous connection based communication
package transport

import (
	"log"
	"reflect"
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
	Name  string
	RType reflect.Type
	Fn    MessageFunc
}

type Socket interface {
	Recv(Register) (*Message, *MessageHandler, error)
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

func NewTransport(proto string, opts ...Option) Transport {
	var options Options

	for _, o := range opts {
		o(&options)
	}

	switch proto {
	case "tcp":
		return &tcpTransport{opts: options}
	case "ws":
		return &wsTransport{opts: options}
	default:
		log.Fatal("unknown transport proto type:", proto)
		return nil
	}
}

func NewTransportRegister() Register {
	return &defaultTransportRegister{msgHandler: make(map[uint32]*MessageHandler, 0)}
}
