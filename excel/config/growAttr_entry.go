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

var growAttrEntries *GrowAttrEntries //GrowAttr.json全局变量

// 属性表 基础配置-拳皇/GrowAttr
type GrowAttrEntry struct {
	Id           int32            `json:"id,omitempty" bson:"id"`                     //##属性id
	Level        int32            `json:"level,omitempty" bson:"level"`               //等级左侧区间
	GodheadScope BonusTargetType  `json:"godheadScope,omitempty" bson:"godheadScope"` //作用对象
	HeroFetterID int32            `json:"heroFetterID,omitempty" bson:"heroFetterID"` //羁绊队伍ID
	Attrs        []*AttrValueCalc `json:"attrs,omitempty" bson:"attrs"`               //对应等级属性
}

// GrowAttr.json属性表集合
type GrowAttrEntries struct {
	Rows   map[string]*GrowAttrEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*GrowAttrEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("GrowAttr.json", (*GrowAttrEntries)(nil))
}

func (e *GrowAttrEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	growAttrEntries = &GrowAttrEntries{
		Rows:   make(map[string]*GrowAttrEntry),
		Arrays: make([]*GrowAttrEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &GrowAttrEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return growAttrEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Level))
		growAttrEntries.Rows[key] = entry
		growAttrEntries.Arrays = append(growAttrEntries.Arrays, entry)
	}

	sort.Slice(growAttrEntries.Arrays, func(i, j int) bool {
		return growAttrEntries.Arrays[i].Id < growAttrEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("growAttr.json", growAttrEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return growAttrEntries, nil

}

func (e *GrowAttrEntries) GetArrays() any {
	return growAttrEntries.Arrays
}

func GetGrowAttrEntry(keys ...string) (*GrowAttrEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := growAttrEntries.Rows[finalKey]
	return entry, ok
}

func GetGrowAttrSize() int {
	return int(len(growAttrEntries.Rows))
}

func GetGrowAttrRows() map[string]*GrowAttrEntry {
	return growAttrEntries.Rows
}

func GetGrowAttrArrays() []*GrowAttrEntry {
	return growAttrEntries.Arrays
}
