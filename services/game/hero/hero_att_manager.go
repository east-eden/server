package hero

import (
	"errors"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	"bitbucket.org/east-eden/server/internal/att"
	"bitbucket.org/east-eden/server/utils"
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

			m.SetPercentAtt(n, add)
		}
	}
}

func (m *HeroAttManager) CalcPromote() {

}

//////////////////////////////////////////////
// 计算所有装备栏属性
func (m *HeroAttManager) CalcEquipBar() {
	// equip bar
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

	// rune box
	for n = 0; n < define.Rune_PositionEnd; n++ {
		r := m.hero.runeBox.GetRuneByPos(n)
		if r == nil {
			continue
		}

		r.GetAttManager().Reset()
		r.CalcAtt()
		m.ModAttManager(r.GetAttManager())
	}
}
