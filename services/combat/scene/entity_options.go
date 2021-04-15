package scene

import (
	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/internal/att"
	"bitbucket.org/funplus/server/utils"
	"github.com/willf/bitset"
)

type EntityOption func(*EntityOptions)
type EntityOptions struct {
	TypeId     int32
	PosX       int32
	PosZ       int32
	AtbValue   int32
	AttManager *att.AttManager
	Scene      *Scene
	ActionCtrl *ActionCtrl
	CombatCtrl *CombatCtrl
	MoveCtrl   *MoveCtrl
	Entry      *auto.HeroEntry

	State    *utils.CountableBitset
	Immunity [define.ImmunityType_End]*bitset.BitSet
}

func DefaultUnitOptions() *EntityOptions {
	o := &EntityOptions{
		TypeId:     -1,
		PosX:       0,
		PosZ:       0,
		AtbValue:   0,
		AttManager: att.NewAttManager(),
		Scene:      nil,
		State:      utils.NewCountableBitset(uint(define.HeroState_End)),
	}

	for k := range o.Immunity {
		o.Immunity[k] = bitset.New(uint(64))
	}

	return o
}

func WithEntityTypeId(typeId int32) EntityOption {
	return func(o *EntityOptions) {
		o.TypeId = typeId
	}
}

func WithEntityAtbValue(value int32) EntityOption {
	return func(o *EntityOptions) {
		o.AtbValue = value
	}
}

func WithEntityHeroEntry(entry *auto.HeroEntry) EntityOption {
	return func(o *EntityOptions) {
		o.Entry = entry
	}
}

func WithEntityScene(scene *Scene) EntityOption {
	return func(o *EntityOptions) {
		o.Scene = scene
	}
}

func WithEntityActionCtrl(ctrl *ActionCtrl) EntityOption {
	return func(o *EntityOptions) {
		o.ActionCtrl = ctrl
	}
}

func WithEntityCombatCtrl(ctrl *CombatCtrl) EntityOption {
	return func(o *EntityOptions) {
		o.CombatCtrl = ctrl
	}
}

func WithEntityMoveCtrl(ctrl *MoveCtrl) EntityOption {
	return func(o *EntityOptions) {
		o.MoveCtrl = ctrl
	}
}

func WithEntityAttValue(tp int, value int32) EntityOption {
	return func(o *EntityOptions) {
		o.AttManager.SetAttValue(tp, value)
	}
}

func WithEntityAttList(attList []int32) EntityOption {
	return func(o *EntityOptions) {
		for tp := range attList {
			o.AttManager.SetAttValue(tp, attList[tp])
		}
	}
}

func WithEntityPosition(posX, posZ int32) EntityOption {
	return func(o *EntityOptions) {
		o.PosX = posX
		o.PosZ = posZ
	}
}
