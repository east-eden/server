package config

import (
	"github.com/east-eden/server/excel"
)

// 超级克拉克-拳皇/SuperClarkPoly
type PrecentRange struct {
	Precent []float32 `json:"precent,omitempty" bson:"precent"` //
	Weight  float32   `json:"weight,omitempty" bson:"weight"`   //
}

func init() {
	excel.AddEntryNewFunc("*PrecentRange", newPrecentRange)

}

func newPrecentRange() any {
	return &PrecentRange{}
}
