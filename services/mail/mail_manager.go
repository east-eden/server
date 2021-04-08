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
	mailBoxCacheExpire = 10 * time.Minute // 邮箱cache缓存10分钟
	ErrInvalidOwner    = errors.New("invalid owner")
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
	if err := store.GetStore().MigrateDbTable("mail", "owner_id", "mail_list._id"); err != nil {
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
		err := mailbox.Load(ownerId)
		if !utils.ErrCheck(err, "mailbox Load failed when MailManager.getMailBox", ownerId) {
			m.mailBoxPool.Put(mb)
			return nil, err
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
func (m *MailManager) CreateMail(ctx context.Context, receiverId int64, mail *define.Mail) error {
	mb, err := m.getMailBox(receiverId)
	if err != nil {
		return err
	}

	err = mb.AddResultHandler(ctx, func(c context.Context, mailBox *mailbox.MailBox) error {
		return mailBox.AddMail(c, mail)
	})

	return err
}

// 删除邮件
func (m *MailManager) DelMail(ctx context.Context, receiverId int64, mailId int64) error {
	mb, err := m.getMailBox(receiverId)
	if err != nil {
		return err
	}

	err = mb.AddResultHandler(ctx, func(c context.Context, mailBox *mailbox.MailBox) error {
		return mailBox.DelMail(c, mailId)
	})

	return err
}

// 查询玩家邮件
func (m *MailManager) QueryPlayerMails(ctx context.Context, ownerId int64) ([]*define.Mail, error) {
	retMails := make([]*define.Mail, 0)
	mb, err := m.getMailBox(ownerId)
	if err != nil {
		return retMails, err
	}

	err = mb.AddResultHandler(ctx, func(c context.Context, mailBox *mailbox.MailBox) error {
		retMails = mailBox.GetMails(c)
		return nil
	})

	return retMails, err
}

// 读取邮件
func (m *MailManager) ReadMail(ctx context.Context, ownerId int64, mailId int64) error {
	mb, err := m.getMailBox(ownerId)
	if err != nil {
		return err
	}

	err = mb.AddResultHandler(ctx, func(c context.Context, mailBox *mailbox.MailBox) error {
		return mailBox.ReadMail(c, mailId)
	})

	return err
}

// 获取附件
func (m *MailManager) GainAttachments(ctx context.Context, ownerId int64, mailId int64) error {
	mb, err := m.getMailBox(ownerId)
	if err != nil {
		return err
	}

	err = mb.AddResultHandler(ctx, func(c context.Context, mailBox *mailbox.MailBox) error {
		return mailBox.GainAttachments(c, mailId)
	})

	return err
}
