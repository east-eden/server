package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var skillEffectEntries *SkillEffectEntries //SkillEffect.csv全局变量

// SkillEffect.csv属性表
type SkillEffectEntry struct {
	Id                int32             `json:"Id,omitempty"`                // 主键
	Jumptext          string            `json:"Jumptext,omitempty"`          //特效飘字
	IsDebuff          int32             `json:"IsDebuff,omitempty"`          //
	EffectFx          string            `json:"EffectFx,omitempty"`          //
	AddRage           int32             `json:"AddRage,omitempty"`           //
	CutProtection     int32             `json:"CutProtection,omitempty"`     //
	HitFxName         string            `json:"HitFxName,omitempty"`         //受击特效
	HitFxSlot         string            `json:"HitFxSlot,omitempty"`         //受击特效插槽
	HitStopTime       decimal.Decimal   `json:"HitStopTime,omitempty"`       //受击停顿时间,仅普攻
	HitAnimName       string            `json:"HitAnimName,omitempty"`       //受击动作,单体弹道的受击时间为物理碰撞的时间点
	HitHurtTime       decimal.Decimal   `json:"HitHurtTime,omitempty"`       //受伤状态的持续时间,仅普攻,切成受伤
	HitBackDistance   decimal.Decimal   `json:"HitBackDistance,omitempty"`   //受伤状态的持续时间的,击退距离
	IsEffectHit       int32             `json:"IsEffectHit,omitempty"`       //判定类型
	Prob              int32             `json:"Prob,omitempty"`              //触发概率
	EffectType        int32             `json:"EffectType,omitempty"`        //效果类型
	ParameterString   []string          `json:"ParameterString,omitempty"`   //参数1
	ParameterInt      []int32           `json:"ParameterInt,omitempty"`      //参数2
	ParameterNumber   []decimal.Decimal `json:"ParameterNumber,omitempty"`   //
	AttributeNumValue *treemap.Map      `json:"AttributeNumValue,omitempty"` //属性
	SkillLaunch       int32             `json:"SkillLaunch,omitempty"`       //发起类型
	TargetLength      decimal.Decimal   `json:"TargetLength,omitempty"`      //范围长
	TargetWide        decimal.Decimal   `json:"TargetWide,omitempty"`        //范围宽
	RangeType         int32             `json:"RangeType,omitempty"`         //目标范围
	Scope             int32             `json:"Scope,omitempty"`             //作用对象
	TriggerId         []int32           `json:"TriggerId,omitempty"`         //触发器类型
	TriggerRelation   int32             `json:"TriggerRelation,omitempty"`   //触发器,条件关系
}

// SkillEffect.csv属性表集合
type SkillEffectEntries struct {
	Rows map[int32]*SkillEffectEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("SkillEffect.csv", (*SkillEffectEntries)(nil))
}

func (e *SkillEffectEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	skillEffectEntries = &SkillEffectEntries{
		Rows: make(map[int32]*SkillEffectEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &SkillEffectEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		skillEffectEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetSkillEffectEntry(id int32) (*SkillEffectEntry, bool) {
	entry, ok := skillEffectEntries.Rows[id]
	return entry, ok
}

func GetSkillEffectSize() int32 {
	return int32(len(skillEffectEntries.Rows))
}

func GetSkillEffectRows() map[int32]*SkillEffectEntry {
	return skillEffectEntries.Rows
}
