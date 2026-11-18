package config

import (
	"sort"
	"testing"

	"github.com/east-eden/server/excel"
)

// TestMainQuestConditionClosure 校验主线任务引用的子条件都能解析到配置，
// 并输出实际用到的子条件类型分布，作为后续移植任务系统的范围依据。
func TestMainQuestConditionClosure(t *testing.T) {
	excel.ReadAllJsonEntries("../../config/json")

	quests := GetMainQuestRows()
	if len(quests) == 0 {
		t.Skip("no main quest table loaded")
	}

	typeCount := make(map[SubConditionType]int)

	for id, q := range quests {
		sub, ok := GetSubConditionEntry(q.SubConditionId)
		if !ok {
			t.Fatalf("main quest %d references missing sub condition %d", id, q.SubConditionId)
		}

		typeCount[sub.SubCondType]++
	}

	types := make([]int, 0, len(typeCount))
	for tp := range typeCount {
		types = append(types, int(tp))
	}
	sort.Ints(types)

	t.Logf("main quests=%d distinct sub condition types=%d", len(quests), len(types))
	for _, tp := range types {
		t.Logf("  %-36s %d", SubConditionType(tp).String(), typeCount[SubConditionType(tp)])
	}
}
