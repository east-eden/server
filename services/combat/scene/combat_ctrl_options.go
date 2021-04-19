package scene

type CombatCtrlOption func(*CombatCtrlOptions)
type CombatCtrlOptions struct {
	AtbValue int32 // 行动条
}

func DefaultCombatCtrlOptions() *CombatCtrlOptions {
	o := &CombatCtrlOptions{
		AtbValue: 0,
	}

	return o
}

func WithCombatCtrlAtbValue(value int32) CombatCtrlOption {
	return func(o *CombatCtrlOptions) {
		o.AtbValue = value
	}
}
