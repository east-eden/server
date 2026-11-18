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

var heroSoulActiveEntries *HeroSoulActiveEntries //HeroSoulActive.json全局变量

// 属性表 英雄表-拳皇/HeroSoulActive
type HeroSoulActiveEntry struct {
	Id           int32           `json:"id,omitempty" bson:"id"`                     //##id
	Rarity       Rarity          `json:"rarity,omitempty" bson:"rarity"`             //##品质
	Count        int32           `json:"count,omitempty" bson:"count"`               //生效数量
	AttrValues   []*AttrValue    `json:"attrValues,omitempty" bson:"attrValues"`     //激活属性Init
	GodheadScope BonusTargetType `json:"godheadScope,omitempty" bson:"godheadScope"` //作用对象
	IsSwitch     bool            `json:"isSwitch,omitempty" bson:"isSwitch"`         //可切换(TRUE=是)
	Sn           int32           `json:"sn,omitempty" bson:"sn"`                     //词条序号
}

// HeroSoulActive.json属性表集合
type HeroSoulActiveEntries struct {
	Rows   map[string]*HeroSoulActiveEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroSoulActiveEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroSoulActive.json", (*HeroSoulActiveEntries)(nil))
}

func (e *HeroSoulActiveEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroSoulActiveEntries = &HeroSoulActiveEntries{
		Rows:   make(map[string]*HeroSoulActiveEntry),
		Arrays: make([]*HeroSoulActiveEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroSoulActiveEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroSoulActiveEntries, err
		}

		key := fmt.Sprintf("%s_%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Rarity), cast.ToString(entry.Count))
		heroSoulActiveEntries.Rows[key] = entry
		heroSoulActiveEntries.Arrays = append(heroSoulActiveEntries.Arrays, entry)
	}

	sort.Slice(heroSoulActiveEntries.Arrays, func(i, j int) bool {
		return heroSoulActiveEntries.Arrays[i].Id < heroSoulActiveEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroSoulActive.json", heroSoulActiveEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroSoulActiveEntries, nil

}

func (e *HeroSoulActiveEntries) GetArrays() any {
	return heroSoulActiveEntries.Arrays
}

func GetHeroSoulActiveEntry(keys ...string) (*HeroSoulActiveEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroSoulActiveEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroSoulActiveSize() int {
	return int(len(heroSoulActiveEntries.Rows))
}

func GetHeroSoulActiveRows() map[string]*HeroSoulActiveEntry {
	return heroSoulActiveEntries.Rows
}

func GetHeroSoulActiveArrays() []*HeroSoulActiveEntry {
	return heroSoulActiveEntries.Arrays
}
