package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var levelEntries *LevelEntries //Level.json全局变量

// 属性表 关卡表-拳皇/Level
type LevelEntry struct {
	Id                int32           `json:"id,omitempty" bson:"id"`                               //##id
	LevelName         string          `json:"levelName,omitempty" bson:"levelName"`                 //关卡名字
	LevelType         LevelType       `json:"levelType,omitempty" bson:"levelType"`                 //关卡类型，1=箱庭2=泰坦7=站桩9=采集
	Chapter           int32           `json:"chapter,omitempty" bson:"chapter"`                     //章序号
	Stage             int32           `json:"stage,omitempty" bson:"stage"`                         //节序号
	LevelLimit        int32           `json:"levelLimit,omitempty" bson:"levelLimit"`               //等级限制
	ServerBattle      bool            `json:"serverBattle,omitempty" bson:"serverBattle"`           //是否服务器战斗
	StaminaEffect     bool            `json:"staminaEffect,omitempty" bson:"staminaEffect"`         //是否受体力影响
	CostStamina       int32           `json:"costStamina,omitempty" bson:"costStamina"`             //胜利消耗体力
	GoalAward         []*IdCount      `json:"goalAward,omitempty" bson:"goalAward"`                 //目标奖励
	StageComfortLoot  int32           `json:"stageComfortLoot,omitempty" bson:"stageComfortLoot"`   //关卡失败掉落id
	StageLoot         int32           `json:"stageLoot,omitempty" bson:"stageLoot"`                 //关卡掉落ID
	CoinRewards       decimal.Decimal `json:"coinRewards,omitempty" bson:"coinRewards"`             //通关金币奖励
	EnemyPower        []int32         `json:"enemyPower,omitempty" bson:"enemyPower"`               //敌人战斗力
	Monsters          [][]int32       `json:"monsters,omitempty" bson:"monsters"`                   //怪物ids
	MonsterStars      [][]int32       `json:"monsterStars,omitempty" bson:"monsterStars"`           //怪物星级[]解锁技能用
	MonsterStrength   [][]int32       `json:"monsterStrength,omitempty" bson:"monsterStrength"`     //怪物强度类型
	StrengthType      [][]int32       `json:"strengthType,omitempty" bson:"strengthType"`           //怪物强度
	StarParam         float32         `json:"starParam,omitempty" bson:"starParam"`                 //星级系数
	MonsterAPStrength [][]int32       `json:"monsterAPStrength,omitempty" bson:"monsterAPStrength"` //怪物AP强度类型
	StrengthAPType    [][]int32       `json:"strengthAPType,omitempty" bson:"strengthAPType"`       //怪物AP强度对应属性表参数：strenA
	FirstMoveVal      []float32       `json:"firstMoveVal,omitempty" bson:"firstMoveVal"`           //关卡怪物先手值
}

// Level.json属性表集合
type LevelEntries struct {
	Rows   map[int32]*LevelEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*LevelEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Level.json", (*LevelEntries)(nil))
}

func (e *LevelEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	levelEntries = &LevelEntries{
		Rows:   make(map[int32]*LevelEntry),
		Arrays: make([]*LevelEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &LevelEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return levelEntries, err
		}

		levelEntries.Rows[entry.Id] = entry
		levelEntries.Arrays = append(levelEntries.Arrays, entry)
	}

	sort.Slice(levelEntries.Arrays, func(i, j int) bool {
		return levelEntries.Arrays[i].Id < levelEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("level.json", levelEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return levelEntries, nil

}

func (e *LevelEntries) GetArrays() any {
	return levelEntries.Arrays
}

func GetLevelEntry(id int32) (*LevelEntry, bool) {
	entry, ok := levelEntries.Rows[id]
	return entry, ok
}

func GetLevelSize() int {
	return int(len(levelEntries.Rows))
}

func GetLevelRows() map[int32]*LevelEntry {
	return levelEntries.Rows
}

func GetLevelArrays() []*LevelEntry {
	return levelEntries.Arrays
}
