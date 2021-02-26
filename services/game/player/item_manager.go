package player

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	"bitbucket.org/east-eden/server/internal/container"
	pbGlobal "bitbucket.org/east-eden/server/proto/global"
	"bitbucket.org/east-eden/server/services/game/item"
	"bitbucket.org/east-eden/server/services/game/prom"
	"bitbucket.org/east-eden/server/store"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	log "github.com/rs/zerolog/log"
	"github.com/valyala/bytebufferpool"
)

func MakeItemKey(itemId int64, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	b.B = append(b.B, "item_map.id_"...)
	b.B = append(b.B, strconv.Itoa(int(itemId))...)

	for _, f := range fields {
		b.B = append(b.B, "."...)
		b.B = append(b.B, f...)
	}

	return b.String()
}

// item effect mapping function
type effectFunc func(item.Itemface, *Player, *Player) error

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
	define.BaseCostLooter `bson:"-" json:"-"`

	nextUpdate int64                     `bson:"-" json:"-"`
	owner      *Player                   `bson:"-" json:"-"`
	CA         *container.ContainerArray `bson:"item_map" json:"item_map"` // 背包列表 0:材料与消耗 1:装备 2:晶石
}

func NewItemManager(owner *Player) *ItemManager {
	m := &ItemManager{
		nextUpdate: time.Now().Add(itemUpdateInterval).Unix(),
		owner:      owner,
		CA:         container.New(int(define.Container_End)),
	}

	return m
}

func (m *ItemManager) Destroy() {
	m.CA.Range(func(val interface{}) bool {
		it := val.(item.Itemface)
		item.GetItemPool(it.GetType()).Put(it)
		return true
	})
}

// 无效果
func itemEffectNull(i item.Itemface, owner *Player, target *Player) error {
	return nil
}

// 掉落
func itemEffectLoot(i item.Itemface, owner *Player, target *Player) error {
	for _, v := range i.Ops().Entry.EffectValue {
		if err := owner.CostLootManager().CanGain(v); err != nil {
			return err
		}
	}

	for _, v := range i.Ops().Entry.EffectValue {
		if err := owner.CostLootManager().GainLoot(v); err != nil {
			log.Warn().
				Int32("loot_id", v).
				Int32("item_type_id", i.Ops().TypeId).
				Msg("itemEffectLoot failed")
		}
	}

	return nil
}

// 御魂鉴定
func itemEffectRuneDefine(i item.Itemface, owner *Player, target *Player) error {
	typeId := rand.Int31n(define.Rune_PositionEnd) + 1
	if err := owner.RuneManager().AddRuneByTypeID(typeId); err != nil {
		return err
	}

	return nil
}

func (m *ItemManager) createItem(typeId int32, num int32) item.Itemface {
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

	i.Ops().Num = add
	i.Ops().CreateTime = time.Now().Unix()

	fields := map[string]interface{}{
		MakeItemKey(i.Ops().Id): i,
	}
	err := store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
	utils.ErrPrint(err, "save item failed when createItem", typeId, m.owner.ID)

	// prometheus ops
	prom.OpsCreateItemCounter.Inc()

	return i
}

func (m *ItemManager) delItem(id int64) error {
	v, ok := m.CA.Get(id)
	if !ok {
		return ErrItemNotFound
	}

	it := v.(item.Itemface)
	it.OnDelete()
	m.CA.Del(id)

	fieldsName := []string{MakeItemKey(id)}
	err := store.GetStore().DeleteFields(define.StoreType_Item, m.owner.ID, fieldsName)
	item.GetItemPool(it.GetType()).Put(it)

	return err
}

func (m *ItemManager) modifyNum(i item.Itemface, add int32) error {
	i.Ops().Num += add

	fields := map[string]interface{}{
		MakeItemKey(i.Ops().Id, "num"): i.Ops().Num,
	}
	return store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
}

func (m *ItemManager) createEntryItem(entry *auto.ItemEntry) item.Itemface {
	if entry == nil {
		log.Error().Msg("createEntryItem with nil ItemEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Item)
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	i := item.NewItem(define.ItemType(entry.Type))

	// item initial
	i.InitItem(
		item.Id(id),
		item.OwnerId(m.owner.GetID()),
		item.TypeId(entry.Id),
		item.Entry(entry),
	)

	// equip initial
	if i.GetType() == define.Item_TypeEquip {
		e := i.(*item.Equip)
		equipEnchantEntry, _ := auto.GetEquipEnchantEntry(e.GetTypeID())
		e.InitEquip(
			item.EquipEnchantEntry(equipEnchantEntry),
		)

		if e.GetEquipEnchantEntry() != nil {
			e.GetAttManager().SetBaseAttId(int32(e.GetEquipEnchantEntry().AttId))
		}
	}

	m.CA.Add(
		int(item.GetContainerType(define.ItemType(i.Ops().Entry.Type))),
		i.Ops().Id,
		i,
	)
	return i
}

func (m *ItemManager) initLoadedItem(i item.Itemface) error {
	entry, ok := auto.GetItemEntry(i.Ops().TypeId)
	if !ok {
		return fmt.Errorf("item<%d> entry invalid", i.Ops().TypeId)
	}

	i.Ops().Entry = entry

	// equip initial
	if i.GetType() == define.Item_TypeEquip {
		e := i.(*item.Equip)
		equipEnchantEntry, _ := auto.GetEquipEnchantEntry(e.GetTypeID())
		e.InitEquip(
			item.EquipEnchantEntry(equipEnchantEntry),
		)

		if e.GetEquipEnchantEntry() != nil {
			e.GetAttManager().SetBaseAttId(int32(e.GetEquipEnchantEntry().AttId))
		}
	}

	m.CA.Add(
		int(item.GetContainerType(i.GetType())),
		i.Ops().Id,
		i,
	)
	return nil
}

// interface of cost_loot
func (m *ItemManager) GetCostLootType() int32 {
	return define.CostLoot_Item
}

func (m *ItemManager) CanCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.CanCost(typeMisc, num)
	if err != nil {
		return err
	}

	var fixNum int32
	m.CA.Range(func(val interface{}) bool {
		it := val.(item.Itemface)

		// isn't this item
		if it.Ops().TypeId != typeMisc {
			return true
		}

		// equip should takeoff first, ignore locked equip
		if it.GetType() == define.Item_TypeEquip {
			equip := it.(*item.Equip)
			if equip.GetEquipObj() != -1 || equip.Lock {
				return true
			}
		}

		fixNum += it.Ops().Num
		return true
	})

	if fixNum >= num {
		return nil
	}

	return fmt.Errorf("not enough item<%d>, num<%d>", typeMisc, num)
}

func (m *ItemManager) DoCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.DoCost(typeMisc, num)
	if err != nil {
		return err
	}

	return m.CostItemByTypeId(typeMisc, num)
}

func (m *ItemManager) CanGain(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.CanGain(typeMisc, num)
	if err != nil {
		return err
	}

	if !m.CanAddItem(typeMisc, num) {
		return fmt.Errorf("bag not enough, cannot add item<%d>", typeMisc)
	}

	return nil
}

func (m *ItemManager) GainLoot(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.GainLoot(typeMisc, num)
	if err != nil {
		return err
	}

	return m.AddItemByTypeId(typeMisc, num)
}

func (m *ItemManager) LoadAll() error {
	loadItems := struct {
		ItemMap map[string]interface{} `bson:"item_map" json:"item_map"`
	}{
		ItemMap: make(map[string]interface{}),
	}

	err := store.GetStore().LoadObject(define.StoreType_Item, m.owner.ID, &loadItems)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("ItemManager LoadAll failed: %w", err)
	}

	for _, v := range loadItems.ItemMap {
		value := v.(map[string]interface{})

		// item.type_id在rejson中读取出来为json.Number类型，mongodb中读取出来为int32类型
		var typeId int32
		switch value["type_id"].(type) {
		case json.Number:
			id, _ := value["type_id"].(json.Number).Int64()
			typeId = int32(id)
		case int32:
			typeId = value["type_id"].(int32)
		}

		itemEntry, ok := auto.GetItemEntry(int32(typeId))
		if !ok {
			return fmt.Errorf("ItemManager LoadAll failed: cannot find item entry<%d>", typeId)
		}

		i := item.NewItem(define.ItemType(itemEntry.Type))

		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName: "json",
			Squash:  true,
			Result:  i,
		})
		if err != nil {
			return fmt.Errorf("mapstructure NewDecoder failed when ItemManager LoadAll: %w", err)
		}

		if err = decoder.Decode(value); err != nil {
			return fmt.Errorf("mapstructure Decode failed when ItemManager LoadAll: %w", err)
		}

		if err = m.initLoadedItem(i); err != nil {
			return fmt.Errorf("ItemManager LoadAll failed: %w", err)
		}
	}

	return err
}

func (m *ItemManager) Save(id int64) error {
	i, err := m.GetItem(id)
	if err != nil {
		return fmt.Errorf("ItemManager.Save failed: %w", err)
	}

	fields := map[string]interface{}{
		MakeItemKey(id): i,
	}
	return store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
}

func (m *ItemManager) update() {
	if time.Now().Unix() < m.nextUpdate {
		return
	}

	// 设置下次更新时间
	m.nextUpdate = time.Now().Add(itemUpdateInterval).Unix()

	// 遍历容器删除过期物品
	m.CA.Range(func(val interface{}) bool {
		it := val.(item.Itemface)
		if it.Ops().Entry.TimeLife == -1 {
			return true
		}

		// 时限物品开始计时时间
		var startTime time.Time
		if it.Ops().Entry.TimeStartLifeStamp == 0 {
			startTime = time.Unix(it.Ops().CreateTime, 0)
		} else {
			startTime = time.Unix(int64(it.Ops().Entry.TimeStartLifeStamp), 0)
		}

		endTime := startTime.Add(time.Minute * time.Duration(it.Ops().Entry.TimeLife))
		if time.Now().Unix() >= endTime.Unix() {
			log.Info().Int32("type_id", it.Ops().Entry.Id).Msg("item time countdown and deleted")
			err := m.expireItem(it)
			utils.ErrPrint(err, "expireItem failed when update", it.Ops().Id, m.owner.ID)
		}

		return true
	})
}

func (m *ItemManager) GetItem(id int64) (item.Itemface, error) {
	val, ok := m.CA.Get(id)
	if ok {
		return val.(item.Itemface), nil
	} else {
		return nil, ErrItemNotFound
	}
}

func (m *ItemManager) GetItemNums(idx int) int {
	return m.CA.Size(idx)
}

func (m *ItemManager) GetItemList() []item.Itemface {
	list := make([]item.Itemface, 0, 50)

	m.CA.Range(func(val interface{}) bool {
		list = append(list, val.(item.Itemface))
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
		m.CA.RangeByIdx(int(idx), func(val interface{}) bool {
			it := val.(item.Itemface)
			if it.Ops().TypeId == typeId {
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
	if m.CA.Size(int(idx))+int(num) >= globalConfig.GetItemContainerSize(idx) {
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
	m.CA.RangeByIdx(
		int(item.GetContainerType(define.ItemType(entry.Type))),
		func(val interface{}) bool {
			it := val.(item.Itemface)
			if incNum <= 0 {
				return false
			}

			// 静态表中配置为不堆叠
			if entry.MaxStack <= 1 {
				return false
			}

			if it.Ops().Entry.Id == typeId && it.Ops().Num < it.Ops().Entry.MaxStack {
				add := incNum
				if incNum > it.Ops().Entry.MaxStack-it.Ops().Num {
					add = it.Ops().Entry.MaxStack - it.Ops().Num
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

		m.CA.Add(
			int(item.GetContainerType(i.GetType())),
			i.Ops().Id,
			i,
		)
		m.SendItemAdd(i)
		incNum -= i.Ops().Num
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
func (m *ItemManager) expireItem(it item.Itemface) error {
	gainId := it.Ops().Entry.StaleGainId

	// 删除物品
	if err := m.DeleteItem(it.Ops().Id); err != nil {
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

	m.CA.RangeByIdx(
		int(item.GetContainerType(define.ItemType(entry.Type))),
		func(val interface{}) bool {
			it := val.(item.Itemface)
			if decNum <= 0 {
				return false
			}

			// 不是这个物品
			if it.Ops().TypeId != typeId {
				return true
			}

			// 装备还穿在身上
			if it.GetType() == define.Item_TypeEquip && it.(*item.Equip).GetEquipObj() != -1 {
				return true
			}

			if it.Ops().Num > num {
				if errModify := m.modifyNum(it, -num); errModify != nil {
					err = errModify
					utils.ErrPrint(errModify, "modifyNum failed when CostItemByTypeID", typeId, num, m.owner.ID)
					return true
				}

				m.SendItemUpdate(it)
				decNum -= num
				return false
			} else {
				decNum -= it.Ops().Num
				delId := it.Ops().Id
				if errDel := m.delItem(delId); errDel != nil {
					err = errDel
					utils.ErrPrint(err, "delItem failed when CostItemByTypeID", typeId, m.owner.ID)
					return true
				}

				m.SendItemDelete(delId)
				return true
			}
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

	if i.Ops().Num < num {
		return fmt.Errorf("ItemManager.CostItemByID failed: item:%d num:%d not enough, should cost %d", id, i.Ops().Num, num)
	}

	// cost
	if i.Ops().Num == num {
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

func (m *ItemManager) CanUseItem(i item.Itemface) error {
	// 时限物品
	if i.Ops().Entry.TimeLife > 0 {
		// 时限开始物品
		if i.Ops().Entry.TimeStartLifeStamp > 0 {
			startTime := time.Unix(int64(i.Ops().Entry.TimeStartLifeStamp), 0)
			endTime := startTime.Add(time.Duration(i.Ops().Entry.TimeLife) * time.Minute)
			if time.Now().Unix() >= endTime.Unix() || time.Now().Unix() < startTime.Unix() {
				return fmt.Errorf("time life limit")
			}
		} else {
			if time.Now().Unix() >= i.Ops().CreateTime {
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

	if i.Ops().Entry.EffectType == define.Item_Effect_Null {
		return fmt.Errorf("ItemManager.UseItem failed: item<%d> effect null", id)
	}

	// do effect
	if err := itemEffectFuncMapping[i.Ops().Entry.EffectType](i, m.owner, m.owner); err != nil {
		return fmt.Errorf("ItemManager.UseItem failed: %w", err)
	}

	return m.CostItemByID(id, 1)
}

func (m *ItemManager) EquipLevelup(equipId int64) error {
	i, err := m.GetItem(equipId)
	utils.ErrPrint(err, "EquipLevelup failed", equipId, m.owner.ID)

	if i.GetType() != define.Item_TypeEquip {
		return fmt.Errorf("EquipLevelup failed, wrong item<%d> type", equipId)
	}

	equip, ok := i.(*item.Equip)
	if !ok {
		return fmt.Errorf("EquipLevelup failed, cannot assert to equip<%d>", equipId)
	}

	levelUpEntry, ok := auto.GetEquipLevelupEntry(int32(equip.Level) + 1)
	if !ok {
		return fmt.Errorf("EquipLevelup failed, cannot find EquipLevelupEntry<%d>", equip.Level+1)
	}

	if int32(equip.Promote) < levelUpEntry.PromoteLimit {
		return fmt.Errorf("EquipLevelup failed, PromoteLevel<%d> limit", equip.Promote)
	}

	costId := levelUpEntry.CostId[int(equip.Entry().Quality)]
	if err := m.owner.CostLootManager().CanCost(costId); err != nil {
		return fmt.Errorf("EquipLevelup failed, CanCost<%d> error: %w", levelUpEntry.Id, err)
	}

	if err := m.owner.CostLootManager().DoCost(costId); err != nil {
		return fmt.Errorf("EquipLevelup failed, DoCost<%d> error: %w", levelUpEntry.Id, err)
	}

	equip.Level++
	equip.Exp = 0
	equip.GetAttManager().CalcAtt()

	// save
	fields := map[string]interface{}{
		MakeItemKey(equip.GetID(), "level"): equip.Level,
		MakeItemKey(equip.GetID(), "exp"):   equip.Exp,
	}
	err = store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
	utils.ErrPrint(err, "SaveFields failed when EquipLevelup", equip.GetID(), m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)

	return err
}

func (m *ItemManager) SendItemAdd(i item.Itemface) {
	msg := &pbGlobal.S2C_ItemAdd{
		Item: &pbGlobal.Item{
			Id:     i.Ops().Id,
			TypeId: int32(i.Ops().TypeId),
			Num:    int32(i.Ops().Num),
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

func (m *ItemManager) SendItemUpdate(i item.Itemface) {
	msg := &pbGlobal.S2C_ItemUpdate{
		Item: &pbGlobal.Item{
			Id:     i.Ops().Id,
			TypeId: int32(i.Ops().TypeId),
			Num:    int32(i.Ops().Num),
		},
	}

	m.owner.SendProtoMessage(msg)
}

func (m *ItemManager) SendEquipUpdate(e *item.Equip) {
	msg := &pbGlobal.S2C_EquipUpdate{
		EquipId: e.GetID(),
		EquipData: &pbGlobal.EquipData{
			Exp:      e.Exp,
			Level:    int32(e.Level),
			Promote:  int32(e.Promote),
			Lock:     e.Lock,
			EquipObj: e.EquipObj,
		},
	}

	m.owner.SendProtoMessage(msg)
}
