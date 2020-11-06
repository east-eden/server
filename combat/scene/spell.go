package scene

import "github.com/yokaiio/yokai_server/define"

type Spell struct {
	opts *SpellOptions
	mapTarget map[uint64]SceneUnit 		// 目标列表
	mapBeatBack map[uint64]SceneUnit 	// 反击列表

	// todo
	resumeCasterRage bool
	resumeCasterEnerge bool
	killEntity bool
	ragePctMod float32
	procCaster uint32
	procTarget uint32
	procEx uint32
	finalProcCaster uint32
	finalProcEx uint32
	multiple float32[define.SpellEffectNum]
	curPoint int32[define.SpellEffectNum]
}

func (s *Spell) prepareTriggerParamOnInit() {
	s.procCaster = define.AuraEvent_None
	s.procTarget = define.AuraEvent_None
	s.procEx = define.AuraEventEx_None

	switch(s.opts.SpellType) {
	case define.SpellType_Melee:
		fallthrough
	case define.SpellType_TriggerMelee:
		fallthrough
	case define.SpellType_BeatBack:
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

	case define.SpellType_AuraTrigger:
		s.procTarget |= (1 << define.AuraEvent_Taken_Aura_Trigger)

	case define.SpellType_AuraTriggerTwice:
			s.procTarget |= (1 << define.AuraEvent_Taken_Aura_Trigger)
			s.procEx |= (1 << define.AuraEventEx_Internal_Cant_Trigger)
			s.procEx |= (1 << define.AuraEventEx_Internal_Triggered)

	default:
		break;
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

func NewSpell(spellId int32, opts ...SpellOption) *Spell {
	s := &Spell{
		opts:     DefaultSpellOptions(),
		mapTarget: make(map[uint64]SceneUnit, 0),
		mapBeatBack: make(map[uint64]SceneUnit, 0),
	}

	for _, o := range opts {
		o(s.opts)
	}

	s.prepareTriggerParamOnInit();

	return s
}

func (s *Spell) checkCasterLimit() error {
	return nil
	// 判断技能施放者状态
	//DWORD	dwCasterStateCheckFlag = m_pEntry->dwCasterStateCheckFlag;
	//if (dwCasterStateCheckFlag != 0)
	//{
		//if (!VALID(pCaster))
			//return E_Spell_Caster_State_Limit;

		//DWORD	dwCasterState = pCaster->GetState();
		//dwCasterState &= dwCasterStateCheckFlag;

		//if( m_pEntry->dwCasterStateLimit != dwCasterState )
			//return E_Spell_Caster_State_Limit;
	//}

	//// 判断施放者aurastate限制
	//if( VALID(m_pEntry->dwCasterAuraState) )
	//{
		//if (!VALID(pCaster))
			//return E_Spell_Caster_State_Limit;

		//if( !pCaster->GetCombatController().HasAuraStateAny(m_pEntry->dwCasterAuraState) )
			//return E_Spell_Caster_State_Limit;
	//}

	//if( VALID(m_pEntry->dwCasterAuraStateNot) )
	//{
		//if (!VALID(pCaster))
			//return E_Spell_Caster_State_Limit;

		//if( pCaster->GetCombatController().HasAuraStateAny(m_pEntry->dwCasterAuraStateNot) )
			//return E_Spell_Caster_State_Limit;
	//}

	//return E_Success;
}

func (s *Spell) checkTargetLimit() error {
	return nil
	// 选取目标类型不是单体则不判断目标限制
	//if(m_pEntry->eSelectType != ESTT_Enemy_Single)
		//return E_Success;

	//// 判断技能目标状态
	//DWORD	dwTargetStateCheckFlag = m_pEntry->dwTargetStateCheckFlag;
	//if (dwTargetStateCheckFlag != 0)
	//{
		//if (!VALID(pTarget))
			//return E_Spell_Target_State_Limit;

		//DWORD dwTargetState = pTarget->GetState();

		//// 释放者处于鹰眼状态
		//if( m_pCaster->HasState(EHS_AntiHidden) )
		//{
			//dwTargetStateCheckFlag &= ~(1 << EHS_Stealth);
		//}

		//dwTargetState &= dwTargetStateCheckFlag;

		//if( m_pEntry->dwTargetStateLimit != dwTargetState )
			//return E_Spell_Target_State_Limit;
	//}

	//// 判断目标aurastate限制
	//if( VALID(m_pEntry->dwTargetAuraState) )
	//{
		//if (!VALID(pTarget))
			//return E_Spell_Target_State_Limit;

		//if( !pTarget->GetCombatController().HasAuraStateAny(m_pEntry->dwTargetAuraState) )
			//return E_Spell_Target_State_Limit;
	//}

	//if( VALID(m_pEntry->dwTargetAuraStateNot) )
	//{
		//if (!VALID(pTarget))
			//return E_Spell_Target_State_Limit;

		//if( pTarget->GetCombatController().HasAuraStateAny(m_pEntry->dwTargetAuraStateNot) )
			//return E_Spell_Target_State_Limit;
	//}

	//return E_Success;
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
	FindTarget()
	SendCastGO()
	CalcEffect()
	SendCastEnd()
	CastBeatBackSpell()
}

func (s *Spell) FindTarget() {
	s.mapTarget = make(map[uint64]SceneUnit, 0)

	// 混乱状态特殊处理
	if s.opts.SpellType == define.SpellType_Melee {
		if s.opts.caster == nil {
			return
		}

		if s.opts.caster.HasState(define.HeroState_Chaos) {
			s.findTargetChaos()
			return
		}

		if s.opts.caster.HasState(define.HeroState_Taunt) {
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
		s.findTargetEnemyFrontline()

	case define.SelectTarget_Enemy_Supporter:
		s.findTargetEnemySupporter()

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
		s.findTargetNextAttackRow();

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

func (s *Spell) CalcEffect() {
	// 计算效果参数
	if (VALID(m_pCaster))
	{
		m_pCaster->GetCombatController().CalSpellPoint(m_pEntry, m_nCurPoint, m_fMultiple, m_nLevel);
	}

	if( m_eSpellType == ERMT_Rage )
	{
		FLOAT fCurRage = (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_CurRage));
		FLOAT fRageThreshold = (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_RageThreshold));
		if( fCurRage >= fRageThreshold + 70)
		{
			// 超过怒气阈值70点怒气增加60%伤害
			m_fRagePctMod =  0.6f;
		}
		else if( fCurRage >= fRageThreshold + 35 )
		{
			// 超过怒气阈值35点怒气增加30%伤害
			m_fRagePctMod =  0.3f;
		}

		m_pCaster->GetAttController().SetAttValue(EHA_CurRage, 0);
	} 

	// 是否恢复施法者怒气和能量
	m_bResumeCasterRage = (m_eSpellType == ERMT_Melee);
	m_bResumeCasterEnerge = (m_eSpellType == ERMT_Melee || m_eSpellType == ERMT_Rage);

	// 计算效果
	EntityHero* pTarget = NULL;
	TargetList::Iterator it = m_listTarget.Begin();
	while (m_listTarget.PeekNext(it, pTarget))
	{
		if (!VALID(pTarget))
			continue;

		DoEffect(pTarget);
	}

	// 计算释放者怒气消耗
	if( m_bResumeCasterRage && !m_pCaster->HasState(EHS_Seal) )
	{
		m_pCaster->GetAttController().ModAttValue(EHA_CurRage, X_Rage_Resume);
	}

	if( m_bResumeCasterEnerge )
	{
		if( m_bKillEntity )
		{
			m_pCaster->ModeAttEnergy(X_Energe_Dead_Reward);
		}
		
		if( m_eSpellType == ERMT_Rage )
		{
			m_pCaster->ModeAttEnergy(X_Energe_Rage);
		}

		if( m_eSpellType == ERMT_Melee )
		{
			m_pCaster->ModeAttEnergy(X_Energe_Melee);
		}
	}

	if( m_pEntry->nTargetNum == 1 && m_pEntry->eSelectType != ESTT_Self && m_listTarget.Size() == 1)
	{
		m_pTarget = m_listTarget.Front();
	}

	if( VALID(m_pEntry->dwTriggerSpellID) )
	{
		Scene* pScene = m_pCaster->GetScene();
		if( VALID(m_pCaster) && (pScene->GetRandom().Rand(1, 10000) <= m_pEntry->nTriggerSpellProp))
		{
			INT32 eSpellType =  (m_eSpellType < ERMT_TriggerNull) ? (m_eSpellType + ERMT_TriggerNull) : m_eSpellType;

			m_pCaster->GetCombatController().CastSpell(m_pEntry->dwTriggerSpellID, m_pCaster, m_pTarget, FALSE, m_nAmount, (ESpellType)eSpellType, m_nLevel);
		}
	}

	m_pCaster->GetCombatController().TriggerByBehaviour(EBT_SpellFinish, m_pTarget, m_dwFinalProcCaster, m_dwFinalProcEx, m_eSpellType);
}

VOID Spell::DoEffect( EntityHero* pTarget )
{
	if (!VALID(m_pCaster) || !VALID(pTarget))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	// 初始化
	m_nBaseDamage						= 0;
	m_DamageInfo.Reset();
	m_DamageInfo.stCaster.nLocation		= m_pCaster->GetLocation();
	m_DamageInfo.stTarget.nLocation		= pTarget->GetLocation();
	m_DamageInfo.dwSpellID				= m_pEntry->dwID;
	m_DamageInfo.dwProcCaster			= m_dwProcCaster;
	m_DamageInfo.dwProcTarget			= m_dwProcTarget;
	m_DamageInfo.dwProcEx				= m_dwProcEx;

	// 计算技能结果
	CalSpellResult(pTarget, m_DamageInfo);

	m_dwEffectFlag = 0;

	// 未命中或闪避或招架
	if( m_DamageInfo.dwProcEx & (EAEE_Miss | EAEE_Dodge) )
	{
		pScene->SendDamage(m_DamageInfo);
	}
	else
	{
		// 计算效果
		m_dwEffectFlag = 0;
		BOOL bHasEffect = FALSE;

		// 计算技能免疫
		if( !pTarget->HasImmunityAny(EIT_Mechanic, m_pEntry->dwMechanicFlags) )
		{	
			for( INT32 i = 0; i < EFFECT_NUM; ++i )
			{
				ESpellEffectType eff = m_pEntry->eEffect[i];

				if( eff == EEST_Null )	continue;

				bHasEffect = TRUE;

				if( pTarget->HasImmunityAny(EIT_Mechanic, m_pEntry->dwEffectMechanic[i]) )
				{
					continue;
				}

				if (!CheckEffectValid(i, pTarget, 0) || !CheckEffectValid(i, pTarget, 1))
				{
					continue;
				}

				(*this.*SpellEffects[eff])(i, pTarget);
			}
		}

		if( bHasEffect && !m_dwEffectFlag )
		{
			m_DamageInfo.dwProcEx |= EAEE_Immnne;
			pScene->SendDamage(m_DamageInfo);
		}
	}

	if( m_dwEffectFlag && 0 < m_nBaseDamage )
	{
		// 计算伤害
		DealHeal(pTarget, m_nBaseDamage, m_DamageInfo);
		DealDamage(pTarget, m_nBaseDamage, m_DamageInfo);

		// 产生伤害
		pTarget->DoneDamage(m_pCaster, m_DamageInfo);

		// For Trigger
		if( EIFT_Damage == m_DamageInfo.eType && 0 < m_DamageInfo.nDamage)
		{
			m_DamageInfo.dwProcTarget |= EAET_Taken_Any_Damage;

			if (pTarget->IsDead())
			{
				m_DamageInfo.dwProcCaster |= EAET_Kill;
				m_DamageInfo.dwProcTarget |= EAET_Killed;
				m_DamageInfo.dwProcEx |= EAEE_Killed;
				m_bKillEntity = TRUE;
				m_dwFinalProcCaster |= EAET_Kill;
			}
			else
			{
				// 是否触发反击
				if( m_pEntry->bBeatBack && (m_DamageInfo.dwProcEx & EAEE_Block) )
				{
					m_listBeatBack.PushBack(pTarget);
				}
			}
		}

		// 发送伤害
		pScene->SendDamage(m_DamageInfo);
	}
	else
	{

	}

	if( m_DamageInfo.dwProcTarget)
	{
		pTarget->GetCombatController().TriggerBySpellResult(FALSE, m_pCaster, m_DamageInfo);
	}

	if(VALID(m_pCaster) )
	{
		m_pCaster->GetCombatController().TriggerBySpellResult(TRUE, pTarget, m_DamageInfo);
	}	
}

//-------------------------------------------------------------------------------
// 施放反击技能
//-------------------------------------------------------------------------------
VOID Spell::CastBeatBackSpell()
{
	EntityHero* pTarget = NULL;
	TargetList::Iterator it = m_listBeatBack.Begin();
	while (m_listBeatBack.PeekNext(it, pTarget))
	{
		if (!VALID(pTarget) )
			continue;

		pTarget->BeatBack(m_pCaster);
	}
}

BOOL Spell::IsTargetValid( EntityHero* pTarget )
{
	if (!VALID(pTarget))
		return FALSE;

	// 是否包括自己
	if (!m_pEntry->bIncludeSelf && (pTarget == m_pCaster))
		return FALSE;

	// 目标种族检查
	if (!((1 << pTarget->GetEntry()->eRace) & m_pEntry->dwTargetRace))
		return FALSE;

	// 目标状态检查
	DWORD dwTargetState = pTarget->GetState();
	DWORD dwTargetStateCheckFlag = m_pEntry->dwTargetStateCheckFlag;

	// 释放者处于鹰眼状态
	if( m_pCaster->HasState(EHS_AntiHidden) )
	{
		dwTargetStateCheckFlag &= ~(1 << EHS_Stealth);
	}

	dwTargetState &= dwTargetStateCheckFlag;

	if( m_pEntry->dwTargetStateLimit != dwTargetState )
		return FALSE;

	return TRUE;
}

VOID Spell::SendCastGO()
{
	if (!VALID(m_pCaster) || m_listTarget.Empty() || !m_pEntry->bHaveVisual)
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	if( !pScene->IsOnlyRecord() )
	{
		CreateSceneProtoMsg(msg, MS_CastGO,);
		*msg << (UINT32)m_pCaster->GetLocation();
		*msg << (UINT32)m_pEntry->dwID;
		*msg << (INT32)m_listTarget.Size();

		EntityHero* pTarget = NULL;
		TargetList::Iterator it = m_listTarget.Begin();
		while( m_listTarget.PeekNext(it, pTarget) )
		{
			if (!VALID(pTarget))
				continue;

			*msg << (UINT32)pTarget->GetLocation();
		}

		pScene->AddMsgList(msg);
	}
}

VOID Spell::SendCastEnd()
{
	// 发送MS_CastGo和MS_CastEnd判断条件需要相同，因为他们是成对生成的
	if (!VALID(m_pCaster) || m_listTarget.Empty() || !m_pEntry->bHaveVisual)
		return;

	if( !m_pCaster->GetScene()->IsOnlyRecord() )
	{
		CreateSceneProtoMsg(msg, MS_CastEnd,);
		*msg << (UINT32)m_pCaster->GetLocation();
		*msg << (UINT32)m_pEntry->dwID;
		*msg << (INT32)m_pCaster->GetAttController().GetAttValue(EHA_CurRage);
		*msg << (INT32)(static_cast<EntityGroup*>(m_pCaster->GetFather())->GetEnergy());
		m_pCaster->GetScene()->AddMsgList(msg);
	}
}

VOID Spell::CalSpellResult( EntityHero* pTarget, tagCalcDamageInfo &DamageInfo )
{
	// 群体伤害
	if( m_pEntry->bGroupDmg )
	{
		DamageInfo.dwProcEx |= EAEE_GroupDmg;
		m_dwFinalProcEx |= EAEE_GroupDmg;
	}

	// 反击技能直接命中，不计算暴击和格挡
	if( m_eSpellType == ERMT_BeatBack )
	{
		DamageInfo.dwProcEx |= EAEE_Normal_Hit;
		return;
	}

	// 判断技能是否命中
	if( !IsSpellHit(pTarget) )
	{
		DamageInfo.dwProcEx |= EAEE_Miss;
		return;
	}

	if ( IsSpellCrit(pTarget) )
	{
		DamageInfo.dwProcEx |= EAEE_Critical_Hit;
		m_dwFinalProcEx |= EAEE_Critical_Hit;
	}

	// 判断技能是否被格挡
	if( IsSpellBlock(pTarget) )
	{
		DamageInfo.dwProcEx |= EAEE_Block;
		m_dwFinalProcEx |= EAEE_Block;
	}

	// 忽略护甲和抗性
	//if( m_pEntry->bCanNotArmor )
	//{
	//	DamageInfo.dwProcEx |= EAEE_IgnoreArmor;
	//	DamageInfo.dwProcEx |= EAEE_IgnorResistance;
	//}

	// 是否恢复目标怒气和能量
	if( !pTarget->HasState(EHS_UnBeat) && !(pTarget->HasState(EHS_ImmunityGroupDmg) && m_pEntry->bGroupDmg) )
	{
		if( m_eSpellType == ERMT_Melee )
		{
			if( m_pCaster->GetCamp() != pTarget->GetCamp() )
				DamageInfo.dwProcEx |= EAEE_RageResume;
		}

		if( m_eSpellType == ERMT_Melee ||m_eSpellType == ERMT_Rage )
		{
			if( m_pCaster->GetCamp() != pTarget->GetCamp() )
				DamageInfo.dwProcEx |= EAEE_EnergyResume;
		}	 
	}

	DamageInfo.dwProcEx |= EAEE_Normal_Hit;
	m_dwFinalProcEx |= EAEE_Normal_Hit;
}

BOOL Spell::IsSpellHit( EntityHero* pTarget )
{
	if(!VALID(pTarget))
		return FALSE;

	if(m_pEntry->bHitCertainly)
		return TRUE;

	if( !VALID(m_pCaster) )
		return FALSE;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return FALSE;

	// 友方必命中
	if( m_pCaster->GetCamp() == pTarget->GetCamp() )
		return TRUE;

	INT nHitChance = m_pCaster->GetAttController().GetAttValue(EHA_Hit) - pTarget->GetAttController().GetAttValue(EHA_Dodge);
	nHitChance += m_pEntry->nSpellHit;

	if( nHitChance < 5000 )
	{
		nHitChance = nHitChance / 2 + 2500;
	}

	nHitChance = Max(2000, nHitChance);

	return nHitChance >= pScene->GetRandom().Rand(1, 10000);
}

BOOL Spell::IsSpellCrit( EntityHero* pTarget )
{
	if(m_pEntry->bNotCrit)
		return FALSE;

	if(!VALID(pTarget))
		return FALSE;

	if( !VALID(m_pCaster) )
		return FALSE;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return FALSE;

	INT nCritChance = m_pCaster->GetAttController().GetAttValue(EHA_Crit);
	// 敌方才算韧性
	if( m_pCaster->GetCamp() != pTarget->GetCamp() )
		nCritChance -= pTarget->GetAttController().GetAttValue(EHA_Resilience);
	nCritChance += m_pEntry->nSpellCrit;

	if( nCritChance > 5000 )
	{
		nCritChance = nCritChance /2 + 2500;
	}

	nCritChance = Max(0, nCritChance);
	nCritChance = Min(9000, nCritChance);

	return nCritChance >= pScene->GetRandom().Rand(1, 10000);
}

BOOL Spell::IsSpellBlock( EntityHero* pTarget )
{
	if(!VALID(pTarget))
		return FALSE;

	if(m_pEntry->bNotBlock)
		return FALSE;

	if( !VALID(m_pCaster) )
		return FALSE;

	// 友方不格挡
	if( m_pCaster->GetCamp() == pTarget->GetCamp() )
		return FALSE;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return FALSE;

	INT nBlockChance = pTarget->GetAttController().GetAttValue(EHA_Block) - m_pCaster->GetAttController().GetAttValue(EHA_Broken);
	nBlockChance -= m_pEntry->nSpellBroken;

	if( nBlockChance > 5000 )
	{
		nBlockChance = nBlockChance /2 + 2500;
	}

	nBlockChance = Max(0, nBlockChance);
	nBlockChance = Min(9000, nBlockChance);

	return nBlockChance >= pScene->GetRandom().Rand(1, 10000);
}

VOID Spell::CalDamage( INT32 nBaseDamage, tagCalcDamageInfo& DamageInfo, EntityHero* pTarget )
{
	if (!VALID(pTarget))
		return;

	if( m_eSpellType == ERMT_Rune || m_eSpellType == ERMT_Pet || m_pEntry->bCanNotArmor )
	{
		DamageInfo.nDamage = nBaseDamage;
		return;
	}

	nBaseDamage += (m_pCaster->GetAttController().GetAttValue(EHA_DmgInc) - pTarget->GetAttController().GetAttValue(EHA_DmgDec));

	if( m_eSpellType == ERMT_Rage )
	{
		FLOAT fDmgMod = m_fRagePctMod * (FLOAT)(nBaseDamage);
		nBaseDamage += fDmgMod; 
	}

	FLOAT fPctDmgMod = (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PctDmgInc) - pTarget->GetAttController().GetAttValue(EHA_PctDmgDec));

	// PVP百分比伤害加成
	if( m_pCaster->GetScene()->GetStateFlag() & ESSF_PVP )
		fPctDmgMod += (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PVPPctDmgInc) - pTarget->GetAttController().GetAttValue(EHA_PVPPctDmgDec));

	// 伤害类型加成
	if (DamageInfo.eSchool == EIS_Physics)
	{
		fPctDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_DamageDonePctPhysics) + pTarget->GetDmgModAtt(EDM_DamageTakenPctPhysics));
	}
	else if(DamageInfo.eSchool == EIS_Magic)
	{
		fPctDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_DamageDonePctMagic) + pTarget->GetDmgModAtt(EDM_DamageTakenPctMagic));
	}
	
	// 伤害种族加成
	fPctDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_RaceDoneKindom + pTarget->GetEntry()->eRace));
	fPctDmgMod += (FLOAT)(pTarget->GetDmgModAtt(EDM_RaceTakenKindom + m_pCaster->GetEntry()->eRace));

	fPctDmgMod += m_pCaster->GetScene()->GetSceneDmgMod();
	fPctDmgMod += m_pCaster->GetScene()->GetLevelSuppress(m_pCaster, pTarget);

	// 判断百分比下限
	if( fPctDmgMod < -7000.0f )
		fPctDmgMod = -7000.0f;

	nBaseDamage += ((fPctDmgMod / 10000.0f) * (FLOAT)nBaseDamage);

	if (nBaseDamage < 1)
	{
		DamageInfo.nDamage = 1;
	}

	if (DamageInfo.dwProcEx & EAEE_Critical_Hit)
	{
		INT nCrit = m_pCaster->GetAttController().GetAttValue(EHA_CritInc) - pTarget->GetAttController().GetAttValue(EHA_CritDec);
		nBaseDamage *= (Max(10000, 17000 + nCrit) /10000.0f);
	}

	if (DamageInfo.dwProcEx & EAEE_Block)
	{
		nBaseDamage *= 0.5f;
	}

	INT32 nMinDmg = (INT32)((FLOAT)m_pCaster->GetAttController().GetAttValue(EHA_AttackPower) * 0.05f);
	DamageInfo.nDamage = Max(nMinDmg, nBaseDamage);
}

VOID Spell::CalHeal( INT32 nBaseHeal, tagCalcDamageInfo& DamageInfo, EntityHero* pTarget )
{
	if (!VALID(pTarget))
		return;

	// 重伤状态无法加血
	if( pTarget->HasState(EHS_Injury) )
	{
		DamageInfo.nDamage = 0;
		return; 
	}		

	// 中毒状态加血效果减75%
	FLOAT fHealPct		= pTarget->HasState(EHS_Poison) ? 0.25f : 1.0f;

	if( m_eSpellType == ERMT_Rune || m_eSpellType == ERMT_Pet )
	{
		DamageInfo.nDamage = nBaseHeal * fHealPct;
		return;
	}

	if( m_pEntry->bCanNotArmor )
	{
		DamageInfo.nDamage = nBaseHeal;
		return;
	}

	if( m_eSpellType == ERMT_Rage )
	{
		FLOAT fDmgMod = m_fRagePctMod * (FLOAT)(nBaseHeal);
		nBaseHeal += fDmgMod; 
	}

	FLOAT fDmgMod = (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PctDmgInc));
	if( m_pCaster->GetScene()->GetStateFlag() & ESSF_PVP )
		fDmgMod += (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PVPPctDmgInc));

	fDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_DamageDonePctHeal) + pTarget->GetDmgModAtt(EDM_DamageTakenPcttHeal));
//	fDmgMod += (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_HealPctIncDone) + pTarget->GetAttController().GetAttValue(EHA_HealPctIncTaken));
	fDmgMod = fDmgMod / 10000.0f;
	fDmgMod = fDmgMod * nBaseHeal;
	nBaseHeal += fDmgMod;

	if (DamageInfo.dwProcEx & EAEE_Critical_Hit)
	{
		INT nCrit = m_pCaster->GetAttController().GetAttValue(EHA_CritInc);
		nBaseHeal *= (Max(10000, 17000 + nCrit) /10000.0f);
	}

	DamageInfo.nDamage = Max(0, nBaseHeal) * fHealPct;
}

VOID Spell::DealDamage( EntityHero* pTarget, INT32 nBaseDamage, tagCalcDamageInfo& DamageInfo )
{
	if (!VALID(pTarget))
		return;

	if( EIS_Null == DamageInfo.eSchool )	return;

	if( EIFT_Damage != DamageInfo.eType )	return;

	CalDamage(nBaseDamage, DamageInfo, pTarget);

	// 触发双方的伤害接口
	if (VALID(m_pCaster))
	{
		m_pCaster->OnDamage(pTarget, DamageInfo);
	}
	pTarget->OnBeDamaged(m_pCaster, DamageInfo);

	// todo:删除受伤打断的aura
}

VOID Spell::DealHeal( EntityHero* pTarget, INT32 nBaseHeal, tagCalcDamageInfo& DamageInfo )
{
	if (!VALID(pTarget))
		return;

	if( EIS_Null == DamageInfo.eSchool )	return;

	if( EIFT_Heal != DamageInfo.eType )		return;

	CalHeal(nBaseHeal, DamageInfo, pTarget);

	// 触发双方的伤害接口
	if (VALID(m_pCaster))
	{
		m_pCaster->OnDamage(pTarget, DamageInfo);
	}
	pTarget->OnBeDamaged(m_pCaster, DamageInfo);

	// 计算有效治疗
	DamageInfo.nDamage = Min(DamageInfo.nDamage, pTarget->GetAttController().GetAttValue(EHA_MaxHP) - pTarget->GetAttController().GetAttValue(EHA_CurHP));
}





//--------------------------------------------------------------------------------------------------
// 效果是否可作用于目标
//--------------------------------------------------------------------------------------------------
BOOL Spell::CheckEffectValid(INT32 nEffectIndex, EntityHero* pTarget, INT32 nIndex)
{
	if( EEST_Null == m_pEntry->eEffect[nEffectIndex])	return FALSE;

	switch (m_pEntry->eEffectTargetValid[nIndex][nEffectIndex])
	{
	case EETV_Null:
		{
			return TRUE;
		}
		break;

	case EETV_Self:
		{
			if( pTarget == m_pCaster )	return TRUE;
		}
		break;

	case EETV_UnSelf:
		{
			if( pTarget != m_pCaster )	return TRUE;
		}
		break;

	case EETV_Caster_State:
		{
			if( m_pCaster->HasStateAny(m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex]) )	
				return TRUE;

			return FALSE;
		}
		break;

	case EETV_Target_State:
		{
			if( pTarget->HasStateAny(m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex]) )	
				return TRUE;

			return FALSE;
		}
		break;

	case EETV_Caster_HP_Low:
		{
			FLOAT fHPPct = (FLOAT)(m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex]);
			if((fHPPct / 10000.0f) * (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_MaxHP)) > (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_CurHP)))
				return TRUE;
		}
		break;

	case EETV_Target_HP_Low:
		{
			FLOAT fHPPct = (FLOAT)(m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex]);
			if((fHPPct / 10000.0f) * (FLOAT)(pTarget->GetAttController().GetAttValue(EHA_MaxHP)) > (FLOAT)(pTarget->GetAttController().GetAttValue(EHA_CurHP)))
				return TRUE;
		}
		break;

	case EETV_Target_HP_High:
		{
			FLOAT fHPPct = (FLOAT)(m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex]);
			if((fHPPct / 10000.0f) * (FLOAT)(pTarget->GetAttController().GetAttValue(EHA_MaxHP)) < (FLOAT)(pTarget->GetAttController().GetAttValue(EHA_CurHP)))
				return TRUE;
		}
		break;

	case EETV_Pct:
		{
			// 种族概率加成
			INT32 nBasePct = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
			if( m_pEntry->eEffectValidRace[nEffectIndex] == pTarget->GetEntry()->eRace )
				nBasePct += m_pEntry->dwEffectValidRaceMod[nEffectIndex];

			// 等级概率衰减
			if( VALID(m_pEntry->nDecayLevel) && pTarget->GetLevel() > m_pEntry->nDecayLevel )
			{
				INT32 nLevelDiffer = pTarget->GetLevel() - m_pEntry->nDecayLevel;
				nBasePct -= nLevelDiffer * m_pEntry->nDecayRate;
			}

			if( nBasePct > m_pCaster->GetScene()->GetRandom().Rand(1, 10000) )
				return TRUE;
			else
			{
				m_DamageInfo.dwProcEx |= EAEE_Invalid;
				return FALSE;
			}
		}
		break;

	case EETV_Target_AuraNot:
		{
			DWORD dwAuraID = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
			if( VALID(pTarget->GetCombatController().GetAuraByIDCaster(dwAuraID)) )
				return FALSE;

			return TRUE;
		}
		break;

	case EETV_Target_Aura:
		{
			DWORD dwAuraID = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
			if( VALID(pTarget->GetCombatController().GetAuraByIDCaster(dwAuraID)) )
				return TRUE;
		}
		break;

	case EETV_Target_Race:
		{
			// 目标种族限制
			if (VALID(pTarget) && VALID(m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex]))
			{
				if(((1 << pTarget->GetEntry()->eRace) & m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex]))
					return TRUE;
			}
		}
		break;

	case EEIV_Caster_AuraState:
		{
			INT32 nAuraState = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
			if( m_pCaster->GetCombatController().HasAuraState(nAuraState) )
				return TRUE;

			return FALSE;
		}
		break;

	case EEIV_Target_AuraState:
		{
			INT32 nAuraState = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];
			if( pTarget->GetCombatController().HasAuraState(nAuraState) )
				return TRUE;

			return FALSE;
		}
		break;

	case EETV_Target_GT_Level:
		{
			if (VALID(pTarget) && pTarget->GetLevel() > m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex] )
			{
				return TRUE;
			}
		}
		break;

	case EETV_Target_LT_Level:
		{
			if (VALID(pTarget) && pTarget->GetLevel() <= m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex] )
			{
				return TRUE;
			}
		}
		break;

	case EEIV_Caster_AuraPN:
		{
			INT32 nEffectPriority = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];

			INT32	nPosNum		=	0;
			INT32	nNegNum		=	0;
			m_pCaster->GetCombatController().GetPositiveAndNegativeNum(nPosNum, nNegNum);
			if( (nEffectPriority > 0 && nPosNum > 0) || (nEffectPriority < 0 && nNegNum > 0) )
				return TRUE;

			return FALSE;
		}
		break;

	case EEIV_Target_AuraPN:
		{
			INT32 nEffectPriority = m_pEntry->dwEffectValidMiscValue[nIndex][nEffectIndex];

			INT32	nPosNum		=	0;
			INT32	nNegNum		=	0;
			pTarget->GetCombatController().GetPositiveAndNegativeNum(nPosNum, nNegNum);
			if( (nEffectPriority > 0 && nPosNum > 0) || (nEffectPriority < 0 && nNegNum > 0) )
				return TRUE;

			return FALSE;
		}
		break;

	default:
		return FALSE;
	}

	return FALSE;
}
