package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	attEntries     	*AttEntries    	//Att.xlsx全局变量   

// Att.xlsx属性表
type AttEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	Atk            	int32               	`json:"Atk,omitempty"`	//攻击力       
	AtkPercent     	int32               	`json:"AtkPercent,omitempty"`	//攻击力百分比    
	Armor          	int32               	`json:"Armor,omitempty"`	//护甲        
	ArmorPercent   	int32               	`json:"ArmorPercent,omitempty"`	//护甲百分比     
	DmgInc         	int32               	`json:"DmgInc,omitempty"`	//总伤害加成     
	Crit           	int32               	`json:"Crit,omitempty"`	//暴击值       
	CritInc        	int32               	`json:"CritInc,omitempty"`	//暴击倍数加成    
	Heal           	int32               	`json:"Heal,omitempty"`	//治疗强度      
	HealPercent    	int32               	`json:"HealPercent,omitempty"`	//治疗强度百分比   
	RealDmg        	int32               	`json:"RealDmg,omitempty"`	//真实伤害      
	MoveSpeed      	int32               	`json:"MoveSpeed,omitempty"`	//战场移动速度    
	MoveSpeedPercent	int32               	`json:"MoveSpeedPercent,omitempty"`	//战场移动速度百分比 
	AtbSpeed       	int32               	`json:"AtbSpeed,omitempty"`	//时间槽速度     
	AtbSpeedPercent	int32               	`json:"AtbSpeedPercent,omitempty"`	//时间槽速度百分比  
	EffectHit      	int32               	`json:"EffectHit,omitempty"`	//技能效果命中    
	EffectResist   	int32               	`json:"EffectResist,omitempty"`	//技能效果抵抗    
	MaxHP          	int32               	`json:"MaxHP,omitempty"`	//血量上限      
	MaxHPPercent   	int32               	`json:"MaxHPPercent,omitempty"`	//血量上限百分比   
	MaxMP          	int32               	`json:"MaxMP,omitempty"`	//蓝量上限      
	GenMP          	int32               	`json:"GenMP,omitempty"`	//魔法恢复      
	Rage           	int32               	`json:"Rage,omitempty"`	//怒气        
	Hit            	int32               	`json:"Hit,omitempty"`	//命中值       
	Dodge          	int32               	`json:"Dodge,omitempty"`	//闪避值       
	Movedistance   	int32               	`json:"Movedistance,omitempty"`	//移动距离      
	DmgOfType      	[]int32             	`json:"DmgOfType,omitempty"`	//各系伤害加层    
	ResOfType      	[]int32             	`json:"ResOfType,omitempty"`	//各系伤害减免    
}

// Att.xlsx属性表集合
type AttEntries struct {
	Rows           	map[int32]*AttEntry 	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("Att.xlsx", (*AttEntries)(nil))
}

func (e *AttEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	attEntries = &AttEntries{
		Rows: make(map[int32]*AttEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &AttEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	attEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetAttEntry(id int32) (*AttEntry, bool) {
	entry, ok := attEntries.Rows[id]
	return entry, ok
}

func  GetAttSize() int32 {
	return int32(len(attEntries.Rows))
}

func  GetAttRows() map[int32]*AttEntry {
	return attEntries.Rows
}

