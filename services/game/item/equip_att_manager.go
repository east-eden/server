package item

import (
	"errors"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/internal/att"
	"github.com/east-eden/server/utils"
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
		utils.ErrPrint(errors.New("invalid global config"), "equip CalcLevelup failed")
		return
	}

	// 成长率att
	attGrowRatio := att.NewAttManager()
	attGrowRatio.SetBaseAttId(globalConfig.EquipLevelGrowRatioAttId)

	for n := define.Att_Begin; n < define.Att_End; n++ {
		baseAttValue := m.GetAttValue(n)
		growRatioBase := attGrowRatio.GetAttValue(n)
		if baseAttValue != 0 && growRatioBase != 0 {
			// 等级*升级成长率
			add := growRatioBase * int32(m.equip.Level)

			// 品质参数
			qualityRatio := globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)]

			value64 := float64(add+baseAttValue) * (float64(qualityRatio) / float64(define.PercentBase))
			value := int32(utils.Round(value64))
			if value < 0 {
				utils.ErrPrint(att.ErrAttValueOverflow, "equip att calc failed", n, value, m.equip.Id)
			}

			m.SetAttValue(n, value)
		}
	}
}

//////////////////////////////////////////////
// 突破属性 = (突破强度等级*各属性成长率+各属性固定值)*装备品质参数
func (m *EquipAttManager) CalcPromote() {
	globalConfig, ok := auto.GetGlobalConfig()
	if !ok {
		utils.ErrPrint(errors.New("invalid global config"), "equip CalcPromote failed")
		return
	}

	// 成长率att
	attGrowRatio := att.NewAttManager()
	attGrowRatio.SetBaseAttId(m.equip.EquipEnchantEntry.AttPromoteGrowupId)

	// 基础att
	promoteBaseAtt := att.NewAttManager()
	promoteBaseAtt.SetBaseAttId(m.equip.EquipEnchantEntry.AttPromoteBaseId)

	for n := define.Att_Begin; n < define.Att_End; n++ {
		growRatioBase := attGrowRatio.GetAttValue(n)
		promoteBase := promoteBaseAtt.GetAttValue(n)
		if promoteBase != 0 && growRatioBase != 0 {
			// 强度等级*升级成长率 + 属性固定值
			add := growRatioBase*globalConfig.EquipPromoteIntensityRatio[m.equip.Promote] + promoteBase

			// 品质参数
			qualityRatio := globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)]

			value64 := float64(add) * (float64(qualityRatio) / float64(define.PercentBase))
			value := int32(utils.Round(value64))
			if value < 0 {
				utils.ErrPrint(att.ErrAttValueOverflow, "equip att calc failed", n, value, m.equip.Id)
			}

			m.SetAttValue(n, value)
		}
	}
}

func (m *EquipAttManager) CalcStarup() {

}
