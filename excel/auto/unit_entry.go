package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	unitEntries    	*UnitEntries   	//Unit.xlsx全局变量  

// Unit.xlsx属性表
type UnitEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	AttId          	int32               	`json:"AttId,omitempty"`	//属性id      
}

// Unit.xlsx属性表集合
type UnitEntries struct {
	Rows           	map[int32]*UnitEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("Unit.xlsx", (*UnitEntries)(nil))
}

func (e *UnitEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	unitEntries = &UnitEntries{
		Rows: make(map[int32]*UnitEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &UnitEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	unitEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetUnitEntry(id int32) (*UnitEntry, bool) {
	entry, ok := unitEntries.Rows[id]
	return entry, ok
}

func  GetUnitSize() int32 {
	return int32(len(unitEntries.Rows))
}

func  GetUnitRows() map[int32]*UnitEntry {
	return unitEntries.Rows
}

