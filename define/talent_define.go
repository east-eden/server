package define

// talent entry
type TalentEntry struct {
	ID         int32  `json:"_id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	LevelLimit int32  `json:"level_limit"`
	GroupID    int32  `json:"group_id"`
	CostID     int32  `json:"cost_id"`
}
