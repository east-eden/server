package define

// 消耗和掉落类型

const (
	CostLoot_Begin    = iota
	CostLoot_Item     = iota - 1 // 物品
	CostLoot_Token               // 代币
	CostLoot_Hero                // 英雄卡牌
	CostLoot_Player              // 玩家经验
	CostLoot_Fragment            // 卡牌碎片

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
