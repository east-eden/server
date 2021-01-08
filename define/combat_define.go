package define

const (
	Combat_MaxAura = 50 // 战斗管理器最多50个buff

	Combat_DmgModTypeNum = 3
)

var Combat_DmgModType = [Combat_DmgModTypeNum]EAuraEffectType{
	AuraEffectType_DmgMod,
	AuraEffectType_DmgFix,
	AuraEffectType_AbsorbAllDmg,
}

//-------------------------------------------------------------------------------
// 战斗时行动类型
//-------------------------------------------------------------------------------
type ECombatActionType int32

const (
	CombatAction_Begin  ECombatActionType = iota
	CombatAction_Idle   ECombatActionType = iota - 1 // 0 idle
	CombatAction_Attack                              // 1 攻击行动
	CombatAction_Move                                // 2 移动行动

	CombatAction_End
)

//-------------------------------------------------------------------------------
// 行为类型
//-------------------------------------------------------------------------------
type EBehaviourType int32

const (
	BehaviourType_Null         EBehaviourType = -1
	BehaviourType_Begin        EBehaviourType = iota
	BehaviourType_SpellFinish  EBehaviourType = iota - 1 // 0 技能结束触发
	BehaviourType_Start                                  // 1 战斗开始触发
	BehaviourType_BeforeNormal                           // 2 普攻前触发

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
