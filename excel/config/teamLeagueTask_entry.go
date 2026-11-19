package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var teamLeagueTaskEntries *TeamLeagueTaskEntries //TeamLeagueTask.json全局变量

// 属性表 公会团赛表-拳皇/teamLeagueTask
type TeamLeagueTaskEntry struct {
	Id             int32 `json:"id,omitempty" bson:"id"`                         //##id
	SubConditionId int32 `json:"subConditionId,omitempty" bson:"subConditionId"` //子条件id
	ConditionId    int32 `json:"conditionId,omitempty" bson:"conditionId"`       //条件id
	Score          int32 `json:"score,omitempty" bson:"score"`                   //达成奖励
	Weight         int32 `json:"weight,omitempty" bson:"weight"`                 //权重
	CondMaxCnt     int32 `json:"condMaxCnt,omitempty" bson:"condMaxCnt"`         //完成条件数
}

// TeamLeagueTask.json属性表集合
type TeamLeagueTaskEntries struct {
	Rows   map[int32]*TeamLeagueTaskEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*TeamLeagueTaskEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("TeamLeagueTask.json", (*TeamLeagueTaskEntries)(nil))
}

func (e *TeamLeagueTaskEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	teamLeagueTaskEntries = &TeamLeagueTaskEntries{
		Rows:   make(map[int32]*TeamLeagueTaskEntry),
		Arrays: make([]*TeamLeagueTaskEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &TeamLeagueTaskEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return teamLeagueTaskEntries, err
		}

		teamLeagueTaskEntries.Rows[entry.Id] = entry
		teamLeagueTaskEntries.Arrays = append(teamLeagueTaskEntries.Arrays, entry)
	}

	sort.Slice(teamLeagueTaskEntries.Arrays, func(i, j int) bool {
		return teamLeagueTaskEntries.Arrays[i].Id < teamLeagueTaskEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("teamLeagueTask.json", teamLeagueTaskEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return teamLeagueTaskEntries, nil

}

func (e *TeamLeagueTaskEntries) GetArrays() any {
	return teamLeagueTaskEntries.Arrays
}

func GetTeamLeagueTaskEntry(id int32) (*TeamLeagueTaskEntry, bool) {
	entry, ok := teamLeagueTaskEntries.Rows[id]
	return entry, ok
}

func GetTeamLeagueTaskSize() int {
	return int(len(teamLeagueTaskEntries.Rows))
}

func GetTeamLeagueTaskRows() map[int32]*TeamLeagueTaskEntry {
	return teamLeagueTaskEntries.Rows
}

func GetTeamLeagueTaskArrays() []*TeamLeagueTaskEntry {
	return teamLeagueTaskEntries.Arrays
}
