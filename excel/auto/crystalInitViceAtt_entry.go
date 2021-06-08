package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var crystalInitViceAttEntries *CrystalInitViceAttEntries //CrystalInitViceAtt.csv全局变量

// CrystalInitViceAtt.csv属性表
type CrystalInitViceAttEntry struct {
	Id           int32   `json:"Id,omitempty"`           // 主键
	AttNumWeight []int32 `json:"AttNumWeight,omitempty"` //随机多条副属性权重
}

// CrystalInitViceAtt.csv属性表集合
type CrystalInitViceAttEntries struct {
	Rows map[int32]*CrystalInitViceAttEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("CrystalInitViceAtt.csv", (*CrystalInitViceAttEntries)(nil))
}

func (e *CrystalInitViceAttEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	crystalInitViceAttEntries = &CrystalInitViceAttEntries{
		Rows: make(map[int32]*CrystalInitViceAttEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CrystalInitViceAttEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		crystalInitViceAttEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetCrystalInitViceAttEntry(id int32) (*CrystalInitViceAttEntry, bool) {
	entry, ok := crystalInitViceAttEntries.Rows[id]
	return entry, ok
}

func GetCrystalInitViceAttSize() int32 {
	return int32(len(crystalInitViceAttEntries.Rows))
}

func GetCrystalInitViceAttRows() map[int32]*CrystalInitViceAttEntry {
	return crystalInitViceAttEntries.Rows
}
