package player

import (
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
)

type defaultPlayer struct {
	id          int64
	name        string
	itemManager *item.ItemManager
	heroManager *hero.HeroManager
}

func newDefaultPlayer(id int64, name string) Player {
	return &defaultPlayer{
		id:          id,
		name:        name,
		itemManager: item.NewItemManager(),
		heroManager: hero.NewHeroManager(),
	}
}

func (p *defaultPlayer) Init() error {
	return nil
}

func (p *defaultPlayer) ID() int64 {
	return p.id
}

func (p *defaultPlayer) Name() string {
	return p.name
}

func (p *defaultPlayer) HeroManager() *hero.HeroManager {
	return p.heroManager
}

func (p *defaultPlayer) ItemManager() *item.ItemManager {
	return p.itemManager
}
