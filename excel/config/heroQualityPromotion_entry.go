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

var heroQualityPromotionEntries *HeroQualityPromotionEntries //HeroQualityPromotion.json全局变量

// 属性表 英雄表-拳皇/heroQualityPromotion
type HeroQualityPromotionEntry struct {
	Id                Rarity     `json:"id,omitempty" bson:"id"`                               //##id
	Level             int32      `json:"level,omitempty" bson:"level"`                         //品质等级
	PromotionNextCost []*IdCount `json:"promotionNextCost,omitempty" bson:"promotionNextCost"` //晋升下一品质消耗道具
	PromotionLastCost []*IdCount `json:"promotionLastCost,omitempty" bson:"promotionLastCost"` //上一品质回退到此品质所需道具
}

// HeroQualityPromotion.json属性表集合
type HeroQualityPromotionEntries struct {
	Rows   map[string]*HeroQualityPromotionEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroQualityPromotionEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroQualityPromotion.json", (*HeroQualityPromotionEntries)(nil))
}

func (e *HeroQualityPromotionEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroQualityPromotionEntries = &HeroQualityPromotionEntries{
		Rows:   make(map[string]*HeroQualityPromotionEntry),
		Arrays: make([]*HeroQualityPromotionEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroQualityPromotionEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroQualityPromotionEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Level))
		heroQualityPromotionEntries.Rows[key] = entry
		heroQualityPromotionEntries.Arrays = append(heroQualityPromotionEntries.Arrays, entry)
	}

	sort.Slice(heroQualityPromotionEntries.Arrays, func(i, j int) bool {
		return heroQualityPromotionEntries.Arrays[i].Id < heroQualityPromotionEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroQualityPromotion.json", heroQualityPromotionEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroQualityPromotionEntries, nil

}

func (e *HeroQualityPromotionEntries) GetArrays() any {
	return heroQualityPromotionEntries.Arrays
}

func GetHeroQualityPromotionEntry(keys ...string) (*HeroQualityPromotionEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroQualityPromotionEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroQualityPromotionSize() int {
	return int(len(heroQualityPromotionEntries.Rows))
}

func GetHeroQualityPromotionRows() map[string]*HeroQualityPromotionEntry {
	return heroQualityPromotionEntries.Rows
}

func GetHeroQualityPromotionArrays() []*HeroQualityPromotionEntry {
	return heroQualityPromotionEntries.Arrays
}
