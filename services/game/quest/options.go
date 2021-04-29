package quest

import "bitbucket.org/funplus/server/define"

type Option func(*Options)
type Options struct {
	Id      int32       `bson:"_id" json:"_id"`           // 任务id
	OwnerId int64       `bson:"owner_id" json:"owner_id"` // 玩家id
	Objs    []*QuestObj `bson:"objs" json:"objs"`         // 任务目标数据
	State   int32       `bson:"state" json:"state"`       // 任务状态
}

func DefaultOptions() Options {
	return Options{
		Id:      -1,
		OwnerId: -1,
		Objs:    make([]*QuestObj, 0, Quest_Max_Obj),
		State:   define.Quest_State_Type_Opened,
	}
}

func WithId(id int32) Option {
	return func(o *Options) {
		o.Id = id
	}
}

func WithOwnerId(ownerId int64) Option {
	return func(o *Options) {
		o.OwnerId = ownerId
	}
}

func WithObjs(objs []*QuestObj) Option {
	return func(o *Options) {
		o.Objs = objs
	}
}

func WithState(state int32) Option {
	return func(o *Options) {
		o.State = state
	}
}
