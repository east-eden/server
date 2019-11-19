package codec

type JsonMarshaler struct {
}

func NewJsonCodec() Marshaler {
	return &JsonMarshaler{}
}

func (m *JsonMarshaler) Marshal(interface{}) ([]byte, error) {
	return []byte("success"), nil
}

func (m *JsonMarshaler) Unmarshal([]byte, interface{}) (string, error) {
	return "", nil
}

func (m *JsonMarshaler) String() string {
	return "success"
}
