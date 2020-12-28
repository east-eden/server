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

	CanCost(int, int) error
	DoCost(int, int) error

	CanGain(int, int) error
	GainLoot(int, int) error
}
