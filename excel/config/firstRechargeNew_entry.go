package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

var firstRechargeNewEntries *FirstRechargeNewEntries //FirstRechargeNew.json全局变量

// 属性表 礼包表-拳皇/FirstRechargeNew
type FirstRechargeNewEntry struct {
	Id     int32      `json:"id,omitempty" bson:"id"`         //##商品id
	Day    int32      `json:"day,omitempty" bson:"day"`       //第x天
	Awards []*IdCount `json:"awards,omitempty" bson:"awards"` //奖励列表
}

// FirstRechargeNew.json属性表集合
type FirstRechargeNewEntries struct {
	Rows   map[string]*FirstRechargeNewEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*FirstRechargeNewEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("FirstRechargeNew.json", (*FirstRechargeNewEntries)(nil))
}

func (e *FirstRechargeNewEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	firstRechargeNewEntries = &FirstRechargeNewEntries{
		Rows:   make(map[string]*FirstRechargeNewEntry),
		Arrays: make([]*FirstRechargeNewEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &FirstRechargeNewEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return firstRechargeNewEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Day))
		firstRechargeNewEntries.Rows[key] = entry
		firstRechargeNewEntries.Arrays = append(firstRechargeNewEntries.Arrays, entry)
	}

	sort.Slice(firstRechargeNewEntries.Arrays, func(i, j int) bool {
		return firstRechargeNewEntries.Arrays[i].Id < firstRechargeNewEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("firstRechargeNew.json", firstRechargeNewEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return firstRechargeNewEntries, nil

}

func (e *FirstRechargeNewEntries) GetArrays() any {
	return firstRechargeNewEntries.Arrays
}

func GetFirstRechargeNewEntry(keys ...string) (*FirstRechargeNewEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := firstRechargeNewEntries.Rows[finalKey]
	return entry, ok
}

func GetFirstRechargeNewSize() int {
	return int(len(firstRechargeNewEntries.Rows))
}

func GetFirstRechargeNewRows() map[string]*FirstRechargeNewEntry {
	return firstRechargeNewEntries.Rows
}

func GetFirstRechargeNewArrays() []*FirstRechargeNewEntry {
	return firstRechargeNewEntries.Arrays
}
