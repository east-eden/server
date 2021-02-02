package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	talentEntries  	*TalentEntries 	//Talent.xlsx全局变量

// Talent.xlsx属性表
type TalentEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	Name           	string              	`json:"Name,omitempty"`	//天赋名称      
	Desc           	string              	`json:"Desc,omitempty"`	//描述        
	LevelLimit     	int32               	`json:"LevelLimit,omitempty"`	//等级限制      
	GroupId        	int32               	`json:"GroupId,omitempty"`	//天赋组id     
	CostId         	int32               	`json:"CostId,omitempty"`	//消耗id      
}

// Talent.xlsx属性表集合
type TalentEntries struct {
	Rows           	map[int32]*TalentEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("Talent.xlsx", (*TalentEntries)(nil))
}

func (e *TalentEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	talentEntries = &TalentEntries{
		Rows: make(map[int32]*TalentEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &TalentEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	talentEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetTalentEntry(id int32) (*TalentEntry, bool) {
	entry, ok := talentEntries.Rows[id]
	return entry, ok
}

func  GetTalentSize() int32 {
	return int32(len(talentEntries.Rows))
}

func  GetTalentRows() map[int32]*TalentEntry {
	return talentEntries.Rows
}

