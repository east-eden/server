package micro_service

import (
	"errors"

	"go-micro.dev/v4/client"
	"google.golang.org/protobuf/proto"
)

var (
	ErrRpcMethodNotFound = errors.New("rpc method not found")
)

type RpcRegisterCaller func(nodeName string, msg proto.Message, copts ...client.CallOption) (resp proto.Message, err error)

type RpcRegister struct {
	rpcCallers map[string]RpcRegisterCaller
	protos     map[string]proto.Message
}

func NewRpcRegister() *RpcRegister {
	return &RpcRegister{
		rpcCallers: make(map[string]RpcRegisterCaller),
		protos:     make(map[string]proto.Message),
	}
}

func (r *RpcRegister) Register(msg proto.Message, fn RpcRegisterCaller) {
	r.rpcCallers[string(proto.MessageName(msg))] = fn
	r.protos[string(proto.MessageName(msg))] = msg
}

func (r *RpcRegister) GetProto(name string) (proto.Message, bool) {
	msg, ok := r.protos[name]
	return msg, ok
}

func (r *RpcRegister) Call(nodeName string, msg proto.Message, copts ...client.CallOption) (resp proto.Message, err error) {
	method := string(proto.MessageName(msg))
	fn, ok := r.rpcCallers[method]
	if !ok {
		return nil, ErrRpcMethodNotFound
	}
	return fn(nodeName, msg, copts...)
}
