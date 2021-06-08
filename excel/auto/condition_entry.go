package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var conditionEntries *ConditionEntries //Condition.csv全局变量

// Condition.csv属性表
type ConditionEntry struct {
	Id        int32   `json:"Id,omitempty"`        // 主键
	Type      int32   `json:"Type,omitempty"`      //条件类型
	SubTypes  []int32 `json:"SubTypes,omitempty"`  //解锁子条件类型
	SubValues []int32 `json:"SubValues,omitempty"` //解锁子条件数值
}

// Condition.csv属性表集合
type ConditionEntries struct {
	Rows map[int32]*ConditionEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Condition.csv", (*ConditionEntries)(nil))
}

func (e *ConditionEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	conditionEntries = &ConditionEntries{
		Rows: make(map[int32]*ConditionEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &ConditionEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		conditionEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetConditionEntry(id int32) (*ConditionEntry, bool) {
	entry, ok := conditionEntries.Rows[id]
	return entry, ok
}

func GetConditionSize() int32 {
	return int32(len(conditionEntries.Rows))
}

func GetConditionRows() map[int32]*ConditionEntry {
	return conditionEntries.Rows
}
