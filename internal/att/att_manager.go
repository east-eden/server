package att

import (
	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/utils"
)

type AttManager struct {
	baseAttId int32 // 属性id

	attBase     [define.AttNum]int32     // 基础属性32位
	attBasePlus [define.PlusAttNum]int64 // 基础属性64位

	attFinal     [define.AttNum]int32     // 计算后最终属性32位
	attFinalPlus [define.PlusAttNum]int64 // 计算后最终属性64位
}

func NewAttManager(attId int32) *AttManager {
	m := &AttManager{baseAttId: attId}

	m.Reset()

	return m
}

func (m *AttManager) SetBaseAttId(attId int32) {
	m.baseAttId = attId
	m.Reset()
}

func (m *AttManager) GetAttValue(index int) int {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		return int(m.attFinal[index])
	} else if utils.Between(index, define.Att_Plus_Begin, define.Att_Plus_End) {
		return int(m.attFinalPlus[index])
	} else {
		return 0
	}
}

func (m *AttManager) SetAttValue(index int, value int) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attFinal[index] = int32(value)
	} else if utils.Between(index, define.Att_Plus_Begin, define.Att_Plus_End) {
		m.attFinalPlus[index] = int64(value)
	} else {
	}
}

func (m *AttManager) ModAttValue(index int, value int) {
	if utils.Between(index, define.Att_Begin, define.Att_End) {
		m.attFinal[index] += int32(value)
	} else if utils.Between(index, define.Att_Plus_Begin, define.Att_Plus_End) {
		m.attFinalPlus[index] += int64(value)
	} else {
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

	m.attBase[define.Att_Str] = int32(attEntry.Str)
	m.attBase[define.Att_Agl] = int32(attEntry.Agl)
	m.attBase[define.Att_Con] = int32(attEntry.Con)
	m.attBase[define.Att_Int] = int32(attEntry.Int)
	m.attBase[define.Att_AtkSpeed] = int32(attEntry.AtkSpeed)

	m.attBase[define.Att_Atk] = int32(attEntry.Atk)
	m.attBase[define.Att_Def] = int32(attEntry.Def)
	m.attBase[define.Att_CriProb] = int32(attEntry.CriProb)
	m.attBase[define.Att_CriDmg] = int32(attEntry.CriDmg)
	m.attBase[define.Att_EffectHit] = int32(attEntry.EffectHit)
	m.attBase[define.Att_EffectResist] = int32(attEntry.EffectResist)
	m.attBase[define.Att_ConPercent] = int32(attEntry.ConPercent)
	m.attBase[define.Att_AtkPercent] = int32(attEntry.AtkPercent)
	m.attBase[define.Att_DefPercent] = int32(attEntry.DefPercent)

	m.attBasePlus[define.Att_Plus_MaxHP-define.Att_Plus_Begin] = int64(attEntry.MaxHP)
	m.attBasePlus[define.Att_Plus_MaxMP-define.Att_Plus_Begin] = int64(attEntry.MaxMP)
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

	m.attFinal[define.Att_Plus_MaxHP-define.Att_Plus_Begin] = m.attBase[define.Att_Con] + 10000
	m.attFinal[define.Att_Plus_MaxMP-define.Att_Plus_Begin] = m.attBase[define.Att_Int] + 1000
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

func (m *AttManager) ModBaseAtt(idx int, value int) {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		m.attBase[idx] += int32(value)
	} else if utils.Between(idx, define.Att_Plus_Begin, define.Att_Plus_End) {
		m.attBasePlus[idx-define.Att_Plus_Begin] += int64(value)
	} else {
		return
	}
}

func (m *AttManager) SetBaseAtt(idx int, value int) {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		m.attBase[idx] = int32(value)
	} else if utils.Between(idx, define.Att_Plus_Begin, define.Att_Plus_End) {
		m.attBasePlus[idx-define.Att_Plus_Begin] = int64(value)
	} else {
		return
	}
}

func (m *AttManager) ModAttManager(r *AttManager) {
	for k := range m.attBase {
		m.attBase[k] += r.attFinal[k]
	}
}
