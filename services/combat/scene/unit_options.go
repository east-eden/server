package scene

import (
	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/internal/att"
	"bitbucket.org/funplus/server/utils"
	"github.com/willf/bitset"
)

type UnitOption func(*UnitOptions)
type UnitOptions struct {
	TypeId     int32
	AttValue   []int32
	PosX       int32
	PosY       int32
	Entry      *auto.UnitEntry
	AttManager *att.AttManager
	Scene      *Scene
	ActionCtrl *ActionCtrl
	CombatCtrl *CombatCtrl
	MoveCtrl   *MoveCtrl

	State    *utils.CountableBitset
	Immunity [define.ImmunityType_End]*bitset.BitSet
}

func DefaultUnitOptions() *UnitOptions {
	o := &UnitOptions{
		TypeId:     -1,
		PosX:       0,
		PosY:       0,
		Entry:      nil,
		AttManager: nil,
		Scene:      nil,
		State:      utils.NewCountableBitset(uint(define.HeroState_End)),
	}

	for k := range o.Immunity {
		o.Immunity[k] = bitset.New(uint(64))
	}

	return o
}

func WithUnitTypeId(typeId int32) UnitOption {
	return func(o *UnitOptions) {
		o.TypeId = typeId
	}
}

func WithUnitEntry(entry *auto.UnitEntry) UnitOption {
	return func(o *UnitOptions) {
		o.Entry = entry
	}
}

func WithUnitScene(scene *Scene) UnitOption {
	return func(o *UnitOptions) {
		o.Scene = scene
	}
}

func WithUnitActionCtrl(ctrl *ActionCtrl) UnitOption {
	return func(o *UnitOptions) {
		o.ActionCtrl = ctrl
	}
}

func WithUnitCombatCtrl(ctrl *CombatCtrl) UnitOption {
	return func(o *UnitOptions) {
		o.CombatCtrl = ctrl
	}
}

func WithUnitMoveCtrl(ctrl *MoveCtrl) UnitOption {
	return func(o *UnitOptions) {
		o.MoveCtrl = ctrl
	}
}

func WithUnitAttValue(value []int32) UnitOption {
	return func(o *UnitOptions) {
		o.AttValue = value
	}
}

func WithUnitAttList(attList []int32) UnitOption {
	return func(o *UnitOptions) {
		o.AttValue = make([]int32, define.Att_End)

		for k, v := range attList {
			o.AttValue[k] = v
		}
	}
}

func WithUnitPosition(posX, posY int32) UnitOption {
	return func(o *UnitOptions) {
		o.PosX = posX
		o.PosY = posY
	}
}

func WithAttManager(attId int32) UnitOption {
	return func(o *UnitOptions) {
		o.AttManager = att.NewAttManager()
		o.AttManager.SetBaseAttId(attId)
	}
}
