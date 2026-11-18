package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var lDSkillEntries *LDSkillEntries //LDSkill.json全局变量

// 属性表 LD技能表-拳皇/LDSkill
type LDSkillEntry struct {
	Id          int32                  `json:"id,omitempty" bson:"id"`                   //##id
	AttValue    map[LDAttrType]float32 `json:"attValue,omitempty" bson:"attValue"`       //技能属性加成
	MutexSkills []int32                `json:"mutexSkills,omitempty" bson:"mutexSkills"` //互斥技能id
}

// LDSkill.json属性表集合
type LDSkillEntries struct {
	Rows   map[int32]*LDSkillEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*LDSkillEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("LDSkill.json", (*LDSkillEntries)(nil))
}

func (e *LDSkillEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	lDSkillEntries = &LDSkillEntries{
		Rows:   make(map[int32]*LDSkillEntry),
		Arrays: make([]*LDSkillEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &LDSkillEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return lDSkillEntries, err
		}

		lDSkillEntries.Rows[entry.Id] = entry
		lDSkillEntries.Arrays = append(lDSkillEntries.Arrays, entry)
	}

	sort.Slice(lDSkillEntries.Arrays, func(i, j int) bool {
		return lDSkillEntries.Arrays[i].Id < lDSkillEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("lDSkill.json", lDSkillEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return lDSkillEntries, nil

}

func (e *LDSkillEntries) GetArrays() any {
	return lDSkillEntries.Arrays
}

func GetLDSkillEntry(id int32) (*LDSkillEntry, bool) {
	entry, ok := lDSkillEntries.Rows[id]
	return entry, ok
}

func GetLDSkillSize() int {
	return int(len(lDSkillEntries.Rows))
}

func GetLDSkillRows() map[int32]*LDSkillEntry {
	return lDSkillEntries.Rows
}

func GetLDSkillArrays() []*LDSkillEntry {
	return lDSkillEntries.Arrays
}
