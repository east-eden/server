package config

import (
	"github.com/east-eden/server/excel"
)

// 皇城争锋-拳皇/KingWarMapBuilding
type KingWarBuilding struct {
	LandIcon         string `json:"landIcon,omitempty" bson:"landIcon"`                 //
	BuildingIcon     string `json:"buildingIcon,omitempty" bson:"buildingIcon"`         //
	FreeIcon         string `json:"freeIcon,omitempty" bson:"freeIcon"`                 //
	FreeBuildingIcon string `json:"freeBuildingIcon,omitempty" bson:"freeBuildingIcon"` //
}

// 皇城争锋-拳皇/KingWarMapBuilding
type KingWarTime struct {
	Week   int32 `json:"week,omitempty" bson:"week"`     //
	Hour   int32 `json:"hour,omitempty" bson:"hour"`     //
	Minute int32 `json:"minute,omitempty" bson:"minute"` //
	Sec    int32 `json:"sec,omitempty" bson:"sec"`       //
}

func init() {
	excel.AddEntryNewFunc("*KingWarBuilding", newKingWarBuilding)
	excel.AddEntryNewFunc("*KingWarTime", newKingWarTime)

}

func newKingWarBuilding() any {
	return &KingWarBuilding{}
}

func newKingWarTime() any {
	return &KingWarTime{}
}
