package config

import (
	"github.com/east-eden/server/excel"
)

// 关卡表-拳皇/LevelPoly
type MonsterRandomGroup struct {
	MonsterId int32 `json:"monsterId,omitempty" bson:"monsterId"` //
	Weight    int32 `json:"weight,omitempty" bson:"weight"`       //
}

// 关卡表-拳皇/LevelPoly
type MonsterRandomGroupList struct {
	Group []*MonsterRandomGroup `json:"group,omitempty" bson:"group"` //
}

// 关卡表-拳皇/LevelPoly
type RandomMonsterGroup struct {
	Id     int32 `json:"id,omitempty" bson:"id"`         //
	Weight int32 `json:"weight,omitempty" bson:"weight"` //
}

// 关卡表-拳皇/LevelPoly
type TalkGroups struct {
	Ids []int32 `json:"ids,omitempty" bson:"ids"` //
}

// 关卡表-拳皇/LevelPoly
type NpcData struct {
	MonsterId         int32 `json:"monsterId,omitempty" bson:"monsterId"`                 //
	Level             int32 `json:"level,omitempty" bson:"level"`                         //
	Strength          int32 `json:"strength,omitempty" bson:"strength"`                   //
	StrengthType      int32 `json:"strengthType,omitempty" bson:"strengthType"`           //
	MonsterAPStrength int32 `json:"monsterAPStrength,omitempty" bson:"monsterAPStrength"` //
	StrengthAPType    int32 `json:"strengthAPType,omitempty" bson:"strengthAPType"`       //
}

// 关卡表-拳皇/LevelPoly
type MonsterGroup struct {
	MonsterIds []int32 `json:"monsterIds,omitempty" bson:"monsterIds"` //
}

func init() {
	excel.AddEntryNewFunc("*MonsterRandomGroup", newMonsterRandomGroup)
	excel.AddEntryNewFunc("*MonsterRandomGroupList", newMonsterRandomGroupList)
	excel.AddEntryNewFunc("*RandomMonsterGroup", newRandomMonsterGroup)
	excel.AddEntryNewFunc("*TalkGroups", newTalkGroups)
	excel.AddEntryNewFunc("*NpcData", newNpcData)
	excel.AddEntryNewFunc("*MonsterGroup", newMonsterGroup)

}

func newMonsterRandomGroup() any {
	return &MonsterRandomGroup{}
}

func newMonsterRandomGroupList() any {
	return &MonsterRandomGroupList{}
}

func newRandomMonsterGroup() any {
	return &RandomMonsterGroup{}
}

func newTalkGroups() any {
	return &TalkGroups{}
}

func newNpcData() any {
	return &NpcData{}
}

func newMonsterGroup() any {
	return &MonsterGroup{}
}
