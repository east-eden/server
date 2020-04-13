package att

import (
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
)

type AttManager struct {
	Atts   [define.Att_End]int64   //一级属性
	AttExs [define.AttEx_End]int64 //二级属性
}

func NewAttManager(attID int32) *AttManager {
	m := &AttManager{}

	attEntry := global.GetAttEntry(attID)
	if attEntry == nil {
		return m
	}

	m.Atts[define.Att_Str] = attEntry.Str
	m.Atts[define.Att_Agl] = attEntry.Agl
	m.Atts[define.Att_Con] = attEntry.Con
	m.Atts[define.Att_Int] = attEntry.Int
	m.Atts[define.Att_AtkSpeed] = attEntry.AtkSpeed

	m.AttExs[define.AttEx_MaxHP] = attEntry.MaxHP
	m.AttExs[define.AttEx_MaxMP] = attEntry.MaxMP
	m.AttExs[define.AttEx_Atn] = attEntry.Atn
	m.AttExs[define.AttEx_Def] = attEntry.Def
	m.AttExs[define.AttEx_Ats] = attEntry.Ats
	m.AttExs[define.AttEx_Adf] = attEntry.Adf

	return m
}

func (m *AttManager) GetAttValue(index int32) int64 {
	if index < 0 || index >= define.Att_End {
		return 0
	}

	return m.Atts[index]
}

func (m *AttManager) GetAttExValue(index int32) int64 {
	if index < 0 || index >= define.AttEx_End {
		return 0
	}

	return m.AttExs[index]
}

func (m *AttManager) CalcAtt() {
	for k, _ := range m.AttExs {
		m.AttExs[k] = 0
	}

	m.AttExs[define.AttEx_MaxHP] = m.Atts[define.Att_Con] + 10000
	m.AttExs[define.AttEx_MaxMP] = m.Atts[define.Att_Int] + 1000
	m.AttExs[define.AttEx_Atn] = m.Atts[define.Att_Str] + m.Atts[define.Att_Agl]
	m.AttExs[define.AttEx_Def] = m.Atts[define.Att_Con] + 1000
	m.AttExs[define.AttEx_Ats] = m.Atts[define.Att_Int]
	m.AttExs[define.AttEx_Adf] = m.Atts[define.Att_Con] + m.Atts[define.Att_Int]
}

func (m *AttManager) AddAtt(r *AttManager) {
}
