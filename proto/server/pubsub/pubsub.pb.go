// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: server/pubsub/pubsub.proto

package pubsub

import (
	global "e.coding.net/mmstudio/blade/server/proto/global"
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

/////////////////////////////////////////////////
// pub/sub
/////////////////////////////////////////////////
type PubStartGate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Info *global.AccountInfo `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *PubStartGate) Reset() {
	*x = PubStartGate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_pubsub_pubsub_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PubStartGate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PubStartGate) ProtoMessage() {}

func (x *PubStartGate) ProtoReflect() protoreflect.Message {
	mi := &file_server_pubsub_pubsub_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PubStartGate.ProtoReflect.Descriptor instead.
func (*PubStartGate) Descriptor() ([]byte, []int) {
	return file_server_pubsub_pubsub_proto_rawDescGZIP(), []int{0}
}

func (x *PubStartGate) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PubStartGate) GetInfo() *global.AccountInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type PubGateResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Info *global.AccountInfo `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	Win  bool                `protobuf:"varint,3,opt,name=win,proto3" json:"win,omitempty"`
}

func (x *PubGateResult) Reset() {
	*x = PubGateResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_pubsub_pubsub_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PubGateResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PubGateResult) ProtoMessage() {}

func (x *PubGateResult) ProtoReflect() protoreflect.Message {
	mi := &file_server_pubsub_pubsub_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PubGateResult.ProtoReflect.Descriptor instead.
func (*PubGateResult) Descriptor() ([]byte, []int) {
	return file_server_pubsub_pubsub_proto_rawDescGZIP(), []int{1}
}

func (x *PubGateResult) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PubGateResult) GetInfo() *global.AccountInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

func (x *PubGateResult) GetWin() bool {
	if x != nil {
		return x.Win
	}
	return false
}

type PubSyncPlayerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64              `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Info *global.PlayerInfo `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *PubSyncPlayerInfo) Reset() {
	*x = PubSyncPlayerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_pubsub_pubsub_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PubSyncPlayerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PubSyncPlayerInfo) ProtoMessage() {}

func (x *PubSyncPlayerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_server_pubsub_pubsub_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PubSyncPlayerInfo.ProtoReflect.Descriptor instead.
func (*PubSyncPlayerInfo) Descriptor() ([]byte, []int) {
	return file_server_pubsub_pubsub_proto_rawDescGZIP(), []int{2}
}

func (x *PubSyncPlayerInfo) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PubSyncPlayerInfo) GetInfo() *global.PlayerInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type MultiPublishTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *MultiPublishTest) Reset() {
	*x = MultiPublishTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_pubsub_pubsub_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MultiPublishTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MultiPublishTest) ProtoMessage() {}

func (x *MultiPublishTest) ProtoReflect() protoreflect.Message {
	mi := &file_server_pubsub_pubsub_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MultiPublishTest.ProtoReflect.Descriptor instead.
func (*MultiPublishTest) Descriptor() ([]byte, []int) {
	return file_server_pubsub_pubsub_proto_rawDescGZIP(), []int{3}
}

func (x *MultiPublishTest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MultiPublishTest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_server_pubsub_pubsub_proto protoreflect.FileDescriptor

var file_server_pubsub_pubsub_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2f,
	0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x75,
	0x62, 0x73, 0x75, 0x62, 0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2f, 0x64, 0x65, 0x66,
	0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x0c, 0x50, 0x75, 0x62,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x47, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x26, 0x0a, 0x04, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66,
	0x6f, 0x22, 0x59, 0x0a, 0x0d, 0x50, 0x75, 0x62, 0x47, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x26, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x77, 0x69,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x77, 0x69, 0x6e, 0x22, 0x4a, 0x0a, 0x11,
	0x50, 0x75, 0x62, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x25, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x36, 0x0a, 0x10, 0x4d, 0x75, 0x6c, 0x74,
	0x69, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x54, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x42, 0x38, 0x5a, 0x36, 0x65, 0x2e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x6e, 0x65, 0x74,
	0x2f, 0x6d, 0x6d, 0x73, 0x74, 0x75, 0x64, 0x69, 0x6f, 0x2f, 0x62, 0x6c, 0x61, 0x64, 0x65, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2f, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_server_pubsub_pubsub_proto_rawDescOnce sync.Once
	file_server_pubsub_pubsub_proto_rawDescData = file_server_pubsub_pubsub_proto_rawDesc
)

func file_server_pubsub_pubsub_proto_rawDescGZIP() []byte {
	file_server_pubsub_pubsub_proto_rawDescOnce.Do(func() {
		file_server_pubsub_pubsub_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_pubsub_pubsub_proto_rawDescData)
	})
	return file_server_pubsub_pubsub_proto_rawDescData
}

var file_server_pubsub_pubsub_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_server_pubsub_pubsub_proto_goTypes = []interface{}{
	(*PubStartGate)(nil),       // 0: pubsub.PubStartGate
	(*PubGateResult)(nil),      // 1: pubsub.PubGateResult
	(*PubSyncPlayerInfo)(nil),  // 2: pubsub.PubSyncPlayerInfo
	(*MultiPublishTest)(nil),   // 3: pubsub.MultiPublishTest
	(*global.AccountInfo)(nil), // 4: proto.AccountInfo
	(*global.PlayerInfo)(nil),  // 5: proto.PlayerInfo
}
var file_server_pubsub_pubsub_proto_depIdxs = []int32{
	4, // 0: pubsub.PubStartGate.info:type_name -> proto.AccountInfo
	4, // 1: pubsub.PubGateResult.info:type_name -> proto.AccountInfo
	5, // 2: pubsub.PubSyncPlayerInfo.info:type_name -> proto.PlayerInfo
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_server_pubsub_pubsub_proto_init() }
func file_server_pubsub_pubsub_proto_init() {
	if File_server_pubsub_pubsub_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_pubsub_pubsub_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PubStartGate); i {
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
		file_server_pubsub_pubsub_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PubGateResult); i {
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
		file_server_pubsub_pubsub_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PubSyncPlayerInfo); i {
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
		file_server_pubsub_pubsub_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MultiPublishTest); i {
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
			RawDescriptor: file_server_pubsub_pubsub_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_server_pubsub_pubsub_proto_goTypes,
		DependencyIndexes: file_server_pubsub_pubsub_proto_depIdxs,
		MessageInfos:      file_server_pubsub_pubsub_proto_msgTypes,
	}.Build()
	File_server_pubsub_pubsub_proto = out.File
	file_server_pubsub_pubsub_proto_rawDesc = nil
	file_server_pubsub_pubsub_proto_goTypes = nil
	file_server_pubsub_pubsub_proto_depIdxs = nil
}
