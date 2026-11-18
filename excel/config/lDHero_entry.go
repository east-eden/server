package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var lDHeroEntries *LDHeroEntries //LDHero.json全局变量

// 属性表 LD英雄表-拳皇/LDHero
type LDHeroEntry struct {
	Id                   int32              `json:"id,omitempty" bson:"id"`                                     //##id
	Name                 string             `json:"name,omitempty" bson:"name"`                                 //名称（不填则用种族随机名称+随机天赋前缀）
	Rarity               LDRarity           `json:"rarity,omitempty" bson:"rarity"`                             //品质（查询基础配置表的道具类型）
	PropertyId           int32              `json:"propertyId,omitempty" bson:"propertyId"`                     //属性id（索引角色、教官表的战斗属性表）
	StackMax             int32              `json:"stackMax,omitempty" bson:"stackMax"`                         //英雄在同一个格子内可堆叠数量
	Level                int32              `json:"level,omitempty" bson:"level"`                               //英雄初始等级
	LevelMax             int32              `json:"levelMax,omitempty" bson:"levelMax"`                         //英雄等级上限
	DailyFormationWeight int32              `json:"dailyFormationWeight,omitempty" bson:"dailyFormationWeight"` //每日阵容随机权重
	SummonWeight         int32              `json:"summonWeight,omitempty" bson:"summonWeight"`                 //召唤随机权重
	DecWeightHeroNum     int32              `json:"decWeightHeroNum,omitempty" bson:"decWeightHeroNum"`         //场上英雄数达到多少时触发降权(阈值B)
	WeightPerDec         int32              `json:"weightPerDec,omitempty" bson:"weightPerDec"`                 //每次降权下降权重值（差值C）
	MinWeight            int32              `json:"minWeight,omitempty" bson:"minWeight"`                       //降权后最少权重(最小值D)
	SellItem             *IdCount           `json:"sellItem,omitempty" bson:"sellItem"`                         //出售获得道具
	PackGridWeight       float32            `json:"packGridWeight,omitempty" bson:"packGridWeight"`             //整理时英雄所处格子加权系数
	SkillIds             []*RoleSkillUnlock `json:"skillIds,omitempty" bson:"skillIds"`                         //技能列表id1,lv1|id2,lv2
	TalentIds            []*RoleSkillUnlock `json:"talentIds,omitempty" bson:"talentIds"`                       //天赋列表id1,lv1|id2,lv2
	WordsOut             []int32            `json:"wordsOut,omitempty" bson:"wordsOut"`                         //词条战斗外
	InitCnt              int32              `json:"initCnt,omitempty" bson:"initCnt"`                           //初始数量
	BuyCost              []*IdCount         `json:"buyCost,omitempty" bson:"buyCost"`                           //购买消耗
}

// LDHero.json属性表集合
type LDHeroEntries struct {
	Rows   map[int32]*LDHeroEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*LDHeroEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("LDHero.json", (*LDHeroEntries)(nil))
}

func (e *LDHeroEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	lDHeroEntries = &LDHeroEntries{
		Rows:   make(map[int32]*LDHeroEntry),
		Arrays: make([]*LDHeroEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &LDHeroEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return lDHeroEntries, err
		}

		lDHeroEntries.Rows[entry.Id] = entry
		lDHeroEntries.Arrays = append(lDHeroEntries.Arrays, entry)
	}

	sort.Slice(lDHeroEntries.Arrays, func(i, j int) bool {
		return lDHeroEntries.Arrays[i].Id < lDHeroEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("lDHero.json", lDHeroEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return lDHeroEntries, nil

}

func (e *LDHeroEntries) GetArrays() any {
	return lDHeroEntries.Arrays
}

func GetLDHeroEntry(id int32) (*LDHeroEntry, bool) {
	entry, ok := lDHeroEntries.Rows[id]
	return entry, ok
}

func GetLDHeroSize() int {
	return int(len(lDHeroEntries.Rows))
}

func GetLDHeroRows() map[int32]*LDHeroEntry {
	return lDHeroEntries.Rows
}

func GetLDHeroArrays() []*LDHeroEntry {
	return lDHeroEntries.Arrays
}
