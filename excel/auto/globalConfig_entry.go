package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var	globalConfigEntries	*GlobalConfigEntries	//GlobalConfig.xlsx全局变量

// GlobalConfig.xlsx属性表
type GlobalConfigEntry struct {
	Id             	int32               	`json:"Id,omitempty"`	// 主键       
	ArmorRatio     	int32               	`json:"ArmorRatio,omitempty"`	//护甲减免系数    
	DmgRatio       	int32               	`json:"DmgRatio,omitempty"`	//总伤害率系数    
	CritRatio      	int32               	`json:"CritRatio,omitempty"`	//暴击率系数     
	CritIncRatio   	int32               	`json:"CritIncRatio,omitempty"`	//暴击倍数加成系数  
	HealRatio      	int32               	`json:"HealRatio,omitempty"`	//治疗量系数     
	EffectHitRatio 	int32               	`json:"EffectHitRatio,omitempty"`	//效果命中系数    
	EffectResistRatio	int32               	`json:"EffectResistRatio,omitempty"`	//效果抵抗系数    
	ElementDmgRatio	int32               	`json:"ElementDmgRatio,omitempty"`	//各伤害类型加成系数 
	ElementResRatio	int32               	`json:"ElementResRatio,omitempty"`	//各伤害类型抗性系数 
	MaterialContainerMax	int32               	`json:"MaterialContainerMax,omitempty"`	//材料与消耗背包容量上限，超过此容量无法获得物品
	EquipContainerMax	int32               	`json:"EquipContainerMax,omitempty"`	//装备背包容量上限，超过此容量无法获得物品
	CrystalContainerMax	int32               	`json:"CrystalContainerMax,omitempty"`	//晶石容量上限，超过此容量无法获得物品
	EquipPromoteLevelLimit	[]int32             	`json:"EquipPromoteLevelLimit,omitempty"`	//装备突破队伍等级限制
	EquipLevelQualityRatio	[]int32             	`json:"EquipLevelQualityRatio,omitempty"`	//装备升级品质参数，100%=10000
	EquipSwallowExpLoss	int32               	`json:"EquipSwallowExpLoss,omitempty"`	//装备吞噬经验折损率 
	EquipSwallowGoldLoss	int32               	`json:"EquipSwallowGoldLoss,omitempty"`	//装备吞噬金币折损率 
	EquipLevelGrowRatioAttId	int32               	`json:"EquipLevelGrowRatioAttId,omitempty"`	//装备升级成长率attid
	HeroPromoteLevelLimit	[]int32             	`json:"HeroPromoteLevelLimit,omitempty"`	//英雄突破队伍等级限制
	HeroLevelQualityRatio	[]int32             	`json:"HeroLevelQualityRatio,omitempty"`	//英雄升级品质参数，100%=10000
	HeroLevelGrowRatioAttId	int32               	`json:"HeroLevelGrowRatioAttId,omitempty"`	//英雄升级成长率attid
	CrystalSwallowExpLoss	int32               	`json:"CrystalSwallowExpLoss,omitempty"`	//晶石吞噬经验折损率 
	CrystalSwallowGoldLoss	int32               	`json:"CrystalSwallowGoldLoss,omitempty"`	//晶石吞噬金币折损率 
	CrystalLevelupIntensityRatio	int32               	`json:"CrystalLevelupIntensityRatio,omitempty"`	//晶石升级强度系数  
	CrystalLevelupQualityRatio	[]int32             	`json:"CrystalLevelupQualityRatio,omitempty"`	//晶石升级品质系数  
	CrystalLevelupRandRatio	[]int32             	`json:"CrystalLevelupRandRatio,omitempty"`	//晶石升级随机区间系数
}

// GlobalConfig.xlsx属性表集合
type GlobalConfigEntries struct {
	Rows           	map[int32]*GlobalConfigEntry	`json:"Rows,omitempty"`	//          
}

func  init()  {
	excel.AddEntryLoader("GlobalConfig.xlsx", (*GlobalConfigEntries)(nil))
}

func (e *GlobalConfigEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {
	
	globalConfigEntries = &GlobalConfigEntries{
		Rows: make(map[int32]*GlobalConfigEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &GlobalConfigEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

	 	globalConfigEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil
	
}

func  GetGlobalConfigEntry(id int32) (*GlobalConfigEntry, bool) {
	entry, ok := globalConfigEntries.Rows[id]
	return entry, ok
}

func  GetGlobalConfigSize() int32 {
	return int32(len(globalConfigEntries.Rows))
}

func  GetGlobalConfigRows() map[int32]*GlobalConfigEntry {
	return globalConfigEntries.Rows
}

