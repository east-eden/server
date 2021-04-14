package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	unitGroupEntries	*UnitGroupEntries	//UnitGroup.xlsx全局变量

// UnitGroup.xlsx属性表
type UnitGroupEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	Name           	string              	`json:"Name,omitempty"`	//怪物组名字     
	HeroIds        	[]int32             	`json:"HeroIds,omitempty"`	//怪物id      
	Camps          	[]int32             	`json:"Camps,omitempty"`	//所属阵营      
	PosXs          	[]int32             	`json:"PosXs,omitempty"`	//x坐标       
	PosZs          	[]int32             	`json:"PosZs,omitempty"`	//z坐标       
	InitComs       	[]number            	`json:"InitComs,omitempty"`	//初始行动条     
}

// UnitGroup.xlsx属性表集合
type UnitGroupEntries struct {
	Rows           	map[int32]*UnitGroupEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("UnitGroup.xlsx", (*UnitGroupEntries)(nil))
}

func (e *UnitGroupEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	unitGroupEntries = &UnitGroupEntries{
		Rows: make(map[int32]*UnitGroupEntry, 100),
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

func  GetUnitGroupRows() map[int32]*UnitGroupEntry {
	return unitGroupEntries.Rows
}

