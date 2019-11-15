package player

import (
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
)

type defaultPlayer struct {
	id          int64
	name        string
	itemManager ItemManager
	heroManager HeroManager
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
	p.id = 100001
	return nil
}

func (p *defaultPlayer) ID() int64 {
	return p.id
}

func (p *defaultPlayer) Name() string {
	return p.name
}

func (p *defaultPlayer) HeroManager() HeroManager {
	return p.heroManager
}

func (p *defaultPlayer) ItemManager() ItemManager {
	return p.itemManager
}
