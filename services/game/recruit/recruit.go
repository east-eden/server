// Package recruit 实现配置驱动的招募/抽卡。
//
// kof 的 recruit_manager 有 2400 多行，包含免费/广告/高级/超级/基金/心愿/梦境等
// 多套卡池，以及十连特权卡、梦境英雄概率、招募等级等。这里移植其中通用且可配置的部分：
//
//   - 免费招募：按 Config.maxFreeRecruit 每日限次
//   - 高级招募：单抽/十连，消耗高级招募券，券不足时按商城道具价折算钻石
//   - 保底：按 Config 里的大保底/小保底次数与掉落列表触发
//   - 首次/首次十连使用专门的掉落
//
// 广告招募、超级招募、心愿/梦境卡池、招募等级与十连特权卡未包含。
package recruit

import (
	"errors"
	"fmt"
	"time"

	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/services/game/loot"
	"github.com/east-eden/server/utils"
)

var (
	ErrRecruitTimesInvalid  = errors.New("recruit times invalid")
	ErrRecruitFreeLimit     = errors.New("free recruit times reached")
	ErrRecruitCostNotEnough = errors.New("recruit cost not enough")
	ErrRecruitLootMissing   = errors.New("recruit loot not configured")
)

// State 招募状态
type State struct {
	Day            string `bson:"day" json:"day"`                         // 免费次数所属日期
	FreeUsed       int32  `bson:"free_used" json:"free_used"`             // 当日已用免费次数
	SeniorTimes    int32  `bson:"senior_times" json:"senior_times"`       // 高级招募累计次数
	SmallGuarantee int32  `bson:"small_guarantee" json:"small_guarantee"` // 距离小保底的次数
	BigGuarantee   int32  `bson:"big_guarantee" json:"big_guarantee"`     // 距离大保底的次数
	FirstFree      bool   `bson:"first_free" json:"first_free"`           // 是否已用过首次免费
	FirstSenior    bool   `bson:"first_senior" json:"first_senior"`       // 是否已用过首次高级
	FirstTen       bool   `bson:"first_ten" json:"first_ten"`             // 是否已用过首次十连
}

// RecruitType 招募类型
const (
	TypeFree   int32 = 0 // 免费招募
	TypeSenior int32 = 1 // 高级招募
)

// Reward 抽卡产出
type Reward struct {
	ItemId int32 `json:"item_id"`
	Num    int32 `json:"num"`
}

// Owner 招募所需宿主能力
type Owner interface {
	GetItemCount(entryId int32) int64
	CostItems(costs []*config.IdCount, reason string) error
	// GrantLootAndReport 按掉落表发放并返回实际产出
	GrantLootAndReport(lootId int32, reason string) []Reward
	SaveRecruit(state *State) error
}

// Manager 招募管理器
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

// CheckDay 跨天重置免费次数
func (m *Manager) CheckDay(now time.Time) {
	day := now.Format(time.DateOnly)
	if m.State.Day == day {
		return
	}

	m.State.Day = day
	m.State.FreeUsed = 0
}

// Free 免费招募
func (m *Manager) Free(now time.Time) ([]Reward, error) {
	m.CheckDay(now)

	cfg := config.GetConfigEntry()
	if cfg == nil {
		return nil, ErrRecruitLootMissing
	}

	if cfg.MaxFreeRecruit > 0 && m.State.FreeUsed >= cfg.MaxFreeRecruit {
		return nil, ErrRecruitFreeLimit
	}

	lootId := cfg.RecruitADLootId
	if !m.State.FirstFree {
		lootId = cfg.FirstADRecruitLoot
		m.State.FirstFree = true
	}

	if lootId <= 0 {
		return nil, ErrRecruitLootMissing
	}

	rewards := m.owner.GrantLootAndReport(lootId, "recruitFree")
	m.State.FreeUsed++

	if err := m.owner.SaveRecruit(m.State); err != nil {
		return rewards, err
	}

	return rewards, nil
}

// seniorCost 计算本次高级招募的消耗：优先招募券，不足时按商城道具价折算
func (m *Manager) seniorCost(times int32) ([]*config.IdCount, error) {
	cfg := config.GetConfigEntry()
	if cfg == nil || cfg.SeniorRecruitCost == nil {
		return nil, ErrRecruitCostNotEnough
	}

	ticket := cfg.SeniorRecruitCost
	need := int64(ticket.Count) * int64(times)

	if m.owner.GetItemCount(ticket.Id) >= need {
		return []*config.IdCount{{Id: ticket.Id, Count: int32(need)}}, nil
	}

	// 券不足：按高级招募券商品的售价折算，十连按配置折扣
	goods, ok := config.GetShopGoodsEntry(int32(config.PackageId_RecruitTicket))
	if !ok || len(goods.CostItems) == 0 {
		return nil, ErrRecruitCostNotEnough
	}

	costs := make([]*config.IdCount, 0, len(goods.CostItems))
	for _, c := range goods.CostItems {
		count := float32(c.Count) * float32(times)
		if times == 10 && cfg.TenTimesDiscount > 0 {
			count *= cfg.TenTimesDiscount
		}

		costs = append(costs, &config.IdCount{Id: c.Id, Count: int32(count)})
	}

	return costs, nil
}

// seniorPoolLoot 选择本次抽取的掉落：先判保底，再判首次。
// 返回的 isBig/isSmall 表示该次是否由保底触发，调用方据此重置保底计数。
func (m *Manager) seniorPoolLoot(times int32) (lootId int32, isBig, isSmall bool) {
	cfg := config.GetConfigEntry()

	// 大保底
	if len(cfg.SeniorRecruitGuaranteeLootList) > 0 && cfg.SeniorRecruitGuaranteeTimes > 0 &&
		m.State.BigGuarantee+1 >= cfg.SeniorRecruitGuaranteeTimes {
		return cfg.SeniorRecruitGuaranteeLootList[0], true, false
	}

	// 小保底
	if len(cfg.SeniorRecruitSmallGuaranteeList) > 0 && cfg.SeniorRecruitSmallGuaranteeTimes > 0 &&
		m.State.SmallGuarantee+1 >= cfg.SeniorRecruitSmallGuaranteeTimes {
		return cfg.SeniorRecruitSmallGuaranteeList[0], false, true
	}

	// 首次十连 / 首次高级
	if times == 10 && !m.State.FirstTen && cfg.FirstTenRecruitLoot > 0 {
		m.State.FirstTen = true
		return cfg.FirstTenRecruitLoot, false, false
	}

	if !m.State.FirstSenior && cfg.FirstSeniorRecruitLoot > 0 {
		m.State.FirstSenior = true
		return cfg.FirstSeniorRecruitLoot, false, false
	}

	if len(cfg.SeniorRecruitLootIdList) > 0 {
		return cfg.SeniorRecruitLootIdList[0], false, false
	}

	return 0, false, false
}

// Senior 高级招募，times 只支持 1 或 10
func (m *Manager) Senior(times int32) ([]Reward, error) {
	if times != 1 && times != 10 {
		return nil, ErrRecruitTimesInvalid
	}

	costs, err := m.seniorCost(times)
	if err != nil {
		return nil, err
	}

	if err := m.owner.CostItems(costs, "recruitSenior"); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecruitCostNotEnough, err)
	}

	rewards := make([]Reward, 0, times*2)

	for i := int32(0); i < times; i++ {
		lootId, isBig, isSmall := m.seniorPoolLoot(times)
		if lootId <= 0 {
			return rewards, ErrRecruitLootMissing
		}

		rewards = append(rewards, m.owner.GrantLootAndReport(lootId, "recruitSenior")...)

		m.State.SeniorTimes++
		m.State.BigGuarantee++
		m.State.SmallGuarantee++

		// 保底触发后计数清零（保底不属于“未出货”的次数）
		if isBig {
			m.State.BigGuarantee = 0
		}

		if isSmall {
			m.State.SmallGuarantee = 0
		}
	}

	if err := m.owner.SaveRecruit(m.State); err != nil {
		return rewards, err
	}

	return rewards, nil
}

// CalcLoot 按掉落 id 展开结果，供宿主实现 GrantLootById
func CalcLoot(lootId int32) []Reward {
	itemIds, counts, err := loot.CalcLootResult(lootId)
	if err != nil {
		utils.ErrPrint(err, "CalcLootResult failed", lootId)
		return nil
	}

	rewards := make([]Reward, 0, len(itemIds))
	for i, id := range itemIds {
		if i >= len(counts) || counts[i] <= 0 {
			continue
		}

		rewards = append(rewards, Reward{ItemId: id, Num: counts[i]})
	}

	return rewards
}
