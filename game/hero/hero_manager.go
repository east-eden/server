package hero

type HeroManager struct {
	mapHero map[int64]Hero
}

func NewHeroManager() *HeroManager {

	return &HeroManager{
		mapHero: make(map[int64]Hero, 0),
	}
}
