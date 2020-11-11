package define

// 一级属性
const (
	Att_Begin = iota

	Att_FirstBegin = iota - 1 // 一级属性开始
	Att_Str        = iota - 2 // 力量
	Att_Agl                   // 敏捷
	Att_Con                   // 耐力
	Att_Int                   // 智力
	Att_AtkSpeed              // 攻击速度
	Att_FirstEnd              // 一级属性结束

	Att_SecondBegin  = iota - 3 // 二级属性开始
	Att_MaxHP        = iota - 4 // 血量
	Att_MaxMP                   // 蓝量
	Att_Atk                     // 攻击力
	Att_Def                     // 防御力
	Att_CriProb                 // 暴击率
	Att_CriDmg                  // 暴击伤害
	Att_EffectHit               // 效果命中
	Att_EffectResist            // 效果抵抗
	Att_ConPercent              // 体力加层
	Att_AtkPercent              // 攻击力加层
	Att_DefPercent              // 防御力加层
	Att_Rage                    // 怒气
	Att_Hit                     // 命中
	Att_Dodge                   // 闪避
	Att_Block                   // 格挡
	Att_Broken                  // 破击
	Att_DmgInc                  // 伤害加层
	Att_DmgDec                  // 伤害减免
	Att_SecondEnd               // 二级属性结束

	Att_CurHP = iota - 5 // 当前血量
	Att_CurMP            // 当前蓝量

	Att_End
)

// att entry
type AttEntry struct {
	ID   int32  `json:"_id"`
	Desc string `json:"Desc"`

	// 一级属性
	Str      int64 `json:"Str"`
	Agl      int64 `json:"Agl"`
	Con      int64 `json:"Con"`
	Int      int64 `json:"Int"`
	AtkSpeed int64 `json:"AtkSpeed"`

	// 二级属性
	MaxHP        int64 `json:"MaxHP"`
	MaxMP        int64 `json:"MaxMP"`
	Atk          int64 `json:"Atk"`
	Def          int64 `json:"Def"`
	CriProb      int64 `json:"CriProb"`
	CriDmg       int64 `json:"CriDmg"`
	EffectHit    int64 `json:"EffectHit"`
	EffectResist int64 `json:"EffectResist"`
	ConPercent   int64 `json:"ConPercent"`
	AtkPercent   int64 `json:"AtkPercent"`
	DefPercent   int64 `json:"DefPercent"`
	Rage         int64 `json:"Rage"`
	Hit          int64 `json:"Hit"`
	Dodge        int64 `json:"Dodge"`
	Block        int64 `json:"Block"`
	Broken       int64 `json:"Broken"`
	DmgInc       int64 `json:"DmgInc"`
	DmgDec       int64 `json:"DmgDec"`
}
