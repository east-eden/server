package item

import (
	"errors"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	"bitbucket.org/east-eden/server/internal/att"
	"bitbucket.org/east-eden/server/utils"
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
		// base value
		baseAttValue := m.GetBaseAtt(n)
		growRatioBase := attGrowRatio.GetBaseAtt(n)
		if baseAttValue != 0 && growRatioBase != 0 {
			// 等级*升级成长率
			add := growRatioBase * int32(m.equip.Level)

			// 品质参数
			qualityRatio := globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)]

			value64 := float64(add+baseAttValue) * (float64(qualityRatio) / float64(define.AttPercentBase))
			value := int32(utils.Round(value64))
			if value < 0 {
				utils.ErrPrint(att.ErrAttValueOverflow, "equip att calc failed", n, value, m.equip.Id)
			}

			m.SetBaseAtt(n, value)
		}

		// percent value
		percentAttValue := m.GetPercentAtt(n)
		growRatioPercent := attGrowRatio.GetPercentAtt(n)
		if percentAttValue != 0 && growRatioPercent != 0 {
			// 等级*升级成长率
			add := attGrowRatio.GetPercentAtt(n) * int32(m.equip.Level)

			// 品质参数
			qualityRatio := globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)]

			value64 := float64(add+percentAttValue) * (float64(qualityRatio) / float64(define.AttPercentBase))
			value := int32(utils.Round(value64))
			if value < 0 {
				utils.ErrPrint(att.ErrAttValueOverflow, "equip att calc failed", n, value, m.equip.Id)
			}

			m.SetPercentAtt(n, add)
		}
	}
}

func (m *EquipAttManager) CalcPromote() {

}

func (m *EquipAttManager) CalcStarup() {

}
