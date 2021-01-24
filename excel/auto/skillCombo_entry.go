package auto

import (
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"bitbucket.org/east-eden/server/excel"
)

var	skillComboEntries	*SkillComboEntries	//skillCombo.xlsx全局变量

// skillCombo.xlsx属性表
type SkillComboEntry struct {
	Id        	int32     	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名字        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	ClassSequence	[]int32   	`json:"ClassSequence,omitempty"`	//连击属性序列    
	SkillSequence	[]int32   	`json:"SkillSequence,omitempty"`	//连击技能序列    
	Fomula    	int32     	`json:"Fomula,omitempty"`	//伤害公式      
	Buffs     	[]int32   	`json:"Buffs,omitempty"`	//添加Buff    
}

// skillCombo.xlsx属性表集合
type SkillComboEntries struct {
	Rows      	map[int32]*SkillComboEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("skillCombo.xlsx", skillComboEntries)
}

func (e *SkillComboEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	skillComboEntries = &SkillComboEntries{
		Rows: make(map[int32]*SkillComboEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillComboEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	skillComboEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetSkillComboEntry(id int32) (*SkillComboEntry, bool) {
	entry, ok := skillComboEntries.Rows[id]
	return entry, ok
}

func  GetSkillComboSize() int32 {
	return int32(len(skillComboEntries.Rows))
}

