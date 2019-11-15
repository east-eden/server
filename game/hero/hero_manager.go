package hero

import "sync/atomic"

type HeroManager struct {
	idGen   atomic.Value
	mapHero map[int64]Hero
}

func NewHeroManager() *HeroManager {
	m := &HeroManager{
		mapHero: make(map[int64]Hero, 0),
	}

	m.idGen.Store(int64(0))
	return m
}

func (m *HeroManager) GenID() int64 {
	id := m.idGen.Load().(int64) + 1
	m.idGen.Store(id)
	return id
}

func (m *HeroManager) NewHero(typeID int32) Hero {
	id := m.GenID()
	hero := NewHero(id, typeID)
	m.mapHero[hero.ID()] = hero
	return hero
}
