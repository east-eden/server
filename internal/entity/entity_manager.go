// Package entity 维护需要下发给客户端的实体缓存。
//
// 对应 kof 的 EntityManager：按实体地址缓存实体快照，增量模式下用 DiffEntity
// 生成差量，并决定下发方法（$update$ / $update_clean$），缓存条目带过期时间。
// kof 直接依赖 loong 的 space.Entity 与自定义 marshaler，这里抽成 Entity/Notifier
// 接口，逻辑保持一致但不再绑定具体序列化框架。
package entity

import (
	"sync"
	"time"

	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog/log"
)

const (
	// EntityExpireSecond 实体快照的过期时间（秒）
	EntityExpireSecond = 600
)

// Entity 参与同步的实体需要实现的最小接口
type Entity interface {
	// Address 返回实体地址，用于缓存索引与客户端路由
	Address() string
	// Name 返回实体名，仅用于日志
	Name() string
	// Clone 返回可安全缓存的深拷贝
	Clone() Entity
	// Marshal 返回全量负载
	Marshal() ([]byte, error)
	// Diff 返回相对旧快照的增量负载
	Diff(old Entity) []byte
}

// Notifier 把实体负载推给客户端
type Notifier interface {
	SendNotify(payload []byte, address, method string) error
}

// EntityData 缓存条目
type EntityData struct {
	Entity          Entity
	ExpireTimestamp int64
}

func newEntityData(entity Entity) *EntityData {
	return &EntityData{
		Entity:          entity,
		ExpireTimestamp: utils.NowInLocation().Add(time.Second * EntityExpireSecond).UnixMilli(),
	}
}

// EntityManager 实体缓存管理器
type EntityManager struct {
	lock         sync.RWMutex
	m            map[string]*EntityData
	increasement bool // 是否使用增量同步
	cleanUpdate  bool // 是否强制全量覆盖
}

func NewEntityManager(increasement bool) *EntityManager {
	return &EntityManager{
		m:            make(map[string]*EntityData, 64),
		increasement: increasement,
	}
}

// EntityCount 当前缓存条目数
func (m *EntityManager) EntityCount() int {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return len(m.m)
}

func (m *EntityManager) SetIncreasement(increasement bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.increasement = increasement
}

func (m *EntityManager) GetIncreasement() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.increasement
}

func (m *EntityManager) SetCleanUpdate(cleanUpdate bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.cleanUpdate = cleanUpdate
}

func (m *EntityManager) GetCleanUpdate() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.cleanUpdate
}

// Reset 清空缓存
func (m *EntityManager) Reset() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.m = make(map[string]*EntityData, 64)
}

// CleanExpired 异步清理过期实体
func (m *EntityManager) CleanExpired() {
	utils.Go(func() {
		now := utils.NowInLocation().UnixMilli()

		m.lock.Lock()
		expired := 0
		for key, data := range m.m {
			if data.ExpireTimestamp < now {
				delete(m.m, key)
				expired++
			}
		}
		m.lock.Unlock()

		log.Info().Int("expired_entity_num", expired).Msg("clean expired entity")
	})
}

// CalEntityBody 计算实体负载。
//
// fullPayload 为 true 或非增量模式时返回全量并刷新快照；
// 增量模式下与旧快照做 diff，无变化时返回空负载。
func (m *EntityManager) CalEntityBody(entity Entity, fullPayload ...bool) ([]byte, string, error) {
	address := entity.Address()

	m.lock.Lock()
	defer m.lock.Unlock()

	if (len(fullPayload) > 0 && fullPayload[0]) || !m.increasement {
		m.m[address] = newEntityData(entity.Clone())

		payload, err := entity.Marshal()
		return payload, address, err
	}

	old, found := m.m[address]
	if !found {
		m.m[address] = newEntityData(entity.Clone())

		payload, err := entity.Marshal()
		return payload, address, err
	}

	payload := entity.Diff(old.Entity)
	if len(payload) > 0 {
		m.m[address] = newEntityData(entity.Clone())
	}

	return payload, address, nil
}

// GetOldEntity 取缓存里的旧快照
func (m *EntityManager) GetOldEntity(address string) (Entity, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	data, found := m.m[address]
	if !found {
		return nil, false
	}

	return data.Entity, true
}

// SendEntity 计算负载并推送给客户端，返回负载与使用的方法
func (m *EntityManager) SendEntity(entity Entity, notifier Notifier, fullPayload ...bool) ([]byte, string) {
	payload, address, err := m.CalEntityBody(entity, fullPayload...)
	if !utils.ErrCheck(err, "CalEntityBody failed", entity.Name()) {
		return payload, ""
	}

	if len(payload) == 0 || notifier == nil {
		return payload, ""
	}

	method := "$update$"
	switch {
	case m.GetCleanUpdate():
		method = "$update_clean$"
	case len(fullPayload) > 0 && fullPayload[0]:
		method = "$update_clean$"
	case !m.GetIncreasement():
		// 非增量模式始终返回全量实体，必须用 clean 覆盖，
		// 否则 map 字段的 key 删除不会同步到前端
		method = "$update_clean$"
	}

	if err := notifier.SendNotify(payload, address, method); err != nil {
		utils.ErrPrint(err, "SendEntity failed", entity.Name(), address)
		return payload, ""
	}

	return payload, method
}
