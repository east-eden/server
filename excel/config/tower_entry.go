package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var towerEntries *TowerEntries //Tower.json全局变量

// 属性表 爬塔-拳皇/Tower
type TowerEntry struct {
	Id         int32      `json:"id,omitempty" bson:"id"`                 //##塔层数
	LevelName  string     `json:"levelName,omitempty" bson:"levelName"`   //关卡名字
	FloorAward []*IdCount `json:"floorAward,omitempty" bson:"floorAward"` //通关层数奖励
	PhaseAward []*IdCount `json:"phaseAward,omitempty" bson:"phaseAward"` //阶段奖励
	LevelId    int32      `json:"levelId,omitempty" bson:"levelId"`       //对应关卡表id
}

// Tower.json属性表集合
type TowerEntries struct {
	Rows   map[int32]*TowerEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*TowerEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Tower.json", (*TowerEntries)(nil))
}

func (e *TowerEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	towerEntries = &TowerEntries{
		Rows:   make(map[int32]*TowerEntry),
		Arrays: make([]*TowerEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &TowerEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return towerEntries, err
		}

		towerEntries.Rows[entry.Id] = entry
		towerEntries.Arrays = append(towerEntries.Arrays, entry)
	}

	sort.Slice(towerEntries.Arrays, func(i, j int) bool {
		return towerEntries.Arrays[i].Id < towerEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("tower.json", towerEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return towerEntries, nil

}

func (e *TowerEntries) GetArrays() any {
	return towerEntries.Arrays
}

func GetTowerEntry(id int32) (*TowerEntry, bool) {
	entry, ok := towerEntries.Rows[id]
	return entry, ok
}

func GetTowerSize() int {
	return int(len(towerEntries.Rows))
}

func GetTowerRows() map[int32]*TowerEntry {
	return towerEntries.Rows
}

func GetTowerArrays() []*TowerEntry {
	return towerEntries.Arrays
}
