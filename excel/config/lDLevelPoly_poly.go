package config

import (
	"github.com/east-eden/server/excel"
)

// LD关卡表-拳皇/LDLevelPoly
type RoomGrid struct {
	HeroList []int32 `json:"heroList,omitempty" bson:"heroList"` //
}

// LD关卡表-拳皇/LDLevelPoly
type HeroAntiValue struct {
	Cnt   int32   `json:"cnt,omitempty" bson:"cnt"`     //
	Value float32 `json:"value,omitempty" bson:"value"` //
}

// LD关卡表-拳皇/LDLevelPoly
type HeroAnti struct {
	HeroId   int32            `json:"heroId,omitempty" bson:"heroId"`     //
	AttrType LDAttrType       `json:"attrType,omitempty" bson:"attrType"` //
	Values   []*HeroAntiValue `json:"values,omitempty" bson:"values"`     //
}

func init() {
	excel.AddEntryNewFunc("*RoomGrid", newRoomGrid)
	excel.AddEntryNewFunc("*HeroAntiValue", newHeroAntiValue)
	excel.AddEntryNewFunc("*HeroAnti", newHeroAnti)

}

func newRoomGrid() any {
	return &RoomGrid{}
}

func newHeroAntiValue() any {
	return &HeroAntiValue{}
}

func newHeroAnti() any {
	return &HeroAnti{}
}
