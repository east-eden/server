package codec

import "reflect"

// Marshaler is a simple encoding interface used for the broker/transport
// where headers are not supported by the underlying implementation.
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, reflect.Type) (interface{}, error)
	Name() string
}
