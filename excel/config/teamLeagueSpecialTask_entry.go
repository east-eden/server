package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var teamLeagueSpecialTaskEntries *TeamLeagueSpecialTaskEntries //TeamLeagueSpecialTask.json全局变量

// 属性表 公会团赛表-拳皇/teamLeagueSpecialTask
type TeamLeagueSpecialTaskEntry struct {
	Id            int32 `json:"id,omitempty" bson:"id"`                       //##id
	RechargeMoney int32 `json:"rechargeMoney,omitempty" bson:"rechargeMoney"` //完成所需要的充值金额数(仅当期有效)
	Score         int32 `json:"score,omitempty" bson:"score"`                 //积分
	Weight        int32 `json:"weight,omitempty" bson:"weight"`               //权重
	Cnt           int32 `json:"cnt,omitempty" bson:"cnt"`                     //可以完成的次数
}

// TeamLeagueSpecialTask.json属性表集合
type TeamLeagueSpecialTaskEntries struct {
	Rows   map[int32]*TeamLeagueSpecialTaskEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*TeamLeagueSpecialTaskEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("TeamLeagueSpecialTask.json", (*TeamLeagueSpecialTaskEntries)(nil))
}

func (e *TeamLeagueSpecialTaskEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	teamLeagueSpecialTaskEntries = &TeamLeagueSpecialTaskEntries{
		Rows:   make(map[int32]*TeamLeagueSpecialTaskEntry),
		Arrays: make([]*TeamLeagueSpecialTaskEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &TeamLeagueSpecialTaskEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return teamLeagueSpecialTaskEntries, err
		}

		teamLeagueSpecialTaskEntries.Rows[entry.Id] = entry
		teamLeagueSpecialTaskEntries.Arrays = append(teamLeagueSpecialTaskEntries.Arrays, entry)
	}

	sort.Slice(teamLeagueSpecialTaskEntries.Arrays, func(i, j int) bool {
		return teamLeagueSpecialTaskEntries.Arrays[i].Id < teamLeagueSpecialTaskEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("teamLeagueSpecialTask.json", teamLeagueSpecialTaskEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return teamLeagueSpecialTaskEntries, nil

}

func (e *TeamLeagueSpecialTaskEntries) GetArrays() any {
	return teamLeagueSpecialTaskEntries.Arrays
}

func GetTeamLeagueSpecialTaskEntry(id int32) (*TeamLeagueSpecialTaskEntry, bool) {
	entry, ok := teamLeagueSpecialTaskEntries.Rows[id]
	return entry, ok
}

func GetTeamLeagueSpecialTaskSize() int {
	return int(len(teamLeagueSpecialTaskEntries.Rows))
}

func GetTeamLeagueSpecialTaskRows() map[int32]*TeamLeagueSpecialTaskEntry {
	return teamLeagueSpecialTaskEntries.Rows
}

func GetTeamLeagueSpecialTaskArrays() []*TeamLeagueSpecialTaskEntry {
	return teamLeagueSpecialTaskEntries.Arrays
}
