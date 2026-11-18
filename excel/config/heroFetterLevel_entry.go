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

var heroFetterLevelEntries *HeroFetterLevelEntries //HeroFetterLevel.json全局变量

// 属性表 英雄表-拳皇/heroFetterLevel
type HeroFetterLevelEntry struct {
	Id             int32        `json:"id,omitempty" bson:"id"`                         //羁绊id
	Level          int32        `json:"level,omitempty" bson:"level"`                   //等级
	Star           int32        `json:"star,omitempty" bson:"star"`                     //升到此级所需星级
	TotalAttrValue []*AttrValue `json:"totalAttrValue,omitempty" bson:"totalAttrValue"` //累积属性
	Cost           int32        `json:"cost,omitempty" bson:"cost"`                     //升级消耗羁绊点数
}

// HeroFetterLevel.json属性表集合
type HeroFetterLevelEntries struct {
	Rows   map[string]*HeroFetterLevelEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroFetterLevelEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroFetterLevel.json", (*HeroFetterLevelEntries)(nil))
}

func (e *HeroFetterLevelEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroFetterLevelEntries = &HeroFetterLevelEntries{
		Rows:   make(map[string]*HeroFetterLevelEntry),
		Arrays: make([]*HeroFetterLevelEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroFetterLevelEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroFetterLevelEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Level))
		heroFetterLevelEntries.Rows[key] = entry
		heroFetterLevelEntries.Arrays = append(heroFetterLevelEntries.Arrays, entry)
	}

	sort.Slice(heroFetterLevelEntries.Arrays, func(i, j int) bool {
		return heroFetterLevelEntries.Arrays[i].Id < heroFetterLevelEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroFetterLevel.json", heroFetterLevelEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroFetterLevelEntries, nil

}

func (e *HeroFetterLevelEntries) GetArrays() any {
	return heroFetterLevelEntries.Arrays
}

func GetHeroFetterLevelEntry(keys ...string) (*HeroFetterLevelEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroFetterLevelEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroFetterLevelSize() int {
	return int(len(heroFetterLevelEntries.Rows))
}

func GetHeroFetterLevelRows() map[string]*HeroFetterLevelEntry {
	return heroFetterLevelEntries.Rows
}

func GetHeroFetterLevelArrays() []*HeroFetterLevelEntry {
	return heroFetterLevelEntries.Arrays
}
