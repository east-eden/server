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

var heroStarEntries *HeroStarEntries //HeroStar.json全局变量

// 属性表 英雄表-拳皇/heroStar
type HeroStarEntry struct {
	Id             int32        `json:"id,omitempty" bson:"id"`                         //星级
	Rarity         Rarity       `json:"rarity,omitempty" bson:"rarity"`                 //品质
	Cost           int32        `json:"cost,omitempty" bson:"cost"`                     //升到下一星级消耗碎片数
	Attrs          []*AttrValue `json:"attrs,omitempty" bson:"attrs"`                   //升到此星级加成属性
	ExtraAttrs     []*AttrValue `json:"extraAttrs,omitempty" bson:"extraAttrs"`         //升到此星级额外加成属性
	AttackAttrs    []*AttrValue `json:"attackAttrs,omitempty" bson:"attackAttrs"`       //攻击职业升到此星级额外属性
	DefenseAttrs   []*AttrValue `json:"defenseAttrs,omitempty" bson:"defenseAttrs"`     //防御职业升到此星级额外属性
	SkillAttrs     []*AttrValue `json:"skillAttrs,omitempty" bson:"skillAttrs"`         //均衡职业升到此星级额外属性
	AuxiliaryAttrs []*AttrValue `json:"auxiliaryAttrs,omitempty" bson:"auxiliaryAttrs"` //辅助职业升到此星级额外属性
}

// HeroStar.json属性表集合
type HeroStarEntries struct {
	Rows   map[string]*HeroStarEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroStarEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroStar.json", (*HeroStarEntries)(nil))
}

func (e *HeroStarEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroStarEntries = &HeroStarEntries{
		Rows:   make(map[string]*HeroStarEntry),
		Arrays: make([]*HeroStarEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroStarEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroStarEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Rarity))
		heroStarEntries.Rows[key] = entry
		heroStarEntries.Arrays = append(heroStarEntries.Arrays, entry)
	}

	sort.Slice(heroStarEntries.Arrays, func(i, j int) bool {
		return heroStarEntries.Arrays[i].Id < heroStarEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroStar.json", heroStarEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroStarEntries, nil

}

func (e *HeroStarEntries) GetArrays() any {
	return heroStarEntries.Arrays
}

func GetHeroStarEntry(keys ...string) (*HeroStarEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroStarEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroStarSize() int {
	return int(len(heroStarEntries.Rows))
}

func GetHeroStarRows() map[string]*HeroStarEntry {
	return heroStarEntries.Rows
}

func GetHeroStarArrays() []*HeroStarEntry {
	return heroStarEntries.Arrays
}
