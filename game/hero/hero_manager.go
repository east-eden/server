package hero

import (
	"fmt"
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type HeroManager struct {
	Owner        define.PluginObj
	mapHero      map[int64]Hero
	mapEquipHero map[int64]int64 // map[EquipID]HeroID

	ds *db.Datastore
	sync.RWMutex
}

func NewHeroManager(owner define.PluginObj, ds *db.Datastore) *HeroManager {
	m := &HeroManager{
		Owner:        owner,
		ds:           ds,
		mapHero:      make(map[int64]Hero, 0),
		mapEquipHero: make(map[int64]int64, 0),
	}

	return m
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
		h := m.AddHero(typeMisc)
		if h == nil {
			return fmt.Errorf("hero manager gain hero<%d> failed, cannot add new hero<%d>", typeMisc, num)
		}
	}

	return nil
}

func (m *HeroManager) LoadFromDB() {
	l := LoadAll(m.ds, m.Owner.GetID())
	sliceHero := make([]Hero, 0)

	listHero := reflect.ValueOf(l)
	if listHero.Kind() != reflect.Slice {
		logger.Error("load hero returns non-slice type")
		return
	}

	for n := 0; n < listHero.Len(); n++ {
		p := listHero.Index(n)
		sliceHero = append(sliceHero, p.Interface().(Hero))
	}

	for _, v := range sliceHero {
		m.newDBHero(v)

		maxID, err := utils.GeneralIDGet(define.Plugin_Hero)
		if err != nil {
			logger.Error(err)
			return
		}

		if v.GetID() >= maxID {
			utils.GeneralIDSet(define.Plugin_Hero, v.GetID())
		}
	}
}

func (m *HeroManager) newEntryHero(entry *define.HeroEntry) Hero {
	if entry == nil {
		logger.Error("newEntryHero with nil HeroEntry")
		return nil
	}

	id, err := utils.GeneralIDGen(define.Plugin_Hero)
	if err != nil {
		logger.Error(err)
		return nil
	}

	hero := NewHero(id)
	hero.SetOwnerID(m.Owner.GetID())
	hero.SetOwnerType(m.Owner.GetType())
	hero.SetLevel(m.Owner.GetLevel())
	hero.SetTypeID(entry.ID)
	hero.SetEntry(entry)

	m.mapHero[hero.GetID()] = hero

	return hero
}

func (m *HeroManager) newDBHero(h Hero) Hero {
	hero := NewHero(h.GetID())
	hero.SetOwnerID(h.GetOwnerID())
	hero.SetOwnerType(h.GetOwnerType())
	hero.SetLevel(h.GetLevel())
	hero.SetTypeID(h.GetTypeID())
	hero.SetEntry(global.GetHeroEntry(h.GetTypeID()))

	m.mapHero[hero.GetID()] = hero

	return hero
}

func (m *HeroManager) GetHero(id int64) Hero {
	return m.mapHero[id]
}

func (m *HeroManager) GetHeroNums() int {
	return len(m.mapHero)
}

func (m *HeroManager) GetHeroList() []Hero {
	list := make([]Hero, 0)

	m.RLock()
	for _, v := range m.mapHero {
		list = append(list, v)
	}
	m.RUnlock()

	return list
}

func (m *HeroManager) AddHero(typeID int32) Hero {
	heroEntry := global.GetHeroEntry(typeID)
	hero := m.newEntryHero(heroEntry)
	if hero == nil {
		return nil
	}

	m.ds.ORM().Save(hero)
	return hero
}

func (m *HeroManager) DelHero(id int64) {
	h, ok := m.mapHero[id]
	if !ok {
		return
	}

	equipList := h.GetEquips()
	h.BeforeDelete()

	delete(m.mapHero, id)
	for _, v := range equipList {
		delete(m.mapEquipHero, v)
	}

	m.ds.ORM().Delete(h)
}

func (m *HeroManager) HeroSetLevel(level int32) {
	for _, v := range m.mapHero {
		v.SetLevel(level)
		m.ds.ORM().Save(v)
	}
}

func (m *HeroManager) PutonEquip(heroID int64, equipID int64, pos int32) error {
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return fmt.Errorf("invalid pos")
	}

	hero, ok := m.mapHero[heroID]
	if !ok {
		return fmt.Errorf("invalid heroid")
	}

	if id, ok := m.mapEquipHero[equipID]; ok {
		return fmt.Errorf("equip has put on another hero<%d>", id)
	}

	equipList := hero.GetEquips()
	if equipList[pos] != -1 {
		return fmt.Errorf("pos existing equip_id<%d>", equipList[pos])
	}

	hero.SetEquip(equipID, pos)
	m.mapEquipHero[equipID] = heroID
	return nil
}

func (m *HeroManager) TakeoffEquip(heroID int64, pos int32) error {
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return fmt.Errorf("invalid pos")
	}

	hero, ok := m.mapHero[heroID]
	if !ok {
		return fmt.Errorf("invalid heroid")
	}

	equipID := hero.GetEquips()[pos]
	if _, ok := m.mapEquipHero[equipID]; !ok {
		return fmt.Errorf("equip didn't put on this hero<%d>", heroID)
	}

	hero.UnsetEquip(pos)
	delete(m.mapEquipHero, equipID)
	return nil
}
