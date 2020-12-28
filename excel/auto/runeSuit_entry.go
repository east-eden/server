package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	runeSuitEntries	*RuneSuitEntries	//runeSuit.xlsx全局变量

// runeSuit.xlsx属性表
type RuneSuitEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Suit2_AttID	int       	`json:"Suit2_AttID,omitempty"`	//2件套装属性ID  
	Suit3_AttID	int       	`json:"Suit3_AttID,omitempty"`	//3件套装属性ID  
	Suit4_AttID	int       	`json:"Suit4_AttID,omitempty"`	//2件套装属性ID  
	Suit5_AttID	int       	`json:"Suit5_AttID,omitempty"`	//2件套装属性ID  
	Suit6_AttID	int       	`json:"Suit6_AttID,omitempty"`	//2件套装属性ID  
}

// runeSuit.xlsx属性表集合
type RuneSuitEntries struct {
	Rows      	map[int]*RuneSuitEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("runeSuit.xlsx", runeSuitEntries)
}

func (e *RuneSuitEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	runeSuitEntries = &RuneSuitEntries{
		Rows: make(map[int]*RuneSuitEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &RuneSuitEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	runeSuitEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetRuneSuitEntry(id int) (*RuneSuitEntry, bool) {
	entry, ok := runeSuitEntries.Rows[id]
	return entry, ok
}

func  GetRuneSuitSize() int {
	return len(runeSuitEntries.Rows)
}
