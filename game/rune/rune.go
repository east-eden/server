package rune

import (
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
)

type RuneAtt struct {
	AttType  int32 `gorm:"-" bson:"att_type"`
	AttValue int64 `gorm:"-" bson:"att_value"`
}

type Rune struct {
	ID       int64 `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerID  int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	TypeID   int32 `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`
	EquipObj int64 `gorm:"type:bigint(20);column:equip_obj;default:-1;not null" bson:"equip_obj"`

	atts       [define.Rune_AttNum]*RuneAtt `gorm:"-" bson:"atts"`
	entry      *define.RuneEntry            `gorm:"-" bson:"-"`
	attManager *att.AttManager              `gorm:"-" bson:"-"`
}

func NewRune(id int64) *Rune {
	return &Rune{
		ID:       id,
		EquipObj: -1,
	}
}

func (r *Rune) GetID() int64 {
	return r.ID
}

func (r *Rune) GetOwnerID() int64 {
	return r.OwnerID
}

func (r *Rune) GetTypeID() int32 {
	return r.TypeID
}

func (r *Rune) GetEquipObj() int64 {
	return r.EquipObj
}

func (r *Rune) GetAttManager() *att.AttManager {
	return r.attManager
}

func (r *Rune) GetAtt(idx int32) *RuneAtt {
	if idx < 0 || idx >= define.Rune_AttNum {
		return nil
	}

	return r.atts[idx]
}

func (r *Rune) Entry() *define.RuneEntry {
	return r.entry
}

func (r *Rune) SetOwnerID(id int64) {
	r.OwnerID = id
}

func (r *Rune) SetTypeID(id int32) {
	r.TypeID = id
}

func (r *Rune) SetEquipObj(id int64) {
	r.EquipObj = id
}

func (r *Rune) SetEntry(e *define.RuneEntry) {
	r.entry = e
}

func (r *Rune) SetAttManager(m *att.AttManager) {
	r.attManager = m
}

func (r *Rune) SetAtt(idx int32, att *RuneAtt) {
	if idx < 0 || idx >= define.Rune_AttNum {
		return
	}

	r.atts[idx] = att
}

func (r *Rune) CalcAtt() {
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
