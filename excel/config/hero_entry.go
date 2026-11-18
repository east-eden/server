package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var heroEntries *HeroEntries //Hero.json全局变量

// 属性表 英雄表-拳皇/Hero
type HeroEntry struct {
	Id                   int32               `json:"id,omitempty" bson:"id"`                                     //##id
	Name                 string              `json:"name,omitempty" bson:"name"`                                 //名称（不填则用种族随机名称+随机天赋前缀）
	ImgName              string              `json:"imgName,omitempty" bson:"imgName"`                           //名字的图片资源
	Rarity               Rarity              `json:"rarity,omitempty" bson:"rarity"`                             //品质（查询基础配置表的道具类型）
	Career               Career              `json:"career,omitempty" bson:"career"`                             //职业
	RoleType             BattleRoleType      `json:"roleType,omitempty" bson:"roleType"`                         //角色类型
	ChipId               int32               `json:"chipId,omitempty" bson:"chipId"`                             //碎片id
	PropertyIds          [][]int32           `json:"propertyIds,omitempty" bson:"propertyIds"`                   //属性id（索引角色的战斗属性表）
	DreamHero            bool                `json:"dreamHero,omitempty" bson:"dreamHero"`                       //是否为可许愿英雄
	WishHero             bool                `json:"wishHero,omitempty" bson:"wishHero"`                         //是否为心愿英雄
	EquipBreakAttId      int32               `json:"equipBreakAttId,omitempty" bson:"equipBreakAttId"`           //装备突破属性id
	SoulActiveId         int32               `json:"soulActiveId,omitempty" bson:"soulActiveId"`                 //战魂激活属性 id
	FightOrder           int32               `json:"fightOrder,omitempty" bson:"fightOrder"`                     //站位优先级越小越考前
	AttackDistance       float32             `json:"attackDistance,omitempty" bson:"attackDistance"`             //攻击距离
	SkillIds             []int32             `json:"skillIds,omitempty" bson:"skillIds"`                         //技能列表id1|id2
	AwakenSkills         [][]int32           `json:"awakenSkills,omitempty" bson:"awakenSkills"`                 //觉醒技能列表
	SoulSkill            []*SoulSkillRequire `json:"soulSkill,omitempty" bson:"soulSkill"`                       //战魂绝技
	ArtifactModuleSkills []int32             `json:"artifactModuleSkills,omitempty" bson:"artifactModuleSkills"` //
	HideInit             bool                `json:"hideInit,omitempty" bson:"hideInit"`                         //未获得时羁绊不显示
	IsTest               bool                `json:"isTest,omitempty" bson:"isTest"`                             //状态说明1为测试角色，不填为正式角色
}

// Hero.json属性表集合
type HeroEntries struct {
	Rows   map[int32]*HeroEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*HeroEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Hero.json", (*HeroEntries)(nil))
}

func (e *HeroEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	heroEntries = &HeroEntries{
		Rows:   make(map[int32]*HeroEntry),
		Arrays: make([]*HeroEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &HeroEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return heroEntries, err
		}

		heroEntries.Rows[entry.Id] = entry
		heroEntries.Arrays = append(heroEntries.Arrays, entry)
	}

	sort.Slice(heroEntries.Arrays, func(i, j int) bool {
		return heroEntries.Arrays[i].Id < heroEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("hero.json", heroEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return heroEntries, nil

}

func (e *HeroEntries) GetArrays() any {
	return heroEntries.Arrays
}

func GetHeroEntry(id int32) (*HeroEntry, bool) {
	entry, ok := heroEntries.Rows[id]
	return entry, ok
}

func GetHeroSize() int {
	return int(len(heroEntries.Rows))
}

func GetHeroRows() map[int32]*HeroEntry {
	return heroEntries.Rows
}

func GetHeroArrays() []*HeroEntry {
	return heroEntries.Arrays
}
