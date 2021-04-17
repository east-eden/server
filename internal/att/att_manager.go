package att

import (
	"errors"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/utils"
)

var (
	ErrAttValueOverflow = errors.New("att value overflow")
)

type AttManager struct {
	baseAttId  int32                       // 基础属性id
	attFinal   [define.AttFinalNum]int32   // 属性值
	attBase    [define.AttBaseNum]int32    // 基础属性值
	attPercent [define.AttPercentNum]int32 // 百分比属性值
}

func NewAttManager() *AttManager {
	m := &AttManager{}

	m.Reset()

	return m
}

func (m *AttManager) SetBaseAttId(attId int32) {
	m.baseAttId = attId
	m.Reset()
}

func (m *AttManager) GetBaseAttId() int32 {
	return m.baseAttId
}

func (m *AttManager) GetFinalAttValue(index int) int32 {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		return m.attFinal[index]
	} else {
		return 0
	}
}

func (m *AttManager) SetFinalAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attFinal[index] = value
	}
}

func (m *AttManager) ModFinalAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attFinal[index] += value

		if m.attFinal[index] < 0 {
			m.attFinal[index] = 0
		}
	}
}

func (m *AttManager) GetBaseAttValue(index int) int32 {
	if utils.Between(index, define.Att_Base_Begin, define.Att_Base_End) {
		return m.attBase[index]
	} else {
		return 0
	}
}

func (m *AttManager) SetBaseAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Base_Begin, define.Att_Base_End) {
		m.attBase[index] = value
	}
}

func (m *AttManager) ModBaseAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Base_Begin, define.Att_Base_End) {
		m.attBase[index] += value

		if m.attBase[index] < 0 {
			m.attBase[index] = 0
		}
	}
}

func (m *AttManager) GetPercentAttValue(index int) int32 {
	if utils.Between(index, define.Att_Percent_Begin, define.Att_Percent_End) {
		return m.attPercent[index]
	} else {
		return 0
	}
}

func (m *AttManager) SetPercentAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Percent_Begin, define.Att_Percent_End) {
		m.attPercent[index] = value
	}
}

func (m *AttManager) ModPercentAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Percent_Begin, define.Att_Percent_End) {
		m.attPercent[index] += value

		if m.attPercent[index] < 0 {
			m.attPercent[index] = 0
		}
	}
}

func (m *AttManager) Reset() {
	for k := range m.attFinal {
		m.attFinal[k] = 0
	}

	attEntry, ok := auto.GetAttEntry(m.baseAttId)
	if !ok {
		return
	}

	// 基础属性
	m.attBase[define.Att_AtkBase] = int32(attEntry.Atk)
	m.attBase[define.Att_ArmorBase] = int32(attEntry.Armor)
	m.attBase[define.Att_HealBase] = attEntry.Heal
	m.attBase[define.Att_MoveSpeedBase] = int32(attEntry.MoveSpeed)
	m.attBase[define.Att_AtbSpeedBase] = int32(attEntry.AtbSpeed)
	m.attBase[define.Att_MaxHPBase] = attEntry.MaxHP

	// 百分比属性
	m.attPercent[define.Att_AtkPercent] = int32(attEntry.AtkPercent)
	m.attPercent[define.Att_ArmorPercent] = int32(attEntry.ArmorPercent)
	m.attPercent[define.Att_HealPercent] = int32(attEntry.HealPercent)
	m.attPercent[define.Att_MoveSpeedPercent] = int32(attEntry.MoveSpeedPercent)
	m.attPercent[define.Att_AtbSpeedPercent] = int32(attEntry.AtbSpeedPercent)
	m.attPercent[define.Att_MaxHPPercent] = int32(attEntry.MaxHPPercent)

	// 最终属性
	m.attFinal[define.Att_DmgInc] = attEntry.DmgInc
	m.attFinal[define.Att_Crit] = attEntry.Crit
	m.attFinal[define.Att_CritInc] = attEntry.CritInc
	m.attFinal[define.Att_Tenacity] = attEntry.Tenacity
	m.attFinal[define.Att_RealDmg] = attEntry.RealDmg
	m.attFinal[define.Att_EffectHit] = attEntry.EffectHit
	m.attFinal[define.Att_EffectResist] = attEntry.EffectResist
	m.attFinal[define.Att_MaxMP] = attEntry.MaxMP
	m.attFinal[define.Att_GenMP] = attEntry.GenMP
	m.attFinal[define.Att_Rage] = attEntry.Rage
	m.attFinal[define.Att_GenRagePercent] = int32(attEntry.GenRagePercent)
	m.attFinal[define.Att_InitRage] = attEntry.InitRage
	m.attFinal[define.Att_Hit] = attEntry.Hit
	m.attFinal[define.Att_Dodge] = attEntry.Dodge
	m.attFinal[define.Att_MoveScope] = int32(attEntry.MoveScope)

	for n := 0; n < len(attEntry.DmgOfType); n++ {
		m.attFinal[define.Att_DmgTypeBegin+n] = int32(attEntry.DmgOfType[n])
	}

	for n := 0; n < len(attEntry.ResOfType); n++ {
		m.attFinal[define.Att_ResTypeBegin+n] = int32(attEntry.ResOfType[n])
	}
}

func (m *AttManager) CalcAtt() {
	atk := float64(m.attBase[define.Att_AtkBase]) / float64(define.PercentBase) * float64(define.PercentBase+m.attPercent[define.Att_AtkPercent]) / float64(define.PercentBase)
	m.attFinal[define.Att_Atk] = int32(utils.Round(atk))

	armor := float64(m.attFinal[define.Att_ArmorBase]) / float64(define.PercentBase) * float64(define.PercentBase+m.attFinal[define.Att_ArmorPercent]) / float64(define.PercentBase)
	m.attFinal[define.Att_Armor] = int32(utils.Round(armor))

	heal := float64(m.attFinal[define.Att_HealBase]) * float64(define.PercentBase+m.attFinal[define.Att_HealPercent]) / float64(define.PercentBase)
	m.attFinal[define.Att_Heal] = int32(utils.Round(heal))

	moveSpeed := float64(m.attFinal[define.Att_MoveSpeedBase]) / float64(define.PercentBase) * float64(define.PercentBase+m.attFinal[define.Att_MoveSpeedPercent]) / float64(define.PercentBase)
	m.attFinal[define.Att_MoveSpeed] = int32(utils.Round(moveSpeed))

	atbSpeed := float64(m.attFinal[define.Att_AtbSpeedBase]) / float64(define.PercentBase) * float64(define.PercentBase+m.attFinal[define.Att_AtbSpeedPercent]) / float64(define.PercentBase)
	m.attFinal[define.Att_AtbSpeed] = int32(utils.Round(atbSpeed))

	maxHP := float64(m.attFinal[define.Att_MaxHPBase]) * float64(define.PercentBase+m.attFinal[define.Att_MaxHPPercent]) / float64(define.PercentBase)
	m.attFinal[define.Att_MaxHP] = int32(utils.Round(maxHP))
}

func (m *AttManager) Export() []int32 {
	att := make([]int32, len(m.attFinal))
	for n := define.Att_Begin; n < define.Att_End; n++ {
		att[n] = m.attFinal[n]
	}

	return att
}

func (m *AttManager) ModAttManager(r *AttManager) {
	for n := define.Att_Begin; n < define.Att_End; n++ {
		m.attFinal[n] += r.attFinal[n]
	}
}
