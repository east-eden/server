// Package teamleague 组队联赛任务构造。
//
// 对应 kof 的 NewTeamLeagueTask：特殊任务用充值额度作为目标上限，
// 普通任务由 TeamLeagueTask 的 SubConditionId 或 Condition 展开的子条件组成。
// kof 直接返回 ts 协议里的 QuestInfo，这里用等价的本地结构 + ToPB 适配 protobuf。
package teamleague

import (
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
)

// SubCondition 任务子条件
type SubCondition struct {
	SubConditionId int32 `json:"subConditionId"`
	Count          int32 `json:"count"`
	MaxCount       int32 `json:"maxCount"`
}

// Task 组队联赛任务
type Task struct {
	TaskId        int32                   `json:"taskId"`
	State         int32                   `json:"state"`
	SubConditions map[int32]*SubCondition `json:"subConditions"`
}

// NewTask 构造组队联赛任务
func NewTask(taskId int32) *Task {
	task := &Task{
		TaskId:        taskId,
		State:         int32(config.QuestState_Open),
		SubConditions: make(map[int32]*SubCondition),
	}

	// 特殊任务：目标上限是充值额度
	if specialEntry, found := config.GetTeamLeagueSpecialTaskEntry(taskId); found {
		task.SubConditions[0] = &SubCondition{
			SubConditionId: -1,
			Count:          0,
			MaxCount:       specialEntry.RechargeMoney,
		}

		return task
	}

	taskEntry, found := config.GetTeamLeagueTaskEntry(taskId)
	if !found {
		return task
	}

	maxCnt := taskEntry.CondMaxCnt

	if taskEntry.SubConditionId > 0 {
		task.SubConditions[int32(len(task.SubConditions))] = &SubCondition{
			SubConditionId: taskEntry.SubConditionId,
			Count:          0,
			MaxCount:       maxCnt,
		}
	}

	if taskEntry.ConditionId > 0 {
		if condEntry, ok := config.GetConditionEntry(taskEntry.ConditionId); ok {
			for _, subCondId := range condEntry.SubCondIds {
				task.SubConditions[int32(len(task.SubConditions))] = &SubCondition{
					SubConditionId: subCondId,
					Count:          0,
					MaxCount:       maxCnt,
				}
			}
		}
	}

	return task
}

// ToPB 转成下发用的任务结构。
// 协议里 QuestObj 没有 MaxCount，上限仍以配置为准。
func (t *Task) ToPB() *pbGlobal.Quest {
	pb := &pbGlobal.Quest{
		Id:    t.TaskId,
		State: t.State,
		Objs:  make([]*pbGlobal.QuestObj, 0, len(t.SubConditions)),
	}

	for _, sub := range t.SubConditions {
		pb.Objs = append(pb.Objs, &pbGlobal.QuestObj{
			Type:  sub.SubConditionId,
			Count: sub.Count,
		})
	}

	return pb
}
