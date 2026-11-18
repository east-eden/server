package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroQualityEntries *HeroQualityEntries //HeroQuality.json全局变量

// 属性表 英雄表-拳皇/heroQuality
type HeroQualityEntry struct {
	Id                      Rarity     `json:"id,omitempty" bson:"id"`                                           //##id
	FirstGotAtlasAwards     []*IdCount `json:"firstGotAtlasAwards,omitempty" bson:"firstGotAtlasAwards"`         //首次获得英雄图鉴奖励
	FirstGotAtlasBadgePoint int32      `json:"firstGotAtlasBadgePoint,omitempty" bson:"firstGotAtlasBadgePoint"` //首次获得英雄羁绊点数奖励
	StarAtlasAward          int32      `json:"starAtlasAward,omitempty" bson:"starAtlasAward"`                   //每次英雄升星图鉴奖励
	ComposeChipCnt          int32      `json:"composeChipCnt,omitempty" bson:"composeChipCnt"`                   //合成英雄需要碎片数
	ConvertChipCnt          int32      `json:"convertChipCnt,omitempty" bson:"convertChipCnt"`                   //重复获得卡牌转换碎片数
	PromotionNextCost       []*IdCount `json:"promotionNextCost,omitempty" bson:"promotionNextCost"`             //晋升下一品质消耗道具
	BaseSkillChipCost       []int32    `json:"baseSkillChipCost,omitempty" bson:"baseSkillChipCost"`             //四挡基础技能的本体碎片消耗数
	SpecialSkillChipCost    [][]int32  `json:"specialSkillChipCost,omitempty" bson:"specialSkillChipCost"`       //四挡觉醒特殊技能本体碎片消耗
	SpecialSkillStoneCost   [][]int32  `json:"specialSkillStoneCost,omitempty" bson:"specialSkillStoneCost"`     //四挡特殊技能觉醒石消耗
}

// HeroQuality.json属性表集合
type HeroQualityEntries struct {
	Rows   map[Rarity]*HeroQualityEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroQualityEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroQuality.json", (*HeroQualityEntries)(nil))
}

func (e *HeroQualityEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroQualityEntries = &HeroQualityEntries{
		Rows:   make(map[Rarity]*HeroQualityEntry),
		Arrays: make([]*HeroQualityEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroQualityEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroQualityEntries, err
		}

		heroQualityEntries.Rows[entry.Id] = entry
		heroQualityEntries.Arrays = append(heroQualityEntries.Arrays, entry)
	}

	sort.Slice(heroQualityEntries.Arrays, func(i, j int) bool {
		return heroQualityEntries.Arrays[i].Id < heroQualityEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroQuality.json", heroQualityEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroQualityEntries, nil

}

func (e *HeroQualityEntries) GetArrays() any {
	return heroQualityEntries.Arrays
}

func GetHeroQualityEntry(id Rarity) (*HeroQualityEntry, bool) {
	entry, ok := heroQualityEntries.Rows[id]
	return entry, ok
}

func GetHeroQualitySize() int {
	return int(len(heroQualityEntries.Rows))
}

func GetHeroQualityRows() map[Rarity]*HeroQualityEntry {
	return heroQualityEntries.Rows
}

func GetHeroQualityArrays() []*HeroQualityEntry {
	return heroQualityEntries.Arrays
}
