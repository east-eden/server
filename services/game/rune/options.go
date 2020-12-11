package rune

import "github.com/east-eden/server/define"

type Option func(*Options)

// rune options
type Options struct {
	Id       int64             `bson:"_id" json:"_id"`
	OwnerId  int64             `bson:"owner_id" json:"owner_id"`
	TypeId   int32             `bson:"type_id" json:"type_id"`
	EquipObj int64             `bson:"equip_obj" json:"equip_obj"`
	Entry    *define.RuneEntry `bson:"-" json:"-"`
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
