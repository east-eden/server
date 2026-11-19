package alliance

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/utils"
)

// memStore 内存实现，便于测试
type memStore struct {
	guilds map[int64]*Guild
	names  map[string]int64
}

func newMemStore() *memStore {
	return &memStore{guilds: map[int64]*Guild{}, names: map[string]int64{}}
}

func (s *memStore) Load(ctx context.Context, id int64) (*Guild, error) {
	g, ok := s.guilds[id]
	if !ok {
		return nil, nil
	}

	return g, nil
}

func (s *memStore) Save(ctx context.Context, g *Guild) error {
	s.guilds[g.Id] = g
	return nil
}

func (s *memStore) Delete(ctx context.Context, id int64) error {
	delete(s.guilds, id)
	return nil
}

func (s *memStore) FindByName(ctx context.Context, name string) (int64, bool, error) {
	id, ok := s.names[name]
	return id, ok, nil
}

func (s *memStore) SaveName(ctx context.Context, name string, id int64) error {
	s.names[name] = id
	return nil
}

type fakeOwner struct {
	id     int64
	name   string
	level  int32
	items  map[int32]int64
	awards []*config.IdCount
	state  *PlayerState
}

func (o *fakeOwner) GetId() int64    { return o.id }
func (o *fakeOwner) GetName() string { return o.name }
func (o *fakeOwner) GetLevel() int32 { return o.level }

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

func (o *fakeOwner) SaveAllianceState(state *PlayerState) error { o.state = state; return nil }

func loadConfig(t *testing.T) {
	t.Helper()
	excel.ReadAllJsonEntries("../../../config/json")

	// 创建公会需要雪花 id，测试里初始化一次
	utils.InitMachineID(1, 0, func() {})
}

func newTestManager(t *testing.T, st Store, id int64) (*Manager, *fakeOwner) {
	owner := &fakeOwner{id: id, name: "p", level: 30, items: map[int32]int64{}}
	m := NewManager(owner, st)
	m.LoadAll(nil)
	return m, owner
}

// TestCreateJoinLeave 创建、加入、退出与重名/重复加入校验
func TestCreateJoinLeave(t *testing.T) {
	loadConfig(t)

	cfg := config.GetAllianceConfigEntry()
	if cfg.CreateCost != nil {
		_ = cfg
	}

	st := newMemStore()

	m1, o1 := newTestManager(t, st, 1)
	if cfg.CreateCost != nil {
		o1.items[cfg.CreateCost.Id] = int64(cfg.CreateCost.Count) + 1
	}

	g, err := m1.Create("测试公会")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if g.LeaderId != 1 || g.Members[1] != PositionLeader {
		t.Fatalf("guild leader/member wrong: %+v", g)
	}

	// 重名
	m2, _ := newTestManager(t, st, 2)
	if _, err := m2.Create("测试公会"); !errors.Is(err, ErrAllianceNameExists) {
		t.Fatalf("expect name exists, got %v", err)
	}

	// 加入
	if err := m2.Join(g.Id); err != nil {
		t.Fatalf("Join failed: %v", err)
	}

	if len(g.Members) != 2 || m2.State.Position != PositionMember {
		t.Fatalf("join state wrong: %+v %+v", g.Members, m2.State)
	}

	// 重复加入
	if err := m2.Join(g.Id); !errors.Is(err, ErrAlreadyInAlliance) {
		t.Fatalf("expect already in, got %v", err)
	}

	// 会长不能在有其他成员时退出
	if err := m1.Leave(); !errors.Is(err, ErrNoPermission) {
		t.Fatalf("expect no permission, got %v", err)
	}

	// 成员退出
	if err := m2.Leave(); err != nil {
		t.Fatalf("member leave failed: %v", err)
	}

	if len(g.Members) != 1 {
		t.Fatalf("members after leave: %d", len(g.Members))
	}
}

// TestKick 权限校验
func TestKick(t *testing.T) {
	loadConfig(t)

	st := newMemStore()
	m1, o1 := newTestManager(t, st, 1)
	cfg := config.GetAllianceConfigEntry()
	if cfg.CreateCost != nil {
		o1.items[cfg.CreateCost.Id] = int64(cfg.CreateCost.Count) + 1
	}

	g, _ := m1.Create("踢人公会")

	m2, _ := newTestManager(t, st, 2)
	_ = m2.Join(g.Id)

	// 普通成员不能踢人
	if err := m2.Kick(1); !errors.Is(err, ErrNoPermission) {
		t.Fatalf("expect no permission, got %v", err)
	}

	// 会长可以踢
	if err := m1.Kick(2); err != nil {
		t.Fatalf("kick failed: %v", err)
	}

	if _, ok := g.Members[2]; ok {
		t.Fatal("member not removed")
	}
}

// TestBuild 建设消耗、限次与公会经验
func TestBuild(t *testing.T) {
	loadConfig(t)

	st := newMemStore()
	m, o := newTestManager(t, st, 1)

	cfg := config.GetAllianceConfigEntry()
	if cfg.CreateCost != nil {
		o.items[cfg.CreateCost.Id] = int64(cfg.CreateCost.Count) + 1
	}

	if _, err := m.Create("建设公会"); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// 建设消耗按公会等级取档（可能和创建消耗是同一种道具，创建之后再补足）
	if len(cfg.BuildCost) > 0 && cfg.BuildCost[0] != nil {
		o.items[cfg.BuildCost[0].Id] += int64(cfg.BuildCost[0].Count)*100 + 1000
	}

	now := time.Now()
	for i := 0; i < int(cfg.DailyMaxBuildTimes); i++ {
		if err := m.Build(now); err != nil {
			t.Fatalf("Build %d failed: %v", i, err)
		}
	}

	// 超过每日上限
	if err := m.Build(now); !errors.Is(err, ErrBuildTimesLimit) {
		t.Fatalf("expect build limit, got %v", err)
	}

	if len(o.awards) == 0 {
		t.Fatal("no build awards")
	}

	// 跨天重置
	m.CheckDay(now.AddDate(0, 0, 1))
	if m.State.BuildTimes != 0 {
		t.Fatalf("build times not reset: %d", m.State.BuildTimes)
	}
}
