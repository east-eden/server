package costloot

import (
	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/utils"
)

type CostLootManager struct {
	Owner define.PluginObj
	objs  [define.CostLoot_End]define.CostLooter
}

func NewCostLootManager(owner define.PluginObj, objs ...define.CostLooter) *CostLootManager {
	return &CostLootManager{Owner: owner}
}

func (m *CostLootManager) Init(objs ...define.CostLooter) {
	for _, o := range objs {
		m.objs[o.GetCostLootType()] = o
	}
}

func (m *CostLootManager) CanGain(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return nil
	}

	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		err := m.objs[entry.Type[n]].CanGain(entry.Misc[n], entry.Num[n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *CostLootManager) GainLoot(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return nil
	}

	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		err := m.objs[entry.Type[n]].GainLoot(entry.Misc[n], entry.Num[n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *CostLootManager) CanCost(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return nil
	}

	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		err := m.objs[entry.Type[n]].CanCost(entry.Misc[n], entry.Num[n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *CostLootManager) DoCost(id int32) error {
	entry, ok := auto.GetCostLootEntry(id)
	if !ok {
		return nil
	}

	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		err := m.objs[entry.Type[n]].DoCost(entry.Misc[n], entry.Num[n])
		if err != nil {
			return err
		}
	}

	return nil
}
