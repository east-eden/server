package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	runeSuitEntries	*RuneSuitEntries	//runeSuit.xlsx全局变量

// runeSuit.xlsx属性表
type RuneSuitEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	//id        
	Suit2_AttID    	int32               	`json:"Suit2_AttID,omitempty"`	//2件套装属性ID  
	Suit3_AttID    	int32               	`json:"Suit3_AttID,omitempty"`	//3件套装属性ID  
	Suit4_AttID    	int32               	`json:"Suit4_AttID,omitempty"`	//2件套装属性ID  
	Suit5_AttID    	int32               	`json:"Suit5_AttID,omitempty"`	//2件套装属性ID  
	Suit6_AttID    	int32               	`json:"Suit6_AttID,omitempty"`	//2件套装属性ID  
}

// runeSuit.xlsx属性表集合
type RuneSuitEntries struct {
	Rows           	map[int32]*RuneSuitEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("runeSuit.xlsx", runeSuitEntries)
}

func (e *RuneSuitEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	runeSuitEntries = &RuneSuitEntries{
		Rows: make(map[int32]*RuneSuitEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &RuneSuitEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	runeSuitEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetRuneSuitEntry(id int32) (*RuneSuitEntry, bool) {
	entry, ok := runeSuitEntries.Rows[id]
	return entry, ok
}

func  GetRuneSuitSize() int32 {
	return int32(len(runeSuitEntries.Rows))
}

func  GetRuneSuitRows() map[int32]*RuneSuitEntry {
	return runeSuitEntries.Rows
}

