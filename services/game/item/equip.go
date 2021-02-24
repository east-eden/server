package item

import (
	"bitbucket.org/east-eden/server/excel/auto"
)

// 装备
type Equip struct {
	*Item        `bson:"inline" json:",inline"`
	EquipOptions `bson:"inline" json:",inline"`
	attManager   *EquipAttManager `json:"-" bson:"-"`
}

func (e *Equip) Init(opts ...EquipOption) {
	for _, o := range opts {
		o(&e.EquipOptions)
	}
}

func (e *Equip) OnDelete() {
	e.SetEquipObj(-1)
	e.Item.OnDelete()
}

func (e *Equip) GetAttManager() *EquipAttManager {
	return e.attManager
}

func (e *Equip) GetEquipEnchantEntry() *auto.EquipEnchantEntry {
	return e.EquipOptions.EquipEnchantEntry
}

func (e *Equip) GetEquipObj() int64 {
	return e.EquipOptions.EquipObj
}

func (e *Equip) SetEquipObj(obj int64) {
	e.EquipOptions.EquipObj = obj
}
