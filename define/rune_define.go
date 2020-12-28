package define

const (
	Rune_TypeBegin = iota // 魂石类型
	Rune_TypeA     = iota - 1
	Rune_TypeB
	Rune_TypeC
	Rune_TypeEnd
)

const (
	Rune_PositionBegin = iota     // 魂石位置
	Rune_Position1     = iota - 1 // 1号位
	Rune_Position2                // 2号位
	Rune_Position3                // 3号位
	Rune_Position4                // 4号位
	Rune_Position5                // 5号位
	Rune_Position6                // 6号位
	Rune_PositionEnd
)

const (
	Rune_MaxStar = 6 // 魂石最高星级
	Rune_MainAtt = 1 // 魂石主属性一条
	Rune_ViceAtt = 5 // 魂石副属性五条
	Rune_AttNum  = Rune_MainAtt + Rune_ViceAtt
)
