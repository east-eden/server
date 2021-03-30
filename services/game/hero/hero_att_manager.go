package hero

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/internal/att"
	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog/log"
)

// 英雄属性计算管理
type HeroAttManager struct {
	hero *Hero
	att.AttManager
}

func NewHeroAttManager(hero *Hero) *HeroAttManager {
	m := &HeroAttManager{
		hero: hero,
	}

	return m
}

// 计算英雄属性
func (m *HeroAttManager) CalcAtt() {
	m.Reset()

	// 升级
	m.CalcLevelup()

	// 突破
	m.CalcPromote()

	// 装备
	m.CalcEquipBar()

	// 晶石
	m.CalcCrystalBox()

	// 计算最终值
	m.AttManager.CalcAtt()
}

//////////////////////////////////////////////
// 升级属性 =（卡牌等级*各升级属性成长率+各升级属性固定值）*卡牌品质参数
func (m *HeroAttManager) CalcLevelup() {
	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		log.Error().Caller().Msg("invalid global config")
		return
	}

	// 成长率att
	attGrowRatio := att.NewAttManager()
	attGrowRatio.SetBaseAttId(globalConfig.HeroLevelGrowRatioAttId)

	for n := define.Att_Begin; n < define.Att_End; n++ {
		growRatioBase := attGrowRatio.GetAttValue(n)
		if growRatioBase != 0 {
			// 等级*升级成长率
			add := growRatioBase * int32(m.hero.Level)

			// 品质参数
			qualityRatio := globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)]

			value64 := float64(add) * float64(qualityRatio/define.PercentBase)
			value := int32(utils.Round(value64))
			if value < 0 {
				log.Error().
					Caller().
					Int("att_enum", n).
					Int32("att_value", value).
					Int64("hero_id", m.hero.Id).
					Msg("hero CalcLevelup overflow")
			}

			m.ModAttValue(n, value)
		}
	}
}

//////////////////////////////////////////////
// 突破属性 =（突破强度等级*各升级属性成长率+各升级属性固定值）*卡牌品质参数*卡牌职业参数
func (m *HeroAttManager) CalcPromote() {
	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		log.Error().Caller().Msg("invalid global config")
		return
	}

	// 成长率att
	attGrowRatio := att.NewAttManager()
	attGrowRatio.SetBaseAttId(globalConfig.HeroPromoteGrowupId)

	// 基础att
	promoteBaseAtt := att.NewAttManager()
	promoteBaseAtt.SetBaseAttId(globalConfig.HeroPromoteBaseId)

	// 职业参数
	professionEntry, ok := auto.GetHeroProfessionEntry(m.hero.Entry.Profession)
	if !ok {
		log.Error().Caller().Int32("profession", m.hero.Entry.Profession).Msg("can not find HeroProfessionEntry")
		return
	}

	for n := define.Att_Begin; n < define.Att_End; n++ {
		growRatioBase := attGrowRatio.GetAttValue(n)
		promoteBase := promoteBaseAtt.GetAttValue(n)
		if promoteBase != 0 && growRatioBase != 0 {
			// 强度等级*升级成长率 + 基础值
			add := growRatioBase*globalConfig.HeroPromoteIntensityRatio[m.hero.PromoteLevel] + promoteBase

			// 品质参数
			qualityRatio := globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)]

			// 职业参数
			professionRatio := professionEntry.GetRatio(n)

			value64 := float64(add) * float64(qualityRatio/define.PercentBase) * float64(professionRatio/define.PercentBase)
			value := int32(utils.Round(value64))
			if value < 0 {
				log.Error().
					Caller().
					Int("att_enum", n).
					Int32("att_value", value).
					Int64("hero_id", m.hero.Id).
					Msg("hero CalcPromote overflow")
			}

			m.ModAttValue(n, value)
		}
	}
}

//////////////////////////////////////////////
// 计算所有装备栏属性
func (m *HeroAttManager) CalcEquipBar() {
	var n int32
	for n = 0; n < int32(define.Equip_Pos_End); n++ {
		e := m.hero.equipBar.GetEquipByPos(n)
		if e == nil {
			continue
		}

		e.GetAttManager().Reset()
		e.GetAttManager().CalcAtt()
		m.ModAttManager(&e.GetAttManager().AttManager)
	}
}

//////////////////////////////////////////////
// 计算所有晶石属性
func (m *HeroAttManager) CalcCrystalBox() {
	// crystal box
	var n int32
	for n = 0; n < define.Crystal_PosEnd; n++ {
		c := m.hero.crystalBox.GetCrystalByPos(n)
		if c == nil {
			continue
		}

		c.GetAttManager().Reset()
		c.GetAttManager().CalcAtt()
		m.ModAttManager(&c.GetAttManager().AttManager)
	}
}
