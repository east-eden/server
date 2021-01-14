package define

const (
	CostLoot_Item = iota
	CostLoot_Token
	CostLoot_Hero
	CostLoot_Player
	CostLoot_Blade
	CostLoot_Rune

	CostLoot_End
)

type CostLootObj interface {
	GetCostLootType() int32

	CanCost(int32, int32) error
	DoCost(int32, int32) error

	CanGain(int32, int32) error
	GainLoot(int32, int32) error
}
