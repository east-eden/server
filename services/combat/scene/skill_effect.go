package scene

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/utils/random"
	log "github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

// 技能效果处理函数
type SkillEffectsHandler func(*Skill, *auto.SkillEffectEntry, *SceneEntity)

var (
	skillEffectsHandlers map[int32]SkillEffectsHandler
)

func init() {
	skillEffectsHandlers = make(map[int32]SkillEffectsHandler)

	register()
}

func register() {
	skillEffectsHandlers[define.SkillEffectDamage] = effectDamage
	skillEffectsHandlers[define.SkillEffectHeal] = effectHeal
	skillEffectsHandlers[define.SkillEffectInterrupt] = effectInterrupt
	skillEffectsHandlers[define.SkillEffectGather] = effectGather
	skillEffectsHandlers[define.SkillEffectAddBuff] = effectAddBuff
}

func handleSkillEffect(s *Skill, effectEntry *auto.SkillEffectEntry, target *SceneEntity) {
	h, ok := skillEffectsHandlers[effectEntry.EffectType]
	if !ok {
		log.Error().Caller().Int32("effect_entry", effectEntry.Id).Msg("invalid skill effect type")
		return
	}

	h(s, effectEntry, target)
}

// 101 造成伤害
func effectDamage(s *Skill, effectEntry *auto.SkillEffectEntry, target *SceneEntity) {
	globalConfig, _ := auto.GetGlobalConfig()

	// 伤害类型
	damageType := int(effectEntry.ParameterInt[0])

	// 伤害百分比
	damagePercent := effectEntry.ParameterNumber[0]

	// 伤害固定值
	damageBase := effectEntry.ParameterNumber[1]

	// 忽略防御百分比
	ignoreDefence := effectEntry.ParameterNumber[2]

	// 真实伤害固定值
	realDamageBase := effectEntry.ParameterNumber[3]

	// 元素伤害加成
	dmgInc := s.opts.Caster.GetAttManager().GetFinalAttValue(int(define.Att_DmgTypeBegin) + damageType)

	// 元素伤害抗性
	dmgRes := target.GetAttManager().GetFinalAttValue(int(define.Att_ResTypeBegin) + damageType)

	// 攻击力
	atk := s.opts.Caster.GetAttManager().GetFinalAttValue(define.Att_Atk)

	// 护甲
	armor := target.GetAttManager().GetFinalAttValue(define.Att_Armor)

	// 护甲伤害减免 = 防御方最终面板护甲*(1-忽略防御%)/(防御方最终面板护甲*(1-忽略防御%)+攻击方最终攻击*护甲减免常数)
	armorDec := func() decimal.Decimal {
		armorPartA := armor.Mul(decimal.NewFromInt32(1).Sub(ignoreDefence))
		armorPartB := armorPartA.Add(atk.Mul(decimal.NewFromInt32(globalConfig.ArmorRatio)))
		return armorPartA.Div(armorPartB)
	}()

	// 最终伤害=(攻击力*技能伤害%+技能固定值) * (1+元素伤害加成%) * (1-护甲伤害减免) * (1-元素伤害抗性%) * 总伤害系数 * 伤害浮动系数 * (1+暴击伤害%) + 真实伤害固定值
	partA := atk.Mul(damagePercent).Add(damageBase)
	partB := decimal.NewFromInt32(1).Add(dmgInc)
	partC := decimal.NewFromInt32(1).Sub(armorDec)
	partD := decimal.NewFromInt32(1).Sub(dmgRes)
	partE := s.opts.Caster.GetAttManager().GetFinalAttValue(define.Att_SelfDmgInc)
	partF := random.DecimalFake(globalConfig.DamageRange[0], globalConfig.DamageRange[1], s.GetScene().GetRand())
	partG := func() decimal.Decimal {
		critDamage := decimal.NewFromInt32(1)
		if s.damageInfo.Crit {
			critInc := s.opts.Caster.GetAttManager().GetFinalAttValue(define.Att_CritInc)
			critDamage.Add(critInc)
		}
		return critDamage
	}()
	partH := realDamageBase

	s.baseDamage = partA.Mul(partB).Mul(partC).Mul(partD).Mul(partE).Mul(partF).Mul(partG).Add(partH).Round(0).IntPart()
}

// 201 治疗效果
func effectHeal(s *Skill, effectEntry *auto.SkillEffectEntry, target *SceneEntity) {
}

// 301 打断效果
func effectInterrupt(s *Skill, effectEntry *auto.SkillEffectEntry, target *SceneEntity) {
}

// 401 聚集效果
func effectGather(s *Skill, effectEntry *auto.SkillEffectEntry, target *SceneEntity) {
}

// 501 加buff
func effectAddBuff(s *Skill, effectEntry *auto.SkillEffectEntry, target *SceneEntity) {
}
