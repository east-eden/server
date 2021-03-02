package define

const (
	Crystal_TypeBegin = iota // 魂石类型
	Crystal_TypeA     = iota - 1
	Crystal_TypeB
	Crystal_TypeC
	Crystal_TypeEnd
)

const (
	Crystal_PosBegin = iota     // 魂石位置
	Crystal_Pos1     = iota - 1 // 1号位
	Crystal_Pos2                // 2号位
	Crystal_Pos3                // 3号位
	Crystal_Pos4                // 4号位
	Crystal_Pos5                // 5号位
	Crystal_Pos6                // 6号位
	Crystal_PosEnd
)

const (
	Crystal_MaxStar = 5 // 魂石最高星级
	Crystal_MainAtt = 1 // 魂石主属性一条
	Crystal_ViceAtt = 4 // 魂石副属性五条
	Crystal_AttNum  = Crystal_MainAtt + Crystal_ViceAtt
)
