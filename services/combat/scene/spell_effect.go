package scene

// 技能效果处理函数
type SpellEffectsHandler func(*Skill, int32, *SceneEntity) error

var spellEffectsHandlers []SpellEffectsHandler = []SpellEffectsHandler{
	EffectNull,            // 0 无效果
	EffectDamage,          // 1 伤害
	EffectHeal,            // 2 治疗
	EffectAddAura,         // 3 添加aura
	EffectPlacate,         // 4 安抚
	EffectEnrage,          // 5 激怒
	EffectCastSpell,       // 6 施放技能
	EffectDispel,          // 7 驱散
	EffectModAuraDuration, // 8 强化Aura作用时间
	EffectAverageHP,       // 9 平均血量
	EffectAuraNumDmg,      // 10 根据buff数量计算伤害
	EffectTargetAttDamage, // 11 根据目标某一属性计算伤害
	EffectCasterAttDamage, // 12 根据施放者某一属性计算伤害
	EffectDamageRaceMod,   // 13 种族加成伤害
	EffectDispelAndWeak,   // 14 驱散虚弱
	EffectAddLevelAura,    // 15 根据目标等级添加Aura
	EffectLevelEnrage,     // 16 根据目标等级激怒
	EffectAddStateAura,    // 17 添加状态类Aura,并计算状态抗性
	EffectRandAura,        // 18 添加随机buff
	EffectPetDamage,       // 19 宠物伤害
	EffectPetHeal,         // 20 宠物治疗
	EffectChangeRageSpell, // 21 替换英雄怒气技能
	EffectAddWrapAura,     // 22 添加可叠加buff
	EffectModPctCurHP,     // 23 百分比修改目标当前血量
}

// 0 无效果
func EffectNull(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 1 伤害
func EffectDamage(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 2 治疗
func EffectHeal(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 3 添加aura
func EffectAddAura(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 4 安抚
func EffectPlacate(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 5 激怒
func EffectEnrage(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 6 施放技能
func EffectCastSpell(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 7 驱散
func EffectDispel(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 8 强化aura作用时间
func EffectModAuraDuration(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 9 平均血量
func EffectAverageHP(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 10 根据buff数量计算伤害
func EffectAuraNumDmg(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 11 根据目标某一属性计算伤害
func EffectTargetAttDamage(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 12 根据施放者某一属性计算伤害
func EffectCasterAttDamage(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 13 种族加成伤害
func EffectDamageRaceMod(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 14 驱散虚弱
func EffectDispelAndWeak(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 15 根据目标等级添加Aura
func EffectAddLevelAura(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 16 根据目标等级激怒
func EffectLevelEnrage(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 17 添加状态类Aura,并计算状态抗性
func EffectAddStateAura(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 18 添加随机buff
func EffectRandAura(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 19 宠物伤害
func EffectPetDamage(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 20 宠物治疗
func EffectPetHeal(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 21 替换英雄怒气技能
func EffectChangeRageSpell(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 22 添加可叠加buff
func EffectAddWrapAura(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}

// 23 百分比修改目标当前血量
func EffectModPctCurHP(spell *Skill, index int32, target *SceneEntity) error {

	return nil
}
