package config

import "sort"

// battlePropertiesByID 按属性 id 分组并按等级升序排列，供英雄属性计算按等级取档。
// 生成代码只提供复合键（id_level）查询，这里补一个按 id 的索引。
var battlePropertiesByID map[int32][]*BattlePropertyEntry

// GetBattlePropertiesById 返回某个属性 id 下所有等级档位，按 Level 升序
func GetBattlePropertiesById(id int32) []*BattlePropertyEntry {
	if battlePropertiesByID == nil {
		m := make(map[int32][]*BattlePropertyEntry, 64)
		for _, e := range battlePropertyEntries.Arrays {
			m[e.Id] = append(m[e.Id], e)
		}

		for _, list := range m {
			sort.Slice(list, func(i, j int) bool { return list[i].Level < list[j].Level })
		}

		battlePropertiesByID = m
	}

	return battlePropertiesByID[id]
}

// ResetBattlePropertyIndex 配置重载后重建索引
func ResetBattlePropertyIndex() {
	battlePropertiesByID = nil
}
