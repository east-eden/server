package att

import (
	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	"bitbucket.org/east-eden/server/utils"
)

type AttManager struct {
	baseAttId int32 // 属性id

	attBase [define.AttNum]int32 // 基础属性
	attFlat [define.AttNum]int32 // 平值加成
	attProb [define.AttNum]int32 // 百分比加成

	attFinal [define.AttNum]int32 // 计算后最终属性32位
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
		m.attProb[k] = 0
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

}

func (m *AttManager) CalcAtt() {
	for k := range m.attFinal {
		m.attFinal[k] = 0
	}

	for k := 0; k < define.AttNum; k++ {
		m.attFinal[k] = (m.attBase[k] + m.attFlat[k]) * (define.AttPercentBase + m.attProb[k]) / define.AttPercentBase
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

func (m *AttManager) ModProbAtt(idx int, value int32) {
	if utils.Between(idx, define.Att_Begin, define.Att_End) {
		m.attProb[idx] += value
	} else {
		return
	}
}

func (m *AttManager) ModAttManager(r *AttManager) {
	for k := 0; k < define.AttNum; k++ {
		m.attFinal[k] += r.attFinal[k]
	}
}
