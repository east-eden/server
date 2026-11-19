package recruit

import (
	"errors"
	"testing"
	"time"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
	"google.golang.org/protobuf/proto"
)

type fakeOwner struct {
	items map[int32]int64
	state *State
	sent  int
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

func (o *fakeOwner) GrantLootAndReport(lootId int32, reason string) []Reward {
	rewards := CalcLoot(lootId)
	for _, r := range rewards {
		o.items[r.ItemId] += int64(r.Num)
	}

	return rewards
}

func (o *fakeOwner) SaveRecruit(state *State) error { o.state = state; return nil }

func (o *fakeOwner) SendProtoMessage(msg proto.Message) { o.sent++ }

func loadConfig(t *testing.T) {
	t.Helper()
	excel.ReadAllJsonEntries("../../../config/json")
}

func newManager(o *fakeOwner) *Manager {
	m := NewManager(o)
	m.LoadAll(nil)
	return m
}

// TestFreeRecruit 首次免费使用专属卡池，且受每日次数限制
func TestFreeRecruit(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{items: map[int32]int64{}}
	m := newManager(owner)

	rewards, err := m.Free(time.Now())
	if err != nil {
		t.Fatalf("Free failed: %v", err)
	}

	if len(rewards) == 0 {
		t.Fatal("no rewards from free recruit")
	}

	if m.State.FreeUsed != 1 || !m.State.FirstFree {
		t.Fatalf("state not updated: %+v", m.State)
	}

	// 同日再抽：MaxFreeRecruit=1
	if _, err := m.Free(time.Now()); !errors.Is(err, ErrRecruitFreeLimit) {
		t.Fatalf("expect free limit, got %v", err)
	}

	// 跨天重置
	m.State.Day = "2000-01-01"
	if _, err := m.Free(time.Now()); err != nil {
		t.Fatalf("free after day change failed: %v", err)
	}
}

// TestSeniorRecruitCost 消耗招募券
func TestSeniorRecruitCost(t *testing.T) {
	loadConfig(t)

	cfg := config.GetConfigEntry()
	owner := &fakeOwner{items: map[int32]int64{cfg.SeniorRecruitCost.Id: 10}}
	m := newManager(owner)

	rewards, err := m.Senior(1)
	if err != nil {
		t.Fatalf("Senior failed: %v", err)
	}

	if len(rewards) == 0 {
		t.Fatal("no rewards")
	}

	if got := owner.items[cfg.SeniorRecruitCost.Id]; got != 10-int64(cfg.SeniorRecruitCost.Count) {
		t.Fatalf("ticket left = %d", got)
	}

	if m.State.SeniorTimes != 1 {
		t.Fatalf("senior times = %d", m.State.SeniorTimes)
	}
}

// TestSeniorRecruitTen 十连抽累计次数与首次十连卡池
func TestSeniorRecruitTen(t *testing.T) {
	loadConfig(t)

	cfg := config.GetConfigEntry()
	owner := &fakeOwner{items: map[int32]int64{cfg.SeniorRecruitCost.Id: 100}}
	m := newManager(owner)

	rewards, err := m.Senior(10)
	if err != nil {
		t.Fatalf("Senior ten failed: %v", err)
	}

	if len(rewards) == 0 {
		t.Fatal("no rewards")
	}

	if m.State.SeniorTimes != 10 {
		t.Fatalf("senior times = %d, want 10", m.State.SeniorTimes)
	}

	if !m.State.FirstTen || !m.State.FirstSenior {
		t.Fatalf("first flags not set: %+v", m.State)
	}
}

// TestSeniorGuarantee 小保底触发专属掉落
func TestSeniorGuarantee(t *testing.T) {
	loadConfig(t)

	cfg := config.GetConfigEntry()
	owner := &fakeOwner{items: map[int32]int64{cfg.SeniorRecruitCost.Id: 10}}
	m := newManager(owner)

	// 差一次触发小保底
	m.State.SmallGuarantee = cfg.SeniorRecruitSmallGuaranteeTimes - 1

	if _, err := m.Senior(1); err != nil {
		t.Fatalf("Senior failed: %v", err)
	}

	if m.State.SmallGuarantee != 0 {
		t.Fatalf("small guarantee not reset: %d", m.State.SmallGuarantee)
	}
}

// TestSeniorInvalidTimes 只允许单抽/十连
func TestSeniorInvalidTimes(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{items: map[int32]int64{}}
	m := newManager(owner)

	if _, err := m.Senior(3); !errors.Is(err, ErrRecruitTimesInvalid) {
		t.Fatalf("expect invalid times, got %v", err)
	}
}
