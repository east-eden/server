package player

import (
	"errors"
	"fmt"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	pbCombat "bitbucket.org/funplus/server/proto/server/combat"
	"bitbucket.org/funplus/server/services/game/hero"
	"bitbucket.org/funplus/server/services/game/item"
	"bitbucket.org/funplus/server/services/game/prom"
	"bitbucket.org/funplus/server/store"
	"bitbucket.org/funplus/server/utils"
	json "github.com/json-iterator/go"
	log "github.com/rs/zerolog/log"
)

var (
	ErrHeroNotFound = errors.New("hero not found")
)

type HeroManager struct {
	define.BaseCostLooter `bson:"-" json:"-"`

	owner       *Player              `bson:"-" json:"-"`
	HeroList    map[int64]*hero.Hero `bson:"-" json:"-"` // 卡牌包
	heroTypeSet map[int32]struct{}   `bson:"-" json:"-"` // 已获得卡牌
}

func NewHeroManager(owner *Player) *HeroManager {
	m := &HeroManager{
		owner:       owner,
		HeroList:    make(map[int64]*hero.Hero),
		heroTypeSet: make(map[int32]struct{}),
	}

	return m
}

func (m *HeroManager) Destroy() {
	for _, h := range m.HeroList {
		hero.GetHeroPool().Put(h)
	}
}

func (m *HeroManager) createEntryHero(entry *auto.HeroEntry) *hero.Hero {
	if entry == nil {
		log.Error().Msg("newEntryHero with nil HeroEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Hero)
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	h := hero.NewHero()
	h.Init(
		hero.Id(id),
		hero.OwnerId(m.owner.GetID()),
		hero.OwnerType(m.owner.GetType()),
		hero.Entry(entry),
		hero.TypeId(entry.Id),
	)

	h.GetAttManager().SetBaseAttId(int32(entry.AttId))
	m.HeroList[h.GetOptions().Id] = h
	m.heroTypeSet[h.GetOptions().TypeId] = struct{}{}

	return h
}

func (m *HeroManager) initLoadedHero(h *hero.Hero) error {
	entry, ok := auto.GetHeroEntry(h.GetOptions().TypeId)
	if !ok {
		return fmt.Errorf("HeroManager initLoadedHero: hero<%d> entry invalid", h.GetOptions().TypeId)
	}

	h.GetOptions().Entry = entry
	h.GetAttManager().SetBaseAttId(int32(entry.AttId))

	m.HeroList[h.GetOptions().Id] = h
	m.heroTypeSet[h.GetOptions().TypeId] = struct{}{}

	return nil
}

// interface of cost_loot
func (m *HeroManager) GetCostLootType() int32 {
	return define.CostLoot_Hero
}

func (m *HeroManager) CanCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.CanCost(typeMisc, num)
	if err != nil {
		return err
	}

	var fixNum int32
	for _, v := range m.HeroList {
		if v.GetOptions().TypeId == typeMisc {
			eb := v.GetEquipBar()
			hasEquip := false

			var n int32
			for n = 0; n < int32(define.Equip_Pos_End); n++ {
				if eb.GetEquipByPos(n) != nil {
					hasEquip = true
					break
				}
			}

			if !hasEquip {
				fixNum++
			}
		}
	}

	if fixNum >= num {
		return nil
	}

	return fmt.Errorf("not enough hero<%d>, num<%d>", typeMisc, num)
}

func (m *HeroManager) DoCost(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.DoCost(typeMisc, num)
	if err != nil {
		return err
	}

	var costNum int32
	for _, v := range m.HeroList {
		if v.GetOptions().TypeId == typeMisc {
			eb := v.GetEquipBar()
			hasEquip := false

			var n int32
			for n = 0; n < int32(define.Equip_Pos_End); n++ {
				if eb.GetEquipByPos(n) != nil {
					hasEquip = true
					break
				}
			}

			if !hasEquip {
				m.DelHero(v.GetOptions().Id)
				costNum++
			}
		}
	}

	if costNum < num {
		log.Warn().
			Int32("cost_type_misc", typeMisc).
			Int32("cost_num", num).
			Int32("actual_cost_num", costNum).
			Msg("hero manager cost num error")
		return nil
	}

	return nil
}

func (m *HeroManager) GainLoot(typeMisc int32, num int32) error {
	err := m.BaseCostLooter.GainLoot(typeMisc, num)
	if err != nil {
		return err
	}

	var n int32
	for n = 0; n < num; n++ {
		_ = m.AddHeroByTypeId(typeMisc)
	}

	return nil
}

func (m *HeroManager) LoadAll() error {
	docs, err := store.GetStore().LoadHashAll(define.StoreType_Hero, "owner_id", m.owner.ID)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("HeroManager LoadAll: %w", err)
	}

	mm := docs.(map[string]interface{})
	for _, v := range mm {
		vv := v.([]byte)
		h := hero.NewHero()
		err := json.Unmarshal(vv, h)
		if !utils.ErrCheck(err, "json unmarshal failed", vv) {
			return err
		}

		if err := m.initLoadedHero(h); err != nil {
			return fmt.Errorf("HeroManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *HeroManager) GetHero(id int64) *hero.Hero {
	return m.HeroList[id]
}

func (m *HeroManager) GetHeroByTypeId(typeId int32) *hero.Hero {
	if _, ok := m.heroTypeSet[typeId]; !ok {
		return nil
	}

	for _, h := range m.HeroList {
		if h.Entry.Id == typeId {
			return h
		}
	}

	return nil
}

func (m *HeroManager) GetHeroNums() int {
	return len(m.HeroList)
}

func (m *HeroManager) GetHeroList() []*hero.Hero {
	list := make([]*hero.Hero, 0)

	for _, v := range m.HeroList {
		list = append(list, v)
	}

	return list
}

func (m *HeroManager) AddHeroByTypeId(typeId int32) *hero.Hero {
	heroEntry, ok := auto.GetHeroEntry(typeId)
	if !ok {
		log.Warn().Int32("type_id", typeId).Msg("GetHeroEntry failed")
		return nil
	}

	// 重复获得卡牌，转换为对应碎片
	_, ok = m.heroTypeSet[typeId]
	if ok {
		m.owner.FragmentManager().Inc(typeId, heroEntry.FragmentTransform)
		return nil
	}

	h := m.createEntryHero(heroEntry)
	if h == nil {
		log.Warn().Int32("type_id", typeId).Msg("createEntryHero failed")
		return nil
	}

	err := store.GetStore().SaveHashObject(define.StoreType_Hero, h.OwnerId, h.Id, h)
	if !utils.ErrCheck(err, "SaveObject failed when AddHeroByTypeID", typeId, m.owner.ID) {
		m.delHero(h)
		return nil
	}

	m.SendHeroUpdate(h)

	// prometheus ops
	prom.OpsCreateHeroCounter.Inc()

	return h
}

func (m *HeroManager) delHero(h *hero.Hero) {
	delete(m.HeroList, h.Options.Id)
	delete(m.heroTypeSet, h.Options.TypeId)
	hero.GetHeroPool().Put(h)
}

func (m *HeroManager) DelHero(id int64) {
	h, ok := m.HeroList[id]
	if !ok {
		return
	}

	eb := h.GetEquipBar()
	var n int32
	for n = 0; n < int32(define.Equip_Pos_End); n++ {
		err := eb.TakeoffEquip(n)
		utils.ErrPrint(err, "DelHero TakeoffEquip failed", id, n)
	}

	err := store.GetStore().DeleteHashObject(define.StoreType_Hero, h.OwnerId, h.Id)
	utils.ErrPrint(err, "DelHero DeleteObject failed", id)
	m.delHero(h)
}

func (m *HeroManager) HeroLevelup(heroId int64, stuffItems []int64) error {
	h := m.GetHero(heroId)
	if h == nil {
		return errors.New("hero not found")
	}

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		return errors.New("invalid global config")
	}

	// 经验道具
	expItems := make(map[item.Itemface]int32)

	// 剔除重复的物品
	unrepeatedItemId := make(map[int64]struct{})

	for _, id := range stuffItems {
		it, err := m.owner.ItemManager().GetItem(id)
		if err != nil {
			continue
		}

		// 判断物品类型合法
		if it.GetType() != define.Item_TypeItem {
			continue
		}

		if it.Opts().ItemEntry.SubType != define.Item_SubType_Item_HeroExp {
			continue
		}

		// 重复的id不计入
		if _, ok := unrepeatedItemId[id]; ok {
			continue
		}

		expItems[it] = it.Opts().ItemEntry.PublicMisc[0]
		unrepeatedItemId[id] = struct{}{}
	}

	// 状态是否改变
	changed := false

	// 升级处理
	levelupFn := func(itemId int64, exp int32) bool {
		nextLevelEntry, ok := auto.GetHeroLevelupEntry(int32(h.Level) + 1)
		if !ok {
			return false
		}

		// 突破限制
		if int32(h.PromoteLevel) < nextLevelEntry.PromoteLimit {
			return false
		}

		// 金币限制
		costGold := int32(int64(exp) * int64(globalConfig.HeroLevelupExpGoldRatio))
		if costGold < 0 {
			return false
		}

		if err := m.owner.TokenManager().CanCost(define.Token_Gold, costGold); err != nil {
			return false
		}

		// overflow
		if h.Exp+exp < 0 {
			return false
		}

		h.Exp += exp
		changed = true
		reachLimit := false
		for {
			curLevelEntry, _ := auto.GetHeroLevelupEntry(int32(h.Level))
			nextLevelEntry, ok := auto.GetHeroLevelupEntry(int32(h.Level) + 1)
			if !ok {
				reachLimit = true
				break
			}

			if int32(h.PromoteLevel) < nextLevelEntry.PromoteLimit {
				reachLimit = true
				break
			}

			levelExp := nextLevelEntry.Exp - curLevelEntry.Exp
			if h.Exp < levelExp {
				break
			}

			h.Level++
			h.Exp -= levelExp
		}

		// 消耗
		err := m.owner.TokenManager().DoCost(define.Token_Gold, costGold)
		utils.ErrPrint(err, "TokenManager DoCost failed", costGold)

		err = m.owner.ItemManager().CostItemByID(itemId, 1)
		utils.ErrPrint(err, "ItemManager CostItemByID failed", itemId)

		// 返还处理
		if reachLimit && h.Exp > 0 {
			exp := h.Exp
			h.Exp = 0

			for {
				if exp <= 0 {
					break
				}

				// 没有可补的道具了
				expItem := globalConfig.GetHeroExpItemByExp(exp)
				if expItem == nil {
					break
				}

				err := m.owner.ItemManager().GainLoot(expItem.ItemTypeId, exp/expItem.Exp)
				utils.ErrPrint(err, "gain loot failed when hero levelup return exp items", exp, expItem.ItemTypeId)

				returnToken := exp / expItem.Exp * expItem.Exp * globalConfig.HeroLevelupExpGoldRatio
				err = m.owner.TokenManager().GainLoot(define.Token_Gold, returnToken)
				utils.ErrPrint(err, "gain loot failed when hero levelup return exp items", exp, returnToken)

				exp %= expItem.Exp
			}
		}

		return true
	}

	continueCheck := true
	for it, exp := range expItems {
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
		"level": h.Level,
		"exp":   h.Exp,
	}
	err := store.GetStore().SaveHashObjectFields(define.StoreType_Hero, h.OwnerId, h.Id, h, fields)
	if !utils.ErrCheck(err, "HeroLevelup SaveFields failed", m.owner.ID, h.Level, h.Exp) {
		return err
	}

	m.SendHeroUpdate(h)
	return nil
}

func (m *HeroManager) HeroPromote(heroId int64) error {
	h := m.GetHero(heroId)
	if h == nil {
		return errors.New("hero not found")
	}

	if h.PromoteLevel >= define.Hero_Max_Promote_Times {
		return errors.New("promote level max")
	}

	nextLevelEntry, ok := auto.GetHeroLevelupEntry(int32(h.Level) + 1)
	if !ok {
		return errors.New("hero level max")
	}

	if int32(h.PromoteLevel) >= nextLevelEntry.PromoteLimit {
		return errors.New("hero levelup max, then promote")
	}

	promoteEntry, ok := auto.GetHeroPromoteEntry(h.TypeId)
	if !ok {
		return errors.New("hero promote entry not found")
	}

	costId := promoteEntry.PromoteCostId[h.PromoteLevel+1]
	err := m.owner.CostLootManager().CanCost(costId)
	if err != nil {
		return err
	}

	err = m.owner.CostLootManager().DoCost(costId)
	if !utils.ErrCheck(err, "HeroPromote failed", heroId, costId) {
		return err
	}

	h.PromoteLevel++

	// save
	fields := map[string]interface{}{
		"promote_level": h.PromoteLevel,
	}
	err = store.GetStore().SaveHashObjectFields(define.StoreType_Hero, h.OwnerId, h.Id, h, fields)
	if !utils.ErrCheck(err, "HeroPromote SaveFields failed", m.owner.ID, h.PromoteLevel) {
		return err
	}

	m.SendHeroUpdate(h)
	return nil
}

func (m *HeroManager) PutonEquip(heroId int64, equipId int64) error {
	it, err := m.owner.ItemManager().GetItem(equipId)
	if err != nil {
		return fmt.Errorf("HeroManager.PutonEquip failed: %w", err)
	}

	if it.GetType() != define.Item_TypeEquip {
		return fmt.Errorf("item<%d> is not an equip when PutonEquip", equipId)
	}

	equip := it.(*item.Equip)
	if objId := equip.GetEquipObj(); objId != -1 {
		return fmt.Errorf("equip has put on another hero<%d>", objId)
	}

	if equip.EquipEnchantEntry == nil {
		return fmt.Errorf("cannot find equip_enchant_entry<%d> while PutonEquip", equipId)
	}

	h, ok := m.HeroList[heroId]
	if !ok {
		return errors.New("invalid heroid")
	}

	equipBar := h.GetEquipBar()
	pos := equip.EquipEnchantEntry.EquipPos

	// 英雄能否装备这件武器
	if equip.EquipEnchantEntry.EquipPos == define.Equip_Pos_Weapon &&
		equip.ItemEntry.SubType != h.Entry.WeaponType {
		return errors.New("cannot equip this weapon type")
	}

	// takeoff previous equip
	if pe := equipBar.GetEquipByPos(pos); pe != nil {
		if err := m.TakeoffEquip(heroId, pos); err != nil {
			return err
		}
	}

	// puton this equip
	if err := equipBar.PutonEquip(equip); err != nil {
		return err
	}

	err = m.owner.ItemManager().Save(equip.Opts().Id)
	utils.ErrPrint(err, "PutonEquip Save item failed", equip.Opts().Id)

	m.owner.ItemManager().SendItemUpdate(equip)
	m.SendHeroUpdate(h)

	// att
	equip.GetAttManager().CalcAtt()
	h.GetAttManager().ModAttManager(&equip.GetAttManager().AttManager)
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) TakeoffEquip(heroId int64, pos int32) error {
	if !utils.Between(int(pos), int(define.Equip_Pos_Begin), int(define.Equip_Pos_End)) {
		return fmt.Errorf("invalid pos<%d>", pos)
	}

	h, ok := m.HeroList[heroId]
	if !ok {
		return fmt.Errorf("invalid heroid")
	}

	equipBar := h.GetEquipBar()
	equip := equipBar.GetEquipByPos(pos)
	if equip == nil {
		return fmt.Errorf("cannot find hero<%d> equip by pos<%d> while TakeoffEquip", heroId, pos)
	}

	if objId := equip.GetEquipObj(); objId == -1 {
		return fmt.Errorf("equip<%d> didn't put on this hero<%d> ", equip.Opts().Id, heroId)
	}

	// unequip
	if err := equipBar.TakeoffEquip(pos); err != nil {
		return err
	}

	err := m.owner.ItemManager().Save(equip.Opts().Id)
	utils.ErrPrint(err, "TakeoffEquip Save item failed", equip.Opts().Id)
	m.owner.ItemManager().SendItemUpdate(equip)
	m.SendHeroUpdate(h)

	// att
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

func (m *HeroManager) PutonCrystal(heroId int64, crystalId int64) error {

	i, err := m.owner.ItemManager().GetItem(crystalId)
	if !utils.ErrCheck(err, "PutonCrystal failed", crystalId, m.owner.ID) {
		return err
	}

	if i.GetType() != define.Item_TypeCrystal {
		err := errors.New("item type isn't crystal")
		log.Error().Err(err).Caller().Msg("PutonCrystal failed")
		return err
	}

	c := i.(*item.Crystal)
	if objId := c.CrystalObj; objId != -1 {
		return fmt.Errorf("crystal has put on another obj<%d>", objId)
	}

	pos := c.CrystalEntry.Pos
	if pos < define.Crystal_PosBegin || pos >= define.Crystal_PosEnd {
		return fmt.Errorf("invalid pos<%d>", pos)
	}

	h, ok := m.HeroList[heroId]
	if !ok {
		return fmt.Errorf("invalid heroid<%d>", heroId)
	}

	crystalBox := h.GetCrystalBox()

	// takeoff previous crystal
	if pr := crystalBox.GetCrystalByPos(pos); pr != nil {
		if err := m.TakeoffCrystal(heroId, pos); err != nil {
			return err
		}
	}

	// equip new crystal
	if err := crystalBox.PutonCrystal(c); err != nil {
		return err
	}

	m.owner.ItemManager().SaveCrystalEquiped(c)
	m.owner.ItemManager().SendCrystalUpdate(c)
	m.SendHeroUpdate(h)

	// att
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return err
}

func (m *HeroManager) TakeoffCrystal(heroId int64, pos int32) error {
	if pos < 0 || pos >= define.Crystal_PosEnd {
		return fmt.Errorf("invalid pos<%d>", pos)
	}

	h, ok := m.HeroList[heroId]
	if !ok {
		return fmt.Errorf("invalid heroid<%d>", heroId)
	}

	c := h.GetCrystalBox().GetCrystalByPos(pos)
	if c == nil {
		return fmt.Errorf("cannot find crystal from hero<%d>'s crystalbox pos<%d> while TakeoffCrystal", heroId, pos)
	}

	// unequip
	if err := h.GetCrystalBox().TakeoffCrystal(pos); err != nil {
		return err
	}

	m.owner.ItemManager().SaveCrystalEquiped(c)
	m.owner.ItemManager().SendCrystalUpdate(c)
	m.SendHeroUpdate(h)

	// att
	h.GetAttManager().CalcAtt()
	m.SendHeroAtt(h)

	return nil
}

// gm 改变经验
func (m *HeroManager) GmExpChange(heroId int64, exp int32) error {
	h := m.GetHero(heroId)
	if h == nil {
		return ErrHeroNotFound
	}

	// 升级处理
	h.Exp += exp
	for {
		curLevelEntry, _ := auto.GetHeroLevelupEntry(int32(h.Level))
		nextLevelEntry, ok := auto.GetHeroLevelupEntry(int32(h.Level) + 1)
		if !ok {
			break
		}

		levelExp := nextLevelEntry.Exp - curLevelEntry.Exp
		if h.Exp < levelExp {
			break
		}

		h.Level++
		h.Exp -= levelExp
	}

	m.SendHeroUpdate(h)

	// save
	fields := map[string]interface{}{
		"level": h.Level,
		"exp":   h.Exp,
	}
	return store.GetStore().SaveHashObjectFields(define.StoreType_Hero, h.OwnerId, h.Id, h, fields)
}

// gm 改变等级
func (m *HeroManager) GmLevelChange(heroId int64, level int32) error {
	h := m.GetHero(heroId)
	if h == nil {
		return ErrHeroNotFound
	}

	h.Level += int16(level)
	m.SendHeroUpdate(h)

	// save
	fields := map[string]interface{}{
		"level": h.Level,
		"exp":   h.Exp,
	}
	return store.GetStore().SaveHashObjectFields(define.StoreType_Hero, h.OwnerId, h.Id, h, fields)
}

// gm 突破
func (m *HeroManager) GmPromoteChange(heroId int64, promote int32) error {
	h := m.GetHero(heroId)
	if h == nil {
		return ErrHeroNotFound
	}

	h.PromoteLevel += int8(promote)
	m.SendHeroUpdate(h)

	// save
	fields := map[string]interface{}{
		"promote_level": h.PromoteLevel,
	}
	return store.GetStore().SaveHashObjectFields(define.StoreType_Hero, h.OwnerId, h.Id, h, fields)
}

func (m *HeroManager) GenerateCombatUnitInfo() []*pbCombat.UnitInfo {
	retList := make([]*pbCombat.UnitInfo, 0)

	list := m.GetHeroList()
	for _, hero := range list {
		unitInfo := &pbCombat.UnitInfo{
			UnitTypeId:   int32(hero.GetOptions().TypeId),
			UnitAttValue: make([]int32, define.Att_End),
		}

		for n := define.Att_Begin; n < define.Att_End; n++ {
			unitInfo.UnitAttValue[n] = hero.GetAttManager().GetAttValue(n)
		}

		retList = append(retList, unitInfo)
	}

	return retList
}

func (m *HeroManager) SendHeroUpdate(h *hero.Hero) {
	// send equips update
	reply := &pbGlobal.S2C_HeroInfo{
		Info: &pbGlobal.Hero{
			Id:             h.GetOptions().Id,
			TypeId:         int32(h.GetOptions().TypeId),
			Exp:            h.GetOptions().Exp,
			Level:          int32(h.GetOptions().Level),
			PromoteLevel:   int32(h.GetOptions().PromoteLevel),
			Star:           int32(h.GetOptions().Star),
			NormalSpellId:  h.GetOptions().NormalSpellId,
			SpecialSpellId: h.GetOptions().SpecialSpellId,
			RageSpellId:    h.GetOptions().RageSpellId,
			Friendship:     h.GetOptions().Friendship,
			FashionId:      h.GetOptions().FashionId,
		},
	}

	// equip list
	// eb := h.GetEquipBar()
	// var n int32
	// for n = 0; n < define.Equip_Pos_End; n++ {
	// 	var equipId int64 = -1
	// 	if i := eb.GetEquipByPos(n); i != nil {
	// 		equipId = i.GetOptions().Id
	// 	}

	// 	reply.Info.EquipList = append(reply.Info.EquipList, equipId)
	// }

	// crystal list
	// var pos int32
	// for pos = 0; pos < define.Crystal_PositionEnd; pos++ {
	// 	var crystalId int64 = -1
	// 	if r := h.GetCrystalBox().GetCrystalByPos(pos); r != nil {
	// 		crystalId = r.GetOptions().Id
	// 	}

	// 	reply.Info.CrystalList = append(reply.Info.CrystalList, crystalId)
	// }

	m.owner.SendProtoMessage(reply)
}

func (m *HeroManager) SendHeroAtt(h *hero.Hero) {
	attManager := h.GetAttManager()
	reply := &pbGlobal.S2C_HeroAttUpdate{
		HeroId:   h.GetOptions().Id,
		AttValue: make([]int32, define.Att_End),
	}

	for n := 0; n < define.Att_End; n++ {
		reply.AttValue[n] = attManager.GetAttValue(n)
	}

	m.owner.SendProtoMessage(reply)
}

func (m *HeroManager) GenHeroListPB() []*pbGlobal.Hero {
	heros := make([]*pbGlobal.Hero, 0, len(m.HeroList))
	for _, h := range m.HeroList {
		heros = append(heros, h.GenHeroPB())
	}

	return heros
}
