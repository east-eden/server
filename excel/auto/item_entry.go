package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var itemEntries *ItemEntries //Item.csv全局变量

// Item.csv属性表
type ItemEntry struct {
	Id                 int32   `json:"Id,omitempty"`                 // 主键
	Type               int32   `json:"Type,omitempty"`               //物品类型
	SubType            int32   `json:"SubType,omitempty"`            //物品子类型
	Quality            int32   `json:"Quality,omitempty"`            //品质
	MaxStack           int32   `json:"MaxStack,omitempty"`           //最大堆叠数
	TimeLife           int32   `json:"TimeLife,omitempty"`           //时限（分钟）
	TimeStartLifeStamp int32   `json:"TimeStartLifeStamp,omitempty"` //时限开始时间（unix时间戳）
	CanSell            bool    `json:"CanSell,omitempty"`            //是否可以出售
	SellType           int32   `json:"SellType,omitempty"`           //出售货币类型
	SellPrice          int32   `json:"SellPrice,omitempty"`          //出售价格
	StaleGainId        int32   `json:"StaleGainId,omitempty"`        //过期后转换的掉落id
	EffectType         int32   `json:"EffectType,omitempty"`         //使用效果类型
	EffectValue        int32   `json:"EffectValue,omitempty"`        //使用效果参数
	PublicMisc         []int32 `json:"PublicMisc,omitempty"`         //物品公用参数
}

// Item.csv属性表集合
type ItemEntries struct {
	Rows map[int32]*ItemEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Item.csv", (*ItemEntries)(nil))
}

func (e *ItemEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	itemEntries = &ItemEntries{
		Rows: make(map[int32]*ItemEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &ItemEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		itemEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetItemEntry(id int32) (*ItemEntry, bool) {
	entry, ok := itemEntries.Rows[id]
	return entry, ok
}

func GetItemSize() int32 {
	return int32(len(itemEntries.Rows))
}

func GetItemRows() map[int32]*ItemEntry {
	return itemEntries.Rows
}
