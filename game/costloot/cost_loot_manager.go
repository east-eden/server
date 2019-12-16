package costloot

import (
	"fmt"

	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
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

func (m *CostLootManager) GainLoot(id int32) error {
	entry := global.GetCostLootEntry(id)
	if entry == nil {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].GainLoot(entry.TypeMisc, entry.Num)
}

func (m *CostLootManager) DoCost(id int32) error {
	entry := global.GetCostLootEntry(id)
	if entry == nil {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry, id:%d", id)
	}

	if entry.Type < 0 || entry.Type >= define.CostLoot_End {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry type, id:%d", id)
	}

	return m.objs[entry.Type].DoCost(entry.TypeMisc, entry.Num)
}
