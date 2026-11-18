package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var specialAttrGroupEntries *SpecialAttrGroupEntries //SpecialAttrGroup.json全局变量

// 属性表 音巢科技表-拳皇/SpecialAttrGroup
type SpecialAttrGroupEntry struct {
	AttrType AttrType `json:"attrType,omitempty" bson:"attrType"` //属性类型
	Group    int32    `json:"group,omitempty" bson:"group"`       //组
}

// SpecialAttrGroup.json属性表集合
type SpecialAttrGroupEntries struct {
	Rows   map[AttrType]*SpecialAttrGroupEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*SpecialAttrGroupEntry            `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("SpecialAttrGroup.json", (*SpecialAttrGroupEntries)(nil))
}

func (e *SpecialAttrGroupEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	specialAttrGroupEntries = &SpecialAttrGroupEntries{
		Rows:   make(map[AttrType]*SpecialAttrGroupEntry),
		Arrays: make([]*SpecialAttrGroupEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &SpecialAttrGroupEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return specialAttrGroupEntries, err
		}

		specialAttrGroupEntries.Rows[entry.AttrType] = entry
		specialAttrGroupEntries.Arrays = append(specialAttrGroupEntries.Arrays, entry)
	}

	sort.Slice(specialAttrGroupEntries.Arrays, func(i, j int) bool {
		return specialAttrGroupEntries.Arrays[i].AttrType < specialAttrGroupEntries.Arrays[j].AttrType
	})

	excel.AddArrayGetter("specialAttrGroup.json", specialAttrGroupEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return specialAttrGroupEntries, nil

}

func (e *SpecialAttrGroupEntries) GetArrays() any {
	return specialAttrGroupEntries.Arrays
}

func GetSpecialAttrGroupEntry(id AttrType) (*SpecialAttrGroupEntry, bool) {
	entry, ok := specialAttrGroupEntries.Rows[id]
	return entry, ok
}

func GetSpecialAttrGroupSize() int {
	return int(len(specialAttrGroupEntries.Rows))
}

func GetSpecialAttrGroupRows() map[AttrType]*SpecialAttrGroupEntry {
	return specialAttrGroupEntries.Rows
}

func GetSpecialAttrGroupArrays() []*SpecialAttrGroupEntry {
	return specialAttrGroupEntries.Arrays
}
