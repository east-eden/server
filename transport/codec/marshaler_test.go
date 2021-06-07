package codec

import (
	"fmt"
	"reflect"
	"testing"

	"e.coding.net/mmstudio/blade/gate/msg"
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
)

// test cases
var (
	cases = map[string]struct {
		codec  CodecType
		input  *pbGlobal.Hero
		option cmp.Option
	}{
		"hero_json": {
			codec: Codec_Json,
			input: &pbGlobal.Hero{
				Id:     3211475,
				TypeId: 1001,
				Exp:    120,
				Level:  99,
			},
			option: cmp.Comparer(func(x, y *pbGlobal.Hero) bool {
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
			input: &pbGlobal.Hero{
				Id:     1928884,
				TypeId: 2001,
				Exp:    99183,
				Level:  13,
			},
			option: cmp.Comparer(func(x, y *pbGlobal.Hero) bool {
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

			h, err := marshalers[Codec_Protobuf].Unmarshal(data, reflect.TypeOf(&pbGlobal.Hero{}))
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

			h, err = marshalers[Codec_Json].Unmarshal(data, reflect.TypeOf(&pbGlobal.Hero{}))
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

func TestProtoUnmarshal(t *testing.T) {
	handshake := &msg.Handshake{
		Cmd:          msg.CmdType_NEW,
		Src:          msg.SrcType_CLIENT,
		ServiceName:  "game",
		ClientAddr:   "127.0.0.1:13888",
		UserID:       "user_id1",
		ClientVer:    "0.0.1",
		ClientResVer: "0.0.1",
		Meta:         make([]*msg.MapFieldEntry, 0),
	}

	data, err := proto.Marshal(handshake)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(data)

	data2 := []byte{26, 4, 103, 97, 109, 101, 58, 15, 49, 50, 55, 46, 48, 46, 48, 46, 49, 58, 49, 51, 56, 56, 56, 66, 8, 117, 115, 101, 114, 95, 105, 100, 49, 74, 5, 48, 46, 48, 46, 49, 82, 5, 48, 46, 48, 46, 49}
	handshake2 := &msg.Handshake{}
	err = proto.Unmarshal(data2, handshake2)
	if err != nil {
		t.Fatal(err)
	}
}
