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
		func() {
			baseAttValue := m.GetBaseAtt(n)

			if baseAttValue == 0 {
				return
			}

			// 等级*升级成长率
			add := attGrowRatio.GetBaseAtt(n) * int32(m.equip.Level)

			// 品质参数
			qualityRatio := globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)]

			add64 := float64(add) * (float64(qualityRatio) / float64(define.AttPercentBase))
			add = int32(utils.Round(add64))
			if add < 0 {
				utils.ErrPrint(att.ErrAttValueOverflow, "equip att calc failed", n, add, m.equip.Id)
			}

			m.ModBaseAtt(n, add)
		}()

		// percent value
		func() {
			percentAttValue := m.GetPercentAtt(n)

			if percentAttValue == 0 {
				return
			}

			// 基础值+等级*升级成长率
			add := percentAttValue + attGrowRatio.GetPercentAtt(n)*int32(m.equip.Level)

			// 品质参数
			qualityRatio := globalConfig.EquipLevelQualityRatio[int(m.equip.Entry().Quality)]

			add64 := float64(add) * (float64(qualityRatio) / float64(define.AttPercentBase))
			add = int32(utils.Round(add64))
			if add < 0 {
				utils.ErrPrint(att.ErrAttValueOverflow, "equip att calc failed", n, add, m.equip.Id)
			}

			m.ModPercentAtt(n, add)
		}()

	}
}

func (m *EquipAttManager) CalcPromote() {

}

func (m *EquipAttManager) CalcStarup() {

}
