package bsonproto

import (
	"testing"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// roundTrip 走一遍真实的 mongodb-driver BSON 编解码，确保 result 与输入语义等价
func roundTrip(t *testing.T, in proto.Message, out proto.Message) bson.D {
	t.Helper()

	doc, err := Marshal(in)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	raw, err := bson.Marshal(doc)
	if err != nil {
		t.Fatalf("bson.Marshal failed: %v", err)
	}

	var decoded bson.D
	if err := bson.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("bson.Unmarshal failed: %v", err)
	}

	if err := Unmarshal(decoded, out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	return decoded
}

func TestSnakeCase(t *testing.T) {
	cases := map[string]string{
		"Id":                 "id",
		"AccountId":          "account_id",
		"BuyStrengthenTimes": "buy_strengthen_times",
		"TypeId":             "type_id",
		"Name":               "name",
	}

	for in, want := range cases {
		if got := SnakeCase(in); got != want {
			t.Fatalf("SnakeCase(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMarshalScalarAndRepeated(t *testing.T) {
	in := &pbGlobal.PlayerInfo{
		Id:                 1001100,
		AccountId:          2002200,
		Name:               "加百列",
		Exp:                123,
		Level:              66,
		BuyStrengthenTimes: 3,
		BattleArray:        []int64{1, 2, 3},
	}

	out := &pbGlobal.PlayerInfo{}
	doc := roundTrip(t, in, out)

	if !proto.Equal(in, out) {
		t.Fatalf("round trip mismatch: %s", cmp.Diff(in.String(), out.String()))
	}

	// 主键必须落在 _id 上
	var idFound bool
	keys := make([]string, 0, len(doc))
	for _, e := range doc {
		keys = append(keys, e.Key)
		if e.Key == IDField {
			idFound = true
			if got, ok := e.Value.(int64); !ok || got != in.Id {
				t.Fatalf("_id = %v (%T), want %d", e.Value, e.Value, in.Id)
			}
		}
	}

	if !idFound {
		t.Fatalf("_id not found in document keys %v", keys)
	}

	// 其余字段用 snake_case
	for _, want := range []string{"account_id", "buy_strengthen_times", "battle_array"} {
		var found bool
		for _, k := range keys {
			if k == want {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("field %q not found in document keys %v", want, keys)
		}
	}
}

func TestMarshalNestedMessageAndEnum(t *testing.T) {
	in := &pbGlobal.S2C_PlayerInitInfo{
		Info: &pbGlobal.PlayerInfo{
			Id:        9001,
			AccountId: 7001,
			Name:      "加百列",
		},
		Items: []*pbGlobal.Item{
			{Id: 1, TypeId: 10100000, Num: 100},
			{Id: 2, TypeId: 10200000, Num: 5},
		},
	}

	out := &pbGlobal.S2C_PlayerInitInfo{}
	doc := roundTrip(t, in, out)

	if !proto.Equal(in, out) {
		t.Fatalf("round trip mismatch: %s", cmp.Diff(in.String(), out.String()))
	}

	// 嵌套消息应保持为子文档
	var infoDoc bson.D
	for _, e := range doc {
		if e.Key != "info" {
			continue
		}

		d, ok := e.Value.(bson.D)
		if !ok {
			t.Fatalf("info is %T, want bson.D", e.Value)
		}

		infoDoc = d
	}

	if infoDoc == nil {
		t.Fatalf("nested field info not found in %v", doc)
	}

	// 嵌套层的主键不应被改写为 _id
	for _, e := range infoDoc {
		if e.Key == IDField {
			t.Fatalf("nested id should not become _id")
		}
	}

	// 枚举字段
	tokenIn := &pbGlobal.Handshake{Cmd: pbGlobal.CmdType_NEW, ClientAddr: "127.0.0.1:1"}
	tokenOut := &pbGlobal.Handshake{}
	roundTrip(t, tokenIn, tokenOut)

	if !proto.Equal(tokenIn, tokenOut) {
		t.Fatalf("enum round trip mismatch: %s", cmp.Diff(tokenIn.String(), tokenOut.String()))
	}
}

func TestMarshalMap(t *testing.T) {
	md := mapDescriptor(t, "StringKeyMap", descriptorpb.FieldDescriptorProto_TYPE_STRING, descriptorpb.FieldDescriptorProto_TYPE_STRING)
	field := md.Fields().ByName("Items")

	in := dynamicpb.NewMessage(md)
	mp := in.Mutable(field).Map()
	mp.Set(protoreflect.ValueOfString("trace").MapKey(), protoreflect.ValueOfString("abc"))
	mp.Set(protoreflect.ValueOfString("zone").MapKey(), protoreflect.ValueOfString("1"))

	out := dynamicpb.NewMessage(md)
	roundTrip(t, in, out)

	if !proto.Equal(in, out) {
		t.Fatalf("map round trip mismatch: %s", cmp.Diff(in.String(), out.String()))
	}
}

func TestMarshalSkipsUnsetFields(t *testing.T) {
	// proto3 零值不写库，保证文档只保存真实数据
	doc, err := Marshal(&pbGlobal.PlayerInfo{Id: 1})
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if len(doc) != 1 || doc[0].Key != IDField {
		t.Fatalf("want only _id, got %v", doc)
	}
}

func TestUnmarshalIgnoresUnknownFields(t *testing.T) {
	doc := bson.D{
		{Key: IDField, Value: int64(42)},
		{Key: "removed_field", Value: "legacy"},
	}

	out := &pbGlobal.PlayerInfo{}
	if err := Unmarshal(doc, out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if out.Id != 42 {
		t.Fatalf("Id = %d, want 42", out.Id)
	}
}

// intKeyMapDescriptor 构造一个 map<int64, int32> 的消息描述符，
// 用来验证非 string 键的 map 能正确往返（kof 的邮件/英雄数据结构大量使用 map[int64]）。
func mapDescriptor(t *testing.T, msgName string, keyType, valueType descriptorpb.FieldDescriptorProto_Type) protoreflect.MessageDescriptor {
	t.Helper()

	fdp := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("bsonproto/test/" + msgName + ".proto"),
		Package: proto.String("bsonprototest"),
		Syntax:  proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{{
			Name: proto.String(msgName),
			Field: []*descriptorpb.FieldDescriptorProto{{
				Name:     proto.String("Items"),
				Number:   proto.Int32(1),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
				TypeName: proto.String(".bsonprototest." + msgName + ".ItemsEntry"),
				JsonName: proto.String("Items"),
			}},
			NestedType: []*descriptorpb.DescriptorProto{{
				Name: proto.String("ItemsEntry"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("key"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     keyType.Enum(),
						JsonName: proto.String("key"),
					},
					{
						Name:     proto.String("value"),
						Number:   proto.Int32(2),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     valueType.Enum(),
						JsonName: proto.String("value"),
					},
				},
				Options: &descriptorpb.MessageOptions{MapEntry: proto.Bool(true)},
			}},
		}},
	}

	file, err := protodesc.NewFile(fdp, nil)
	if err != nil {
		t.Fatalf("build descriptor failed: %v", err)
	}

	return file.Messages().Get(0)
}

func TestMarshalMapWithIntKey(t *testing.T) {
	md := mapDescriptor(t, "IntKeyMap", descriptorpb.FieldDescriptorProto_TYPE_INT64, descriptorpb.FieldDescriptorProto_TYPE_INT32)
	field := md.Fields().ByName("Items")

	in := dynamicpb.NewMessage(md)
	mp := in.Mutable(field).Map()
	mp.Set(protoreflect.ValueOfInt64(10100000).MapKey(), protoreflect.ValueOfInt32(5))
	mp.Set(protoreflect.ValueOfInt64(-1).MapKey(), protoreflect.ValueOfInt32(7))

	out := dynamicpb.NewMessage(md)
	roundTrip(t, in, out)

	if !proto.Equal(in, out) {
		t.Fatalf("int key map round trip mismatch: %s", cmp.Diff(in.String(), out.String()))
	}
}
