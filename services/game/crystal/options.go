package crystal

import "bitbucket.org/funplus/server/excel/auto"

type Option func(*Options)

// crystal options
type Options struct {
	Id           int64              `bson:"_id" json:"_id"`
	OwnerId      int64              `bson:"owner_id" json:"owner_id"`
	TypeId       int32              `bson:"type_id" json:"type_id"`
	EquipObj     int64              `bson:"equip_obj" json:"equip_obj"`
	Level        int8               `bson:"level" json:"level"`
	Exp          int32              `bson:"exp" json:"exp"`
	ItemEntry    *auto.ItemEntry    `bson:"-" json:"-"`
	CrystalEntry *auto.CrystalEntry `bson:"-" json:"-"`
}

func DefaultOptions() Options {
	return Options{
		Id:           -1,
		OwnerId:      -1,
		TypeId:       -1,
		EquipObj:     -1,
		Level:        0,
		Exp:          0,
		ItemEntry:    nil,
		CrystalEntry: nil,
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

func ItemEntry(entry *auto.ItemEntry) Option {
	return func(o *Options) {
		o.ItemEntry = entry
	}
}

func CrystalEntry(entry *auto.CrystalEntry) Option {
	return func(o *Options) {
		o.CrystalEntry = entry
	}
}

func Level(lv int8) Option {
	return func(o *Options) {
		o.Level = lv
	}
}

func Exp(exp int32) Option {
	return func(o *Options) {
		o.Exp = exp
	}
}
