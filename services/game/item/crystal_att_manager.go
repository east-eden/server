package item

import (
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/internal/att"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
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
	baseAtt := att.NewAttManager()
	baseAtt.SetBaseAttId(mainAttRepoEntry.AttId)

	// 属性成长率
	growAtt := att.NewAttManager()
	growAtt.SetBaseAttId(mainAttRepoEntry.AttGrowRatioId)

	// 晶石强化等级*强度系数
	intensity := decimal.NewFromInt32(int32(m.c.Level) * globalConfig.CrystalLevelupIntensityRatio)

	// 品质系数
	qualityRatio := globalConfig.CrystalLevelupMainQualityRatio[m.c.ItemEntry.Quality]

	growAtt.Mul(intensity).ModAttManager(baseAtt).Mul(qualityRatio).Round()
	m.ModAttManager(growAtt)
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

		viceAttManager := att.NewAttManager()
		viceAttManager.SetBaseAttId(attRepoEntry.AttId)

		// 品质系数
		qualityRatio := globalConfig.CrystalLevelupViceQualityRatio[m.c.ItemEntry.Quality]

		// 随机区间系数
		randRatio := viceAtt.AttRandRatio

		viceAttManager.Mul(qualityRatio).Mul(randRatio).Round()
		m.ModAttManager(viceAttManager)
	}
}
