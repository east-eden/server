package costloot

import (
	"errors"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/utils"
	"e.coding.net/mmstudio/blade/server/utils/random"
	"github.com/rs/zerolog/log"
)

var (
	ErrCostLootInvalidType = errors.New("cost loot type invalid")
)

type gainLootHandler func(*auto.CostLootEntry) error

type CostLootManager struct {
	Owner   define.PluginObj
	objs    [define.CostLoot_End]define.CostLooter
	handler map[int32]gainLootHandler
}

func NewCostLootManager(owner define.PluginObj, objs ...define.CostLooter) *CostLootManager {
	m := &CostLootManager{
		Owner:   owner,
		handler: make(map[int32]gainLootHandler, 8),
	}

	m.handler[define.LootKind_Fixed] = m.gainLootFixed
	m.handler[define.LootKind_RandProb] = m.gainLootRandProb
	m.handler[define.LootKind_RandWeight] = m.gainLootRandWeight
	m.handler[define.LootKind_Assemble] = m.gainLootAssemble

	return m
}

// 固定掉落
func (m *CostLootManager) gainLootFixed(entry *auto.CostLootEntry) error {
	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		// sub loot
		if entry.Type[n] == define.CostLoot_SubLoot {
			_ = m.GainLoot(entry.Misc[n])
			continue
		}

		num := random.Int32(entry.NumMin[n], entry.NumMax[n])
		err := m.objs[entry.Type[n]].GainLoot(entry.Misc[n], num*entry.LootTimes)
		if err != nil {
			return err
		}
	}

	return nil
}

// 随机机率掉落
func (m *CostLootManager) gainLootRandProb(entry *auto.CostLootEntry) error {
	for times := 0; times < int(entry.LootTimes); times++ {
		for n := range entry.Type {
			if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
				continue
			}

			randVal := random.Int32(1, define.PercentBase)
			if randVal <= entry.Prob[n] {
				num := random.Int32(entry.NumMin[n], entry.NumMax[n])
				err := m.objs[entry.Type[n]].GainLoot(entry.Misc[n], num)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// 随机权重掉落
func (m *CostLootManager) gainLootRandWeight(entry *auto.CostLootEntry) error {
	randomList := auto.GetCostLootRandomList(entry)

	// 可重复
	if entry.CanRepeat {
		for n := 0; n < int(entry.LootTimes); n++ {
			it, err := random.PickOne(randomList, nil)
			if !utils.ErrCheck(err, "PickOne failed when CostLootManager.gainLootRandWeight", randomList) {
				return err
			}

			lootData := it.(*auto.LootElement)
			err = m.objs[lootData.LootType].GainLoot(lootData.LootMisc, random.Int32(lootData.LootNumMin, lootData.LootNumMax))
			if err != nil {
				return err
			}
		}
	} else {
		// 不可重复
		its, err := random.PickUnrepeated(randomList, int(entry.LootTimes), nil)
		if !utils.ErrCheck(err, "PickUnrepeated failed when CostLootManager.gainLootRandWeight", randomList) {
			return err
		}

		for _, it := range its {
			lootData := it.(*auto.LootElement)
			err = m.objs[lootData.LootType].GainLoot(lootData.LootMisc, random.Int32(lootData.LootNumMin, lootData.LootNumMax))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 集合掉落
func (m *CostLootManager) gainLootAssemble(entry *auto.CostLootEntry) error {
	randomList := auto.GetCostLootRandomList(entry)

	// 可重复
	if entry.CanRepeat {
		for n := 0; n < int(entry.LootTimes); n++ {
			it, err := random.PickOne(randomList, nil)
			if !utils.ErrCheck(err, "PickOne failed when CostLootManager.gainLootRandWeight", randomList) {
				return err
			}

			lootData := it.(*auto.LootElement)
			if err := m.GainLoot(lootData.LootMisc); err != nil {
				return err
			}
		}
	} else {
		// 不可重复
		its, err := random.PickUnrepeated(randomList, int(entry.LootTimes), nil)
		if !utils.ErrCheck(err, "PickUnrepeated failed when CostLootManager.gainLootRandWeight", randomList) {
			return err
		}

		for _, it := range its {
			lootData := it.(*auto.LootElement)

			if err := m.GainLoot(lootData.LootMisc); err != nil {
				return err
			}
		}
	}

	return nil
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

	if entry.LootKind != define.LootKind_Fixed {
		return nil
	}

	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		num := random.Int32(entry.NumMin[n], entry.NumMax[n])
		err := m.objs[entry.Type[n]].CanGain(entry.Misc[n], num*entry.LootTimes)
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

	h, ok := m.handler[entry.LootKind]
	if !ok {
		return errors.New("loot kind invalid")
	}

	return h(entry)
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

	if entry.LootKind != define.LootKind_Fixed {
		return ErrCostLootInvalidType
	}

	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		num := random.Int32(entry.NumMin[n], entry.NumMax[n])
		err := m.objs[entry.Type[n]].CanCost(entry.Misc[n], num)
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

	if entry.LootKind != define.LootKind_Fixed {
		return ErrCostLootInvalidType
	}

	for n := range entry.Type {
		if !utils.BetweenInt32(entry.Type[n], define.CostLoot_Start, define.CostLoot_End) {
			continue
		}

		num := random.Int32(entry.NumMin[n], entry.NumMax[n])
		err := m.objs[entry.Type[n]].DoCost(entry.Misc[n], num)
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
			LootNum:  random.Int32(entry.NumMin[n], entry.NumMax[n]),
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
