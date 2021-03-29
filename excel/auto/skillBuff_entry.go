package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/emirpasic/gods/maps/treemap"
)

var	skillBuffEntries	*SkillBuffEntries	//SkillBuff.xlsx全局变量

// SkillBuff.xlsx属性表
type SkillBuffEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	Name           	string              	`json:"Name,omitempty"`	//名字        
	Desc           	string              	`json:"Desc,omitempty"`	//buff描述    
	Continueplay   	string              	`json:"Continueplay,omitempty"`	//持续表现      
	Startplay      	string              	`json:"Startplay,omitempty"`	//结束表现      
	LifeTime       	float32             	`json:"LifeTime,omitempty"`	//持续时间      
	BuffLevel      	int32               	`json:"BuffLevel,omitempty"`	//buff等级    
	BuffGroup      	int32               	`json:"BuffGroup,omitempty"`	//buff分组    
	BuffOverlap    	int32               	`json:"BuffOverlap,omitempty"`	//是否重置      
	MaxLimit       	int32               	`json:"MaxLimit,omitempty"`	//层数限制      
	Actplay        	string              	`json:"Actplay,omitempty"`	//触发表演      
	Cd             	float32             	`json:"Cd,omitempty"`	//冷却时间      
	Prob           	int32               	`json:"Prob,omitempty"`	//触发概率      
	BuffEffectType 	int32               	`json:"BuffEffectType,omitempty"`	//效果类型      
	A              	int32               	`json:"A,omitempty"`	//          
	B              	int32               	`json:"B,omitempty"`	//          
	C              	int32               	`json:"C,omitempty"`	//          
	D              	int32               	`json:"D,omitempty"`	//          
	E              	int32               	`json:"E,omitempty"`	//          
	AttributeNumValue1	*treemap.Map        	`json:"AttributeNumValue1,omitempty"`	//属性1       
	AttributeNumValue2	*treemap.Map        	`json:"AttributeNumValue2,omitempty"`	//属性2       
	AttributeNumValue3	*treemap.Map        	`json:"AttributeNumValue3,omitempty"`	//属性3       
	RangeType      	int32               	`json:"RangeType,omitempty"`	//目标范围      
	TargetType     	int32               	`json:"TargetType,omitempty"`	//目标类型      
	TargetLength   	int32               	`json:"TargetLength,omitempty"`	//范围长       
	TargetWide     	int32               	`json:"TargetWide,omitempty"`	//范围宽       
	Scope          	int32               	`json:"Scope,omitempty"`	//作用对象      
}

// SkillBuff.xlsx属性表集合
type SkillBuffEntries struct {
	Rows           	map[int32]*SkillBuffEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("SkillBuff.xlsx", (*SkillBuffEntries)(nil))
}

func (e *SkillBuffEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	skillBuffEntries = &SkillBuffEntries{
		Rows: make(map[int32]*SkillBuffEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillBuffEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	skillBuffEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetSkillBuffEntry(id int32) (*SkillBuffEntry, bool) {
	entry, ok := skillBuffEntries.Rows[id]
	return entry, ok
}

func  GetSkillBuffSize() int32 {
	return int32(len(skillBuffEntries.Rows))
}

func  GetSkillBuffRows() map[int32]*SkillBuffEntry {
	return skillBuffEntries.Rows
}
