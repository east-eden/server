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

var activitySignInEntries *ActivitySignInEntries //ActivitySignIn.json全局变量

// 属性表 活动表-拳皇/ActivitySignIn
type ActivitySignInEntry struct {
	ActivityId string     `json:"activityId,omitempty" bson:"activityId"` //##活动id
	Day        int32      `json:"day,omitempty" bson:"day"`               //##签到天数
	Awards     []*IdCount `json:"awards,omitempty" bson:"awards"`         //奖励列表
}

// ActivitySignIn.json属性表集合
type ActivitySignInEntries struct {
	Rows   map[string]*ActivitySignInEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ActivitySignInEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("ActivitySignIn.json", (*ActivitySignInEntries)(nil))
}

func (e *ActivitySignInEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	activitySignInEntries = &ActivitySignInEntries{
		Rows:   make(map[string]*ActivitySignInEntry),
		Arrays: make([]*ActivitySignInEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ActivitySignInEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return activitySignInEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.ActivityId), cast.ToString(entry.Day))
		activitySignInEntries.Rows[key] = entry
		activitySignInEntries.Arrays = append(activitySignInEntries.Arrays, entry)
	}

	sort.Slice(activitySignInEntries.Arrays, func(i, j int) bool {
		return activitySignInEntries.Arrays[i].ActivityId < activitySignInEntries.Arrays[j].ActivityId
	})

	excel.AddArrayGetter("activitySignIn.json", activitySignInEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return activitySignInEntries, nil

}

func (e *ActivitySignInEntries) GetArrays() any {
	return activitySignInEntries.Arrays
}

func GetActivitySignInEntry(keys ...string) (*ActivitySignInEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := activitySignInEntries.Rows[finalKey]
	return entry, ok
}

func GetActivitySignInSize() int {
	return int(len(activitySignInEntries.Rows))
}

func GetActivitySignInRows() map[string]*ActivitySignInEntry {
	return activitySignInEntries.Rows
}

func GetActivitySignInArrays() []*ActivitySignInEntry {
	return activitySignInEntries.Arrays
}
