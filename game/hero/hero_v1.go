package hero

import (
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/rune"
)

type HeroV1 struct {
	Opts       *Options        `bson:"inline"`
	equipBar   *item.EquipBar  `gorm:"-" bson:"-"`
	attManager *att.AttManager `gorm:"-" bson:"-"`
	runeBox    *rune.RuneBox   `gorm:"-" bson:"-"`
}

func newPoolHeroV1() interface{} {
	h := &HeroV1{
		Opts: DefaultOptions(),
	}

	h.equipBar = item.NewEquipBar(h)
	h.attManager = att.NewAttManager(-1)
	h.runeBox = rune.NewRuneBox(h)

	return h
}

func (h *HeroV1) Options() *Options {
	return h.Opts
}

func (h *HeroV1) GetType() int32 {
	return define.Plugin_Hero
}

func (h *HeroV1) GetID() int64 {
	return h.Opts.Id
}

func (h *HeroV1) GetLevel() int32 {
	return h.Opts.Level
}

func (h *HeroV1) GetAttManager() *att.AttManager {
	return h.attManager
}

func (h *HeroV1) GetEquipBar() *item.EquipBar {
	return h.equipBar
}

func (h *HeroV1) GetRuneBox() *rune.RuneBox {
	return h.runeBox
}

func (h *HeroV1) AddExp(exp int64) int64 {
	h.Opts.Exp += exp
	return h.Opts.Exp
}

func (h *HeroV1) AddLevel(level int32) int32 {
	h.Opts.Level += level
	return h.Opts.Level
}

func (h *HeroV1) BeforeDelete() {

}

func (h *HeroV1) CalcAtt() {
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
