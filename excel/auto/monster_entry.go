package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var monsterEntries *MonsterEntries //Monster.csv全局变量

// Monster.csv属性表
type MonsterEntry struct {
	Id                int32           `json:"Id,omitempty"`                // 主键
	ModelID           int32           `json:"ModelID,omitempty"`           //
	Name              string          `json:"Name,omitempty"`              //怪物名称
	Desc              string          `json:"Desc,omitempty"`              //怪物简介
	Type              int32           `json:"Type,omitempty"`              //类型
	Race              int32           `json:"Race,omitempty"`              //种族
	Profession        int32           `json:"Profession,omitempty"`        //职业
	AttId             int32           `json:"AttId,omitempty"`             //属性id
	Skill1            int32           `json:"Skill1,omitempty"`            //普攻
	Skill2            []int32         `json:"Skill2,omitempty"`            //主动技能
	Skill3            []int32         `json:"Skill3,omitempty"`            //被动技能
	Protection        int32           `json:"Protection,omitempty"`        //防护罩值
	ProtectionBreakFx string          `json:"ProtectionBreakFx,omitempty"` //
	BreakTime         decimal.Decimal `json:"BreakTime,omitempty"`         //break,持续时间(秒)
	NormalTime        decimal.Decimal `json:"NormalTime,omitempty"`        //normal,持续时间(秒)
}

// Monster.csv属性表集合
type MonsterEntries struct {
	Rows map[int32]*MonsterEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Monster.csv", (*MonsterEntries)(nil))
}

func (e *MonsterEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	monsterEntries = &MonsterEntries{
		Rows: make(map[int32]*MonsterEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &MonsterEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		monsterEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetMonsterEntry(id int32) (*MonsterEntry, bool) {
	entry, ok := monsterEntries.Rows[id]
	return entry, ok
}

func GetMonsterSize() int32 {
	return int32(len(monsterEntries.Rows))
}

func GetMonsterRows() map[int32]*MonsterEntry {
	return monsterEntries.Rows
}
