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

	for n := define.Att_Begin; n < define.Att_End; n++ {

		// base value
		func() {
			baseAttValue := m.GetBaseAtt(n)

			if baseAttValue == 0 {
				return
			}

			value := m.GetBaseAtt(n) + attGrowRatio.GetBaseAtt(n)*int32(m.hero.Level)
			m.SetBaseAtt(n, value*(1+globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)]))
		}()

		// percent value
		func() {
			percentAttValue := m.GetPercentAtt(n)

			if percentAttValue == 0 {
				return
			}

			m.ModPercentAtt(n, attGrowRatio.GetPercentAtt(n)*int32(m.hero.Level))
			m.SetPercentAtt(n, m.GetPercentAtt(n)*globalConfig.HeroLevelQualityRatio[int(m.hero.Entry.Quality)])
		}()

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

		e.GetAttManager().CalcAtt()
		m.ModAttManager(&e.GetAttManager().AttManager)
	}

	// rune box
	for n = 0; n < define.Rune_PositionEnd; n++ {
		r := m.hero.runeBox.GetRuneByPos(n)
		if r == nil {
			continue
		}

		r.CalcAtt()
		m.ModAttManager(r.GetAttManager())
	}
}
