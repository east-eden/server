package config

import (
	"github.com/east-eden/server/excel"
	"github.com/shopspring/decimal"
)

// 技能表-拳皇/SkillPoly
type SkillAction struct {
	ActionName string   `json:"actionName,omitempty" bson:"actionName"` //
	KeyTime    []string `json:"keyTime,omitempty" bson:"keyTime"`       //
}

// 技能表-拳皇/SkillPoly
type SkillPerfrom struct {
	Type      SkillPerfromType `json:"type,omitempty" bson:"type"`           //
	TotalTime string           `json:"totalTime,omitempty" bson:"totalTime"` //
	Key       string           `json:"key,omitempty" bson:"key"`             //
}

// 技能表-拳皇/SkillPoly
type RougeSkillCondi struct {
	SkillId int32 `json:"skillId,omitempty" bson:"skillId"` //
	Level   int32 `json:"level,omitempty" bson:"level"`     //
}

// 技能表-拳皇/SkillPoly
type EffectGroup struct {
	EffectName []string `json:"effectName,omitempty" bson:"effectName"` //
}

// 技能表-拳皇/SkillPoly
type SkillSequenceTrigger struct {
	Type         SequenceTriggerType `json:"type,omitempty" bson:"type"`                 //
	Param        float32             `json:"param,omitempty" bson:"param"`               //
	FrocePerform bool                `json:"frocePerform,omitempty" bson:"frocePerform"` //
}

// 技能表-拳皇/SkillPoly
type SkillExcuteCheck struct {
	Type  SkillSequenceExcuteAndCDType `json:"type,omitempty" bson:"type"`   //
	Param []float64                    `json:"param,omitempty" bson:"param"` //
}

// 技能表-拳皇/SkillPoly
type BuffEffect struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// 技能表-拳皇/SkillPoly
type Attack struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
}

// 技能表-拳皇/SkillPoly
type Heal struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          []float64           `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	MyAttrFactor      AttrType            `json:"myAttrFactor,omitempty" bson:"myAttrFactor"`           //
	TargetAttrFactor  AttrType            `json:"targetAttrFactor,omitempty" bson:"targetAttrFactor"`   //
}

// 技能表-拳皇/SkillPoly
type MeleeAttack struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	BringBuffs        []int32             `json:"bringBuffs,omitempty" bson:"bringBuffs"`               //
	TmpChange         []AttrType          `json:"tmpChange,omitempty" bson:"tmpChange"`                 //
	Add               []string            `json:"add,omitempty" bson:"add"`                             //
}

// 技能表-拳皇/SkillPoly
type Shoot struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	BulletId          []string            `json:"bulletId,omitempty" bson:"bulletId"`                   //
}

// 技能表-拳皇/SkillPoly
type Strafing struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	BulletId          []string            `json:"bulletId,omitempty" bson:"bulletId"`                   //
	InitAngle         decimal.Decimal     `json:"initAngle,omitempty" bson:"initAngle"`                 //
	DeltaAngle        decimal.Decimal     `json:"deltaAngle,omitempty" bson:"deltaAngle"`               //
	MaxAngle          decimal.Decimal     `json:"maxAngle,omitempty" bson:"maxAngle"`                   //
}

// 技能表-拳皇/SkillPoly
type ProcessBuffEffect struct {
	EffectName      string  `json:"effectName,omitempty" bson:"effectName"`           //
	Interval        float32 `json:"interval,omitempty" bson:"interval"`               //
	NeedKeep        bool    `json:"needKeep,omitempty" bson:"needKeep"`               //
	DurationPrecent float32 `json:"durationPrecent,omitempty" bson:"durationPrecent"` //
}

// 技能表-拳皇/SkillPoly
type SummonAttrExtendsPair struct {
	AttrType AttrType `json:"attrType,omitempty" bson:"attrType"` //
	Precent  float32  `json:"precent,omitempty" bson:"precent"`   //
}

// 技能表-拳皇/SkillPoly
type DOT struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	Duration          float32             `json:"duration,omitempty" bson:"duration"`                   //
	Interval          float32             `json:"interval,omitempty" bson:"interval"`                   //
	CanCrital         bool                `json:"canCrital,omitempty" bson:"canCrital"`                 //
	FixedDamage       bool                `json:"fixedDamage,omitempty" bson:"fixedDamage"`             //
}

// 技能表-拳皇/SkillPoly
type HOT struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          []float64           `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	CanCrital         bool                `json:"canCrital,omitempty" bson:"canCrital"`                 //
}

// 技能表-拳皇/SkillPoly
type AOE struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	AoeId             string              `json:"aoeId,omitempty" bson:"aoeId"`                         //
	NoDamage          bool                `json:"noDamage,omitempty" bson:"noDamage"`                   //
	OnlyCheckX        bool                `json:"onlyCheckX,omitempty" bson:"onlyCheckX"`               //
}

// #位移类效果
// 技能表-拳皇/SkillPoly
type Move struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	MoveId            string              `json:"moveId,omitempty" bson:"moveId"`                       //
}

// #位移类效果
// 技能表-拳皇/SkillPoly
type AvoidMove struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	MoveId            string              `json:"moveId,omitempty" bson:"moveId"`                       //
	IsWalk            bool                `json:"isWalk,omitempty" bson:"isWalk"`                       //
}

// #冲击波
// 技能表-拳皇/SkillPoly
type ShockWave struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	ShockId           []string            `json:"shockId,omitempty" bson:"shockId"`                     //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Dizzy struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Snare struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Stun struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Freeze struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ExcludeSkllids    []int32             `json:"excludeSkllids,omitempty" bson:"excludeSkllids"`       //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Frozen struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Sleep struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Silence struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Charm struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Disarm struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Taunt struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #控制类效果
// 技能表-拳皇/SkillPoly
type Imprison struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Range             float32             `json:"range,omitempty" bson:"range"`                         //
}

// #改变属性
// 技能表-拳皇/SkillPoly
type ChangeAttr struct {
	EffectRoleType          EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`                   //
	ReplaceTargetType       ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"`             //
	TargetSelectParam       []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"`             //
	TargetsCount            string              `json:"targetsCount,omitempty" bson:"targetsCount"`                       //
	SiftType                SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                               //
	SelectRoleType          EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`                   //
	ExcludeSelf             bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`                         //
	MountEffects            []string            `json:"mountEffects,omitempty" bson:"mountEffects"`                       //
	CenterEffects           []string            `json:"centerEffects,omitempty" bson:"centerEffects"`                     //
	SelfDead                bool                `json:"selfDead,omitempty" bson:"selfDead"`                               //
	SiftParam               []float64           `json:"siftParam,omitempty" bson:"siftParam"`                             //
	Rate                    int32               `json:"rate,omitempty" bson:"rate"`                                       //
	CacheCnt                bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                               //
	NoBaseSelect            bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`                       //
	AttrType                AttrType            `json:"attrType,omitempty" bson:"attrType"`                               //
	Add                     string              `json:"add,omitempty" bson:"add"`                                         //
	Mult                    string              `json:"mult,omitempty" bson:"mult"`                                       //
	EffectTime              int32               `json:"effectTime,omitempty" bson:"effectTime"`                           //
	Min                     string              `json:"min,omitempty" bson:"min"`                                         //
	Max                     string              `json:"max,omitempty" bson:"max"`                                         //
	TotalAddMax             string              `json:"totalAddMax,omitempty" bson:"totalAddMax"`                         //
	ForceEffectSelf         bool                `json:"forceEffectSelf,omitempty" bson:"forceEffectSelf"`                 //
	ExtraAddRuntimeAttrType AttrType            `json:"extraAddRuntimeAttrType,omitempty" bson:"extraAddRuntimeAttrType"` //
	ExtraAddRuntimeValue    string              `json:"extraAddRuntimeValue,omitempty" bson:"extraAddRuntimeValue"`       //
}

// #改变基础属性
// 技能表-拳皇/SkillPoly
type ChangeBaseAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AttrType          AttrType            `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
	EffectTime        int32               `json:"effectTime,omitempty" bson:"effectTime"`               //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type ChangeGrowAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AttrType          AttrType            `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
	EffectTime        int32               `json:"effectTime,omitempty" bson:"effectTime"`               //
}

// #改变特殊属性
// 技能表-拳皇/SkillPoly
type ChangeSpecialAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AttrType          AttrType            `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
	EffectTime        int32               `json:"effectTime,omitempty" bson:"effectTime"`               //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ForceChangAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AttrType          AttrType            `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
	EffectTime        int32               `json:"effectTime,omitempty" bson:"effectTime"`               //
	LimitCalculate    string              `json:"limitCalculate,omitempty" bson:"limitCalculate"`       //
	FroceChange       bool                `json:"froceChange,omitempty" bson:"froceChange"`             //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ChangeGlobalAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AttrType          AttrType            `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
	EffectTime        int32               `json:"effectTime,omitempty" bson:"effectTime"`               //
	CampType          BattleRoleCamp      `json:"campType,omitempty" bson:"campType"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type BuffTrigger struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenBeHurt struct {
	CheckCount      int32   `json:"checkCount,omitempty" bson:"checkCount"`           //
	ExculdeDOT      bool    `json:"exculdeDOT,omitempty" bson:"exculdeDOT"`           //
	ExculdeHeal     bool    `json:"exculdeHeal,omitempty" bson:"exculdeHeal"`         //
	ExculdeSelfCamp bool    `json:"exculdeSelfCamp,omitempty" bson:"exculdeSelfCamp"` //
	SpecialHit      string  `json:"specialHit,omitempty" bson:"specialHit"`           //
	Precent         float32 `json:"precent,omitempty" bson:"precent"`                 //
	CheckPreDamage  bool    `json:"checkPreDamage,omitempty" bson:"checkPreDamage"`   //
	SpecialHit2     string  `json:"specialHit2,omitempty" bson:"specialHit2"`         //
	IgnoreTime      bool    `json:"ignoreTime,omitempty" bson:"ignoreTime"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenDeadlyHurt struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type RandomTrigger struct {
	Rate float32 `json:"rate,omitempty" bson:"rate"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AttackAfterEffect struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ListenEnemyAttr struct {
	TargetAttr           AttrType `json:"targetAttr,omitempty" bson:"targetAttr"`                     //
	ExculdeSelf          bool     `json:"exculdeSelf,omitempty" bson:"exculdeSelf"`                   //
	MaxAttr              AttrType `json:"maxAttr,omitempty" bson:"maxAttr"`                           //
	ComparePrecent       bool     `json:"comparePrecent,omitempty" bson:"comparePrecent"`             //
	TargetValue          float32  `json:"targetValue,omitempty" bson:"targetValue"`                   //
	CompareFunc          string   `json:"compareFunc,omitempty" bson:"compareFunc"`                   //
	CheckSelf            bool     `json:"checkSelf,omitempty" bson:"checkSelf"`                       //
	CheckTargetCampRoles bool     `json:"checkTargetCampRoles,omitempty" bson:"checkTargetCampRoles"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenKill struct {
	NeedKillbyMe bool            `json:"needKillbyMe,omitempty" bson:"needKillbyMe"` //
	SkillId      decimal.Decimal `json:"skillId,omitempty" bson:"skillId"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WillBuildDamage struct {
	JudgeSameTarget      bool    `json:"judgeSameTarget,omitempty" bson:"judgeSameTarget"`           //
	JudgeDefferentTarget bool    `json:"judgeDefferentTarget,omitempty" bson:"judgeDefferentTarget"` //
	JudgeNormal          bool    `json:"judgeNormal,omitempty" bson:"judgeNormal"`                   //
	JudgeUltimate        bool    `json:"judgeUltimate,omitempty" bson:"judgeUltimate"`               //
	WhiteList            []int32 `json:"whiteList,omitempty" bson:"whiteList"`                       //
	IgnoreCollectGas     bool    `json:"ignoreCollectGas,omitempty" bson:"ignoreCollectGas"`         //
	IgnoreTime           bool    `json:"ignoreTime,omitempty" bson:"ignoreTime"`                     //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenBeHurtDelay struct {
	Delay float32 `json:"delay,omitempty" bson:"delay"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenBackAttack struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenWillBeHurt struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenShieldBroken struct {
	ShieldIds []string `json:"shieldIds,omitempty" bson:"shieldIds"` //
	InitTrue  bool     `json:"initTrue,omitempty" bson:"initTrue"`   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenRefractionBroken struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenAfterBackAttack struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenUseSkill struct {
	SkillId       int32        `json:"skillId,omitempty" bson:"skillId"`             //
	CheckCount    int32        `json:"checkCount,omitempty" bson:"checkCount"`       //
	SkillType     SkillTagType `json:"skillType,omitempty" bson:"skillType"`         //
	CheckUltimate bool         `json:"checkUltimate,omitempty" bson:"checkUltimate"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenWillNormalAttack struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenSummonDispear struct {
	SummonId string `json:"summonId,omitempty" bson:"summonId"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AfterSendDamage struct {
	CheckCount int32  `json:"checkCount,omitempty" bson:"checkCount"` //
	SpecialHit string `json:"specialHit,omitempty" bson:"specialHit"` //
	SkillId    int32  `json:"skillId,omitempty" bson:"skillId"`       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AfterSendNormalDamage struct {
	CheckCount   int32   `json:"checkCount,omitempty" bson:"checkCount"`     //
	CompareValue float32 `json:"compareValue,omitempty" bson:"compareValue"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenAttacherDead struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenRoleDead struct {
	CheckType BattleRoleCamp `json:"checkType,omitempty" bson:"checkType"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenLayerMoreThan struct {
	Layer      int32 `json:"layer,omitempty" bson:"layer"`           //
	IgnoreTime bool  `json:"ignoreTime,omitempty" bson:"ignoreTime"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenAttrChanged struct {
	AttrType    AttrType `json:"attrType,omitempty" bson:"attrType"`       //
	CompareSign string   `json:"compareSign,omitempty" bson:"compareSign"` //
	TargetValue float32  `json:"targetValue,omitempty" bson:"targetValue"` //
	PerCheck    float32  `json:"perCheck,omitempty" bson:"perCheck"`       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenCampAttrChanged struct {
	IsSelf      bool     `json:"isSelf,omitempty" bson:"isSelf"`           //
	MaxAttrType AttrType `json:"maxAttrType,omitempty" bson:"maxAttrType"` //
	AttrType    AttrType `json:"attrType,omitempty" bson:"attrType"`       //
	CompareSign string   `json:"compareSign,omitempty" bson:"compareSign"` //
	TargetValue float32  `json:"targetValue,omitempty" bson:"targetValue"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SelfCenterAoe struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	AoeId             string              `json:"aoeId,omitempty" bson:"aoeId"`                         //
	NoDamage          bool                `json:"noDamage,omitempty" bson:"noDamage"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type BuffPair struct {
	EffectType EffectRoleType `json:"effectType,omitempty" bson:"effectType"` //
	BuffId     int32          `json:"buffId,omitempty" bson:"buffId"`         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type Clash struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	MoveId            string              `json:"moveId,omitempty" bson:"moveId"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type MountBuff struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	BuffId            []int32             `json:"buffId,omitempty" bson:"buffId"`                       //
	Count             string              `json:"count,omitempty" bson:"count"`                         //
	FromAttacher      bool                `json:"fromAttacher,omitempty" bson:"fromAttacher"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SubBuff struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	BuffIds           []int32             `json:"buffIds,omitempty" bson:"buffIds"`                     //
	Count             string              `json:"count,omitempty" bson:"count"`                         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type CounterInjury struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	CounterParam      string              `json:"counterParam,omitempty" bson:"counterParam"`           //
	BringBuffs        []string            `json:"bringBuffs,omitempty" bson:"bringBuffs"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AddFlag struct {
	EffectRoleType      EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`           //
	ReplaceTargetType   ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"`     //
	TargetSelectParam   []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"`     //
	TargetsCount        string              `json:"targetsCount,omitempty" bson:"targetsCount"`               //
	SiftType            SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                       //
	SelectRoleType      EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`           //
	ExcludeSelf         bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`                 //
	MountEffects        []string            `json:"mountEffects,omitempty" bson:"mountEffects"`               //
	CenterEffects       []string            `json:"centerEffects,omitempty" bson:"centerEffects"`             //
	SelfDead            bool                `json:"selfDead,omitempty" bson:"selfDead"`                       //
	SiftParam           []float64           `json:"siftParam,omitempty" bson:"siftParam"`                     //
	Rate                int32               `json:"rate,omitempty" bson:"rate"`                               //
	CacheCnt            bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                       //
	NoBaseSelect        bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`               //
	FlagType            BuffFlagType        `json:"flagType,omitempty" bson:"flagType"`                       //
	Param               []int32             `json:"param,omitempty" bson:"param"`                             //
	RemoveFlagWhenClear bool                `json:"removeFlagWhenClear,omitempty" bson:"removeFlagWhenClear"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AddShield struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ShieldId          string              `json:"shieldId,omitempty" bson:"shieldId"`                   //
	EffectId          string              `json:"effectId,omitempty" bson:"effectId"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type RefractionShield struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ShieldId          string              `json:"shieldId,omitempty" bson:"shieldId"`                   //
	EffectId          string              `json:"effectId,omitempty" bson:"effectId"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type Relive struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ReliveRate        float32             `json:"reliveRate,omitempty" bson:"reliveRate"`               //
	Heal              string              `json:"heal,omitempty" bson:"heal"`                           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type Slash struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	BringBuffs        []int32             `json:"bringBuffs,omitempty" bson:"bringBuffs"`               //
	Precent           float32             `json:"precent,omitempty" bson:"precent"`                     //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type Support struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ResistDeath struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Duration          float32             `json:"duration,omitempty" bson:"duration"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type Dispel struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	BuffTypes         []BuffType          `json:"buffTypes,omitempty" bson:"buffTypes"`                 //
	Froce             bool                `json:"froce,omitempty" bson:"froce"`                         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type HpChange struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type HealWhenKill struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          []float64           `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	MyAttrFactor      AttrType            `json:"myAttrFactor,omitempty" bson:"myAttrFactor"`           //
	TargetAttrFactor  AttrType            `json:"targetAttrFactor,omitempty" bson:"targetAttrFactor"`   //
	HealPrecent       float32             `json:"healPrecent,omitempty" bson:"healPrecent"`             //
	JudgeFactor       string              `json:"judgeFactor,omitempty" bson:"judgeFactor"`             //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type TmpFixedAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	FixedType         BuffFixedDamageType `json:"fixedType,omitempty" bson:"fixedType"`                 //
	AttrType          AttrType            `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type Stealth struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	FullStealth       bool                `json:"fullStealth,omitempty" bson:"fullStealth"`             //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AttackByMountBuffLayer struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	DependBuff        int32               `json:"dependBuff,omitempty" bson:"dependBuff"`               //
	BringBuffs        []int32             `json:"bringBuffs,omitempty" bson:"bringBuffs"`               //
	AnotherBuff       int32               `json:"anotherBuff,omitempty" bson:"anotherBuff"`             //
	BuffTypes         []BuffType          `json:"buffTypes,omitempty" bson:"buffTypes"`                 //
	BuffIDs           []int32             `json:"buffIDs,omitempty" bson:"buffIDs"`                     //
	MaxCnt            int32               `json:"maxCnt,omitempty" bson:"maxCnt"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SummonBase struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	SummoId           string              `json:"summoId,omitempty" bson:"summoId"`                     //
	SummonCount       string              `json:"summonCount,omitempty" bson:"summonCount"`             //
	ShieldId          string              `json:"shieldId,omitempty" bson:"shieldId"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SummonPitfall struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	SummoId           string              `json:"summoId,omitempty" bson:"summoId"`                     //
	SummonCount       string              `json:"summonCount,omitempty" bson:"summonCount"`             //
	ShieldId          string              `json:"shieldId,omitempty" bson:"shieldId"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SelfBoom struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	AoeId             string              `json:"aoeId,omitempty" bson:"aoeId"`                         //
	NoDamage          bool                `json:"noDamage,omitempty" bson:"noDamage"`                   //
	OnlyCheckX        bool                `json:"onlyCheckX,omitempty" bson:"onlyCheckX"`               //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ConditionEffect struct {
	EffectRoleType       EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`             //
	ReplaceTargetType    ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"`       //
	TargetSelectParam    []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"`       //
	TargetsCount         string              `json:"targetsCount,omitempty" bson:"targetsCount"`                 //
	SiftType             SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                         //
	SelectRoleType       EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`             //
	ExcludeSelf          bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`                   //
	MountEffects         []string            `json:"mountEffects,omitempty" bson:"mountEffects"`                 //
	CenterEffects        []string            `json:"centerEffects,omitempty" bson:"centerEffects"`               //
	SelfDead             bool                `json:"selfDead,omitempty" bson:"selfDead"`                         //
	SiftParam            []float64           `json:"siftParam,omitempty" bson:"siftParam"`                       //
	Rate                 int32               `json:"rate,omitempty" bson:"rate"`                                 //
	CacheCnt             bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                         //
	NoBaseSelect         bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`                 //
	If                   []*EffectCondition  `json:"if,omitempty" bson:"if"`                                     //
	IfEffects            []string            `json:"ifEffects,omitempty" bson:"ifEffects"`                       //
	ElseEffects          []string            `json:"elseEffects,omitempty" bson:"elseEffects"`                   //
	ElseIf               []*EffectCondition  `json:"elseIf,omitempty" bson:"elseIf"`                             //
	ElseIfEffects        []string            `json:"elseIfEffects,omitempty" bson:"elseIfEffects"`               //
	Ifand                []*EffectCondition  `json:"ifand,omitempty" bson:"ifand"`                               //
	ElseIfand            []*EffectCondition  `json:"elseIfand,omitempty" bson:"elseIfand"`                       //
	ExtendsParentTargets bool                `json:"extendsParentTargets,omitempty" bson:"extendsParentTargets"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type EffectCondition struct {
	ConditionName string   `json:"conditionName,omitempty" bson:"conditionName"` //
	Param         []string `json:"param,omitempty" bson:"param"`                 //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ConditionEffectPair struct {
	Condition *EffectCondition `json:"condition,omitempty" bson:"condition"` //
	Effects   *BuffEffect      `json:"effects,omitempty" bson:"effects"`     //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type EffectPair struct {
	EffectId string  `json:"effectId,omitempty" bson:"effectId"` //
	Interval float32 `json:"interval,omitempty" bson:"interval"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type MountHalo struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HaloId            string              `json:"haloId,omitempty" bson:"haloId"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type LinkMagic struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	LinkId            string              `json:"linkId,omitempty" bson:"linkId"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AddAttackRange struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AddCalculate      string              `json:"addCalculate,omitempty" bson:"addCalculate"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SkillSiftTargetConfig struct {
	Type  SkillSiftTargetType `json:"type,omitempty" bson:"type"`   //
	Param []float64           `json:"param,omitempty" bson:"param"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type RemoveBuff struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	BuffId            int32               `json:"buffId,omitempty" bson:"buffId"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type RemoveHalo struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HaloName          string              `json:"haloName,omitempty" bson:"haloName"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type Vector struct {
	X float32 `json:"x,omitempty" bson:"x"` //
	Y float32 `json:"y,omitempty" bson:"y"` //
	Z float32 `json:"z,omitempty" bson:"z"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ExtraAttr struct {
	Camp  BattleRoleCamp `json:"camp,omitempty" bson:"camp"`   //
	Type  AttrType       `json:"type,omitempty" bson:"type"`   //
	Value float32        `json:"value,omitempty" bson:"value"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type KeyTimes struct {
	Keys []float64 `json:"keys,omitempty" bson:"keys"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ChangeExpRate struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Formula           string              `json:"formula,omitempty" bson:"formula"`                     //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ChangeStage struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Stage             int32               `json:"stage,omitempty" bson:"stage"`                         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type TalkEffects struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	TalkIds           []int32             `json:"talkIds,omitempty" bson:"talkIds"`                     //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type EnterGatherState struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ChangeShape struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ShapeId           int32               `json:"shapeId,omitempty" bson:"shapeId"`                     //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ForceOrderSummon struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	SummonId          string              `json:"summonId,omitempty" bson:"summonId"`                   //
	SkillId           int32               `json:"skillId,omitempty" bson:"skillId"`                     //
	UseMasterTarget   bool                `json:"useMasterTarget,omitempty" bson:"useMasterTarget"`     //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ChangeScale struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Magnification     string              `json:"magnification,omitempty" bson:"magnification"`         //
	MaxLimit          float32             `json:"maxLimit,omitempty" bson:"maxLimit"`                   //
	MinLimit          float32             `json:"minLimit,omitempty" bson:"minLimit"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ChangeNormalAttack struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	SkillId           int32               `json:"skillId,omitempty" bson:"skillId"`                     //
	Level             int32               `json:"level,omitempty" bson:"level"`                         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type StealAttr struct {
	EffectRoleType          EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`                   //
	ReplaceTargetType       ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"`             //
	TargetSelectParam       []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"`             //
	TargetsCount            string              `json:"targetsCount,omitempty" bson:"targetsCount"`                       //
	SiftType                SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                               //
	SelectRoleType          EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`                   //
	ExcludeSelf             bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`                         //
	MountEffects            []string            `json:"mountEffects,omitempty" bson:"mountEffects"`                       //
	CenterEffects           []string            `json:"centerEffects,omitempty" bson:"centerEffects"`                     //
	SelfDead                bool                `json:"selfDead,omitempty" bson:"selfDead"`                               //
	SiftParam               []float64           `json:"siftParam,omitempty" bson:"siftParam"`                             //
	Rate                    int32               `json:"rate,omitempty" bson:"rate"`                                       //
	CacheCnt                bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                               //
	NoBaseSelect            bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`                       //
	AttrType                AttrType            `json:"attrType,omitempty" bson:"attrType"`                               //
	Add                     string              `json:"add,omitempty" bson:"add"`                                         //
	Mult                    string              `json:"mult,omitempty" bson:"mult"`                                       //
	EffectTime              int32               `json:"effectTime,omitempty" bson:"effectTime"`                           //
	Min                     string              `json:"min,omitempty" bson:"min"`                                         //
	Max                     string              `json:"max,omitempty" bson:"max"`                                         //
	TotalAddMax             string              `json:"totalAddMax,omitempty" bson:"totalAddMax"`                         //
	ForceEffectSelf         bool                `json:"forceEffectSelf,omitempty" bson:"forceEffectSelf"`                 //
	ExtraAddRuntimeAttrType AttrType            `json:"extraAddRuntimeAttrType,omitempty" bson:"extraAddRuntimeAttrType"` //
	ExtraAddRuntimeValue    string              `json:"extraAddRuntimeValue,omitempty" bson:"extraAddRuntimeValue"`       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ExchangeHpPrecent struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SkillEnterCd struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	SkillId           int32               `json:"skillId,omitempty" bson:"skillId"`                     //
	SkillType         SkillTagType        `json:"skillType,omitempty" bson:"skillType"`                 //
	IsPrecent         bool                `json:"isPrecent,omitempty" bson:"isPrecent"`                 //
	SubValue          float32             `json:"subValue,omitempty" bson:"subValue"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ComboSkill struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	SkillId           int32               `json:"skillId,omitempty" bson:"skillId"`                     //
	ComboRate         float32             `json:"comboRate,omitempty" bson:"comboRate"`                 //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SubShield struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Precent           float32             `json:"precent,omitempty" bson:"precent"`                     //
	ForceRemove       []string            `json:"forceRemove,omitempty" bson:"forceRemove"`             //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenEnterLevel struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenEnterBattle struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenDamageMoreThan struct {
	TargetVal float32 `json:"targetVal,omitempty" bson:"targetVal"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenEnterStage struct {
	TargetStage int32 `json:"targetStage,omitempty" bson:"targetStage"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenStageTimeMoreThan struct {
	TargetTime  float32 `json:"targetTime,omitempty" bson:"targetTime"`   //
	TargetStage int32   `json:"targetStage,omitempty" bson:"targetStage"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenDead struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenBattleTimeMoreThan struct {
	TargetTime float32 `json:"targetTime,omitempty" bson:"targetTime"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenNormalAttack struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenBehitMoreThan struct {
	TargetCount int32 `json:"targetCount,omitempty" bson:"targetCount"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenLockBySkill struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenCoolDown struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenEnterBattleAndCoolDown struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenBeHit struct {
	SpecialHit string  `json:"specialHit,omitempty" bson:"specialHit"` //
	Precent    float32 `json:"precent,omitempty" bson:"precent"`       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type QuickSettleTypeBuff struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Type              BuffType            `json:"type,omitempty" bson:"type"`                           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenRoleAddFlag struct {
	FlagType  BuffFlagType   `json:"flagType,omitempty" bson:"flagType"`   //
	CheckType BattleRoleCamp `json:"checkType,omitempty" bson:"checkType"` //
	CheckSelf bool           `json:"checkSelf,omitempty" bson:"checkSelf"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type DeductHp struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenSupporting struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenSupportingEnd struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AddHelpValue struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Value             float32             `json:"value,omitempty" bson:"value"`                         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ChangeCameraViewScale struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Time              float32             `json:"time,omitempty" bson:"time"`                           //
	Scale             float32             `json:"scale,omitempty" bson:"scale"`                         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type LifeLink struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Radius            float32             `json:"radius,omitempty" bson:"radius"`                       //
	CampType          BattleRoleCamp      `json:"campType,omitempty" bson:"campType"`                   //
	ShareCount        int32               `json:"shareCount,omitempty" bson:"shareCount"`               //
	SharePrecent      float32             `json:"sharePrecent,omitempty" bson:"sharePrecent"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type RandomWeightEffect struct {
	EffectRoleType    EffectRoleType            `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType         `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string                  `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string                    `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType       `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType            `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                      `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string                  `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string                  `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                      `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64                 `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32                     `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                      `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                      `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	WeightEffects     []*RandomWeightEffectUnit `json:"weightEffects,omitempty" bson:"weightEffects"`         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type RandomWeightEffectUnit struct {
	EffectName string  `json:"effectName,omitempty" bson:"effectName"` //
	Weight     float32 `json:"weight,omitempty" bson:"weight"`         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SkillAttRandom struct {
	AttrType       AttrType `json:"attrType,omitempty" bson:"attrType"`             //
	AttrMaxValue   float32  `json:"attrMaxValue,omitempty" bson:"attrMaxValue"`     //
	AttrMinValue   float32  `json:"attrMinValue,omitempty" bson:"attrMinValue"`     //
	Weights        []int32  `json:"weights,omitempty" bson:"weights"`               //
	AttrMaxFormula string   `json:"attrMaxFormula,omitempty" bson:"attrMaxFormula"` //
	AttrMinFormula string   `json:"attrMinFormula,omitempty" bson:"attrMinFormula"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type CameraShake struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Time              float32             `json:"time,omitempty" bson:"time"`                           //
	Times             int32               `json:"times,omitempty" bson:"times"`                         //
	XAmplitude        float32             `json:"xAmplitude,omitempty" bson:"xAmplitude"`               //
	YAmplitude        float32             `json:"yAmplitude,omitempty" bson:"yAmplitude"`               //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type SkillEffectViewPair struct {
	EffectName string  `json:"effectName,omitempty" bson:"effectName"` //
	Delay      float32 `json:"delay,omitempty" bson:"delay"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenBuffAdd struct {
	BuffTypes []BuffType      `json:"buffTypes,omitempty" bson:"buffTypes"` //
	BuffId    decimal.Decimal `json:"buffId,omitempty" bson:"buffId"`       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AddRoleSubDamage struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ChangeBuffLayerWithParam struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ParamKey          string              `json:"paramKey,omitempty" bson:"paramKey"`                   //
	BuffId            int32               `json:"buffId,omitempty" bson:"buffId"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenCallConnection struct {
	ConnectId int32 `json:"connectId,omitempty" bson:"connectId"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type CallConnection struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ConnectId         int32               `json:"connectId,omitempty" bson:"connectId"`                 //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type HitStun struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HitStunID         int32               `json:"hitStunID,omitempty" bson:"hitStunID"`                 //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type FaceMove struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	MoveId            string              `json:"moveId,omitempty" bson:"moveId"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type HitJudge struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ComboId           int32               `json:"comboId,omitempty" bson:"comboId"`                     //
	ComboBackUp       int32               `json:"comboBackUp,omitempty" bson:"comboBackUp"`             //
	JudgeCD           bool                `json:"judgeCD,omitempty" bson:"judgeCD"`                     //
	RefindTarget      string              `json:"refindTarget,omitempty" bson:"refindTarget"`           //
	UseBFirst         bool                `json:"useBFirst,omitempty" bson:"useBFirst"`                 //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type Grab struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	GrabAction        string              `json:"grabAction,omitempty" bson:"grabAction"`               //
	GrabHitAction     string              `json:"grabHitAction,omitempty" bson:"grabHitAction"`         //
	GrabType          GrabPointType       `json:"grabType,omitempty" bson:"grabType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	FaceDir           string              `json:"faceDir,omitempty" bson:"faceDir"`                     //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type Release struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ChangeGrabTargetAction struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Action            string              `json:"action,omitempty" bson:"action"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ConnectBulletTime struct {
	KeyTime  float32 `json:"keyTime,omitempty" bson:"keyTime"`   //
	Scale    float32 `json:"scale,omitempty" bson:"scale"`       //
	KeepTime float32 `json:"keepTime,omitempty" bson:"keepTime"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenEffectCountReached struct {
	EffectName  string `json:"effectName,omitempty" bson:"effectName"`   //
	EffectName2 string `json:"effectName2,omitempty" bson:"effectName2"` //
	Count       int32  `json:"count,omitempty" bson:"count"`             //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenDropDown struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type AfterSendHeal struct {
	CheckCount int32  `json:"checkCount,omitempty" bson:"checkCount"` //
	NotToMe    bool   `json:"notToMe,omitempty" bson:"notToMe"`       //
	SpecialHit string `json:"specialHit,omitempty" bson:"specialHit"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenCounter struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type CounterAttack struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ComboAttack struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenCombo struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenSelfEnterHitStun struct {
	CheckCount int32 `json:"checkCount,omitempty" bson:"checkCount"` //
	EndStun    bool  `json:"endStun,omitempty" bson:"endStun"`       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type BeforeDamage struct {
	WhiteList                []int32      `json:"whiteList,omitempty" bson:"whiteList"`                               //
	SkillType                SkillTagType `json:"skillType,omitempty" bson:"skillType"`                               //
	CancleWhenSkillEnd       bool         `json:"cancleWhenSkillEnd,omitempty" bson:"cancleWhenSkillEnd"`             //
	CancleWhenSkillEffectEnd bool         `json:"cancleWhenSkillEffectEnd,omitempty" bson:"cancleWhenSkillEffectEnd"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ConnectionTargetSelect struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenAllShieldDispose struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenAddShield struct {
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenRoleOnBoard struct {
	CheckLeave bool `json:"checkLeave,omitempty" bson:"checkLeave"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenBloodSuck struct {
	Count int32 `json:"count,omitempty" bson:"count"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenComboGroupEffect struct {
	GroupId  decimal.Decimal `json:"groupId,omitempty" bson:"groupId"`   //
	IsCancle bool            `json:"isCancle,omitempty" bson:"isCancle"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type ComboFontShowData struct {
	Show      string           `json:"show,omitempty" bson:"show"`           //
	KeyTime   float32          `json:"keyTime,omitempty" bson:"keyTime"`     //
	CenterPos []float64        `json:"centerPos,omitempty" bson:"centerPos"` //
	BaseSacle float32          `json:"baseSacle,omitempty" bson:"baseSacle"` //
	ShowType  ComboFntShowType `json:"showType,omitempty" bson:"showType"`   //
	ShowParam []float64        `json:"showParam,omitempty" bson:"showParam"` //
	Delay     float32          `json:"delay,omitempty" bson:"delay"`         //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenAttrPerChanged struct {
	AttrType    AttrType `json:"attrType,omitempty" bson:"attrType"`       //
	IsPrecent   bool     `json:"isPrecent,omitempty" bson:"isPrecent"`     //
	MaxAttType  AttrType `json:"maxAttType,omitempty" bson:"maxAttType"`   //
	CompareFunc string   `json:"compareFunc,omitempty" bson:"compareFunc"` //
	TargetValue float32  `json:"targetValue,omitempty" bson:"targetValue"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type MoveEndAction struct {
	Action string  `json:"action,omitempty" bson:"action"` //
	Keep   float32 `json:"keep,omitempty" bson:"keep"`     //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenSelfAreaChange struct {
	Areas []BattleAreaType `json:"areas,omitempty" bson:"areas"` //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type HarmImmunity struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ExcludeheroIDs    []decimal.Decimal   `json:"excludeheroIDs,omitempty" bson:"excludeheroIDs"`       //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type RoleShowFlash struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Min               float32             `json:"min,omitempty" bson:"min"`                             //
	Max               float32             `json:"max,omitempty" bson:"max"`                             //
	BaseColor         string              `json:"baseColor,omitempty" bson:"baseColor"`                 //
	Interval          float32             `json:"interval,omitempty" bson:"interval"`                   //
}

// ##改变动态属性用
// 技能表-拳皇/SkillPoly
type WhenHurtDamageMoreThan struct {
	Formula     string  `json:"formula,omitempty" bson:"formula"`         //
	CompareFunc string  `json:"compareFunc,omitempty" bson:"compareFunc"` //
	CacheTime   float32 `json:"cacheTime,omitempty" bson:"cacheTime"`     //
}

// ###LD
// 技能表-拳皇/SkillPoly
type LDChangeGlobalAttr struct {
	AttrType       LDAttrType      `json:"attrType,omitempty" bson:"attrType"`             //
	Add            string          `json:"add,omitempty" bson:"add"`                       //
	Mult           string          `json:"mult,omitempty" bson:"mult"`                     //
	EffectTime     int32           `json:"effectTime,omitempty" bson:"effectTime"`         //
	Min            float32         `json:"min,omitempty" bson:"min"`                       //
	Max            float32         `json:"max,omitempty" bson:"max"`                       //
	NeedKeep       bool            `json:"needKeep,omitempty" bson:"needKeep"`             //
	GlobalAttrType LDAttrType      `json:"globalAttrType,omitempty" bson:"globalAttrType"` //
	CampType       BattleRoleCamp  `json:"campType,omitempty" bson:"campType"`             //
	MemberType     LevelMemberCamp `json:"memberType,omitempty" bson:"memberType"`         //
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenWeatherIs struct {
	Weather WeatherTimeType `json:"weather,omitempty" bson:"weather"` //
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenWeatherNot struct {
	Weather WeatherTimeType `json:"weather,omitempty" bson:"weather"` //
}

// ###LD
// 技能表-拳皇/SkillPoly
type ReqServerEffect struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	EffectId          int32               `json:"effectId,omitempty" bson:"effectId"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
}

// ###LD
// 技能表-拳皇/SkillPoly
type ReqRoleEnchant struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenRoleEnhant struct {
	Force        bool `json:"force,omitempty" bson:"force"`               //
	CheckFailded bool `json:"checkFailded,omitempty" bson:"checkFailded"` //
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenRaceCountChange struct {
	MemberType LevelMemberCamp `json:"memberType,omitempty" bson:"memberType"` //
	Race       RoleRaceType    `json:"race,omitempty" bson:"race"`             //
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenHeroMove struct {
}

// ###LD
// 技能表-拳皇/SkillPoly
type PathAoe struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AoeId             string              `json:"aoeId,omitempty" bson:"aoeId"`                         //
	NeighbourCount    int32               `json:"neighbourCount,omitempty" bson:"neighbourCount"`       //
	AoeCount          string              `json:"aoeCount,omitempty" bson:"aoeCount"`                   //
	Random            bool                `json:"random,omitempty" bson:"random"`                       //
	Interval          float32             `json:"interval,omitempty" bson:"interval"`                   //
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenKillCountChange struct {
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenGirdParamChanged struct {
	ListenKey string `json:"listenKey,omitempty" bson:"listenKey"` //
	Force     bool   `json:"force,omitempty" bson:"force"`         //
}

// ###LD
// 技能表-拳皇/SkillPoly
type LDChangeBuffLayerWithParam struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ParamKey          string              `json:"paramKey,omitempty" bson:"paramKey"`                   //
	BuffId            int32               `json:"buffId,omitempty" bson:"buffId"`                       //
}

// ###LD
// 技能表-拳皇/SkillPoly
type AutoEnchant struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// ###LD
// 技能表-拳皇/SkillPoly
type SaleSelf struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenHeroCountLessThan struct {
	HeroId       int32 `json:"heroId,omitempty" bson:"heroId"`             //
	MaxHeroCount int32 `json:"maxHeroCount,omitempty" bson:"maxHeroCount"` //
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenClickEnchantBtn struct {
}

// ###LD
// 技能表-拳皇/SkillPoly
type WhenDiceShapeChange struct {
}

// ###LD
// 技能表-拳皇/SkillPoly
type MountBuffWithDiceShape struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	ShapeId           int32               `json:"shapeId,omitempty" bson:"shapeId"`                     //
	BuffId            int32               `json:"buffId,omitempty" bson:"buffId"`                       //
}

// ###LD
// 技能表-拳皇/SkillPoly
type SelfRushShock struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	RushCount         int32               `json:"rushCount,omitempty" bson:"rushCount"`                 //
	ShockId           string              `json:"shockId,omitempty" bson:"shockId"`                     //
}

// ###LD
// 技能表-拳皇/SkillPoly
type AttackBack struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Distance          float32             `json:"distance,omitempty" bson:"distance"`                   //
	Time              float32             `json:"time,omitempty" bson:"time"`                           //
}

// ###LD
// 技能表-拳皇/SkillPoly
type LDWhenAttrChanged struct {
	AttrType    LDAttrType `json:"attrType,omitempty" bson:"attrType"`       //
	CompareSign string     `json:"compareSign,omitempty" bson:"compareSign"` //
	TargetValue float32    `json:"targetValue,omitempty" bson:"targetValue"` //
}

// #改变属性
// 技能表-拳皇/SkillPoly
type LDChangeAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AttrType          LDAttrType          `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
	EffectTime        int32               `json:"effectTime,omitempty" bson:"effectTime"`               //
	Min               float32             `json:"min,omitempty" bson:"min"`                             //
	Max               float32             `json:"max,omitempty" bson:"max"`                             //
	NeedKeep          bool                `json:"needKeep,omitempty" bson:"needKeep"`                   //
}

// #改变基础属性
// 技能表-拳皇/SkillPoly
type LDChangeBaseAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AttrType          LDAttrType          `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
	EffectTime        int32               `json:"effectTime,omitempty" bson:"effectTime"`               //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type LDChangeGrowAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AttrType          LDAttrType          `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
	EffectTime        int32               `json:"effectTime,omitempty" bson:"effectTime"`               //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type LDForceChangAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AttrType          LDAttrType          `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
	EffectTime        int32               `json:"effectTime,omitempty" bson:"effectTime"`               //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type LDChangeScale struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Magnification     string              `json:"magnification,omitempty" bson:"magnification"`         //
	MaxLimit          float32             `json:"maxLimit,omitempty" bson:"maxLimit"`                   //
	MinLimit          float32             `json:"minLimit,omitempty" bson:"minLimit"`                   //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type LDTmpFixedAttr struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	FixedType         BuffFixedDamageType `json:"fixedType,omitempty" bson:"fixedType"`                 //
	AttrType          LDAttrType          `json:"attrType,omitempty" bson:"attrType"`                   //
	Add               string              `json:"add,omitempty" bson:"add"`                             //
	Mult              string              `json:"mult,omitempty" bson:"mult"`                           //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type LDListenEnemyAttr struct {
	TargetAttr           LDAttrType `json:"targetAttr,omitempty" bson:"targetAttr"`                     //
	ExculdeSelf          bool       `json:"exculdeSelf,omitempty" bson:"exculdeSelf"`                   //
	MaxAttr              LDAttrType `json:"maxAttr,omitempty" bson:"maxAttr"`                           //
	ComparePrecent       bool       `json:"comparePrecent,omitempty" bson:"comparePrecent"`             //
	TargetValue          float32    `json:"targetValue,omitempty" bson:"targetValue"`                   //
	CompareFunc          string     `json:"compareFunc,omitempty" bson:"compareFunc"`                   //
	CheckSelf            bool       `json:"checkSelf,omitempty" bson:"checkSelf"`                       //
	CheckTargetCampRoles bool       `json:"checkTargetCampRoles,omitempty" bson:"checkTargetCampRoles"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type ChangeCameraFocus struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type PlayUIEffect struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	EffectName        string              `json:"effectName,omitempty" bson:"effectName"`               //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenRoleEnter struct {
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type AddBuffTypeTime struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	BuffType          BuffType            `json:"buffType,omitempty" bson:"buffType"`                   //
	AddValue          float32             `json:"addValue,omitempty" bson:"addValue"`                   //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type ChangeCounterCDParam struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Value             float32             `json:"value,omitempty" bson:"value"`                         //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type ImmunitySkills struct {
	EffectRoleType       EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`             //
	ReplaceTargetType    ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"`       //
	TargetSelectParam    []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"`       //
	TargetsCount         string              `json:"targetsCount,omitempty" bson:"targetsCount"`                 //
	SiftType             SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                         //
	SelectRoleType       EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`             //
	ExcludeSelf          bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`                   //
	MountEffects         []string            `json:"mountEffects,omitempty" bson:"mountEffects"`                 //
	CenterEffects        []string            `json:"centerEffects,omitempty" bson:"centerEffects"`               //
	SelfDead             bool                `json:"selfDead,omitempty" bson:"selfDead"`                         //
	SiftParam            []float64           `json:"siftParam,omitempty" bson:"siftParam"`                       //
	Rate                 int32               `json:"rate,omitempty" bson:"rate"`                                 //
	CacheCnt             bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                         //
	NoBaseSelect         bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`                 //
	SkillIds             []decimal.Decimal   `json:"skillIds,omitempty" bson:"skillIds"`                         //
	ImmuintyTriggerSkill bool                `json:"immuintyTriggerSkill,omitempty" bson:"immuintyTriggerSkill"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenTypeSkillEnd struct {
	SkillType SkillTagType `json:"skillType,omitempty" bson:"skillType"` //
	Delay     float32      `json:"delay,omitempty" bson:"delay"`         //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenBattleStart struct {
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type AddEffectParam struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Type              EffectAddParam      `json:"type,omitempty" bson:"type"`                           //
	AddValue          float32             `json:"addValue,omitempty" bson:"addValue"`                   //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type ChangeToPuppet struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type AddPuppetShield struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type PuppetShieldBroken struct {
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenTotalDmgCompare struct {
	SpecialHit  string `json:"specialHit,omitempty" bson:"specialHit"`   //
	TargetValue string `json:"targetValue,omitempty" bson:"targetValue"` //
	Compare     string `json:"compare,omitempty" bson:"compare"`         //
	SkillId     int32  `json:"skillId,omitempty" bson:"skillId"`         //
	BuffId      int32  `json:"buffId,omitempty" bson:"buffId"`           //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenTeamMateBehurt struct {
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenCounterSuccess struct {
	Count int32 `json:"count,omitempty" bson:"count"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenPropertyNeedTransfer struct {
	AttackTrigger bool `json:"attackTrigger,omitempty" bson:"attackTrigger"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type TransferProperities struct {
	EffectRoleType    EffectRoleType       `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType    `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string             `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string               `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType  `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType       `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                 `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string             `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string             `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                 `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64            `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32                `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                 `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                 `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	TargetType        string               `json:"targetType,omitempty" bson:"targetType"`               //
	CalculateType     DmgCalculateRateType `json:"calculateType,omitempty" bson:"calculateType"`         //
	CompareValue      string               `json:"compareValue,omitempty" bson:"compareValue"`           //
	Compare           string               `json:"compare,omitempty" bson:"compare"`                     //
	Limit             string               `json:"limit,omitempty" bson:"limit"`                         //
	AddValue          string               `json:"addValue,omitempty" bson:"addValue"`                   //
	TargetAttrType    AttrType             `json:"targetAttrType,omitempty" bson:"targetAttrType"`       //
	TargetValueMax    float32              `json:"targetValueMax,omitempty" bson:"targetValueMax"`       //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenRelive struct {
	CanReliveCnt int32   `json:"canReliveCnt,omitempty" bson:"canReliveCnt"` //
	ReliveRate   float32 `json:"reliveRate,omitempty" bson:"reliveRate"`     //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type AddReliveCnt struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	AddCnt            int32               `json:"addCnt,omitempty" bson:"addCnt"`                       //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type SkillParamPair struct {
	Key   string `json:"key,omitempty" bson:"key"`     //
	Value string `json:"value,omitempty" bson:"value"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type AddSkillParam struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	SkillId           int32               `json:"skillId,omitempty" bson:"skillId"`                     //
	ParamMap          []*SkillParamPair   `json:"paramMap,omitempty" bson:"paramMap"`                   //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type SupportSummon struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	SummoId           string              `json:"summoId,omitempty" bson:"summoId"`                     //
	SummonCount       string              `json:"summonCount,omitempty" bson:"summonCount"`             //
	ShieldId          string              `json:"shieldId,omitempty" bson:"shieldId"`                   //
	EnterType         SummonAppearType    `json:"enterType,omitempty" bson:"enterType"`                 //
	SkillIdSqunence   []int32             `json:"skillIdSqunence,omitempty" bson:"skillIdSqunence"`     //
	SummonHeroId      int32               `json:"summonHeroId,omitempty" bson:"summonHeroId"`           //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type CampUseSkill struct {
	CampType   string       `json:"campType,omitempty" bson:"campType"`     //
	SkillType  SkillTagType `json:"skillType,omitempty" bson:"skillType"`   //
	Count      int32        `json:"count,omitempty" bson:"count"`           //
	SkillIds   []int32      `json:"skillIds,omitempty" bson:"skillIds"`     //
	SkillLevel int32        `json:"skillLevel,omitempty" bson:"skillLevel"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type CampMemberCntReach struct {
	CampType string `json:"campType,omitempty" bson:"campType"` //
	Count    int32  `json:"count,omitempty" bson:"count"`       //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenFirstMove struct {
	MoveType string `json:"moveType,omitempty" bson:"moveType"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenAssisKllOrKill struct {
	Count int32 `json:"count,omitempty" bson:"count"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenSupportLeave struct {
	HeroId int32 `json:"heroId,omitempty" bson:"heroId"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenCampRoleDead struct {
	CheckType string `json:"checkType,omitempty" bson:"checkType"` //
	DeadNum   int32  `json:"deadNum,omitempty" bson:"deadNum"`     //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type AnyOneAttrReach struct {
	AttrType    AttrType `json:"attrType,omitempty" bson:"attrType"`       //
	MaxAttrType AttrType `json:"maxAttrType,omitempty" bson:"maxAttrType"` //
	CompareFun  string   `json:"compareFun,omitempty" bson:"compareFun"`   //
	TargetValue float32  `json:"targetValue,omitempty" bson:"targetValue"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type Boomerang struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	HurtType          DamageType          `json:"hurtType,omitempty" bson:"hurtType"`                   //
	Calculator        string              `json:"calculator,omitempty" bson:"calculator"`               //
	SubType           SubDamageType       `json:"subType,omitempty" bson:"subType"`                     //
	Subparam          string              `json:"subparam,omitempty" bson:"subparam"`                   //
	SkillExtra        float32             `json:"skillExtra,omitempty" bson:"skillExtra"`               //
	Max               string              `json:"max,omitempty" bson:"max"`                             //
	Spatter           string              `json:"spatter,omitempty" bson:"spatter"`                     //
	SubMax            string              `json:"subMax,omitempty" bson:"subMax"`                       //
	BoomerangId       string              `json:"boomerangId,omitempty" bson:"boomerangId"`             //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type BoomerangCountEffects struct {
	Count   int32    `json:"count,omitempty" bson:"count"`     //
	Effects []string `json:"effects,omitempty" bson:"effects"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type LimitDamage struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Min               string              `json:"min,omitempty" bson:"min"`                             //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type LimitTotalDamage struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	Limit             string              `json:"limit,omitempty" bson:"limit"`                         //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type WhenGetDispel struct {
	Result bool `json:"result,omitempty" bson:"result"` //
}

// #改变培养属性
// 技能表-拳皇/SkillPoly
type AddBuffTime struct {
	EffectRoleType    EffectRoleType      `json:"effectRoleType,omitempty" bson:"effectRoleType"`       //
	ReplaceTargetType ReplaceTargetType   `json:"replaceTargetType,omitempty" bson:"replaceTargetType"` //
	TargetSelectParam []string            `json:"targetSelectParam,omitempty" bson:"targetSelectParam"` //
	TargetsCount      string              `json:"targetsCount,omitempty" bson:"targetsCount"`           //
	SiftType          SkillSiftTargetType `json:"siftType,omitempty" bson:"siftType"`                   //
	SelectRoleType    EffectRoleType      `json:"selectRoleType,omitempty" bson:"selectRoleType"`       //
	ExcludeSelf       bool                `json:"excludeSelf,omitempty" bson:"excludeSelf"`             //
	MountEffects      []string            `json:"mountEffects,omitempty" bson:"mountEffects"`           //
	CenterEffects     []string            `json:"centerEffects,omitempty" bson:"centerEffects"`         //
	SelfDead          bool                `json:"selfDead,omitempty" bson:"selfDead"`                   //
	SiftParam         []float64           `json:"siftParam,omitempty" bson:"siftParam"`                 //
	Rate              int32               `json:"rate,omitempty" bson:"rate"`                           //
	CacheCnt          bool                `json:"cacheCnt,omitempty" bson:"cacheCnt"`                   //
	NoBaseSelect      bool                `json:"noBaseSelect,omitempty" bson:"noBaseSelect"`           //
	BuffId            int32               `json:"buffId,omitempty" bson:"buffId"`                       //
	KeepTimeAdd       string              `json:"keepTimeAdd,omitempty" bson:"keepTimeAdd"`             //
}

func init() {
	excel.AddEntryNewFunc("*SkillAction", newSkillAction)
	excel.AddEntryNewFunc("*SkillPerfrom", newSkillPerfrom)
	excel.AddEntryNewFunc("*RougeSkillCondi", newRougeSkillCondi)
	excel.AddEntryNewFunc("*EffectGroup", newEffectGroup)
	excel.AddEntryNewFunc("*SkillSequenceTrigger", newSkillSequenceTrigger)
	excel.AddEntryNewFunc("*SkillExcuteCheck", newSkillExcuteCheck)
	excel.AddEntryNewFunc("*BuffEffect", newBuffEffect)
	excel.AddEntryNewFunc("*Attack", newAttack)
	excel.AddEntryNewFunc("*Heal", newHeal)
	excel.AddEntryNewFunc("*MeleeAttack", newMeleeAttack)
	excel.AddEntryNewFunc("*Shoot", newShoot)
	excel.AddEntryNewFunc("*Strafing", newStrafing)
	excel.AddEntryNewFunc("*ProcessBuffEffect", newProcessBuffEffect)
	excel.AddEntryNewFunc("*SummonAttrExtendsPair", newSummonAttrExtendsPair)
	excel.AddEntryNewFunc("*DOT", newDOT)
	excel.AddEntryNewFunc("*HOT", newHOT)
	excel.AddEntryNewFunc("*AOE", newAOE)
	excel.AddEntryNewFunc("*Move", newMove)
	excel.AddEntryNewFunc("*AvoidMove", newAvoidMove)
	excel.AddEntryNewFunc("*ShockWave", newShockWave)
	excel.AddEntryNewFunc("*Dizzy", newDizzy)
	excel.AddEntryNewFunc("*Snare", newSnare)
	excel.AddEntryNewFunc("*Stun", newStun)
	excel.AddEntryNewFunc("*Freeze", newFreeze)
	excel.AddEntryNewFunc("*Frozen", newFrozen)
	excel.AddEntryNewFunc("*Sleep", newSleep)
	excel.AddEntryNewFunc("*Silence", newSilence)
	excel.AddEntryNewFunc("*Charm", newCharm)
	excel.AddEntryNewFunc("*Disarm", newDisarm)
	excel.AddEntryNewFunc("*Taunt", newTaunt)
	excel.AddEntryNewFunc("*Imprison", newImprison)
	excel.AddEntryNewFunc("*ChangeAttr", newChangeAttr)
	excel.AddEntryNewFunc("*ChangeBaseAttr", newChangeBaseAttr)
	excel.AddEntryNewFunc("*ChangeGrowAttr", newChangeGrowAttr)
	excel.AddEntryNewFunc("*ChangeSpecialAttr", newChangeSpecialAttr)
	excel.AddEntryNewFunc("*ForceChangAttr", newForceChangAttr)
	excel.AddEntryNewFunc("*ChangeGlobalAttr", newChangeGlobalAttr)
	excel.AddEntryNewFunc("*BuffTrigger", newBuffTrigger)
	excel.AddEntryNewFunc("*WhenBeHurt", newWhenBeHurt)
	excel.AddEntryNewFunc("*WhenDeadlyHurt", newWhenDeadlyHurt)
	excel.AddEntryNewFunc("*RandomTrigger", newRandomTrigger)
	excel.AddEntryNewFunc("*AttackAfterEffect", newAttackAfterEffect)
	excel.AddEntryNewFunc("*ListenEnemyAttr", newListenEnemyAttr)
	excel.AddEntryNewFunc("*WhenKill", newWhenKill)
	excel.AddEntryNewFunc("*WillBuildDamage", newWillBuildDamage)
	excel.AddEntryNewFunc("*WhenBeHurtDelay", newWhenBeHurtDelay)
	excel.AddEntryNewFunc("*WhenBackAttack", newWhenBackAttack)
	excel.AddEntryNewFunc("*WhenWillBeHurt", newWhenWillBeHurt)
	excel.AddEntryNewFunc("*WhenShieldBroken", newWhenShieldBroken)
	excel.AddEntryNewFunc("*WhenRefractionBroken", newWhenRefractionBroken)
	excel.AddEntryNewFunc("*WhenAfterBackAttack", newWhenAfterBackAttack)
	excel.AddEntryNewFunc("*WhenUseSkill", newWhenUseSkill)
	excel.AddEntryNewFunc("*WhenWillNormalAttack", newWhenWillNormalAttack)
	excel.AddEntryNewFunc("*WhenSummonDispear", newWhenSummonDispear)
	excel.AddEntryNewFunc("*AfterSendDamage", newAfterSendDamage)
	excel.AddEntryNewFunc("*AfterSendNormalDamage", newAfterSendNormalDamage)
	excel.AddEntryNewFunc("*WhenAttacherDead", newWhenAttacherDead)
	excel.AddEntryNewFunc("*WhenRoleDead", newWhenRoleDead)
	excel.AddEntryNewFunc("*WhenLayerMoreThan", newWhenLayerMoreThan)
	excel.AddEntryNewFunc("*WhenAttrChanged", newWhenAttrChanged)
	excel.AddEntryNewFunc("*WhenCampAttrChanged", newWhenCampAttrChanged)
	excel.AddEntryNewFunc("*SelfCenterAoe", newSelfCenterAoe)
	excel.AddEntryNewFunc("*BuffPair", newBuffPair)
	excel.AddEntryNewFunc("*Clash", newClash)
	excel.AddEntryNewFunc("*MountBuff", newMountBuff)
	excel.AddEntryNewFunc("*SubBuff", newSubBuff)
	excel.AddEntryNewFunc("*CounterInjury", newCounterInjury)
	excel.AddEntryNewFunc("*AddFlag", newAddFlag)
	excel.AddEntryNewFunc("*AddShield", newAddShield)
	excel.AddEntryNewFunc("*RefractionShield", newRefractionShield)
	excel.AddEntryNewFunc("*Relive", newRelive)
	excel.AddEntryNewFunc("*Slash", newSlash)
	excel.AddEntryNewFunc("*Support", newSupport)
	excel.AddEntryNewFunc("*ResistDeath", newResistDeath)
	excel.AddEntryNewFunc("*Dispel", newDispel)
	excel.AddEntryNewFunc("*HpChange", newHpChange)
	excel.AddEntryNewFunc("*HealWhenKill", newHealWhenKill)
	excel.AddEntryNewFunc("*TmpFixedAttr", newTmpFixedAttr)
	excel.AddEntryNewFunc("*Stealth", newStealth)
	excel.AddEntryNewFunc("*AttackByMountBuffLayer", newAttackByMountBuffLayer)
	excel.AddEntryNewFunc("*SummonBase", newSummonBase)
	excel.AddEntryNewFunc("*SummonPitfall", newSummonPitfall)
	excel.AddEntryNewFunc("*SelfBoom", newSelfBoom)
	excel.AddEntryNewFunc("*ConditionEffect", newConditionEffect)
	excel.AddEntryNewFunc("*EffectCondition", newEffectCondition)
	excel.AddEntryNewFunc("*ConditionEffectPair", newConditionEffectPair)
	excel.AddEntryNewFunc("*EffectPair", newEffectPair)
	excel.AddEntryNewFunc("*MountHalo", newMountHalo)
	excel.AddEntryNewFunc("*LinkMagic", newLinkMagic)
	excel.AddEntryNewFunc("*AddAttackRange", newAddAttackRange)
	excel.AddEntryNewFunc("*SkillSiftTargetConfig", newSkillSiftTargetConfig)
	excel.AddEntryNewFunc("*RemoveBuff", newRemoveBuff)
	excel.AddEntryNewFunc("*RemoveHalo", newRemoveHalo)
	excel.AddEntryNewFunc("*Vector", newVector)
	excel.AddEntryNewFunc("*ExtraAttr", newExtraAttr)
	excel.AddEntryNewFunc("*KeyTimes", newKeyTimes)
	excel.AddEntryNewFunc("*ChangeExpRate", newChangeExpRate)
	excel.AddEntryNewFunc("*ChangeStage", newChangeStage)
	excel.AddEntryNewFunc("*TalkEffects", newTalkEffects)
	excel.AddEntryNewFunc("*EnterGatherState", newEnterGatherState)
	excel.AddEntryNewFunc("*ChangeShape", newChangeShape)
	excel.AddEntryNewFunc("*ForceOrderSummon", newForceOrderSummon)
	excel.AddEntryNewFunc("*ChangeScale", newChangeScale)
	excel.AddEntryNewFunc("*ChangeNormalAttack", newChangeNormalAttack)
	excel.AddEntryNewFunc("*StealAttr", newStealAttr)
	excel.AddEntryNewFunc("*ExchangeHpPrecent", newExchangeHpPrecent)
	excel.AddEntryNewFunc("*SkillEnterCd", newSkillEnterCd)
	excel.AddEntryNewFunc("*ComboSkill", newComboSkill)
	excel.AddEntryNewFunc("*SubShield", newSubShield)
	excel.AddEntryNewFunc("*WhenEnterLevel", newWhenEnterLevel)
	excel.AddEntryNewFunc("*WhenEnterBattle", newWhenEnterBattle)
	excel.AddEntryNewFunc("*WhenDamageMoreThan", newWhenDamageMoreThan)
	excel.AddEntryNewFunc("*WhenEnterStage", newWhenEnterStage)
	excel.AddEntryNewFunc("*WhenStageTimeMoreThan", newWhenStageTimeMoreThan)
	excel.AddEntryNewFunc("*WhenDead", newWhenDead)
	excel.AddEntryNewFunc("*WhenBattleTimeMoreThan", newWhenBattleTimeMoreThan)
	excel.AddEntryNewFunc("*WhenNormalAttack", newWhenNormalAttack)
	excel.AddEntryNewFunc("*WhenBehitMoreThan", newWhenBehitMoreThan)
	excel.AddEntryNewFunc("*WhenLockBySkill", newWhenLockBySkill)
	excel.AddEntryNewFunc("*WhenCoolDown", newWhenCoolDown)
	excel.AddEntryNewFunc("*WhenEnterBattleAndCoolDown", newWhenEnterBattleAndCoolDown)
	excel.AddEntryNewFunc("*WhenBeHit", newWhenBeHit)
	excel.AddEntryNewFunc("*QuickSettleTypeBuff", newQuickSettleTypeBuff)
	excel.AddEntryNewFunc("*WhenRoleAddFlag", newWhenRoleAddFlag)
	excel.AddEntryNewFunc("*DeductHp", newDeductHp)
	excel.AddEntryNewFunc("*WhenSupporting", newWhenSupporting)
	excel.AddEntryNewFunc("*WhenSupportingEnd", newWhenSupportingEnd)
	excel.AddEntryNewFunc("*AddHelpValue", newAddHelpValue)
	excel.AddEntryNewFunc("*ChangeCameraViewScale", newChangeCameraViewScale)
	excel.AddEntryNewFunc("*LifeLink", newLifeLink)
	excel.AddEntryNewFunc("*RandomWeightEffect", newRandomWeightEffect)
	excel.AddEntryNewFunc("*RandomWeightEffectUnit", newRandomWeightEffectUnit)
	excel.AddEntryNewFunc("*SkillAttRandom", newSkillAttRandom)
	excel.AddEntryNewFunc("*CameraShake", newCameraShake)
	excel.AddEntryNewFunc("*SkillEffectViewPair", newSkillEffectViewPair)
	excel.AddEntryNewFunc("*WhenBuffAdd", newWhenBuffAdd)
	excel.AddEntryNewFunc("*AddRoleSubDamage", newAddRoleSubDamage)
	excel.AddEntryNewFunc("*ChangeBuffLayerWithParam", newChangeBuffLayerWithParam)
	excel.AddEntryNewFunc("*WhenCallConnection", newWhenCallConnection)
	excel.AddEntryNewFunc("*CallConnection", newCallConnection)
	excel.AddEntryNewFunc("*HitStun", newHitStun)
	excel.AddEntryNewFunc("*FaceMove", newFaceMove)
	excel.AddEntryNewFunc("*HitJudge", newHitJudge)
	excel.AddEntryNewFunc("*Grab", newGrab)
	excel.AddEntryNewFunc("*Release", newRelease)
	excel.AddEntryNewFunc("*ChangeGrabTargetAction", newChangeGrabTargetAction)
	excel.AddEntryNewFunc("*ConnectBulletTime", newConnectBulletTime)
	excel.AddEntryNewFunc("*WhenEffectCountReached", newWhenEffectCountReached)
	excel.AddEntryNewFunc("*WhenDropDown", newWhenDropDown)
	excel.AddEntryNewFunc("*AfterSendHeal", newAfterSendHeal)
	excel.AddEntryNewFunc("*WhenCounter", newWhenCounter)
	excel.AddEntryNewFunc("*CounterAttack", newCounterAttack)
	excel.AddEntryNewFunc("*ComboAttack", newComboAttack)
	excel.AddEntryNewFunc("*WhenCombo", newWhenCombo)
	excel.AddEntryNewFunc("*WhenSelfEnterHitStun", newWhenSelfEnterHitStun)
	excel.AddEntryNewFunc("*BeforeDamage", newBeforeDamage)
	excel.AddEntryNewFunc("*ConnectionTargetSelect", newConnectionTargetSelect)
	excel.AddEntryNewFunc("*WhenAllShieldDispose", newWhenAllShieldDispose)
	excel.AddEntryNewFunc("*WhenAddShield", newWhenAddShield)
	excel.AddEntryNewFunc("*WhenRoleOnBoard", newWhenRoleOnBoard)
	excel.AddEntryNewFunc("*WhenBloodSuck", newWhenBloodSuck)
	excel.AddEntryNewFunc("*WhenComboGroupEffect", newWhenComboGroupEffect)
	excel.AddEntryNewFunc("*ComboFontShowData", newComboFontShowData)
	excel.AddEntryNewFunc("*WhenAttrPerChanged", newWhenAttrPerChanged)
	excel.AddEntryNewFunc("*MoveEndAction", newMoveEndAction)
	excel.AddEntryNewFunc("*WhenSelfAreaChange", newWhenSelfAreaChange)
	excel.AddEntryNewFunc("*HarmImmunity", newHarmImmunity)
	excel.AddEntryNewFunc("*RoleShowFlash", newRoleShowFlash)
	excel.AddEntryNewFunc("*WhenHurtDamageMoreThan", newWhenHurtDamageMoreThan)
	excel.AddEntryNewFunc("*LDChangeGlobalAttr", newLDChangeGlobalAttr)
	excel.AddEntryNewFunc("*WhenWeatherIs", newWhenWeatherIs)
	excel.AddEntryNewFunc("*WhenWeatherNot", newWhenWeatherNot)
	excel.AddEntryNewFunc("*ReqServerEffect", newReqServerEffect)
	excel.AddEntryNewFunc("*ReqRoleEnchant", newReqRoleEnchant)
	excel.AddEntryNewFunc("*WhenRoleEnhant", newWhenRoleEnhant)
	excel.AddEntryNewFunc("*WhenRaceCountChange", newWhenRaceCountChange)
	excel.AddEntryNewFunc("*WhenHeroMove", newWhenHeroMove)
	excel.AddEntryNewFunc("*PathAoe", newPathAoe)
	excel.AddEntryNewFunc("*WhenKillCountChange", newWhenKillCountChange)
	excel.AddEntryNewFunc("*WhenGirdParamChanged", newWhenGirdParamChanged)
	excel.AddEntryNewFunc("*LDChangeBuffLayerWithParam", newLDChangeBuffLayerWithParam)
	excel.AddEntryNewFunc("*AutoEnchant", newAutoEnchant)
	excel.AddEntryNewFunc("*SaleSelf", newSaleSelf)
	excel.AddEntryNewFunc("*WhenHeroCountLessThan", newWhenHeroCountLessThan)
	excel.AddEntryNewFunc("*WhenClickEnchantBtn", newWhenClickEnchantBtn)
	excel.AddEntryNewFunc("*WhenDiceShapeChange", newWhenDiceShapeChange)
	excel.AddEntryNewFunc("*MountBuffWithDiceShape", newMountBuffWithDiceShape)
	excel.AddEntryNewFunc("*SelfRushShock", newSelfRushShock)
	excel.AddEntryNewFunc("*AttackBack", newAttackBack)
	excel.AddEntryNewFunc("*LDWhenAttrChanged", newLDWhenAttrChanged)
	excel.AddEntryNewFunc("*LDChangeAttr", newLDChangeAttr)
	excel.AddEntryNewFunc("*LDChangeBaseAttr", newLDChangeBaseAttr)
	excel.AddEntryNewFunc("*LDChangeGrowAttr", newLDChangeGrowAttr)
	excel.AddEntryNewFunc("*LDForceChangAttr", newLDForceChangAttr)
	excel.AddEntryNewFunc("*LDChangeScale", newLDChangeScale)
	excel.AddEntryNewFunc("*LDTmpFixedAttr", newLDTmpFixedAttr)
	excel.AddEntryNewFunc("*LDListenEnemyAttr", newLDListenEnemyAttr)
	excel.AddEntryNewFunc("*ChangeCameraFocus", newChangeCameraFocus)
	excel.AddEntryNewFunc("*PlayUIEffect", newPlayUIEffect)
	excel.AddEntryNewFunc("*WhenRoleEnter", newWhenRoleEnter)
	excel.AddEntryNewFunc("*AddBuffTypeTime", newAddBuffTypeTime)
	excel.AddEntryNewFunc("*ChangeCounterCDParam", newChangeCounterCDParam)
	excel.AddEntryNewFunc("*ImmunitySkills", newImmunitySkills)
	excel.AddEntryNewFunc("*WhenTypeSkillEnd", newWhenTypeSkillEnd)
	excel.AddEntryNewFunc("*WhenBattleStart", newWhenBattleStart)
	excel.AddEntryNewFunc("*AddEffectParam", newAddEffectParam)
	excel.AddEntryNewFunc("*ChangeToPuppet", newChangeToPuppet)
	excel.AddEntryNewFunc("*AddPuppetShield", newAddPuppetShield)
	excel.AddEntryNewFunc("*PuppetShieldBroken", newPuppetShieldBroken)
	excel.AddEntryNewFunc("*WhenTotalDmgCompare", newWhenTotalDmgCompare)
	excel.AddEntryNewFunc("*WhenTeamMateBehurt", newWhenTeamMateBehurt)
	excel.AddEntryNewFunc("*WhenCounterSuccess", newWhenCounterSuccess)
	excel.AddEntryNewFunc("*WhenPropertyNeedTransfer", newWhenPropertyNeedTransfer)
	excel.AddEntryNewFunc("*TransferProperities", newTransferProperities)
	excel.AddEntryNewFunc("*WhenRelive", newWhenRelive)
	excel.AddEntryNewFunc("*AddReliveCnt", newAddReliveCnt)
	excel.AddEntryNewFunc("*SkillParamPair", newSkillParamPair)
	excel.AddEntryNewFunc("*AddSkillParam", newAddSkillParam)
	excel.AddEntryNewFunc("*SupportSummon", newSupportSummon)
	excel.AddEntryNewFunc("*CampUseSkill", newCampUseSkill)
	excel.AddEntryNewFunc("*CampMemberCntReach", newCampMemberCntReach)
	excel.AddEntryNewFunc("*WhenFirstMove", newWhenFirstMove)
	excel.AddEntryNewFunc("*WhenAssisKllOrKill", newWhenAssisKllOrKill)
	excel.AddEntryNewFunc("*WhenSupportLeave", newWhenSupportLeave)
	excel.AddEntryNewFunc("*WhenCampRoleDead", newWhenCampRoleDead)
	excel.AddEntryNewFunc("*AnyOneAttrReach", newAnyOneAttrReach)
	excel.AddEntryNewFunc("*Boomerang", newBoomerang)
	excel.AddEntryNewFunc("*BoomerangCountEffects", newBoomerangCountEffects)
	excel.AddEntryNewFunc("*LimitDamage", newLimitDamage)
	excel.AddEntryNewFunc("*LimitTotalDamage", newLimitTotalDamage)
	excel.AddEntryNewFunc("*WhenGetDispel", newWhenGetDispel)
	excel.AddEntryNewFunc("*AddBuffTime", newAddBuffTime)

}

func newSkillAction() any {
	return &SkillAction{}
}

func newSkillPerfrom() any {
	return &SkillPerfrom{}
}

func newRougeSkillCondi() any {
	return &RougeSkillCondi{}
}

func newEffectGroup() any {
	return &EffectGroup{}
}

func newSkillSequenceTrigger() any {
	return &SkillSequenceTrigger{}
}

func newSkillExcuteCheck() any {
	return &SkillExcuteCheck{}
}

func newBuffEffect() any {
	return &BuffEffect{}
}

func newAttack() any {
	return &Attack{}
}

func newHeal() any {
	return &Heal{}
}

func newMeleeAttack() any {
	return &MeleeAttack{}
}

func newShoot() any {
	return &Shoot{}
}

func newStrafing() any {
	return &Strafing{}
}

func newProcessBuffEffect() any {
	return &ProcessBuffEffect{}
}

func newSummonAttrExtendsPair() any {
	return &SummonAttrExtendsPair{}
}

func newDOT() any {
	return &DOT{}
}

func newHOT() any {
	return &HOT{}
}

func newAOE() any {
	return &AOE{}
}

func newMove() any {
	return &Move{}
}

func newAvoidMove() any {
	return &AvoidMove{}
}

func newShockWave() any {
	return &ShockWave{}
}

func newDizzy() any {
	return &Dizzy{}
}

func newSnare() any {
	return &Snare{}
}

func newStun() any {
	return &Stun{}
}

func newFreeze() any {
	return &Freeze{}
}

func newFrozen() any {
	return &Frozen{}
}

func newSleep() any {
	return &Sleep{}
}

func newSilence() any {
	return &Silence{}
}

func newCharm() any {
	return &Charm{}
}

func newDisarm() any {
	return &Disarm{}
}

func newTaunt() any {
	return &Taunt{}
}

func newImprison() any {
	return &Imprison{}
}

func newChangeAttr() any {
	return &ChangeAttr{}
}

func newChangeBaseAttr() any {
	return &ChangeBaseAttr{}
}

func newChangeGrowAttr() any {
	return &ChangeGrowAttr{}
}

func newChangeSpecialAttr() any {
	return &ChangeSpecialAttr{}
}

func newForceChangAttr() any {
	return &ForceChangAttr{}
}

func newChangeGlobalAttr() any {
	return &ChangeGlobalAttr{}
}

func newBuffTrigger() any {
	return &BuffTrigger{}
}

func newWhenBeHurt() any {
	return &WhenBeHurt{}
}

func newWhenDeadlyHurt() any {
	return &WhenDeadlyHurt{}
}

func newRandomTrigger() any {
	return &RandomTrigger{}
}

func newAttackAfterEffect() any {
	return &AttackAfterEffect{}
}

func newListenEnemyAttr() any {
	return &ListenEnemyAttr{}
}

func newWhenKill() any {
	return &WhenKill{}
}

func newWillBuildDamage() any {
	return &WillBuildDamage{}
}

func newWhenBeHurtDelay() any {
	return &WhenBeHurtDelay{}
}

func newWhenBackAttack() any {
	return &WhenBackAttack{}
}

func newWhenWillBeHurt() any {
	return &WhenWillBeHurt{}
}

func newWhenShieldBroken() any {
	return &WhenShieldBroken{}
}

func newWhenRefractionBroken() any {
	return &WhenRefractionBroken{}
}

func newWhenAfterBackAttack() any {
	return &WhenAfterBackAttack{}
}

func newWhenUseSkill() any {
	return &WhenUseSkill{}
}

func newWhenWillNormalAttack() any {
	return &WhenWillNormalAttack{}
}

func newWhenSummonDispear() any {
	return &WhenSummonDispear{}
}

func newAfterSendDamage() any {
	return &AfterSendDamage{}
}

func newAfterSendNormalDamage() any {
	return &AfterSendNormalDamage{}
}

func newWhenAttacherDead() any {
	return &WhenAttacherDead{}
}

func newWhenRoleDead() any {
	return &WhenRoleDead{}
}

func newWhenLayerMoreThan() any {
	return &WhenLayerMoreThan{}
}

func newWhenAttrChanged() any {
	return &WhenAttrChanged{}
}

func newWhenCampAttrChanged() any {
	return &WhenCampAttrChanged{}
}

func newSelfCenterAoe() any {
	return &SelfCenterAoe{}
}

func newBuffPair() any {
	return &BuffPair{}
}

func newClash() any {
	return &Clash{}
}

func newMountBuff() any {
	return &MountBuff{}
}

func newSubBuff() any {
	return &SubBuff{}
}

func newCounterInjury() any {
	return &CounterInjury{}
}

func newAddFlag() any {
	return &AddFlag{}
}

func newAddShield() any {
	return &AddShield{}
}

func newRefractionShield() any {
	return &RefractionShield{}
}

func newRelive() any {
	return &Relive{}
}

func newSlash() any {
	return &Slash{}
}

func newSupport() any {
	return &Support{}
}

func newResistDeath() any {
	return &ResistDeath{}
}

func newDispel() any {
	return &Dispel{}
}

func newHpChange() any {
	return &HpChange{}
}

func newHealWhenKill() any {
	return &HealWhenKill{}
}

func newTmpFixedAttr() any {
	return &TmpFixedAttr{}
}

func newStealth() any {
	return &Stealth{}
}

func newAttackByMountBuffLayer() any {
	return &AttackByMountBuffLayer{}
}

func newSummonBase() any {
	return &SummonBase{}
}

func newSummonPitfall() any {
	return &SummonPitfall{}
}

func newSelfBoom() any {
	return &SelfBoom{}
}

func newConditionEffect() any {
	return &ConditionEffect{}
}

func newEffectCondition() any {
	return &EffectCondition{}
}

func newConditionEffectPair() any {
	return &ConditionEffectPair{}
}

func newEffectPair() any {
	return &EffectPair{}
}

func newMountHalo() any {
	return &MountHalo{}
}

func newLinkMagic() any {
	return &LinkMagic{}
}

func newAddAttackRange() any {
	return &AddAttackRange{}
}

func newSkillSiftTargetConfig() any {
	return &SkillSiftTargetConfig{}
}

func newRemoveBuff() any {
	return &RemoveBuff{}
}

func newRemoveHalo() any {
	return &RemoveHalo{}
}

func newVector() any {
	return &Vector{}
}

func newExtraAttr() any {
	return &ExtraAttr{}
}

func newKeyTimes() any {
	return &KeyTimes{}
}

func newChangeExpRate() any {
	return &ChangeExpRate{}
}

func newChangeStage() any {
	return &ChangeStage{}
}

func newTalkEffects() any {
	return &TalkEffects{}
}

func newEnterGatherState() any {
	return &EnterGatherState{}
}

func newChangeShape() any {
	return &ChangeShape{}
}

func newForceOrderSummon() any {
	return &ForceOrderSummon{}
}

func newChangeScale() any {
	return &ChangeScale{}
}

func newChangeNormalAttack() any {
	return &ChangeNormalAttack{}
}

func newStealAttr() any {
	return &StealAttr{}
}

func newExchangeHpPrecent() any {
	return &ExchangeHpPrecent{}
}

func newSkillEnterCd() any {
	return &SkillEnterCd{}
}

func newComboSkill() any {
	return &ComboSkill{}
}

func newSubShield() any {
	return &SubShield{}
}

func newWhenEnterLevel() any {
	return &WhenEnterLevel{}
}

func newWhenEnterBattle() any {
	return &WhenEnterBattle{}
}

func newWhenDamageMoreThan() any {
	return &WhenDamageMoreThan{}
}

func newWhenEnterStage() any {
	return &WhenEnterStage{}
}

func newWhenStageTimeMoreThan() any {
	return &WhenStageTimeMoreThan{}
}

func newWhenDead() any {
	return &WhenDead{}
}

func newWhenBattleTimeMoreThan() any {
	return &WhenBattleTimeMoreThan{}
}

func newWhenNormalAttack() any {
	return &WhenNormalAttack{}
}

func newWhenBehitMoreThan() any {
	return &WhenBehitMoreThan{}
}

func newWhenLockBySkill() any {
	return &WhenLockBySkill{}
}

func newWhenCoolDown() any {
	return &WhenCoolDown{}
}

func newWhenEnterBattleAndCoolDown() any {
	return &WhenEnterBattleAndCoolDown{}
}

func newWhenBeHit() any {
	return &WhenBeHit{}
}

func newQuickSettleTypeBuff() any {
	return &QuickSettleTypeBuff{}
}

func newWhenRoleAddFlag() any {
	return &WhenRoleAddFlag{}
}

func newDeductHp() any {
	return &DeductHp{}
}

func newWhenSupporting() any {
	return &WhenSupporting{}
}

func newWhenSupportingEnd() any {
	return &WhenSupportingEnd{}
}

func newAddHelpValue() any {
	return &AddHelpValue{}
}

func newChangeCameraViewScale() any {
	return &ChangeCameraViewScale{}
}

func newLifeLink() any {
	return &LifeLink{}
}

func newRandomWeightEffect() any {
	return &RandomWeightEffect{}
}

func newRandomWeightEffectUnit() any {
	return &RandomWeightEffectUnit{}
}

func newSkillAttRandom() any {
	return &SkillAttRandom{}
}

func newCameraShake() any {
	return &CameraShake{}
}

func newSkillEffectViewPair() any {
	return &SkillEffectViewPair{}
}

func newWhenBuffAdd() any {
	return &WhenBuffAdd{}
}

func newAddRoleSubDamage() any {
	return &AddRoleSubDamage{}
}

func newChangeBuffLayerWithParam() any {
	return &ChangeBuffLayerWithParam{}
}

func newWhenCallConnection() any {
	return &WhenCallConnection{}
}

func newCallConnection() any {
	return &CallConnection{}
}

func newHitStun() any {
	return &HitStun{}
}

func newFaceMove() any {
	return &FaceMove{}
}

func newHitJudge() any {
	return &HitJudge{}
}

func newGrab() any {
	return &Grab{}
}

func newRelease() any {
	return &Release{}
}

func newChangeGrabTargetAction() any {
	return &ChangeGrabTargetAction{}
}

func newConnectBulletTime() any {
	return &ConnectBulletTime{}
}

func newWhenEffectCountReached() any {
	return &WhenEffectCountReached{}
}

func newWhenDropDown() any {
	return &WhenDropDown{}
}

func newAfterSendHeal() any {
	return &AfterSendHeal{}
}

func newWhenCounter() any {
	return &WhenCounter{}
}

func newCounterAttack() any {
	return &CounterAttack{}
}

func newComboAttack() any {
	return &ComboAttack{}
}

func newWhenCombo() any {
	return &WhenCombo{}
}

func newWhenSelfEnterHitStun() any {
	return &WhenSelfEnterHitStun{}
}

func newBeforeDamage() any {
	return &BeforeDamage{}
}

func newConnectionTargetSelect() any {
	return &ConnectionTargetSelect{}
}

func newWhenAllShieldDispose() any {
	return &WhenAllShieldDispose{}
}

func newWhenAddShield() any {
	return &WhenAddShield{}
}

func newWhenRoleOnBoard() any {
	return &WhenRoleOnBoard{}
}

func newWhenBloodSuck() any {
	return &WhenBloodSuck{}
}

func newWhenComboGroupEffect() any {
	return &WhenComboGroupEffect{}
}

func newComboFontShowData() any {
	return &ComboFontShowData{}
}

func newWhenAttrPerChanged() any {
	return &WhenAttrPerChanged{}
}

func newMoveEndAction() any {
	return &MoveEndAction{}
}

func newWhenSelfAreaChange() any {
	return &WhenSelfAreaChange{}
}

func newHarmImmunity() any {
	return &HarmImmunity{}
}

func newRoleShowFlash() any {
	return &RoleShowFlash{}
}

func newWhenHurtDamageMoreThan() any {
	return &WhenHurtDamageMoreThan{}
}

func newLDChangeGlobalAttr() any {
	return &LDChangeGlobalAttr{}
}

func newWhenWeatherIs() any {
	return &WhenWeatherIs{}
}

func newWhenWeatherNot() any {
	return &WhenWeatherNot{}
}

func newReqServerEffect() any {
	return &ReqServerEffect{}
}

func newReqRoleEnchant() any {
	return &ReqRoleEnchant{}
}

func newWhenRoleEnhant() any {
	return &WhenRoleEnhant{}
}

func newWhenRaceCountChange() any {
	return &WhenRaceCountChange{}
}

func newWhenHeroMove() any {
	return &WhenHeroMove{}
}

func newPathAoe() any {
	return &PathAoe{}
}

func newWhenKillCountChange() any {
	return &WhenKillCountChange{}
}

func newWhenGirdParamChanged() any {
	return &WhenGirdParamChanged{}
}

func newLDChangeBuffLayerWithParam() any {
	return &LDChangeBuffLayerWithParam{}
}

func newAutoEnchant() any {
	return &AutoEnchant{}
}

func newSaleSelf() any {
	return &SaleSelf{}
}

func newWhenHeroCountLessThan() any {
	return &WhenHeroCountLessThan{}
}

func newWhenClickEnchantBtn() any {
	return &WhenClickEnchantBtn{}
}

func newWhenDiceShapeChange() any {
	return &WhenDiceShapeChange{}
}

func newMountBuffWithDiceShape() any {
	return &MountBuffWithDiceShape{}
}

func newSelfRushShock() any {
	return &SelfRushShock{}
}

func newAttackBack() any {
	return &AttackBack{}
}

func newLDWhenAttrChanged() any {
	return &LDWhenAttrChanged{}
}

func newLDChangeAttr() any {
	return &LDChangeAttr{}
}

func newLDChangeBaseAttr() any {
	return &LDChangeBaseAttr{}
}

func newLDChangeGrowAttr() any {
	return &LDChangeGrowAttr{}
}

func newLDForceChangAttr() any {
	return &LDForceChangAttr{}
}

func newLDChangeScale() any {
	return &LDChangeScale{}
}

func newLDTmpFixedAttr() any {
	return &LDTmpFixedAttr{}
}

func newLDListenEnemyAttr() any {
	return &LDListenEnemyAttr{}
}

func newChangeCameraFocus() any {
	return &ChangeCameraFocus{}
}

func newPlayUIEffect() any {
	return &PlayUIEffect{}
}

func newWhenRoleEnter() any {
	return &WhenRoleEnter{}
}

func newAddBuffTypeTime() any {
	return &AddBuffTypeTime{}
}

func newChangeCounterCDParam() any {
	return &ChangeCounterCDParam{}
}

func newImmunitySkills() any {
	return &ImmunitySkills{}
}

func newWhenTypeSkillEnd() any {
	return &WhenTypeSkillEnd{}
}

func newWhenBattleStart() any {
	return &WhenBattleStart{}
}

func newAddEffectParam() any {
	return &AddEffectParam{}
}

func newChangeToPuppet() any {
	return &ChangeToPuppet{}
}

func newAddPuppetShield() any {
	return &AddPuppetShield{}
}

func newPuppetShieldBroken() any {
	return &PuppetShieldBroken{}
}

func newWhenTotalDmgCompare() any {
	return &WhenTotalDmgCompare{}
}

func newWhenTeamMateBehurt() any {
	return &WhenTeamMateBehurt{}
}

func newWhenCounterSuccess() any {
	return &WhenCounterSuccess{}
}

func newWhenPropertyNeedTransfer() any {
	return &WhenPropertyNeedTransfer{}
}

func newTransferProperities() any {
	return &TransferProperities{}
}

func newWhenRelive() any {
	return &WhenRelive{}
}

func newAddReliveCnt() any {
	return &AddReliveCnt{}
}

func newSkillParamPair() any {
	return &SkillParamPair{}
}

func newAddSkillParam() any {
	return &AddSkillParam{}
}

func newSupportSummon() any {
	return &SupportSummon{}
}

func newCampUseSkill() any {
	return &CampUseSkill{}
}

func newCampMemberCntReach() any {
	return &CampMemberCntReach{}
}

func newWhenFirstMove() any {
	return &WhenFirstMove{}
}

func newWhenAssisKllOrKill() any {
	return &WhenAssisKllOrKill{}
}

func newWhenSupportLeave() any {
	return &WhenSupportLeave{}
}

func newWhenCampRoleDead() any {
	return &WhenCampRoleDead{}
}

func newAnyOneAttrReach() any {
	return &AnyOneAttrReach{}
}

func newBoomerang() any {
	return &Boomerang{}
}

func newBoomerangCountEffects() any {
	return &BoomerangCountEffects{}
}

func newLimitDamage() any {
	return &LimitDamage{}
}

func newLimitTotalDamage() any {
	return &LimitTotalDamage{}
}

func newWhenGetDispel() any {
	return &WhenGetDispel{}
}

func newAddBuffTime() any {
	return &AddBuffTime{}
}
