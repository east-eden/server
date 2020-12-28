package excel

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	heroEntries	*HeroEntries	//hero.xlsx全局变量

// hero.xlsx属性表
type HeroEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名字        
	AttID     	int       	`json:"AttID,omitempty"`	//属性id      
	Quality   	int       	`json:"Quality,omitempty"`	//品质        
}

// hero.xlsx属性表集合
type HeroEntries struct {
	Rows      	map[int]*HeroEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	AddEntries("hero.xlsx", heroEntries)
}

func (e *HeroEntries) load(excelFileRaw *ExcelFileRaw) error {
	
	heroEntries = &HeroEntries{
		Rows: make(map[int]*HeroEntry),
	}

	for _, v := range excelFileRaw.cellData {
		entry := &HeroEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	heroEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.filename).Msg("excel load success")
	return nil
	
}

func  GetHeroEntry(id int) (*HeroEntry, bool) {
	entry, ok := heroEntries.Rows[id]
	return entry, ok
}

