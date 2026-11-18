package config

import (
	"github.com/east-eden/server/excel"
)

// 音巢挑战-拳皇/nestPoly
type BuffTagEffect struct {
	HeroMark HeroMarkType `json:"heroMark,omitempty" bson:"heroMark"` //
	Tag      int32        `json:"tag,omitempty" bson:"tag"`           //
	Count    int32        `json:"count,omitempty" bson:"count"`       //
}

func init() {
	excel.AddEntryNewFunc("*BuffTagEffect", newBuffTagEffect)

}

func newBuffTagEffect() any {
	return &BuffTagEffect{}
}
