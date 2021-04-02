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
	mailQueryInterval = time.Minute * 30 // 每30分钟拉取一次最新的邮件数据
)

type MailManager struct {
	nextUpdate int64                  `bson:"-" json:"-"` // 下次更新时间
	owner      *Player                `bson:"-" json:"-"`
	Mails      map[int64]*define.Mail `bson:"mail_list" json:"mail_list"` // 邮件缓存
}

func NewMailManager(owner *Player) *MailManager {
	m := &MailManager{
		nextUpdate: time.Now().Add(time.Second * time.Duration(rand.Int31n(5))).Unix(),
		owner:      owner,
		Mails:      make(map[int64]*define.Mail),
	}

	return m
}

func (m *MailManager) update() {
	if m.nextUpdate < time.Now().Unix() {
		return
	}

	// 请求邮件列表
	rsp, err := m.owner.acct.rpcCaller.CallQueryPlayerMails(&pbMail.QueryPlayerMailsRq{
		OwnerId: m.owner.ID,
	})

	// 请求失败5秒后再试
	if !utils.ErrCheck(err, "CallQueryPlayerMails failed when MailManager.update", m.owner.ID) {
		m.nextUpdate = time.Now().Add(time.Second * 5).Unix()
		return
	}

	m.Mails = make(map[int64]*define.Mail)
	for _, pb := range rsp.GetMails() {
		newMail := &define.Mail{}
		newMail.FromPB(pb)
		m.Mails[newMail.Id] = newMail
	}

	// 请求成功半小时后再同步
	m.nextUpdate = time.Now().Add(mailQueryInterval).Unix()
	log.Info().Int64("player_id", m.owner.ID).Msg("rpc query mail list success")
}

func (m *MailManager) GetMail(mailId int64) (*define.Mail, bool) {
	mail, ok := m.Mails[mailId]
	return mail, ok
}
