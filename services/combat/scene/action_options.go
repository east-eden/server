package scene

import (
	"e.coding.net/mmstudio/blade/server/define"
)

type ActionOption func(*ActionOptions)
type ActionOptions struct {
	Type       define.ECombatActionType // 行动类型
	Owner      *SceneUnit               // 拥有者
	TargetId   int64                    // 行动目标单位id
	TargetPosX int32                    // 行动目标坐标
	TargetPosY int32
}

func DefaultActionOptions() *ActionOptions {
	return &ActionOptions{
		Type:       define.CombatAction_Idle,
		Owner:      nil,
		TargetId:   -1,
		TargetPosX: 0,
		TargetPosY: 0,
	}
}

func WithActionType(tp define.ECombatActionType) ActionOption {
	return func(o *ActionOptions) {
		o.Type = tp
	}
}

func WithActionOwner(owner *SceneUnit) ActionOption {
	return func(o *ActionOptions) {
		o.Owner = owner
	}
}

func WithActionTargetId(targetId int64) ActionOption {
	return func(o *ActionOptions) {
		o.TargetId = targetId
	}
}

func WithActionTargetPosX(x int32) ActionOption {
	return func(o *ActionOptions) {
		o.TargetPosX = x
	}
}

func WithActionTargetPosY(y int32) ActionOption {
	return func(o *ActionOptions) {
		o.TargetPosY = y
	}
}
