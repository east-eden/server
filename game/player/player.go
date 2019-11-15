package player

import (
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
)

type Player interface {
	Init() error

	ID() int64
	Name() string

	HeroManager() hero.HeroManager
	ItemManager() item.ItemManager
}

func NewPlayer(id int64, name string) Player {
	return newDefaultPlayer(id, name)
}
