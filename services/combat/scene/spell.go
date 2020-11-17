package scene

import (
	"container/list"
	"errors"

	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/define"
)

//-------------------------------------------------------------------------------
// 伤害信息
//-------------------------------------------------------------------------------
type CalcDamageInfo struct {
	Type define.EDmgInfoType // 伤害方式
	//tagHeroLocation			stCaster;
	//tagHeroLocation			stTarget;
	SchoolType define.ESchoolType // 伤害类型
	Damage     int64              // 伤害量
	SpellId    uint32             // 技能ID
	ProcCaster uint32
	ProcTarget uint32
	ProcEx     uint32 // 技能结果类型掩码

}

func (d *CalcDamageInfo) Reset() {
	d.Type = define.DmgInfo_Null
	d.SchoolType = define.SchoolType_Null
	d.Damage = 0
	d.SpellId = 0
	d.ProcCaster = 0
	d.ProcTarget = 0
	d.ProcEx = 0
}

type Spell struct {
	opts         *SpellOptions
	listTargets  *list.List // 目标列表list<SceneUnit>
	listBeatBack *list.List // 反击列表list<SceneUnit>

	// todo
	baseDamage         int64          // 基础伤害
	damageInfo         CalcDamageInfo // 伤害信息
	effectFlag         uint32         // 效果掩码
	resumeCasterRage   bool
	resumeCasterEnerge bool
	killEntity         bool
	ragePctMod         float32
	procCaster         uint32
	procTarget         uint32
	procEx             uint32
	finalProcCaster    uint32
	finalProcEx        uint32
	multiple           [define.SpellEffectNum]float32
	curPoint           [define.SpellEffectNum]int32
}

func (s *Spell) prepareTriggerParamOnInit() {
	s.procCaster = uint32(define.AuraEvent_None)
	s.procTarget = uint32(define.AuraEvent_None)
	s.procEx = uint32(define.AuraEventEx_Null)

	switch s.opts.SpellType {
	case define.SpellType_Melee:
		fallthrough
	case define.SpellType_TriggerMelee:
		fallthrough
	case define.SpellType_TriggerBeatBack:
		s.procCaster |= (1 << define.AuraEvent_Hit)
		s.procTarget |= (1 << define.AuraEvent_Taken_Hit)

	case define.SpellType_Rage:
		fallthrough
	case define.SpellType_TriggerRage:
		s.procCaster |= (1 << define.AuraEvent_Rage_Hit)
		s.procTarget |= (1 << define.AuraEvent_Taken_RageHit)

	case define.SpellType_Rune:
		fallthrough
	case define.SpellType_TriggerRune:
		s.procCaster |= (1 << define.AuraEvent_Rune_Hit)
		s.procTarget |= (1 << define.AuraEvent_Taken_Rune_Hit)

	case define.SpellType_TriggerAura:
		s.procTarget |= (1 << define.AuraEvent_Taken_Aura_Trigger)

	case define.SpellType_TriggerAuraTwice:
		s.procTarget |= (1 << define.AuraEvent_Taken_Aura_Trigger)
		s.procEx |= (1 << define.AuraEventEx_Internal_Cant_Trigger)
		s.procEx |= (1 << define.AuraEventEx_Internal_Triggered)

	default:
		break
	}

	if s.opts.Entry.SchoolType == define.SchoolType_Null {
		s.procEx |= (1 << define.AuraEventEx_Not_Active_Spell)
	} else {
		s.procEx |= (1 << define.AuraEventEx_Only_Active_Spell)
	}

	if s.opts.Triggered {
		s.procEx |= (1 << define.AuraEventEx_Internal_Triggered)
	}

	s.finalProcEx = s.procEx
	s.finalProcCaster = s.procCaster
}

func NewSpell(opts ...SpellOption) *Spell {
	s := &Spell{
		opts:         DefaultSpellOptions(),
		listTargets:  list.New(),
		listBeatBack: list.New(),
	}

	for _, o := range opts {
		o(s.opts)
	}

	s.prepareTriggerParamOnInit()

	return s
}

func (s *Spell) checkCasterLimit() error {

	// 判断技能施放者状态
	if s.opts.Entry.CasterStateCheckFlag != 0 {
		if s.opts.Caster == nil {
			return errors.New("spell caster state limit")
		}

		casterState := s.opts.Caster.GetState64() & s.opts.Entry.CasterStateCheckFlag
		if s.opts.Entry.CasterStateLimit != casterState {
			return errors.New("spell caster state limit")
		}
	}

	// 判断施放者aurastate限制
	if s.opts.Entry.CasterAuraState != 0 {
		if s.opts.Caster == nil {
			return errors.New("spell caster state limit")
		}

		if !s.opts.Caster.CombatCtrl().HasAuraStateAny(s.opts.Entry.CasterAuraState) {
			return errors.New("spell caster state limit")
		}
	}

	if s.opts.Entry.CasterAuraStateNot != 0 {
		if s.opts.Caster == nil {
			return errors.New("spell caster state limit")
		}

		if s.opts.Caster.CombatCtrl().HasAuraStateAny(s.opts.Entry.CasterAuraStateNot) {
			return errors.New("spell caster state limit")
		}
	}

	return nil
}

func (s *Spell) checkTargetLimit() error {
	// 选取目标类型不是单体则不判断目标限制
	if s.opts.Entry.SelectType != define.SelectTarget_Enemy_Single {
		return nil
	}

	// 判断技能目标状态
	targetStateCheckFlag := s.opts.Entry.TargetStateCheckFlag
	if targetStateCheckFlag != 0 {
		if s.opts.Target == nil {
			return errors.New("spell target state limit")
		}

		targetState := s.opts.Target.GetState64()

		// 释放者处于鹰眼状态
		if s.opts.Caster.HasState(define.HeroState_AntiHidden) {
			targetStateCheckFlag &= ^uint64(1 << define.HeroState_Stealth)
		}

		targetState &= targetStateCheckFlag
		if s.opts.Entry.TargetStateLimit != targetState {
			return errors.New("spell target state limit")
		}
	}

	// 判断目标aurastate限制
	if s.opts.Entry.TargetAuraState != 0 {
		if s.opts.Target == nil {
			return errors.New("spell target state limit")
		}

		if !s.opts.Target.CombatCtrl().HasAuraStateAny(s.opts.Entry.TargetAuraState) {
			return errors.New("spell target state limit")
		}
	}

	if s.opts.Entry.TargetAuraStateNot != 0 {
		if s.opts.Target == nil {
			return errors.New("spell target state limit")
		}

		if s.opts.Target.CombatCtrl().HasAuraStateAny(s.opts.Entry.TargetAuraStateNot) {
			return errors.New("spell target state limit")
		}
	}

	return nil
}

//-------------------------------------------------------------------------------
// 施放检查
//-------------------------------------------------------------------------------
func (s *Spell) CanCast() error {
	if err := s.checkTargetLimit(); err != nil {
		return err
	}

	if err := s.checkCasterLimit(); err != nil {
		return err
	}

	return nil
}

//-------------------------------------------------------------------------------
// 技能施放
//-------------------------------------------------------------------------------
func (s *Spell) Cast() {
	s.findTarget()
	s.sendCastGO()
	s.calcEffect()
	s.sendCastEnd()
	s.castBeatBackSpell()
}

func (s *Spell) findTarget() {
	s.listTargets.Init()

	// 混乱状态特殊处理
	if s.opts.SpellType == define.SpellType_Melee {
		if s.opts.Caster == nil {
			return
		}

		if s.opts.Caster.HasState(define.HeroState_Chaos) {
			s.findTargetChaos()
			return
		}

		if s.opts.Caster.HasState(define.HeroState_Taunt) {
			s.findTargetEnemySingle()
			return
		}
	}

	switch s.opts.Entry.SelectType {
	case define.SelectTarget_Null:

	case define.SelectTarget_Self:
		s.findTargetSelf()

	case define.SelectTarget_Enemy_Single:
		s.findTargetEnemySingle()

	case define.SelectTarget_Enemy_Single_Back:
		s.findTargetEnemySingleBack()

	case define.SelectTarget_Friend_HP_Min:
		s.findTargetFrienHPMin()

	case define.SelectTarget_Enemy_HP_Max:
		s.findTargetEnemyHPMax()

	case define.SelectTarget_Enemy_Rage_Max:
		s.findTargetEnemyRageMax()

	case define.SelectTarget_Enemy_Column:
		s.findTargetEnemyColumn()

	case define.SelectTarget_Enemy_Frontline:
		s.findTargetEnemyFrontline(true)

	case define.SelectTarget_Enemy_Supporter:
		s.findTargetEnemySupporter(true)

	case define.SelectTarget_Friend_Random:
		s.findTargetFriendRandom()

	case define.SelectTarget_Enemy_Random:
		s.findTargetEnemyRandom()

	case define.SelectTarget_Friend_All:
		s.findTargetFriendAll()

	case define.SelectTarget_Enemy_All:
		s.findTargetEnemyAll()

	case define.SelectTarget_Enemy_Rune:
		s.findTargetEnemyRune()

	case define.SelectTarget_Friend_Rune:
		s.findTargetFriendRune()

	case define.SelectTarget_Next_Attack:
		s.findTargetNextAttack()

	case define.SelectTarget_Friend_Rage_Min:
		s.findTargetFriendRageMin()

	case define.SelectTarget_Enemy_Frontline_Random:
		s.findTargetEnemyFrontLineRandom()

	case define.SelectTarget_Enemy_Backline_Random:
		s.findTargetEnemyBackLineRandom()

	case define.SelectTarget_Friend_Frontline_Random:
		s.findTargetFriendFrontLineRandom()

	case define.SelectTarget_Friend_Backline_Random:
		s.findTargetFriendBackLineRandom()

	case define.SelectTarget_Next_Attack_Row:
		s.findTargetNextAttackRow()

	case define.SelectTarget_Next_Attack_Column:
		s.findTargetNextAttackConlumn()

	case define.SelectTarget_Next_Attack_Border:
		s.findTargetNextAttackBorder()

	case define.SelectTarget_Next_Attack_Explode:
		s.findTargetNextAttackExplode()

	case define.SelectTarget_Caster_Max_Attack:
		s.findCasterMaxAttack()

	case define.SelectTarget_Target_Max_Attack:
		s.findTargetMaxAttack()

	case define.SelectTarget_Enemy_HP_Min:
		s.findEnemyHPMin()
	}
}

func (s *Spell) calcEffect() {
	if s.opts.Caster != nil {
		s.opts.Caster.CombatCtrl().CalSpellPoint(&s.opts.Entry.SpellBase, s.curPoint[:], s.multiple[:], s.opts.Level)
	}

	if s.opts.SpellType == define.SpellType_Rage {
		curRage := s.opts.Caster.Opts().AttManager.GetAttValue(define.Att_Rage)
		var rageThreshold int64 = 100
		if curRage >= rageThreshold+70 {
			s.ragePctMod = 0.6
		} else if curRage >= rageThreshold+35 {
			s.ragePctMod = 0.3
		} else {
			s.ragePctMod = 0.0
		}

		s.opts.Caster.Opts().AttManager.SetBaseAtt(define.Att_Rage, 0)
	}

	// 是否恢复施法者怒气和能量
	if s.opts.SpellType == define.SpellType_Melee {
		s.resumeCasterRage = true
	}

	if s.opts.SpellType == define.SpellType_Melee || s.opts.SpellType == define.SpellType_Rage {
		s.resumeCasterEnerge = true
	}

	// 计算效果
	for target := s.listTargets.Front(); target != nil; target = target.Next() {
		s.doEffect(target.Value.(SceneUnit))
	}

	// 回复怒气
	if s.resumeCasterRage && !s.opts.Caster.HasState(define.HeroState_Seal) {
		s.opts.Caster.Opts().AttManager.ModBaseAtt(define.Att_Rage, 35)
	}

	// 回复符文能量
	if s.resumeCasterEnerge {
		if s.killEntity {
			//m_pCaster->ModeAttEnergy(X_Energe_Dead_Reward);
		}

		if s.opts.SpellType == define.SpellType_Rage {
			//m_pCaster->ModeAttEnergy(X_Energe_Rage);
		}

		if s.opts.SpellType == define.SpellType_Melee {
			//m_pCaster->ModeAttEnergy(X_Energe_Melee);
		}
	}

	if s.opts.Entry.TargetNum == 1 && s.opts.Entry.SelectType != define.SelectTarget_Self && s.listTargets.Len() == 1 {
		s.opts.Target = s.listTargets.Front().Value.(SceneUnit)
	}

	// 触发子技能
	if s.opts.Entry.TriggerSpellId > 0 {
		if s.opts.Caster != nil {
			scene := s.opts.Caster.GetScene()
			if scene.Rand(1, 10000) <= int(s.opts.Entry.TriggerSpellProp) {
				spellType := s.opts.SpellType
				if s.opts.SpellType < define.SpellType_TriggerNull {
					spellType += define.SpellType_TriggerNull
				}

				s.opts.Caster.CombatCtrl().CastSpell(s.opts.Entry.TriggerSpellId, s.opts.Caster, s.opts.Target, false)
			}
		}
	}

	// 技能施放后触发
	s.opts.Caster.CombatCtrl().TriggerByBehaviour(
		define.BehaviourType_SpellFinish,
		s.opts.Target,
		s.finalProcCaster,
		s.finalProcEx,
		s.opts.SpellType,
	)
}

func (s *Spell) doEffect(target SceneUnit) {
	if s.opts.Caster == nil {
		log.Warn().Uint32("spell_id", s.opts.Entry.ID).Msg("spell doEffect failed with no caster")
		return
	}

	if target == nil {
		log.Warn().Uint32("spell_id", s.opts.Entry.ID).Msg("spell doEffect failed with no target")
		return
	}

	scene := s.opts.Caster.GetScene()
	if scene == nil {
		log.Warn().Uint32("spell_id", s.opts.Entry.ID).Msg("spell doEffect failed with cannot get caster's scene")
		return
	}

	// 初始化
	s.baseDamage = 0
	s.damageInfo.Reset()
	//m_DamageInfo.stCaster.nLocation		= m_pCaster->GetLocation();
	//m_DamageInfo.stTarget.nLocation		= pTarget->GetLocation();
	s.damageInfo.SpellId = s.opts.Entry.ID
	s.damageInfo.ProcCaster = s.procCaster
	s.damageInfo.ProcTarget = s.procTarget
	s.damageInfo.ProcEx = s.procEx

	// 计算技能结果
	s.calSpellResult(target)

	s.effectFlag = 0

	// 未命中或闪避或招架
	if s.damageInfo.ProcEx&(1<<define.AuraEventEx_Miss|1<<define.AuraEventEx_Dodge) != 0 {
		scene.SendDamage(&s.damageInfo)
	} else {

		// 计算效果
		hasEffect := false

		// 计算技能免疫
		if !target.HasImmunityAny(define.ImmunityType_Mechanic, s.opts.Entry.MechanicFlags) {
			for i := 0; i < define.SpellEffectNum; i++ {
				eff := s.opts.Entry.Effects[i]

				if eff < define.SpellEffectType_Null || eff >= define.SpellEffectType_End {
					continue
				}

				hasEffect = true

				if target.HasImmunityAny(define.ImmunityType_Mechanic, s.opts.Entry.EffectsMechanic[i]) {
					continue
				}

				if !s.checkEffectValid(int32(i), target, 0) || !s.checkEffectValid(int32(i), target, 1) {
					continue
				}

				// 技能效果处理
				spellEffectsHandlers[eff](s, int32(i), target)
			}
		}

		if hasEffect && s.effectFlag != 0 {
			s.damageInfo.ProcEx |= (1 << define.AuraEventEx_Immnne)
			scene.SendDamage(&s.damageInfo)
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
				if s.opts.Entry.BeatBack && (s.damageInfo.ProcEx&(1<<define.AuraEventEx_Block) != 0) {
					s.listBeatBack.PushBack(target)
				}
			}
		}

		// 发送伤害
		scene.SendDamage(&s.damageInfo)
	}

	if s.damageInfo.ProcTarget != 0 {
		target.CombatCtrl().TriggerBySpellResult(false, s.opts.Caster, &s.damageInfo)
	}

	if s.opts.Caster != nil {
		s.opts.Caster.CombatCtrl().TriggerBySpellResult(true, target, &s.damageInfo)
	}
}

//-------------------------------------------------------------------------------
// 施放反击技能
//-------------------------------------------------------------------------------
func (s *Spell) castBeatBackSpell() {
	for e := s.listBeatBack.Front(); e != nil; e = e.Next() {
		target := e.Value.(SceneUnit)
		target.BeatBack(s.opts.Caster)
	}
}

func (s *Spell) isTargetValid(target SceneUnit) bool {
	if target == nil {
		return false
	}

	// 是否包括自己
	if !s.opts.Entry.IncludeSelf && target == s.opts.Caster {
		return false
	}

	// 目标种族检查
	if ((1 << target.Opts().Entry.Race) & s.opts.Entry.TargetRace) == 0 {
		return false
	}

	// 目标状态检查
	targetState := target.GetState64()
	targetStateCheckFlag := s.opts.Entry.TargetStateCheckFlag

	// 释放者处于鹰眼状态
	if s.opts.Caster.HasState(define.HeroState_AntiHidden) {
		targetStateCheckFlag &= ^uint64(1 << define.HeroState_Stealth)
	}

	targetState &= targetStateCheckFlag

	if s.opts.Entry.TargetStateLimit != targetState {
		return false
	}

	return true
}

func (s *Spell) sendCastGO() {
	if s.opts.Caster == nil || s.listTargets.Len() == 0 || !s.opts.Entry.HaveVisual {
		return
	}

	scene := s.opts.Caster.GetScene()
	if scene == nil {
		return
	}

	if scene.IsOnlyRecord() {
		return
	}

	// todo send message
	//CreateSceneProtoMsg(msg, MS_CastGO,);
	//*msg << (UINT32)m_pCaster->GetLocation();
	//*msg << (UINT32)m_pEntry->dwID;
	//*msg << (INT32)m_listTarget.Size();

	//EntityHero* pTarget = NULL;
	//TargetList::Iterator it = m_listTarget.Begin();
	//while( m_listTarget.PeekNext(it, pTarget) )
	//{
	//if (!VALID(pTarget))
	//continue;

	//*msg << (UINT32)pTarget->GetLocation();
	//}

	//pScene->AddMsgList(msg);
}

func (s *Spell) sendCastEnd() {
	// 发送MS_CastGo和MS_CastEnd判断条件需要相同，因为他们是成对生成的
	if s.opts.Caster == nil || s.listTargets.Len() == 0 || !s.opts.Entry.HaveVisual {
		return
	}

	scene := s.opts.Caster.GetScene()
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

func (s *Spell) calSpellResult(target SceneUnit) {
	// 群体伤害
	if s.opts.Entry.GroupDmg {
		s.damageInfo.ProcEx |= (1 << define.AuraEventEx_GroupDmg)
		s.finalProcEx |= (1 << define.AuraEventEx_GroupDmg)
	}

	// 反击技能直接命中，不计算暴击和格挡
	if s.opts.SpellType == define.SpellType_TriggerBeatBack {
		s.damageInfo.ProcEx |= (1 << define.AuraEventEx_Normal_Hit)
		return
	}

	// 判断技能是否命中
	if !s.isSpellHit(target) {
		s.damageInfo.ProcEx |= (1 << define.AuraEventEx_Miss)
		return
	}

	if s.isSpellCrit(target) {
		s.damageInfo.ProcEx |= (1 << define.AuraEventEx_Critical_Hit)
		s.finalProcEx |= (1 << define.AuraEventEx_Critical_Hit)
	}

	// 判断技能是否被格挡
	if s.isSpellBlock(target) {
		s.damageInfo.ProcEx |= (1 << define.AuraEventEx_Block)
		s.finalProcEx |= (1 << define.AuraEventEx_Block)
	}

	// 忽略护甲和抗性
	//if( m_pEntry->bCanNotArmor )
	//{
	//	DamageInfo.dwProcEx |= EAEE_IgnoreArmor;
	//	DamageInfo.dwProcEx |= EAEE_IgnorResistance;
	//}

	// 是否恢复目标怒气和能量
	if !target.HasState(define.HeroState_UnBeat) &&
		!(target.HasState(define.HeroState_ImmunityGroupDmg) && s.opts.Entry.GroupDmg) {
		if s.opts.SpellType == define.SpellType_Melee {
			if s.opts.Caster.GetCamp() != target.GetCamp() {
				s.damageInfo.ProcEx |= (1 << define.AuraEventEx_RageResume)
			}
		}

		if s.opts.SpellType == define.SpellType_Melee || s.opts.SpellType == define.SpellType_Rage {
			if s.opts.Caster.GetCamp() != target.GetCamp() {
				s.damageInfo.ProcEx |= (1 << define.AuraEventEx_EnergyResume)
			}
		}
	}

	s.damageInfo.ProcEx |= (1 << define.AuraEventEx_Normal_Hit)
	s.finalProcEx |= (1 << define.AuraEventEx_Normal_Hit)
}

func (s *Spell) isSpellHit(target SceneUnit) bool {
	if target == nil {
		return false
	}

	if s.opts.Entry.HitCertainly {
		return true
	}

	if s.opts.Caster == nil {
		return false
	}

	scene := s.opts.Caster.GetScene()
	if scene == nil {
		return false
	}

	// 友方必命中
	if s.opts.Caster.GetCamp() == target.GetCamp() {
		return true
	}

	hitChance := s.opts.Caster.Opts().AttManager.GetAttValue(define.Att_Hit) - target.Opts().AttManager.GetAttValue(define.Att_Dodge)
	hitChance += int64(s.opts.Entry.SpellHit)

	if hitChance < 5000 {
		hitChance = hitChance/2 + 2500
	}

	// 保底命中率
	if hitChance < 2000 {
		hitChance = 2000
	}

	return int(hitChance) >= scene.Rand(1, 10000)
}

func (s *Spell) isSpellCrit(target SceneUnit) bool {
	if s.opts.Entry.NotCrit {
		return false
	}

	if target == nil {
		return false
	}

	if s.opts.Caster == nil {
		return false
	}

	scene := s.opts.Caster.GetScene()
	if scene == nil {
		return false
	}

	critChance := s.opts.Caster.Opts().AttManager.GetAttValue(define.Att_CriProb)

	// todo 韧性属性
	// 敌方才算韧性
	//if( m_pCaster->GetCamp() != pTarget->GetCamp() )
	//nCritChance -= pTarget->GetAttController().GetAttValue(EHA_Resilience);
	//nCritChance += m_pEntry->nSpellCrit;

	//if( nCritChance > 5000 )
	//{
	//nCritChance = nCritChance /2 + 2500;
	//}

	//nCritChance = Max(0, nCritChance);
	//nCritChance = Min(9000, nCritChance);

	return int(critChance) >= scene.Rand(1, 10000)
}

func (s *Spell) isSpellBlock(target SceneUnit) bool {
	if target == nil {
		return false
	}

	if s.opts.Entry.NotBlock {
		return false
	}

	if s.opts.Caster == nil {
		return false
	}

	// 友方不格挡
	if s.opts.Caster.GetCamp() == target.GetCamp() {
		return false
	}

	scene := s.opts.Caster.GetScene()
	if scene == nil {
		return false
	}

	blockChance := target.Opts().AttManager.GetAttValue(define.Att_Block) - s.opts.Caster.Opts().AttManager.GetAttValue(define.Att_Broken)
	blockChance -= int64(s.opts.Entry.SpellBroken)
	if blockChance > 5000 {
		blockChance = blockChance/2 + 2500
	}

	if blockChance < 0 {
		blockChance = 0
	}

	if blockChance > 9000 {
		blockChance = 9000
	}

	return int(blockChance) >= scene.Rand(1, 10000)
}

func (s *Spell) calDamage(baseDamage int64, damageInfo *CalcDamageInfo, target SceneUnit) {
	if target == nil {
		return
	}

	if s.opts.SpellType == define.SpellType_Rune || s.opts.Entry.CanNotArmor {
		damageInfo.Damage = int64(baseDamage)
		return
	}

	baseDamage += int64(s.opts.Caster.Opts().AttManager.GetAttValue(define.Att_DmgInc)) - int64(s.opts.Target.Opts().AttManager.GetAttValue(define.Att_DmgDec))

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

	minDmg := int64(float64(s.opts.Caster.Opts().AttManager.GetAttValue(define.Att_Atk)) * 0.05)
	if baseDamage < minDmg {
		damageInfo.Damage = minDmg
	}
}

func (s *Spell) calHeal(baseHeal int64, damageInfo *CalcDamageInfo, target SceneUnit) {
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

	if s.opts.Entry.CanNotArmor {
		damageInfo.Damage = baseHeal
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

func (s *Spell) dealDamage(target SceneUnit, baseDamage int64, damageInfo *CalcDamageInfo) {
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

func (s *Spell) dealHeal(target SceneUnit, baseHeal int64, damageInfo *CalcDamageInfo) {
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
	maxHeal := target.Opts().AttManager.GetAttValue(define.Att_MaxHP) - target.Opts().AttManager.GetAttValue(define.Att_CurHP)
	if maxHeal < damageInfo.Damage {
		damageInfo.Damage = maxHeal
	}
}

//--------------------------------------------------------------------------------------------------
// 效果是否可作用于目标
//--------------------------------------------------------------------------------------------------
func (s *Spell) checkEffectValid(effectIndex int32, target SceneUnit, index int32) bool {
	if s.opts.Entry.Effects[index] == define.SpellEffectType_Null {
		return false
	}

	switch s.opts.Entry.EffectsTargetLimit[index][effectIndex] {
	case define.EffectTargetLimit_Null:
		return true

	case define.EffectTargetLimit_Self:
		return target == s.opts.Caster

	case define.EffectTargetLimit_UnSelf:
		return target != s.opts.Caster

	case define.EffectTargetLimit_Caster_State:
		if s.opts.Caster.HasStateAny(s.opts.Entry.EffectsValidMiscValue[index][effectIndex]) {
			return true
		}

	case define.EffectTargetLimit_Target_State:
		if target.HasStateAny(s.opts.Entry.EffectsValidMiscValue[index][effectIndex]) {
			return true
		}

	case define.EffectTargetLimit_Caster_HP_Low:
		hpPct := s.opts.Entry.EffectsValidMiscValue[index][effectIndex]
		if (float64(hpPct) / float64(10000.0) * float64(s.opts.Caster.Opts().AttManager.GetAttValue(define.Att_MaxHP))) > float64(s.opts.Caster.Opts().AttManager.GetAttValue(define.Att_CurHP)) {
			return true
		}

	case define.EffectTargetLimit_Target_HP_Low:
		hpPct := s.opts.Entry.EffectsValidMiscValue[index][effectIndex]
		if (float64(hpPct) / 10000.0 * float64(target.Opts().AttManager.GetAttValue(define.Att_MaxHP))) > float64(target.Opts().AttManager.GetAttValue(define.Att_CurHP)) {
			return true
		}

	case define.EffectTargetLimit_Target_HP_High:
		hpPct := s.opts.Entry.EffectsValidMiscValue[index][effectIndex]
		if (float64(hpPct) / 10000.0 * float64(target.Opts().AttManager.GetAttValue(define.Att_MaxHP))) < float64(target.Opts().AttManager.GetAttValue(define.Att_CurHP)) {
			return true
		}

	case define.EffectTargetLimit_Pct:
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
	}

	return false
}
