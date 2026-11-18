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

var heroCareerSoulEntries *HeroCareerSoulEntries //HeroCareerSoul.json全局变量

// 属性表 英雄表-拳皇/HeroCareerSoul
type HeroCareerSoulEntry struct {
	Id      Career       `json:"id,omitempty" bson:"id"`           //##职业
	Pos     int32        `json:"pos,omitempty" bson:"pos"`         //格子索引
	Type    HeroSoulType `json:"type,omitempty" bson:"type"`       //战魂类型待删除0530,下版本改用新字段
	NewType HeroSoulType `json:"newType,omitempty" bson:"newType"` //新战魂类型
}

// HeroCareerSoul.json属性表集合
type HeroCareerSoulEntries struct {
	Rows   map[string]*HeroCareerSoulEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroCareerSoulEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroCareerSoul.json", (*HeroCareerSoulEntries)(nil))
}

func (e *HeroCareerSoulEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroCareerSoulEntries = &HeroCareerSoulEntries{
		Rows:   make(map[string]*HeroCareerSoulEntry),
		Arrays: make([]*HeroCareerSoulEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroCareerSoulEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroCareerSoulEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Pos))
		heroCareerSoulEntries.Rows[key] = entry
		heroCareerSoulEntries.Arrays = append(heroCareerSoulEntries.Arrays, entry)
	}

	sort.Slice(heroCareerSoulEntries.Arrays, func(i, j int) bool {
		return heroCareerSoulEntries.Arrays[i].Id < heroCareerSoulEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroCareerSoul.json", heroCareerSoulEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroCareerSoulEntries, nil

}

func (e *HeroCareerSoulEntries) GetArrays() any {
	return heroCareerSoulEntries.Arrays
}

func GetHeroCareerSoulEntry(keys ...string) (*HeroCareerSoulEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroCareerSoulEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroCareerSoulSize() int {
	return int(len(heroCareerSoulEntries.Rows))
}

func GetHeroCareerSoulRows() map[string]*HeroCareerSoulEntry {
	return heroCareerSoulEntries.Rows
}

func GetHeroCareerSoulArrays() []*HeroCareerSoulEntry {
	return heroCareerSoulEntries.Arrays
}
