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

var heroLevelNewEntries *HeroLevelNewEntries //HeroLevelNew.json全局变量

// 属性表 英雄表-拳皇/HeroLevelNew
type HeroLevelNewEntry struct {
	Id       Rarity     `json:"id,omitempty" bson:"id"`             //##品质
	Level    int32      `json:"level,omitempty" bson:"level"`       //等级
	NextCost []*IdCount `json:"nextCost,omitempty" bson:"nextCost"` //下一级消耗
}

// HeroLevelNew.json属性表集合
type HeroLevelNewEntries struct {
	Rows   map[string]*HeroLevelNewEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroLevelNewEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroLevelNew.json", (*HeroLevelNewEntries)(nil))
}

func (e *HeroLevelNewEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroLevelNewEntries = &HeroLevelNewEntries{
		Rows:   make(map[string]*HeroLevelNewEntry),
		Arrays: make([]*HeroLevelNewEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroLevelNewEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroLevelNewEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Level))
		heroLevelNewEntries.Rows[key] = entry
		heroLevelNewEntries.Arrays = append(heroLevelNewEntries.Arrays, entry)
	}

	sort.Slice(heroLevelNewEntries.Arrays, func(i, j int) bool {
		return heroLevelNewEntries.Arrays[i].Id < heroLevelNewEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroLevelNew.json", heroLevelNewEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroLevelNewEntries, nil

}

func (e *HeroLevelNewEntries) GetArrays() any {
	return heroLevelNewEntries.Arrays
}

func GetHeroLevelNewEntry(keys ...string) (*HeroLevelNewEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroLevelNewEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroLevelNewSize() int {
	return int(len(heroLevelNewEntries.Rows))
}

func GetHeroLevelNewRows() map[string]*HeroLevelNewEntry {
	return heroLevelNewEntries.Rows
}

func GetHeroLevelNewArrays() []*HeroLevelNewEntry {
	return heroLevelNewEntries.Arrays
}
