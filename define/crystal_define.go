package define

const (
	Crystal_TypeBegin = iota // 魂石类型
	Crystal_TypeA     = iota - 1
	Crystal_TypeB
	Crystal_TypeC
	Crystal_TypeEnd
)

type Crystal_PosType int32

const (
	Crystal_PosBegin Crystal_PosType = iota     // 魂石位置
	Crystal_Pos1     Crystal_PosType = iota - 1 // 1号位
	Crystal_Pos2                                // 2号位
	Crystal_Pos3                                // 3号位
	Crystal_Pos4                                // 4号位
	Crystal_Pos5                                // 5号位
	Crystal_Pos6                                // 6号位
	Crystal_PosEnd
)

type Crystal_AttType int32

const (
	Crystal_AttTypeBegin Crystal_AttType = iota
	Crystal_AttTypeMain  Crystal_AttType = iota - 1 // 晶石主属性
	Crystal_AttTypeVice                             // 晶石副属性
	Crystal_AttTypeEnd
)

const (
	Crystal_MaxStar = 5 // 魂石最高星级
	Crystal_MainAtt = 1 // 魂石主属性一条
	Crystal_ViceAtt = 4 // 魂石副属性五条
	Crystal_AttNum  = Crystal_MainAtt + Crystal_ViceAtt
)
