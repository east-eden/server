package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	runeEntries	*RuneEntries	//rune.xlsx全局变量

// rune.xlsx属性表
type RuneEntry struct {
	Id        	int32     	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名称        
	Type      	int32     	`json:"Type,omitempty"`	//类型        
	Pos       	int32     	`json:"Pos,omitempty"`	//位置        
	Quality   	int32     	`json:"Quality,omitempty"`	//品质        
	SuitID    	int32     	`json:"SuitID,omitempty"`	//套装id      
}

// rune.xlsx属性表集合
type RuneEntries struct {
	Rows      	map[int32]*RuneEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("rune.xlsx", runeEntries)
}

func (e *RuneEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	runeEntries = &RuneEntries{
		Rows: make(map[int32]*RuneEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &RuneEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	runeEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetRuneEntry(id int32) (*RuneEntry, bool) {
	entry, ok := runeEntries.Rows[id]
	return entry, ok
}

func  GetRuneSize() int32 {
	return int32(len(runeEntries.Rows))
}

