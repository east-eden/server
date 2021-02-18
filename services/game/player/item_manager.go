package player

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	"bitbucket.org/east-eden/server/internal/container"
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/services/game/item"
	"bitbucket.org/east-eden/server/services/game/prom"
	"bitbucket.org/east-eden/server/store"
	"bitbucket.org/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

// item effect mapping function
type effectFunc func(*item.Item, *Player, *Player) error

// 物品使用效果
var itemEffectFuncMapping = map[int32]effectFunc{
	define.Item_Effect_Null:       itemEffectNull,
	define.Item_Effect_Loot:       itemEffectLoot,
	define.Item_Effect_RuneDefine: itemEffectRuneDefine,
}

var (
	itemUpdateInterval = time.Second * 5 // 物品每5秒更新一次
	ErrItemNotFound    = errors.New("item not found")
)

// 物品管理
type ItemManager struct {
	nextUpdate int64
	owner      *Player
	ca         *container.ContainerArray // 背包列表 0:材料与消耗 1:装备 2:晶石
}

func NewItemManager(owner *Player) *ItemManager {
	m := &ItemManager{
		nextUpdate: time.Now().Add(itemUpdateInterval).Unix(),
		owner:      owner,
		ca:         container.New(int(define.Container_End)),
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
	err := store.GetStore().SaveObject(define.StoreType_Item, i.GetObjID(), i)
	utils.ErrPrint(err, "save item failed when createItem", typeId, m.owner.ID)

	// prometheus ops
	prom.OpsCreateItemCounter.Inc()

	return i
}

func (m *ItemManager) delItem(id int64) error {
	v, ok := m.ca.Get(id)
	if !ok {
		return ErrItemNotFound
	}

	it := v.(*item.Item)
	it.SetEquipObj(-1)
	m.ca.Del(id)
	err := store.GetStore().DeleteObject(define.StoreType_Item, it)
	item.ReleasePoolItem(it)

	return err
}

func (m *ItemManager) modifyNum(i *item.Item, add int32) error {
	i.GetOptions().Num += add
	return store.GetStore().SaveObject(define.StoreType_Item, i.GetObjID(), i)
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

	m.ca.Add(
		int(item.GetContainerType(define.ItemType(i.Entry().Type))),
		i.Id,
		i,
	)
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

	m.ca.Add(
		int(item.GetContainerType(define.ItemType(i.Entry().Type))),
		i.Id,
		i,
	)
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
	m.ca.Range(func(val interface{}) bool {
		it := val.(*item.Item)
		if it.TypeId == typeMisc && it.GetEquipObj() == -1 {
			fixNum += it.Num
		}

		return true
	})

	if fixNum >= num {
		return nil
	}

	return fmt.Errorf("not enough item<%d>, num<%d>", typeMisc, num)
}

func (m *ItemManager) DoCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("item manager cost item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	return m.CostItemByTypeId(typeMisc, num)
}

func (m *ItemManager) CanGain(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("item manager check gain item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	// todo bag max item
	if !m.CanAddItem(typeMisc, num) {
		return fmt.Errorf("bag not enough, cannot add item<%d>", typeMisc)
	}

	return nil
}

func (m *ItemManager) GainLoot(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("item manager gain item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	return m.AddItemByTypeId(typeMisc, num)
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

	return store.GetStore().SaveObject(define.StoreType_Item, i.GetObjID(), i)
}

func (m *ItemManager) update() {
	if time.Now().Unix() < m.nextUpdate {
		return
	}

	// 设置下次更新时间
	m.nextUpdate = time.Now().Add(itemUpdateInterval).Unix()

	// 遍历容器删除过期物品
	m.ca.Range(func(val interface{}) bool {
		it := val.(*item.Item)
		if it.Entry().TimeLife == -1 {
			return true
		}

		// 时限物品开始计时时间
		var startTime time.Time
		if it.Entry().TimeStartLifeStamp == 0 {
			startTime = time.Unix(it.CreateTime, 0)
		} else {
			startTime = time.Unix(int64(it.Entry().TimeStartLifeStamp), 0)
		}

		endTime := startTime.Add(time.Minute * time.Duration(it.Entry().TimeLife))
		if time.Now().Unix() >= endTime.Unix() {
			log.Info().Int32("type_id", it.Entry().Id).Msg("item time countdown and deleted")
			err := m.expireItem(it)
			utils.ErrPrint(err, "expireItem failed when update", it.Id, m.owner.ID)
		}

		return true
	})
}

func (m *ItemManager) GetItem(id int64) (*item.Item, error) {
	val, ok := m.ca.Get(id)
	if ok {
		return val.(*item.Item), nil
	} else {
		return nil, ErrItemNotFound
	}
}

func (m *ItemManager) GetItemNums(idx int) int {
	return m.ca.Size(idx)
}

func (m *ItemManager) GetItemList() []*item.Item {
	list := make([]*item.Item, 0, 50)

	m.ca.Range(func(val interface{}) bool {
		list = append(list, val.(*item.Item))
		return true
	})

	return list
}

func (m *ItemManager) CanAddItem(typeId, num int32) bool {
	itemEntry, ok := auto.GetItemEntry(typeId)
	if !ok {
		return false
	}

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return false
	}

	var canAdd bool
	idx := item.GetContainerType(define.ItemType(itemEntry.Type))

	// 背包中有相同typeId的物品，并且是可叠加的，一定成功
	if itemEntry.MaxStack > 1 {
		m.ca.RangeByIdx(int(idx), func(val interface{}) bool {
			it := val.(*item.Item)
			if it.TypeId == typeId {
				canAdd = true
				return false
			}

			return true
		})
	}

	if canAdd {
		return true
	}

	// 无可叠加物品，并且背包容量不够
	if m.ca.Size(int(idx))+int(num) >= globalConfig.GetItemContainerSize(idx) {
		return false
	}

	return true
}

func (m *ItemManager) AddItemByTypeId(typeId int32, num int32) error {
	if num <= 0 {
		return nil
	}

	entry, ok := auto.GetItemEntry(typeId)
	if !ok {
		return fmt.Errorf("GetItemEntry<%d> failed", typeId)
	}

	incNum := num

	// add to existing item stack first
	var err error
	m.ca.RangeByIdx(
		int(item.GetContainerType(define.ItemType(entry.Type))),
		func(val interface{}) bool {
			it := val.(*item.Item)
			if incNum <= 0 {
				return false
			}

			if it.Entry().Id == typeId && it.Num < it.Entry().MaxStack {
				add := incNum
				if incNum > it.Entry().MaxStack-it.Num {
					add = it.Entry().MaxStack - it.Num
				}

				if errModify := m.modifyNum(it, add); errModify != nil {
					err = errModify
					utils.ErrPrint(errModify, "modifyNum failed when AddItemByTypeID", typeId, num, m.owner.ID)
					return true
				}

				m.SendItemUpdate(it)
				incNum -= add
			}

			return true
		},
	)

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

// 物品过期处理
func (m *ItemManager) expireItem(it *item.Item) error {
	gainId := it.Entry().StaleGainId

	// 删除物品
	if err := m.DeleteItem(it.Id); err != nil {
		return err
	}

	if gainId == -1 {
		return nil
	}

	// 添加转换后的物品
	if err := m.owner.costLootManager.CanGain(gainId); err != nil {
		return err
	}

	return m.owner.costLootManager.GainLoot(gainId)
}

func (m *ItemManager) CostItemByTypeId(typeId int32, num int32) error {
	if num < 0 {
		return fmt.Errorf("dec item error, invalid number:%d", num)
	}

	entry, ok := auto.GetItemEntry(typeId)
	if !ok {
		return fmt.Errorf("GetItemEntry<%d> failed", typeId)
	}

	var err error
	decNum := num

	m.ca.RangeByIdx(
		int(item.GetContainerType(define.ItemType(entry.Type))),
		func(val interface{}) bool {
			it := val.(*item.Item)
			if decNum <= 0 {
				return false
			}

			if it.Entry().Id == typeId && it.GetEquipObj() == -1 {
				if it.Num > num {
					if errModify := m.modifyNum(it, -num); errModify != nil {
						err = errModify
						utils.ErrPrint(errModify, "modifyNum failed when CostItemByTypeID", typeId, num, m.owner.ID)
						return true
					}

					m.SendItemUpdate(it)
					decNum -= num
					return false
				} else {
					decNum -= it.Num
					delId := it.Id
					if errDel := m.delItem(delId); errDel != nil {
						err = errDel
						utils.ErrPrint(err, "delItem failed when CostItemByTypeID", typeId, m.owner.ID)
						return true
					}

					m.SendItemDelete(delId)
					return true
				}
			}
			return true
		},
	)

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

func (m *ItemManager) CanUseItem(i *item.Item) error {
	// 时限物品
	if i.Entry().TimeLife > 0 {
		// 时限开始物品
		if i.Entry().TimeStartLifeStamp > 0 {
			startTime := time.Unix(int64(i.Entry().TimeStartLifeStamp), 0)
			endTime := startTime.Add(time.Duration(i.Entry().TimeLife) * time.Minute)
			if time.Now().Unix() >= endTime.Unix() || time.Now().Unix() < startTime.Unix() {
				return fmt.Errorf("time life limit")
			}
		} else {
			if time.Now().Unix() >= i.CreateTime {
				return fmt.Errorf("time life limit")
			}
		}
	}

	return nil
}

func (m *ItemManager) UseItem(id int64) error {
	i, err := m.GetItem(id)
	if err != nil {
		return fmt.Errorf("ItemManager.UseItem failed: %w", err)
	}

	if err := m.CanUseItem(i); err != nil {
		return err
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
