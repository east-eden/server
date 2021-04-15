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
	baseAttId int32                // 基础属性id
	attValue  [define.AttNum]int32 // 属性值
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
		return m.attValue[index]
	} else {
		return 0
	}
}

func (m *AttManager) SetAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attValue[index] = value
	}
}

func (m *AttManager) ModAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attValue[index] += value
	}
}

func (m *AttManager) Reset() {
	for k := range m.attValue {
		m.attValue[k] = 0
	}

	attEntry, ok := auto.GetAttEntry(m.baseAttId)
	if !ok {
		return
	}

	m.attValue[define.Att_AtkBase] = int32(attEntry.Atk)
	m.attValue[define.Att_AtkPercent] = int32(attEntry.AtkPercent)
	m.attValue[define.Att_ArmorBase] = int32(attEntry.Armor)
	m.attValue[define.Att_ArmorPercent] = int32(attEntry.ArmorPercent)
	m.attValue[define.Att_DmgInc] = attEntry.DmgInc
	m.attValue[define.Att_Crit] = attEntry.Crit
	m.attValue[define.Att_CritInc] = attEntry.CritInc
	m.attValue[define.Att_Tenacity] = attEntry.Tenacity
	m.attValue[define.Att_HealBase] = attEntry.Heal
	m.attValue[define.Att_HealPercent] = int32(attEntry.HealPercent)
	m.attValue[define.Att_RealDmg] = attEntry.RealDmg
	m.attValue[define.Att_MoveSpeedBase] = int32(attEntry.MoveSpeed)
	m.attValue[define.Att_MoveSpeedPercent] = int32(attEntry.MoveSpeedPercent)
	m.attValue[define.Att_AtbSpeedBase] = int32(attEntry.AtbSpeed)
	m.attValue[define.Att_AtbSpeedPercent] = int32(attEntry.AtbSpeedPercent)
	m.attValue[define.Att_EffectHit] = attEntry.EffectHit
	m.attValue[define.Att_EffectResist] = attEntry.EffectResist
	m.attValue[define.Att_MaxHPBase] = attEntry.MaxHP
	m.attValue[define.Att_MaxHPPercent] = int32(attEntry.MaxHPPercent)
	m.attValue[define.Att_MaxMP] = attEntry.MaxMP
	m.attValue[define.Att_GenMP] = attEntry.GenMP
	m.attValue[define.Att_Rage] = attEntry.Rage
	m.attValue[define.Att_GenRagePercent] = int32(attEntry.GenRagePercent)
	m.attValue[define.Att_InitRage] = attEntry.InitRage
	m.attValue[define.Att_Hit] = attEntry.Hit
	m.attValue[define.Att_Dodge] = attEntry.Dodge
	m.attValue[define.Att_MoveScope] = int32(attEntry.MoveScope)

	for n := 0; n < len(attEntry.DmgOfType); n++ {
		m.attValue[define.Att_DmgTypeBegin+n] = int32(attEntry.DmgOfType[n])
	}

	for n := 0; n < len(attEntry.ResOfType); n++ {
		m.attValue[define.Att_ResTypeBegin+n] = int32(attEntry.ResOfType[n])
	}
}

func (m *AttManager) CalcAtt() {
	atk := int64(m.attValue[define.Att_AtkBase]) * int64(define.PercentBase+m.attValue[define.Att_AtkPercent]) / int64(define.PercentBase)
	m.attValue[define.Att_AtkFinal] = int32(atk)

	armor := int64(m.attValue[define.Att_ArmorBase]) * int64(define.PercentBase+m.attValue[define.Att_ArmorPercent]) / int64(define.PercentBase)
	m.attValue[define.Att_ArmorFinal] = int32(armor)

	heal := int64(m.attValue[define.Att_HealBase]) * int64(define.PercentBase+m.attValue[define.Att_HealPercent]) / int64(define.PercentBase)
	m.attValue[define.Att_HealFinal] = int32(heal)

	moveSpeed := int64(m.attValue[define.Att_MoveSpeedBase]) * int64(define.PercentBase+m.attValue[define.Att_MoveSpeedPercent]) / int64(define.PercentBase)
	m.attValue[define.Att_MoveSpeedFinal] = int32(moveSpeed)

	atbSpeed := int64(m.attValue[define.Att_AtbSpeedBase]) * int64(define.PercentBase+m.attValue[define.Att_AtbSpeedPercent]) / int64(define.PercentBase)
	m.attValue[define.Att_AtbSpeedFinal] = int32(atbSpeed)

	maxHP := int64(m.attValue[define.Att_MaxHPBase]) * int64(define.PercentBase+m.attValue[define.Att_MaxHPPercent]) / int64(define.PercentBase)
	m.attValue[define.Att_MaxHPFinal] = int32(maxHP)
}

func (m *AttManager) Export() []int32 {
	att := make([]int32, len(m.attValue))
	for n := define.Att_Begin; n < define.Att_End; n++ {
		att[n] = m.attValue[n]
	}

	return att
}

func (m *AttManager) ModAttManager(r *AttManager) {
	for n := define.Att_Begin; n < define.Att_End; n++ {
		m.attValue[n] += r.attValue[n]
	}
}
