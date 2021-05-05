package costloot

import (
	"errors"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog/log"
)

var (
	ErrCostLootInvalidType = errors.New("cost loot type invalid")
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

func (m *CostLootManager) GainLootByList(list []*define.LootData) error {
	for _, data := range list {
		if !utils.BetweenInt32(data.LootType, define.CostLoot_Start, define.CostLoot_End) {
			log.Error().Caller().Interface("loot_data", data).Msg("invalid loot type")
			continue
		}

		err := m.objs[data.LootType].GainLoot(data.LootMisc, data.LootNum)
		if !utils.ErrCheck(err, "GainLoot failed when CostLootManager.GainLootByList") {
			continue
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
