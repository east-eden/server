package auto

import (
	"bitbucket.org/funplus/server/define"
	"github.com/shopspring/decimal"
)

// 获取职业系数加成属性
func (e *HeroProfessionEntry) GetBaseRatio(attType int) decimal.Decimal {
	switch attType {
	case define.Att_AtkBase:
		return e.AtkRatio
	case define.Att_ArmorBase:
		return e.ArmorRatio
	case define.Att_MaxHPBase:
		return e.MaxHPRatio
	default:
		return decimal.NewFromInt(1)
	}
}

// 获取职业系数加成属性
func (e *HeroProfessionEntry) GetPercentRatio(attType int) decimal.Decimal {
	switch attType {
	case define.Att_AtkPercent:
		return e.AtkRatio
	case define.Att_ArmorPercent:
		return e.ArmorRatio
	case define.Att_MaxHPPercent:
		return e.MaxHPRatio
	default:
		return decimal.NewFromInt(1)
	}
}

// 获取职业系数加成属性
func (e *HeroProfessionEntry) GetFinalRatio(attType int) decimal.Decimal {
	switch attType {
	case define.Att_Atk:
		return e.AtkRatio
	case define.Att_Armor:
		return e.ArmorRatio
	case define.Att_MaxHP:
		return e.MaxHPRatio
	default:
		return decimal.NewFromInt(1)
	}
}
