package hero

import (
	"testing"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
)

func loadConfig(t *testing.T) {
	t.Helper()
	excel.ReadAllJsonEntries("../../../config/json")
}

// TestCardBaseAttrs 用真实配置验证卡牌基础属性计算
func TestCardBaseAttrs(t *testing.T) {
	loadConfig(t)

	card, err := NewCard(2090101)
	if err != nil {
		t.Fatalf("NewCard failed: %v", err)
	}

	propertyId, ok := card.BasePropertyId()
	if !ok {
		t.Fatal("BasePropertyId failed")
	}

	t.Logf("card=%d rarity=%d propertyId=%d", card.ConfigId, card.Rarity, propertyId)

	if err := card.CalcBaseAttrs(1); err != nil {
		t.Fatalf("CalcBaseAttrs failed: %v", err)
	}

	hp := card.GetAttr(config.AttrType_MAXHP)
	atk := card.GetAttr(config.AttrType_ATTACK)

	if hp <= 0 || atk <= 0 {
		t.Fatalf("attrs invalid: hp=%v atk=%v", hp, atk)
	}

	t.Logf("level 1: hp=%v atk=%v def=%v", hp, atk, card.GetAttr(config.AttrType_DEFENCE))

	// 升级后属性应增长
	card.Level = 2
	if err := card.CalcBaseAttrs(1); err != nil {
		t.Fatalf("CalcBaseAttrs(level2) failed: %v", err)
	}

	if card.GetAttr(config.AttrType_MAXHP) <= hp {
		t.Fatalf("hp should grow with level: level1=%v level2=%v", hp, card.GetAttr(config.AttrType_MAXHP))
	}
}

// TestCardInvalidConfig 未知英雄配置应返回错误而不是空卡
func TestCardInvalidConfig(t *testing.T) {
	loadConfig(t)

	if _, err := NewCard(999999999); err == nil {
		t.Fatal("expect error for unknown hero config")
	}
}

// TestCardRecalc 验证完整属性链：基础 + 星级 → 百分比折算 → 战力
func TestCardRecalc(t *testing.T) {
	loadConfig(t)

	card, err := NewCard(2090101)
	if err != nil {
		t.Fatalf("NewCard failed: %v", err)
	}

	if err := card.Recalc(1); err != nil {
		t.Fatalf("Recalc failed: %v", err)
	}

	power1 := card.BattlePower
	hp1 := card.GetAttr(config.AttrType_MAXHP)

	if power1 <= 0 {
		t.Fatalf("battle power should be positive, got %d", power1)
	}

	// 升星后战力应提升
	card.Star = 3
	if err := card.Recalc(1); err != nil {
		t.Fatalf("Recalc after starup failed: %v", err)
	}

	if card.BattlePower <= power1 {
		t.Fatalf("battle power should grow with star: %d -> %d", power1, card.BattlePower)
	}

	if card.GetAttr(config.AttrType_MAXHP) < hp1 {
		t.Fatalf("hp should not drop after starup: %v -> %v", hp1, card.GetAttr(config.AttrType_MAXHP))
	}

	t.Logf("star1 power=%d hp=%v; star3 power=%d hp=%v", power1, hp1, card.BattlePower, card.GetAttr(config.AttrType_MAXHP))
}
