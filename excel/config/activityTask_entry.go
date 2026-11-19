package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var activityTaskEntries *ActivityTaskEntries //ActivityTask.json全局变量

// 属性表 活动表-拳皇/ActivityTask
type ActivityTaskEntry struct {
	Id             int32      `json:"id,omitempty" bson:"id"`                         //##id
	ActivityId     string     `json:"activityId,omitempty" bson:"activityId"`         //所属活动id
	SubConditionId int32      `json:"subConditionId,omitempty" bson:"subConditionId"` //子条件id
	Order          int32      `json:"order,omitempty" bson:"order"`                   //排序
	ConditionId    int32      `json:"conditionId,omitempty" bson:"conditionId"`       //条件id
	Awards         []*IdCount `json:"awards,omitempty" bson:"awards"`                 //达成奖励
	ScoreAwards    []*IdCount `json:"scoreAwards,omitempty" bson:"scoreAwards"`       //积分奖励
}

// ActivityTask.json属性表集合
type ActivityTaskEntries struct {
	Rows   map[int32]*ActivityTaskEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ActivityTaskEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("ActivityTask.json", (*ActivityTaskEntries)(nil))
}

func (e *ActivityTaskEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	activityTaskEntries = &ActivityTaskEntries{
		Rows:   make(map[int32]*ActivityTaskEntry),
		Arrays: make([]*ActivityTaskEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ActivityTaskEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return activityTaskEntries, err
		}

		activityTaskEntries.Rows[entry.Id] = entry
		activityTaskEntries.Arrays = append(activityTaskEntries.Arrays, entry)
	}

	sort.Slice(activityTaskEntries.Arrays, func(i, j int) bool {
		return activityTaskEntries.Arrays[i].Id < activityTaskEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("activityTask.json", activityTaskEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return activityTaskEntries, nil

}

func (e *ActivityTaskEntries) GetArrays() any {
	return activityTaskEntries.Arrays
}

func GetActivityTaskEntry(id int32) (*ActivityTaskEntry, bool) {
	entry, ok := activityTaskEntries.Rows[id]
	return entry, ok
}

func GetActivityTaskSize() int {
	return int(len(activityTaskEntries.Rows))
}

func GetActivityTaskRows() map[int32]*ActivityTaskEntry {
	return activityTaskEntries.Rows
}

func GetActivityTaskArrays() []*ActivityTaskEntry {
	return activityTaskEntries.Arrays
}
