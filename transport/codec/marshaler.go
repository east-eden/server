package codec

import "reflect"

// Marshaler is a simple encoding interface used for the broker/transport
// where headers are not supported by the underlying implementation.
type Marshaler interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, reflect.Type) (any, error)
	Name() string
}
