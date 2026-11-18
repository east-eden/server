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

var heroSoulAttrEntries *HeroSoulAttrEntries //HeroSoulAttr.json全局变量

// 属性表 英雄表-拳皇/HeroSoulAttr
type HeroSoulAttrEntry struct {
	Id            AttrType      `json:"id,omitempty" bson:"id"`                       //属性
	Name          HeroSoulGrade `json:"name,omitempty" bson:"name"`                   //属性评级
	Min           float32       `json:"min,omitempty" bson:"min"`                     //Min
	Max           float32       `json:"max,omitempty" bson:"max"`                     //Max
	BaseMark      float32       `json:"baseMark,omitempty" bson:"baseMark"`           //词条基础分
	RandomMarkMax float32       `json:"randomMarkMax,omitempty" bson:"randomMarkMax"` //词条随机最大分值
}

// HeroSoulAttr.json属性表集合
type HeroSoulAttrEntries struct {
	Rows   map[string]*HeroSoulAttrEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroSoulAttrEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroSoulAttr.json", (*HeroSoulAttrEntries)(nil))
}

func (e *HeroSoulAttrEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroSoulAttrEntries = &HeroSoulAttrEntries{
		Rows:   make(map[string]*HeroSoulAttrEntry),
		Arrays: make([]*HeroSoulAttrEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroSoulAttrEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroSoulAttrEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Name))
		heroSoulAttrEntries.Rows[key] = entry
		heroSoulAttrEntries.Arrays = append(heroSoulAttrEntries.Arrays, entry)
	}

	sort.Slice(heroSoulAttrEntries.Arrays, func(i, j int) bool {
		return heroSoulAttrEntries.Arrays[i].Id < heroSoulAttrEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroSoulAttr.json", heroSoulAttrEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroSoulAttrEntries, nil

}

func (e *HeroSoulAttrEntries) GetArrays() any {
	return heroSoulAttrEntries.Arrays
}

func GetHeroSoulAttrEntry(keys ...string) (*HeroSoulAttrEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroSoulAttrEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroSoulAttrSize() int {
	return int(len(heroSoulAttrEntries.Rows))
}

func GetHeroSoulAttrRows() map[string]*HeroSoulAttrEntry {
	return heroSoulAttrEntries.Rows
}

func GetHeroSoulAttrArrays() []*HeroSoulAttrEntry {
	return heroSoulAttrEntries.Arrays
}
