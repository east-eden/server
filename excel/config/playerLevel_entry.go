package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var playerLevelEntries *PlayerLevelEntries //PlayerLevel.json全局变量

// 属性表 玩家表-拳皇/playerLevel
type PlayerLevelEntry struct {
	Id           int32   `json:"id,omitempty" bson:"id"`                     //##玩家等级
	Exp          int32   `json:"exp,omitempty" bson:"exp"`                   //升下一级经验
	Power        int32   `json:"power,omitempty" bson:"power"`               //玩家等级对应战力
	FirstMoveVal float32 `json:"firstMoveVal,omitempty" bson:"firstMoveVal"` //队伍先手值
}

// PlayerLevel.json属性表集合
type PlayerLevelEntries struct {
	Rows   map[int32]*PlayerLevelEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*PlayerLevelEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("PlayerLevel.json", (*PlayerLevelEntries)(nil))
}

func (e *PlayerLevelEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	playerLevelEntries = &PlayerLevelEntries{
		Rows:   make(map[int32]*PlayerLevelEntry),
		Arrays: make([]*PlayerLevelEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &PlayerLevelEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return playerLevelEntries, err
		}

		playerLevelEntries.Rows[entry.Id] = entry
		playerLevelEntries.Arrays = append(playerLevelEntries.Arrays, entry)
	}

	sort.Slice(playerLevelEntries.Arrays, func(i, j int) bool {
		return playerLevelEntries.Arrays[i].Id < playerLevelEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("playerLevel.json", playerLevelEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return playerLevelEntries, nil

}

func (e *PlayerLevelEntries) GetArrays() any {
	return playerLevelEntries.Arrays
}

func GetPlayerLevelEntry(id int32) (*PlayerLevelEntry, bool) {
	entry, ok := playerLevelEntries.Rows[id]
	return entry, ok
}

func GetPlayerLevelSize() int {
	return int(len(playerLevelEntries.Rows))
}

func GetPlayerLevelRows() map[int32]*PlayerLevelEntry {
	return playerLevelEntries.Rows
}

func GetPlayerLevelArrays() []*PlayerLevelEntry {
	return playerLevelEntries.Arrays
}
