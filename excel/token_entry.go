package excel

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	tokenEntries	*TokenEntries	//token.xlsx全局变量

// token.xlsx属性表
type TokenEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名字        
	MaxHold   	int       	`json:"MaxHold,omitempty"`	//持有上限      
}

// token.xlsx属性表集合
type TokenEntries struct {
	Rows      	map[int]*TokenEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	AddEntries("token.xlsx", heroEntries)
}

func (e *TokenEntries) load(excelFileRaw *ExcelFileRaw) error {
	
	tokenEntries = &TokenEntries{
		Rows: make(map[int]*TokenEntry),
	}

	for _, v := range excelFileRaw.cellData {
		entry := &TokenEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	tokenEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.filename).Msg("excel load success")
	return nil
	
}

func  GetTokenEntry(id int) (*TokenEntry, bool) {
	entry, ok := tokenEntries.Rows[id]
	return entry, ok
}

