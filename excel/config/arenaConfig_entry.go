package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var arenaConfigEntries *ArenaConfigEntries //ArenaConfig.json全局变量

// 属性表 竞技场-拳皇/arenaConfig
type ArenaConfigEntry struct {
	Test                            int32   `json:"test,omitempty" bson:"test"`                                                       //##测试数据
	DefaultScore                    int32   `json:"defaultScore,omitempty" bson:"defaultScore"`                                       //竞技场默认积分
	RefreshTargetInterval           int32   `json:"refreshTargetInterval,omitempty" bson:"refreshTargetInterval"`                     //竞技场刷新对手间隔时间（秒）
	DefaultTargetNum                int32   `json:"defaultTargetNum,omitempty" bson:"defaultTargetNum"`                               //竞技场默认对手数
	DefaultDefFormation             []int32 `json:"defaultDefFormation,omitempty" bson:"defaultDefFormation"`                         //竞技场默认防守阵容
	RefreshCost                     []int32 `json:"refreshCost,omitempty" bson:"refreshCost"`                                         //竞技场刷新消耗
	ChallengeCost                   []int32 `json:"challengeCost,omitempty" bson:"challengeCost"`                                     //竞技场挑战消耗
	WeeklyLogTime                   string  `json:"weeklyLogTime,omitempty" bson:"weeklyLogTime"`                                     //竞技场周排行记录日志时间
	WeeklyRankResetTime             string  `json:"weeklyRankResetTime,omitempty" bson:"weeklyRankResetTime"`                         //竞技场周排行重置时间
	DailySettleTimeBegin            string  `json:"dailySettleTimeBegin,omitempty" bson:"dailySettleTimeBegin"`                       //竞技场每日结算开始时间
	DailySettleTimeEnd              string  `json:"dailySettleTimeEnd,omitempty" bson:"dailySettleTimeEnd"`                           //竞技场每日结算结束时间
	ScoreSection                    []int32 `json:"scoreSection,omitempty" bson:"scoreSection"`                                       //竞技场金钱奖励积分区域
	ScoreSectionRewardLoots         []int32 `json:"scoreSectionRewardLoots,omitempty" bson:"scoreSectionRewardLoots"`                 //竞技场战斗胜利掉落id（<1200|1200-1499|1500-1999|>2000)
	CrushScoreDisparity             int32   `json:"crushScoreDisparity,omitempty" bson:"crushScoreDisparity"`                         //竞技场积分碾压差值
	CrushPowerDisparity             float32 `json:"crushPowerDisparity,omitempty" bson:"crushPowerDisparity"`                         //竞技场形成战力碾压低于自身的比例
	ArenaShieldDuration             int32   `json:"arenaShieldDuration,omitempty" bson:"arenaShieldDuration"`                         //竞技场保护时间（秒）
	NoRobotScore                    int32   `json:"noRobotScore,omitempty" bson:"noRobotScore"`                                       //不补充机器人的竞技场分数
	ArenaWeekAwardsBXWeeks          int32   `json:"arenaWeekAwardsBXWeeks,omitempty" bson:"arenaWeekAwardsBXWeeks"`                   //开服前x周结算奖励B
	ArenaWeekAwardsBWechatServerIds []int32 `json:"arenaWeekAwardsBWechatServerIds,omitempty" bson:"arenaWeekAwardsBWechatServerIds"` //竞技场前x周结算奖励B生效微信服务器id
	ArenaWeekAwardsBOtherServerIds  []int32 `json:"arenaWeekAwardsBOtherServerIds,omitempty" bson:"arenaWeekAwardsBOtherServerIds"`   //竞技场前x周结算奖励B生效其他平台服务器id
}

// ArenaConfig.json属性表集合
type ArenaConfigEntries struct {
	Rows   map[int32]*ArenaConfigEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ArenaConfigEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("ArenaConfig.json", (*ArenaConfigEntries)(nil))
}

func (e *ArenaConfigEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	arenaConfigEntries = &ArenaConfigEntries{
		Rows:   make(map[int32]*ArenaConfigEntry),
		Arrays: make([]*ArenaConfigEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ArenaConfigEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return arenaConfigEntries, err
		}

		arenaConfigEntries.Rows[entry.Test] = entry
		arenaConfigEntries.Arrays = append(arenaConfigEntries.Arrays, entry)
	}

	sort.Slice(arenaConfigEntries.Arrays, func(i, j int) bool {
		return arenaConfigEntries.Arrays[i].Test < arenaConfigEntries.Arrays[j].Test
	})

	excel.AddArrayGetter("arenaConfig.json", arenaConfigEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return arenaConfigEntries, nil

}

func (e *ArenaConfigEntries) GetArrays() any {
	return arenaConfigEntries.Arrays
}

func GetArenaConfigEntry() *ArenaConfigEntry {
	entry := arenaConfigEntries.Rows[0]
	return entry
}

func GetArenaConfigSize() int {
	return int(len(arenaConfigEntries.Rows))
}

func GetArenaConfigRows() map[int32]*ArenaConfigEntry {
	return arenaConfigEntries.Rows
}

func GetArenaConfigArrays() []*ArenaConfigEntry {
	return arenaConfigEntries.Arrays
}
