// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: global/fragment.proto

package global

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type C2S_QueryFragments struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *C2S_QueryFragments) Reset() {
	*x = C2S_QueryFragments{}
	if protoimpl.UnsafeEnabled {
		mi := &file_global_fragment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2S_QueryFragments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2S_QueryFragments) ProtoMessage() {}

func (x *C2S_QueryFragments) ProtoReflect() protoreflect.Message {
	mi := &file_global_fragment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2S_QueryFragments.ProtoReflect.Descriptor instead.
func (*C2S_QueryFragments) Descriptor() ([]byte, []int) {
	return file_global_fragment_proto_rawDescGZIP(), []int{0}
}

type S2C_FragmentsList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Frags []*Fragment `protobuf:"bytes,1,rep,name=frags,proto3" json:"frags,omitempty"`
}

func (x *S2C_FragmentsList) Reset() {
	*x = S2C_FragmentsList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_global_fragment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2C_FragmentsList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2C_FragmentsList) ProtoMessage() {}

func (x *S2C_FragmentsList) ProtoReflect() protoreflect.Message {
	mi := &file_global_fragment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2C_FragmentsList.ProtoReflect.Descriptor instead.
func (*S2C_FragmentsList) Descriptor() ([]byte, []int) {
	return file_global_fragment_proto_rawDescGZIP(), []int{1}
}

func (x *S2C_FragmentsList) GetFrags() []*Fragment {
	if x != nil {
		return x.Frags
	}
	return nil
}

// 碎片更新
type S2C_FragmentsUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Frags []*Fragment `protobuf:"bytes,1,rep,name=frags,proto3" json:"frags,omitempty"`
}

func (x *S2C_FragmentsUpdate) Reset() {
	*x = S2C_FragmentsUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_global_fragment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *S2C_FragmentsUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*S2C_FragmentsUpdate) ProtoMessage() {}

func (x *S2C_FragmentsUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_global_fragment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use S2C_FragmentsUpdate.ProtoReflect.Descriptor instead.
func (*S2C_FragmentsUpdate) Descriptor() ([]byte, []int) {
	return file_global_fragment_proto_rawDescGZIP(), []int{2}
}

func (x *S2C_FragmentsUpdate) GetFrags() []*Fragment {
	if x != nil {
		return x.Frags
	}
	return nil
}

// 碎片合成
type C2S_FragmentsCompose struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FragId int32 `protobuf:"varint,1,opt,name=FragId,proto3" json:"FragId,omitempty"`
}

func (x *C2S_FragmentsCompose) Reset() {
	*x = C2S_FragmentsCompose{}
	if protoimpl.UnsafeEnabled {
		mi := &file_global_fragment_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *C2S_FragmentsCompose) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*C2S_FragmentsCompose) ProtoMessage() {}

func (x *C2S_FragmentsCompose) ProtoReflect() protoreflect.Message {
	mi := &file_global_fragment_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use C2S_FragmentsCompose.ProtoReflect.Descriptor instead.
func (*C2S_FragmentsCompose) Descriptor() ([]byte, []int) {
	return file_global_fragment_proto_rawDescGZIP(), []int{3}
}

func (x *C2S_FragmentsCompose) GetFragId() int32 {
	if x != nil {
		return x.FragId
	}
	return 0
}

var File_global_fragment_proto protoreflect.FileDescriptor

var file_global_fragment_proto_rawDesc = []byte{
	0x0a, 0x15, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2f, 0x66, 0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x1a,
	0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x14, 0x0a, 0x12, 0x43, 0x32, 0x53, 0x5f, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x46, 0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x3b, 0x0a, 0x11, 0x53, 0x32,
	0x43, 0x5f, 0x46, 0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12,
	0x26, 0x0a, 0x05, 0x66, 0x72, 0x61, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x46, 0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x05, 0x66, 0x72, 0x61, 0x67, 0x73, 0x22, 0x3d, 0x0a, 0x13, 0x53, 0x32, 0x43, 0x5f, 0x46,
	0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x26,
	0x0a, 0x05, 0x66, 0x72, 0x61, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x46, 0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x05, 0x66, 0x72, 0x61, 0x67, 0x73, 0x22, 0x2e, 0x0a, 0x14, 0x43, 0x32, 0x53, 0x5f, 0x46, 0x72,
	0x61, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x46, 0x72, 0x61, 0x67, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x46, 0x72, 0x61, 0x67, 0x49, 0x64, 0x42, 0x34, 0x5a, 0x29, 0x62, 0x69, 0x74, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x66, 0x75, 0x6e, 0x70, 0x6c, 0x75, 0x73, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6c, 0x6f,
	0x62, 0x61, 0x6c, 0xaa, 0x02, 0x06, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_global_fragment_proto_rawDescOnce sync.Once
	file_global_fragment_proto_rawDescData = file_global_fragment_proto_rawDesc
)

func file_global_fragment_proto_rawDescGZIP() []byte {
	file_global_fragment_proto_rawDescOnce.Do(func() {
		file_global_fragment_proto_rawDescData = protoimpl.X.CompressGZIP(file_global_fragment_proto_rawDescData)
	})
	return file_global_fragment_proto_rawDescData
}

var file_global_fragment_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_global_fragment_proto_goTypes = []interface{}{
	(*C2S_QueryFragments)(nil),   // 0: global.C2S_QueryFragments
	(*S2C_FragmentsList)(nil),    // 1: global.S2C_FragmentsList
	(*S2C_FragmentsUpdate)(nil),  // 2: global.S2C_FragmentsUpdate
	(*C2S_FragmentsCompose)(nil), // 3: global.C2S_FragmentsCompose
	(*Fragment)(nil),             // 4: global.Fragment
}
var file_global_fragment_proto_depIdxs = []int32{
	4, // 0: global.S2C_FragmentsList.frags:type_name -> global.Fragment
	4, // 1: global.S2C_FragmentsUpdate.frags:type_name -> global.Fragment
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_global_fragment_proto_init() }
func file_global_fragment_proto_init() {
	if File_global_fragment_proto != nil {
		return
	}
	file_global_define_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_global_fragment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2S_QueryFragments); i {
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
		file_global_fragment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2C_FragmentsList); i {
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
		file_global_fragment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*S2C_FragmentsUpdate); i {
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
		file_global_fragment_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*C2S_FragmentsCompose); i {
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
			RawDescriptor: file_global_fragment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_global_fragment_proto_goTypes,
		DependencyIndexes: file_global_fragment_proto_depIdxs,
		MessageInfos:      file_global_fragment_proto_msgTypes,
	}.Build()
	File_global_fragment_proto = out.File
	file_global_fragment_proto_rawDesc = nil
	file_global_fragment_proto_goTypes = nil
	file_global_fragment_proto_depIdxs = nil
}
