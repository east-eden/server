package player

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/item"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/store/db"
	"github.com/yokaiio/yokai_server/utils"
)

// item effect mapping function
type effectFunc func(item.Item) error

type ItemManager struct {
	itemEffectMapping map[int32]effectFunc // item effect mapping function

	owner   *Player
	mapItem map[int64]item.Item

	sync.RWMutex
}

func NewItemManager(owner *Player) *ItemManager {
	m := &ItemManager{
		itemEffectMapping: make(map[int32]effectFunc, 0),
		owner:             owner,
		mapItem:           make(map[int64]item.Item, 0),
	}

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
				"item_type_id": i.GetOptions().TypeId,
			}).Warn("itemEffectLoot failed")
		}
	}

	return nil
}

// 御魂鉴定
func (m *ItemManager) itemEffectRuneDefine(i item.Item) error {
	typeId := rand.Int31n(define.Rune_PositionEnd) + 1
	if err := m.owner.RuneManager().AddRuneByTypeID(typeId); err != nil {
		return err
	}

	return nil
}

func (m *ItemManager) initEffectMapping() {
	m.itemEffectMapping[define.Item_Effect_Null] = m.itemEffectNull
	m.itemEffectMapping[define.Item_Effect_Loot] = m.itemEffectLoot
	m.itemEffectMapping[define.Item_Effect_RuneDefine] = m.itemEffectRuneDefine
}

func (m *ItemManager) createItem(typeID int32, num int32) item.Item {
	itemEntry := entries.GetItemEntry(typeID)
	i := m.createEntryItem(itemEntry)
	if i == nil {
		logger.Warning("new item failed when AddItem:", typeID)
		return nil
	}

	add := num
	if num > itemEntry.MaxStack {
		add = itemEntry.MaxStack
	}

	i.GetOptions().Num = add
	store.GetStore().SaveObject(store.StoreType_Item, i)

	return i
}

func (m *ItemManager) delItem(id int64) {
	i, ok := m.mapItem[id]
	if !ok {
		return
	}

	i.SetEquipObj(-1)
	delete(m.mapItem, id)
	store.GetStore().DeleteObjectFromCacheAndDB(store.StoreType_Item, i)
	item.ReleasePoolItem(i)
}

func (m *ItemManager) modifyNum(i item.Item, add int32) {
	i.GetOptions().Num += add
	store.GetStore().SaveObject(store.StoreType_Item, i)
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

	i := item.NewItem(
		item.Id(id),
		item.OwnerId(m.owner.GetID()),
		item.TypeId(entry.ID),
		item.Entry(entry),
	)

	if entry.EquipEnchantID != -1 {
		i.GetOptions().EquipEnchantEntry = entries.GetEquipEnchantEntry(entry.EquipEnchantID)
		i.GetAttManager().SetBaseAttId(i.EquipEnchantEntry().AttID)
	}

	m.mapItem[i.GetOptions().Id] = i

	return i
}

func (m *ItemManager) initLoadedItem(i item.Item) error {
	entry := entries.GetItemEntry(i.GetOptions().TypeId)

	if entry == nil {
		return fmt.Errorf("item<%d> entry invalid", i.GetOptions().TypeId)
	}

	i.GetOptions().Entry = entry

	if entry.EquipEnchantID != -1 {
		i.GetOptions().EquipEnchantEntry = entries.GetEquipEnchantEntry(entry.EquipEnchantID)
		i.GetAttManager().SetBaseAttId(i.GetOptions().EquipEnchantEntry.AttID)
	}

	m.mapItem[i.GetOptions().Id] = i
	return nil
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
		if v.GetOptions().TypeId == typeMisc && v.GetEquipObj() == -1 {
			fixNum += v.GetOptions().Num
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

func (m *ItemManager) LoadAll() error {
	itemList, err := store.GetStore().LoadArray(store.StoreType_Item, "owner_id", m.owner.GetID(), item.GetItemPool())
	if errors.Is(err, db.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("ItemManager LoadAll: %w", err)
	}

	for _, i := range itemList {
		err := m.initLoadedItem(i.(item.Item))
		if err != nil {
			return fmt.Errorf("ItemManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *ItemManager) Save(id int64) {
	if i := m.GetItem(id); i != nil {
		store.GetStore().SaveObject(store.StoreType_Item, i)
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

		if v.Entry().ID == typeID && v.GetOptions().Num < v.Entry().MaxStack {
			add := incNum
			if incNum > v.Entry().MaxStack-v.GetOptions().Num {
				add = v.Entry().MaxStack - v.GetOptions().Num
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

		m.SendItemAdd(i)
		incNum -= i.GetOptions().Num
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

		if v.Entry().ID == typeID && v.GetEquipObj() == -1 {
			if v.GetOptions().Num > num {
				m.modifyNum(v, -num)
				m.SendItemUpdate(v)
				decNum -= num
				break
			} else {
				decNum -= v.GetOptions().Num
				delId := v.GetOptions().Id
				m.delItem(delId)
				m.SendItemDelete(delId)
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

	if i.GetOptions().Num < num {
		return fmt.Errorf("item:%d num:%d not enough, should cost %d", id, i.GetOptions().Num, num)
	}

	// cost
	if i.GetOptions().Num == num {
		m.delItem(id)
		m.SendItemDelete(id)
	} else {
		m.modifyNum(i, -num)
		m.SendItemUpdate(i)
	}

	return nil
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

func (m *ItemManager) SendItemAdd(i item.Item) {
	msg := &pbGame.M2C_ItemAdd{
		Item: &pbGame.Item{
			Id:     i.GetOptions().Id,
			TypeId: i.GetOptions().TypeId,
			Num:    i.GetOptions().Num,
		},
	}

	m.owner.SendProtoMessage(msg)
}

func (m *ItemManager) SendItemDelete(id int64) {
	msg := &pbGame.M2C_DelItem{
		ItemId: id,
	}

	m.owner.SendProtoMessage(msg)
}

func (m *ItemManager) SendItemUpdate(i item.Item) {
	msg := &pbGame.M2C_ItemUpdate{
		Item: &pbGame.Item{
			Id:     i.GetOptions().Id,
			TypeId: i.GetOptions().TypeId,
			Num:    i.GetOptions().Num,
		},
	}

	m.owner.SendProtoMessage(msg)
}
