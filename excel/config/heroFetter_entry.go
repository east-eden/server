package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroFetterEntries *HeroFetterEntries //HeroFetter.json全局变量

// 属性表 英雄表-拳皇/heroFetter
type HeroFetterEntry struct {
	Id                    int32       `json:"id,omitempty" bson:"id"`                                       //羁绊id
	HeroIds               []int32     `json:"heroIds,omitempty" bson:"heroIds"`                             //羁绊所需英雄id
	Name                  string      `json:"name,omitempty" bson:"name"`                                   //羁绊名字
	Leaders               []int32     `json:"leaders,omitempty" bson:"leaders"`                             //队长Id
	ActiveGotOne          bool        `json:"activeGotOne,omitempty" bson:"activeGotOne"`                   //获得任意英雄激活羁绊
	FetterSkill           int32       `json:"fetterSkill,omitempty" bson:"fetterSkill"`                     //连携技能ID
	AllianceTechInfoID    int32       `json:"allianceTechInfoID,omitempty" bson:"allianceTechInfoID"`       //联盟科技ID
	FetterSkillUnlockArgs [][]float32 `json:"fetterSkillUnlockArgs,omitempty" bson:"fetterSkillUnlockArgs"` //连携技解锁条件
}

// HeroFetter.json属性表集合
type HeroFetterEntries struct {
	Rows   map[int32]*HeroFetterEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroFetterEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroFetter.json", (*HeroFetterEntries)(nil))
}

func (e *HeroFetterEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroFetterEntries = &HeroFetterEntries{
		Rows:   make(map[int32]*HeroFetterEntry),
		Arrays: make([]*HeroFetterEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroFetterEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroFetterEntries, err
		}

		heroFetterEntries.Rows[entry.Id] = entry
		heroFetterEntries.Arrays = append(heroFetterEntries.Arrays, entry)
	}

	sort.Slice(heroFetterEntries.Arrays, func(i, j int) bool {
		return heroFetterEntries.Arrays[i].Id < heroFetterEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroFetter.json", heroFetterEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroFetterEntries, nil

}

func (e *HeroFetterEntries) GetArrays() any {
	return heroFetterEntries.Arrays
}

func GetHeroFetterEntry(id int32) (*HeroFetterEntry, bool) {
	entry, ok := heroFetterEntries.Rows[id]
	return entry, ok
}

func GetHeroFetterSize() int {
	return int(len(heroFetterEntries.Rows))
}

func GetHeroFetterRows() map[int32]*HeroFetterEntry {
	return heroFetterEntries.Rows
}

func GetHeroFetterArrays() []*HeroFetterEntry {
	return heroFetterEntries.Arrays
}
