package player

import (
	"errors"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
)

var (
	ErrConditionLimit = errors.New("condition limit")
)

type ConditionHandler func(value int32) bool
type ConditionLimiter func() bool

type ConditionManager struct {
	owner    *Player                    `bson:"-" json:"-"`
	handlers map[int32]ConditionHandler `bson:"-" json:"-"`
}

func NewConditionManager(owner *Player) *ConditionManager {
	m := &ConditionManager{
		owner:    owner,
		handlers: make(map[int32]ConditionHandler),
	}

	m.initHandlers()
	return m
}

func (m *ConditionManager) initHandlers() {
	m.handlers[define.Condition_SubType_TeamLevel_Achieve] = m.handleTeamLevelAchieve
	m.handlers[define.Condition_SubType_KillAllEnemy] = m.handleKillAllEnemy
	m.handlers[define.Condition_SubType_OurUnitAllDead] = m.handleOurUnitAllDead
	m.handlers[define.Condition_SubType_OurUnitDeadLessThan] = m.handleOurUnitDeadLessThan
	m.handlers[define.Condition_SubType_InterruptEnemySkill] = m.handleInterruptEnemySkill
	m.handlers[define.Condition_SubType_OurUnitCastUltimateSkill] = m.handleOurUnitCastUltimateSkill
	m.handlers[define.Condition_SubType_CombatPassTimeLessThan] = m.handleCombatPassTimeLessThan
	m.handlers[define.Condition_SubType_KillEnemyTypeIdFirst] = m.handleKillEnemyTypeIdFirst
}

// 检查条件是否满足
func (m *ConditionManager) CheckCondition(conditionId int32, limiters ...ConditionLimiter) bool {
	if conditionId == -1 {
		return true
	}

	entry, ok := auto.GetConditionEntry(conditionId)
	if !ok {
		return false
	}

	switch entry.Type {
	// 满足所有条件
	case define.Condition_Type_And:
		for k, tp := range entry.SubTypes {
			if tp == -1 {
				continue
			}

			h, ok := m.handlers[tp]
			if !ok {
				return false
			}

			for _, limiter := range limiters {
				if !limiter() {
					return false
				}
			}

			if !h(entry.SubValues[k]) {
				return false
			}
		}
		return true

		// 满足一个条件
	case define.Condition_Type_Or:
		for k, tp := range entry.SubTypes {
			if tp == -1 {
				continue
			}

			h, ok := m.handlers[tp]
			if !ok {
				continue
			}

			for _, limiter := range limiters {
				if limiter() {
					return true
				}
			}

			if h(entry.SubValues[k]) {
				return true
			}
		}
	}

	return false
}

// 队伍等级达到**级
func (m *ConditionManager) handleTeamLevelAchieve(value int32) bool {
	return m.owner.Level >= value
}

// 击杀所有敌方单位
func (m *ConditionManager) handleKillAllEnemy(value int32) bool {
	return true
}

// 己方单位全部死亡
func (m *ConditionManager) handleOurUnitAllDead(value int32) bool {
	return true
}

// 己方单位死亡人数小于*
func (m *ConditionManager) handleOurUnitDeadLessThan(value int32) bool {
	return true
}

// 成功打断*次敌方技能
func (m *ConditionManager) handleInterruptEnemySkill(value int32) bool {
	return true
}

// 己方单位成功使用*次奥义技能
func (m *ConditionManager) handleOurUnitCastUltimateSkill(value int32) bool {
	return true
}

// 通关时间小于*秒
func (m *ConditionManager) handleCombatPassTimeLessThan(value int32) bool {
	return true
}

// 优先击杀id为*的敌方单位
func (m *ConditionManager) handleKillEnemyTypeIdFirst(value int32) bool {
	return true
}
