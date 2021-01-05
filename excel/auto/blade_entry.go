package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	bladeEntries	*BladeEntries	//blade.xlsx全局变量

// blade.xlsx属性表
type BladeEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名字        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	Star      	int       	`json:"Star,omitempty"`	//星级        
	Level     	int       	`json:"Level,omitempty"`	//等级        
	NextLevel 	int       	`json:"NextLevel,omitempty"`	//下个等级      
	CD        	float32   	`json:"CD,omitempty"`	//切换冷却时间    
	Class     	int       	`json:"Class,omitempty"`	//异刃属性      
	AttrName  	[]string  	`json:"AttrName,omitempty"`	//属性名       
	AttrValue 	[]int     	`json:"AttrValue,omitempty"`	//属性值       
	AttackType	int       	`json:"AttackType,omitempty"`	//攻击方式      
	Skill     	[]int     	`json:"Skill,omitempty"`	//拥有技能      
}

// blade.xlsx属性表集合
type BladeEntries struct {
	Rows      	map[int]*BladeEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("blade.xlsx", bladeEntries)
}

func (e *BladeEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	bladeEntries = &BladeEntries{
		Rows: make(map[int]*BladeEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &BladeEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	bladeEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetBladeEntry(id int) (*BladeEntry, bool) {
	entry, ok := bladeEntries.Rows[id]
	return entry, ok
}

func  GetBladeSize() int {
	return len(bladeEntries.Rows)
}

