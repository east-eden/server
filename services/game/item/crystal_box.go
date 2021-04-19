package item

import (
	"fmt"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
)

type CrystalBox struct {
	owner       define.PluginObj
	crystalList [define.Crystal_PosEnd]*Crystal
}

func NewCrystalBox(owner define.PluginObj) *CrystalBox {
	m := &CrystalBox{
		owner: owner,
	}

	for k := range m.crystalList {
		m.crystalList[k] = nil
	}

	return m
}

func (cb *CrystalBox) GetCrystalByPos(pos int32) *Crystal {
	if pos < define.Crystal_PosBegin || pos >= define.Crystal_PosEnd {
		return nil
	}

	return cb.crystalList[pos]
}

func (cb *CrystalBox) GetSkills() []int32 {
	rows := auto.GetCrystalSkillRows()

	elemNum := make(map[int32]int32)
	for _, c := range cb.crystalList {
		if c == nil {
			continue
		}

		elemNum[c.Opts().ItemEntry.SubType]++
	}

	crystalSkills := make([]int32, 0, len(elemNum))
	for tp, num := range elemNum {
		if num <= 0 {
			continue
		}

		skillId := rows[tp].SkillId[num-1]
		crystalSkills = append(crystalSkills, skillId)
	}

	return crystalSkills
}

func (cb *CrystalBox) PutonCrystal(c *Crystal) error {
	pos := c.CrystalEntry.Pos
	if pos < define.Crystal_PosBegin || pos >= define.Crystal_PosEnd {
		return fmt.Errorf("puton crystal error: invalid pos<%d>", pos)
	}

	if cb.GetCrystalByPos(pos) != nil {
		return fmt.Errorf("puton crystal error: cannot recover crystal on this pos<%d>", pos)
	}

	cb.crystalList[pos] = c
	c.CrystalObj = cb.owner.GetID()
	return nil
}

func (cb *CrystalBox) TakeoffCrystal(pos int32) error {
	if pos < define.Crystal_PosBegin || pos >= define.Crystal_PosEnd {
		return fmt.Errorf("takeoff crystal error: invalid pos<%d>", pos)
	}

	if c := cb.crystalList[pos]; c != nil {
		c.CrystalObj = -1
	}

	cb.crystalList[pos] = nil
	return nil
}
