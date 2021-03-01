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

	Att_DmgTypeBegin = 17 // 17 各系伤害加成begin
	Att_DmgPhysics   = 17 // 17 物理系伤害加成
	Att_DmgEarth     = 18 // 18 地系伤害加成
	Att_DmgWater     = 19 // 19 水系伤害加成
	Att_DmgFire      = 20 // 20 火系伤害加成
	Att_DmgWind      = 21 // 21 风系伤害加成
	Att_DmgTime      = 22 // 22 时系伤害加成
	Att_DmgSpace     = 23 // 23 空系伤害加成
	Att_DmgMirage    = 24 // 24 幻系伤害加成
	Att_DmgLight     = 25 // 25 光系伤害加成
	Att_DmgDark      = 26 // 26 暗系伤害加成
	Att_DmgTypeEnd   = 27 // 27 各系伤害加成end

	Att_ResTypeBegin = 27 // 27 各系伤害抗性
	Att_ResPhysics   = 27 // 27 物理系伤害抗性
	Att_ResEarth     = 28 // 28 地系伤害抗性
	Att_ResWater     = 29 // 29 水系伤害抗性
	Att_ResFire      = 30 // 30 火系伤害抗性
	Att_ResWind      = 31 // 31 风系伤害抗性
	Att_ResTime      = 32 // 32 时系伤害抗性
	Att_ResSpace     = 33 // 33 空系伤害抗性
	Att_ResMirage    = 34 // 34 幻系伤害抗性
	Att_ResLight     = 35 // 35 光系伤害抗性
	Att_ResDark      = 36 // 36 暗系伤害抗性
	Att_ResTypeEnd   = 37 // 37 各系伤害抗性

	Att_SecondEnd = 37 // 二级属性结束
	Att_End       = 37 // 32位属性结束
)

const (
	AttNum         int = Att_End - Att_Begin // 32位属性数量
	AttPercentBase     = 10000               // 百分比基数

)
