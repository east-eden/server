package loot

import (
	"fmt"
	"math/rand"

	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/utils"
)

func CalcLootListResult(lootIds []int32) (itemIds []int32, counts []int32) {
	for _, lootId := range lootIds {
		subItems, subCounts, err := CalcLootResult(lootId)
		if !utils.ErrCheck(err, "CalcLootListResult with lootId") {
			continue
		}

		itemIds = append(itemIds, subItems...)
		counts = append(counts, subCounts...)
	}
	return
}

func CalcLootResult(lootId int32) (itemIds []int32, counts []int32, err error) {
	lootConf, found := config.GetLootEntry(lootId)
	if !found {
		err = fmt.Errorf("invalid loot id <%d>", lootId)
		return
	}

	if len(lootConf.WeightItemWeight) > 0 {
		index := utils.RollInt32(lootConf.WeightItemWeight)
		if index >= len(lootConf.WeightItemNum) {
			err = fmt.Errorf("rand WeightItem index<%d> out of range<%d %d", index, len(lootConf.WeightItemNum), lootId)
			return
		}

		if lootConf.WeightItemNum[index] > 0 {
			itemIds = append(itemIds, lootConf.WeightItemId[index])
			counts = append(counts, lootConf.WeightItemNum[index])
		}
	}

	if len(lootConf.WeightLootWeight) > 0 {
		index := utils.RollInt32(lootConf.WeightLootWeight)
		if index >= len(lootConf.WeightLootNum) {
			err = fmt.Errorf("rand WeightLoot index<%d> out of range<%d %d", index, len(lootConf.WeightLootNum), lootId)
			return
		}

		if lootConf.WeightLootNum[index] > 0 {
			var suberr error
			subItemIds, subItemCounts, suberr := CalcLootResult(lootConf.WeightLootId[index])
			if !utils.ErrCheck(err, "CalcLootListResult with lootId", lootId) {
				err = suberr
				return
			}

			itemIds = append(itemIds, subItemIds...)
			counts = append(counts, subItemCounts...)
		}
	}

	for i := 0; i < len(lootConf.RandomItemId); i++ {
		itemId := lootConf.RandomItemId[i]
		count := lootConf.RandomItemNum[i]
		if count <= 0 {
			continue
		}
		probability := lootConf.RandomItemProbability[i]
		if rand.Float32() < probability {
			itemIds = append(itemIds, itemId)
			counts = append(counts, count)
		}
	}

	for i := 0; i < len(lootConf.RandomLootId); i++ {
		subLootId := lootConf.RandomLootId[i]
		count := lootConf.RandomLootNum[i]
		if count <= 0 {
			continue
		}
		probability := lootConf.RandomLootProbability[i]
		if rand.Float32() < probability {
			for j := 0; j < int(count); j++ {
				subItemIds, subItemCounts, err := CalcLootResult(subLootId)
				if !utils.ErrCheck(err, "CalcLootListResult with lootId", lootId) {
					continue
				}

				itemIds = append(itemIds, subItemIds...)
				counts = append(counts, subItemCounts...)
			}
		}
	}
	return
}

// 特定掉落不会掉出指定道具
//func CalcLootResultWithException(lootId int32, exceptItemIds map[int32]map[int32]bool) (itemIds []int32, counts []int32, err error) {
//	lootConf, found := config.GetLootEntry(lootId)
//	if !found {
//		err = fmt.Errorf("invalid loot id <%d>", lootId)
//		return
//	}
//
//	if len(lootConf.WeightItemWeight) > 0 {
//
//		modifiedItemWeights := make([]int32, 0, len(lootConf.WeightItemWeight))
//
//		itemMap, ok := exceptItemIds[lootId]
//		if ok && len(itemMap) > 0 {
//			for i := 0; i < len(lootConf.WeightItemId); i++ {
//				if _, has := itemMap[lootConf.WeightItemId[i]]; has {
//					modifiedItemWeights = append(modifiedItemWeights, 0)
//				} else {
//					modifiedItemWeights = append(modifiedItemWeights, lootConf.WeightItemWeight[i])
//				}
//			}
//		} else {
//			modifiedItemWeights = append(modifiedItemWeights, lootConf.WeightItemWeight...)
//		}
//
//		index := utils.RollInt32(modifiedItemWeights)
//		if index >= len(lootConf.WeightItemNum) {
//			err = fmt.Errorf("rand WeightItem index<%d> out of range<%d %d", index, len(lootConf.WeightItemNum), lootId)
//			return
//		}
//
//		if lootConf.WeightItemNum[index] > 0 {
//			itemIds = append(itemIds, lootConf.WeightItemId[index])
//			counts = append(counts, lootConf.WeightItemNum[index])
//		}
//	}
//
//	if len(lootConf.WeightLootWeight) > 0 {
//		index := utils.RollInt32(lootConf.WeightLootWeight)
//		if index >= len(lootConf.WeightLootNum) {
//			err = fmt.Errorf("rand WeightLoot index<%d> out of range<%d %d", index, len(lootConf.WeightLootNum), lootId)
//			return
//		}
//
//		if lootConf.WeightLootNum[index] > 0 {
//			var suberr error
//			subItemIds, subItemCounts, suberr := CalcLootResultWithException(lootConf.WeightLootId[index], exceptItemIds)
//			if !utils.ErrCheck(err, "CalcLootListResult with lootId", lootId) {
//				err = suberr
//				return
//			}
//
//			itemIds = append(itemIds, subItemIds...)
//			counts = append(counts, subItemCounts...)
//		}
//	}
//
//	for i := 0; i < len(lootConf.RandomItemId); i++ {
//		itemId := lootConf.RandomItemId[i]
//		count := lootConf.RandomItemNum[i]
//		if count <= 0 {
//			continue
//		}
//		probability := lootConf.RandomItemProbability[i]
//		if rand.Float32() < probability {
//			itemIds = append(itemIds, itemId)
//			counts = append(counts, count)
//		}
//	}
//
//	for i := 0; i < len(lootConf.RandomLootId); i++ {
//		subLootId := lootConf.RandomLootId[i]
//		count := lootConf.RandomLootNum[i]
//		if count <= 0 {
//			continue
//		}
//		probability := lootConf.RandomLootProbability[i]
//		if rand.Float32() < probability {
//			for j := 0; j < int(count); j++ {
//				subItemIds, subItemCounts, err := CalcLootResultWithException(subLootId, exceptItemIds)
//				if !utils.ErrCheck(err, "CalcLootListResult with lootId", lootId) {
//					continue
//				}
//
//				itemIds = append(itemIds, subItemIds...)
//				counts = append(counts, subItemCounts...)
//			}
//		}
//	}
//	return
//}

// 特定掉落会额外掉出道具
func CalcLootResultWithFilter(lootId int32, filter func(itemTypeId int32) bool) (itemIds []int32, counts []int32, err error) {
	lootConf, found := config.GetLootEntry(lootId)
	if !found {
		err = fmt.Errorf("invalid loot id <%d>", lootId)
		return
	}

	filtedItemIds := make([]int32, 0, len(lootConf.WeightItemId))
	filtedItemWeights := make([]int32, 0, len(lootConf.WeightItemWeight))
	filtedItemNums := make([]int32, 0, len(lootConf.WeightItemNum))
	for idx := range lootConf.WeightItemId {
		if filter(lootConf.WeightItemId[idx]) {
			filtedItemIds = append(filtedItemIds, lootConf.WeightItemId[idx])
			filtedItemNums = append(filtedItemNums, lootConf.WeightItemNum[idx])
			filtedItemWeights = append(filtedItemWeights, lootConf.WeightItemWeight[idx])
		}
	}
	if len(filtedItemWeights) > 0 {
		index := utils.RollInt32(filtedItemWeights)
		if index >= len(filtedItemWeights) {
			err = fmt.Errorf("rand WeightItem with filter index<%d> out of range<%d", index, len(filtedItemWeights))
			return
		}

		if filtedItemNums[index] > 0 && filtedItemWeights[index] > 0 {
			itemIds = append(itemIds, filtedItemIds[index])
			counts = append(counts, filtedItemNums[index])
		}
	}

	if len(lootConf.WeightLootWeight) > 0 {
		index := utils.RollInt32(lootConf.WeightLootWeight)
		if index >= len(lootConf.WeightLootNum) {
			err = fmt.Errorf("rand WeightLoot index<%d> out of range<%d", index, len(lootConf.WeightLootNum))
			return
		}

		if lootConf.WeightLootNum[index] > 0 {
			var suberr error
			subItemIds, subItemCounts, suberr := CalcLootResultWithFilter(lootConf.WeightLootId[index], filter)
			if !utils.ErrCheck(err, "CalcLootListResult with lootId") {
				err = suberr
				return
			}

			itemIds = append(itemIds, subItemIds...)
			counts = append(counts, subItemCounts...)
		}
	}

	filtedItemIds = make([]int32, 0, len(lootConf.RandomItemId))
	filtedItemProbs := make([]float32, 0, len(lootConf.RandomItemProbability))
	filtedItemNums = make([]int32, 0, len(lootConf.RandomItemNum))
	for idx := range lootConf.RandomItemId {
		if filter(lootConf.RandomItemId[idx]) {
			filtedItemIds = append(filtedItemIds, lootConf.RandomItemId[idx])
			filtedItemProbs = append(filtedItemProbs, lootConf.RandomItemProbability[idx])
			filtedItemNums = append(filtedItemNums, lootConf.RandomItemNum[idx])
		}
	}
	for i := 0; i < len(filtedItemIds); i++ {
		itemId := filtedItemIds[i]
		count := filtedItemNums[i]
		if count <= 0 {
			continue
		}
		probability := filtedItemProbs[i]
		if rand.Float32() < probability {
			itemIds = append(itemIds, itemId)
			counts = append(counts, count)
		}
	}

	for i := 0; i < len(lootConf.RandomLootId); i++ {
		subLootId := lootConf.RandomLootId[i]
		count := lootConf.RandomLootNum[i]
		if count <= 0 {
			continue
		}
		probability := lootConf.RandomLootProbability[i]
		if rand.Float32() < probability {
			for j := 0; j < int(count); j++ {
				subItemIds, subItemCounts, err := CalcLootResultWithFilter(subLootId, filter)
				if !utils.ErrCheck(err, "CalcLootListResult with lootId") {
					continue
				}

				itemIds = append(itemIds, subItemIds...)
				counts = append(counts, subItemCounts...)
			}
		}
	}
	return
}
