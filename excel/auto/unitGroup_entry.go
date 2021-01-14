package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	unitGroupEntries	*UnitGroupEntries	//unitGroup.xlsx全局变量

// unitGroup.xlsx属性表
type UnitGroupEntry struct {
	Id        	int32     	`json:"Id,omitempty"`	//怪物组id     
	Name      	string    	`json:"Name,omitempty"`	//怪物组名      
	UnitTypeId	[]int32   	`json:"UnitTypeId,omitempty"`	//怪物id      
	Position  	string    	`json:"Position,omitempty"`	//位置坐标      
}

// unitGroup.xlsx属性表集合
type UnitGroupEntries struct {
	Rows      	map[int32]*UnitGroupEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("unitGroup.xlsx", unitGroupEntries)
}

func (e *UnitGroupEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	unitGroupEntries = &UnitGroupEntries{
		Rows: make(map[int32]*UnitGroupEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &UnitGroupEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	unitGroupEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetUnitGroupEntry(id int32) (*UnitGroupEntry, bool) {
	entry, ok := unitGroupEntries.Rows[id]
	return entry, ok
}

func  GetUnitGroupSize() int32 {
	return int32(len(unitGroupEntries.Rows))
}

