package item

import (
	"github.com/east-eden/server/excel/config"
)

type ItemOption func(*ItemOptions)

// ItemOptions 物品（可堆叠账本）数据。
//
// 与 kof 一致：一个配置 id 只对应一条堆叠记录，Id 即配置 id（EntryId），
// 不再有“每件装备一个实例 id”的概念。
type ItemOptions struct {
	Id         int64              `bson:"_id" json:"_id"`
	OwnerId    int64              `bson:"owner_id" json:"owner_id"`
	EntryId    int32              `bson:"entry_id" json:"entry_id"`
	Num        int64              `bson:"num" json:"num"`
	CreateTime int64              `bson:"create_time" json:"create_time"`
	GoodsEntry *config.GoodsEntry `bson:"-" json:"-"`
}

func DefaultItemOptions() ItemOptions {
	return ItemOptions{
		Id:         -1,
		OwnerId:    -1,
		EntryId:    -1,
		Num:        0,
		CreateTime: 0,
		GoodsEntry: nil,
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

func EntryId(id int32) ItemOption {
	return func(o *ItemOptions) {
		o.EntryId = id
	}
}

func Num(n int64) ItemOption {
	return func(o *ItemOptions) {
		o.Num = n
	}
}

func CreateTime(t int64) ItemOption {
	return func(o *ItemOptions) {
		o.CreateTime = t
	}
}

func GoodsEntry(entry *config.GoodsEntry) ItemOption {
	return func(o *ItemOptions) {
		o.GoodsEntry = entry
	}
}
