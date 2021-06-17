package codec

import (
	"fmt"
	"reflect"

	json "github.com/json-iterator/go"
)

type JsonMarshaler struct {
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

func (m *JsonMarshaler) Name() string {
	return "json"
}
