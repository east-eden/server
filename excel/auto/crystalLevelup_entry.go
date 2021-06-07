package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var crystalLevelupEntries *CrystalLevelupEntries //CrystalLevelup.xlsx全局变量

// CrystalLevelup.xlsx属性表
type CrystalLevelupEntry struct {
	Id                int32   `json:"Id,omitempty"`                // 主键
	Exp               []int32 `json:"Exp,omitempty"`               //所需经验值
	AccountLevelLimit int32   `json:"AccountLevelLimit,omitempty"` //账号等级限制
}

// CrystalLevelup.xlsx属性表集合
type CrystalLevelupEntries struct {
	Rows map[int32]*CrystalLevelupEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("CrystalLevelup.xlsx", (*CrystalLevelupEntries)(nil))
}

func (e *CrystalLevelupEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	crystalLevelupEntries = &CrystalLevelupEntries{
		Rows: make(map[int32]*CrystalLevelupEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CrystalLevelupEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		crystalLevelupEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetCrystalLevelupEntry(id int32) (*CrystalLevelupEntry, bool) {
	entry, ok := crystalLevelupEntries.Rows[id]
	return entry, ok
}

func GetCrystalLevelupSize() int32 {
	return int32(len(crystalLevelupEntries.Rows))
}

func GetCrystalLevelupRows() map[int32]*CrystalLevelupEntry {
	return crystalLevelupEntries.Rows
}
