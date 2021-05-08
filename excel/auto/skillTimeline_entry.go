package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var skillTimelineEntries *SkillTimelineEntries //SkillTimeline.xlsx全局变量

// SkillTimeline.xlsx属性表
type SkillTimelineEntry struct {
	Id           int32           `json:"Id,omitempty"`           // 主键
	ShowType     int32           `json:"ShowType,omitempty"`     //1-战斗场景,2-虚拟环境
	Type         int32           `json:"Type,omitempty"`         //1：固定范围(普攻),2：单体弹道(普攻),3：技能表演
	DurationTime decimal.Decimal `json:"DurationTime,omitempty"` //持续时间
	AnimName     string          `json:"AnimName,omitempty"`     //动作名称
	Effects      []int32         `json:"Effects,omitempty"`      //触发效果逻辑,ID的顺序由track表中的开始时间决定
	FxName       string          `json:"FxName,omitempty"`       //可能为多个,做到一个文件里,直接挂在人身上,或武器上
	BulletFx     string          `json:"BulletFx,omitempty"`     //
	BulletSpeed  decimal.Decimal `json:"BulletSpeed,omitempty"`  //
	BulletOffset decimal.Decimal `json:"BulletOffset,omitempty"` //
	BulletHeight decimal.Decimal `json:"BulletHeight,omitempty"` //
	EffectTime   decimal.Decimal `json:"EffectTime,omitempty"`   //固定范围-是,受击的时间点,单体弹道-是,发出的时间点
}

// SkillTimeline.xlsx属性表集合
type SkillTimelineEntries struct {
	Rows map[int32]*SkillTimelineEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("SkillTimeline.xlsx", (*SkillTimelineEntries)(nil))
}

func (e *SkillTimelineEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	skillTimelineEntries = &SkillTimelineEntries{
		Rows: make(map[int32]*SkillTimelineEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillTimelineEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		skillTimelineEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetSkillTimelineEntry(id int32) (*SkillTimelineEntry, bool) {
	entry, ok := skillTimelineEntries.Rows[id]
	return entry, ok
}

func GetSkillTimelineSize() int32 {
	return int32(len(skillTimelineEntries.Rows))
}

func GetSkillTimelineRows() map[int32]*SkillTimelineEntry {
	return skillTimelineEntries.Rows
}
