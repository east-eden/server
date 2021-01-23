package scene

import (
	"errors"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/utils"
)

var (
	ErrAction_TargetNotFound = errors.New("cannot find target")
)

type ActionHandle func(*Action) error

type Action struct {
	opts      *ActionOptions // 行动参数
	handler   ActionHandle   // 行动对应的处理
	completed bool           // 行动是否结束
	count     int32          // 执行次数
}

func NewAction() *Action {
	return &Action{
		opts:      DefaultActionOptions(),
		completed: false,
		count:     0,
	}
}

func (a *Action) Init(opts ...ActionOption) {
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
func (a *Action) getScene() *Scene {
	return a.opts.Owner.scene
}

func (a *Action) getCamp() *SceneCamp {
	return a.opts.Owner.camp
}

func (a *Action) getEnemyCamp() (*SceneCamp, bool) {
	return a.getScene().GetSceneCamp(a.getCamp().GetOtherCamp())
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
	enemyCamp, ok := a.getEnemyCamp()
	if !ok {
		return errors.New("cannot get enemy camp")
	}

	target, ok := enemyCamp.GetUnit(a.opts.TargetId)
	if !ok {
		return ErrAction_TargetNotFound
	}

	owner := a.opts.Owner
	err := owner.CombatCtrl().CastSpell(owner.normalSpell, owner, target, false)
	if !utils.ErrCheck(err, "Action CastSpell failed", a.opts.Owner.id, a.opts.TargetId) {
		return err
	}

	a.Complete()
	return nil
}

// 移动行动处理
func (a *Action) handleMove() error {

	return nil
}
