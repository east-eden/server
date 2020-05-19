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
	ID              int32  `json:"_id"`
	Name            string `json:"Name"`
	Type            int32  `json:"Type"`
	CreatureGroupID int32  `json:"CreatureGroupID"`
}
