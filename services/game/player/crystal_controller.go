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
	"bitbucket.org/funplus/server/utils/random"
	log "github.com/rs/zerolog/log"
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
	c.MainAtt.AttRandRatio = random.Int32(int32(globalConfig.CrystalLevelupRandRatio[0]), int32(globalConfig.CrystalLevelupRandRatio[1]))

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
		log.Error().Err(err).Int64("crystal_id", c.Id).Msg("pick unrepeated crystal vice att failed")
		return
	}

	for _, v := range viceAttItems {
		viceAttRepoEntry := v.(*auto.CrystalAttRepoEntry)
		c.ViceAtts = append(c.ViceAtts, item.CrystalAtt{
			AttRepoId:    viceAttRepoEntry.Id,
			AttRandRatio: random.Int32(int32(globalConfig.CrystalLevelupRandRatio[0]), int32(globalConfig.CrystalLevelupRandRatio[1])),
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
	if pass := utils.ErrCheck(err, "pick one vice att failed", c.Id); !pass {
		return
	}

	attRepoEntry := it.(*auto.CrystalAttRepoEntry)
	c.ViceAtts = append(c.ViceAtts, item.CrystalAtt{
		AttRepoId:    attRepoEntry.Id,
		AttRandRatio: random.Int32(int32(globalConfig.CrystalLevelupRandRatio[0]), int32(globalConfig.CrystalLevelupRandRatio[1])),
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

	// 限制器：只能强化晶石已有的副属性
	limiter := func(item random.Item) bool {
		if times, ok := attType[item.GetId()]; ok {
			// 同一条副属性最多只能随机到n次
			return times <= int(globalConfig.CrystalLevelupAssistantNumber)
		}
		return false
	}

	viceAttRepoList := auto.GetCrystalAttRepoList(c.CrystalEntry.Pos, define.Crystal_AttTypeVice)
	it, err := random.PickOne(viceAttRepoList, limiter)
	if pass := utils.ErrCheck(err, "pick one vice att failed", c.Id); !pass {
		return
	}

	viceAttRepoEntry := it.(*auto.CrystalAttRepoEntry)
	c.ViceAtts = append(c.ViceAtts, item.CrystalAtt{
		AttRepoId:    viceAttRepoEntry.Id,
		AttRandRatio: random.Int32(int32(globalConfig.CrystalLevelupRandRatio[0]), int32(globalConfig.CrystalLevelupRandRatio[1])),
	})
}

// 晶石升级
func (m *ItemManager) CrystalLevelup(crystalId int64, stuffItems, expItems []int64) error {
	it, err := m.GetItem(crystalId)
	utils.ErrPrint(err, "CrystalLevelup failed", crystalId, m.owner.ID)

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return errors.New("invalid global config")
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
		if pass := utils.ErrCheck(err, "cannot find item", id); !pass {
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
		itemExps[it] = int32(int64(crystalLv1Exp) + int64(crystallvTotalExp)*int64(globalConfig.CrystalSwallowExpLoss)/int64(define.PercentBase))
		unrepeatedItemId[id] = struct{}{}
	}

	// 经验道具
	for _, id := range expItems {
		it, err := m.GetItem(id)
		if pass := utils.ErrCheck(err, "cannot find item", id); !pass {
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
	fields := map[string]interface{}{
		MakeItemKey(c): c,
	}
	err = store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
	if pass := utils.ErrCheck(err, "CrystalLevelup SaveFields failed", m.owner.ID, c.Level, c.Exp); !pass {
		return err
	}

	m.SendCrystalUpdate(c)
	return nil
}

func (m *ItemManager) SaveCrystalEquiped(c *item.Crystal) {
	fields := map[string]interface{}{
		MakeItemKey(c, "crystal_obj"): c.CrystalObj,
	}

	err := store.GetStore().SaveFields(define.StoreType_Item, m.owner.ID, fields)
	utils.ErrPrint(err, "SaveCrystalEquiped failed", c.Id)
}

func (m *ItemManager) SendCrystalAttUpdate(c *item.Crystal) {
	msg := &pbGlobal.S2C_CrystalAttUpdate{
		CrystalId: c.Id,
		AttValue:  make([]int32, define.Att_End),
	}

	for n := 0; n < define.Att_End; n++ {
		msg.AttValue[n] = c.GetAttManager().GetAttValue(n)
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
				AttRandRatio: c.MainAtt.AttRandRatio,
			},
			ViceAtts: make([]*pbGlobal.CrystalAtt, len(c.ViceAtts)),
		},
	}

	for n, att := range c.ViceAtts {
		msg.CrystalData.ViceAtts[n] = &pbGlobal.CrystalAtt{
			AttRepoId:    att.AttRepoId,
			AttRandRatio: att.AttRandRatio,
		}
	}

	m.owner.SendProtoMessage(msg)
}
