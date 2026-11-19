package config

import "sort"

// activityBattlePassLevels 按活动 id 聚合战斗通行证阶段，按 stage 升序
var activityBattlePassLevels map[string][]*ActivityBattlePassEntry

func buildActivityBattlePassIndex() {
	if activityBattlePassLevels != nil {
		return
	}

	m := make(map[string][]*ActivityBattlePassEntry, 32)
	for _, e := range activityBattlePassEntries.Arrays {
		m[e.Id] = append(m[e.Id], e)
	}

	for _, list := range m {
		sort.Slice(list, func(i, j int) bool { return list[i].Stage < list[j].Stage })
	}

	activityBattlePassLevels = m
}

// GetActivityBattlePassLevels 返回某活动的通行证阶段，按 stage 升序
func GetActivityBattlePassLevels(activityId string) []*ActivityBattlePassEntry {
	buildActivityBattlePassIndex()
	return activityBattlePassLevels[activityId]
}

// ResetActivityBattlePassIndex 配置重载后重建索引
func ResetActivityBattlePassIndex() {
	activityBattlePassLevels = nil
}
