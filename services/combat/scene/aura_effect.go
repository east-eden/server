package scene

import "bitbucket.org/funplus/server/define"

// 技能效果处理函数
type AuraEffectsHandler func(*Aura, define.EAuraEffectStep, int32, interface{}, interface{}) define.EAuraAddResult

var auraEffectsHandlers []AuraEffectsHandler = []AuraEffectsHandler{
	AuraEffectNull,               // 0
	AuraEffectPeriodDamage,       // 1周期伤害
	AuraEffectModAtt,             // 2属性改变
	AuraEffectSpell,              // 3施放技能
	AuraEffectState,              // 4状态变更
	AuraEffectImmunity,           // 5免疫变更
	AuraEffectDmgMod,             // 6伤害转换(总量一定)
	AuraEffectNewDmg,             // 7生成伤害(原伤害不变)
	AuraEffectDmgFix,             // 8限量伤害(改变原伤害)
	AuraEffectChgMelee,           // 9替换普通攻击
	AuraEffectShield,             // 10血盾
	AuraEffectDmgAttMod,          // 11改变伤害属性
	AuraEffectAbsorbAllDmg,       // 12全伤害吸收
	AuraEffectDmgAccumulate,      // 13累计伤害
	AuraEffectMeleeSpell,         // 14释放普通攻击
	AuraEffectLimitAttack,        // 15限制攻击力
	AuraEffectPeriodHeal,         // 16周期治疗
	AuraEffectModAttByAlive,      // 15根据当前友方存活人数,计算属性改变
	AuraEffectModAttByEnemyAlive, // 16根据当前敌方存活人数,计算属性改变
}

func AuraEffectNull(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 周期伤害
//-------------------------------------------------------------------------------
func AuraEffectPeriodDamage(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//ESchool eSchool		=	m_pAuraEntry->eSchool;
	//INT32	nDmgPctMin  =	m_pAuraEntry->nMiscValue1[nIndex];
	//INT32	nDmgPctMax  =	m_pAuraEntry->nMiscValue2[nIndex];
	//FLOAT	fDmgPct		=	(FLOAT)m_pCaster->GetScene()->GetRandom().Rand(nDmgPctMin, nDmgPctMax);

	//ASSERT( eSchool == EIS_Physics || eSchool == EIS_Magic );

	//switch (eStep)
	//{
	//case EAES_Apply:
	//{
	//if (!VALID(m_pCaster))	return EAAR_Immunity;

	//FLOAT fDamage = (FLOAT)m_pCaster->GetAttController().GetAttValue(EHA_AttackPower);
	//if( fDamage > m_pOwner->GetEntry()->dwAttackReduceThreshold )
	//{
	//fDamage = (fDamage - m_pOwner->GetEntry()->dwAttackReduceThreshold) / 2 + m_pOwner->GetEntry()->dwAttackReduceThreshold;
	//}
	//fDamage -= m_pOwner->GetAttController().GetAttValue(EHA_DefenceMelee+eSchool);
	//fDamage = Max(1.0f, fDamage);
	//fDamage *= fDmgPct / 10000.0f;
	//fDamage += (FLOAT)m_nCurPoint[nIndex];
	//fDamage += (m_nAmount * (m_pAuraEntry->nAmountEffect[nIndex] / 10000.0f));
	//fDamage *= m_fMultiple[nIndex] / 10000.0f;

	//fDamage = Max(1.0f, fDamage);

	//m_nEffMisc[nIndex] = (INT32)fDamage;
	//}
	//break;

	//case EAES_Effect:
	//{
	//Scene* pScene = m_pOwner->GetScene();
	//if (!VALID(pScene))
	//{
	//return EAAR_Immunity;
	//}

	//tagCalcDamageInfo DamageInfo;
	//DamageInfo.Reset();

	//// 检查免疫
	//if (m_pOwner->HasImmunity(EIT_Damage, eSchool))
	//{
	//DamageInfo.dwProcEx	|= EAEE_Immnne;
	//pScene->SendDamage(DamageInfo);
	//return EAAR_Immunity;
	//}

	//DamageInfo.nDamage				= m_nEffMisc[nIndex];
	//DamageInfo.stCaster.nLocation	= m_pCaster->GetLocation();
	//DamageInfo.stTarget.nLocation	= m_pOwner->GetLocation();
	//DamageInfo.dwSpellID			= m_pAuraEntry->dwID;
	//DamageInfo.eType				= EIFT_Damage;
	//DamageInfo.eSchool				= eSchool;
	//DamageInfo.dwProcEx				|= EAEE_IgnoreArmor;
	//DamageInfo.dwProcEx				|= EAEE_Normal_Hit;
	//DamageInfo.dwProcCaster			|= EAET_On_Do_Periodic;
	//DamageInfo.dwProcTarget			|= EAET_On_Take_Periodic;

	//CalDamage(DamageInfo.nDamage, DamageInfo, m_pOwner);

	//m_pCaster->OnDamage(m_pOwner, DamageInfo);
	//m_pOwner->OnBeDamaged(m_pCaster, DamageInfo);

	//// todo:删除受伤打断的aura

	//// 发送伤害
	//pScene->SendDamage(DamageInfo);

	//// 产生伤害
	//m_pOwner->DoneDamage(m_pCaster, DamageInfo);

	//// For Trigger
	//if( 0 < DamageInfo.nDamage)
	//{
	//DamageInfo.dwProcTarget |= EAET_Taken_Any_Damage;

	//if (m_pOwner->IsDead())
	//{
	//DamageInfo.dwProcCaster |= EAET_Kill;
	//DamageInfo.dwProcTarget |= EAET_Killed;
	//}
	//}

	//// Trigger
	//if (VALID(m_pCaster))
	//{
	//m_pCaster->GetCombatController().TriggerBySpellResult(TRUE, m_pOwner, DamageInfo);
	//}
	//m_pOwner->GetCombatController().TriggerBySpellResult(FALSE, m_pCaster, DamageInfo);
	//}
	//break;

	//case EAES_Check:
	//// 检查免疫
	//if (m_pOwner->HasImmunity(EIT_Damage, eSchool))
	//return EAAR_Immunity;
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 属性改变
//-------------------------------------------------------------------------------
func AuraEffectModAtt(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//INT nAttType1			= m_pAuraEntry->nMiscType1[nIndex];
	//INT nAttType2			= m_pAuraEntry->nMiscType2[nIndex];
	//SpellValueType eModType	= (SpellValueType)m_pAuraEntry->nMiscValue1[nIndex];

	//// 映射属性
	//INT32 nAttType			= sEntityAttType.GetType(nAttType1, nAttType2);
	//if (!VALID(nAttType))
	//return EAAR_Immunity;

	//m_pOwner->GetAttController().LockRecal();

	//switch (eStep)
	//{
	//case EAES_Apply:
	//switch (eModType)
	//{
	//case SMT_FLAT:
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, -m_nEffMisc[nIndex]);
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, m_nCurPoint[nIndex]);
	//break;

	//case SMT_PCT:
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, -m_nEffMisc[nIndex]);
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, m_nCurPoint[nIndex]);
	//break;
	//}
	//m_nEffMisc[nIndex] = m_nCurPoint[nIndex];
	//break;

	//case EAES_Remove:
	//switch (eModType)
	//{
	//case SMT_FLAT:
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, -m_nEffMisc[nIndex]);
	//break;

	//case SMT_PCT:
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, -m_nEffMisc[nIndex]);
	//break;
	//}
	//break;
	//}

	//m_pOwner->GetAttController().RecalAtt();

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 施放技能
//-------------------------------------------------------------------------------
func AuraEffectSpell(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//EntityHero* pCaster	= NULL;
	//EntityHero* pTarget	= NULL;

	//switch (m_pAuraEntry->nMiscType1[nIndex])
	//{
	//case ECU_Caster:
	//pCaster	= m_pCaster;
	//break;

	//case ECU_Owner:
	//pCaster = m_pOwner;
	//break;

	//case ECU_Target:
	//if (!VALID((pCaster = (EntityHero*)pParam2)))
	//{
	//return EAAR_Immunity;
	//}
	//break;
	//}

	//switch (m_pAuraEntry->nMiscType2[nIndex])
	//{
	//case ECU_Caster:
	//pTarget	= m_pCaster;
	//break;

	//case ECU_Owner:
	//pTarget	= m_pOwner;
	//break;

	//case ECU_Target:
	//pTarget	= (EntityHero*)pParam2;
	//break;
	//}

	//switch (eStep)
	//{
	//case EAES_Apply:
	//if (VALID(m_pAuraEntry->nMiscValue1[nIndex]))
	//{
	//if (E_Success != m_pOwner->GetCombatController().CastSpell(m_pAuraEntry->nMiscValue1[nIndex], pCaster, pTarget, TRUE, 0, ERMT_AuraTrigger))
	//return EAAR_Immunity;
	//}
	//break;

	//case EAES_Effect:
	//if (VALID(m_pAuraEntry->nMiscValue2[nIndex]))
	//{
	//if( VALID(pParam1)  )
	//{
	//MTransPtr(pDmg, pParam1, tagCalcDamageInfo);
	//if( (pDmg->dwProcEx & EAEE_Internal_Cant_Trigger) )
	//{
	//return EAAR_Immunity;
	//}

	//if( (pDmg->dwProcEx & EAEE_Internal_Triggered) )
	//{
	//if (E_Success != m_pOwner->GetCombatController().CastSpell(m_pAuraEntry->nMiscValue2[nIndex], pCaster, pTarget, TRUE, 0, ERMT_AuraTriggerTwice))
	//return EAAR_Immunity;
	//else
	//return EAAR_Success;
	//}
	//}

	//if (E_Success != m_pOwner->GetCombatController().CastSpell(m_pAuraEntry->nMiscValue2[nIndex], pCaster, pTarget, TRUE, 0, ERMT_AuraTrigger))
	//return EAAR_Immunity;
	//}
	//break;

	//case EAES_Remove:
	//if (VALID(m_pAuraEntry->nMiscValue3[nIndex]))
	//{
	//if (E_Success != m_pOwner->GetCombatController().CastSpell(m_pAuraEntry->nMiscValue3[nIndex], pCaster, pTarget, TRUE, 0, ERMT_AuraTrigger))
	//return EAAR_Immunity;
	//}
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 改变状态
//-------------------------------------------------------------------------------
func AuraEffectState(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//DWORD dwServentState = m_pAuraEntry->nMiscValue1[nIndex];

	//switch (eStep)
	//{
	//case EAES_Apply:
	//// 设置状态
	//m_pOwner->AddState(dwServentState);

	//if( dwServentState & EHSF_Taunt )
	//{
	//m_pOwner->SetTauntTarget(m_pCaster->GetPos());
	//}

	//break;

	//case EAES_Remove:
	//// 去除状态
	//m_pOwner->DecState(dwServentState);
	//break;

	//default:
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 免疫变更
//-------------------------------------------------------------------------------
func AuraEffectImmunity(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//EImmunityType eType		= (EImmunityType)m_pAuraEntry->nMiscType1[nIndex];
	//DWORD dwUnitImmunity	= m_pAuraEntry->nMiscValue1[nIndex];

	//switch (eStep)
	//{
	//case EAES_Apply:
	//m_pOwner->AddImmunity(eType, dwUnitImmunity);
	//break;

	//case EAES_Remove:
	//m_pOwner->DecImmunity(eType, dwUnitImmunity);
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 伤害转换(总量一定)
//-------------------------------------------------------------------------------
func AuraEffectDmgMod(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//DWORD		dwSpellID	= m_pAuraEntry->dwMiscID[nIndex];
	//ESchoolMask eSrcSchool	= (ESchoolMask)m_pAuraEntry->nMiscType1[nIndex];
	//SpellValueType eType	= (SpellValueType)m_pAuraEntry->nMiscType2[nIndex];
	//DmgInfoType eDmgType	= (DmgInfoType)m_pAuraEntry->nMiscValue1[nIndex];

	//switch (eStep)
	//{
	//case EAES_Effect:
	//{
	//if( !VALID(pParam1) || !VALID(pParam2) )
	//return EAAR_Immunity;

	//MTransPtr(pDmg, pParam1, tagCalcDamageInfo);
	//if( (pDmg->dwProcEx & EAEE_Internal_Triggered) || (pDmg->dwProcEx & EAEE_Internal_Cant_Trigger) )
	//return EAAR_Immunity;

	//if ((eDmgType == pDmg->eType) && (eSrcSchool & (1 << pDmg->eSchool)))
	//{
	//INT32 nDamage = 0;
	//switch (eType)
	//{
	//case SMT_FLAT:
	//nDamage = m_nCurPoint[nIndex];
	//break;

	//case SMT_PCT:
	//nDamage = INT32((FLOAT)pDmg->nDamage * ((FLOAT)m_nCurPoint[nIndex] / 10000.0f));
	//break;
	//}

	//// 原伤害
	//pDmg->nDamage = (nDamage > pDmg->nDamage) ? 0 : pDmg->nDamage - nDamage;

	//// 生成新伤害
	//MTransPtr(pTarget, pParam2, EntityHero);
	//m_pOwner->GetCombatController().CastSpell(dwSpellID, m_pOwner, pTarget, TRUE, nDamage, ERMT_AuraTrigger);
	//}
	//}
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 生成伤害(原伤害不变)
//-------------------------------------------------------------------------------
func AuraEffectNewDmg(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//DWORD		dwSpellID	= m_pAuraEntry->dwMiscID[nIndex];
	//ESchoolMask eSrcSchool	= (ESchoolMask)m_pAuraEntry->nMiscType1[nIndex];
	//SpellValueType eType	= (SpellValueType)m_pAuraEntry->nMiscType2[nIndex];
	//DmgInfoType eDmgType	= (DmgInfoType)m_pAuraEntry->nMiscValue1[nIndex];

	//switch (eStep)
	//{
	//case EAES_Effect:
	//{
	//if( !VALID(pParam1) || !VALID(pParam2) )
	//return EAAR_Immunity;

	//MTransPtr(pDmg, pParam1, tagCalcDamageInfo);
	//if( (pDmg->dwProcEx & EAEE_Internal_Triggered) || (pDmg->dwProcEx & EAEE_Internal_Cant_Trigger) )
	//return EAAR_Immunity;

	//if ((eDmgType == pDmg->eType) && (eSrcSchool & (1 << pDmg->eSchool)))
	//{
	//INT32 nDamage = 0;
	//switch (eType)
	//{
	//case SMT_FLAT:
	//nDamage = m_nCurPoint[nIndex];
	//break;

	//case SMT_PCT:
	//nDamage = INT32((FLOAT)pDmg->nDamage * ((FLOAT)m_nCurPoint[nIndex] / 10000.0f));
	//break;
	//}

	//// 生成新伤害
	//MTransPtr(pTarget, pParam2, EntityHero);
	//m_pOwner->GetCombatController().CastSpell(dwSpellID, m_pOwner, pTarget, TRUE, nDamage, ERMT_AuraTrigger);
	//}
	//}
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 限量伤害(改变原伤害)
//-------------------------------------------------------------------------------
func AuraEffectDmgFix(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//DWORD		dwSpellID	= m_pAuraEntry->dwMiscID[nIndex];
	//ESchoolMask eSrcSchool	= (ESchoolMask)m_pAuraEntry->nMiscType1[nIndex];
	//SpellValueType eType	= (SpellValueType)m_pAuraEntry->nMiscType2[nIndex];
	//DmgInfoType eDmgType	= (DmgInfoType)m_pAuraEntry->nMiscValue1[nIndex];

	//switch (eStep)
	//{
	//case EAES_Effect:
	//{
	//if( !VALID(pParam1) || !VALID(pParam2) )
	//return EAAR_Immunity;

	//MTransPtr(pDmg, pParam1, tagCalcDamageInfo);
	//if( (pDmg->dwProcEx & EAEE_Internal_Triggered) || (pDmg->dwProcEx & EAEE_Internal_Cant_Trigger) )
	//return EAAR_Immunity;

	//if ((eDmgType == pDmg->eType) && (eSrcSchool & (1 << pDmg->eSchool)))
	//{
	//INT32 nDamage = 0;
	//switch (eType)
	//{
	//case SMT_FLAT:
	//nDamage = m_nCurPoint[nIndex];
	//break;

	//case SMT_PCT:
	//nDamage = INT32((FLOAT)pDmg->nDamage * ((FLOAT)m_nCurPoint[nIndex] / 10000.0f));
	//break;
	//}

	//// 原伤害
	//pDmg->nDamage = (nDamage > pDmg->nDamage) ? pDmg->nDamage : nDamage;

	//// 生成新伤害
	//if (VALID(dwSpellID))
	//{
	//MTransPtr(pTarget, pParam2, EntityHero);
	//m_pOwner->GetCombatController().CastSpell(dwSpellID, m_pOwner, pTarget, TRUE, pDmg->nDamage, ERMT_AuraTrigger);
	//}
	//}
	//}
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 替换普通攻击
//-------------------------------------------------------------------------------
func AuraEffectChgMelee(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//DWORD dwSpellID	= m_pAuraEntry->dwMiscID[nIndex];

	//if (!VALID(m_pOwner))
	//return EAAR_Immunity;

	//switch (eStep)
	//{
	//case EAES_Apply:
	//static_cast<EntityHero*>(m_pOwner)->SetMeleeSpell(dwSpellID);
	//break;

	//case EAES_Remove:
	//static_cast<EntityHero*>(m_pOwner)->SetMeleeSpell(INVALID);
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 血盾
//-------------------------------------------------------------------------------
func AuraEffectShield(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//m_pOwner->GetAttController().LockRecal();

	//switch (eStep)
	//{
	//case EAES_Apply:
	//m_pOwner->GetAttController().ModAttFlatModValue(EHA_MaxHP, -m_nEffMisc[nIndex]);
	//m_pOwner->GetAttController().ModAttFlatModValue(EHA_MaxHP, m_nCurPoint[nIndex]);
	//m_pOwner->GetAttController().ModAttValue(EHA_CurHP, m_nCurPoint[nIndex]);
	//m_nEffMisc[nIndex] = m_nCurPoint[nIndex];
	//break;

	//case EAES_Remove:
	//m_pOwner->GetAttController().ModAttFlatModValue(EHA_MaxHP, -m_nEffMisc[nIndex]);
	//break;
	//}

	//m_pOwner->GetAttController().RecalAtt();

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 改变伤害属性
//-------------------------------------------------------------------------------
func AuraEffectDmgAttMod(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//INT32 nDmgAttMask	= m_pAuraEntry->nMiscType1[nIndex];

	//switch (eStep)
	//{
	//case EAES_Apply:
	//{
	//for( INT32 i = 0; i < EDM_End; ++i )
	//{
	//if( nDmgAttMask & (1 << i) )
	//{
	//m_pOwner->ModDmgModAtt(EDamageMod(i), -m_nEffMisc[nIndex]);
	//m_pOwner->ModDmgModAtt(EDamageMod(i), m_nCurPoint[nIndex]);
	//}
	//}
	//m_nEffMisc[nIndex] = m_nCurPoint[nIndex];
	//}

	//break;

	//case EAES_Remove:
	//{
	//for( INT32 i = 0; i < EDM_End; ++i )
	//{
	//if( nDmgAttMask & (1 << i) )
	//{
	//m_pOwner->ModDmgModAtt(EDamageMod(i), -m_nCurPoint[nIndex]);
	//}
	//}
	//}

	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 全类型伤害吸收
//-------------------------------------------------------------------------------
func AuraEffectAbsorbAllDmg(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//FLOAT fAbsorbFlat		= (FLOAT)m_pAuraEntry->nMiscValue1[nIndex];
	//SpellValueType eType	= (SpellValueType)m_pAuraEntry->nMiscValue2[nIndex];
	//BOOL bCasterHP			= (BOOL)(m_pAuraEntry->nMiscValue3[nIndex]);

	//switch (eStep)
	//{
	//case EAES_Apply:
	//{
	//switch (eType)
	//{
	//case SMT_FLAT:
	//// 给效果变量赋值
	//m_nEffMisc[nIndex] = (INT32)fAbsorbFlat;
	//break;

	//case SMT_PCT:
	//FLOAT fAbsorb = 0.0f;
	//if( bCasterHP )
	//{
	//fAbsorb = (FLOAT)m_pCaster->GetAttController().GetAttValue(EHA_MaxHP);
	//}
	//else
	//{
	//fAbsorb = (FLOAT)m_pOwner->GetAttController().GetAttValue(EHA_MaxHP);
	//}

	//fAbsorb *= (FLOAT)m_nCurPoint[nIndex] / 10000.0f;
	//fAbsorb += (FLOAT)fAbsorbFlat;
	//fAbsorb += (m_nAmount * (m_pAuraEntry->nAmountEffect[nIndex] / 10000.0f));
	//fAbsorb *= (m_fMultiple[nIndex] / 10000.0f);

	//// 给效果变量赋值
	//m_nEffMisc[nIndex] = (INT32)fAbsorb;
	//break;
	//}
	//}
	//break;

	//case EAES_Effect:
	//// 受伤事件
	//{
	//// 优先触发无敌
	//if( m_pOwner->HasState(EHS_UnBeat) )
	//break;

	//MTransPtr(pDmg, pParam1, tagCalcDamageInfo);
	//if ((pDmg->nDamage > 0) && (pDmg->eType == EIFT_Damage))
	//{
	//// 计算可吸收伤害
	//INT32 nMaxAbsorb	= min(m_nEffMisc[nIndex], pDmg->nDamage);
	//if (nMaxAbsorb > 0)
	//{
	//pDmg->dwProcEx |= EAEE_Absorb;
	//}

	//// 伤害改变
	//if (nMaxAbsorb >= pDmg->nDamage)
	//{
	//m_nEffMisc[nIndex] -= pDmg->nDamage;
	//pDmg->nDamage = 0;
	//}
	//else
	//{
	//pDmg->nDamage -= nMaxAbsorb;
	//m_nEffMisc[nIndex] = 0;
	//}

	//if (m_nEffMisc[nIndex] == 0)
	//{
	//// 删除Aura
	//m_pOwner->GetCombatController().RemoveAura(this, EARM_Default);
	//}
	//}
	//}
	//break;

	//case EAES_Remove:
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 累计伤害
//-------------------------------------------------------------------------------
func AuraEffectDmgAccumulate(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//DWORD dwSpellID				= m_pAuraEntry->dwMiscID[nIndex];
	//FLOAT fAccumulatePct		= (FLOAT)m_pAuraEntry->nMiscValue1[nIndex];

	//switch (eStep)
	//{
	//case EAES_Apply:
	//{
	//// 给效果变量赋值
	//m_nEffMisc[nIndex] = 0;
	//}
	//break;

	//case EAES_Effect:
	//// 受伤事件
	//{
	//MTransPtr(pDmg, pParam1, tagCalcDamageInfo);
	//if ((pDmg->nDamage > 0) && (pDmg->eType == EIFT_Damage))
	//{
	//m_nEffMisc[nIndex] += pDmg->nDamage;
	//}
	//}
	//break;

	//case EAES_Remove:

	//if (EARM_Dispel != (EAuraRemoveMode)*(INT*)pParam1)
	//{
	//INT32 nDamage = (INT32)((FLOAT)m_nEffMisc[nIndex] * (fAccumulatePct / 10000.0f));
	//// 生成新伤害
	//m_pOwner->GetCombatController().CastSpell(dwSpellID, m_pOwner, m_pOwner, TRUE, nDamage, ERMT_AuraTrigger);
	//}

	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 释放普通攻击
//-------------------------------------------------------------------------------
func AuraEffectMeleeSpell(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//EntityHero* pTarget	= (EntityHero*)pParam2;

	//switch (eStep)
	//{
	//case EAES_Apply:
	//break;

	//case EAES_Effect:
	//if (VALID(m_pAuraEntry->nMiscValue2[nIndex]))
	//{
	//m_pOwner->GetCombatController().CastSpell(m_pAuraEntry->nMiscValue2[nIndex], m_pOwner, pTarget, FALSE, 0, ERMT_Melee);
	//}
	//break;

	//case EAES_Remove:
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 限制攻击力
//-------------------------------------------------------------------------------
func AuraEffectLimitAttack(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//INT nAttackLimit		= m_pAuraEntry->nMiscType1[nIndex];				// 攻击力上限

	//m_pOwner->GetAttController().LockRecal();

	//switch (eStep)
	//{
	//case EAES_Apply:
	//{
	//INT32 nCurAttack = m_pOwner->GetAttController().GetAttValue(EHA_AttackPower);
	//if( nCurAttack > nAttackLimit )
	//{
	//m_nEffMisc[nIndex] = nCurAttack - nAttackLimit;
	//m_pOwner->GetAttController().ModAttFlatModValue(EHA_AttackPower, -m_nEffMisc[nIndex]);
	//}
	//}
	//break;

	//case EAES_Remove:
	//{
	//m_pOwner->GetAttController().ModAttFlatModValue(EHA_AttackPower, m_nEffMisc[nIndex]);
	//}
	//break;
	//}

	//m_pOwner->GetAttController().RecalAtt();

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 周期治疗
//-------------------------------------------------------------------------------
func AuraEffectPeriodHeal(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//INT32	nDmgPctMin  =	m_pAuraEntry->nMiscValue1[nIndex];
	//INT32	nDmgPctMax  =	m_pAuraEntry->nMiscValue2[nIndex];
	//FLOAT	fHealPct	=	(FLOAT)m_pCaster->GetScene()->GetRandom().Rand(nDmgPctMin, nDmgPctMax);

	//switch (eStep)
	//{
	//case EAES_Apply:
	//{
	//if (!VALID(m_pCaster))	return EAAR_Immunity;

	//FLOAT fHeal = (FLOAT)m_pCaster->GetAttController().GetAttValue(EHA_AttackPower);
	//fHeal *= fHealPct / 10000.0f;
	//fHeal += (FLOAT)m_nCurPoint[nIndex];
	//fHeal += (m_nAmount * (m_pAuraEntry->nAmountEffect[nIndex] / 10000.0f));
	//fHeal *= m_fMultiple[nIndex] / 10000.0f;

	//m_nEffMisc[nIndex] = (INT32)fHeal;
	//}
	//break;

	//case EAES_Effect:
	//{
	//Scene* pScene = m_pOwner->GetScene();
	//if (!VALID(pScene))
	//{
	//return EAAR_Immunity;
	//}

	//tagCalcDamageInfo DamageInfo;
	//DamageInfo.Reset();

	//DamageInfo.nDamage				= m_nEffMisc[nIndex];
	//DamageInfo.stCaster.nLocation	= m_pCaster->GetLocation();
	//DamageInfo.stTarget.nLocation	= m_pOwner->GetLocation();
	//DamageInfo.dwSpellID			= m_pAuraEntry->dwID;
	//DamageInfo.eType				= EIFT_Heal;
	//DamageInfo.eSchool				= m_pAuraEntry->eSchool;
	//DamageInfo.dwProcCaster			|= EAET_On_Do_Periodic;
	//DamageInfo.dwProcTarget			|= EAET_On_Take_Periodic;

	//CalHeal(DamageInfo.nDamage, DamageInfo, m_pOwner);

	//m_pCaster->OnDamage(m_pOwner, DamageInfo);
	//m_pOwner->OnBeDamaged(m_pCaster, DamageInfo);

	//// todo:删除受伤打断的aura

	//// 发送伤害
	//pScene->SendDamage(DamageInfo);

	//// 产生伤害
	//m_pOwner->DoneDamage(m_pCaster, DamageInfo);

	//// Trigger
	//if (VALID(m_pCaster))
	//{
	//m_pCaster->GetCombatController().TriggerBySpellResult(TRUE, m_pOwner, DamageInfo);
	//}
	//m_pOwner->GetCombatController().TriggerBySpellResult(FALSE, m_pCaster, DamageInfo);
	//}
	//break;
	//}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 根据当前友方存活人数,计算属性改
//-------------------------------------------------------------------------------
func AuraEffectModAttByAlive(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//INT nAttType1			= m_pAuraEntry->nMiscType1[nIndex];
	//INT nAttType2			= m_pAuraEntry->nMiscType2[nIndex];
	//SpellValueType eModType	= (SpellValueType)m_pAuraEntry->nMiscValue1[nIndex];
	//INT32 nModValue         = (SpellValueType)m_pAuraEntry->nMiscValue2[nIndex];

	//// 映射属性
	//INT32 nAttType			= sEntityAttType.GetType(nAttType1, nAttType2);
	//if (!VALID(nAttType))
	//return EAAR_Immunity;

	//Scene* pScene = m_pOwner->GetScene();
	//if (!VALID(pScene))
	//return EAAR_Immunity;

	//EntityGroup& group = pScene->GetGroup(m_pOwner->GetCamp());

	//INT32 nTotalAttMod = nModValue * group.GetValidNum();

	//m_pOwner->GetAttController().LockRecal();

	//switch (eStep)
	//{
	//case EAES_Apply:
	//switch (eModType)
	//{
	//case SMT_FLAT:
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, nTotalAttMod);
	//break;

	//case SMT_PCT:
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, nTotalAttMod);
	//break;
	//}
	//m_nEffMisc[nIndex] = nTotalAttMod;
	//break;

	//case EAES_Effect:
	//switch (eModType)
	//{
	//case SMT_FLAT:
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, -m_nEffMisc[nIndex]);
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, nTotalAttMod);
	//break;

	//case SMT_PCT:
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, -m_nEffMisc[nIndex]);
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, nTotalAttMod);
	//break;
	//}
	//m_nEffMisc[nIndex] = nTotalAttMod;
	//break;

	//case EAES_Remove:
	//switch (eModType)
	//{
	//case SMT_FLAT:
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, -m_nEffMisc[nIndex]);
	//break;

	//case SMT_PCT:
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, -m_nEffMisc[nIndex]);
	//break;
	//}
	//break;
	//}

	//m_pOwner->GetAttController().RecalAtt();

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 根据当前敌方存活人数,计算属性改
//-------------------------------------------------------------------------------
func AuraEffectModAttByEnemyAlive(aura *Aura, step define.EAuraEffectStep, index int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	//INT nAttType1			= m_pAuraEntry->nMiscType1[nIndex];
	//INT nAttType2			= m_pAuraEntry->nMiscType2[nIndex];
	//SpellValueType eModType	= (SpellValueType)m_pAuraEntry->nMiscValue1[nIndex];
	//INT32 nModValue         = (SpellValueType)m_pAuraEntry->nMiscValue2[nIndex];

	//// 映射属性
	//INT32 nAttType			= sEntityAttType.GetType(nAttType1, nAttType2);
	//if (!VALID(nAttType))
	//return EAAR_Immunity;

	//Scene* pScene = m_pOwner->GetScene();
	//if (!VALID(pScene))
	//return EAAR_Immunity;

	//EntityGroup& group = pScene->GetGroup(m_pOwner->GetOtherCamp());

	//INT32 nTotalAttMod = nModValue * (X_Max_Summon_Num-group.GetValidNum());

	//m_pOwner->GetAttController().LockRecal();

	//switch (eStep)
	//{
	//case EAES_Apply:
	//switch (eModType)
	//{
	//case SMT_FLAT:
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, nTotalAttMod);
	//break;

	//case SMT_PCT:
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, nTotalAttMod);
	//break;
	//}
	//m_nEffMisc[nIndex] = nTotalAttMod;
	//break;

	//case EAES_Effect:
	//switch (eModType)
	//{
	//case SMT_FLAT:
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, -m_nEffMisc[nIndex]);
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, nTotalAttMod);
	//break;

	//case SMT_PCT:
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, -m_nEffMisc[nIndex]);
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, nTotalAttMod);
	//break;
	//}
	//m_nEffMisc[nIndex] = nTotalAttMod;
	//break;

	//case EAES_Remove:
	//switch (eModType)
	//{
	//case SMT_FLAT:
	//m_pOwner->GetAttController().ModAttFlatModValue(nAttType, -m_nEffMisc[nIndex]);
	//break;

	//case SMT_PCT:
	//m_pOwner->GetAttController().ModAttPctModValue(nAttType, -m_nEffMisc[nIndex]);
	//break;
	//}
	//break;
	//}

	//m_pOwner->GetAttController().RecalAtt();

	return define.AuraAddResult_Success
}
