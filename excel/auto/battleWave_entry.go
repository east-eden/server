package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var battleWaveEntries *BattleWaveEntries //BattleWave.xlsx全局变量

// BattleWave.xlsx属性表
type BattleWaveEntry struct {
	Id          int32 `json:"Id,omitempty"`          // 主键
	SceneId     int32 `json:"SceneId,omitempty"`     //场景id
	BattleView  int32 `json:"BattleView,omitempty"`  //战斗区域ID
	UnitGroupId int32 `json:"UnitGroupId,omitempty"` //怪物组id
}

// BattleWave.xlsx属性表集合
type BattleWaveEntries struct {
	Rows map[int32]*BattleWaveEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("BattleWave.xlsx", (*BattleWaveEntries)(nil))
}

func (e *BattleWaveEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	battleWaveEntries = &BattleWaveEntries{
		Rows: make(map[int32]*BattleWaveEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &BattleWaveEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		battleWaveEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetBattleWaveEntry(id int32) (*BattleWaveEntry, bool) {
	entry, ok := battleWaveEntries.Rows[id]
	return entry, ok
}

func GetBattleWaveSize() int32 {
	return int32(len(battleWaveEntries.Rows))
}

func GetBattleWaveRows() map[int32]*BattleWaveEntry {
	return battleWaveEntries.Rows
}
