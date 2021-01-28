package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	bladeEntries   	*BladeEntries  	//blade.xlsx全局变量 

// blade.xlsx属性表
type BladeEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	//id        
	Name           	string              	`json:"Name,omitempty"`	//名字        
	Desc           	string              	`json:"Desc,omitempty"`	//描述        
	Star           	int32               	`json:"Star,omitempty"`	//星级        
	Level          	int32               	`json:"Level,omitempty"`	//等级        
	NextLevel      	int32               	`json:"NextLevel,omitempty"`	//下个等级      
	CD             	float32             	`json:"CD,omitempty"`	//切换冷却时间    
	Class          	int32               	`json:"Class,omitempty"`	//异刃属性      
	AttrName       	[]string            	`json:"AttrName,omitempty"`	//属性名       
	AttrValue      	[]int32             	`json:"AttrValue,omitempty"`	//属性值       
	AttackType     	int32               	`json:"AttackType,omitempty"`	//攻击方式      
	Skill          	[]int32             	`json:"Skill,omitempty"`	//拥有技能      
}

// blade.xlsx属性表集合
type BladeEntries struct {
	Rows           	map[int32]*BladeEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("blade.xlsx", bladeEntries)
}

func (e *BladeEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	bladeEntries = &BladeEntries{
		Rows: make(map[int32]*BladeEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &BladeEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	bladeEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetBladeEntry(id int32) (*BladeEntry, bool) {
	entry, ok := bladeEntries.Rows[id]
	return entry, ok
}

func  GetBladeSize() int32 {
	return int32(len(bladeEntries.Rows))
}

