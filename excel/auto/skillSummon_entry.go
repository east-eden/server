package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var skillSummonEntries *SkillSummonEntries //SkillSummon.xlsx全局变量

// SkillSummon.xlsx属性表
type SkillSummonEntry struct {
	Id                int32           `json:"Id,omitempty"`                // 主键
	Name              string          `json:"Name,omitempty"`              //名称
	LoopFx            string          `json:"LoopFx,omitempty"`            //生成持续表现
	EndFx             string          `json:"EndFx,omitempty"`             //结束表现
	Icon              string          `json:"Icon,omitempty"`              //图标
	LifeTime          decimal.Decimal `json:"LifeTime,omitempty"`          //持续时间
	EffectCD          decimal.Decimal `json:"EffectCD,omitempty"`          //冷却时间
	Type              int32           `json:"Type,omitempty"`              //召唤物类型
	VW1XT5G753TykGiQp decimal.Decimal `json:"VW1XT5G753TykGiQp,omitempty"` //
	TimelineID        int32           `json:"TimelineID,omitempty"`        //召唤物效果
	Limit             int32           `json:"Limit,omitempty"`             //数量上限
	IsDieClear        int32           `json:"IsDieClear,omitempty"`        //死亡清除
	Cleartrigger      int32           `json:"Cleartrigger,omitempty"`      //清除触发器
}

// SkillSummon.xlsx属性表集合
type SkillSummonEntries struct {
	Rows map[int32]*SkillSummonEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("SkillSummon.xlsx", (*SkillSummonEntries)(nil))
}

func (e *SkillSummonEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	skillSummonEntries = &SkillSummonEntries{
		Rows: make(map[int32]*SkillSummonEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillSummonEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		skillSummonEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetSkillSummonEntry(id int32) (*SkillSummonEntry, bool) {
	entry, ok := skillSummonEntries.Rows[id]
	return entry, ok
}

func GetSkillSummonSize() int32 {
	return int32(len(skillSummonEntries.Rows))
}

func GetSkillSummonRows() map[int32]*SkillSummonEntry {
	return skillSummonEntries.Rows
}
