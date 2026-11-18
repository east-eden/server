package condition

import (
	"errors"

	"github.com/east-eden/server/excel/config"
	log "github.com/rs/zerolog/log"
	"golang.org/x/exp/constraints"
)

var (
	ErrConditionLimit = errors.New("condition limit")
)

type ConditionHandler[T constraints.Ordered] func(cur, target T) bool
type ConditionLimiter func() bool

func getOperatorHandler[T constraints.Ordered](operator config.SubConditionOperationType) ConditionHandler[T] {
	switch operator {
	case config.SubConditionOperationType_Equal:
		return handleEqual[T]
	case config.SubConditionOperationType_Reach:
		return handleReach[T]
	case config.SubConditionOperationType_Less:
		return handleLess[T]
	case config.SubConditionOperationType_Lte:
		return handleLte[T]
	default:
		log.Error().Str("operator", operator.String()).Msg("cannot find handler by operator")
		return nil
	}
}

// 检查单个条件是否满足
func CheckSingleCondition[T constraints.Ordered](operator config.SubConditionOperationType, cur, target T, limiters ...ConditionLimiter) bool {
	operHandle := getOperatorHandler[T](operator)

	for _, limiter := range limiters {
		if !limiter() {
			return false
		}
	}

	return operHandle(cur, target)
}

func CheckMultiCondition[T constraints.Ordered](tp config.SubConditionRelation, operators []config.SubConditionOperationType, cur, target []T, limiters ...ConditionLimiter) bool {
	switch tp {
	// 满足所有条件
	case config.SubConditionRelation_And:
		for index, oper := range operators {
			if !CheckSingleCondition(oper, cur[index], target[index], limiters...) {
				return false
			}
		}
		return true

		// 满足一个条件
	case config.SubConditionRelation_Or:
		for index, oper := range operators {
			if CheckSingleCondition(oper, cur[index], target[index], limiters...) {
				return true
			}
		}
		return false
	default:
		log.Error().Int32("condition_type", int32(tp)).Msg("invalid condition type")
		return false
	}
}

// 相等
func handleEqual[T constraints.Ordered](cur, target T) bool {
	return cur == target
}

// 达到
func handleReach[T constraints.Ordered](cur, target T) bool {
	return cur >= target
}

// 小于
func handleLess[T constraints.Ordered](cur, target T) bool {
	return cur < target
}

// 小于等于
func handleLte[T constraints.Ordered](cur, target T) bool {
	return cur <= target
}
