package config

import (
	"github.com/east-eden/server/excel"
)

// 领地表-拳皇/ManoPoly
type ManorResShowParam struct {
	GoodsId int32 `json:"goodsId,omitempty" bson:"goodsId"` //
	Min     int32 `json:"min,omitempty" bson:"min"`         //
	Max     int32 `json:"max,omitempty" bson:"max"`         //
}

// 领地表-拳皇/ManoPoly
type ManorResourceCount struct {
	Id    ManorResourceType `json:"id,omitempty" bson:"id"`       //
	Count int32             `json:"count,omitempty" bson:"count"` //
}

func init() {
	excel.AddEntryNewFunc("*ManorResShowParam", newManorResShowParam)
	excel.AddEntryNewFunc("*ManorResourceCount", newManorResourceCount)

}

func newManorResShowParam() any {
	return &ManorResShowParam{}
}

func newManorResourceCount() any {
	return &ManorResourceCount{}
}
