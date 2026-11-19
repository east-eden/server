package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var allianceLevelEntries *AllianceLevelEntries //AllianceLevel.json全局变量

// 属性表 联盟-拳皇/AllianceLevel
type AllianceLevelEntry struct {
	Id                   int32 `json:"id,omitempty" bson:"id"`                                     //##等级
	NextNeedExp          int32 `json:"nextNeedExp,omitempty" bson:"nextNeedExp"`                   //升到下一级需要的经验
	MaxMemberCount       int32 `json:"maxMemberCount,omitempty" bson:"maxMemberCount"`             //公会人数上限
	MaxDeputyLeaderCount int32 `json:"maxDeputyLeaderCount,omitempty" bson:"maxDeputyLeaderCount"` //最大副会长数
	MaxEliteCount        int32 `json:"maxEliteCount,omitempty" bson:"maxEliteCount"`               //最大精英数
}

// AllianceLevel.json属性表集合
type AllianceLevelEntries struct {
	Rows   map[int32]*AllianceLevelEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*AllianceLevelEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("AllianceLevel.json", (*AllianceLevelEntries)(nil))
}

func (e *AllianceLevelEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	allianceLevelEntries = &AllianceLevelEntries{
		Rows:   make(map[int32]*AllianceLevelEntry),
		Arrays: make([]*AllianceLevelEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &AllianceLevelEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return allianceLevelEntries, err
		}

		allianceLevelEntries.Rows[entry.Id] = entry
		allianceLevelEntries.Arrays = append(allianceLevelEntries.Arrays, entry)
	}

	sort.Slice(allianceLevelEntries.Arrays, func(i, j int) bool {
		return allianceLevelEntries.Arrays[i].Id < allianceLevelEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("allianceLevel.json", allianceLevelEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return allianceLevelEntries, nil

}

func (e *AllianceLevelEntries) GetArrays() any {
	return allianceLevelEntries.Arrays
}

func GetAllianceLevelEntry(id int32) (*AllianceLevelEntry, bool) {
	entry, ok := allianceLevelEntries.Rows[id]
	return entry, ok
}

func GetAllianceLevelSize() int {
	return int(len(allianceLevelEntries.Rows))
}

func GetAllianceLevelRows() map[int32]*AllianceLevelEntry {
	return allianceLevelEntries.Rows
}

func GetAllianceLevelArrays() []*AllianceLevelEntry {
	return allianceLevelEntries.Arrays
}
