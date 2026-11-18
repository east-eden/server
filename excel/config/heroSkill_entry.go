package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroSkillEntries *HeroSkillEntries //HeroSkill.json全局变量

// 属性表 英雄表-拳皇/heroSkill
type HeroSkillEntry struct {
	Id                   int32                `json:"id,omitempty" bson:"id"`                                     //##id
	UnlockStar           int32                `json:"unlockStar,omitempty" bson:"unlockStar"`                     //解锁星级
	UnlockRarity         Rarity               `json:"unlockRarity,omitempty" bson:"unlockRarity"`                 //解锁品质
	SkillIds             []int32              `json:"skillIds,omitempty" bson:"skillIds"`                         //战斗技能id
	Type                 HeroSkillType        `json:"type,omitempty" bson:"type"`                                 //技能类型
	UnlockRaritys        []Rarity             `json:"unlockRaritys,omitempty" bson:"unlockRaritys"`               //神器品质解锁
	UnlockSkillIds       [][]int32            `json:"unlockSkillIds,omitempty" bson:"unlockSkillIds"`             //神器解锁技能id
	UnlockListenSkillIds [][]int32            `json:"unlockListenSkillIds,omitempty" bson:"unlockListenSkillIds"` //神器解锁监听技能id
	UnlockDescs          []string             `json:"unlockDescs,omitempty" bson:"unlockDescs"`                   //神器解锁技能描述
	AttrList             map[AttrType]float32 `json:"attrList,omitempty" bson:"attrList"`                         //被动属性效果列表|攻击百分比,0.2|血量百分比,0.2
	SoulPos              int32                `json:"soulPos,omitempty" bson:"soulPos"`                           //八门位置限制
	SoulBonus            float32              `json:"soulBonus,omitempty" bson:"soulBonus"`                       //加成百分比
}

// HeroSkill.json属性表集合
type HeroSkillEntries struct {
	Rows   map[int32]*HeroSkillEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroSkillEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("HeroSkill.json", (*HeroSkillEntries)(nil))
}

func (e *HeroSkillEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroSkillEntries = &HeroSkillEntries{
		Rows:   make(map[int32]*HeroSkillEntry),
		Arrays: make([]*HeroSkillEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroSkillEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroSkillEntries, err
		}

		heroSkillEntries.Rows[entry.Id] = entry
		heroSkillEntries.Arrays = append(heroSkillEntries.Arrays, entry)
	}

	sort.Slice(heroSkillEntries.Arrays, func(i, j int) bool {
		return heroSkillEntries.Arrays[i].Id < heroSkillEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("heroSkill.json", heroSkillEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroSkillEntries, nil

}

func (e *HeroSkillEntries) GetArrays() any {
	return heroSkillEntries.Arrays
}

func GetHeroSkillEntry(id int32) (*HeroSkillEntry, bool) {
	entry, ok := heroSkillEntries.Rows[id]
	return entry, ok
}

func GetHeroSkillSize() int {
	return int(len(heroSkillEntries.Rows))
}

func GetHeroSkillRows() map[int32]*HeroSkillEntry {
	return heroSkillEntries.Rows
}

func GetHeroSkillArrays() []*HeroSkillEntry {
	return heroSkillEntries.Arrays
}
