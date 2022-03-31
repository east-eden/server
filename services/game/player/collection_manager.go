package player

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/collection"
	"github.com/east-eden/server/services/game/event"
	"github.com/east-eden/server/services/game/quest"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	ErrCollectionNotFound        = errors.New("collection not found")
	ErrCollectionStarMax         = errors.New("collection star max")
	ErrCollectionAlreadyActived  = errors.New("collection already actived")
	ErrCollectionAlreadyWakeuped = errors.New("collection already wakeuped")
)

type CollectionManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	owner          *Player                          `bson:"-" json:"-"`
	CollectionList map[int32]*collection.Collection `bson:"-" json:"-"` // 收集品列表
}

func NewCollectionManager(owner *Player) *CollectionManager {
	m := &CollectionManager{
		owner:          owner,
		CollectionList: make(map[int32]*collection.Collection),
	}

	return m
}

func (m *CollectionManager) Destroy() {
	for _, c := range m.CollectionList {
		collection.GetCollectionPool().Put(c)
	}
}

func (m *CollectionManager) createEntryCollection(entry *auto.CollectionEntry) *collection.Collection {
	if entry == nil {
		log.Error().Msg("newEntryCollection with nil CollectionEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Collection)
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	c := collection.NewCollection()
	c.Init(
		collection.Id(id),
		collection.TypeId(entry.Id),
		collection.OwnerId(m.owner.GetId()),
		collection.Entry(entry),
		collection.EventManager(m.owner.EventManager()),
		collection.QuestUpdateCb(func(q *quest.Quest) {
			m.SendCollectionUpdate(c)
		}),
	)
	c.InitQuestManager()

	m.CollectionList[c.GetOptions().TypeId] = c

	return c
}

func (m *CollectionManager) initLoadedCollection(c *collection.Collection) error {
	entry, ok := auto.GetCollectionEntry(c.TypeId)
	if !ok {
		return fmt.Errorf("CollectionManager initLoadedCollection: collection<%d> entry invalid", c.GetOptions().TypeId)
	}

	c.Entry = entry
	c.InitQuestManager()

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

		c.Init(
			collection.EventManager(m.owner.EventManager()),
			collection.QuestUpdateCb(func(q *quest.Quest) {
				m.SendCollectionUpdate(c)
			}),
		)

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

func (m *CollectionManager) GetCollection(typeId int32) *collection.Collection {
	return m.CollectionList[typeId]
}

// 收集品满星需要多少碎片
func (m *CollectionManager) GetCollectionMaxStarNeedFragments(typeId int32) int32 {
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
			Miscs: []any{typeId},
		})
	}()

	// 重复获得收集品，转换为对应碎片
	c, ok := m.CollectionList[typeId]
	if ok {
		_ = m.owner.FragmentManager().CollectionFragmentManager.GainLoot(typeId, collectionEntry.FragmentTransform)
		return c
	}

	c = m.createEntryCollection(collectionEntry)
	if c == nil {
		log.Warn().Int32("type_id", typeId).Msg("createEntryCollection failed")
		return nil
	}

	err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Collection, c.Id, c)
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

	filter := map[string]any{
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

	err := m.owner.CostLootManager().CanCost(c.Entry.ActiveCostId)
	if err != nil {
		return err
	}

	err = m.owner.CostLootManager().DoCost(c.Entry.ActiveCostId)
	utils.ErrPrint(err, "DoCost failed when CollectionManager.CollectionActive", m.owner.ID, typeId)

	c.Active = true
	c.CalcScore()

	fields := map[string]any{
		"active": c.Active,
	}

	err = store.GetStore().UpdateOne(context.Background(), define.StoreType_Collection, c.Id, fields)
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
	c.CalcScore()

	// save
	filter := map[string]any{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}

	fields := map[string]any{
		"star": c.Star,
	}
	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Collection, filter, fields)
	if !utils.ErrCheck(err, "UpdateFields failed when CollectionManager.CollectionStarup", m.owner.ID, c.Star) {
		return err
	}

	m.SendCollectionUpdate(c)
	return nil
}

func (m *CollectionManager) CollectionWakeup(typeId int32) error {
	c := m.GetCollection(typeId)
	if c == nil {
		return ErrCollectionNotFound
	}

	if c.Wakeup {
		return ErrCollectionAlreadyWakeuped
	}

	err := m.owner.CostLootManager().CanCost(c.Entry.WakeupCostId)
	if err != nil {
		return err
	}

	err = m.owner.CostLootManager().DoCost(c.Entry.WakeupCostId)
	utils.ErrPrint(err, "DoCost failed when CollectionManager.CollectionWakeup", m.owner.ID, typeId)

	c.Wakeup = true
	c.CalcScore()

	// save
	filter := map[string]any{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}

	fields := map[string]any{
		"wakeup": c.Wakeup,
	}
	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Collection, filter, fields)
	if !utils.ErrCheck(err, "UpdateFields failed when CollectionManager.CollectionWakeup", m.owner.ID, c.TypeId) {
		return err
	}

	m.SendCollectionUpdate(c)
	return nil
}

// gm 激活
func (m *CollectionManager) GmCollectionActive(typeId int32) error {
	c := m.GetCollection(typeId)
	if c == nil {
		return ErrCollectionNotFound
	}

	if c.Active {
		return ErrCollectionAlreadyActived
	}

	c.Active = true
	c.CalcScore()

	// save
	filter := map[string]any{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}

	fields := map[string]any{
		"active": c.Active,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Collection, filter, fields)
	utils.ErrPrint(err, "UpdateFields failed when CollectionManager.GmActive")
	m.SendCollectionUpdate(c)
	return nil
}

// gm 觉醒
func (m *CollectionManager) GmCollectionWakeup(typeId int32) error {
	c := m.GetCollection(typeId)
	if c == nil {
		return ErrCollectionNotFound
	}

	if c.Wakeup {
		return ErrCollectionAlreadyWakeuped
	}

	c.Wakeup = true
	c.CalcScore()

	// save
	filter := map[string]any{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}

	fields := map[string]any{
		"wakeup": c.Wakeup,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Collection, filter, fields)
	utils.ErrPrint(err, "UpdateFields failed when CollectionManager.GmWakeup")
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
	c.CalcScore()

	// save
	filter := map[string]any{
		"type_id":  c.TypeId,
		"owner_id": c.OwnerId,
	}

	fields := map[string]any{
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
