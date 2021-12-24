package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var raceEntries *RaceEntries //Race.csv全局变量

// Race.csv属性表
type RaceEntry struct {
	Id       int32  `json:"Id,omitempty"`       // 主键
	Racename string `json:"Racename,omitempty"` //种族名称
	Racedesc string `json:"Racedesc,omitempty"` //种族简介
	Icon     string `json:"Icon,omitempty"`     //模型战斗头像
}

// Race.csv属性表集合
type RaceEntries struct {
	Rows map[int32]*RaceEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Race.csv", (*RaceEntries)(nil))
}

func (e *RaceEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	raceEntries = &RaceEntries{
		Rows: make(map[int32]*RaceEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &RaceEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		raceEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetRaceEntry(id int32) (*RaceEntry, bool) {
	entry, ok := raceEntries.Rows[id]
	return entry, ok
}

func GetRaceSize() int32 {
	return int32(len(raceEntries.Rows))
}

func GetRaceRows() map[int32]*RaceEntry {
	return raceEntries.Rows
}
