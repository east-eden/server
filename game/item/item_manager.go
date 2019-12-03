package item

import (
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/define"
	"github.com/yokaiio/yokai_server/game/global"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type ItemManager struct {
	OwnerID int64
	mapItem map[int64]Item

	ds *db.Datastore
	sync.RWMutex
}

func NewItemManager(ownerID int64, ds *db.Datastore) *ItemManager {
	m := &ItemManager{
		OwnerID: ownerID,
		ds:      ds,
		mapItem: make(map[int64]Item, 0),
	}

	return m
}

func (m *ItemManager) LoadFromDB() {
	l := LoadAll(m.ds, m.OwnerID)
	sliceItem := make([]Item, 0)

	listItem := reflect.ValueOf(l)
	if listItem.Kind() != reflect.Slice {
		logger.Error("load item returns non-slice type")
		return
	}

	for n := 0; n < listItem.Len(); n++ {
		p := listItem.Index(n)
		sliceItem = append(sliceItem, p.Interface().(Item))
	}

	for _, v := range sliceItem {
		m.newDBItem(v)

		maxID, err := utils.GeneralIDGet(define.Plugin_Item)
		if err != nil {
			logger.Fatal(err)
			return
		}

		if v.GetID() >= maxID {
			utils.GeneralIDSet(define.Plugin_Item, v.GetID())
		}
	}
}

func (m *ItemManager) Save(i Item) {
	m.ds.ORM().Save(i)
}

func (m *ItemManager) newEntryItem(entry *define.ItemEntry) Item {
	if entry == nil {
		logger.Error("newEntryItem with nil ItemEntry")
		return nil
	}

	id, err := utils.GeneralIDGen(define.Plugin_Item)
	if err != nil {
		logger.Error(err)
		return nil
	}

	item := NewItem(id)
	item.SetOwnerID(m.OwnerID)
	item.SetTypeID(entry.ID)
	item.SetEntry(entry)

	m.Lock()
	defer m.Unlock()
	m.mapItem[item.GetID()] = item

	return item
}

func (m *ItemManager) newDBItem(i Item) Item {
	item := NewItem(i.GetID())
	item.SetOwnerID(i.GetOwnerID())
	item.SetTypeID(i.GetTypeID())
	item.SetEntry(global.GetItemEntry(i.GetTypeID()))

	m.Lock()
	defer m.Unlock()
	m.mapItem[item.GetID()] = item

	return item
}

func (m *ItemManager) GetItem(id int64) Item {
	m.RLock()
	defer m.RUnlock()
	return m.mapItem[id]
}

func (m *ItemManager) GetItemList() []Item {
	list := make([]Item, len(m.mapItem))

	m.RLock()
	for _, v := range m.mapItem {
		list = append(list, v)
	}
	m.RUnlock()

	return list
}

func (m *ItemManager) AddItem(typeID int32) Item {
	itemEntry := global.GetItemEntry(typeID)
	item := m.newEntryItem(itemEntry)
	if item == nil {
		return nil
	}

	m.Save(item)
	return item
}
