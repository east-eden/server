package item

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ItemManager struct {
	Owner          define.PluginObj
	mapItem        map[int64]Item
	mapEquipedList map[int64]int64 // map[itemID]heroID

	ds   *db.Datastore
	coll *mongo.Collection
	sync.RWMutex
}

func NewItemManager(owner define.PluginObj, ds *db.Datastore) *ItemManager {
	m := &ItemManager{
		Owner:          owner,
		ds:             ds,
		mapItem:        make(map[int64]Item, 0),
		mapEquipedList: make(map[int64]int64, 0),
	}

	m.coll = ds.Database().Collection(m.TableName())

	return m
}

func (m *ItemManager) TableName() string {
	return "item"
}

// interface of cost_loot
func (m *ItemManager) GetCostLootType() int32 {
	return define.CostLoot_Item
}

func (m *ItemManager) CanCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("item manager check item<%d> cost failed, wrong number<%d>", typeMisc, num)
	}

	var fixNum int32 = 0
	for _, v := range m.mapItem {
		if v.GetTypeID() == typeMisc {
			_, ok := m.mapEquipedList[v.GetID()]
			if !ok {
				fixNum += v.GetNum()
			}
		}
	}

	if fixNum >= num {
		return nil
	}

	return fmt.Errorf("not enough item<%d>, num<%d>", typeMisc, num)
}

func (m *ItemManager) DoCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("item manager cost item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	return m.CostItemByTypeID(typeMisc, num)
}

func (m *ItemManager) CanGain(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("item manager check gain item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	// todo bag max item

	return nil
}

func (m *ItemManager) GainLoot(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("item manager gain item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	return m.AddItemByTypeID(typeMisc, num)
}

func (m *ItemManager) LoadFromDB() {
	l := LoadAll(m.ds, m.Owner.GetID(), m.TableName())
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
	}
}

func (m *ItemManager) newEntryItem(entry *define.ItemEntry) Item {
	if entry == nil {
		logger.Error("newEntryItem with nil ItemEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Item)
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
	item.SetNum(i.GetNum())
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

func (m *ItemManager) save(i Item) {
	go func() {
		filter := bson.D{{"_id", i.GetID()}}
		update := bson.D{{"$set", i}}
		op := options.Update().SetUpsert(true)
		m.coll.UpdateOne(context.Background(), filter, update, op)
	}()
}

func (m *ItemManager) delete(i Item) {
	go func() {
		m.coll.DeleteOne(context.Background(), bson.D{{"_id", i.GetID()}})
	}()
}

func (m *ItemManager) AddItemByTypeID(typeID int32, num int32) error {
	if num <= 0 {
		return nil
	}

	incNum := num
	itemEntry := global.GetItemEntry(typeID)

	// add to existing item stack first
	for _, v := range m.mapItem {
		if incNum <= 0 {
			break
		}

		if v.Entry().ID == typeID && v.GetNum() < v.Entry().MaxStack {
			add := incNum
			if incNum > v.Entry().MaxStack-v.GetNum() {
				add = v.Entry().MaxStack - v.GetNum()
			}

			v.SetNum(v.GetNum() + add)
			m.save(v)
			incNum -= add
		}
	}

	// new item to add
	for {
		if incNum <= 0 {
			break
		}

		item := m.newEntryItem(itemEntry)
		if item == nil {
			return fmt.Errorf("new item failed when AddItem:%d", typeID)
		}

		add := incNum
		if incNum > itemEntry.MaxStack {
			add = itemEntry.MaxStack
		}

		item.SetNum(add)
		m.save(item)
		incNum -= add
	}

	return nil
}

func (m *ItemManager) DelItem(id int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(-1)
	delete(m.mapEquipedList, id)
	delete(m.mapItem, id)

	m.delete(i)
}

func (m *ItemManager) CostItemByTypeID(typeID int32, num int32) error {
	if num < 0 {
		return fmt.Errorf("dec item error, invalid number:%d", num)
	}

	decNum := num
	for _, v := range m.mapItem {
		if decNum <= 0 {
			break
		}

		if v.Entry().ID == typeID {
			if v.GetNum() > num {
				v.SetNum(v.GetNum() - num)
				decNum -= num
				m.save(v)
				break
			} else {
				decNum -= v.GetNum()
				m.DelItem(v.GetID())
				continue
			}
		}
	}

	if decNum > 0 {
		logger.WithFields(logger.Fields{
			"need_dec":   num,
			"actual_dec": num - decNum,
		})
	}

	return nil
}

func (m *ItemManager) SetItemEquiped(id int64, objID int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(objID)
	m.mapEquipedList[id] = objID
	m.save(i)
}

func (m *ItemManager) SetItemUnEquiped(id int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(-1)
	delete(m.mapEquipedList, id)
	m.save(i)
}
