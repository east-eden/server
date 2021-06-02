package transport

import (
	"errors"
	"fmt"
	"hash/crc32"
	"reflect"

	"bitbucket.org/funplus/server/transport/codec"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUnregistedMessage = errors.New("unregisted message")
)

type Register interface {
	RegisterProtobufMessage(proto.Message, MessageFunc) error
	RegisterJsonMessage(codec.JsonCodec, MessageFunc) error
	GetHandler(uint32) (*MessageHandler, error)
}

type defaultTransportRegister struct {
	msgHandler map[uint32]*MessageHandler
}

func (t *defaultTransportRegister) RegisterProtobufMessage(p proto.Message, f MessageFunc) error {
	protoName := p.ProtoReflect().Descriptor().Name()
	id := crc32.ChecksumIEEE([]byte(protoName))
	if _, ok := t.msgHandler[id]; ok {
		return fmt.Errorf("register protobuf message name existed:%s", protoName)
	}

	tp := reflect.TypeOf(p)
	t.msgHandler[id] = &MessageHandler{Name: string(protoName), RType: tp, Fn: f}
	return nil
}

func (t *defaultTransportRegister) RegisterJsonMessage(p codec.JsonCodec, f MessageFunc) error {
	id := crc32.ChecksumIEEE([]byte(p.GetName()))
	if _, ok := t.msgHandler[id]; ok {
		return fmt.Errorf("register json message name existed:%s", p.GetName())
	}

	tp := reflect.TypeOf(p)
	t.msgHandler[id] = &MessageHandler{Name: p.GetName(), RType: tp, Fn: f}
	return nil
}

func (t *defaultTransportRegister) GetHandler(id uint32) (*MessageHandler, error) {
	h, ok := t.msgHandler[id]
	if ok {
		return h, nil
	}
	return nil, fmt.Errorf("error<%w>, message_id<%d>", ErrUnregistedMessage, id)
}
