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
	TypeId        int32
	PosX          int32
	PosZ          int32
	Rotate        int32
	InitAtbValue  int32
	AttManager    *att.AttManager
	Scene         *Scene
	SceneCamp     *SceneCamp
	Entry         *auto.HeroEntry
	CrystalSkills []*auto.SkillBaseEntry // 残响技能
	PassiveSkills []*auto.SkillBaseEntry // 被动技能列表

	State    *utils.CountableBitset
	Immunity [define.ImmunityType_End]*bitset.BitSet
}

func DefaultEntityOptions() *EntityOptions {
	o := &EntityOptions{
		TypeId:        -1,
		PosX:          0,
		PosZ:          0,
		Rotate:        0,
		InitAtbValue:  0,
		Entry:         nil,
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
		o.AttManager.SetAttValue(tp, value)
	}
}

func WithEntityAttList(attList []int32) EntityOption {
	return func(o *EntityOptions) {
		for tp := range attList {
			o.AttManager.SetAttValue(tp, attList[tp])
		}
	}
}

func WithEntityPosition(posX, posZ, rotate int32) EntityOption {
	return func(o *EntityOptions) {
		o.PosX = posX
		o.PosZ = posZ
		o.Rotate = rotate
	}
}

func WithEntityInitAtbValue(value int32) EntityOption {
	return func(o *EntityOptions) {
		o.InitAtbValue = value
	}
}
