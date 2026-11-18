package config

import (
	"github.com/east-eden/server/excel"
)

// 副玩法关卡表-拳皇/TDPoly
type TDPathGroup struct {
	Ids []int32 `json:"ids,omitempty" bson:"ids"` //
}

// 副玩法关卡表-拳皇/TDPoly
type ExtraRandomRogueGroups struct {
	List []*ExtraRandomRogueGroup `json:"list,omitempty" bson:"list"` //
}

// 副玩法关卡表-拳皇/TDPoly
type ExtraRandomRogueGroup struct {
	Id     int32   `json:"id,omitempty" bson:"id"`         //
	Weight float32 `json:"weight,omitempty" bson:"weight"` //
}

func init() {
	excel.AddEntryNewFunc("*TDPathGroup", newTDPathGroup)
	excel.AddEntryNewFunc("*ExtraRandomRogueGroups", newExtraRandomRogueGroups)
	excel.AddEntryNewFunc("*ExtraRandomRogueGroup", newExtraRandomRogueGroup)

}

func newTDPathGroup() any {
	return &TDPathGroup{}
}

func newExtraRandomRogueGroups() any {
	return &ExtraRandomRogueGroups{}
}

func newExtraRandomRogueGroup() any {
	return &ExtraRandomRogueGroup{}
}
