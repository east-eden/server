package config

import (
	"github.com/east-eden/server/excel"
)

// LD技能表-拳皇/LDSkillPoly
type SkillDescParam struct {
	Type SkillDescAttrType `json:"type,omitempty" bson:"type"` //
	Val  float32           `json:"val,omitempty" bson:"val"`   //
}

func init() {
	excel.AddEntryNewFunc("*SkillDescParam", newSkillDescParam)

}

func newSkillDescParam() any {
	return &SkillDescParam{}
}
