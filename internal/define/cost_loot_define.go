package define

const (
	CostLoot_Item = iota
	CostLoot_Token
	CostLoot_Hero
	CostLoot_Blade

	CostLoot_End
)

// cost_loot entry
type CostLootEntry struct {
	ID       int32 `json:"id"`
	Type     int32 `json:"type"`
	TypeMisc int32 `json:"type_misc"`
	Num      int32 `json:"num"`
}

type CostLootObj interface {
	GetCostLootType() int32

	CanCost(int32, int32) error
	DoCost(int32, int32) error

	CanGain(int32, int32) error
	GainLoot(int32, int32) error
}
