package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	heroPromoteEntries	*HeroPromoteEntries	//HeroPromote.xlsx全局变量

// HeroPromote.xlsx属性表
type HeroPromoteEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	PromoteCostId1 	int32               	`json:"PromoteCostId1,omitempty"`	//突破1消耗id   
	PromoteCostId2 	int32               	`json:"PromoteCostId2,omitempty"`	//突破2消耗id   
	PromoteCostId3 	int32               	`json:"PromoteCostId3,omitempty"`	//突破3消耗id   
	PromoteCostId4 	int32               	`json:"PromoteCostId4,omitempty"`	//突破4消耗id   
	PromoteCostId5 	int32               	`json:"PromoteCostId5,omitempty"`	//突破5消耗id   
	PromoteCostId6 	int32               	`json:"PromoteCostId6,omitempty"`	//突破6消耗id   
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

