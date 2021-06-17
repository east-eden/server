package player

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/internal/container"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/item"
	"github.com/east-eden/server/services/game/prom"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

// func getItemKeyName(it item.Itemface) string {
// 	var keyName string
// 	switch e := it.GetType(); e {
// 	case define.Item_TypePresent:
// 		fallthrough
// 	case define.Item_TypeItem:
// 		keyName = "item_list"
// 	case define.Item_TypeEquip:
// 		keyName = "equip_list"
// 	case define.Item_TypeCrystal:
// 		keyName = "crystal_list"
// 	default:
// 		log.Error().Caller().Int32("item_type", e).Msg("makeItemKey failed, invalid item type")
// 		return ""
// 	}

// 	return keyName
// }

// func makeItemKey(it item.Itemface, fields ...string) string {
// 	b := bytebufferpool.Get()
// 	defer bytebufferpool.Put(b)

// 	keyName := getItemKeyName(it)

// 	_, _ = b.WriteString(keyName)
// 	_, _ = b.WriteString(".$")

// 	for _, f := range fields {
// 		_, _ = b.WriteString(".")
// 		_, _ = b.WriteString(f)
// 	}

// 	return b.String()
// }

// item effect mapping function
type effectFunc func(item.Itemface, *Player, *Player) error

// 物品使用效果
var itemEffectFuncMapping = map[int32]effectFunc{
	define.Item_Effect_Null: itemEffectNull,
	define.Item_Effect_Loot: itemEffectLoot,
}

var (
	itemUpdateInterval       = time.Second * 5 // 物品每5秒更新一次
	ErrItemNotFound          = errors.New("item not found")
	ErrItemNotEnough         = errors.New("item not enough")
	ErrItemInvalidType       = errors.New("invalid item type")
	ErrEquipPromoteTimesFull = errors.New("promote times full")
	ErrEquipStarTimesFull    = errors.New("star times full")
)

// 物品管理
type ItemManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	nextUpdate int64                     `bson:"-" json:"-"`
	owner      *Player                   `bson:"-" json:"-"`
	ca         *container.ContainerArray `bson:"-" json:"-"` // 背包列表 0:材料与消耗 1:装备 2:晶石
}

func NewItemManager(owner *Player) *ItemManager {
	m := &ItemManager{
		nextUpdate: time.Now().Add(itemUpdateInterval).Unix(),
		owner:      owner,
		ca:         container.New(int(define.Container_End)),
	}

	return m
}

func (m *ItemManager) Destroy() {
	m.ca.Range(func(val interface{}) bool {
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
	lootId := i.Opts().ItemEntry.EffectValue
	if err := owner.CostLootManager().CanGain(lootId); err != nil {
		return err
	}

	err := owner.CostLootManager().GainLoot(lootId)
	utils.ErrPrint(err, "itemEffectLoot failed", lootId, i.Opts().TypeId)
	return err
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

	i.Opts().Num = add
	i.Opts().CreateTime = time.Now().Unix()

	return i
}

func (m *ItemManager) delItem(id int64) error {
	v, ok := m.ca.Get(id)
	if !ok {
		return ErrItemNotFound
	}

	it := v.(item.Itemface)
	m.ca.Del(id)

	err := store.GetStore().DeleteOne(context.Background(), define.StoreType_Item, id)
	utils.ErrPrint(err, "DeleteOne failed when ItemManager.delItem", m.owner.ID, id)

	item.GetItemPool(it.GetType()).Put(it)

	return err
}

func (m *ItemManager) modifyNum(i item.Itemface, add int32) error {
	i.Opts().Num += add

	fields := map[string]interface{}{
		"num": i.Opts().Num,
	}
	return store.GetStore().UpdateFields(context.Background(), define.StoreType_Item, i.Opts().Id, fields)
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

	i := item.NewItem(entry.Type)

	// item initial
	i.InitItem(
		item.Id(id),
		item.OwnerId(m.owner.GetId()),
		item.TypeId(entry.Id),
		item.ItemEntry(entry),
	)

	switch i.GetType() {

	// equip initial
	case define.Item_TypeEquip:
		e := i.(*item.Equip)
		equipEnchantEntry, ok := auto.GetEquipEnchantEntry(e.GetTypeID())
		if !ok {
			log.Error().Caller().Int32("type_id", e.GetTypeID()).Msg("createEntryItem failed")
			return nil
		}

		e.InitEquip(
			item.EquipEnchantEntry(equipEnchantEntry),
		)

		if e.EquipEnchantEntry != nil {
			e.GetAttManager().SetBaseAttId(int32(e.EquipEnchantEntry.AttId))
		}

	// crystal initial
	case define.Item_TypeCrystal:
		c := i.(*item.Crystal)
		crystalEntry, ok := auto.GetCrystalEntry(c.GetTypeID())
		if !ok {
			log.Error().Caller().Int32("type_id", c.GetTypeID()).Msg("createEntryItem failed")
			return nil
		}

		c.InitCrystal(
			item.CrystalEntry(crystalEntry),
		)

		// 生成初始属性
		c.GetAttManager().SetBaseAttId(-1)
		m.initCrystalAtt(c)
	}

	m.ca.Add(
		int(item.GetContainerType(i.Opts().ItemEntry.Type)),
		i.Opts().Id,
		i,
	)
	return i
}

func (m *ItemManager) initLoadedItem(i item.Itemface) error {
	entry, ok := auto.GetItemEntry(i.Opts().TypeId)
	if !ok {
		return fmt.Errorf("item<%d> entry invalid", i.Opts().TypeId)
	}

	i.Opts().ItemEntry = entry

	switch i.GetType() {
	// equip initial
	case define.Item_TypeEquip:
		e := i.(*item.Equip)
		equipEnchantEntry, ok := auto.GetEquipEnchantEntry(e.GetTypeID())
		if !ok {
			err := errors.New("invalid equip enchant entry")
			log.Error().Err(err).Caller().Int32("type_id", e.GetTypeID()).Send()
			return err
		}

		e.InitEquip(
			item.EquipEnchantEntry(equipEnchantEntry),
		)

		if e.EquipEnchantEntry != nil {
			e.GetAttManager().SetBaseAttId(int32(e.EquipEnchantEntry.AttId))
		}

	// crystal initial
	case define.Item_TypeCrystal:
		c := i.(*item.Crystal)
		crystalEntry, ok := auto.GetCrystalEntry(c.GetTypeID())
		if !ok {
			err := errors.New("invalid cyrstal entry")
			log.Error().Err(err).Caller().Int32("type_id", c.GetTypeID()).Send()
			return err
		}

		c.InitCrystal(
			item.CrystalEntry(crystalEntry),
		)

		c.GetAttManager().SetBaseAttId(-1)
	}

	m.ca.Add(
		int(item.GetContainerType(i.GetType())),
		i.Opts().Id,
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
	m.ca.Range(func(val interface{}) bool {
		it := val.(item.Itemface)

		// isn't this item
		if it.Opts().TypeId != typeMisc {
			return true
		}

		// equip should takeoff first, ignore locked equip
		if it.GetType() == define.Item_TypeEquip {
			equip := it.(*item.Equip)
			if equip.GetEquipObj() != -1 || equip.Lock {
				return true
			}
		}

		fixNum += it.Opts().Num
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
	res, err := store.GetStore().FindAll(context.Background(), define.StoreType_Item, "owner_id", m.owner.ID)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if !utils.ErrCheck(err, "FindAll failed when ItemManager.LoadAll", m.owner.ID) {
		return err
	}

	for _, v := range res {
		var opts item.ItemOptions
		vv := v.([]byte)
		if err := json.Unmarshal(vv, &opts); err != nil {
			return err
		}

		itemEntry, ok := auto.GetItemEntry(opts.TypeId)
		if !ok {
			return fmt.Errorf("ItemManager LoadAll failed: cannot find item entry<%v>", opts.TypeId)
		}

		i := item.NewItem(itemEntry.Type)
		err = json.Unmarshal(vv, i)
		if !utils.ErrCheck(err, "mapstructure NewDecoder failed", v) {
			return err
		}

		err = m.initLoadedItem(i)
		if !utils.ErrCheck(err, "initLoadedItem failed") {
			return err
		}
	}

	return err
}

func (m *ItemManager) Save(id int64) error {
	it, err := m.GetItem(id)
	if err != nil {
		return fmt.Errorf("ItemManager.Save failed: %w", err)
	}

	return store.GetStore().UpdateOne(context.Background(), define.StoreType_Item, it.Opts().Id, it)
}

func (m *ItemManager) update() {
	if time.Now().Unix() < m.nextUpdate {
		return
	}

	// 设置下次更新时间
	m.nextUpdate = time.Now().Add(itemUpdateInterval).Unix()

	// 遍历容器删除过期物品
	m.ca.Range(func(val interface{}) bool {
		it := val.(item.Itemface)
		if it.Opts().ItemEntry.TimeLife == -1 {
			return true
		}

		// 时限物品开始计时时间
		var startTime time.Time
		if it.Opts().ItemEntry.TimeStartLifeStamp == 0 {
			startTime = time.Unix(it.Opts().CreateTime, 0)
		} else {
			startTime = time.Unix(int64(it.Opts().ItemEntry.TimeStartLifeStamp), 0)
		}

		endTime := startTime.Add(time.Minute * time.Duration(it.Opts().ItemEntry.TimeLife))
		if time.Now().Unix() >= endTime.Unix() {
			log.Info().Int32("type_id", it.Opts().ItemEntry.Id).Msg("item time countdown and deleted")
			err := m.expireItem(it)
			utils.ErrPrint(err, "expireItem failed when update", it.Opts().Id, m.owner.ID)
		}

		return true
	})
}

func (m *ItemManager) GetItem(id int64) (item.Itemface, error) {
	val, ok := m.ca.Get(id)
	if ok {
		return val.(item.Itemface), nil
	} else {
		return nil, ErrItemNotFound
	}
}

func (m *ItemManager) GetItemByTypeId(typeId int32) item.Itemface {
	var retIt item.Itemface
	m.ca.Range(func(val interface{}) bool {
		it := val.(item.Itemface)
		if it.Opts().TypeId == typeId {
			retIt = it
			return false
		}
		return true
	})

	return retIt
}

func (m *ItemManager) GetItemNums(idx int) int {
	return m.ca.Size(idx)
}

func (m *ItemManager) GetItemList() []item.Itemface {
	list := make([]item.Itemface, 0, 50)

	m.ca.Range(func(val interface{}) bool {
		list = append(list, val.(item.Itemface))
		return true
	})

	return list
}

func (m *ItemManager) RangeByType(tp int, fn func(v interface{}) bool) {
	m.ca.RangeByIdx(tp, fn)
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
	idx := item.GetContainerType(itemEntry.Type)

	// 背包中有相同typeId的物品，并且是可叠加的，一定成功
	if itemEntry.MaxStack > 1 {
		m.ca.RangeByIdx(int(idx), func(val interface{}) bool {
			it := val.(item.Itemface)
			if it.Opts().TypeId == typeId {
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
		int(item.GetContainerType(entry.Type)),
		func(val interface{}) bool {
			it := val.(item.Itemface)
			if incNum <= 0 {
				return false
			}

			// 静态表中配置为不堆叠
			if entry.MaxStack <= 1 {
				return false
			}

			if it.Opts().ItemEntry.Id == typeId && it.Opts().Num < it.Opts().ItemEntry.MaxStack {
				add := incNum
				if incNum > it.Opts().ItemEntry.MaxStack-it.Opts().Num {
					add = it.Opts().ItemEntry.MaxStack - it.Opts().Num
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

		it := m.createItem(typeId, incNum)
		if it == nil {
			break
		}

		err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Item, it.Opts().Id, it)
		utils.ErrPrint(err, "UpdateManual failed when ItemManager.AddItemByTypeId", typeId, m.owner.ID)

		// prometheus ops
		prom.OpsCreateItemCounter.Inc()

		m.SendItemAdd(it)
		incNum -= it.Opts().Num
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
	gainId := it.Opts().ItemEntry.StaleGainId

	// 删除物品
	if err := m.DeleteItem(it.Opts().Id); err != nil {
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
		int(item.GetContainerType(entry.Type)),
		func(val interface{}) bool {
			it := val.(item.Itemface)
			if decNum <= 0 {
				return false
			}

			// 不是这个物品
			if it.Opts().TypeId != typeId {
				return true
			}

			// 装备还穿在身上
			if it.GetType() == define.Item_TypeEquip && it.(*item.Equip).GetEquipObj() != -1 {
				return true
			}

			if it.Opts().Num > num {
				if errModify := m.modifyNum(it, -num); errModify != nil {
					err = errModify
					utils.ErrPrint(errModify, "modifyNum failed when CostItemByTypeID", typeId, num, m.owner.ID)
					return true
				}

				m.SendItemUpdate(it)
				decNum -= num
				return false
			} else {
				decNum -= it.Opts().Num
				delId := it.Opts().Id
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

	if i.Opts().Num < num {
		return fmt.Errorf("ItemManager.CostItemByID failed: item:%d num:%d not enough, should cost %d", id, i.Opts().Num, num)
	}

	// cost
	if i.Opts().Num == num {
		err = m.delItem(id)
		if !utils.ErrCheck(err, "delItem failed when CostItemByID", id, num, m.owner.ID) {
			return err
		}

		m.SendItemDelete(id)
	} else {
		err = m.modifyNum(i, -num)
		if !utils.ErrCheck(err, "modifyNum failed when CostItemByID", id, num, m.owner.ID) {
			return err
		}

		m.SendItemUpdate(i)
	}

	return nil
}

func (m *ItemManager) CanUseItem(i item.Itemface) error {
	// 时限物品
	if i.Opts().ItemEntry.TimeLife > 0 {
		// 时限开始物品
		if i.Opts().ItemEntry.TimeStartLifeStamp > 0 {
			startTime := time.Unix(int64(i.Opts().ItemEntry.TimeStartLifeStamp), 0)
			endTime := startTime.Add(time.Duration(i.Opts().ItemEntry.TimeLife) * time.Minute)
			if time.Now().Unix() >= endTime.Unix() || time.Now().Unix() < startTime.Unix() {
				return fmt.Errorf("time life limit")
			}
		} else {
			if time.Now().Unix() >= i.Opts().CreateTime {
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

	if i.Opts().ItemEntry.EffectType == define.Item_Effect_Null {
		return fmt.Errorf("ItemManager.UseItem failed: item<%d> effect null", id)
	}

	// do effect
	if err := itemEffectFuncMapping[i.Opts().ItemEntry.EffectType](i, m.owner, m.owner); err != nil {
		return fmt.Errorf("ItemManager.UseItem failed: %w", err)
	}

	return m.CostItemByID(id, 1)
}

func (m *ItemManager) GenItemListPB() []*pbGlobal.Item {
	items := make([]*pbGlobal.Item, 0, m.GetItemNums(int(define.Container_Material)))
	m.ca.RangeByIdx(int(define.Container_Material), func(val interface{}) bool {
		it, ok := val.(*item.Item)
		if !ok {
			return true
		}

		items = append(items, it.GenItemPB())
		return true
	})

	return items
}

func (m *ItemManager) SendItemAdd(i item.Itemface) {
	msg := &pbGlobal.S2C_ItemAdd{
		Item: &pbGlobal.Item{
			Id:     i.Opts().Id,
			TypeId: int32(i.Opts().TypeId),
			Num:    int32(i.Opts().Num),
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
			Id:     i.Opts().Id,
			TypeId: int32(i.Opts().TypeId),
			Num:    int32(i.Opts().Num),
		},
	}

	m.owner.SendProtoMessage(msg)
}
