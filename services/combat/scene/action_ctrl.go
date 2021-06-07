package scene

import (
	"container/list"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/utils"
	log "github.com/rs/zerolog/log"
)

type ActionCtrl struct {
	owner      *SceneEntity // 拥有者
	actionList *list.List   // 行动列表
}

func NewActionCtrl(owner *SceneEntity) *ActionCtrl {
	c := &ActionCtrl{
		owner:      owner,
		actionList: list.New(),
	}

	return c
}

func (c *ActionCtrl) Update() {
	log.Info().Int64("owner_id", c.owner.id).Msg("ActionCtrl update")

	c.updateActionList()
}

func (c *ActionCtrl) updateActionList() {
	// 删除完成的行动
	var next *list.Element
	for e := c.actionList.Front(); e != nil; e = next {
		next = e.Next()
		if e.Value.(*Action).IsCompleted() {
			c.actionList.Remove(e)
		}
	}

	// 需要产生新行动
	if c.actionList.Len() == 0 {
		c.createNewAction()
	}

	// 执行当前行动
	curAction := c.actionList.Front().Value.(*Action)
	err := curAction.Handle()
	utils.ErrPrint(err, "action handle failed", curAction.opts.Type, c.owner.id)
}

// 创建新行动
func (c *ActionCtrl) createNewAction() {
	// 还有敌人
	if target, ok := c.findTarget(); ok {
		action := NewAction()
		action.Init(
			c.owner,
			WithActionType(define.CombatAction_Attack),
			WithActionTargetId(target.id),
		)

		c.actionList.PushBack(action)
		return
	}

	// 无事可做，添加空闲行动
	action := NewAction()
	action.Init(
		c.owner,
		WithActionType(define.CombatAction_Idle),
	)

	c.actionList.PushBack(action)
}

// 寻找敌人
func (c *ActionCtrl) findTarget() (*SceneEntity, bool) {
	return c.owner.GetScene().findEnemyEntityByHead(c.owner.GetCamp().camp)
}
