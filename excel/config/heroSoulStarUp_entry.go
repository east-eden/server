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

var heroSoulStarUpEntries *HeroSoulStarUpEntries //HeroSoulStarUp.json全局变量

// 属性表 英雄表-拳皇/HeroSoulStarUp
type HeroSoulStarUpEntry struct {
	Id            HeroSoulType        `json:"id,omitempty" bson:"id"`                       //战魂类型
	Star          int32               `json:"star,omitempty" bson:"star"`                   //星级
	Material1     *SoulStarUpMaterial `json:"material1,omitempty" bson:"material1"`         //材料1
	Material2     *SoulStarUpMaterial `json:"material2,omitempty" bson:"material2"`         //材料2
	Material3     *SoulStarUpMaterial `json:"material3,omitempty" bson:"material3"`         //材料3
	RequiredItems []*IdCount          `json:"requiredItems,omitempty" bson:"requiredItems"` //额外道具需求
}

// HeroSoulStarUp.json属性表集合
type HeroSoulStarUpEntries struct {
	Rows   map[string]*HeroSoulStarUpEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroSoulStarUpEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroSoulStarUp.json", (*HeroSoulStarUpEntries)(nil))
}

func (e *HeroSoulStarUpEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroSoulStarUpEntries = &HeroSoulStarUpEntries{
		Rows:   make(map[string]*HeroSoulStarUpEntry),
		Arrays: make([]*HeroSoulStarUpEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroSoulStarUpEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroSoulStarUpEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Star))
		heroSoulStarUpEntries.Rows[key] = entry
		heroSoulStarUpEntries.Arrays = append(heroSoulStarUpEntries.Arrays, entry)
	}

	sort.Slice(heroSoulStarUpEntries.Arrays, func(i, j int) bool {
		return heroSoulStarUpEntries.Arrays[i].Id < heroSoulStarUpEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroSoulStarUp.json", heroSoulStarUpEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroSoulStarUpEntries, nil

}

func (e *HeroSoulStarUpEntries) GetArrays() any {
	return heroSoulStarUpEntries.Arrays
}

func GetHeroSoulStarUpEntry(keys ...string) (*HeroSoulStarUpEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroSoulStarUpEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroSoulStarUpSize() int {
	return int(len(heroSoulStarUpEntries.Rows))
}

func GetHeroSoulStarUpRows() map[string]*HeroSoulStarUpEntry {
	return heroSoulStarUpEntries.Rows
}

func GetHeroSoulStarUpArrays() []*HeroSoulStarUpEntry {
	return heroSoulStarUpEntries.Arrays
}
