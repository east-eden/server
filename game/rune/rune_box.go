package rune

import (
	"fmt"
	"sync"

	"github.com/yokaiio/yokai_server/internal/define"
)

type RuneBox struct {
	owner    define.PluginObj
	runeList [define.Rune_PositionEnd]*Rune

	sync.RWMutex
}

func NewRuneBox(owner define.PluginObj) *RuneBox {
	m := &RuneBox{
		owner: owner,
	}

	return m
}

func (rb *RuneBox) GetRuneByPos(pos int32) *Rune {
	if pos < define.Rune_PositionBegin || pos >= define.Rune_PositionEnd {
		return nil
	}

	return rb.runeList[pos]
}

func (rb *RuneBox) PutonRune(r *Rune, pos int32) error {
	if pos < define.Rune_PositionBegin || pos >= define.Rune_PositionEnd {
		return fmt.Errorf("puton rune error: invalid pos<%d>", pos)
	}

	if rb.GetRuneByPos(pos) != nil {
		return fmt.Errorf("puton rune error: cannot recover rune on this pos<%d>", pos)
	}

	rb.runeList[pos] = r
	r.SetEquipObj(rb.owner.GetID())
	return nil
}

func (rb *RuneBox) TakeoffRune(pos int32) error {
	if pos < define.Rune_PositionBegin || pos >= define.Rune_PositionEnd {
		return fmt.Errorf("takeoff rune error: invalid pos<%d>", pos)
	}

	if r := rb.runeList[pos]; r != nil {
		r.SetEquipObj(-1)
	}

	rb.runeList[pos] = nil
	return nil
}
