package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var lDlootEntries *LDlootEntries //LDloot.json全局变量

// 属性表 LD宝箱及掉落-拳皇/LDloot
type LDlootEntry struct {
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

// LDloot.json属性表集合
type LDlootEntries struct {
	Rows   map[int32]*LDlootEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*LDlootEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("LDloot.json", (*LDlootEntries)(nil))
}

func (e *LDlootEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	lDlootEntries = &LDlootEntries{
		Rows:   make(map[int32]*LDlootEntry),
		Arrays: make([]*LDlootEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &LDlootEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return lDlootEntries, err
		}

		lDlootEntries.Rows[entry.Id] = entry
		lDlootEntries.Arrays = append(lDlootEntries.Arrays, entry)
	}

	sort.Slice(lDlootEntries.Arrays, func(i, j int) bool {
		return lDlootEntries.Arrays[i].Id < lDlootEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("lDloot.json", lDlootEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return lDlootEntries, nil

}

func (e *LDlootEntries) GetArrays() any {
	return lDlootEntries.Arrays
}

func GetLDlootEntry(id int32) (*LDlootEntry, bool) {
	entry, ok := lDlootEntries.Rows[id]
	return entry, ok
}

func GetLDlootSize() int {
	return int(len(lDlootEntries.Rows))
}

func GetLDlootRows() map[int32]*LDlootEntry {
	return lDlootEntries.Rows
}

func GetLDlootArrays() []*LDlootEntry {
	return lDlootEntries.Arrays
}
