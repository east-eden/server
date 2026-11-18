package config

import (
	"github.com/east-eden/server/excel"
)

// LD反作弊配置表-拳皇/AntiPoly
type GuanYuInhibitFactor struct {
	LevelId int32   `json:"levelId,omitempty" bson:"levelId"` //
	Num     int32   `json:"num,omitempty" bson:"num"`         //
	Factor  float32 `json:"factor,omitempty" bson:"factor"`   //
}

func init() {
	excel.AddEntryNewFunc("*GuanYuInhibitFactor", newGuanYuInhibitFactor)

}

func newGuanYuInhibitFactor() any {
	return &GuanYuInhibitFactor{}
}
