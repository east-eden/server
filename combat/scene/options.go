package scene

import (
	"log"
	"strconv"
	"strings"

	"github.com/yokaiio/yokai_server/define"
	pbCombat "github.com/yokaiio/yokai_server/proto/combat"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
)

type SceneOption func(*SceneOptions)
type SceneOptions struct {
	AttackId        int64
	DefenceId       int64
	AttackUnitList  []*pbCombat.UnitInfo
	DefenceUnitList []*pbCombat.UnitInfo
	Entry           *define.SceneEntry
}

func DefaultSceneOptions() *SceneOptions {
	return &SceneOptions{
		AttackId:  -1,
		DefenceId: -1,
		Entry:     nil,
	}
}

func WithSceneAttackId(id int64) SceneOption {
	return func(o *SceneOptions) {
		o.AttackId = id
	}
}

func WithSceneDefenceId(id int64) SceneOption {
	return func(o *SceneOptions) {
		o.DefenceId = id
	}
}

func WithSceneAttackUnitList(list []*pbCombat.UnitInfo) SceneOption {
	return func(o *SceneOptions) {
		o.AttackUnitList = list
	}
}

func WithSceneDefenceUnitList(list []*pbCombat.UnitInfo) SceneOption {
	return func(o *SceneOptions) {
		o.DefenceUnitList = list
	}
}

func WithSceneEntry(entry *define.SceneEntry) SceneOption {
	return func(o *SceneOptions) {
		o.Entry = entry
	}
}

type UnitOption func(*UnitOptions)
type UnitOptions struct {
	TypeId   int32
	AttValue []int64
	Position [3]float32
	Entry    *define.UnitEntry
}

func DefaultUnitOptions() *UnitOptions {
	return &UnitOptions{
		TypeId:   -1,
		Position: [3]float32{0, 0, 0},
		Entry:    nil,
	}
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
