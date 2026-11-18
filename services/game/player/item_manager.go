package player

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/item"
	"github.com/east-eden/server/services/game/prom"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	ErrItemNotFound    = errors.New("item not found")
	ErrItemNotEnough   = errors.New("item not enough")
	ErrItemInvalidType = errors.New("invalid item type")
)

// ItemManager 物品账本。
//
// 采用 kof 的模型：一个道具配置 id 只对应一条堆叠记录（Id == EntryId），
// 不再区分装备/晶石实例。所有增减都通过 OnItemAdd / OnItemCost 收口，
// 便于统一接入事件、任务与埋点。
type ItemManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	Items      map[int32]*item.Item `bson:"items" json:"items"`
	nextUpdate int64                `bson:"-" json:"-"`
	owner      *Player              `bson:"-" json:"-"`
}

func NewItemManager(owner *Player) *ItemManager {
	return &ItemManager{
		owner: owner,
		Items: make(map[int32]*item.Item, 64),
	}
}

func (m *ItemManager) Destroy() {
	m.Items = make(map[int32]*item.Item, 64)
}

func (m *ItemManager) GetCostLootType() int32 {
	return define.CostLoot_Item
}

// OnItemAdd 增加道具：按 StackMax 截断堆叠，返回实际入账数量。
func (m *ItemManager) OnItemAdd(entryId int32, num int64, reason string) (int64, error) {
	if num <= 0 {
		return 0, nil
	}

	entry, ok := config.GetGoodsEntry(entryId)
	if !ok {
		return 0, fmt.Errorf("goods entry<%d> not found", entryId)
	}

	// 英雄类道具不进物品账本，直接进英雄系统（kof 的 heroItemAdd 口径）
	if hit, err := m.heroGoodsAdd(entry, num); hit {
		return num, err
	}

	if m.Items == nil {
		m.Items = make(map[int32]*item.Item, 64)
	}

	it, found := m.Items[entryId]
	if found && entry.StackMax > 0 && it.Num >= entry.StackMax {
		return 0, nil
	}

	if !found {
		it = item.NewItem()
		it.InitItem(
			item.Id(int64(entryId)),
			item.OwnerId(m.owner.ID),
			item.EntryId(entryId),
			item.GoodsEntry(entry),
		)

		m.Items[entryId] = it
	}

	add := num
	if entry.StackMax > 0 && it.Num+add > entry.StackMax {
		add = entry.StackMax - it.Num
	}

	if add <= 0 {
		return 0, nil
	}

	it.Num += add

	if err := m.save(it); err != nil {
		return 0, err
	}

	m.SendItemAdd(it)

	prom.OpsCreateItemCounter.Inc()
	log.Debug().
		Int32("entry_id", entryId).
		Int64("num", add).
		Str("reason", reason).
		Msg("item gain")

	return add, nil
}

// OnItemCost 扣除道具。调用方需先用 OnCheckEnough 校验余量。
func (m *ItemManager) OnItemCost(entryId int32, num int64, reason string) error {
	if num <= 0 {
		return nil
	}

	// 英雄碎片从 chip 计数里扣
	if entry, ok := config.GetGoodsEntry(entryId); ok && entry.Type == config.GoodsType_Hero {
		if hit, err := m.heroGoodsCost(entry, num); hit {
			return err
		}
	}

	it, ok := m.Items[entryId]
	if !ok {
		return ErrItemNotFound
	}

	if it.Num < num {
		return ErrItemNotEnough
	}

	it.Num -= num

	if it.Num == 0 {
		delete(m.Items, entryId)
		m.SendItemDelete(it.Opts().Id)

		if err := store.GetStore().DeleteOne(context.Background(), define.StoreType_Item, it.Opts().Id); err != nil {
			return err
		}
	} else {
		m.SendItemUpdate(it)

		if err := m.save(it); err != nil {
			return err
		}
	}

	log.Debug().
		Int32("entry_id", entryId).
		Int64("num", num).
		Str("reason", reason).
		Msg("item cost")

	return nil
}

// OnCheckEnough 校验道具是否足够
func (m *ItemManager) OnCheckEnough(entryId int32, num int64) error {
	// 英雄碎片按 chip 计数判断
	if entry, ok := config.GetGoodsEntry(entryId); ok && entry.Type == config.GoodsType_Hero {
		if entry.SubType == config.GoodsSubType_HeroChip {
			if have := m.owner.HeroManager().GetHeroChipCount(entryId); int64(have) < num {
				return fmt.Errorf("%w: chip<%d> has %d, need %d", ErrItemNotEnough, entryId, have, num)
			}

			return nil
		}
	}

	it, ok := m.Items[entryId]
	if !ok {
		return fmt.Errorf("%w: entry<%d>", ErrItemNotEnough, entryId)
	}

	if it.Num < num {
		return fmt.Errorf("%w: entry<%d> has %d, need %d", ErrItemNotEnough, entryId, it.Num, num)
	}

	return nil
}

func (m *ItemManager) GetItemCount(entryId int32) int64 {
	if it, ok := m.Items[entryId]; ok {
		return it.Num
	}

	return 0
}

// SetItemCount 把道具数量直接调整到目标值
func (m *ItemManager) SetItemCount(entryId int32, target int64, reason string) error {
	cur := m.GetItemCount(entryId)
	switch {
	case target > cur:
		_, err := m.OnItemAdd(entryId, target-cur, reason)
		return err
	case target < cur:
		return m.OnItemCost(entryId, cur-target, reason)
	default:
		return nil
	}
}

////////////////////////////////////////////////////////
// CostLooter 适配

func (m *ItemManager) CanCost(typeMisc int32, num int32) error {
	return m.OnCheckEnough(typeMisc, int64(num))
}

func (m *ItemManager) DoCost(typeMisc int32, num int32) error {
	return m.OnItemCost(typeMisc, int64(num), "cost")
}

func (m *ItemManager) CanGain(typeMisc int32, num int32) error {
	if _, ok := config.GetGoodsEntry(typeMisc); !ok {
		return fmt.Errorf("goods entry<%d> not found", typeMisc)
	}

	return nil
}

func (m *ItemManager) GainLoot(typeMisc int32, num int32) error {
	_, err := m.OnItemAdd(typeMisc, int64(num), "loot")
	return err
}

////////////////////////////////////////////////////////
// 查询

func (m *ItemManager) GetItem(id int64) (item.Itemface, error) {
	it, ok := m.Items[int32(id)]
	if !ok {
		return nil, ErrItemNotFound
	}

	return it, nil
}

func (m *ItemManager) GetItemByTypeId(typeId int32) item.Itemface {
	if it, ok := m.Items[typeId]; ok {
		return it
	}

	return nil
}

// GetItemNums 返回道具种类数（旧的按容器统计已不再适用）
func (m *ItemManager) GetItemNums() int {
	return len(m.Items)
}

func (m *ItemManager) GetItemList() []item.Itemface {
	list := make([]item.Itemface, 0, len(m.Items))
	for _, it := range m.Items {
		list = append(list, it)
	}

	return list
}

// RangeByType 遍历所有道具。旧的按容器索引遍历已合并为单一账本。
func (m *ItemManager) RangeByType(fn func(v item.Itemface) bool) {
	for _, it := range m.Items {
		if !fn(it) {
			return
		}
	}
}

func (m *ItemManager) CanAddItem(typeId int32, num int32) bool {
	entry, ok := config.GetGoodsEntry(typeId)
	if !ok {
		return false
	}

	if entry.StackMax <= 0 {
		return true
	}

	it, found := m.Items[typeId]
	if !found {
		return int64(num) <= entry.StackMax
	}

	return it.Num+int64(num) <= entry.StackMax
}

////////////////////////////////////////////////////////
// 兼容旧调用方的按类型接口

func (m *ItemManager) AddItemByTypeId(typeId int32, num int32) error {
	_, err := m.OnItemAdd(typeId, int64(num), "add")
	return err
}

func (m *ItemManager) DeleteItem(id int64) error {
	entryId := int32(id)
	it, ok := m.Items[entryId]
	if !ok {
		return ErrItemNotFound
	}

	delete(m.Items, entryId)
	m.SendItemDelete(it.Opts().Id)

	return store.GetStore().DeleteOne(context.Background(), define.StoreType_Item, id)
}

func (m *ItemManager) CostItemByTypeId(typeId int32, num int32) error {
	return m.OnItemCost(typeId, int64(num), "cost")
}

func (m *ItemManager) CostItemByID(id int64, num int32) error {
	return m.OnItemCost(int32(id), int64(num), "cost")
}

// UseItem 使用一个道具。
//
// 使用效果分发（kof 的 FnUse 表）依赖充值/活动/挂机等尚未移植的子系统，
// 这里先只做数量扣减，效果分发在对应子系统接入时补齐。
func (m *ItemManager) UseItem(id int64) error {
	entryId := int32(id)

	if _, ok := config.GetGoodsEntry(entryId); !ok {
		return ErrItemInvalidType
	}

	if err := m.OnCheckEnough(entryId, 1); err != nil {
		return err
	}

	return m.OnItemCost(entryId, 1, "use")
}

////////////////////////////////////////////////////////
// 持久化

func (m *ItemManager) save(it *item.Item) error {
	return store.GetStore().UpdateOne(context.Background(), define.StoreType_Item, it.Opts().Id, it)
}

func (m *ItemManager) LoadAll() error {
	res, err := store.GetStore().FindAll(context.Background(), define.StoreType_Item, "owner_id", m.owner.ID)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if !utils.ErrCheck(err, "FindAll failed when ItemManager.LoadAll", m.owner.ID) {
		return err
	}

	m.Items = make(map[int32]*item.Item, len(res))

	for _, v := range res {
		var opts item.ItemOptions
		vv, ok := v.([]byte)
		if !ok {
			continue
		}

		if err := json.Unmarshal(vv, &opts); err != nil {
			return err
		}

		entry, ok := config.GetGoodsEntry(opts.EntryId)
		if !ok {
			// 配置里已删除的道具，直接丢弃并记录
			log.Warn().Int32("entry_id", opts.EntryId).Int64("owner_id", m.owner.ID).Msg("item entry missing, dropped")
			continue
		}

		it := item.NewItem()
		it.InitItem(
			item.Id(opts.Id),
			item.OwnerId(opts.OwnerId),
			item.EntryId(opts.EntryId),
			item.Num(opts.Num),
			item.CreateTime(opts.CreateTime),
			item.GoodsEntry(entry),
		)

		m.Items[opts.EntryId] = it
	}

	return nil
}

// update 保留 tick 入口，供后续过期道具等逻辑接入
func (m *ItemManager) update() {}

////////////////////////////////////////////////////////
// 协议

func (m *ItemManager) GenItemListPB() []*pbGlobal.Item {
	list := make([]*pbGlobal.Item, 0, len(m.Items))
	for _, it := range m.Items {
		list = append(list, it.GenItemPB())
	}

	return list
}

func (m *ItemManager) SendItemAdd(it item.Itemface) {
	pb := &pbGlobal.S2C_ItemAdd{Item: it.GenItemPB()}
	m.owner.SendProtoMessage(pb)
}

func (m *ItemManager) SendItemDelete(id int64) {
	pb := &pbGlobal.S2C_DelItem{ItemId: id}
	m.owner.SendProtoMessage(pb)
}

func (m *ItemManager) SendItemUpdate(it item.Itemface) {
	pb := &pbGlobal.S2C_ItemUpdate{Item: it.GenItemPB()}
	m.owner.SendProtoMessage(pb)
}

// heroGoodsAdd 英雄类道具入账：卡牌直接获得英雄，碎片累加到 chip 计数。
// 返回 hit=true 表示已由英雄系统处理，不再进物品账本。
func (m *ItemManager) heroGoodsAdd(entry *config.GoodsEntry, num int64) (bool, error) {
	if entry.Type != config.GoodsType_Hero {
		return false, nil
	}

	switch entry.SubType {
	case config.GoodsSubType_HeroCard:
		for i := int64(0); i < num; i++ {
			if _, err := m.owner.HeroManager().AddHeroByTypeId(entry.Id); err != nil {
				return true, err
			}
		}

		return true, nil

	case config.GoodsSubType_HeroChip:
		if err := m.owner.HeroManager().AddHeroChip(entry.Id, num); err != nil {
			return true, err
		}

		return true, nil
	}

	return false, nil
}

// heroGoodsCost 英雄类道具扣除，目前只有碎片可扣
func (m *ItemManager) heroGoodsCost(entry *config.GoodsEntry, num int64) (bool, error) {
	if entry.Type != config.GoodsType_Hero {
		return false, nil
	}

	if entry.SubType == config.GoodsSubType_HeroChip {
		return true, m.owner.HeroManager().CostHeroChip(entry.Id, int32(num))
	}

	return false, nil
}
