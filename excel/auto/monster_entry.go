package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var monsterEntries *MonsterEntries //Monster.xlsx全局变量

// Monster.xlsx属性表
type MonsterEntry struct {
	Id            int32           `json:"Id,omitempty"`            // 主键
	Name          string          `json:"Name,omitempty"`          //怪物名称
	Desc          string          `json:"Desc,omitempty"`          //怪物简介
	Type          int32           `json:"Type,omitempty"`          //类型
	Race          int32           `json:"Race,omitempty"`          //种族
	Profession    int32           `json:"Profession,omitempty"`    //职业
	ModelResource string          `json:"ModelResource,omitempty"` //模型资源
	Modelscope    decimal.Decimal `json:"Modelscope,omitempty"`    //模型范围
	Icon          string          `json:"Icon,omitempty"`          //头像资源
	AttId         int32           `json:"AttId,omitempty"`         //属性id
	Skill1        int32           `json:"Skill1,omitempty"`        //怪物普攻
	Skill2        int32           `json:"Skill2,omitempty"`        //怪物技能1
	Skill3        int32           `json:"Skill3,omitempty"`        //怪物技能2
	Skill4        int32           `json:"Skill4,omitempty"`        //怪物技能3
	PassiveSkill1 int32           `json:"PassiveSkill1,omitempty"` //怪物被动1
	PassiveSkill2 int32           `json:"PassiveSkill2,omitempty"` //怪物被动2
	PassiveSkill3 int32           `json:"PassiveSkill3,omitempty"` //怪物被动3
	PassiveSkill4 int32           `json:"PassiveSkill4,omitempty"` //怪物被动4
	PassiveSkill5 int32           `json:"PassiveSkill5,omitempty"` //怪物被动5
}

// Monster.xlsx属性表集合
type MonsterEntries struct {
	Rows map[int32]*MonsterEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Monster.xlsx", (*MonsterEntries)(nil))
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
