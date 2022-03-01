package mail

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/utils/cache"
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
	mu             sync.Mutex
}

func NewMailManager(ctx *cli.Context, m *Mail) *MailManager {
	manager := &MailManager{
		m:              m,
		cacheMailBoxes: cache.New(mailBoxCacheExpire, mailBoxCacheExpire),
	}

	// 邮箱池
	manager.mailBoxPool.New = NewMailBox

	// 邮箱缓存删除时处理
	manager.cacheMailBoxes.OnEvicted(func(k, v interface{}) {
		v.(*MailBox).Stop()
		manager.mailBoxPool.Put(v)
	})

	// 初始化db
	store.GetStore().AddStoreInfo(define.StoreType_Mail, "mail", "_id")
	if err := store.GetStore().MigrateDbTable("mail", "owner_id"); err != nil {
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

func (m *MailManager) KickAllMailBox() {
	m.cacheMailBoxes.DeleteAll()
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

		mb.(*MailBox).Stop()
		store.GetStore().Flush()
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
func (m *MailManager) getMailBox(ownerId int64) (*MailBox, error) {
	if ownerId == -1 {
		return nil, ErrInvalidOwner
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	cache, ok := m.cacheMailBoxes.Get(ownerId)

	if ok {
		mb := cache.(*MailBox)
		if mb.IsTaskRunning() {
			return mb, nil
		}

	} else {

		// 缓存没有，从db加载
		cache = m.mailBoxPool.Get()
		mailbox := cache.(*MailBox)
		mailbox.Init(m.m.ID, m.m.rpcHandler)
		err := mailbox.Load(ownerId)
		if !utils.ErrCheck(err, "mailbox Load failed when MailManager.getMailBox", ownerId) {
			m.mailBoxPool.Put(cache)
			return nil, err
		}

		// 踢掉上一个节点的缓存
		if mailbox.LastSaveNodeId != -1 && mailbox.LastSaveNodeId != int32(m.m.ID) {
			err := m.KickMailBox(mailbox.OwnerId, mailbox.LastSaveNodeId)
			if !utils.ErrCheck(err, "kick mailbox failed", mailbox.OwnerId, mailbox.LastSaveNodeId, m.m.ID) {
				return nil, err
			}
		}

		m.cacheMailBoxes.Set(ownerId, cache, mailBoxCacheExpire)
	}

	mb := cache.(*MailBox)
	mb.InitTask()
	m.wg.Wrap(func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			// 立即删除缓存
			m.cacheMailBoxes.Delete(cache.(*MailBox).OwnerId)
		}()

		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
		for {
			err := cache.(*MailBox).TaskRun(ctx)
			utils.ErrPrint(err, "mailbox run failed", cache.(*MailBox).OwnerId)

			// pull up goroutine when task panic
			if errors.Is(err, task.ErrTaskPanic) {
				continue
			} else {
				break
			}
		}
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
	return m.AddTask(
		ctx,
		ownerId,
		func(c context.Context, p ...interface{}) error {
			mailBox := p[0].(*MailBox)
			return mailBox.AddMail(c, mail)
		},
	)
}

// 删除邮件
func (m *MailManager) DelMail(ctx context.Context, ownerId int64, mailId int64) error {
	return m.AddTask(
		ctx,
		ownerId,
		func(c context.Context, p ...interface{}) error {
			mailBox := p[0].(*MailBox)
			return mailBox.DelMail(c, mailId)
		},
	)
}

// 查询玩家邮件
func (m *MailManager) QueryPlayerMails(ctx context.Context, ownerId int64) (mails []define.Mail, err error) {
	err = m.AddTask(
		ctx,
		ownerId,
		func(c context.Context, p ...interface{}) error {
			mailBox := p[0].(*MailBox)
			mails = mailBox.GetMails(c)
			return nil
		},
	)

	return
}

// 读取邮件
func (m *MailManager) ReadMail(ctx context.Context, ownerId int64, mailId int64) error {
	return m.AddTask(
		ctx,
		ownerId,
		func(c context.Context, p ...interface{}) error {
			mailBox := p[0].(*MailBox)
			return mailBox.ReadMail(c, mailId)
		},
	)
}

// 获取附件
func (m *MailManager) GainAttachments(ctx context.Context, ownerId int64, mailId int64) error {
	return m.AddTask(
		ctx,
		ownerId,
		func(c context.Context, p ...interface{}) error {
			mailBox := p[0].(*MailBox)
			return mailBox.GainAttachments(c, mailId)
		},
	)
}
