package scene

import (
	"log"
	"strconv"
	"strings"

	"github.com/willf/bitset"
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/internal/att"
	pbGame "github.com/east-eden/server/proto/game"
)

type UnitOption func(*UnitOptions)
type UnitOptions struct {
	TypeId     int32
	AttValue   []int64
	Position   [3]float32
	Entry      *define.UnitEntry
	AttManager *att.AttManager
	Scene      *Scene
	CombatCtrl *CombatCtrl

	State    *bitset.BitSet
	Immunity [define.ImmunityType_End]*bitset.BitSet
}

func DefaultUnitOptions() *UnitOptions {
	o := &UnitOptions{
		TypeId:     -1,
		Position:   [3]float32{0, 0, 0},
		Entry:      nil,
		AttManager: nil,
		Scene:      nil,
		CombatCtrl: nil,
		State:      bitset.New(uint(define.HeroState_End)),
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

func WithUnitEntry(entry *define.UnitEntry) UnitOption {
	return func(o *UnitOptions) {
		o.Entry = entry
	}
}

func WithUnitScene(scene *Scene) UnitOption {
	return func(o *UnitOptions) {
		o.Scene = scene
	}
}

func WithUnitCombatCtrl(ctrl *CombatCtrl) UnitOption {
	return func(o *UnitOptions) {
		o.CombatCtrl = ctrl
	}
}

func WithUnitAttValue(value []int64) UnitOption {
	return func(o *UnitOptions) {
		o.AttValue = value
	}
}

func WithUnitAttList(attList []*pbGame.Att) UnitOption {
	return func(o *UnitOptions) {
		o.AttValue = make([]int64, define.Att_End)

		for _, v := range attList {
			o.AttValue[v.AttType] = v.AttValue
		}
	}
}

func WithUnitPosition(pos []float32) UnitOption {
	return func(o *UnitOptions) {
		if len(pos) != 3 {
			return
		}

		for k, v := range pos {
			o.Position[k] = v
		}
	}
}

func WithUnitPositionString(pos string) UnitOption {
	return func(o *UnitOptions) {
		arr := strings.Split(pos, ",")
		if len(arr) != 3 {
			return
		}

		for k, v := range arr {
			p, err := strconv.ParseFloat(v, 32)
			if err != nil {
				log.Println("parsing position string error:", err)
				return
			}

			o.Position[k] = float32(p)
		}
	}
}

func WithAttManager(attId int32) UnitOption {
	return func(o *UnitOptions) {
		o.AttManager = att.NewAttManager(attId)
	}
}
