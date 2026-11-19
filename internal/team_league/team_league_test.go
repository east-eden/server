package teamleague

import (
	"testing"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
)

func loadConfig(t *testing.T) {
	t.Helper()
	excel.ReadAllJsonEntries("../../config/json")
}

// TestNewTaskWithRealConfig 用真实配置构造普通任务
func TestNewTaskWithRealConfig(t *testing.T) {
	loadConfig(t)

	var taskId int32
	for id := range config.GetTeamLeagueTaskRows() {
		taskId = id
		break
	}

	if taskId == 0 {
		t.Skip("no team league task config")
	}

	task := NewTask(taskId)
	if task.TaskId != taskId {
		t.Fatalf("task id = %d, want %d", task.TaskId, taskId)
	}

	if task.State != int32(config.QuestState_Open) {
		t.Fatalf("state = %d", task.State)
	}

	if len(task.SubConditions) == 0 {
		t.Fatalf("task %d has no sub conditions", taskId)
	}

	pb := task.ToPB()
	if pb.Id != taskId || len(pb.Objs) != len(task.SubConditions) {
		t.Fatalf("pb mismatch: %+v", pb)
	}
}

// TestNewTaskUnknown 未知任务返回空任务而不是 panic
func TestNewTaskUnknown(t *testing.T) {
	loadConfig(t)

	task := NewTask(999999999)
	if len(task.SubConditions) != 0 {
		t.Fatalf("expect empty sub conditions, got %d", len(task.SubConditions))
	}
}
