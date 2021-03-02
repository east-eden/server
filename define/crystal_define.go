package define

const (
	Crystal_TypeBegin = iota // 魂石类型
	Crystal_TypeA     = iota - 1
	Crystal_TypeB
	Crystal_TypeC
	Crystal_TypeEnd
)

const (
	Crystal_PositionBegin = iota     // 魂石位置
	Crystal_Position1     = iota - 1 // 1号位
	Crystal_Position2                // 2号位
	Crystal_Position3                // 3号位
	Crystal_Position4                // 4号位
	Crystal_Position5                // 5号位
	Crystal_Position6                // 6号位
	Crystal_PositionEnd
)

const (
	Crystal_MaxStar = 6 // 魂石最高星级
	Crystal_MainAtt = 1 // 魂石主属性一条
	Crystal_ViceAtt = 5 // 魂石副属性五条
	Crystal_AttNum  = Crystal_MainAtt + Crystal_ViceAtt
)
