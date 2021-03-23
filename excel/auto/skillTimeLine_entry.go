package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	skillTimeLineEntries	*SkillTimeLineEntries	//SkillTimeLine.xlsx全局变量

// SkillTimeLine.xlsx属性表
type SkillTimeLineEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	TrackID        	int32               	`json:"TrackID,omitempty"`	//trackID   
	StartTime      	number              	`json:"StartTime,omitempty"`	//开始时间      
	ShowType       	string              	`json:"ShowType,omitempty"`	//track类型   
	Sort           	string              	`json:"Sort,omitempty"`	//sort      
	DurationTime   	number              	`json:"DurationTime,omitempty"`	//持续时间      
	IsBullet       	int32               	`json:"IsBullet,omitempty"`	//是否弹道      
	BulletLocate   	string              	`json:"BulletLocate,omitempty"`	//弹道初始      
	BulletSpeed    	number              	`json:"BulletSpeed,omitempty"`	//子弹速度      
	BeforeSkillAct 	string              	`json:"BeforeSkillAct,omitempty"`	//动作前摇      
	AnimName       	string              	`json:"AnimName,omitempty"`	//动作名称      
	FxName         	string              	`json:"FxName,omitempty"`	//特效名称      
	AfterSkillAct  	string              	`json:"AfterSkillAct,omitempty"`	//动作后摇      
	SlotID         	[]string            	`json:"SlotID,omitempty"`	//插槽ID      
	BeHitPoints    	string              	`json:"BeHitPoints,omitempty"`	//命中点       
	BeAtkShow      	[]string            	`json:"BeAtkShow,omitempty"`	//          
}

// SkillTimeLine.xlsx属性表集合
type SkillTimeLineEntries struct {
	Rows           	map[int32]*SkillTimeLineEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("SkillTimeLine.xlsx", (*SkillTimeLineEntries)(nil))
}

func (e *SkillTimeLineEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	skillTimeLineEntries = &SkillTimeLineEntries{
		Rows: make(map[int32]*SkillTimeLineEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillTimeLineEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	skillTimeLineEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetSkillTimeLineEntry(id int32) (*SkillTimeLineEntry, bool) {
	entry, ok := skillTimeLineEntries.Rows[id]
	return entry, ok
}

func  GetSkillTimeLineSize() int32 {
	return int32(len(skillTimeLineEntries.Rows))
}

func  GetSkillTimeLineRows() map[int32]*SkillTimeLineEntry {
	return skillTimeLineEntries.Rows
}

