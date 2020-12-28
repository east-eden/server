package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	formulaEntries	*FormulaEntries	//formula.xlsx全局变量

// formula.xlsx属性表
type FormulaEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Type      	int       	`json:"Type,omitempty"`	//公式类型      
	Formula   	string    	`json:"Formula,omitempty"`	//伤害公式      
}

// formula.xlsx属性表集合
type FormulaEntries struct {
	Rows      	map[int]*FormulaEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("formula.xlsx", formulaEntries)
}

func (e *FormulaEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	formulaEntries = &FormulaEntries{
		Rows: make(map[int]*FormulaEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &FormulaEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	formulaEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetFormulaEntry(id int) (*FormulaEntry, bool) {
	entry, ok := formulaEntries.Rows[id]
	return entry, ok
}

func  GetFormulaSize() int {
	return len(formulaEntries.Rows)
}

