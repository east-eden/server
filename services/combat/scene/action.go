package scene

import (
	"errors"

	"github.com/east-eden/server/define"
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
	// 空闲执行10次后结束
	if a.count >= 10 {
		a.completed = true
	}
	return nil
}

// 攻击行动处理
func (a *Action) handleAttack() error {

	return nil
}

// 移动行动处理
func (a *Action) handleMove() error {

	return nil
}
