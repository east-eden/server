package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	globalConfigEntries	*GlobalConfigEntries	//GlobalConfig.xlsx全局变量

// GlobalConfig.xlsx属性表
type GlobalConfigEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	HeroMaxHP      	int32               	`json:"HeroMaxHP,omitempty"`	// 英雄血量上限   
	HeroMaxMP      	int32               	`json:"HeroMaxMP,omitempty"`	//英雄蓝量上限    
	HeroMaxRage    	int32               	`json:"HeroMaxRage,omitempty"`	//英雄怒气值上限   
	HeroMaxAtk     	int32               	`json:"HeroMaxAtk,omitempty"`	//英雄攻击上限    
	HeroDefName    	string              	`json:"HeroDefName,omitempty"`	//英雄默认名字    
	PlayerMaxName  	int32               	`json:"PlayerMaxName,omitempty"`	//角色名字字数上限  
}

// GlobalConfig.xlsx属性表集合
type GlobalConfigEntries struct {
	Rows           	map[int32]*GlobalConfigEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("GlobalConfig.xlsx", (*GlobalConfigEntries)(nil))
}

func (e *GlobalConfigEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	globalConfigEntries = &GlobalConfigEntries{
		Rows: make(map[int32]*GlobalConfigEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &GlobalConfigEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	globalConfigEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetGlobalConfigEntry(id int32) (*GlobalConfigEntry, bool) {
	entry, ok := globalConfigEntries.Rows[id]
	return entry, ok
}

func  GetGlobalConfigSize() int32 {
	return int32(len(globalConfigEntries.Rows))
}

func  GetGlobalConfigRows() map[int32]*GlobalConfigEntry {
	return globalConfigEntries.Rows
}

