package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var towerConfigEntries *TowerConfigEntries //TowerConfig.json全局变量

// 属性表 爬塔-拳皇/TowerConfig
type TowerConfigEntry struct {
	Test                  int32 `json:"test,omitempty" bson:"test"`                                   //##测试数据
	TicketRestoreMaxCount int32 `json:"ticketRestoreMaxCount,omitempty" bson:"ticketRestoreMaxCount"` //爬塔挑战券最多恢复数量
	TicketRestoreInterval int32 `json:"ticketRestoreInterval,omitempty" bson:"ticketRestoreInterval"` //爬塔挑战券恢复间隔(分钟)
	FloorPerPhase         int32 `json:"floorPerPhase,omitempty" bson:"floorPerPhase"`                 //每几层领奖
	RankFloorLimit        int32 `json:"rankFloorLimit,omitempty" bson:"rankFloorLimit"`               //爬塔通过几层后进入排行榜
}

// TowerConfig.json属性表集合
type TowerConfigEntries struct {
	Rows   map[int32]*TowerConfigEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*TowerConfigEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("TowerConfig.json", (*TowerConfigEntries)(nil))
}

func (e *TowerConfigEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	towerConfigEntries = &TowerConfigEntries{
		Rows:   make(map[int32]*TowerConfigEntry),
		Arrays: make([]*TowerConfigEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &TowerConfigEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return towerConfigEntries, err
		}

		towerConfigEntries.Rows[entry.Test] = entry
		towerConfigEntries.Arrays = append(towerConfigEntries.Arrays, entry)
	}

	sort.Slice(towerConfigEntries.Arrays, func(i, j int) bool {
		return towerConfigEntries.Arrays[i].Test < towerConfigEntries.Arrays[j].Test
	})

	excel.AddArrayGetter("towerConfig.json", towerConfigEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return towerConfigEntries, nil

}

func (e *TowerConfigEntries) GetArrays() any {
	return towerConfigEntries.Arrays
}

func GetTowerConfigEntry() *TowerConfigEntry {
	entry := towerConfigEntries.Rows[0]
	return entry
}

func GetTowerConfigSize() int {
	return int(len(towerConfigEntries.Rows))
}

func GetTowerConfigRows() map[int32]*TowerConfigEntry {
	return towerConfigEntries.Rows
}

func GetTowerConfigArrays() []*TowerConfigEntry {
	return towerConfigEntries.Arrays
}
