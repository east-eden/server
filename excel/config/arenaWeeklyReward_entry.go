package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var arenaWeeklyRewardEntries *ArenaWeeklyRewardEntries //ArenaWeeklyReward.json全局变量

// 属性表 竞技场-拳皇/arenaWeeklyReward
type ArenaWeeklyRewardEntry struct {
	Id                 int32   `json:"id,omitempty" bson:"id"`                                 //索引id
	RewardItemIds      []int32 `json:"rewardItemIds,omitempty" bson:"rewardItemIds"`           //奖励道具id
	RewardItemIdsPlanB []int32 `json:"rewardItemIdsPlanB,omitempty" bson:"rewardItemIdsPlanB"` //奖励道具id planB
	RewardCounts       []int32 `json:"rewardCounts,omitempty" bson:"rewardCounts"`             //奖励数量
	RewardCountsPlanB  []int32 `json:"rewardCountsPlanB,omitempty" bson:"rewardCountsPlanB"`   //奖励数量 planB
	RankRange          int32   `json:"rankRange,omitempty" bson:"rankRange"`                   //排名区间起点
}

// ArenaWeeklyReward.json属性表集合
type ArenaWeeklyRewardEntries struct {
	Rows   map[int32]*ArenaWeeklyRewardEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ArenaWeeklyRewardEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("ArenaWeeklyReward.json", (*ArenaWeeklyRewardEntries)(nil))
}

func (e *ArenaWeeklyRewardEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	arenaWeeklyRewardEntries = &ArenaWeeklyRewardEntries{
		Rows:   make(map[int32]*ArenaWeeklyRewardEntry),
		Arrays: make([]*ArenaWeeklyRewardEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ArenaWeeklyRewardEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return arenaWeeklyRewardEntries, err
		}

		arenaWeeklyRewardEntries.Rows[entry.Id] = entry
		arenaWeeklyRewardEntries.Arrays = append(arenaWeeklyRewardEntries.Arrays, entry)
	}

	sort.Slice(arenaWeeklyRewardEntries.Arrays, func(i, j int) bool {
		return arenaWeeklyRewardEntries.Arrays[i].Id < arenaWeeklyRewardEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("arenaWeeklyReward.json", arenaWeeklyRewardEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return arenaWeeklyRewardEntries, nil

}

func (e *ArenaWeeklyRewardEntries) GetArrays() any {
	return arenaWeeklyRewardEntries.Arrays
}

func GetArenaWeeklyRewardEntry(id int32) (*ArenaWeeklyRewardEntry, bool) {
	entry, ok := arenaWeeklyRewardEntries.Rows[id]
	return entry, ok
}

func GetArenaWeeklyRewardSize() int {
	return int(len(arenaWeeklyRewardEntries.Rows))
}

func GetArenaWeeklyRewardRows() map[int32]*ArenaWeeklyRewardEntry {
	return arenaWeeklyRewardEntries.Rows
}

func GetArenaWeeklyRewardArrays() []*ArenaWeeklyRewardEntry {
	return arenaWeeklyRewardEntries.Arrays
}
