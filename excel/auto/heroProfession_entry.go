package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroProfessionEntries *HeroProfessionEntries //HeroProfession.xlsx全局变量

// HeroProfession.xlsx属性表
type HeroProfessionEntry struct {
	Id         int32 `json:"Id,omitempty"`         // 主键
	AtkRatio   int32 `json:"AtkRatio,omitempty"`   //攻击系数
	ArmorRatio int32 `json:"ArmorRatio,omitempty"` //护甲系数
	MaxHPRatio int32 `json:"MaxHPRatio,omitempty"` //血量系数
}

// HeroProfession.xlsx属性表集合
type HeroProfessionEntries struct {
	Rows map[int32]*HeroProfessionEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("HeroProfession.xlsx", (*HeroProfessionEntries)(nil))
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
