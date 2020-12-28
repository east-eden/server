package auto

import (
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/east-eden/server/excel"
)

var	attEntries	*AttEntries	//att.xlsx全局变量

// att.xlsx属性表
type AttEntry struct {
	Id        	int       	`json:"Id,omitempty"`	//id        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	Str       	int       	`json:"Str,omitempty"`	//力量        
	Agl       	int       	`json:"Agl,omitempty"`	//敏捷        
	Con       	int       	`json:"Con,omitempty"`	//耐力        
	Int       	int       	`json:"Int,omitempty"`	//智力        
	AtkSpeed  	int       	`json:"AtkSpeed,omitempty"`	//攻击速度      
	MaxHP     	int       	`json:"MaxHP,omitempty"`	//血量        
	MaxMP     	int       	`json:"MaxMP,omitempty"`	//蓝量        
	Atk       	int       	`json:"Atk,omitempty"`	//攻击力       
	Def       	int       	`json:"Def,omitempty"`	//物理防御力     
	CriProb   	int       	`json:"CriProb,omitempty"`	//暴击率       
	CriDmg    	int       	`json:"CriDmg,omitempty"`	//暴击伤害      
	EffectHit 	int       	`json:"EffectHit,omitempty"`	//效果命中      
	EffectResist	int       	`json:"EffectResist,omitempty"`	//效果抵抗      
	ConPercent	int       	`json:"ConPercent,omitempty"`	//体力百分比     
	AtkPercent	int       	`json:"AtkPercent,omitempty"`	//攻击百分比     
	DefPercent	int       	`json:"DefPercent,omitempty"`	//防御百分比     
}

// att.xlsx属性表集合
type AttEntries struct {
	Rows      	map[int]*AttEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("att.xlsx", attEntries)
}

func (e *AttEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	attEntries = &AttEntries{
		Rows: make(map[int]*AttEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &AttEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if utils.ErrCheck(err, "decode excel data to struct failed", v) {
	 		return err
	 	}

	 	attEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetAttEntry(id int) (*AttEntry, bool) {
	entry, ok := attEntries.Rows[id]
	return entry, ok
}

func  GetAttSize() int {
	return len(attEntries.Rows)
}

