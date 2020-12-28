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
	Id        	int       	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名称        
	Type      	int       	`json:"Type,omitempty"`	//类型        
	Pos       	int       	`json:"Pos,omitempty"`	//位置        
	Quality   	int       	`json:"Quality,omitempty"`	//品质        
	SuitID    	int       	`json:"SuitID,omitempty"`	//套装id      
}

// rune.xlsx属性表集合
type RuneEntries struct {
	Rows      	map[int]*RuneEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("rune.xlsx", runeEntries)
}

func (e *RuneEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	runeEntries = &RuneEntries{
		Rows: make(map[int]*RuneEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &RuneEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	runeEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetRuneEntry(id int) (*RuneEntry, bool) {
	entry, ok := runeEntries.Rows[id]
	return entry, ok
}

func  GetRuneSize() int {
	return len(runeEntries.Rows)
}

