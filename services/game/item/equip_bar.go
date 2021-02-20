package item

import (
	"fmt"

	"bitbucket.org/east-eden/server/define"
)

type EquipBar struct {
	owner     define.PluginObj
	equipList [define.Equip_Pos_End]*Item
}

func NewEquipBar(owner define.PluginObj) *EquipBar {
	m := &EquipBar{
		owner: owner,
	}

	return m
}

func (eb *EquipBar) GetEquipByPos(pos int32) *Item {
	if pos < 0 || pos >= define.Equip_Pos_End {
		return nil
	}

	return eb.equipList[pos]
}

func (eb *EquipBar) PutonEquip(i *Item) error {
	pos := i.EquipEnchantEntry().EquipPos
	if pos < 0 || pos >= define.Equip_Pos_End {
		return fmt.Errorf("puton equip error: invalid pos<%d>", pos)
	}

	if eb.GetEquipByPos(pos) != nil {
		return fmt.Errorf("puton equip error: existing equip on this pos<%d>", pos)
	}

	eb.equipList[pos] = i
	i.SetEquipObj(eb.owner.GetID())
	return nil
}

func (eb *EquipBar) TakeoffEquip(pos int32) error {
	if pos < 0 || pos >= define.Equip_Pos_End {
		return fmt.Errorf("takeoff equip error: invalid pos<%d>", pos)
	}

	if i := eb.equipList[pos]; i != nil {
		i.SetEquipObj(-1)
	}

	eb.equipList[pos] = nil
	return nil
}
