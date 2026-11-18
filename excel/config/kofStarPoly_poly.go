package config

import (
	"github.com/east-eden/server/excel"
)

// 拳皇大明星-拳皇/kofStarPoly
type KofStarRandomPair struct {
	RandomType KofStarRandomType `json:"randomType,omitempty" bson:"randomType"` //
	Weights    float32           `json:"weights,omitempty" bson:"weights"`       //
}

// 拳皇大明星-拳皇/kofStarPoly
type KofStarHeroRandomPair struct {
	HeroId  int32   `json:"heroId,omitempty" bson:"heroId"`   //
	Weights float32 `json:"weights,omitempty" bson:"weights"` //
}

// 拳皇大明星-拳皇/kofStarPoly
type KofStarRandomCondition struct {
	Type  KofRogueCardType `json:"type,omitempty" bson:"type"`   //
	Id    int32            `json:"id,omitempty" bson:"id"`       //
	Level int32            `json:"level,omitempty" bson:"level"` //
}

func init() {
	excel.AddEntryNewFunc("*KofStarRandomPair", newKofStarRandomPair)
	excel.AddEntryNewFunc("*KofStarHeroRandomPair", newKofStarHeroRandomPair)
	excel.AddEntryNewFunc("*KofStarRandomCondition", newKofStarRandomCondition)

}

func newKofStarRandomPair() any {
	return &KofStarRandomPair{}
}

func newKofStarHeroRandomPair() any {
	return &KofStarHeroRandomPair{}
}

func newKofStarRandomCondition() any {
	return &KofStarRandomCondition{}
}
