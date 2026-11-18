package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroAtlasBadgeEntries *HeroAtlasBadgeEntries //HeroAtlasBadge.json全局变量

// 属性表 英雄表-拳皇/heroAtlasBadge
type HeroAtlasBadgeEntry struct {
	Id     int32      `json:"id,omitempty" bson:"id"`         //进度id
	Score  int32      `json:"score,omitempty" bson:"score"`   //需要积分
	Awards []*IdCount `json:"awards,omitempty" bson:"awards"` //奖励
}

// HeroAtlasBadge.json属性表集合
type HeroAtlasBadgeEntries struct {
	Rows   map[int32]*HeroAtlasBadgeEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroAtlasBadgeEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroAtlasBadge.json", (*HeroAtlasBadgeEntries)(nil))
}

func (e *HeroAtlasBadgeEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroAtlasBadgeEntries = &HeroAtlasBadgeEntries{
		Rows:   make(map[int32]*HeroAtlasBadgeEntry),
		Arrays: make([]*HeroAtlasBadgeEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroAtlasBadgeEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroAtlasBadgeEntries, err
		}

		heroAtlasBadgeEntries.Rows[entry.Id] = entry
		heroAtlasBadgeEntries.Arrays = append(heroAtlasBadgeEntries.Arrays, entry)
	}

	sort.Slice(heroAtlasBadgeEntries.Arrays, func(i, j int) bool {
		return heroAtlasBadgeEntries.Arrays[i].Id < heroAtlasBadgeEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroAtlasBadge.json", heroAtlasBadgeEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroAtlasBadgeEntries, nil

}

func (e *HeroAtlasBadgeEntries) GetArrays() any {
	return heroAtlasBadgeEntries.Arrays
}

func GetHeroAtlasBadgeEntry(id int32) (*HeroAtlasBadgeEntry, bool) {
	entry, ok := heroAtlasBadgeEntries.Rows[id]
	return entry, ok
}

func GetHeroAtlasBadgeSize() int {
	return int(len(heroAtlasBadgeEntries.Rows))
}

func GetHeroAtlasBadgeRows() map[int32]*HeroAtlasBadgeEntry {
	return heroAtlasBadgeEntries.Rows
}

func GetHeroAtlasBadgeArrays() []*HeroAtlasBadgeEntry {
	return heroAtlasBadgeEntries.Arrays
}
