package config

import (
	"github.com/east-eden/server/excel"
)

// LD创意工坊-拳皇 client_skip/WorkShopPoly
type PathGroup struct {
	PathId int32 `json:"pathId,omitempty" bson:"pathId"` //
	Weight int32 `json:"weight,omitempty" bson:"weight"` //
}

// LD创意工坊-拳皇 client_skip/WorkShopPoly
type Pos struct {
	X int32 `json:"x,omitempty" bson:"x"` //
	Y int32 `json:"y,omitempty" bson:"y"` //
}

// LD创意工坊-拳皇 client_skip/WorkShopPoly
type WorkShopMonsterGroup struct {
	MonsterId int32 `json:"monsterId,omitempty" bson:"monsterId"` //
	Weight    int32 `json:"weight,omitempty" bson:"weight"`       //
}

// LD创意工坊-拳皇 client_skip/WorkShopPoly
type WorkShopReward struct {
	ItemId int32 `json:"itemId,omitempty" bson:"itemId"` //
	Count  int32 `json:"count,omitempty" bson:"count"`   //
}

// LD创意工坊-拳皇 client_skip/WorkShopPoly
type WsMusic struct {
	Name string `json:"name,omitempty" bson:"name"` //
	Res  string `json:"res,omitempty" bson:"res"`   //
}

func init() {
	excel.AddEntryNewFunc("*PathGroup", newPathGroup)
	excel.AddEntryNewFunc("*Pos", newPos)
	excel.AddEntryNewFunc("*WorkShopMonsterGroup", newWorkShopMonsterGroup)
	excel.AddEntryNewFunc("*WorkShopReward", newWorkShopReward)
	excel.AddEntryNewFunc("*WsMusic", newWsMusic)

}

func newPathGroup() any {
	return &PathGroup{}
}

func newPos() any {
	return &Pos{}
}

func newWorkShopMonsterGroup() any {
	return &WorkShopMonsterGroup{}
}

func newWorkShopReward() any {
	return &WorkShopReward{}
}

func newWsMusic() any {
	return &WsMusic{}
}
