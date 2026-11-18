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

var heroSoulHideAttrEntries *HeroSoulHideAttrEntries //HeroSoulHideAttr.json全局变量

// 属性表 英雄表-拳皇/HeroSoulHideAttr
type HeroSoulHideAttrEntry struct {
	Id                int32              `json:"id,omitempty" bson:"id"`                               //隐藏属性id
	Rarity            Rarity             `json:"rarity,omitempty" bson:"rarity"`                       //品质
	HideClass         int32              `json:"hideClass,omitempty" bson:"hideClass"`                 //隐藏类别
	AttrType          AttrType           `json:"attrType,omitempty" bson:"attrType"`                   //数值属性类型
	EffectType        HeroSoulEffectType `json:"effectType,omitempty" bson:"effectType"`               //效果属性
	EffectRange       *FloatPair         `json:"effectRange,omitempty" bson:"effectRange"`             //效果范围
	HideBaseMark      float32            `json:"hideBaseMark,omitempty" bson:"hideBaseMark"`           //词条基础分
	HideRandomMarkMax float32            `json:"hideRandomMarkMax,omitempty" bson:"hideRandomMarkMax"` //词条随机最大分值
	TargetType        BonusTargetType    `json:"targetType,omitempty" bson:"targetType"`               //作用对象
	TargetId          int32              `json:"targetId,omitempty" bson:"targetId"`                   //目标对象
}

// HeroSoulHideAttr.json属性表集合
type HeroSoulHideAttrEntries struct {
	Rows   map[string]*HeroSoulHideAttrEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroSoulHideAttrEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroSoulHideAttr.json", (*HeroSoulHideAttrEntries)(nil))
}

func (e *HeroSoulHideAttrEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroSoulHideAttrEntries = &HeroSoulHideAttrEntries{
		Rows:   make(map[string]*HeroSoulHideAttrEntry),
		Arrays: make([]*HeroSoulHideAttrEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroSoulHideAttrEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroSoulHideAttrEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Rarity))
		heroSoulHideAttrEntries.Rows[key] = entry
		heroSoulHideAttrEntries.Arrays = append(heroSoulHideAttrEntries.Arrays, entry)
	}

	sort.Slice(heroSoulHideAttrEntries.Arrays, func(i, j int) bool {
		return heroSoulHideAttrEntries.Arrays[i].Id < heroSoulHideAttrEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroSoulHideAttr.json", heroSoulHideAttrEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroSoulHideAttrEntries, nil

}

func (e *HeroSoulHideAttrEntries) GetArrays() any {
	return heroSoulHideAttrEntries.Arrays
}

func GetHeroSoulHideAttrEntry(keys ...string) (*HeroSoulHideAttrEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := heroSoulHideAttrEntries.Rows[finalKey]
	return entry, ok
}

func GetHeroSoulHideAttrSize() int {
	return int(len(heroSoulHideAttrEntries.Rows))
}

func GetHeroSoulHideAttrRows() map[string]*HeroSoulHideAttrEntry {
	return heroSoulHideAttrEntries.Rows
}

func GetHeroSoulHideAttrArrays() []*HeroSoulHideAttrEntry {
	return heroSoulHideAttrEntries.Arrays
}
