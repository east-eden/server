package config

import (
	"testing"

	"github.com/east-eden/server/excel"
)

// TestLoadJSONEntries 验证 kof 体系的 json 配置表能正确加载，
// 并且枚举字段被解析成对应的枚举类型。
func TestLoadJSONEntries(t *testing.T) {
	excel.ReadAllJsonEntries("../../config/json")

	// 金币（kof Goods 表里的数值道具）
	gold, ok := GetGoodsEntry(10100000)
	if !ok {
		t.Fatal("goods 10100000 (金币) not loaded")
	}

	if gold.Type != GoodsType_Value {
		t.Fatalf("gold type = %v, want %v", gold.Type, GoodsType_Value)
	}

	// 枚举表：GoodsType 至少应包含数值/英雄卡两类
	if GoodsType_Value == 0 {
		t.Fatal("GoodsType enum not loaded")
	}
}
