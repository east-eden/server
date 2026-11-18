package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var lDcollectionLevelEntries *LDcollectionLevelEntries //LDcollectionLevel.json全局变量

// 属性表 LD收藏表-拳皇/LDcollectionLevel
type LDcollectionLevelEntry struct {
	Id        int32    `json:"id,omitempty" bson:"id"`               //##等级
	CostChip  int32    `json:"costChip,omitempty" bson:"costChip"`   //升到下一级需要碎片数
	CostItems *IdCount `json:"costItems,omitempty" bson:"costItems"` //升级消耗道具
}

// LDcollectionLevel.json属性表集合
type LDcollectionLevelEntries struct {
	Rows   map[int32]*LDcollectionLevelEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*LDcollectionLevelEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("LDcollectionLevel.json", (*LDcollectionLevelEntries)(nil))
}

func (e *LDcollectionLevelEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	lDcollectionLevelEntries = &LDcollectionLevelEntries{
		Rows:   make(map[int32]*LDcollectionLevelEntry),
		Arrays: make([]*LDcollectionLevelEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &LDcollectionLevelEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return lDcollectionLevelEntries, err
		}

		lDcollectionLevelEntries.Rows[entry.Id] = entry
		lDcollectionLevelEntries.Arrays = append(lDcollectionLevelEntries.Arrays, entry)
	}

	sort.Slice(lDcollectionLevelEntries.Arrays, func(i, j int) bool {
		return lDcollectionLevelEntries.Arrays[i].Id < lDcollectionLevelEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("lDcollectionLevel.json", lDcollectionLevelEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return lDcollectionLevelEntries, nil

}

func (e *LDcollectionLevelEntries) GetArrays() any {
	return lDcollectionLevelEntries.Arrays
}

func GetLDcollectionLevelEntry(id int32) (*LDcollectionLevelEntry, bool) {
	entry, ok := lDcollectionLevelEntries.Rows[id]
	return entry, ok
}

func GetLDcollectionLevelSize() int {
	return int(len(lDcollectionLevelEntries.Rows))
}

func GetLDcollectionLevelRows() map[int32]*LDcollectionLevelEntry {
	return lDcollectionLevelEntries.Rows
}

func GetLDcollectionLevelArrays() []*LDcollectionLevelEntry {
	return lDcollectionLevelEntries.Arrays
}
