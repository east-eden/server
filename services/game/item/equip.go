package item

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
