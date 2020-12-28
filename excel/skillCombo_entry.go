package excel

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
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
	AddEntries("skillCombo.xlsx", heroEntries)
}

func (e *SkillComboEntries) load(excelFileRaw *ExcelFileRaw) error {
	
	skillComboEntries = &SkillComboEntries{
		Rows: make(map[int]*SkillComboEntry),
	}

	for _, v := range excelFileRaw.cellData {
		entry := &SkillComboEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	skillComboEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.filename).Msg("excel load success")
	return nil
	
}

func  GetSkillComboEntry(id int) (*SkillComboEntry, bool) {
	entry, ok := skillComboEntries.Rows[id]
	return entry, ok
}

