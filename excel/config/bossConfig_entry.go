package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var bossConfigEntries *BossConfigEntries //BossConfig.json全局变量

// 属性表 世界BOSS-拳皇/bossConfig
type BossConfigEntry struct {
	Test           int32   `json:"test,omitempty" bson:"test"`                     //##测试数据
	FirstBossTime  string  `json:"firstBossTime,omitempty" bson:"firstBossTime"`   //第一期BOSS开始时间
	BossDuration   int32   `json:"bossDuration,omitempty" bson:"bossDuration"`     //每期Boss持续时间（秒）
	FreeTimes      int32   `json:"freeTimes,omitempty" bson:"freeTimes"`           //每天免费次数
	BuyTimes       int32   `json:"buyTimes,omitempty" bson:"buyTimes"`             //每天购买次数
	CostId         int32   `json:"costId,omitempty" bson:"costId"`                 //购买次数所需道具ID
	CostCnt        []int32 `json:"costCnt,omitempty" bson:"costCnt"`               //购买消耗(超过取最后)
	GroupRotation  []int32 `json:"groupRotation,omitempty" bson:"groupRotation"`   //组轮替
	SettleDuration int32   `json:"settleDuration,omitempty" bson:"settleDuration"` //轮提前多少秒进结算时间
}

// BossConfig.json属性表集合
type BossConfigEntries struct {
	Rows   map[int32]*BossConfigEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*BossConfigEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("BossConfig.json", (*BossConfigEntries)(nil))
}

func (e *BossConfigEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	bossConfigEntries = &BossConfigEntries{
		Rows:   make(map[int32]*BossConfigEntry),
		Arrays: make([]*BossConfigEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &BossConfigEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return bossConfigEntries, err
		}

		bossConfigEntries.Rows[entry.Test] = entry
		bossConfigEntries.Arrays = append(bossConfigEntries.Arrays, entry)
	}

	sort.Slice(bossConfigEntries.Arrays, func(i, j int) bool {
		return bossConfigEntries.Arrays[i].Test < bossConfigEntries.Arrays[j].Test
	})

	excel.AddArrayGetter("bossConfig.json", bossConfigEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return bossConfigEntries, nil

}

func (e *BossConfigEntries) GetArrays() any {
	return bossConfigEntries.Arrays
}

func GetBossConfigEntry() *BossConfigEntry {
	entry := bossConfigEntries.Rows[0]
	return entry
}

func GetBossConfigSize() int {
	return int(len(bossConfigEntries.Rows))
}

func GetBossConfigRows() map[int32]*BossConfigEntry {
	return bossConfigEntries.Rows
}

func GetBossConfigArrays() []*BossConfigEntry {
	return bossConfigEntries.Arrays
}
