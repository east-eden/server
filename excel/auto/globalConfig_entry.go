package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var globalConfigEntries *GlobalConfigEntries //GlobalConfig.csv全局变量

// GlobalConfig.csv属性表
type GlobalConfigEntry struct {
	Id                             int32             `json:"Id,omitempty"`                             // 主键
	DemoSceneRadius                decimal.Decimal   `json:"DemoSceneRadius,omitempty"`                //Demo战斗场景半径，单位米(Demo完成后删除)
	ActRestFrameTime               decimal.Decimal   `json:"ActRestFrameTime,omitempty"`               //COM确定行动时静帧时间(移动无效)
	ActFailComReset                decimal.Decimal   `json:"ActFailComReset,omitempty"`                //行动失败时，COM条重置位置%(COM左侧作为起始)
	RaceAdvantage                  decimal.Decimal   `json:"RaceAdvantage,omitempty"`                  //种族有利时伤害加成
	MoveReadyTime                  decimal.Decimal   `json:"MoveReadyTime,omitempty"`                  //开始移动前的准备时间
	RageLimitEachLayer             int32             `json:"RageLimitEachLayer,omitempty"`             //每层的奥义怒气上限值
	ArmorRatio                     int32             `json:"ArmorRatio,omitempty"`                     //护甲减免系数
	DmgRatio                       int32             `json:"DmgRatio,omitempty"`                       //总伤害率系数
	CritRatio                      int32             `json:"CritRatio,omitempty"`                      //暴击率系数
	CritIncRatio                   int32             `json:"CritIncRatio,omitempty"`                   //暴击倍数加成系数
	TenacityRatio                  int32             `json:"TenacityRatio,omitempty"`                  //韧性率系数
	HealRatio                      int32             `json:"HealRatio,omitempty"`                      //治疗量系数
	EffectHitRatio                 int32             `json:"EffectHitRatio,omitempty"`                 //效果命中系数
	EffectResistRatio              int32             `json:"EffectResistRatio,omitempty"`              //效果抵抗系数
	ElementDmgRatio                int32             `json:"ElementDmgRatio,omitempty"`                //各伤害类型加成系数
	ElementResRatio                int32             `json:"ElementResRatio,omitempty"`                //各伤害类型抗性系数
	DamageRange                    []decimal.Decimal `json:"DamageRange,omitempty"`                    //伤害浮动区间系数
	MaterialContainerMax           int32             `json:"MaterialContainerMax,omitempty"`           //材料与消耗背包容量上限，超过此容量无法获得物品
	EquipContainerMax              int32             `json:"EquipContainerMax,omitempty"`              //装备背包容量上限，超过此容量无法获得物品
	CrystalContainerMax            int32             `json:"CrystalContainerMax,omitempty"`            //晶石容量上限，超过此容量无法获得物品
	EquipPromoteLevelLimit         []int32           `json:"EquipPromoteLevelLimit,omitempty"`         //装备突破[0-7]队伍等级限制
	EquipPromoteIntensityRatio     []int32           `json:"EquipPromoteIntensityRatio,omitempty"`     //装备突破[0-7]每次强度等级
	EquipLevelQualityRatio         []decimal.Decimal `json:"EquipLevelQualityRatio,omitempty"`         //装备升级和突破品质参数
	EquipSwallowExpLoss            decimal.Decimal   `json:"EquipSwallowExpLoss,omitempty"`            //装备吞噬经验折损率
	EquipExpItems                  []int32           `json:"EquipExpItems,omitempty"`                  //装备经验道具id
	EquipStarupCostItemNum         []int32           `json:"EquipStarupCostItemNum,omitempty"`         //装备升星次数消耗同id物品数量
	EquipStarupCostGold            []int32           `json:"EquipStarupCostGold,omitempty"`            //装备升星次数消耗金币
	EquipQualityStarupTimes        []int32           `json:"EquipQualityStarupTimes,omitempty"`        //装备不同品质升星次数限制[白,绿,蓝,紫,橙,红]
	HeroLevelupExpGoldRatio        int32             `json:"HeroLevelupExpGoldRatio,omitempty"`        //英雄升级经验对应消耗金币比例
	HeroExpItems                   []int32           `json:"HeroExpItems,omitempty"`                   //英雄经验道具id
	EquipLevelupExpGoldRatio       int32             `json:"EquipLevelupExpGoldRatio,omitempty"`       //装备升级经验对应消耗金币比例
	EquipLevelGrowRatioAttId       int32             `json:"EquipLevelGrowRatioAttId,omitempty"`       //装备升级成长率attid
	HeroPromoteLevelLimit          []int32           `json:"HeroPromoteLevelLimit,omitempty"`          //英雄突破[0-6]队伍等级限制
	HeroPromoteIntensityRatio      []int32           `json:"HeroPromoteIntensityRatio,omitempty"`      //英雄突破[0-6]每次强度等级
	HeroPromoteGrowupId            int32             `json:"HeroPromoteGrowupId,omitempty"`            //英雄突破成长率attid
	HeroPromoteBaseId              int32             `json:"HeroPromoteBaseId,omitempty"`              //英雄突破固定值attid
	HeroLevelQualityRatio          []decimal.Decimal `json:"HeroLevelQualityRatio,omitempty"`          //英雄升级和突破品质参数：N，R，SR，SSR，UR
	HeroLevelGrowRatioAttId        int32             `json:"HeroLevelGrowRatioAttId,omitempty"`        //英雄升级成长率attid
	CrystalSwallowExpLoss          decimal.Decimal   `json:"CrystalSwallowExpLoss,omitempty"`          //晶石吞噬经验折损率
	CrystalLevelupExpGoldRatio     int32             `json:"CrystalLevelupExpGoldRatio,omitempty"`     //晶石升级经验对应消耗金币比例
	CrystalExpItems                []int32           `json:"CrystalExpItems,omitempty"`                //晶石经验道具id
	CrystalLevelupIntensityRatio   int32             `json:"CrystalLevelupIntensityRatio,omitempty"`   //晶石升级强度系数
	CrystalLevelupMainQualityRatio []decimal.Decimal `json:"CrystalLevelupMainQualityRatio,omitempty"` //晶石升级主属性品质系数
	CrystalLevelupViceQualityRatio []decimal.Decimal `json:"CrystalLevelupViceQualityRatio,omitempty"` //晶石升级副属性品质系数
	CrystalLevelupRandRatio        []decimal.Decimal `json:"CrystalLevelupRandRatio,omitempty"`        //晶石升级副属性随机区间系数
	CrystalViceAttAddLevel         []int32           `json:"CrystalViceAttAddLevel,omitempty"`         //晶石升级到3，6，9，12，15级时强化副属性
	CrystalLevelupQualityLimit     []int32           `json:"CrystalLevelupQualityLimit,omitempty"`     //各品质晶石强化等级上限
	CrystalLevelupAssistantNumber  int32             `json:"CrystalLevelupAssistantNumber,omitempty"`  //晶石副属性随机到相同属性的次数上限
	SweepStageItem                 int32             `json:"SweepStageItem,omitempty"`                 //扫荡券物品id
	MailExpireTime                 int32             `json:"MailExpireTime,omitempty"`                 //邮件过期时间30天（秒）
	CollectionMaxStar              []int32           `json:"CollectionMaxStar,omitempty"`              //不同品质收集品星级上限(绿、蓝、紫、黄)
	CollectionQualityScoreFactor   []decimal.Decimal `json:"CollectionQualityScoreFactor,omitempty"`   //收集品评分品质系数(绿、蓝、紫、黄)
	CollectionStarScoreFactor      []decimal.Decimal `json:"CollectionStarScoreFactor,omitempty"`      //收集品评分星级系数(0星开始，最多6星)
	CollectionWakeupScoreFactor    decimal.Decimal   `json:"CollectionWakeupScoreFactor,omitempty"`    //收集品评分觉醒系数
}

// GlobalConfig.csv属性表集合
type GlobalConfigEntries struct {
	Rows map[int32]*GlobalConfigEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("GlobalConfig.csv", (*GlobalConfigEntries)(nil))
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

func GetGlobalConfigEntry(id int32) (*GlobalConfigEntry, bool) {
	entry, ok := globalConfigEntries.Rows[id]
	return entry, ok
}

func GetGlobalConfigSize() int32 {
	return int32(len(globalConfigEntries.Rows))
}

func GetGlobalConfigRows() map[int32]*GlobalConfigEntry {
	return globalConfigEntries.Rows
}
