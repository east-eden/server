package costloot

import (
	"errors"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/utils"
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

// generate loot list from entry id
func (m *CostLootManager) GenLootList(lootId int32) []*define.LootData {
	ret := make([]*define.LootData, 0, 10)
	entry, ok := auto.GetCostLootEntry(lootId)
	if !ok {
		return ret
	}

	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		ret = append(ret, &define.LootData{
			LootType: entry.Type[n],
			LootMisc: entry.Misc[n],
			LootNum:  entry.Num[n],
		})
	}

	return ret
}

// pack loot list
func (m *CostLootManager) PackLootList(lootList []*define.LootData) []*define.LootData {
	packedMap := make(map[int64]int32, 8)
	for _, data := range lootList {
		packId := utils.PackId(data.LootType, data.LootMisc)
		packedMap[packId] += data.LootNum
	}

	ret := make([]*define.LootData, 0, 10)
	for key, num := range packedMap {
		lootType := utils.HighId(key)
		lootMisc := utils.LowId(key)
		ret = append(ret, &define.LootData{
			LootType: lootType,
			LootMisc: lootMisc,
			LootNum:  num,
		})
	}

	return ret
}
