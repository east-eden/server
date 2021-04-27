// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: server/mail/mail.proto

package mail

import (
	common "bitbucket.org/funplus/server/proto/global/common"
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

type CreateMailRq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReceiverId  int64              `protobuf:"varint,1,opt,name=ReceiverId,proto3" json:"ReceiverId,omitempty"`          // 收件人id
	SenderId    int64              `protobuf:"varint,2,opt,name=SenderId,proto3" json:"SenderId,omitempty"`              // 发件人id
	Type        common.MailType    `protobuf:"varint,3,opt,name=Type,proto3,enum=common.MailType" json:"Type,omitempty"` // 邮件类型
	SenderName  string             `protobuf:"bytes,4,opt,name=SenderName,proto3" json:"SenderName,omitempty"`           // 发件人姓名
	Title       string             `protobuf:"bytes,5,opt,name=Title,proto3" json:"Title,omitempty"`                     // 邮件标题
	Content     string             `protobuf:"bytes,6,opt,name=Content,proto3" json:"Content,omitempty"`                 // 邮件内容
	Attachments []*common.LootData `protobuf:"bytes,7,rep,name=Attachments,proto3" json:"Attachments,omitempty"`         // 邮件附件
}

func (x *CreateMailRq) Reset() {
	*x = CreateMailRq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMailRq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMailRq) ProtoMessage() {}

func (x *CreateMailRq) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMailRq.ProtoReflect.Descriptor instead.
func (*CreateMailRq) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{0}
}

func (x *CreateMailRq) GetReceiverId() int64 {
	if x != nil {
		return x.ReceiverId
	}
	return 0
}

func (x *CreateMailRq) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *CreateMailRq) GetType() common.MailType {
	if x != nil {
		return x.Type
	}
	return common.MailType_System
}

func (x *CreateMailRq) GetSenderName() string {
	if x != nil {
		return x.SenderName
	}
	return ""
}

func (x *CreateMailRq) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreateMailRq) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *CreateMailRq) GetAttachments() []*common.LootData {
	if x != nil {
		return x.Attachments
	}
	return nil
}

type CreateMailRs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MailId int64 `protobuf:"varint,1,opt,name=MailId,proto3" json:"MailId,omitempty"` // 邮件id
}

func (x *CreateMailRs) Reset() {
	*x = CreateMailRs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMailRs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMailRs) ProtoMessage() {}

func (x *CreateMailRs) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMailRs.ProtoReflect.Descriptor instead.
func (*CreateMailRs) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{1}
}

func (x *CreateMailRs) GetMailId() int64 {
	if x != nil {
		return x.MailId
	}
	return 0
}

type QueryPlayerMailsRq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerId int64 `protobuf:"varint,1,opt,name=OwnerId,proto3" json:"OwnerId,omitempty"` // 邮箱主人id
}

func (x *QueryPlayerMailsRq) Reset() {
	*x = QueryPlayerMailsRq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryPlayerMailsRq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryPlayerMailsRq) ProtoMessage() {}

func (x *QueryPlayerMailsRq) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryPlayerMailsRq.ProtoReflect.Descriptor instead.
func (*QueryPlayerMailsRq) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{2}
}

func (x *QueryPlayerMailsRq) GetOwnerId() int64 {
	if x != nil {
		return x.OwnerId
	}
	return 0
}

type QueryPlayerMailsRs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mails []*common.Mail `protobuf:"bytes,1,rep,name=Mails,proto3" json:"Mails,omitempty"`
}

func (x *QueryPlayerMailsRs) Reset() {
	*x = QueryPlayerMailsRs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryPlayerMailsRs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryPlayerMailsRs) ProtoMessage() {}

func (x *QueryPlayerMailsRs) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryPlayerMailsRs.ProtoReflect.Descriptor instead.
func (*QueryPlayerMailsRs) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{3}
}

func (x *QueryPlayerMailsRs) GetMails() []*common.Mail {
	if x != nil {
		return x.Mails
	}
	return nil
}

type ReadMailRq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerId int64 `protobuf:"varint,1,opt,name=OwnerId,proto3" json:"OwnerId,omitempty"` // 邮箱主人id
	MailId  int64 `protobuf:"varint,2,opt,name=MailId,proto3" json:"MailId,omitempty"`   // 邮件id
}

func (x *ReadMailRq) Reset() {
	*x = ReadMailRq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadMailRq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadMailRq) ProtoMessage() {}

func (x *ReadMailRq) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadMailRq.ProtoReflect.Descriptor instead.
func (*ReadMailRq) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{4}
}

func (x *ReadMailRq) GetOwnerId() int64 {
	if x != nil {
		return x.OwnerId
	}
	return 0
}

func (x *ReadMailRq) GetMailId() int64 {
	if x != nil {
		return x.MailId
	}
	return 0
}

type ReadMailRs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MailId int64             `protobuf:"varint,1,opt,name=MailId,proto3" json:"MailId,omitempty"`                        // 邮件id
	Status common.MailStatus `protobuf:"varint,2,opt,name=Status,proto3,enum=common.MailStatus" json:"Status,omitempty"` // 邮件状态
}

func (x *ReadMailRs) Reset() {
	*x = ReadMailRs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadMailRs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadMailRs) ProtoMessage() {}

func (x *ReadMailRs) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadMailRs.ProtoReflect.Descriptor instead.
func (*ReadMailRs) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{5}
}

func (x *ReadMailRs) GetMailId() int64 {
	if x != nil {
		return x.MailId
	}
	return 0
}

func (x *ReadMailRs) GetStatus() common.MailStatus {
	if x != nil {
		return x.Status
	}
	return common.MailStatus_Unread
}

type GainAttachmentsRq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerId int64 `protobuf:"varint,1,opt,name=OwnerId,proto3" json:"OwnerId,omitempty"` // 邮箱主人id
	MailId  int64 `protobuf:"varint,2,opt,name=MailId,proto3" json:"MailId,omitempty"`   // 邮件id
}

func (x *GainAttachmentsRq) Reset() {
	*x = GainAttachmentsRq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GainAttachmentsRq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GainAttachmentsRq) ProtoMessage() {}

func (x *GainAttachmentsRq) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GainAttachmentsRq.ProtoReflect.Descriptor instead.
func (*GainAttachmentsRq) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{6}
}

func (x *GainAttachmentsRq) GetOwnerId() int64 {
	if x != nil {
		return x.OwnerId
	}
	return 0
}

func (x *GainAttachmentsRq) GetMailId() int64 {
	if x != nil {
		return x.MailId
	}
	return 0
}

type GainAttachmentsRs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MailId int64 `protobuf:"varint,1,opt,name=MailId,proto3" json:"MailId,omitempty"`
}

func (x *GainAttachmentsRs) Reset() {
	*x = GainAttachmentsRs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GainAttachmentsRs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GainAttachmentsRs) ProtoMessage() {}

func (x *GainAttachmentsRs) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GainAttachmentsRs.ProtoReflect.Descriptor instead.
func (*GainAttachmentsRs) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{7}
}

func (x *GainAttachmentsRs) GetMailId() int64 {
	if x != nil {
		return x.MailId
	}
	return 0
}

type DelMailRq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerId int64 `protobuf:"varint,1,opt,name=OwnerId,proto3" json:"OwnerId,omitempty"` // 邮箱主人id
	MailId  int64 `protobuf:"varint,2,opt,name=MailId,proto3" json:"MailId,omitempty"`   // 邮件id
}

func (x *DelMailRq) Reset() {
	*x = DelMailRq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelMailRq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelMailRq) ProtoMessage() {}

func (x *DelMailRq) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelMailRq.ProtoReflect.Descriptor instead.
func (*DelMailRq) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{8}
}

func (x *DelMailRq) GetOwnerId() int64 {
	if x != nil {
		return x.OwnerId
	}
	return 0
}

func (x *DelMailRq) GetMailId() int64 {
	if x != nil {
		return x.MailId
	}
	return 0
}

type DelMailRs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MailId int64 `protobuf:"varint,1,opt,name=MailId,proto3" json:"MailId,omitempty"`
}

func (x *DelMailRs) Reset() {
	*x = DelMailRs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_mail_mail_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelMailRs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelMailRs) ProtoMessage() {}

func (x *DelMailRs) ProtoReflect() protoreflect.Message {
	mi := &file_server_mail_mail_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelMailRs.ProtoReflect.Descriptor instead.
func (*DelMailRs) Descriptor() ([]byte, []int) {
	return file_server_mail_mail_proto_rawDescGZIP(), []int{9}
}

func (x *DelMailRs) GetMailId() int64 {
	if x != nil {
		return x.MailId
	}
	return 0
}

var File_server_mail_mail_proto protoreflect.FileDescriptor

var file_server_mail_mail_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x2f, 0x6d, 0x61,
	0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d, 0x61, 0x69, 0x6c, 0x1a, 0x1a,
	0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x64, 0x65,
	0x66, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf4, 0x01, 0x0a, 0x0c, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x71, 0x12, 0x1e, 0x0a, 0x0a, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x53,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x53,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d,
	0x61, 0x69, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x69,
	0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x32, 0x0a,
	0x0b, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x07, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4c, 0x6f, 0x6f, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x0b, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x22, 0x26, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x69, 0x6c, 0x52,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x22, 0x2e, 0x0a, 0x12, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x71, 0x12,
	0x18, 0x0a, 0x07, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x22, 0x38, 0x0a, 0x12, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x73, 0x12,
	0x22, 0x0a, 0x05, 0x4d, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x05, 0x4d, 0x61,
	0x69, 0x6c, 0x73, 0x22, 0x3e, 0x0a, 0x0a, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x52,
	0x71, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x4d,
	0x61, 0x69, 0x6c, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x4d, 0x61, 0x69,
	0x6c, 0x49, 0x64, 0x22, 0x50, 0x0a, 0x0a, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x52,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x45, 0x0a, 0x11, 0x47, 0x61, 0x69, 0x6e, 0x41, 0x74, 0x74,
	0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x77,
	0x6e, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x4f, 0x77, 0x6e,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x22, 0x2b, 0x0a, 0x11,
	0x47, 0x61, 0x69, 0x6e, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x22, 0x3d, 0x0a, 0x09, 0x44, 0x65, 0x6c,
	0x4d, 0x61, 0x69, 0x6c, 0x52, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x22, 0x23, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x4d,
	0x61, 0x69, 0x6c, 0x52, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x32, 0xb7, 0x02,
	0x0a, 0x0b, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a,
	0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x2e, 0x6d, 0x61,
	0x69, 0x6c, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x71, 0x1a,
	0x12, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x69,
	0x6c, 0x52, 0x73, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x10, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x18, 0x2e, 0x6d, 0x61, 0x69, 0x6c,
	0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x69, 0x6c,
	0x73, 0x52, 0x71, 0x1a, 0x18, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x73, 0x22, 0x00, 0x12,
	0x30, 0x0a, 0x08, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x12, 0x10, 0x2e, 0x6d, 0x61,
	0x69, 0x6c, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x71, 0x1a, 0x10, 0x2e,
	0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x73, 0x22,
	0x00, 0x12, 0x45, 0x0a, 0x0f, 0x47, 0x61, 0x69, 0x6e, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x17, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x47, 0x61, 0x69, 0x6e,
	0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x71, 0x1a, 0x17, 0x2e,
	0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x47, 0x61, 0x69, 0x6e, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x52, 0x73, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x4d,
	0x61, 0x69, 0x6c, 0x12, 0x0f, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x44, 0x65, 0x6c, 0x4d, 0x61,
	0x69, 0x6c, 0x52, 0x71, 0x1a, 0x0f, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x44, 0x65, 0x6c, 0x4d,
	0x61, 0x69, 0x6c, 0x52, 0x73, 0x22, 0x00, 0x42, 0x30, 0x5a, 0x2e, 0x62, 0x69, 0x74, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x66, 0x75, 0x6e, 0x70, 0x6c, 0x75, 0x73,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_server_mail_mail_proto_rawDescOnce sync.Once
	file_server_mail_mail_proto_rawDescData = file_server_mail_mail_proto_rawDesc
)

func file_server_mail_mail_proto_rawDescGZIP() []byte {
	file_server_mail_mail_proto_rawDescOnce.Do(func() {
		file_server_mail_mail_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_mail_mail_proto_rawDescData)
	})
	return file_server_mail_mail_proto_rawDescData
}

var file_server_mail_mail_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_server_mail_mail_proto_goTypes = []interface{}{
	(*CreateMailRq)(nil),       // 0: mail.CreateMailRq
	(*CreateMailRs)(nil),       // 1: mail.CreateMailRs
	(*QueryPlayerMailsRq)(nil), // 2: mail.QueryPlayerMailsRq
	(*QueryPlayerMailsRs)(nil), // 3: mail.QueryPlayerMailsRs
	(*ReadMailRq)(nil),         // 4: mail.ReadMailRq
	(*ReadMailRs)(nil),         // 5: mail.ReadMailRs
	(*GainAttachmentsRq)(nil),  // 6: mail.GainAttachmentsRq
	(*GainAttachmentsRs)(nil),  // 7: mail.GainAttachmentsRs
	(*DelMailRq)(nil),          // 8: mail.DelMailRq
	(*DelMailRs)(nil),          // 9: mail.DelMailRs
	(common.MailType)(0),       // 10: common.MailType
	(*common.LootData)(nil),    // 11: common.LootData
	(*common.Mail)(nil),        // 12: common.Mail
	(common.MailStatus)(0),     // 13: common.MailStatus
}
var file_server_mail_mail_proto_depIdxs = []int32{
	10, // 0: mail.CreateMailRq.Type:type_name -> common.MailType
	11, // 1: mail.CreateMailRq.Attachments:type_name -> common.LootData
	12, // 2: mail.QueryPlayerMailsRs.Mails:type_name -> common.Mail
	13, // 3: mail.ReadMailRs.Status:type_name -> common.MailStatus
	0,  // 4: mail.mailService.CreateMail:input_type -> mail.CreateMailRq
	2,  // 5: mail.mailService.QueryPlayerMails:input_type -> mail.QueryPlayerMailsRq
	4,  // 6: mail.mailService.ReadMail:input_type -> mail.ReadMailRq
	6,  // 7: mail.mailService.GainAttachments:input_type -> mail.GainAttachmentsRq
	8,  // 8: mail.mailService.DelMail:input_type -> mail.DelMailRq
	1,  // 9: mail.mailService.CreateMail:output_type -> mail.CreateMailRs
	3,  // 10: mail.mailService.QueryPlayerMails:output_type -> mail.QueryPlayerMailsRs
	5,  // 11: mail.mailService.ReadMail:output_type -> mail.ReadMailRs
	7,  // 12: mail.mailService.GainAttachments:output_type -> mail.GainAttachmentsRs
	9,  // 13: mail.mailService.DelMail:output_type -> mail.DelMailRs
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_server_mail_mail_proto_init() }
func file_server_mail_mail_proto_init() {
	if File_server_mail_mail_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_mail_mail_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMailRq); i {
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
		file_server_mail_mail_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMailRs); i {
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
		file_server_mail_mail_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryPlayerMailsRq); i {
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
		file_server_mail_mail_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryPlayerMailsRs); i {
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
		file_server_mail_mail_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadMailRq); i {
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
		file_server_mail_mail_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadMailRs); i {
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
		file_server_mail_mail_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GainAttachmentsRq); i {
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
		file_server_mail_mail_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GainAttachmentsRs); i {
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
		file_server_mail_mail_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelMailRq); i {
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
		file_server_mail_mail_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelMailRs); i {
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
			RawDescriptor: file_server_mail_mail_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_server_mail_mail_proto_goTypes,
		DependencyIndexes: file_server_mail_mail_proto_depIdxs,
		MessageInfos:      file_server_mail_mail_proto_msgTypes,
	}.Build()
	File_server_mail_mail_proto = out.File
	file_server_mail_mail_proto_rawDesc = nil
	file_server_mail_mail_proto_goTypes = nil
	file_server_mail_mail_proto_depIdxs = nil
}
