package codec

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
)

type ProtoBufMarshaler struct {
}

func (m *ProtoBufMarshaler) Marshal(v any) ([]byte, error) {
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

func (m *ProtoBufMarshaler) Unmarshal(data []byte, rtype reflect.Type) (any, error) {
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

func (m *ProtoBufMarshaler) Name() string {
	return "protobuf"
}
