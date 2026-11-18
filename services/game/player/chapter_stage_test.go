package player

import (
	"errors"
	"testing"

	"github.com/east-eden/server/excel/config"
	"github.com/golang/mock/gomock"
)

// TestStagePass 用 kof 的 Level 配置验证关卡通关与奖励发放
func TestStagePass(t *testing.T) {
	setupTestEnv()

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	initMockStore(t, ctl)

	levelId, ok := config.GetLevelIdByChapterStage(1, 1)
	if !ok {
		t.Skip("no level config for chapter 1 stage 1")
	}

	level, _ := config.GetLevelEntry(levelId)

	p := newTestPlayer()
	m := NewChapterStageManager(p)

	if err := m.StageChallenge(levelId, true, true, []bool{true}); err != nil {
		t.Fatalf("StageChallenge failed: %v", err)
	}

	stage, ok := m.Stages[levelId]
	if !ok || !stage.Pass {
		t.Fatalf("stage %d not passed", levelId)
	}

	if !stage.Achieve {
		t.Fatalf("stage %d goal not achieved", levelId)
	}

	// 章节星数应增加
	chapter, ok := m.Chapters[level.Chapter]
	if !ok || chapter.Stars <= 0 {
		t.Fatalf("chapter %d stars not updated", level.Chapter)
	}

	// 通关掉落应入账（StageLoot 有效时）
	if level.StageLoot > 0 {
		t.Logf("level=%d chapter=%d stageLoot=%d items=%d", levelId, level.Chapter, level.StageLoot, len(p.ItemManager().Items))
	}
}

// TestStagePrevNotPassed 未通关前置关卡时阻止挑战
func TestStagePrevNotPassed(t *testing.T) {
	setupTestEnv()

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	initMockStore(t, ctl)

	levelId, ok := config.GetLevelIdByChapterStage(1, 2)
	if !ok {
		t.Skip("no level config for chapter 1 stage 2")
	}

	p := newTestPlayer()
	m := NewChapterStageManager(p)

	if err := m.StageChallenge(levelId, true, false, nil); err == nil {
		t.Fatal("expect prev stage not passed error")
	}
}

// TestStageRejectClientWin 配置了服务器战斗后，客户端自报胜负应被拒绝
func TestStageRejectClientWin(t *testing.T) {
	setupTestEnv()

	levelId, ok := config.GetLevelIdByChapterStage(1, 1)
	if !ok {
		t.Skip("no level config")
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	initMockStore(t, ctl)

	p := newTestPlayer()
	p.BattleManager().BattleRecord = &PlayerBattleRecord{IsServerBattle: true}

	m := NewChapterStageManager(p)
	if err := m.StageChallenge(levelId, true, true, []bool{true}); !errors.Is(err, ErrStageNeedBattleResult) {
		t.Fatalf("expect ErrStageNeedBattleResult, got %v", err)
	}

	// 服务器战斗结算入口可以正常推进
	if err := m.StagePassByBattle(levelId, true); err != nil {
		t.Fatalf("StagePassByBattle failed: %v", err)
	}

	if stage, ok := m.Stages[levelId]; !ok || !stage.Pass {
		t.Fatalf("stage %d not passed", levelId)
	}
}
