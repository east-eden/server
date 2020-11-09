package costloot

import (
	"fmt"

	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
)

type CostLootManager struct {
	Owner define.PluginObj
	objs  [define.CostLoot_End]define.CostLootObj
}

func NewCostLootManager(owner define.PluginObj, objs ...define.CostLootObj) *CostLootManager {
	m := &CostLootManager{Owner: owner}

	for _, o := range objs {
		m.objs[o.GetCostLootType()] = o
	}

	return m
}

func (m *CostLootManager) CanGain(id int32) error {
	entry := entries.GetCostLootEntry(id)
	if entry == nil {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].CanGain(entry.Misc, entry.Num)
}

func (m *CostLootManager) GainLoot(id int32) error {
	entry := entries.GetCostLootEntry(id)
	if entry == nil {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].GainLoot(entry.Misc, entry.Num)
}

func (m *CostLootManager) CanCost(id int32) error {
	entry := entries.GetCostLootEntry(id)
	if entry == nil {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].CanCost(entry.Misc, entry.Num)
}

func (m *CostLootManager) DoCost(id int32) error {
	entry := entries.GetCostLootEntry(id)
	if entry == nil {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].DoCost(entry.Misc, entry.Num)
}
