package player

import (
	"testing"

	"github.com/east-eden/server/excel/config"
	"github.com/golang/mock/gomock"
)

// TestFirstLogonConfig 首次登录按 kof 的 Config.initGold / initItems 发放
func TestFirstLogonConfig(t *testing.T) {
	setupTestEnv()

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	initMockStore(t, ctl)

	p := newTestPlayer()
	p.OnFirstLogon()

	cfg := config.GetConfigEntry()
	if cfg == nil {
		t.Fatal("config entry not loaded")
	}

	// 初始金币
	if cfg.InitGold > 0 {
		got := p.ItemManager().GetItemCount(int32(config.GoodsConfigId_Gold))
		if got != int64(cfg.InitGold) {
			t.Fatalf("init gold = %d, want %d", got, cfg.InitGold)
		}
	}

	// 初始道具里的英雄卡应被路由成英雄
	var heroGranted bool
	for _, item := range cfg.InitItems {
		entry, ok := config.GetGoodsEntry(item.Id)
		if !ok || entry.Type != config.GoodsType_Hero || entry.SubType != config.GoodsSubType_HeroCard {
			continue
		}

		if p.HeroManager().GetHeroByTypeId(item.Id) != nil {
			heroGranted = true
		}
	}

	if !heroGranted {
		t.Fatalf("no starter hero granted from initItems %+v", cfg.InitItems)
	}

	// 英雄卡不应残留在物品账本里
	for _, item := range cfg.InitItems {
		entry, ok := config.GetGoodsEntry(item.Id)
		if ok && entry.Type == config.GoodsType_Hero {
			if got := p.ItemManager().GetItemCount(item.Id); got != 0 {
				t.Fatalf("hero goods %d should not be in item ledger, got %d", item.Id, got)
			}
		}
	}
}

// TestHeroChipModel 英雄碎片走 chip 计数，升星从 chip 计数扣
func TestHeroChipModel(t *testing.T) {
	setupTestEnv()

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	initMockStore(t, ctl)

	p := newTestPlayer()

	card, err := p.HeroManager().AddHeroByTypeId(2090101)
	if err != nil {
		t.Fatalf("AddHeroByTypeId failed: %v", err)
	}

	entry, _ := card.Entry()
	chipId := entry.ChipId

	// 碎片通过物品账本入口发放，应进 chip 计数而不是物品背包
	if _, err := p.ItemManager().OnItemAdd(chipId, 100, "test"); err != nil {
		t.Fatalf("add chip failed: %v", err)
	}

	if got := p.HeroManager().GetHeroChipCount(chipId); got != 100 {
		t.Fatalf("chip count = %d, want 100", got)
	}

	if got := p.ItemManager().GetItemCount(chipId); got != 0 {
		t.Fatalf("chip should not be in item ledger, got %d", got)
	}

	// 升星消耗碎片
	starEntry, ok := config.GetHeroStarEntry("2", entry.Rarity.String())
	if !ok {
		t.Skip("no hero star entry for star 2")
	}

	if err := p.HeroManager().HeroStarup(int64(card.ConfigId)); err != nil {
		t.Fatalf("HeroStarup failed: %v", err)
	}

	if card.Star != 2 {
		t.Fatalf("star = %d, want 2", card.Star)
	}

	want := 100 - starEntry.Cost
	if got := p.HeroManager().GetHeroChipCount(chipId); got != want {
		t.Fatalf("chip after starup = %d, want %d", got, want)
	}
}
