package hero

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/att"
	"github.com/east-eden/server/services/game/item"
	"github.com/east-eden/server/services/game/rune"
)

type HeroV1 struct {
	Options    `bson:"inline" json:",inline"`
	equipBar   *item.EquipBar  `bson:"-" json:"-"`
	attManager *att.AttManager `bson:"-" json:"-"`
	runeBox    *rune.RuneBox   `bson:"-" json:"-"`
}

func newPoolHeroV1() interface{} {
	h := &HeroV1{
		Options: DefaultOptions(),
	}

	h.equipBar = item.NewEquipBar(h)
	h.attManager = att.NewAttManager(-1)
	h.runeBox = rune.NewRuneBox(h)

	return h
}

func (h *HeroV1) GetOptions() *Options {
	return &h.Options
}

// store.StoreObjector interface
func (h *HeroV1) AfterLoad() error {
	return nil
}

func (h *HeroV1) GetObjID() int64 {
	return h.Options.Id
}

func (h *HeroV1) GetStoreIndex() int64 {
	return h.Options.OwnerId
}

func (h *HeroV1) GetType() int32 {
	return define.Plugin_Hero
}

func (h *HeroV1) GetID() int64 {
	return h.Options.Id
}

func (h *HeroV1) GetLevel() int32 {
	return h.Options.Level
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
	h.Exp += exp
	return h.Exp
}

func (h *HeroV1) AddLevel(level int32) int32 {
	h.Level += level
	return h.Level
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
