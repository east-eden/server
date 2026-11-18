package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var attrSwitchEntries *AttrSwitchEntries //AttrSwitch.json全局变量

// 属性表 英雄表-拳皇/AttrSwitch
type AttrSwitchEntry struct {
	Id         AttrType   `json:"id,omitempty" bson:"id"`                 //属性Key
	AttrSet    []AttrType `json:"attrSet,omitempty" bson:"attrSet"`       //属性切换库
	SwitchCD   int32      `json:"switchCD,omitempty" bson:"switchCD"`     //切换CD(秒)
	SwitchCost []*IdCount `json:"switchCost,omitempty" bson:"switchCost"` //切换消耗
}

// AttrSwitch.json属性表集合
type AttrSwitchEntries struct {
	Rows   map[AttrType]*AttrSwitchEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*AttrSwitchEntry            `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("AttrSwitch.json", (*AttrSwitchEntries)(nil))
}

func (e *AttrSwitchEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	attrSwitchEntries = &AttrSwitchEntries{
		Rows:   make(map[AttrType]*AttrSwitchEntry),
		Arrays: make([]*AttrSwitchEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &AttrSwitchEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return attrSwitchEntries, err
		}

		attrSwitchEntries.Rows[entry.Id] = entry
		attrSwitchEntries.Arrays = append(attrSwitchEntries.Arrays, entry)
	}

	sort.Slice(attrSwitchEntries.Arrays, func(i, j int) bool {
		return attrSwitchEntries.Arrays[i].Id < attrSwitchEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("attrSwitch.json", attrSwitchEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return attrSwitchEntries, nil

}

func (e *AttrSwitchEntries) GetArrays() any {
	return attrSwitchEntries.Arrays
}

func GetAttrSwitchEntry(id AttrType) (*AttrSwitchEntry, bool) {
	entry, ok := attrSwitchEntries.Rows[id]
	return entry, ok
}

func GetAttrSwitchSize() int {
	return int(len(attrSwitchEntries.Rows))
}

func GetAttrSwitchRows() map[AttrType]*AttrSwitchEntry {
	return attrSwitchEntries.Rows
}

func GetAttrSwitchArrays() []*AttrSwitchEntry {
	return attrSwitchEntries.Arrays
}
