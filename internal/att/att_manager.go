package att

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/utils"
)

type AttManager struct {
	baseAttId  int32                // 基础属性id
	attBase    [define.AttNum]int32 // 基础属性
	attFlat    [define.AttNum]int32 // 平值加成
	attPercent [define.AttNum]int32 // 百分比加成

	attFinal [define.AttNum]int32 // 计算后最终属性32位
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

func (m *AttManager) GetAttValue(index int) int32 {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		return m.attFinal[index]
	} else {
		return 0
	}
}

func (m *AttManager) SetAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attFinal[index] = value
	}
}

func (m *AttManager) ModAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attFinal[index] += value
	}
}

func (m *AttManager) Reset() {
	for k := range m.attFinal {
		m.attBase[k] = 0
		m.attFlat[k] = 0
		m.attPercent[k] = 0
		m.attFinal[k] = 0
	}

	attEntry, ok := auto.GetAttEntry(m.baseAttId)
	if !ok {
		return
	}

	m.attBase[define.Att_Atk] = attEntry.Atk
	m.attBase[define.Att_Armor] = attEntry.Armor
	m.attBase[define.Att_DmgInc] = attEntry.DmgInc
	m.attBase[define.Att_Crit] = attEntry.Crit
	m.attBase[define.Att_CritInc] = attEntry.CritInc
	m.attBase[define.Att_Heal] = attEntry.Heal
	m.attBase[define.Att_RealDmg] = attEntry.RealDmg
	m.attBase[define.Att_MoveSpeed] = attEntry.MoveSpeed
	m.attBase[define.Att_AtbSpeed] = attEntry.AtbSpeed
	m.attBase[define.Att_EffectHit] = attEntry.EffectHit
	m.attBase[define.Att_EffectResist] = attEntry.EffectResist
	m.attBase[define.Att_MaxHP] = attEntry.MaxHP
	m.attBase[define.Att_MaxMP] = attEntry.MaxMP
	m.attBase[define.Att_GenMP] = attEntry.GenMP
	m.attBase[define.Att_Rage] = attEntry.Rage

	for n := 0; n < len(attEntry.DmgOfType); n++ {
		m.attBase[define.Att_DmgTypeBegin+n] = attEntry.DmgOfType[n]
	}

	for n := 0; n < len(attEntry.ResOfType); n++ {
		m.attBase[define.Att_ResTypeBegin+n] = attEntry.ResOfType[n]
	}
}

func (m *AttManager) CalcAtt() {
	for k := range m.attFinal {
		m.attFinal[k] = 0
	}

	for k := 0; k < define.AttNum; k++ {
		m.attFinal[k] = (m.attBase[k] + m.attFlat[k]) * (define.AttPercentBase + m.attPercent[k]) / define.AttPercentBase
	}
}

func (m *AttManager) ModBaseAtt(idx int, value int32) {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		m.attBase[idx] += value
	} else {
		return
	}
}

func (m *AttManager) SetBaseAtt(idx int, value int32) {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		m.attBase[idx] = value
	} else {
		return
	}
}

func (m *AttManager) ModFlatAtt(idx int, value int32) {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		m.attFlat[idx] += value
	} else {
		return
	}
}

func (m *AttManager) ModPercentAtt(idx int, value int32) {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		m.attPercent[idx] += value
	} else {
		return
	}
}

func (m *AttManager) ModAttManager(r *AttManager) {
	for k := 0; k < define.AttNum; k++ {
		m.attFinal[k] += r.attFinal[k]
	}
}
