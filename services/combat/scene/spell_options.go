package scene

import (
	"github.com/east-eden/server/define"
)

type SpellOption func(*SpellOptions)
type SpellOptions struct {
	Caster    *SceneUnit
	Target    *SceneUnit
	Triggered bool
	Amount    int32
	SpellType define.ESpellType
	Level     uint32
	Entry     *define.SpellEntry
}

func DefaultSpellOptions() *SpellOptions {
	return &SpellOptions{
		Caster:    nil,
		Target:    nil,
		Triggered: false,
		Amount:    0,
		SpellType: define.SpellType_Melee,
		Level:     1,
		Entry:     nil,
	}
}

func WithSpellCaster(caster *SceneUnit) SpellOption {
	return func(o *SpellOptions) {
		o.Caster = caster
	}
}

func WithSpellTarget(target *SceneUnit) SpellOption {
	return func(o *SpellOptions) {
		o.Target = target
	}
}

func WithSpellTriggered(triggered bool) SpellOption {
	return func(o *SpellOptions) {
		o.Triggered = triggered
	}
}

func WithSpellAmount(amount int32) SpellOption {
	return func(o *SpellOptions) {
		o.Amount = amount
	}
}

func WithSpellType(tp define.ESpellType) SpellOption {
	return func(o *SpellOptions) {
		o.SpellType = tp
	}
}

func WithSpellLevel(level uint32) SpellOption {
	return func(o *SpellOptions) {
		o.Level = level
	}
}

func WithSpellEntry(entry *define.SpellEntry) SpellOption {
	return func(o *SpellOptions) {
		o.Entry = entry
	}
}
