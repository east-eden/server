package crystal

import (
	"fmt"
	"sync"

	"bitbucket.org/funplus/server/define"
)

type CrystalBox struct {
	owner       define.PluginObj
	crystalList [define.Crystal_PositionEnd]*Crystal

	sync.RWMutex
}

func NewCrystalBox(owner define.PluginObj) *CrystalBox {
	m := &CrystalBox{
		owner: owner,
	}

	return m
}

func (rb *CrystalBox) GetCrystalByPos(pos int32) *Crystal {
	if pos < define.Crystal_PositionBegin || pos >= define.Crystal_PositionEnd {
		return nil
	}

	return rb.crystalList[pos]
}

func (rb *CrystalBox) PutonCrystal(r *Crystal) error {
	pos := r.GetOptions().Entry.Pos
	if pos < define.Crystal_PositionBegin || pos >= define.Crystal_PositionEnd {
		return fmt.Errorf("puton crystal error: invalid pos<%d>", pos)
	}

	if rb.GetCrystalByPos(pos) != nil {
		return fmt.Errorf("puton crystal error: cannot recover crystal on this pos<%d>", pos)
	}

	rb.crystalList[pos] = r
	r.GetOptions().EquipObj = rb.owner.GetID()
	return nil
}

func (rb *CrystalBox) TakeoffCrystal(pos int32) error {
	if pos < define.Crystal_PositionBegin || pos >= define.Crystal_PositionEnd {
		return fmt.Errorf("takeoff crystal error: invalid pos<%d>", pos)
	}

	if r := rb.crystalList[pos]; r != nil {
		r.GetOptions().EquipObj = -1
	}

	rb.crystalList[pos] = nil
	return nil
}
