package auto

import (
	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/utils"
	"e.coding.net/mmstudio/blade/server/utils/random"
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
	list []random.Item
}

func (l *LootRandomList) GetItemList() []random.Item {
	return l.list
}

// 获取掉落随机库列表
func GetCostLootRandomList(entry *CostLootEntry) *LootRandomList {
	ls := &LootRandomList{
		list: make([]random.Item, 0, 10),
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
