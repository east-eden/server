package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var buyStrengthenEntries *BuyStrengthenEntries //BuyStrengthen.xlsx全局变量

// BuyStrengthen.xlsx属性表
type BuyStrengthenEntry struct {
	Id          int32 `json:"Id,omitempty"`          // 主键
	ConditionId int32 `json:"ConditionId,omitempty"` //限制条件id
	Cost        int32 `json:"Cost,omitempty"`        //花费钻石
	Strengthen  int32 `json:"Strengthen,omitempty"`  //获得体力
}

// BuyStrengthen.xlsx属性表集合
type BuyStrengthenEntries struct {
	Rows map[int32]*BuyStrengthenEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("BuyStrengthen.xlsx", (*BuyStrengthenEntries)(nil))
}

func (e *BuyStrengthenEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	buyStrengthenEntries = &BuyStrengthenEntries{
		Rows: make(map[int32]*BuyStrengthenEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &BuyStrengthenEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		buyStrengthenEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetBuyStrengthenEntry(id int32) (*BuyStrengthenEntry, bool) {
	entry, ok := buyStrengthenEntries.Rows[id]
	return entry, ok
}

func GetBuyStrengthenSize() int32 {
	return int32(len(buyStrengthenEntries.Rows))
}

func GetBuyStrengthenRows() map[int32]*BuyStrengthenEntry {
	return buyStrengthenEntries.Rows
}
