package define

// 一级属性
const (
	Att_Begin      = 0 // 32位属性开始
	Att_FirstBegin = 0 // 一级属性开始
	Att_Str        = 0 // 0 力量
	Att_Agl        = 1 // 1 敏捷
	Att_Con        = 2 // 2 耐力
	Att_Int        = 3 // 3 智力
	Att_AtkSpeed   = 4 // 4 攻击速度
	Att_FirstEnd   = 5 // 一级属性结束

	Att_SecondBegin  = 5  // 二级属性开始
	Att_Atk          = 5  // 5 攻击力
	Att_Def          = 6  // 6 防御力
	Att_CriProb      = 7  // 7 暴击率
	Att_CriDmg       = 8  // 8 暴击伤害
	Att_EffectHit    = 9  // 9 效果命中
	Att_EffectResist = 10 // 10 效果抵抗
	Att_ConPercent   = 11 // 11 体力加层
	Att_AtkPercent   = 12 // 12 攻击力加层
	Att_DefPercent   = 13 // 13 防御力加层
	Att_Rage         = 14 // 14 怒气
	Att_Hit          = 15 // 15 命中
	Att_Dodge        = 16 // 16 闪避
	Att_Block        = 17 // 17 格挡
	Att_Broken       = 18 // 18 破击
	Att_DmgInc       = 19 // 19 伤害加层
	Att_DmgDec       = 20 // 20 伤害减免
	Att_SecondEnd    = 21 // 二级属性结束
	Att_End          = 21 // 32位属性结束

	Att_Plus_Begin = 21 // 64位属性开始
	Att_Plus_CurHP = 21 // 当前血量
	Att_Plus_CurMP = 22 // 当前蓝量
	Att_Plus_MaxHP = 23 // 23 血量
	Att_Plus_MaxMP = 24 // 24 蓝量
	Att_Plus_End   = 25 // 64位属性结束
)

const (
	AttNum     int = Att_End - Att_Begin           // 32位属性数量
	PlusAttNum int = Att_Plus_End - Att_Plus_Begin // 64位属性数量
)
