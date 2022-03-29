package auto

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/utils/random"
)

type LootElement struct {
	LootType   int32
	LootMisc   int32
	LootNumMin int32
	LootNumMax int32
	LootWeight int32
}

// random.Item interface
func (e *LootElement) GetId() int {
	return int(e.LootMisc)
}

func (e *LootElement) GetWeight() int {
	return int(e.LootWeight)
}

// RandomPicker interface
type LootRandomList struct {
	list []random.Item[int]
}

func (l *LootRandomList) GetItemList() []random.Item[int] {
	return l.list
}

// 获取掉落随机库列表
func GetCostLootRandomList(entry *CostLootEntry) *LootRandomList {
	ls := &LootRandomList{
		list: make([]random.Item[int], 0, 8),
	}

	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		ls.list = append(ls.list, &LootElement{
			LootType:   entry.Type[n],
			LootMisc:   entry.Misc[n],
			LootNumMin: entry.NumMin[n],
			LootNumMax: entry.NumMax[n],
			LootWeight: entry.Prob[n],
		})
	}

	return ls
}
