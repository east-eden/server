package player

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/collection"
	"bitbucket.org/funplus/server/services/game/event"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	ErrCollectionNotFound       = errors.New("collection not found")
	ErrCollectionStarMax        = errors.New("collection star max")
	ErrCollectionAlreadyActived = errors.New("collection already actived")
)

type CollectionManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`
	event.EventRegister   `bson:"-" json:"-"`

	owner          *Player                          `bson:"-" json:"-"`
	CollectionList map[int32]*collection.Collection `bson:"-" json:"-"` // 收集品列表
}

func NewCollectionManager(owner *Player) *CollectionManager {
	m := &CollectionManager{
		owner:          owner,
		CollectionList: make(map[int32]*collection.Collection),
	}

	m.RegisterEvent()

	return m
}

func (m *CollectionManager) Destroy() {
	for _, c := range m.CollectionList {
		collection.GetCollectionPool().Put(c)
	}
}

// 事件注册
func (m *CollectionManager) RegisterEvent() {
	m.owner.eventManager.Register(define.Event_Type_Sign, m.onEventSign)
}

func (m *CollectionManager) onEventSign(e *event.Event) error {
	log.Info().Interface("event", e).Msg("CollectionManager.onEventSign")
	return nil
}

func (m *CollectionManager) createEntryCollection(entry *auto.CollectionEntry) *collection.Collection {
	if entry == nil {
		log.Error().Msg("newEntryCollection with nil CollectionEntry")
		return nil
	}

	c := collection.NewCollection()
	c.Init(
		collection.TypeId(entry.Id),
		collection.OwnerId(m.owner.GetId()),
		collection.Entry(entry),
	)

	m.CollectionList[c.GetOptions().TypeId] = c

	return c
}

func (m *CollectionManager) initLoadedCollection(c *collection.Collection) error {
	entry, ok := auto.GetCollectionEntry(c.GetOptions().TypeId)
	if !ok {
		return fmt.Errorf("CollectionManager initLoadedCollection: collection<%d> entry invalid", c.GetOptions().TypeId)
	}

	c.GetOptions().Entry = entry

	m.CollectionList[c.GetOptions().TypeId] = c

	return nil
}

// interface of cost_loot
func (m *CollectionManager) GetCostLootType() int32 {
	return define.CostLoot_Collection
}

func (m *CollectionManager) CanCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.CanCost(typeMisc, num)
	if err != nil {
		return err
	}

	var fixNum int32
	for _, v := range m.CollectionList {
		if v.GetOptions().TypeId == typeMisc {
			fixNum++
		}
	}

	if fixNum >= num {
		return nil
	}

	return fmt.Errorf("not enough hero<%d>, num<%d>", typeMisc, num)
}

func (m *CollectionManager) DoCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.DoCost(typeMisc, num)
	if err != nil {
		return err
	}

	var costNum int32
	for _, v := range m.CollectionList {
		if v.GetOptions().TypeId == typeMisc {
			m.DelCollection(v.GetOptions().TypeId)
			costNum++
		}
	}

	if costNum < num {
		log.Warn().
			Int32("cost_type_misc", typeMisc).
			Int32("cost_num", num).
			Int32("actual_cost_num", costNum).
			Msg("collection manager cost num error")
		return nil
	}

	return nil
}

func (m *CollectionManager) GainLoot(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.GainLoot(typeMisc, num)
	if err != nil {
		return err
	}

	var n int32
	for n = 0; n < num; n++ {
		_ = m.AddCollectionByTypeId(typeMisc)
	}

	return nil
}

func (m *CollectionManager) LoadAll() error {
	res, err := store.GetStore().FindAll(context.Background(), define.StoreType_Collection, "owner_id", m.owner.ID)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if !utils.ErrCheck(err, "FindAll failed when CollectionManager.LoadAll", m.owner.ID) {
		return err
	}

	for _, v := range res {
		vv := v.([]byte)
		c := collection.NewCollection()
		c.Init()
		err := json.Unmarshal(vv, c)
		if !utils.ErrCheck(err, "Unmarshal failed when CollectionManager.LoadAll") {
			continue
		}

		if err := m.initLoadedCollection(c); err != nil {
			return fmt.Errorf("CollectionManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *CollectionManager) AfterLoad() {
}

func (m *CollectionManager) GetCollection(typeId int32) *collection.Collection {
	return m.CollectionList[typeId]
}

func (m *CollectionManager) GetCollectionNums() int {
	return len(m.CollectionList)
}

// 获取收集品可转换碎片的数量上限
func (m *CollectionManager) getCollectionTransformLimit(typeId int32) int32 {
	globalConfig, _ := auto.GetGlobalConfig()
	collectionEntry, ok := auto.GetCollectionEntry(typeId)
	if !ok {
		log.Warn().Int32("type_id", typeId).Msg("GetCollectionEntry failed")
		return 0
	}

	var curStar int32 = 0
	collection, ok := m.CollectionList[typeId]
	if ok {
		curStar = int32(collection.Star)
	}

	maxStar := globalConfig.CollectionMaxStar[collectionEntry.Quality]
	var totalFragments int32 = 0
	for n := curStar; n < maxStar; n++ {
		totalFragments += collectionEntry.StarCostFragments[n]
	}

	return totalFragments
}

func (m *CollectionManager) AddCollectionByTypeId(typeId int32) *collection.Collection {
	collectionEntry, ok := auto.GetCollectionEntry(typeId)
	if !ok {
		log.Warn().Int32("type_id", typeId).Msg("GetCollectionEntry failed")
		return nil
	}

	// send event
	defer func() {
		m.owner.eventManager.AddEvent(&event.Event{
			Type:  define.Event_Type_CollectionGain,
			Miscs: []interface{}{typeId},
		})
	}()

	// 重复获得收集品，转换为对应碎片
	_, ok = m.CollectionList[typeId]
	if ok {
		curNum := m.owner.FragmentManager().CollectionFragmentManager.FragmentList[typeId]
		limit := m.getCollectionTransformLimit(typeId)
		add := collectionEntry.FragmentTransform

		// 超过满星碎片，转换为对应品质通用碎片代币
		if collectionEntry.FragmentTransform+curNum > limit {
			add = limit - curNum
			_ = m.owner.TokenManager().TokenInc(define.Token_CollectionBegin+collectionEntry.Quality, collectionEntry.FragmentTransform-add)
		}

		m.owner.FragmentManager().CollectionFragmentManager.Inc(typeId, add)
		return nil
	}

	c := m.createEntryCollection(collectionEntry)
	if c == nil {
		log.Warn().Int32("type_id", typeId).Msg("createEntryCollection failed")
		return nil
	}

	filter := map[string]interface{}{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}
	err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Collection, filter, c)
	if !utils.ErrCheck(err, "UpdateOne failed when AddCollectionByTypeID", typeId, m.owner.ID) {
		m.delCollection(c)
		return nil
	}

	m.SendCollectionUpdate(c)

	return c
}

func (m *CollectionManager) delCollection(c *collection.Collection) {
	delete(m.CollectionList, c.TypeId)
	collection.GetCollectionPool().Put(c)
}

func (m *CollectionManager) DelCollection(typeId int32) {
	c, ok := m.CollectionList[typeId]
	if !ok {
		return
	}

	filter := map[string]interface{}{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}
	err := store.GetStore().DeleteOne(context.Background(), define.StoreType_Collection, filter)
	utils.ErrPrint(err, "DeleteOne failed when CollectionManager.DelCollection", typeId)
	m.delCollection(c)

	m.SendCollectionDelete(typeId)
}

func (m *CollectionManager) CollectionActive(typeId int32) error {
	c := m.GetCollection(typeId)
	if c == nil {
		return ErrCollectionNotFound
	}

	if c.Active {
		return ErrCollectionAlreadyActived
	}

	// todo cost

	c.Active = true

	// save
	filter := map[string]interface{}{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}

	fields := map[string]interface{}{
		"active": c.Active,
	}

	err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Collection, filter, fields)
	utils.ErrPrint(err, "UpdateOne failed when CollectionManager.CollectionActive", m.owner.ID, typeId)

	m.SendCollectionUpdate(c)
	return err
}

func (m *CollectionManager) CollectionStarup(typeId int32) error {
	c := m.GetCollection(typeId)
	if c == nil {
		return ErrCollectionNotFound
	}

	globalConfig, _ := auto.GetGlobalConfig()

	if int32(c.Star) >= globalConfig.CollectionMaxStar[c.Entry.Quality] {
		return ErrCollectionStarMax
	}

	nextStarFragments := c.Entry.StarCostFragments[c.Star]

	// 碎片不足
	err := m.owner.FragmentManager().CollectionFragmentManager.CanCost(c.TypeId, nextStarFragments)
	if err != nil {
		return err
	}

	err = m.owner.FragmentManager().CollectionFragmentManager.DoCost(c.TypeId, nextStarFragments)
	utils.ErrPrint(err, "CollectionStarup failed", m.owner.ID, c.TypeId, c.Star)

	c.Star++

	// save
	filter := map[string]interface{}{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}

	fields := map[string]interface{}{
		"star": c.Star,
	}
	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Collection, filter, fields)
	if !utils.ErrCheck(err, "UpdateFields failed when CollectionManager.CollectionStarup", m.owner.ID, c.Star) {
		return err
	}

	m.SendCollectionUpdate(c)
	return nil
}

// gm 升星
func (m *CollectionManager) GmCollectionStarup(typeId int32, star int32) error {
	c := m.GetCollection(typeId)
	if c == nil {
		return ErrCollectionNotFound
	}

	globalConfig, _ := auto.GetGlobalConfig()
	maxStar := globalConfig.CollectionMaxStar[c.Entry.Quality]
	if star > maxStar {
		star = maxStar
	}

	c.Star += int8(star)

	// save
	filter := map[string]interface{}{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}

	fields := map[string]interface{}{
		"star": c.Star,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Collection, filter, fields)
	utils.ErrPrint(err, "UpdateFields failed when CollectionManager.GmStarup")
	m.SendCollectionUpdate(c)

	return err
}

func (m *CollectionManager) SendCollectionUpdate(c *collection.Collection) {
	reply := &pbGlobal.S2C_CollectionInfo{
		Info: c.GenCollectionPB(),
	}

	m.owner.SendProtoMessage(reply)
}

func (m *CollectionManager) SendCollectionDelete(typeId int32) {
	// msg := &pbGlobal.S2C_DelCollection{
	// 	TypeId: typeId,
	// }
	// m.owner.SendProtoMessage(msg)
}

func (m *CollectionManager) GenCollectionListPB() []*pbGlobal.Collection {
	collections := make([]*pbGlobal.Collection, 0, len(m.CollectionList))
	for _, c := range m.CollectionList {
		collections = append(collections, c.GenCollectionPB())
	}

	return collections
}