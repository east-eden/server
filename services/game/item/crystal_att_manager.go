package item

import (
	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/internal/att"
	"bitbucket.org/funplus/server/utils"
	"github.com/rs/zerolog/log"
)

// crystal属性计算管理
type CrystalAttManager struct {
	c *Crystal
	att.AttManager
}

func NewCrystalAttManager(c *Crystal) *CrystalAttManager {
	m := &CrystalAttManager{
		c: c,
	}

	m.AttManager.SetBaseAttId(-1)
	return m
}

// 计算晶石属性
func (m *CrystalAttManager) CalcAtt() {
	// 主属性
	m.CalcMainAtt()

	// 副属性
	m.CalcViceAtts()
}

//////////////////////////////////////////////
// 主属性 = （晶石强化等级*强度系数*各属性成长率+各属性固定值）*品质系数
func (m *CrystalAttManager) CalcMainAtt() {
	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		log.Error().Caller().Err(auto.ErrGlobalConfigInvalid).Send()
		return
	}

	mainAttRepoEntry, ok := auto.GetCrystalAttRepoEntry(m.c.MainAtt.AttRepoId)
	if !ok {
		log.Error().Caller().Int32("id", m.c.MainAtt.AttRepoId).Msg("cannot find CrystalAttRepoEntry")
		return
	}

	// 属性固定值
	baseAtt := att.AttManager{}
	baseAtt.SetBaseAttId(mainAttRepoEntry.AttId)

	// 属性成长率
	growAtt := att.AttManager{}
	growAtt.SetBaseAttId(mainAttRepoEntry.AttGrowRatioId)

	for n := define.Att_Begin; n < define.Att_End; n++ {
		// base value
		baseAttValue := baseAtt.GetBaseAttValue(n)
		growRatioBase := growAtt.GetBaseAttValue(n)
		if baseAttValue != 0 && growRatioBase != 0 {
			// 晶石强化等级*强度系数*成长率
			add := int32(m.c.Level) * globalConfig.CrystalLevelupIntensityRatio * growRatioBase

			// 品质系数
			qualityRatio := globalConfig.CrystalLevelupMainQualityRatio[m.c.ItemEntry.Quality]

			value64 := float64(add+baseAttValue) * (float64(qualityRatio) / float64(define.PercentBase))
			value := int32(utils.Round(value64))
			if value < 0 {
				utils.ErrPrint(att.ErrAttValueOverflow, "crystal main att calc failed", n, value, m.c.Id)
			}

			m.ModBaseAttValue(n, value)
		}
	}
}

//////////////////////////////////////////////
// 副属性 = 各属性固定值*品质系数*随机区间系数
func (m *CrystalAttManager) CalcViceAtts() {
	if len(m.c.ViceAtts) == 0 {
		return
	}

	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		log.Error().Caller().Err(auto.ErrGlobalConfigInvalid).Send()
		return
	}

	for _, viceAtt := range m.c.ViceAtts {
		attRepoEntry, ok := auto.GetCrystalAttRepoEntry(viceAtt.AttRepoId)
		if !ok {
			log.Error().Caller().Int32("repo_id", viceAtt.AttRepoId).Msg("invalid crystal att repo entry")
			return
		}

		viceAttManager := att.AttManager{}
		viceAttManager.SetBaseAttId(attRepoEntry.AttId)

		// 品质系数
		qualityRatio := globalConfig.CrystalLevelupViceQualityRatio[m.c.ItemEntry.Quality]

		// 随机区间系数
		randRatio := viceAtt.AttRandRatio

		for n := define.Att_Begin; n < define.Att_End; n++ {
			baseAttValue := viceAttManager.GetBaseAttValue(n)

			if baseAttValue != 0 {
				value64 := float64(baseAttValue) * (float64(qualityRatio) / float64(define.PercentBase)) * (float64(randRatio) / float64(define.PercentBase))
				value := int32(utils.Round(value64))
				if value < 0 {
					utils.ErrPrint(att.ErrAttValueOverflow, "crystal vice att calc failed", n, value, m.c.Id)
				}

				m.ModBaseAttValue(n, value)
			}
		}
	}
}
