package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var crystalSuitEntries *CrystalSuitEntries //CrystalSuit.xlsx全局变量

// CrystalSuit.xlsx属性表
type CrystalSuitEntry struct {
	Id          int32 `json:"Id,omitempty"`          // 主键
	Suit2_AttID int32 `json:"Suit2_AttID,omitempty"` //2件套装属性ID
	Suit3_AttID int32 `json:"Suit3_AttID,omitempty"` //3件套装属性ID
	Suit4_AttID int32 `json:"Suit4_AttID,omitempty"` //2件套装属性ID
	Suit5_AttID int32 `json:"Suit5_AttID,omitempty"` //2件套装属性ID
	Suit6_AttID int32 `json:"Suit6_AttID,omitempty"` //2件套装属性ID
}

// CrystalSuit.xlsx属性表集合
type CrystalSuitEntries struct {
	Rows map[int32]*CrystalSuitEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("CrystalSuit.xlsx", (*CrystalSuitEntries)(nil))
}

func (e *CrystalSuitEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	crystalSuitEntries = &CrystalSuitEntries{
		Rows: make(map[int32]*CrystalSuitEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CrystalSuitEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		crystalSuitEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetCrystalSuitEntry(id int32) (*CrystalSuitEntry, bool) {
	entry, ok := crystalSuitEntries.Rows[id]
	return entry, ok
}

func GetCrystalSuitSize() int32 {
	return int32(len(crystalSuitEntries.Rows))
}

func GetCrystalSuitRows() map[int32]*CrystalSuitEntry {
	return crystalSuitEntries.Rows
}
