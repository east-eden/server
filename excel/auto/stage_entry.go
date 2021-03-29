package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	stageEntries   	*StageEntries  	//Stage.xlsx全局变量 

// Stage.xlsx属性表
type StageEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	PrevStageId    	int32               	`json:"PrevStageId,omitempty"`	//前置关卡id    
	ConditionId    	int32               	`json:"ConditionId,omitempty"`	//解锁条件id    
	CostStrength   	int32               	`json:"CostStrength,omitempty"`	//消耗体力      
	ChapterId      	int32               	`json:"ChapterId,omitempty"`	//所属章节id    
	Type           	int32               	`json:"Type,omitempty"`	//关卡类型      
	StarConditionIds	[]int32             	`json:"StarConditionIds,omitempty"`	//星级挑战条件列表  
	FirstRewardLootId	int32               	`json:"FirstRewardLootId,omitempty"`	//首次通关掉落id  
	RewardLootId   	int32               	`json:"RewardLootId,omitempty"`	//通关掉落id    
	DailyTimes     	int32               	`json:"DailyTimes,omitempty"`	//每日挑战次数    
}

// Stage.xlsx属性表集合
type StageEntries struct {
	Rows           	map[int32]*StageEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("Stage.xlsx", (*StageEntries)(nil))
}

func (e *StageEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	stageEntries = &StageEntries{
		Rows: make(map[int32]*StageEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &StageEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	stageEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetStageEntry(id int32) (*StageEntry, bool) {
	entry, ok := stageEntries.Rows[id]
	return entry, ok
}

func  GetStageSize() int32 {
	return int32(len(stageEntries.Rows))
}

func  GetStageRows() map[int32]*StageEntry {
	return stageEntries.Rows
}
