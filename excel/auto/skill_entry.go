package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	skillEntries	*SkillEntries	//skill.xlsx全局变量

// skill.xlsx属性表
type SkillEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名字        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	Icon      	string    	`json:"Icon,omitempty"`	//图标        
	NextLevel 	int       	`json:"NextLevel,omitempty"`	//下个等级      
	CD        	float32   	`json:"CD,omitempty"`	//技能CD      
	SkillCombo	int       	`json:"SkillCombo,omitempty"`	//连击类型      
	SkillType 	int       	`json:"SkillType,omitempty"`	//技能类型      
	Precondition	int       	`json:"Precondition,omitempty"`	//前置条件      
	TargetType	int       	`json:"TargetType,omitempty"`	//释放目标类型    
	TargetRule	[]int     	`json:"TargetRule,omitempty"`	//释放目标规则    
	AttackType	int       	`json:"AttackType,omitempty"`	//攻击方式      
	SkillBlocks	[]int     	`json:"SkillBlocks,omitempty"`	//包含的技能块    
	EffectType	int       	`json:"EffectType,omitempty"`	//特效类型      
	Resource  	string    	`json:"Resource,omitempty"`	//特效路径      
}

// skill.xlsx属性表集合
type SkillEntries struct {
	Rows      	map[int]*SkillEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("skill.xlsx", skillEntries)
}

func (e *SkillEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	skillEntries = &SkillEntries{
		Rows: make(map[int]*SkillEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	skillEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetSkillEntry(id int) (*SkillEntry, bool) {
	entry, ok := skillEntries.Rows[id]
	return entry, ok
}

func  GetSkillSize() int {
	return len(skillEntries.Rows)
}

