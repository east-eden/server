package scene

import (
	"errors"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/utils"
)

var (
	ErrAction_TargetNotFound = errors.New("cannot find target")
)

type ActionHandle func(*Action) error

type Action struct {
	owner     *SceneEntity
	opts      *ActionOptions // 行动参数
	handler   ActionHandle   // 行动对应的处理
	completed bool           // 行动是否结束
	count     int32          // 执行次数
}

func (a *Action) Init(owner *SceneEntity, opts ...ActionOption) {
	a.owner = owner
	a.opts = DefaultActionOptions()
	a.completed = false
	a.count = 0

	for _, o := range opts {
		o(a.opts)
	}
}

func (a *Action) Complete() {
	a.completed = true
}

func (a *Action) IsCompleted() bool {
	return a.completed
}

// helper functions
func (a *Action) GetScene() *Scene {
	return a.owner.GetScene()
}

func (a *Action) GetCamp() *SceneCamp {
	return a.owner.GetCamp()
}

func (a *Action) GetEnemyCamp() (*SceneCamp, bool) {
	enemyCamp := define.GetEnemyCamp(a.GetCamp().camp)
	return a.GetScene().GetSceneCamp(enemyCamp)
}

// 执行行动
func (a *Action) Handle() error {
	a.count++

	switch a.opts.Type {
	case define.CombatAction_Idle:
		return a.handleIdle()
	case define.CombatAction_Attack:
		return a.handleAttack()
	case define.CombatAction_Move:
		return a.handleMove()
	}
	return errors.New("invalid action type")
}

// 空闲行动处理
func (a *Action) handleIdle() error {
	// 空闲执行10次后结束 100ms一次handle？
	if a.count >= 10 {
		a.Complete()
	}
	return nil
}

// 攻击行动处理
func (a *Action) handleAttack() error {
	target, ok := a.GetScene().GetEntity(a.opts.TargetId)
	if !ok {
		return ErrAction_TargetNotFound
	}

	err := a.owner.CombatCtrl.CastSkill(a.owner.opts.NormalSkill, target, false)
	if !utils.ErrCheck(err, "Action CastSpell failed", a.owner.id, a.opts.TargetId) {
		return err
	}

	a.Complete()
	return nil
}

// 移动行动处理
func (a *Action) handleMove() error {

	return nil
}
