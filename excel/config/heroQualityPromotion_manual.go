package config

import "github.com/spf13/cast"

// GetNextQualityPromotionEntry 查找下一档英雄晋升配置。
//
// 规则与 kof 一致：先找同品质的下一个晋升等级，没有就取下一品质的 0 级。
func GetNextQualityPromotionEntry(rarity Rarity, level int32) (*HeroQualityPromotionEntry, bool) {
	if _, found := GetHeroQualityPromotionEntry(rarity.String(), cast.ToString(level)); !found {
		return nil, false
	}

	if next, found := GetHeroQualityPromotionEntry(rarity.String(), cast.ToString(level+1)); found {
		return next, true
	}

	if next, found := GetHeroQualityPromotionEntry((rarity + 1).String(), "0"); found {
		return next, true
	}

	return nil, false
}
