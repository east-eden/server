// Package subcondition 提供子条件（SubCondition）的运行时模型与事件判定。
//
// kof 里子条件判定分散在任务、活动、成就等处，这里抽出共享实现：
// 任务与活动共用同一套模型与判定逻辑，避免两边各写一份导致口径不一致。
package subcondition

import (
	"errors"
	"fmt"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/event"
)

var (
	ErrEventParamInvalid  = errors.New("sub condition event params invalid")
	ErrConditionUnsupport = errors.New("sub condition type not supported")
)

// SubCondition 子条件运行时状态
type SubCondition struct {
	Id        int32                              `bson:"id" json:"id"`               // 子条件id
	Type      config.SubConditionType            `bson:"type" json:"type"`           // 子条件类型
	Params    []int32                            `bson:"params" json:"params"`       // 子条件参数
	OpTypes   []config.SubConditionOperationType `bson:"opTypes" json:"opTypes"`     // 比较方式
	Count     int32                              `bson:"count" json:"count"`         // 当前计数
	MaxCount  int32                              `bson:"maxCount" json:"maxCount"`   // 目标计数
	Completed bool                               `bson:"completed" json:"completed"` // 是否达成
}

// GenPB 生成下发用的目标结构
func (sc *SubCondition) GenPB() *pbGlobal.QuestObj {
	return &pbGlobal.QuestObj{
		Type:      int32(sc.Type),
		Count:     sc.Count,
		Completed: sc.Completed,
	}
}

// SupportedEventTypes 已接入的子条件类型与其对应的事件类型。
//
// kof 的任务与活动会用到大量子条件类型，这里只接当前已有的系统
// （角色等级、关卡、爬塔），其余类型在判定时会显式返回不支持，避免“看起来完成了”。
var SupportedEventTypes = map[config.SubConditionType]int{
	config.SubConditionType_PlayerLevel:    define.Event_Type_PlayerLevelup,
	config.SubConditionType_ChapterStage:   define.Event_Type_StagePass,
	config.SubConditionType_TowerFloorPass: define.Event_Type_TowerPass,
}

// New 按子条件配置构造运行时模型
func New(id int32) (*SubCondition, bool) {
	entry, ok := config.GetSubConditionEntry(id)
	if !ok {
		return nil, false
	}

	return &SubCondition{
		Id:      entry.Id,
		Type:    entry.SubCondType,
		Params:  entry.SubCondParams,
		OpTypes: entry.OperationType,
	}, true
}

// BuildFromConditionId 按主条件展开子条件。
//
// 主条件可能是复合条件（Condition）也可能是直接子条件，统一展开成子条件列表；
// 任一缺失就返回空列表，由调用方判定为“未配置”。
func BuildFromConditionId(conditionId int32) []*SubCondition {
	subs := make([]*SubCondition, 0, 2)

	if cond, ok := config.GetConditionEntry(conditionId); ok {
		for _, sid := range cond.SubCondIds {
			if sc, ok := New(sid); ok {
				subs = append(subs, sc)
			}
		}

		return subs
	}

	if sc, ok := New(conditionId); ok {
		subs = append(subs, sc)
	}

	return subs
}

// Eval 用事件判定子条件，返回是否达成。playerLevel 由调用方提供（子条件可能按玩家等级判定）。
//
// 已支持的类型会更新 sc.Count / sc.MaxCount，调用方据返回值决定是否置为完成。
func Eval(sc *SubCondition, e *event.Event, playerLevel int32) (bool, error) {
	switch sc.Type {
	case config.SubConditionType_PlayerLevel:
		if len(sc.Params) == 0 {
			return false, ErrEventParamInvalid
		}

		sc.Count = playerLevel
		sc.MaxCount = sc.Params[0]

		return playerLevel >= sc.Params[0], nil

	case config.SubConditionType_ChapterStage:
		if len(sc.Params) == 0 || len(e.Miscs) == 0 {
			return false, ErrEventParamInvalid
		}

		stageId, ok := asInt32(e.Miscs[0])
		if !ok {
			return false, ErrEventParamInvalid
		}

		var chapterId int32
		if len(e.Miscs) > 1 {
			chapterId, _ = asInt32(e.Miscs[1])
		}

		sc.MaxCount = 1

		// 配置里通常给 [章节id, 关卡id]；单参数时按关卡id匹配，-1 表示不限。
		switch len(sc.Params) {
		case 1:
			return sc.Params[0] == -1 || sc.Params[0] == stageId, nil
		default:
			if sc.Params[0] != -1 && sc.Params[0] != chapterId {
				return false, nil
			}

			return sc.Params[1] == -1 || sc.Params[1] == stageId, nil
		}

	case config.SubConditionType_TowerFloorPass:
		if len(sc.Params) == 0 || len(e.Miscs) == 0 {
			return false, ErrEventParamInvalid
		}

		floor, ok := asInt32(e.Miscs[0])
		if !ok {
			return false, ErrEventParamInvalid
		}

		sc.MaxCount = 1

		return floor >= sc.Params[0], nil

	default:
		return false, fmt.Errorf("%w: %v", ErrConditionUnsupport, sc.Type)
	}
}

func asInt32(v any) (int32, bool) {
	switch n := v.(type) {
	case int32:
		return n, true
	case int:
		return int32(n), true
	case int64:
		return int32(n), true
	case uint32:
		return int32(n), true
	default:
		return 0, false
	}
}
