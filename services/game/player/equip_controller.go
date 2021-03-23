package player

import (
	"errors"
	"fmt"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/services/game/item"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
)

// 装备升级
func (m *ItemManager) EquipLevelup(equipId int64, stuffItems, expItems []int64) error {
	i, err := m.GetItem(equipId)
	utils.ErrPrint(err, "EquipLevelup failed", equipId, m.owner.ID)

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return errors.New("invalid global config")
	}

	if i.GetType() != define.Item_TypeEquip {
		return fmt.Errorf("EquipLevelup failed, wrong item<%d> type", i.Opts().TypeId)
	}

	equip, ok := i.(*item.Equip)
	if !ok {
		return fmt.Errorf("EquipLevelup failed, cannot assert to equip<%d>", equipId)
	}

	levelUpEntry, ok := auto.GetEquipLevelupEntry(int32(equip.Level) + 1)
	if !ok {
		return fmt.Errorf("EquipLevelup failed, cannot find EquipLevelupEntry<%d>", equip.Level+1)
	}

	if int32(equip.Promote) < levelUpEntry.PromoteLimit {
		return fmt.Errorf("EquipLevelup failed, PromoteLevel<%d> limit", equip.Promote)
	}

	// 所有合法的消耗物品及对应的经验值
	itemExps := make(map[item.Itemface]int32)

	// 剔除重复的物品id
	unrepeatedItemId := make(map[int64]struct{})

	// 吞噬材料
	for _, id := range stuffItems {
		it, err := m.GetItem(id)
		if !utils.ErrCheck(err, "cannot find item", id) {
			continue
		}

		if it.Opts().ItemEntry.Type != define.Item_TypeEquip {
			continue
		}

		// 重复的id不计入
		if _, ok := unrepeatedItemId[id]; ok {
			continue
		}

		stuffEquip := it.(*item.Equip)

		// 1级经验不算折损率
		equipLv1Entry, ok := auto.GetEquipLevelupEntry(1)
		if !ok {
			log.Error().Caller().Msg("can not find equip levelup 1 entry")
			continue
		}
		equipLv1Exp := equipLv1Entry.Exp[stuffEquip.ItemEntry.Quality]

		// 已升级累计的经验
		equipLvEntry, ok := auto.GetEquipLevelupEntry(int32(stuffEquip.Level))
		if !ok {
			log.Error().Caller().Int8("level", stuffEquip.Level).Msg("can not find equip levelup entry")
			continue
		}
		equiplvTotalExp := equipLvEntry.Exp[stuffEquip.ItemEntry.Quality] + stuffEquip.Exp - equipLv1Exp

		// 物品总经验 = 物品1级经验 + 已消耗所有经验 * 经验折损率
		itemExps[it] = int32(int64(equipLv1Exp) + int64(equiplvTotalExp)*int64(globalConfig.EquipSwallowExpLoss)/int64(define.PercentBase))
		unrepeatedItemId[id] = struct{}{}
	}

	// 经验道具
	for _, id := range expItems {
		it, err := m.GetItem(id)
		if !utils.ErrCheck(err, "cannot find item", id) {
			continue
		}

		if it.GetType() != define.Item_TypeItem {
			continue
		}

		if it.Opts().ItemEntry.SubType != define.Item_SubType_Item_EquipExp {
			continue
		}

		if _, ok := unrepeatedItemId[id]; ok {
			continue
		}

		itemExps[it] = it.Opts().ItemEntry.PublicMisc[0]
		unrepeatedItemId[id] = struct{}{}
	}

	// 状态是否改变
	changed := false

	// 升级处理
	levelupFn := func(itemId int64, exp int32) bool {
		nextLevelEntry, ok := auto.GetEquipLevelupEntry(int32(equip.Level) + 1)
		if !ok {
			return false
		}

		// 金币限制
		costGold := int32(int64(exp) * int64(globalConfig.EquipLevelupExpGoldRatio))
		if costGold < 0 {
			return false
		}

		// 突破限制
		if int32(equip.Promote) < nextLevelEntry.PromoteLimit {
			return false
		}

		if err := m.owner.TokenManager().CanCost(define.Token_Gold, costGold); err != nil {
			return false
		}

		// overflow
		if equip.Exp+exp < 0 {
			return false
		}

		equip.Exp += exp
		changed = true
		reachLimit := false
		for {
			curLevelEntry, _ := auto.GetEquipLevelupEntry(int32(equip.Level))
			nextLevelEntry, ok := auto.GetEquipLevelupEntry(int32(equip.Level) + 1)
			if !ok {
				reachLimit = true
				break
			}

			// 突破限制
			if int32(equip.Promote) < nextLevelEntry.PromoteLimit {
				reachLimit = true
				break
			}

			levelExp := nextLevelEntry.Exp[equip.ItemEntry.Quality] - curLevelEntry.Exp[equip.ItemEntry.Quality]
			if equip.Exp < levelExp {
				break
			}

			equip.Level++
			equip.Exp -= levelExp
		}

		// 消耗
		err := m.owner.TokenManager().DoCost(define.Token_Gold, costGold)
		utils.ErrPrint(err, "TokenManager DoCost failed", costGold)

		err = m.owner.ItemManager().CostItemByID(itemId, 1)
		utils.ErrPrint(err, "ItemManager CostItemByID failed", itemId)

		// 返还处理
		if reachLimit && equip.Exp > 0 {
			exp := equip.Exp
			equip.Exp = 0

			for {
				if exp <= 0 {
					break
				}

				// 没有可补的道具了
				expItem := globalConfig.GetEquipExpItemByExp(exp)
				if expItem == nil {
					break
				}

				err := m.owner.ItemManager().GainLoot(expItem.ItemTypeId, exp/expItem.Exp)
				utils.ErrPrint(err, "gain loot failed when equip levelup return exp items", exp, expItem.ItemTypeId)

				returnToken := exp / expItem.Exp * expItem.Exp * globalConfig.EquipLevelupExpGoldRatio
				err = m.owner.TokenManager().GainLoot(define.Token_Gold, returnToken)
				utils.ErrPrint(err, "gain loot failed when equip levelup return exp items", exp, returnToken)

				exp %= expItem.Exp
			}
		}

		return true
	}

	continueCheck := true
	for it, exp := range itemExps {
		if !continueCheck {
			break
		}

		var n int32
		for n = 0; n < it.Opts().Num; n++ {
			continueCheck = levelupFn(it.Opts().Id, exp)
			if !continueCheck {
				break
			}
		}
	}

	// 经验等级道具均没有改变
	if !changed {
		return nil
	}

	// save
	fields := map[string]interface{}{
		"level": equip.Level,
		"exp":   equip.Exp,
	}
	err = store.GetStore().SaveHashObjectFields(define.StoreType_Item, equip.OwnerId, equip.Id, equip, fields)
	utils.ErrPrint(err, "SaveFields failed when EquipLevelup", equip.GetID(), m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)

	return err
}

// 装备突破
func (m *ItemManager) EquipPromote(equipId int64) error {
	it, err := m.GetItem(equipId)
	utils.ErrPrint(err, "EquipPromote failed", equipId, m.owner.ID)

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return errors.New("invalid global config")
	}

	if it.GetType() != define.Item_TypeEquip {
		return errors.New("invalid item type")
	}

	equip := it.(*item.Equip)
	if equip.Promote >= define.Equip_Max_Promote_Times {
		return errors.New("promote times full")
	}

	// 队伍等级
	if m.owner.Level < globalConfig.EquipPromoteLevelLimit[equip.Promote+1] {
		return errors.New("team level limit")
	}

	// 装备是否达到等级上限
	levelupEntry, ok := auto.GetEquipLevelupEntry(int32(equip.Level) + 1)
	if !ok {
		return errors.New("equip reach max level")
	}

	if int32(equip.Promote) >= levelupEntry.PromoteLimit {
		return errors.New("equip has not levelup to max")
	}

	// 消耗
	costId := equip.EquipEnchantEntry.PromoteCostId[equip.Promote+1]
	err = m.owner.CostLootManager().CanCost(costId)
	if !utils.ErrCheck(err, "EquipPromote can cost failed", equipId, costId, m.owner.ID) {
		return err
	}

	err = m.owner.CostLootManager().DoCost(costId)
	if !utils.ErrCheck(err, "EquipPromote do cost failed", equipId, costId, m.owner.ID) {
		return err
	}

	equip.Promote++

	// save
	fields := map[string]interface{}{
		"promote": equip.Promote,
	}
	err = store.GetStore().SaveHashObjectFields(define.StoreType_Item, equip.OwnerId, equip.Id, equip, fields)
	utils.ErrPrint(err, "SaveFields failed when EquipPromote", equip.Id, m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)
	return err
}

func (m *ItemManager) SendEquipUpdate(e *item.Equip) {
	msg := &pbGlobal.S2C_EquipUpdate{
		EquipId: e.GetID(),
		EquipData: &pbGlobal.EquipData{
			Exp:      e.Exp,
			Level:    int32(e.Level),
			Promote:  int32(e.Promote),
			Lock:     e.Lock,
			EquipObj: e.EquipObj,
		},
	}

	m.owner.SendProtoMessage(msg)
}
