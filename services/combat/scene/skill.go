package scene

import (
	"container/list"

	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	log "github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

//-------------------------------------------------------------------------------
// 伤害信息
//-------------------------------------------------------------------------------
type CalcDamageInfo struct {
	Type       define.EDmgInfoType // 伤害方式
	SchoolType define.ESchoolType  // 伤害类型
	Damage     int64               // 伤害量
	SpellId    int32               // 技能ID
	ProcCaster int32               // 释放者技能效果掩码
	ProcTarget int32               // 目标技能效果掩码
	ProcEx     int32               // 技能结果类型掩码
	Hit        bool                // 技能是否命中
	Crit       bool                // 技能是否暴击

}

func (d *CalcDamageInfo) Reset() {
	d.Type = define.DmgInfo_Null
	d.SchoolType = define.SchoolType_Null
	d.Damage = 0
	d.SpellId = 0
	d.ProcCaster = 0
	d.ProcTarget = 0
	d.ProcEx = 0
	d.Hit = false
	d.Crit = false
}

type Skill struct {
	opts         *SkillOptions
	scene        *Scene
	listTargets  *list.List // 目标列表list<*SceneUnit>
	listBeatBack *list.List // 反击列表list<*SceneUnit>

	// todo
	baseDamage         int64          // 基础伤害
	damageInfo         CalcDamageInfo // 伤害信息
	effectFlag         uint32         // 效果掩码
	resumeCasterRage   bool
	resumeCasterEnerge bool
	killEntity         bool
	ragePctMod         float32
	procCaster         int32 // 释放者技能效果类型掩码
	procTarget         int32 // 目标技能效果类型掩码
	procEx             int32 // 技能结果掩码
	finalProcCaster    int32
	finalProcEx        int32
	multiple           [define.SpellEffectNum]float32
	curPoint           [define.SpellEffectNum]int32

	// todo
	completed bool // 是否作用结束
}

func (s *Skill) Init(scene *Scene, opts ...SkillOption) {
	s.scene = scene
	s.opts = DefaultSkillOptions()
	s.listTargets = list.New()
	s.listBeatBack = list.New()

	for _, o := range opts {
		o(s.opts)
	}

}

func (s *Skill) GetScene() *Scene {
	return s.scene
}

func (s *Skill) Complete() {
	s.completed = true
}

func (s *Skill) IsCompleted() bool {
	return s.completed
}

func (s *Skill) Update() {
	if s.IsCompleted() {
		return
	}

	// todo
}

// 技能施放
func (s *Skill) Cast() {
	s.findTarget()
	s.sendCastGO()
	s.calcEffect()
	s.sendCastEnd()
	s.castBeatBackSpell()
}

// 技能目标合法检查
func (s *Skill) checkTargetsValid() {
	var next *list.Element
	for e := s.listTargets.Front(); e != nil; e = next {
		next = e.Next()

		// todo 目标种族检查
		// if ((1 << target.opts.Entry.Race) & s.opts.Entry.TargetRace) == 0 {
		// 	return false
		// }

		// 目标状态检查
		// targetState := target.GetState64()
		// targetStateCheckFlag := s.opts.Entry.TargetStateCheckFlag

		// // 释放者处于鹰眼状态
		// if s.opts.Caster.HasState(define.HeroState_AntiHidden) {
		// 	targetStateCheckFlag &= ^uint64(1 << define.HeroState_Stealth)
		// }

		// targetState &= targetStateCheckFlag

		// if s.opts.Entry.TargetStateLimit != targetState {
		// 	s.listTargets.Remove(e)
		// }
	}
}

func (s *Skill) findTarget() {
	s.listTargets.Init()

	// 先把所有单位放进目标列表中，然后再筛选
	entities := s.GetScene().GetEntityMap()
	it := entities.Iterator()
	for it.Next() {
		s.listTargets.PushBack(it.Value())
	}

	// 筛选目标列表
	s.selectTargets()

	// 检查目标是否合法
	s.checkTargetsValid()
}

func (s *Skill) calcEffect() {
	// if s.opts.Caster != nil {
	// 	s.opts.Caster.getCombatCtrl().CalSpellPoint(&s.opts.Entry.SpellBase, s.curPoint[:], s.multiple[:], s.opts.Level)
	// }

	// 计算效果
	for target := s.listTargets.Front(); target != nil; target = target.Next() {
		s.doEffect(target.Value.(*SceneEntity))
	}

	// 回复怒气
	rage := s.opts.Entry.Rage
	if rage > 0 {
		genRagePercent := s.opts.Caster.GetAttManager().GetFinalAttValue(define.Att_GenRagePercent)
		add := genRagePercent.Mul(decimal.NewFromInt32(rage)).Round(0)
		s.opts.Caster.GetAttManager().ModFinalAttValue(define.Att_MaxRage, add)
	}

	// 消耗怒气
	if rage < 0 {
		s.opts.Caster.GetAttManager().ModFinalAttValue(define.Att_MaxRage, decimal.NewFromInt32(rage))
	}

	// 触发子技能
	// if s.opts.Entry.TriggerSpellId > 0 {
	// 	if s.opts.Caster != nil {
	// 		scene := s.GetScene()
	// 		if scene.Rand(1, 10000) <= int(s.opts.Entry.TriggerSpellProp) {
	// 			spellType := s.opts.SpellType
	// 			if s.opts.SpellType < define.SpellType_TriggerNull {
	// 				spellType += define.SpellType_TriggerNull
	// 			}

	// 			// compile comment
	// 			// spellEntry, _ := auto.GetSpellEntry(s.opts.Entry.TriggerSpellId)
	// 			// s.opts.Caster.CombatCtrl().CastSpell(spellEntry, s.opts.Caster, s.opts.Target, false)
	// 		}
	// 	}
	// }

	// 技能施放后触发
	// s.opts.Caster.getCombatCtrl().TriggerByBehaviour(
	// 	define.BehaviourType_SpellFinish,
	// 	s.opts.Target,
	// 	s.finalProcCaster,
	// 	s.finalProcEx,
	// 	s.opts.SpellType,
	// )
}

// 技能效果
func (s *Skill) doEffect(target *SceneEntity) {
	if s.opts.Caster == nil {
		log.Warn().Int32("skill_id", s.opts.Entry.Id).Msg("skill doEffect failed with no caster")
		return
	}

	if target == nil {
		log.Warn().Int32("skill_id", s.opts.Entry.Id).Msg("skill doEffect failed with no target")
		return
	}

	scene := s.GetScene()
	if scene == nil {
		log.Warn().Int32("skill_id", s.opts.Entry.Id).Msg("skill doEffect failed with cannot get caster's scene")
		return
	}

	// 初始化
	s.baseDamage = 0
	s.damageInfo.Reset()
	s.damageInfo.SpellId = int32(s.opts.Entry.Id)
	s.damageInfo.ProcCaster = s.procCaster
	s.damageInfo.ProcTarget = s.procTarget
	s.damageInfo.ProcEx = s.procEx

	// 计算技能结果
	s.calSpellResult(target)

	s.effectFlag = 0

	// 计算技能效果
	for _, timelineId := range s.opts.Entry.TimelineID {
		skillTimelineEntry, ok := auto.GetSkillTimelineEntry(timelineId)
		if !ok {
			log.Error().Caller().Int32("timeline_id", timelineId).Msg("cannot find SkillTimelineEntry")
			continue
		}

		for _, effectId := range skillTimelineEntry.Effects {
			effectEntry, ok := auto.GetSkillEffectEntry(effectId)
			if !ok {
				log.Error().Caller().Int32("effect_id", effectId).Msg("cannot find SkillEffectEntry")
				continue
			}

			// 效果是否可作用于目标
			if !s.checkEffectValid(effectEntry, target) {
				continue
			}

			// 效果命中抵抗概率
			if effectEntry.IsEffectHit == define.SkillEffectHitResistProb {
				effectHit := s.opts.Caster.GetAttManager().GetFinalAttValue(define.Att_EffectHit)
				effectResist := target.GetAttManager().GetFinalAttValue(define.Att_EffectResist)
				hit := effectHit.Sub(effectResist).Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()
				if s.GetScene().Rand(1, define.PercentBase) > int(hit) {
					continue
				}
			} else {
				// effect静态表效果命中概率
				if s.GetScene().Rand(1, define.PercentBase) > int(effectEntry.Prob) {
					continue
				}
			}

			// 技能效果处理
			handleSkillEffect(s, effectEntry, target)
		}
	}

	if s.effectFlag != 0 && s.baseDamage > 0 {
		// 计算伤害
		s.dealHeal(target, s.baseDamage, &s.damageInfo)
		s.dealDamage(target, s.baseDamage, &s.damageInfo)

		// 产生伤害
		target.DoneDamage(s.opts.Caster, &s.damageInfo)

		// 触发信息改变
		if define.DmgInfo_Damage == s.damageInfo.Type && s.damageInfo.Damage > 0 {
			s.damageInfo.ProcTarget |= (1 << define.AuraEvent_Taken_Any_Damage)

			if target.HasState(define.HeroState_Dead) {
				s.damageInfo.ProcCaster |= (1 << define.AuraEvent_Kill)
				s.damageInfo.ProcTarget |= (1 << define.AuraEvent_Killed)
				s.damageInfo.ProcEx |= (1 << define.AuraEvent_Killed)
				s.killEntity = true
				s.finalProcCaster |= (1 << define.AuraEvent_Kill)
			} else {
				// 是否触发反击
				// if s.opts.Entry.BeatBack && (s.damageInfo.ProcEx&(1<<define.AuraEventEx_Block) != 0) {
				// 	s.listBeatBack.PushBack(target)
				// }
			}
		}

		// 发送伤害
		scene.SendDamage(&s.damageInfo)
	}

	if s.damageInfo.ProcTarget != 0 {
		// target.getCombatCtrl().TriggerBySpellResult(false, s.opts.Caster, &s.damageInfo)
	}

	if s.opts.Caster != nil {
		// s.opts.Caster.getCombatCtrl().TriggerBySpellResult(true, target, &s.damageInfo)
	}
}

//-------------------------------------------------------------------------------
// 施放反击技能
//-------------------------------------------------------------------------------
func (s *Skill) castBeatBackSpell() {
	for e := s.listBeatBack.Front(); e != nil; e = e.Next() {
		target := e.Value.(*SceneEntity)
		target.BeatBack(s.opts.Caster)
	}
}

func (s *Skill) sendCastGO() {

}

func (s *Skill) sendCastEnd() {
	// 发送MS_CastGo和MS_CastEnd判断条件需要相同，因为他们是成对生成的
	if s.opts.Caster == nil || s.listTargets.Len() == 0 {
		return
	}

	scene := s.GetScene()
	if scene == nil {
		return
	}

	if scene.IsOnlyRecord() {
		return
	}

	// todo send message
	//CreateSceneProtoMsg(msg, MS_CastEnd,);
	//*msg << (UINT32)m_pCaster->GetLocation();
	//*msg << (UINT32)m_pEntry->dwID;
	//*msg << (INT32)m_pCaster->GetAttController().GetAttValue(EHA_CurRage);
	//*msg << (INT32)(static_cast<EntityGroup*>(m_pCaster->GetFather())->GetEnergy());
	//m_pCaster->GetScene()->AddMsgList(msg);
}

func (s *Skill) calSpellResult(target *SceneEntity) {
	// 反击技能直接命中，不计算暴击和格挡
	if s.opts.SpellType == define.SpellType_TriggerBeatBack {
		s.damageInfo.ProcEx |= (1 << define.AuraEventEx_Normal_Hit)
		return
	}

	// 判断技能是否命中
	if !s.checkSkillHit(target) {
		return
	}

	_ = s.checkSkillCrit(target)
}

func (s *Skill) checkSkillHit(target *SceneEntity) bool {
	// 友方必命中
	if s.opts.Caster.GetCamp() == target.GetCamp() {
		s.damageInfo.Hit = true
		return s.damageInfo.Hit
	}

	// 敌方判断命中和闪避
	hit := s.opts.Caster.GetAttManager().GetFinalAttValue(define.Att_Hit)
	doge := target.GetAttManager().GetFinalAttValue(define.Att_Dodge)
	hitChance := hit.Sub(doge).Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()

	// 保底命中率
	if hitChance < 2000 {
		hitChance = 2000
	}

	s.damageInfo.Hit = int(hitChance) >= s.GetScene().Rand(1, define.PercentBase)
	return s.damageInfo.Hit
}

func (s *Skill) checkSkillCrit(target *SceneEntity) bool {
	critChance := s.opts.Caster.Opts().AttManager.GetFinalAttValue(define.Att_Crit)

	// 敌方计算韧性
	if target.GetCamp().camp != s.opts.Caster.GetCamp().camp {
		critChance.Sub(target.GetAttManager().GetFinalAttValue(define.Att_Tenacity))
	}

	crit := critChance.Mul(decimal.NewFromInt(define.PercentBase)).Round(0).IntPart()
	if crit < 0 {
		crit = 0
	}

	s.damageInfo.Crit = int(crit) >= s.GetScene().Rand(1, define.PercentBase)
	return s.damageInfo.Crit
}

func (s *Skill) calDamage(baseDamage int64, damageInfo *CalcDamageInfo, target *SceneEntity) {
	if target == nil {
		return
	}

	dmgInc := s.opts.Caster.Opts().AttManager.GetFinalAttValue(define.Att_SelfDmgInc)
	baseDamage = dmgInc.Mul(decimal.NewFromInt(baseDamage)).Round(0).IntPart()

	if s.opts.SpellType == define.SpellType_Rage {
		dmgMod := int64(float64(s.ragePctMod) * float64(baseDamage))
		baseDamage += dmgMod
	}

	//FLOAT fPctDmgMod = (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PctDmgInc) - pTarget->GetAttController().GetAttValue(EHA_PctDmgDec));

	//// PVP百分比伤害加成
	//if( m_pCaster->GetScene()->GetStateFlag() & ESSF_PVP )
	//fPctDmgMod += (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PVPPctDmgInc) - pTarget->GetAttController().GetAttValue(EHA_PVPPctDmgDec));

	//// 伤害类型加成
	//if (DamageInfo.eSchool == EIS_Physics)
	//{
	//fPctDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_DamageDonePctPhysics) + pTarget->GetDmgModAtt(EDM_DamageTakenPctPhysics));
	//}
	//else if(DamageInfo.eSchool == EIS_Magic)
	//{
	//fPctDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_DamageDonePctMagic) + pTarget->GetDmgModAtt(EDM_DamageTakenPctMagic));
	//}

	//// 伤害种族加成
	//fPctDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_RaceDoneKindom + pTarget->GetEntry()->eRace));
	//fPctDmgMod += (FLOAT)(pTarget->GetDmgModAtt(EDM_RaceTakenKindom + m_pCaster->GetEntry()->eRace));

	//fPctDmgMod += m_pCaster->GetScene()->GetSceneDmgMod();
	//fPctDmgMod += m_pCaster->GetScene()->GetLevelSuppress(m_pCaster, pTarget);

	//// 判断百分比下限
	//if( fPctDmgMod < -7000.0f )
	//fPctDmgMod = -7000.0f;

	//nBaseDamage += ((fPctDmgMod / 10000.0f) * (FLOAT)nBaseDamage);

	if baseDamage < 1 {
		damageInfo.Damage = 1
	}

	//if (DamageInfo.dwProcEx & EAEE_Critical_Hit)
	//{
	//INT nCrit = m_pCaster->GetAttController().GetAttValue(EHA_CritInc) - pTarget->GetAttController().GetAttValue(EHA_CritDec);
	//nBaseDamage *= (Max(10000, 17000 + nCrit) /10000.0f);
	//}

	//if (DamageInfo.dwProcEx & EAEE_Block)
	//{
	//nBaseDamage *= 0.5f;
	//}

}

func (s *Skill) calHeal(baseHeal int64, damageInfo *CalcDamageInfo, target *SceneEntity) {
	if target == nil {
		return
	}

	// 重伤状态无法加血
	if target.HasState(define.HeroState_Injury) {
		damageInfo.Damage = 0
		return
	}

	// 中毒状态加血效果减75%
	healPct := 1.0
	if target.HasState(define.HeroState_Poison) {
		healPct = 0.25
	}

	if s.opts.SpellType == define.SpellType_Rune {
		damageInfo.Damage = int64(float64(baseHeal) * healPct)
		return
	}

	if s.opts.SpellType == define.SpellType_Rage {
		baseHeal += int64(float64(s.ragePctMod) * float64(baseHeal))
	}

	//FLOAT fDmgMod = (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PctDmgInc));
	//if( m_pCaster->GetScene()->GetStateFlag() & ESSF_PVP )
	//fDmgMod += (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PVPPctDmgInc));

	//fDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_DamageDonePctHeal) + pTarget->GetDmgModAtt(EDM_DamageTakenPcttHeal));
	////	fDmgMod += (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_HealPctIncDone) + pTarget->GetAttController().GetAttValue(EHA_HealPctIncTaken));
	//fDmgMod = fDmgMod / 10000.0f;
	//fDmgMod = fDmgMod * nBaseHeal;
	//nBaseHeal += fDmgMod;

	//if (DamageInfo.dwProcEx & EAEE_Critical_Hit)
	//{
	//INT nCrit = m_pCaster->GetAttController().GetAttValue(EHA_CritInc);
	//nBaseHeal *= (Max(10000, 17000 + nCrit) /10000.0f);
	//}

	if baseHeal < 0 {
		baseHeal = 0
	}

	damageInfo.Damage = int64(float64(baseHeal) * float64(healPct))
}

func (s *Skill) dealDamage(target *SceneEntity, baseDamage int64, damageInfo *CalcDamageInfo) {
	if target == nil {
		return
	}

	if damageInfo.SchoolType == define.SchoolType_Null {
		return
	}

	if damageInfo.Type != define.DmgInfo_Damage {
		return
	}

	s.calDamage(baseDamage, damageInfo, target)

	// 触发双方的伤害接口
	if s.opts.Caster != nil {
		s.opts.Caster.OnDamage(target, damageInfo)
	}

	target.OnBeDamaged(s.opts.Caster, damageInfo)
}

func (s *Skill) dealHeal(target *SceneEntity, baseHeal int64, damageInfo *CalcDamageInfo) {
	if target != nil {
		return
	}

	if damageInfo.SchoolType == define.SchoolType_Null {
		return
	}

	if damageInfo.Type != define.DmgInfo_Heal {
		return
	}

	s.calHeal(baseHeal, damageInfo, target)

	// 触发双方的伤害接口
	if s.opts.Caster != nil {
		s.opts.Caster.OnDamage(target, damageInfo)
	}

	target.OnBeDamaged(s.opts.Caster, damageInfo)

	// 计算有效治疗
	maxHeal := target.opts.AttManager.GetFinalAttValue(define.Att_MaxHP).Sub(target.opts.AttManager.GetFinalAttValue(define.Att_CurHP)).IntPart()
	if maxHeal < damageInfo.Damage {
		damageInfo.Damage = maxHeal
	}
}

//--------------------------------------------------------------------------------------------------
// 效果是否可作用于目标
//--------------------------------------------------------------------------------------------------
func (s *Skill) checkEffectValid(effectEntry *auto.SkillEffectEntry, target *SceneEntity) bool {
	// if s.opts.Entry.Effects[index] == define.SpellEffectType_Null {
	// 	return false
	// }

	// switch s.opts.Entry.EffectsTargetLimit[index][effectIndex] {
	// case define.EffectTargetLimit_Null:
	// 	return true

	// case define.EffectTargetLimit_Self:
	// 	return target == s.opts.Caster

	// case define.EffectTargetLimit_UnSelf:
	// 	return target != s.opts.Caster

	// case define.EffectTargetLimit_Caster_State:
	// 	if s.opts.Caster.HasStateAny(s.opts.Entry.EffectsValidMiscValue[index][effectIndex]) {
	// 		return true
	// 	}

	// case define.EffectTargetLimit_Target_State:
	// 	if target.HasStateAny(s.opts.Entry.EffectsValidMiscValue[index][effectIndex]) {
	// 		return true
	// 	}

	// case define.EffectTargetLimit_Caster_HP_Low:
	// 	hpPct := s.opts.Entry.EffectsValidMiscValue[index][effectIndex]
	// 	if (float64(hpPct) / float64(10000.0) * float64(s.opts.Caster.opts.AttManager.GetAttValue(define.Att_MaxHPBase))) > float64(s.opts.Caster.opts.AttManager.GetAttValue(define.Att_CurHP)) {
	// 		return true
	// 	}

	// case define.EffectTargetLimit_Target_HP_Low:
	// 	hpPct := s.opts.Entry.EffectsValidMiscValue[index][effectIndex]
	// 	if (float64(hpPct) / 10000.0 * float64(target.opts.AttManager.GetAttValue(define.Att_MaxHPBase))) > float64(target.opts.AttManager.GetAttValue(define.Att_CurHP)) {
	// 		return true
	// 	}

	// case define.EffectTargetLimit_Target_HP_High:
	// 	hpPct := s.opts.Entry.EffectsValidMiscValue[index][effectIndex]
	// 	if (float64(hpPct) / 10000.0 * float64(target.opts.AttManager.GetAttValue(define.Att_MaxHPBase))) < float64(target.Opts().AttManager.GetAttValue(define.Att_CurHP)) {
	// 		return true
	// 	}

	// case define.EffectTargetLimit_Pct:
	// 种族概率加成
	//INT32 nBasePct = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
	//if( m_pEntry->eEffectValidRace[nEffectIndex] == pTarget->GetEntry()->eRace )
	//nBasePct += m_pEntry->dwEffectValidRaceMod[nEffectIndex];

	//// 等级概率衰减
	//if( VALID(m_pEntry->nDecayLevel) && pTarget->GetLevel() > m_pEntry->nDecayLevel )
	//{
	//INT32 nLevelDiffer = pTarget->GetLevel() - m_pEntry->nDecayLevel;
	//nBasePct -= nLevelDiffer * m_pEntry->nDecayRate;
	//}

	//if( nBasePct > m_pCaster->GetScene()->GetRandom().Rand(1, 10000) )
	//return TRUE;
	//else
	//{
	//m_DamageInfo.dwProcEx |= EAEE_Invalid;
	//return FALSE;
	//}

	//case EETV_Target_AuraNot:
	//{
	//DWORD dwAuraID = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
	//if( VALID(pTarget->GetCombatController().GetAuraByIDCaster(dwAuraID)) )
	//return FALSE;

	//return TRUE;
	//}
	//break;

	//case EETV_Target_Aura:
	//{
	//DWORD dwAuraID = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
	//if( VALID(pTarget->GetCombatController().GetAuraByIDCaster(dwAuraID)) )
	//return TRUE;
	//}
	//break;

	//case EETV_Target_Race:
	//{
	//// 目标种族限制
	//if (VALID(pTarget) && VALID(m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex]))
	//{
	//if(((1 << pTarget->GetEntry()->eRace) & m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex]))
	//return TRUE;
	//}
	//}
	//break;

	//case EEIV_Caster_AuraState:
	//{
	//INT32 nAuraState = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
	//if( m_pCaster->GetCombatController().HasAuraState(nAuraState) )
	//return TRUE;

	//return FALSE;
	//}
	//break;

	//case EEIV_Target_AuraState:
	//{
	//INT32 nAuraState = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
	//if( pTarget->GetCombatController().HasAuraState(nAuraState) )
	//return TRUE;

	//return FALSE;
	//}
	//break;

	//case EETV_Target_GT_Level:
	//{
	//if (VALID(pTarget) && pTarget->GetLevel() > m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex] )
	//{
	//return TRUE;
	//}
	//}
	//break;

	//case EETV_Target_LT_Level:
	//{
	//if (VALID(pTarget) && pTarget->GetLevel() <= m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex] )
	//{
	//return TRUE;
	//}
	//}
	//break;

	//case EEIV_Caster_AuraPN:
	//{
	//INT32 nEffectPriority = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];

	//INT32	nPosNum		=	0;
	//INT32	nNegNum		=	0;
	//m_pCaster->GetCombatController().GetPositiveAndNegativeNum(nPosNum, nNegNum);
	//if( (nEffectPriority > 0 && nPosNum > 0) || (nEffectPriority < 0 && nNegNum > 0) )
	//return TRUE;

	//return FALSE;
	//}
	//break;

	//case EEIV_Target_AuraPN:
	//{
	//INT32 nEffectPriority = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];

	//INT32	nPosNum		=	0;
	//INT32	nNegNum		=	0;
	//pTarget->GetCombatController().GetPositiveAndNegativeNum(nPosNum, nNegNum);
	//if( (nEffectPriority > 0 && nPosNum > 0) || (nEffectPriority < 0 && nNegNum > 0) )
	//return TRUE;

	//return FALSE;
	//}
	//break;
	// }

	return true
}
