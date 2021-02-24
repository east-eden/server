package define

const (
	Scene_TypeBegin = iota
	Scene_TypeStage = iota - 1 // 关卡场景
	Scene_TypeArena            // 竞技场场景
	Scene_TypeEnd
)

// 阵营
type SceneCampType int

const (
	Scene_Camp_Begin   SceneCampType = iota
	Scene_Camp_Attack  SceneCampType = iota - 1 // 0 进攻阵营
	Scene_Camp_Defence                          // 1 防守阵营
	Scene_Camp_End
)

const (
	Scene_MaxNumPerCombat = 5000 // 每台combat最多跑5000个场景
	Scene_MaxUnitPerScene = 1000 // 每个scene最多跑1000个unit
)

// 朝向
type Vector2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// 场景属性表
// type SceneEntry struct {
// 	ID          int32  `json:"_id"`
// 	Name        string `json:"Name"`
// 	Type        int32  `json:"Type"`
// 	UnitGroupID int32  `json:"UnitGroupID"`
// }

// // 怪物组属性表
// type UnitGroupEntry struct {
// 	ID         int32    `json:"_id"`
// 	Name       string   `json:"Name"`
// 	UnitTypeID []int32  `json:"UnitTypeID"`
// 	Position   []string `json:"Position"`
// }

// // unit属性表
// type UnitEntry struct {
// 	ID        int32    `json:"_id"`
// 	AttrName  []string `json:"AttrName"`
// 	AttrValue []int64  `json:"AttrValue"`
// 	Race      int32    `json:"Race"`
// }
