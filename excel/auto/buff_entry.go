package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	buffEntries	*BuffEntries	//buff.xlsx全局变量

// buff.xlsx属性表
type BuffEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	BuffType  	int       	`json:"BuffType,omitempty"`	//buff类型    


	Level     	int       	`json:"Level,omitempty"`	//等级        
	NextLevel 	int       	`json:"NextLevel,omitempty"`	//下个等级      
	CD        	float32   	`json:"CD,omitempty"`	//冷却时间(秒)   
	LifeTime  	float32   	`json:"LifeTime,omitempty"`	//持续时间(秒)   
	BuffOverlap	int       	`json:"BuffOverlap,omitempty"`	//叠加类型      
	MaxLimit  	int       	`json:"MaxLimit,omitempty"`	//限制        
	Params_StrValue	[]string  	`json:"Params_StrValue,omitempty"`	//参数列表，目标属性 
	Params_Formula	[]string  	`json:"Params_Formula,omitempty"`	//公式        
	Params_NumValue	[]int     	`json:"Params_NumValue,omitempty"`	//参数列表，固定数值 
	Effect    	string    	`json:"Effect,omitempty"`	//显示特效      
}

// buff.xlsx属性表集合
type BuffEntries struct {
	Rows      	map[int]*BuffEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("buff.xlsx", buffEntries)
}

func (e *BuffEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	buffEntries = &BuffEntries{
		Rows: make(map[int]*BuffEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &BuffEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	buffEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetBuffEntry(id int) (*BuffEntry, bool) {
	entry, ok := buffEntries.Rows[id]
	return entry, ok
}

func  GetBuffSize() int {
	return len(buffEntries.Rows)
}

