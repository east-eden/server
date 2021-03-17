package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	crystalAttRepoEntries	*CrystalAttRepoEntries	//CrystalAttRepo.xlsx全局变量

// CrystalAttRepo.xlsx属性表
type CrystalAttRepoEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	Desc           	string              	`json:"Desc,omitempty"`	//属性描述      
	Pos            	int32               	`json:"Pos,omitempty"`	//晶石位置      
	Type           	int32               	`json:"Type,omitempty"`	//属性库类型     
	AttId          	int32               	`json:"AttId,omitempty"`	//属性id      
	AttGrowRatioId 	int32               	`json:"AttGrowRatioId,omitempty"`	//属性成长率     
	AttWeight      	int32               	`json:"AttWeight,omitempty"`	//属性权重      
}

// CrystalAttRepo.xlsx属性表集合
type CrystalAttRepoEntries struct {
	Rows           	map[int32]*CrystalAttRepoEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("CrystalAttRepo.xlsx", (*CrystalAttRepoEntries)(nil))
}

func (e *CrystalAttRepoEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	crystalAttRepoEntries = &CrystalAttRepoEntries{
		Rows: make(map[int32]*CrystalAttRepoEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CrystalAttRepoEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	crystalAttRepoEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetCrystalAttRepoEntry(id int32) (*CrystalAttRepoEntry, bool) {
	entry, ok := crystalAttRepoEntries.Rows[id]
	return entry, ok
}

func  GetCrystalAttRepoSize() int32 {
	return int32(len(crystalAttRepoEntries.Rows))
}

func  GetCrystalAttRepoRows() map[int32]*CrystalAttRepoEntry {
	return crystalAttRepoEntries.Rows
}

