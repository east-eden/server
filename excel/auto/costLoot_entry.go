package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var costLootEntries *CostLootEntries //CostLoot.xlsx全局变量

// CostLoot.xlsx属性表
type CostLootEntry struct {
	Id        int32   `json:"Id,omitempty"`        // 主键
	LootKind  int32   `json:"LootKind,omitempty"`  //掉落种类
	LootTimes int32   `json:"LootTimes,omitempty"` //掉落次数
	CanRepeat bool    `json:"CanRepeat,omitempty"` //是否可重复掉落
	Type      []int32 `json:"Type,omitempty"`      //消耗和掉落类型
	Misc      []int32 `json:"Misc,omitempty"`      //消耗和掉落参数
	NumMin    []int32 `json:"NumMin,omitempty"`    //最小数量
	NumMax    []int32 `json:"NumMax,omitempty"`    //最大数量
	Prob      []int32 `json:"Prob,omitempty"`      //概率、权重
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
