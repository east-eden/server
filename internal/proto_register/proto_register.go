package proto_register

import (
	"context"
	"errors"
	"reflect"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

var (
	errHandlerNotRegisted     = errors.New("handler not registed")
	errReflectNewObjectFailed = errors.New("reflect new object failed")
	errProtoUnmarshalFailed   = errors.New("proto unmarshal failed")
)

type ProtoHandler func(context.Context, proto.Message) (proto.Message, error)

type ProtoRegister struct {
	handlers map[string]ProtoHandler
	rTypes   map[string]reflect.Type
}

func NewProtoRegister() *ProtoRegister {
	return &ProtoRegister{
		handlers: make(map[string]ProtoHandler),
		rTypes:   make(map[string]reflect.Type),
	}
}

func (p *ProtoRegister) RegisterProto(msg proto.Message, fn ProtoHandler) {
	if p.handlers == nil {
		p.handlers = make(map[string]ProtoHandler)
	}

	if p.rTypes == nil {
		p.rTypes = make(map[string]reflect.Type)
	}

	msgName := string(proto.MessageName(msg))
	p.rTypes[msgName] = reflect.TypeOf(msg)
	p.handlers[msgName] = fn
}

func (p *ProtoRegister) GetRegisterProto(pbName string) (reflect.Type, ProtoHandler, bool) {
	rtype, found := p.rTypes[pbName]
	if !found {
		return nil, nil, false
	}

	fn, found := p.handlers[pbName]
	if !found {
		return rtype, nil, true
	}

	return rtype, fn, true
}

func (p *ProtoRegister) UnmarshalRegisterProto(pbName string, pbData []byte) proto.Message {
	rtype, _, found := p.GetRegisterProto(pbName)
	if !found {
		return nil
	}

	msg, ok := reflect.New(rtype.Elem()).Interface().(proto.Message)
	if !ok {
		return nil
	}

	if err := proto.Unmarshal(pbData, msg); err != nil {
		return nil
	}

	return msg
}

func (p *ProtoRegister) Execute(ctx context.Context, pbName string, pbData []byte) (proto.Message, error) {
	rtype, fn, found := p.GetRegisterProto(pbName)
	if !found {
		log.Info().Str("name", pbName).Msg("cannot find ProtoRegister's registered rtype or handler")
		return nil, errHandlerNotRegisted
	}

	msg, ok := reflect.New(rtype.Elem()).Interface().(proto.Message)
	if !ok {
		log.Error().Str("name", pbName).Msg("reflect new proto message failed")
		return nil, errReflectNewObjectFailed
	}

	if err := proto.Unmarshal(pbData, msg); err != nil {
		log.Error().Str("name", pbName).Msg("proto.Unmarshal failed")
		return nil, errProtoUnmarshalFailed
	}

	if fn == nil {
		log.Error().Str("name", pbName).Msg("handler is nil")
		return nil, errHandlerNotRegisted
	}

	return fn(ctx, msg)
}
