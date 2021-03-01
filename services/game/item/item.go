package item

import (
	"sync"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
)

// 物品接口
type Itemface interface {
	InitItem(opts ...ItemOption)
	GetType() define.ItemType
	Ops() *ItemOptions
	OnDelete()
}

// item create pool
var itemPool = &sync.Pool{
	New: func() interface{} {
		return &Item{
			ItemOptions: DefaultItemOptions(),
		}
	},
}

// equip create pool
var equipPool = &sync.Pool{
	New: func() interface{} {
		e := &Equip{
			Item: Item{
				ItemOptions: DefaultItemOptions(),
			},
			EquipOptions: DefaultEquipOptions(),
		}
		e.attManager = NewEquipAttManager(e)
		return e
	},
}

func NewPoolItem(tp define.ItemType) Itemface {
	if tp == define.Item_TypeEquip {
		return equipPool.Get().(Itemface)
	}

	return itemPool.Get().(Itemface)
}

func GetItemPool(tp define.ItemType) *sync.Pool {
	if tp == define.Item_TypeEquip {
		return equipPool
	}

	return itemPool
}

func NewItem(tp define.ItemType) Itemface {
	return NewPoolItem(tp)
}

func GetContainerType(tp define.ItemType) define.ContainerType {
	switch tp {
	case define.Item_TypeItem:
		fallthrough
	case define.Item_TypePresent:
		return define.Container_Material

	case define.Item_TypeEquip:
		return define.Container_Equip
	}

	return define.Container_Null
}

type Item struct {
	ItemOptions `bson:"inline" json:",inline"`
}

func (i *Item) InitItem(opts ...ItemOption) {
	for _, o := range opts {
		o(&i.ItemOptions)
	}
}

func (i *Item) GetType() define.ItemType {
	return define.ItemType(i.Entry().Type)
}

func (i *Item) OnDelete() {

}

func (i *Item) Ops() *ItemOptions {
	return &i.ItemOptions
}

func (i *Item) GetID() int64 {
	return i.ItemOptions.Id
}

func (i *Item) GetOwnerID() int64 {
	return i.ItemOptions.OwnerId
}

func (i *Item) GetTypeID() int32 {
	return i.ItemOptions.TypeId
}

func (i *Item) Entry() *auto.ItemEntry {
	return i.ItemOptions.Entry
}
