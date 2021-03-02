package crystal

import (
	"fmt"

	"bitbucket.org/funplus/server/define"
)

type CrystalBox struct {
	owner       define.PluginObj
	crystalList [define.Crystal_PosEnd]*Crystal
}

func NewCrystalBox(owner define.PluginObj) *CrystalBox {
	m := &CrystalBox{
		owner: owner,
	}

	return m
}

func (cb *CrystalBox) GetCrystalByPos(pos int32) *Crystal {
	if pos < define.Crystal_PosBegin || pos >= define.Crystal_PosEnd {
		return nil
	}

	return cb.crystalList[pos]
}

func (cb *CrystalBox) PutonCrystal(c *Crystal) error {
	pos := c.GetOptions().Entry.Pos
	if pos < define.Crystal_PosBegin || pos >= define.Crystal_PosEnd {
		return fmt.Errorf("puton crystal error: invalid pos<%d>", pos)
	}

	if cb.GetCrystalByPos(pos) != nil {
		return fmt.Errorf("puton crystal error: cannot recover crystal on this pos<%d>", pos)
	}

	cb.crystalList[pos] = c
	c.GetOptions().EquipObj = cb.owner.GetID()
	return nil
}

func (cb *CrystalBox) TakeoffCrystal(pos int32) error {
	if pos < define.Crystal_PosBegin || pos >= define.Crystal_PosEnd {
		return fmt.Errorf("takeoff crystal error: invalid pos<%d>", pos)
	}

	if c := cb.crystalList[pos]; c != nil {
		c.GetOptions().EquipObj = -1
	}

	cb.crystalList[pos] = nil
	return nil
}
