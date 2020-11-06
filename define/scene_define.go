package define

const (
	Scene_TypeBegin = iota
	Scene_TypeStage = iota - 1 // 关卡场景
	Scene_TypeArena            // 竞技场场景
	Scene_TypeEnd
)

const (
	Scene_MaxNumPerCombat = 5000 // 每台combat最多跑5000个场景
	Scene_MaxUnitPerScene = 1000 // 每个scene最多跑1000个unit
)

// 场景属性表
type SceneEntry struct {
	ID          int32  `json:"_id"`
	Name        string `json:"Name"`
	Type        int32  `json:"Type"`
	UnitGroupID int32  `json:"UnitGroupID"`
}

// 怪物组属性表
type UnitGroupEntry struct {
	ID         int32    `json:"_id"`
	Name       string   `json:"Name"`
	UnitTypeID []int32  `json:"UnitTypeID"`
	Position   []string `json:"Position"`
}

// unit属性表
type UnitEntry struct {
	ID        int32    `json:"_id"`
	AttrName  []string `json:"AttrName"`
	AttrValue []int64  `json:"AttrValue"`
	Race      int32    `json:"Race"`
}
