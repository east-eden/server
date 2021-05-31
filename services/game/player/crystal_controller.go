package player

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/east-eden/server/services/game/item"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	"github.com/east-eden/server/utils/random"
	log "github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

// 初始化晶石属性
func (m *ItemManager) initCrystalAtt(c *item.Crystal) {
	globalConfig, _ := auto.GetGlobalConfig()

	// 初始主属性
	mainAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, define.Crystal_AttTypeMain)
	mainAttItem, err := random.PickOne(mainAttRepoList, func(random.Item) bool {
		return true
	})
	if err != nil {
		log.Error().Err(err).Int64("crystal_id", c.Id).Msg("pick crystal main att failed")
		return
	}

	// 记录主属性库id
	mainAttRepoEntry := mainAttItem.(*auto.CrystalAttRepoEntry)
	c.MainAtt.AttRepoId = mainAttRepoEntry.Id
	c.MainAtt.AttRandRatio = random.Decimal(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1])

	// 随机几条副属性
	viceAttNum := auto.GetCrystalInitViceAttNum(c.ItemEntry.Quality)

	// 初始副属性
	viceAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, define.Crystal_AttTypeVice)
	viceAttItems, err := random.PickUnrepeated(viceAttRepoList, viceAttNum, func(random.Item) bool {
		return true
	})

	if errors.Is(err, random.ErrNoResult) {
		return
	}

	if err != nil {
		log.Error().Err(err).Caller().Int64("crystal_id", c.Id).Msg("pick unrepeated crystal vice att failed")
		return
	}

	for _, v := range viceAttItems {
		viceAttRepoEntry := v.(*auto.CrystalAttRepoEntry)
		c.ViceAtts = append(c.ViceAtts, item.CrystalAtt{
			AttRepoId:    viceAttRepoEntry.Id,
			AttRandRatio: random.Decimal(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1]),
		})
	}
}

// 新增副属性
func (m *ItemManager) generateCrystalViceAtt(c *item.Crystal) {
	if c == nil {
		return
	}

	globalConfig, _ := auto.GetGlobalConfig()

	attType := make(map[int]struct{}, 20)
	for _, att := range c.ViceAtts {
		attType[int(att.AttRepoId)] = struct{}{}
	}

	// 副属性已满4条
	if len(attType) >= define.Crystal_ViceAttNum {
		return
	}

	// 不满4条，则随机一条未曾有过的属性类型
	limiter := func(it random.Item) bool {
		if _, ok := attType[it.GetId()]; ok {
			return false
		}
		return true
	}
	viceAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, define.Crystal_AttTypeVice)
	it, err := random.PickOne(viceAttRepoList, limiter)
	if !utils.ErrCheck(err, "pick one vice att failed", c.Id) {
		return
	}

	attRepoEntry := it.(*auto.CrystalAttRepoEntry)
	c.ViceAtts = append(c.ViceAtts, item.CrystalAtt{
		AttRepoId:    attRepoEntry.Id,
		AttRandRatio: random.Decimal(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1]),
	})
}

// 强化副属性
func (m *ItemManager) enforceCrystalViceAtt(c *item.Crystal) {
	if c == nil {
		return
	}

	globalConfig, _ := auto.GetGlobalConfig()

	// 所有副属性种类对应强化次数
	attType := make(map[int]int, 20)
	for _, att := range c.ViceAtts {
		attType[int(att.AttRepoId)]++
	}

	var viceAttRepoEntry *auto.CrystalAttRepoEntry

	// 如果已有4条副属性，则强化概率皆为1/4
	if len(attType) >= 4 {
		rd := rand.Intn(len(attType))
		for attRepoId := range attType {
			if rd <= 0 {
				viceAttRepoEntry, _ = auto.GetCrystalAttRepoEntry(int32(attRepoId))
				break
			}
			rd--
		}
	} else {
		// 继续按权重随机强化升级

		// 限制器：只能强化晶石已有的副属性
		limiter := func(item random.Item) bool {
			if times, ok := attType[item.GetId()]; ok {
				// 同一条副属性最多只能随机到n次
				return times < int(globalConfig.CrystalLevelupAssistantNumber)
			}
			return false
		}

		viceAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, define.Crystal_AttTypeVice)
		it, err := random.PickOne(viceAttRepoList, limiter)
		if !utils.ErrCheck(err, "pick one vice att failed", c.Id) {
			return
		}

		viceAttRepoEntry = it.(*auto.CrystalAttRepoEntry)
	}

	if viceAttRepoEntry == nil {
		log.Error().Int64("player_id", m.owner.ID).Msg("enforceCrystalViceAtt failed")
		return
	}

	// 添加副属性
	c.ViceAtts = append(c.ViceAtts, item.CrystalAtt{
		AttRepoId:    viceAttRepoEntry.Id,
		AttRandRatio: random.Decimal(globalConfig.CrystalLevelupRandRatio[0], globalConfig.CrystalLevelupRandRatio[1]),
	})
}

// 晶石升级
func (m *ItemManager) CrystalLevelup(crystalId int64, stuffItems, expItems []int64) error {
	it, err := m.GetItem(crystalId)
	utils.ErrPrint(err, "CrystalLevelup failed", crystalId, m.owner.ID)

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return auto.ErrGlobalConfigInvalid
	}

	if it.GetType() != define.Item_TypeCrystal {
		return fmt.Errorf("CrystalLevelup failed, wrong item<%d> type", it.Opts().TypeId)
	}

	c := it.(*item.Crystal)
	_, ok = auto.GetCrystalLevelupEntry(int32(c.Level) + 1)
	if !ok {
		return fmt.Errorf("CyrstalLevelup failed, cannot find crystal levelup entry<%d>", c.Level+1)
	}

	// 品质限制等级上限
	if int32(c.Level) >= globalConfig.CrystalLevelupQualityLimit[c.ItemEntry.Quality] {
		return errors.New("crystal quality limit")
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

		if it.Opts().ItemEntry.Type != define.Item_TypeCrystal {
			continue
		}

		// 重复的id不计入
		if _, ok := unrepeatedItemId[id]; ok {
			continue
		}

		stuffCrystal := it.(*item.Crystal)

		// 1级经验不算折损率
		crystalLv1Entry, ok := auto.GetCrystalLevelupEntry(1)
		if !ok {
			log.Error().Caller().Msg("can not find crystal levelup 1 entry")
			continue
		}
		crystalLv1Exp := crystalLv1Entry.Exp[stuffCrystal.ItemEntry.Quality]

		// 已升级累计的经验
		crystalLvEntry, ok := auto.GetCrystalLevelupEntry(int32(stuffCrystal.Level))
		if !ok {
			log.Error().Caller().Int8("level", stuffCrystal.Level).Msg("can not find crystal levelup entry")
			continue
		}
		crystallvTotalExp := crystalLvEntry.Exp[stuffCrystal.ItemEntry.Quality] + stuffCrystal.Exp - crystalLv1Exp

		// 物品总经验 = 物品1级经验 + 已消耗所有经验 * 经验折损率
		itemExps[it] = int32(globalConfig.CrystalSwallowExpLoss.Mul(decimal.NewFromInt32(crystallvTotalExp)).Round(0).IntPart()) + crystalLv1Exp

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

		if it.Opts().ItemEntry.SubType != define.Item_SubType_Item_CrystalExp {
			continue
		}

		if _, ok := unrepeatedItemId[id]; ok {
			continue
		}

		itemExps[it] = it.Opts().ItemEntry.PublicMisc[0]
		unrepeatedItemId[id] = struct{}{}
	}

	// 状态改变
	changed := false

	// 升级处理
	levelupFn := func(itemId int64, exp int32) bool {
		_, ok := auto.GetCrystalLevelupEntry(int32(c.Level) + 1)
		if !ok {
			return false
		}

		// 判断金币
		costGold := int32(int64(exp) * int64(globalConfig.CrystalLevelupExpGoldRatio))
		if costGold < 0 {
			return false
		}

		// 品质限制等级上限
		if int32(c.Level) >= globalConfig.CrystalLevelupQualityLimit[c.ItemEntry.Quality] {
			return false
		}

		if err := m.owner.TokenManager().CanCost(define.Token_Gold, costGold); err != nil {
			return false
		}

		// overflow
		if c.Exp+exp < 0 {
			return false
		}

		c.Exp += exp
		changed = true
		reachLimit := false
		for {
			curLevelEntry, _ := auto.GetCrystalLevelupEntry(int32(c.Level))
			nextLevelEntry, ok := auto.GetCrystalLevelupEntry(int32(c.Level) + 1)
			if !ok {
				reachLimit = true
				break
			}

			// 品质限制等级上限
			if int32(c.Level) >= globalConfig.CrystalLevelupQualityLimit[c.ItemEntry.Quality] {
				reachLimit = true
				break
			}

			levelExp := nextLevelEntry.Exp[c.ItemEntry.Quality] - curLevelEntry.Exp[c.ItemEntry.Quality]
			if c.Exp < levelExp {
				break
			}

			c.Level++
			c.Exp -= levelExp
			for _, level := range globalConfig.CrystalViceAttAddLevel {
				if int32(c.Level) == level {
					// 增加新的副属性直到满4条
					m.generateCrystalViceAtt(c)

					// 强化副属性
					m.enforceCrystalViceAtt(c)
					// c.GetAttManager().CalcAtt()
					// m.SendCrystalAttUpdate(c)
					break
				}
			}
		}

		// 消耗材料
		err = m.CostItemByID(itemId, 1)
		utils.ErrPrint(err, "ItemManager CostItemByID failed", itemId)

		// 消耗金币
		err = m.owner.TokenManager().DoCost(define.Token_Gold, costGold)
		utils.ErrPrint(err, "TokenManager DoCost failed", costGold)

		// 返还处理
		if reachLimit && c.Exp > 0 {
			exp := c.Exp
			c.Exp = 0

			for {
				if exp <= 0 {
					break
				}

				// 没有可补的道具了
				expItem := globalConfig.GetCrystalExpItemByExp(exp)
				if expItem == nil {
					break
				}

				err := m.owner.ItemManager().GainLoot(expItem.ItemTypeId, exp/expItem.Exp)
				utils.ErrPrint(err, "gain loot failed when crystal levelup return exp items", exp, expItem.ItemTypeId)

				returnToken := exp / expItem.Exp * expItem.Exp * globalConfig.CrystalLevelupExpGoldRatio
				err = m.owner.TokenManager().GainLoot(define.Token_Gold, returnToken)
				utils.ErrPrint(err, "gain loot failed when crystal levelup return exp items", exp, returnToken)

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
	err = store.GetStore().UpdateOne(context.Background(), define.StoreType_Item, c.Id, c)
	if !utils.ErrCheck(err, "UpdateOne failed when ItemManager.CrystalLevelup", m.owner.ID, c.Level, c.Exp) {
		return err
	}

	m.SendCrystalUpdate(c)
	return nil
}

// gm晶石升级
func (m *ItemManager) GmCrystalLevelup(crystalTypeId int32, level int32, exp int32) error {
	it := m.GetItemByTypeId(crystalTypeId)
	if it == nil {
		return ErrItemNotFound
	}

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return auto.ErrGlobalConfigInvalid
	}

	if it.GetType() != define.Item_TypeCrystal {
		return ErrItemInvalidType
	}

	c := it.(*item.Crystal)
	_, ok = auto.GetCrystalLevelupEntry(level)
	if !ok {
		return fmt.Errorf("CyrstalLevelup failed, cannot find crystal levelup entry<%d>", c.Level+1)
	}

	// 品质限制等级上限
	if level >= globalConfig.CrystalLevelupQualityLimit[c.ItemEntry.Quality] {
		return errors.New("crystal quality limit")
	}

	if level < 0 {
		level = int32(c.Level)
	}

	if exp < 0 {
		exp = c.Exp
	}

	// 属性生成
	for n := c.Level; n < int8(level); n++ {
		for _, lv := range globalConfig.CrystalViceAttAddLevel {
			if int32(n) == lv {
				// 增加新的副属性直到满4条
				m.generateCrystalViceAtt(c)

				// 强化副属性
				m.enforceCrystalViceAtt(c)
				break
			}
		}
	}

	c.Level = int8(level)
	c.Exp = exp

	// save
	err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Item, c.Id, c)
	if !utils.ErrCheck(err, "UpdateOne failed when ItemManager.GmCrystalLevelup", m.owner.ID, c.Level, c.Exp) {
		return err
	}

	m.SendCrystalUpdate(c)
	return nil
}

// 测试接口，不得用于正常逻辑
func (m *ItemManager) CrystalBulkRandom(num int32) error {
	itemRows := auto.GetItemRows()
	crystalEntries := make([]*auto.ItemEntry, 0, 100)
	for _, entry := range itemRows {
		if entry.Type == define.Item_TypeCrystal {
			crystalEntries = append(crystalEntries, entry)
		}
	}

	globalConfig, _ := auto.GetGlobalConfig()
	generatedCrystals := make([]*item.Crystal, 0, num)

	var n int32
	for n = 0; n < num; n++ {
		entry := crystalEntries[rand.Intn(len(crystalEntries))]
		it := m.createItem(entry.Id, 1)
		if it == nil {
			log.Error().Caller().Int32("type_id", entry.Id).Msg("createItem failed")
			continue
		}

		crystal := it.(*item.Crystal)
		crystal.Level = 15
		for range globalConfig.CrystalViceAttAddLevel {
			// 增加新的副属性直到满4条
			m.generateCrystalViceAtt(crystal)

			// 强化副属性
			m.enforceCrystalViceAtt(crystal)
		}

		generatedCrystals = append(generatedCrystals, crystal)

		err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Item, crystal.Id, crystal)
		utils.ErrPrint(err, "UpdateOne failed when CrystalBulkRandom", it.Opts().TypeId, m.owner.ID)
	}
	log.Info().Int32("num", num).Msg("CrystalBulkRandom success")

	for _, c := range generatedCrystals {
		mapViceAtts := make(map[int32]int32)
		for _, att := range c.ViceAtts {
			mapViceAtts[att.AttRepoId]++
		}

		event := log.Info()
		event.Int32("晶石id", c.TypeId)
		attString := make([]string, 0, 10)
		for attRepoId, num := range mapViceAtts {
			entry, _ := auto.GetCrystalAttRepoEntry(attRepoId)
			attString = append(attString, fmt.Sprintf("%s, 出现次数:%d", entry.Desc, num))
		}
		event.Strs("副属性", attString).Send()
	}

	msg := &pbGlobal.S2C_TestCrystalRandomReport{}

	// mapMainAttRepo := make(map[int32]int32)
	// mapViceAttRepo := make(map[int32]int32)
	// for _, c := range generatedCrystals {
	// 	for _, att := range c.ViceAtts {
	// 		mapViceAttRepo[att.AttRepoId]++
	// 	}

	// 	mapMainAttRepo[c.MainAtt.AttRepoId]++
	// }

	// msg.Report = make([]string, 0, 100)

	// // 主属性统计
	// var mainNum int32
	// for repoId, num := range mapMainAttRepo {
	// 	attRepoEntry, ok := auto.GetCrystalAttRepoEntry(repoId)
	// 	if !ok {
	// 		continue
	// 	}

	// 	report := fmt.Sprintf("主属性描述<%s> att_id<%d> 权重<%d> 出现次数<%d>", attRepoEntry.Desc, attRepoEntry.AttId, attRepoEntry.AttWeight, num)
	// 	msg.Report = append(msg.Report, report)
	// 	mainNum += num
	// }

	// msg.Report = append(msg.Report, fmt.Sprintf("总主属性条数<%d>", mainNum))

	// // 副属性统计
	// var viceNum int32
	// for repoId, num := range mapViceAttRepo {
	// 	attRepoEntry, ok := auto.GetCrystalAttRepoEntry(repoId)
	// 	if !ok {
	// 		continue
	// 	}

	// 	report := fmt.Sprintf("副属性描述<%s> att_id<%d> 权重<%d> 出现次数<%d>", attRepoEntry.Desc, attRepoEntry.AttId, attRepoEntry.AttWeight, num)
	// 	msg.Report = append(msg.Report, report)
	// 	viceNum += num
	// }

	// msg.Report = append(msg.Report, fmt.Sprintf("总副属性条数<%d>", viceNum))

	m.owner.SendProtoMessage(msg)

	return nil
}

func (m *ItemManager) SaveCrystalEquiped(c *item.Crystal) {
	fields := map[string]interface{}{
		"crystal_obj": c.CrystalObj,
	}

	err := store.GetStore().UpdateFields(context.Background(), define.StoreType_Item, c.Id, fields)
	utils.ErrPrint(err, "UpdateArray failed when ItemManager.SaveCrystalEquiped failed", c.Id)
}

func (m *ItemManager) SendCrystalAttUpdate(c *item.Crystal) {
	msg := &pbGlobal.S2C_CrystalAttUpdate{
		CrystalId: c.Id,
		AttValue:  make([]int32, define.Att_End),
	}

	for n := 0; n < define.Att_End; n++ {
		msg.AttValue[n] = int32(c.GetAttManager().GetFinalAttValue(n).Round(0).IntPart())
	}

	m.owner.SendProtoMessage(msg)
}

func (m *ItemManager) SendCrystalUpdate(c *item.Crystal) {
	msg := &pbGlobal.S2C_CrystalUpdate{
		CrystalId: c.Id,
		CrystalData: &pbGlobal.CrystalData{
			Level:      int32(c.Level),
			Exp:        c.Exp,
			CrystalObj: c.CrystalObj,
			MainAtt: &pbGlobal.CrystalAtt{
				AttRepoId:    c.MainAtt.AttRepoId,
				AttRandRatio: int32(c.MainAtt.AttRandRatio.Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()),
			},
			ViceAtts: make([]*pbGlobal.CrystalAtt, len(c.ViceAtts)),
		},
	}

	for n, att := range c.ViceAtts {
		msg.CrystalData.ViceAtts[n] = &pbGlobal.CrystalAtt{
			AttRepoId:    att.AttRepoId,
			AttRandRatio: int32(att.AttRandRatio.Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()),
		}
	}

	m.owner.SendProtoMessage(msg)
}

func (m *ItemManager) GenCrystalListPB() []*pbGlobal.Crystal {
	crystals := make([]*pbGlobal.Crystal, 0, m.GetItemNums(int(define.Container_Crystal)))
	m.ca.RangeByIdx(int(define.Container_Crystal), func(val interface{}) bool {
		it, ok := val.(*item.Crystal)
		if !ok {
			return true
		}

		crystals = append(crystals, it.GenCrystalPB())
		return true
	})

	return crystals
}
