package item

import (
	"sync"

	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/internal/att"
	"e.coding.net/mmstudio/blade/server/store"
)

// item create pool
var itemPool = &sync.Pool{New: newPoolItemV1}

func NewPoolItem() Item {
	return itemPool.Get().(Item)
}

func GetItemPool() *sync.Pool {
	return itemPool
}

func ReleasePoolItem(x interface{}) {
	itemPool.Put(x)
}

type Item interface {
	store.StoreObjector

	GetOptions() *Options
	Entry() *auto.ItemEntry
	EquipEnchantEntry() *auto.EquipEnchantEntry
	GetAttManager() *att.AttManager

	GetEquipObj() int64
	SetEquipObj(int64)

	CalcAtt()
}

func NewItem(opts ...Option) Item {
	i := NewPoolItem()

	for _, o := range opts {
		o(i.GetOptions())
	}

	return i
}
