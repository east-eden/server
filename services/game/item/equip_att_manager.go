package item

import (
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/internal/att"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

// 装备属性计算管理
type EquipAttManager struct {
	equip *Equip
	att.AttManager
}

func NewEquipAttManager(equip *Equip) *EquipAttManager {
	return &EquipAttManager{
		equip: equip,
	}
}

// 计算装备属性
func (m *EquipAttManager) CalcAtt() {
	// 升级
	m.CalcLevelup()

	// 突破
	m.CalcPromote()

	// 升星
	m.CalcStarup()
}

//////////////////////////////////////////////
// 升级属性 = (装备等级*各属性成长率+各属性固定值)*装备品质参数
func (m *EquipAttManager) CalcLevelup() {
	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		log.Error().Caller().Err(auto.ErrGlobalConfigInvalid).Send()
		return
	}

	// 属性固定值
	baseAtt := att.NewAttManager()
	baseAtt.SetBaseAttId(m.equip.EquipEnchantEntry.AttId)

	// 成长率att
	attGrowRatio := att.NewAttManager()
	attGrowRatio.SetBaseAttId(globalConfig.EquipLevelGrowRatioAttId)

	// 等级*升级成长率
	attGrowRatio.Mul(decimal.NewFromInt32(int32(m.equip.Level)))

	// 品质参数
	qualityRatio := globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)]

	baseAtt.ModAttManager(attGrowRatio)
	baseAtt.Mul(qualityRatio).Round()
	m.ModAttManager(baseAtt)
}

//////////////////////////////////////////////
// 突破属性 = (突破强度等级*各属性成长率+各属性固定值)*装备品质参数
func (m *EquipAttManager) CalcPromote() {
	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		log.Error().Caller().Err(auto.ErrGlobalConfigInvalid).Send()
		return
	}

	// 成长率att
	attGrowRatio := att.NewAttManager()
	attGrowRatio.SetBaseAttId(m.equip.EquipEnchantEntry.AttPromoteGrowupId)

	// 基础att
	promoteBaseAtt := att.NewAttManager()
	promoteBaseAtt.SetBaseAttId(m.equip.EquipEnchantEntry.AttPromoteBaseId)

	// 品质参数
	qualityRatio := globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)]

	attGrowRatio.Mul(decimal.NewFromInt32(globalConfig.EquipPromoteIntensityRatio[m.equip.Promote]))
	attGrowRatio.ModAttManager(promoteBaseAtt)
	attGrowRatio.Mul(qualityRatio).Round()
	m.ModAttManager(attGrowRatio)
}

func (m *EquipAttManager) CalcStarup() {

}
