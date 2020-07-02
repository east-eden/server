package codec

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

// test cases
var (
	cases = map[string]struct {
		codec  CodecType
		input  *pbGame.Hero
		option cmp.Option
	}{
		"hero_json": {
			codec: Codec_Json,
			input: &pbGame.Hero{
				Id:     3211475,
				TypeId: 1001,
				Exp:    120,
				Level:  99,
			},
			option: cmp.Comparer(func(x, y *pbGame.Hero) bool {
				if !cmp.Equal(x.Id, y.Id) {
					return false
				}

				if !cmp.Equal(x.TypeId, y.TypeId) {
					return false
				}

				if !cmp.Equal(x.Exp, y.Exp) {
					return false
				}

				if !cmp.Equal(x.Level, y.Level) {
					return false
				}

				return true
			}),
		},

		"hero_protobuf": {
			codec: Codec_Protobuf,
			input: &pbGame.Hero{
				Id:     1928884,
				TypeId: 2001,
				Exp:    99183,
				Level:  13,
			},
			option: cmp.Comparer(func(x, y *pbGame.Hero) bool {
				if !cmp.Equal(x.Id, y.Id) {
					return false
				}

				if !cmp.Equal(x.TypeId, y.TypeId) {
					return false
				}

				if !cmp.Equal(x.Exp, y.Exp) {
					return false
				}

				if !cmp.Equal(x.Level, y.Level) {
					return false
				}

				return true
			}),
		},
	}
)

func TestTransportCodec(t *testing.T) {
	marshalers := []Marshaler{&ProtoBufMarshaler{}, &JsonMarshaler{}}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			data, err := marshalers[Codec_Protobuf].Marshal(cs.input)
			if err != nil {
				t.Fatal(err)
			}

			h, err := marshalers[Codec_Protobuf].Unmarshal(data, reflect.TypeOf(&pbGame.Hero{}))
			if err != nil {
				t.Fatal(err)
			}

			diff := cmp.Diff(cs.input, h, cs.option)
			if diff != "" {
				t.Fatalf(diff)
			}

			data, err = marshalers[Codec_Json].Marshal(cs.input)
			if err != nil {
				t.Fatal(err)
			}

			h, err = marshalers[Codec_Json].Unmarshal(data, reflect.TypeOf(&pbGame.Hero{}))
			if err != nil {
				t.Fatal(err)
			}

			diff = cmp.Diff(cs.input, h, cs.option)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
