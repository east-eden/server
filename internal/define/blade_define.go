package define

const (
	Blade_MaxBlade = 4
)

// blade entry
type BladeEntry struct {
	ID        int32   `json:"id"`
	AttID     int32   `json:"att_id"`
	Quality   int32   `json:"quality"`
	SpellList []int32 `json:"spell_list"`
}
