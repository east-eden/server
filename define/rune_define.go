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

// 魂石属性表
type RuneEntry struct {
	ID      int32  `json:"_id"`
	Name    string `json:"Name"`
	Type    int32  `json:"Type"`
	Pos     int32  `json:"Pos"`
	Quality int32  `json:"Quality"`
	SuitID  int32  `json:"SuitID"`
}

// 魂石套装属性表
type RuneSuitEntry struct {
	ID          int32 `json:"_id"`
	Suit2_AttID int32 `json:"Suit2_AttID"`
	Suit3_AttID int32 `json:"Suit3_AttID"`
	Suit4_AttID int32 `json:"Suit4_AttID"`
	Suit5_AttID int32 `json:"Suit5_AttID"`
	Suit6_AttID int32 `json:"Suit6_AttID"`
}
