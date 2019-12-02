package hero

import (
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/define"
	"github.com/yokaiio/yokai_server/game/global"
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
		m.NewDBHero(v)

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

func (m *HeroManager) Save(h Hero) {
	m.ds.ORM().Save(h)
}

func (m *HeroManager) NewHero(entry *define.HeroEntry) Hero {
	if entry == nil {
		logger.Error("NewHero with nil HeroEntry")
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

func (m *HeroManager) NewDBHero(h Hero) Hero {
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
