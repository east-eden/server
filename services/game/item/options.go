package item

import (
	"github.com/east-eden/server/excel/auto"
)

type ItemOption func(*ItemOptions)
type EquipOption func(*EquipOptions)

// item options
type ItemOptions struct {
	Id         int64           `bson:"_id" json:"_id"`
	OwnerId    int64           `bson:"owner_id" json:"owner_id"`
	TypeId     int32           `bson:"type_id" json:"type_id"`
	Num        int32           `bson:"num" json:"num"`
	CreateTime int64           `bson:"create_time" json:"create_time"`
	Entry      *auto.ItemEntry `bson:"-" json:"-"`
}

func DefaultItemOptions() ItemOptions {
	return ItemOptions{
		Id:         -1,
		OwnerId:    -1,
		TypeId:     -1,
		Num:        0,
		CreateTime: 0,
		Entry:      nil,
	}
}

func Id(id int64) ItemOption {
	return func(o *ItemOptions) {
		o.Id = id
	}
}

func OwnerId(id int64) ItemOption {
	return func(o *ItemOptions) {
		o.OwnerId = id
	}
}

func TypeId(id int32) ItemOption {
	return func(o *ItemOptions) {
		o.TypeId = id
	}
}

func Num(n int32) ItemOption {
	return func(o *ItemOptions) {
		o.Num = n
	}
}

func Entry(entry *auto.ItemEntry) ItemOption {
	return func(o *ItemOptions) {
		o.Entry = entry
	}
}

// equip options
type EquipOptions struct {
	Exp               int32                   `bson:"exp" json:"exp"`
	Level             int8                    `bson:"level" json:"level"`
	Promote           int8                    `bson:"promote" json:"promote"`
	Lock              bool                    `bson:"lock" json:"lock"`
	EquipObj          int64                   `bson:"equip_obj" json:"equip_obj"`
	EquipEnchantEntry *auto.EquipEnchantEntry `bson:"-" json:"-"`
}

func DefaultEquipOptions() EquipOptions {
	return EquipOptions{
		Exp:               0,
		Level:             1,
		Promote:           0,
		Lock:              false,
		EquipObj:          -1,
		EquipEnchantEntry: nil,
	}
}

func Level(lv int8) EquipOption {
	return func(o *EquipOptions) {
		o.Level = lv
	}
}

func Exp(exp int32) EquipOption {
	return func(o *EquipOptions) {
		o.Exp = exp
	}
}

func Promote(prom int8) EquipOption {
	return func(o *EquipOptions) {
		o.Promote = prom
	}
}

func Lock(lock bool) EquipOption {
	return func(o *EquipOptions) {
		o.Lock = lock
	}
}

func EquipObj(obj int64) EquipOption {
	return func(o *EquipOptions) {
		o.EquipObj = obj
	}
}

func EquipEnchantEntry(entry *auto.EquipEnchantEntry) EquipOption {
	return func(o *EquipOptions) {
		o.EquipEnchantEntry = entry
	}
}
