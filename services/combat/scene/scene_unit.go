package scene

import (
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	log "github.com/rs/zerolog/log"
	"github.com/willf/bitset"
)

const (
	Unit_Energy_OnBeDamaged = 2 // 受伤害增加能量
)

type SceneUnit struct {
	id           int64
	level        uint32
	posX         int16              // x坐标
	posY         int16              // y坐标
	TauntId      int64              // 被嘲讽目标
	v2           define.Vector2     // 朝向
	scene        *Scene             // 场景
	camp         *SceneCamp         // 场景阵营
	normalSpell  *define.SpellEntry // 普攻技能
	specialSpell *define.SpellEntry // 特殊技能

	// 伤害统计
	totalDmgRecv int64 // 总共受到的伤害
	totalDmgDone int64 // 总共造成的伤害
	totalHeal    int64 // 总共产生治疗
	attackNum    int   // 攻击次数

	opts *UnitOptions
}

func (s *SceneUnit) Guid() int64 {
	return s.id
}

func (s *SceneUnit) GetLevel() uint32 {
	return s.level
}

func (s *SceneUnit) GetScene() *Scene {
	return s.opts.Scene
}

func (s *SceneUnit) GetCamp() int32 {
	return 0
}

func (s *SceneUnit) CombatCtrl() *CombatCtrl {
	return s.opts.CombatCtrl
}

func (s *SceneUnit) Opts() *UnitOptions {
	return s.opts
}

func (s *SceneUnit) UpdateSpell() {
	log.Info().
		Int64("id", s.id).
		Int32("type_id", s.opts.TypeId).
		Floats32("pos", s.opts.Position[:]).
		Msg("creature start UpdateSpell")

	s.CombatCtrl().Update()
}

func (s *SceneUnit) AddState(e define.EHeroState) {
	s.opts.State.Set(uint(e))
}

func (s *SceneUnit) HasState(e define.EHeroState) bool {
	return s.opts.State.Test(uint(e))
}

func (s *SceneUnit) HasStateAny(flag uint32) bool {
	compare := bitset.From([]uint64{uint64(flag)})
	return s.opts.State.Intersection(compare).Any()
}

func (s *SceneUnit) GetState64() uint64 {
	return s.opts.State.Bytes()[0]
}

func (s *SceneUnit) HasImmunityAny(tp define.EImmunityType, flag uint32) bool {
	compare := bitset.From([]uint64{uint64(flag)})
	return s.opts.Immunity[tp].Intersection(compare).Any()
}

//-----------------------------------------------------------------------------
// 初始化
//-----------------------------------------------------------------------------
func (s *SceneUnit) InitByScene(scene *Scene, camp *SceneCamp, posX, posY int16, v2 define.Vector2) {
	s.posX = posX
	s.posY = posY
	s.scene = scene
	s.camp = camp

	s.initSpell()
	s.initAura()
}

//-----------------------------------------------------------------------------
// 进攻
//-----------------------------------------------------------------------------
func (s *SceneUnit) Attack(target *SceneUnit) {
	// 死亡状态
	if s.HasState(define.HeroState_Dead) {
		return
	}

	// 无法行动状态
	if s.HasStateAny(1<<define.HeroState_Freeze | 1<<define.HeroState_Solid | 1<<define.HeroState_Stun | 1<<define.HeroState_Paralyzed) {
		return
	}

	// todo 释放特殊技能
	// if (GetAttController().GetAttValue(EHA_CurRage) >= GetAttController().GetAttValue(EHA_RageThreshold) && !HasState(EHS_Seal) && !HasState(EHS_Chaos) &&!HasState(EHS_Taunt) )
	// {
	// 	GetCombatController().CastSpell(m_pSpellEntry, this, pTarget, FALSE, 0, ERMT_Rage);
	// }

	// 普通攻击技能 -- 处于封印、混乱、被嘲讽状态时
	if s.HasStateAny(1<<define.HeroState_Seal | 1<<define.HeroState_Chaos | 1<<define.HeroState_Taunt) {
		if s.HasState(define.HeroState_Taunt) {
			targetCamp, ok := s.scene.GetSceneCamp(s.camp.GetOtherCamp())
			if !ok {
				log.Error().Int("target_camp", int(s.camp.GetOtherCamp())).Msg("cannot get target camp")
				return
			}

			var pass bool
			target, pass = targetCamp.GetUnit(s.TauntId)
			if !pass {
				log.Error().Int64("taunt_id", s.TauntId).Msg("cannot get target")
				return
			}
		}

		s.opts.CombatCtrl.CastSpell(s.normalSpell, s, target, false)

		// 普通攻击技能
	} else {
		if s.opts.CombatCtrl.TriggerByBehaviour(define.BehaviourType_BeforeNormal, target, -1, -1, define.SpellType_Null) == 0 {
			s.opts.CombatCtrl.CastSpell(s.normalSpell, s, target, false)
		}
	}
}

//-----------------------------------------------------------------------------
// 反击
//-----------------------------------------------------------------------------
func (s *SceneUnit) BeatBack(target *SceneUnit) {
	if s.HasState(define.HeroState_Dead) {
		return
	}

	if !s.HasStateAny(1<<define.HeroState_Freeze | 1<<define.HeroState_Solid | 1<<define.HeroState_Stun | 1<<define.HeroState_Paralyzed) {
		s.opts.CombatCtrl.CastSpell(s.normalSpell, s, target, false)
	}
}

//-----------------------------------------------------------------------------
// 死亡
//-----------------------------------------------------------------------------
func (s *SceneUnit) OnDead(caster *SceneUnit, spellId int32) {
	if s.HasState(define.HeroState_Dead) {
		return
	}

	s.camp.OnUnitDead(s)

	// 清空当前值
	s.opts.AttManager.SetAttValue(define.Att_Plus_CurHP, 0)

	// 设置为死亡状态
	s.AddState(define.HeroState_Dead)
}

//-----------------------------------------------------------------------------
// 造成伤害
//-----------------------------------------------------------------------------
func (s *SceneUnit) OnDamage(target *SceneUnit, dmgInfo *CalcDamageInfo) {
	s.opts.CombatCtrl.TriggerByDmgMod(true, target, dmgInfo)
}

//-----------------------------------------------------------------------------
// 改变符文能量
//-----------------------------------------------------------------------------
func (s *SceneUnit) ModAttEnergy(mod int32) {
	s.camp.ModAttEnergy(mod)
}

//-----------------------------------------------------------------------------
// 承受伤害
//-----------------------------------------------------------------------------
func (s *SceneUnit) OnBeDamaged(caster *SceneUnit, dmgInfo *CalcDamageInfo) {
	s.opts.CombatCtrl.TriggerByDmgMod(false, caster, dmgInfo)

	if define.DmgInfo_Damage == dmgInfo.Type {
		//// 计算怒气恢复
		//if( (DmgInfo.dwProcEx & EAEE_RageResume) && !HasState(EHS_Seal))
		//{
		//	GetAttController().ModAttValue(EHA_CurRage, X_Rage_Resume);
		//}

		// 计算能量恢复
		if dmgInfo.ProcEx&(1<<define.AuraEventEx_EnergyResume) != 0 {
			s.ModAttEnergy(Unit_Energy_OnBeDamaged)
		}
	}
}

//-----------------------------------------------------------------------------
// 处理伤害
//-----------------------------------------------------------------------------
func (s *SceneUnit) DoneDamage(caster *SceneUnit, dmgInfo *CalcDamageInfo) {
	if dmgInfo.Damage <= 0 {
		return
	}

	if dmgInfo.Type == define.DmgInfo_Null {
		return
	}

	switch dmgInfo.Type {
	// 伤害
	case define.DmgInfo_Damage:
		if s.HasState(define.HeroState_UnBeat) || s.HasState(define.HeroState_ImmunityGroupDmg) && (dmgInfo.ProcEx&1<<define.AuraEventEx_GroupDmg != 0) {
			dmgInfo.Damage = 0
			dmgInfo.ProcEx |= (1 << define.AuraEventEx_Immnne)
		} else if s.HasState(define.HeroState_UnDead) {
			if int64(s.opts.AttManager.GetAttValue(define.Att_Plus_CurHP)) <= dmgInfo.Damage {
				dmgInfo.Damage = int64(s.opts.AttManager.GetAttValue(define.Att_Plus_CurHP) - 1)
				s.opts.AttManager.SetAttValue(define.Att_Plus_CurHP, 1)

				// 伤害统计
				s.totalDmgRecv += dmgInfo.Damage
				caster.totalDmgDone += dmgInfo.Damage

				dmgInfo.ProcTarget |= (1 << define.AuraEvent_Taken_Any_Damage)
				dmgInfo.ProcEx |= (1 << define.AuraEventEx_UnDead)
			} else {
				// 伤害统计
				s.totalDmgRecv += dmgInfo.Damage
				caster.totalDmgDone += dmgInfo.Damage

				s.opts.AttManager.ModAttValue(define.Att_Plus_CurHP, int(-dmgInfo.Damage))
			}
		} else {
			// 伤害统计
			s.totalDmgRecv += dmgInfo.Damage
			caster.totalDmgDone += dmgInfo.Damage

			s.opts.AttManager.ModAttValue(define.Att_Plus_CurHP, int(-dmgInfo.Damage))

			if s.opts.AttManager.GetAttValue(define.Att_Plus_CurHP) <= 0 {
				// 刚刚死亡
				s.OnDead(caster, dmgInfo.SpellId)
			}
		}

		// 治疗
	case define.DmgInfo_Heal:
		s.opts.AttManager.ModAttValue(define.Att_Plus_CurHP, int(dmgInfo.Damage))

		// 治疗统计
		s.totalHeal += dmgInfo.Damage

		// 安抚
	case define.DmgInfo_Placate:
		// 减少怒气
		// s.opts.AttManager.ModAttValue(define.Att_Plus_CurRage, -dmgInfo.Damage)

		// 激怒
	case define.DmgInfo_Enrage:
		if !s.HasState(define.HeroState_Seal) {
			// s.opts.AttManager.ModAttValue(define.Att_CurRage, dmgInfo.Damage)
		}
	}
}

//-----------------------------------------------------------------------------
// 进入状态
//-----------------------------------------------------------------------------
func (s *SceneUnit) AddToState(state define.EHeroState) {
	s.opts.CombatCtrl.TriggerByServentState(state, true)
}

//-----------------------------------------------------------------------------
// 脱离状态
//-----------------------------------------------------------------------------
func (s *SceneUnit) EscFromState(state define.EHeroState) {
	s.opts.CombatCtrl.TriggerByServentState(state, false)
}

//-----------------------------------------------------------------------------
// 免疫
//-----------------------------------------------------------------------------
func (s *SceneUnit) AddToImmunity(immunityType define.EImmunityType, immunity int) {
	switch immunityType {
	case define.ImmunityType_Mechanic:
		// 删除指定类型的Aura
		//s.opts.CombatCtrl.RemoveAura(immunity)
	}
}

//-----------------------------------------------------------------------------
// 初始化伤害加成
//-----------------------------------------------------------------------------
func (s *SceneUnit) InitDmgModAtt() {
	// memcpy(m_nDmgModAtt, static_cast<EntityGroup*>(m_pFather)->GetDmgModAtt(), sizeof(m_nDmgModAtt));

	// 	switch( m_pEntry->eClass )
	// 	{
	// 	case EHC_Tank:
	// 		{
	// 			m_nDmgModAtt[EDM_DamageDonePctPhysics] -= 2000;
	// 			m_nDmgModAtt[EDM_DamageDonePctMagic] -= 2000;
	// 			m_nDmgModAtt[EDM_DamageTakenPctPhysics] -= 3000;
	// 			m_nDmgModAtt[EDM_DamageTakenPctMagic] -= 3000;
	// 		}
	// 		break;

	// 	case EHC_Berserker:
	// 		{
	// 			m_nDmgModAtt[EDM_DamageDonePctPhysics] += 1000;
	// 			m_nDmgModAtt[EDM_DamageDonePctMagic] += 1000;
	// 		}
	// 		break;

	// 	case EHC_Assassin:
	// 		{
	// 			m_nDmgModAtt[EDM_DamageDonePctPhysics] += 2000;
	// 			m_nDmgModAtt[EDM_DamageDonePctMagic] += 2000;
	// 			m_nDmgModAtt[EDM_DamageTakenPctPhysics] += 2000;
	// 			m_nDmgModAtt[EDM_DamageTakenPctMagic] += 2000;
	// 		}
	// 		break;

	// 	case EHC_Elementer:
	// 		{
	// 			m_nDmgModAtt[EDM_DamageDonePctPhysics] += 1000;
	// 			m_nDmgModAtt[EDM_DamageDonePctMagic] += 1000;
	// 		}
	// 		break;

	// 	case EHC_Healer:
	// 		{
	// 			m_nDmgModAtt[EDM_DamageDonePctHeal] += 1000;
	// 			m_nDmgModAtt[EDM_DamageTakenPctPhysics] += 1000;
	// 			m_nDmgModAtt[EDM_DamageTakenPctMagic] += 1000;
	// 		}
	// 		break;

	// 	default:
	// 		break;
	// 	}
}

//-----------------------------------------------------------------------------
// 属性初始化
//-----------------------------------------------------------------------------
func (s *SceneUnit) InitAttribute(heroInfo *define.HeroInfo) {
	// todo 读取静态表中的状态掩码
	// s.opts.State = bitset.From([]uint64{uint64(s.opts.Entry.dwStateMask)})

	// todo 免疫
	for n := 0; n < define.ImmunityType_End; n++ {
		// s.opts.Immunity[n] = bitset.From([]uint64{uint64(s.opts.Entry.dwImmunity[n])})
	}

	// todo AttEntry
	// auto.GetAttEntry(s.opts.Entry.BaseAttId)
	heroEntry, ok := auto.GetHeroEntry(int(heroInfo.TypeId))
	if !ok {
		log.Warn().Int32("type_id", heroInfo.TypeId).Msg("cannot find hero entry")
		return
	}

	s.opts.AttManager.SetBaseAttId(int32(heroEntry.AttID))
	s.opts.AttManager.CalcAtt()
	s.opts.AttManager.SetAttValue(define.Att_Plus_CurHP, s.opts.AttManager.GetAttValue(define.Att_Plus_MaxHP))
}

// //-----------------------------------------------------------------------------
// // 属性初始化
// //-----------------------------------------------------------------------------
// VOID EntityHero::InitAttribute(const tagSquareBeast* pBeast, DWORD dwAttID/* = INVALID */)
// {
// 	if (!VALID(m_pEntry))
// 		return;

// 	// 状态
// 	m_State.Import(m_pEntry->dwStateMask);

// 	// 免疫
// 	for (INT n = 0; n < EIT_End; ++n)
// 	{
// 		m_Immunity[n].Import(m_pEntry->dwImmunity[n]);
// 	}

// 	DWORD dwBaseAttID = VALID(dwAttID) ? dwAttID : pBeast->pEntry->dwBaseAttID;
// 	const tagAttEntry* pAttEntry = sAttEntry(dwBaseAttID);

// 	INT64 n64PlayerID = static_cast<EntityGroup*>(GetFather())->GetPlayerID();
// 	m_AttController.InitAttribute(pAttEntry, pBeast, n64PlayerID);

// 	if (VALID(pBeast->nCurHP) && !VALID(dwAttID))
// 	{
// 		GetAttController().SetAttValue(EHA_CurHP, pBeast->nCurHP);
// 	}
// 	else
// 	{
// 		GetAttController().SetAttValue(EHA_CurHP, GetAttController().GetAttValue(EHA_MaxHP));
// 	}
// }

// //-----------------------------------------------------------------------------
// // 技能初始化
// //-----------------------------------------------------------------------------
// VOID EntityHero::InitSpell()
// {
// 	if (!VALID(m_pEntry))
// 		return;

// 	// 设置初始技能
// 	m_pMeleeEntry = NULL;
// 	m_pSpellEntry = NULL;

// 	m_pMeleeEntry = m_pEntry->pMeleeSpell;
// 	m_pSpellEntry = m_pEntry->pRageSpell;

// 	// 时装技能
// 	const tagFashionEntry* pFashionEntry = sResMgr.GetFashionEntry(m_nFashionID);
// 	if(VALID(pFashionEntry))
// 	{
// 		const tagSpellEntry* pFashionMeleeEntry = sResMgr.GetSpellEntry(pFashionEntry->dwMeleeSpellID);
// 		if(VALID(pFashionMeleeEntry))
// 			m_pMeleeEntry = pFashionMeleeEntry;

// 		const tagSpellEntry* pFashionRageEntry = sResMgr.GetSpellEntry(pFashionEntry->dwRageSpellID);
// 		if(VALID(pFashionRageEntry))
// 			m_pSpellEntry = pFashionRageEntry;
// 	}

// 	// 被动技能
// 	for(INT32 i = 0; i < X_Passive_Spell_Num; ++i )
// 	{
// 		if( !m_AttController.CastPassiveSpell(i) )
// 			break;
// 	}
// }

// //-----------------------------------------------------------------------------
// // 初始化被动技能
// //-----------------------------------------------------------------------------
// VOID EntityHero::InitAura()
// {
// 	// 增加初始被动Aura
// 	for( INT32 n = 0; n < X_Hero_Aura_Init; ++n )
// 	{
// 		if( !VALID(m_pEntry->dwPassiveAuraID[n]) )
// 			break;

// 		GetCombatController().AddAura(m_pEntry->dwPassiveAuraID[n], this);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 设置普通攻击
// //-----------------------------------------------------------------------------
// VOID EntityHero::SetMeleeSpell(DWORD dwSpellID)
// {
// 	const tagSpellEntry* pSpell = sSpellEntry(dwSpellID);
// 	m_pMeleeEntry = pSpell;
// }

// //-------------------------------------------------------------------------------
// // 状态
// //-------------------------------------------------------------------------------
// VOID EntityHero::AddState( EHeroState eState, INT nCount /*= 1*/ )
// {
// 	bool bNewState = !HasState(eState);

// 	m_State.Set(eState, nCount);

// 	if (bNewState)
// 	{
// 		Scene* pScene = GetScene();
// 		if (VALID(pScene) && !pScene->IsOnlyRecord() )
// 		{
// 			CreateSceneProtoMsg(msg, MS_SetState,);
// 			*msg << (UINT32)GetLocation();
// 			*msg << (UINT32)eState;
// 			pScene->AddMsgList(msg);
// 		}

// 		// 追加状态处理
// 		AddToState(eState);
// 	}
// }

// VOID EntityHero::DecState( EHeroState eState, INT nCount /*= 1*/ )
// {
// 	if( !HasState(eState) )
// 		return;

// 	m_State.Unset(eState, nCount);

// 	if( !HasState(eState) )
// 	{
// 		Scene* pScene = GetScene();
// 		if (VALID(pScene) && !pScene->IsOnlyRecord() )
// 		{
// 			CreateSceneProtoMsg(msg, MS_UnsetState, );
// 			*msg << (UINT32)GetLocation();
// 			*msg << (UINT32)eState;
// 			pScene->AddMsgList(msg);
// 		}

// 		EscFromState(eState);
// 	}
// }

// //-------------------------------------------------------------------------------
// // 保存录像
// //-------------------------------------------------------------------------------
// VOID EntityHero::Save2DB(tagHeroRecord* pRecord)
// {
// 	pRecord->dwEntityID = m_pEntry->dwTypeID;
// 	pRecord->nFashionID = m_nFashionID;
// 	pRecord->dwMountTypeID = m_dwMountTypeID;
// 	pRecord->nStateFlag = m_n16HeroState;
// 	pRecord->nFlyUp = m_nFly_Up;
// 	pRecord->nLevel = m_nLevel;
// 	pRecord->nRageLevel = m_n16RageLevel;
// 	pRecord->nStarLevel = m_nStar;
// 	pRecord->nQuality = m_nQuality;
// 	memcpy(pRecord->nAtt, m_AttRecord.ExportAtt(), sizeof(pRecord->nAtt));
// 	memcpy(pRecord->nBaseAtt, m_AttRecord.ExportBaseAtt(), sizeof(pRecord->nBaseAtt));
// 	memcpy(pRecord->nBaseAttModPct, m_AttRecord.ExportBaseAttModPct(), sizeof(pRecord->nBaseAttModPct));
// 	memcpy(pRecord->nAttMod, m_AttRecord.ExportAttMod(), sizeof(pRecord->nAttMod));
// 	memcpy(pRecord->nAttModPct, m_AttRecord.ExportAttModPct(), sizeof(pRecord->nAttModPct));
// 	memcpy(pRecord->dwPassiveSpell, m_AttRecord.ExportPassiveSpell(), sizeof(pRecord->dwPassiveSpell));
// }

// //-------------------------------------------------------------------------------
// // 保存录像
// //-------------------------------------------------------------------------------
// VOID EntityHero::Save2DmgDB(tagGroupRecord* pRecord, INT16 n16Index)
// {
// 	for ( int i = EHDM_RaceDoneKindom; i < EHDM_End; i++ )
// 	{
// 		pRecord->nHeroDmgModAtt[n16Index][i] = m_nHeroDmgModAtt[EDM_RaceDoneKindom + i];
// 	}
// }

// //-------------------------------------------------------------------------------
// // 保存录像
// //-------------------------------------------------------------------------------
// VOID EntityHero::Save2DB(tagBeastRecord* pRecord)
// {
// 	pRecord->dwTypeID = m_dwBeastTypeID;
// 	pRecord->nStepLevel = m_nBeastStepLevel;
// 	pRecord->dwEntityID = m_pEntry->dwTypeID;
// 	pRecord->nLevel = m_nLevel;
// 	pRecord->nQuality = m_nQuality;
// 	memcpy(pRecord->nAtt, m_AttRecord.ExportAtt(), sizeof(pRecord->nAtt));
// 	memcpy(pRecord->nBaseAtt, m_AttRecord.ExportBaseAtt(), sizeof(pRecord->nBaseAtt));
// 	memcpy(pRecord->nBaseAttModPct, m_AttRecord.ExportBaseAttModPct(), sizeof(pRecord->nBaseAttModPct));
// 	memcpy(pRecord->nAttMod, m_AttRecord.ExportAttMod(), sizeof(pRecord->nAttMod));
// 	memcpy(pRecord->nAttModPct, m_AttRecord.ExportAttModPct(), sizeof(pRecord->nAttModPct));
// 	memcpy(pRecord->dwPassiveSpell, m_AttRecord.ExportPassiveSpell(), sizeof(pRecord->dwPassiveSpell));
// }

// //-----------------------------------------------------------------------------
// // 初始化伤害加成
// //-----------------------------------------------------------------------------
// VOID EntityHero::InitHeroDmgModAtt(const HeroData* pHeroData,INT8 nPos)
// {
// 	ZeroMemory(m_nHeroDmgModAtt,sizeof(m_nHeroDmgModAtt));

// 	if ( !VALID(pHeroData) ) return;

// 	INT64 n64PlayerID = static_cast<EntityGroup*>(GetFather())->GetPlayerID();
// 	Player* pPlayer = sPlayerMgr.GetPlayerByGUID(n64PlayerID);
// 	if(!VALID(pPlayer)) return;

// 	MahjongData* pMahjongData = NULL;
// 	MahjongContainer& conMahjong = pPlayer->GetMahjongContainer();

// 	pMahjongData = conMahjong.GetMahjongGroup(pHeroData->GetGroupPos());
// 	if (VALID(pMahjongData)) // 有宠物守护
// 	{
// 		const DWORD* dwLinkID = pMahjongData->GetLinkID();
// 		//todo
// 		for ( INT i=0; i < EMAHT_End; i++ )
// 		{
// 			const tagMahjongLinkEntry* pLinkEntry = sMahjongLinkEntry(dwLinkID[i]);
// 			if(VALID(pLinkEntry))
// 			{
// 				const tagMahjongLinkAttEntry* pLinkAttEntry = sMahjongLinkAttEntry(pLinkEntry->dwLinkAttID);
// 				if ( VALID(pLinkAttEntry) )
// 				{
// 					for ( INT i = EHDM_RaceDoneKindom; i < EHDM_End; i++ )
// 					{
// 						if( pLinkAttEntry->nHeroValue[i] > 0 )
// 						{
// 							if ( i >= EHDM_RaceTakenKindom )
// 							{
// 								m_nDmgModAtt[i+EDM_RaceDoneKindom] -= pLinkAttEntry->nHeroValue[i];
// 								m_nHeroDmgModAtt[i+EDM_RaceDoneKindom] -= pLinkAttEntry->nHeroValue[i];
// 							}
// 							else
// 							{
// 								m_nDmgModAtt[i+EDM_RaceDoneKindom] += pLinkAttEntry->nHeroValue[i];
// 								m_nHeroDmgModAtt[i+EDM_RaceDoneKindom] += pLinkAttEntry->nHeroValue[i];
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 		//麻将牌加成
// 		for ( int i=0; i < 8; i++ )
// 		{
// 			const tagMahjongInfoEntry* pMahjonInfoEntry = sMahjongInfoEntry(pMahjongData->GetMahjongSize(i));
// 			if(VALID(pMahjonInfoEntry))
// 			{
// 				for ( INT i = EHDM_RaceDoneKindom; i < EHDM_End; i++ )
// 				{
// 					if( pMahjonInfoEntry->nHeroValue[i] > 0 )
// 					{
// 						if ( i >= EHDM_RaceTakenKindom )
// 						{
// 							m_nDmgModAtt[i+EDM_RaceDoneKindom] -= pMahjonInfoEntry->nHeroValue[i];
// 							m_nHeroDmgModAtt[i+EDM_RaceDoneKindom] -= pMahjonInfoEntry->nHeroValue[i];
// 						}
// 						else
// 						{
// 							m_nDmgModAtt[i+EDM_RaceDoneKindom] += pMahjonInfoEntry->nHeroValue[i];
// 							m_nHeroDmgModAtt[i+EDM_RaceDoneKindom] += pMahjonInfoEntry->nHeroValue[i];
// 						}
// 					}

// 				}
// 			}
// 		}

// 		//清一色
// 		if ( pMahjongData->IsLinkAll() )
// 		{
// 			DWORD dwSameSuitAttID = pMahjongData->GetSameSuitID();
// 			if ( VALID(dwSameSuitAttID) )
// 			{
// 				const tagMahjongLinkEntry* pLinkEntry = sMahjongLinkEntry(dwSameSuitAttID);
// 				if ( VALID(pLinkEntry) )
// 				{
// 					const tagMahjongLinkAttEntry* pLinkAttEntry = sMahjongLinkAttEntry(pLinkEntry->dwLinkAttID);
// 					if ( VALID(pLinkAttEntry) )
// 					{
// 						for ( INT i = EHDM_RaceDoneKindom; i < EHDM_End; i++ )
// 						{
// 							if( pLinkAttEntry->nHeroValue[i] > 0 )
// 							{
// 								if ( i >= EHDM_RaceTakenKindom )
// 								{
// 									m_nDmgModAtt[i+EDM_RaceDoneKindom] -= pLinkAttEntry->nHeroValue[i];
// 									m_nHeroDmgModAtt[i+EDM_RaceDoneKindom] -= pLinkAttEntry->nHeroValue[i];
// 								}
// 								else
// 								{
// 									m_nDmgModAtt[i+EDM_RaceDoneKindom] += pLinkAttEntry->nHeroValue[i];
// 									m_nHeroDmgModAtt[i+EDM_RaceDoneKindom] += pLinkAttEntry->nHeroValue[i];
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}

// 	}

// }

// VOID EntityHero::InitHeroRecordDmgModAtt(const tagGroupRecord* pRecord,INT8 nPos)
// {
// 	ZeroMemory(m_nHeroDmgModAtt,sizeof(m_nHeroDmgModAtt));
// 	// 录像
// 	if ( VALID(pRecord) )
// 	{
// 		for ( INT i = EHDM_RaceDoneKindom; i < EHDM_End; i++ )
// 		{
// 			m_nDmgModAtt[i+EDM_RaceDoneKindom] += pRecord->nHeroDmgModAtt[nPos][i];
// 			m_nHeroDmgModAtt[i+EDM_RaceDoneKindom] += pRecord->nHeroDmgModAtt[nPos][i];
// 		}
// 	}
// }
