package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	costLootEntries	*CostLootEntries	//costLoot.xlsx全局变量

// costLoot.xlsx属性表
type CostLootEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	Type      	int       	`json:"Type,omitempty"`	//类型        
	Misc      	int       	`json:"Misc,omitempty"`	//参数        
	Num       	int       	`json:"Num,omitempty"`	//数量        
}

// costLoot.xlsx属性表集合
type CostLootEntries struct {
	Rows      	map[int]*CostLootEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("costLoot.xlsx", costLootEntries)
}

func (e *CostLootEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	costLootEntries = &CostLootEntries{
		Rows: make(map[int]*CostLootEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CostLootEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	costLootEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetCostLootEntry(id int) (*CostLootEntry, bool) {
	entry, ok := costLootEntries.Rows[id]
	return entry, ok
}

func  GetCostLootSize() int {
	return len(costLootEntries.Rows)
}

