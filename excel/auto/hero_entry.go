package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
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
	excel.AddEntries("hero.xlsx", heroEntries)
}

func (e *HeroEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	heroEntries = &HeroEntries{
		Rows: make(map[int]*HeroEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &HeroEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	heroEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetHeroEntry(id int) (*HeroEntry, bool) {
	entry, ok := heroEntries.Rows[id]
	return entry, ok
}

func  GetHeroSize() int {
	return len(heroEntries.Rows)
}

