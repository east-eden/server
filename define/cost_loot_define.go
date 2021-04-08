package define

import "fmt"

const (
	CostLoot_Invalid  int32 = -1
	CostLoot_Start    int32 = 0
	CostLoot_Item     int32 = 0
	CostLoot_Token    int32 = 1
	CostLoot_Hero     int32 = 2
	CostLoot_Player   int32 = 3
	CostLoot_Fragment int32 = 4

	CostLoot_End int32 = 5
)

// 掉落数据
type LootData struct {
	LootType int32 `bson:"loot_type" json:"loot_type"` // 掉落类型
	LootMisc int32 `bson:"loot_misc" json:"loot_misc"` // 掉落参数
	LootNum  int32 `bson:"loot_num" json:"loot_num"`   // 掉落数量
}

type CostLooter interface {
	GetCostLootType() int32

	CanCost(int32, int32) error
	DoCost(int32, int32) error

	CanGain(int32, int32) error
	GainLoot(int32, int32) error
}

type BaseCostLooter struct {
	CostLooter
}

func (bc *BaseCostLooter) GetCostLootType() int32 {
	return CostLoot_Invalid
}

func (bc *BaseCostLooter) CanCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("BaseCostLooter CanCost<id:%d> failed, wrong number<%d>", typeMisc, num)
	}

	return nil
}

func (bc *BaseCostLooter) DoCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("BaseCostLooter DoCost<id:%d> failed, wrong number<%d>", typeMisc, num)
	}

	return nil
}

func (bc *BaseCostLooter) CanGain(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("BaseCostLooter CanGain<id:%d> failed, wrong number<%d>", typeMisc, num)
	}

	return nil
}

func (bc *BaseCostLooter) GainLoot(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("BaseCostLooter GainLoot<id:%d> failed, wrong number<%d>", typeMisc, num)
	}

	return nil
}
