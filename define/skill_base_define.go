package define

// 技能类型
const (
	SkillType_Begin    int32 = 1
	SkillType_General  int32 = 1 // 普通攻击
	SkillType_Normal   int32 = 2 // 一般技能
	SkillType_Ultimate int32 = 3 // 奥义技能
	SkillType_Crystal  int32 = 4 // 残响技能
	SkillType_Passive  int32 = 5 // 被动技能
	SkillType_Channel  int32 = 6 // 引导技能
	SkillType_End      int32 = 7
)

// 目标类型
const (
	SkillTargetType_Begin          int32 = 0
	SkillTargetType_SelfRound      int32 = 0 // 自身周围
	SkillTargetType_SelectRound    int32 = 1 // 选定空间
	SkillTargetType_FriendlySingle int32 = 2 // 友军单体
	SkillTargetType_EnemySingle    int32 = 3 // 敌军单体
	SkillTargetType_End                  = 4
)

// 范围类型
const (
	SkillRangeType_Begin     int32 = 0
	SkillRangeType_Single    int32 = 0 // 默认单体
	SkillRangeType_Circle    int32 = 1 // 圆形
	SkillRangeType_Rectangle int32 = 2 // 矩形
	SkillRangeType_Fan       int32 = 3 // 扇形
	SkillRangeType_End       int32 = 4
)

// 发起类型
const (
	SkillLaunchType_Begin  int32 = 1
	SkillLaunchType_Self   int32 = 1 // 以自身发起
	SkillLaunchType_Target int32 = 2 // 以目标发起
	SkillLaunchType_End    int32 = 3
)

// 技能作用范围
const (
	SkillScopeType_Begin             int32 = 1
	SkillScopeType_SelectTarget      int32 = 1 // 选定目标
	SkillScopeType_FriendlyTroops    int32 = 2 // 友军(除目标)
	SkillScopeType_AllFriendlyTroops int32 = 3 // 所有友军
	SkillScopeType_EnemyTroops       int32 = 4 // 敌军(除目标)
	SkillScopeType_AllEnemyTroops    int32 = 5 // 所有敌军
	SkillScopeType_End               int32 = 6
)

// 伤害类型
const (
	SkillDamageType_Begin   int32 = 0
	SkillDamageType_Physics int32 = 0 // 物理属性伤害
	SkillDamageType_Earth   int32 = 1 // 地属性伤害
	SkillDamageType_Water   int32 = 2 // 水属性伤害
	SkillDamageType_Fire    int32 = 3 // 火属性伤害
	SkillDamageType_Wind    int32 = 4 // 风属性伤害
	SkillDamageType_Thunder int32 = 5 // 雷属性伤害
	SkillDamageType_Time    int32 = 6 // 时属性伤害
	SkillDamageType_Space   int32 = 7 // 空属性伤害
	SkillDamageType_End           = 8
)
