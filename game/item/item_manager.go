package item

import "sync/atomic"

type ItemManager struct {
	idGen   atomic.Value
	mapItem map[int64]Item
}

func NewItemManager() *ItemManager {
	m := &ItemManager{
		mapItem: make(map[int64]Item, 0),
	}

	m.idGen.Store(int64(0))
	return m
}

func (m *ItemManager) GenID() int64 {
	id := m.idGen.Load().(int64) + 1
	m.idGen.Store(id)
	return id
}

func (m *ItemManager) NewItem(typeID int32) Item {
	id := m.GenID()
	item := NewItem(id, typeID)
	m.mapItem[item.ID()] = item
	return item
}
