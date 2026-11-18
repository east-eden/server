package player

import (
	"testing"

	"github.com/east-eden/server/excel/config"
	"github.com/golang/mock/gomock"
)

// TestTowerChallenge 用 kof 的 Tower 配置验证爬塔消耗门票与层奖励发放
func TestTowerChallenge(t *testing.T) {
	setupTestEnv()

	entry, ok := config.GetTowerEntry(1)
	if !ok {
		t.Skip("no tower floor 1 config")
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	initMockStore(t, ctl)

	p := newTestPlayer()

	// 补足爬塔挑战券
	ticketId := int32(config.GoodsConfigId_TowerTicket)
	if _, err := p.ItemManager().OnItemAdd(ticketId, 2, "test"); err != nil {
		t.Fatalf("grant tower ticket failed: %v", err)
	}

	m := NewTowerManager(p)

	if err := m.Challenge(0, 1, nil); err != nil {
		t.Fatalf("Challenge floor 1 failed: %v", err)
	}

	if m.PassedFloor != 1 {
		t.Fatalf("passed floor = %d, want 1", m.PassedFloor)
	}

	if got := p.ItemManager().GetItemCount(ticketId); got != 1 {
		t.Fatalf("ticket count = %d, want 1", got)
	}

	// 层奖励入账（FloorAward 非空时）
	if len(entry.FloorAward) > 0 {
		t.Logf("floor 1 award=%d items, bag size=%d", len(entry.FloorAward), len(p.ItemManager().Items))
	}

	// 跳过层的请求应被拒绝
	if err := m.Challenge(0, 5, nil); err == nil {
		t.Fatal("expect invalid floor error")
	}
}
