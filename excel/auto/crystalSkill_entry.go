package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var crystalSkillEntries *CrystalSkillEntries //CrystalSkill.csv全局变量

// CrystalSkill.csv属性表
type CrystalSkillEntry struct {
	Id      int32   `json:"Id,omitempty"`      // 主键
	SkillId []int32 `json:"SkillId,omitempty"` //元素技能id
}

// CrystalSkill.csv属性表集合
type CrystalSkillEntries struct {
	Rows map[int32]*CrystalSkillEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("CrystalSkill.csv", (*CrystalSkillEntries)(nil))
}

func (e *CrystalSkillEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	crystalSkillEntries = &CrystalSkillEntries{
		Rows: make(map[int32]*CrystalSkillEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CrystalSkillEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		crystalSkillEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetCrystalSkillEntry(id int32) (*CrystalSkillEntry, bool) {
	entry, ok := crystalSkillEntries.Rows[id]
	return entry, ok
}

func GetCrystalSkillSize() int32 {
	return int32(len(crystalSkillEntries.Rows))
}

func GetCrystalSkillRows() map[int32]*CrystalSkillEntry {
	return crystalSkillEntries.Rows
}
