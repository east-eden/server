package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var conditionEntries *ConditionEntries //Condition.json全局变量

// 属性表 任务表-拳皇/Condition
type ConditionEntry struct {
	Id              int32                `json:"id,omitempty" bson:"id"`                           //条件id
	Desc            string               `json:"desc,omitempty" bson:"desc"`                       //描述
	SubCondRelation SubConditionRelation `json:"subCondRelation,omitempty" bson:"subCondRelation"` //子条件关系
	SubCondIds      []int32              `json:"subCondIds,omitempty" bson:"subCondIds"`           //子条件id列表
}

// Condition.json属性表集合
type ConditionEntries struct {
	Rows   map[int32]*ConditionEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ConditionEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Condition.json", (*ConditionEntries)(nil))
}

func (e *ConditionEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	conditionEntries = &ConditionEntries{
		Rows:   make(map[int32]*ConditionEntry),
		Arrays: make([]*ConditionEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ConditionEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return conditionEntries, err
		}

		conditionEntries.Rows[entry.Id] = entry
		conditionEntries.Arrays = append(conditionEntries.Arrays, entry)
	}

	sort.Slice(conditionEntries.Arrays, func(i, j int) bool {
		return conditionEntries.Arrays[i].Id < conditionEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("condition.json", conditionEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return conditionEntries, nil

}

func (e *ConditionEntries) GetArrays() any {
	return conditionEntries.Arrays
}

func GetConditionEntry(id int32) (*ConditionEntry, bool) {
	entry, ok := conditionEntries.Rows[id]
	return entry, ok
}

func GetConditionSize() int {
	return int(len(conditionEntries.Rows))
}

func GetConditionRows() map[int32]*ConditionEntry {
	return conditionEntries.Rows
}

func GetConditionArrays() []*ConditionEntry {
	return conditionEntries.Arrays
}
