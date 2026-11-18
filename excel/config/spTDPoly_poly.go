package config

import (
	"github.com/east-eden/server/excel"
)

// SpTD关卡表-拳皇/SpTDPoly
type GlobalAttrs struct {
	Attrs []*AttrPair `json:"attrs,omitempty" bson:"attrs"` //
}

func init() {
	excel.AddEntryNewFunc("*GlobalAttrs", newGlobalAttrs)

}

func newGlobalAttrs() any {
	return &GlobalAttrs{}
}
