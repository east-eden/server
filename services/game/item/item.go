package item

import (
	"sync"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	"bitbucket.org/east-eden/server/internal/att"
)

// item create pool
var itemPool = &sync.Pool{New: newPoolItem}

func NewPoolItem() *Item {
	return itemPool.Get().(*Item)
}

func GetItemPool() *sync.Pool {
	return itemPool
}

func ReleasePoolItem(x interface{}) {
	itemPool.Put(x)
}

func NewItem(opts ...Option) *Item {
	i := NewPoolItem()

	for _, o := range opts {
		o(i.GetOptions())
	}

	return i
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
	Options    `bson:"inline" json:",inline"`
	attManager *att.AttManager `json:"-" bson:"-"`
}

func newPoolItem() interface{} {
	h := &Item{
		Options: DefaultOptions(),
	}

	h.attManager = att.NewAttManager(-1)

	return h
}

// StoreObjector interface
func (i *Item) AfterLoad() error {
	return nil
}

func (i *Item) GetObjID() int64 {
	return i.Options.Id
}

func (i *Item) GetStoreIndex() int64 {
	return i.Options.OwnerId
}

func (i *Item) GetOptions() *Options {
	return &i.Options
}

func (i *Item) GetID() int64 {
	return i.Options.Id
}

func (i *Item) GetOwnerID() int64 {
	return i.Options.OwnerId
}

func (i *Item) GetTypeID() int32 {
	return i.Options.TypeId
}

func (i *Item) GetAttManager() *att.AttManager {
	return i.attManager
}

func (i *Item) Entry() *auto.ItemEntry {
	return i.Options.Entry
}

func (i *Item) EquipEnchantEntry() *auto.EquipEnchantEntry {
	return i.Options.EquipEnchantEntry
}

func (i *Item) GetEquipObj() int64 {
	return i.Options.EquipObj
}

func (i *Item) SetEquipObj(obj int64) {
	i.Options.EquipObj = obj
}

func (i *Item) CalcAtt() {
	i.attManager.Reset()
}
