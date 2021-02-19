package hero

import (
	"sync"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/internal/att"
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

func ReleasePoolHero(x interface{}) {
	heroPool.Put(x)
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
	attManager *att.AttManager `bson:"-" json:"-"`
	runeBox    *rune.RuneBox   `bson:"-" json:"-"`
}

func newPoolHero() interface{} {
	h := &Hero{
		Options: DefaultOptions(),
	}

	h.equipBar = item.NewEquipBar(h)
	h.attManager = att.NewAttManager(-1)
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
	return h.Options.Level
}

func (h *Hero) GetAttManager() *att.AttManager {
	return h.attManager
}

func (h *Hero) GetEquipBar() *item.EquipBar {
	return h.equipBar
}

func (h *Hero) GetRuneBox() *rune.RuneBox {
	return h.runeBox
}

func (h *Hero) AddExp(exp int64) int64 {
	h.Exp += exp
	return h.Exp
}

func (h *Hero) AddLevel(level int32) int32 {
	h.Level += level
	return h.Level
}

func (h *Hero) BeforeDelete() {

}

func (h *Hero) CalcAtt() {
	h.attManager.Reset()

	// equip bar
	var n int32
	for n = 0; n < define.Hero_MaxEquip; n++ {
		i := h.equipBar.GetEquipByPos(n)
		if i == nil {
			continue
		}

		i.CalcAtt()
		h.attManager.ModAttManager(i.GetAttManager())
	}

	// rune box
	for n = 0; n < define.Rune_PositionEnd; n++ {
		r := h.runeBox.GetRuneByPos(n)
		if r == nil {
			continue
		}

		r.CalcAtt()
		h.attManager.ModAttManager(r.GetAttManager())
	}

	h.attManager.CalcAtt()
}
