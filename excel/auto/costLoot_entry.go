package auto

import (
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"e.coding.net/mmstudio/blade/server/excel"
)

var	costLootEntries	*CostLootEntries	//costLoot.xlsx全局变量

// costLoot.xlsx属性表
type CostLootEntry struct {
	Id        	int32     	`json:"Id,omitempty"`	//id        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	Type      	int32     	`json:"Type,omitempty"`	//类型        
	Misc      	int32     	`json:"Misc,omitempty"`	//参数        
	Num       	int32     	`json:"Num,omitempty"`	//数量        
}

// costLoot.xlsx属性表集合
type CostLootEntries struct {
	Rows      	map[int32]*CostLootEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("costLoot.xlsx", costLootEntries)
}

func (e *CostLootEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	costLootEntries = &CostLootEntries{
		Rows: make(map[int32]*CostLootEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CostLootEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	costLootEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetCostLootEntry(id int32) (*CostLootEntry, bool) {
	entry, ok := costLootEntries.Rows[id]
	return entry, ok
}

func  GetCostLootSize() int32 {
	return int32(len(costLootEntries.Rows))
}

