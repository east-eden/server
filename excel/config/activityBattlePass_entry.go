package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

var activityBattlePassEntries *ActivityBattlePassEntries //ActivityBattlePass.json全局变量

// 属性表 活动表-拳皇/ActivityBattlePass
type ActivityBattlePassEntry struct {
	Id              string     `json:"id,omitempty" bson:"id"`                           //活动id
	Stage           int32      `json:"stage,omitempty" bson:"stage"`                     //阶段
	Score           int32      `json:"score,omitempty" bson:"score"`                     //所需积分
	AwardItems      []*IdCount `json:"awardItems,omitempty" bson:"awardItems"`           //普通奖励道具列表(免费)
	HightAwardItems []*IdCount `json:"hightAwardItems,omitempty" bson:"hightAwardItems"` //高级奖励道具列表(付费)
}

// ActivityBattlePass.json属性表集合
type ActivityBattlePassEntries struct {
	Rows   map[string]*ActivityBattlePassEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ActivityBattlePassEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("ActivityBattlePass.json", (*ActivityBattlePassEntries)(nil))
}

func (e *ActivityBattlePassEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	activityBattlePassEntries = &ActivityBattlePassEntries{
		Rows:   make(map[string]*ActivityBattlePassEntry),
		Arrays: make([]*ActivityBattlePassEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ActivityBattlePassEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return activityBattlePassEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Stage))
		activityBattlePassEntries.Rows[key] = entry
		activityBattlePassEntries.Arrays = append(activityBattlePassEntries.Arrays, entry)
	}

	sort.Slice(activityBattlePassEntries.Arrays, func(i, j int) bool {
		return activityBattlePassEntries.Arrays[i].Id < activityBattlePassEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("activityBattlePass.json", activityBattlePassEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return activityBattlePassEntries, nil

}

func (e *ActivityBattlePassEntries) GetArrays() any {
	return activityBattlePassEntries.Arrays
}

func GetActivityBattlePassEntry(keys ...string) (*ActivityBattlePassEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := activityBattlePassEntries.Rows[finalKey]
	return entry, ok
}

func GetActivityBattlePassSize() int {
	return int(len(activityBattlePassEntries.Rows))
}

func GetActivityBattlePassRows() map[string]*ActivityBattlePassEntry {
	return activityBattlePassEntries.Rows
}

func GetActivityBattlePassArrays() []*ActivityBattlePassEntry {
	return activityBattlePassEntries.Arrays
}
