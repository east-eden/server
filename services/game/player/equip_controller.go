package player

import (
	"context"
	"errors"
	"fmt"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/item"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/shopspring/decimal"
)

// 计算装备经验
func GetEquipReturnData(equip *item.Equip) int32 {
	if equip == nil {
		return 0
	}

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return 0
	}

	// 1级经验不算折损率
	equipLv1Entry, ok := auto.GetEquipLevelupEntry(1)
	if !ok {
		return 0
	}
	equipLv1Exp := equipLv1Entry.Exp[equip.ItemEntry.Quality]

	// 已升级累计的经验
	levelTotalExp := auto.GetEquipLevelTotalExp(int32(equip.Level), equip.ItemEntry.Quality)
	if !ok {
		return 0
	}
	equiplvTotalExp := levelTotalExp + equip.Exp - equipLv1Exp

	// 物品总经验 = 物品1级经验 + 已消耗所有经验 * 经验折损率
	return int32(globalConfig.EquipSwallowExpLoss.Mul(decimal.NewFromInt32(equiplvTotalExp)).Round(0).IntPart()) + equipLv1Exp
}

// 通过经验计算返还经验道具和金币
func GetItemsByExp(exp int32) (items map[int32]int32, gold int32) {
	items = make(map[int32]int32)
	gold = 0

	globalConfig, _ := auto.GetGlobalConfig()
	for {
		if exp <= 0 {
			break
		}

		// 没有可补的道具了
		expItem := globalConfig.GetEquipExpItemByExp(exp)
		if expItem == nil {
			break
		}

		items[expItem.ItemTypeId] = exp / expItem.Exp
		gold += exp / expItem.Exp * expItem.Exp * globalConfig.EquipLevelupExpGoldRatio

		exp %= expItem.Exp
	}

	return
}

// 装备升级
func (m *ItemManager) EquipLevelup(equipId int64, stuffItems, expItems []int64) error {
	i, err := m.GetItem(equipId)
	utils.ErrPrint(err, "EquipLevelup failed", equipId, m.owner.ID)

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return auto.ErrGlobalConfigInvalid
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
		itemExps[it] = GetEquipReturnData(stuffEquip)
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

			levelExp := nextLevelEntry.Exp[equip.ItemEntry.Quality]
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

			items, gold := GetItemsByExp(exp)
			for id, num := range items {
				err := m.owner.ItemManager().GainLoot(id, num)
				utils.ErrPrint(err, "gain loot failed when equip levelup return exp items", id, num)
			}

			err = m.owner.TokenManager().GainLoot(define.Token_Gold, gold)
			utils.ErrPrint(err, "gain loot failed when equip levelup return exp items", exp, gold)
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
	fields := map[string]any{
		"level": equip.Level,
		"exp":   equip.Exp,
	}
	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Item, equip.Id, fields)
	utils.ErrPrint(err, "UpdateFields failed when ItemManager.EquipLevelup", equip.GetID(), m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)

	return err
}

// gm装备升级
func (m *ItemManager) GmEquipLevelup(equipTypeId int32, level int32, exp int32) error {
	it := m.GetItemByTypeId(equipTypeId)
	if it == nil {
		return ErrItemNotFound
	}

	if it.GetType() != define.Item_TypeEquip {
		return fmt.Errorf("GmEquipLevelup failed, wrong item<%d> type", it.Opts().TypeId)
	}

	equip, ok := it.(*item.Equip)
	if !ok {
		return fmt.Errorf("GmEquipLevelup failed, cannot assert to equip<%d>", it.Opts().TypeId)
	}

	if level < 0 {
		level = int32(equip.Level)
	}

	if exp < 0 {
		exp = int32(equip.Exp)
	}

	_, ok = auto.GetEquipLevelupEntry(level)
	if !ok {
		return fmt.Errorf("GmEquipLevelup failed, cannot find EquipLevelupEntry<%d>", equip.Level+1)
	}

	equip.Level = int8(level)
	equip.Exp = exp

	// save
	fields := map[string]any{
		"level": equip.Level,
		"exp":   equip.Exp,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Item, equip.Id, fields)
	utils.ErrPrint(err, "UpdateFields failed when ItemManager.GmEquipLevelup", equip.GetID(), m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)

	return err
}

// 装备突破
func (m *ItemManager) EquipPromote(equipId int64) error {
	it, err := m.GetItem(equipId)
	if !utils.ErrCheck(err, "EquipPromote failed", equipId, m.owner.ID) {
		return err
	}

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return auto.ErrGlobalConfigInvalid
	}

	if it.GetType() != define.Item_TypeEquip {
		return ErrItemInvalidType
	}

	equip := it.(*item.Equip)
	if equip.Promote >= define.Equip_Max_Promote_Times {
		return ErrEquipPromoteTimesFull
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
	fields := map[string]any{
		"promote": equip.Promote,
	}
	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Item, equip.Id, fields)
	utils.ErrPrint(err, "UpdateFields failed when ItemManager.EquipPromote", equip.Id, m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)
	return err
}

// gm装备突破
func (m *ItemManager) GmEquipPromote(equipTypeId int32, promote int32) error {
	it := m.GetItemByTypeId(equipTypeId)
	if it == nil {
		return ErrItemNotFound
	}

	if it.GetType() != define.Item_TypeEquip {
		return ErrItemInvalidType
	}

	equip := it.(*item.Equip)
	if equip.Promote >= define.Equip_Max_Promote_Times {
		return ErrEquipPromoteTimesFull
	}

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return auto.ErrGlobalConfigInvalid
	}

	if int(equip.Promote) >= len(globalConfig.EquipPromoteLevelLimit)-1 {
		return errors.New("equip has not levelup to max")
	}

	equip.Promote = int8(promote)

	// save
	fields := map[string]any{
		"promote": equip.Promote,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Item, equip.Id, fields)
	utils.ErrPrint(err, "UpdateFields failed when ItemManager.GmEquipPromote", equip.Id, m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)
	return err
}

// 装备升星
func (m *ItemManager) EquipStarup(equipId int64, stuffIds []int64) error {
	it, err := m.GetItem(equipId)
	if !utils.ErrCheck(err, "EquipStarup failed", equipId, m.owner.ID) {
		return err
	}

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return auto.ErrGlobalConfigInvalid
	}

	if it.GetType() != define.Item_TypeEquip {
		return ErrItemInvalidType
	}

	equip := it.(*item.Equip)

	// 升星次数限制
	if int32(equip.Star) >= globalConfig.EquipQualityStarupTimes[equip.ItemEntry.Quality+1] {
		return ErrEquipStarTimesFull
	}

	// 材料
	stuffItems := make(map[int64]*item.Equip)
	for _, id := range stuffIds {
		stuff, err := m.GetItem(id)
		if err != nil {
			continue
		}

		// 材料必须是id相同的装备
		if stuff.Opts().TypeId != equip.TypeId {
			continue
		}

		// 不能是升星的物品
		if stuff.Opts().Id == equip.Id {
			continue
		}

		stuffItems[stuff.Opts().Id] = stuff.(*item.Equip)
	}

	nextStar := equip.Star + 1

	// 金币消耗
	costGold := globalConfig.EquipStarupCostGold[nextStar]
	err = m.owner.TokenManager().CanCost(define.Token_Gold, costGold)
	if err != nil {
		return err
	}

	// 材料消耗
	costItemNum := globalConfig.EquipStarupCostItemNum[nextStar]
	if len(stuffItems) < int(costItemNum) {
		return errors.New("not enough stuff equips")
	}

	err = m.owner.TokenManager().DoCost(define.Token_Gold, costGold)
	if !utils.ErrCheck(err, "TokenManager.DoCost failed when ItemManager.EquipStarup", equipId, costGold, m.owner.ID) {
		return err
	}

	// 消耗同名物品并返还材料
	for _, stuff := range stuffItems {
		// 没有强化过的装备不计算返还
		if stuff.Level == 1 && stuff.Exp == 0 {
			continue
		}

		returnExp := GetEquipReturnData(stuff)
		items, gold := GetItemsByExp(returnExp)
		for id, num := range items {
			err := m.owner.ItemManager().GainLoot(id, num)
			utils.ErrPrint(err, "ItemManager.GainLoot failed when ItemManager.EquipStarup", id, num)
		}

		err := m.owner.TokenManager().GainLoot(define.Token_Gold, gold)
		utils.ErrPrint(err, "ItemManager.GainLoot failed when ItemManager.EquipStarup", stuff, gold)
	}

	equip.Star = nextStar

	// save
	fields := map[string]any{
		"star": equip.Star,
	}
	err = store.GetStore().UpdateFields(context.Background(), define.StoreType_Item, equip.Id, fields)
	utils.ErrPrint(err, "UpdateFields failed when ItemManager.EquipStarup", equip.Id, m.owner.ID)

	// send client
	m.SendEquipUpdate(equip)
	return err
}

// gm 升星
func (m *ItemManager) GmEquipStarup(typeId int32, star int32) error {
	it := m.GetItemByTypeId(typeId)
	if it == nil {
		return ErrItemNotFound
	}

	if it.GetType() != define.Item_TypeEquip {
		return ErrItemInvalidType
	}

	equip := it.(*item.Equip)
	if equip.Star >= define.Equip_Max_Starup_Times {
		return ErrEquipStarTimesFull
	}

	equip.Star = int8(star)

	// save
	fields := map[string]any{
		"star": equip.Star,
	}
	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Item, equip.Id, fields)
	utils.ErrPrint(err, "UpdateFields failed when ItemManager.GmEquipStarup", equip.Id, m.owner.ID)

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

func (m *ItemManager) GenEquipListPB() []*pbGlobal.Equip {
	equips := make([]*pbGlobal.Equip, 0, m.GetItemNums(int(define.Container_Equip)))
	m.ca.RangeByIdx(int(define.Container_Equip), func(val item.Itemface) bool {
		it, ok := val.(*item.Equip)
		if !ok {
			return true
		}

		equips = append(equips, it.GenEquipPB())
		return true
	})

	return equips
}
