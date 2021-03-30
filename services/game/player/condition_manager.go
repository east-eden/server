package player

import (
	"errors"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
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
	m.handlers[define.Condition_SubType_VipLevel_Achieve] = m.handleVipLevelAchieve
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

// vip等级达到**级
func (m *ConditionManager) handleVipLevelAchieve(value int32) bool {
	return m.owner.VipLevel >= value
}
