package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	unitGroupEntries	*UnitGroupEntries	//unitGroup.xlsx全局变量

// unitGroup.xlsx属性表
type UnitGroupEntry struct {
	Id        	int32     	`json:"Id,omitempty"`	//怪物组id     
	Name      	string    	`json:"Name,omitempty"`	//怪物组名      
	UnitTypeId	[]int32   	`json:"UnitTypeId,omitempty"`	//怪物id      
	Position  	interface{}	`json:"Position,omitempty"`	//位置坐标      
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
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
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

