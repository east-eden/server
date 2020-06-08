package rune

import "github.com/yokaiio/yokai_server/define"

type Option func(*Options)

// rune options
type Options struct {
	Id       int64             `bson:"_id" redis:"_id"`
	OwnerId  int64             `bson:"owner_id" redis:"owner_id"`
	TypeId   int32             `bson:"type_id" redis:"type_id"`
	EquipObj int64             `bson:"equip_obj" redis:"equip_obj"`
	Entry    *define.RuneEntry `bson:"-" redis:"-"`
}

func DefaultOptions() Options {
	return Options{
		Id:       -1,
		OwnerId:  -1,
		TypeId:   -1,
		EquipObj: -1,
		Entry:    nil,
	}
}

func Id(id int64) Option {
	return func(o *Options) {
		o.Id = id
	}
}

func OwnerId(id int64) Option {
	return func(o *Options) {
		o.OwnerId = id
	}
}

func TypeId(id int32) Option {
	return func(o *Options) {
		o.TypeId = id
	}
}

func EquipObj(obj int64) Option {
	return func(o *Options) {
		o.EquipObj = obj
	}
}

func Entry(entry *define.RuneEntry) Option {
	return func(o *Options) {
		o.Entry = entry
	}
}
