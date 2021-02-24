package hero

import (
	"sync"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/services/game/item"
	"bitbucket.org/east-eden/server/services/game/rune"
)

// hero create pool
var heroPool = &sync.Pool{New: newPoolHero}

func NewPoolHero() *Hero {
	return heroPool.Get().(*Hero)
}

func GetHeroPool() *sync.Pool {
	return heroPool
}

func NewHero(opts ...Option) *Hero {
	h := NewPoolHero()

	for _, o := range opts {
		o(h.GetOptions())
	}

	return h
}

type Hero struct {
	Options    `bson:"inline" json:",inline"`
	equipBar   *item.EquipBar  `bson:"-" json:"-"`
	attManager *HeroAttManager `bson:"-" json:"-"`
	runeBox    *rune.RuneBox   `bson:"-" json:"-"`
}

func newPoolHero() interface{} {
	h := &Hero{
		Options: DefaultOptions(),
	}

	h.equipBar = item.NewEquipBar(h)
	h.attManager = NewHeroAttManager(h)
	h.runeBox = rune.NewRuneBox(h)

	return h
}

func (h *Hero) GetOptions() *Options {
	return &h.Options
}

func (h *Hero) GetStoreIndex() int64 {
	return h.Options.OwnerId
}

func (h *Hero) GetType() int32 {
	return define.Plugin_Hero
}

func (h *Hero) GetID() int64 {
	return h.Options.Id
}

func (h *Hero) GetLevel() int32 {
	return int32(h.Options.Level)
}

func (h *Hero) GetAttManager() *HeroAttManager {
	return h.attManager
}

func (h *Hero) GetEquipBar() *item.EquipBar {
	return h.equipBar
}

func (h *Hero) GetRuneBox() *rune.RuneBox {
	return h.runeBox
}

func (h *Hero) AddExp(exp int32) int32 {
	h.Exp += exp
	return h.Exp
}

func (h *Hero) AddLevel(level int8) int8 {
	h.Level += level
	return h.Level
}

func (h *Hero) BeforeDelete() {

}
