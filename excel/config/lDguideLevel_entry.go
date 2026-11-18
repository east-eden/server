package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var lDguideLevelEntries *LDguideLevelEntries //LDguideLevel.json全局变量

// 属性表 LD关卡表-拳皇/LDguideLevel
type LDguideLevelEntry struct {
	Test                         int32           `json:"test,omitempty" bson:"test"`                                                 //##保留主键id
	GuideNextCertainFormation    map[int32]int32 `json:"guideNextCertainFormation,omitempty" bson:"guideNextCertainFormation"`       //下一次阵容刷新出指定神话英雄和等级
	GuideSummonCertainHero       map[int32]int32 `json:"guideSummonCertainHero,omitempty" bson:"guideSummonCertainHero"`             //指定召唤次数出指定ID的英雄
	NewbieGuideCertainFormation  map[int32]int32 `json:"newbieGuideCertainFormation,omitempty" bson:"newbieGuideCertainFormation"`   //新手引导下一次阵容刷新出指定神话英雄和等级
	NewbieGuideSummonCertainHero map[int32]int32 `json:"newbieGuideSummonCertainHero,omitempty" bson:"newbieGuideSummonCertainHero"` //新手引导指定召唤次数出指定ID的英雄
}

// LDguideLevel.json属性表集合
type LDguideLevelEntries struct {
	Rows   map[int32]*LDguideLevelEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*LDguideLevelEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("LDguideLevel.json", (*LDguideLevelEntries)(nil))
}

func (e *LDguideLevelEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	lDguideLevelEntries = &LDguideLevelEntries{
		Rows:   make(map[int32]*LDguideLevelEntry),
		Arrays: make([]*LDguideLevelEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &LDguideLevelEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return lDguideLevelEntries, err
		}

		lDguideLevelEntries.Rows[entry.Test] = entry
		lDguideLevelEntries.Arrays = append(lDguideLevelEntries.Arrays, entry)
	}

	sort.Slice(lDguideLevelEntries.Arrays, func(i, j int) bool {
		return lDguideLevelEntries.Arrays[i].Test < lDguideLevelEntries.Arrays[j].Test
	})

	excel.AddArrayGetter("lDguideLevel.json", lDguideLevelEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return lDguideLevelEntries, nil

}

func (e *LDguideLevelEntries) GetArrays() any {
	return lDguideLevelEntries.Arrays
}

func GetLDguideLevelEntry() *LDguideLevelEntry {
	entry := lDguideLevelEntries.Rows[0]
	return entry
}

func GetLDguideLevelSize() int {
	return int(len(lDguideLevelEntries.Rows))
}

func GetLDguideLevelRows() map[int32]*LDguideLevelEntry {
	return lDguideLevelEntries.Rows
}

func GetLDguideLevelArrays() []*LDguideLevelEntry {
	return lDguideLevelEntries.Arrays
}
