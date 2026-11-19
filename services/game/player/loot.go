package player

import (
	"context"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/services/game/arena"
	"github.com/east-eden/server/services/game/loot"
	"github.com/east-eden/server/services/game/recharge"
	"github.com/east-eden/server/services/game/recruit"
	"github.com/east-eden/server/services/game/shop"
	"github.com/east-eden/server/store"
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

// GetItemCount 返回某道具数量（实现 shop.Owner 等宿主接口）
func (p *Player) GetItemCount(entryId int32) int64 {
	return p.ItemManager().GetItemCount(entryId)
}

// CostItems 批量扣除道具：先整体校验，再逐个扣除
func (p *Player) CostItems(costs []*config.IdCount, reason string) error {
	itemMgr := p.ItemManager()

	for _, c := range costs {
		if c == nil || c.Id <= 0 || c.Count <= 0 {
			continue
		}

		if err := itemMgr.OnCheckEnough(c.Id, int64(c.Count)); err != nil {
			return err
		}
	}

	for _, c := range costs {
		if c == nil || c.Id <= 0 || c.Count <= 0 {
			continue
		}

		if err := itemMgr.OnItemCost(c.Id, int64(c.Count), reason); err != nil {
			return err
		}
	}

	return nil
}

// GrantItems 批量发放道具
func (p *Player) GrantItems(items []*config.IdCount, reason string) {
	for _, it := range items {
		if it == nil || it.Id <= 0 || it.Count <= 0 {
			continue
		}

		if _, err := p.ItemManager().OnItemAdd(it.Id, int64(it.Count), reason); err != nil {
			utils.ErrPrint(err, "grant item failed", it.Id, it.Count, reason)
		}
	}
}

// SaveShops 持久化商店状态（实现 shop.Owner）
func (p *Player) SaveShops(data map[int32]*shop.ShopState) error {
	p.Shops = data

	return store.GetStore().UpdateFields(
		context.Background(),
		define.StoreType_Player,
		p.ID,
		map[string]any{"shops": data},
	)
}

// GrantLootAndReport 发放掉落并返回实际产出（实现 recruit.Owner）
func (p *Player) GrantLootAndReport(lootId int32, reason string) []recruit.Reward {
	rewards := recruit.CalcLoot(lootId)
	for _, r := range rewards {
		if _, err := p.ItemManager().OnItemAdd(r.ItemId, int64(r.Num), reason); err != nil {
			utils.ErrPrint(err, "grant recruit reward failed", r.ItemId, r.Num, reason)
		}
	}

	return rewards
}

// SaveRecruit 持久化招募状态（实现 recruit.Owner）
func (p *Player) SaveRecruit(state *recruit.State) error {
	p.Recruit = state

	return store.GetStore().UpdateFields(
		context.Background(),
		define.StoreType_Player,
		p.ID,
		map[string]any{"recruit": state},
	)
}

// SaveRecharge 持久化充值状态（实现 recharge.Owner）
func (p *Player) SaveRecharge(state *recharge.State) error {
	p.Recharge = state

	return store.GetStore().UpdateFields(
		context.Background(),
		define.StoreType_Player,
		p.ID,
		map[string]any{"recharge": state},
	)
}

// SaveArena 持久化竞技场状态（实现 arena.Owner）
func (p *Player) SaveArena(state *arena.State) error {
	p.Arena = state

	return store.GetStore().UpdateFields(
		context.Background(),
		define.StoreType_Player,
		p.ID,
		map[string]any{"arena": state},
	)
}
