package scene

import (
	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"github.com/shopspring/decimal"
)

type SkillOption func(*SkillOptions)
type SkillOptions struct {
	Caster         *SceneEntity
	Target         *SceneEntity
	TargetPosition *Position

	Triggered bool
	Amount    int32
	SpellType define.ESpellType
	Level     uint32
	Entry     *auto.SkillBaseEntry
}

func DefaultSkillOptions() *SkillOptions {
	return &SkillOptions{
		Caster:    nil,
		Target:    nil,
		Triggered: false,
		Amount:    0,
		SpellType: define.SpellType_Melee,
		Level:     1,
		Entry:     nil,
	}
}

func WithSkillCaster(caster *SceneEntity) SkillOption {
	return func(o *SkillOptions) {
		o.Caster = caster
	}
}

func WithSkillTarget(target *SceneEntity) SkillOption {
	return func(o *SkillOptions) {
		o.Target = target
	}
}

func WithSkillTargetPos(x, z, rotate decimal.Decimal) SkillOption {
	return func(o *SkillOptions) {
		o.TargetPosition.X = x
		o.TargetPosition.Z = z
		o.TargetPosition.Rotate = rotate
	}
}

func WithSkillTriggered(triggered bool) SkillOption {
	return func(o *SkillOptions) {
		o.Triggered = triggered
	}
}

func WithSpellAmount(amount int32) SkillOption {
	return func(o *SkillOptions) {
		o.Amount = amount
	}
}

func WithSpellType(tp define.ESpellType) SkillOption {
	return func(o *SkillOptions) {
		o.SpellType = tp
	}
}

func WithSpellLevel(level uint32) SkillOption {
	return func(o *SkillOptions) {
		o.Level = level
	}
}

func WithSkilEntry(entry *auto.SkillBaseEntry) SkillOption {
	return func(o *SkillOptions) {
		o.Entry = entry
	}
}
