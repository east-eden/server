package rune

import (
	"sync"
	"time"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/internal/att"
)

// rune create pool
var runePool = &sync.Pool{New: newPoolRune}

func NewPoolRune() *Rune {
	return runePool.Get().(*Rune)
}

func GetRunePool() *sync.Pool {
	return runePool
}

func NewRune(opts ...Option) *Rune {
	r := NewPoolRune()

	for _, o := range opts {
		o(r.GetOptions())
	}

	return r
}

type RuneAtt struct {
	AttType  int32 `bson:"att_type" json:"att_type"`
	AttValue int32 `bson:"att_value" json:"att_value"`
}

type Rune struct {
	Options    `bson:"inline" json:",inline"`
	Atts       [define.Rune_AttNum]*RuneAtt `bson:"atts" json:"atts"`
	attManager *att.AttManager              `bson:"-" json:"-"`
}

func newPoolRune() interface{} {
	r := &Rune{
		Options: DefaultOptions(),
	}

	r.attManager = att.NewAttManager()

	return r
}

func (r *Rune) GetExpire() *time.Timer {
	return nil
}

func (r *Rune) GetStoreIndex() int64 {
	return r.Options.OwnerId
}

func (r *Rune) GetOptions() *Options {
	return &r.Options
}

func (r *Rune) GetType() int32 {
	return define.Plugin_Rune
}

func (r *Rune) GetID() int64 {
	return r.Options.Id
}

func (r *Rune) GetOwnerID() int64 {
	return r.Options.OwnerId
}

func (r *Rune) GetTypeID() int32 {
	return r.Options.TypeId
}

func (r *Rune) GetEquipObj() int64 {
	return r.Options.EquipObj
}

func (r *Rune) GetAttManager() *att.AttManager {
	return r.attManager
}

func (r *Rune) GetAtt(idx int32) *RuneAtt {
	if idx < 0 || idx >= define.Rune_AttNum {
		return nil
	}

	return r.Atts[idx]
}

func (r *Rune) SetAtt(idx int32, att *RuneAtt) {
	if idx < 0 || idx >= define.Rune_AttNum {
		return
	}

	r.Atts[idx] = att
}

func (r *Rune) CalcAtt() {
	r.attManager.Reset()

	var n int32
	for n = 0; n < define.Rune_AttNum; n++ {
		att := r.Atts[n]
		if att == nil {
			continue
		}

		r.attManager.ModBaseAtt(int(att.AttType), att.AttValue)
	}

	r.attManager.CalcAtt()
}
