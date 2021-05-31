package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var battleViewEntries *BattleViewEntries //BattleView.xlsx全局变量

// BattleView.xlsx属性表
type BattleViewEntry struct {
	Id             int32             `json:"Id,omitempty"`             // 主键
	Radius         decimal.Decimal   `json:"Radius,omitempty"`         //战斗区域半径
	Height         decimal.Decimal   `json:"Height,omitempty"`         //战斗区域矩形长
	UsPositionX    []decimal.Decimal `json:"UsPositionX,omitempty"`    //友军单位x坐标
	UsPositionZ    []decimal.Decimal `json:"UsPositionZ,omitempty"`    //友军单位z坐标
	UsInitalCom    []decimal.Decimal `json:"UsInitalCom,omitempty"`    //友军单位z坐标
	UsRotation     []decimal.Decimal `json:"UsRotation,omitempty"`     //友军单位旋转值
	EnemyPositionX []decimal.Decimal `json:"EnemyPositionX,omitempty"` //敌军单位x坐标
	EnemyPositionZ []decimal.Decimal `json:"EnemyPositionZ,omitempty"` //敌军单位z坐标
	EnemyInitalCom []decimal.Decimal `json:"EnemyInitalCom,omitempty"` //敌军单位z坐标
	EnemyRotation  []decimal.Decimal `json:"EnemyRotation,omitempty"`  //敌军单位旋转值
}

// BattleView.xlsx属性表集合
type BattleViewEntries struct {
	Rows map[int32]*BattleViewEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("BattleView.xlsx", (*BattleViewEntries)(nil))
}

func (e *BattleViewEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	battleViewEntries = &BattleViewEntries{
		Rows: make(map[int32]*BattleViewEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &BattleViewEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		battleViewEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetBattleViewEntry(id int32) (*BattleViewEntry, bool) {
	entry, ok := battleViewEntries.Rows[id]
	return entry, ok
}

func GetBattleViewSize() int32 {
	return int32(len(battleViewEntries.Rows))
}

func GetBattleViewRows() map[int32]*BattleViewEntry {
	return battleViewEntries.Rows
}
