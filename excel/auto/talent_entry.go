package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var talentEntries *TalentEntries //Talent.xlsx全局变量

// Talent.xlsx属性表
type TalentEntry struct {
	Id                 int32   `json:"Id,omitempty"`                 // 主键
	Type               int32   `json:"Type,omitempty"`               //天赋类型
	SubType            int32   `json:"SubType,omitempty"`            //子类型
	MaxLevel           int32   `json:"MaxLevel,omitempty"`           //天赋最大等级
	CostId             int32   `json:"CostId,omitempty"`             //天赋升级消耗ID
	OwnerId            int32   `json:"OwnerId,omitempty"`            //拥有者id
	Step               int32   `json:"Step,omitempty"`               //所属层级
	PrevStepLevelLimit int32   `json:"PrevStepLevelLimit,omitempty"` //上一层激活天赋总等级限制
	PrevTalentId1      int32   `json:"PrevTalentId1,omitempty"`      //前置天赋id1
	PrevTalentLevel1   int32   `json:"PrevTalentLevel1,omitempty"`   //前置天赋等级1
	PrevTalentId2      int32   `json:"PrevTalentId2,omitempty"`      //前置天赋id2
	PrevTalentLevel2   int32   `json:"PrevTalentLevel2,omitempty"`   //前置天赋等级2
	AttIds             []int32 `json:"AttIds,omitempty"`             //天赋等级对应的属性
	PassiveSkillId     []int32 `json:"PassiveSkillId,omitempty"`     //被动技能id
}

// Talent.xlsx属性表集合
type TalentEntries struct {
	Rows map[int32]*TalentEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Talent.xlsx", (*TalentEntries)(nil))
}

func (e *TalentEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	talentEntries = &TalentEntries{
		Rows: make(map[int32]*TalentEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &TalentEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		talentEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetTalentEntry(id int32) (*TalentEntry, bool) {
	entry, ok := talentEntries.Rows[id]
	return entry, ok
}

func GetTalentSize() int32 {
	return int32(len(talentEntries.Rows))
}

func GetTalentRows() map[int32]*TalentEntry {
	return talentEntries.Rows
}
