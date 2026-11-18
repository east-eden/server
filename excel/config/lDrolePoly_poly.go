package config

import (
	"github.com/east-eden/server/excel"
)

// LD英雄表-拳皇/LDrolePoly
type HeroParamCount struct {
	ParamKey string `json:"paramKey,omitempty" bson:"paramKey"` //
	Count    int32  `json:"count,omitempty" bson:"count"`       //
}

// LD英雄表-拳皇/LDrolePoly
type BondsUnitCheckPair struct {
	Pression LDRoleProfessionType `json:"pression,omitempty" bson:"pression"` //
	Count    int32                `json:"count,omitempty" bson:"count"`       //
	Rarity   LDRarity             `json:"rarity,omitempty" bson:"rarity"`     //
}

// LD英雄表-拳皇/LDrolePoly
type BondsShopLevelCheckUnit struct {
	Rarity LDRarity `json:"rarity,omitempty" bson:"rarity"` //
	Level  int32    `json:"level,omitempty" bson:"level"`   //
}

func init() {
	excel.AddEntryNewFunc("*HeroParamCount", newHeroParamCount)
	excel.AddEntryNewFunc("*BondsUnitCheckPair", newBondsUnitCheckPair)
	excel.AddEntryNewFunc("*BondsShopLevelCheckUnit", newBondsShopLevelCheckUnit)

}

func newHeroParamCount() any {
	return &HeroParamCount{}
}

func newBondsUnitCheckPair() any {
	return &BondsUnitCheckPair{}
}

func newBondsShopLevelCheckUnit() any {
	return &BondsShopLevelCheckUnit{}
}
