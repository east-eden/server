package quest

import (
	"testing"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/services/game/event"
)

// TestMainQuestChain 用真实 MainQuest 配置验证：
// 首个主线任务能正确展开子条件，且在事件驱动下完成并推进到下一个任务。
func TestMainQuestChain(t *testing.T) {
	excel.ReadAllJsonEntries("../../../config/json")

	cfg := config.GetConfigEntry()
	if cfg == nil || cfg.FirstMainQuestId <= 0 {
		t.Skip("config table not loaded")
	}

	em := event.NewEventManager()
	em.Register(define.Event_Type_StagePass, func(e *event.Event) error { return nil })

	var level int32 = 1

	m := NewQuestManager()
	m.Init(
		WithManagerOwnerId(1),
		WithManagerEventManager(em),
		WithManagerOwnerLevel(func() int32 { return level }),
	)

	if q := m.AddMainQuest(cfg.FirstMainQuestId); q == nil {
		t.Fatal("AddMainQuest returned nil")
	}

	q := m.CurMainQuest()
	if q == nil || q.Id != cfg.FirstMainQuestId {
		t.Fatalf("cur main quest = %v, want %d", q, cfg.FirstMainQuestId)
	}

	if len(q.SubConditions) == 0 {
		t.Fatal("first main quest has no sub conditions")
	}

	t.Logf("quest %d sub conditions:", q.Id)
	for _, sc := range q.SubConditions {
		t.Logf("  id=%d type=%v params=%v", sc.Id, sc.Type, sc.Params)
	}

	// 用真实事件驱动完成第一个子条件（章节1-关卡3）
	for _, sc := range q.SubConditions {
		ok, err := m.evalSubCondition(sc, &event.Event{
			Type:  define.Event_Type_StagePass,
			Miscs: []any{sc.Params[len(sc.Params)-1], sc.Params[0]},
		})
		if err != nil {
			t.Fatalf("eval sub condition failed: %v", err)
		}

		if !ok {
			t.Fatalf("stage pass event should satisfy sub condition id=%d params=%v", sc.Id, sc.Params)
		}

		sc.Completed = true
	}

	if !q.CanComplete() {
		t.Fatal("quest should be completable after all sub conditions completed")
	}

	q.Complete()

	if err := m.QuestReward(q.Id); err != nil {
		t.Fatalf("QuestReward failed: %v", err)
	}

	if !q.IsRewarded() {
		t.Fatal("quest should be rewarded")
	}

	if q.Entry.NextId > 0 {
		if next := m.CurMainQuest(); next == nil || next.Id != q.Entry.NextId {
			t.Fatalf("next main quest not advanced, got %v want %d", next, q.Entry.NextId)
		}
	}
}

// TestPlayerLevelCondition 验证角色等级类子条件的判定
func TestPlayerLevelCondition(t *testing.T) {
	excel.ReadAllJsonEntries("../../../config/json")

	em := event.NewEventManager()
	var level int32 = 1

	m := NewQuestManager()
	m.Init(WithManagerEventManager(em), WithManagerOwnerLevel(func() int32 { return level }))

	sc := &SubCondition{
		Id:       1,
		Type:     config.SubConditionType_PlayerLevel,
		Params:   []int32{10},
		MaxCount: 10,
	}

	if ok, _ := m.evalSubCondition(sc, &event.Event{Type: define.Event_Type_PlayerLevelup}); ok {
		t.Fatal("level 1 should not satisfy level 10")
	}

	level = 10
	if ok, _ := m.evalSubCondition(sc, &event.Event{Type: define.Event_Type_PlayerLevelup}); !ok {
		t.Fatal("level 10 should satisfy level 10")
	}
}
