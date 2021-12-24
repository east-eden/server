package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var crystalEntries *CrystalEntries //Crystal.csv全局变量

// Crystal.csv属性表
type CrystalEntry struct {
	Id  int32 `json:"Id,omitempty"`  // 主键
	Pos int32 `json:"Pos,omitempty"` //位置
}

// Crystal.csv属性表集合
type CrystalEntries struct {
	Rows map[int32]*CrystalEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Crystal.csv", (*CrystalEntries)(nil))
}

func (e *CrystalEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	crystalEntries = &CrystalEntries{
		Rows: make(map[int32]*CrystalEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CrystalEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		crystalEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetCrystalEntry(id int32) (*CrystalEntry, bool) {
	entry, ok := crystalEntries.Rows[id]
	return entry, ok
}

func GetCrystalSize() int32 {
	return int32(len(crystalEntries.Rows))
}

func GetCrystalRows() map[int32]*CrystalEntry {
	return crystalEntries.Rows
}
