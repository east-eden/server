package hero

import (
	"sync"
)

type HeroManager struct {
	OwnerID int64
	mapHero map[int64]Hero
	sync.RWMutex
}

func NewHeroManager(ownerID int64) *HeroManager {
	m := &HeroManager{
		OwnerID: ownerID,
		mapHero: make(map[int64]Hero, 0),
	}

	return m
}

func (m *HeroManager) LoadFromDB() {

}

/*func (m *HeroManager) GenID() int64 {*/
//id := m.idGen.Load().(int64) + 1
//m.idGen.Store(id)
//return id
/*}*/

func (m *HeroManager) NewHero(typeID int32) Hero {
	id := m.GenID()
	hero := NewHero(id, m.OwnerID, typeID)

	m.Lock()
	m.mapHero[hero.GetID()] = hero
	m.Unlock()
	return hero
}

func (m *HeroManager) GetHero(id int64) Hero {
	m.RLock()
	defer m.RUnlock()
	return m.mapHero[id]
}
