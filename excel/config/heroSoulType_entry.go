package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroSoulTypeEntries *HeroSoulTypeEntries //HeroSoulType.json全局变量

// 属性表 英雄表-拳皇/HeroSoulType
type HeroSoulTypeEntry struct {
	Id                    HeroSoulType    `json:"id,omitempty" bson:"id"`                                       //战魂类型
	RandomAttrList        []*RandomAttr   `json:"randomAttrList,omitempty" bson:"randomAttrList"`               //随机属性列表（属性类型，权重）
	RainbowRandomAttrList []*RandomAttr   `json:"rainbowRandomAttrList,omitempty" bson:"rainbowRandomAttrList"` //永恒随机属性列表
	RedHideIdWeightList   []*RandomId     `json:"redHideIdWeightList,omitempty" bson:"redHideIdWeightList"`     //红色隐藏属性列表
	RedHideIdCountMax     map[int32]int32 `json:"redHideIdCountMax,omitempty" bson:"redHideIdCountMax"`         //红色隐藏属性最大数量
	HideIdWeightList      []*RandomId     `json:"hideIdWeightList,omitempty" bson:"hideIdWeightList"`           //隐藏属性列表
	HideIdCountMax        map[int32]int32 `json:"hideIdCountMax,omitempty" bson:"hideIdCountMax"`               //隐藏属性最大数量
}

// HeroSoulType.json属性表集合
type HeroSoulTypeEntries struct {
	Rows   map[HeroSoulType]*HeroSoulTypeEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroSoulTypeEntry                `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroSoulType.json", (*HeroSoulTypeEntries)(nil))
}

func (e *HeroSoulTypeEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroSoulTypeEntries = &HeroSoulTypeEntries{
		Rows:   make(map[HeroSoulType]*HeroSoulTypeEntry),
		Arrays: make([]*HeroSoulTypeEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroSoulTypeEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroSoulTypeEntries, err
		}

		heroSoulTypeEntries.Rows[entry.Id] = entry
		heroSoulTypeEntries.Arrays = append(heroSoulTypeEntries.Arrays, entry)
	}

	sort.Slice(heroSoulTypeEntries.Arrays, func(i, j int) bool {
		return heroSoulTypeEntries.Arrays[i].Id < heroSoulTypeEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroSoulType.json", heroSoulTypeEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroSoulTypeEntries, nil

}

func (e *HeroSoulTypeEntries) GetArrays() any {
	return heroSoulTypeEntries.Arrays
}

func GetHeroSoulTypeEntry(id HeroSoulType) (*HeroSoulTypeEntry, bool) {
	entry, ok := heroSoulTypeEntries.Rows[id]
	return entry, ok
}

func GetHeroSoulTypeSize() int {
	return int(len(heroSoulTypeEntries.Rows))
}

func GetHeroSoulTypeRows() map[HeroSoulType]*HeroSoulTypeEntry {
	return heroSoulTypeEntries.Rows
}

func GetHeroSoulTypeArrays() []*HeroSoulTypeEntry {
	return heroSoulTypeEntries.Arrays
}
