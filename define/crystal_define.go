package define

const (
	Crystal_TypeBegin = iota // 魂石类型
	Crystal_TypeA     = iota - 1
	Crystal_TypeB
	Crystal_TypeC
	Crystal_TypeEnd
)

const (
	Crystal_PosBegin int32 = iota     // 魂石位置
	Crystal_Pos1     int32 = iota - 1 // 1号位
	Crystal_Pos2                      // 2号位
	Crystal_Pos3                      // 3号位
	Crystal_Pos4                      // 4号位
	Crystal_Pos5                      // 5号位
	Crystal_Pos6                      // 6号位
	Crystal_PosEnd
)

const (
	Crystal_AttTypeBegin int32 = iota
	Crystal_AttTypeMain  int32 = iota - 1 // 晶石主属性
	Crystal_AttTypeVice                   // 晶石副属性
	Crystal_AttTypeEnd
)

const (
	Crystal_MaxStar    = 5 // 魂石最高星级
	Crystal_MainAttNum = 1 // 魂石主属性1条
	Crystal_ViceAttNum = 4 // 魂石副属性4条
	Crystal_AttNum     = Crystal_MainAttNum + Crystal_ViceAttNum
)
