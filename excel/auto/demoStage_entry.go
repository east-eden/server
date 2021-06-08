package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var demoStageEntries *DemoStageEntries //DemoStage.csv全局变量

// DemoStage.csv属性表
type DemoStageEntry struct {
	Id         int32           `json:"Id,omitempty"`         // 主键
	Camp       int32           `json:"Camp,omitempty"`       //阵营
	X          int32           `json:"X,omitempty"`          //坐标
	Z          int32           `json:"Z,omitempty"`          //坐标
	HeroID     int32           `json:"HeroID,omitempty"`     //单位ID
	InitialCom decimal.Decimal `json:"InitialCom,omitempty"` //初始行动条
}

// DemoStage.csv属性表集合
type DemoStageEntries struct {
	Rows map[int32]*DemoStageEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("DemoStage.csv", (*DemoStageEntries)(nil))
}

func (e *DemoStageEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	demoStageEntries = &DemoStageEntries{
		Rows: make(map[int32]*DemoStageEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &DemoStageEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		demoStageEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetDemoStageEntry(id int32) (*DemoStageEntry, bool) {
	entry, ok := demoStageEntries.Rows[id]
	return entry, ok
}

func GetDemoStageSize() int32 {
	return int32(len(demoStageEntries.Rows))
}

func GetDemoStageRows() map[int32]*DemoStageEntry {
	return demoStageEntries.Rows
}
