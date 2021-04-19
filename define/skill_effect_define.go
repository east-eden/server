package define

// 技能效果类型
const (
	// 伤害效果
	SkillEffectDamage_Begin int32 = 101
	SkillEffectDamage       int32 = 101
	SkillEffectDamage_End   int32 = 102

	// 治疗效果
	SkillEffectHeal_Begin int32 = 201
	SkillEffectHeal       int32 = 201
	SkillEffectHeal_End   int32 = 202

	// 打断效果
	SkillEffectInterrupt_Begin int32 = 301
	SkillEffectInterrupt       int32 = 301
	SkillEffectInterrupt_End   int32 = 302

	// 聚集效果
	SkillEffectGather_Begin int32 = 401
	SkillEffectGather       int32 = 401
	SkillEffectGather_End         = 402

	// 加buff
	SkillEffectBuff_Begin int32 = 501
	SkillEffectAddBuff    int32 = 501
	SkillEffectBuff_End         = 502
)
