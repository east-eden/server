package protobuf

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"

	"e.coding.net/mmstudio/blade/golib/encoding2"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var Codec = &protoCodec{usingPool: false, name: "proto"}
var CodecUsingPool = &protoCodec{usingPool: true, name: "proto_pool"}

func init() {
	encoding2.RegisterCodec(Codec)
	encoding2.RegisterCodec(CodecUsingPool)
}

// protoCodec is a Codec implementation with protobuf. It is the default codec.
type protoCodec struct {
	usingPool bool
	name      string
}

func (p *protoCodec) Marshal(v interface{}) ([]byte, error) {
	if pm, ok := v.(proto.Marshaler); ok {
		// object can marshal itself, no need for buffer
		return pm.Marshal()
	}
	if pm, ok := v.(proto.Message); ok {
		if p.usingPool {
			cb := protoBufferPool.Get().(*cachedProtoBuffer)
			out, err := marshal(pm, cb)
			// put back buffer and lose the ref to the slice
			cb.SetBuf(nil)
			protoBufferPool.Put(cb)
			return out, err
		}
		return proto.Marshal(pm)
	}
	return nil, fmt.Errorf("%T is not a proto.Marshaler", v)
}

func (protoCodec) Uri(t interface{}) string     { return proto.MessageName(t.(proto.Message)) }
func (protoCodec) Type(uri string) reflect.Type { return proto.MessageType(uri) }

func (p *protoCodec) Unmarshal(data []byte, v interface{}) error {
	if pu, ok := v.(proto.Unmarshaler); ok {
		// object can unmarshal itself, no need for buffer
		return pu.Unmarshal(data)
	}

	if m, ok := v.(proto.Message); ok {
		m.Reset()
		if p.usingPool {
			cb := protoBufferPool.Get().(*cachedProtoBuffer)
			cb.SetBuf(data)
			err := cb.Unmarshal(m)
			cb.SetBuf(nil)
			protoBufferPool.Put(cb)
			return err
		}
		return proto.Unmarshal(data, m)
	}

	return fmt.Errorf("%T is not a proto.Unmarshaler", v)
}

func (protoCodec) JSONMarshal(obj interface{}) ([]byte, error) {
	if pm, ok := obj.(proto.Message); ok {
		m := jsonpb.Marshaler{EmitDefaults: false}
		var buf bytes.Buffer
		if err := m.Marshal(&buf, pm); err != nil {
			return buf.Bytes(), err
		}
	}
	return nil, errors.New("not proto message")
}

func (p *protoCodec) Name() string { return p.name }

func marshal(pm proto.Message, cb *cachedProtoBuffer) ([]byte, error) {
	newSlice := make([]byte, 0, cb.lastMarshaledSize)

	cb.SetBuf(newSlice)
	cb.Reset()
	if err := cb.Marshal(pm); err != nil {
		return nil, err
	}
	out := cb.Bytes()
	cb.lastMarshaledSize = capToMaxInt32(len(out))
	return out, nil
}

func capToMaxInt32(val int) uint32 {
	if val > math.MaxInt32 {
		return uint32(math.MaxInt32)
	}
	return uint32(val)
}

type cachedProtoBuffer struct {
	lastMarshaledSize uint32
	proto.Buffer
}

var protoBufferPool = &sync.Pool{
	New: func() interface{} {
		return &cachedProtoBuffer{
			Buffer:            proto.Buffer{},
			lastMarshaledSize: 16,
		}
	},
}
