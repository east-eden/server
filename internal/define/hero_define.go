package define

const (
	Hero_MaxEquip = 4 // how many equips can put on per hero
)

// hero entry
type HeroEntry struct {
	ID        int32   `json:"id"`
	Name      string  `json:"name"`
	AttID     int32   `json:"att_id"`
	Quality   int32   `json:"quality"`
	SpellList []int32 `json:"spell_list"`
}
