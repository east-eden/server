package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroArtifactEntries *HeroArtifactEntries //HeroArtifact.json全局变量

// 属性表 英雄表-拳皇/HeroArtifact
type HeroArtifactEntry struct {
	Id      int32        `json:"id,omitempty" bson:"id"`           //##神器id
	Quality Rarity       `json:"quality,omitempty" bson:"quality"` //品质
	HeroId  int32        `json:"heroId,omitempty" bson:"heroId"`   //专属英雄
	Attr    []*AttrValue `json:"attr,omitempty" bson:"attr"`       //属性
}

// HeroArtifact.json属性表集合
type HeroArtifactEntries struct {
	Rows   map[int32]*HeroArtifactEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroArtifactEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroArtifact.json", (*HeroArtifactEntries)(nil))
}

func (e *HeroArtifactEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroArtifactEntries = &HeroArtifactEntries{
		Rows:   make(map[int32]*HeroArtifactEntry),
		Arrays: make([]*HeroArtifactEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroArtifactEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroArtifactEntries, err
		}

		heroArtifactEntries.Rows[entry.Id] = entry
		heroArtifactEntries.Arrays = append(heroArtifactEntries.Arrays, entry)
	}

	sort.Slice(heroArtifactEntries.Arrays, func(i, j int) bool {
		return heroArtifactEntries.Arrays[i].Id < heroArtifactEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroArtifact.json", heroArtifactEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroArtifactEntries, nil

}

func (e *HeroArtifactEntries) GetArrays() any {
	return heroArtifactEntries.Arrays
}

func GetHeroArtifactEntry(id int32) (*HeroArtifactEntry, bool) {
	entry, ok := heroArtifactEntries.Rows[id]
	return entry, ok
}

func GetHeroArtifactSize() int {
	return int(len(heroArtifactEntries.Rows))
}

func GetHeroArtifactRows() map[int32]*HeroArtifactEntry {
	return heroArtifactEntries.Rows
}

func GetHeroArtifactArrays() []*HeroArtifactEntry {
	return heroArtifactEntries.Arrays
}
