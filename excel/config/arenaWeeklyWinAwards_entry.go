package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var arenaWeeklyWinAwardsEntries *ArenaWeeklyWinAwardsEntries //ArenaWeeklyWinAwards.json全局变量

// 属性表 竞技场-拳皇/arenaWeeklyWinAwards
type ArenaWeeklyWinAwardsEntry struct {
	Id     int32      `json:"id,omitempty" bson:"id"`         //进度id
	Score  int32      `json:"score,omitempty" bson:"score"`   //胜利次数
	Awards []*IdCount `json:"awards,omitempty" bson:"awards"` //奖励
}

// ArenaWeeklyWinAwards.json属性表集合
type ArenaWeeklyWinAwardsEntries struct {
	Rows   map[int32]*ArenaWeeklyWinAwardsEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ArenaWeeklyWinAwardsEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("ArenaWeeklyWinAwards.json", (*ArenaWeeklyWinAwardsEntries)(nil))
}

func (e *ArenaWeeklyWinAwardsEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	arenaWeeklyWinAwardsEntries = &ArenaWeeklyWinAwardsEntries{
		Rows:   make(map[int32]*ArenaWeeklyWinAwardsEntry),
		Arrays: make([]*ArenaWeeklyWinAwardsEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ArenaWeeklyWinAwardsEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return arenaWeeklyWinAwardsEntries, err
		}

		arenaWeeklyWinAwardsEntries.Rows[entry.Id] = entry
		arenaWeeklyWinAwardsEntries.Arrays = append(arenaWeeklyWinAwardsEntries.Arrays, entry)
	}

	sort.Slice(arenaWeeklyWinAwardsEntries.Arrays, func(i, j int) bool {
		return arenaWeeklyWinAwardsEntries.Arrays[i].Id < arenaWeeklyWinAwardsEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("arenaWeeklyWinAwards.json", arenaWeeklyWinAwardsEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return arenaWeeklyWinAwardsEntries, nil

}

func (e *ArenaWeeklyWinAwardsEntries) GetArrays() any {
	return arenaWeeklyWinAwardsEntries.Arrays
}

func GetArenaWeeklyWinAwardsEntry(id int32) (*ArenaWeeklyWinAwardsEntry, bool) {
	entry, ok := arenaWeeklyWinAwardsEntries.Rows[id]
	return entry, ok
}

func GetArenaWeeklyWinAwardsSize() int {
	return int(len(arenaWeeklyWinAwardsEntries.Rows))
}

func GetArenaWeeklyWinAwardsRows() map[int32]*ArenaWeeklyWinAwardsEntry {
	return arenaWeeklyWinAwardsEntries.Rows
}

func GetArenaWeeklyWinAwardsArrays() []*ArenaWeeklyWinAwardsEntry {
	return arenaWeeklyWinAwardsEntries.Arrays
}
