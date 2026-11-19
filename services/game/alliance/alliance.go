// Package alliance 实现公会（联盟）核心。
//
// kof 的 alliance_manager 有 4350 行，另有组队联赛、对战、科技等 2000 多行，覆盖
// 申请审批、职位、科技、boss、对战、旗帜/图标、二维码等完整社交体系。这里移植
// 公会本身的核心状态机：
//
//   - 创建（消耗配置道具）、加入、退出、踢人、会长转让
//   - 成员上限按 AllianceLevel 表随等级提升
//   - 每日建设：消耗配置道具、给公会加经验、发建设奖励、按日限次
//   - 公会升级：经验达到 AllianceLevel.nextNeedExp
//
// 申请审批、科技、boss、对战等不在本包内。
package alliance

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/utils"
)

var (
	ErrAllianceNotFound    = errors.New("alliance not found")
	ErrAllianceNameInvalid = errors.New("alliance name invalid")
	ErrAllianceNameExists  = errors.New("alliance name exists")
	ErrAllianceFull        = errors.New("alliance is full")
	ErrAlreadyInAlliance   = errors.New("already in an alliance")
	ErrNotInAlliance       = errors.New("not in an alliance")
	ErrNoPermission        = errors.New("no alliance permission")
	ErrBuildTimesLimit     = errors.New("alliance build times limit")
	ErrBuildCostNotEnough  = errors.New("alliance build cost not enough")
	ErrExitCooldown        = errors.New("alliance exit cooldown")
)

// 职位
const (
	PositionMember int32 = 0 // 成员
	PositionElite  int32 = 1 // 精英
	PositionDeputy int32 = 2 // 副会长
	PositionLeader int32 = 3 // 会长
)

// Guild 公会
type Guild struct {
	Id        int64           `bson:"_id" json:"_id"`
	Name      string          `bson:"name" json:"name"`
	LeaderId  int64           `bson:"leader_id" json:"leader_id"`
	Members   map[int64]int32 `bson:"members" json:"members"` // memberId -> position
	Level     int32           `bson:"level" json:"level"`
	Exp       int32           `bson:"exp" json:"exp"`
	Notice    string          `bson:"notice" json:"notice"`
	CreatedAt int64           `bson:"created_at" json:"created_at"`
}

// Store 公会持久化接口（Mongo 实现与内存实现都满足）
type Store interface {
	Load(ctx context.Context, id int64) (*Guild, error)
	Save(ctx context.Context, g *Guild) error
	Delete(ctx context.Context, id int64) error
	FindByName(ctx context.Context, name string) (int64, bool, error)
	SaveName(ctx context.Context, name string, id int64) error
}

// PlayerState 玩家侧公会状态
type PlayerState struct {
	AllianceId int64  `bson:"alliance_id" json:"alliance_id"`
	Position   int32  `bson:"position" json:"position"`
	BuildTimes int32  `bson:"build_times" json:"build_times"` // 今日已建设次数
	BuildDay   string `bson:"build_day" json:"build_day"`
	ExitAt     int64  `bson:"exit_at" json:"exit_at"` // 退出冷却结束时间(毫秒)
}

// Owner 公会所需宿主能力
type Owner interface {
	GetId() int64
	GetName() string
	GetLevel() int32
	CostItems(costs []*config.IdCount, reason string) error
	GrantAwardList(awards []*config.IdCount, reason string)
	SaveAllianceState(state *PlayerState) error
}

// Manager 公会管理器（玩家侧入口，共享 Store）
type Manager struct {
	owner Owner
	store Store
	State *PlayerState
}

func NewManager(owner Owner, store Store) *Manager {
	return &Manager{owner: owner, store: store, State: &PlayerState{}}
}

// LoadAll 接上读档数据
func (m *Manager) LoadAll(state *PlayerState) {
	if state != nil {
		m.State = state
	}

	if m.State == nil {
		m.State = &PlayerState{}
	}
}

func allianceConfig() *config.AllianceConfigEntry {
	return config.GetAllianceConfigEntry()
}

// CheckDay 跨天重置建设次数
func (m *Manager) CheckDay(now time.Time) {
	day := now.Format(time.DateOnly)
	if m.State.BuildDay == day {
		return
	}

	m.State.BuildDay = day
	m.State.BuildTimes = 0
}

// Create 创建公会
func (m *Manager) Create(name string) (*Guild, error) {
	cfg := allianceConfig()
	if cfg != nil && cfg.AllianceNameMaxLimit > 0 && len([]rune(name)) > int(cfg.AllianceNameMaxLimit) {
		return nil, ErrAllianceNameInvalid
	}

	if name == "" {
		return nil, ErrAllianceNameInvalid
	}

	if m.State.AllianceId > 0 {
		return nil, ErrAlreadyInAlliance
	}

	ctx := context.Background()

	if _, exists, err := m.store.FindByName(ctx, name); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrAllianceNameExists
	}

	if cfg != nil && cfg.CreateCost != nil {
		if err := m.owner.CostItems([]*config.IdCount{cfg.CreateCost}, "allianceCreate"); err != nil {
			return nil, err
		}
	}

	id, err := utils.NextID(define.SnowFlake_Alliance)
	if err != nil {
		return nil, err
	}

	g := &Guild{
		Id:        id,
		Name:      name,
		LeaderId:  m.owner.GetId(),
		Members:   map[int64]int32{m.owner.GetId(): PositionLeader},
		Level:     1,
		CreatedAt: time.Now().UnixMilli(),
	}

	if err := m.store.Save(ctx, g); err != nil {
		return nil, err
	}

	if err := m.store.SaveName(ctx, name, id); err != nil {
		return nil, err
	}

	if err := m.setState(g.Id, PositionLeader); err != nil {
		return nil, err
	}

	return g, nil
}

// Join 加入公会
func (m *Manager) Join(allianceId int64) error {
	if m.State.AllianceId > 0 {
		return ErrAlreadyInAlliance
	}

	if m.State.ExitAt > time.Now().UnixMilli() {
		return ErrExitCooldown
	}

	ctx := context.Background()
	g, err := m.store.Load(ctx, allianceId)
	if err != nil {
		return err
	}

	if g == nil {
		return ErrAllianceNotFound
	}

	if len(g.Members) >= maxMembers(g.Level) {
		return ErrAllianceFull
	}

	g.Members[m.owner.GetId()] = PositionMember
	if err := m.store.Save(ctx, g); err != nil {
		return err
	}

	return m.setState(g.Id, PositionMember)
}

// Leave 退出公会
func (m *Manager) Leave() error {
	if m.State.AllianceId == 0 {
		return ErrNotInAlliance
	}

	ctx := context.Background()
	g, err := m.store.Load(ctx, m.State.AllianceId)
	if err != nil {
		return err
	}

	if g == nil {
		return ErrAllianceNotFound
	}

	memberId := m.owner.GetId()

	// 会长只有在自己是最后一人时才能退出
	if g.LeaderId == memberId && len(g.Members) > 1 {
		return ErrNoPermission
	}

	delete(g.Members, memberId)

	if len(g.Members) == 0 {
		if err := m.store.Delete(ctx, g.Id); err != nil {
			return err
		}
	} else if err := m.store.Save(ctx, g); err != nil {
		return err
	}

	cooldown := int64(0)
	if cfg := allianceConfig(); cfg != nil && cfg.ExitCD > 0 {
		cooldown = time.Now().Add(time.Duration(cfg.ExitCD) * time.Hour).UnixMilli()
	}

	return m.setState(0, PositionMember, cooldown)
}

// Kick 踢出成员（会长/副会长）
func (m *Manager) Kick(memberId int64) error {
	if m.State.Position < PositionDeputy {
		return ErrNoPermission
	}

	ctx := context.Background()
	g, err := m.store.Load(ctx, m.State.AllianceId)
	if err != nil {
		return err
	}

	if g == nil {
		return ErrAllianceNotFound
	}

	if _, ok := g.Members[memberId]; !ok {
		return ErrAllianceNotFound
	}

	if memberId == g.LeaderId {
		return ErrNoPermission
	}

	delete(g.Members, memberId)
	return m.store.Save(ctx, g)
}

// Transfer 转让会长
func (m *Manager) Transfer(memberId int64) error {
	if m.State.Position != PositionLeader {
		return ErrNoPermission
	}

	ctx := context.Background()
	g, err := m.store.Load(ctx, m.State.AllianceId)
	if err != nil {
		return err
	}

	if g == nil {
		return ErrAllianceNotFound
	}

	if _, ok := g.Members[memberId]; !ok {
		return ErrAllianceNotFound
	}

	g.LeaderId = memberId
	g.Members[memberId] = PositionLeader
	g.Members[m.owner.GetId()] = PositionMember

	return m.store.Save(ctx, g)
}

// Build 每日建设：消耗配置道具、发奖励、给公会加经验并可能升级
func (m *Manager) Build(now time.Time) error {
	if m.State.AllianceId == 0 {
		return ErrNotInAlliance
	}

	m.CheckDay(now)

	cfg := allianceConfig()
	if cfg == nil {
		return ErrAllianceNotFound
	}

	if cfg.DailyMaxBuildTimes > 0 && m.State.BuildTimes >= cfg.DailyMaxBuildTimes {
		return ErrBuildTimesLimit
	}

	ctx := context.Background()
	g, err := m.store.Load(ctx, m.State.AllianceId)
	if err != nil {
		return err
	}

	if g == nil {
		return ErrAllianceNotFound
	}

	buildCost := buildCostForLevel(cfg, g.Level)
	if buildCost != nil {
		if err := m.owner.CostItems([]*config.IdCount{buildCost}, "allianceBuild"); err != nil {
			return ErrBuildCostNotEnough
		}
	}

	if len(cfg.BuildAwards) > 0 {
		m.owner.GrantAwardList(cfg.BuildAwards, "allianceBuild")
	}

	// 公会经验与升级
	if entry, ok := config.GetAllianceLevelEntry(g.Level); ok && entry.NextNeedExp > 0 {
		g.Exp += 10
		for g.Exp >= entry.NextNeedExp {
			g.Exp -= entry.NextNeedExp

			next, ok := config.GetAllianceLevelEntry(g.Level + 1)
			if !ok {
				g.Exp = 0
				break
			}

			g.Level++
			entry = next
		}
	}

	if err := m.store.Save(ctx, g); err != nil {
		return err
	}

	m.State.BuildTimes++

	return m.owner.SaveAllianceState(m.State)
}

// Info 返回当前公会信息
func (m *Manager) Info() (*Guild, error) {
	if m.State.AllianceId == 0 {
		return nil, ErrNotInAlliance
	}

	return m.store.Load(context.Background(), m.State.AllianceId)
}

// MemberList 按职位倒序返回成员 id
func (m *Manager) MemberList(g *Guild) []int64 {
	ids := make([]int64, 0, len(g.Members))
	for id := range g.Members {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool { return g.Members[ids[i]] > g.Members[ids[j]] })
	return ids
}

func (m *Manager) setState(allianceId int64, position int32, exitAt ...int64) error {
	m.State.AllianceId = allianceId
	m.State.Position = position
	if len(exitAt) > 0 {
		m.State.ExitAt = exitAt[0]
	}

	return m.owner.SaveAllianceState(m.State)
}

func maxMembers(level int32) int {
	if entry, ok := config.GetAllianceLevelEntry(level); ok && entry.MaxMemberCount > 0 {
		return int(entry.MaxMemberCount)
	}

	return 15
}

// buildCostForLevel 按公会等级取建设消耗，等级超出配置时取最后一项
func buildCostForLevel(cfg *config.AllianceConfigEntry, level int32) *config.IdCount {
	if len(cfg.BuildCost) == 0 {
		return nil
	}

	idx := int(level) - 1
	if idx < 0 {
		idx = 0
	}

	if idx >= len(cfg.BuildCost) {
		idx = len(cfg.BuildCost) - 1
	}

	return cfg.BuildCost[idx]
}
