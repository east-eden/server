package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	crystalSpellEntries	*CrystalSpellEntries	//CrystalSpell.xlsx全局变量

// CrystalSpell.xlsx属性表
type CrystalSpellEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	SpellId        	[]int32             	`json:"SpellId,omitempty"`	//元素技能id    
}

// CrystalSpell.xlsx属性表集合
type CrystalSpellEntries struct {
	Rows           	map[int32]*CrystalSpellEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("CrystalSpell.xlsx", (*CrystalSpellEntries)(nil))
}

func (e *CrystalSpellEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	crystalSpellEntries = &CrystalSpellEntries{
		Rows: make(map[int32]*CrystalSpellEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &CrystalSpellEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	crystalSpellEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetCrystalSpellEntry(id int32) (*CrystalSpellEntry, bool) {
	entry, ok := crystalSpellEntries.Rows[id]
	return entry, ok
}

func  GetCrystalSpellSize() int32 {
	return int32(len(crystalSpellEntries.Rows))
}

func  GetCrystalSpellRows() map[int32]*CrystalSpellEntry {
	return crystalSpellEntries.Rows
}

