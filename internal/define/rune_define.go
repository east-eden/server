package define

const (
	Rune_TypeBegin = iota // 魂石类型
	Rune_TypeA     = iota - 1
	Rune_TypeB
	Rune_TypeC
	Rune_TypeEnd
)

const (
	Rune_PositionBegin = iota // 魂石位置
	Rune_Position0     = iota - 1
	Rune_Position1
	Rune_Position2
	Rune_Position3
	Rune_Position4
	Rune_Position5
	Rune_PositionEnd
)

const (
	Rune_MaxStar = 6 // 魂石最高星级
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
	ID    int32 `json:"_id"`
	AttID int32 `json:"AttID"`
}
