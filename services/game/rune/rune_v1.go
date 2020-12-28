package rune

import (
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/att"
)

type RuneAtt struct {
	AttType  int32 `bson:"att_type" json:"att_type"`
	AttValue int64 `bson:"att_value" json:"att_value"`
}

type RuneV1 struct {
	Options    `bson:"inline" json:",inline"`
	Atts       [define.Rune_AttNum]*RuneAtt `bson:"atts" json:"atts"`
	attManager *att.AttManager              `bson:"-" json:"-"`
}

func newPoolRuneV1() interface{} {
	r := &RuneV1{
		Options: DefaultOptions(),
	}

	r.attManager = att.NewAttManager(-1)

	return r
}

// StoreObjector interface
func (r *RuneV1) AfterLoad() error {
	return nil
}

func (r *RuneV1) GetExpire() *time.Timer {
	return nil
}

func (r *RuneV1) GetObjID() int64 {
	return r.Options.Id
}

func (r *RuneV1) GetStoreIndex() int64 {
	return r.Options.OwnerId
}

func (r *RuneV1) GetOptions() *Options {
	return &r.Options
}

func (r *RuneV1) GetType() int32 {
	return define.Plugin_Rune
}

func (r *RuneV1) GetID() int64 {
	return r.Options.Id
}

func (r *RuneV1) GetOwnerID() int64 {
	return r.Options.OwnerId
}

func (r *RuneV1) GetTypeID() int {
	return r.Options.TypeId
}

func (r *RuneV1) GetEquipObj() int64 {
	return r.Options.EquipObj
}

func (r *RuneV1) GetAttManager() *att.AttManager {
	return r.attManager
}

func (r *RuneV1) GetAtt(idx int32) *RuneAtt {
	if idx < 0 || idx >= define.Rune_AttNum {
		return nil
	}

	return r.Atts[idx]
}

func (r *RuneV1) SetAtt(idx int32, att *RuneAtt) {
	if idx < 0 || idx >= define.Rune_AttNum {
		return
	}

	r.Atts[idx] = att
}

func (r *RuneV1) CalcAtt() {
	r.attManager.Reset()

	var n int32
	for n = 0; n < define.Rune_AttNum; n++ {
		att := r.Atts[n]
		if att == nil {
			continue
		}

		r.attManager.ModBaseAtt(att.AttType, att.AttValue)
	}

	r.attManager.CalcAtt()
}
