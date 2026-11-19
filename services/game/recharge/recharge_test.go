package recharge

import (
	"errors"
	"testing"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
)

type fakeOwner struct {
	items map[int32]int64
	state *State
	loots []int32
}

func (o *fakeOwner) GrantItems(items []*config.IdCount, reason string) {
	for _, it := range items {
		o.items[it.Id] += int64(it.Count)
	}
}

func (o *fakeOwner) GrantLootById(lootId int32, reason string) { o.loots = append(o.loots, lootId) }

func (o *fakeOwner) SaveRecharge(state *State) error { o.state = state; return nil }

func loadConfig(t *testing.T) {
	t.Helper()
	excel.ReadAllJsonEntries("../../../config/json")
}

// TestVipExpLevelup 按 VipLevel 表累计经验升级
func TestVipExpLevelup(t *testing.T) {
	loadConfig(t)

	entry, ok := config.GetVipLevelEntry(0)
	if !ok || entry.NextExp <= 0 {
		t.Skip("no vip level 0 config")
	}

	owner := &fakeOwner{items: map[int32]int64{}}
	m := NewManager(owner)
	m.LoadAll(nil)

	if err := m.AddVipExp(int64(entry.NextExp)); err != nil {
		t.Fatalf("AddVipExp failed: %v", err)
	}

	if m.State.VipLevel != 1 {
		t.Fatalf("vip level = %d, want 1", m.State.VipLevel)
	}

	if m.State.VipExp != 0 {
		t.Fatalf("vip exp = %d, want 0", m.State.VipExp)
	}
}

// TestCreateAndDeliverOrder 下单、发货、幂等与 VIP 经验
func TestCreateAndDeliverOrder(t *testing.T) {
	loadConfig(t)

	// goods 1214: 获得 12500001 x10
	owner := &fakeOwner{items: map[int32]int64{}}
	m := NewManager(owner)
	m.LoadAll(nil)

	order, err := m.CreateOrder(1214, 1)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	if order.Status != OrderPending {
		t.Fatalf("status = %d", order.Status)
	}

	if err := m.DeliverOrder(order.Id, 60); err != nil {
		t.Fatalf("DeliverOrder failed: %v", err)
	}

	if owner.items[12500001] != 10 {
		t.Fatalf("goods not granted: %d", owner.items[12500001])
	}

	if !m.State.FirstPaid {
		t.Fatal("first paid flag not set")
	}

	// 重复发货
	if err := m.DeliverOrder(order.Id, 60); !errors.Is(err, ErrOrderAlreadyDone) {
		t.Fatalf("expect already done, got %v", err)
	}

	// 不存在的订单
	if err := m.DeliverOrder("not-exist", 0); !errors.Is(err, ErrOrderNotFound) {
		t.Fatalf("expect not found, got %v", err)
	}
}

// TestFirstRecharge 首充前不可领，领取按天递增且不可重复
func TestFirstRecharge(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{items: map[int32]int64{}}
	m := NewManager(owner)
	m.LoadAll(nil)

	// 找一条首充配置
	var rechargeId int32
	for _, e := range config.GetFirstRechargeNewRows() {
		rechargeId = e.Id
		break
	}

	if rechargeId == 0 {
		t.Skip("no first recharge config")
	}

	if err := m.ClaimFirstRecharge(rechargeId, 1); !errors.Is(err, ErrFirstRechargeLocked) {
		t.Fatalf("expect locked, got %v", err)
	}

	m.State.FirstPaid = true

	if err := m.ClaimFirstRecharge(rechargeId, 1); err != nil {
		t.Fatalf("claim day 1 failed: %v", err)
	}

	if got := m.State.FirstRechargeDays[rechargeId]; got != 1 {
		t.Fatalf("claimed day = %d, want 1", got)
	}

	if err := m.ClaimFirstRecharge(rechargeId, 1); !errors.Is(err, ErrFirstRechargeDone) {
		t.Fatalf("expect already done, got %v", err)
	}
}
