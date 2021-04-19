package hero

import (
	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/internal/att"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
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
		log.Error().Caller().Err(auto.ErrGlobalConfigInvalid).Send()
		return
	}

	// 成长率att
	attGrowRatio := att.NewAttManager()
	attGrowRatio.SetBaseAttId(globalConfig.HeroLevelGrowRatioAttId)

	// 品质参数
	qualityRatio := globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)]

	// 等级*升级成长率*品质参数
	attGrowRatio.Mul(decimal.NewFromInt32(int32(m.hero.Level)).Mul(qualityRatio))
	attGrowRatio.Round()
	m.ModAttManager(attGrowRatio)
}

//////////////////////////////////////////////
// 突破属性 =（突破强度等级*各升级属性成长率+各升级属性固定值）*卡牌品质参数*卡牌职业参数
func (m *HeroAttManager) CalcPromote() {
	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		log.Error().Caller().Err(auto.ErrGlobalConfigInvalid).Send()
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

	// 基础值
	for n := define.Att_Base_Begin; n < define.Att_Base_End; n++ {
		growRatioBase := attGrowRatio.GetBaseAttValue(n)
		promoteBase := promoteBaseAtt.GetBaseAttValue(n)
		if !promoteBase.Equal(decimal.NewFromInt32(0)) && !growRatioBase.Equal(decimal.NewFromInt32(0)) {
			// 强度等级*升级成长率 + 基础值
			add := growRatioBase.Mul(decimal.NewFromInt32(globalConfig.HeroPromoteIntensityRatio[m.hero.PromoteLevel])).Add(promoteBase)

			// 品质参数
			qualityRatio := globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)]

			// 职业参数
			professionRatio := professionEntry.GetBaseRatio(n)

			value := add.Mul(qualityRatio).Mul(professionRatio).Round(0)
			m.ModBaseAttValue(n, value)
		}
	}

	// 百分比值
	for n := define.Att_Percent_Begin; n < define.Att_Percent_End; n++ {
		growRatioPercent := attGrowRatio.GetPercentAttValue(n)
		promotePercent := promoteBaseAtt.GetPercentAttValue(n)
		if !promotePercent.Equal(decimal.NewFromInt32(0)) && !growRatioPercent.Equal(decimal.NewFromInt32(0)) {
			// 强度等级*升级成长率 + 基础值
			add := growRatioPercent.Mul(decimal.NewFromInt32(globalConfig.HeroPromoteIntensityRatio[m.hero.PromoteLevel])).Add(promotePercent)

			// 品质参数
			qualityRatio := globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)]

			// 职业参数
			professionRatio := professionEntry.GetPercentRatio(n)

			value := add.Mul(qualityRatio).Mul(professionRatio).Round(4)
			m.ModPercentAttValue(n, value)
		}
	}

	// 最终值
	for n := define.Att_Begin; n < define.Att_End; n++ {
		growRatioFinal := attGrowRatio.GetFinalAttValue(n)
		promoteFinal := promoteBaseAtt.GetFinalAttValue(n)
		if !promoteFinal.Equal(decimal.NewFromInt32(0)) && !growRatioFinal.Equal(decimal.NewFromInt32(0)) {
			// 强度等级*升级成长率 + 基础值
			add := growRatioFinal.Mul(decimal.NewFromInt32(globalConfig.HeroPromoteIntensityRatio[m.hero.PromoteLevel])).Add(promoteFinal)

			// 品质参数
			qualityRatio := globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)]

			// 职业参数
			professionRatio := professionEntry.GetFinalRatio(n)

			value := add.Mul(qualityRatio).Mul(professionRatio).Round(0)
			m.ModFinalAttValue(n, value)
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
