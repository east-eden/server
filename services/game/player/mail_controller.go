package player

import (
	"math/rand"
	"time"

	"bitbucket.org/funplus/server/define"
	pbMail "bitbucket.org/funplus/server/proto/server/mail"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	mailUpdateInterval = time.Second * 5  // 每5秒更新一次
	mailQueryInterval  = time.Minute * 30 // 每30分钟拉取一次最新的邮件数据
)

type MailController struct {
	nextUpdate int64                  `bson:"-" json:"-"` // 下次更新时间
	nextQuery  int64                  `bson:"-" json:"-"` // 下次请求邮件列表时间
	owner      *Player                `bson:"-" json:"-"`
	Mails      map[int64]*define.Mail `bson:"mail_list" json:"mail_list"` // 邮件缓存
}

func NewMailManager(owner *Player) *MailController {
	m := &MailController{
		nextUpdate: time.Now().Add(time.Second * time.Duration(rand.Int31n(5))).Unix(),
		nextQuery:  time.Now().Add(time.Second * time.Duration(rand.Int31n(5))).Unix(),
		owner:      owner,
		Mails:      make(map[int64]*define.Mail),
	}

	return m
}

func (m *MailController) update() {
	if m.nextUpdate > time.Now().Unix() {
		return
	}

	m.nextUpdate = time.Now().Add(mailUpdateInterval).Unix()

	// 请求所有邮件列表
	m.updateQueryMails()

	// 更新过期邮件
	m.updateExpiredMails()
}

func (m *MailController) updateQueryMails() {
	if m.nextQuery > time.Now().Unix() {
		return
	}

	// 请求邮件列表
	req := &pbMail.QueryPlayerMailsRq{
		OwnerId: m.owner.ID,
	}
	rsp, err := m.owner.acct.rpcCaller.CallQueryPlayerMails(req)
	if !utils.ErrCheck(err, "CallQueryPlayerMails failed when MailManager.updateQueryMails", req) {
		return
	}

	m.Mails = make(map[int64]*define.Mail)
	for _, pb := range rsp.GetMails() {
		newMail := &define.Mail{}
		newMail.FromPB(pb)
		m.Mails[newMail.Id] = newMail
	}

	m.nextQuery = time.Now().Add(mailQueryInterval).Unix()
	log.Info().Int64("player_id", m.owner.ID).Interface("response", rsp).Msg("rpc query mail list success")
}

func (m *MailController) updateExpiredMails() {
	for _, mail := range m.Mails {
		if !mail.IsExpired() {
			continue
		}

		req := &pbMail.DelMailRq{
			OwnerId: m.owner.ID,
			MailId:  mail.Id,
		}
		_, err := m.owner.acct.rpcCaller.CallDelMail(req)
		if utils.ErrCheck(err, "CallDelMail failed when MailManager.updateExpiredMails", req) {
			delete(m.Mails, mail.Id)
		}
	}
}

////////////////////////////////////////////////////
// user interface
func (m *MailController) GetMail(mailId int64) (*define.Mail, bool) {
	mail, ok := m.Mails[mailId]
	return mail, ok
}

func (m *MailController) ReadAllMail() error {
	for _, mail := range m.Mails {
		if !mail.CanRead() {
			continue
		}

		req := &pbMail.ReadMailRq{
			OwnerId: m.owner.ID,
			MailId:  mail.Id,
		}
		_, err := m.owner.acct.rpcCaller.CallReadMail(req)
		if utils.ErrCheck(err, "CallReadMail failed when MailManager.ReadAllMail", req) {
			mail.Status = define.Mail_Status_Readed
		}
	}

	return nil
}

func (m *MailController) GainAllMailsAttachments() error {
	for _, mail := range m.Mails {
		if !mail.CanGainAttachments() {
			continue
		}

		req := &pbMail.GainAttachmentsRq{
			OwnerId: m.owner.ID,
			MailId:  mail.Id,
		}
		_, err := m.owner.acct.rpcCaller.CallGainAttachments(req)
		if utils.ErrCheck(err, "CallGainAttachments failed when MailManager.GainAllMailsAttachments", req) {
			mail.Status = define.Mail_Status_GainedAttachments
			_ = m.owner.CostLootManager().GainLootByList(mail.Attachments)
		}
	}

	return nil
}

func (m *MailController) DelAllMails() error {
	for _, mail := range m.Mails {
		if !mail.CanDel() {
			continue
		}

		req := &pbMail.DelMailRq{
			OwnerId: m.owner.ID,
			MailId:  mail.Id,
		}
		_, err := m.owner.acct.rpcCaller.CallDelMail(req)
		if utils.ErrCheck(err, "CallDelMail failed when MailManager.DelAllMails", req) {
			delete(m.Mails, mail.Id)
		}
	}

	return nil
}
