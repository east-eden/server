package config

import (
	"github.com/east-eden/server/excel"
)

// 阵营战-拳皇/CampBattleMapPoly
type CampBattleMapGrid struct {
	MapGridType EMCampBattleMapGridType `json:"mapGridType,omitempty" bson:"mapGridType"` //
}

func init() {
	excel.AddEntryNewFunc("*CampBattleMapGrid", newCampBattleMapGrid)

}

func newCampBattleMapGrid() any {
	return &CampBattleMapGrid{}
}
