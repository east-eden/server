package auto

import (
	"sort"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel"
	"github.com/rs/zerolog/log"
)

// 经验道具
type ExpItem struct {
	ItemTypeId int32 // 物品typeid
	Exp        int32 // 物品经验值
}

var (
	heroExpItemSection    = make([]*ExpItem, 0, 10) // 英雄经验道具档位
	equipExpItemSection   = make([]*ExpItem, 0, 10) // 装备经验道具档位
	crystalExpItemSection = make([]*ExpItem, 0, 10) // 晶石经验道具档位
)

func init() {
	excel.AddEntryManualLoader("GlobalConfig.xlsx", (*GlobalConfigEntries)(nil))
}

// ManualLoader
func (e *GlobalConfigEntries) ManualLoad(*excel.ExcelFileRaw) error {
	// 经验道具按经验值排序
	add := func(typeId int32, list []*ExpItem) []*ExpItem {
		itemEntry, ok := GetItemEntry(typeId)
		if !ok {
			log.Error().Caller().Int32("type_id", typeId).Msg("can not find item entry")
			return list
		}

		list = append(list, &ExpItem{
			ItemTypeId: typeId,
			Exp:        itemEntry.PublicMisc[0],
		})

		return list
	}

	g, _ := GetGlobalConfig()

	for _, typeId := range g.HeroExpItems {
		heroExpItemSection = add(typeId, heroExpItemSection)
	}
	sort.Slice(heroExpItemSection, func(i, j int) bool {
		return heroExpItemSection[i].Exp < heroExpItemSection[j].Exp
	})

	for _, typeId := range g.EquipExpItems {
		equipExpItemSection = add(typeId, equipExpItemSection)
	}
	sort.Slice(equipExpItemSection, func(i, j int) bool {
		return equipExpItemSection[i].Exp < equipExpItemSection[j].Exp
	})

	for _, typeId := range g.CrystalExpItems {
		crystalExpItemSection = add(typeId, crystalExpItemSection)
	}
	sort.Slice(crystalExpItemSection, func(i, j int) bool {
		return crystalExpItemSection[i].Exp < crystalExpItemSection[j].Exp
	})

	return nil
}

func GetGlobalConfig() (*GlobalConfigEntry, bool) {
	return GetGlobalConfigEntry(1)
}

// 获取背包容量上限
func (g *GlobalConfigEntry) GetItemContainerSize(tp int32) int {
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

// 获取英雄经验道具 -- 返还经验道具时用
func (g *GlobalConfigEntry) GetHeroExpItemByExp(exp int32) *ExpItem {
	var ret *ExpItem
	if exp <= 0 {
		return nil
	}

	for _, expItem := range heroExpItemSection {
		if exp < expItem.Exp {
			break
		}

		ret = expItem
	}

	return ret
}

// 获取装备经验道具 -- 返还经验道具时用
func (g *GlobalConfigEntry) GetEquipExpItemByExp(exp int32) *ExpItem {
	var ret *ExpItem
	if exp <= 0 {
		return nil
	}

	for _, expItem := range equipExpItemSection {
		if exp < expItem.Exp {
			break
		}

		ret = expItem
	}

	return ret
}

// 获取晶石经验道具 -- 返还经验道具时用
func (g *GlobalConfigEntry) GetCrystalExpItemByExp(exp int32) *ExpItem {
	var ret *ExpItem
	if exp <= 0 {
		return nil
	}

	for _, expItem := range crystalExpItemSection {
		if exp < expItem.Exp {
			break
		}

		ret = expItem
	}

	return ret
}
