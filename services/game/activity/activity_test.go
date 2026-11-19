package activity

import (
	"errors"
	"testing"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/event"
	"google.golang.org/protobuf/proto"
)

type fakeOwner struct {
	level  int32
	awards []*config.IdCount
	saved  map[string]*ActivityState
	sent   int
}

func (o *fakeOwner) GetLevel() int32 { return o.level }

func (o *fakeOwner) GrantAwardList(awards []*config.IdCount, reason string) {
	o.awards = append(o.awards, awards...)
}

func (o *fakeOwner) SaveActivities(data map[string]*ActivityState) error {
	o.saved = data
	return nil
}

func (o *fakeOwner) SendProtoMessage(msg proto.Message) { o.sent++ }

func loadConfig(t *testing.T) {
	t.Helper()
	excel.ReadAllJsonEntries("../../../config/json")
}

// TestCheckAllOpensCreatetimeActivity createtime 模式的活动在创建时间 + start 之后开启
func TestCheckAllOpensCreatetimeActivity(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{level: 30}
	m := NewManager(owner)

	// sevenDayTask4: createtime +72h 开启，duration 192h 结束
	createdAt := time.Now().Add(-100 * time.Hour)
	m.CheckAll(time.Now(), createdAt)

	state, ok := m.Data["sevenDayTask4"]
	if !ok {
		t.Fatal("sevenDayTask4 not created")
	}

	if !state.Open {
		t.Fatalf("sevenDayTask4 should be open: start=%v end=%v", state.StartAt, state.EndAt)
	}

	if len(state.Tasks) == 0 {
		t.Fatal("activity tasks not initialized")
	}
}

// TestActivityTaskEventAndClaim 事件推进任务进度并领奖
func TestActivityTaskEventAndClaim(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{level: 30}
	m := NewManager(owner)
	m.CheckAll(time.Now(), time.Now().Add(-100*time.Hour))

	// 任务 43：sevenDayTask4，ChapterStage 条件 [章节8, 关卡20]
	task, ok := m.Data["sevenDayTask4"].Tasks[43]
	if !ok {
		t.Fatal("task 43 not initialized")
	}

	// 未完成时不能领奖
	if err := m.ClaimTask("sevenDayTask4", 43); !errors.Is(err, ErrTaskNotCompleted) {
		t.Fatalf("expect not completed, got %v", err)
	}

	// 派发通关事件（miscs = [stageId, chapterId]）
	m.OnEvent(&event.Event{Type: define.Event_Type_StagePass, Miscs: []any{int32(20), int32(8)}})

	if !m.TaskCompleted(task) {
		t.Fatal("task should be completed after event")
	}

	if err := m.ClaimTask("sevenDayTask4", 43); err != nil {
		t.Fatalf("ClaimTask failed: %v", err)
	}

	if len(owner.awards) == 0 {
		t.Fatal("awards not granted")
	}

	if !task.Claimed {
		t.Fatal("task not marked claimed")
	}

	// 重复领取
	if err := m.ClaimTask("sevenDayTask4", 43); !errors.Is(err, ErrTaskAlreadyClaimed) {
		t.Fatalf("expect already claimed, got %v", err)
	}
}

// TestSignIn 签到按天推进，重复/跳天被拒绝
func TestSignIn(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{level: 30}
	m := NewManager(owner)

	// 构造一个开启中的签到活动
	m.Data["wheelIsland"] = &ActivityState{
		Id:    "wheelIsland",
		Open:  true,
		Tasks: map[int32]*TaskState{},
	}

	if err := m.SignIn("wheelIsland", 1); err != nil {
		t.Fatalf("first sign in failed: %v", err)
	}

	if got := m.Data["wheelIsland"].SignDays; got != 1 {
		t.Fatalf("sign days = %d, want 1", got)
	}

	if len(owner.awards) == 0 {
		t.Fatal("sign in awards not granted")
	}

	// 同一天重复签到
	if err := m.SignIn("wheelIsland", 1); !errors.Is(err, ErrSignAlreadyDone) {
		t.Fatalf("expect already done, got %v", err)
	}

	// 跳天
	if err := m.SignIn("wheelIsland", 5); !errors.Is(err, ErrSignInvalidDay) {
		t.Fatalf("expect invalid day, got %v", err)
	}

	// 关闭的活动
	m.Data["wheelIsland"].Open = false
	if err := m.SignIn("wheelIsland", 2); !errors.Is(err, ErrActivityClosed) {
		t.Fatalf("expect closed, got %v", err)
	}
}

// TestGenPB 只下发开启中的活动
func TestGenPB(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{level: 1}
	m := NewManager(owner)
	m.Data["wheelIsland"] = &ActivityState{Id: "wheelIsland", Open: true, SignDays: 2, Tasks: map[int32]*TaskState{}}
	m.Data["closedOne"] = &ActivityState{Id: "closedOne", Open: false, Tasks: map[int32]*TaskState{}}

	list := m.GenPB()
	if len(list) != 1 {
		t.Fatalf("pb list = %d, want 1", len(list))
	}

	if list[0].GetId() != "wheelIsland" || list[0].GetSignDays() != 2 {
		t.Fatalf("unexpected pb: %+v", list[0])
	}

	var _ *pbGlobal.ActivityInfo = list[0]
}
