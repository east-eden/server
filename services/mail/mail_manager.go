package mail

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/services/mail/mailbox"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"bitbucket.org/funplus/server/utils/cache"
	"github.com/hellodudu/task"
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
		v.(*mailbox.MailBox).Stop()
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
	log.Info().Msg("MailManager context done...")
	return nil
}

func (m *MailManager) Exit(ctx *cli.Context) {
	m.wg.Wait()
	log.Info().Msg("MailManager exit...")
}

// 提掉邮箱缓存
func (m *MailManager) KickMailBox(ownerId int64, mailNodeId int32) error {
	if ownerId == -1 {
		return nil
	}

	// 踢掉本服mailbox
	if mailNodeId == int32(m.m.ID) {
		mb, ok := m.cacheMailBoxes.Get(ownerId)
		if !ok {
			return nil
		}

		mb.(*mailbox.MailBox).Stop()
		return nil

	} else {
		// mail节点不存在的话不用发送rpc
		nodeId := fmt.Sprintf("mail-%d", mailNodeId)
		srvs, err := m.m.mi.srv.Server().Options().Registry.GetService("mail")
		if err != nil {
			return nil
		}

		hit := false
		for _, srv := range srvs {
			for _, node := range srv.Nodes {
				if node.Id == nodeId {
					hit = true
					break
				}
			}
		}

		if !hit {
			return nil
		}

		// 发送rpc踢掉其他服mailbox
		rs, err := m.m.rpcHandler.CallKickMailBox(ownerId, mailNodeId)
		if !utils.ErrCheck(err, "kick mail box failed", ownerId, mailNodeId, rs) {
			return err
		}

		// rpc调用成功
		if rs.GetOwnerId() == ownerId {
			return nil
		}

		return errors.New("kick mail box invalid error")
	}
}

// 获取邮箱数据
func (m *MailManager) getMailBox(ownerId int64) (*mailbox.MailBox, error) {
	if ownerId == -1 {
		return nil, ErrInvalidOwner
	}

	cache, ok := m.cacheMailBoxes.Get(ownerId)

	if ok {
		mb := cache.(*mailbox.MailBox)
		if mb.IsTaskRunning() {
			return mb, nil
		}

	} else {

		// 缓存没有，从db加载
		cache = m.mailBoxPool.Get()
		mailbox := cache.(*mailbox.MailBox)
		mailbox.Init(m.m.ID)
		err := mailbox.Load(ownerId)
		if !utils.ErrCheck(err, "mailbox Load failed when MailManager.getMailBox", ownerId) {
			m.mailBoxPool.Put(cache)
			return nil, err
		}

		// 踢掉上一个节点的缓存
		if mailbox.LastSaveNodeId != -1 && mailbox.LastSaveNodeId != int32(m.m.ID) {
			err := m.KickMailBox(mailbox.Id, mailbox.LastSaveNodeId)
			if !utils.ErrCheck(err, "kick mailbox failed", mailbox.Id, mailbox.LastSaveNodeId, m.m.ID) {
				return nil, err
			}
		}

		m.cacheMailBoxes.Set(ownerId, cache, mailBoxCacheExpire)
	}

	mb := cache.(*mailbox.MailBox)
	mb.InitTask()
	mb.ResetTaskTimeout()
	m.wg.Wrap(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			// 立即删除缓存
			m.cacheMailBoxes.Delete(cache.(*mailbox.MailBox).Id)
		}()

		ctx := utils.WithSignaledCancel(context.Background())
		err := cache.(*mailbox.MailBox).TaskRun(ctx)
		utils.ErrPrint(err, "mailbox run failed", cache.(*mailbox.MailBox).Id)
	})

	return mb, nil
}

func (m *MailManager) AddTask(ctx context.Context, ownerId int64, fn task.TaskHandler) error {
	mb, err := m.getMailBox(ownerId)
	if err != nil {
		return err
	}

	return mb.AddTask(ctx, fn, mb)
}

// 创建新邮件
func (m *MailManager) CreateMail(ctx context.Context, ownerId int64, mail *define.Mail) error {
	fn := func(c context.Context, p ...interface{}) error {
		mailBox := p[0].(*mailbox.MailBox)

		return mailBox.AddMail(c, mail)
	}

	return m.AddTask(ctx, ownerId, fn)
}

// 删除邮件
func (m *MailManager) DelMail(ctx context.Context, ownerId int64, mailId int64) error {
	fn := func(c context.Context, p ...interface{}) error {
		mailBox := p[0].(*mailbox.MailBox)

		return mailBox.DelMail(c, mailId)
	}

	return m.AddTask(ctx, ownerId, fn)
}

// 查询玩家邮件
func (m *MailManager) QueryPlayerMails(ctx context.Context, ownerId int64) ([]*define.Mail, error) {
	retMails := make([]*define.Mail, 0)

	fn := func(c context.Context, p ...interface{}) error {
		mailBox := p[0].(*mailbox.MailBox)

		retMails = mailBox.GetMails(c)
		return nil
	}

	err := m.AddTask(ctx, ownerId, fn)
	return retMails, err
}

// 读取邮件
func (m *MailManager) ReadMail(ctx context.Context, ownerId int64, mailId int64) error {
	fn := func(c context.Context, p ...interface{}) error {
		mailBox := p[0].(*mailbox.MailBox)

		return mailBox.ReadMail(c, mailId)
	}

	return m.AddTask(ctx, ownerId, fn)
}

// 获取附件
func (m *MailManager) GainAttachments(ctx context.Context, ownerId int64, mailId int64) error {
	fn := func(c context.Context, p ...interface{}) error {
		mailBox := p[0].(*mailbox.MailBox)

		return mailBox.GainAttachments(c, mailId)
	}

	return m.AddTask(ctx, ownerId, fn)
}
