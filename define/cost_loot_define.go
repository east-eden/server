package define

import "fmt"

const (
	CostLoot_Invalid = -1
	CostLoot_Start   = iota
	CostLoot_Item    = iota - 1
	CostLoot_Token
	CostLoot_Hero
	CostLoot_Player
	CostLoot_Fragment

	CostLoot_Blade
	CostLoot_Rune

	CostLoot_End
)

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
