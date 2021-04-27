package item

import (
	"sync"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbCommon "bitbucket.org/funplus/server/proto/global/common"
	"github.com/shopspring/decimal"
)

// 物品接口
type Itemface interface {
	InitItem(opts ...ItemOption)
	GetType() int32
	Opts() *ItemOptions
}

// item create pool
var itemPool = &sync.Pool{
	New: func() interface{} {
		return &Item{}
	},
}

// equip create pool
var equipPool = &sync.Pool{
	New: func() interface{} {
		return &Equip{}
	},
}

// crystal create pool
var crystalPool = &sync.Pool{
	New: func() interface{} {
		return &Crystal{}
	},
}

func NewPoolItem(tp int32) Itemface {
	switch tp {
	case define.Item_TypeEquip:
		e := equipPool.Get().(*Equip)
		e.Item.ItemOptions = DefaultItemOptions()
		e.EquipOptions = DefaultEquipOptions()
		e.attManager = NewEquipAttManager(e)
		return e

	case define.Item_TypeCrystal:
		c := crystalPool.Get().(*Crystal)
		c.Item.ItemOptions = DefaultItemOptions()
		c.CrystalOptions = DefaultCrystalOptions()
		c.MainAtt.AttRepoId = -1
		c.MainAtt.AttRandRatio = decimal.NewFromInt32(0)
		c.ViceAtts = make([]CrystalAtt, 0, 20)
		c.attManager = NewCrystalAttManager(c)
		return c

	default:
		i := itemPool.Get().(*Item)
		i.ItemOptions = DefaultItemOptions()
		return i
	}
}

func GetItemPool(tp int32) *sync.Pool {
	switch tp {
	case define.Item_TypeEquip:
		return equipPool
	case define.Item_TypeCrystal:
		return crystalPool
	default:
		return itemPool
	}
}

func NewItem(tp int32) Itemface {
	return NewPoolItem(tp)
}

func GetContainerType(tp int32) int32 {
	switch tp {
	case define.Item_TypeItem, define.Item_TypePresent:
		return define.Container_Material

	case define.Item_TypeEquip:
		return define.Container_Equip
	case define.Item_TypeCrystal:
		return define.Container_Crystal
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

func (i *Item) GetType() int32 {
	return i.Entry().Type
}

func (i *Item) Opts() *ItemOptions {
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
	return i.ItemOptions.ItemEntry
}

func (i *Item) GenItemPB() *pbCommon.Item {
	pb := &pbCommon.Item{
		Id:     i.Id,
		TypeId: i.TypeId,
		Num:    i.Num,
	}

	return pb
}
