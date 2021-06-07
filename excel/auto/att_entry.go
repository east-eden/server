package auto

import (
	"e.coding.net/mmstudio/blade/server/excel"
	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var attEntries *AttEntries //Att.xlsx全局变量

// Att.xlsx属性表
type AttEntry struct {
	Id               int32             `json:"Id,omitempty"`               // 主键
	Atk              decimal.Decimal   `json:"Atk,omitempty"`              //攻击力
	AtkPercent       decimal.Decimal   `json:"AtkPercent,omitempty"`       //攻击力百分比
	Armor            decimal.Decimal   `json:"Armor,omitempty"`            //护甲
	ArmorPercent     decimal.Decimal   `json:"ArmorPercent,omitempty"`     //护甲百分比
	SelfDmgInc       decimal.Decimal   `json:"SelfDmgInc,omitempty"`       //自身伤害加成
	EnemyWoundInc    decimal.Decimal   `json:"EnemyWoundInc,omitempty"`    //敌方受伤加成
	SelfDmgDec       decimal.Decimal   `json:"SelfDmgDec,omitempty"`       //自身伤害消减
	EnemyWoundDec    decimal.Decimal   `json:"EnemyWoundDec,omitempty"`    //敌方受伤消减
	Crit             decimal.Decimal   `json:"Crit,omitempty"`             //暴击率
	CritInc          decimal.Decimal   `json:"CritInc,omitempty"`          //暴击倍数加成
	Tenacity         decimal.Decimal   `json:"Tenacity,omitempty"`         //韧性值
	Heal             int32             `json:"Heal,omitempty"`             //治疗强度
	HealPercent      decimal.Decimal   `json:"HealPercent,omitempty"`      //治疗效果加成百分比
	GethealPercent   decimal.Decimal   `json:"GethealPercent,omitempty"`   //受治疗效果加成百分比
	RealDmg          int32             `json:"RealDmg,omitempty"`          //真实伤害
	MoveSpeed        decimal.Decimal   `json:"MoveSpeed,omitempty"`        //战场移动速度
	MoveSpeedPercent decimal.Decimal   `json:"MoveSpeedPercent,omitempty"` //战场移动速度百分比
	AtbSpeed         decimal.Decimal   `json:"AtbSpeed,omitempty"`         //时间槽速度
	AtbSpeedPercent  decimal.Decimal   `json:"AtbSpeedPercent,omitempty"`  //时间槽速度百分比
	EffectHit        int32             `json:"EffectHit,omitempty"`        //技能效果命中
	EffectResist     int32             `json:"EffectResist,omitempty"`     //技能效果抵抗
	MaxHP            int32             `json:"MaxHP,omitempty"`            //血量上限
	MaxHPPercent     decimal.Decimal   `json:"MaxHPPercent,omitempty"`     //血量上限百分比
	MaxMP            int32             `json:"MaxMP,omitempty"`            //MP上限
	MaxMPPercent     decimal.Decimal   `json:"MaxMPPercent,omitempty"`     //MP上限百分比
	GenMP            int32             `json:"GenMP,omitempty"`            //MP恢复
	GenMPPercent     decimal.Decimal   `json:"GenMPPercent,omitempty"`     //MP恢复百分比
	MaxRage          int32             `json:"MaxRage,omitempty"`          //怒气上限
	GenRagePercent   decimal.Decimal   `json:"GenRagePercent,omitempty"`   //怒气增长提高百分比
	InitRage         int32             `json:"InitRage,omitempty"`         //初始怒气
	Hit              int32             `json:"Hit,omitempty"`              //命中值
	Dodge            int32             `json:"Dodge,omitempty"`            //闪避值
	MoveScope        decimal.Decimal   `json:"MoveScope,omitempty"`        //移动范围
	DmgOfType        []decimal.Decimal `json:"DmgOfType,omitempty"`        //各系伤害加层
	ResOfType        []decimal.Decimal `json:"ResOfType,omitempty"`        //各系伤害减免
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
