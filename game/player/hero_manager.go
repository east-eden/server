package player

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/rune"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HeroManager struct {
	owner   *Player
	mapHero map[int64]hero.Hero

	ctx    context.Context
	cancel context.CancelFunc
	ds     *db.Datastore
	coll   *mongo.Collection
	sync.RWMutex
}

func NewHeroManager(ctx context.Context, owner *Player, ds *db.Datastore) *HeroManager {
	m := &HeroManager{
		owner:   owner,
		ds:      ds,
		mapHero: make(map[int64]hero.Hero, 0),
	}

	m.ctx, m.cancel = context.WithCancel(ctx)
	m.coll = ds.Database().Collection(m.TableName())

	return m
}

func (m *HeroManager) TableName() string {
	return "hero"
}

func (m *HeroManager) save(h hero.Hero) {
	filter := bson.D{{"_id", h.GetID()}}
	update := bson.D{{"$set", h}}
	op := options.Update().SetUpsert(true)
	if _, err := m.coll.UpdateOne(m.ctx, filter, update, op); err != nil {
		logger.WithFields(logger.Fields{
			"id":    h.GetID(),
			"error": err,
		}).Warning("hero save failed")
	}
}

func (m *HeroManager) saveField(h hero.Hero, up *bson.D) {
	filter := bson.D{{"_id", h.GetID()}}
	if _, err := m.coll.UpdateOne(m.ctx, filter, *up); err != nil {
		logger.WithFields(logger.Fields{
			"id":    h.GetID(),
			"level": h.GetLevel(),
			"error": err,
		}).Warning("hero save level failed")
	}
}

func (m *HeroManager) delete(h hero.Hero, filter *bson.D) {
	if _, err := m.coll.DeleteOne(m.ctx, *filter); err != nil {
		logger.WithFields(logger.Fields{
			"id":    h.GetID(),
			"error": err,
		}).Warning("hero delete level failed")
	}
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

	h := hero.NewHero(id)
	h.SetOwnerID(m.owner.GetID())
	h.SetOwnerType(m.owner.GetType())
	h.SetLevel(m.owner.GetLevel())
	h.SetTypeID(entry.ID)
	h.SetEntry(entry)

	attManager := att.NewAttManager(entry.AttID)
	h.SetAttManager(attManager)

	equipBar := item.NewEquipBar(h)
	h.SetEquipBar(equipBar)

	runeBox := rune.NewRuneBox(h)
	h.SetRuneBox(runeBox)

	m.mapHero[h.GetID()] = h

	return h
}

func (m *HeroManager) createDBHero(h hero.Hero) hero.Hero {
	newHero := hero.NewHero(h.GetID())
	newHero.SetOwnerID(h.GetOwnerID())
	newHero.SetOwnerType(h.GetOwnerType())
	newHero.SetLevel(h.GetLevel())
	newHero.SetTypeID(h.GetTypeID())

	entry := global.GetHeroEntry(h.GetTypeID())
	newHero.SetEntry(entry)

	attManager := att.NewAttManager(entry.AttID)
	newHero.SetAttManager(attManager)

	equipBar := item.NewEquipBar(newHero)
	newHero.SetEquipBar(equipBar)

	runeBox := rune.NewRuneBox(newHero)
	newHero.SetRuneBox(runeBox)

	m.mapHero[newHero.GetID()] = newHero

	return newHero
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
		if v.GetTypeID() == typeMisc {
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
		if v.GetTypeID() == typeMisc {
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
				m.DelHero(v.GetID())
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

func (m *HeroManager) LoadFromDB() {
	l := hero.LoadAll(m.ds, m.owner.GetID(), m.TableName())
	sliceHero := make([]hero.Hero, 0)

	listHero := reflect.ValueOf(l)
	if listHero.Kind() != reflect.Slice {
		logger.Error("load hero returns non-slice type")
		return
	}

	for n := 0; n < listHero.Len(); n++ {
		p := listHero.Index(n)
		sliceHero = append(sliceHero, p.Interface().(hero.Hero))
	}

	for _, v := range sliceHero {
		m.createDBHero(v)
	}
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
	heroEntry := global.GetHeroEntry(typeID)
	h := m.createEntryHero(heroEntry)
	if h == nil {
		return nil
	}

	m.save(h)
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
	m.delete(h, &bson.D{{"_id", id}})
}

func (m *HeroManager) HeroSetLevel(level int32) {
	for _, v := range m.mapHero {
		v.SetLevel(level)

		update := &bson.D{{"$set",
			bson.D{
				{"level", v.GetLevel()},
			},
		}}

		m.saveField(v, update)
	}
}

func (m *HeroManager) PutonEquip(heroID int64, equipID int64) error {

	equip := m.owner.ItemManager().GetItem(equipID)
	if equip == nil {
		return fmt.Errorf("cannot find equip<%d> while PutonEquip", equipID)
	}

	if objID := equip.GetEquipObj(); objID != -1 {
		return fmt.Errorf("equip has put on another hero<%d>", objID)
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

	m.owner.ItemManager().Save(equip.GetID())
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
		return fmt.Errorf("equip<%d> didn't put on this hero<%d> ", equip.GetID(), heroID)
	}

	// unequip
	if err := equipBar.TakeoffEquip(pos); err != nil {
		return err
	}

	m.owner.ItemManager().Save(equip.GetID())
	m.owner.ItemManager().SendItemUpdate(equip)
	m.SendHeroUpdate(h)

	// att
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) PutonRune(heroID int64, runeID int64) error {

	r := m.owner.RuneManager().GetRune(runeID)
	if r == nil {
		return fmt.Errorf("cannot find rune<%d> while PutonRune", runeID)
	}

	if objID := r.GetEquipObj(); objID != -1 {
		return fmt.Errorf("rune has put on another obj<%d>", objID)
	}

	pos := r.Entry().Pos
	if pos < define.Rune_PositionBegin || pos >= define.Rune_PositionEnd {
		return fmt.Errorf("invalid pos<%d>", pos)
	}

	h, ok := m.mapHero[heroID]
	if !ok {
		return fmt.Errorf("invalid heroid<%d>", heroID)
	}

	runeBox := h.GetRuneBox()

	// takeoff previous rune
	if pr := runeBox.GetRuneByPos(pos); pr != nil {
		if err := m.TakeoffRune(heroID, pos); err != nil {
			return err
		}
	}

	// equip new rune
	if err := runeBox.PutonRune(r); err != nil {
		return err
	}

	m.owner.RuneManager().Save(runeID)
	m.owner.RuneManager().SendRuneUpdate(r)
	m.SendHeroUpdate(h)

	// att
	r.GetAttManager().CalcAtt()
	h.GetAttManager().ModAttManager(r.GetAttManager())
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) TakeoffRune(heroID int64, pos int32) error {
	if pos < 0 || pos >= define.Rune_PositionEnd {
		return fmt.Errorf("invalid pos<%d>", pos)
	}

	h, ok := m.mapHero[heroID]
	if !ok {
		return fmt.Errorf("invalid heroid<%d>", heroID)
	}

	r := h.GetRuneBox().GetRuneByPos(pos)
	if r == nil {
		return fmt.Errorf("cannot find rune from hero<%d>'s runebox pos<%d> while TakeoffRune", heroID, pos)
	}

	// unequip
	if err := h.GetRuneBox().TakeoffRune(pos); err != nil {
		return err
	}

	m.owner.RuneManager().Save(r.GetID())
	m.owner.RuneManager().SendRuneUpdate(r)
	m.SendHeroUpdate(h)

	// att
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) SendHeroUpdate(h hero.Hero) {
	// send equips update
	reply := &pbGame.M2C_HeroInfo{
		Info: &pbGame.Hero{
			Id:     h.GetID(),
			TypeId: h.GetTypeID(),
			Exp:    h.GetExp(),
			Level:  h.GetLevel(),
		},
	}

	// equip list
	eb := h.GetEquipBar()
	var n int32
	for n = 0; n < define.Hero_MaxEquip; n++ {
		var equipId int64 = -1
		if i := eb.GetEquipByPos(n); i != nil {
			equipId = i.GetID()
		}

		reply.Info.EquipList = append(reply.Info.EquipList, equipId)
	}

	// rune list
	var pos int32
	for pos = 0; pos < define.Rune_PositionEnd; pos++ {
		var runeId int64 = -1
		if r := h.GetRuneBox().GetRuneByPos(pos); r != nil {
			runeId = r.GetID()
		}

		reply.Info.RuneList = append(reply.Info.RuneList, runeId)
	}

	m.owner.SendProtoMessage(reply)
}

func (m *HeroManager) SendHeroAtt(h hero.Hero) {
	attManager := h.GetAttManager()
	reply := &pbGame.M2C_HeroAttUpdate{
		HeroId:   h.GetID(),
		AttValue: &pbGame.Att{},
	}

	for k := int32(0); k < define.Att_End; k++ {
		reply.AttValue.Value = append(reply.AttValue.Value, attManager.GetAttValue(k))
	}

	m.owner.SendProtoMessage(reply)
}
