package item

import "github.com/yokaiio/yokai_server/define"

type Option func(*Options)

// item options
type Options struct {
	Id      int64 `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerId int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	TypeId  int32 `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`
	Num     int32 `gorm:"type:int(10);column:num;default:0;not null" bson:"num"`

	EquipObj          int64                     `gorm:"type:bigint(20);column:equip_obj;default:-1;not null" bson:"equip_obj"`
	Entry             *define.ItemEntry         `gorm:"-" bson:"-"`
	EquipEnchantEntry *define.EquipEnchantEntry `gorm:"-" bson:"-"`
}

func DefaultOptions() *Options {
	return &Options{
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

func Entry(entry *define.ItemEntry) Option {
	return func(o *Options) {
		o.Entry = entry
	}
}

func EquipEnchantEntry(entry *define.EquipEnchantEntry) Option {
	return func(o *Options) {
		o.EquipEnchantEntry = entry
	}
}
