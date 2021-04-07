package mailbox

import (
	"context"
	"errors"
	"strconv"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
	"github.com/valyala/bytebufferpool"
)

var (
	ErrInvalidMail       = errors.New("invalid mail")
	ErrInvalidMailStatus = errors.New("invalid mail status")
	ErrAddExistMail      = errors.New("add exist mail")

	MailBoxHandlerNum       = 100             // 邮箱channel处理缓存
	MailBoxResultHandlerNum = 100             // 邮箱带返回channel处理缓存
	channelHandleTimeout    = 5 * time.Second // channel处理超时
)

type MailBoxHandler func(*MailBox) error
type MailBoxResultHandler struct {
	F MailBoxHandler
	E chan<- error
}

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
	MailOwnerInfo `json:",inline" bson:"inline"` // 邮件主人信息
	NodeId        int16                          `json:"-" bson:"-"`                 // 当前节点id
	Mails         map[int64]*define.Mail         `json:"mail_list" bson:"mail_list"` // 邮件
	Handles       chan MailBoxHandler            `json:"-" bson:"-"`
	ResultHandles chan *MailBoxResultHandler     `json:"-" bson:"-"`
}

func NewMailBox() interface{} {
	return &MailBox{}
}

func (b *MailBox) Init(nodeId int16) {
	b.Id = -1
	b.LastSaveNodeId = -1
	b.NodeId = nodeId
	b.Mails = make(map[int64]*define.Mail)
	b.Handles = make(chan MailBoxHandler, MailBoxHandlerNum)
	b.ResultHandles = make(chan *MailBoxResultHandler, MailBoxHandlerNum)
}

func (b *MailBox) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Info().Int64("mailbox_id", b.Id).Msg("mail box context done...")
			return nil

		case h, ok := <-b.Handles:
			if !ok {
				log.Info().Int64("mailbox_id", b.Id).Msg("mail box handler channel closed")
				return nil
			} else {
				// 每次回调前先检查是否需要重新加载db
				if err := b.checkAvaliable(); err != nil {
					return err
				}

				err := h(b)
				if !utils.ErrCheck(err, "mailbox handler failed", b.Id) {
					return err
				}
			}

		case rh, ok := <-b.ResultHandles:
			if !ok {
				log.Info().Int64("mailbox_id", b.Id).Msg("mail box result handler channel closed")
				return nil
			} else {
				// 每次回调前先检查是否需要重新加载db
				if err := b.checkAvaliable(); err != nil {
					return err
				}

				err := rh.F(b)
				rh.E <- err
				if !utils.ErrCheck(err, "mailbox handler failed", b.Id) {
					return err
				}
			}
		}
	}
}

func (b *MailBox) AddHandler(h MailBoxHandler) {
	b.Handles <- h
}

func (b *MailBox) AddResultHandler(h MailBoxHandler) error {
	timeout, cancel := context.WithTimeout(context.Background(), channelHandleTimeout)
	defer cancel()

	e := make(chan error, 1)
	b.ResultHandles <- &MailBoxResultHandler{
		F: h,
		E: e,
	}

	select {
	case err := <-e:
		return err
	case <-timeout.Done():
		return timeout.Err()
	}
}

func (b *MailBox) checkAvaliable() error {
	// 读取最后存储时节点id
	var ownerInfo MailOwnerInfo
	err := store.GetStore().FindOne(define.StoreType_Mail, b.Id, &ownerInfo)
	if !utils.ErrCheck(err, "LoadObject failed when MailBox.checkAvaliable", b.Id) {
		return err
	}

	// 上次存储不在当前节点
	if int16(ownerInfo.LastSaveNodeId) != b.NodeId {
		err := store.GetStore().FindOne(define.StoreType_Mail, b.Id, b)
		if !utils.ErrCheck(err, "LoadObject failed when MailBox.checkAvaliable", b.Id) {
			return err
		}
	}

	return nil
}

func (b *MailBox) ReadMail(mailId int64) error {
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
	err := store.GetStore().UpdateFields(define.StoreType_Mail, b.Id, fields)
	utils.ErrPrint(err, "SaveObjectFields failed when MailBox.ReadMail", b.Id, mail.Id)
	return err
}

func (b *MailBox) GainAttachments(mailId int64) error {
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
	err := store.GetStore().UpdateFields(define.StoreType_Mail, b.Id, fields)
	utils.ErrPrint(err, "SaveObjectFields failed when MailBox.GainAttachments", b.Id, mail.Id)
	return err
}

func (b *MailBox) AddMail(mail *define.Mail) error {
	_, ok := b.Mails[mail.Id]
	if ok {
		return ErrAddExistMail
	}

	b.Mails[mail.Id] = mail
	fields := map[string]interface{}{
		makeMailKey(mail.Id): mail,
	}
	err := store.GetStore().UpdateFields(define.StoreType_Mail, b.Id, fields)
	utils.ErrPrint(err, "SaveobjectFields failed when MailBox.AddMail", b.Id, mail.Id)
	return err
}

func (b *MailBox) DelMail(mailId int64) error {
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

func (b *MailBox) GetMails() []*define.Mail {
	r := make([]*define.Mail, 0, len(b.Mails))
	for _, mail := range b.Mails {
		r = append(r, mail)
	}

	return r
}
