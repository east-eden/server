// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.14.0
// source: global/player.proto

package global

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 创建角色
type C2S_CreatePlayer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *C2S_CreatePlayer) Reset() {
	*x = C2S_CreatePlayer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_global_player_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2S_CreatePlayer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2S_CreatePlayer) ProtoMessage() {}

func (x *C2S_CreatePlayer) ProtoReflect() protoreflect.Message {
	mi := &file_global_player_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2S_CreatePlayer.ProtoReflect.Descriptor instead.
func (*C2S_CreatePlayer) Descriptor() ([]byte, []int) {
	return file_global_player_proto_rawDescGZIP(), []int{0}
}

func (x *C2S_CreatePlayer) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type S2C_CreatePlayer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *PlayerInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *S2C_CreatePlayer) Reset() {
	*x = S2C_CreatePlayer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_global_player_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2C_CreatePlayer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2C_CreatePlayer) ProtoMessage() {}

func (x *S2C_CreatePlayer) ProtoReflect() protoreflect.Message {
	mi := &file_global_player_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2C_CreatePlayer.ProtoReflect.Descriptor instead.
func (*S2C_CreatePlayer) Descriptor() ([]byte, []int) {
	return file_global_player_proto_rawDescGZIP(), []int{1}
}

func (x *S2C_CreatePlayer) GetInfo() *PlayerInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

// 玩家登陆初始信息
type S2C_PlayerInitInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info     *PlayerInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`         // 玩家基本信息
	Heros    []*Hero     `protobuf:"bytes,2,rep,name=heros,proto3" json:"heros,omitempty"`       // 英雄数据
	Items    []*Item     `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`       // 物品数据
	Equips   []*Equip    `protobuf:"bytes,4,rep,name=equips,proto3" json:"equips,omitempty"`     // 装备数据
	Crystals []*Crystal  `protobuf:"bytes,5,rep,name=crystals,proto3" json:"crystals,omitempty"` // 晶石数据
	Frags    []*Fragment `protobuf:"bytes,6,rep,name=frags,proto3" json:"frags,omitempty"`       // 碎片数据
}

func (x *S2C_PlayerInitInfo) Reset() {
	*x = S2C_PlayerInitInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_global_player_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2C_PlayerInitInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2C_PlayerInitInfo) ProtoMessage() {}

func (x *S2C_PlayerInitInfo) ProtoReflect() protoreflect.Message {
	mi := &file_global_player_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2C_PlayerInitInfo.ProtoReflect.Descriptor instead.
func (*S2C_PlayerInitInfo) Descriptor() ([]byte, []int) {
	return file_global_player_proto_rawDescGZIP(), []int{2}
}

func (x *S2C_PlayerInitInfo) GetInfo() *PlayerInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

func (x *S2C_PlayerInitInfo) GetHeros() []*Hero {
	if x != nil {
		return x.Heros
	}
	return nil
}

func (x *S2C_PlayerInitInfo) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *S2C_PlayerInitInfo) GetEquips() []*Equip {
	if x != nil {
		return x.Equips
	}
	return nil
}

func (x *S2C_PlayerInitInfo) GetCrystals() []*Crystal {
	if x != nil {
		return x.Crystals
	}
	return nil
}

func (x *S2C_PlayerInitInfo) GetFrags() []*Fragment {
	if x != nil {
		return x.Frags
	}
	return nil
}

type S2C_ExpUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exp   int32 `protobuf:"varint,1,opt,name=Exp,proto3" json:"Exp,omitempty"`
	Level int32 `protobuf:"varint,2,opt,name=Level,proto3" json:"Level,omitempty"`
}

func (x *S2C_ExpUpdate) Reset() {
	*x = S2C_ExpUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_global_player_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2C_ExpUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2C_ExpUpdate) ProtoMessage() {}

func (x *S2C_ExpUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_global_player_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2C_ExpUpdate.ProtoReflect.Descriptor instead.
func (*S2C_ExpUpdate) Descriptor() ([]byte, []int) {
	return file_global_player_proto_rawDescGZIP(), []int{3}
}

func (x *S2C_ExpUpdate) GetExp() int32 {
	if x != nil {
		return x.Exp
	}
	return 0
}

func (x *S2C_ExpUpdate) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

// gm 命令:
// gm player level 10
// gm hero add 1
// gm item add 1 10
type C2S_GmCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd string `protobuf:"bytes,1,opt,name=cmd,proto3" json:"cmd,omitempty"`
}

func (x *C2S_GmCmd) Reset() {
	*x = C2S_GmCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_global_player_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2S_GmCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2S_GmCmd) ProtoMessage() {}

func (x *C2S_GmCmd) ProtoReflect() protoreflect.Message {
	mi := &file_global_player_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2S_GmCmd.ProtoReflect.Descriptor instead.
func (*C2S_GmCmd) Descriptor() ([]byte, []int) {
	return file_global_player_proto_rawDescGZIP(), []int{4}
}

func (x *C2S_GmCmd) GetCmd() string {
	if x != nil {
		return x.Cmd
	}
	return ""
}

var File_global_player_proto protoreflect.FileDescriptor

var file_global_player_proto_rawDesc = []byte{
	0x0a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x1a, 0x13, 0x67,
	0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x26, 0x0a, 0x10, 0x43, 0x32, 0x53, 0x5f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x0a, 0x10, 0x53, 0x32,
	0x43, 0x5f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x26,
	0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67,
	0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x80, 0x02, 0x0a, 0x12, 0x53, 0x32, 0x43, 0x5f, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x69, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x26, 0x0a,
	0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x6c,
	0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x04, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x22, 0x0a, 0x05, 0x68, 0x65, 0x72, 0x6f, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x48, 0x65,
	0x72, 0x6f, 0x52, 0x05, 0x68, 0x65, 0x72, 0x6f, 0x73, 0x12, 0x22, 0x0a, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61,
	0x6c, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x25, 0x0a,
	0x06, 0x65, 0x71, 0x75, 0x69, 0x70, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x45, 0x71, 0x75, 0x69, 0x70, 0x52, 0x06, 0x65, 0x71,
	0x75, 0x69, 0x70, 0x73, 0x12, 0x2b, 0x0a, 0x08, 0x63, 0x72, 0x79, 0x73, 0x74, 0x61, 0x6c, 0x73,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e,
	0x43, 0x72, 0x79, 0x73, 0x74, 0x61, 0x6c, 0x52, 0x08, 0x63, 0x72, 0x79, 0x73, 0x74, 0x61, 0x6c,
	0x73, 0x12, 0x26, 0x0a, 0x05, 0x66, 0x72, 0x61, 0x67, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x46, 0x72, 0x61, 0x67, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x05, 0x66, 0x72, 0x61, 0x67, 0x73, 0x22, 0x37, 0x0a, 0x0d, 0x53, 0x32, 0x43,
	0x5f, 0x45, 0x78, 0x70, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x45, 0x78,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x45, 0x78, 0x70, 0x12, 0x14, 0x0a, 0x05,
	0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x4c, 0x65, 0x76,
	0x65, 0x6c, 0x22, 0x1d, 0x0a, 0x09, 0x43, 0x32, 0x53, 0x5f, 0x47, 0x6d, 0x43, 0x6d, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x6d,
	0x64, 0x42, 0x34, 0x5a, 0x29, 0x62, 0x69, 0x74, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x6f,
	0x72, 0x67, 0x2f, 0x66, 0x75, 0x6e, 0x70, 0x6c, 0x75, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0xaa, 0x02,
	0x06, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_global_player_proto_rawDescOnce sync.Once
	file_global_player_proto_rawDescData = file_global_player_proto_rawDesc
)

func file_global_player_proto_rawDescGZIP() []byte {
	file_global_player_proto_rawDescOnce.Do(func() {
		file_global_player_proto_rawDescData = protoimpl.X.CompressGZIP(file_global_player_proto_rawDescData)
	})
	return file_global_player_proto_rawDescData
}

var file_global_player_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_global_player_proto_goTypes = []interface{}{
	(*C2S_CreatePlayer)(nil),   // 0: global.C2S_CreatePlayer
	(*S2C_CreatePlayer)(nil),   // 1: global.S2C_CreatePlayer
	(*S2C_PlayerInitInfo)(nil), // 2: global.S2C_PlayerInitInfo
	(*S2C_ExpUpdate)(nil),      // 3: global.S2C_ExpUpdate
	(*C2S_GmCmd)(nil),          // 4: global.C2S_GmCmd
	(*PlayerInfo)(nil),         // 5: global.PlayerInfo
	(*Hero)(nil),               // 6: global.Hero
	(*Item)(nil),               // 7: global.Item
	(*Equip)(nil),              // 8: global.Equip
	(*Crystal)(nil),            // 9: global.Crystal
	(*Fragment)(nil),           // 10: global.Fragment
}
var file_global_player_proto_depIdxs = []int32{
	5,  // 0: global.S2C_CreatePlayer.info:type_name -> global.PlayerInfo
	5,  // 1: global.S2C_PlayerInitInfo.info:type_name -> global.PlayerInfo
	6,  // 2: global.S2C_PlayerInitInfo.heros:type_name -> global.Hero
	7,  // 3: global.S2C_PlayerInitInfo.items:type_name -> global.Item
	8,  // 4: global.S2C_PlayerInitInfo.equips:type_name -> global.Equip
	9,  // 5: global.S2C_PlayerInitInfo.crystals:type_name -> global.Crystal
	10, // 6: global.S2C_PlayerInitInfo.frags:type_name -> global.Fragment
	7,  // [7:7] is the sub-list for method output_type
	7,  // [7:7] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_global_player_proto_init() }
func file_global_player_proto_init() {
	if File_global_player_proto != nil {
		return
	}
	file_global_define_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_global_player_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2S_CreatePlayer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_global_player_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2C_CreatePlayer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_global_player_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2C_PlayerInitInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_global_player_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2C_ExpUpdate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_global_player_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2S_GmCmd); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_global_player_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_global_player_proto_goTypes,
		DependencyIndexes: file_global_player_proto_depIdxs,
		MessageInfos:      file_global_player_proto_msgTypes,
	}.Build()
	File_global_player_proto = out.File
	file_global_player_proto_rawDesc = nil
	file_global_player_proto_goTypes = nil
	file_global_player_proto_depIdxs = nil
}
