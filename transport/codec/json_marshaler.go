package codec

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type JsonMarshaler struct {
}

func NewJsonCodec() Marshaler {
	return &JsonMarshaler{}
}

func (m *JsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *JsonMarshaler) Unmarshal(data []byte, rtype reflect.Type) (interface{}, error) {
	// prepare proto struct to be unmarshaled in
	msg := reflect.New(rtype.Elem()).Interface()
	if err := json.Unmarshal(data, msg); err != nil {
		return "", fmt.Errorf("json cannot unmarshal data:%v", err)
	}

	return msg, nil
}

func (m *JsonMarshaler) String() string {
	return "json"
}
