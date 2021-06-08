package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var heroProfessionEntries *HeroProfessionEntries //HeroProfession.csv全局变量

// HeroProfession.csv属性表
type HeroProfessionEntry struct {
	Id         int32           `json:"Id,omitempty"`         // 主键
	AtkRatio   decimal.Decimal `json:"AtkRatio,omitempty"`   //攻击系数
	ArmorRatio decimal.Decimal `json:"ArmorRatio,omitempty"` //护甲系数
	MaxHPRatio decimal.Decimal `json:"MaxHPRatio,omitempty"` //血量系数
}

// HeroProfession.csv属性表集合
type HeroProfessionEntries struct {
	Rows map[int32]*HeroProfessionEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("HeroProfession.csv", (*HeroProfessionEntries)(nil))
}

func (e *HeroProfessionEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	heroProfessionEntries = &HeroProfessionEntries{
		Rows: make(map[int32]*HeroProfessionEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &HeroProfessionEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		heroProfessionEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetHeroProfessionEntry(id int32) (*HeroProfessionEntry, bool) {
	entry, ok := heroProfessionEntries.Rows[id]
	return entry, ok
}

func GetHeroProfessionSize() int32 {
	return int32(len(heroProfessionEntries.Rows))
}

func GetHeroProfessionRows() map[int32]*HeroProfessionEntry {
	return heroProfessionEntries.Rows
}
