package scene

import (
	"github.com/east-eden/server/excel/auto"
	pbCombat "github.com/east-eden/server/proto/server/combat"
)

type SceneOption func(*SceneOptions)
type SceneOptions struct {
	AttackId        int64
	DefenceId       int64
	AttackUnitList  []*pbCombat.UnitInfo
	DefenceUnitList []*pbCombat.UnitInfo
	Entry           *auto.SceneEntry
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

func WithSceneEntry(entry *auto.SceneEntry) SceneOption {
	return func(o *SceneOptions) {
		o.Entry = entry
	}
}
