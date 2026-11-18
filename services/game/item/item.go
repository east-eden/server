package item

import (
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
)

// Itemface 物品接口
type Itemface interface {
	InitItem(opts ...ItemOption)
	Opts() *ItemOptions
	GenItemPB() *pbGlobal.Item
}

func NewItem() *Item {
	i := &Item{}
	i.ItemOptions = DefaultItemOptions()
	return i
}

// Item 可堆叠物品
type Item struct {
	ItemOptions `bson:"inline" json:",inline"`
}

func (i *Item) InitItem(opts ...ItemOption) {
	for _, o := range opts {
		o(&i.ItemOptions)
	}
}

func (i *Item) Opts() *ItemOptions {
	return &i.ItemOptions
}

// Entry 返回道具配置。配置表热更后需要重新回填。
func (i *Item) Entry() *config.GoodsEntry {
	return i.ItemOptions.GoodsEntry
}

func (i *Item) GenItemPB() *pbGlobal.Item {
	return &pbGlobal.Item{
		Id:     i.Id,
		TypeId: i.EntryId,
		Num:    i.Num,
	}
}
