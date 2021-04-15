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
	attBase   [define.AttNum]int32 // 基础属性值
	attFinal  [define.AttNum]int32 // 最终属性值
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

func (m *AttManager) GetBaseAttValue(index int) int32 {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		return m.attBase[index]
	} else {
		return 0
	}
}

func (m *AttManager) SetBaseAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attBase[index] = value
	}
}

func (m *AttManager) ModBaseAttValue(index int, value int32) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attBase[index] += value
	}
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
	}
}

func (m *AttManager) Reset() {
	for k := range m.attBase {
		m.attBase[k] = 0
	}

	attEntry, ok := auto.GetAttEntry(m.baseAttId)
	if !ok {
		return
	}

	m.attBase[define.Att_Atk] = int32(attEntry.Atk)
	m.attBase[define.Att_AtkPercent] = int32(attEntry.AtkPercent)
	m.attBase[define.Att_Armor] = int32(attEntry.Armor)
	m.attBase[define.Att_ArmorPercent] = int32(attEntry.ArmorPercent)
	m.attBase[define.Att_DmgInc] = attEntry.DmgInc
	m.attBase[define.Att_Crit] = attEntry.Crit
	m.attBase[define.Att_CritInc] = attEntry.CritInc
	m.attBase[define.Att_Tenacity] = attEntry.Tenacity
	m.attBase[define.Att_Heal] = attEntry.Heal
	m.attBase[define.Att_HealPercent] = int32(attEntry.HealPercent)
	m.attBase[define.Att_RealDmg] = attEntry.RealDmg
	m.attBase[define.Att_MoveSpeed] = int32(attEntry.MoveSpeed)
	m.attBase[define.Att_MoveSpeedPercent] = int32(attEntry.MoveSpeedPercent)
	m.attBase[define.Att_AtbSpeed] = int32(attEntry.AtbSpeed)
	m.attBase[define.Att_AtbSpeedPercent] = int32(attEntry.AtbSpeedPercent)
	m.attBase[define.Att_EffectHit] = attEntry.EffectHit
	m.attBase[define.Att_EffectResist] = attEntry.EffectResist
	m.attBase[define.Att_MaxHP] = attEntry.MaxHP
	m.attBase[define.Att_MaxHPPercent] = int32(attEntry.MaxHPPercent)
	m.attBase[define.Att_MaxMP] = attEntry.MaxMP
	m.attBase[define.Att_GenMP] = attEntry.GenMP
	m.attBase[define.Att_Rage] = attEntry.Rage
	m.attBase[define.Att_GenRagePercent] = int32(attEntry.GenRagePercent)
	m.attBase[define.Att_InitRage] = attEntry.InitRage
	m.attBase[define.Att_Hit] = attEntry.Hit
	m.attBase[define.Att_Dodge] = attEntry.Dodge
	m.attBase[define.Att_MoveScope] = int32(attEntry.MoveScope)

	for n := 0; n < len(attEntry.DmgOfType); n++ {
		m.attBase[define.Att_DmgTypeBegin+n] = int32(attEntry.DmgOfType[n])
	}

	for n := 0; n < len(attEntry.ResOfType); n++ {
		m.attBase[define.Att_ResTypeBegin+n] = int32(attEntry.ResOfType[n])
	}
}

func (m *AttManager) CalcAtt() {
	copy(m.attFinal[:], m.attBase[:])

	atk := int64(m.attBase[define.Att_Atk]) * int64(define.PercentBase+m.attBase[define.Att_AtkPercent]) / int64(define.PercentBase)
	m.attFinal[define.Att_Atk] = int32(atk)

	armor := int64(m.attBase[define.Att_Armor]) * int64(define.PercentBase+m.attBase[define.Att_ArmorPercent]) / int64(define.PercentBase)
	m.attFinal[define.Att_Armor] = int32(armor)

	heal := int64(m.attBase[define.Att_Heal]) * int64(define.PercentBase+m.attBase[define.Att_HealPercent]) / int64(define.PercentBase)
	m.attFinal[define.Att_Heal] = int32(heal)

	moveSpeed := int64(m.attBase[define.Att_MoveSpeed]) * int64(define.PercentBase+m.attBase[define.Att_MoveSpeedPercent]) / int64(define.PercentBase)
	m.attFinal[define.Att_MoveSpeed] = int32(moveSpeed)

	atbSpeed := int64(m.attBase[define.Att_AtbSpeed]) * int64(define.PercentBase+m.attBase[define.Att_AtbSpeedPercent]) / int64(define.PercentBase)
	m.attFinal[define.Att_AtbSpeed] = int32(atbSpeed)

	maxHP := int64(m.attBase[define.Att_MaxHP]) * int64(define.PercentBase+m.attBase[define.Att_MaxHPPercent]) / int64(define.PercentBase)
	m.attFinal[define.Att_MaxHP] = int32(maxHP)
}

func (m *AttManager) ExportBase() []int32 {
	att := make([]int32, len(m.attBase))
	for n := define.Att_Begin; n < define.Att_End; n++ {
		att[n] = m.attBase[n]
	}

	return att
}

func (m *AttManager) ModBaseAttManager(r *AttManager) {
	for n := define.Att_Begin; n < define.Att_End; n++ {
		m.attBase[n] += r.attBase[n]
	}
}

func (m *AttManager) ModFinalAttManager(r *AttManager) {
	for n := define.Att_Begin; n < define.Att_End; n++ {
		m.attFinal[n] += r.attFinal[n]
	}
}
