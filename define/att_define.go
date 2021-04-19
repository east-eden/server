package define

// 一级属性
const (
	Att_Begin       = 0 // 32位属性开始
	Att_SecondBegin = 0 // 二级属性开始

	Att_Atk            = 0  // 0 攻击力最终值
	Att_Armor          = 1  // 1 护甲最终值
	Att_DmgInc         = 2  // 2 总伤害加成
	Att_Crit           = 3  // 3 暴击值
	Att_CritInc        = 4  // 4 暴击倍数加层
	Att_Tenacity       = 5  // 5 韧性
	Att_Heal           = 6  // 6 治疗强度最终值
	Att_RealDmg        = 7  // 7 真实伤害
	Att_MoveSpeed      = 8  // 8 移动速度最终值
	Att_AtbSpeed       = 9  // 9 时间槽速度最终值
	Att_EffectHit      = 10 // 10 技能效果命中
	Att_EffectResist   = 11 // 11 技能效果抵抗
	Att_MaxHP          = 12 // 12 生命值上限最终值
	Att_CurHP          = 13 // 13 当前生命值
	Att_MaxMP          = 14 // 14 蓝量上限
	Att_CurMP          = 15 // 15 当前蓝量
	Att_GenMP          = 16 // 16 mp恢复值
	Att_MaxRage        = 17 // 17 怒气值上限
	Att_GenRagePercent = 18 // 18 回怒百分比
	Att_InitRage       = 19 // 19 初始怒气
	Att_Hit            = 20 // 20 命中
	Att_Dodge          = 21 // 21 闪避
	Att_MoveScope      = 22 // 22 移动范围

	Att_DmgTypeBegin = 23 // 23 各系伤害加成begin
	Att_DmgPhysics   = 23 // 23 物理系伤害加成
	Att_DmgEarth     = 24 // 24 地系伤害加成
	Att_DmgWater     = 25 // 25 水系伤害加成
	Att_DmgFire      = 26 // 26 火系伤害加成
	Att_DmgWind      = 27 // 27 风系伤害加成
	Att_DmgThunder   = 28 // 28 雷系伤害加成
	Att_DmgTime      = 29 // 29 时系伤害加成
	Att_DmgSpace     = 30 // 30 空系伤害加成
	Att_DmgSteel     = 31 // 31 钢系伤害加成
	Att_DmgDeath     = 32 // 32 灭系伤害加成
	Att_DmgTypeEnd   = 33 // 33 各系伤害加成end

	Att_ResTypeBegin = 33 // 33 各系伤害抗性
	Att_ResPhysics   = 33 // 33 物理系伤害抗性
	Att_ResEarth     = 34 // 34 地系伤害抗性
	Att_ResWater     = 35 // 35 水系伤害抗性
	Att_ResFire      = 36 // 36 火系伤害抗性
	Att_ResWind      = 37 // 37 风系伤害抗性
	Att_ResThunder   = 38 // 38 雷系伤害抗性
	Att_ResTime      = 39 // 39 时系伤害抗性
	Att_ResSpace     = 40 // 40 空系伤害抗性
	Att_ResSteel     = 41 // 41 钢系伤害抗性
	Att_ResDeath     = 42 // 42 灭系伤害抗性
	Att_ResTypeEnd   = 43 // 43 各系伤害抗性

	Att_SecondEnd = 43 // 二级属性结束
	Att_End       = 43 // 32位属性结束
)

// 基础属性枚举
const (
	Att_Base_Begin    = 0
	Att_AtkBase       = 0 // 0 攻击力基础值
	Att_ArmorBase     = 1 // 1 护甲基础值
	Att_HealBase      = 2 // 2 治疗强度基础值
	Att_MoveSpeedBase = 3 // 3 移动速度基础值
	Att_AtbSpeedBase  = 4 // 4 时间槽速度基础值
	Att_MaxHPBase     = 5 // 5 生命值上限基础值
	Att_Base_End      = 6
)

// 百分比属性枚举
const (
	Att_Percent_Begin    = 0
	Att_AtkPercent       = 0 // 0 攻击力百分比
	Att_ArmorPercent     = 1 // 1 护甲百分比
	Att_HealPercent      = 2 // 2 治疗强度百分比
	Att_MoveSpeedPercent = 3 // 3 移动速度百分比
	Att_AtbSpeedPercent  = 4 // 4 时间槽速度百分比
	Att_MaxHPPercent     = 5 // 5 生命值上限百分比
	Att_Percent_End      = 6
)

const (
	AttFinalNum   = Att_End - Att_Begin // 32位属性数量
	AttBaseNum    = Att_Base_End - Att_Base_Begin
	AttPercentNum = Att_Percent_End - Att_Percent_Begin
)

// combat legacy: do not edit!
const (
	Att_Block  = 9  // 格挡
	Att_Broken = 10 // 破击
	Att_DmgDec = 14 // 伤害减免
)

// 属性名
var AttNames = [Att_End]string{
	Att_Atk:            "攻击力",
	Att_DmgInc:         "总伤害加成",
	Att_Crit:           "暴击值",
	Att_CritInc:        "暴击倍数加层",
	Att_Tenacity:       "韧性",
	Att_Heal:           "治疗强度",
	Att_RealDmg:        "真实伤害",
	Att_MoveSpeed:      "战场移动速度",
	Att_AtbSpeed:       "时间槽速度",
	Att_EffectHit:      "技能效果命中",
	Att_EffectResist:   "技能效果抵抗",
	Att_MaxHP:          "生命值上限",
	Att_CurHP:          "当前生命值",
	Att_MaxMP:          "蓝量上限",
	Att_CurMP:          "当前蓝量",
	Att_GenMP:          "mp恢复值",
	Att_MaxRage:        "怒气值",
	Att_GenRagePercent: "怒气增长提高百分比",
	Att_InitRage:       "初始怒气",
	Att_Hit:            "命中",
	Att_Dodge:          "闪避",
	Att_MoveScope:      "移动范围",

	Att_DmgPhysics: "物理系伤害加成",
	Att_DmgEarth:   "地系伤害加成",
	Att_DmgWater:   "水系伤害加成",
	Att_DmgFire:    "火系伤害加成",
	Att_DmgWind:    "风系伤害加成",
	Att_DmgThunder: "雷系伤害加成",
	Att_DmgTime:    "时系伤害加成",
	Att_DmgSpace:   "空系伤害加成",
	Att_DmgSteel:   "钢系伤害加成",
	Att_DmgDeath:   "灭系伤害加成",

	Att_ResPhysics: "物理系伤害抗性",
	Att_ResEarth:   "地系伤害抗性",
	Att_ResWater:   "水系伤害抗性",
	Att_ResFire:    "火系伤害抗性",
	Att_ResWind:    "风系伤害抗性",
	Att_ResThunder: "雷系伤害抗性",
	Att_ResTime:    "时系伤害抗性",
	Att_ResSpace:   "空系伤害抗性",
	Att_ResSteel:   "钢系伤害抗性",
	Att_ResDeath:   "灭系伤害抗性",
}

// 基础属性名
var AttBaseNames = [Att_Base_End]string{
	Att_AtkBase:       "攻击力",
	Att_ArmorBase:     "护甲",
	Att_HealBase:      "治疗强度",
	Att_MoveSpeedBase: "战场移动速度",
	Att_AtbSpeedBase:  "时间槽速度",
	Att_MaxHPBase:     "生命值上限",
}

// 百分比属性名
var AttPercentNames = [Att_Percent_End]string{
	Att_AtkPercent:       "攻击力百分比",
	Att_ArmorPercent:     "护甲百分比",
	Att_HealPercent:      "治疗强度百分比",
	Att_MoveSpeedPercent: "战场移动速度百分比",
	Att_AtbSpeedPercent:  "时间槽速度百分比",
	Att_MaxHPPercent:     "生命值上限百分比",
}
