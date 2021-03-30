package auto

import (
	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var attEntries *AttEntries //Att.xlsx全局变量

// Att.xlsx属性表
type AttEntry struct {
	Id               int32    `json:"Id,omitempty"`               // 主键
	Atk              number   `json:"Atk,omitempty"`              //攻击力
	AtkPercent       number   `json:"AtkPercent,omitempty"`       //攻击力百分比
	Armor            number   `json:"Armor,omitempty"`            //护甲
	ArmorPercent     number   `json:"ArmorPercent,omitempty"`     //护甲百分比
	DmgInc           int32    `json:"DmgInc,omitempty"`           //总伤害加成
	Crit             int32    `json:"Crit,omitempty"`             //暴击值
	CritInc          int32    `json:"CritInc,omitempty"`          //暴击倍数加成
	Tenacity         int32    `json:"Tenacity,omitempty"`         //韧性值
	Heal             int32    `json:"Heal,omitempty"`             //治疗强度
	HealPercent      number   `json:"HealPercent,omitempty"`      //治疗强度百分比
	RealDmg          int32    `json:"RealDmg,omitempty"`          //真实伤害
	MoveSpeed        number   `json:"MoveSpeed,omitempty"`        //战场移动速度
	MoveSpeedPercent number   `json:"MoveSpeedPercent,omitempty"` //战场移动速度百分比
	AtbSpeed         number   `json:"AtbSpeed,omitempty"`         //时间槽速度
	AtbSpeedPercent  number   `json:"AtbSpeedPercent,omitempty"`  //时间槽速度百分比
	EffectHit        int32    `json:"EffectHit,omitempty"`        //技能效果命中
	EffectResist     int32    `json:"EffectResist,omitempty"`     //技能效果抵抗
	MaxHP            int32    `json:"MaxHP,omitempty"`            //血量上限
	MaxHPPercent     number   `json:"MaxHPPercent,omitempty"`     //血量上限百分比
	MaxMP            int32    `json:"MaxMP,omitempty"`            //MP上限
	MaxMPPercent     number   `json:"MaxMPPercent,omitempty"`     //MP上限百分比
	GenMP            int32    `json:"GenMP,omitempty"`            //MP恢复
	GenMPPercent     number   `json:"GenMPPercent,omitempty"`     //MP恢复百分比
	Rage             int32    `json:"Rage,omitempty"`             //怒气上限
	GenRagePercent   number   `json:"GenRagePercent,omitempty"`   //怒气增长提高百分比
	InitRage         int32    `json:"InitRage,omitempty"`         //初始怒气
	Hit              int32    `json:"Hit,omitempty"`              //命中值
	Dodge            int32    `json:"Dodge,omitempty"`            //闪避值
	MoveScope        number   `json:"MoveScope,omitempty"`        //移动范围
	DmgOfType        []number `json:"DmgOfType,omitempty"`        //各系伤害加层
	ResOfType        []number `json:"ResOfType,omitempty"`        //各系伤害减免
}

// Att.xlsx属性表集合
type AttEntries struct {
	Rows map[int32]*AttEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Att.xlsx", (*AttEntries)(nil))
}

func (e *AttEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	attEntries = &AttEntries{
		Rows: make(map[int32]*AttEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &AttEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		attEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetAttEntry(id int32) (*AttEntry, bool) {
	entry, ok := attEntries.Rows[id]
	return entry, ok
}

func GetAttSize() int32 {
	return int32(len(attEntries.Rows))
}

func GetAttRows() map[int32]*AttEntry {
	return attEntries.Rows
}
