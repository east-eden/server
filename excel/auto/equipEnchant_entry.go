package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var equipEnchantEntries *EquipEnchantEntries //EquipEnchant.xlsx全局变量

// EquipEnchant.xlsx属性表
type EquipEnchantEntry struct {
	Id             int32   `json:"Id,omitempty"`             // 主键
	AttId          int32   `json:"AttId,omitempty"`          //基础属性id
	EquipPos       int32   `json:"EquipPos,omitempty"`       //装备位置
	PromoteAttId   []int32 `json:"PromoteAttId,omitempty"`   //装备突破属性
	PromoteSpellId []int32 `json:"PromoteSpellId,omitempty"` //装备突破被动技能
	PromoteCostId  []int32 `json:"PromoteCostId,omitempty"`  //装备突破消耗
}

// EquipEnchant.xlsx属性表集合
type EquipEnchantEntries struct {
	Rows map[int32]*EquipEnchantEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("EquipEnchant.xlsx", (*EquipEnchantEntries)(nil))
}

func (e *EquipEnchantEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	equipEnchantEntries = &EquipEnchantEntries{
		Rows: make(map[int32]*EquipEnchantEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &EquipEnchantEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		equipEnchantEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetEquipEnchantEntry(id int32) (*EquipEnchantEntry, bool) {
	entry, ok := equipEnchantEntries.Rows[id]
	return entry, ok
}

func GetEquipEnchantSize() int32 {
	return int32(len(equipEnchantEntries.Rows))
}

func GetEquipEnchantRows() map[int32]*EquipEnchantEntry {
	return equipEnchantEntries.Rows
}
