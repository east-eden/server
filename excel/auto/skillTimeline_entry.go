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
	Type           	int32               	`json:"Type,omitempty"`	//类型        
	DurationTime   	number              	`json:"DurationTime,omitempty"`	//持续时间      
	AnimName       	string              	`json:"AnimName,omitempty"`	//动作名称      
	FxName         	string              	`json:"FxName,omitempty"`	//特效名称,做到一个文件里,直接挂在人身上
	BulletFx       	string              	`json:"BulletFx,omitempty"`	//          
	BulletSpeed    	number              	`json:"BulletSpeed,omitempty"`	//          
	EffectTime     	number              	`json:"EffectTime,omitempty"`	//固定范围是受击的时间点,单体弹道是发出的时间点
	EffectIndex    	int32               	`json:"EffectIndex,omitempty"`	//第几个effect 
	HitAnimName    	string              	`json:"HitAnimName,omitempty"`	//受击动作      
	HitFxName      	string              	`json:"HitFxName,omitempty"`	//受击特效      
	HitFxSlot      	string              	`json:"HitFxSlot,omitempty"`	//受击特效插槽    
	HitStopTime    	number              	`json:"HitStopTime,omitempty"`	//受击停顿时间    
	HitBackDistance	number              	`json:"HitBackDistance,omitempty"`	//击退距离      
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

