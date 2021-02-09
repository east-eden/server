package auto

import "bitbucket.org/east-eden/server/define"

func GetGlobalConfig() (*GlobalConfigEntry, bool) {
	return GetGlobalConfigEntry(1)
}

// 获取背包容量上限
func (g *GlobalConfigEntry) GetItemContainerSize(tp define.ContainerType) int {
	switch tp {
	case define.Container_Material:
		return int(g.MaterialContainerMax)
	case define.Container_Equip:
		return int(g.EquipContainerMax)
	case define.Container_Crystal:
		return int(g.CrystalContainerMax)
	default:
		return 0
	}
}
