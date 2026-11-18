package config

import (
	"github.com/east-eden/server/excel"
)

// 格斗家礼物-拳皇/BoxGiftPoly
type BoxGiftShow struct {
	Spine      string  `json:"spine,omitempty" bson:"spine"`           //
	Action     string  `json:"action,omitempty" bson:"action"`         //
	ActionTime float32 `json:"actionTime,omitempty" bson:"actionTime"` //
}

func init() {
	excel.AddEntryNewFunc("*BoxGiftShow", newBoxGiftShow)

}

func newBoxGiftShow() any {
	return &BoxGiftShow{}
}
