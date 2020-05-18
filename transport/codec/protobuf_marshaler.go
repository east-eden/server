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

func (m *ProtoBufMarshaler) Unmarshal(data []byte, rtype reflect.Type) (interface{}, error) {
	// prepare proto struct to be unmarshaled in
	msg, ok := reflect.New(rtype.Elem()).Interface().(proto.Message)
	if !ok {
		return nil, fmt.Errorf("protobuf new elem interface failed:%v", rtype)
	}

	if err := proto.Unmarshal(data, msg); err != nil {
		return "", fmt.Errorf("protobuf cannot unmarshal data:%v", err)
	}

	return msg, nil
}

func (m *ProtoBufMarshaler) String() string {
	return "protobuf"
}
