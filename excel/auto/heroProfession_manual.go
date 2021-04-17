package auto

import (
	"bitbucket.org/funplus/server/define"
	"github.com/shopspring/decimal"
)

// 获取职业系数加成属性
func (e *HeroProfessionEntry) GetRatio(attType int) decimal.Decimal {
	switch attType {
	case define.Att_AtkBase:
		return e.AtkRatio
	case define.Att_ArmorBase:
		return e.ArmorRatio
	case define.Att_MaxHPBase:
		return e.MaxHPRatio
	default:
		return decimal.NewFromInt(define.PercentBase)
	}
}
