package mailbox

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	ErrInvalidMail       = errors.New("invalid mail")
	ErrInvalidMailStatus = errors.New("invalid mail status")
	ErrAddExistMail      = errors.New("add exist mail")

	MailBoxHandlerNum        = 100             // 邮箱channel处理缓存
	MailBoxResultHandlerNum  = 100             // 邮箱带返回channel处理缓存
	MailChannelResultTimeout = 5 * time.Second // 邮箱channel处理超时
)

type MailBoxHandler func(context.Context, *MailBox) error
type MailBoxResultHandler struct {
	F MailBoxHandler
	E chan<- error
	C context.Context
}

type MailOwnerInfo struct {
	Id             int64 `json:"_id" bson:"_id"`                             // 邮箱主人id
	LastSaveNodeId int32 `json:"last_save_node_id" bson:"last_save_node_id"` // 最后一次存储时所在节点的id
}

// 邮件箱
type MailBox struct {
	MailOwnerInfo `json:",inline" bson:"inline"` // 邮件主人信息
	NodeId        int16                          `json:"-" bson:"-"` // 当前节点id
	Mails         map[int64]*define.Mail         `json:"-" bson:"-"` // 邮件
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

func (b *MailBox) Load(ownerId int64) error {
	// 加载邮箱信息
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Mail, ownerId, b)

	// 创建新邮箱数据
	if errors.Is(err, store.ErrNoResult) {
		b.Id = ownerId
		b.LastSaveNodeId = int32(b.NodeId)
		errSave := store.GetStore().UpdateOne(context.Background(), define.StoreType_Mail, ownerId, b)
		utils.ErrPrint(errSave, "UpdateOne failed when MailBox.Load", ownerId)
		return errSave
	}

	if !utils.ErrCheck(err, "FindOne failed when MailBox.Load", ownerId) {
		return err
	}

	// 加载所有邮件
	res, errMails := store.GetStore().FindAll(context.Background(), define.StoreType_Mail, "owner_id", ownerId)
	if !utils.ErrCheck(errMails, "FindAll failed when MailBox.Load", ownerId) {
		return errMails
	}

	for _, v := range res {
		vv := v.([]byte)
		mail := &define.Mail{}
		err := json.Unmarshal(vv, mail)
		if !utils.ErrCheck(err, "json.Unmarshal failed when MailBox.Load", ownerId) {
			continue
		}

		b.Mails[mail.Id] = mail
	}

	return nil
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
				if err := b.checkAvaliable(ctx); err != nil {
					return err
				}

				err := h(ctx, b)
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
				if err := b.checkAvaliable(rh.C); err != nil {
					return err
				}

				err := rh.F(rh.C, b)
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

func (b *MailBox) AddResultHandler(ctx context.Context, h MailBoxHandler) error {
	subCtx, cancel := utils.WithTimeoutContext(ctx, MailChannelResultTimeout)
	defer cancel()

	e := make(chan error, 1)
	b.ResultHandles <- &MailBoxResultHandler{
		F: h,
		E: e,
		C: subCtx,
	}

	select {
	case err := <-e:
		return err
	case <-subCtx.Done():
		return subCtx.Err()
	}
}

func (b *MailBox) checkAvaliable(ctx context.Context) error {
	// 读取最后存储时节点id
	var ownerInfo MailOwnerInfo
	err := store.GetStore().FindOne(ctx, define.StoreType_Mail, b.Id, &ownerInfo)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if !utils.ErrCheck(err, "LoadObject failed when MailBox.checkAvaliable", b.Id) {
		return err
	}

	// 上次存储不在当前节点
	if int16(ownerInfo.LastSaveNodeId) != b.NodeId {
		err := store.GetStore().FindOne(ctx, define.StoreType_Mail, b.Id, b)
		if !utils.ErrCheck(err, "LoadObject failed when MailBox.checkAvaliable", b.Id) {
			return err
		}
	}

	return nil
}

func (b *MailBox) ReadMail(ctx context.Context, mailId int64) error {
	mail, ok := b.Mails[mailId]
	if !ok {
		return ErrInvalidMail
	}

	if !mail.CanRead() {
		return ErrInvalidMailStatus
	}

	mail.Status = define.Mail_Status_Readed
	fields := map[string]interface{}{
		"status": define.Mail_Status_Readed,
	}
	err := store.GetStore().UpdateFields(ctx, define.StoreType_Mail, mail.Id, fields)
	utils.ErrPrint(err, "UpdateFields failed when MailBox.ReadMail", b.Id, mail.Id)

	return err
}

func (b *MailBox) GainAttachments(ctx context.Context, mailId int64) error {
	mail, ok := b.Mails[mailId]
	if !ok {
		return ErrInvalidMail
	}

	// 已领取过附件
	if !mail.CanGainAttachments() {
		return ErrInvalidMailStatus
	}

	mail.Status = define.Mail_Status_GainedAttachments
	fields := map[string]interface{}{
		"status": define.Mail_Status_GainedAttachments,
	}
	err := store.GetStore().UpdateFields(ctx, define.StoreType_Mail, mail.Id, fields)
	utils.ErrPrint(err, "UpdateFields failed when MailBox.GainAttachments", b.Id, mail.Id)

	return err
}

func (b *MailBox) AddMail(ctx context.Context, mail *define.Mail) error {
	_, ok := b.Mails[mail.Id]
	if ok {
		return ErrAddExistMail
	}

	b.Mails[mail.Id] = mail

	err := store.GetStore().UpdateOne(ctx, define.StoreType_Mail, mail.Id, mail)
	utils.ErrPrint(err, "UpdateOne failed when MailBox.AddMail", b.Id, mail.Id)

	return err
}

func (b *MailBox) DelMail(ctx context.Context, mailId int64) error {
	_, ok := b.Mails[mailId]
	if !ok {
		return ErrInvalidMail
	}

	delete(b.Mails, mailId)
	err := store.GetStore().DeleteOne(ctx, define.StoreType_Mail, mailId)
	utils.ErrPrint(err, "DeleteObjectFields failed when MailBox.DeleteMail", b.Id, mailId)

	return err
}

func (b *MailBox) GetMails(ctx context.Context) []*define.Mail {
	r := make([]*define.Mail, 0, len(b.Mails))
	for _, mail := range b.Mails {
		r = append(r, mail)
	}

	return r
}

// test interface
func (b *MailBox) BenchAddMail(ctx context.Context, mail *define.Mail) error {
	_, ok := b.Mails[mail.Id]
	if ok {
		return ErrAddExistMail
	}

	b.Mails[mail.Id] = mail

	// fields := map[string]interface{}{
	// 	makeMailKey(mail.Id): mail,
	// }
	// err := store.GetStore().UpdateFields(ctx, define.StoreType_Mail, b.Id, fields)

	// err := store.GetStore().UpdateOne(ctx, define.StoreType_Mail, mail.Id, mail)

	err := store.GetStore().PushArray(ctx, define.StoreType_Mail, b.Id, "mail_list", mail)

	utils.ErrPrint(err, "UpdateOne failed when MailBox.AddMail", b.Id, mail.Id)

	return err
}
