package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var extraItemLootEntries *ExtraItemLootEntries //ExtraItemLoot.json全局变量

// 属性表 基础配置-拳皇/ExtraItemLoot
type ExtraItemLootEntry struct {
	Id                    int32               `json:"id,omitempty" bson:"id"`                                       //##道具 id
	ExtraAwardMax         int32               `json:"extraAwardMax,omitempty" bson:"extraAwardMax"`                 //额外掉落奖励最大限制(背包拥有的+此次额外获得的)
	DayExtraAwardMax      int32               `json:"dayExtraAwardMax,omitempty" bson:"dayExtraAwardMax"`           //每天额外掉落奖励最大限制
	ExtraAwardProbability []*ValueProbability `json:"extraAwardProbability,omitempty" bson:"extraAwardProbability"` //额外奖励掉落奖励概率(已掉落个数#概率)
}

// ExtraItemLoot.json属性表集合
type ExtraItemLootEntries struct {
	Rows   map[int32]*ExtraItemLootEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ExtraItemLootEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("ExtraItemLoot.json", (*ExtraItemLootEntries)(nil))
}

func (e *ExtraItemLootEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	extraItemLootEntries = &ExtraItemLootEntries{
		Rows:   make(map[int32]*ExtraItemLootEntry),
		Arrays: make([]*ExtraItemLootEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ExtraItemLootEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return extraItemLootEntries, err
		}

		extraItemLootEntries.Rows[entry.Id] = entry
		extraItemLootEntries.Arrays = append(extraItemLootEntries.Arrays, entry)
	}

	sort.Slice(extraItemLootEntries.Arrays, func(i, j int) bool {
		return extraItemLootEntries.Arrays[i].Id < extraItemLootEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("extraItemLoot.json", extraItemLootEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return extraItemLootEntries, nil

}

func (e *ExtraItemLootEntries) GetArrays() any {
	return extraItemLootEntries.Arrays
}

func GetExtraItemLootEntry(id int32) (*ExtraItemLootEntry, bool) {
	entry, ok := extraItemLootEntries.Rows[id]
	return entry, ok
}

func GetExtraItemLootSize() int {
	return int(len(extraItemLootEntries.Rows))
}

func GetExtraItemLootRows() map[int32]*ExtraItemLootEntry {
	return extraItemLootEntries.Rows
}

func GetExtraItemLootArrays() []*ExtraItemLootEntry {
	return extraItemLootEntries.Arrays
}
