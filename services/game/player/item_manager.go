package player

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/internal/container"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/item"
	"bitbucket.org/funplus/server/services/game/prom"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	"bitbucket.org/funplus/server/utils/random"
	"github.com/mitchellh/mapstructure"
	log "github.com/rs/zerolog/log"
	"github.com/valyala/bytebufferpool"
)

func MakeItemKey(it item.Itemface, fields ...string) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	var keyName string
	switch e := it.GetType(); e {
	case define.Item_TypePresent:
		fallthrough
	case define.Item_TypeItem:
		keyName = "item_list"
	case define.Item_TypeEquip:
		keyName = "equip_list"
	case define.Item_TypeCrystal:
		keyName = "crystal_list"
	default:
		log.Error().Caller().Int32("item_type", e).Msg("MakeItemKey failed, invalid item type")
		return ""
	}

	_, _ = b.WriteString(keyName)
	_, _ = b.WriteString(".id_")
	_, _ = b.WriteString(strconv.Itoa(int(it.Opts().Id)))

	for _, f := range fields {
		_, _ = b.WriteString(".")
		_, _ = b.WriteString(f)
	}

	return b.String()
}

// item effect mapping function
type effectFunc func(item.Itemface, *Player, *Player) error

// 物品使用效果
var itemEffectFuncMapping = map[int32]effectFunc{
	define.Item_Effect_Null: itemEffectNull,
	define.Item_Effect_Loot: itemEffectLoot,
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
	CA         *container.ContainerArray `bson:"-" json:"-"` // 背包列表 0:材料与消耗 1:装备 2:晶石

	// 反射生成表结构，无其他用处
	KeepItemList    map[int64]interface{} `bson:"item_list" json:"item_list"`
	KeepEquipList   map[int64]interface{} `bson:"equip_list" json:"equip_list"`
	KeepCrystalList map[int64]interface{} `bson:"crystal_list" json:"crystal_list"`
}

func NewItemManager(owner *Player) *ItemManager {
	m := &ItemManager{
		nextUpdate: time.Now().Add(itemUpdateInterval).Unix(),
		owner:      owner,
		CA:         container.New(int(define.Container_End)),

		KeepItemList:    make(map[int64]interface{}),
		KeepEquipList:   make(map[int64]interface{}),
		KeepCrystalList: make(map[int64]interface{}),
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
	for _, v := range i.Opts().ItemEntry.EffectValue {
		if err := owner.CostLootManager().CanGain(v); err != nil {
			return err
		}
	}

	for _, v := range i.Opts().ItemEntry.EffectValue {
		if err := owner.CostLootManager().GainLoot(v); err != nil {
			log.Warn().
				Int32("loot_id", v).
				Int32("item_type_id", i.Opts().TypeId).
				Msg("itemEffectLoot failed")
		}
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

	i.Opts().Num = add
	i.Opts().CreateTime = time.Now().Unix()

	fields := map[string]interface{}{
		MakeItemKey(i): i,
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

	fieldsName := []string{MakeItemKey(it)}
	err := store.GetStore().DeleteFields(define.StoreType_Item, m.owner.ID, fieldsName)
	item.GetItemPool(it.GetType()).Put(it)

	return err
}

func (m *ItemManager) modifyNum(i item.Itemface, add int32) error {
	i.Opts().Num += add

	fields := map[string]interface{}{
		MakeItemKey(i, "num"): i.Opts().Num,
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

	i := item.NewItem(entry.Type)

	// item initial
	i.InitItem(
		item.Id(id),
		item.OwnerId(m.owner.GetID()),
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

	m.CA.Add(
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

	m.CA.Add(
		int(item.GetContainerType(i.GetType())),
		i.Opts().Id,
		i,
	)
	return nil
}

// 初始化晶石属性
func (m *ItemManager) initCrystalAtt(c *item.Crystal) {
	globalConfig, _ := auto.GetGlobalConfig()

	// 初始主属性
	mainAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, define.Crystal_AttTypeMain)
	mainAttItem, err := random.PickOne(mainAttRepoList, func(random.Item) bool {
		return true
	})
	if err != nil {
		log.Error().Err(err).Int64("crystal_id", c.Id).Msg("pick crystal main att failed")
		return
	}

	// 记录主属性库id
	mainAttRepoEntry := mainAttItem.(*auto.CrystalAttRepoEntry)
	c.MainAtt.AttRepoId = mainAttRepoEntry.Id
	c.MainAtt.AttRandRatio = random.Int32(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1])

	// 随机几条副属性
	viceAttNum := auto.GetCrystalInitViceAttNum(c.ItemEntry.Quality)

	// 初始副属性
	viceAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, define.Crystal_AttTypeVice)
	viceAttItems, err := random.PickUnrepeated(viceAttRepoList, viceAttNum, func(random.Item) bool {
		return true
	})

	if errors.Is(err, random.ErrNoResult) {
		return
	}

	if err != nil {
		log.Error().Err(err).Int64("crystal_id", c.Id).Msg("pick unrepeated crystal vice att failed")
		return
	}

	for _, v := range viceAttItems {
		viceAttRepoEntry := v.(*auto.CrystalAttRepoEntry)
		c.ViceAtts = append(c.ViceAtts, item.CrystalAtt{
			AttRepoId:    viceAttRepoEntry.Id,
			AttRandRatio: random.Int32(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1]),
		})
	}
}

// 新增副属性
func (m *ItemManager) generateViceAtt(c *item.Crystal) {
	if c == nil {
		return
	}

	globalConfig, _ := auto.GetGlobalConfig()

	attType := make(map[int]struct{}, 20)
	for _, att := range c.ViceAtts {
		attType[int(att.AttRepoId)] = struct{}{}
	}

	// 副属性已满4条
	if len(attType) >= define.Crystal_ViceAttNum {
		return
	}

	// 不满4条，则随机一条未曾有过的属性类型
	limiter := func(it random.Item) bool {
		if _, ok := attType[it.GetId()]; ok {
			return false
		}
		return true
	}
	viceAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, define.Crystal_AttTypeVice)
	it, err := random.PickOne(viceAttRepoList, limiter)
	if pass := utils.ErrCheck(err, "pick one vice att failed", c.Id); !pass {
		return
	}

	attRepoEntry := it.(*auto.CrystalAttRepoEntry)
	c.ViceAtts = append(c.ViceAtts, item.CrystalAtt{
		AttRepoId:    attRepoEntry.Id,
		AttRandRatio: random.Int32(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1]),
	})
}

// 强化副属性
func (m *ItemManager) enforceViceAtt(c *item.Crystal) {
	if c == nil {
		return
	}

	globalConfig, _ := auto.GetGlobalConfig()

	// 所有副属性种类
	attType := make(map[int]struct{}, 20)
	for _, att := range c.ViceAtts {
		attType[int(att.AttRepoId)] = struct{}{}
	}

	// 限制器：只能强化晶石已有的副属性
	limiter := func(item random.Item) bool {
		if _, ok := attType[item.GetId()]; ok {
			return true
		}
		return false
	}

	viceAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, define.Crystal_AttTypeVice)
	it, err := random.PickOne(viceAttRepoList, limiter)
	if pass := utils.ErrCheck(err, "pick one vice att failed", c.Id); !pass {
		return
	}

	viceAttRepoEntry := it.(*auto.CrystalAttRepoEntry)
	c.ViceAtts = append(c.ViceAtts, item.CrystalAtt{
		AttRepoId:    viceAttRepoEntry.Id,
		AttRandRatio: random.Int32(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1]),
	})
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
	loadItems := struct {
		ItemList    map[string]interface{} `bson:"item_list" json:"item_list"`
		EquipList   map[string]interface{} `bson:"equip_list" json:"equip_list"`
		CrystalList map[string]interface{} `bson:"crystal_list" json:"crystal_list"`
	}{
		ItemList:    make(map[string]interface{}),
		EquipList:   make(map[string]interface{}),
		CrystalList: make(map[string]interface{}),
	}

	err := store.GetStore().LoadObject(define.StoreType_Item, m.owner.ID, &loadItems)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("ItemManager LoadAll failed: %w", err)
	}

	loadFn := func(list map[string]interface{}) error {
		for _, v := range list {
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

			i := item.NewItem(itemEntry.Type)

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

		return nil
	}

	err = loadFn(loadItems.ItemList)
	utils.ErrPrint(err, "ItemManager load items failed", m.owner.ID)

	err = loadFn(loadItems.EquipList)
	utils.ErrPrint(err, "ItemManager load equips failed", m.owner.ID)

	err = loadFn(loadItems.CrystalList)
	utils.ErrPrint(err, "ItemManager load crystals failed", m.owner.ID)

	return err
}

func (m *ItemManager) Save(id int64) error {
	it, err := m.GetItem(id)
	if err != nil {
		return fmt.Errorf("ItemManager.Save failed: %w", err)
	}

	fields := map[string]interface{}{
		MakeItemKey(it): it,
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
	idx := item.GetContainerType(itemEntry.Type)

	// 背包中有相同typeId的物品，并且是可叠加的，一定成功
	if itemEntry.MaxStack > 1 {
		m.CA.RangeByIdx(int(idx), func(val interface{}) bool {
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

		i := m.createItem(typeId, incNum)
		if i == nil {
			break
		}

		m.CA.Add(
			int(item.GetContainerType(i.GetType())),
			i.Opts().Id,
			i,
		)
		m.SendItemAdd(i)
		incNum -= i.Opts().Num
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

	m.CA.RangeByIdx(
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

// 装备升级
func (m *ItemManager) EquipLevelup(equipId int64, stuffItems, expItems []int64) error {
	i, err := m.GetItem(equipId)
	utils.ErrPrint(err, "EquipLevelup failed", equipId, m.owner.ID)

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return errors.New("invalid global config")
	}

	if i.GetType() != define.Item_TypeEquip {
		return fmt.Errorf("EquipLevelup failed, wrong item<%d> type", i.Opts().TypeId)
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

	// 所有合法的消耗物品及对应的经验值
	itemExps := make(map[int64]int32)

	// 吞噬材料
	for _, id := range stuffItems {
		it, err := m.GetItem(id)
		if pass := utils.ErrCheck(err, "cannot find item", id); !pass {
			continue
		}

		if it.Opts().ItemEntry.Type != define.Item_TypeEquip {
			continue
		}

		stuffEquip := it.(*item.Equip)

		// 1级经验不算折损率
		equipLv1Entry, ok := auto.GetEquipLevelupEntry(1)
		if !ok {
			log.Error().Caller().Msg("can not find equip levelup 1 entry")
			continue
		}
		equipLv1Exp := equipLv1Entry.Exp[stuffEquip.ItemEntry.Quality]

		// 已升级累计的经验
		equipLvEntry, ok := auto.GetEquipLevelupEntry(int32(stuffEquip.Level))
		if !ok {
			log.Error().Caller().Int8("level", stuffEquip.Level).Msg("can not find equip levelup entry")
			continue
		}
		equiplvTotalExp := equipLvEntry.Exp[stuffEquip.ItemEntry.Quality] + stuffEquip.Exp - equipLv1Exp

		// 物品总经验 = 物品1级经验 + 已消耗所有经验 * 经验折损率
		itemExps[id] = int32(int64(equipLv1Exp) + int64(equiplvTotalExp)*int64(globalConfig.EquipSwallowExpLoss)/int64(define.PercentBase))
	}

	// 经验道具
	for _, id := range expItems {
		it, err := m.GetItem(id)
		if pass := utils.ErrCheck(err, "cannot find item", id); !pass {
			continue
		}

		if it.GetType() != define.Item_TypeItem {
			continue
		}

		if it.Opts().ItemEntry.SubType != define.Item_SubType_Item_EquipExp {
			continue
		}

		itemExps[id] = it.Opts().ItemEntry.PublicMisc[0]
	}

	// 升级处理
	levelupFn := func(itemId int64, exp int32) bool {
		_, ok := auto.GetEquipLevelupEntry(int32(equip.Level) + 1)
		if !ok {
			return false
		}

		// 金币限制
		costGold := int32(int64(exp) * int64(globalConfig.EquipLevelupExpGoldRatio) / int64(define.PercentBase))
		if costGold < 0 {
			return false
		}

		if err := m.owner.TokenManager().CanCost(define.Token_Gold, costGold); err != nil {
			return false
		}

		// overflow
		if equip.Exp+exp < 0 {
			return false
		}

		equip.Exp += exp
		for {
			curLevelEntry, _ := auto.GetEquipLevelupEntry(int32(equip.Level))
			nextLevelEntry, ok := auto.GetEquipLevelupEntry(int32(equip.Level) + 1)
			if !ok {
				break
			}

			levelExp := nextLevelEntry.Exp[equip.ItemEntry.Quality] - curLevelEntry.Exp[equip.ItemEntry.Quality]
			if equip.Exp < levelExp {
				break
			}

			equip.Level++
			equip.Exp -= levelExp
		}

		// 消耗
		err := m.owner.TokenManager().DoCost(define.Token_Gold, costGold)
		utils.ErrPrint(err, "TokenManager DoCost failed", costGold)

		err = m.owner.ItemManager().CostItemByID(itemId, 1)
		utils.ErrPrint(err, "ItemManager CostItemByID failed", itemId)
		return true
	}

	modified := false
	for itemId, exp := range itemExps {
		if !levelupFn(itemId, exp) {
			break
		}

		modified = true
	}

	// 经验等级道具均没有改变
	if !modified {
		return nil
	}

	// save
	fields := map[string]interface{}{
		MakeItemKey(equip, "level"): equip.Level,
		MakeItemKey(equip, "exp"):   equip.Exp,
	}
	err = store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
	utils.ErrPrint(err, "SaveFields failed when EquipLevelup", equip.GetID(), m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)

	return err
}

// 装备突破
func (m *ItemManager) EquipPromote(equipId int64) error {
	it, err := m.GetItem(equipId)
	utils.ErrPrint(err, "EquipPromote failed", equipId, m.owner.ID)

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return errors.New("invalid global config")
	}

	if it.GetType() != define.Item_TypeEquip {
		return errors.New("invalid item type")
	}

	equip := it.(*item.Equip)
	if equip.Promote >= define.Equip_Max_Promote_Times {
		return errors.New("promote times full")
	}

	// 队伍等级
	if m.owner.Level < globalConfig.EquipPromoteLevelLimit[equip.Promote+1] {
		return errors.New("team level limit")
	}

	// 装备是否达到等级上限
	levelupEntry, ok := auto.GetEquipLevelupEntry(int32(equip.Level) + 1)
	if !ok {
		return errors.New("equip reach max level")
	}

	if int32(equip.Promote) >= levelupEntry.PromoteLimit {
		return errors.New("equip has not levelup to max")
	}

	// 消耗
	costId := equip.EquipEnchantEntry.PromoteCostId[equip.Promote+1]
	err = m.owner.CostLootManager().CanCost(costId)
	if pass := utils.ErrCheck(err, "EquipPromote can cost failed", equipId, costId, m.owner.ID); !pass {
		return err
	}

	err = m.owner.CostLootManager().DoCost(costId)
	if pass := utils.ErrCheck(err, "EquipPromote do cost failed", equipId, costId, m.owner.ID); !pass {
		return err
	}

	equip.Promote++

	// save
	fields := map[string]interface{}{
		MakeItemKey(equip, "promote"): equip.Promote,
	}
	err = store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
	utils.ErrPrint(err, "SaveFields failed when EquipPromote", equip.GetID(), m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)
	return err
}

// 晶石升级
func (m *ItemManager) CrystalLevelup(crystalId int64, stuffItems, expItems []int64) error {
	it, err := m.GetItem(crystalId)
	utils.ErrPrint(err, "CrystalLevelup failed", crystalId, m.owner.ID)

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return errors.New("invalid global config")
	}

	if it.GetType() != define.Item_TypeCrystal {
		return fmt.Errorf("CrystalLevelup failed, wrong item<%d> type", it.Opts().TypeId)
	}

	c := it.(*item.Crystal)
	_, ok = auto.GetCrystalLevelupEntry(int32(c.Level) + 1)
	if !ok {
		return fmt.Errorf("CyrstalLevelup failed, cannot find crystal levelup entry<%d>", c.Level+1)
	}

	// 品质限制等级上限
	if int32(c.Level) >= globalConfig.CrystalLevelupQualityLimit[c.ItemEntry.Quality] {
		return errors.New("crystal quality limit")
	}

	// 所有合法的消耗物品及对应的经验值
	itemExps := make(map[int64]int32)

	// 吞噬材料
	for _, id := range stuffItems {
		it, err := m.GetItem(id)
		if pass := utils.ErrCheck(err, "cannot find item", id); !pass {
			continue
		}

		if it.Opts().ItemEntry.Type != define.Item_TypeCrystal {
			continue
		}

		stuffCrystal := it.(*item.Crystal)

		// 1级经验不算折损率
		crystalLv1Entry, ok := auto.GetCrystalLevelupEntry(1)
		if !ok {
			log.Error().Caller().Msg("can not find crystal levelup 1 entry")
			continue
		}
		crystalLv1Exp := crystalLv1Entry.Exp[stuffCrystal.ItemEntry.Quality]

		// 已升级累计的经验
		crystalLvEntry, ok := auto.GetCrystalLevelupEntry(int32(stuffCrystal.Level))
		if !ok {
			log.Error().Caller().Int8("level", stuffCrystal.Level).Msg("can not find crystal levelup entry")
			continue
		}
		crystallvTotalExp := crystalLvEntry.Exp[stuffCrystal.ItemEntry.Quality] + stuffCrystal.Exp - crystalLv1Exp

		// 物品总经验 = 物品1级经验 + 已消耗所有经验 * 经验折损率
		itemExps[id] = int32(int64(crystalLv1Exp) + int64(crystallvTotalExp)*int64(globalConfig.CrystalSwallowExpLoss)/int64(define.PercentBase))
	}

	// 经验道具
	for _, id := range expItems {
		it, err := m.GetItem(id)
		if pass := utils.ErrCheck(err, "cannot find item", id); !pass {
			continue
		}

		if it.GetType() != define.Item_TypeItem {
			continue
		}

		if it.Opts().ItemEntry.SubType != define.Item_SubType_Item_CrystalExp {
			continue
		}

		itemExps[id] = it.Opts().ItemEntry.PublicMisc[0]
	}

	// 升级处理
	levelupFn := func(itemId int64, exp int32) bool {
		_, ok := auto.GetCrystalLevelupEntry(int32(c.Level) + 1)
		if !ok {
			return false
		}

		// 判断金币
		costGold := int32(int64(exp) * int64(globalConfig.CrystalLevelupExpGoldRatio) / int64(define.PercentBase))
		if costGold < 0 {
			return false
		}

		// 品质限制等级上限
		if int32(c.Level) >= globalConfig.CrystalLevelupQualityLimit[c.ItemEntry.Quality] {
			return false
		}

		if err := m.owner.TokenManager().CanCost(define.Token_Gold, costGold); err != nil {
			return false
		}

		// overflow
		if c.Exp+exp < 0 {
			return false
		}

		c.Exp += exp
		for {
			curLevelEntry, _ := auto.GetCrystalLevelupEntry(int32(c.Level))
			nextLevelEntry, ok := auto.GetCrystalLevelupEntry(int32(c.Level) + 1)
			if !ok {
				break
			}

			levelExp := nextLevelEntry.Exp[c.ItemEntry.Quality] - curLevelEntry.Exp[c.ItemEntry.Quality]
			if c.Exp < levelExp {
				break
			}

			c.Level++
			c.Exp -= levelExp
			for _, level := range globalConfig.CrystalViceAttAddLevel {
				if int32(c.Level) == level {
					// 增加新的副属性直到满4条
					m.generateViceAtt(c)

					// 强化副属性
					m.enforceViceAtt(c)
					// c.GetAttManager().CalcAtt()
					// m.SendCrystalAttUpdate(c)
					break
				}
			}
		}

		// 消耗材料
		err = m.CostItemByID(itemId, 1)
		utils.ErrPrint(err, "ItemManager CostItemByID failed", itemId)

		// 消耗金币
		err = m.owner.TokenManager().DoCost(define.Token_Gold, costGold)
		utils.ErrPrint(err, "TokenManager DoCost failed", costGold)
		return true
	}

	modified := false
	for itemId, exp := range itemExps {
		if !levelupFn(itemId, exp) {
			break
		}

		modified = true
	}

	// 经验等级道具均没有改变
	if !modified {
		return nil
	}

	// save
	fields := map[string]interface{}{
		MakeItemKey(c): c,
	}
	err = store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
	if pass := utils.ErrCheck(err, "CrystalLevelup SaveFields failed", m.owner.ID, c.Level, c.Exp); !pass {
		return err
	}

	m.SendCrystalUpdate(c)
	return nil
}

func (m *ItemManager) SaveCrystalEquiped(c *item.Crystal) {
	fields := map[string]interface{}{
		MakeItemKey(c, "crystal_obj"): c.CrystalObj,
	}

	err := store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
	utils.ErrPrint(err, "SaveCrystalEquiped failed", c.Id)
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

func (m *ItemManager) SendCrystalAttUpdate(c *item.Crystal) {
	msg := &pbGlobal.S2C_CrystalAttUpdate{
		CrystalId: c.Id,
		AttValue:  make([]int32, define.Att_End),
	}

	for n := 0; n < define.Att_End; n++ {
		msg.AttValue[n] = c.GetAttManager().GetAttValue(n)
	}

	m.owner.SendProtoMessage(msg)
}

func (m *ItemManager) SendCrystalUpdate(c *item.Crystal) {
	msg := &pbGlobal.S2C_CrystalUpdate{
		CrystalId: c.Id,
		CrystalData: &pbGlobal.CrystalData{
			Level:      int32(c.Level),
			Exp:        c.Exp,
			CrystalObj: c.CrystalObj,
			MainAtt: &pbGlobal.CrystalAtt{
				AttRepoId:    c.MainAtt.AttRepoId,
				AttRandRatio: c.MainAtt.AttRandRatio,
			},
			ViceAtts: make([]*pbGlobal.CrystalAtt, len(c.ViceAtts)),
		},
	}

	for n, att := range c.ViceAtts {
		msg.CrystalData.ViceAtts[n] = &pbGlobal.CrystalAtt{
			AttRepoId:    att.AttRepoId,
			AttRandRatio: att.AttRandRatio,
		}
	}

	m.owner.SendProtoMessage(msg)
}
