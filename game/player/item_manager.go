package player

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// item effect mapping function
type effectFunc func(item.Item) error

type ItemManager struct {
	itemEffectMapping map[int32]effectFunc // item effect mapping function

	owner   *Player
	mapItem map[int64]item.Item

	ds   *db.Datastore
	coll *mongo.Collection
	sync.RWMutex
}

func NewItemManager(owner *Player, ds *db.Datastore) *ItemManager {
	m := &ItemManager{
		itemEffectMapping: make(map[int32]effectFunc, 0),
		owner:             owner,
		ds:                ds,
		mapItem:           make(map[int64]item.Item, 0),
	}

	m.coll = ds.Database().Collection(m.TableName())
	m.initEffectMapping()

	return m
}

// 无效果
func (m *ItemManager) itemEffectNull(i item.Item) error {
	return nil
}

// 掉落
func (m *ItemManager) itemEffectLoot(i item.Item) error {
	for _, v := range i.Entry().EffectValue {
		if err := m.owner.CostLootManager().CanGain(v); err != nil {
			return err
		}
	}

	for _, v := range i.Entry().EffectValue {
		if err := m.owner.CostLootManager().GainLoot(v); err != nil {
			logger.WithFields(logger.Fields{
				"loot_id":      v,
				"item_type_id": i.GetTypeID(),
			}).Warn("itemEffectLoot failed")
		}
	}

	return nil
}

// 御魂鉴定
func (m *ItemManager) itemEffectRuneDefine(i item.Item) error {
	// todo
	return nil
}

func (m *ItemManager) initEffectMapping() {
	m.itemEffectMapping[define.Item_Effect_Null] = m.itemEffectNull
	m.itemEffectMapping[define.Item_Effect_Loot] = m.itemEffectLoot
	m.itemEffectMapping[define.Item_Effect_RuneDefine] = m.itemEffectRuneDefine
}

func (m *ItemManager) save(i item.Item) {
	go func() {
		filter := bson.D{{"_id", i.GetID()}}
		update := bson.D{{"$set", i}}
		op := options.Update().SetUpsert(true)
		m.coll.UpdateOne(context.Background(), filter, update, op)
	}()
}

func (m *ItemManager) delete(id int64) {
	go func() {
		m.coll.DeleteOne(context.Background(), bson.D{{"_id", id}})
	}()
}

func (m *ItemManager) createItem(typeID int32, num int32) item.Item {
	itemEntry := global.GetItemEntry(typeID)
	i := m.createEntryItem(itemEntry)
	if i == nil {
		logger.Warning("new item failed when AddItem:", typeID)
		return nil
	}

	add := num
	if num > itemEntry.MaxStack {
		add = itemEntry.MaxStack
	}

	i.SetNum(add)
	m.save(i)

	return i
}

func (m *ItemManager) delItem(id int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(-1)
	delete(m.mapItem, id)
	m.delete(id)
}

func (m *ItemManager) modifyNum(i item.Item, add int32) {
	i.SetNum(i.GetNum() + add)
	m.save(i)
}

func (m *ItemManager) createEntryItem(entry *define.ItemEntry) item.Item {
	if entry == nil {
		logger.Error("createEntryItem with nil ItemEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Item)
	if err != nil {
		logger.Error(err)
		return nil
	}

	i := item.NewItem(id)
	i.SetOwnerID(m.owner.GetID())
	i.SetTypeID(entry.ID)
	i.SetEntry(entry)

	if entry.EquipEnchantID != -1 {
		i.SetEquipEnchantEntry(global.GetEquipEnchantEntry(entry.EquipEnchantID))

		attManager := att.NewAttManager(i.EquipEnchantEntry().AttID)
		i.SetAttManager(attManager)
	}

	m.mapItem[i.GetID()] = i

	return i
}

func (m *ItemManager) createDBItem(i item.Item) item.Item {
	newItem := item.NewItem(i.GetID())
	newItem.SetOwnerID(i.GetOwnerID())
	newItem.SetTypeID(i.GetTypeID())
	newItem.SetNum(i.GetNum())
	newItem.SetEquipObj(i.GetEquipObj())

	entry := global.GetItemEntry(i.GetTypeID())
	newItem.SetEntry(entry)

	if entry.EquipEnchantID != -1 {
		newItem.SetEquipEnchantEntry(global.GetEquipEnchantEntry(entry.EquipEnchantID))

		attManager := att.NewAttManager(i.EquipEnchantEntry().AttID)
		i.SetAttManager(attManager)
	}

	m.mapItem[newItem.GetID()] = newItem

	return i
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
		if v.GetTypeID() == typeMisc && v.GetEquipObj() != -1 {
			fixNum += v.GetNum()
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
	l := item.LoadAll(m.ds, m.owner.GetID(), m.TableName())
	sliceItem := make([]item.Item, 0)

	listItem := reflect.ValueOf(l)
	if listItem.Kind() != reflect.Slice {
		logger.Error("load item returns non-slice type")
		return
	}

	for n := 0; n < listItem.Len(); n++ {
		p := listItem.Index(n)
		sliceItem = append(sliceItem, p.Interface().(item.Item))
	}

	for _, v := range sliceItem {
		m.createDBItem(v)
	}
}

func (m *ItemManager) GetItem(id int64) item.Item {
	return m.mapItem[id]
}

func (m *ItemManager) GetItemNums() int {
	return len(m.mapItem)
}

func (m *ItemManager) GetItemList() []item.Item {
	list := make([]item.Item, 0)

	for _, v := range m.mapItem {
		list = append(list, v)
	}

	return list
}

func (m *ItemManager) AddItemByTypeID(typeID int32, num int32) error {
	if num <= 0 {
		return nil
	}

	incNum := num

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

			m.modifyNum(v, add)
			m.SendItemUpdate(v)
			incNum -= add
		}
	}

	// new item to add
	for {
		if incNum <= 0 {
			break
		}

		i := m.createItem(typeID, incNum)
		if i == nil {
			break
		}

		incNum -= i.GetNum()
	}

	return nil
}

func (m *ItemManager) DeleteItem(id int64) error {
	if i := m.GetItem(id); i == nil {
		return fmt.Errorf("cannot find item<%d> while DeleteItem", id)
	}

	m.delItem(id)
	m.SendItemDelete(id)

	return nil
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
				m.modifyNum(v, -num)
				m.SendItemUpdate(v)
				decNum -= num
				break
			} else {
				decNum -= v.GetNum()
				delID := v.GetID()
				m.delItem(delID)
				m.SendItemDelete(delID)
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

func (m *ItemManager) CostItemByID(id int64, num int32) error {
	if num < 0 {
		return fmt.Errorf("dec item error, invalid number:%d", num)
	}

	i := m.GetItem(id)
	if i == nil {
		return fmt.Errorf("cannot find item by id:%d", id)
	}

	if i.GetNum() < num {
		return fmt.Errorf("item:%d num:%d not enough, should cost %d", id, i.GetNum(), num)
	}

	// cost
	if i.GetNum() == num {
		m.delItem(id)
		m.SendItemDelete(id)
	} else {
		m.modifyNum(i, -num)
		m.SendItemUpdate(i)
	}

	return nil
}

func (m *ItemManager) SetItemEquiped(id int64, objID int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(objID)
	m.save(i)
}

func (m *ItemManager) SetItemUnEquiped(id int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(-1)
	m.save(i)
}

func (m *ItemManager) UseItem(id int64) error {
	i := m.GetItem(id)
	if i == nil {
		return fmt.Errorf("cannot find item:%d", id)
	}

	if i.Entry().EffectType == define.Item_Effect_Null {
		return fmt.Errorf("item effect null:%d", id)
	}

	// do effect
	if err := m.itemEffectMapping[i.Entry().EffectType](i); err != nil {
		return err
	}

	return m.CostItemByID(id, 1)
}

func (m *ItemManager) SendItemDelete(id int64) {
	msg := &pbGame.MS_DelItem{
		ItemId: id,
	}

	m.owner.SendProtoMessage(msg)
}

func (m *ItemManager) SendItemUpdate(i item.Item) {
	msg := &pbGame.MS_ItemUpdate{
		Item: &pbGame.Item{
			Id:     i.GetID(),
			TypeId: i.GetTypeID(),
			Num:    i.GetNum(),
		},
	}

	m.owner.SendProtoMessage(msg)
}
