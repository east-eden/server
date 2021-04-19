package scene

import (
	"bitbucket.org/funplus/server/define"
)

type AuraOption func(*BuffOptions)
type BuffOptions struct {
	Owner        *SceneEntity
	Caster       *SceneEntity
	Amount       int32
	Level        uint32
	RagePctMod   float32
	CurWrapTimes uint8
	SpellType    define.ESpellType
	SlotIndex    int8 // Aura栏位序号
	Entry        *define.AuraEntry
}

func DefaultBuffOptions() *BuffOptions {
	return &BuffOptions{
		Owner:        nil,
		Caster:       nil,
		Amount:       0,
		Level:        1,
		RagePctMod:   0.0,
		CurWrapTimes: 0,
		SpellType:    define.SpellType_Melee,
		SlotIndex:    -1,
		Entry:        nil,
	}
}

func WithAuraCaster(caster *SceneEntity) AuraOption {
	return func(o *BuffOptions) {
		o.Caster = caster
	}
}

func WithAuraOwner(owner *SceneEntity) AuraOption {
	return func(o *BuffOptions) {
		o.Owner = owner
	}
}

func WithAuraAmount(amount int32) AuraOption {
	return func(o *BuffOptions) {
		o.Amount = amount
	}
}

func WithAuraLevel(level uint32) AuraOption {
	return func(o *BuffOptions) {
		o.Level = level
	}
}

func WithAuraRagePctMod(ragePctMod float32) AuraOption {
	return func(o *BuffOptions) {
		o.RagePctMod = ragePctMod
	}
}

func WithAuraCurWrapTimes(curWrapTimes uint8) AuraOption {
	return func(o *BuffOptions) {
		o.CurWrapTimes = curWrapTimes
	}
}

func WithAuraSpellType(tp define.ESpellType) AuraOption {
	return func(o *BuffOptions) {
		o.SpellType = tp
	}
}

func WithAuraEntry(entry *define.AuraEntry) AuraOption {
	return func(o *BuffOptions) {
		o.Entry = entry
	}
}

func WithAuraSlotIndex(index int8) AuraOption {
	return func(o *BuffOptions) {
		o.SlotIndex = index
	}
}
