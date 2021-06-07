package collection

import (
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/services/game/event"
	"e.coding.net/mmstudio/blade/server/services/game/quest"
)

type Option func(*Options)

// collection options
type Options struct {
	Id            int64                 `bson:"_id" json:"_id"`
	TypeId        int32                 `bson:"type_id" json:"type_id"`
	OwnerId       int64                 `bson:"owner_id" json:"owner_id"`
	Active        bool                  `bson:"active" json:"active"`
	Wakeup        bool                  `bson:"wakeup" json:"wakeup"`
	Star          int8                  `bson:"star" json:"star"`
	BoxId         int32                 `bson:"box_id" json:"box_id"`
	Entry         *auto.CollectionEntry `bson:"-" json:"-"`
	eventManager  *event.EventManager   `bson:"-" json:"-"`
	questUpdateCb func(*quest.Quest)    `bson:"-" json:"-"`
}

func DefaultOptions() Options {
	return Options{
		Id:            -1,
		TypeId:        -1,
		OwnerId:       -1,
		Active:        false,
		Wakeup:        false,
		Star:          0,
		BoxId:         -1,
		Entry:         nil,
		eventManager:  nil,
		questUpdateCb: func(*quest.Quest) {},
	}
}

func Id(id int64) Option {
	return func(o *Options) {
		o.Id = id
	}
}

func TypeId(typeId int32) Option {
	return func(o *Options) {
		o.TypeId = typeId
	}
}

func OwnerId(id int64) Option {
	return func(o *Options) {
		o.OwnerId = id
	}
}

func Star(star int8) Option {
	return func(o *Options) {
		o.Star = star
	}
}

func Entry(entry *auto.CollectionEntry) Option {
	return func(o *Options) {
		o.Entry = entry
	}
}

func EventManager(m *event.EventManager) Option {
	return func(o *Options) {
		o.eventManager = m
	}
}

func QuestUpdateCb(cb func(*quest.Quest)) Option {
	return func(o *Options) {
		o.questUpdateCb = cb
	}
}
