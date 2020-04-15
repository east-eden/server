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
	h.SetAttManager(attManager)

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
			equips := v.GetEquips()
			hasEquip := false
			for i := 0; i < define.Hero_MaxEquip; i++ {
				if equips[i] != -1 {
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
			equips := v.GetEquips()
			hasEquip := false
			for i := 0; i < define.Hero_MaxEquip; i++ {
				if equips[i] != -1 {
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

	equipList := h.GetEquips()
	for _, v := range equipList {
		if i := m.owner.ItemManager().GetItem(v); i != nil {
			i.SetEquipObj(-1)
		}
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

	if equip.EquipEnchantEntry() == nil {
		return fmt.Errorf("cannot find equip_enchant_entry<%d> while PutonEquip", equipID)
	}

	if objID := equip.GetEquipObj(); objID != -1 {
		return fmt.Errorf("equip has put on another hero<%d>", objID)
	}

	pos := equip.EquipEnchantEntry().EquipPos
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return fmt.Errorf("invalid pos")
	}

	h, ok := m.mapHero[heroID]
	if !ok {
		return fmt.Errorf("invalid heroid")
	}

	equipList := h.GetEquips()
	if equipList[pos] != -1 {
		return fmt.Errorf("pos existing equip_id<%d>", equipList[pos])
	}

	// equip
	h.SetEquip(equipID, pos)
	m.owner.ItemManager().SetItemEquiped(equipID, heroID)
	m.SendHeroEquips(h)

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

	equipID := h.GetEquips()[pos]
	equip := m.owner.ItemManager().GetItem(equipID)
	if equip == nil {
		return fmt.Errorf("cannot find equip<%d> while TakeoffEquip", equipID)
	}

	if objID := equip.GetEquipObj(); objID == -1 {
		return fmt.Errorf("equip didn't put on this hero<%d> pos<%d>", heroID, pos)
	}

	// unequip
	h.UnsetEquip(pos)
	m.owner.ItemManager().SetItemUnEquiped(equipID)
	m.SendHeroEquips(h)

	// att
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) SendHeroEquips(h hero.Hero) {
	// send equips update
	reply := &pbGame.M2C_HeroEquips{
		HeroId: h.GetID(),
		Equips: make([]*pbGame.Item, 0),
	}

	equips := h.GetEquips()
	for _, v := range equips {
		if it := m.owner.ItemManager().GetItem(v); it != nil {
			i := &pbGame.Item{
				Id:     v,
				TypeId: it.GetTypeID(),
			}
			reply.Equips = append(reply.Equips, i)
		}
	}

	m.owner.SendProtoMessage(reply)
}

func (m *HeroManager) SendHeroInfo(h hero.Hero) {
	// send equips update
	reply := &pbGame.M2C_HeroInfo{
		Info: &pbGame.Hero{
			Id:        h.GetID(),
			TypeId:    h.GetTypeID(),
			Exp:       h.GetExp(),
			Level:     h.GetLevel(),
			EquipList: make([]int64, 0),
		},
	}

	equips := h.GetEquips()
	for _, v := range equips {
		reply.Info.EquipList = append(reply.Info.EquipList, v)
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
		reply.AttValue.Value[k] = attManager.GetAttValue(k)
	}

	for k := int32(0); k < define.AttEx_End; k++ {
		reply.AttValue.ExValue[k] = attManager.GetAttExValue(k)
	}

	m.owner.SendProtoMessage(reply)
}
