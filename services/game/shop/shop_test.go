package shop

import (
	"errors"
	"testing"
	"time"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
	"google.golang.org/protobuf/proto"
)

type fakeOwner struct {
	level int32
	items map[int32]int64
	loots []int32
	saved map[int32]*ShopState
	sent  int
}

func (o *fakeOwner) GetLevel() int32 { return o.level }

func (o *fakeOwner) GetItemCount(entryId int32) int64 { return o.items[entryId] }

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

func (o *fakeOwner) GrantItems(items []*config.IdCount, reason string) {
	for _, it := range items {
		o.items[it.Id] += int64(it.Count)
	}
}

func (o *fakeOwner) GrantLootById(lootId int32, reason string) { o.loots = append(o.loots, lootId) }

func (o *fakeOwner) SaveShops(data map[int32]*ShopState) error {
	o.saved = data
	return nil
}

func (o *fakeOwner) SendProtoMessage(msg proto.Message) { o.sent++ }

func loadConfig(t *testing.T) {
	t.Helper()
	excel.ReadAllJsonEntries("../../../config/json")
}

// TestShopBuy 用真实配置验证购买扣费、发货、限购
func TestShopBuy(t *testing.T) {
	loadConfig(t)

	// goods 1214: 消耗 10100023 x15，获得 12500001 x10，限购 1
	owner := &fakeOwner{level: 30, items: map[int32]int64{10100023: 30}}
	m := NewManager(owner)
	m.LoadAll(make(map[int32]*ShopState))

	if err := m.Buy(31, 1214, 1); err != nil {
		t.Fatalf("Buy failed: %v", err)
	}

	if owner.items[10100023] != 15 {
		t.Fatalf("cost not applied: %d", owner.items[10100023])
	}

	if owner.items[12500001] != 10 {
		t.Fatalf("gain not applied: %d", owner.items[12500001])
	}

	if got := m.GoodsRemain(31, 1214); got != 0 {
		t.Fatalf("remain = %d, want 0", got)
	}

	// 超过限购
	if err := m.Buy(31, 1214, 1); !errors.Is(err, ErrGoodsSoldOut) {
		t.Fatalf("expect sold out, got %v", err)
	}

	if m.Shops[31].Goods[1214].Bought != 1 {
		t.Fatalf("bought = %d", m.Shops[31].Goods[1214].Bought)
	}
}

// TestShopBuyNotEnough 货币不足时不扣费也不计数
func TestShopBuyNotEnough(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{level: 30, items: map[int32]int64{10100023: 5}}
	m := NewManager(owner)

	if err := m.Buy(31, 1214, 1); !errors.Is(err, ErrGoodsCostNotEnough) {
		t.Fatalf("expect cost not enough, got %v", err)
	}

	if owner.items[10100023] != 5 {
		t.Fatalf("item should not be touched: %d", owner.items[10100023])
	}

	if got := m.GoodsRemain(31, 1214); got != 1 {
		t.Fatalf("remain = %d, want 1", got)
	}
}

// TestShopReset 商店级 cron 重置会把已购数量清零
func TestShopReset(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{level: 30, items: map[int32]int64{10100023: 30}}
	m := NewManager(owner)
	m.LoadAll(make(map[int32]*ShopState))

	if err := m.Buy(31, 1214, 1); err != nil {
		t.Fatalf("Buy failed: %v", err)
	}

	// 上次重置时间置 0，下一次 CheckAll 必然跨过 cron 点
	m.Shops[31].ResetAt = 0
	m.CheckAll(time.Now())

	if got := m.Shops[31].Goods[1214].Bought; got != 0 {
		t.Fatalf("bought after reset = %d, want 0", got)
	}

	if got := m.GoodsRemain(31, 1214); got != 1 {
		t.Fatalf("remain after reset = %d, want 1", got)
	}
}

// TestShopGenPB 等级不足的商店不下发
func TestShopGenPB(t *testing.T) {
	loadConfig(t)

	owner := &fakeOwner{level: 1, items: map[int32]int64{}}
	m := NewManager(owner)
	m.LoadAll(make(map[int32]*ShopState))

	for _, info := range m.GenPB() {
		if entry, ok := config.GetShopEntry(info.Id); ok && entry.OpenLv > owner.level {
			t.Fatalf("shop %d should be hidden at level %d", info.Id, owner.level)
		}
	}
}
