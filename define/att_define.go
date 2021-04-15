package define

// 一级属性
const (
	Att_Begin       = 0 // 32位属性开始
	Att_SecondBegin = 0 // 二级属性开始

	Att_Atk              = 0  // 0 攻击力
	Att_AtkPercent       = 1  // 1 攻击力百分比
	Att_Armor            = 2  // 2 护甲
	Att_ArmorPercent     = 3  // 3 护甲百分比
	Att_DmgInc           = 4  // 4 总伤害加成
	Att_Crit             = 5  // 5 暴击值
	Att_CritInc          = 6  // 6 暴击倍数加层
	Att_Tenacity         = 7  // 7 韧性
	Att_Heal             = 8  // 8 治疗强度
	Att_HealPercent      = 9  // 9 治疗强度百分比
	Att_RealDmg          = 10 // 10 真实伤害
	Att_MoveSpeed        = 11 // 11 战场移动速度
	Att_MoveSpeedPercent = 12 // 12 移动速度百分比
	Att_AtbSpeed         = 13 // 13 时间槽速度
	Att_AtbSpeedPercent  = 14 // 14 时间槽速度百分比
	Att_EffectHit        = 15 // 15 技能效果命中
	Att_EffectResist     = 16 // 16 技能效果抵抗
	Att_MaxHP            = 17 // 17 生命值上限
	Att_MaxHPPercent     = 18 // 18 生命值上限百分比
	Att_CurHP            = 19 // 19 当前生命值
	Att_MaxMP            = 20 // 20 蓝量上限
	Att_CurMP            = 21 // 21 当前蓝量
	Att_GenMP            = 22 // 22 mp恢复值
	Att_Rage             = 23 // 23 怒气值
	Att_GenRagePercent   = 24 // 24 回怒百分比
	Att_InitRage         = 25 // 25 初始怒气
	Att_Hit              = 26 // 26 命中
	Att_Dodge            = 27 // 27 闪避
	Att_MoveScope        = 28 // 28 移动范围
	Att_MoveTime         = 29 // 29 移动时间限制

	Att_DmgTypeBegin = 30 // 30 各系伤害加成begin
	Att_DmgPhysics   = 30 // 30 物理系伤害加成
	Att_DmgEarth     = 31 // 31 地系伤害加成
	Att_DmgWater     = 32 // 32 水系伤害加成
	Att_DmgFire      = 33 // 33 火系伤害加成
	Att_DmgWind      = 34 // 34 风系伤害加成
	Att_DmgTime      = 35 // 35 时系伤害加成
	Att_DmgSpace     = 36 // 36 空系伤害加成
	Att_DmgMirage    = 37 // 37 幻系伤害加成
	Att_DmgLight     = 38 // 38 光系伤害加成
	Att_DmgDark      = 39 // 39 暗系伤害加成
	Att_DmgTypeEnd   = 40 // 40 各系伤害加成end

	Att_ResTypeBegin = 40 // 40 各系伤害抗性
	Att_ResPhysics   = 40 // 40 物理系伤害抗性
	Att_ResEarth     = 41 // 41 地系伤害抗性
	Att_ResWater     = 42 // 42 水系伤害抗性
	Att_ResFire      = 43 // 43 火系伤害抗性
	Att_ResWind      = 44 // 44 风系伤害抗性
	Att_ResTime      = 45 // 45 时系伤害抗性
	Att_ResSpace     = 46 // 46 空系伤害抗性
	Att_ResMirage    = 47 // 47 幻系伤害抗性
	Att_ResLight     = 48 // 48 光系伤害抗性
	Att_ResDark      = 49 // 49 暗系伤害抗性
	Att_ResTypeEnd   = 50 // 50 各系伤害抗性

	Att_SecondEnd = 50 // 二级属性结束
	Att_End       = 50 // 32位属性结束
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
	Att_Atk:              "攻击力",
	Att_AtkPercent:       "攻击力百分比",
	Att_Armor:            "护甲",
	Att_ArmorPercent:     "护甲百分比",
	Att_DmgInc:           "总伤害加成",
	Att_Crit:             "暴击值",
	Att_CritInc:          "暴击倍数加层",
	Att_Tenacity:         "韧性",
	Att_Heal:             "治疗强度",
	Att_HealPercent:      "治疗强度百分比",
	Att_RealDmg:          "真实伤害",
	Att_MoveSpeed:        "战场移动速度",
	Att_MoveSpeedPercent: "战场移动速度百分比",
	Att_AtbSpeed:         "时间槽速度",
	Att_AtbSpeedPercent:  "时间槽速度百分比",
	Att_EffectHit:        "技能效果命中",
	Att_EffectResist:     "技能效果抵抗",
	Att_MaxHP:            "生命值上限",
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
