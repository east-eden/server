package auto

import (
	"fmt"
	"strings"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

var towerEntries *TowerEntries //Tower.xlsx全局变量

// Tower.xlsx属性表
type TowerEntry struct {
	Id            int32 `json:"Id,omitempty"`            // 主键
	Floor         int32 `json:"Floor,omitempty"`         // 多主键之一
	LevelLimit    int32 `json:"LevelLimit,omitempty"`    //队伍等级限制
	FirstRewardId int32 `json:"FirstRewardId,omitempty"` //首通奖励Id
	DailyRewardId int32 `json:"DailyRewardId,omitempty"` //每日结算奖励id
	SceneId       int32 `json:"SceneId,omitempty"`       //场景id
}

// Tower.xlsx属性表集合
type TowerEntries struct {
	Rows map[string]*TowerEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Tower.xlsx", (*TowerEntries)(nil))
}

func (e *TowerEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	towerEntries = &TowerEntries{
		Rows: make(map[string]*TowerEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &TowerEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		key := fmt.Sprintf("%d+%d", entry.Id, entry.Floor)
		towerEntries.Rows[key] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetTowerEntry(keys ...int32) (*TowerEntry, bool) {
	keyName := make([]string, 0, len(keys))
	for _, key := range keys {
		keyName = append(keyName, cast.ToString(key))
	}

	finalKey := strings.Join(keyName, "+")
	entry, ok := towerEntries.Rows[finalKey]
	return entry, ok
}

func GetTowerSize() int32 {
	return int32(len(towerEntries.Rows))
}

func GetTowerRows() map[string]*TowerEntry {
	return towerEntries.Rows
}
