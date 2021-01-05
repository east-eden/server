package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	skillComboEntries	*SkillComboEntries	//skillCombo.xlsx全局变量

// skillCombo.xlsx属性表
type SkillComboEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名字        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	ClassSequence	[]int     	`json:"ClassSequence,omitempty"`	//连击属性序列    
	SkillSequence	[]int     	`json:"SkillSequence,omitempty"`	//连击技能序列    
	Fomula    	int       	`json:"Fomula,omitempty"`	//伤害公式      
	Buffs     	[]int     	`json:"Buffs,omitempty"`	//添加Buff    
}

// skillCombo.xlsx属性表集合
type SkillComboEntries struct {
	Rows      	map[int]*SkillComboEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("skillCombo.xlsx", skillComboEntries)
}

func (e *SkillComboEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	skillComboEntries = &SkillComboEntries{
		Rows: make(map[int]*SkillComboEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillComboEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	skillComboEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetSkillComboEntry(id int) (*SkillComboEntry, bool) {
	entry, ok := skillComboEntries.Rows[id]
	return entry, ok
}

func  GetSkillComboSize() int {
	return len(skillComboEntries.Rows)
}

