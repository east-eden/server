package hero

import (
	"sync"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/services/game/item"
)

// hero create pool
var heroPool = &sync.Pool{New: newPoolHero}

func GetHeroPool() *sync.Pool {
	return heroPool
}

func NewHero() *Hero {
	return heroPool.Get().(*Hero)
}

type Hero struct {
	Options    `bson:"inline" json:",inline"`
	equipBar   *item.EquipBar   `bson:"-" json:"-"`
	attManager *HeroAttManager  `bson:"-" json:"-"`
	crystalBox *item.CrystalBox `bson:"-" json:"-"`
}

func newPoolHero() interface{} {
	h := &Hero{
		Options: DefaultOptions(),
	}

	h.equipBar = item.NewEquipBar(h)
	h.attManager = NewHeroAttManager(h)
	h.crystalBox = item.NewCrystalBox(h)

	return h
}

func (h *Hero) Init(opts ...Option) {
	for _, o := range opts {
		o(h.GetOptions())
	}

	h.attManager.SetBaseAttId(h.Entry.AttId)
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

func (h *Hero) GetCrystalBox() *item.CrystalBox {
	return h.crystalBox
}

func (h *Hero) AddExp(exp int32) int32 {
	h.Exp += exp
	return h.Exp
}

func (h *Hero) AddLevel(level int8) int8 {
	h.Level += level
	return h.Level
}
