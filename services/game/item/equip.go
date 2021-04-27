package item

import (
	pbCommon "bitbucket.org/funplus/server/proto/global/common"
)

// 装备
type Equip struct {
	Item         `bson:"inline" json:",inline"`
	EquipOptions `bson:"inline" json:",inline"`
	attManager   *EquipAttManager `json:"-" bson:"-"`
}

func (e *Equip) InitEquip(opts ...EquipOption) {
	for _, o := range opts {
		o(&e.EquipOptions)
	}
}

func (e *Equip) GetAttManager() *EquipAttManager {
	return e.attManager
}

func (e *Equip) GetEquipObj() int64 {
	return e.EquipObj
}

func (e *Equip) GenEquipDataPB() *pbCommon.EquipData {
	pb := &pbCommon.EquipData{
		Exp:      e.Exp,
		Level:    int32(e.Level),
		Promote:  int32(e.Promote),
		Star:     int32(e.Star),
		Lock:     e.Lock,
		EquipObj: e.EquipObj,
	}

	return pb
}

func (e *Equip) GenEquipPB() *pbCommon.Equip {
	pb := &pbCommon.Equip{
		Item:      e.GenItemPB(),
		EquipData: e.GenEquipDataPB(),
	}

	return pb
}
