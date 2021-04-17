package scene

import (
	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	log "github.com/rs/zerolog/log"
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
	// 伤害类型
	// damageType := int32(utils.Round(float64(effectEntry.ParameterA)))

	// // 伤害百分比
	// damagePercent := float64(effectEntry.ParameterB) / float64(define.PercentBase)

	// // 伤害固定值
	// damageBase := effectEntry.ParameterC

	// // 忽略防御百分比
	// ignoreDefence := effectEntry.ParameterD

	// // 真实伤害固定值
	// realDamageBase := effectEntry.ParameterE

	// // 最终伤害=(攻击力*技能伤害%+技能固定值) * (1+元素伤害加成%) * (1-护甲伤害减免%(计算忽略防御)) * (1-元素伤害抗性%) * 总伤害系数 * 伤害浮动系数 * (1+) + 真实伤害固定值
	// damage := (float64(s.opts.Caster.GetAttManager().GetFinalAttValue(define.Att_Atk))*
	// 	float64(damagePercent)/float64(define.PercentBase) +
	// 	float64(damageBase))

	// decimal, _ := decimal.NewFromString("136.02")

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
