package scene

import (
	"e.coding.net/mmstudio/blade/server/excel/auto"
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
)

type SceneOption func(*SceneOptions)
type SceneOptions struct {
	AttackId          int64
	DefenceId         int64
	AttackEntityList  []*pbGlobal.EntityInfo
	DefenceEntityList []*pbGlobal.EntityInfo
	SceneEntry        *auto.SceneEntry
	UnitGroupEntry    *auto.UnitGroupEntry
}

func DefaultSceneOptions() *SceneOptions {
	return &SceneOptions{
		AttackId:          -1,
		DefenceId:         -1,
		AttackEntityList:  make([]*pbGlobal.EntityInfo, 0, 10),
		DefenceEntityList: make([]*pbGlobal.EntityInfo, 0, 10),
		SceneEntry:        nil,
		UnitGroupEntry:    nil,
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

func WithSceneUnitGroupEntry(e *auto.UnitGroupEntry) SceneOption {
	return func(o *SceneOptions) {
		o.UnitGroupEntry = e
	}
}
