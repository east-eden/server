package player

import (
	"errors"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/hero"
	pbCombat "github.com/yokaiio/yokai_server/proto/combat"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/store/db"
	"github.com/yokaiio/yokai_server/utils"
)

type HeroManager struct {
	owner   *Player
	mapHero map[int64]hero.Hero

	sync.RWMutex
}

func NewHeroManager(owner *Player) *HeroManager {
	m := &HeroManager{
		owner:   owner,
		mapHero: make(map[int64]hero.Hero, 0),
	}

	return m
}

func (m *HeroManager) TableName() string {
	return "hero"
}

func (m *HeroManager) createEntryHero(entry *define.HeroEntry) hero.Hero {
	if entry == nil {
		logger.Error("newEntryHero with nil HeroEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Hero)
	if err != nil {
		logger.Error(err)
		return nil
	}

	h := hero.NewHero(
		hero.Id(id),
		hero.OwnerId(m.owner.GetID()),
		hero.OwnerType(m.owner.GetType()),
		hero.Entry(entry),
		hero.TypeId(entry.ID),
	)

	h.GetAttManager().SetBaseAttId(entry.AttID)
	m.mapHero[h.GetOptions().Id] = h
	store.GetStore().SaveObject(store.StoreType_Hero, h)

	h.GetAttManager().CalcAtt()

	return h
}

func (m *HeroManager) initLoadedHero(h hero.Hero) error {
	entry := entries.GetHeroEntry(h.GetOptions().TypeId)

	if entry == nil {
		return fmt.Errorf("HeroManager initLoadedHero: hero<%d> entry invalid", h.GetOptions().TypeId)
	}

	h.GetOptions().Entry = entry
	h.GetAttManager().SetBaseAttId(entry.AttID)

	m.mapHero[h.GetOptions().Id] = h
	h.CalcAtt()
	return nil
}

// interface of cost_loot
func (m *HeroManager) GetCostLootType() int32 {
	return define.CostLoot_Hero
}

func (m *HeroManager) CanCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("hero manager check hero<%d> cost failed, wrong number<%d>", typeMisc, num)
	}

	var fixNum int32 = 0
	for _, v := range m.mapHero {
		if v.GetOptions().TypeId == typeMisc {
			eb := v.GetEquipBar()
			hasEquip := false
			var n int32
			for n = 0; n < define.Hero_MaxEquip; n++ {
				if eb.GetEquipByPos(n) != nil {
					hasEquip = true
					break
				}
			}

			if !hasEquip {
				fixNum++
			}
		}
	}

	if fixNum >= num {
		return nil
	}

	return fmt.Errorf("not enough hero<%d>, num<%d>", typeMisc, num)
}

func (m *HeroManager) DoCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("hero manager cost hero<%d> failed, wrong number<%d>", typeMisc, num)
	}

	var costNum int32 = 0
	for _, v := range m.mapHero {
		if v.GetOptions().TypeId == typeMisc {
			eb := v.GetEquipBar()
			hasEquip := false
			var n int32
			for n = 0; n < define.Hero_MaxEquip; n++ {
				if eb.GetEquipByPos(n) != nil {
					hasEquip = true
					break
				}
			}

			if !hasEquip {
				m.DelHero(v.GetOptions().Id)
				costNum++
			}
		}
	}

	if costNum < num {
		logger.WithFields(logger.Fields{
			"cost_type_misc":  typeMisc,
			"cost_num":        num,
			"actual_cost_num": costNum,
		}).Warn("hero manager cost num error")
		return nil
	}

	return nil
}

func (m *HeroManager) CanGain(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("hero manager check hero<%d> gain failed, wrong number<%d>", typeMisc, num)
	}

	// todo max hero num
	return nil
}

func (m *HeroManager) GainLoot(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("hero manager gain hero<%d> failed, wrong number<%d>", typeMisc, num)
	}

	var n int32 = 0
	for ; n < num; n++ {
		h := m.AddHeroByTypeID(typeMisc)
		if h == nil {
			return fmt.Errorf("hero manager gain hero<%d> failed, cannot add new hero<%d>", typeMisc, num)
		}
	}

	return nil
}

func (m *HeroManager) LoadAll() error {
	heroList, err := store.GetStore().LoadArray(store.StoreType_Hero, "owner_id", m.owner.GetID(), hero.GetHeroPool())
	if errors.Is(err, db.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("HeroManager LoadAll: %w", err)
	}

	for _, i := range heroList {
		if err := m.initLoadedHero(i.(hero.Hero)); err != nil {
			return fmt.Errorf("HeroManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *HeroManager) GetHero(id int64) hero.Hero {
	return m.mapHero[id]
}

func (m *HeroManager) GetHeroNums() int {
	return len(m.mapHero)
}

func (m *HeroManager) GetHeroList() []hero.Hero {
	list := make([]hero.Hero, 0)

	m.RLock()
	for _, v := range m.mapHero {
		list = append(list, v)
	}
	m.RUnlock()

	return list
}

func (m *HeroManager) AddHeroByTypeID(typeID int32) hero.Hero {
	heroEntry := entries.GetHeroEntry(typeID)
	h := m.createEntryHero(heroEntry)
	if h == nil {
		return nil
	}

	return h
}

func (m *HeroManager) DelHero(id int64) {
	h, ok := m.mapHero[id]
	if !ok {
		return
	}

	eb := h.GetEquipBar()
	var n int32
	for n = 0; n < define.Hero_MaxEquip; n++ {
		eb.TakeoffEquip(n)
	}
	h.BeforeDelete()

	delete(m.mapHero, id)
	store.GetStore().DeleteObject(store.StoreType_Hero, h)
	hero.ReleasePoolHero(h)
}

func (m *HeroManager) HeroSetLevel(level int32) {
	for _, v := range m.mapHero {
		v.GetOptions().Level = level

		fields := map[string]interface{}{
			"level": v.GetOptions().Level,
		}
		store.GetStore().SaveFields(store.StoreType_Hero, v, fields)
	}
}

func (m *HeroManager) PutonEquip(heroID int64, equipID int64) error {

	equip := m.owner.ItemManager().GetItem(equipID)
	if equip == nil {
		return fmt.Errorf("cannot find equip<%d> while PutonEquip", equipID)
	}

	if objId := equip.GetOptions().EquipObj; objId != -1 {
		return fmt.Errorf("equip has put on another hero<%d>", objId)
	}

	if equip.EquipEnchantEntry() == nil {
		return fmt.Errorf("cannot find equip_enchant_entry<%d> while PutonEquip", equipID)
	}

	h, ok := m.mapHero[heroID]
	if !ok {
		return fmt.Errorf("invalid heroid")
	}

	equipBar := h.GetEquipBar()
	pos := equip.EquipEnchantEntry().EquipPos

	// takeoff previous equip
	if pe := equipBar.GetEquipByPos(pos); pe != nil {
		if err := m.TakeoffEquip(heroID, pos); err != nil {
			return err
		}
	}

	// puton this equip
	if err := equipBar.PutonEquip(equip); err != nil {
		return err
	}

	m.owner.ItemManager().Save(equip.GetOptions().Id)
	m.owner.ItemManager().SendItemUpdate(equip)
	m.SendHeroUpdate(h)

	// att
	equip.GetAttManager().CalcAtt()
	h.GetAttManager().ModAttManager(equip.GetAttManager())
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) TakeoffEquip(heroID int64, pos int32) error {
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return fmt.Errorf("invalid pos")
	}

	h, ok := m.mapHero[heroID]
	if !ok {
		return fmt.Errorf("invalid heroid")
	}

	equipBar := h.GetEquipBar()
	equip := equipBar.GetEquipByPos(pos)
	if equip == nil {
		return fmt.Errorf("cannot find hero<%d> equip by pos<%d> while TakeoffEquip", heroID, pos)
	}

	if objID := equip.GetEquipObj(); objID == -1 {
		return fmt.Errorf("equip<%d> didn't put on this hero<%d> ", equip.GetOptions().Id, heroID)
	}

	// unequip
	if err := equipBar.TakeoffEquip(pos); err != nil {
		return err
	}

	m.owner.ItemManager().Save(equip.GetOptions().Id)
	m.owner.ItemManager().SendItemUpdate(equip)
	m.SendHeroUpdate(h)

	// att
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) PutonRune(heroId int64, runeId int64) error {

	r := m.owner.RuneManager().GetRune(runeId)
	if r == nil {
		return fmt.Errorf("cannot find rune<%d> while PutonRune", runeId)
	}

	if objId := r.GetOptions().EquipObj; objId != -1 {
		return fmt.Errorf("rune has put on another obj<%d>", objId)
	}

	pos := r.GetOptions().Entry.Pos
	if pos < define.Rune_PositionBegin || pos >= define.Rune_PositionEnd {
		return fmt.Errorf("invalid pos<%d>", pos)
	}

	h, ok := m.mapHero[heroId]
	if !ok {
		return fmt.Errorf("invalid heroid<%d>", heroId)
	}

	runeBox := h.GetRuneBox()

	// takeoff previous rune
	if pr := runeBox.GetRuneByPos(pos); pr != nil {
		if err := m.TakeoffRune(heroId, pos); err != nil {
			return err
		}
	}

	// equip new rune
	if err := runeBox.PutonRune(r); err != nil {
		return err
	}

	m.owner.RuneManager().Save(runeId)
	m.owner.RuneManager().SendRuneUpdate(r)
	m.SendHeroUpdate(h)

	// att
	r.GetAttManager().CalcAtt()
	h.GetAttManager().ModAttManager(r.GetAttManager())
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) TakeoffRune(heroId int64, pos int32) error {
	if pos < 0 || pos >= define.Rune_PositionEnd {
		return fmt.Errorf("invalid pos<%d>", pos)
	}

	h, ok := m.mapHero[heroId]
	if !ok {
		return fmt.Errorf("invalid heroid<%d>", heroId)
	}

	r := h.GetRuneBox().GetRuneByPos(pos)
	if r == nil {
		return fmt.Errorf("cannot find rune from hero<%d>'s runebox pos<%d> while TakeoffRune", heroId, pos)
	}

	// unequip
	if err := h.GetRuneBox().TakeoffRune(pos); err != nil {
		return err
	}

	m.owner.RuneManager().Save(r.GetOptions().Id)
	m.owner.RuneManager().SendRuneUpdate(r)
	m.SendHeroUpdate(h)

	// att
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) GenerateCombatUnitInfo() []*pbCombat.UnitInfo {
	retList := make([]*pbCombat.UnitInfo, 0)

	list := m.GetHeroList()
	for _, hero := range list {
		unitInfo := &pbCombat.UnitInfo{
			UnitTypeId: hero.GetOptions().TypeId,
		}

		for n := define.Att_Begin; n < define.Att_End; n++ {
			unitInfo.UnitAttList = append(unitInfo.UnitAttList, &pbGame.Att{
				AttType:  int32(n),
				AttValue: hero.GetAttManager().GetAttValue(int32(n)),
			})
		}

		retList = append(retList, unitInfo)
	}

	return retList
}

func (m *HeroManager) SendHeroUpdate(h hero.Hero) {
	// send equips update
	reply := &pbGame.M2C_HeroInfo{
		Info: &pbGame.Hero{
			Id:     h.GetOptions().Id,
			TypeId: h.GetOptions().TypeId,
			Exp:    h.GetOptions().Exp,
			Level:  h.GetOptions().Level,
		},
	}

	// equip list
	eb := h.GetEquipBar()
	var n int32
	for n = 0; n < define.Hero_MaxEquip; n++ {
		var equipId int64 = -1
		if i := eb.GetEquipByPos(n); i != nil {
			equipId = i.GetOptions().Id
		}

		reply.Info.EquipList = append(reply.Info.EquipList, equipId)
	}

	// rune list
	var pos int32
	for pos = 0; pos < define.Rune_PositionEnd; pos++ {
		var runeId int64 = -1
		if r := h.GetRuneBox().GetRuneByPos(pos); r != nil {
			runeId = r.GetOptions().Id
		}

		reply.Info.RuneList = append(reply.Info.RuneList, runeId)
	}

	m.owner.SendProtoMessage(reply)
}

func (m *HeroManager) SendHeroAtt(h hero.Hero) {
	attManager := h.GetAttManager()
	reply := &pbGame.M2C_HeroAttUpdate{
		HeroId: h.GetOptions().Id,
	}

	for k := int32(0); k < define.Att_End; k++ {
		att := &pbGame.Att{
			AttType:  k,
			AttValue: attManager.GetAttValue(k),
		}
		reply.AttList = append(reply.AttList, att)
	}

	m.owner.SendProtoMessage(reply)
}
