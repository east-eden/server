package att

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
)

type AttManager struct {
	baseAttId int // 属性id

	attFinal [define.Att_End]int //计算后最终属性
	attBase  [define.Att_End]int
}

func NewAttManager(attId int) *AttManager {
	m := &AttManager{baseAttId: attId}

	m.Reset()

	return m
}

func (m *AttManager) SetBaseAttId(attId int) {
	m.baseAttId = attId
	m.Reset()
}

func (m *AttManager) GetAttValue(index int) int {
	if index < 0 || index >= define.Att_End {
		return 0
	}

	return m.attFinal[index]
}

func (m *AttManager) Reset() {
	for k := range m.attFinal {
		m.attFinal[k] = 0
	}

	attEntry, ok := auto.GetAttEntry(m.baseAttId)
	if !ok {
		return
	}

	m.attBase[define.Att_Str] = attEntry.Str
	m.attBase[define.Att_Agl] = attEntry.Agl
	m.attBase[define.Att_Con] = attEntry.Con
	m.attBase[define.Att_Int] = attEntry.Int
	m.attBase[define.Att_AtkSpeed] = attEntry.AtkSpeed

	m.attBase[define.Att_MaxHP] = attEntry.MaxHP
	m.attBase[define.Att_MaxMP] = attEntry.MaxMP
	m.attBase[define.Att_Atk] = attEntry.Atk
	m.attBase[define.Att_Def] = attEntry.Def
	m.attBase[define.Att_CriProb] = attEntry.CriProb
	m.attBase[define.Att_CriDmg] = attEntry.CriDmg
	m.attBase[define.Att_EffectHit] = attEntry.EffectHit
	m.attBase[define.Att_EffectResist] = attEntry.EffectResist
	m.attBase[define.Att_ConPercent] = attEntry.ConPercent
	m.attBase[define.Att_AtkPercent] = attEntry.AtkPercent
	m.attBase[define.Att_DefPercent] = attEntry.DefPercent
}

func (m *AttManager) CalcAtt() {
	for k := range m.attFinal {
		m.attFinal[k] = 0
	}

	m.attFinal[define.Att_Str] = m.attBase[define.Att_Str]
	m.attFinal[define.Att_Agl] = m.attBase[define.Att_Agl]
	m.attFinal[define.Att_Con] = m.attBase[define.Att_Con]
	m.attFinal[define.Att_Int] = m.attBase[define.Att_Int]
	m.attFinal[define.Att_AtkSpeed] = m.attBase[define.Att_AtkSpeed]

	m.attFinal[define.Att_MaxHP] = m.attBase[define.Att_Con] + 10000
	m.attFinal[define.Att_MaxMP] = m.attBase[define.Att_Int] + 1000
	m.attFinal[define.Att_Atk] = m.attBase[define.Att_Str] + m.attBase[define.Att_Agl]
	m.attFinal[define.Att_Def] = m.attBase[define.Att_Con] + 1000
	m.attFinal[define.Att_CriProb] = m.attBase[define.Att_CriProb]
	m.attFinal[define.Att_CriDmg] = m.attBase[define.Att_CriDmg]
	m.attFinal[define.Att_EffectHit] = m.attBase[define.Att_EffectHit]
	m.attFinal[define.Att_EffectResist] = m.attBase[define.Att_EffectResist]
	m.attFinal[define.Att_ConPercent] = m.attBase[define.Att_ConPercent]
	m.attFinal[define.Att_AtkPercent] = m.attBase[define.Att_AtkPercent]
	m.attFinal[define.Att_DefPercent] = m.attBase[define.Att_DefPercent]
}

func (m *AttManager) ModBaseAtt(idx int32, value int64) {
	if idx < define.Att_Begin || idx >= define.Att_End {
		return
	}

	m.attBase[idx] += int(value)
}

func (m *AttManager) SetBaseAtt(index int32, value int64) {
	if index < 0 || index >= define.Att_End {
		return
	}

	m.attBase[index] = int(value)
}

func (m *AttManager) ModAttManager(r *AttManager) {
	for k := range m.attBase {
		m.attBase[k] += r.attFinal[k]
	}
}
