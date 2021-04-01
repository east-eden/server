package mail

import (
	"errors"
	"strconv"
	"sync"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"github.com/valyala/bytebufferpool"
)

var (
	ErrInvalidMail       = errors.New("invalid mail")
	ErrInvalidMailStatus = errors.New("invalid mail status")
	ErrAddExistMail      = errors.New("add exist mail")
)

func makeMailKey(mailId int64, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	_, _ = b.WriteString("mail_list.")
	_, _ = b.WriteString(strconv.Itoa(int(mailId)))

	for _, f := range fields {
		_, _ = b.WriteString(".")
		_, _ = b.WriteString(f)
	}

	return b.String()
}

type MailOwnerInfo struct {
	Id             int64 `json:"_id" bson:"_id"`                             // 邮箱主人id
	LastSaveNodeId int32 `json:"last_save_node_id" bson:"last_save_node_id"` // 最后一次存储时所在节点的id
}

// 邮件箱
type MailBox struct {
	sync.RWMutex  `json:"-" bson:"-"`
	MailOwnerInfo `json:",inline" bson:"inline"` // 邮件主人信息
	Mails         map[int64]*define.Mail         `json:"mail_list" bson:"mail_list"` // 邮件
}

func NewMailBox() interface{} {
	return &MailBox{}
}

func (b *MailBox) Init() {
	b.Id = -1
	b.LastSaveNodeId = -1
	b.Mails = make(map[int64]*define.Mail)
}

func (b *MailBox) GetMail(mailId int64) (*define.Mail, bool) {
	b.RLock()
	defer b.RUnlock()

	mail, ok := b.Mails[mailId]
	return mail, ok
}

func (b *MailBox) ReadMail(mailId int64) error {
	b.Lock()
	defer b.Unlock()

	mail, ok := b.Mails[mailId]
	if !ok {
		return ErrInvalidMail
	}

	if mail.Status == define.Mail_Status_Readed {
		return ErrInvalidMailStatus
	}

	mail.Status = define.Mail_Status_Readed
	fields := map[string]interface{}{
		makeMailKey(mail.Id, "status"): mail.Status,
	}
	err := store.GetStore().SaveObjectFields(define.StoreType_Mail, b.Id, nil, fields)
	utils.ErrPrint(err, "SaveObjectFields failed when MailBox.ReadMail", b.Id, mail.Id)
	return err
}

func (b *MailBox) GainAttachments(mailId int64) error {
	b.Lock()
	defer b.Unlock()

	mail, ok := b.Mails[mailId]
	if !ok {
		return ErrInvalidMail
	}

	// 已领取过附件
	if mail.Status == define.Mail_Status_GainedAttachments {
		return ErrInvalidMailStatus
	}

	mail.Status = define.Mail_Status_GainedAttachments
	fields := map[string]interface{}{
		makeMailKey(mail.Id, "status"): mail.Status,
	}
	err := store.GetStore().SaveObjectFields(define.StoreType_Mail, b.Id, nil, fields)
	utils.ErrPrint(err, "SaveObjectFields failed when MailBox.GainAttachments", b.Id, mail.Id)
	return err
}

func (b *MailBox) AddMail(mail *define.Mail) error {
	b.Lock()
	defer b.Unlock()

	_, ok := b.Mails[mail.Id]
	if ok {
		return ErrAddExistMail
	}

	b.Mails[mail.Id] = mail
	fields := map[string]interface{}{
		makeMailKey(mail.Id): mail,
	}
	err := store.GetStore().SaveObjectFields(define.StoreType_Mail, b.Id, nil, fields)
	utils.ErrPrint(err, "SaveobjectFields failed when MailBox.AddMail", b.Id, mail.Id)
	return err
}

func (b *MailBox) DeleteMail(mailId int64) error {
	b.Lock()
	defer b.Unlock()

	_, ok := b.Mails[mailId]
	if !ok {
		return ErrInvalidMail
	}

	fields := []string{
		makeMailKey(mailId),
	}
	err := store.GetStore().DeleteObjectFields(define.StoreType_Mail, b.Id, nil, fields)
	utils.ErrPrint(err, "DeleteObjectFields failed when MailBox.DeleteMail", b.Id, mailId)
	return err
}
