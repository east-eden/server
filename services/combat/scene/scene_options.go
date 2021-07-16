package scene

import (
	"github.com/east-eden/server/excel/auto"
	pbGlobal "github.com/east-eden/server/proto/global"
)

type SceneOption func(*SceneOptions)
type SceneOptions struct {
	AttackId          int64
	DefenceId         int64
	AttackEntityList  []*pbGlobal.EntityInfo
	DefenceEntityList []*pbGlobal.EntityInfo
	SceneEntry        *auto.SceneEntry
	BattleWaveEntries []*auto.BattleWaveEntry
}

func DefaultSceneOptions() *SceneOptions {
	return &SceneOptions{
		AttackId:          -1,
		DefenceId:         -1,
		AttackEntityList:  make([]*pbGlobal.EntityInfo, 0, 10),
		DefenceEntityList: make([]*pbGlobal.EntityInfo, 0, 10),
		SceneEntry:        nil,
		BattleWaveEntries: make([]*auto.BattleWaveEntry, 0, 3),
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

func WithSceneAttackUnitList(list []*pbGlobal.EntityInfo) SceneOption {
	return func(o *SceneOptions) {
		o.AttackEntityList = append(o.AttackEntityList, list...)
	}
}

func WithSceneDefenceUnitList(list []*pbGlobal.EntityInfo) SceneOption {
	return func(o *SceneOptions) {
		o.DefenceEntityList = append(o.DefenceEntityList, list...)
	}
}

func WithSceneEntry(e *auto.SceneEntry) SceneOption {
	return func(o *SceneOptions) {
		o.SceneEntry = e
	}
}

func WithSceneBattleWaveEntries(e ...*auto.BattleWaveEntry) SceneOption {
	return func(o *SceneOptions) {
		o.BattleWaveEntries = append(o.BattleWaveEntries, e...)
	}
}
