package att

import (
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
)

type AttManager struct {
	BaseAttID int32 // 属性id

	AttFinal   [define.Att_End]int64   //一级属性
	AttExFinal [define.AttEx_End]int64 //二级属性

	AttBase   [define.Att_End]int64
	AttExBase [define.AttEx_End]int64
	AttMod    [define.Att_End]int64
	AttExMod  [define.AttEx_End]int64
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

func (m *AttManager) GetAttExValue(index int32) int64 {
	if index < 0 || index >= define.AttEx_End {
		return 0
	}

	return m.AttExFinal[index]
}

func (m *AttManager) Reset() {
	for k, _ := range m.AttFinal {
		m.AttFinal[k] = 0
	}

	for k, _ := range m.AttExFinal {
		m.AttExFinal[k] = 0
	}

	attEntry := global.GetAttEntry(m.BaseAttID)
	if attEntry == nil {
		return
	}

	m.AttBase[define.Att_Str] = attEntry.Str
	m.AttBase[define.Att_Agl] = attEntry.Agl
	m.AttBase[define.Att_Con] = attEntry.Con
	m.AttBase[define.Att_Int] = attEntry.Int
	m.AttBase[define.Att_AtkSpeed] = attEntry.AtkSpeed

	m.AttExBase[define.AttEx_MaxHP] = attEntry.MaxHP
	m.AttExBase[define.AttEx_MaxMP] = attEntry.MaxMP
	m.AttExBase[define.AttEx_Atn] = attEntry.Atn
	m.AttExBase[define.AttEx_Def] = attEntry.Def
	m.AttExBase[define.AttEx_Ats] = attEntry.Ats
	m.AttExBase[define.AttEx_Adf] = attEntry.Adf
}

func (m *AttManager) CalcAtt() {
	for k, _ := range m.AttFinal {
		m.AttFinal[k] = 0
	}

	for k, _ := range m.AttExFinal {
		m.AttExFinal[k] = 0
	}

	m.AttFinal[define.Att_Str] = m.AttBase[define.Att_Str]
	m.AttFinal[define.Att_Agl] = m.AttBase[define.Att_Agl]
	m.AttFinal[define.Att_Con] = m.AttBase[define.Att_Con]
	m.AttFinal[define.Att_Int] = m.AttBase[define.Att_Int]
	m.AttFinal[define.Att_AtkSpeed] = m.AttBase[define.Att_AtkSpeed]

	m.AttExFinal[define.AttEx_MaxHP] = m.AttExBase[define.Att_Con] + 10000
	m.AttExFinal[define.AttEx_MaxMP] = m.AttExBase[define.Att_Int] + 1000
	m.AttExFinal[define.AttEx_Atn] = m.AttExBase[define.Att_Str] + m.AttExBase[define.Att_Agl]
	m.AttExFinal[define.AttEx_Def] = m.AttExBase[define.Att_Con] + 1000
	m.AttExFinal[define.AttEx_Ats] = m.AttExBase[define.Att_Int]
	m.AttExFinal[define.AttEx_Adf] = m.AttExBase[define.Att_Con] + m.AttExBase[define.Att_Int]
}

func (m *AttManager) ModBaseAtt(idx int32, value int64) {
	if idx < define.Att_Begin || idx >= define.Att_End {
		return
	}

	m.AttBase[idx] += value
}

func (m *AttManager) ModBaseAttEx(idx int32, value int64) {
	if idx < define.AttEx_Begin || idx >= define.AttEx_End {
		return
	}

	m.AttExBase[idx] += value
}

func (m *AttManager) ModAttManager(r *AttManager) {
	m.Reset()

	for k, _ := range m.AttBase {
		m.AttBase[k] += r.AttBase[k]
	}

	for k, _ := range m.AttExBase {
		m.AttExBase[k] += r.AttExBase[k]
	}
}
