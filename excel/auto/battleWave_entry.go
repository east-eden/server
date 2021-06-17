package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var battleWaveEntries *BattleWaveEntries //BattleWave.csv全局变量

// BattleWave.csv属性表
type BattleWaveEntry struct {
	Id        int32             `json:"Id,omitempty"`        // 主键
	MonsterID []int32           `json:"MonsterID,omitempty"` //怪物组ID
	PositionX []decimal.Decimal `json:"PositionX,omitempty"` //单位x坐标
	PositionZ []decimal.Decimal `json:"PositionZ,omitempty"` //单位z坐标
	InitalCom []decimal.Decimal `json:"InitalCom,omitempty"` //单位z坐标
	Rotation  []decimal.Decimal `json:"Rotation,omitempty"`  //单位旋转值
}

// BattleWave.csv属性表集合
type BattleWaveEntries struct {
	Rows map[int32]*BattleWaveEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("BattleWave.csv", (*BattleWaveEntries)(nil))
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
