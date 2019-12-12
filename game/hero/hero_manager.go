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
	OwnerID      int64
	mapHero      map[int64]Hero
	mapEquipHero map[int64]int64 // map[EquipID]HeroID

	ds *db.Datastore
	sync.RWMutex
}

func NewHeroManager(ownerID int64, ds *db.Datastore) *HeroManager {
	m := &HeroManager{
		OwnerID:      ownerID,
		ds:           ds,
		mapHero:      make(map[int64]Hero, 0),
		mapEquipHero: make(map[int64]int64, 0),
	}

	return m
}

func (m *HeroManager) LoadFromDB() {
	l := LoadAll(m.ds, m.OwnerID)
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
	hero.SetOwnerID(m.OwnerID)
	hero.SetTypeID(entry.ID)
	hero.SetEntry(entry)

	m.mapHero[hero.GetID()] = hero

	return hero
}

func (m *HeroManager) newDBHero(h Hero) Hero {
	hero := NewHero(h.GetID())
	hero.SetOwnerID(h.GetOwnerID())
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

func (m *HeroManager) GetHeroList(list []Hero) {
	m.RLock()
	for _, v := range m.mapHero {
		list = append(list, v)
	}
	m.RUnlock()
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

func (m *HeroManager) HeroAddExp(id int64, exp int64) {
	hero, ok := m.mapHero[id]

	if ok {
		hero.AddExp(exp)
		m.ds.ORM().Save(hero)
	}
}

func (m *HeroManager) HeroAddLevel(id int64, level int32) {
	hero, ok := m.mapHero[id]

	if ok {
		hero.AddLevel(level)
		m.ds.ORM().Save(hero)
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
		return fmt.Errorf("equip has put on another hero", id)
	}

	equipList := hero.GetEquips()
	if equipList[pos] != -1 {
		return fmt.Errorf("pos existing equip_id", equipList[pos])
	}

	hero.SetEquip(equipID, pos)
	m.mapEquipHero[equipID] = heroID
	return nil
}

func (m *HeroManager) TakeoffEquip(heroID int64, equipID int64) error {
	hero, ok := m.mapHero[heroID]
	if !ok {
		return fmt.Errorf("invalid heroid")
	}

	if _, ok := m.mapEquipHero[equipID]; !ok {
		return fmt.Errorf("equip didn't put on this hero", heroID)
	}

	hero.UnsetEquip(equipID)
	delete(m.mapEquipHero, equipID)
	return nil
}
