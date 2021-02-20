package costloot

import (
	"fmt"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	"bitbucket.org/east-eden/server/utils"
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

	var err error
	for n := range entry.Type {
		if !utils.Between(int(entry.Type[n]), define.CostLoot_Begin, define.CostLoot_End) {
			return fmt.Errorf("gain loot error, non-existing cost_loot_entry type, id:%d", id)
		}

		if err != nil {
			return err
		}

		err = m.objs[entry.Type[n]].CanGain(entry.Misc[n], entry.Num[n])
	}

	return nil
}

func (m *CostLootManager) GainLoot(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return fmt.Errorf("gain loot error, non-existing cost_loot_entry, id:%d", id)
	}

	var err error
	for n := range entry.Type {
		if !utils.Between(int(entry.Type[n]), define.CostLoot_Begin, define.CostLoot_End) {
			return fmt.Errorf("gain loot error, non-existing cost_loot_entry type, id:%d", id)
		}

		if err != nil {
			return err
		}

		err = m.objs[entry.Type[n]].GainLoot(entry.Misc[n], entry.Num[n])
	}

	return nil
}

func (m *CostLootManager) CanCost(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry, id:%d", id)
	}

	var err error
	for n := range entry.Type {
		if !utils.Between(int(entry.Type[n]), define.CostLoot_Begin, define.CostLoot_End) {
			return fmt.Errorf("do cost error, non-existing cost_loot_entry type, id:%d", id)
		}

		if err != nil {
			return err
		}

		err = m.objs[entry.Type[n]].CanCost(entry.Misc[n], entry.Num[n])
	}

	return nil
}

func (m *CostLootManager) DoCost(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return fmt.Errorf("do cost error, non-existing cost_loot_entry, id:%d", id)
	}

	var err error
	for n := range entry.Type {
		if !utils.Between(int(entry.Type[n]), define.CostLoot_Begin, define.CostLoot_End) {
			return fmt.Errorf("do cost error, non-existing cost_loot_entry type, id:%d", id)
		}

		if err != nil {
			return err
		}

		err = m.objs[entry.Type[n]].DoCost(entry.Misc[n], entry.Num[n])
	}

	return nil
}
