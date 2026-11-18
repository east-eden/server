package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var subConditionEntries *SubConditionEntries //SubCondition.json全局变量

// 属性表 任务表-拳皇/SubCondition
type SubConditionEntry struct {
	Id            int32                       `json:"id,omitempty" bson:"id"`                       //子条件id
	SubCondType   SubConditionType            `json:"subCondType,omitempty" bson:"subCondType"`     //条件类型
	SubCondParams []int32                     `json:"subCondParams,omitempty" bson:"subCondParams"` //条件参数
	OperationType []SubConditionOperationType `json:"operationType,omitempty" bson:"operationType"` //子条件参数比较方式
	CounterType   CounterType                 `json:"counterType,omitempty" bson:"counterType"`     //计数方式
}

// SubCondition.json属性表集合
type SubConditionEntries struct {
	Rows   map[int32]*SubConditionEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*SubConditionEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("SubCondition.json", (*SubConditionEntries)(nil))
}

func (e *SubConditionEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	subConditionEntries = &SubConditionEntries{
		Rows:   make(map[int32]*SubConditionEntry),
		Arrays: make([]*SubConditionEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &SubConditionEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return subConditionEntries, err
		}

		subConditionEntries.Rows[entry.Id] = entry
		subConditionEntries.Arrays = append(subConditionEntries.Arrays, entry)
	}

	sort.Slice(subConditionEntries.Arrays, func(i, j int) bool {
		return subConditionEntries.Arrays[i].Id < subConditionEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("subCondition.json", subConditionEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return subConditionEntries, nil

}

func (e *SubConditionEntries) GetArrays() any {
	return subConditionEntries.Arrays
}

func GetSubConditionEntry(id int32) (*SubConditionEntry, bool) {
	entry, ok := subConditionEntries.Rows[id]
	return entry, ok
}

func GetSubConditionSize() int {
	return int(len(subConditionEntries.Rows))
}

func GetSubConditionRows() map[int32]*SubConditionEntry {
	return subConditionEntries.Rows
}

func GetSubConditionArrays() []*SubConditionEntry {
	return subConditionEntries.Arrays
}
