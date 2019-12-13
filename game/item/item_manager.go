package item

import (
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type ItemManager struct {
	Owner          define.PluginObj
	mapItem        map[int64]Item
	mapEquipedList map[int64]int64 // map[itemID]heroID

	ds *db.Datastore
	sync.RWMutex
}

func NewItemManager(owner define.PluginObj, ds *db.Datastore) *ItemManager {
	m := &ItemManager{
		Owner:          owner,
		ds:             ds,
		mapItem:        make(map[int64]Item, 0),
		mapEquipedList: make(map[int64]int64, 0),
	}

	return m
}

func (m *ItemManager) LoadFromDB() {
	l := LoadAll(m.ds, m.Owner.GetID())
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
	item.SetOwnerID(m.Owner.GetID())
	item.SetTypeID(entry.ID)
	item.SetEntry(entry)

	m.mapItem[item.GetID()] = item

	return item
}

func (m *ItemManager) newDBItem(i Item) Item {
	item := NewItem(i.GetID())
	item.SetOwnerID(i.GetOwnerID())
	item.SetTypeID(i.GetTypeID())
	item.SetEquipObj(i.GetEquipObj())
	item.SetEntry(global.GetItemEntry(i.GetTypeID()))

	m.mapItem[item.GetID()] = item

	return item
}

func (m *ItemManager) GetItem(id int64) Item {
	return m.mapItem[id]
}

func (m *ItemManager) GetItemNums() int {
	return len(m.mapItem)
}

func (m *ItemManager) GetItemList() []Item {
	list := make([]Item, 0)

	for _, v := range m.mapItem {
		list = append(list, v)
	}

	return list
}

func (m *ItemManager) AddItem(typeID int32) Item {
	itemEntry := global.GetItemEntry(typeID)
	item := m.newEntryItem(itemEntry)
	if item == nil {
		return nil
	}

	m.ds.ORM().Save(item)
	return item
}

func (m *ItemManager) DelItem(id int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(-1)
	delete(m.mapEquipedList, id)
	delete(m.mapItem, id)

	m.ds.ORM().Delete(i)
}

func (m *ItemManager) SetItemEquiped(id int64, objID int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(objID)
	m.mapEquipedList[id] = objID
	m.ds.ORM().Save(i)
}

func (m *ItemManager) SetItemUnEquiped(id int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(-1)
	delete(m.mapEquipedList, id)
	m.ds.ORM().Save(i)
}
