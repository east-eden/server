package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var mainQuestEntries *MainQuestEntries //MainQuest.json全局变量

// 属性表 任务表-拳皇/MainQuest
type MainQuestEntry struct {
	Id             int32      `json:"id,omitempty" bson:"id"`                         //索引
	Desc           string     `json:"desc,omitempty" bson:"desc"`                     //描述
	SubConditionId int32      `json:"subConditionId,omitempty" bson:"subConditionId"` //子条件id
	AwardList      []*IdCount `json:"awardList,omitempty" bson:"awardList"`           //奖励列表
	NextId         int32      `json:"nextId,omitempty" bson:"nextId"`                 //下个任务
}

// MainQuest.json属性表集合
type MainQuestEntries struct {
	Rows   map[int32]*MainQuestEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*MainQuestEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("MainQuest.json", (*MainQuestEntries)(nil))
}

func (e *MainQuestEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	mainQuestEntries = &MainQuestEntries{
		Rows:   make(map[int32]*MainQuestEntry),
		Arrays: make([]*MainQuestEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &MainQuestEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return mainQuestEntries, err
		}

		mainQuestEntries.Rows[entry.Id] = entry
		mainQuestEntries.Arrays = append(mainQuestEntries.Arrays, entry)
	}

	sort.Slice(mainQuestEntries.Arrays, func(i, j int) bool {
		return mainQuestEntries.Arrays[i].Id < mainQuestEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("mainQuest.json", mainQuestEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return mainQuestEntries, nil

}

func (e *MainQuestEntries) GetArrays() any {
	return mainQuestEntries.Arrays
}

func GetMainQuestEntry(id int32) (*MainQuestEntry, bool) {
	entry, ok := mainQuestEntries.Rows[id]
	return entry, ok
}

func GetMainQuestSize() int {
	return int(len(mainQuestEntries.Rows))
}

func GetMainQuestRows() map[int32]*MainQuestEntry {
	return mainQuestEntries.Rows
}

func GetMainQuestArrays() []*MainQuestEntry {
	return mainQuestEntries.Arrays
}
