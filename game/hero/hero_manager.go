package hero

import (
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type HeroManager struct {
	OwnerID int64
	mapHero map[int64]Hero

	ds *db.Datastore
	sync.RWMutex
}

func NewHeroManager(ownerID int64, ds *db.Datastore) *HeroManager {
	m := &HeroManager{
		OwnerID: ownerID,
		ds:      ds,
		mapHero: make(map[int64]Hero, 0),
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

	m.Lock()
	defer m.Unlock()
	m.mapHero[hero.GetID()] = hero

	return hero
}

func (m *HeroManager) newDBHero(h Hero) Hero {
	hero := NewHero(h.GetID())
	hero.SetOwnerID(h.GetOwnerID())
	hero.SetTypeID(h.GetTypeID())
	hero.SetEntry(global.GetHeroEntry(h.GetTypeID()))

	m.Lock()
	defer m.Unlock()
	m.mapHero[hero.GetID()] = hero

	return hero
}

func (m *HeroManager) GetHero(id int64) Hero {
	m.RLock()
	defer m.RUnlock()
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
	m.Lock()
	h, ok := m.mapHero[id]
	if !ok {
		m.Unlock()
		return
	}

	delete(m.mapHero, id)
	m.Unlock()

	m.ds.ORM().Delete(h)
}

func (m *HeroManager) HeroAddExp(id int64, exp int64) {
	m.RLock()
	hero, ok := m.mapHero[id]
	m.RUnlock()

	if ok {
		hero.AddExp(exp)
		m.ds.ORM().Save(hero)
	}
}
