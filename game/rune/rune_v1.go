package rune

import (
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
)

type RuneAtt struct {
	AttType  int32 `bson:"att_type" redis:"att_type"`
	AttValue int64 `bson:"att_value" redis:"att_value"`
}

type RuneV1 struct {
	Opts       *Options                     `bson:"inline" redis:"inline"`
	atts       [define.Rune_AttNum]*RuneAtt `bson:"atts" redis:"atts"`
	attManager *att.AttManager              `bson:"-" redis:"-"`
}

func newPoolRuneV1() interface{} {
	r := &RuneV1{
		Opts: DefaultOptions(),
	}

	r.attManager = att.NewAttManager(-1)

	return r
}

func (r *RuneV1) Options() *Options {
	return r.Opts
}

func (r *RuneV1) GetType() int32 {
	return define.Plugin_Rune
}

func (r *RuneV1) GetID() int64 {
	return r.Opts.Id
}

func (r *RuneV1) GetOwnerID() int64 {
	return r.Opts.OwnerId
}

func (r *RuneV1) GetTypeID() int32 {
	return r.Opts.TypeId
}

func (r *RuneV1) GetEquipObj() int64 {
	return r.Opts.EquipObj
}

func (r *RuneV1) GetAttManager() *att.AttManager {
	return r.attManager
}

func (r *RuneV1) GetAtt(idx int32) *RuneAtt {
	if idx < 0 || idx >= define.Rune_AttNum {
		return nil
	}

	return r.atts[idx]
}

func (r *RuneV1) SetAtt(idx int32, att *RuneAtt) {
	if idx < 0 || idx >= define.Rune_AttNum {
		return
	}

	r.atts[idx] = att
}

func (r *RuneV1) CalcAtt() {
	r.attManager.Reset()

	var n int32
	for n = 0; n < define.Rune_AttNum; n++ {
		att := r.atts[n]
		if att == nil {
			continue
		}

		r.attManager.ModBaseAtt(att.AttType, att.AttValue)
	}

	r.attManager.CalcAtt()
}
