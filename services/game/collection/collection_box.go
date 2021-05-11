package collection

import (
	"fmt"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/utils"
)

// 收集品放置管理
type CollectionBox struct {
	tp             int32
	collectionList map[int32]*Collection
	Entry          *auto.CollectionBoxEntry
}

func NewCollectionBox(tp int32) *CollectionBox {
	m := &CollectionBox{
		tp:             tp,
		collectionList: make(map[int32]*Collection),
	}

	m.Entry, _ = auto.GetCollectionBoxEntry(tp)
	return m
}

func (cb *CollectionBox) PutonCollection(c *Collection) error {
	pos := e.EquipEnchantEntry.EquipPos
	if !utils.Between(int(pos), int(define.Equip_Pos_Begin), int(define.Equip_Pos_End)) {
		return fmt.Errorf("puton equip error: invalid pos<%d>", pos)
	}

	if eb.GetEquipByPos(pos) != nil {
		return fmt.Errorf("puton equip error: existing equip on this pos<%d>", pos)
	}

	eb.equipList[pos] = e
	e.EquipObj = eb.owner.GetId()
	return nil
}

func (cb *CollectionBox) TakeoffCollection(c *Collection) error {
	if !utils.Between(int(pos), int(define.Equip_Pos_Begin), int(define.Equip_Pos_End)) {
		return fmt.Errorf("takeoff equip error: invalid pos<%d>", pos)
	}

	if e := eb.equipList[pos]; e != nil {
		e.EquipObj = -1
	}

	eb.equipList[pos] = nil
	return nil
}
