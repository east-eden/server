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
	baseAttId int32 // 基础属性id
	// attBase    [define.AttNum]int32 // 基础属性
	// attFlat    [define.AttNum]int32 // 平值加成
	// attPercent [define.AttNum]int32 // 百分比加成

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
		// m.attBase[k] = 0
		// m.attFlat[k] = 0
		// m.attPercent[k] = 0
		m.attFinal[k] = 0
	}

	attEntry, ok := auto.GetAttEntry(m.baseAttId)
	if !ok {
		return
	}

	m.attFinal[define.Att_Atk] = attEntry.Atk
	m.attFinal[define.Att_AtkPercent] = int32(attEntry.AtkPercent)
	m.attFinal[define.Att_Armor] = attEntry.Armor
	m.attFinal[define.Att_ArmorPercent] = int32(attEntry.ArmorPercent)
	m.attFinal[define.Att_DmgInc] = attEntry.DmgInc
	m.attFinal[define.Att_Crit] = attEntry.Crit
	m.attFinal[define.Att_CritInc] = attEntry.CritInc
	m.attFinal[define.Att_Tenacity] = attEntry.Tenacity
	m.attFinal[define.Att_Heal] = attEntry.Heal
	m.attFinal[define.Att_HealPercent] = int32(attEntry.HealPercent)
	m.attFinal[define.Att_RealDmg] = attEntry.RealDmg
	m.attFinal[define.Att_MoveSpeed] = int32(attEntry.MoveSpeed)
	m.attFinal[define.Att_MoveSpeedPercent] = int32(attEntry.MoveSpeedPercent)
	m.attFinal[define.Att_AtbSpeed] = int32(attEntry.AtbSpeed)
	m.attFinal[define.Att_AtbSpeedPercent] = int32(attEntry.AtbSpeedPercent)
	m.attFinal[define.Att_EffectHit] = attEntry.EffectHit
	m.attFinal[define.Att_EffectResist] = attEntry.EffectResist
	m.attFinal[define.Att_MaxHP] = attEntry.MaxHP
	m.attFinal[define.Att_MaxHPPercent] = int32(attEntry.MaxHPPercent)
	m.attFinal[define.Att_MaxMP] = attEntry.MaxMP
	m.attFinal[define.Att_GenMP] = attEntry.GenMP
	m.attFinal[define.Att_Rage] = attEntry.Rage
	m.attFinal[define.Att_GenRagePercent] = int32(attEntry.GenRagePercent)
	m.attFinal[define.Att_InitRage] = attEntry.InitRage
	m.attFinal[define.Att_Hit] = attEntry.Hit
	m.attFinal[define.Att_Dodge] = attEntry.Dodge
	m.attFinal[define.Att_MoveScope] = int32(attEntry.MoveScope)
	m.attFinal[define.Att_MoveTime] = int32(attEntry.MoveTime)

	for n := 0; n < len(attEntry.DmgOfType); n++ {
		m.attFinal[define.Att_DmgTypeBegin+n] = int32(attEntry.DmgOfType[n])
	}

	for n := 0; n < len(attEntry.ResOfType); n++ {
		m.attFinal[define.Att_ResTypeBegin+n] = int32(attEntry.ResOfType[n])
	}
}

func (m *AttManager) CalcAtt() {
	// for n := range m.attFinal {
	// 	m.attFinal[n] = 0
	// }

	// for n := define.Att_Begin; n < define.Att_End; n++ {
	// 	value64 := float64(m.attBase[n]+m.attFlat[n]) * float64(float64(define.PercentBase+m.attPercent[n])/float64(define.PercentBase))
	// 	m.attFinal[n] = int32(utils.Round(value64))
	// }
}

// func (m *AttManager) ModBaseAtt(idx int, value int32) {
// 	if utils.Between(idx, define.Att_Begin, define.Att_End) {
// 		m.attBase[idx] += value
// 	} else {
// 		return
// 	}
// }

// func (m *AttManager) SetBaseAtt(idx int, value int32) {
// 	if utils.Between(idx, define.Att_Begin, define.Att_End) {
// 		m.attBase[idx] = value
// 	} else {
// 		return
// 	}
// }

// func (m *AttManager) GetBaseAtt(idx int) int32 {
// 	if utils.Between(idx, define.Att_Begin, define.Att_End) {
// 		return m.attBase[idx]
// 	} else {
// 		return 0
// 	}
// }

// func (m *AttManager) ModFlatAtt(idx int, value int32) {
// 	if utils.Between(idx, define.Att_Begin, define.Att_End) {
// 		m.attFlat[idx] += value
// 	} else {
// 		return
// 	}
// }

// func (m *AttManager) ModPercentAtt(idx int, value int32) {
// 	if utils.Between(idx, define.Att_Begin, define.Att_End) {
// 		m.attPercent[idx] += value
// 	} else {
// 		return
// 	}
// }

// func (m *AttManager) SetPercentAtt(idx int, value int32) {
// 	if utils.Between(idx, define.Att_Begin, define.Att_End) {
// 		m.attPercent[idx] = value
// 	} else {
// 		return
// 	}
// }

// func (m *AttManager) GetPercentAtt(idx int) int32 {
// 	if utils.Between(idx, define.Att_Begin, define.Att_End) {
// 		return m.attPercent[idx]
// 	} else {
// 		return 0
// 	}
// }

func (m *AttManager) ModAttManager(r *AttManager) {
	for n := define.Att_Begin; n < define.Att_End; n++ {
		// m.attBase[n] += r.attBase[n]
		// m.attFlat[n] += r.attFlat[n]
		// m.attPercent[n] += r.attPercent[n]
		m.attFinal[n] += r.attFinal[n]
	}
}
