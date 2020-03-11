package define

const (
	CostLoot_Item = iota
	CostLoot_Token
	CostLoot_Hero
	CostLoot_Player
	CostLoot_Blade

	CostLoot_End
)

// cost_loot entry
type CostLootEntry struct {
	ID   int32 `json:"_id"`
	Type int32 `json:"type"`
	Misc int32 `json:"misc"`
	Num  int32 `json:"num"`
}

type CostLootObj interface {
	GetCostLootType() int32

	CanCost(int32, int32) error
	DoCost(int32, int32) error

	CanGain(int32, int32) error
	GainLoot(int32, int32) error
}
