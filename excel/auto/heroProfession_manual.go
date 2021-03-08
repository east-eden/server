package auto

import "bitbucket.org/funplus/server/define"

// 获取职业系数加成属性
func (e *HeroProfessionEntry) GetRatio(attType int) int32 {
	switch attType {
	case define.Att_Atk:
		return e.AtkRatio
	case define.Att_Armor:
		return e.ArmorRatio
	case define.Att_MaxHP:
		return e.MaxHPRatio
	default:
		return define.PercentBase
	}
}
