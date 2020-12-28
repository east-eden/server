package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	equipEnchantEntries	*EquipEnchantEntries	//equipEnchant.xlsx全局变量

// equipEnchant.xlsx属性表
type EquipEnchantEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	EquipPos  	int       	`json:"EquipPos,omitempty"`	//装备位置      
	AttId     	int       	`json:"AttId,omitempty"`	//属性id      
}

// equipEnchant.xlsx属性表集合
type EquipEnchantEntries struct {
	Rows      	map[int]*EquipEnchantEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("equipEnchant.xlsx", equipEnchantEntries)
}

func (e *EquipEnchantEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	equipEnchantEntries = &EquipEnchantEntries{
		Rows: make(map[int]*EquipEnchantEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &EquipEnchantEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	equipEnchantEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetEquipEnchantEntry(id int) (*EquipEnchantEntry, bool) {
	entry, ok := equipEnchantEntries.Rows[id]
	return entry, ok
}

func  GetEquipEnchantSize() int {
	return len(equipEnchantEntries.Rows)
}

