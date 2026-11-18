package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var lootEntries *LootEntries //Loot.json全局变量

// 属性表 宝箱及掉落-拳皇/loot
type LootEntry struct {
	Id                    int32     `json:"id,omitempty" bson:"id"`                                       //##id
	Hint                  string    `json:"hint,omitempty" bson:"hint"`                                   //备注
	WeightItemId          []int32   `json:"weightItemId,omitempty" bson:"weightItemId"`                   //权重掉落 物品ID
	WeightItemNum         []int32   `json:"weightItemNum,omitempty" bson:"weightItemNum"`                 //权重掉落 数量
	WeightItemWeight      []int32   `json:"weightItemWeight,omitempty" bson:"weightItemWeight"`           //权重掉落 权重
	WeightLootId          []int32   `json:"weightLootId,omitempty" bson:"weightLootId"`                   //权重掉落 掉落包ID
	WeightLootNum         []int32   `json:"weightLootNum,omitempty" bson:"weightLootNum"`                 //权重掉落 数量
	WeightLootWeight      []int32   `json:"weightLootWeight,omitempty" bson:"weightLootWeight"`           //权重掉落 权重
	RandomItemId          []int32   `json:"randomItemId,omitempty" bson:"randomItemId"`                   //概率掉落 物品id
	RandomItemNum         []int32   `json:"randomItemNum,omitempty" bson:"randomItemNum"`                 //概率掉落 数量
	RandomItemProbability []float32 `json:"randomItemProbability,omitempty" bson:"randomItemProbability"` //概率掉落 概率
	RandomLootId          []int32   `json:"randomLootId,omitempty" bson:"randomLootId"`                   //概率掉落 掉落包
	RandomLootNum         []int32   `json:"randomLootNum,omitempty" bson:"randomLootNum"`                 //概率掉落 数量
	RandomLootProbability []float32 `json:"randomLootProbability,omitempty" bson:"randomLootProbability"` //概率掉落 概率
}

// Loot.json属性表集合
type LootEntries struct {
	Rows   map[int32]*LootEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*LootEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Loot.json", (*LootEntries)(nil))
}

func (e *LootEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	lootEntries = &LootEntries{
		Rows:   make(map[int32]*LootEntry),
		Arrays: make([]*LootEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &LootEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return lootEntries, err
		}

		lootEntries.Rows[entry.Id] = entry
		lootEntries.Arrays = append(lootEntries.Arrays, entry)
	}

	sort.Slice(lootEntries.Arrays, func(i, j int) bool {
		return lootEntries.Arrays[i].Id < lootEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("loot.json", lootEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return lootEntries, nil

}

func (e *LootEntries) GetArrays() any {
	return lootEntries.Arrays
}

func GetLootEntry(id int32) (*LootEntry, bool) {
	entry, ok := lootEntries.Rows[id]
	return entry, ok
}

func GetLootSize() int {
	return int(len(lootEntries.Rows))
}

func GetLootRows() map[int32]*LootEntry {
	return lootEntries.Rows
}

func GetLootArrays() []*LootEntry {
	return lootEntries.Arrays
}
