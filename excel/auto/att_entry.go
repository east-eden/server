package auto

import (
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"e.coding.net/mmstudio/blade/server/excel"
)

var	attEntries	*AttEntries	//att.xlsx全局变量

// att.xlsx属性表
type AttEntry struct {
	Id        	int32     	`json:"Id,omitempty"`	//id        
	Desc      	string    	`json:"Desc,omitempty"`	//描述        
	Str       	int32     	`json:"Str,omitempty"`	//力量        
	Agl       	int32     	`json:"Agl,omitempty"`	//敏捷        
	Con       	int32     	`json:"Con,omitempty"`	//耐力        
	Int       	int32     	`json:"Int,omitempty"`	//智力        
	AtkSpeed  	int32     	`json:"AtkSpeed,omitempty"`	//攻击速度      
	MaxHP     	int32     	`json:"MaxHP,omitempty"`	//血量        
	MaxMP     	int32     	`json:"MaxMP,omitempty"`	//蓝量        
	Atk       	int32     	`json:"Atk,omitempty"`	//攻击力       
	Def       	int32     	`json:"Def,omitempty"`	//物理防御力     
	CriProb   	int32     	`json:"CriProb,omitempty"`	//暴击率       
	CriDmg    	int32     	`json:"CriDmg,omitempty"`	//暴击伤害      
	EffectHit 	int32     	`json:"EffectHit,omitempty"`	//效果命中      
	EffectResist	int32     	`json:"EffectResist,omitempty"`	//效果抵抗      
	ConPercent	int32     	`json:"ConPercent,omitempty"`	//体力百分比     
	AtkPercent	int32     	`json:"AtkPercent,omitempty"`	//攻击百分比     
	DefPercent	int32     	`json:"DefPercent,omitempty"`	//防御百分比     
}

// att.xlsx属性表集合
type AttEntries struct {
	Rows      	map[int32]*AttEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("att.xlsx", attEntries)
}

func (e *AttEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	attEntries = &AttEntries{
		Rows: make(map[int32]*AttEntry),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &AttEntry{}
	 	err := mapstructure.Decode(v, entry)
	 	if event, pass := utils.ErrCheck(err, v); !pass {
			event.Msg("decode excel data to struct failed")
	 		return err
	 	}

	 	attEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetAttEntry(id int32) (*AttEntry, bool) {
	entry, ok := attEntries.Rows[id]
	return entry, ok
}

func  GetAttSize() int32 {
	return int32(len(attEntries.Rows))
}

