package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var attrEntries *AttrEntries //Attr.json全局变量

// 属性表 基础配置-拳皇/Attr
type AttrEntry struct {
	Id        AttrType      `json:"id,omitempty" bson:"id"`               //##id
	TitleType AttrTitleType `json:"titleType,omitempty" bson:"titleType"` //属性标题类型
	IsPercent bool          `json:"isPercent,omitempty" bson:"isPercent"` //是否百分比
	Value     float32       `json:"value,omitempty" bson:"value"`         //战力价值
}

// Attr.json属性表集合
type AttrEntries struct {
	Rows   map[AttrType]*AttrEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*AttrEntry            `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Attr.json", (*AttrEntries)(nil))
}

func (e *AttrEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	attrEntries = &AttrEntries{
		Rows:   make(map[AttrType]*AttrEntry),
		Arrays: make([]*AttrEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &AttrEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return attrEntries, err
		}

		attrEntries.Rows[entry.Id] = entry
		attrEntries.Arrays = append(attrEntries.Arrays, entry)
	}

	sort.Slice(attrEntries.Arrays, func(i, j int) bool {
		return attrEntries.Arrays[i].Id < attrEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("attr.json", attrEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return attrEntries, nil

}

func (e *AttrEntries) GetArrays() any {
	return attrEntries.Arrays
}

func GetAttrEntry(id AttrType) (*AttrEntry, bool) {
	entry, ok := attrEntries.Rows[id]
	return entry, ok
}

func GetAttrSize() int {
	return int(len(attrEntries.Rows))
}

func GetAttrRows() map[AttrType]*AttrEntry {
	return attrEntries.Rows
}

func GetAttrArrays() []*AttrEntry {
	return attrEntries.Arrays
}
