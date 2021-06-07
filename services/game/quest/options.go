package quest

import (
	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/services/game/event"
)

//////////////////////////////////////////////////////
// Quest Options
type Option func(*Options)
type Options struct {
	QuestId int32 `bson:"quest_id" json:"quest_id"` // 任务id
	// OwnerId int64            `bson:"owner_id" json:"owner_id"` // 玩家id
	Objs  []*QuestObj      `bson:"objs" json:"objs"`   // 任务目标数据
	State int32            `bson:"state" json:"state"` // 任务状态
	Entry *auto.QuestEntry `bson:"-" json:"-"`
}

func DefaultOptions() Options {
	return Options{
		QuestId: -1,
		// OwnerId: -1,
		Objs:  make([]*QuestObj, 0, Quest_Max_Obj),
		State: define.Quest_State_Type_Opened,
		Entry: nil,
	}
}

func WithId(id int32) Option {
	return func(o *Options) {
		o.QuestId = id
	}
}

func WithState(state int32) Option {
	return func(o *Options) {
		o.State = state
	}
}

func WithEntry(entry *auto.QuestEntry) Option {
	return func(o *Options) {
		o.Entry = entry

		for k := range entry.ObjTypes {
			if entry.ObjTypes[k] == -1 {
				break
			}

			o.Objs = append(o.Objs, &QuestObj{
				Type:      entry.ObjTypes[k],
				Count:     0,
				Completed: false,
			})
		}
	}
}

/////////////////////////////////////////////////////
// QuestManager Options
type ManagerOption func(*ManagerOptions)
type ManagerOptions struct {
	ownerId           int64               `bson:"-" json:"-"` // 任务挂载者id
	ownerType         int32               `bson:"-" json:"-"` // 任务挂载者类型
	storeType         int                 `bson:"-" json:"-"` // 数据库存储类型
	additionalQuestId []int32             `bson:"-" json:"-"` // 额外任务id
	eventManager      *event.EventManager `bson:"-" json:"-"` // 事件管理器
	questChangedCb    func(*Quest)        `bson:"-" json:"-"` // 任务变更回调
	questRewardCb     func(*Quest)        `bson:"-" json:"-"` // 任务奖励回调
}

func DefaultManagerOptions() ManagerOptions {
	return ManagerOptions{
		ownerId:           -1,
		ownerType:         define.QuestOwner_Type_Player,
		storeType:         -1,
		additionalQuestId: make([]int32, 0, 1),
		eventManager:      nil,
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

func WithManagerStoreType(storeType int) ManagerOption {
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
