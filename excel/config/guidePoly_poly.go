package config

import (
	"github.com/east-eden/server/excel"
)

// 引导表-拳皇/guidePoly
type GuideTipsParam struct {
	View   GuideViewKey  `json:"view,omitempty" bson:"view"`     //
	Pos    GuideFramePos `json:"pos,omitempty" bson:"pos"`       //
	Offset float32       `json:"offset,omitempty" bson:"offset"` //
}

// 引导表-拳皇/guidePoly
type GuideHandParam struct {
	Action GuideHandAction `json:"action,omitempty" bson:"action"` //
	Views  []GuideViewKey  `json:"views,omitempty" bson:"views"`   //
}

// 引导表-拳皇/guidePoly
type GuideActionParam struct {
	Action GuideCustomAction `json:"action,omitempty" bson:"action"` //
	Datas  []string          `json:"datas,omitempty" bson:"datas"`   //
}

// 引导表-拳皇/guidePoly
type GuideCondition struct {
	Type   GuideCondiType `json:"type,omitempty" bson:"type"`     //
	Params []string       `json:"params,omitempty" bson:"params"` //
}

// 引导表-拳皇/guidePoly
type GuideTalkParam struct {
	Direction GuideTalkDirection `json:"direction,omitempty" bson:"direction"` //
	HeroId    int32              `json:"heroId,omitempty" bson:"heroId"`       //
	Scale     float32            `json:"scale,omitempty" bson:"scale"`         //
	X         float32            `json:"x,omitempty" bson:"x"`                 //
	Y         float32            `json:"y,omitempty" bson:"y"`                 //
}

func init() {
	excel.AddEntryNewFunc("*GuideTipsParam", newGuideTipsParam)
	excel.AddEntryNewFunc("*GuideHandParam", newGuideHandParam)
	excel.AddEntryNewFunc("*GuideActionParam", newGuideActionParam)
	excel.AddEntryNewFunc("*GuideCondition", newGuideCondition)
	excel.AddEntryNewFunc("*GuideTalkParam", newGuideTalkParam)

}

func newGuideTipsParam() any {
	return &GuideTipsParam{}
}

func newGuideHandParam() any {
	return &GuideHandParam{}
}

func newGuideActionParam() any {
	return &GuideActionParam{}
}

func newGuideCondition() any {
	return &GuideCondition{}
}

func newGuideTalkParam() any {
	return &GuideTalkParam{}
}
