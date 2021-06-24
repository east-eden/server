package scene

import "github.com/shopspring/decimal"

type CombatCtrlOption func(*CombatCtrlOptions)
type CombatCtrlOptions struct {
	AtbValue decimal.Decimal // 行动条
}

func DefaultCombatCtrlOptions() *CombatCtrlOptions {
	o := &CombatCtrlOptions{}

	return o
}

func WithCombatCtrlAtbValue(value decimal.Decimal) CombatCtrlOption {
	return func(o *CombatCtrlOptions) {
		o.AtbValue = value
	}
}
