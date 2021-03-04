package hero

import (
	"errors"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/internal/att"
	"bitbucket.org/funplus/server/utils"
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
// 升级属性 =（卡牌等级*各升级属性成长率+各升级属性固定值）*卡牌品质参数*卡牌职业参数
func (m *HeroAttManager) CalcLevelup() {
	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		utils.ErrPrint(errors.New("invalid global config"), "hero CalcLevelup failed")
		return
	}

	// 成长率att
	attGrowRatio := att.NewAttManager()
	attGrowRatio.SetBaseAttId(globalConfig.HeroLevelGrowRatioAttId)

	// 职业参数
	professionEntry, ok := auto.GetHeroProfessionEntry(m.hero.Entry.Profession)
	if !ok {
		utils.ErrPrint(errors.New("invalid profession"), "hero CalcLevelup failed", m.hero.Entry.Profession)
		return
	}

	for n := define.Att_Begin; n < define.Att_End; n++ {
		// base value
		baseAttValue := m.GetBaseAtt(n)
		growRatioBase := attGrowRatio.GetBaseAtt(n)
		if growRatioBase != 0 {
			// 等级*升级成长率
			add := growRatioBase * int32(m.hero.Level)

			// 品质参数
			qualityRatio := globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)]

			// 职业参数
			professionRatio := professionEntry.GetRatio(n)

			value64 := float64(add+baseAttValue) * (float64(qualityRatio) / float64(define.AttPercentBase)) * (float64(professionRatio) / float64(define.AttPercentBase))
			value := int32(utils.Round(value64))
			if value < 0 {
				utils.ErrPrint(att.ErrAttValueOverflow, "hero att calc failed", n, value, m.hero.Id)
			}

			m.SetBaseAtt(n, value)
		}

		// percent value
		percentAttValue := m.GetPercentAtt(n)
		growRatioPercent := attGrowRatio.GetPercentAtt(n)
		if growRatioPercent != 0 {
			// 等级*升级成长率
			add := growRatioPercent * int32(m.hero.Level)

			// 品质参数
			qualityRatio := globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)]

			// 职业参数
			professionRatio := professionEntry.GetRatio(n)

			value64 := float64(add+percentAttValue) * (float64(qualityRatio) / float64(define.AttPercentBase)) * (float64(professionRatio) / float64(define.AttPercentBase))
			value := int32(utils.Round(value64))
			if value < 0 {
				utils.ErrPrint(att.ErrAttValueOverflow, "hero att calc failed", n, value, m.hero.Id)
			}

			m.SetPercentAtt(n, value)
		}
	}
}

func (m *HeroAttManager) CalcPromote() {

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
