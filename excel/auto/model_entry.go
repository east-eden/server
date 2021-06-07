package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var modelEntries *ModelEntries //Model.xlsx全局变量

// Model.xlsx属性表
type ModelEntry struct {
	Id            int32           `json:"Id,omitempty"`            // 主键
	ModelResource string          `json:"ModelResource,omitempty"` //模型资源ID
	Model_scale   string          `json:"Model_scale,omitempty"`   //怪物简介
	Modelscope    decimal.Decimal `json:"Modelscope,omitempty"`    //模型半径
	Icon          string          `json:"Icon,omitempty"`          //模型战斗头像
}

// Model.xlsx属性表集合
type ModelEntries struct {
	Rows map[int32]*ModelEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Model.xlsx", (*ModelEntries)(nil))
}

func (e *ModelEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	modelEntries = &ModelEntries{
		Rows: make(map[int32]*ModelEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &ModelEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		modelEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetModelEntry(id int32) (*ModelEntry, bool) {
	entry, ok := modelEntries.Rows[id]
	return entry, ok
}

func GetModelSize() int32 {
	return int32(len(modelEntries.Rows))
}

func GetModelRows() map[int32]*ModelEntry {
	return modelEntries.Rows
}
