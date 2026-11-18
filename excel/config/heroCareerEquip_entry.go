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

var heroCareerEquipEntries *HeroCareerEquipEntries //HeroCareerEquip.json全局变量

// 属性表 英雄表-拳皇/heroCareerEquip
type HeroCareerEquipEntry struct {
	Career Career       `json:"career,omitempty" bson:"career"` //职业
	Pos    HeroEquipPos `json:"pos,omitempty" bson:"pos"`       //位置
	AttrId int32        `json:"attrId,omitempty" bson:"attrId"` //英雄装备升级属性Id
}

// HeroCareerEquip.json属性表集合
type HeroCareerEquipEntries struct {
	Rows   map[string]*HeroCareerEquipEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroCareerEquipEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroCareerEquip.json", (*HeroCareerEquipEntries)(nil))
}

func (e *HeroCareerEquipEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroCareerEquipEntries = &HeroCareerEquipEntries{
		Rows:   make(map[string]*HeroCareerEquipEntry),
		Arrays: make([]*HeroCareerEquipEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroCareerEquipEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroCareerEquipEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Career), cast.ToString(entry.Pos))
		heroCareerEquipEntries.Rows[key] = entry
		heroCareerEquipEntries.Arrays = append(heroCareerEquipEntries.Arrays, entry)
	}

	sort.Slice(heroCareerEquipEntries.Arrays, func(i, j int) bool {
		return heroCareerEquipEntries.Arrays[i].Career < heroCareerEquipEntries.Arrays[j].Career
	})

	excel.AddArrayGetter("heroCareerEquip.json", heroCareerEquipEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroCareerEquipEntries, nil

}

func (e *HeroCareerEquipEntries) GetArrays() any {
	return heroCareerEquipEntries.Arrays
}

func GetHeroCareerEquipEntry(keys ...string) (*HeroCareerEquipEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroCareerEquipEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroCareerEquipSize() int {
	return int(len(heroCareerEquipEntries.Rows))
}

func GetHeroCareerEquipRows() map[string]*HeroCareerEquipEntry {
	return heroCareerEquipEntries.Rows
}

func GetHeroCareerEquipArrays() []*HeroCareerEquipEntry {
	return heroCareerEquipEntries.Arrays
}
