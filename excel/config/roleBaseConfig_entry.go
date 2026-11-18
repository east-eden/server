package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var roleBaseConfigEntries *RoleBaseConfigEntries //RoleBaseConfig.json全局变量

// 属性表 英雄表-拳皇/roleBaseConfig
type RoleBaseConfigEntry struct {
	Test                              int32                    `json:"test,omitempty" bson:"test"`                                                           //##保留主键id
	HeroLinePosCount                  int32                    `json:"heroLinePosCount,omitempty" bson:"heroLinePosCount"`                                   //最大上阵英雄数量
	DEF_ArgA                          float32                  `json:"dEF_ArgA,omitempty" bson:"dEF_ArgA"`                                                   //防御免伤参数A
	DEF_ArgC                          float32                  `json:"dEF_ArgC,omitempty" bson:"dEF_ArgC"`                                                   //防御免伤参数C
	Init_Blast                        float32                  `json:"init_Blast,omitempty" bson:"init_Blast"`                                               //初始暴击倍率
	Block_HpSub                       float32                  `json:"block_HpSub,omitempty" bson:"block_HpSub"`                                             //格挡免伤
	Stun_Time                         float32                  `json:"stun_Time,omitempty" bson:"stun_Time"`                                                 //硬直时间
	Init_SpRecover                    float32                  `json:"init_SpRecover,omitempty" bson:"init_SpRecover"`                                       //初始怒气回复速度
	DamageFloat                       float32                  `json:"damageFloat,omitempty" bson:"damageFloat"`                                             //伤害波动系数
	MapUnitConvert                    float32                  `json:"mapUnitConvert,omitempty" bson:"mapUnitConvert"`                                       //距离转换常数
	Block_Time                        float32                  `json:"block_Time,omitempty" bson:"block_Time"`                                               //格挡时间
	Suck_HpPrecent                    float32                  `json:"suck_HpPrecent,omitempty" bson:"suck_HpPrecent"`                                       //固定吸血比例
	Combo_param                       float32                  `json:"combo_param,omitempty" bson:"combo_param"`                                             //连击衰减系数
	Init_HealBlast                    float32                  `json:"init_HealBlast,omitempty" bson:"init_HealBlast"`                                       //初始治疗暴击率
	HeroLevelupLimitStart             int32                    `json:"heroLevelupLimitStart,omitempty" bson:"heroLevelupLimitStart"`                         //英雄至少可以升到 x 级
	HeroSoulNums                      int32                    `json:"heroSoulNums,omitempty" bson:"heroSoulNums"`                                           //英雄战魂数量
	HeroEquipsUnlockLevel             map[HeroEquipPos]int32   `json:"heroEquipsUnlockLevel,omitempty" bson:"heroEquipsUnlockLevel"`                         //英雄装备位置解锁所需等级
	HeroAttrAddPerStar                []*AttrValue             `json:"heroAttrAddPerStar,omitempty" bson:"heroAttrAddPerStar"`                               //英雄每次升星属性加成
	HeroSoulPosUnlockLevel            map[int32]int32          `json:"heroSoulPosUnlockLevel,omitempty" bson:"heroSoulPosUnlockLevel"`                       //英雄战魂位置索引解锁所需等级
	HeroSoulNames                     []string                 `json:"heroSoulNames,omitempty" bson:"heroSoulNames"`                                         //英雄战魂名称
	HeroSoulComposeNum                int32                    `json:"heroSoulComposeNum,omitempty" bson:"heroSoulComposeNum"`                               //英雄战魂合成所需数量
	HeroSoulRefiningExpendNum         int32                    `json:"heroSoulRefiningExpendNum,omitempty" bson:"heroSoulRefiningExpendNum"`                 //英雄战魂洗练消耗数量
	HeroSoulMaxQuickAddRarity         Rarity                   `json:"heroSoulMaxQuickAddRarity,omitempty" bson:"heroSoulMaxQuickAddRarity"`                 //英雄战魂快速添加最高品质
	HeroSoulMaxCanComposeRarity       Rarity                   `json:"heroSoulMaxCanComposeRarity,omitempty" bson:"heroSoulMaxCanComposeRarity"`             //最高可用于合成的英雄战魂品质
	HeroSoulSuperCount                int32                    `json:"heroSoulSuperCount,omitempty" bson:"heroSoulSuperCount"`                               //逆天属性条目数量
	HeroSoulSuperBonus                float32                  `json:"heroSoulSuperBonus,omitempty" bson:"heroSoulSuperBonus"`                               //逆天属性本源加成
	HeroSoulStarNumBonus              []*IdFloatVal            `json:"heroSoulStarNumBonus,omitempty" bson:"heroSoulStarNumBonus"`                           //本源星级属性加成
	HeroSoulMaxRecastBonus            float32                  `json:"heroSoulMaxRecastBonus,omitempty" bson:"heroSoulMaxRecastBonus"`                       //本源重铸最大加成
	HeroSoulRecastPerBonus            float32                  `json:"heroSoulRecastPerBonus,omitempty" bson:"heroSoulRecastPerBonus"`                       //本源每次重铸加成
	HeroSoulRecastCost                []*IdCount               `json:"heroSoulRecastCost,omitempty" bson:"heroSoulRecastCost"`                               //本源重铸消耗
	HeroSoulLevelUp                   Rarity                   `json:"heroSoulLevelUp,omitempty" bson:"heroSoulLevelUp"`                                     //战魂升绝品质
	HeroSoulRecastCnt                 int32                    `json:"heroSoulRecastCnt,omitempty" bson:"heroSoulRecastCnt"`                                 //可以重铸的属性条目数
	ArtifactBreakExpendCnt            int32                    `json:"artifactBreakExpendCnt,omitempty" bson:"artifactBreakExpendCnt"`                       //神器突破消耗数量
	DecomposeChipCnt                  []int32                  `json:"decomposeChipCnt,omitempty" bson:"decomposeChipCnt"`                                   //神器分解碎片数量(根据品质)
	AstrologyExchangeScore            int32                    `json:"astrologyExchangeScore,omitempty" bson:"astrologyExchangeScore"`                       //神器占星单次兑奖分值消耗
	AstrologyExchangeAward            *IdCount                 `json:"astrologyExchangeAward,omitempty" bson:"astrologyExchangeAward"`                       //神器占星兑奖给予奖励
	AstrologyFreeCnt                  int32                    `json:"astrologyFreeCnt,omitempty" bson:"astrologyFreeCnt"`                                   //神器占星免费次数
	AstrologyMaxSave                  int32                    `json:"astrologyMaxSave,omitempty" bson:"astrologyMaxSave"`                                   //占星记录最大存储条数
	AstrologyMaxPage                  int32                    `json:"astrologyMaxPage,omitempty" bson:"astrologyMaxPage"`                                   //占星记录最大页数
	AstrologyPageLogCount             int32                    `json:"astrologyPageLogCount,omitempty" bson:"astrologyPageLogCount"`                         //占星记录每页记录条数
	AstrologyMaxAutoCount             int32                    `json:"astrologyMaxAutoCount,omitempty" bson:"astrologyMaxAutoCount"`                         //每次自动占星停止次数
	GoodsBreakArtifactRarity          Rarity                   `json:"goodsBreakArtifactRarity,omitempty" bson:"goodsBreakArtifactRarity"`                   //可以用道具突破的神器品质
	MaxArtifactRarity                 Rarity                   `json:"maxArtifactRarity,omitempty" bson:"maxArtifactRarity"`                                 //神器最高品质
	HeroQuickStarupCount              int32                    `json:"heroQuickStarupCount,omitempty" bson:"heroQuickStarupCount"`                           //大于等于X（设置值）个英雄可升星
	FirstMoveAllStart                 float32                  `json:"firstMoveAllStart,omitempty" bson:"firstMoveAllStart"`                                 //先手值进行对拳
	FirstMoveJudeLimit                float32                  `json:"firstMoveJudeLimit,omitempty" bson:"firstMoveJudeLimit"`                               //先手值判定阈值
	AwakenRarity                      Rarity                   `json:"awakenRarity,omitempty" bson:"awakenRarity"`                                           //觉醒需要的英雄品质
	AwakenStarRequire                 []int32                  `json:"awakenStarRequire,omitempty" bson:"awakenStarRequire"`                                 //觉醒档位星级要求
	AwakenBaseSkillProperty1          []*AttrValue             `json:"awakenBaseSkillProperty1,omitempty" bson:"awakenBaseSkillProperty1"`                   //档位1基础技能属性
	AwakenBaseSkillProperty2          []*AttrValue             `json:"awakenBaseSkillProperty2,omitempty" bson:"awakenBaseSkillProperty2"`                   //档位2基础技能属性
	AwakenBaseSkillProperty3          []*AttrValue             `json:"awakenBaseSkillProperty3,omitempty" bson:"awakenBaseSkillProperty3"`                   //档位3基础技能属性
	AwakenBaseSkillProperty4          []*AttrValue             `json:"awakenBaseSkillProperty4,omitempty" bson:"awakenBaseSkillProperty4"`                   //档位4基础技能属性
	AwakenBaseSkillProperty5          []*AttrValue             `json:"awakenBaseSkillProperty5,omitempty" bson:"awakenBaseSkillProperty5"`                   //档位5基础技能属性
	AwakenBaseSkillProperty6          []*AttrValue             `json:"awakenBaseSkillProperty6,omitempty" bson:"awakenBaseSkillProperty6"`                   //档位6基础技能属性
	AwakenBaseSkillProperty7          []*AttrValue             `json:"awakenBaseSkillProperty7,omitempty" bson:"awakenBaseSkillProperty7"`                   //档位7基础技能属性
	AwakenBaseSkillProperty8          []*AttrValue             `json:"awakenBaseSkillProperty8,omitempty" bson:"awakenBaseSkillProperty8"`                   //档位8基础技能属性
	AwakenStoneId                     map[Career]int32         `json:"awakenStoneId,omitempty" bson:"awakenStoneId"`                                         //每类职业对应的觉醒石
	HeroSoulSuitNum                   []int32                  `json:"heroSoulSuitNum,omitempty" bson:"heroSoulSuitNum"`                                     //战魂套装数量
	FirstRechargeToRightsItems        map[int32]int32          `json:"firstRechargeToRightsItems,omitempty" bson:"firstRechargeToRightsItems"`               //首充八神商品id对应权益道具id
	ArtifactModuleUnlockStar          []int32                  `json:"artifactModuleUnlockStar,omitempty" bson:"artifactModuleUnlockStar"`                   //神器模块解锁所需星级
	ArtifactModuleRefine              *IdCount                 `json:"artifactModuleRefine,omitempty" bson:"artifactModuleRefine"`                           //重构消耗道具id
	ArtifactModuleLockCosts           []*IdCount               `json:"artifactModuleLockCosts,omitempty" bson:"artifactModuleLockCosts"`                     //重构锁定属性消耗道具
	ArtifactModuleRefineUnlockCost    []*IdCount               `json:"artifactModuleRefineUnlockCost,omitempty" bson:"artifactModuleRefineUnlockCost"`       //重构属性解锁消耗
	ArtifactModuleRarityWeights       []*ArtifactRaritySection `json:"artifactModuleRarityWeights,omitempty" bson:"artifactModuleRarityWeights"`             //重构词条品质对应权重
	ArtifactModuleBurstProb           float32                  `json:"artifactModuleBurstProb,omitempty" bson:"artifactModuleBurstProb"`                     //重构爆发概率
	ArtifactModuleBurstAttrs          []*RandomAttrValue       `json:"artifactModuleBurstAttrs,omitempty" bson:"artifactModuleBurstAttrs"`                   //重构爆发属性类型及属性区间
	ArtifactModuleBurstMinTimes       int32                    `json:"artifactModuleBurstMinTimes,omitempty" bson:"artifactModuleBurstMinTimes"`             //重构爆发需要最少重构次数
	ArtifactModuleBurstMaxTimes       int32                    `json:"artifactModuleBurstMaxTimes,omitempty" bson:"artifactModuleBurstMaxTimes"`             //重构爆发保底重构次数
	ArtifactModuleSwitchAttrTypes     [][]AttrType             `json:"artifactModuleSwitchAttrTypes,omitempty" bson:"artifactModuleSwitchAttrTypes"`         //重构属性切换库
	ArtifactModuleSwitchCD            int32                    `json:"artifactModuleSwitchCD,omitempty" bson:"artifactModuleSwitchCD"`                       //重构属性切换cd（秒）
	ArtifactModuleSwitchCost          *IdCount                 `json:"artifactModuleSwitchCost,omitempty" bson:"artifactModuleSwitchCost"`                   //重构属性切换消耗
	ArtifactModuleAttrGroupLimitTimes map[int32]int32          `json:"artifactModuleAttrGroupLimitTimes,omitempty" bson:"artifactModuleAttrGroupLimitTimes"` //重构属性组限制出现次数
	HeroSoulPushAttrCntLimit          int32                    `json:"heroSoulPushAttrCntLimit,omitempty" bson:"heroSoulPushAttrCntLimit"`                   //八门属性多少条才广播(大于等于)
	DepromotionInitRarity             Rarity                   `json:"depromotionInitRarity,omitempty" bson:"depromotionInitRarity"`                         //降品最低品质
	HeroSoulSpecialCanRecast          bool                     `json:"heroSoulSpecialCanRecast,omitempty" bson:"heroSoulSpecialCanRecast"`                   //特殊八本是否可以重铸
}

// RoleBaseConfig.json属性表集合
type RoleBaseConfigEntries struct {
	Rows   map[int32]*RoleBaseConfigEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*RoleBaseConfigEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("RoleBaseConfig.json", (*RoleBaseConfigEntries)(nil))
}

func (e *RoleBaseConfigEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	roleBaseConfigEntries = &RoleBaseConfigEntries{
		Rows:   make(map[int32]*RoleBaseConfigEntry),
		Arrays: make([]*RoleBaseConfigEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &RoleBaseConfigEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return roleBaseConfigEntries, err
		}

		roleBaseConfigEntries.Rows[entry.Test] = entry
		roleBaseConfigEntries.Arrays = append(roleBaseConfigEntries.Arrays, entry)
	}

	sort.Slice(roleBaseConfigEntries.Arrays, func(i, j int) bool {
		return roleBaseConfigEntries.Arrays[i].Test < roleBaseConfigEntries.Arrays[j].Test
	})

	excel.AddArrayGetter("roleBaseConfig.json", roleBaseConfigEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return roleBaseConfigEntries, nil

}

func (e *RoleBaseConfigEntries) GetArrays() any {
	return roleBaseConfigEntries.Arrays
}

func GetRoleBaseConfigEntry() *RoleBaseConfigEntry {
	entry := roleBaseConfigEntries.Rows[0]
	return entry
}

func GetRoleBaseConfigSize() int {
	return int(len(roleBaseConfigEntries.Rows))
}

func GetRoleBaseConfigRows() map[int32]*RoleBaseConfigEntry {
	return roleBaseConfigEntries.Rows
}

func GetRoleBaseConfigArrays() []*RoleBaseConfigEntry {
	return roleBaseConfigEntries.Arrays
}
