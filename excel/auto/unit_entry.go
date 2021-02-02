package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	unitEntries    	*UnitEntries   	//Unit.xlsx全局变量  

// Unit.xlsx属性表
type UnitEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	Name           	string              	`json:"Name,omitempty"`	//名字        
	Desc           	string              	`json:"Desc,omitempty"`	//描述        
	Level          	int32               	`json:"Level,omitempty"`	//等级        
	NextLevel      	int32               	`json:"NextLevel,omitempty"`	//下个等级      
	AttrName       	[]string            	`json:"AttrName,omitempty"`	//属性名       
	AttrValue      	[]int32             	`json:"AttrValue,omitempty"`	//属性值       
	Resource       	string              	`json:"Resource,omitempty"`	//资源路径      
	NormalSpellId  	int32               	`json:"NormalSpellId,omitempty"`	//普攻技能id    
	SpecialSpellId 	int32               	`json:"SpecialSpellId,omitempty"`	//特殊技能id    
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

