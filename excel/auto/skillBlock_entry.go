package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	skillBlockEntries	*SkillBlockEntries	//skillBlock.xlsx全局变量

// skillBlock.xlsx属性表
type SkillBlockEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Index     	int       	`json:"Index,omitempty"`	//技能块的索引    
	Condition 	[]string  	`json:"Condition,omitempty"`	//触发条件      
	Buffs     	[]int     	`json:"Buffs,omitempty"`	//添加Buff    
	Formula   	int       	`json:"Formula,omitempty"`	//伤害公式      
	Ratio     	[]float32 	`json:"Ratio,omitempty"`	//伤害系数      
}

// skillBlock.xlsx属性表集合
type SkillBlockEntries struct {
	Rows      	map[int]*SkillBlockEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("skillBlock.xlsx", skillBlockEntries)
}

func (e *SkillBlockEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	skillBlockEntries = &SkillBlockEntries{
		Rows: make(map[int]*SkillBlockEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillBlockEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	skillBlockEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetSkillBlockEntry(id int) (*SkillBlockEntry, bool) {
	entry, ok := skillBlockEntries.Rows[id]
	return entry, ok
}

func  GetSkillBlockSize() int {
	return len(skillBlockEntries.Rows)
}
