package costloot

import (
	"fmt"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
)

type CostLootManager struct {
	Owner define.PluginObj
	objs  [define.CostLoot_End]define.CostLootObj
}

func NewCostLootManager(owner define.PluginObj, objs ...define.CostLootObj) *CostLootManager {
	return &CostLootManager{Owner: owner}
}

func (m *CostLootManager) Init(objs ...define.CostLootObj) {
	for _, o := range objs {
		m.objs[o.GetCostLootType()] = o
	}
}

func (m *CostLootManager) CanGain(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].CanGain(entry.Misc, entry.Num)
}

func (m *CostLootManager) GainLoot(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].GainLoot(entry.Misc, entry.Num)
}

func (m *CostLootManager) CanCost(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].CanCost(entry.Misc, entry.Num)
}

func (m *CostLootManager) DoCost(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].DoCost(entry.Misc, entry.Num)
}
