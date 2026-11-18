package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroRightsEntries *HeroRightsEntries //HeroRights.json全局变量

// 属性表 英雄表-拳皇/HeroRights
type HeroRightsEntry struct {
	Id                  int32          `json:"id,omitempty" bson:"id"`                                   //权益id
	ItemId              int32          `json:"itemId,omitempty" bson:"itemId"`                           //对应道具id
	HeroFashionRightsId int32          `json:"heroFashionRightsId,omitempty" bson:"heroFashionRightsId"` //对应权益时装id
	RepeatOffset        []*IdCount     `json:"repeatOffset,omitempty" bson:"repeatOffset"`               //重复购买补偿
	HeroId              int32          `json:"heroId,omitempty" bson:"heroId"`                           //对应英雄
	Type                HeroRightsType `json:"type,omitempty" bson:"type"`                               //权益类型
}

// HeroRights.json属性表集合
type HeroRightsEntries struct {
	Rows   map[int32]*HeroRightsEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroRightsEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroRights.json", (*HeroRightsEntries)(nil))
}

func (e *HeroRightsEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroRightsEntries = &HeroRightsEntries{
		Rows:   make(map[int32]*HeroRightsEntry),
		Arrays: make([]*HeroRightsEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroRightsEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroRightsEntries, err
		}

		heroRightsEntries.Rows[entry.Id] = entry
		heroRightsEntries.Arrays = append(heroRightsEntries.Arrays, entry)
	}

	sort.Slice(heroRightsEntries.Arrays, func(i, j int) bool {
		return heroRightsEntries.Arrays[i].Id < heroRightsEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroRights.json", heroRightsEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroRightsEntries, nil

}

func (e *HeroRightsEntries) GetArrays() any {
	return heroRightsEntries.Arrays
}

func GetHeroRightsEntry(id int32) (*HeroRightsEntry, bool) {
	entry, ok := heroRightsEntries.Rows[id]
	return entry, ok
}

func GetHeroRightsSize() int {
	return int(len(heroRightsEntries.Rows))
}

func GetHeroRightsRows() map[int32]*HeroRightsEntry {
	return heroRightsEntries.Rows
}

func GetHeroRightsArrays() []*HeroRightsEntry {
	return heroRightsEntries.Arrays
}
