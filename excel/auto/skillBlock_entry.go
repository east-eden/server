package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var skillBlockEntries *SkillBlockEntries //SkillBlock.xlsx全局变量

// SkillBlock.xlsx属性表
type SkillBlockEntry struct {
	Id        int32     `json:"Id,omitempty"`        // 主键
	Index     int32     `json:"Index,omitempty"`     //技能块的索引
	Condition []string  `json:"Condition,omitempty"` //触发条件
	Buffs     []int32   `json:"Buffs,omitempty"`     //添加Buff
	Formula   int32     `json:"Formula,omitempty"`   //伤害公式
	Ratio     []float32 `json:"Ratio,omitempty"`     //伤害系数
}

// SkillBlock.xlsx属性表集合
type SkillBlockEntries struct {
	Rows map[int32]*SkillBlockEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("SkillBlock.xlsx", (*SkillBlockEntries)(nil))
}

func (e *SkillBlockEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	skillBlockEntries = &SkillBlockEntries{
		Rows: make(map[int32]*SkillBlockEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillBlockEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		skillBlockEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetSkillBlockEntry(id int32) (*SkillBlockEntry, bool) {
	entry, ok := skillBlockEntries.Rows[id]
	return entry, ok
}

func GetSkillBlockSize() int32 {
	return int32(len(skillBlockEntries.Rows))
}

func GetSkillBlockRows() map[int32]*SkillBlockEntry {
	return skillBlockEntries.Rows
}
