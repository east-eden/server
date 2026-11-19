package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var arenaDailyRewardEntries *ArenaDailyRewardEntries //ArenaDailyReward.json全局变量

// 属性表 竞技场-拳皇/arenaDailyReward
type ArenaDailyRewardEntry struct {
	Id            int32   `json:"id,omitempty" bson:"id"`                       //索引id
	RewardItemIds []int32 `json:"rewardItemIds,omitempty" bson:"rewardItemIds"` //奖励道具id
	RewardCounts  []int32 `json:"rewardCounts,omitempty" bson:"rewardCounts"`   //奖励数量
	RankRange     int32   `json:"rankRange,omitempty" bson:"rankRange"`         //排名区间起点
}

// ArenaDailyReward.json属性表集合
type ArenaDailyRewardEntries struct {
	Rows   map[int32]*ArenaDailyRewardEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ArenaDailyRewardEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("ArenaDailyReward.json", (*ArenaDailyRewardEntries)(nil))
}

func (e *ArenaDailyRewardEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	arenaDailyRewardEntries = &ArenaDailyRewardEntries{
		Rows:   make(map[int32]*ArenaDailyRewardEntry),
		Arrays: make([]*ArenaDailyRewardEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ArenaDailyRewardEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return arenaDailyRewardEntries, err
		}

		arenaDailyRewardEntries.Rows[entry.Id] = entry
		arenaDailyRewardEntries.Arrays = append(arenaDailyRewardEntries.Arrays, entry)
	}

	sort.Slice(arenaDailyRewardEntries.Arrays, func(i, j int) bool {
		return arenaDailyRewardEntries.Arrays[i].Id < arenaDailyRewardEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("arenaDailyReward.json", arenaDailyRewardEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return arenaDailyRewardEntries, nil

}

func (e *ArenaDailyRewardEntries) GetArrays() any {
	return arenaDailyRewardEntries.Arrays
}

func GetArenaDailyRewardEntry(id int32) (*ArenaDailyRewardEntry, bool) {
	entry, ok := arenaDailyRewardEntries.Rows[id]
	return entry, ok
}

func GetArenaDailyRewardSize() int {
	return int(len(arenaDailyRewardEntries.Rows))
}

func GetArenaDailyRewardRows() map[int32]*ArenaDailyRewardEntry {
	return arenaDailyRewardEntries.Rows
}

func GetArenaDailyRewardArrays() []*ArenaDailyRewardEntry {
	return arenaDailyRewardEntries.Arrays
}
