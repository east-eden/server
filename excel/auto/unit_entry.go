package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	unitEntries	*UnitEntries	//unit.xlsx全局变量

// unit.xlsx属性表
type UnitEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名字        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	Level     	int       	`json:"Level,omitempty"`	//等级        
	NextLevel 	int       	`json:"NextLevel,omitempty"`	//下个等级      
	AttrName  	[]string  	`json:"AttrName,omitempty"`	//属性名       
	AttrValue 	[]int     	`json:"AttrValue,omitempty"`	//属性值       
	Resource  	string    	`json:"Resource,omitempty"`	//资源路径      
}

// unit.xlsx属性表集合
type UnitEntries struct {
	Rows      	map[int]*UnitEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("unit.xlsx", unitEntries)
}

func (e *UnitEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	unitEntries = &UnitEntries{
		Rows: make(map[int]*UnitEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &UnitEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	unitEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetUnitEntry(id int) (*UnitEntry, bool) {
	entry, ok := unitEntries.Rows[id]
	return entry, ok
}

func  GetUnitSize() int {
	return len(unitEntries.Rows)
}

