package encoding2

import (
	"context"
)

type codecKeyType struct{}

// for ark context
func (*codecKeyType) String() string { return "encoding2-codec—key" }

var KeyForContext = codecKeyType(struct{}{})

func WithContext(ctx context.Context, c Codec) context.Context {
	return context.WithValue(ctx, KeyForContext, c)
}

func FromContext(ctx context.Context) Codec {
	c := ctx.Value(KeyForContext)
	if c == nil {
		c = ctx.Value(KeyForContext.String())
	}
	if c == nil {
		return nil
	}
	return c.(Codec)
}

// Codec defines the interface link uses to encode and decode messages.
type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// String returns the name of the Codec implementation.
	Name() string
}

var registeredCodecs = make(map[string]Codec)

func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	registeredCodecs[codec.Name()] = codec
}

func GetCodec(contentSubtype string) Codec {
	return registeredCodecs[contentSubtype]
}
