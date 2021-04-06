package mail

import (
	"context"
	"errors"
	"sync"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/services/mail/mailbox"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"bitbucket.org/funplus/server/utils/cache"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	mailBoxCacheExpire   = 10 * time.Minute // 邮箱cache缓存10分钟
	channelHandleTimeout = 5 * time.Second  // channel处理超时
	ErrInvalidOwner      = errors.New("invalid owner")
	ErrTimeout           = errors.New("mail manager handle time out")
)

type MailManager struct {
	m              *Mail
	cacheMailBoxes *cache.Cache
	mailBoxPool    sync.Pool
	wg             utils.WaitGroupWrapper
}

func NewMailManager(ctx *cli.Context, m *Mail) *MailManager {
	manager := &MailManager{
		m:              m,
		cacheMailBoxes: cache.New(mailBoxCacheExpire, mailBoxCacheExpire),
	}

	// 邮箱池
	manager.mailBoxPool.New = mailbox.NewMailBox

	// 邮箱缓存删除时处理
	manager.cacheMailBoxes.OnEvicted(func(k, v interface{}) {
		manager.mailBoxPool.Put(v)
	})

	// 初始化db
	store.GetStore().AddStoreInfo(define.StoreType_Mail, "mail", "_id")
	if err := store.GetStore().MigrateDbTable("mail", "player_id"); err != nil {
		log.Fatal().Err(err).Msg("migrate collection mail failed")
	}

	log.Info().Msg("MailManager init ok ...")
	return manager
}

func (m *MailManager) Run(ctx *cli.Context) error {
	<-ctx.Done()
	m.wg.Wait()
	log.Info().Msg("MailManager context done...")
	return nil
}

// 获取邮箱数据
func (m *MailManager) getMailBox(ownerId int64) (*mailbox.MailBox, error) {
	if ownerId == -1 {
		return nil, ErrInvalidOwner
	}

	mb, ok := m.cacheMailBoxes.Get(ownerId)

	// 缓存没有，从db加载
	if !ok {
		mb = m.mailBoxPool.Get()
		mailbox := mb.(*mailbox.MailBox)
		mailbox.Init(m.m.ID)
		err := store.GetStore().LoadObject(define.StoreType_Mail, ownerId, mailbox)

		// 创建新邮箱数据
		if errors.Is(err, store.ErrNoResult) {
			mailbox.Id = ownerId
			mailbox.LastSaveNodeId = int32(m.m.ID)
			errSave := store.GetStore().SaveObject(define.StoreType_Mail, ownerId, mailbox)
			utils.ErrPrint(errSave, "SaveObject failed when MailManager.getMailBox", ownerId)
		} else {
			if !utils.ErrCheck(err, "LoadObject failed when MailManager.getMailBox", ownerId) {
				m.mailBoxPool.Put(mb)
				return nil, err
			}
		}

		m.cacheMailBoxes.Set(ownerId, mb, mailBoxCacheExpire)
	}

	m.wg.Wrap(func() {
		defer utils.CaptureException()

		err := mb.(*mailbox.MailBox).Run(context.Background())
		utils.ErrPrint(err, "mailbox run failed", mb.(*mailbox.MailBox).Id)

		// 删除缓存
		m.cacheMailBoxes.Delete(mb.(*mailbox.MailBox).Id)
	})

	return mb.(*mailbox.MailBox), nil
}

// 创建新邮件
func (m *MailManager) CreateMail(receiverId int64, mail *define.Mail) error {
	mb, err := m.getMailBox(receiverId)
	if err != nil {
		return err
	}

	result := make(chan error, 1)
	timeout, cancel := context.WithTimeout(context.Background(), channelHandleTimeout)
	defer cancel()

	mb.AddResultHandler(func(mailBox *mailbox.MailBox) error {
		return mailBox.AddMail(mail)
	}, result)

	select {
	case err := <-result:
		return err
	case <-timeout.Done():
		return ErrTimeout
	}
}

// 删除邮件
func (m *MailManager) DelMail(receiverId int64, mailId int64) error {
	mb, err := m.getMailBox(receiverId)
	if err != nil {
		return err
	}

	result := make(chan error, 1)
	timeout, cancel := context.WithTimeout(context.Background(), channelHandleTimeout)
	defer cancel()

	mb.AddResultHandler(func(mailBox *mailbox.MailBox) error {
		return mailBox.DelMail(mailId)
	}, result)

	select {
	case err := <-result:
		return err
	case <-timeout.Done():
		return ErrTimeout
	}
}

// 查询玩家邮件
func (m *MailManager) QueryPlayerMails(ownerId int64) ([]*define.Mail, error) {
	retMails := make([]*define.Mail, 0)
	mb, err := m.getMailBox(ownerId)
	if err != nil {
		return retMails, err
	}

	result := make(chan error, 1)
	timeout, cancel := context.WithTimeout(context.Background(), channelHandleTimeout)
	defer cancel()

	mb.AddResultHandler(func(mailBox *mailbox.MailBox) error {
		retMails = mailBox.GetMails()
		return nil
	}, result)

	select {
	case err := <-result:
		return retMails, err
	case <-timeout.Done():
		return retMails, ErrTimeout
	}
}

// 读取邮件
func (m *MailManager) ReadMail(ownerId int64, mailId int64) error {
	mb, err := m.getMailBox(ownerId)
	if err != nil {
		return err
	}

	result := make(chan error, 1)
	timeout, cancel := context.WithTimeout(context.Background(), channelHandleTimeout)
	defer cancel()

	mb.AddResultHandler(func(mailBox *mailbox.MailBox) error {
		return mailBox.ReadMail(mailId)
	}, result)

	select {
	case err := <-result:
		return err
	case <-timeout.Done():
		return ErrTimeout
	}
}

// 获取附件
func (m *MailManager) GainAttachments(ownerId int64, mailId int64) error {
	mb, err := m.getMailBox(ownerId)
	if err != nil {
		return err
	}

	result := make(chan error, 1)
	timeout, cancel := context.WithTimeout(context.Background(), channelHandleTimeout)
	defer cancel()

	mb.AddResultHandler(func(mailBox *mailbox.MailBox) error {
		return mailBox.GainAttachments(mailId)
	}, result)

	select {
	case err := <-result:
		return err
	case <-timeout.Done():
		return ErrTimeout
	}
}
