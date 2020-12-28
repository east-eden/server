package excel

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	itemEntries	*ItemEntries	//item.xlsx全局变量

// item.xlsx属性表
type ItemEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名字        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	Icon      	string    	`json:"Icon,omitempty"`	//图标        
	Type      	int       	`json:"Type,omitempty"`	//类型        
	SubType   	int       	`json:"SubType,omitempty"`	//子类型       
	Quality   	int       	`json:"Quality,omitempty"`	//品质        
	MaxStack  	int       	`json:"MaxStack,omitempty"`	//最大堆叠数     
	EffectType	int       	`json:"EffectType,omitempty"`	//使用效果      
	EffectValue	[]int     	`json:"EffectValue,omitempty"`	//效果数值      
	EquipEnchantID	int       	`json:"EquipEnchantID,omitempty"`	//装备强化id    
}

// item.xlsx属性表集合
type ItemEntries struct {
	Rows      	map[int]*ItemEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	AddEntries("item.xlsx", heroEntries)
}

func (e *ItemEntries) load(excelFileRaw *ExcelFileRaw) error {
	
	itemEntries = &ItemEntries{
		Rows: make(map[int]*ItemEntry),
	}

	for _, v := range excelFileRaw.cellData {
		entry := &ItemEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	itemEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.filename).Msg("excel load success")
	return nil
	
}

func  GetItemEntry(id int) (*ItemEntry, bool) {
	entry, ok := itemEntries.Rows[id]
	return entry, ok
}

