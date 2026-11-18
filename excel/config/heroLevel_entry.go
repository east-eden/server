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

var heroLevelEntries *HeroLevelEntries //HeroLevel.json全局变量

// 属性表 英雄表-拳皇/HeroLevel
type HeroLevelEntry struct {
	Id       Rarity     `json:"id,omitempty" bson:"id"`             //##品质
	Level    int32      `json:"level,omitempty" bson:"level"`       //等级
	NextCost []*IdCount `json:"nextCost,omitempty" bson:"nextCost"` //下一级消耗
}

// HeroLevel.json属性表集合
type HeroLevelEntries struct {
	Rows   map[string]*HeroLevelEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroLevelEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroLevel.json", (*HeroLevelEntries)(nil))
}

func (e *HeroLevelEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroLevelEntries = &HeroLevelEntries{
		Rows:   make(map[string]*HeroLevelEntry),
		Arrays: make([]*HeroLevelEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroLevelEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroLevelEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Level))
		heroLevelEntries.Rows[key] = entry
		heroLevelEntries.Arrays = append(heroLevelEntries.Arrays, entry)
	}

	sort.Slice(heroLevelEntries.Arrays, func(i, j int) bool {
		return heroLevelEntries.Arrays[i].Id < heroLevelEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroLevel.json", heroLevelEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroLevelEntries, nil

}

func (e *HeroLevelEntries) GetArrays() any {
	return heroLevelEntries.Arrays
}

func GetHeroLevelEntry(keys ...string) (*HeroLevelEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroLevelEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroLevelSize() int {
	return int(len(heroLevelEntries.Rows))
}

func GetHeroLevelRows() map[string]*HeroLevelEntry {
	return heroLevelEntries.Rows
}

func GetHeroLevelArrays() []*HeroLevelEntry {
	return heroLevelEntries.Arrays
}
