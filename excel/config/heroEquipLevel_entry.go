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

var heroEquipLevelEntries *HeroEquipLevelEntries //HeroEquipLevel.json全局变量

// 属性表 英雄表-拳皇/heroEquipLevel
type HeroEquipLevelEntry struct {
	Id             int32        `json:"id,omitempty" bson:"id"`                         //属性id
	Level          int32        `json:"level,omitempty" bson:"level"`                   //等级
	NextCost       []*IdCount   `json:"nextCost,omitempty" bson:"nextCost"`             //升到下一级消耗
	TotalAttrValue []*AttrValue `json:"totalAttrValue,omitempty" bson:"totalAttrValue"` //累积属性
}

// HeroEquipLevel.json属性表集合
type HeroEquipLevelEntries struct {
	Rows   map[string]*HeroEquipLevelEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroEquipLevelEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroEquipLevel.json", (*HeroEquipLevelEntries)(nil))
}

func (e *HeroEquipLevelEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroEquipLevelEntries = &HeroEquipLevelEntries{
		Rows:   make(map[string]*HeroEquipLevelEntry),
		Arrays: make([]*HeroEquipLevelEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroEquipLevelEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroEquipLevelEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Level))
		heroEquipLevelEntries.Rows[key] = entry
		heroEquipLevelEntries.Arrays = append(heroEquipLevelEntries.Arrays, entry)
	}

	sort.Slice(heroEquipLevelEntries.Arrays, func(i, j int) bool {
		return heroEquipLevelEntries.Arrays[i].Id < heroEquipLevelEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroEquipLevel.json", heroEquipLevelEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroEquipLevelEntries, nil

}

func (e *HeroEquipLevelEntries) GetArrays() any {
	return heroEquipLevelEntries.Arrays
}

func GetHeroEquipLevelEntry(keys ...string) (*HeroEquipLevelEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroEquipLevelEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroEquipLevelSize() int {
	return int(len(heroEquipLevelEntries.Rows))
}

func GetHeroEquipLevelRows() map[string]*HeroEquipLevelEntry {
	return heroEquipLevelEntries.Rows
}

func GetHeroEquipLevelArrays() []*HeroEquipLevelEntry {
	return heroEquipLevelEntries.Arrays
}
