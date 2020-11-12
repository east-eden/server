package define

const (
	Combat_MaxSpell = 50 // 战斗管理器最多50个技能
	Combat_MaxAura  = 50 // 战斗管理器最多50个buff

	Combat_DmgModTypeNum = 3
)

var Combat_DmgModType = [Combat_DmgModTypeNum]EAuraEffectType{
	AuraEffectType_DmgMod,
	AuraEffectType_DmgFix,
	AuraEffectType_AbsorbAllDmg,
}

//-------------------------------------------------------------------------------
// 行为类型
//-------------------------------------------------------------------------------
type EBehaviourType int32

const (
	BehaviourType_Null        EBehaviourType = -1
	BehaviourType_Begin       EBehaviourType = iota
	BehaviourType_SpellFinish EBehaviourType = iota - 1 // 0 技能结束触发
	BehaviourType_Start                                 // 1 战斗开始触发
	BehaviourType_BeforeMelee                           // 2 普通前触发

	BehaviourType_End
)

//-------------------------------------------------------------------------------
// 免疫类型
//-------------------------------------------------------------------------------
type EImmunityType int32

const (
	ImmunityType_Begin    EImmunityType = iota
	ImmunityType_Damage                 = iota - 1 // 0 伤害免疫
	ImmunityType_Mechanic                          // 1 状态免疫
	ImmunityType_Dispel                            // 2 免疫驱散

	ImmunityType_End
)
