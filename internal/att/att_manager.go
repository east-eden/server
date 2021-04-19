package att

import (
	"errors"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/utils"
	"github.com/shopspring/decimal"
)

var (
	ErrAttValueOverflow = errors.New("att value overflow")
)

func NewAttValue(i interface{}) decimal.Decimal {
	switch v := i.(type) {
	case int32:
		return decimal.NewFromInt32(v)
	case (decimal.Decimal):
		return decimal.NewFromBigInt(v.Coefficient(), v.Exponent())
	default:
		return decimal.NewFromInt32(0)
	}
}

type AttManager struct {
	baseAttId  int32                                 // 基础属性id
	attFinal   [define.AttFinalNum]decimal.Decimal   // 属性值
	attBase    [define.AttBaseNum]decimal.Decimal    // 基础属性值
	attPercent [define.AttPercentNum]decimal.Decimal // 百分比属性值
}

func NewAttManager() *AttManager {
	m := &AttManager{}

	return m
}

func (m *AttManager) SetBaseAttId(attId int32) {
	m.baseAttId = attId
	m.Reset()
}

func (m *AttManager) GetBaseAttId() int32 {
	return m.baseAttId
}

func (m *AttManager) GetFinalAttValue(index int) decimal.Decimal {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		return m.attFinal[index]
	} else {
		return decimal.NewFromInt32(0)
	}
}

func (m *AttManager) SetFinalAttValue(index int, value decimal.Decimal) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attFinal[index] = value
	}
}

func (m *AttManager) ModFinalAttValue(index int, value decimal.Decimal) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attFinal[index] = m.attFinal[index].Add(value)

		if m.attFinal[index].LessThan(decimal.NewFromInt32(0)) {
			m.attFinal[index] = decimal.NewFromInt32(0)
		}
	}
}

func (m *AttManager) GetBaseAttValue(index int) decimal.Decimal {
	if utils.Between(index, define.Att_Base_Begin, define.Att_Base_End) {
		return m.attBase[index]
	} else {
		return decimal.NewFromInt32(0)
	}
}

func (m *AttManager) SetBaseAttValue(index int, value decimal.Decimal) {
	if utils.Between(index, define.Att_Base_Begin, define.Att_Base_End) {
		m.attBase[index] = value
	}
}

func (m *AttManager) ModBaseAttValue(index int, value decimal.Decimal) {
	if utils.Between(index, define.Att_Base_Begin, define.Att_Base_End) {
		m.attBase[index] = m.attBase[index].Add(value)

		if m.attBase[index].LessThan(decimal.NewFromInt32(0)) {
			m.attBase[index] = decimal.NewFromInt32(0)
		}
	}
}

func (m *AttManager) GetPercentAttValue(index int) decimal.Decimal {
	if utils.Between(index, define.Att_Percent_Begin, define.Att_Percent_End) {
		return m.attPercent[index]
	} else {
		return decimal.NewFromInt32(0)
	}
}

func (m *AttManager) SetPercentAttValue(index int, value decimal.Decimal) {
	if utils.Between(index, define.Att_Percent_Begin, define.Att_Percent_End) {
		m.attPercent[index] = value
	}
}

func (m *AttManager) ModPercentAttValue(index int, value decimal.Decimal) {
	if utils.Between(index, define.Att_Percent_Begin, define.Att_Percent_End) {
		m.attPercent[index] = m.attPercent[index].Add(value)

		if m.attPercent[index].LessThan(decimal.NewFromInt32(0)) {
			m.attPercent[index] = decimal.NewFromInt32(0)
		}
	}
}

func (m *AttManager) Reset() {
	for k := range m.attBase {
		m.attBase[k] = decimal.NewFromInt32(0)
	}
	for k := range m.attPercent {
		m.attPercent[k] = decimal.NewFromInt32(0)
	}
	for k := range m.attFinal {
		m.attFinal[k] = decimal.NewFromInt32(0)
	}

	attEntry, ok := auto.GetAttEntry(m.baseAttId)
	if !ok {
		return
	}

	// 基础属性
	m.attBase[define.Att_AtkBase] = NewAttValue(attEntry.Atk)
	m.attBase[define.Att_ArmorBase] = NewAttValue(attEntry.Armor)
	m.attBase[define.Att_HealBase] = NewAttValue(attEntry.Heal)
	m.attBase[define.Att_MoveSpeedBase] = NewAttValue(attEntry.MoveSpeed)
	m.attBase[define.Att_AtbSpeedBase] = NewAttValue(attEntry.AtbSpeed)
	m.attBase[define.Att_MaxHPBase] = NewAttValue(attEntry.MaxHP)

	// 百分比属性
	m.attPercent[define.Att_AtkPercent] = NewAttValue(attEntry.AtkPercent)
	m.attPercent[define.Att_ArmorPercent] = NewAttValue(attEntry.ArmorPercent)
	m.attPercent[define.Att_HealPercent] = NewAttValue(attEntry.HealPercent)
	m.attPercent[define.Att_MoveSpeedPercent] = NewAttValue(attEntry.MoveSpeedPercent)
	m.attPercent[define.Att_AtbSpeedPercent] = NewAttValue(attEntry.AtbSpeedPercent)
	m.attPercent[define.Att_MaxHPPercent] = NewAttValue(attEntry.MaxHPPercent)

	// 最终属性
	m.attFinal[define.Att_DmgInc] = NewAttValue(attEntry.DmgInc)
	m.attFinal[define.Att_Crit] = NewAttValue(attEntry.Crit)
	m.attFinal[define.Att_CritInc] = NewAttValue(attEntry.CritInc)
	m.attFinal[define.Att_Tenacity] = NewAttValue(attEntry.Tenacity)
	m.attFinal[define.Att_RealDmg] = NewAttValue(attEntry.RealDmg)
	m.attFinal[define.Att_EffectHit] = NewAttValue(attEntry.EffectHit)
	m.attFinal[define.Att_EffectResist] = NewAttValue(attEntry.EffectResist)
	m.attFinal[define.Att_MaxMP] = NewAttValue(attEntry.MaxMP)
	m.attFinal[define.Att_GenMP] = NewAttValue(attEntry.GenMP)
	m.attFinal[define.Att_MaxRage] = NewAttValue(attEntry.Rage)
	m.attFinal[define.Att_GenRagePercent] = NewAttValue(attEntry.GenRagePercent)
	m.attFinal[define.Att_InitRage] = NewAttValue(attEntry.InitRage)
	m.attFinal[define.Att_Hit] = NewAttValue(attEntry.Hit)
	m.attFinal[define.Att_Dodge] = NewAttValue(attEntry.Dodge)
	m.attFinal[define.Att_MoveScope] = NewAttValue(attEntry.MoveScope)

	for n := 0; n < len(attEntry.DmgOfType); n++ {
		m.attFinal[define.Att_DmgTypeBegin+n] = NewAttValue(attEntry.DmgOfType[n])
	}

	for n := 0; n < len(attEntry.ResOfType); n++ {
		m.attFinal[define.Att_ResTypeBegin+n] = NewAttValue(attEntry.ResOfType[n])
	}
}

func (m *AttManager) CalcAtt() {
	// final = base * (1 + percent)
	m.attFinal[define.Att_Atk] = m.attBase[define.Att_AtkBase].Mul(m.attPercent[define.Att_AtkPercent].Add(decimal.NewFromInt32(1)))

	m.attFinal[define.Att_Armor] = m.attBase[define.Att_ArmorBase].Mul(m.attPercent[define.Att_ArmorPercent].Add(decimal.NewFromInt32(1)))

	m.attFinal[define.Att_Heal] = m.attBase[define.Att_HealBase].Mul(m.attPercent[define.Att_HealPercent].Add(decimal.NewFromInt32(1)))

	m.attFinal[define.Att_MoveSpeed] = m.attBase[define.Att_MoveSpeedBase].Mul(m.attPercent[define.Att_MoveSpeedPercent].Add(decimal.NewFromInt32(1)))

	m.attFinal[define.Att_AtbSpeed] = m.attBase[define.Att_AtbSpeedBase].Mul(m.attPercent[define.Att_AtbSpeedPercent].Add(decimal.NewFromInt32(1)))

	m.attFinal[define.Att_MaxHP] = m.attBase[define.Att_MaxHPBase].Mul(m.attPercent[define.Att_MaxHPPercent].Add(decimal.NewFromInt32(1)))
}

func (m *AttManager) ExportInt32() []int32 {
	att := make([]int32, len(m.attFinal))
	for n := define.Att_Begin; n < define.Att_End; n++ {
		att[n] = int32(m.attFinal[n].Round(0).IntPart())
	}

	return att
}

func (m *AttManager) ModAttManager(r *AttManager) *AttManager {
	for n := define.Att_Base_Begin; n < define.Att_Base_End; n++ {
		m.attBase[n] = m.attBase[n].Add(r.attBase[n])
	}
	for n := define.Att_Percent_Begin; n < define.Att_Percent_End; n++ {
		m.attPercent[n] = m.attPercent[n].Add(r.attPercent[n])
	}
	for n := define.Att_Begin; n < define.Att_End; n++ {
		m.attFinal[n] = m.attFinal[n].Add(r.attFinal[n])
	}
	return m
}

func (m *AttManager) Mul(d2 decimal.Decimal) *AttManager {
	for n := define.Att_Base_Begin; n < define.Att_Base_End; n++ {
		m.attBase[n] = m.attBase[n].Mul(d2)
	}
	for n := define.Att_Percent_Begin; n < define.Att_Percent_End; n++ {
		m.attPercent[n] = m.attPercent[n].Mul(d2)
	}
	for n := define.Att_Begin; n < define.Att_End; n++ {
		m.attFinal[n] = m.attFinal[n].Mul(d2)
	}
	return m
}

func (m *AttManager) Add(d2 decimal.Decimal) *AttManager {
	for n := define.Att_Base_Begin; n < define.Att_Base_End; n++ {
		m.attBase[n] = m.attBase[n].Add(d2)
	}
	for n := define.Att_Percent_Begin; n < define.Att_Percent_End; n++ {
		m.attPercent[n] = m.attPercent[n].Add(d2)
	}
	for n := define.Att_Begin; n < define.Att_End; n++ {
		m.attFinal[n] = m.attFinal[n].Add(d2)
	}
	return m
}

func (m *AttManager) Round() *AttManager {
	for n := define.Att_Base_Begin; n < define.Att_Base_End; n++ {
		m.attBase[n] = m.attBase[n].Round(0)
	}
	for n := define.Att_Percent_Begin; n < define.Att_Percent_End; n++ {
		m.attPercent[n] = m.attPercent[n].Round(0)
	}
	for n := define.Att_Begin; n < define.Att_End; n++ {
		m.attFinal[n] = m.attFinal[n].Round(4)
	}
	return m
}
