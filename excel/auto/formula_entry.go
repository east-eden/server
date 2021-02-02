package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	formulaEntries 	*FormulaEntries	//Formula.xlsx全局变量

// Formula.xlsx属性表
type FormulaEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	Type           	int32               	`json:"Type,omitempty"`	//公式类型      
	Formula        	string              	`json:"Formula,omitempty"`	//伤害公式      
}

// Formula.xlsx属性表集合
type FormulaEntries struct {
	Rows           	map[int32]*FormulaEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("Formula.xlsx", (*FormulaEntries)(nil))
}

func (e *FormulaEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	formulaEntries = &FormulaEntries{
		Rows: make(map[int32]*FormulaEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &FormulaEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	formulaEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetFormulaEntry(id int32) (*FormulaEntry, bool) {
	entry, ok := formulaEntries.Rows[id]
	return entry, ok
}

func  GetFormulaSize() int32 {
	return int32(len(formulaEntries.Rows))
}

func  GetFormulaRows() map[int32]*FormulaEntry {
	return formulaEntries.Rows
}

