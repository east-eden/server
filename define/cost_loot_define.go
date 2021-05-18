package define

import (
	"fmt"

	pbGlobal "bitbucket.org/funplus/server/proto/global"
)

const (
	CostLoot_Invalid            int32 = -1
	CostLoot_Start              int32 = iota - 1
	CostLoot_Item               int32 = iota - 2 // 0 物品
	CostLoot_Token                               // 1 代币
	CostLoot_Hero                                // 2 英雄
	CostLoot_Player                              // 3 玩家
	CostLoot_HeroFragment                        // 4 英雄碎片
	CostLoot_CollectionFragment                  // 5 收集品碎片
	CostLoot_Collection                          // 6 收集品

	CostLoot_End
)

// 掉落数据
type LootData struct {
	LootType int32 `bson:"loot_type" json:"loot_type"` // 掉落类型
	LootMisc int32 `bson:"loot_misc" json:"loot_misc"` // 掉落参数
	LootNum  int32 `bson:"loot_num" json:"loot_num"`   // 掉落数量
}

func (d *LootData) GenPB() *pbGlobal.LootData {
	pb := &pbGlobal.LootData{
		Type: pbGlobal.LootType(d.LootType),
		Misc: d.LootMisc,
		Num:  d.LootNum,
	}

	return pb
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
