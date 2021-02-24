package hero

import (
	"sync"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/att"
	"github.com/east-eden/server/services/game/item"
	"github.com/east-eden/server/services/game/rune"
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
	attManager *att.AttManager `bson:"-" json:"-"`
	runeBox    *rune.RuneBox   `bson:"-" json:"-"`
}

func newPoolHero() interface{} {
	h := &Hero{
		Options: DefaultOptions(),
	}

	h.equipBar = item.NewEquipBar(h)
	h.attManager = att.NewAttManager()
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

func (h *Hero) GetAttManager() *att.AttManager {
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

func (h *Hero) CalcAtt() {
	h.attManager.Reset()

	// equip bar
	var n int32
	for n = 0; n < int32(define.Equip_Pos_End); n++ {
		e := h.equipBar.GetEquipByPos(n)
		if e == nil {
			continue
		}

		e.GetAttManager().CalcAtt()
		h.attManager.ModAttManager(e.GetAttManager())
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
