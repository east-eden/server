package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroSoulRarityEntries *HeroSoulRarityEntries //HeroSoulRarity.json全局变量

// 属性表 英雄表-拳皇/HeroSoulRarity
type HeroSoulRarityEntry struct {
	Id                     Rarity      `json:"id,omitempty" bson:"id"`                                         //战魂品质
	GoodsId                int32       `json:"goodsId,omitempty" bson:"goodsId"`                               //通用道具Id
	SpecialProb            float32     `json:"specialProb,omitempty" bson:"specialProb"`                       //出现特殊本源概率
	SpecialLoot            int32       `json:"specialLoot,omitempty" bson:"specialLoot"`                       //特殊本源掉落id
	SpecialMustNoBind      bool        `json:"specialMustNoBind,omitempty" bson:"specialMustNoBind"`           //特殊一定不绑定
	AttrNumWeight          []*RandomId `json:"attrNumWeight,omitempty" bson:"attrNumWeight"`                   //属性条目权重
	AttrNum                int32       `json:"attrNum,omitempty" bson:"attrNum"`                               //属性条数
	RandomAttrGradeList    []*RandomId `json:"randomAttrGradeList,omitempty" bson:"randomAttrGradeList"`       //属性评级权重（评级，权重）
	BindHideAttrNumWeight  []*RandomId `json:"bindHideAttrNumWeight,omitempty" bson:"bindHideAttrNumWeight"`   //绑定本源隐藏属性条目权重
	BindHideAttrRarityList []*RandomId `json:"bindHideAttrRarityList,omitempty" bson:"bindHideAttrRarityList"` //绑定本源隐藏属性品质权重
	HideAttrNumWeight      []*RandomId `json:"hideAttrNumWeight,omitempty" bson:"hideAttrNumWeight"`           //隐藏属性条目权重
	HideAttrRarityList     []*RandomId `json:"hideAttrRarityList,omitempty" bson:"hideAttrRarityList"`         //隐藏属性品质权重
	Cost                   []*IdCount  `json:"cost,omitempty" bson:"cost"`                                     //合成消耗道具
	ComposeProbability     float32     `json:"composeProbability,omitempty" bson:"composeProbability"`         //合成成功率
	IsBind                 bool        `json:"isBind,omitempty" bson:"isBind"`                                 //是否可以绑定
}

// HeroSoulRarity.json属性表集合
type HeroSoulRarityEntries struct {
	Rows   map[Rarity]*HeroSoulRarityEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroSoulRarityEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroSoulRarity.json", (*HeroSoulRarityEntries)(nil))
}

func (e *HeroSoulRarityEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroSoulRarityEntries = &HeroSoulRarityEntries{
		Rows:   make(map[Rarity]*HeroSoulRarityEntry),
		Arrays: make([]*HeroSoulRarityEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroSoulRarityEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroSoulRarityEntries, err
		}

		heroSoulRarityEntries.Rows[entry.Id] = entry
		heroSoulRarityEntries.Arrays = append(heroSoulRarityEntries.Arrays, entry)
	}

	sort.Slice(heroSoulRarityEntries.Arrays, func(i, j int) bool {
		return heroSoulRarityEntries.Arrays[i].Id < heroSoulRarityEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroSoulRarity.json", heroSoulRarityEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroSoulRarityEntries, nil

}

func (e *HeroSoulRarityEntries) GetArrays() any {
	return heroSoulRarityEntries.Arrays
}

func GetHeroSoulRarityEntry(id Rarity) (*HeroSoulRarityEntry, bool) {
	entry, ok := heroSoulRarityEntries.Rows[id]
	return entry, ok
}

func GetHeroSoulRaritySize() int {
	return int(len(heroSoulRarityEntries.Rows))
}

func GetHeroSoulRarityRows() map[Rarity]*HeroSoulRarityEntry {
	return heroSoulRarityEntries.Rows
}

func GetHeroSoulRarityArrays() []*HeroSoulRarityEntry {
	return heroSoulRarityEntries.Arrays
}
