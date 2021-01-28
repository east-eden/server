package auto

import (
	"bitbucket.org/east-eden/server/excel"
	"bitbucket.org/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"fmt"
	"strconv"
	"strings"
	"github.com/emirpasic/gods/maps/treemap"
)

var	buffEntries    	*BuffEntries   	//buff.xlsx全局变量  

// buff.xlsx属性表
type BuffEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	//id 多主键之一  
	BuffType       	int32               	`json:"BuffType,omitempty"`	//buff类型 多主键之一


	Level          	int32               	`json:"Level,omitempty"`	//等级        
	NextLevel      	int32               	`json:"NextLevel,omitempty"`	//下个等级      
	CD             	float32             	`json:"CD,omitempty"`	//冷却时间(秒)   
	LifeTime       	float32             	`json:"LifeTime,omitempty"`	//持续时间(秒)   
	BuffOverlap    	[]int32             	`json:"BuffOverlap,omitempty"`	//叠加类型      
	MaxLimit       	int32               	`json:"MaxLimit,omitempty"`	//限制        
	Params_StrValue	[]string            	`json:"Params_StrValue,omitempty"`	//参数列表，目标属性 
	Params_Formula 	*treemap.Map        	`json:"Params_Formula,omitempty"`	//公式        
	Params_NumValue	interface{}         	`json:"Params_NumValue,omitempty"`	//参数列表，固定数值 

	TestField      	string              	`json:"TestField,omitempty"`	//测试字段      

}

// buff.xlsx属性表集合
type BuffEntries struct {
	Rows           	map[string]*BuffEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntries("buff.xlsx", buffEntries)
}

func (e *BuffEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	buffEntries = &BuffEntries{
		Rows: make(map[string]*BuffEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &BuffEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		key := fmt.Sprintf("%d+%d", entry.Id, entry.BuffType)
	 	buffEntries.Rows[key] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetBuffEntry(keys ...int32) (*BuffEntry, bool) {
	keyName := make([]string, 0, len(keys))
	for _, key := range keys {
		keyName = append(keyName, strconv.Itoa(int(key)))
	}

	finalKey := strings.Join(keyName, "+")
	entry, ok := buffEntries.Rows[finalKey]
	return entry, ok 
}

func  GetBuffSize() int32 {
	return int32(len(buffEntries.Rows))
}

func  GetBuffRows() map[string]*BuffEntry {
	return buffEntries.Rows
}

