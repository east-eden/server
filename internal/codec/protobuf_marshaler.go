package codec

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type ProtoBufMarshaler struct {
}

func NewProtobufCodec() Marshaler {
	return &ProtoBufMarshaler{}
}

func (m *ProtoBufMarshaler) Marshal(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("protobuf cannot marshal data:%v", v)
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *ProtoBufMarshaler) Unmarshal(data []byte, name string) (interface{}, error) {
	pType := proto.MessageType(name)
	if pType == nil {
		return nil, fmt.Errorf("protobuf unmarshal failed with name:%s", name)
	}

	// prepare proto struct to be unmarshaled in
	msg, ok := reflect.New(pType.Elem()).Interface().(proto.Message)
	if !ok {
		return nil, fmt.Errorf("protobuf new elem interface failed:%s", name)
	}

	if err := proto.Unmarshal(data, msg); err != nil {
		return "", fmt.Errorf("protobuf cannot unmarshal data:%v", err)
	}

	return msg, nil
}

func (m *ProtoBufMarshaler) String() string {
	return "protobuf"
}
