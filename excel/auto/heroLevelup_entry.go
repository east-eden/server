package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	heroLevelupEntries	*HeroLevelupEntries	//HeroLevelup.xlsx全局变量

// HeroLevelup.xlsx属性表
type HeroLevelupEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	Exp            	int32               	`json:"Exp,omitempty"`	//品质        
	PromoteLimit   	int32               	`json:"PromoteLimit,omitempty"`	//突破次数限制    
}

// HeroLevelup.xlsx属性表集合
type HeroLevelupEntries struct {
	Rows           	map[int32]*HeroLevelupEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("HeroLevelup.xlsx", (*HeroLevelupEntries)(nil))
}

func (e *HeroLevelupEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	heroLevelupEntries = &HeroLevelupEntries{
		Rows: make(map[int32]*HeroLevelupEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &HeroLevelupEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	heroLevelupEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetHeroLevelupEntry(id int32) (*HeroLevelupEntry, bool) {
	entry, ok := heroLevelupEntries.Rows[id]
	return entry, ok
}

func  GetHeroLevelupSize() int32 {
	return int32(len(heroLevelupEntries.Rows))
}

func  GetHeroLevelupRows() map[int32]*HeroLevelupEntry {
	return heroLevelupEntries.Rows
}

