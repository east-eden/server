package config

import (
	"github.com/east-eden/server/excel"
)

// LD异变表-拳皇/MutationPoly
type MutationBossDynamicFactor struct {
	LeftTime int32   `json:"leftTime,omitempty" bson:"leftTime"` //
	Factor   float32 `json:"factor,omitempty" bson:"factor"`     //
}

func init() {
	excel.AddEntryNewFunc("*MutationBossDynamicFactor", newMutationBossDynamicFactor)

}

func newMutationBossDynamicFactor() any {
	return &MutationBossDynamicFactor{}
}
