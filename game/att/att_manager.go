package att

import (
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
)

type AttManager struct {
	BaseAttID int32 // 属性id

	AttFinal [define.Att_End]int64 //计算后最终属性
	AttBase  [define.Att_End]int64
	AttMod   [define.Att_End]int64
}

func NewAttManager(attID int32) *AttManager {
	m := &AttManager{BaseAttID: attID}

	m.Reset()

	return m
}

func (m *AttManager) GetAttValue(index int32) int64 {
	if index < 0 || index >= define.Att_End {
		return 0
	}

	return m.AttFinal[index]
}

func (m *AttManager) Reset() {
	for k, _ := range m.AttFinal {
		m.AttFinal[k] = 0
	}

	attEntry := entries.GetAttEntry(m.BaseAttID)
	if attEntry == nil {
		return
	}

	m.AttBase[define.Att_Str] = attEntry.Str
	m.AttBase[define.Att_Agl] = attEntry.Agl
	m.AttBase[define.Att_Con] = attEntry.Con
	m.AttBase[define.Att_Int] = attEntry.Int
	m.AttBase[define.Att_AtkSpeed] = attEntry.AtkSpeed

	m.AttBase[define.Att_MaxHP] = attEntry.MaxHP
	m.AttBase[define.Att_MaxMP] = attEntry.MaxMP
	m.AttBase[define.Att_Atk] = attEntry.Atk
	m.AttBase[define.Att_Def] = attEntry.Def
	m.AttBase[define.Att_CriProb] = attEntry.CriProb
	m.AttBase[define.Att_CriDmg] = attEntry.CriDmg
	m.AttBase[define.Att_EffectHit] = attEntry.EffectHit
	m.AttBase[define.Att_EffectResist] = attEntry.EffectResist
	m.AttBase[define.Att_ConPercent] = attEntry.ConPercent
	m.AttBase[define.Att_AtkPercent] = attEntry.AtkPercent
	m.AttBase[define.Att_DefPercent] = attEntry.DefPercent
}

func (m *AttManager) CalcAtt() {
	for k, _ := range m.AttFinal {
		m.AttFinal[k] = 0
	}

	m.AttFinal[define.Att_Str] = m.AttBase[define.Att_Str]
	m.AttFinal[define.Att_Agl] = m.AttBase[define.Att_Agl]
	m.AttFinal[define.Att_Con] = m.AttBase[define.Att_Con]
	m.AttFinal[define.Att_Int] = m.AttBase[define.Att_Int]
	m.AttFinal[define.Att_AtkSpeed] = m.AttBase[define.Att_AtkSpeed]

	m.AttFinal[define.Att_MaxHP] = m.AttBase[define.Att_Con] + 10000
	m.AttFinal[define.Att_MaxMP] = m.AttBase[define.Att_Int] + 1000
	m.AttFinal[define.Att_Atk] = m.AttBase[define.Att_Str] + m.AttBase[define.Att_Agl]
	m.AttFinal[define.Att_Def] = m.AttBase[define.Att_Con] + 1000
	m.AttFinal[define.Att_CriProb] = m.AttBase[define.Att_CriProb]
	m.AttFinal[define.Att_CriDmg] = m.AttBase[define.Att_CriDmg]
	m.AttFinal[define.Att_EffectHit] = m.AttBase[define.Att_EffectHit]
	m.AttFinal[define.Att_EffectResist] = m.AttBase[define.Att_EffectResist]
	m.AttFinal[define.Att_ConPercent] = m.AttBase[define.Att_ConPercent]
	m.AttFinal[define.Att_AtkPercent] = m.AttBase[define.Att_AtkPercent]
	m.AttFinal[define.Att_DefPercent] = m.AttBase[define.Att_DefPercent]
}

func (m *AttManager) ModBaseAtt(idx int32, value int64) {
	if idx < define.Att_Begin || idx >= define.Att_End {
		return
	}

	m.AttBase[idx] += value
}

func (m *AttManager) ModAttManager(r *AttManager) {
	for k, _ := range m.AttBase {
		m.AttBase[k] += r.AttFinal[k]
	}
}
