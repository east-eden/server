package scene

import "github.com/east-eden/server/define"

const (
	Camp_Max_Unit   = 30  // 每个阵营最多20个单位
	Camp_Max_Spell  = 10  // 每个阵营所属技能最多10个
	Camp_Max_Energy = 100 // 阵营符文能量最大值
)

type SceneCamp struct {
	scene        *Scene
	unitList     []*SceneUnit         // 战斗unit列表
	unitMap      map[int64]*SceneUnit // 战斗unit查询列表
	actionIdx    int                  // 当前行动unit索引
	camp         define.SceneCampType // 阵营
	aliveUnitNum int32                // 存活的单位数
	playerId     int64                // 所属玩家id
	playerLevel  int32                // 玩家等级
	playerScore  int64                // 玩家战力
	playerName   string               // 玩家名字
	serverName   string               // 服务器名字
	guildName    string               // 工会名字
	guildId      int64                // 工会id
	portrait     int32                // 玩家头像id
	// INT32					m_nMasterIndex;							// 主角索引

	// 阵营所属技能
	energy    int32                // 能量
	spellList []*define.SpellEntry // 技能列表
	spellCd   []int                // 技能cd
}

func NewSceneCamp(scene *Scene, camp define.SceneCampType) *SceneCamp {
	return &SceneCamp{
		scene:     scene,
		unitList:  make([]*SceneUnit, 0, Camp_Max_Unit),
		unitMap:   make(map[int64]*SceneUnit),
		actionIdx: 0,
		camp:      camp,

		spellList: make([]*define.SpellEntry, 0, Camp_Max_Spell),
		spellCd:   make([]int, 0, Camp_Max_Spell),
	}
}

// 获取对方阵营
func (c *SceneCamp) GetOtherCamp() define.SceneCampType {
	if c.camp == define.Scene_Camp_Attack {
		return define.Scene_Camp_Defence
	} else {
		return define.Scene_Camp_Attack
	}
}

// 获取战斗单位
func (c *SceneCamp) GetUnit(id int64) (*SceneUnit, bool) {
	unit, ok := c.unitMap[id]
	return unit, ok
}

// 战斗单位死亡
func (c *SceneCamp) OnUnitDead(u *SceneUnit) {
	c.aliveUnitNum--
	c.scene.OnUnitDead(u)
}

//-----------------------------------------------------------------------------
// 目标优先级顺序
//-----------------------------------------------------------------------------
// const INT32 XFrontTarget_Priority[X_Max_Summon_Num][X_Max_Summon_Num] =
// {
// 	{0, 1, 2, 3, 4, 5},
// 	{1, 0, 2, 4, 3, 5},
// 	{2, 1, 0, 5, 4, 3},
// 	{0, 1, 2, 3, 4, 5},
// 	{1, 0, 2, 4, 3, 5},
// 	{2, 1, 0, 5, 4, 3}
// };

// const INT32 XBackTarget_Priority[X_Max_Summon_Num][X_Max_Summon_Num] =
// {
// 	{3, 4, 5, 0, 1, 2},
// 	{4, 3, 5, 1, 0, 2},
// 	{5, 4, 3, 2, 1, 0},
// 	{3, 4, 5, 0, 1, 2},
// 	{4, 3, 5, 1, 0, 2},
// 	{5, 4, 3, 2, 1, 0}
// };

// // 英雄星级对符文的等级加成
// const INT32 X_RuneLevelAddByHero[X_Hero_Max_Star+1] = {1, 1, 1, 1, 1 , 2, 2, 2, 2, 2, 2, 3, 3, 3,3,4};
// const INT32 X_RuneLevelAddByHeroStep[X_Hero_Step_Max+1] = {0, 0, 0, 0, 0, 1,1,2, 2, 2, 2};
// const INT32 X_RuneLevelAddByHeroFly[X_Hero_FlyUp_Jie+1] = {0, 0, 0, 0, 0};

// //-----------------------------------------------------------------------------
// // 初始化
// //-----------------------------------------------------------------------------
// BOOL EntityGroup::Init(Scene* pScene, ECamp eCamp)
// {
// 	Entity::Init();

// 	m_pFather			=	NULL;
// 	m_nMaxEntityNum		=	0;
// 	m_nValidEntityNum	=	0;
// 	m_n16LoopIndex		=	0;
// 	m_pScene			=	pScene;
// 	m_n64PlayerID		=	INVALID;
// 	m_eCamp				=	eCamp;
// 	m_nLocation			=	(eCamp == ESC_Attack) ? 0 : 1;
// 	m_nEnergy			=	sConstParam->nEnergyInit;
// 	m_nMasterIndex		=	INVALID;
// 	m_nPlayerScore		=	0;
// 	ZeroMemory(m_dwMountTypeID, sizeof(m_dwMountTypeID));
// 	ZeroMemory(m_szPlayerName, sizeof(m_szPlayerName) );
// 	ZeroMemory(m_nDmgModAtt, sizeof(m_nDmgModAtt) );

// 	m_setRune.clear();

// 	ZeroMemory(m_ArrayHero, sizeof(m_ArrayHero));
// 	ZeroMemory(m_pRuneEntry, sizeof(m_pRuneEntry));
// 	ZeroMemory(m_pRuneSpellEntry, sizeof(m_pRuneSpellEntry));

// 	return TRUE;
// }

// //-----------------------------------------------------------------------------
// // 更新
// //-----------------------------------------------------------------------------
// VOID EntityGroup::Update()
// {
// 	for( INT32 i = 0; i < X_Max_Summon_Num; ++i )
// 	{
// 		if( VALID(m_ArrayHero[i]) && m_ArrayHero[i]->IsValid() )
// 		{
// 			m_ArrayHero[i]->GetCombatController().Update();
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 销毁
// //-----------------------------------------------------------------------------
// VOID EntityGroup::Destroy()
// {
// 	for (INT n = 0; n < X_Max_Summon_Num; n++)
// 	{
// 		if (VALID(m_ArrayHero[n]))
// 		{
// 			sSceneMgr.DestroyEntity(m_ArrayHero[n]);
// 		}
// 	}

// 	m_pFather = NULL;
// 	m_pScene = NULL;
// }

// //-----------------------------------------------------------------------------
// // 清空所有英雄
// //-----------------------------------------------------------------------------
// VOID EntityGroup::ClearEntityHero()
// {
// 	m_nValidEntityNum = 0;
// 	m_n16LoopIndex = 0;

// 	for (INT n = 0; n < X_Max_Summon_Num; n++)
// 	{
// 		if (VALID(m_ArrayHero[n]))
// 		{
// 			sSceneMgr.DestroyEntity(m_ArrayHero[n]);
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 加入英雄
// //-----------------------------------------------------------------------------
// BOOL EntityGroup::AddEntityHero( INT nIndex, EntityHero* pHero, Player* pPlayer)
// {
// 	ASSERT( !VALID(m_ArrayHero[nIndex]) );

// 	m_ArrayHero[nIndex]		= pHero;
// 	m_nValidEntityNum++;

// 	m_nMaxEntityNum = m_nValidEntityNum;

// 	return TRUE;
// }

// BOOL EntityGroup::AddRecordEntityHero(INT nIndex, EntityHero* pHero, BOOL bRecord)
// {
// 	ASSERT( !VALID(m_ArrayHero[nIndex]) );

// 	m_ArrayHero[nIndex]		= pHero;

// 	if( !pHero->IsDead() )
// 	{
// 		m_nValidEntityNum++;
// 	}

// 	m_nMaxEntityNum++;

// 	return TRUE;
// }

// //-----------------------------------------------------------------------------
// // 获取英雄
// //-----------------------------------------------------------------------------
// EntityHero* EntityGroup::GetEntityHero( INT nIndex )
// {
// 	if (!MIsBetween(nIndex, 0, X_Max_Summon_Num))
// 		return NULL;

// 	return m_ArrayHero[nIndex];
// }

// //-----------------------------------------------------------------------------
// // 加入符文
// //-----------------------------------------------------------------------------
// VOID EntityGroup::AddRune(Player* pPlayer, INT32 nIndex, DWORD dwTypeID, INT8 n8Level)
// {
// 	// 关联英雄符文等级加成
// 	if( VALID(pPlayer) )
// 	{
// 		const tagRuneEntry* pEntry = sRuneEntry(dwTypeID);
// 		if( VALID(pEntry) )
// 		{
// 			HeroData* pData = pPlayer->GetHeroContainer().GetHeroByTypeID(pEntry->dwHeroRelation);
// 			if( VALID(pData) )
// 			{
// 				INT32 nLevelAdd = X_RuneLevelAddByHero[pData->GetStar()];
// 				INT32 nStepAdd = X_RuneLevelAddByHeroStep[pData->GetStep()];

// 				dwTypeID += nLevelAdd;
// 				n8Level += nLevelAdd;
// 				dwTypeID += nStepAdd;
// 				n8Level += nStepAdd;

// 				if ( pData->GetEntry()->eClass >= EHQ_Yellow )
// 				{
// 					BYTE Jie	=	pData->GetFlyUp() % 100 * 0.1;	//阶
// 					INT32 nFlyAdd = X_RuneLevelAddByHeroFly[Jie];
// 					dwTypeID += nFlyAdd;
// 					n8Level += nFlyAdd;
// 				}

// 			}
// 		}
// 	}

// 	m_pRuneEntry[nIndex]	=	sRuneEntry(dwTypeID);

// 	if( VALID(m_pRuneEntry[nIndex]) )
// 	{
// 		m_pRuneSpellEntry[nIndex] = sResMgr.GetSpellEntry(m_pRuneEntry[nIndex]->dwSpellID);
// 	}

// 	m_n8RuneLevel[nIndex] = n8Level;
// 	m_n8RuneWeight[nIndex] = 0;
// 	m_n8RuneCD[nIndex] = 0;
// }

// //-----------------------------------------------------------------------------
// // 攻击
// //-----------------------------------------------------------------------------
// VOID EntityGroup::Attack(Entity* pEntity)
// {
// 	EntityGroup* pTarget = static_cast<EntityGroup*>(pEntity);
// 	BOOL bBreak = FALSE;
// 	for( INT32 i = m_n16LoopIndex; i < X_Max_Summon_Num; ++i )
// 	{
// 		++m_n16LoopIndex;

// 		if( VALID(m_ArrayHero[i]) && m_ArrayHero[i]->IsValid() )
// 		{
// 			EntityHero* pHero = FindTargetByPriority(i, pTarget, TRUE);

// 			if( VALID(pHero) )
// 			{
// 				m_ArrayHero[i]->Attack(pHero);
// 				m_ArrayHero[i]->GetCombatController().CalAuraEffect(GetScene()->GetCurRound());

// 				// 风怒状态
// 				if( m_ArrayHero[i]->HasState(EHS_Anger) )
// 				{
// 					EntityHero* pHero = FindTargetByPriority(i, pTarget, TRUE);
// 					if( VALID(pHero) )
// 					{
// 						m_ArrayHero[i]->Attack(pHero);
// 					}
// 				}

// 				AddAttackNum();
// 				bBreak = TRUE;
// 			}
// 		}

// 		if( bBreak )
// 			break;
// 	}
// }

// //-----------------------------------------------------------------------------
// // 查找攻击优先级最高的目标
// //-----------------------------------------------------------------------------
// EntityHero* EntityGroup::FindTargetByPriority(INT32 nIndex, EntityGroup* pTarget, BOOL bFront)
// {
// 	EntityHero* pHero = NULL;

// 	for( INT32 i = 0; i < X_Max_Summon_Num; ++i )
// 	{
// 		pHero = pTarget->GetEntityHero(bFront ? XFrontTarget_Priority[nIndex][i] : XBackTarget_Priority[nIndex][i]);
// 		if( VALID(pHero) && pHero->IsValid() )
// 		{
// 			if( !pHero->HasState(EHS_Stealth) )
// 			{
// 				return pHero;
// 			}
// 			else
// 			{
// 				EntityHero* pCasterHero = GetEntityHero(nIndex);
// 				if( pCasterHero->HasState(EHS_AntiHidden) )
// 					return pHero;
// 			}
// 		}
// 	}

// 	for( INT32 i = 0; i < X_Max_Summon_Num; ++i )
// 	{
// 		pHero = pTarget->GetEntityHero(bFront ? XFrontTarget_Priority[nIndex][i] : XBackTarget_Priority[nIndex][i]);
// 		if( VALID(pHero) && pHero->IsValid() )
// 			return pHero;
// 	}

// 	return NULL;
// }

// //-----------------------------------------------------------------------------
// // 死亡
// //-----------------------------------------------------------------------------
// VOID EntityGroup::OnHeroDead(EntityHero* pEntity)
// {
// 	m_nValidEntityNum--;
// 	GetScene()->OnHeroDead(pEntity);
// }

// //-----------------------------------------------------------------------------
// // 释放符文技能
// //-----------------------------------------------------------------------------
// VOID EntityGroup::CastRuneSpell()
// {
// 	if( m_setRune.size() == 0 )
// 		return;

// 	if( !VALID(m_nMasterIndex) )
// 		return;

// 	RuneSet::iterator it = m_setRune.begin();

// 	INT32 nRuneIndex = (*it) / 10000;

// 	// 判断能量是否足够
// 	if( VALID(m_pRuneSpellEntry[nRuneIndex]) && m_nEnergy > m_pRuneSpellEntry[nRuneIndex]->nEnergyCost)
// 	{
// 		ModeAttEnergy(-(m_pRuneSpellEntry[nRuneIndex]->nEnergyCost));

// 		// 释放技能
// 		if( VALID(m_ArrayHero[m_nMasterIndex]) )
// 		{
// 			EntityGroup& group = GetScene()->GetGroup(GetOtherCamp());
// 			EntityHero* pTarget = FindTargetByPriority(m_nMasterIndex, &group, FALSE);

// 			m_ArrayHero[m_nMasterIndex]->CastRuneSpell(m_pRuneSpellEntry[nRuneIndex], pTarget, m_n8RuneLevel[nRuneIndex]);
// 		}

// 		m_n8RuneWeight[nRuneIndex]+= m_pRuneSpellEntry[nRuneIndex]->nRuneCD;
// 		m_n8RuneCD[nRuneIndex] += m_pRuneSpellEntry[nRuneIndex]->nRuneCD;

// 		m_setRune.erase(it);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 更新符文技能CD
// //-----------------------------------------------------------------------------
// VOID EntityGroup::UpdateRuneCD()
// {
// 	for( INT32 i = 0; i < X_Rune_Max_Group; ++i )
// 	{
// 		--m_n8RuneCD[i];
// 		m_n8RuneCD[i] = (0 > m_n8RuneCD[i]) ? 0 : m_n8RuneCD[i];

// 		if( m_n8RuneCD[i] == 0 && VALID(m_pRuneEntry[i]) )
// 		{
// 			m_setRune.insert(i * 10000 + m_n8RuneWeight[i]);
// 		}
// 	}

// }

//-----------------------------------------------------------------------------
// 改变阵营内符文能量
//-----------------------------------------------------------------------------
func (c *SceneCamp) ModAttEnergy(mod int32) {
	c.energy += mod
	if c.energy > Camp_Max_Energy {
		c.energy = Camp_Max_Energy
	}
}

// //-----------------------------------------------------------------------------
// // 同步客户端战斗实体基本属性
// //-----------------------------------------------------------------------------
// VOID EntityGroup::FillEntityInfo(OUT fxMessage& msg)
// {
// 	INT32 nEnityNum = 0;
// 	for( INT32 i = 0; i < X_Max_Summon_Num; ++i )
// 	{
// 		if( VALID(m_ArrayHero[i])  )
// 		{
// 			nEnityNum++;
// 		}
// 	}

// 	CreateProtoMsg(data, EntityInfo,);
// 	msg << (INT32)nEnityNum;

// 	EntityHero* pHero = NULL;
// 	for( INT32 i = 0; i < X_Max_Summon_Num; ++i )
// 	{
// 		if( VALID(m_ArrayHero[i])  )
// 		{
// 			m_ArrayHero[i]->FillEntityInfo(data);
// 			msg << data;
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 战斗开始时触发
// //-----------------------------------------------------------------------------
// VOID EntityGroup::TriggerByStartBehaviour()
// {
// 	for( INT32 i = 0; i < X_Max_Summon_Num; ++i )
// 	{
// 		if( VALID(m_ArrayHero[i])  )
// 		{
// 			m_ArrayHero[i]->GetCombatController().TriggerByBehaviour(EBT_Start, m_ArrayHero[i]);
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 计算帮会和符文产生的伤害改变属性
// //-----------------------------------------------------------------------------
// VOID EntityGroup::CalDmgModAtt(Player* pPlayer)
// {
// 	m_nDmgModAtt[EDM_RaceDoneKindom] += pPlayer->GetScienceSkillValue(ESCS_DoneKindom);
// 	m_nDmgModAtt[EDM_RaceTakenKindom] -= pPlayer->GetScienceSkillValue(ESCS_TakenKindom);
// 	m_nDmgModAtt[EDM_RaceDoneHell] += pPlayer->GetScienceSkillValue(ESCS_DoneHell);
// 	m_nDmgModAtt[EDM_RaceTakenHell] -= pPlayer->GetScienceSkillValue(ESCS_TakenHell);
// 	m_nDmgModAtt[EDM_RaceDoneForest] += pPlayer->GetScienceSkillValue(ESCS_DoneForest);
// 	m_nDmgModAtt[EDM_RaceTakenForest] -= pPlayer->GetScienceSkillValue(ESCS_TakenForest);
// 	m_nDmgModAtt[EDM_RaceDoneWild] += pPlayer->GetScienceSkillValue(ESCS_DoneWild);
// 	m_nDmgModAtt[EDM_RaceTakenWild] -= pPlayer->GetScienceSkillValue(ESCS_TakenWild);
// 	m_nDmgModAtt[EDM_RaceDoneOther] += pPlayer->GetScienceSkillValue(ESCS_DoneForest);
// 	m_nDmgModAtt[EDM_RaceTakenOther] -= pPlayer->GetScienceSkillValue(ESCS_TakenForest);

// 	//RuneData* pRuneData = NULL;
// 	//RuneContainer& conRune = pPlayer->GetRuneContainer();
// 	//RuneContainer::BagRune::Iterator it = conRune->Begin();
// 	//while(conRune->PeekNext(it, pRuneData))
// 	//{
// 	//	if(pRuneData->IsActive())
// 	//	{
// 	//		const tagRuneEntry* pEntry = sResMgr.GetRuneEntry(pRuneData->GetTypeID());
// 	//		if(!VALID(pEntry))
// 	//			continue;

// 	//		if( VALID(pEntry->eDmgModType) )
// 	//		{
// 	//			m_nDmgModAtt[EDM_RaceTakenWild] += pEntry->nDmgModValue;
// 	//		}
// 	//	}
// 	//}
// }

// //-----------------------------------------------------------------------------
// // 同步客户端符文
// //-----------------------------------------------------------------------------
// VOID EntityGroup::FillRuneInfo(OUT fxMessage& msg)
// {
// 	msg << (INT32)X_Rune_Max_Group;
// 	for( INT32 i = 0; i < X_Rune_Max_Group; ++i )
// 	{
// 		msg << (UINT32)(VALID(m_pRuneEntry[i]) ? m_pRuneEntry[i]->dwID : INVALID);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 导出成员ID
// //-----------------------------------------------------------------------------
// INT EntityGroup::ExportEntityID( DWORD dwTypeID[] )
// {
// 	INT nNum = 0;
// 	if ( VALID(m_nMasterIndex) && VALID(m_ArrayHero[m_nMasterIndex]))
// 	{
// 		dwTypeID[nNum++] = m_ArrayHero[m_nMasterIndex]->GetEntry()->dwTypeID;
// 	}
// 	for (INT n = 0; n < X_Max_Summon_Num; n++)
// 	{
// 		if (!VALID(m_ArrayHero[n]) || n == m_nMasterIndex)
// 			continue;

// 		dwTypeID[nNum++] = m_ArrayHero[n]->GetEntry()->dwTypeID;
// 	}

// 	return nNum;
// }

// //-----------------------------------------------------------------------------
// // 导出成员状态flag
// //-----------------------------------------------------------------------------
// INT EntityGroup::ExportEntityStateFlag(INT32 nStateFlag[])
// {
// 	INT nNum = 0;
// 	if ( VALID(m_nMasterIndex) && VALID(m_ArrayHero[m_nMasterIndex]))
// 	{
// 		nStateFlag[nNum++] = (INT32)m_ArrayHero[m_nMasterIndex]->GetStateFlag();
// 	}
// 	for (INT n = 0; n < X_Max_Summon_Num; n++)
// 	{
// 		if (!VALID(m_ArrayHero[n]) || n == m_nMasterIndex)
// 			continue;

// 		nStateFlag[nNum++] = (INT32)m_ArrayHero[n]->GetStateFlag();
// 	}

// 	return nNum;
// }

// //-----------------------------------------------------------------------------
// // 保存录像
// //-----------------------------------------------------------------------------
// VOID EntityGroup::Save2DB(tagGroupRecord* pRecord)
// {
// 	pRecord->n64PlayerID = m_n64PlayerID;
// 	pRecord->nLevel		  = m_nPlayerLevel;
// 	pRecord->nPlayerScore = m_nPlayerScore;
// 	memcpy(pRecord->szName, m_szPlayerName, sizeof(pRecord->szName) );
// 	memcpy(pRecord->nDmgModAtt, m_nDmgModAtt, sizeof(pRecord->nDmgModAtt) );

// 	EntityHero* pHero = NULL;
// 	for( INT32 i = 0; i < X_Max_Summon_Num; ++i )
// 	{
// 		if( VALID(m_ArrayHero[i])  )
// 		{
// 			m_ArrayHero[i]->Save2DB(&(pRecord->stHeroRecord[i]));
// 			m_ArrayHero[i]->Save2DmgDB(pRecord,i);
// 		}
// 	}

// 	// 保存符文
// 	for( INT32 i = 0; i < X_Rune_Max_Group; ++i )
// 	{
// 		if( VALID(m_pRuneEntry[i]) )
// 		{
// 			pRecord->dwRuneID[i] = m_pRuneEntry[i]->dwID;
// 			pRecord->n8RuneLevel[i] = m_n8RuneLevel[i];
// 		}
// 	}

// 	Player* pPlayer = sPlayerMgr.GetPlayerByGUID(m_n64PlayerID);
// 	if(!VALID(pPlayer))
// 	{
// 		pRecord->nLevel		  = m_nPlayerLevel;
// 		pRecord->nPlayerScore = m_nPlayerScore;
// 		pRecord->n16HeadProtrait = m_nProtrait;
// 		pRecord->n8HeadQuality = m_nHeadQuality;
// 		pRecord->n64GuildID = m_n64GuildID;
// 		memcpy(pRecord->szName, m_szPlayerName, sizeof(m_szPlayerName) );
// 		memcpy(pRecord->nDmgModAtt, m_nDmgModAtt, sizeof(m_nDmgModAtt) );
// 		memcpy(pRecord->szWorldName, m_szWorldName, sizeof(m_szWorldName) );
// 		memcpy(pRecord->szGuildName, m_szGuildName, sizeof(m_szGuildName) );

// 		return;
// 	}

// 	pRecord->nLevel = pPlayer->GetLevel();
// 	pRecord->n16HeadProtrait = pPlayer->GetPlayerInfo()->n16HeadProtrait;
// 	pRecord->n8HeadQuality = pPlayer->GetPlayerInfo()->n32HeadQuality;
// 	pRecord->n8VipLevel = pPlayer->GetPlayerInfo()->nVipLevel;
// 	pRecord->n8Flag = 0;
// 	pRecord->n64GuildID = pPlayer->GetGuildID();

// 	memcpy(pRecord->szName, pPlayer->GetPlayerName(), sizeof(pRecord->szName) );
// 	memcpy(pRecord->szWorldName, sServer.GetWorldName(), sizeof(pRecord->szWorldName) );
// 	memcpy(pRecord->szGuildName, pPlayer->GetPlayerInfo()->pGuildMem->szGuildName, sizeof(pRecord->szGuildName) );
// }

// //-----------------------------------------------------------------------------
// // 保存录像
// //-----------------------------------------------------------------------------
// VOID EntityGroup::SaveBeastGroupInfo(tagBeastGroupRecord* pRecord)
// {
// 	pRecord->n64PlayerID = m_n64PlayerID;
// 	pRecord->dwWorldID = sServer.GetWorldID();

// 	Player* pPlayer = sPlayerMgr.GetPlayerByGUID(m_n64PlayerID);
// 	if (!VALID(pPlayer))
// 	{
// 		pRecord->nLevel = m_nPlayerLevel;
// 		pRecord->nPlayerScore = m_nPlayerScore;
// 		pRecord->n16HeadProtrait = m_nProtrait;
// 		pRecord->n8HeadQuality = m_nHeadQuality;
// 		pRecord->n64GuildID = m_n64GuildID;
// 		memcpy(pRecord->szName, m_szPlayerName, sizeof(m_szPlayerName));
// 		memcpy(pRecord->szWorldName, m_szWorldName, sizeof(m_szWorldName));
// 		memcpy(pRecord->szGuildName, m_szGuildName, sizeof(m_szGuildName));

// 		return;
// 	}

// 	pRecord->nLevel = pPlayer->GetLevel();
// 	pRecord->nPlayerScore = pPlayer->GetPlayerScore();
// 	pRecord->n16HeadProtrait = pPlayer->GetPlayerInfo()->n16HeadProtrait;
// 	pRecord->n8HeadQuality = pPlayer->GetPlayerInfo()->n32HeadQuality;
// 	pRecord->n8VipLevel = pPlayer->GetPlayerInfo()->nVipLevel;
// 	pRecord->n8Flag = 0;
// 	pRecord->n64GuildID = pPlayer->GetGuildID();
// 	memcpy(pRecord->szName, pPlayer->GetPlayerName(), sizeof(pRecord->szName));
// 	memcpy(pRecord->szWorldName, sServer.GetWorldName(), sizeof(pRecord->szWorldName));
// 	memcpy(pRecord->szGuildName, pPlayer->GetPlayerInfo()->pGuildMem->szGuildName, sizeof(pRecord->szGuildName));
// }

// //-----------------------------------------------------------------------------
// // 保存录像
// //-----------------------------------------------------------------------------
// VOID EntityGroup::SaveBeastRecord(tagBeastRecord* pRecord)
// {
// 	EntityHero* pHero = NULL;
// 	if (!VALID(m_ArrayHero[1]))
// 		return;

// 	m_ArrayHero[1]->Save2DB(pRecord);
// }

// //导出成员的等级
// VOID EntityGroup::ExportEntityLevel(DWORD dwLevel[])
// {
// 	INT nNum = 0;
// 	if ( VALID(m_nMasterIndex) && VALID(m_ArrayHero[m_nMasterIndex]))
// 	{
// 		dwLevel[nNum++] = m_ArrayHero[m_nMasterIndex]->GetLevel();
// 	}
// 	for (INT n = 0; n < X_Max_Summon_Num; n++)
// 	{
// 		if (!VALID(m_ArrayHero[n]) || n == m_nMasterIndex)
// 			continue;
// 		dwLevel[nNum++] = m_ArrayHero[n]->GetLevel();
// 	}
// }

// //导出成员的星级
// VOID EntityGroup::ExportEntityStar(DWORD dwStar[])
// {
// 	INT nNum = 0;
// 	if ( VALID(m_nMasterIndex) && VALID(m_ArrayHero[m_nMasterIndex]))
// 	{
// 		dwStar[nNum++] = m_ArrayHero[m_nMasterIndex]->GetStar();
// 	}
// 	for (INT n = 0; n < X_Max_Summon_Num; n++)
// 	{
// 		if (!VALID(m_ArrayHero[n]) || n == m_nMasterIndex)
// 			continue;
// 		dwStar[nNum++] = m_ArrayHero[n]->GetStar();
// 	}
// }

// //导出成员的品质
// VOID EntityGroup::ExportEntityQuality(DWORD dwQuality[])
// {
// 	INT nNum = 0;
// 	if ( VALID(m_nMasterIndex) && VALID(m_ArrayHero[m_nMasterIndex]))
// 	{
// 		dwQuality[nNum++] = m_ArrayHero[m_nMasterIndex]->GetQuality();
// 	}
// 	for (INT n = 0; n < X_Max_Summon_Num; n++)
// 	{
// 		if (!VALID(m_ArrayHero[n]) || n == m_nMasterIndex)
// 			continue;
// 		dwQuality[nNum++] = m_ArrayHero[n]->GetQuality();
// 	}
// }
