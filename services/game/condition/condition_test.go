package condition

import (
	"testing"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
)

func TestCheckSingleCondition(t *testing.T) {
	cases := []struct {
		op     config.SubConditionOperationType
		cur    int32
		target int32
		want   bool
	}{
		{config.SubConditionOperationType_Equal, 5, 5, true},
		{config.SubConditionOperationType_Equal, 4, 5, false},
		{config.SubConditionOperationType_Reach, 6, 5, true},
		{config.SubConditionOperationType_Reach, 4, 5, false},
		{config.SubConditionOperationType_Reach, 5, 5, true},
		{config.SubConditionOperationType_Less, 4, 5, true},
		{config.SubConditionOperationType_Less, 5, 5, false},
		{config.SubConditionOperationType_Lte, 5, 5, true},
		{config.SubConditionOperationType_Lte, 6, 5, false},
	}

	for _, c := range cases {
		if got := CheckSingleCondition(c.op, c.cur, c.target); got != c.want {
			t.Fatalf("op=%v cur=%d target=%d: got %v, want %v", c.op, c.cur, c.target, got, c.want)
		}
	}
}

// TestConditionConfigLoaded 用真实 Condition/SubCondition 表验证配置可用
func TestConditionConfigLoaded(t *testing.T) {
	excel.ReadAllJsonEntries("../../../config/json")

	rows := config.GetConditionRows()
	if len(rows) == 0 {
		t.Skip("no condition table loaded")
	}

	for id, entry := range rows {
		for _, subId := range entry.SubCondIds {
			if _, ok := config.GetSubConditionEntry(subId); !ok {
				t.Fatalf("condition %d references missing sub condition %d", id, subId)
			}
		}

		return
	}
}
