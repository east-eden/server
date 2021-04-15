package define

// 一级属性
const (
	Att_Begin       = 0 // 32位属性开始
	Att_SecondBegin = 0 // 二级属性开始

	Att_AtkBase          = 0  // 0 攻击力基础值
	Att_AtkPercent       = 1  // 1 攻击力百分比
	Att_AtkFinal         = 2  // 2 攻击力最终值
	Att_ArmorBase        = 3  // 3 护甲基础值
	Att_ArmorPercent     = 4  // 4 护甲百分比
	Att_ArmorFinal       = 5  // 5 护甲最终值
	Att_DmgInc           = 6  // 6 总伤害加成
	Att_Crit             = 7  // 7 暴击值
	Att_CritInc          = 8  // 8 暴击倍数加层
	Att_Tenacity         = 9  // 9 韧性
	Att_HealBase         = 10 // 10 治疗强度基础值
	Att_HealPercent      = 11 // 11 治疗强度百分比
	Att_HealFinal        = 12 // 12 治疗强度最终值
	Att_RealDmg          = 13 // 13 真实伤害
	Att_MoveSpeedBase    = 14 // 14 移动速度基础值
	Att_MoveSpeedPercent = 15 // 15 移动速度百分比
	Att_MoveSpeedFinal   = 16 // 16 移动速度最终值
	Att_AtbSpeedBase     = 17 // 17 时间槽速度基础值
	Att_AtbSpeedPercent  = 18 // 18 时间槽速度百分比
	Att_AtbSpeedFinal    = 19 // 19 时间槽速度最终值
	Att_EffectHit        = 20 // 20 技能效果命中
	Att_EffectResist     = 21 // 21 技能效果抵抗
	Att_MaxHPBase        = 22 // 22 生命值上限基础值
	Att_MaxHPPercent     = 23 // 23 生命值上限百分比
	Att_MaxHPFinal       = 24 // 24 生命值上限最终值
	Att_CurHP            = 25 // 25 当前生命值
	Att_MaxMP            = 26 // 26 蓝量上限
	Att_CurMP            = 27 // 27 当前蓝量
	Att_GenMP            = 28 // 28 mp恢复值
	Att_Rage             = 29 // 29 怒气值
	Att_GenRagePercent   = 30 // 30 回怒百分比
	Att_InitRage         = 31 // 31 初始怒气
	Att_Hit              = 32 // 32 命中
	Att_Dodge            = 33 // 33 闪避
	Att_MoveScope        = 34 // 34 移动范围
	Att_MoveTime         = 35 // 35 移动时间限制

	Att_DmgTypeBegin = 36 // 36 各系伤害加成begin
	Att_DmgPhysics   = 36 // 36 物理系伤害加成
	Att_DmgEarth     = 37 // 37 地系伤害加成
	Att_DmgWater     = 38 // 38 水系伤害加成
	Att_DmgFire      = 39 // 39 火系伤害加成
	Att_DmgWind      = 40 // 40 风系伤害加成
	Att_DmgTime      = 41 // 41 时系伤害加成
	Att_DmgSpace     = 42 // 42 空系伤害加成
	Att_DmgMirage    = 43 // 43 幻系伤害加成
	Att_DmgLight     = 44 // 44 光系伤害加成
	Att_DmgDark      = 45 // 45 暗系伤害加成
	Att_DmgTypeEnd   = 46 // 46 各系伤害加成end

	Att_ResTypeBegin = 46 // 46 各系伤害抗性
	Att_ResPhysics   = 46 // 46 物理系伤害抗性
	Att_ResEarth     = 47 // 47 地系伤害抗性
	Att_ResWater     = 48 // 48 水系伤害抗性
	Att_ResFire      = 49 // 49 火系伤害抗性
	Att_ResWind      = 50 // 50 风系伤害抗性
	Att_ResTime      = 51 // 51 时系伤害抗性
	Att_ResSpace     = 52 // 52 空系伤害抗性
	Att_ResMirage    = 53 // 53 幻系伤害抗性
	Att_ResLight     = 54 // 54 光系伤害抗性
	Att_ResDark      = 55 // 55 暗系伤害抗性
	Att_ResTypeEnd   = 56 // 56 各系伤害抗性

	Att_SecondEnd = 56 // 二级属性结束
	Att_End       = 56 // 32位属性结束
)

const (
	AttNum int = Att_End - Att_Begin // 32位属性数量
)

// combat legacy: do not edit!
const (
	Att_Block  = 9  // 格挡
	Att_Broken = 10 // 破击
	Att_DmgDec = 14 // 伤害减免
)

// 属性名
var AttNames = [Att_End]string{
	Att_AtkBase:          "攻击力",
	Att_AtkPercent:       "攻击力百分比",
	Att_ArmorBase:        "护甲",
	Att_ArmorPercent:     "护甲百分比",
	Att_DmgInc:           "总伤害加成",
	Att_Crit:             "暴击值",
	Att_CritInc:          "暴击倍数加层",
	Att_Tenacity:         "韧性",
	Att_HealBase:         "治疗强度",
	Att_HealPercent:      "治疗强度百分比",
	Att_RealDmg:          "真实伤害",
	Att_MoveSpeedBase:    "战场移动速度",
	Att_MoveSpeedPercent: "战场移动速度百分比",
	Att_AtbSpeedBase:     "时间槽速度",
	Att_AtbSpeedPercent:  "时间槽速度百分比",
	Att_EffectHit:        "技能效果命中",
	Att_EffectResist:     "技能效果抵抗",
	Att_MaxHPBase:        "生命值上限",
	Att_MaxHPPercent:     "生命值上限百分比",
	Att_CurHP:            "当前生命值",
	Att_MaxMP:            "蓝量上限",
	Att_CurMP:            "当前蓝量",
	Att_GenMP:            "mp恢复值",
	Att_Rage:             "怒气值",
	Att_GenRagePercent:   "怒气增长提高百分比",
	Att_InitRage:         "初始怒气",
	Att_Hit:              "命中",
	Att_Dodge:            "闪避",
	Att_MoveScope:        "移动范围",
	Att_MoveTime:         "移动时间限制",

	Att_DmgPhysics: "物理系伤害加成",
	Att_DmgEarth:   "地系伤害加成",
	Att_DmgWater:   "水系伤害加成",
	Att_DmgFire:    "火系伤害加成",
	Att_DmgWind:    "风系伤害加成",
	Att_DmgTime:    "时系伤害加成",
	Att_DmgSpace:   "空系伤害加成",
	Att_DmgMirage:  "幻系伤害加成",
	Att_DmgLight:   "光系伤害加成",
	Att_DmgDark:    "暗系伤害加成",

	Att_ResPhysics: "物理系伤害抗性",
	Att_ResEarth:   "地系伤害抗性",
	Att_ResWater:   "水系伤害抗性",
	Att_ResFire:    "火系伤害抗性",
	Att_ResWind:    "风系伤害抗性",
	Att_ResTime:    "时系伤害抗性",
	Att_ResSpace:   "空系伤害抗性",
	Att_ResMirage:  "幻系伤害抗性",
	Att_ResLight:   "光系伤害抗性",
	Att_ResDark:    "暗系伤害抗性",
}
