package rune

import (
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
