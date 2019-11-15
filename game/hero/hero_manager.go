package hero

import "github.com/yokaiio/yokai_server/game/player"

type HeroManager struct {
	mapHero map[int64]player.Hero
}

func NewHeroManager() HeroManager {

	return &HeroManager{
		mapHero: make(map[int64]Hero, 0),
	}
}
