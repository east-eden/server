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

func (m *AttManager) GetBaseAttId() int32 {
	return m.baseAttId
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
	m.attBase[define.Att_Hit] = attEntry.Hit
	m.attBase[define.Att_Dodge] = attEntry.Dodge

	m.attPercent[define.Att_Atk] = attEntry.AtkPercent
	m.attPercent[define.Att_Armor] = attEntry.ArmorPercent
	m.attPercent[define.Att_Heal] = attEntry.HealPercent
	m.attPercent[define.Att_MoveSpeed] = attEntry.MoveSpeedPercent
	m.attPercent[define.Att_AtbSpeed] = attEntry.AtbSpeedPercent
	m.attPercent[define.Att_MaxHP] = attEntry.MaxHPPercent

	for n := 0; n < len(attEntry.DmgOfType); n++ {
		m.attBase[define.Att_DmgTypeBegin+n] = attEntry.DmgOfType[n]
	}

	for n := 0; n < len(attEntry.ResOfType); n++ {
		m.attBase[define.Att_ResTypeBegin+n] = attEntry.ResOfType[n]
	}
}

func (m *AttManager) CalcAtt() {
	for n := range m.attFinal {
		m.attFinal[n] = 0
	}

	for n := define.Att_Begin; n < define.Att_End; n++ {
		value64 := float64(m.attBase[n]+m.attFlat[n]) * float64(float64(define.PercentBase+m.attPercent[n])/float64(define.PercentBase))
		m.attFinal[n] = int32(utils.Round(value64))
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

func (m *AttManager) GetBaseAtt(idx int) int32 {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		return m.attBase[idx]
	} else {
		return 0
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

func (m *AttManager) SetPercentAtt(idx int, value int32) {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		m.attPercent[idx] = value
	} else {
		return
	}
}

func (m *AttManager) GetPercentAtt(idx int) int32 {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		return m.attPercent[idx]
	} else {
		return 0
	}
}

func (m *AttManager) ModAttManager(r *AttManager) {
	for n := define.Att_Begin; n < define.Att_End; n++ {
		m.attBase[n] += r.attBase[n]
		m.attFlat[n] += r.attFlat[n]
		m.attPercent[n] += r.attPercent[n]
	}
}
