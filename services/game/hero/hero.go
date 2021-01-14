package hero

import (
	"sync"

	"e.coding.net/mmstudio/blade/server/internal/att"
	"e.coding.net/mmstudio/blade/server/services/game/item"
	"e.coding.net/mmstudio/blade/server/services/game/rune"
	"e.coding.net/mmstudio/blade/server/store"
)

// hero create pool
var heroPool = &sync.Pool{New: newPoolHeroV1}

func NewPoolHero() Hero {
	return heroPool.Get().(Hero)
}

func GetHeroPool() *sync.Pool {
	return heroPool
}

func ReleasePoolHero(x interface{}) {
	heroPool.Put(x)
}

type Hero interface {
	store.StoreObjector

	GetOptions() *Options
	GetEquipBar() *item.EquipBar
	GetAttManager() *att.AttManager
	GetRuneBox() *rune.RuneBox

	AddExp(int64) int64
	AddLevel(int32) int32
	BeforeDelete()
	CalcAtt()
}

func NewHero(opts ...Option) Hero {
	h := NewPoolHero()

	for _, o := range opts {
		o(h.GetOptions())
	}

	return h
}
