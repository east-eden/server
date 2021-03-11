package scene

import (
	pbCombat "bitbucket.org/funplus/server/proto/server/combat"
)

type SceneOption func(*SceneOptions)
type SceneOptions struct {
	AttackId        int64
	DefenceId       int64
	AttackUnitList  []*pbCombat.UnitInfo
	DefenceUnitList []*pbCombat.UnitInfo
}

func DefaultSceneOptions() *SceneOptions {
	return &SceneOptions{
		AttackId:  -1,
		DefenceId: -1,
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
