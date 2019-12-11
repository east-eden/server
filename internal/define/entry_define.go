package define

// hero entry
type HeroEntry struct {
	ID        int32   `json:"id"`
	Name      string  `json:"name"`
	AttID     int32   `json:"att_id"`
	Quality   int32   `json:"quality"`
	SpellList []int32 `json:"spell_list"`
}

// item entry
type ItemEntry struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	ItemType int32  `json:"item_type"`
	Quality  int32  `json:"quality"`
	Price    int32  `json:"price"`
}

// token entry
type TokenEntry struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	MaxHold int32  `json:"max_hold"`
}

// talent entry
type TalentEntry struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	PrevID  int32  `json:"prev_id"`
	MutexID int32  `json:"mutex_id"`
	CostID  int32  `json:"cost_id"`
}
