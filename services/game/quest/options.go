package quest

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
)

type Option func(*Options)
type Options struct {
	QuestId int32            `bson:"quest_id" json:"quest_id"` // 任务id
	OwnerId int64            `bson:"owner_id" json:"owner_id"` // 玩家id
	Objs    []*QuestObj      `bson:"objs" json:"objs"`         // 任务目标数据
	State   int32            `bson:"state" json:"state"`       // 任务状态
	Entry   *auto.QuestEntry `bson:"-" json:"-"`
}

func DefaultOptions() Options {
	return Options{
		QuestId: -1,
		OwnerId: -1,
		Objs:    make([]*QuestObj, 0, Quest_Max_Obj),
		State:   define.Quest_State_Type_Opened,
		Entry:   nil,
	}
}

func WithId(id int32) Option {
	return func(o *Options) {
		o.QuestId = id
	}
}

func WithOwnerId(ownerId int64) Option {
	return func(o *Options) {
		o.OwnerId = ownerId
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
