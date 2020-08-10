// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gate/gate.proto

package gate

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	game "github.com/yokaiio/yokai_server/proto/game"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GateStatus struct {
	GateId               int32    `protobuf:"varint,1,opt,name=gate_id,json=gateId,proto3" json:"gate_id,omitempty"`
	Health               int32    `protobuf:"varint,2,opt,name=health,proto3" json:"health,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GateStatus) Reset()         { *m = GateStatus{} }
func (m *GateStatus) String() string { return proto.CompactTextString(m) }
func (*GateStatus) ProtoMessage()    {}
func (*GateStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_d3a53c0d96333ed5, []int{0}
}

func (m *GateStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GateStatus.Unmarshal(m, b)
}
func (m *GateStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GateStatus.Marshal(b, m, deterministic)
}
func (m *GateStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GateStatus.Merge(m, src)
}
func (m *GateStatus) XXX_Size() int {
	return xxx_messageInfo_GateStatus.Size(m)
}
func (m *GateStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_GateStatus.DiscardUnknown(m)
}

var xxx_messageInfo_GateStatus proto.InternalMessageInfo

func (m *GateStatus) GetGateId() int32 {
	if m != nil {
		return m.GateId
	}
	return 0
}

func (m *GateStatus) GetHealth() int32 {
	if m != nil {
		return m.Health
	}
	return 0
}

type UserInfo struct {
	UserId               int64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AccountId            int64    `protobuf:"varint,2,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	GameId               int32    `protobuf:"varint,3,opt,name=game_id,json=gameId,proto3" json:"game_id,omitempty"`
	PlayerId             int64    `protobuf:"varint,4,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	PlayerName           string   `protobuf:"bytes,5,opt,name=player_name,json=playerName,proto3" json:"player_name,omitempty"`
	PlayerLevel          int32    `protobuf:"varint,6,opt,name=player_level,json=playerLevel,proto3" json:"player_level,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserInfo) Reset()         { *m = UserInfo{} }
func (m *UserInfo) String() string { return proto.CompactTextString(m) }
func (*UserInfo) ProtoMessage()    {}
func (*UserInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_d3a53c0d96333ed5, []int{1}
}

func (m *UserInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserInfo.Unmarshal(m, b)
}
func (m *UserInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserInfo.Marshal(b, m, deterministic)
}
func (m *UserInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserInfo.Merge(m, src)
}
func (m *UserInfo) XXX_Size() int {
	return xxx_messageInfo_UserInfo.Size(m)
}
func (m *UserInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_UserInfo.DiscardUnknown(m)
}

var xxx_messageInfo_UserInfo proto.InternalMessageInfo

func (m *UserInfo) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *UserInfo) GetAccountId() int64 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *UserInfo) GetGameId() int32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *UserInfo) GetPlayerId() int64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

func (m *UserInfo) GetPlayerName() string {
	if m != nil {
		return m.PlayerName
	}
	return ""
}

func (m *UserInfo) GetPlayerLevel() int32 {
	if m != nil {
		return m.PlayerLevel
	}
	return 0
}

type GateEmptyMessage struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GateEmptyMessage) Reset()         { *m = GateEmptyMessage{} }
func (m *GateEmptyMessage) String() string { return proto.CompactTextString(m) }
func (*GateEmptyMessage) ProtoMessage()    {}
func (*GateEmptyMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_d3a53c0d96333ed5, []int{2}
}

func (m *GateEmptyMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GateEmptyMessage.Unmarshal(m, b)
}
func (m *GateEmptyMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GateEmptyMessage.Marshal(b, m, deterministic)
}
func (m *GateEmptyMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GateEmptyMessage.Merge(m, src)
}
func (m *GateEmptyMessage) XXX_Size() int {
	return xxx_messageInfo_GateEmptyMessage.Size(m)
}
func (m *GateEmptyMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_GateEmptyMessage.DiscardUnknown(m)
}

var xxx_messageInfo_GateEmptyMessage proto.InternalMessageInfo

type GetGateStatusReply struct {
	Status               *GateStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetGateStatusReply) Reset()         { *m = GetGateStatusReply{} }
func (m *GetGateStatusReply) String() string { return proto.CompactTextString(m) }
func (*GetGateStatusReply) ProtoMessage()    {}
func (*GetGateStatusReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_d3a53c0d96333ed5, []int{3}
}

func (m *GetGateStatusReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetGateStatusReply.Unmarshal(m, b)
}
func (m *GetGateStatusReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetGateStatusReply.Marshal(b, m, deterministic)
}
func (m *GetGateStatusReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetGateStatusReply.Merge(m, src)
}
func (m *GetGateStatusReply) XXX_Size() int {
	return xxx_messageInfo_GetGateStatusReply.Size(m)
}
func (m *GetGateStatusReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetGateStatusReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetGateStatusReply proto.InternalMessageInfo

func (m *GetGateStatusReply) GetStatus() *GateStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

type UpdateUserInfoRequest struct {
	Info                 *UserInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *UpdateUserInfoRequest) Reset()         { *m = UpdateUserInfoRequest{} }
func (m *UpdateUserInfoRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateUserInfoRequest) ProtoMessage()    {}
func (*UpdateUserInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d3a53c0d96333ed5, []int{4}
}

func (m *UpdateUserInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateUserInfoRequest.Unmarshal(m, b)
}
func (m *UpdateUserInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateUserInfoRequest.Marshal(b, m, deterministic)
}
func (m *UpdateUserInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateUserInfoRequest.Merge(m, src)
}
func (m *UpdateUserInfoRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateUserInfoRequest.Size(m)
}
func (m *UpdateUserInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateUserInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateUserInfoRequest proto.InternalMessageInfo

func (m *UpdateUserInfoRequest) GetInfo() *UserInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

type SyncPlayerInfoRequest struct {
	UserId               int64            `protobuf:"varint,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Info                 *game.PlayerInfo `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SyncPlayerInfoRequest) Reset()         { *m = SyncPlayerInfoRequest{} }
func (m *SyncPlayerInfoRequest) String() string { return proto.CompactTextString(m) }
func (*SyncPlayerInfoRequest) ProtoMessage()    {}
func (*SyncPlayerInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d3a53c0d96333ed5, []int{5}
}

func (m *SyncPlayerInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncPlayerInfoRequest.Unmarshal(m, b)
}
func (m *SyncPlayerInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncPlayerInfoRequest.Marshal(b, m, deterministic)
}
func (m *SyncPlayerInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncPlayerInfoRequest.Merge(m, src)
}
func (m *SyncPlayerInfoRequest) XXX_Size() int {
	return xxx_messageInfo_SyncPlayerInfoRequest.Size(m)
}
func (m *SyncPlayerInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncPlayerInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SyncPlayerInfoRequest proto.InternalMessageInfo

func (m *SyncPlayerInfoRequest) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *SyncPlayerInfoRequest) GetInfo() *game.PlayerInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

type SyncPlayerInfoReply struct {
	Info                 *UserInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *SyncPlayerInfoReply) Reset()         { *m = SyncPlayerInfoReply{} }
func (m *SyncPlayerInfoReply) String() string { return proto.CompactTextString(m) }
func (*SyncPlayerInfoReply) ProtoMessage()    {}
func (*SyncPlayerInfoReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_d3a53c0d96333ed5, []int{6}
}

func (m *SyncPlayerInfoReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncPlayerInfoReply.Unmarshal(m, b)
}
func (m *SyncPlayerInfoReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncPlayerInfoReply.Marshal(b, m, deterministic)
}
func (m *SyncPlayerInfoReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncPlayerInfoReply.Merge(m, src)
}
func (m *SyncPlayerInfoReply) XXX_Size() int {
	return xxx_messageInfo_SyncPlayerInfoReply.Size(m)
}
func (m *SyncPlayerInfoReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncPlayerInfoReply.DiscardUnknown(m)
}

var xxx_messageInfo_SyncPlayerInfoReply proto.InternalMessageInfo

func (m *SyncPlayerInfoReply) GetInfo() *UserInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func init() {
	proto.RegisterType((*GateStatus)(nil), "yokai_gate.GateStatus")
	proto.RegisterType((*UserInfo)(nil), "yokai_gate.UserInfo")
	proto.RegisterType((*GateEmptyMessage)(nil), "yokai_gate.GateEmptyMessage")
	proto.RegisterType((*GetGateStatusReply)(nil), "yokai_gate.GetGateStatusReply")
	proto.RegisterType((*UpdateUserInfoRequest)(nil), "yokai_gate.UpdateUserInfoRequest")
	proto.RegisterType((*SyncPlayerInfoRequest)(nil), "yokai_gate.SyncPlayerInfoRequest")
	proto.RegisterType((*SyncPlayerInfoReply)(nil), "yokai_gate.SyncPlayerInfoReply")
}

func init() { proto.RegisterFile("gate/gate.proto", fileDescriptor_d3a53c0d96333ed5) }

var fileDescriptor_d3a53c0d96333ed5 = []byte{
	// 452 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x6d, 0xd2, 0xd6, 0x34, 0x13, 0x28, 0x30, 0xd0, 0x10, 0x85, 0x8f, 0xb6, 0x7b, 0x8a, 0x2a,
	0xe4, 0x48, 0xe5, 0x8c, 0x10, 0x08, 0x54, 0x59, 0xe2, 0x4b, 0x8e, 0xca, 0x01, 0x0e, 0xd1, 0xd6,
	0x9e, 0x3a, 0x16, 0x5e, 0xaf, 0xf1, 0xae, 0x23, 0xf9, 0x37, 0xf0, 0x97, 0xf8, 0x71, 0x68, 0x77,
	0x9d, 0xd6, 0x0e, 0x55, 0xc5, 0x25, 0xf1, 0xbc, 0xf7, 0xe6, 0xcd, 0xcc, 0x93, 0x16, 0xee, 0x27,
	0x5c, 0xd3, 0xcc, 0xfc, 0xf8, 0x45, 0x29, 0xb5, 0x44, 0xa8, 0xe5, 0x4f, 0x9e, 0x2e, 0x0c, 0x32,
	0x79, 0x98, 0x70, 0x41, 0xb3, 0x22, 0xe3, 0x35, 0x95, 0x8e, 0x66, 0xaf, 0x01, 0xce, 0xb8, 0xa6,
	0xb9, 0xe6, 0xba, 0x52, 0xf8, 0x04, 0xee, 0x18, 0xe1, 0x22, 0x8d, 0xc7, 0xbd, 0xa3, 0xde, 0x74,
	0x37, 0xf4, 0x4c, 0x19, 0xc4, 0x38, 0x02, 0x6f, 0x49, 0x3c, 0xd3, 0xcb, 0x71, 0xdf, 0xe1, 0xae,
	0x62, 0x7f, 0x7a, 0xb0, 0x77, 0xae, 0xa8, 0x0c, 0xf2, 0x4b, 0x69, 0xba, 0x2b, 0x45, 0xe5, 0xba,
	0x7b, 0x3b, 0xf4, 0x4c, 0x19, 0xc4, 0xf8, 0x1c, 0x80, 0x47, 0x91, 0xac, 0x72, 0x6d, 0xb8, 0xbe,
	0xe5, 0x06, 0x0d, 0x12, 0xc4, 0x6e, 0xaa, 0xb0, 0x53, 0xb7, 0xd7, 0x53, 0x85, 0x99, 0xfa, 0x14,
	0x06, 0x6e, 0x59, 0x43, 0xed, 0xd8, 0xb6, 0x3d, 0x07, 0x04, 0x31, 0x1e, 0xc2, 0xb0, 0x21, 0x73,
	0x2e, 0x68, 0xbc, 0x7b, 0xd4, 0x9b, 0x0e, 0x42, 0x70, 0xd0, 0x67, 0x2e, 0x08, 0x8f, 0xe1, 0x6e,
	0x23, 0xc8, 0x68, 0x45, 0xd9, 0xd8, 0xb3, 0xde, 0x4d, 0xd3, 0x47, 0x03, 0x31, 0x84, 0x07, 0xe6,
	0xfa, 0x0f, 0xa2, 0xd0, 0xf5, 0x27, 0x52, 0x8a, 0x27, 0xc4, 0xde, 0x03, 0x9e, 0x91, 0xbe, 0x0e,
	0x25, 0xa4, 0x22, 0xab, 0xd1, 0x07, 0x4f, 0xd9, 0xd2, 0x9e, 0x36, 0x3c, 0x1d, 0xf9, 0xd7, 0xb9,
	0xfa, 0x2d, 0x71, 0xa3, 0x62, 0x6f, 0xe1, 0xe0, 0xbc, 0x88, 0xb9, 0xa6, 0x75, 0x3a, 0x21, 0xfd,
	0xaa, 0x48, 0x69, 0x9c, 0xc2, 0x4e, 0x9a, 0x5f, 0xca, 0xc6, 0xe6, 0x71, 0xdb, 0xe6, 0x4a, 0x6a,
	0x15, 0xec, 0x07, 0x1c, 0xcc, 0xeb, 0x3c, 0xfa, 0xea, 0x0e, 0x6e, 0x59, 0x8c, 0xc0, 0xb3, 0xd2,
	0xab, 0x98, 0x5d, 0x85, 0x27, 0x8d, 0x75, 0x7f, 0x63, 0x43, 0x41, 0x7e, 0xcb, 0xc4, 0x99, 0xbf,
	0x81, 0x47, 0x9b, 0xe6, 0xe6, 0xcc, 0xff, 0xde, 0xee, 0xf4, 0x77, 0x1f, 0x86, 0x06, 0x9f, 0x53,
	0xb9, 0x4a, 0x23, 0xc2, 0x2f, 0x70, 0xaf, 0x13, 0x1b, 0x3e, 0xdb, 0x4c, 0xa8, 0x9d, 0xf2, 0xe4,
	0x45, 0x87, 0xfd, 0x27, 0x6f, 0xb6, 0x85, 0x73, 0xd8, 0xef, 0x26, 0x88, 0xc7, 0x9d, 0x75, 0x6e,
	0x4a, 0x77, 0x72, 0xeb, 0x50, 0xb6, 0x85, 0xdf, 0x60, 0xbf, 0x7b, 0x76, 0xd7, 0xf4, 0xc6, 0xbc,
	0x27, 0x87, 0xb7, 0x49, 0xec, 0xb2, 0xef, 0x5e, 0x7e, 0x3f, 0x49, 0x52, 0xbd, 0xac, 0x2e, 0xfc,
	0x48, 0x8a, 0x99, 0x95, 0xa7, 0xd2, 0xfd, 0x2f, 0x14, 0x95, 0x2b, 0x2a, 0x67, 0xf6, 0xbd, 0xd9,
	0x97, 0x79, 0xe1, 0xd9, 0xef, 0x57, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x15, 0xe9, 0xca, 0x7b,
	0xad, 0x03, 0x00, 0x00,
}
