package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	heroPromoteEntries	*HeroPromoteEntries	//HeroPromote.xlsx全局变量

// HeroPromote.xlsx属性表
type HeroPromoteEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	PromoteCostId  	[]int32             	`json:"PromoteCostId,omitempty"`	//突破消耗id    
	PromoteAttId   	[]int32             	`json:"PromoteAttId,omitempty"`	//英雄突破属性    
}

// HeroPromote.xlsx属性表集合
type HeroPromoteEntries struct {
	Rows           	map[int32]*HeroPromoteEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("HeroPromote.xlsx", (*HeroPromoteEntries)(nil))
}

func (e *HeroPromoteEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	heroPromoteEntries = &HeroPromoteEntries{
		Rows: make(map[int32]*HeroPromoteEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &HeroPromoteEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	heroPromoteEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetHeroPromoteEntry(id int32) (*HeroPromoteEntry, bool) {
	entry, ok := heroPromoteEntries.Rows[id]
	return entry, ok
}

func  GetHeroPromoteSize() int32 {
	return int32(len(heroPromoteEntries.Rows))
}

func  GetHeroPromoteRows() map[int32]*HeroPromoteEntry {
	return heroPromoteEntries.Rows
}

