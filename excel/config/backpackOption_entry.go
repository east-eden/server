package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var backpackOptionEntries *BackpackOptionEntries //BackpackOption.json全局变量

// 属性表 基础配置-拳皇/BackpackOption
type BackpackOptionEntry struct {
	Id         int32   `json:"id,omitempty" bson:"id"`                 //自选道具id
	GoodsList  []int32 `json:"goodsList,omitempty" bson:"goodsList"`   //自选道具包含的道具列表
	NumberList []int32 `json:"numberList,omitempty" bson:"numberList"` //各道具数量
}

// BackpackOption.json属性表集合
type BackpackOptionEntries struct {
	Rows   map[int32]*BackpackOptionEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*BackpackOptionEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("BackpackOption.json", (*BackpackOptionEntries)(nil))
}

func (e *BackpackOptionEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	backpackOptionEntries = &BackpackOptionEntries{
		Rows:   make(map[int32]*BackpackOptionEntry),
		Arrays: make([]*BackpackOptionEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &BackpackOptionEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return backpackOptionEntries, err
		}

		backpackOptionEntries.Rows[entry.Id] = entry
		backpackOptionEntries.Arrays = append(backpackOptionEntries.Arrays, entry)
	}

	sort.Slice(backpackOptionEntries.Arrays, func(i, j int) bool {
		return backpackOptionEntries.Arrays[i].Id < backpackOptionEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("backpackOption.json", backpackOptionEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return backpackOptionEntries, nil

}

func (e *BackpackOptionEntries) GetArrays() any {
	return backpackOptionEntries.Arrays
}

func GetBackpackOptionEntry(id int32) (*BackpackOptionEntry, bool) {
	entry, ok := backpackOptionEntries.Rows[id]
	return entry, ok
}

func GetBackpackOptionSize() int {
	return int(len(backpackOptionEntries.Rows))
}

func GetBackpackOptionRows() map[int32]*BackpackOptionEntry {
	return backpackOptionEntries.Rows
}

func GetBackpackOptionArrays() []*BackpackOptionEntry {
	return backpackOptionEntries.Arrays
}
