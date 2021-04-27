package define

import (
	"time"

	pbCommon "bitbucket.org/funplus/server/proto/global/common"
)

// 邮件状态
const (
	Mail_Status_Unread            int32 = 0 // 未读
	Mail_Status_Readed            int32 = 1 // 已读
	Mail_Status_GainedAttachments int32 = 2 // 已获取附件
)

// 邮件类型
const (
	Mail_Type_System int32 = 0 // 系统邮件
	Mail_Type_Player int32 = 1 // 玩家寄送
)

// 邮件上下文
type MailContext struct {
	Id         int64  `bson:"_id" json:"_id"`                 // 邮件id
	OwnerId    int64  `bson:"owner_id" json:"owner_id"`       // 拥有者id
	SenderId   int64  `bson:"sender_id" json:"sender_id"`     // 发件人id
	Status     int32  `bson:"status" json:"status"`           // 邮件状态
	Type       int32  `bson:"type" json:"type"`               // 邮件类型
	Date       int32  `bson:"date" json:"date"`               // 寄件日期
	ExpireDate int32  `bson:"expire_date" json:"expire_date"` // 邮件过期日期
	SenderName string `bson:"sender_name" json:"from"`        // 寄件人名字
	Title      string `bson:"title" json:"title"`             // 邮件标题
	Content    string `bson:"context" json:"content"`         // 邮件内容
}

// 邮件
type Mail struct {
	MailContext `bson:"inline" json:",inline"` // 邮件上下文
	Attachments []*LootData                    `bson:"attachments" json:"attachments"` // 附件
}

func (m *Mail) Init() {
	m.Id = -1
	m.OwnerId = -1
	m.SenderId = -1
	m.Status = Mail_Status_Unread
	m.Type = Mail_Type_System
	m.Date = int32(time.Now().Unix())
	m.ExpireDate = -1
}

func (m *Mail) CanRead() bool {
	return m.Status == Mail_Status_Unread
}

func (m *Mail) CanGainAttachments() bool {
	return m.Status != Mail_Status_GainedAttachments && len(m.Attachments) > 0
}

func (m *Mail) CanDel() bool {
	if m.CanRead() {
		return false
	}

	if m.CanGainAttachments() {
		return false
	}

	return true
}

func (m *Mail) IsExpired() bool {
	if m.ExpireDate == -1 {
		return false
	}

	return m.ExpireDate < int32(time.Now().Unix())
}

func (m *Mail) ToPB() *pbCommon.Mail {
	pb := &pbCommon.Mail{
		Context: &pbCommon.MailContext{
			Id:         m.Id,
			SenderId:   m.SenderId,
			Status:     pbCommon.MailStatus(m.Status),
			Type:       pbCommon.MailType(m.Type),
			Date:       m.Date,
			ExpireDate: m.ExpireDate,
			SenderName: m.SenderName,
			Title:      m.Title,
			Content:    m.Content,
		},
		Attachments: make([]*pbCommon.LootData, 0, len(m.Attachments)),
	}

	for _, attachment := range m.Attachments {
		pb.Attachments = append(pb.Attachments, &pbCommon.LootData{
			Type: pbCommon.LootType(attachment.LootType),
			Misc: attachment.LootMisc,
			Num:  attachment.LootNum,
		})
	}

	return pb
}

func (m *Mail) FromPB(pb *pbCommon.Mail) {
	m.Id = pb.Context.GetId()
	m.SenderId = pb.Context.GetSenderId()
	m.Status = int32(pb.GetContext().Status)
	m.Type = int32(pb.Context.GetType())
	m.Date = pb.Context.GetDate()
	m.ExpireDate = pb.Context.GetExpireDate()
	m.SenderName = pb.Context.GetSenderName()
	m.Title = pb.Context.GetTitle()
	m.Content = pb.Context.GetContent()
	m.Attachments = make([]*LootData, 0, len(pb.GetAttachments()))

	for _, attachment := range pb.GetAttachments() {
		newAttachment := &LootData{
			LootType: int32(attachment.GetType()),
			LootMisc: attachment.GetMisc(),
			LootNum:  attachment.GetNum(),
		}

		m.Attachments = append(m.Attachments, newAttachment)
	}
}
