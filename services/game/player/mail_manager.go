package player

import (
	"math/rand"
	"time"

	"bitbucket.org/funplus/server/define"
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

	m.nextUpdate = time.Now().Add(mailQueryInterval).Unix()

	// todo rpc query mails
}

func (m *MailManager) GetMail(mailId int64) (*define.Mail, bool) {
	mail, ok := m.Mails[mailId]
	return mail, ok
}
