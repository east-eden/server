package codec

import (
	"fmt"

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

func (m *ProtoBufMarshaler) Unmarshal(data []byte, v interface{}) (string, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return "", fmt.Errorf("protobuf cannot marshal data:%v", v)
	}

	if err := proto.Unmarshal(data, msg); err != nil {
		return "", fmt.Errorf("protobuf cannot unmarshal data:%v", err)
	}

	name := proto.MessageName(msg)
	return name, nil
}

func (m *ProtoBufMarshaler) String() string {
	return "protobuf"
}
