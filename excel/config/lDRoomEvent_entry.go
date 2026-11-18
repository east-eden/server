package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var lDRoomEventEntries *LDRoomEventEntries //LDRoomEvent.json全局变量

// 属性表 LD收藏表-拳皇/LDRoomEvent
type LDRoomEventEntry struct {
	Id             int32      `json:"id,omitempty" bson:"id"`                         //##id
	Name           string     `json:"name,omitempty" bson:"name"`                     //##事件名
	SubConditionId int32      `json:"subConditionId,omitempty" bson:"subConditionId"` //子条件id
	Awards         []*IdCount `json:"awards,omitempty" bson:"awards"`                 //达成奖励
	IsTogether     bool       `json:"isTogether,omitempty" bson:"isTogether"`         //是否共同计数
}

// LDRoomEvent.json属性表集合
type LDRoomEventEntries struct {
	Rows   map[int32]*LDRoomEventEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*LDRoomEventEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("LDRoomEvent.json", (*LDRoomEventEntries)(nil))
}

func (e *LDRoomEventEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	lDRoomEventEntries = &LDRoomEventEntries{
		Rows:   make(map[int32]*LDRoomEventEntry),
		Arrays: make([]*LDRoomEventEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &LDRoomEventEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return lDRoomEventEntries, err
		}

		lDRoomEventEntries.Rows[entry.Id] = entry
		lDRoomEventEntries.Arrays = append(lDRoomEventEntries.Arrays, entry)
	}

	sort.Slice(lDRoomEventEntries.Arrays, func(i, j int) bool {
		return lDRoomEventEntries.Arrays[i].Id < lDRoomEventEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("lDRoomEvent.json", lDRoomEventEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return lDRoomEventEntries, nil

}

func (e *LDRoomEventEntries) GetArrays() any {
	return lDRoomEventEntries.Arrays
}

func GetLDRoomEventEntry(id int32) (*LDRoomEventEntry, bool) {
	entry, ok := lDRoomEventEntries.Rows[id]
	return entry, ok
}

func GetLDRoomEventSize() int {
	return int(len(lDRoomEventEntries.Rows))
}

func GetLDRoomEventRows() map[int32]*LDRoomEventEntry {
	return lDRoomEventEntries.Rows
}

func GetLDRoomEventArrays() []*LDRoomEventEntry {
	return lDRoomEventEntries.Arrays
}
