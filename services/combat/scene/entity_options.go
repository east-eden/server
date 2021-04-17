package scene

import (
	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/internal/att"
	"bitbucket.org/funplus/server/utils"
	"github.com/willf/bitset"
)

type EntityOption func(*EntityOptions)
type EntityOptions struct {
	TypeId       int32
	Pos          *Position
	InitAtbValue int32
	AttManager   *att.AttManager
	Scene        *Scene
	SceneCamp    *SceneCamp
	Entry        *auto.HeroEntry

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
		TypeId: -1,
		Pos: &Position{
			Pos:    Pos{X: 0, Z: 0},
			Rotate: 0,
		},
		InitAtbValue: 0,
		Entry:        nil,

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

func WithEntityTypeId(typeId int32) EntityOption {
	return func(o *EntityOptions) {
		o.TypeId = typeId
	}
}

func WithEntityHeroEntry(entry *auto.HeroEntry) EntityOption {
	return func(o *EntityOptions) {
		o.Entry = entry

		o.GeneralSkill, _ = auto.GetSkillBaseEntry(entry.Skill1)
		o.NormalSkill, _ = auto.GetSkillBaseEntry(entry.Skill2)
		o.UltimateSkill, _ = auto.GetSkillBaseEntry(entry.Skill3)
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
		o.AttManager.SetFinalAttValue(tp, value)
	}
}

func WithEntityAttList(attList []int32) EntityOption {
	return func(o *EntityOptions) {
		for tp := range attList {
			o.AttManager.SetFinalAttValue(tp, attList[tp])
		}
	}
}

func WithEntityPosition(posX, posZ, rotate int32) EntityOption {
	return func(o *EntityOptions) {
		o.Pos.X = posX
		o.Pos.Z = posZ
		o.Pos.Rotate = rotate
	}
}

func WithEntityInitAtbValue(value int32) EntityOption {
	return func(o *EntityOptions) {
		o.InitAtbValue = value
	}
}
