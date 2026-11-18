package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroSoulEntries *HeroSoulEntries //HeroSoul.json全局变量

// 属性表 英雄表-拳皇/HeroSoul
type HeroSoulEntry struct {
	Id                     int32                 `json:"id,omitempty" bson:"id"`                                         //##id
	Type                   HeroSoulType          `json:"type,omitempty" bson:"type"`                                     //战魂类型
	Rarity                 Rarity                `json:"rarity,omitempty" bson:"rarity"`                                 //稀有度
	Weight                 int32                 `json:"weight,omitempty" bson:"weight"`                                 //品质内权重
	IsSpecial              bool                  `json:"isSpecial,omitempty" bson:"isSpecial"`                           //是否为特殊本源
	Fix                    bool                  `json:"fix,omitempty" bson:"fix"`                                       //是否固定属性
	FixAttrs               []*AttrValue          `json:"fixAttrs,omitempty" bson:"fixAttrs"`                             //固定属性列表
	GradeList              []HeroSoulGrade       `json:"gradeList,omitempty" bson:"gradeList"`                           //属性评级
	FixAttrTypes           []AttrType            `json:"fixAttrTypes,omitempty" bson:"fixAttrTypes"`                     //固定属性列表
	RandomGradeList        []*RandomId           `json:"randomGradeList,omitempty" bson:"randomGradeList"`               //属性评级随机率
	HideAttrIdList         []int32               `json:"hideAttrIdList,omitempty" bson:"hideAttrIdList"`                 //隐藏属性id列表
	RandomRarityList       []*RandomId           `json:"randomRarityList,omitempty" bson:"randomRarityList"`             //隐藏属性品质随机
	CustomAttrNum          int32                 `json:"customAttrNum,omitempty" bson:"customAttrNum"`                   //定制隐藏属性条目
	CustomHideAttrIdList   []*SoulCustomHideAttr `json:"customHideAttrIdList,omitempty" bson:"customHideAttrIdList"`     //定制隐藏属性ID
	CustomRandomRarityList []*RandomId           `json:"customRandomRarityList,omitempty" bson:"customRandomRarityList"` //定制隐藏属性品质随机
}

// HeroSoul.json属性表集合
type HeroSoulEntries struct {
	Rows   map[int32]*HeroSoulEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroSoulEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroSoul.json", (*HeroSoulEntries)(nil))
}

func (e *HeroSoulEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroSoulEntries = &HeroSoulEntries{
		Rows:   make(map[int32]*HeroSoulEntry),
		Arrays: make([]*HeroSoulEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroSoulEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroSoulEntries, err
		}

		heroSoulEntries.Rows[entry.Id] = entry
		heroSoulEntries.Arrays = append(heroSoulEntries.Arrays, entry)
	}

	sort.Slice(heroSoulEntries.Arrays, func(i, j int) bool {
		return heroSoulEntries.Arrays[i].Id < heroSoulEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroSoul.json", heroSoulEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroSoulEntries, nil

}

func (e *HeroSoulEntries) GetArrays() any {
	return heroSoulEntries.Arrays
}

func GetHeroSoulEntry(id int32) (*HeroSoulEntry, bool) {
	entry, ok := heroSoulEntries.Rows[id]
	return entry, ok
}

func GetHeroSoulSize() int {
	return int(len(heroSoulEntries.Rows))
}

func GetHeroSoulRows() map[int32]*HeroSoulEntry {
	return heroSoulEntries.Rows
}

func GetHeroSoulArrays() []*HeroSoulEntry {
	return heroSoulEntries.Arrays
}
