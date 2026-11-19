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

var gameplayCrossDivideEntries *GameplayCrossDivideEntries //GameplayCrossDivide.json全局变量

// 属性表 跨服配置-拳皇/gameplayCrossDivide
type GameplayCrossDivideEntry struct {
	Id         GameplayType            `json:"id,omitempty" bson:"id"`                 //玩法类型
	Round      int32                   `json:"round,omitempty" bson:"round"`           //期数
	DivideType GameplayCrossDivideType `json:"divideType,omitempty" bson:"divideType"` //分组类型
}

// GameplayCrossDivide.json属性表集合
type GameplayCrossDivideEntries struct {
	Rows   map[string]*GameplayCrossDivideEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*GameplayCrossDivideEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("GameplayCrossDivide.json", (*GameplayCrossDivideEntries)(nil))
}

func (e *GameplayCrossDivideEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	gameplayCrossDivideEntries = &GameplayCrossDivideEntries{
		Rows:   make(map[string]*GameplayCrossDivideEntry),
		Arrays: make([]*GameplayCrossDivideEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &GameplayCrossDivideEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return gameplayCrossDivideEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Round))
		gameplayCrossDivideEntries.Rows[key] = entry
		gameplayCrossDivideEntries.Arrays = append(gameplayCrossDivideEntries.Arrays, entry)
	}

	sort.Slice(gameplayCrossDivideEntries.Arrays, func(i, j int) bool {
		return gameplayCrossDivideEntries.Arrays[i].Id < gameplayCrossDivideEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("gameplayCrossDivide.json", gameplayCrossDivideEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return gameplayCrossDivideEntries, nil

}

func (e *GameplayCrossDivideEntries) GetArrays() any {
	return gameplayCrossDivideEntries.Arrays
}

func GetGameplayCrossDivideEntry(keys ...string) (*GameplayCrossDivideEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := gameplayCrossDivideEntries.Rows[finalKey]
	return entry, ok
}

func GetGameplayCrossDivideSize() int {
	return int(len(gameplayCrossDivideEntries.Rows))
}

func GetGameplayCrossDivideRows() map[string]*GameplayCrossDivideEntry {
	return gameplayCrossDivideEntries.Rows
}

func GetGameplayCrossDivideArrays() []*GameplayCrossDivideEntry {
	return gameplayCrossDivideEntries.Arrays
}
