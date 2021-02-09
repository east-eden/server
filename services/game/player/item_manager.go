package player

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/services/game/item"
	"bitbucket.org/east-eden/server/services/game/prom"
	"bitbucket.org/east-eden/server/store"
	"bitbucket.org/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

// item effect mapping function
type effectFunc func(*item.Item, *Player, *Player) error

var itemEffectFuncMapping = map[int32]effectFunc{
	define.Item_Effect_Null:       itemEffectNull,
	define.Item_Effect_Loot:       itemEffectLoot,
	define.Item_Effect_RuneDefine: itemEffectRuneDefine,
}

var (
	itemUpdateInterval = time.Second * 5 // 物品每5秒更新一次
)

type ItemManager struct {
	nextUpdate int64
	owner      *Player
	mapItem    map[int64]*item.Item
}

func NewItemManager(owner *Player) *ItemManager {
	m := &ItemManager{
		nextUpdate: time.Now().Add(itemUpdateInterval).Unix(),
		owner:      owner,
		mapItem:    make(map[int64]*item.Item),
	}

	return m
}

// 无效果
func itemEffectNull(i *item.Item, owner *Player, target *Player) error {
	return nil
}

// 掉落
func itemEffectLoot(i *item.Item, owner *Player, target *Player) error {
	for _, v := range i.Entry().EffectValue {
		if err := owner.CostLootManager().CanGain(v); err != nil {
			return err
		}
	}

	for _, v := range i.Entry().EffectValue {
		if err := owner.CostLootManager().GainLoot(v); err != nil {
			log.Warn().
				Int32("loot_id", v).
				Int32("item_type_id", i.GetOptions().TypeId).
				Msg("itemEffectLoot failed")
		}
	}

	return nil
}

// 御魂鉴定
func itemEffectRuneDefine(i *item.Item, owner *Player, target *Player) error {
	typeId := rand.Int31n(define.Rune_PositionEnd) + 1
	if err := owner.RuneManager().AddRuneByTypeID(typeId); err != nil {
		return err
	}

	return nil
}

func (m *ItemManager) createItem(typeId int32, num int32) *item.Item {
	itemEntry, ok := auto.GetItemEntry(typeId)
	if !ok {
		return nil
	}

	i := m.createEntryItem(itemEntry)
	if i == nil {
		log.Warn().
			Int32("item_type_id", typeId).
			Msg("new item failed when AddItem")
		return nil
	}

	add := num
	if num > itemEntry.MaxStack {
		add = itemEntry.MaxStack
	}

	i.GetOptions().Num = add
	i.GetOptions().CreateTime = time.Now().Unix()
	err := store.GetStore().SaveObject(define.StoreType_Item, i)
	utils.ErrPrint(err, "save item failed when createItem", typeId, m.owner.ID)

	// prometheus ops
	prom.OpsCreateItemCounter.Inc()

	return i
}

func (m *ItemManager) delItem(id int64) error {
	i, ok := m.mapItem[id]
	if !ok {
		return nil
	}

	i.SetEquipObj(-1)
	delete(m.mapItem, id)
	err := store.GetStore().DeleteObject(define.StoreType_Item, i)
	item.ReleasePoolItem(i)

	return err
}

func (m *ItemManager) modifyNum(i *item.Item, add int32) error {
	i.GetOptions().Num += add
	return store.GetStore().SaveObject(define.StoreType_Item, i)
}

func (m *ItemManager) createEntryItem(entry *auto.ItemEntry) *item.Item {
	if entry == nil {
		log.Error().Msg("createEntryItem with nil ItemEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Item)
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	i := item.NewItem(
		item.Id(id),
		item.OwnerId(m.owner.GetID()),
		item.TypeId(entry.Id),
		item.Entry(entry),
	)

	if entry.EquipEnchantId != -1 {
		var ok bool
		if i.GetOptions().EquipEnchantEntry, ok = auto.GetEquipEnchantEntry(entry.EquipEnchantId); ok {
			i.GetAttManager().SetBaseAttId(int32(i.EquipEnchantEntry().AttId))
		}
	}

	m.mapItem[i.GetOptions().Id] = i

	return i
}

func (m *ItemManager) initLoadedItem(i *item.Item) error {
	entry, ok := auto.GetItemEntry(i.GetOptions().TypeId)
	if !ok {
		return fmt.Errorf("item<%d> entry invalid", i.GetOptions().TypeId)
	}

	i.GetOptions().Entry = entry

	if entry.EquipEnchantId != -1 {
		i.GetOptions().EquipEnchantEntry, _ = auto.GetEquipEnchantEntry(entry.EquipEnchantId)
		i.GetAttManager().SetBaseAttId(int32(i.GetOptions().EquipEnchantEntry.AttId))
	}

	m.mapItem[i.GetOptions().Id] = i
	return nil
}

// interface of cost_loot
func (m *ItemManager) GetCostLootType() int32 {
	return define.CostLoot_Item
}

func (m *ItemManager) CanCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("item manager check item<%d> cost failed, wrong number<%d>", typeMisc, num)
	}

	var fixNum int32
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
	itemList, err := store.GetStore().LoadArray(define.StoreType_Item, m.owner.GetID(), item.GetItemPool())
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("ItemManager LoadAll: %w", err)
	}

	for _, i := range itemList {
		err := m.initLoadedItem(i.(*item.Item))
		if err != nil {
			return fmt.Errorf("ItemManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *ItemManager) Save(id int64) error {
	i, err := m.GetItem(id)
	if err != nil {
		return fmt.Errorf("ItemManager.Save failed: %w", err)
	}

	return store.GetStore().SaveObject(define.StoreType_Item, i)
}

func (m *ItemManager) update() {
	if time.Now().Unix() < m.nextUpdate {
		return
	}

	// 设置下次更新时间
	m.nextUpdate = time.Now().Add(itemUpdateInterval).Unix()

	for k, v := range m.mapItem {
		if v.Entry().TimeLife == -1 {
			continue
		}

		// 时限物品开始计时时间
		var startTime time.Time
		if v.Entry().TimeStartLifeStamp == 0 {
			startTime = time.Unix(v.CreateTime, 0)
		} else {
			startTime = time.Unix(int64(v.Entry().TimeStartLifeStamp), 0)
		}

		endTime := startTime.Add(time.Minute * time.Duration(v.Entry().TimeLife))
		if time.Now().Unix() >= endTime.Unix() {
			log.Info().Int32("type_id", v.Entry().Id).Msg("item time countdown and deleted")
			err := m.DeleteItem(k)
			utils.ErrPrint(err, "DeleteItem failed when update", k, m.owner.ID)
		}
	}
}

func (m *ItemManager) GetItem(id int64) (*item.Item, error) {
	if i, ok := m.mapItem[id]; ok {
		return i, nil
	} else {
		return nil, fmt.Errorf("invalid item id<%d>", id)
	}
}

func (m *ItemManager) GetItemNums() int {
	return len(m.mapItem)
}

func (m *ItemManager) GetItemList() []*item.Item {
	list := make([]*item.Item, 0)

	for _, v := range m.mapItem {
		list = append(list, v)
	}

	return list
}

func (m *ItemManager) AddItemByTypeID(typeId int32, num int32) error {
	if num <= 0 {
		return nil
	}

	incNum := num

	// add to existing item stack first
	var err error
	for _, v := range m.mapItem {
		if incNum <= 0 {
			break
		}

		if v.Entry().Id == typeId && v.GetOptions().Num < v.Entry().MaxStack {
			add := incNum
			if incNum > v.Entry().MaxStack-v.GetOptions().Num {
				add = v.Entry().MaxStack - v.GetOptions().Num
			}

			if errModify := m.modifyNum(v, add); errModify != nil {
				err = errModify
				utils.ErrPrint(errModify, "modifyNum failed when AddItemByTypeID", typeId, num, m.owner.ID)
				continue
			}

			m.SendItemUpdate(v)
			incNum -= add
		}
	}

	// new item to add
	for {
		if incNum <= 0 {
			break
		}

		i := m.createItem(typeId, incNum)
		if i == nil {
			break
		}

		m.SendItemAdd(i)
		incNum -= i.GetOptions().Num
	}

	return err
}

func (m *ItemManager) DeleteItem(id int64) error {
	_, err := m.GetItem(id)
	if err != nil {
		return fmt.Errorf("ItemManager.DeleteItem failed: %w", err)
	}

	err = m.delItem(id)
	m.SendItemDelete(id)
	return err
}

func (m *ItemManager) CostItemByTypeID(typeID int32, num int32) error {
	if num < 0 {
		return fmt.Errorf("dec item error, invalid number:%d", num)
	}

	var err error
	decNum := num
	for _, v := range m.mapItem {
		if decNum <= 0 {
			break
		}

		if v.Entry().Id == typeID && v.GetEquipObj() == -1 {
			if v.GetOptions().Num > num {
				if errModify := m.modifyNum(v, -num); errModify != nil {
					err = errModify
					utils.ErrPrint(errModify, "modifyNum failed when CostItemByTypeID", typeID, num, m.owner.ID)
					continue
				}

				m.SendItemUpdate(v)
				decNum -= num
				break
			} else {
				decNum -= v.GetOptions().Num
				delId := v.GetOptions().Id
				if errDel := m.delItem(delId); errDel != nil {
					err = errDel
					utils.ErrPrint(err, "delItem failed when CostItemByTypeID", typeID, m.owner.ID)
					continue
				}

				m.SendItemDelete(delId)
				continue
			}
		}
	}

	if decNum > 0 {
		log.Warn().
			Int32("need_dec", num).
			Int32("actual_dec", num-decNum).
			Msg("cost item num not enough")
	}

	return err
}

func (m *ItemManager) CostItemByID(id int64, num int32) error {
	if num < 0 {
		return fmt.Errorf("dec item error, invalid number:%d", num)
	}

	i, err := m.GetItem(id)
	if err != nil {
		return fmt.Errorf("ItemManager.CostItemByID failed: %w", err)
	}

	if i.GetOptions().Num < num {
		return fmt.Errorf("ItemManager.CostItemByID failed: item:%d num:%d not enough, should cost %d", id, i.GetOptions().Num, num)
	}

	// cost
	if i.GetOptions().Num == num {
		err = m.delItem(id)
		if pass := utils.ErrCheck(err, "delItem failed when CostItemByID", id, num, m.owner.ID); !pass {
			return err
		}

		m.SendItemDelete(id)
	} else {
		err = m.modifyNum(i, -num)
		if pass := utils.ErrCheck(err, "modifyNum failed when CostItemByID", id, num, m.owner.ID); !pass {
			return err
		}

		m.SendItemUpdate(i)
	}

	return nil
}

func (m *ItemManager) UseItem(id int64) error {
	i, err := m.GetItem(id)
	if err != nil {
		return fmt.Errorf("ItemManager.UseItem failed: %w", err)
	}

	if i.Entry().EffectType == define.Item_Effect_Null {
		return fmt.Errorf("ItemManager.UseItem failed: item<%d> effect null", id)
	}

	// do effect
	if err := itemEffectFuncMapping[i.Entry().EffectType](i, m.owner, m.owner); err != nil {
		return fmt.Errorf("ItemManager.UseItem failed: %w", err)
	}

	return m.CostItemByID(id, 1)
}

func (m *ItemManager) SendItemAdd(i *item.Item) {
	msg := &pbGlobal.S2C_ItemAdd{
		Item: &pbGlobal.Item{
			Id:     i.GetOptions().Id,
			TypeId: int32(i.GetOptions().TypeId),
			Num:    int32(i.GetOptions().Num),
		},
	}

	m.owner.SendProtoMessage(msg)
}

func (m *ItemManager) SendItemDelete(id int64) {
	msg := &pbGlobal.S2C_DelItem{
		ItemId: id,
	}

	m.owner.SendProtoMessage(msg)
}

func (m *ItemManager) SendItemUpdate(i *item.Item) {
	msg := &pbGlobal.S2C_ItemUpdate{
		Item: &pbGlobal.Item{
			Id:     i.GetOptions().Id,
			TypeId: int32(i.GetOptions().TypeId),
			Num:    int32(i.GetOptions().Num),
		},
	}

	m.owner.SendProtoMessage(msg)
}
