package auto

import (
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"bitbucket.org/east-eden/server/excel"
)

var	attEntries	*AttEntries	//att.xlsx全局变量

// att.xlsx属性表
type AttEntry struct {
	Id        	int32     	`json:"Id,omitempty"`	//id        

	Atk       	int32     	`json:"Atk,omitempty"`	//攻击力       
	Armor     	int32     	`json:"Armor,omitempty"`	//护甲        
	DmgInc    	int32     	`json:"DmgInc,omitempty"`	//总伤害加成     
	Crit      	int32     	`json:"Crit,omitempty"`	//暴击值       
	CritInc   	int32     	`json:"CritInc,omitempty"`	//暴击倍数加成    
	Heal      	int32     	`json:"Heal,omitempty"`	//治疗        
	RealDmg   	int32     	`json:"RealDmg,omitempty"`	//真实伤害      
	MoveSpeed 	int32     	`json:"MoveSpeed,omitempty"`	//战场移动速度    
	AtbSpeed  	int32     	`json:"AtbSpeed,omitempty"`	//时间槽速度     
	EffectHit 	int32     	`json:"EffectHit,omitempty"`	//技能效果命中    
	EffectResist	int32     	`json:"EffectResist,omitempty"`	//技能效果抵抗    
	MaxHP     	int32     	`json:"MaxHP,omitempty"`	//血量上限      
	GenMP     	int32     	`json:"GenMP,omitempty"`	//魔法恢复      
	Rage      	int32     	`json:"Rage,omitempty"`	//怒气        
	DmgTypes  	[]int32   	`json:"DmgTypes,omitempty"`	//各系伤害加成    
	ResTypes  	[]int32   	`json:"ResTypes,omitempty"`	//各系伤害抗性    
}

// att.xlsx属性表集合
type AttEntries struct {
	Rows      	map[int32]*AttEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("att.xlsx", attEntries)
}

func (e *AttEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	attEntries = &AttEntries{
		Rows: make(map[int32]*AttEntry),
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

