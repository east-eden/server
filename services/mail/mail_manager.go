package mail

import (
	"errors"
	"sync"
	"time"

	"bitbucket.org/funplus/server/define"
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

// 获取邮箱数据
func (m *MailManager) getMailBox(ownerId int64) (*MailBox, error) {
	if ownerId == -1 {
		return nil, ErrInvalidOwner
	}

	mb, ok := m.cacheMailBoxes.Get(ownerId)

	// 读取最后存储时节点id
	var ownerInfo MailOwnerInfo
	err := store.GetStore().LoadObject(define.StoreType_Mail, ownerId, &ownerInfo)

	// 是否新玩家数据
	var isNew bool

	// 新玩家没有记录
	if errors.Is(err, store.ErrNoResult) {
		isNew = true
	}

	// 老玩家读取邮件失败
	if !isNew && !utils.ErrCheck(err, "LoadObject failed when MailManager.getMailBox", ownerId) {
		return nil, err
	}

	// 如果没有缓存或者最后存储时节点id不是当前节点，删除缓存，重新从store load
	if !ok || int16(ownerInfo.LastSaveNodeId) != m.m.ID {
		m.cacheMailBoxes.Delete(ownerId)

		mb = m.mailBoxPool.Get()

		// 新玩家初始化
		if isNew {
			mailbox := mb.(*MailBox)
			mailbox.Init()
			mailbox.Id = ownerId
			mailbox.LastSaveNodeId = int32(m.m.ID)
			errSave := store.GetStore().SaveObject(define.StoreType_Mail, ownerId, mailbox)
			utils.ErrPrint(errSave, "SaveObject failed when MailManager.getMailBox", ownerId)
		} else {
			// 老玩家加载
			err := store.GetStore().LoadObject(define.StoreType_Mail, ownerId, mb)
			if !utils.ErrCheck(err, "LoadObject failed when MailManager.getMailBox", ownerId) {
				m.mailBoxPool.Put(mb)
				return nil, err
			}
		}

		m.cacheMailBoxes.Set(ownerId, mb, mailBoxCacheExpire)
	}

	return mb.(*MailBox), nil
}

// 创建新邮件
func (m *MailManager) CreateMail(receiverId int64, mail *define.Mail) error {
	mb, err := m.getMailBox(receiverId)
	if err != nil {
		return err
	}

	return mb.AddMail(mail)
}
