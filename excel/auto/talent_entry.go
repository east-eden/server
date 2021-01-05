package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	talentEntries	*TalentEntries	//talent.xlsx全局变量

// talent.xlsx属性表
type TalentEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//天赋名称      
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	LevelLimit	int       	`json:"LevelLimit,omitempty"`	//等级限制      
	GroupId   	int       	`json:"GroupId,omitempty"`	//天赋组id     
	CostId    	int       	`json:"CostId,omitempty"`	//消耗id      
}

// talent.xlsx属性表集合
type TalentEntries struct {
	Rows      	map[int]*TalentEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("talent.xlsx", talentEntries)
}

func (e *TalentEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	talentEntries = &TalentEntries{
		Rows: make(map[int]*TalentEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &TalentEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	talentEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetTalentEntry(id int) (*TalentEntry, bool) {
	entry, ok := talentEntries.Rows[id]
	return entry, ok
}

func  GetTalentSize() int {
	return len(talentEntries.Rows)
}

