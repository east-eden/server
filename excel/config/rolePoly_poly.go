package config

import (
	"github.com/east-eden/server/excel"
)

// 英雄表-拳皇/rolePoly
type HateInfo struct {
	Range     float32 `json:"range,omitempty" bson:"range"`         //
	BaseValue float32 `json:"baseValue,omitempty" bson:"baseValue"` //
}

// 英雄表-拳皇/rolePoly
type RoleSkillUnlock struct {
	Id          int32 `json:"id,omitempty" bson:"id"`                   //
	UnlockLevel int32 `json:"unlockLevel,omitempty" bson:"unlockLevel"` //
}

// 英雄表-拳皇/rolePoly
type ArtifactRaritySection struct {
	Rarity     Rarity `json:"rarity,omitempty" bson:"rarity"`         //
	MaxSection int32  `json:"maxSection,omitempty" bson:"maxSection"` //
	Weight     int32  `json:"weight,omitempty" bson:"weight"`         //
}

func init() {
	excel.AddEntryNewFunc("*HateInfo", newHateInfo)
	excel.AddEntryNewFunc("*RoleSkillUnlock", newRoleSkillUnlock)
	excel.AddEntryNewFunc("*ArtifactRaritySection", newArtifactRaritySection)

}

func newHateInfo() any {
	return &HateInfo{}
}

func newRoleSkillUnlock() any {
	return &RoleSkillUnlock{}
}

func newArtifactRaritySection() any {
	return &ArtifactRaritySection{}
}
