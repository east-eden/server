package define

import (
	"time"

	"github.com/willf/bitset"
)

// Positive + Negative
const (
	Aura_HeroSyncNum = 30 + 30             // hero增益+减益aura数量
	Aura_MaxPassive  = 5                   // 被动aura数量
	Aura_MaxDuration = 20 * 24 * time.Hour // aura持续最多20天
)

//-------------------------------------------------------------------------------
// aura施放类型
//-------------------------------------------------------------------------------
type EAuraCastingType int32

const (
	AuraCasting_Begin   EAuraCastingType = iota
	AuraCasting_Inter   EAuraCastingType = iota - 1 // 0 间隔作用
	AuraCasting_Persist                             // 1 持续作用
	AuraCasting_Times                               // 2 作用次数限量
	AuraCasting_End
)

//-------------------------------------------------------------------------------
// Aura作用阶段
//-------------------------------------------------------------------------------
type EAuraEffectStep int32

const (
	AuraEffectStep_Begin  EAuraEffectStep = iota
	AuraEffectStep_Apply  EAuraEffectStep = iota - 1 // 0 获得Aura
	AuraEffectStep_Effect                            // 1 Aura作用
	AuraEffectStep_Remove                            // 2 移除Aura
	AuraEffectStep_Check                             // 3 检查效果
	AuraEffect_End
)

//-------------------------------------------------------------------------------
// aura添加结果
//-------------------------------------------------------------------------------
type EAuraAddResult int32

const (
	AuraAddResult_Null EAuraAddResult = -1 // 无效

	AuraAddResult_Begin      EAuraAddResult = iota
	AuraAddResult_Success    EAuraAddResult = iota - 1 // 0 成功
	AuraAddResult_Inferior                             // 1 有更强的Aura
	AuraAddResult_Full                                 // 2 没有空位
	AuraAddResult_Immunity                             // 3 免疫
	AuraAddResult_Resistance                           // 4 抵抗
	AuraAddResult_End
)

//-------------------------------------------------------------------------------
// aura对伤害的影响
//-------------------------------------------------------------------------------
type EAuraEffectForDmg int32

const (
	AuraEffectForDmg_Begin     EAuraEffectForDmg = iota
	AuraEffectForDmg_Null      EAuraEffectForDmg = iota - 1 // 0 无
	AuraEffectForDmg_Transform                              // 1 伤害类型转化
	AuraEffectForDmg_ValueMod                               // 2 伤害值加成减免
	AuraEffectForDmg_Absorb                                 // 3 伤害吸收
	AuraEffectForDmg_End
)

//-------------------------------------------------------------------------------
// aura slot
//-------------------------------------------------------------------------------
type EAuraSlotGroup int32

const (
	AuraSlotGroup_Begin    EAuraSlotGroup = iota
	AuraSlotGroup_Positive EAuraSlotGroup = iota - 1 // 0 增益aura
	AuraSlotGroup_Negative                           // 1 减益aura
	AuraSlotGroup_Passive                            // 2 被动aura
	AuraSlotGroup_End
)

//-------------------------------------------------------------------------------
// aura效果类型
//-------------------------------------------------------------------------------
type EAuraEffectType int32

const (
	AuraEffectType_Begin              EAuraEffectType = iota
	AuraEffectType_Null               EAuraEffectType = iota - 1 // 0 无效果
	AuraEffectType_PeriodDamage                                  // 1 周期伤害
	AuraEffectType_ModAtt                                        // 2 属性改变
	AuraEffectType_Spell                                         // 3 施放技能
	AuraEffectType_State                                         // 4 状态变更
	AuraEffectType_Immunity                                      // 5 免疫变更
	AuraEffectType_DmgMod                                        // 6 伤害转换(总量一定)
	AuraEffectType_NewDmg                                        // 7 生成伤害(原伤害不变)
	AuraEffectType_DmgFix                                        // 8 限量伤害(改变原伤害)
	AuraEffectType_ChgMelee                                      // 9 替换普通攻击
	AuraEffectType_Shield                                        // 10 血盾
	AuraEffectType_DmgAttMod                                     // 11 改变伤害属性
	AuraEffectType_AbsorbAllDmg                                  // 12 全伤害吸收
	AuraEffectType_DmgAccumulate                                 // 13 累计伤害
	AuraEffectType_MeleeSpell                                    // 14 释放普通攻击
	AuraEffectType_LimitAttack                                   // 15 限制攻击力
	AuraEffectType_PerioHeal                                     // 16 周期治疗
	AuraEffectType_ModAttByAlive                                 // 17 根据当前友方存活人数,计算属性改变
	AuraEffectType_ModAttByEnemyAlive                            // 18 根据当前敌方存活人数,计算属性改变

	AuraEffectType_End
)

//-------------------------------------------------------------------------------
// Aura叠加结果
//-------------------------------------------------------------------------------
type EAuraWrapResult int32

const (
	AuraWrapResult_Begin   EAuraWrapResult = iota
	AuraWrapResult_Add     EAuraWrapResult = iota - 1 // 0 追加
	AuraWrapResult_Replace                            // 1 替换
	AuraWrapResult_Wrap                               // 2 叠加
	AuraWrapResult_Invalid                            // 3 加不上

	AuraWrapResult_End
)

//-------------------------------------------------------------------------------
// aura删除方式
//-------------------------------------------------------------------------------
type EAuraRemoveMode int32

const (
	AuraRemoveMode_Begin EAuraRemoveMode = iota
	AuraRemoveMode_Null  EAuraRemoveMode = iota - 1 // 0 无
	AuraRemoveMode_Running
	AuraRemoveMode_Default
	AuraRemoveMode_Replace
	AuraRemoveMode_Cancel
	AuraRemoveMode_Dispel
	AuraRemoveMode_Consume
	AuraRemoveMode_Delete
	AuraRemoveMode_Interrupt
	AuraRemoveMode_Destroy
	AuraRemoveMode_Registered
	AuraRemoveMode_Hangup
	AuraRemoveMode_End
)

const AuraRemoveMode_Removed = 1<<AuraRemoveMode_Default |
	1<<AuraRemoveMode_Replace |
	1<<AuraRemoveMode_Cancel |
	1<<AuraRemoveMode_Dispel |
	1<<AuraRemoveMode_Consume |
	1<<AuraRemoveMode_Delete |
	1<<AuraRemoveMode_Interrupt |
	1<<AuraRemoveMode_Destroy

//-------------------------------------------------------------------------------
// aura sync step
//-------------------------------------------------------------------------------
type EAuraSyncStep int32

const (
	AuraSyncStep_Begin EAuraSyncStep = iota
	AuraSyncStep_Add   EAuraSyncStep = iota - 1
	AuraSyncStep_Update
	AuraSyncStep_Remove
	AuraSyncStep_End
)

//-------------------------------------------------------------------------------
// aura静态属性
//-------------------------------------------------------------------------------
type AuraEntry struct {
	SpellBase     `json:",inline"`
	DispelFlags   uint32 `json:"dispel_flags"`   // 技能效果的所属驱散类型
	DurationFlags uint32 `json:"duration_flags"` // 技能效果所属强化buff时间类型

	OwnerStateCheckFlag   uint32         `json:"owner_state_check_flag"` // 是否判断所有者状态标示
	OwnerStateCheckBitSet *bitset.BitSet `json:"-"`

	OwnerStateLimit       uint32         `json:"owner_state_limit"` // 所有者状态限制
	OwnerStateLimitBitSet *bitset.BitSet `json:"-"`

	AuraCastType    EAuraCastingType `json:"aura_cast_type"` // aura效果施放类型
	AuraState       int32            `json:"aura_state"`     // aura state
	Duration        int32            `json:"duration"`
	EffectTimes     int32            `json:"effect_times"`
	DependCaster    bool             `json:"depend_caster"`   // 依赖施法者
	EffectPriority  int32            `json:"effect_priority"` // 同类效果作用优先
	MultiWrap       bool             `json:"multi_wrap"`
	DecByTarget     bool             `json:"dec_by_target"`     // 是否根据目标等级衰减效果
	RoundUpdateMask uint32           `json:"round_update_mask"` // aura更新回合数掩码(0-31回合)

	Effects   [SpellEffectNum]EAuraEffectType `json:"effects"`    // 技能效果类型
	TriggerId [SpellEffectNum]uint32          `json:"trigger_id"` // 触发器ID

	RemoveEffect [SpellEffectNum]uint32 `json:"remove_effect"` // aura按删除方式作用移除时效果
	TriggerCount [SpellEffectNum]uint32 `json:"trigger_count"` // 可触发次数
	TriggerCd    [SpellEffectNum]int32  `json:"trigger_cd"`    // 触发CD
}

//-------------------------------------------------------------------------------
// aura同步属性
//-------------------------------------------------------------------------------
type AuraInfo struct {
	CasterID  int64  `json:"caster_id"`
	AuraID    uint32 `json:"aura_id"`
	CurTime   int32  `json:"cur_time"`
	TotalTime int32  `json:"total_time"`
	WrapTimes uint32 `json:"wrap_times"`
}
