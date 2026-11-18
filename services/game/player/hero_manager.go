package player

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/event"
	"github.com/east-eden/server/services/game/hero"
	"github.com/east-eden/server/services/game/prom"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

var (
	heroUpdateInterval = time.Second * 5

	ErrHeroNotFound      = errors.New("hero not found")
	ErrHeroMaxLevel      = errors.New("hero reach max level")
	ErrHeroMaxStar       = errors.New("hero reach max star")
	ErrHeroMaxPromotion  = errors.New("hero reach max promotion")
	ErrHeroChipNotEnough = errors.New("hero chip not enough")
	ErrHeroRepeated      = errors.New("hero already owned")
)

// HeroManager 英雄管理器。
//
// 采用 kof 的卡牌模型：一个英雄配置 id 对应一张卡，存放在 map 里，
// 等级/星级/晋升决定属性，属性与战力由 Card.Recalc 统一重算。
type HeroManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	Heroes     map[int32]*hero.Card `bson:"heroes" json:"heroes"` // 英雄配置 id -> 卡牌
	nextUpdate int64                `bson:"-" json:"-"`
	owner      *Player              `bson:"-" json:"-"`
}

func NewHeroManager(owner *Player) *HeroManager {
	return &HeroManager{
		owner:  owner,
		Heroes: make(map[int32]*hero.Card),
	}
}

func (m *HeroManager) Destroy() {
	m.Heroes = make(map[int32]*hero.Card)
}

func (m *HeroManager) RegisterEvent() {
	m.owner.eventManager.Register(define.Event_Type_PlayerLevelup, m.onEventHeroLevelup)
}

// onEventHeroLevelup 玩家升级后英雄的 playerLevel 相关属性需要重算
func (m *HeroManager) onEventHeroLevelup(e *event.Event) error {
	for _, card := range m.Heroes {
		if err := card.Recalc(m.owner.GetLevel()); err != nil {
			utils.ErrPrint(err, "hero recalc failed", card.ConfigId)
			continue
		}
	}

	return nil
}

////////////////////////////////////////////////////////
// CostLooter

func (m *HeroManager) GetCostLootType() int32 {
	return define.CostLoot_Hero
}

func (m *HeroManager) CanCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return nil
	}

	if _, ok := m.Heroes[typeMisc]; !ok {
		return fmt.Errorf("%w: config<%d>", ErrHeroNotFound, typeMisc)
	}

	return nil
}

func (m *HeroManager) DoCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return nil
	}

	card, ok := m.Heroes[typeMisc]
	if !ok {
		return fmt.Errorf("%w: config<%d>", ErrHeroNotFound, typeMisc)
	}

	return m.removeHero(card)
}

func (m *HeroManager) GainLoot(typeMisc int32, num int32) error {
	for i := int32(0); i < num; i++ {
		if _, err := m.AddHeroByTypeId(typeMisc); err != nil {
			return err
		}
	}

	return nil
}

////////////////////////////////////////////////////////
// 生命周期

func (m *HeroManager) LoadAll() error {
	res, err := store.GetStore().FindAll(context.Background(), define.StoreType_Hero, "owner_id", m.owner.ID)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if !utils.ErrCheck(err, "FindAll failed when HeroManager.LoadAll", m.owner.ID) {
		return err
	}

	m.Heroes = make(map[int32]*hero.Card, len(res))

	for _, v := range res {
		raw, ok := v.([]byte)
		if !ok {
			continue
		}

		var card hero.Card
		if err := json.Unmarshal(raw, &card); err != nil {
			return err
		}

		if _, ok := config.GetHeroEntry(card.ConfigId); !ok {
			log.Warn().Int32("config_id", card.ConfigId).Int64("owner_id", m.owner.ID).Msg("hero config missing, dropped")
			continue
		}

		if err := card.Recalc(m.owner.GetLevel()); err != nil {
			utils.ErrPrint(err, "hero recalc failed on load", card.ConfigId)
		}

		m.Heroes[card.ConfigId] = &card
	}

	return nil
}

func (m *HeroManager) AfterLoad() {
	for _, card := range m.Heroes {
		if err := card.Recalc(m.owner.GetLevel()); err != nil {
			utils.ErrPrint(err, "hero recalc failed after load", card.ConfigId)
		}
	}
}

func (m *HeroManager) update() {}

func (m *HeroManager) save(card *hero.Card) error {
	return store.GetStore().UpdateOne(context.Background(), define.StoreType_Hero, int64(card.ConfigId), card)
}

////////////////////////////////////////////////////////
// 查询

func (m *HeroManager) GetHero(id int64) *hero.Card {
	return m.Heroes[int32(id)]
}

func (m *HeroManager) GetHeroByTypeId(typeId int32) *hero.Card {
	return m.Heroes[typeId]
}

// GetHeroList 按配置 id 升序返回所有英雄
func (m *HeroManager) GetHeroList() []*hero.Card {
	list := make([]*hero.Card, 0, len(m.Heroes))
	for _, card := range m.Heroes {
		list = append(list, card)
	}

	sort.Slice(list, func(i, j int) bool { return list[i].ConfigId < list[j].ConfigId })
	return list
}

////////////////////////////////////////////////////////
// 养成

// AddHeroByTypeId 获得英雄。已拥有时按配置转换成英雄碎片。
func (m *HeroManager) AddHeroByTypeId(typeId int32) (*hero.Card, error) {
	entry, ok := config.GetHeroEntry(typeId)
	if !ok {
		return nil, fmt.Errorf("hero config<%d> not found", typeId)
	}

	if card, ok := m.Heroes[typeId]; ok {
		// 重复获得：按品质配置转换成碎片（kof chip 模型，碎片直接计数）
		if quality, ok := config.GetHeroQualityEntry(entry.Rarity); ok && quality.ConvertChipCnt > 0 {
			if err := m.AddHeroChip(entry.ChipId, int64(quality.ConvertChipCnt)); err != nil {
				return card, err
			}
		}

		return card, nil
	}

	card, err := hero.NewCard(typeId)
	if err != nil {
		return nil, err
	}

	if err := card.Recalc(m.owner.GetLevel()); err != nil {
		return nil, err
	}

	m.Heroes[typeId] = card

	if err := m.save(card); err != nil {
		return nil, err
	}

	prom.OpsCreateHeroCounter.Inc()
	m.SendHeroUpdate(card)

	m.owner.eventManager.DispatchEvent(&event.Event{
		Type:  define.Event_Type_HeroGain,
		Miscs: []any{typeId},
	})

	return card, nil
}

func (m *HeroManager) removeHero(card *hero.Card) error {
	delete(m.Heroes, card.ConfigId)
	m.SendHeroDelete(int64(card.ConfigId))

	return store.GetStore().DeleteOne(context.Background(), define.StoreType_Hero, int64(card.ConfigId))
}

func (m *HeroManager) DelHero(id int64) {
	card, ok := m.Heroes[int32(id)]
	if !ok {
		return
	}

	_ = m.removeHero(card)
}

// HeroLevelup 升到目标等级，逐级消耗 HeroLevelNew 配置的道具
func (m *HeroManager) HeroLevelup(heroId int64, targetLevel int32) error {
	card, ok := m.Heroes[int32(heroId)]
	if !ok {
		return ErrHeroNotFound
	}

	if targetLevel <= card.Level {
		return nil
	}

	entry, ok := card.Entry()
	if !ok {
		return ErrHeroNotFound
	}

	itemMgr := m.owner.ItemManager()

	for lv := card.Level; lv < targetLevel; lv++ {
		levelEntry, ok := config.GetHeroLevelNewEntry(entry.Rarity.String(), fmt.Sprintf("%d", lv))
		if !ok {
			return ErrHeroMaxLevel
		}

		for _, cost := range levelEntry.NextCost {
			if err := itemMgr.OnCheckEnough(cost.Id, int64(cost.Count)); err != nil {
				return err
			}
		}

		for _, cost := range levelEntry.NextCost {
			if err := itemMgr.OnItemCost(cost.Id, int64(cost.Count), "heroLevelup"); err != nil {
				return err
			}
		}
	}

	card.Level = targetLevel

	if err := card.Recalc(m.owner.GetLevel()); err != nil {
		return err
	}

	if err := m.save(card); err != nil {
		return err
	}

	m.SendHeroLevelup(card)

	m.owner.eventManager.DispatchEvent(&event.Event{
		Type:  define.Event_Type_HeroLevelup,
		Miscs: []any{card.ConfigId, card.Level},
	})

	return nil
}

// HeroStarup 升一星
func (m *HeroManager) HeroStarup(heroId int64) error {
	card, ok := m.Heroes[int32(heroId)]
	if !ok {
		return ErrHeroNotFound
	}

	entry, ok := card.Entry()
	if !ok {
		return ErrHeroNotFound
	}

	nextStar := card.Star + 1
	starEntry, ok := config.GetHeroStarEntry(fmt.Sprintf("%d", nextStar), entry.Rarity.String())
	if !ok {
		return ErrHeroMaxStar
	}

	// 升星消耗英雄碎片（heroEntry.ChipId 即碎片物品 id）
	if starEntry.Cost > 0 {
		if err := m.CostHeroChip(entry.ChipId, starEntry.Cost); err != nil {
			return err
		}
	}

	card.Star = nextStar

	if err := card.Recalc(m.owner.GetLevel()); err != nil {
		return err
	}

	if err := m.save(card); err != nil {
		return err
	}

	m.SendHeroStarup(card)

	m.owner.eventManager.DispatchEvent(&event.Event{
		Type:  define.Event_Type_HeroStarup,
		Miscs: []any{card.ConfigId, card.Star},
	})

	return nil
}

// HeroPromote 晋升一档品质
func (m *HeroManager) HeroPromote(heroId int64) error {
	card, ok := m.Heroes[int32(heroId)]
	if !ok {
		return ErrHeroNotFound
	}

	entry, ok := card.Entry()
	if !ok {
		return ErrHeroNotFound
	}

	rarity := config.Rarity(card.Rarity)
	if int32(entry.Rarity) > card.Rarity {
		rarity = entry.Rarity
	}

	curEntry, ok := config.GetHeroQualityPromotionEntry(rarity.String(), fmt.Sprintf("%d", card.PromotionLv))
	if !ok {
		return ErrHeroMaxPromotion
	}

	nextEntry, ok := config.GetNextQualityPromotionEntry(rarity, card.PromotionLv)
	if !ok {
		return ErrHeroMaxPromotion
	}

	itemMgr := m.owner.ItemManager()

	for _, cost := range curEntry.PromotionNextCost {
		if err := itemMgr.OnCheckEnough(cost.Id, int64(cost.Count)); err != nil {
			return err
		}
	}

	for _, cost := range curEntry.PromotionNextCost {
		if err := itemMgr.OnItemCost(cost.Id, int64(cost.Count), "heroPromotion"); err != nil {
			return err
		}
	}

	card.Rarity = int32(nextEntry.Id)
	card.PromotionLv = nextEntry.Level

	if err := card.Recalc(m.owner.GetLevel()); err != nil {
		return err
	}

	if err := m.save(card); err != nil {
		return err
	}

	m.SendHeroPromote(card)

	m.owner.eventManager.DispatchEvent(&event.Event{
		Type:  define.Event_Type_HeroPromotion,
		Miscs: []any{card.ConfigId, card.PromotionLv},
	})

	return nil
}

////////////////////////////////////////////////////////
// GM

func (m *HeroManager) GmExpChange(heroId int64, exp int32) error {
	card, ok := m.Heroes[int32(heroId)]
	if !ok {
		return ErrHeroNotFound
	}

	card.Exp = exp
	return m.save(card)
}

func (m *HeroManager) GmLevelChange(heroId int64, level int32) error {
	card, ok := m.Heroes[int32(heroId)]
	if !ok {
		return ErrHeroNotFound
	}

	card.Level = level
	if err := card.Recalc(m.owner.GetLevel()); err != nil {
		return err
	}

	return m.save(card)
}

func (m *HeroManager) GmPromoteChange(heroId int64, promote int32) error {
	card, ok := m.Heroes[int32(heroId)]
	if !ok {
		return ErrHeroNotFound
	}

	card.PromotionLv = promote
	if err := card.Recalc(m.owner.GetLevel()); err != nil {
		return err
	}

	return m.save(card)
}

////////////////////////////////////////////////////////
// 协议

func (m *HeroManager) GenCombatEntityInfo() []*pbGlobal.EntityInfo {
	list := make([]*pbGlobal.EntityInfo, 0, len(m.Heroes))
	for _, card := range m.GetHeroList() {
		list = append(list, &pbGlobal.EntityInfo{
			HeroTypeId: card.ConfigId,
			AttValue:   exportAttValues(card),
		})
	}

	return list
}

func exportAttValues(card *hero.Card) []float32 {
	const attrMax = 512

	values := make([]float32, attrMax)
	for attrType, v := range card.Attrs {
		if attrType < 0 || attrType >= attrMax {
			continue
		}

		values[attrType] = v
	}

	return values
}

func (m *HeroManager) cardPB(card *hero.Card) *pbGlobal.Hero {
	return &pbGlobal.Hero{
		Id:           int64(card.ConfigId),
		TypeId:       card.ConfigId,
		Exp:          card.Exp,
		Level:        card.Level,
		PromoteLevel: card.PromotionLv,
		Star:         card.Star,
		AttValues:    exportAttValues(card),
	}
}

func (m *HeroManager) SendHeroUpdate(card *hero.Card) {
	m.owner.SendProtoMessage(&pbGlobal.S2C_HeroInfo{Info: m.cardPB(card)})
}

func (m *HeroManager) SendHeroLevelup(card *hero.Card) {
	m.SendHeroUpdate(card)
}

func (m *HeroManager) SendHeroStarup(card *hero.Card) {
	m.SendHeroUpdate(card)
}

func (m *HeroManager) SendHeroPromote(card *hero.Card) {
	m.SendHeroUpdate(card)
}

func (m *HeroManager) SendHeroDelete(id int64) {
	m.owner.SendProtoMessage(&pbGlobal.S2C_DelHero{Id: id})
}

func (m *HeroManager) SendHeroAttUpdate(card *hero.Card) {
	m.SendHeroUpdate(card)
}

func (m *HeroManager) GenHeroListPB() []*pbGlobal.Hero {
	list := make([]*pbGlobal.Hero, 0, len(m.Heroes))
	for _, card := range m.GetHeroList() {
		list = append(list, m.cardPB(card))
	}

	return list
}

////////////////////////////////////////////////////////
// 英雄碎片（kof chip 模型）

// HeroChips 返回碎片计数表（保证非 nil）
func (m *HeroManager) HeroChips() map[int32]int32 {
	if m.owner.HeroChips == nil {
		m.owner.HeroChips = make(map[int32]int32)
	}

	return m.owner.HeroChips
}

// GetHeroChipCount 返回某碎片数量
func (m *HeroManager) GetHeroChipCount(chipId int32) int32 {
	return m.HeroChips()[chipId]
}

// AddHeroChip 增加碎片
func (m *HeroManager) AddHeroChip(chipId int32, num int64) error {
	if chipId <= 0 || num <= 0 {
		return nil
	}

	m.HeroChips()[chipId] += int32(num)

	if err := m.saveChips(); err != nil {
		return err
	}

	m.SendHeroChipUpdate()
	return nil
}

// CostHeroChip 消耗碎片
func (m *HeroManager) CostHeroChip(chipId int32, num int32) error {
	if num <= 0 {
		return nil
	}

	chips := m.HeroChips()
	if chips[chipId] < num {
		return fmt.Errorf("%w: chip<%d> has %d, need %d", ErrHeroChipNotEnough, chipId, chips[chipId], num)
	}

	chips[chipId] -= num
	if chips[chipId] == 0 {
		delete(chips, chipId)
	}

	if err := m.saveChips(); err != nil {
		return err
	}

	m.SendHeroChipUpdate()
	return nil
}

// ComposeHero 用碎片合成英雄
func (m *HeroManager) ComposeHero(heroConfigId int32) error {
	entry, ok := config.GetHeroEntry(heroConfigId)
	if !ok {
		return fmt.Errorf("hero config<%d> not found", heroConfigId)
	}

	if _, exists := m.Heroes[heroConfigId]; exists {
		return fmt.Errorf("%w: config<%d> already owned", ErrHeroRepeated, heroConfigId)
	}

	quality, ok := config.GetHeroQualityEntry(entry.Rarity)
	if !ok {
		return fmt.Errorf("hero quality<%s> not found", entry.Rarity.String())
	}

	if quality.ComposeChipCnt <= 0 {
		return fmt.Errorf("hero<%d> has no compose chip cost", heroConfigId)
	}

	if err := m.CostHeroChip(entry.ChipId, quality.ComposeChipCnt); err != nil {
		return err
	}

	_, err := m.AddHeroByTypeId(heroConfigId)
	return err
}

func (m *HeroManager) saveChips() error {
	return store.GetStore().UpdateFields(
		context.Background(),
		define.StoreType_Player,
		m.owner.ID,
		map[string]any{"hero_chips": m.HeroChips()},
	)
}

func (m *HeroManager) SendHeroChipUpdate() {
	m.owner.SendProtoMessage(&pbGlobal.S2C_HeroChipUpdate{Chips: m.HeroChips()})
}
