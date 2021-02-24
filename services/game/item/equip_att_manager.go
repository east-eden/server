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
	m := &EquipAttManager{
		equip: equip,
	}

	m.AttManager.SetBaseAttId(equip.GetEquipEnchantEntry().AttId)
	return m
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
		func() {
			baseAttValue := m.GetBaseAtt(n)

			if baseAttValue == 0 {
				return
			}

			value := m.GetBaseAtt(n) + attGrowRatio.GetBaseAtt(n)*int32(m.equip.Level)
			m.SetBaseAtt(n, value*(1+globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)]))
		}()

		// percent value
		func() {
			percentAttValue := m.GetPercentAtt(n)

			if percentAttValue == 0 {
				return
			}

			m.ModPercentAtt(n, attGrowRatio.GetPercentAtt(n)*int32(m.equip.Level))
			m.SetPercentAtt(n, m.GetPercentAtt(n)*globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)])
		}()

	}
}

func (m *EquipAttManager) CalcPromote() {

}

func (m *EquipAttManager) CalcStarup() {

}
