package auto

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var unitGroupEntries *UnitGroupEntries //UnitGroup.xlsx全局变量

// UnitGroup.xlsx属性表
type UnitGroupEntry struct {
	Id     int32 `json:"Id,omitempty"`     // 主键
	UnitId int32 `json:"UnitId,omitempty"` // 多主键之一
	PosX   int32 `json:"PosX,omitempty"`   //怪物x坐标
	PosY   int32 `json:"PosY,omitempty"`   //怪物y坐标
	PosZ   int32 `json:"PosZ,omitempty"`   //怪物z坐标
	Rotate int32 `json:"Rotate,omitempty"` //怪物朝向
}

// UnitGroup.xlsx属性表集合
type UnitGroupEntries struct {
	Rows map[string]*UnitGroupEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("UnitGroup.xlsx", (*UnitGroupEntries)(nil))
}

func (e *UnitGroupEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	unitGroupEntries = &UnitGroupEntries{
		Rows: make(map[string]*UnitGroupEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &UnitGroupEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		key := fmt.Sprintf("%d+%d", entry.Id, entry.UnitId)
		unitGroupEntries.Rows[key] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetUnitGroupEntry(keys ...int32) (*UnitGroupEntry, bool) {
	keyName := make([]string, 0, len(keys))
	for _, key := range keys {
		keyName = append(keyName, strconv.Itoa(int(key)))
	}

	finalKey := strings.Join(keyName, "+")
	entry, ok := unitGroupEntries.Rows[finalKey]
	return entry, ok
}

func GetUnitGroupSize() int32 {
	return int32(len(unitGroupEntries.Rows))
}

func GetUnitGroupRows() map[string]*UnitGroupEntry {
	return unitGroupEntries.Rows
}
