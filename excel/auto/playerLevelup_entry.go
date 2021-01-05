package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	playerLevelupEntries	*PlayerLevelupEntries	//playerLevelup.xlsx全局变量

// playerLevelup.xlsx属性表
type PlayerLevelupEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//等级        
	Exp       	int       	`json:"Exp,omitempty"`	//达到此等级需要的经验值
}

// playerLevelup.xlsx属性表集合
type PlayerLevelupEntries struct {
	Rows      	map[int]*PlayerLevelupEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("playerLevelup.xlsx", playerLevelupEntries)
}

func (e *PlayerLevelupEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	playerLevelupEntries = &PlayerLevelupEntries{
		Rows: make(map[int]*PlayerLevelupEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &PlayerLevelupEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	playerLevelupEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetPlayerLevelupEntry(id int) (*PlayerLevelupEntry, bool) {
	entry, ok := playerLevelupEntries.Rows[id]
	return entry, ok
}

func  GetPlayerLevelupSize() int {
	return len(playerLevelupEntries.Rows)
}

