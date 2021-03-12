package auto

import "bitbucket.org/funplus/server/define"

// 获取职业系数加成属性
func (e *HeroProfessionEntry) GetRatio(attType int) int32 {
	switch attType {
	case define.Att_Atk:
		return int32(e.AtkRatio)
	case define.Att_Armor:
		return int32(e.ArmorRatio)
	case define.Att_MaxHP:
		return int32(e.MaxHPRatio)
	default:
		return define.PercentBase
	}
}
