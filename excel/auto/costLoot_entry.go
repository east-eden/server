package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var costLootEntries *CostLootEntries //CostLoot.xlsx全局变量

// CostLoot.xlsx属性表
type CostLootEntry struct {
	Id   int32   `json:"Id,omitempty"`   // 主键
	Type []int32 `json:"Type,omitempty"` //消耗和掉落类型
	Misc []int32 `json:"Misc,omitempty"` //消耗和掉落参数
	Num  []int32 `json:"Num,omitempty"`  //数量
}

// CostLoot.xlsx属性表集合
type CostLootEntries struct {
	Rows map[int32]*CostLootEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("CostLoot.xlsx", (*CostLootEntries)(nil))
}

func (e *CostLootEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	costLootEntries = &CostLootEntries{
		Rows: make(map[int32]*CostLootEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CostLootEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		costLootEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetCostLootEntry(id int32) (*CostLootEntry, bool) {
	entry, ok := costLootEntries.Rows[id]
	return entry, ok
}

func GetCostLootSize() int32 {
	return int32(len(costLootEntries.Rows))
}

func GetCostLootRows() map[int32]*CostLootEntry {
	return costLootEntries.Rows
}
