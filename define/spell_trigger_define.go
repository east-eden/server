package define

//-------------------------------------------------------------------------------
// 触发器类型
//-------------------------------------------------------------------------------
type EAuraTriggerType int32

const (
	AuraTrigger_Begin       EAuraTriggerType = iota
	AuraTrigger_SpellResult EAuraTriggerType = iota - 1
	AuraTrigger_State
	AuraTrigger_Behaviour
	AuraTrigger_AuraState

	AuraTrigger_End
)

//-------------------------------------------------------------------------------
// 技能效果相应的事件类型
//-------------------------------------------------------------------------------
type EAuraEventType int32

const (
	AuraEvent_Begin EAuraEventType = iota
	AuraEvent_None  EAuraEventType = iota - 1

	AuraEvent_Killed // 1 被杀
	AuraEvent_Kill   // 2 击杀

	AuraEvent_Hit       // 3 普通攻击命中
	AuraEvent_Taken_Hit // 4 被普通攻击命中

	AuraEvent_On_Do_Periodic   // 5 造成周期伤害/治疗
	AuraEvent_On_Take_Periodic // 6 受到周期伤害/治疗

	AuraEvent_Rage_Hit      // 7 怒气技能命中
	AuraEvent_Taken_RageHit // 8 被怒气技能命中

	AuraEvent_Rune_Hit       // 9 符文技能命中
	AuraEvent_Taken_Rune_Hit // 10 被符文技能命中

	AuraEvent_Taken_Aura_Trigger // 11 被buff触发技能命中

	AuraEvent_Pet_Hit       // 12 宠物技能命中
	AuraEvent_Taken_Pet_Hit // 13 被宠物技能命中

	AuraEvent_Taken_Any_Damage // 14受到任意伤害

	EAET_End
)

//-------------------------------------------------------------------------------
// 触发技能效果所需状态的类型
//-------------------------------------------------------------------------------
type EAuraEventConditionType int32

const (
	AuraEventCondition_None               EAuraEventConditionType = -1
	AuraEventCondition_Begin              EAuraEventConditionType = iota
	AuraEventCondition_HPLowerFlat        EAuraEventConditionType = iota - 1 // 0 体力平值低于
	AuraEventCondition_HPLowerPct                                            // 1 体力百分比低于
	AuraEventCondition_HPHigherFlat                                          // 2 体力平值高于
	AuraEventCondition_HPHigherPct                                           // 3 体力百分比高于
	AuraEventCondition_AnyUnitState                                          // 4 处于任一状态
	AuraEventCondition_AllUnitState                                          // 5 处于所有状态
	AuraEventCondition_AuraState                                             // 6 处于Aura状态
	AuraEventCondition_WeaponMask                                            // 7 穿着武器类型
	AuraEventCondition_HaveShield                                            // 8 穿着盾牌
	AuraEventCondition_TargetClass                                           // 9 目标类型
	AuraEventCondition_StrongTarget                                          // 10 目标血量比自己多
	AuraEventCondition_TargetAuraState                                       // 11 目标处于Aura状态
	AuraEventCondition_TargetAllUnitState                                    // 12 目标处于所有状态
	AuraEventCondition_TargetAnyUnitState                                    // 13 目标处于任意状态

	AuraEventCondition_End
)

//-------------------------------------------------------------------------------
// 状态转换
//-------------------------------------------------------------------------------
type EStateChangeMode int32

const (
	StateChangeMode_Begin EStateChangeMode = iota
	StateChangeMode_Add   EStateChangeMode = iota - 1
	StateChangeMode_Remove

	StateChangeMode_End
)

//-------------------------------------------------------------------------------
// 触发器
//-------------------------------------------------------------------------------
type AuraTriggerEntry struct {
	ID uint32 `json:"_id"` // 触发器ID

	SpellId       uint32        `json:"spell_id"`        // 触发技能ID
	FamilyMask    uint32        `json:"family_mask"`     // 技能所属分组
	FamilyRace    EHeroRaceType `json:"family_race"`     // 触发技能种族
	DmgInfoType   uint32        `json:"dmg_info_type"`   // 伤害方式掩码
	SchoolType    ESchoolType   `json:"school_type"`     // 伤害类型
	SpellTypeMask uint32        `json:"spell_type_mask"` // 技能类型掩码

	TriggerType  EAuraTriggerType `json:"trigger_type"`  // 触发器类型
	TriggerMisc1 uint32           `json:"trigger_misc1"` // 触发器参数1
	TriggerMisc2 uint32           `json:"trigger_misc2"` // 触发器参数2

	EventProp int32 `json:"event_prop"` // 触发几率

	ConditionType  EAuraEventConditionType `json:"condition_type"`  // 状态条件类型
	ConditionMisc1 int32                   `json:"condition_misc1"` // 条件参数1
	ConditionMisc2 int32                   `json:"condition_misc2"` // 条件参数2

	AddHead bool `json:"add_head"` // 触发器添加到队列头部,不再继续相应相同事件的触发(防止死循环)
}
