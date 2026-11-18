package player

import (
	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/services/game/loot"
	"github.com/east-eden/server/utils"
)

// GrantLootById 按 kof 的 Loot 表发放奖励到物品账本。
//
// 与旧的 CostLootManager 不同，kof 的 Loot 支持权重/概率/子掉落，
// 这里统一用 loot.CalcLootResult 展开后再入账。
func (p *Player) GrantLootById(lootId int32, reason string) {
	if lootId <= 0 {
		return
	}

	itemIds, counts, err := loot.CalcLootResult(lootId)
	if err != nil {
		utils.ErrPrint(err, "CalcLootResult failed", lootId, reason)
		return
	}

	for i, itemId := range itemIds {
		if i >= len(counts) || counts[i] <= 0 {
			continue
		}

		if _, err := p.ItemManager().OnItemAdd(itemId, int64(counts[i]), reason); err != nil {
			utils.ErrPrint(err, "grant loot item failed", lootId, itemId, reason)
		}
	}
}

// GrantAwardList 发放 IdCount 形式的固定奖励（关卡目标奖、任务奖）
func (p *Player) GrantAwardList(awards []*config.IdCount, reason string) {
	for _, award := range awards {
		if award == nil || award.Id <= 0 || award.Count <= 0 {
			continue
		}

		if _, err := p.ItemManager().OnItemAdd(award.Id, int64(award.Count), reason); err != nil {
			utils.ErrPrint(err, "grant award failed", award.Id, award.Count, reason)
		}
	}
}

// GrantHeroChip 按英雄配置 id 发放碎片（GM/补偿用，内部换算成 chipId）
func (p *Player) GrantHeroChip(heroConfigId int32, num int32) {
	entry, ok := config.GetHeroEntry(heroConfigId)
	if !ok || entry.ChipId <= 0 || num <= 0 {
		return
	}

	if err := p.HeroManager().AddHeroChip(entry.ChipId, int64(num)); err != nil {
		utils.ErrPrint(err, "grant hero chip failed", heroConfigId, num)
	}
}
