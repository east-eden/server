package define

// talent entry
type TalentEntry struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	PrevID  int32  `json:"prev_id"`
	MutexID int32  `json:"mutex_id"`
	CostID  int32  `json:"cost_id"`
}
