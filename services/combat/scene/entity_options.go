package scene

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/internal/att"
	"github.com/east-eden/server/utils"
	"github.com/shopspring/decimal"
	"github.com/willf/bitset"
)

type EntityOption func(*EntityOptions)
type EntityOptions struct {
	MonsterId    int32
	HeroId       int32
	Pos          *Position
	InitAtbValue decimal.Decimal
	AttManager   *att.AttManager
	Scene        *Scene
	SceneCamp    *SceneCamp

	MonsterEntry *auto.MonsterEntry
	HeroEntry    *auto.HeroEntry
	ModelEntry   *auto.ModelEntry

	GeneralSkill  *auto.SkillBaseEntry   // 普攻技能
	NormalSkill   *auto.SkillBaseEntry   // 一般技能
	UltimateSkill *auto.SkillBaseEntry   // 特殊技能
	CrystalSkills []*auto.SkillBaseEntry // 残响技能
	PassiveSkills []*auto.SkillBaseEntry // 被动技能列表

	State    *utils.CountableBitset
	Immunity [define.ImmunityType_End]*bitset.BitSet
}

func DefaultEntityOptions() *EntityOptions {
	o := &EntityOptions{
		MonsterId:    -1,
		InitAtbValue: decimal.NewFromInt32(0),
		MonsterEntry: nil,
		HeroEntry:    nil,
		ModelEntry:   nil,

		GeneralSkill:  nil,
		NormalSkill:   nil,
		UltimateSkill: nil,
		CrystalSkills: make([]*auto.SkillBaseEntry, 0, define.Crystal_Subtype_Num),
		PassiveSkills: make([]*auto.SkillBaseEntry, 0, define.Skill_PassiveNum),
		AttManager:    att.NewAttManager(),
		Scene:         nil,
		SceneCamp:     nil,
		State:         utils.NewCountableBitset(uint(define.HeroState_End)),
	}

	for k := range o.Immunity {
		o.Immunity[k] = bitset.New(uint(64))
	}

	return o
}

func WithEntityMonsterId(typeId int32) EntityOption {
	return func(o *EntityOptions) {
		o.MonsterId = typeId
	}
}

func WithEntityHeroId(heroId int32) EntityOption {
	return func(o *EntityOptions) {
		o.HeroId = heroId
	}
}

func WithEntityHeroEntry(entry *auto.HeroEntry) EntityOption {
	return func(o *EntityOptions) {
		o.HeroEntry = entry

		o.GeneralSkill, _ = auto.GetSkillBaseEntry(entry.Skill1)
		o.NormalSkill, _ = auto.GetSkillBaseEntry(entry.Skill2)
		o.UltimateSkill, _ = auto.GetSkillBaseEntry(entry.Skill3)
	}
}

func WithEntityMonsterEntry(entry *auto.MonsterEntry) EntityOption {
	return func(o *EntityOptions) {
		o.MonsterEntry = entry
	}
}

func WithEntityModelEntry(entry *auto.ModelEntry) EntityOption {
	return func(o *EntityOptions) {
		o.ModelEntry = entry
	}
}

func WithEntityScene(scene *Scene) EntityOption {
	return func(o *EntityOptions) {
		o.Scene = scene
	}
}

func WithEntitySceneCamp(camp *SceneCamp) EntityOption {
	return func(o *EntityOptions) {
		o.SceneCamp = camp
	}
}

func WithCrystalSkills(skills []*auto.SkillBaseEntry) EntityOption {
	return func(o *EntityOptions) {
		for _, v := range skills {
			o.CrystalSkills = append(o.CrystalSkills, v)
		}
	}
}

func WithPassiveSkills(skills []*auto.SkillBaseEntry) EntityOption {
	return func(o *EntityOptions) {
		for _, v := range skills {
			o.PassiveSkills = append(o.PassiveSkills, v)
		}
	}
}

func WithEntityAttValue(tp int, value int32) EntityOption {
	return func(o *EntityOptions) {
		o.AttManager.SetFinalAttValue(tp, decimal.NewFromInt32(value))
	}
}

func WithEntityAttList(attList []float32) EntityOption {
	return func(o *EntityOptions) {
		for tp := range attList {
			o.AttManager.SetFinalAttValue(tp, decimal.NewFromFloat32(attList[tp]).Round(2))
		}
	}
}

func WithEntityPosition(posX, posZ, rotate decimal.Decimal) EntityOption {
	return func(o *EntityOptions) {
		o.Pos.X = posX
		o.Pos.Z = posZ
		o.Pos.Rotate = rotate
	}
}

func WithEntityInitAtbValue(value decimal.Decimal) EntityOption {
	return func(o *EntityOptions) {
		o.InitAtbValue = value
	}
}
