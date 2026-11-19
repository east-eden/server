// Package arena 实现竞技场核心：积分（ELO）、挑战消耗、日/周排名奖励。
//
// kof 的 arena_manager 有 1880 行，依赖跨服排名服务、机器人池、战斗校验与结算窗口。
// 这里移植与具体匹配实现无关的部分：
//
//   - ELO 积分公式（与 kof 一致）
//   - 挑战/刷新消耗与每日次数
//   - 日榜/周榜奖励按排名区间发放
//
// 对手列表由排名服务提供，战斗校验走 BattleManager，均不在本包内。
package arena

import (
	"errors"
	"math"
	"time"

	"github.com/east-eden/server/excel/config"
)

var (
	ErrChallengeCostNotEnough = errors.New("arena challenge cost not enough")
	ErrDailyRewardClaimed     = errors.New("arena daily reward already claimed")
	ErrWeeklyRewardClaimed    = errors.New("arena weekly reward already claimed")
	ErrRewardNotFound         = errors.New("arena reward not found")
)

// State 竞技场状态
type State struct {
	Score         int32  `bson:"score" json:"score"`
	Rank          int32  `bson:"rank" json:"rank"`
	WeekWins      int32  `bson:"week_wins" json:"week_wins"`
	DailyClaimed  bool   `bson:"daily_claimed" json:"daily_claimed"`
	WeeklyClaimed bool   `bson:"weekly_claimed" json:"weekly_claimed"`
	Day           string `bson:"day" json:"day"`
	Week          string `bson:"week" json:"week"`
}

// Owner 竞技场所需宿主能力
type Owner interface {
	GetItemCount(entryId int32) int64
	CostItems(costs []*config.IdCount, reason string) error
	GrantAwardList(awards []*config.IdCount, reason string)
	SaveArena(state *State) error
}

// Manager 竞技场管理器
type Manager struct {
	owner Owner
	State *State
}

func NewManager(owner Owner) *Manager {
	return &Manager{owner: owner, State: &State{}}
}

// LoadAll 接上读档数据
func (m *Manager) LoadAll(state *State) {
	if state != nil {
		m.State = state
	}

	if m.State == nil {
		m.State = &State{}
	}
}

// InitScore 初始积分取配置的默认值
func (m *Manager) InitScore() {
	if m.State.Score > 0 {
		return
	}

	if cfg := config.GetArenaConfigEntry(); cfg != nil && cfg.DefaultScore > 0 {
		m.State.Score = cfg.DefaultScore
	} else {
		m.State.Score = 1000
	}
}

// CalcNewScore ELO 积分计算，与 kof 的实现一致
func CalcNewScore(score, enemyScore int32, win bool) int32 {
	var k float32

	switch {
	case score <= 1200 && enemyScore <= 1200:
		if win {
			k = 30
		} else {
			k = 10
		}
	case score <= 1500:
		k = 30
	case score <= 2400:
		k = 63.3 - float32(score)/45
	default:
		k = 10
	}

	g := float32(1)
	if !win {
		g = 0
	}

	ge := float32(1 / (1 + math.Pow(10, float64(enemyScore-score)/400)))
	rn := float32(score) + k*(g-ge)
	if rn < 1 {
		rn = 1
	}

	return int32(math.Round(float64(rn)))
}

// Challenge 校验并扣除挑战消耗，返回是否成功
func (m *Manager) Challenge() error {
	cfg := config.GetArenaConfigEntry()
	if cfg == nil || len(cfg.ChallengeCost) < 2 {
		return nil
	}

	// KV 配置形如 "10900000,1"，生成器按 list:int 展开成 [id, count, id, count...]
	costs := make([]*config.IdCount, 0, len(cfg.ChallengeCost)/2)
	for i := 0; i+1 < len(cfg.ChallengeCost); i += 2 {
		costs = append(costs, &config.IdCount{Id: cfg.ChallengeCost[i], Count: cfg.ChallengeCost[i+1]})
	}

	return m.owner.CostItems(costs, "arenaChallenge")
}

// ApplyResult 应用战斗结果，返回积分变化
func (m *Manager) ApplyResult(enemyScore int32, win bool) int32 {
	old := m.State.Score
	m.State.Score = CalcNewScore(old, enemyScore, win)

	if win {
		m.State.WeekWins++
	}

	_ = m.owner.SaveArena(m.State)

	return m.State.Score - old
}

// ResetIfNeeded 按天/周重置领取状态
func (m *Manager) ResetIfNeeded(now time.Time) {
	day := now.Format(time.DateOnly)
	if m.State.Day != day {
		m.State.Day = day
		m.State.DailyClaimed = false
	}

	year, week := now.ISOWeek()
	weekKey := time.Date(year, 1, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, (week-1)*7).Format(time.DateOnly)
	if m.State.Week != weekKey {
		m.State.Week = weekKey
		m.State.WeeklyClaimed = false
		m.State.WeekWins = 0
	}
}

// dailyAward 按排名取日榜奖励
func dailyAward(rank int32) ([]*config.IdCount, bool) {
	var (
		matched *config.ArenaDailyRewardEntry
	)

	for _, e := range config.GetArenaDailyRewardRows() {
		if rank >= e.RankRange && (matched == nil || e.RankRange > matched.RankRange) {
			matched = e
		}
	}

	if matched == nil {
		return nil, false
	}

	awards := make([]*config.IdCount, 0, len(matched.RewardItemIds))
	for i, id := range matched.RewardItemIds {
		if i >= len(matched.RewardCounts) {
			break
		}

		awards = append(awards, &config.IdCount{Id: id, Count: matched.RewardCounts[i]})
	}

	return awards, true
}

// ClaimDaily 领取日榜奖励
func (m *Manager) ClaimDaily(rank int32) error {
	if m.State.DailyClaimed {
		return ErrDailyRewardClaimed
	}

	awards, ok := dailyAward(rank)
	if !ok {
		return ErrRewardNotFound
	}

	m.owner.GrantAwardList(awards, "arenaDaily")
	m.State.DailyClaimed = true

	return m.owner.SaveArena(m.State)
}

// ClaimWeekly 领取周榜奖励
func (m *Manager) ClaimWeekly(rank int32) error {
	if m.State.WeeklyClaimed {
		return ErrWeeklyRewardClaimed
	}

	var matched *config.ArenaWeeklyRewardEntry
	for _, e := range config.GetArenaWeeklyRewardRows() {
		if rank >= e.RankRange && (matched == nil || e.RankRange > matched.RankRange) {
			matched = e
		}
	}

	if matched == nil {
		return ErrRewardNotFound
	}

	awards := make([]*config.IdCount, 0, len(matched.RewardItemIds))
	for i, id := range matched.RewardItemIds {
		if i >= len(matched.RewardCounts) {
			break
		}

		awards = append(awards, &config.IdCount{Id: id, Count: matched.RewardCounts[i]})
	}

	m.owner.GrantAwardList(awards, "arenaWeekly")
	m.State.WeeklyClaimed = true

	return m.owner.SaveArena(m.State)
}
