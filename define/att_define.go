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
	Att_Heal         = 5  // 5 治疗强度
	Att_RealDmg      = 6  // 6 真实伤害
	Att_MoveSpeed    = 7  // 7 战场移动速度
	Att_AtbSpeed     = 8  // 8 时间槽速度
	Att_EffectHit    = 9  // 9 技能效果命中
	Att_EffectResist = 10 // 10 技能效果抵抗
	Att_MaxHP        = 11 // 11 生命值上限
	Att_CurHP        = 12 // 12 当前生命值
	Att_MaxMP        = 13 // 13 蓝量上限
	Att_CurMP        = 14 // 14 当前蓝量
	Att_GenMP        = 15 // 15 mp恢复值
	Att_Rage         = 16 // 16 怒气值
	Att_Hit          = 17 // 17 命中
	Att_Dodge        = 18 // 18 闪避

	Att_DmgTypeBegin = 19 // 19 各系伤害加成begin
	Att_DmgPhysics   = 19 // 19 物理系伤害加成
	Att_DmgEarth     = 20 // 20 地系伤害加成
	Att_DmgWater     = 21 // 21 水系伤害加成
	Att_DmgFire      = 22 // 22 火系伤害加成
	Att_DmgWind      = 23 // 23 风系伤害加成
	Att_DmgTime      = 24 // 24 时系伤害加成
	Att_DmgSpace     = 25 // 25 空系伤害加成
	Att_DmgMirage    = 26 // 26 幻系伤害加成
	Att_DmgLight     = 27 // 27 光系伤害加成
	Att_DmgDark      = 28 // 28 暗系伤害加成
	Att_DmgTypeEnd   = 29 // 29 各系伤害加成end

	Att_ResTypeBegin = 29 // 29 各系伤害抗性
	Att_ResPhysics   = 29 // 29 物理系伤害抗性
	Att_ResEarth     = 30 // 30 地系伤害抗性
	Att_ResWater     = 31 // 31 水系伤害抗性
	Att_ResFire      = 32 // 32 火系伤害抗性
	Att_ResWind      = 33 // 33 风系伤害抗性
	Att_ResTime      = 34 // 34 时系伤害抗性
	Att_ResSpace     = 35 // 35 空系伤害抗性
	Att_ResMirage    = 36 // 36 幻系伤害抗性
	Att_ResLight     = 37 // 37 光系伤害抗性
	Att_ResDark      = 38 // 38 暗系伤害抗性
	Att_ResTypeEnd   = 39 // 39 各系伤害抗性

	Att_SecondEnd = 39 // 二级属性结束
	Att_End       = 39 // 32位属性结束
)

const (
	AttNum int = Att_End - Att_Begin // 32位属性数量

)
