package quest

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/services/game/event"
)

/////////////////////////////////////////////////////
// QuestManager Options

type ManagerOption func(*ManagerOptions)

// ManagerOptions 任务管理器配置
type ManagerOptions struct {
	ownerId           int64               `bson:"-" json:"-"` // 任务挂载者id
	ownerType         int32               `bson:"-" json:"-"` // 任务挂载者类型
	storeType         define.StoreType    `bson:"-" json:"-"` // 数据库存储类型
	additionalQuestId []int32             `bson:"-" json:"-"` // 额外任务id
	eventManager      *event.EventManager `bson:"-" json:"-"` // 事件管理器
	ownerLevel        func() int32        `bson:"-" json:"-"` // 获取挂载者等级（角色等级类条件用）
	questChangedCb    func(*Quest)        `bson:"-" json:"-"` // 任务变更回调
	questRewardCb     func(*Quest)        `bson:"-" json:"-"` // 任务奖励回调
}

func DefaultManagerOptions() ManagerOptions {
	return ManagerOptions{
		ownerId:           -1,
		ownerType:         define.QuestOwner_Type_Player,
		storeType:         define.StoreType_Player,
		additionalQuestId: make([]int32, 0, 1),
		eventManager:      nil,
		ownerLevel:        func() int32 { return 1 },
		questChangedCb:    func(*Quest) {},
		questRewardCb:     func(*Quest) {},
	}
}

func WithManagerOwnerId(ownerId int64) ManagerOption {
	return func(o *ManagerOptions) {
		o.ownerId = ownerId
	}
}

func WithManagerOwnerType(tp int32) ManagerOption {
	return func(o *ManagerOptions) {
		o.ownerType = tp
	}
}

func WithManagerStoreType(storeType define.StoreType) ManagerOption {
	return func(o *ManagerOptions) {
		o.storeType = storeType
	}
}

func WithManagerAdditionalQuestId(ids ...int32) ManagerOption {
	return func(o *ManagerOptions) {
		o.additionalQuestId = append(o.additionalQuestId, ids...)
	}
}

func WithManagerEventManager(m *event.EventManager) ManagerOption {
	return func(o *ManagerOptions) {
		o.eventManager = m
	}
}

func WithManagerOwnerLevel(fn func() int32) ManagerOption {
	return func(o *ManagerOptions) {
		o.ownerLevel = fn
	}
}

func WithManagerQuestChangedCb(cb func(*Quest)) ManagerOption {
	return func(o *ManagerOptions) {
		o.questChangedCb = cb
	}
}

func WithManagerQuestRewardCb(cb func(*Quest)) ManagerOption {
	return func(o *ManagerOptions) {
		o.questRewardCb = cb
	}
}
