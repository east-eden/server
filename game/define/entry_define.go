package define

// hero entry
type HeroEntry struct {
	ID        int32   `json:"id"`
	AttID     int32   `json:"att_id"`
	Quality   int32   `json:"quality"`
	SpellList []int32 `json:"spell_list"`
}

// item entry
type ItemEntry struct {
	ID       int32 `json:"id"`
	ItemType int32 `json:"item_type"`
	Quality  int32 `json:"quality"`
	Price    int32 `json:"price"`
}
