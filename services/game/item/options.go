package item

import (
	"e.coding.net/mmstudio/blade/server/excel/auto"
)

type Option func(*Options)

// item options
type Options struct {
	Id      int64 `bson:"_id" json:"_id"`
	OwnerId int64 `bson:"owner_id" json:"owner_id"`
	TypeId  int32 `bson:"type_id" json:"type_id"`
	Num     int32 `bson:"num" json:"num"`

	EquipObj          int64                   `bson:"equip_obj" json:"equip_obj"`
	Entry             *auto.ItemEntry         `bson:"-" json:"-"`
	EquipEnchantEntry *auto.EquipEnchantEntry `bson:"-" json:"-"`
}

func DefaultOptions() Options {
	return Options{
		Id:                -1,
		OwnerId:           -1,
		TypeId:            -1,
		Num:               0,
		EquipObj:          -1,
		Entry:             nil,
		EquipEnchantEntry: nil,
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

func Num(n int32) Option {
	return func(o *Options) {
		o.Num = n
	}
}

func EquipObj(obj int64) Option {
	return func(o *Options) {
		o.EquipObj = obj
	}
}

func Entry(entry *auto.ItemEntry) Option {
	return func(o *Options) {
		o.Entry = entry
	}
}

func EquipEnchantEntry(entry *auto.EquipEnchantEntry) Option {
	return func(o *Options) {
		o.EquipEnchantEntry = entry
	}
}
