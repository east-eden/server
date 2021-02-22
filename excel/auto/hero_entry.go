package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	heroEntries    	*HeroEntries   	//Hero.xlsx全局变量  

// Hero.xlsx属性表
type HeroEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	Quality        	int32               	`json:"Quality,omitempty"`	//品质        
	Race           	int32               	`json:"Race,omitempty"`	//种族        
	AttId          	int32               	`json:"AttId,omitempty"`	//属性id      
	FragmentCompose	int32               	`json:"FragmentCompose,omitempty"`	//合成卡牌所需碎片  
	FragmentTransform	int32               	`json:"FragmentTransform,omitempty"`	//重复获得卡牌转化碎片数
}

// Hero.xlsx属性表集合
type HeroEntries struct {
	Rows           	map[int32]*HeroEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("Hero.xlsx", (*HeroEntries)(nil))
}

func (e *HeroEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	heroEntries = &HeroEntries{
		Rows: make(map[int32]*HeroEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &HeroEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	heroEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetHeroEntry(id int32) (*HeroEntry, bool) {
	entry, ok := heroEntries.Rows[id]
	return entry, ok
}

func  GetHeroSize() int32 {
	return int32(len(heroEntries.Rows))
}

func  GetHeroRows() map[int32]*HeroEntry {
	return heroEntries.Rows
}

