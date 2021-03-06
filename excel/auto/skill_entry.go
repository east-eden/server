package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var skillEntries *SkillEntries //Skill.xlsx全局变量

// Skill.xlsx属性表
type SkillEntry struct {
	Id           int32   `json:"Id,omitempty"`           // 主键
	Name         string  `json:"Name,omitempty"`         //名字
	Desc         string  `json:"Desc,omitempty"`         //描述
	Icon         string  `json:"Icon,omitempty"`         //图标
	NextLevel    int32   `json:"NextLevel,omitempty"`    //下个等级
	CD           float32 `json:"CD,omitempty"`           //技能CD
	SkillCombo   int32   `json:"SkillCombo,omitempty"`   //连击类型
	SkillType    int32   `json:"SkillType,omitempty"`    //技能类型
	Precondition int32   `json:"Precondition,omitempty"` //前置条件
	TargetType   int32   `json:"TargetType,omitempty"`   //释放目标类型
	TargetRule   []int32 `json:"TargetRule,omitempty"`   //释放目标规则
	AttackType   int32   `json:"AttackType,omitempty"`   //攻击方式
	SkillBlocks  []int32 `json:"SkillBlocks,omitempty"`  //包含的技能块
	EffectType   int32   `json:"EffectType,omitempty"`   //特效类型
	Resource     string  `json:"Resource,omitempty"`     //特效路径
}

// Skill.xlsx属性表集合
type SkillEntries struct {
	Rows map[int32]*SkillEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Skill.xlsx", (*SkillEntries)(nil))
}

func (e *SkillEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	skillEntries = &SkillEntries{
		Rows: make(map[int32]*SkillEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		skillEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetSkillEntry(id int32) (*SkillEntry, bool) {
	entry, ok := skillEntries.Rows[id]
	return entry, ok
}

func GetSkillSize() int32 {
	return int32(len(skillEntries.Rows))
}

func GetSkillRows() map[int32]*SkillEntry {
	return skillEntries.Rows
}
