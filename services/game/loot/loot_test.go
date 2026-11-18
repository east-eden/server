package loot

import (
	"sort"
	"testing"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/excel/config"
)

// TestCalcLootResult 用真实 Loot.json 数据验证掉落计算能跑通，
// 且产出的道具 id 都能在 Goods 表里找到。
func TestCalcLootResult(t *testing.T) {
	excel.ReadAllJsonEntries("../../../config/json")

	rows := config.GetLootRows()
	if len(rows) == 0 {
		t.Skip("no loot table loaded")
	}

	// 固定抽样前 200 个 loot id，避免 map 遍历顺序导致结果不可复现
	ids := make([]int32, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	if len(ids) > 200 {
		ids = ids[:200]
	}

	for _, lootID := range ids {
		itemIds, counts, err := CalcLootResult(lootID)
		if err != nil {
			t.Fatalf("CalcLootResult(%d) failed: %v", lootID, err)
		}

		if len(itemIds) != len(counts) {
			t.Fatalf("loot %d: itemIds/counts length mismatch: %d vs %d", lootID, len(itemIds), len(counts))
		}

		for i, itemId := range itemIds {
			if _, ok := config.GetGoodsEntry(itemId); !ok {
				t.Fatalf("loot %d produced unknown goods %d (count %d)", lootID, itemId, counts[i])
			}
		}
	}
}
