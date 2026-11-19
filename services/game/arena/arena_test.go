package arena

import (
	"errors"
	"testing"
	"time"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
)

type fakeOwner struct {
	items  map[int32]int64
	awards []*config.IdCount
	state  *State
}

func (o *fakeOwner) GetItemCount(id int32) int64 { return o.items[id] }

func (o *fakeOwner) CostItems(costs []*config.IdCount, reason string) error {
	for _, c := range costs {
		if o.items[c.Id] < int64(c.Count) {
			return errors.New("not enough")
		}
	}

	for _, c := range costs {
		o.items[c.Id] -= int64(c.Count)
	}

	return nil
}

func (o *fakeOwner) GrantAwardList(awards []*config.IdCount, reason string) {
	o.awards = append(o.awards, awards...)
}

func (o *fakeOwner) SaveArena(state *State) error { o.state = state; return nil }

func loadConfig(t *testing.T) {
	t.Helper()
	excel.ReadAllJsonEntries("../../../config/json")
}

// TestCalcNewScore ELO 公式：赢涨、输跌、不低于 1
func TestCalcNewScore(t *testing.T) {
	win := CalcNewScore(1000, 1000, true)
	if win <= 1000 {
		t.Fatalf("win score = %d, want > 1000", win)
	}

	lose := CalcNewScore(1000, 1000, false)
	if lose >= 1000 {
		t.Fatalf("lose score = %d, want < 1000", lose)
	}

	// 越级挑战赢加分更多
	upset := CalcNewScore(1000, 2000, true)
	normal := CalcNewScore(1000, 1000, true)
	if upset <= normal {
		t.Fatalf("upset %d should beat normal %d", upset, normal)
	}

	// 下限保护
	if got := CalcNewScore(1, 5000, false); got < 1 {
		t.Fatalf("score floor violated: %d", got)
	}
}

// TestChallengeCost 按 ArenaConfig 消耗挑战券
func TestChallengeCost(t *testing.T) {
	loadConfig(t)

	cfg := config.GetArenaConfigEntry()
	if len(cfg.ChallengeCost) < 2 {
		t.Skip("no challenge cost")
	}

	id, count := cfg.ChallengeCost[0], cfg.ChallengeCost[1]
	owner := &fakeOwner{items: map[int32]int64{id: int64(count) + 1}}
	m := NewManager(owner)
	m.LoadAll(nil)

	if err := m.Challenge(); err != nil {
		t.Fatalf("Challenge failed: %v", err)
	}

	if got := owner.items[id]; got != 1 {
		t.Fatalf("cost not applied: %d", got)
	}

	// 第二次刚好够
	if err := m.Challenge(); err != nil {
		t.Fatalf("second Challenge failed: %v", err)
	}

	// 第三次不足
	if err := m.Challenge(); err == nil {
		t.Fatal("expect cost error")
	}
}

// TestApplyResult 结果影响积分与周胜场
func TestApplyResult(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{items: map[int32]int64{}}
	m := NewManager(owner)
	m.LoadAll(&State{Score: 1000})

	delta := m.ApplyResult(1000, true)
	if delta <= 0 || m.State.Score != 1000+delta {
		t.Fatalf("delta = %d, score = %d", delta, m.State.Score)
	}

	if m.State.WeekWins != 1 {
		t.Fatalf("week wins = %d", m.State.WeekWins)
	}

	// 同分输掉会掉分（打高分对手输掉时 ELO 变化接近 0，属于预期）
	delta = m.ApplyResult(1000, false)
	if delta >= 0 {
		t.Fatalf("lose delta = %d", delta)
	}

	if m.State.WeekWins != 1 {
		t.Fatalf("week wins should not change on lose: %d", m.State.WeekWins)
	}
}

// TestClaimRewards 日榜/周榜奖励领取与重复领取
func TestClaimRewards(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{items: map[int32]int64{}}
	m := NewManager(owner)
	m.LoadAll(&State{Score: 1000})

	if err := m.ClaimDaily(1); err != nil {
		t.Fatalf("ClaimDaily failed: %v", err)
	}

	if len(owner.awards) == 0 {
		t.Fatal("no daily awards")
	}

	if err := m.ClaimDaily(1); !errors.Is(err, ErrDailyRewardClaimed) {
		t.Fatalf("expect daily claimed, got %v", err)
	}

	if err := m.ClaimWeekly(1); err != nil {
		t.Fatalf("ClaimWeekly failed: %v", err)
	}

	if err := m.ClaimWeekly(1); !errors.Is(err, ErrWeeklyRewardClaimed) {
		t.Fatalf("expect weekly claimed, got %v", err)
	}
}

// TestResetIfNeeded 跨天/跨周重置领取状态
func TestResetIfNeeded(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{items: map[int32]int64{}}
	m := NewManager(owner)
	m.LoadAll(&State{Score: 1000, DailyClaimed: true, WeeklyClaimed: true, WeekWins: 3, Day: "2000-01-01", Week: "2000-01-03"})

	m.ResetIfNeeded(time.Now())

	if m.State.DailyClaimed || m.State.WeeklyClaimed {
		t.Fatalf("claims not reset: %+v", m.State)
	}

	if m.State.WeekWins != 0 {
		t.Fatalf("week wins = %d, want 0", m.State.WeekWins)
	}
}
