package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var battleLayoutEntries *BattleLayoutEntries //BattleLayout.csv全局变量

// BattleLayout.csv属性表
type BattleLayoutEntry struct {
	Id        int32             `json:"Id,omitempty"`        // 主键
	PositionX []decimal.Decimal `json:"PositionX,omitempty"` //友军单位x坐标
	PositionZ []decimal.Decimal `json:"PositionZ,omitempty"` //友军单位z坐标
	Rotation  []decimal.Decimal `json:"Rotation,omitempty"`  //友军单位旋转值
}

// BattleLayout.csv属性表集合
type BattleLayoutEntries struct {
	Rows map[int32]*BattleLayoutEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("BattleLayout.csv", (*BattleLayoutEntries)(nil))
}

func (e *BattleLayoutEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	battleLayoutEntries = &BattleLayoutEntries{
		Rows: make(map[int32]*BattleLayoutEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &BattleLayoutEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		battleLayoutEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetBattleLayoutEntry(id int32) (*BattleLayoutEntry, bool) {
	entry, ok := battleLayoutEntries.Rows[id]
	return entry, ok
}

func GetBattleLayoutSize() int32 {
	return int32(len(battleLayoutEntries.Rows))
}

func GetBattleLayoutRows() map[int32]*BattleLayoutEntry {
	return battleLayoutEntries.Rows
}
