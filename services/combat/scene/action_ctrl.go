package scene

import (
	"container/list"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

type ActionCtrl struct {
	owner      *SceneUnit // 拥有者
	actionList *list.List // 行动列表
}

func NewActionCtrl(owner *SceneUnit) *ActionCtrl {
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
	// 需要产生新行动
	if c.actionList.Len() == 0 {
		c.createNewAction()
	}

	// 执行当前行动
	curAction := c.actionList.Front().Value.(*Action)
	err := curAction.Handle()
	if event, pass := utils.ErrCheck(err, curAction.opts.Type, c.owner.id); !pass {
		event.Msg("action handle failed")
	}
}

// 创建新行动
func (c *ActionCtrl) createNewAction() {
	// 还有敌人
	if target, ok := c.findTarget(); ok {
		action := c.owner.scene.CreateAction()
		action.Init(
			WithActionOwner(c.owner),
			WithActionType(define.CombatAction_Attack),
			WithActionTarget(target),
		)

		c.actionList.PushBack(action)
		return
	}

	// 无事可做，添加空闲行动
	action := c.owner.scene.CreateAction()
	action.Init(
		WithActionOwner(c.owner),
		WithActionType(define.CombatAction_Idle),
	)

	c.actionList.PushBack(action)
}

// 寻找敌人
func (c *ActionCtrl) findTarget() (*SceneUnit, bool) {
	enemyCamp, ok := c.owner.scene.GetSceneCamp(c.owner.camp.GetOtherCamp())
	if ok && enemyCamp.GetUnitsLen() > 0 {
		return enemyCamp.FindUnitByHead()
	}

	return nil, false
}

//-----------------------------------------------------------------------------
// 攻击
//-----------------------------------------------------------------------------
// func (c *SceneCamp) Attack(Entity* pEntity)
// {
// 	EntityGroup* pTarget = static_cast<EntityGroup*>(pEntity);
// 	BOOL bBreak = FALSE;
// 	for( INT32 i = m_n16LoopIndex; i < X_Max_Summon_Num; ++i )
// 	{
// 		++m_n16LoopIndex;

// 		if( VALID(m_ArrayHero[i]) && m_ArrayHero[i]->IsValid() )
// 		{
// 			EntityHero* pHero = FindTargetByPriority(i, pTarget, TRUE);

// 			if( VALID(pHero) )
// 			{
// 				m_ArrayHero[i]->Attack(pHero);
// 				m_ArrayHero[i]->GetCombatController().CalAuraEffect(GetScene()->GetCurRound());

// 				// 风怒状态
// 				if( m_ArrayHero[i]->HasState(EHS_Anger) )
// 				{
// 					EntityHero* pHero = FindTargetByPriority(i, pTarget, TRUE);
// 					if( VALID(pHero) )
// 					{
// 						m_ArrayHero[i]->Attack(pHero);
// 					}
// 				}

// 				AddAttackNum();
// 				bBreak = TRUE;
// 			}
// 		}

// 		if( bBreak )
// 			break;
// 	}
// }
