// Package recharge 实现充值订单、VIP 与首充的抽象核心。
//
// kof 的 recharge_manager 有 2400 多行，强依赖公司支付平台（下单、回调验签、渠道对账）
// 以及月卡/特权卡/权益体系。这里只移植与具体支付渠道无关的部分：
//
//   - 充值订单：创建、幂等发货、状态记录
//   - VIP：按 VipLevel 表累计经验并升级
//   - 首充：按 FirstRechargeNew 表按天领取
//
// 支付回调由外部渠道调用 DeliverOrder，渠道验签与对账不在本包内；
// 月卡/特权卡/权益（VipLevelRights）需要单独的权益系统，暂未包含。
package recharge

import (
	"errors"
	"fmt"
	"time"

	"github.com/east-eden/server/excel/config"
	"github.com/google/uuid"
)

var (
	ErrOrderNotFound       = errors.New("recharge order not found")
	ErrOrderAlreadyDone    = errors.New("recharge order already delivered")
	ErrGoodsNotFound       = errors.New("recharge goods not found")
	ErrFirstRechargeLocked = errors.New("first recharge not unlocked")
	ErrFirstRechargeDone   = errors.New("first recharge day already claimed")
	ErrVipMaxLevel         = errors.New("vip reach max level")
)

// OrderStatus 订单状态
type OrderStatus int32

const (
	OrderPending   OrderStatus = 0 // 待支付
	OrderDelivered OrderStatus = 1 // 已发货
)

// Order 充值订单
type Order struct {
	Id          string      `bson:"id" json:"id"`
	GoodsId     int32       `bson:"goods_id" json:"goods_id"`
	Count       int32       `bson:"count" json:"count"`
	Status      OrderStatus `bson:"status" json:"status"`
	CreatedAt   int64       `bson:"created_at" json:"created_at"`
	DeliveredAt int64       `bson:"delivered_at" json:"delivered_at"`
}

// State 充值状态
type State struct {
	VipLevel  int32             `bson:"vip_level" json:"vip_level"`
	VipExp    int32             `bson:"vip_exp" json:"vip_exp"`
	Orders    map[string]*Order `bson:"orders" json:"orders"`
	FirstPaid bool              `bson:"first_paid" json:"first_paid"` // 是否已完成首充
	// FirstRechargeDays 首充活动已领取的天数（key 为活动 id）
	FirstRechargeDays map[int32]int32 `bson:"first_recharge_days" json:"first_recharge_days"`
}

// Owner 充值所需宿主能力
type Owner interface {
	GrantItems(items []*config.IdCount, reason string)
	GrantLootById(lootId int32, reason string)
	SaveRecharge(state *State) error
}

// Manager 充值管理器
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

	if m.State.Orders == nil {
		m.State.Orders = make(map[string]*Order, 8)
	}

	if m.State.FirstRechargeDays == nil {
		m.State.FirstRechargeDays = make(map[int32]int32, 4)
	}
}

// CreateOrder 创建充值订单，支付由外部渠道完成
func (m *Manager) CreateOrder(goodsId, count int32) (*Order, error) {
	entry, ok := config.GetShopGoodsEntry(goodsId)
	if !ok || len(entry.GainItems) == 0 && entry.LootId <= 0 {
		return nil, ErrGoodsNotFound
	}

	if count <= 0 {
		count = 1
	}

	order := &Order{
		Id:        uuid.NewString(),
		GoodsId:   goodsId,
		Count:     count,
		Status:    OrderPending,
		CreatedAt: time.Now().UnixMilli(),
	}

	m.State.Orders[order.Id] = order

	if err := m.owner.SaveRecharge(m.State); err != nil {
		return nil, err
	}

	return order, nil
}

// DeliverOrder 订单发货（支付回调入口），幂等
func (m *Manager) DeliverOrder(orderId string, vipExp int64) error {
	order, ok := m.State.Orders[orderId]
	if !ok {
		return ErrOrderNotFound
	}

	if order.Status == OrderDelivered {
		return ErrOrderAlreadyDone
	}

	entry, ok := config.GetShopGoodsEntry(order.GoodsId)
	if !ok {
		return ErrGoodsNotFound
	}

	// 发货
	if len(entry.GainItems) > 0 {
		m.owner.GrantItems(scale(entry.GainItems, order.Count), "recharge")
	}

	for i := int32(0); i < order.Count; i++ {
		if entry.LootId > 0 {
			m.owner.GrantLootById(entry.LootId, "recharge")
		}
	}

	order.Status = OrderDelivered
	order.DeliveredAt = time.Now().UnixMilli()
	m.State.FirstPaid = true

	if err := m.AddVipExp(vipExp); err != nil && !errors.Is(err, ErrVipMaxLevel) {
		return err
	}

	return m.owner.SaveRecharge(m.State)
}

// AddVipExp 增加 VIP 经验并按表升级
func (m *Manager) AddVipExp(exp int64) error {
	if exp <= 0 {
		return nil
	}

	cur, ok := config.GetVipLevelEntry(m.State.VipLevel)
	if !ok || cur.NextExp <= 0 {
		return ErrVipMaxLevel
	}

	m.State.VipExp += int32(exp)

	for {
		cur, ok = config.GetVipLevelEntry(m.State.VipLevel)
		if !ok || cur.NextExp <= 0 {
			m.State.VipExp = 0
			return ErrVipMaxLevel
		}

		if m.State.VipExp < cur.NextExp {
			break
		}

		if _, ok := config.GetVipLevelEntry(m.State.VipLevel + 1); !ok {
			m.State.VipExp = cur.NextExp
			return ErrVipMaxLevel
		}

		m.State.VipExp -= cur.NextExp
		m.State.VipLevel++
	}

	return nil
}

// ClaimFirstRecharge 领取首充奖励，day 必须按顺序递增
func (m *Manager) ClaimFirstRecharge(rechargeId int32, day int32) error {
	if !m.State.FirstPaid {
		return ErrFirstRechargeLocked
	}

	cur := m.State.FirstRechargeDays[rechargeId]
	if day <= cur {
		return ErrFirstRechargeDone
	}

	if day != cur+1 {
		return fmt.Errorf("%w: day %d, expect %d", ErrOrderNotFound, day, cur+1)
	}

	entry, ok := config.GetFirstRechargeNewEntry(fmt.Sprintf("%d", rechargeId), fmt.Sprintf("%d", day))
	if !ok {
		return ErrGoodsNotFound
	}

	m.owner.GrantItems(entry.Awards, "firstRecharge")
	m.State.FirstRechargeDays[rechargeId] = day

	return m.owner.SaveRecharge(m.State)
}

func scale(items []*config.IdCount, count int32) []*config.IdCount {
	out := make([]*config.IdCount, 0, len(items))
	for _, it := range items {
		if it == nil || it.Id <= 0 || it.Count <= 0 {
			continue
		}

		out = append(out, &config.IdCount{Id: it.Id, Count: it.Count * count})
	}

	return out
}
