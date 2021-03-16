package define

// 一级属性
const (
	Att_Begin       = 0 // 32位属性开始
	Att_SecondBegin = 0 // 二级属性开始

	Att_Atk          = 0  // 0 攻击力
	Att_Armor        = 1  // 1 护甲
	Att_DmgInc       = 2  // 2 总伤害加成
	Att_Crit         = 3  // 3 暴击值
	Att_CritInc      = 4  // 4 暴击倍数加层
	Att_Tenacity     = 5  // 5 韧性
	Att_Heal         = 6  // 6 治疗强度
	Att_RealDmg      = 7  // 7 真实伤害
	Att_MoveSpeed    = 8  // 8 战场移动速度
	Att_AtbSpeed     = 9  // 9 时间槽速度
	Att_EffectHit    = 10 // 10 技能效果命中
	Att_EffectResist = 11 // 11 技能效果抵抗
	Att_MaxHP        = 12 // 12 生命值上限
	Att_CurHP        = 13 // 13 当前生命值
	Att_MaxMP        = 14 // 14 蓝量上限
	Att_CurMP        = 15 // 15 当前蓝量
	Att_GenMP        = 16 // 16 mp恢复值
	Att_Rage         = 17 // 17 怒气值
	Att_Hit          = 18 // 18 命中
	Att_Dodge        = 19 // 19 闪避

	Att_DmgTypeBegin = 20 // 20 各系伤害加成begin
	Att_DmgPhysics   = 20 // 20 物理系伤害加成
	Att_DmgEarth     = 21 // 21 地系伤害加成
	Att_DmgWater     = 22 // 22 水系伤害加成
	Att_DmgFire      = 23 // 23 火系伤害加成
	Att_DmgWind      = 24 // 24 风系伤害加成
	Att_DmgTime      = 25 // 25 时系伤害加成
	Att_DmgSpace     = 26 // 26 空系伤害加成
	Att_DmgMirage    = 27 // 27 幻系伤害加成
	Att_DmgLight     = 28 // 28 光系伤害加成
	Att_DmgDark      = 29 // 29 暗系伤害加成
	Att_DmgTypeEnd   = 30 // 30 各系伤害加成end

	Att_ResTypeBegin = 30 // 30 各系伤害抗性
	Att_ResPhysics   = 30 // 30 物理系伤害抗性
	Att_ResEarth     = 31 // 31 地系伤害抗性
	Att_ResWater     = 32 // 32 水系伤害抗性
	Att_ResFire      = 33 // 33 火系伤害抗性
	Att_ResWind      = 34 // 34 风系伤害抗性
	Att_ResTime      = 35 // 35 时系伤害抗性
	Att_ResSpace     = 36 // 36 空系伤害抗性
	Att_ResMirage    = 37 // 37 幻系伤害抗性
	Att_ResLight     = 38 // 38 光系伤害抗性
	Att_ResDark      = 39 // 39 暗系伤害抗性
	Att_ResTypeEnd   = 40 // 40 各系伤害抗性

	Att_SecondEnd = 40 // 二级属性结束
	Att_End       = 40 // 32位属性结束
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
	Att_Atk:          "攻击力",
	Att_Armor:        "护甲",
	Att_DmgInc:       "总伤害加成",
	Att_Crit:         "暴击值",
	Att_CritInc:      "暴击倍数加层",
	Att_Tenacity:     "韧性",
	Att_Heal:         "治疗强度",
	Att_RealDmg:      "真实伤害",
	Att_MoveSpeed:    "战场移动速度",
	Att_AtbSpeed:     "时间槽速度",
	Att_EffectHit:    "技能效果命中",
	Att_EffectResist: "技能效果抵抗",
	Att_MaxHP:        "生命值上限",
	Att_CurHP:        "当前生命值",
	Att_MaxMP:        "蓝量上限",
	Att_CurMP:        "当前蓝量",
	Att_GenMP:        "mp恢复值",
	Att_Rage:         "怒气值",
	Att_Hit:          "命中",
	Att_Dodge:        "闪避",

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
