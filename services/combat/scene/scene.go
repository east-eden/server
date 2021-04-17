package scene

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"bitbucket.org/funplus/server/utils"
	"bitbucket.org/funplus/server/utils/random"
	"github.com/emirpasic/gods/maps/treemap"
	god_utils "github.com/emirpasic/gods/utils"
	"github.com/hellodudu/task"
	log "github.com/rs/zerolog/log"
)

var (
	Scene_InitSpellNum = 50 // 每个scene初始有50个技能正在作用
)

type Scene struct {
	opts *SceneOptions
	*task.Tasker

	id          int64
	entityIdGen int64
	entityMap   *treemap.Map // 战斗unit列表
	curRound    int32
	maxRound    int32
	result      chan bool
	rand        *random.FakeRandom
	camps       [define.Scene_Camp_End]*SceneCamp

	comFinishList *list.List // com条结束的entity列表
	spellList     *list.List // 场景内技能列表

	wg utils.WaitGroupWrapper
	sync.RWMutex
}

func (s *Scene) Init(sceneId int64, opts ...SceneOption) *Scene {
	s.id = sceneId
	s.entityMap = treemap.NewWith(god_utils.Int64Comparator)
	s.comFinishList = list.New()
	s.spellList = list.New()
	s.result = make(chan bool, 1)
	s.opts = DefaultSceneOptions()
	s.rand = random.NewFakeRandom(int(time.Now().Unix()))
	s.Tasker = task.NewTasker(1)

	for n := define.Scene_Camp_Begin; n < define.Scene_Camp_End; n++ {
		s.camps[n] = NewSceneCamp(s, n)
	}

	for _, o := range opts {
		o(s.opts)
	}

	// add attack unit list
	for _, unitInfo := range s.opts.AttackEntityList {
		err := s.AddEntityByPB(s.camps[define.Scene_Camp_Attack], unitInfo)
		utils.ErrPrint(err, "AddEntityByPB failed when Scene.Init", sceneId, s.opts.SceneEntry.Id, unitInfo.HeroTypeId)
	}

	// add defence unit list
	for _, unitInfo := range s.opts.DefenceEntityList {
		err := s.AddEntityByPB(s.camps[define.Scene_Camp_Defence], unitInfo)
		_ = utils.ErrCheck(err, "AddEntityByPB failed when Scene.Init", sceneId, s.opts.SceneEntry.Id, unitInfo.HeroTypeId)
	}

	// unit group
	unitGroupEntry := s.opts.UnitGroupEntry
	if unitGroupEntry != nil {
		for idx := range unitGroupEntry.HeroIds {
			// hero id invalid
			if unitGroupEntry.HeroIds[idx] == -1 {
				continue
			}

			heroEntry, ok := auto.GetHeroEntry(unitGroupEntry.HeroIds[idx])
			if !ok {
				continue
			}

			// camp invalid
			if !utils.BetweenInt32(unitGroupEntry.Camps[idx], define.Scene_Camp_Begin, define.Scene_Camp_End) {
				continue
			}

			err := s.AddEntityByOptions(
				s.camps[unitGroupEntry.Camps[idx]],
				WithEntityTypeId(unitGroupEntry.HeroIds[idx]),
				WithEntityHeroEntry(heroEntry),
				WithEntityPosition(unitGroupEntry.PosXs[idx], unitGroupEntry.PosZs[idx], int32(unitGroupEntry.Rotates[idx])),
				WithEntityInitAtbValue(int32(unitGroupEntry.InitComs[idx])),
			)

			_ = utils.ErrCheck(err, "AddEntityByOptions failed when Scene.Init", unitGroupEntry.HeroIds[idx])
		}
	}

	// tasker init
	s.Tasker.Init(
		task.WithStartFn(func() {
			s.start()
		}),

		task.WithContextDoneFn(func() {
			log.Info().
				Int32("scene_type_id", s.opts.SceneEntry.Id).
				Int64("scene_id", s.GetId()).
				Msg("scene context done...")
		}),

		task.WithUpdateFn(func() {
			s.update()
		}),

		task.WithSleep(time.Millisecond*100),
	)

	return s
}

func (s *Scene) start() {
	it := s.entityMap.Iterator()
	for it.Next() {
		it.Value().(*SceneEntity).OnSceneStart()
	}
}

func (s *Scene) update() {
	s.updateEntities()
	s.updateCamps()
}

func (s *Scene) Run(ctx context.Context) error {
	return s.Tasker.Run(ctx)
}

func (s *Scene) Exit(ctx context.Context) {
	s.wg.Wait()
}

func (s *Scene) GetId() int64 {
	return s.id
}

func (s *Scene) GetCamp(camp int32) *SceneCamp {
	return s.camps[camp]
}

func (s *Scene) GetEntity(id int64) (*SceneEntity, bool) {
	val, ok := s.entityMap.Get(id)
	if ok {
		return val.(*SceneEntity), ok
	}

	return nil, ok
}

func (s *Scene) GetEntityMap() *treemap.Map {
	return s.entityMap
}

// 寻找单位
func (s *Scene) findEnemyEntityByHead(camp int32) (*SceneEntity, bool) {
	if s.entityMap.Size() == 0 {
		return nil, false
	}

	it := s.entityMap.Iterator()
	for it.Next() {
		e := it.Value().(*SceneEntity)
		if e.GetCamp().camp != camp {
			return e, true
		}
	}

	return nil, false
}

func (s *Scene) GetResult() bool {
	return <-s.result
}

func (s *Scene) GetSceneCamp(camp int32) (*SceneCamp, bool) {
	if !utils.BetweenInt32(camp, define.Scene_Camp_Begin, define.Scene_Camp_End) {
		return nil, false
	}

	return s.camps[camp], true
}

// todo
func (s *Scene) IsOnlyRecord() bool {
	return false
}

func (s *Scene) Rand(min, max int) int {
	return s.rand.RandSection(min, max)
}

func (s *Scene) SendDamage(dmgInfo *CalcDamageInfo) {
	// todo
}

func (s *Scene) updateCamps() {

	// 是否攻击方先手
	bAttackFirst := true

	for ; s.curRound+1 <= s.maxRound; s.curRound++ {
		bEnterNextRound := false

		s.updateRuneCD()

		nActionRount := 0
		for !bEnterNextRound {
			nActionRount++

			if bAttackFirst {
				s.camps[int(define.Scene_Camp_Attack)].Update()

				// 本轮攻击没有结束
				if !s.camps[int(define.Scene_Camp_Attack)].IsLoopEnd() {
					s.camps[int(define.Scene_Camp_Attack)].Attack(s.camps[int(define.Scene_Camp_Defence)])
				}

				s.camps[int(define.Scene_Camp_Defence)].Update()

				// 本轮攻击没有结束
				if !s.camps[int(define.Scene_Camp_Defence)].IsLoopEnd() {
					s.camps[int(define.Scene_Camp_Defence)].Attack(s.camps[int(define.Scene_Camp_Attack)])
				}
			} else {
				s.camps[int(define.Scene_Camp_Defence)].Update()

				// 本轮攻击没有结束
				if !s.camps[int(define.Scene_Camp_Defence)].IsLoopEnd() {
					s.camps[int(define.Scene_Camp_Defence)].Attack(s.camps[int(define.Scene_Camp_Defence)])
				}

				s.camps[int(define.Scene_Camp_Attack)].Update()

				// 本轮攻击没有结束
				if !s.camps[int(define.Scene_Camp_Attack)].IsLoopEnd() {
					s.camps[int(define.Scene_Camp_Attack)].Attack(s.camps[int(define.Scene_Camp_Defence)])
				}
			}

			if s.camps[int(define.Scene_Camp_Attack)].IsLoopEnd() &&
				s.camps[int(define.Scene_Camp_Defence)].IsLoopEnd() {
				for i := nActionRount; i < Camp_Max_Unit; i++ {
					s.camps[int(define.Scene_Camp_Defence)].Update()
					s.camps[int(define.Scene_Camp_Attack)].Update()
				}

				nActionRount = 0
				bEnterNextRound = true
			}

			// 释放符文技能
			// compile comment
			// s.UpdateRuneSpell(bAttackFirst);
		}

		// 重置攻击顺序
		s.camps[int(define.Scene_Camp_Attack)].ResetLoopIndex()
		s.camps[int(define.Scene_Camp_Defence)].ResetLoopIndex()

		// 战斗结束
		if !s.camps[int(define.Scene_Camp_Attack)].IsValid() ||
			!s.camps[int(define.Scene_Camp_Defence)].IsValid() {
			break
		}
	}
}

// 更新场景内技能
func (s *Scene) updateSpells() {
	var next *list.Element
	for e := s.spellList.Front(); e != nil; e = next {
		next = e.Next()

		skill := e.Value.(*Skill)
		skill.Update()

		// 删除已作用玩的技能
		if skill.IsCompleted() {
			s.spellList.Remove(e)
		}
	}
}

func (s *Scene) updateEntities() {
	it := s.entityMap.Iterator()
	for it.Next() {
		it.Value().(*SceneEntity).Update()
	}
}

func (s *Scene) AddEntityByPB(camp *SceneCamp, unitInfo *pbGlobal.EntityInfo) error {
	entry, ok := auto.GetHeroEntry(unitInfo.HeroTypeId)
	if !ok {
		return fmt.Errorf("GetUnitEntry failed: type_id<%d>", unitInfo.HeroTypeId)
	}

	id := atomic.AddInt64(&s.entityIdGen, 1)
	e, err := NewSceneEntity(
		id,
		WithEntityTypeId(unitInfo.HeroTypeId),
		WithEntityAttList(unitInfo.AttValue),
		WithEntityHeroEntry(entry),
	)

	if err != nil {
		return err
	}

	s.entityMap.Put(id, e)

	return nil
}

func (s *Scene) AddEntityByOptions(camp *SceneCamp, opts ...EntityOption) error {
	id := atomic.AddInt64(&s.entityIdGen, 1)
	opts = append(opts, WithEntitySceneCamp(camp))
	e, err := NewSceneEntity(id, opts...)
	if err != nil {
		return err
	}

	s.entityMap.Put(id, e)
	return nil
}

func (s *Scene) ClearEntities() {
	s.entityMap.Clear()
}

// //-----------------------------------------------------------------------------
// // 更新加载怪物
// //-----------------------------------------------------------------------------
// VOID Scene::UpdateEntityLoad()
// {
// 	INT nMaxWave = VALID(m_pEntityGroupEntry) ? m_pEntityGroupEntry->nMaxWave : 1;
// 	if( m_nWave >= nMaxWave )
// 	{
// 		ASSERT(0);
// 		return;
// 	}

// 	// 夫妻塔特殊处理
// 	if(m_nWave == 0 && m_pEntry->eType == ESBT_CoupleTower)
// 	{
// 		Player* pLeader = sPlayerMgr.GetPlayerByGUID(m_stMisc.couple_tower.n64LeaderID);
// 		if(!VALID(pLeader))
// 			return;

// 		Marriage* pMarriage = sMarriageMgr.GetMarriage(pLeader->GetMarriageID());
// 		if(!VALID(pMarriage))
// 			return;

// 		Player* pSpouse = sPlayerMgr.GetPlayerByGUID(pMarriage->GetSpouseID(pLeader));
// 		if(!VALID(pSpouse))
// 			return;

// 		if(VALID(m_stGroupRecord[ESC_Attack]) && m_stGroupRecord[ESC_Attack]->IsValid())
// 		{
// 			AddRecordPlayer(ESC_Attack, m_stGroupRecord[ESC_Attack]);
// 		}
// 		else
// 		{
// 			AddCouple(ESC_Attack, pMarriage, pLeader, pSpouse);
// 		}
// 	}

// 	// 飞升特殊处理
// 	else if(m_nWave == 0 && m_pEntry->eType == ESBT_FlyUp)
// 	{
// 		Player* pPlayer = sPlayerMgr.GetPlayerByGUID(m_n64Creator);
// 		if( !VALID(pPlayer) )
// 			return;
// 		DWORD			pHero[X_JA_Hero_Max];			// 选择的英雄
// 		memset(pHero, 0xff, sizeof(pHero) );
// 		pHero[m_stMisc.flyup.HeroIdx]	=	m_stMisc.flyup.dwHeroSerial;
// 		pHero[m_stMisc.flyup.AidIdx]	=	m_stMisc.flyup.dwAidSerial;
// 		AddPlayer(ESC_Attack, pPlayer, X_Max_Summon_Num, pHero);
// 		m_MuitlGroup[ESC_Attack].SetMasterIndex(m_stMisc.flyup.HeroIdx);
// 	}

// 	// 山海特殊处理
// 	else if (m_pEntry->eType == ESBT_Square)
// 	{
// 		// 没有录像的话重新加载
// 		if (!m_MuitlGroup[ESC_Attack].IsValid())
// 		{
// 			Player* pAttack = sPlayerMgr.GetPlayerByGUID(m_MuitlGroup[ESC_Attack].GetPlayerID());
// 			if (!VALID(pAttack))
// 				return;

// 			tagSquareBeast* pBeast = pAttack->GetBeastContainer().GetBeast(m_stMisc.square.dwAttackBeastTypeID);
// 			if (!VALID(pBeast))
// 				return;

// 			if (VALID(pBeast->nCurHP) && pBeast->nCurHP <= 0)
// 				return;

// 			// 进攻阵容
// 			AddBeast(ESC_Attack, pAttack, pBeast);
// 		}

// 		if (!m_MuitlGroup[ESC_Defence].IsValid())
// 		{
// 			// 玩家防守阵容
// 			EntityHero* pEntity = NULL;
// 			DWORD dwEntityID = INVALID;
// 			if (!VALID(m_pEntityGroupEntry))
// 			{
// 				tagSquareSceneData* pData = sSceneMgr.GetSquareSceneData(m_dwSquareSerial);
// 				if (!VALID(pData))
// 					return;

// 				// 防守录像
// 				tagGroupRecord* pDefRecord = sAuxSquare.GetEntityRecord(pData->stDefenceBeast.n64PlayerID, pData->nDefenceIndex);
// 				if (!VALID(pDefRecord))
// 					return;

// 				m_stGroupRecord[ESC_Defence] = pDefRecord;
// 				AddRecordPlayer(ESC_Defence, m_stGroupRecord[ESC_Defence]);
// 			}
// 		}
// 	}

// 	// 生成双方战斗阵容
// 	else if( m_nWave == 0 && !IsRecord() )
// 	{
// 		// 使用玩家阵容备份创建队伍
// 		if( VALID(m_stGroupRecord[ESC_Attack]) )
// 		{
// 			AddRecordPlayer(ESC_Attack, m_stGroupRecord[ESC_Attack]);
// 		}
// 		else  // 使用玩家当前阵容创建队伍
// 		{
// 			if( VALID(m_MuitlGroup[ESC_Attack].GetPlayerID() ) )
// 			{
// 				Player* pAttack = sPlayerMgr.GetPlayerByGUID(m_MuitlGroup[ESC_Attack].GetPlayerID());
// 				if( VALID(pAttack) )
// 				{
// 					AddPlayer(ESC_Attack, pAttack);
// 				}
// 			}

// 		}

// 		// 使用玩家阵容备份创建队伍
// 		if( VALID(m_stGroupRecord[ESC_Defence]) )
// 		{
// 			AddRecordPlayer(ESC_Defence, m_stGroupRecord[ESC_Defence]);
// 		}
// 		else	// 使用玩家当前阵容创建队伍
// 		{
// 			if( VALID(m_MuitlGroup[ESC_Defence].GetPlayerID() ) )
// 			{
// 				Player* pDefence = sPlayerMgr.GetPlayerByGUID(m_MuitlGroup[ESC_Defence].GetPlayerID());
// 				if( VALID(pDefence) )
// 				{
// 					AddPlayer(ESC_Defence, pDefence);
// 				}
// 			}

// 		}
// 	}

// 	m_nCurRound = 0;

// 	// 加载怪物
// 	if( !VALID(m_MuitlGroup[ESC_Defence].GetPlayerID()) && !VALID(m_stGroupRecord[ESC_Defence]) && !IsRecord())
// 	{
// 		if(!VALID(m_pEntityGroupEntry))
// 			return;

// 		m_MuitlGroup[ESC_Defence].ClearEntityHero();

// 		// 加载当前轮数怪物
// 		EntityHero* pEntity = NULL;
// 		DWORD dwEntityID = INVALID;
// 		DWORD dwRuneTypeID = INVALID;
// 		for( INT32 i = 0 ; i < X_Max_Summon_Num; ++i )
// 		{
// 			dwEntityID = m_pEntityGroupEntry->stGroup[m_nWave].dwEntityID[i];
// 			if( VALID(dwEntityID) )
// 			{
// 				// 山海怪物boss处理
// 				if (m_pEntry->eType == ESBT_Square)
// 				{
// 					tagSquareSceneData* pSquareData = sSceneMgr.GetSquareSceneData(m_dwSquareSerial);
// 					if (!VALID(pSquareData))
// 						continue;

// 					INT32 nCurHP = pSquareData->stDefenceBeast.nEntityCurHP[i];
// 					if (nCurHP == 0)
// 						continue;
// 				}

// 				pEntity = sSceneMgr.CreateEntity();
// 				if( VALID(pEntity) )
// 				{
// 					pEntity->SetFather(&m_MuitlGroup[ESC_Defence]);

// 					if( pEntity->Init(NULL, dwEntityID, i) )
// 					{
// 						m_MuitlGroup[ESC_Defence].AddEntityHero(i, pEntity, NULL);
// 					}
// 					else
// 					{
// 						sSceneMgr.DestroyEntity(pEntity);
// 					}
// 				}
// 			}

// 			// 队长
// 			if (VALID(pEntity) && (i == m_pEntityGroupEntry->stGroup[m_nWave].nMasterIndex))
// 			{
// 				switch (m_pEntry->eType)
// 				{
// 				case ESBT_Devil:
// 					pEntity->GetAttController().SetAttValue(EHA_CurHP, m_stMisc.devil_offer.nCurHP);
// 					break;
// 				}
// 			}
// 		}

// 		for( INT32 i = 0 ; i < X_Rune_Max_Group; ++i )
// 		{
// 			// 加载符文
// 			dwRuneTypeID = m_pEntityGroupEntry->stGroup[m_nWave].dwRuneTypeID[i];
// 			if( VALID(dwRuneTypeID) )
// 			{
// 				m_MuitlGroup[ESC_Defence].AddRune(NULL, i, dwRuneTypeID, 0);
// 			}
// 		}

// 		m_MuitlGroup[ESC_Defence].SetMasterIndex(m_pEntityGroupEntry->stGroup[m_nWave].nMasterIndex);
// 	}

// 	// 山海设置野怪战力和boss血量
// 	if (m_pEntry->eType == ESBT_Square && VALID(m_stMisc.square.dwEntityGroupID))
// 	{
// 		INT32 nPlayerScore = 0;
// 		for (INT n = 0; n < X_Max_Summon_Num; ++n)
// 		{
// 			EntityHero* pEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(n);
// 			if (VALID(pEntity))
// 			{
// 				nPlayerScore += (pEntity->GetAttController().GetAttValue(EHA_AttackPower) * 2 + pEntity->GetAttController().GetAttValue(EHA_MaxHP) / 4);
// 			}
// 		}

// 		m_MuitlGroup[ESC_Defence].SetPlayerScore(nPlayerScore);

// 		if (m_stMisc.square.bIsBoss)
// 		{
// 			tagSquareSceneData* pSquareData = sSceneMgr.GetSquareSceneData(m_dwSquareSerial);
// 			if (VALID(pSquareData))
// 			{
// 				for (INT n = 0; n < X_Max_Summon_Num; ++n)
// 				{
// 					EntityHero* pEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(n);
// 					INT32 nCurHP = pSquareData->stDefenceBeast.nEntityCurHP[n];
// 					if (VALID(pEntity) && VALID(nCurHP))
// 					{
// 						pEntity->GetAttController().SetAttValue(EHA_CurHP, nCurHP);
// 						pEntity->GetAttRecord().SetAttValue(EHA_CurHP, nCurHP);
// 					}
// 				}
// 			}
// 		}
// 	}

// 	m_nWave++;

// 	if( !m_bOnlyRecord )
// 	{
// 		// 同步客户端对阵双方数据
// 		CreateProtoMsg(msg, MS_WaveInfo,);
// 		msg << (UINT)(m_nWave * 100 + nMaxWave);
// 		msg << (UINT)INVALID;
// 		msg << (UINT)INVALID;
// 		m_MuitlGroup[ESC_Attack].FillRuneInfo(msg);
// 		m_MuitlGroup[ESC_Defence].FillRuneInfo(msg);
// 		m_MuitlGroup[ESC_Attack].FillEntityInfo(msg);
// 		m_MuitlGroup[ESC_Defence].FillEntityInfo(msg);

// 		SendSceneMessage(NULL, msg);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 更新符文技能释放
// //-----------------------------------------------------------------------------
// VOID Scene::UpdateRuneSpell(BOOL bAttackFirst)
// {
// 	if( !m_MuitlGroup[ESC_Attack].IsValid() ||
// 		!m_MuitlGroup[ESC_Defence].IsValid() )
// 		return;

// 	if( bAttackFirst)
// 	{
// 		m_MuitlGroup[ESC_Attack].CastRuneSpell();

// 		if( !m_MuitlGroup[ESC_Defence].IsValid() )
// 			return;

// 		m_MuitlGroup[ESC_Defence].CastRuneSpell();
// 	}
// 	else
// 	{
// 		m_MuitlGroup[ESC_Defence].CastRuneSpell();

// 		if( !m_MuitlGroup[ESC_Attack].IsValid() )
// 			return;

// 		m_MuitlGroup[ESC_Attack].CastRuneSpell();
// 	}
// }

//-----------------------------------------------------------------------------
// 更新符文技CD
//-----------------------------------------------------------------------------
func (s *Scene) updateRuneCD() {
	// compile comment
	// m_MuitlGroup[ESC_Attack].updateRuneCD();
	// m_MuitlGroup[ESC_Defence].updateRuneCD();
}

// //-----------------------------------------------------------------------------
// // 销毁
// //-----------------------------------------------------------------------------
// VOID Scene::Destroy()
// {
// 	if( sCfgMgr.GetInt32(EGV_DmgStatistics) > 0 )
// 	{
// 		MsgLog("=============================================\r\n");

// 		for (INT n = 0; n < X_Max_Summon_Num; n++)
// 		{
// 			EntityHero* pEntity = m_MuitlGroup[ESC_Attack].GetEntityHero(n);
// 			if( !VALID(pEntity) )
// 			{
// 				MsgLog("%d: empty\r\n", n);
// 				continue;
// 			}

// 			const tagEntityEntry* pEntry = pEntity->GetEntry();
// 			if( VALID(pEntity) )
// 			{
// 				WCHAR tstrWorldName[SHORT_STRING] = {0};
// 				MultiByteToWideChar(CP_UTF8, 0, sResMgr.GetHeroName(pEntry->dwTypeID).c_str(), -1, tstrWorldName, SHORT_STRING);

// 				CHAR szANSIWorldName[SHORT_STRING] = {0};
// 				WideCharToMultiByte(CP_ACP, 0, tstrWorldName, -1, szANSIWorldName, SHORT_STRING, 0, 0);

// 				MsgLog("%d: %s damage is %d\r\n", n, szANSIWorldName, pEntity->GetTotalDmgDone());
// 			}
// 		}

// 		MsgLog("attack total damage is %d\r\n", m_MuitlGroup[ESC_Attack].GetTotalDmgDone());

// 		for (INT n = 0; n < X_Max_Summon_Num; n++)
// 		{
// 			EntityHero* pEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(n);
// 			if( !VALID(pEntity) )
// 			{
// 				MsgLog("%d: empty\r\n", n);
// 				continue;
// 			}

// 			const tagEntityEntry* pEntry = pEntity->GetEntry();
// 			if( VALID(pEntity) )
// 			{
// 				if( sResMgr.GetHeroName(pEntry->dwTypeID).empty() )
// 				{
// 					MsgLog("%d: robot damage is %d\r\n", n, pEntity->GetTotalDmgDone());
// 				}
// 				else
// 				{
// 					WCHAR tstrWorldName[SHORT_STRING] = {0};
// 					MultiByteToWideChar(CP_UTF8, 0, sResMgr.GetHeroName(pEntry->dwTypeID).c_str(), -1, tstrWorldName, SHORT_STRING);

// 					CHAR szANSIWorldName[SHORT_STRING] = {0};
// 					WideCharToMultiByte(CP_ACP, 0, tstrWorldName, -1, szANSIWorldName, SHORT_STRING, 0, 0);

// 					MsgLog("%d: %s damage is %d\r\n", n, szANSIWorldName, pEntity->GetTotalDmgDone());
// 				}
// 			}
// 		}

// 		MsgLog("defence total damage is %d\r\n", m_MuitlGroup[ESC_Defence].GetTotalDmgDone());

// 		MsgLog("=============================================\r\n");
// 	}

// 	m_MuitlGroup[ESC_Attack].Destroy();
// 	m_MuitlGroup[ESC_Defence].Destroy();
// }

// //-----------------------------------------------------------------------------
// // 缓存消息
// //-----------------------------------------------------------------------------
// VOID Scene::AddMsgList( fxMessage* msg )
// {
// 	if (MIsBetween(m_nCurRound, 0, X_Max_Battle_Round + 1))
// 	{
// 		m_listMsg[m_nCurRound].push_back(msg);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 清空缓存消息
// //-----------------------------------------------------------------------------
// VOID Scene::ClrMsgList()
// {
// 	list<fxMessage*>::iterator it;
// 	for (INT n = 0; n < X_Max_Battle_Round+1; n++)
// 	{
// 		it = m_listMsg[n].begin();
// 		while( it != m_listMsg[n].end() )
//  		{
// 			SAFE_DEL((*it));
// 			++it;
// 		}

// 		m_listMsg[n].clear();
// 	}
// }

// //-----------------------------------------------------------------------------
// // 发送缓存消息
// //-----------------------------------------------------------------------------
// VOID Scene::SendMsgList()
// {
// 	if( VALID(m_dwTeamSerial) || IsNeedBroadCast() )
// 	{
// 		for (INT n = 0; n < X_Max_Battle_Round+1; n++)
// 		{
// 			list<fxMessage*>::iterator it = m_listMsg[n].begin();
// 			for (; it != m_listMsg[n].end(); ++it)
// 			{
// 				fxMessage* pCacheMsg = (*it);

// 				SendSceneMessage(NULL, *pCacheMsg);

// 				SAFE_DEL((*it));
// 			}

// 			m_listMsg[n].clear();
// 		}
// 	}
// 	// 本服玩家
// 	else if (!VALID(m_dwServerID) || (m_dwServerID == sServer.GetWorldID()))
// 	{
// 		Player* pPlayer = sPlayerMgr.GetPlayerByGUID(m_n64Creator);
// 		if(VALID(pPlayer) && pPlayer->IsOnline() )
// 		{
// 			for (INT n = 0; n < X_Max_Battle_Round+1; n++)
// 			{
// 				list<fxMessage*>::iterator it = m_listMsg[n].begin();
// 				for (; it != m_listMsg[n].end(); ++it)
// 				{
// 					fxMessage* pCacheMsg = (*it);

// 					SendSceneMessage(pPlayer, *pCacheMsg);

// 					SAFE_DEL((*it));
// 				}

// 				m_listMsg[n].clear();
// 			}
// 		}
// 		else
// 		{
// 			ClrMsgList();
// 		}
// 	}
// 	else
// 	{
// 		for (INT n = 0; n < X_Max_Battle_Round+1; n++)
// 		{
// 			list<fxMessage*>::iterator it = m_listMsg[n].begin();
// 			for (; it != m_listMsg[n].end(); ++it)
// 			{
// 				fxMessage* pCacheMsg = (*it);

// 				SendSceneMessage(NULL, *pCacheMsg);

// 				SAFE_DEL((*it));
// 			}

// 			m_listMsg[n].clear();
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 玩家加入场景
// //-----------------------------------------------------------------------------
// VOID Scene::AddPlayer(ECamp eCamp, Player* pPlayer)
// {
// 	HeroData* pHeroData = NULL;
// 	EntityHero* pEntity = NULL;

// 	m_MuitlGroup[eCamp].SetPlayerLevel(pPlayer->GetLevel());
// 	m_MuitlGroup[eCamp].SetProtrait(pPlayer->GetPlayerInfo()->n16HeadProtrait);
// 	m_MuitlGroup[eCamp].SetHeadQuality(pPlayer->GetPlayerInfo()->n32HeadQuality);
// 	m_MuitlGroup[eCamp].SetPlayerScore(pPlayer->GetPlayerScore());
// 	m_MuitlGroup[eCamp].SetPlayerName(pPlayer->GetPlayerName());
// 	m_MuitlGroup[eCamp].SetMountTypeID(pPlayer->GetPlayerInfo()->dwMountTypeID);
// 	m_MuitlGroup[eCamp].SetWorldName(sServer.GetWorldName() );
// 	m_MuitlGroup[eCamp].SetGuildName(pPlayer->GetPlayerInfo()->pGuildMem->szGuildName);
// 	m_MuitlGroup[eCamp].SetGuildID(pPlayer->GetPlayerInfo()->pGuildMem->n64GuildID);

// 	// 计算帮会和符文产生的伤害改变属性
// 	m_MuitlGroup[eCamp].CalDmgModAtt(pPlayer);

// 	for( INT32 i = 0 ; i < X_Max_Summon_Num; ++i )
// 	{
// 		pHeroData = pPlayer->GetHeroContainer().GetHeroByField(i);

// 		if( VALID(pHeroData) )
// 		{
// 			pEntity = sSceneMgr.CreateEntity();
// 			if( VALID(pEntity) )
// 			{
// 				pEntity->SetFather(&m_MuitlGroup[eCamp]);
// 				if( pEntity->Init(pHeroData, pHeroData->GetTypeID(), i) )
// 				{
// 					m_MuitlGroup[eCamp].AddEntityHero(i, pEntity, pPlayer);

// 					if( MIsMaster(pHeroData->GetTypeID()) )
// 					{
// 						m_MuitlGroup[eCamp].SetMasterIndex(i);
// 					}
// 				}
// 				else
// 				{
// 					sSceneMgr.DestroyEntity(pEntity);
// 				}
// 			}
// 		}
// 	}

// 	for( INT32 i = 0 ; i < X_Rune_Max_Group; ++i )
// 	{
// 		// 加载符文
// 		const RuneData* pRuneData = pPlayer->GetRuneContainer().GetRuneGroup(i);
// 		if( VALID(pRuneData) )
// 		{
// 			m_MuitlGroup[eCamp].AddRune(pPlayer, i, pRuneData->GetTypeID(), pRuneData->GetLevel());
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 夫妻加入场景
// //-----------------------------------------------------------------------------
// VOID Scene::AddCouple(ECamp eCamp, Marriage* pMarriage, Player* pLeader, Player* pSpouse)
// {
// 	if(!VALID(pMarriage) || !VALID(pLeader) || !VALID(pSpouse))
// 		return;

// 	HeroData* pHeroData = NULL;
// 	EntityHero* pEntity = NULL;

// 	const tagCPTowerData& stCPData = pMarriage->GetCPTowerData();
// 	INT32 nPlayerLevel = max(pLeader->GetLevel(), pSpouse->GetLevel());

// 	m_MuitlGroup[eCamp].SetPlayerLevel(nPlayerLevel);
// 	m_MuitlGroup[eCamp].SetProtrait(pLeader->GetPlayerInfo()->n16HeadProtrait);
// 	m_MuitlGroup[eCamp].SetHeadQuality(pLeader->GetPlayerInfo()->n32HeadQuality);
// 	m_MuitlGroup[eCamp].SetPlayerScore(pLeader->GetPlayerScore());
// 	m_MuitlGroup[eCamp].SetPlayerName(pLeader->GetPlayerName());
// 	m_MuitlGroup[eCamp].SetMountTypeID(pLeader->GetPlayerInfo()->dwMountTypeID);
// 	m_MuitlGroup[eCamp].SetWorldName(sServer.GetWorldName() );
// 	m_MuitlGroup[eCamp].SetGuildName(pLeader->GetPlayerInfo()->pGuildMem->szGuildName);
// 	m_MuitlGroup[eCamp].SetGuildID(pLeader->GetPlayerInfo()->pGuildMem->n64GuildID);

// 	// 计算帮会和符文产生的伤害改变属性
// 	m_MuitlGroup[eCamp].CalDmgModAtt(pLeader);

// 	for( INT32 i = 0 ; i < X_Max_Summon_Num; ++i )
// 	{
// 		if(!stCPData.stHero[i].Valid())
// 			continue;

// 		Player* pPlayer = NULL;
// 		if(stCPData.stHero[i].n64OwnerID == pLeader->GetID())
// 			pPlayer = pLeader;
// 		if(stCPData.stHero[i].n64OwnerID == pSpouse->GetID())
// 			pPlayer = pSpouse;

// 		if(!VALID(pPlayer))
// 			return;

// 		pHeroData = pPlayer->GetHeroContainer().GetHero(stCPData.stHero[i].stData.dwHeroID);
// 		if( VALID(pHeroData) )
// 		{
// 			pEntity = sSceneMgr.CreateEntity();
// 			if( VALID(pEntity) )
// 			{
// 				pEntity->SetFather(&m_MuitlGroup[eCamp]);
// 				if( pEntity->Init(pHeroData, pHeroData->GetTypeID(), i) )
// 				{
// 					m_MuitlGroup[eCamp].AddEntityHero(i, pEntity, pLeader);

// 					if( MIsMaster(pHeroData->GetTypeID()) && stCPData.stHero[i].n64OwnerID == pLeader->GetID() )
// 					{
// 						m_MuitlGroup[eCamp].SetMasterIndex(i);
// 					}
// 				}
// 				else
// 				{
// 					sSceneMgr.DestroyEntity(pEntity);
// 				}
// 			}
// 		}
// 	}

// 	for( INT32 i = 0 ; i < X_Rune_Max_Group; ++i )
// 	{
// 		// 加载符文
// 		const RuneData* pRuneData = pLeader->GetRuneContainer().GetRuneGroup(i);
// 		if( VALID(pRuneData) )
// 		{
// 			m_MuitlGroup[eCamp].AddRune(pLeader, i, pRuneData->GetTypeID(), pRuneData->GetLevel());
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 异兽加入场景
// //-----------------------------------------------------------------------------
// VOID Scene::AddBeast(ECamp eCamp, Player* pPlayer, tagSquareBeast* pBeast, DWORD dwAttID/* = INVALID */)
// {
// 	if (!VALID(pPlayer) || !VALID(pBeast))
// 		return;

// 	// 异兽默认站立前排中央
// 	INT32 nPos = 1;

// 	m_MuitlGroup[eCamp].SetPlayerID(pPlayer->GetID());
// 	m_MuitlGroup[eCamp].SetPlayerLevel(pPlayer->GetLevel());
// 	m_MuitlGroup[eCamp].SetProtrait(pPlayer->GetPlayerInfo()->n16HeadProtrait);
// 	m_MuitlGroup[eCamp].SetHeadQuality(pPlayer->GetPlayerInfo()->n32HeadQuality);
// 	m_MuitlGroup[eCamp].SetPlayerName(pPlayer->GetPlayerName());
// 	m_MuitlGroup[eCamp].SetMountTypeID(pPlayer->GetPlayerInfo()->dwMountTypeID);
// 	m_MuitlGroup[eCamp].SetWorldName(sServer.GetWorldName());
// 	m_MuitlGroup[eCamp].SetGuildName(pPlayer->GetPlayerInfo()->pGuildMem->szGuildName);
// 	m_MuitlGroup[eCamp].SetGuildID(pPlayer->GetPlayerInfo()->pGuildMem->n64GuildID);

// 	EntityHero* pEntity = sSceneMgr.CreateEntity();
// 	if (VALID(pEntity))
// 	{
// 		pEntity->SetFather(&m_MuitlGroup[eCamp]);
// 		if (pEntity->InitBeast(pBeast, pBeast->pEntry->dwEntityID, nPos, dwAttID))
// 		{
// 			m_MuitlGroup[eCamp].AddEntityHero(nPos, pEntity, pPlayer);
// 			m_MuitlGroup[eCamp].SetMasterIndex(nPos);

// 			INT32 nBeastScore = (pEntity->GetAttController().GetAttValue(EHA_AttackPower) * 2 + pEntity->GetAttController().GetAttValue(EHA_MaxHP) / 4);
// 			m_MuitlGroup[eCamp].SetPlayerScore(nBeastScore);
// 		}
// 		else
// 		{
// 			sSceneMgr.DestroyEntity(pEntity);
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 根据临时阵容创建战斗数据
// //-----------------------------------------------------------------------------
// VOID Scene::AddPlayer(ECamp eCamp, Player* pPlayer,INT32 nHeroNum,const DWORD *pHero)
// {
// 	HeroData* pHeroData = NULL;
// 	EntityHero* pEntity = NULL;

// 	m_MuitlGroup[eCamp].SetPlayerLevel(pPlayer->GetLevel());
// 	m_MuitlGroup[eCamp].SetProtrait(pPlayer->GetPlayerInfo()->n16HeadProtrait);
// 	m_MuitlGroup[eCamp].SetHeadQuality(pPlayer->GetPlayerInfo()->n32HeadQuality);
// 	m_MuitlGroup[eCamp].SetPlayerScore(pPlayer->GetPlayerScore());
// 	m_MuitlGroup[eCamp].SetPlayerName(pPlayer->GetPlayerName());
// 	m_MuitlGroup[eCamp].SetMountTypeID(pPlayer->GetPlayerInfo()->dwMountTypeID);
// 	m_MuitlGroup[eCamp].SetWorldName(sServer.GetWorldName() );
// 	m_MuitlGroup[eCamp].SetGuildName(pPlayer->GetPlayerInfo()->pGuildMem->szGuildName);
// 	m_MuitlGroup[eCamp].SetGuildID(pPlayer->GetPlayerInfo()->pGuildMem->n64GuildID);

// 	// 计算帮会和符文产生的伤害改变属性
// 	m_MuitlGroup[eCamp].CalDmgModAtt(pPlayer);

// 	for( INT32 i = 0 ; i < nHeroNum; ++i )
// 	{
// 		if( VALID(pHero[i]) )
// 		{
// 			pHeroData = pPlayer->GetHeroContainer().GetHero(pHero[i]);
// 			if( VALID(pHeroData) )
// 			{
// 				pEntity = sSceneMgr.CreateEntity();
// 				if( VALID(pEntity) )
// 				{
// 					pEntity->SetFather(&m_MuitlGroup[eCamp]);
// 					if( pEntity->Init(pHeroData, pHeroData->GetTypeID(), i) )
// 					{
// 						m_MuitlGroup[eCamp].AddEntityHero(i, pEntity, pPlayer);

// 						if( MIsMaster(pHeroData->GetTypeID()) )
// 						{
// 							m_MuitlGroup[eCamp].SetMasterIndex(i);
// 						}
// 					}
// 					else
// 					{
// 						sSceneMgr.DestroyEntity(pEntity);
// 					}
// 				}
// 			}
// 		}
// 	}

// 	for( INT32 i = 0 ; i < X_Rune_Max_Group; ++i )
// 	{
// 		// 加载符文
// 		const RuneData* pRuneData = pPlayer->GetRuneContainer().GetRuneGroup(i);
// 		if( VALID(pRuneData) )
// 		{
// 			m_MuitlGroup[eCamp].AddRune(pPlayer, i, pRuneData->GetTypeID(), pRuneData->GetLevel());
// 		}
// 	}
// }
// //-----------------------------------------------------------------------------
// // 根据给定阵容创建战斗数据
// //-----------------------------------------------------------------------------
// VOID Scene::AddPlayer(ECamp eCamp, INT64 n64PlayerID, INT32 nHeroNum, INT32 nRuneNum, const DWORD *pHero, const DWORD *pRune)
// {
// 	if ( !VALID(pHero) || !VALID(pRune) )
// 		return;

// 	tagPlayerInfo *pInfo = sPlayerMgr.GetPlayerInfoByGUID( n64PlayerID );
// 	if ( !VALID(pInfo) )
// 		return;

// 	m_MuitlGroup[eCamp].SetPlayerLevel(pInfo->nLevel);
// 	m_MuitlGroup[eCamp].SetProtrait(pInfo->n16HeadProtrait);
// 	m_MuitlGroup[eCamp].SetHeadQuality(pInfo->n32HeadQuality);
// 	m_MuitlGroup[eCamp].SetPlayerScore(pInfo->nPlayerScore);
// 	m_MuitlGroup[eCamp].SetPlayerName(pInfo->szPlayerName);
// 	m_MuitlGroup[eCamp].SetWorldName(sServer.GetWorldName() );
// 	m_MuitlGroup[eCamp].SetGuildName(pInfo->pGuildMem->szGuildName);
// 	m_MuitlGroup[eCamp].SetGuildID(pInfo->pGuildMem->n64GuildID);

// 	m_MuitlGroup[ESC_Attack].ClearEntityHero();

// 	// 加载当前轮数怪物
// 	EntityHero* pEntity = NULL;
// 	DWORD dwEntityID = INVALID;
// 	DWORD dwRuneTypeID = INVALID;
// 	for( INT32 i = 0 ; i < nHeroNum; ++i )
// 	{
// 		dwEntityID = pHero[i];
// 		if( VALID(dwEntityID) )
// 		{
// 			pEntity = sSceneMgr.CreateEntity();
// 			if( VALID(pEntity) )
// 			{
// 				pEntity->SetFather(&m_MuitlGroup[eCamp]);

// 				if( pEntity->Init(NULL, dwEntityID, i) )
// 				{
// 					m_MuitlGroup[eCamp].AddEntityHero(i, pEntity, NULL);

// 					if( MIsMaster(dwEntityID) || IsJustMaster(dwEntityID) )
// 					{
// 						m_MuitlGroup[eCamp].SetMasterIndex(i);
// 					}
// 				}
// 				else
// 				{
// 					sSceneMgr.DestroyEntity(pEntity);
// 				}
// 			}
// 		}
// 	}

// 	for( INT32 i = 0 ; i < nRuneNum; ++i )
// 	{
// 		// 加载符文
// 		dwRuneTypeID = pRune[i];
// 		if( VALID(dwRuneTypeID) )
// 		{
// 			m_MuitlGroup[eCamp].AddRune(NULL, i, dwRuneTypeID, 0);
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 加入录像中玩家数据
// //-----------------------------------------------------------------------------
// VOID Scene::AddRecordPlayer(ECamp eCamp, const tagGroupRecord* pRecord)
// {
// 	EntityHero* pEntity = NULL;
// 	const tagHeroRecord* pHeroRecord = NULL;

// 	m_MuitlGroup[eCamp].SetPlayerID(pRecord->n64PlayerID);
// 	m_MuitlGroup[eCamp].SetPlayerLevel(pRecord->nLevel);
// 	m_MuitlGroup[eCamp].SetPlayerScore(pRecord->nPlayerScore);
// 	m_MuitlGroup[eCamp].SetPlayerName(pRecord->szName);
// 	m_MuitlGroup[eCamp].SetDmgModAtt(pRecord->nDmgModAtt);
// 	m_MuitlGroup[eCamp].SetWorldName(pRecord->szWorldName);
// 	m_MuitlGroup[eCamp].SetGuildName(pRecord->szGuildName);
// 	m_MuitlGroup[eCamp].SetProtrait(pRecord->n16HeadProtrait);
// 	m_MuitlGroup[eCamp].SetHeadQuality(pRecord->n8HeadQuality);
// 	m_MuitlGroup[eCamp].SetGuildID(pRecord->n64GuildID);

// 	for( INT32 i = 0 ; i < X_Max_Summon_Num; ++i )
// 	{
// 		pHeroRecord = &(pRecord->stHeroRecord[i]);

// 		if( VALID(pHeroRecord->dwEntityID) )
// 		{
// 			pEntity = sSceneMgr.CreateEntity();
// 			if( VALID(pEntity) )
// 			{
// 				pEntity->SetFather(&m_MuitlGroup[eCamp]);

// 				if( pEntity->InitRecord(pHeroRecord, i) )
// 				{
// 					pEntity->InitHeroRecordDmgModAtt(pRecord,i);

// 					m_MuitlGroup[eCamp].AddRecordEntityHero(i, pEntity, m_bRecord);
// 					if( MIsMaster(pHeroRecord->dwEntityID) || IsJustMaster(pHeroRecord->dwEntityID) )
// 					{
// 						m_MuitlGroup[eCamp].SetMasterIndex(i);
// 					}
// 				}
// 				else
// 				{
// 					sSceneMgr.DestroyEntity(pEntity);
// 				}
// 			}
// 		}
// 	}

// 	for( INT32 i = 0 ; i < X_Rune_Max_Group; ++i )
// 	{
// 		// 加载符文
// 		if( VALID(pRecord->dwRuneID[i]) )
// 		{
// 			m_MuitlGroup[eCamp].AddRune(NULL, i, pRecord->dwRuneID[i], pRecord->n8RuneLevel[i]);
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 发送伤害消息
// //-----------------------------------------------------------------------------
// VOID Scene::SendDamage( tagCalcDamageInfo& dmgInfo )
// {
// 	if( IsOnlyRecord() )
// 		return;

// 	if( EIFT_NULL == dmgInfo.eType )
// 	{
// 		CreateSceneProtoMsg(msg, MS_Miss,);
// 		*msg << (UINT32)dmgInfo.stCaster.nLocation;
// 		*msg << (UINT32)dmgInfo.stTarget.nLocation;
// 		*msg << (UINT32)dmgInfo.dwSpellID;
// 		*msg << (UINT32)dmgInfo.dwProcEx;
// 		AddMsgList(msg);

// 		return;
// 	}

// 	if( EIFT_Damage == dmgInfo.eType )
// 	{
// 		CreateSceneProtoMsg(msg, MS_Damage,);
// 		*msg << (UINT32)dmgInfo.stCaster.nLocation;
// 		*msg << (UINT32)dmgInfo.stTarget.nLocation;
// 		*msg << (UINT32)dmgInfo.dwSpellID;
// 		*msg << (UINT32)dmgInfo.dwProcEx;
// 		*msg << (UINT32)dmgInfo.eSchool;
// 		*msg << (UINT32)dmgInfo.nDamage;
// 		AddMsgList(msg);

// 		return;
// 	}

// 	if( EIFT_Heal == dmgInfo.eType )
// 	{
// 		CreateSceneProtoMsg(msg, MS_Heal,);
// 		*msg << (UINT32)dmgInfo.stCaster.nLocation;
// 		*msg << (UINT32)dmgInfo.stTarget.nLocation;
// 		*msg << (UINT32)dmgInfo.dwSpellID;
// 		*msg << (UINT32)dmgInfo.dwProcEx;
// 		*msg << (UINT32)dmgInfo.eSchool;
// 		*msg << (UINT32)dmgInfo.nDamage;
// 		AddMsgList(msg);

// 		return;
// 	}

// 	if( EIFT_Placate == dmgInfo.eType )
// 	{
// 		CreateSceneProtoMsg(msg, MS_Placate,);
// 		*msg << (UINT32)dmgInfo.stCaster.nLocation;
// 		*msg << (UINT32)dmgInfo.stTarget.nLocation;
// 		*msg << (UINT32)dmgInfo.dwSpellID;
// 		*msg << (UINT32)dmgInfo.dwProcEx;
// 		*msg << (UINT32)dmgInfo.eSchool;
// 		*msg << (UINT32)dmgInfo.nDamage;
// 		AddMsgList(msg);

// 		return;
// 	}

// 	if( EIFT_Enrage == dmgInfo.eType )
// 	{
// 		CreateSceneProtoMsg(msg, MS_Enrage,);
// 		*msg << (UINT32)dmgInfo.stCaster.nLocation;
// 		*msg << (UINT32)dmgInfo.stTarget.nLocation;
// 		*msg << (UINT32)dmgInfo.dwSpellID;
// 		*msg << (UINT32)dmgInfo.dwProcEx;
// 		*msg << (UINT32)dmgInfo.eSchool;
// 		*msg << (UINT32)dmgInfo.nDamage;
// 		AddMsgList(msg);

// 		return;
// 	}

// 	if( EIFT_AverageHP == dmgInfo.eType )
// 	{
// 		CreateSceneProtoMsg(msg, MS_AverageHP,);
// 		*msg << (UINT32)dmgInfo.stCaster.nLocation;
// 		*msg << (UINT32)dmgInfo.stTarget.nLocation;
// 		*msg << (UINT32)dmgInfo.dwSpellID;
// 		*msg << (UINT32)dmgInfo.dwProcEx;
// 		*msg << (UINT32)dmgInfo.eSchool;
// 		*msg << (UINT32)dmgInfo.nDamage;
// 		AddMsgList(msg);

// 		return;
// 	}

// 	if( EIFT_SetHP == dmgInfo.eType )
// 	{
// 		CreateSceneProtoMsg(msg, MS_SetHP,);
// 		*msg << (UINT32)dmgInfo.stCaster.nLocation;
// 		*msg << (UINT32)dmgInfo.stTarget.nLocation;
// 		*msg << (UINT32)dmgInfo.dwSpellID;
// 		*msg << (UINT32)dmgInfo.dwProcEx;
// 		*msg << (UINT32)dmgInfo.eSchool;
// 		*msg << (UINT32)dmgInfo.nDamage;
// 		AddMsgList(msg);

// 		return;
// 	}
// }

// //-----------------------------------------------------------------------------
// // 判断战斗是否结束
// //-----------------------------------------------------------------------------
// VOID Scene::CheckCombatFinish(BOOL bAttackFirst)
// {
// 	// pvp
// 	if( m_pEntry->eType == ESBT_Cup ||
// 		m_pEntry->eType == ESBT_Cross3v3Team ||
// 		m_pEntry->eType == ESBT_KingSea ||
// 		m_pEntry->eType == ESBT_KingRank ||
// 		m_pEntry->eType == ESBT_KingCup ||
// 		m_pEntry->eType == ESBT_GuildBump ||
// 		m_pEntry->eType == ESBT_JustArena ||
// 		m_pEntry->eType == ESBT_SeaTrade ||
// 		m_pEntry->eType == ESBT_Arena ||
// 		m_pEntry->eType == ESBT_GuildWheelWar ||
// 		m_pEntry->eType == ESBT_IslandBattle ||
// 		m_pEntry->eType == ESBT_SeaTradeCrossBattle ||
// 		m_pEntry->eType == ESBT_FriendPK)
// 	{
// 		if ( m_MuitlGroup[ESC_Attack].GetValidNum() > m_MuitlGroup[ESC_Defence].GetValidNum() )
// 		{
// 			m_eWinner = ESC_Attack;
// 		}
// 		else if ( m_MuitlGroup[ESC_Attack].GetValidNum() < m_MuitlGroup[ESC_Defence].GetValidNum() )
// 		{
// 			m_eWinner = ESC_Defence;
// 		}
// 		else
// 		{
// 			m_eWinner = bAttackFirst ? ESC_Attack : ESC_Defence;
// 		}

// 	}

// 	// pve
// 	else
// 	{
// 		if( m_MuitlGroup[ESC_Defence].IsValid() )
// 		{
// 			m_eWinner = ESC_Defence;
// 		}
// 		else
// 		{
// 			m_eWinner = ESC_Attack;
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnCombatFinish()
// {
// 	m_bDelete = TRUE;

// 	Player* pCreator = sPlayerMgr.GetPlayerByGUID(m_n64Creator);

// 	if(m_bRecord)
// 	{
// 		if (m_pEntry->eType == ESBT_Square)
// 		{
// 			AuxSquare::SquareNextRecord stNextRecord;
// 			sAuxSquare.GetNextSquareRecrod(m_n64RecordID, stNextRecord);
// 			if (VALID(stNextRecord.n64RecordID))
// 			{
// 				sRecordMgr.RequestRecord(pCreator, stNextRecord.n64RecordID, FALSE);

// 				CreateProtoMsg(msg, MS_TeamNextRound, );
// 				msg << (INT32)m_eWinner;
// 				msg << stNextRecord.nIndex;
// 				SendSceneMessage(NULL, msg);

// 				return;
// 			}
// 		}

// 		OnRecordFinish(pCreator);
// 		return;

// 	}

// 	switch(m_pEntry->eType)
// 	{
// 	case ESBT_Stage:
// 		//OnStageFinish(pCreator);
// 		break;

// 	case ESBT_Grab:
// 		OnGrabFinish(pCreator);
// 		break;

// 	case ESBT_Arena:
// 		OnArenaFinish(pCreator);
// 		break;

// 	case ESBT_Tower:
// 		OnTowerFinish(pCreator);
// 		break;

// 	case ESBT_Script:
// 		OnScriptFinish();
// 		break;

// 	case ESBT_Mine:
// 		OnMineFinish(pCreator);
// 		break;

// 	case ESBT_SeaTrade:
// 		OnSeaTradeFinish(pCreator);
// 		break;

// 	case ESBT_Devil:
// 		OnDevilOfferFinish(pCreator);
// 		break;

// 	case ESBT_BHonour:
// 		OnBHonourFinish(pCreator);
// 		break;

// 	case ESBT_Castle:
// 		OnCastleBattleFinish(pCreator);
// 		break;

// 	case ESBT_Cup:
// 		OnCupFinish();
// 		break;

// 	case ESBT_GuildBoss:
// 		OnGuildBossFinish(pCreator);
// 		break;

// 	case ESBT_FlyUp:
// 		OnFlyUpFinish(pCreator);
// 		break;

// 	default:
// 		break;
// 	}

// 	// 掉落奖励
// 	while (!m_listLoot.empty())
// 	{
// 		tagLootData& stData = m_listLoot.front();

// 		if (VALID(pCreator))
// 		{
// 			pCreator->GainLoot(stData, ELCID_Stage_Reward);
// 		}

// 		m_listLoot.pop_front();
// 	}
// }

// //-----------------------------------------------------------------------------
// // 队伍战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnTeamCombatFinish()
// {
// 	m_bDelete = TRUE;

// 	tagTeamSceneData* pData = sSceneMgr.GetTeamSceneData(m_dwTeamSerial);
// 	if( !VALID(pData) )
// 	{
// 		return;
// 	}

// 	ECamp eLoser = (m_eWinner == ESC_Attack) ? ESC_Defence : ESC_Attack;

// 	// 重置失败方数据
// 	pData->stGroupRecord[eLoser].Init();
// 	pData->bRecordValid[eLoser] = FALSE;
// 	pData->bRecordValid[m_eWinner] = TRUE;
// 	pData->nBattleIndex[eLoser]++;
// 	INT64* n64LoserMem = (eLoser == ESC_Attack) ? pData->stEnterInfo.n64Attack : pData->stEnterInfo.n64Defence;

// 	if( (pData->nBattleIndex[eLoser] < MAX_TEAM_MEM) && VALID(n64LoserMem[pData->nBattleIndex[eLoser]]) )
// 	{
// 		// 保存获胜方战斗数据
// 		m_MuitlGroup[m_eWinner].Save2DB(&(pData->stGroupRecord[m_eWinner]));

// 		// 记录当前血量
// 		for( INT32 i = 0 ; i < X_Hero_Max_Group; ++i )
// 		{
// 			EntityHero* pEntity = m_MuitlGroup[m_eWinner].GetEntityHero(i);
// 			if( VALID(pEntity) )
// 			{
// 				if( pEntity->IsDead() )
// 				{
// 					pData->stGroupRecord[m_eWinner].stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 				}
// 				else
// 				{
// 					pData->stGroupRecord[m_eWinner].stHeroRecord[i].nAtt[EHA_CurHP] = pEntity->GetAttController().GetAttValue(EHA_CurHP);
// 				}
// 			}
// 			else
// 			{
// 				pData->stGroupRecord[m_eWinner].stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 			}
// 		}

// 		// 同步客户端进入下一轮
// 		CreateProtoMsg(msg, MS_TeamNextRound, );
// 		msg << (INT32)m_eWinner;
// 		msg << (INT32)(pData->nBattleIndex[0] + pData->nBattleIndex[1]);
// 		SendSceneMessage(NULL, msg);

// 		return;
// 	}

// 	// 战斗结束发送奖励
// 	switch(m_pEntry->eType)
// 	{
// 	case ESBT_TeamStage:
// 		OnTeamRoomSingleFinish(pData);
// 		break;

// 	case ESBT_RaidBoss:
// 		OnRaidBossFinish(pData);
// 		break;
// 	case ESBT_3v3Team:
// 		{
// 			On3v3Finish();
// 		}
// 		break;
// 	case ESBT_CoupleTower:
// 		OnCoupleTowerFinish(pData);
// 		break;
// 	case ESBT_Remains:
// 		OnRemainsSingleFinish(pData);
// 		break;

// 	default:
// 		break;
// 	}

// 	sSceneMgr.RemoveTeamSceneData(pData);
// }

// //-----------------------------------------------------------------------------
// // 跨服战场结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnCrossCombatFinish()
// {
// 	m_bDelete = TRUE;

// 	tagCrossSceneData* pData = sSceneMgr.GetCrossTeamData(m_dwCrossSerial);
// 	if( !VALID(pData) )
// 	{
// 		return;
// 	}

// 	ECamp eLoser = (m_eWinner == ESC_Attack) ? ESC_Defence : ESC_Attack;
// 	pData->eLoser = eLoser;

// 	// 进入下一轮
// 	INT32 nLoserIdx = ++pData->nBattleIndex[eLoser];
// 	INT32 nWinnerIdx = pData->nBattleIndex[m_eWinner];

// 	if( (nLoserIdx < MAX_TEAM_MEM) &&
// 		(nWinnerIdx < MAX_TEAM_MEM ) &&
// 		VALID(pData->stTeam[eLoser].stGroupRecord[nLoserIdx].n64PlayerID) )
// 	{
// 		INT64 n64ID= pData->stTeam[m_eWinner].stGroupRecord[nWinnerIdx].n64PlayerID;
// 		// 保存获胜方战斗数据
// 		m_MuitlGroup[m_eWinner].Save2DB( &(pData->stTeam[m_eWinner].stGroupRecord[nWinnerIdx]) );
// 		pData->stTeam[m_eWinner].stGroupRecord[nWinnerIdx].n64PlayerID = n64ID;

// 		// 记录当前血量
// 		for( INT32 i = 0 ; i < X_Hero_Max_Group; ++i )
// 		{
// 			EntityHero* pEntity = m_MuitlGroup[m_eWinner].GetEntityHero(i);
// 			if( VALID(pEntity) )
// 			{
// 				if( pEntity->IsDead() )
// 				{
// 					pData->stTeam[m_eWinner].stGroupRecord[nWinnerIdx].stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 				}
// 				else
// 				{
// 					pData->stTeam[m_eWinner].stGroupRecord[nWinnerIdx].stHeroRecord[i].nAtt[EHA_CurHP] = pEntity->GetAttController().GetAttValue(EHA_CurHP);
// 				}
// 			}
// 			else
// 			{
// 				pData->stTeam[m_eWinner].stGroupRecord[nWinnerIdx].stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 			}
// 		}

// 		if ( m_pEntry->eType == ESBT_Cross3v3Team || m_pEntry->eType == ESBT_TeamStage || m_pEntry->eType == ESBT_RaidBoss || m_pEntry->eType == ESBT_Remains)
// 		{
// 			CreateProtoMsg(msg, MS_TeamNextRound, );
// 			msg << (INT32)m_eWinner;
// 			msg << (INT32)(pData->nBattleIndex[0] + pData->nBattleIndex[1]);
// 			SendSceneMessage(NULL, msg);
// 		}

// 		return;
// 	}

// 	// 战斗结束发送奖励
// 	switch(m_pEntry->eType)
// 	{
// 	case ESBT_TeamStage:
// 		{
// 			OnTeamRoomFinish(pData);
// 		}
// 		break;
// 	case ESBT_Cross3v3Team:
// 		{
// 			OnCross3v3Finish(pData);
// 		}
// 		break;
// 	case ESBT_KingSea:
// 		{
// 			OnKingSeaFinish(pData);
// 		}
// 		break;
// 	case ESBT_KingRank:
// 		{
// 			OnKingRankFinish(pData);
// 		}
// 		break;
// 	case ESBT_KingCup:
// 		{
// 			OnKingCupFinish(pData);
// 		}
// 		break;
// 	case ESBT_JustArena:
// 		{
// 			OnJustArenaFinish(pData);
// 		}
// 		break;

// 	case ESBT_GuildBump:
// 		{
// 			OnGuildBumpFinish(pData);
// 		}
// 		break;

// 	case ESBT_GuildWheelWar:
// 		{
// 			OnGuildWheelFinish(pData);
// 		}
// 		break;

// 	case ESBT_IslandBattle:
// 		{
// 			OnIslandBattleFinish(pData);
// 		}
// 		break;
// 	case ESBT_SeaTradeCrossBattle:
// 		{
// 			OnSeaTradeCrossFinish(pData);
// 		}
// 		break;

// 	case ESBT_MakinoBattle:
// 		{
// 			OnMakinoBattleFinish(pData);
// 		}
// 		break;

// 	case ESBT_Remains:
// 		{
// 			OnRemainsBattleFinish(pData);
// 		}
// 		break;
// 	case ESBT_FriendPK:
// 		{
// 			OnFriendPkFinish(pData);
// 		}
// 		break;

// 	default:
// 		break;
// 	}

// 	sSceneMgr.RemoveCrossTeamData(pData);
// }

// //-----------------------------------------------------------------------------
// // 跨服战场结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnSquareCombatFinish()
// {
// 	m_bDelete = TRUE;

// 	tagSquareSceneData* pData = sSceneMgr.GetSquareSceneData(m_dwSquareSerial);
// 	if (!VALID(pData))
// 	{
// 		return;
// 	}

// 	Player* pAttack = sPlayerMgr.GetPlayerByGUID(pData->n64AttackID);
// 	if (!VALID(pAttack))
// 	{
// 		sSceneMgr.RemoveSquareSceneData(pData);
// 		return;
// 	}

// 	tagSquareBeast* pAttackBeast = pAttack->GetBeastContainer().GetBeast(m_stMisc.square.dwAttackBeastTypeID);
// 	if (!VALID(pAttackBeast))
// 	{
// 		sSceneMgr.RemoveSquareSceneData(pData);
// 		return;
// 	}

// 	EntityHero* pAttackEntity = m_MuitlGroup[ESC_Attack].GetEntityHero(1);
// 	if (!VALID(pAttackEntity))
// 	{
// 		sSceneMgr.RemoveSquareSceneData(pData);
// 		return;
// 	}

// 	ECamp eLoser = (m_eWinner == ESC_Attack) ? ESC_Defence : ESC_Attack;
// 	pData->eLoser = eLoser;

// 	// 记录进攻者血量
// 	pAttackBeast->nCurHP = pAttackEntity->GetAttController().GetAttValue(EHA_CurHP);

// 	// 记录进攻异兽速度
// 	pData->nAttackBeastSpeed = pAttack->GetSquareController().CalcBeastSpeed(pAttackBeast);

// 	// 记录进攻异兽当前精力
// 	pData->stMisc.square.nAttackBeastSpirit = pAttackBeast->nCurSpiritPower;

// 	// 记录中立boss血量
// 	if (VALID(m_stMisc.square.dwEntityGroupID) && m_stMisc.square.bIsBoss)
// 	{
// 		for (INT n = 0; n < X_Max_Summon_Num; ++n)
// 		{
// 			EntityHero* pBossEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(n);
// 			if (!VALID(pBossEntity))
// 				continue;

// 			pData->stDefenceBeast.nEntityCurHP[n] = pBossEntity->GetAttController().GetAttValue(EHA_CurHP);
// 		}
// 	}

// 	// 记录防守者血量
// 	if (!VALID(m_stMisc.square.dwEntityGroupID))
// 	{
// 		EntityHero* pDefenceEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(1);
// 		if (VALID(pDefenceEntity))
// 		{
// 			pData->stDefenceBeast.stBeastRecord[pData->nDefenceIndex].nAtt[EHA_CurHP] = pDefenceEntity->GetAttController().GetAttValue(EHA_CurHP);
// 		}
// 	}

// 	// pve进攻胜利
// 	list<tagLootData> listLoot;
// 	const tagSquareMapMonsterEntry* pMonsterEntry = sResMgr.GetSquareMapMonsterEntry(pData->dwMonsterTypeID);
// 	if (m_eWinner == ESC_Attack && VALID(m_pEntityGroupEntry) && VALID(pMonsterEntry))
// 	{
// 		sResMgr.GenLootData(pMonsterEntry->dwLootID, listLoot);

// 		// 获得奖励
// 		for (auto it = listLoot.begin(); it != listLoot.end(); ++it)
// 		{
// 			tagLootData& stData = *it;
// 			pAttack->GainLoot(stData, ELCID_SquareAttackMonster);
// 		}
// 		INT32 nExp = pMonsterEntry->nExp;
// 		//pAttack->GetSquareController().GetScienceValue(ESST_15, nExp);
// 		for (INT32 i = 0; i < SquareBeast_Passive_Skill_Num; i++)
// 		{
// 			pAttack->GetSquareController().GetScienceValueByID(pAttackBeast->dwPassiveSkillID[i], ESST_34, nExp);
// 		}

// 		pAttack->GetSquareController().GainBeastExpChange(pAttackBeast->dwTypeID, nExp);

// 		evtSquareKillMonster evt;
// 		pAttack->SendEvent(evt);

// 	}

// 	// pvp进攻胜利
// 	if (m_eWinner == ESC_Attack && VALID(m_stMisc.square.n64Defencer))
// 	{
// 		// 找到下一个防守者
// 		pData->nDefenceIndex++;
// 		while (pData->nDefenceIndex < SquareBeast_LineUp_Num)
// 		{
// 			tagGroupRecord* pRecord = sAuxSquare.GetEntityRecord(m_stMisc.square.n64Defencer, pData->nDefenceIndex);
// 			if (VALID(pRecord) && pRecord->GetCurHP() > 0)
// 			{
// 				break;
// 			}

// 			pData->nDefenceIndex++;
// 		}
// 	}

// 	if (VALID(m_stMisc.square.dwEntityGroupID) || pData->nDefenceIndex >= SquareBeast_LineUp_Num || m_eWinner == ESC_Defence)
// 	{
// 		OnSquareBattleFinish(pAttack, pData, listLoot);
// 		sSceneMgr.RemoveSquareSceneData(pData);
// 		return;
// 	}

// 	/*CreateProtoMsg(msg, MS_TeamNextRound, );
// 	msg << (INT32)m_eWinner;
// 	msg << (INT32)(pData->nDefenceIndex);
// 	SendSceneMessage(NULL, msg);*/

// 	//
// }

//-----------------------------------------------------------------------------
// 英雄死亡
//-----------------------------------------------------------------------------
func (s *Scene) OnUnitDead(unit *SceneEntity) {
	if unit == nil {
		return
	}

	// todo 计算英雄掉落
	// if (pHero->GetCamp() == ESC_Defence)
	// {
	// 	list<tagLootData> listLoot;
	// 	sResMgr.GenLootData(pHero->GetEntry()->dwLootID, listLoot, &pHero->GetScene()->GetRandom());

	// 	// 通知客户端
	// 	while (!listLoot.empty())
	// 	{
	// 		tagLootData& stData = listLoot.front();

	// 		if( !IsOnlyRecord() )
	// 		{
	// 			CreateSceneProtoMsg(msg, MS_ItemLoot, );
	// 			*msg << (UINT32)pHero->GetLocation();

	// 			CreateSceneProtoMsg(loot_data, LootData, );
	// 			*loot_data << (INT32)stData.eType;
	// 			*loot_data << (UINT32)stData.dwTypeMisc;
	// 			*loot_data << (INT32)stData.nNum;

	// 			*msg << *loot_data;

	// 			AddMsgList(msg);
	// 		}

	// 		// 整合相同掉落
	// 		if (stData.CanPack())
	// 		{
	// 			list<tagLootData>::iterator it = m_listLoot.begin();
	// 			for (; it != m_listLoot.end(); ++it)
	// 			{
	// 				if (it->eType == stData.eType && it->dwTypeMisc == stData.dwTypeMisc)
	// 				{
	// 					it->nNum += stData.nNum;
	// 					break;
	// 				}
	// 			}

	// 			if (it == m_listLoot.end())
	// 			{
	// 				m_listLoot.push_back(stData);
	// 			}
	// 		}
	// 		else
	// 		{
	// 			m_listLoot.push_back(stData);
	// 		}

	// 		listLoot.pop_front();
	// 	}
	// }
}

//-----------------------------------------------------------------------------
// 观看录像结束
//-----------------------------------------------------------------------------
// VOID Scene::OnRecordFinish(Player* pCreator)
// {
// 	DWORD dwEntityID[X_Max_Summon_Num];
// 	INT nEntityNum = 0;
// 	INT32 nStateFlag[X_Max_Summon_Num];
// 	INT nStateNum = 0;

// 	CreateProtoMsg(msg, MS_RecordFinish, );
// 	msg << (bool)(m_eWinner == ESC_Attack);
// 	msg << (UINT32)sServer.GetWorldID();
// 	msg << string( m_MuitlGroup[ESC_Attack].GetPlayerName() );
// 	msg << m_MuitlGroup[ESC_Attack].GetPlayerLevel();
// 	msg << m_MuitlGroup[ESC_Attack].GetPlayerScore();
// 	memset(dwEntityID, 0xff, sizeof(dwEntityID));
// 	nEntityNum  = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 	InsertArrayProtoMsg(msg, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	DWORD dwLevel[X_Max_Summon_Num];
// 	memset(dwLevel, 0xff, sizeof(dwLevel));
// 	 m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(msg, dwLevel, nEntityNum, INT32);
// 	DWORD dwStar[X_Max_Summon_Num];
// 	memset(dwStar, 0xff, sizeof(dwStar));
// 	 m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(msg, dwStar, nEntityNum, INT32);

// 	memset(nStateFlag, 0, sizeof(nStateFlag));
// 	nStateNum = m_MuitlGroup[ESC_Attack].ExportEntityStateFlag(nStateFlag);
// 	InsertArrayProtoMsg(msg, nStateFlag, nStateNum, INT32);
// 	DWORD dwQuality[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(msg, dwQuality, nEntityNum, INT32);

// 	msg << string( m_MuitlGroup[ESC_Defence].GetPlayerName() );
// 	msg << m_MuitlGroup[ESC_Defence].GetPlayerLevel();
// 	msg << m_MuitlGroup[ESC_Defence].GetPlayerScore();
// 	memset(dwEntityID, 0xff, sizeof(dwEntityID));
// 	nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 	InsertArrayProtoMsg(msg, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	memset(dwLevel, 0xff, sizeof(dwLevel));
// 	m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(msg, dwLevel, nEntityNum, INT32);
// 	memset(dwStar, 0xff, sizeof(dwStar));
// 	m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(msg, dwStar, nEntityNum, INT32);
// 	memset(nStateFlag, 0, sizeof(nStateFlag));
// 	nStateNum = m_MuitlGroup[ESC_Defence].ExportEntityStateFlag(nStateFlag);
// 	InsertArrayProtoMsg(msg, nStateFlag, nStateNum, INT32);
// 	memset(dwQuality, 0xff, sizeof(dwQuality));
// 	m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(msg, dwQuality, nEntityNum, INT32);

// 	SendSceneMessage(NULL, msg);
// }

// //-----------------------------------------------------------------------------
// // 关卡战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnStageFinish(Player* pCreator)
// {
// 	/*
// 	if (!VALID(pCreator))
// 		return;

// 	// 通关
// 	pCreator->StagePass(GetStageID(), GetStageGrade());

// 	// 奖励
// 	list<tagLootData> listLoot;
// 	const tagStageEntry* pStageEntry = sStageEntry(GetStageID());
// 	if (VALID(pStageEntry))
// 	{
// 		// 扣除体力
// 		INT32 nThewOrig = pCreator->GetAttController().GetAttValue(EPA_Thew);
// 		pCreator->GetAttController().ModAttValueWithLog(EPA_Thew, -pStageEntry->nThewCost);
// 		INT32 nThewCurr = pCreator->GetAttController().GetAttValue(EPA_Thew);

// 		sLogMgr.LogThew(pCreator, nThewOrig, nThewCurr, ELCID_StageCost, INVALID);

// 		pCreator->PrayFallenLuck();

// 		pCreator->GetCurrencyMgr().Inc(EMT_Gold, ELCID_Stage_Reward, INVALID);
// 		//pCreator->GetCurrencyMgr().Inc(EMT_Soul, pStageEntry->stDetail[GetStageGrade()].nSoul, ELCID_Stage_Reward, INVALID);

// 		// 经验奖励
// 		if (pStageEntry->nThewCost > 0)
// 		{
// 			const tagPlayerLevelUpInfo* pLevelUpEntry = sPlayerLevelUpInfo(pCreator->GetLevel());
// 			if (VALID(pLevelUpEntry))
// 			{
// 				for(INT n = 0; n < X_Max_Summon_Num; ++n)
// 				{
// 					HeroData* pHeroData = pCreator->GetHeroContainer().GetHeroGroup(n);
// 					if(VALID(pHeroData))
// 						pCreator->GainHeroExp(pHeroData, (INT32)(pCreator->GetLevelExpFactor() * (FLOAT)pLevelUpEntry->nPvEExp));
// 				}
// 			}
// 		}

// 		// 随机掉落
// 		sResMgr.GenLootData(pStageEntry->dwLootID, listLoot, &GetRandom());
// 	}

// 	m_listLoot.insert(m_listLoot.end(), listLoot.begin(), listLoot.end());

// 	evtStageBattleEnd evt;
// 	evt.n64Winner = (m_eWinner == ESC_Attack) ? pCreator->GetID() : INVALID;
// 	evt.eMode = pStageEntry->eMode;
// 	evt.dwStageID = pStageEntry->dwID;
// 	pCreator->SendEvent(evt);

// 	// 通知客户端
// 	CreateProtoMsg(msg, MS_StageFinish, );
// 	msg << (bool)(m_eWinner == ESC_Attack);

// 	msg << (INT32)listLoot.size();
// 	while (!listLoot.empty())
// 	{
// 		tagLootData& stData = listLoot.front();

// 		CreateProtoMsg(loot, LootData, );
// 		loot << (INT32)stData.eType;
// 		loot << (UINT32)stData.dwTypeMisc;
// 		loot << (INT32)stData.nNum;

// 		msg << loot;

// 		listLoot.pop_front();
// 	}

// 	pCreator->SendMessage(msg);

// 	sLogMgr.LogFunctionPlaying(pCreator->GetID(), ELCID_Stage_Reward, 2, GetStageID(), GetStageGrade());
// 	*/
// }

// //-----------------------------------------------------------------------------
// // 夺宝战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnGrabFinish(Player* pCreator)
// {
// 	if (!VALID(pCreator))
// 		return;

// 	// 扣除耐力
// 	pCreator->GetAttController().ModAttValueWithLog(EPA_Stamina, -2, ELCID_Amulet_Grab);

// 	// 固定奖励
// 	const tagPlayerLevelUpInfo* pLevelUpInfo = sPlayerLevelUpInfo(pCreator->GetHeroContainer().GetMaster()->GetLevel());

// 	DWORD dwGrabTargetID = pCreator->GetGrabController().GetTargetID();
// 	ASSERT(VALID(dwGrabTargetID));

// 	INT64 n64PreyID = pCreator->GetGrabController().GetPreyID();

// 	BOOL bGrabSuccess = FALSE;

// 	if (m_eWinner == ESC_Attack)
// 	{
// 		DWORD dwScriptResult = INVALID;
// 		if( VALID(sScriptMgr.GetScriptPlayer()) )
// 		{
// 			dwScriptResult = sScriptMgr.GetScriptPlayer()->OnAmuletGrab(pCreator, n64PreyID, dwGrabTargetID);
// 		}

// 		// 玩家第一次必得碎片
// 		if (VALID(dwScriptResult))
// 		{
// 			bGrabSuccess = (BOOL)dwScriptResult;
// 		}
// 		else
// 		{
// 			const tagAmuletGrabRule* pGrabRule = sAmuletGrabRule(dwGrabTargetID);
// 			if (VALID(pGrabRule))
// 			{
// 				bGrabSuccess = Rand(0, 10000) <= (m_stMisc.grab.bPVP ? pGrabRule->n16PropPlayer : pGrabRule->n16PropRobot);

// 				// 对非玩家有抢夺最少次数限制逻辑
// 				if( !VALID(n64PreyID))
// 				{
// 					const tagAmuletChipEntry* pChipEntry = sAmuletChipEntry(dwGrabTargetID);
// 					if (VALID(pChipEntry))
// 					{
// 						const tagItemEntry* pItemEntry = sItemEntry(pChipEntry->dwAmuletID);

// 						if (VALID(pItemEntry))
// 						{
// 							INT8 nTimesOld = pCreator->GetAmuletGrabTimes(dwGrabTargetID);
// 							bGrabSuccess = ( nTimesOld > X_Amulet_Grab_LimitMinNum[(EItemQuality)pItemEntry->byQuality]) ? bGrabSuccess : FALSE;
// 						}
// 					}
// 				}
// 			}
// 		}

// 		// 等级奖励
// 		if (VALID(pLevelUpInfo))
// 		{
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, pLevelUpInfo->nPvPGold, ELCID_Amulet_Grab, INVALID);
// 			pCreator->GetCurrencyMgr().Inc(EMT_Honour, pLevelUpInfo->nPvPHonour, ELCID_Amulet_Grab, INVALID);

// 			for(INT n = 0; n < X_Max_Summon_Num; ++n)
// 			{
// 				HeroData* pHeroData = pCreator->GetHeroContainer().GetHeroGroup(n);
// 				if(VALID(pHeroData))
// 					pCreator->GainHeroExp(pHeroData, (INT32)(pCreator->GetLevelExpFactor() * (FLOAT)pLevelUpInfo->nPvPExp));
// 			}
// 		}

// 		// 掉落修正
// 		if (bGrabSuccess)
// 		{
// 			// PVP掉落
// 			if (VALID(n64PreyID))
// 			{
// 				Player* pPrey = sPlayerMgr.GetPlayerByGUID(n64PreyID);
// 				if (!VALID(pPrey) || pPrey->GetAmuletChipMgr().GetCur(dwGrabTargetID) < 1
// 					|| !sAvyMgr.IsAvyStart(X_Amulet_Grab_ActivityID)
// 					|| pPrey->GetGrabController().IsProtecting()
// 					|| pCreator->GetAmuletChipMgr().GetCur(dwGrabTargetID) > 0)
// 				{
// 					bGrabSuccess = FALSE;
// 				}
// 				else
// 				{
// 					pPrey->GetAmuletChipMgr().Dec(dwGrabTargetID, 1, ELCID_Amulet_Grab, INVALID);
// 					CreateProtoMsg(msgNotice, MS_BeGrabedNotice, );
// 					msgNotice << pCreator->GetID();
// 					msgNotice << (UINT32)dwGrabTargetID;
// 					pPrey->SendMessage(msgNotice);
// 				}
// 			}
// 			// PVE掉落
// 			else
// 			{
// 				// 紫色碎片
// 				const tagAmuletChipEntry* pChipEntry = sAmuletChipEntry(dwGrabTargetID);
// 				if (VALID(pChipEntry))
// 				{
// 					const tagItemEntry* pItemEntry = sItemEntry(pChipEntry->dwAmuletID);
// 					if (VALID(pItemEntry) && pItemEntry->byQuality >= EIQ_Purple)
// 					{
// 						pCreator->GetGrabController().SetPurpleGrabTimes(0);
// 					}
// 				}
// 			}
// 		}
// 		else
// 		{
// 			// PVE掉落
// 			if (!VALID(n64PreyID))
// 			{
// 				// 紫色碎片
// 				const tagAmuletChipEntry* pChipEntry = sAmuletChipEntry(dwGrabTargetID);
// 				if (VALID(pChipEntry))
// 				{
// 					const tagItemEntry* pItemEntry = sItemEntry(pChipEntry->dwAmuletID);
// 					if (VALID(pItemEntry) && pItemEntry->byQuality >= EIQ_Purple)
// 					{
// 						INT nGrabTimes = pCreator->GetGrabController().GetPurpleGrabTimes();
// 						if (pCreator->GetGrabController().GetPurpleGrabTimes() >= X_Amulet_Grab_LeastNum - 1)
// 						{
// 							pCreator->GetGrabController().SetPurpleGrabTimes(0);
// 							bGrabSuccess = TRUE;
// 						}
// 						else
// 						{
// 							pCreator->GetGrabController().SetPurpleGrabTimes(nGrabTimes + 1);
// 						}
// 					}
// 				}
// 			}
// 		}

// 		// 掉落碎片
// 		if (bGrabSuccess)
// 		{
// 			pCreator->GetAmuletChipMgr().Inc(dwGrabTargetID, 1, ELCID_Amulet_Grab, INVALID);
// 			pCreator->GetAmuletChipMgr().SetCanGrabTime(dwGrabTargetID, IncTime(UCLOCK->CurrentClock(), X_Amulet_Grab_ProtectPeriod) );
// 			pCreator->SetAmuletGrabTimes(dwGrabTargetID, 0);
// 		}
// 		else
// 		{
// 			INT8 nTimesOld = pCreator->GetAmuletGrabTimes(dwGrabTargetID);
// 			pCreator->SetAmuletGrabTimes(dwGrabTargetID, nTimesOld + 1);
// 		}
// 	}
// 	else
// 	{
// 		if (VALID(pLevelUpInfo))
// 		{
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, pLevelUpInfo->nPvPGold * 0.5, ELCID_Amulet_Grab, INVALID);
// 			pCreator->GetCurrencyMgr().Inc(EMT_Honour, pLevelUpInfo->nPvPHonour * 0.5, ELCID_Amulet_Grab, INVALID);

// 			for(INT n = 0; n < X_Max_Summon_Num; ++n)
// 			{
// 				HeroData* pHeroData = pCreator->GetHeroContainer().GetHeroGroup(n);
// 				if(VALID(pHeroData))
// 					pCreator->GainHeroExp(pHeroData, (INT32)(pCreator->GetLevelExpFactor() * (FLOAT)pLevelUpInfo->nPvPExp * 0.5));
// 			}
// 		}
// 	}

// 	// 保存录像
// 	Save2DB(TRUE, bGrabSuccess);

// 	evtGrabBattleEnd evt;
// 	evt.n64Winner = (m_eWinner == ESC_Attack) ? pCreator->GetID() : n64PreyID;
// 	evt.n64Attacker = pCreator->GetID();
// 	pCreator->SendEvent(evt);

// 	CreateProtoMsg(msg, MS_GrabFinish, );

// 	CreateProtoMsg(result, PVPResult, );
// 	result << (bool)(m_eWinner == ESC_Attack);
// 	result << (UINT32)m_dwServerID;
// 	result << m_n64RecordID;

// 	DWORD dwEntityID[X_Max_Summon_Num];
// 	INT nEntityNum = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 	result << string(pCreator->GetPlayerInfo()->szPlayerName);
// 	result << (INT32)pCreator->GetLevel();
// 	result << (INT32)pCreator->GetPlayerScore();
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	DWORD dwLevel[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	DWORD dwStar[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	DWORD dwQuality[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);
// 	if (VALID(n64PreyID))
// 	{
// 		const tagPlayerInfo* pPreyInfo = sPlayerMgr.GetPlayerInfoByGUID(n64PreyID);
// 		result << (VALID(pPreyInfo) ? string(pPreyInfo->szPlayerName) : string());
// 		result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nLevel) : INT32(1));
// 		result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nPlayerScore) : INT32(0));
// 		nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 		InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 		InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 		InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 		InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);
// 	}
// 	else
// 	{
// 		result << string();
// 		result << (INT32)1;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 	}

// 	msg << result;

// 	msg << (bool)!!bGrabSuccess;
// 	msg << (UINT32)dwGrabTargetID;

// 	pCreator->SendMessage(msg);
// }

// //-----------------------------------------------------------------------------
// // 竞技场战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnArenaFinish(Player* pCreator)
// {
// 	if (!VALID(pCreator))
// 		return;

// 	// gm命令战斗
// 	if( m_stMisc.arena.bTest )
// 	{
// 		CreateProtoMsg(msg, MS_ArenaFinish, );

// 		CreateProtoMsg(result, PVPResult, );
// 		result << (bool)(m_eWinner == ESC_Attack);
// 		result << (UINT32)m_dwServerID;
// 		result << m_n64RecordID;

// 		DWORD dwEntityID[X_Max_Summon_Num];
// 		INT nEntityNum = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 		result << string(pCreator->GetPlayerInfo()->szPlayerName);
// 		result << (INT32)pCreator->GetLevel();
// 		result << (INT32)pCreator->GetPlayerScore();
// 		InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 		//添加英雄的等级，星级，战斗结束界面新需求
// 		DWORD dwLevel[X_Max_Summon_Num];
// 		m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 		InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 		DWORD dwStar[X_Max_Summon_Num];
// 		m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 		InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 		DWORD dwQuality[X_Max_Summon_Num];
// 		m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 		InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);
// 		if (m_stMisc.arena.bPVP)
// 		{
// 			const tagPlayerInfo* pPreyInfo = sPlayerMgr.GetPlayerInfoByGUID(m_stMisc.arena.n64DstPlayerID);
// 			result << (VALID(pPreyInfo) ? string(pPreyInfo->szPlayerName) : string());
// 			result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nLevel) : INT32(1));
// 			result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nPlayerScore) : INT32(0));
// 			nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 			InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 			m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 			InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 			m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 			InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 			m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 			InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);
// 		}
// 		else
// 		{
// 			result << string();
// 			result << (INT32)1;
// 			result << (INT32)0;
// 			result << (INT32)0;
// 			result << (INT32)0;
// 			result << (INT32)0;
// 			result << (INT32)0;
// 		}

// 		msg << result;

// 		msg << (INT32)0;
// 		msg << (INT32)0;
// 		msg << (INT32)0;
// 		pCreator->SendMessage(msg);
// 		return;
// 	}

// 	const tagPlayerLevelUpInfo* pLevelUpInfo = sPlayerLevelUpInfo(pCreator->GetHeroContainer().GetMaster()->GetLevel());

// 	const tagArenaPlayer* pArenaAttacker = sArenaMgr.GetArenaPlayer(pCreator->GetID());
// 	const tagArenaPlayer* pArenaDefencer = sArenaMgr.GetArenaPlayer(m_stMisc.arena.n64DstPlayerID);
// 	ASSERT(VALID(pArenaAttacker) && VALID(pArenaDefencer));

// 	INT nOldRank = pArenaAttacker->stData.nMaxRank;
// 	INT nRankReward = 0;

// 	// 游戏事件
// 	EVENT_CREATE(pEvent, FinishArena);
// 	pEvent->n64Attacker	= m_n64Creator;
// 	pEvent->n64Defencer	= m_stMisc.arena.n64DstPlayerID;

// 	// 固定奖励
// 	if (m_eWinner == ESC_Attack)
// 	{
// 		// 扣除耐力
// 		pCreator->GetAttController().ModAttValueWithLog(EPA_Stamina, -2, ELCID_Arena_Finish);
// 	   //固定奖励
// 		pCreator->GetCurrencyMgr().Inc(EMT_Honour, 50, ELCID_Arena_Finish, INVALID);
// 		if (VALID(pLevelUpInfo))
// 		{
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, pLevelUpInfo->nPvPGold, ELCID_Arena_Finish, INVALID);
// 			for(INT n = 0; n < X_Max_Summon_Num; ++n)
// 			{
// 				HeroData* pHeroData = pCreator->GetHeroContainer().GetHeroGroup(n);
// 				if(VALID(pHeroData))
// 					pCreator->GainHeroExp(pHeroData, (INT32)(pCreator->GetLevelExpFactor() * (FLOAT)pLevelUpInfo->nPvPExp));
// 			}
// 		}

// 		if (pArenaAttacker->stData.nCurRank > pArenaDefencer->stData.nCurRank)
// 		{
// 			pEvent->n8Result = 1;
// 			nRankReward = sArenaMgr.Swap(pEvent->n64Attacker, pEvent->n64Defencer);
// 		}
// 		else
// 		{
// 			pEvent->n8Result = INVALID;
// 		}
// 		pCreator->OnActivityAvy(EAAT_Arena, 1);
// 		pCreator->OnScoreAvy(ESAT_Arena, 1);
// 	}
// 	else
// 	{
// 		//增加cd
// 		tagDateTime dwCDTime = IncTime(UCLOCK->CurrentClock(), 120);
// 		pCreator->SetArenaCooldown(dwCDTime);
// 		CreateProtoMsg(cd, MS_StartArenaBattle,);
// 		pCreator->SendMessage(cd);

// 		pEvent->n8Result = 0;
// 	}

// 	EVENT_DISPATCH(pEvent);

// 	evtArenaBattleEnd evt;
// 	evt.n64Winner = (m_eWinner == ESC_Attack) ? pEvent->n64Attacker : pEvent->n64Defencer;
// 	evt.n64Attacker = pEvent->n64Attacker;
// 	pCreator->SendEvent(evt);

// 	CreateProtoMsg(msg, MS_ArenaFinish, );

// 	CreateProtoMsg(result, PVPResult, );
// 	result << (bool)(m_eWinner == ESC_Attack);
// 	result << (UINT32)m_dwServerID;
// 	result << m_n64RecordID;

// 	DWORD dwEntityID[X_Max_Summon_Num];
// 	INT nEntityNum = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 	result << string(pCreator->GetPlayerInfo()->szPlayerName);
// 	result << (INT32)pCreator->GetLevel();
// 	result << (INT32)pCreator->GetPlayerScore();
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	DWORD dwLevel[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	DWORD dwStar[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	DWORD dwQuality[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);
// 	if (m_stMisc.arena.bPVP)
// 	{
// 		const tagPlayerInfo* pPreyInfo = sPlayerMgr.GetPlayerInfoByGUID(m_stMisc.arena.n64DstPlayerID);
// 		result << (VALID(pPreyInfo) ? string(pPreyInfo->szPlayerName) : string());
// 		result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nLevel) : INT32(1));
// 		result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nPlayerScore) : INT32(0));
// 		nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 		InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 		InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 		InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 		InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);
// 	}
// 	else
// 	{
// 		result << string();
// 		result << (INT32)1;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 	}

// 	msg << result;

// 	msg << (INT32)nOldRank;
// 	msg << (INT32)pArenaAttacker->stData.nMaxRank;
// 	msg << (INT32)nRankReward;
// 	pCreator->SendMessage(msg);

// }
// //-----------------------------------------------------------------------------
// // 英雄飞升战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnFlyUpFinish(Player* pCreator)
// {
// 	if (!VALID(pCreator))
// 		return;

// 	CreateProtoMsg(msg, MS_FlyUpFinish, );
// 	// 固定奖励
// 	if (m_eWinner == ESC_Attack)
// 	{
// 		pCreator->HeroFlyUpBattleEnd(m_stMisc.flyup.dwHeroSerial,m_stMisc.flyup.IsGod);
// 	}
// 	else
// 	{
// 	}
// 	msg << (bool)(m_eWinner == ESC_Attack);
// 	msg << (UINT32)m_stMisc.flyup.dwHeroSerial;
// 	HeroData* pHero = pCreator->GetHeroContainer().GetHero(m_stMisc.flyup.dwHeroSerial);
// 	if(VALID(pHero))
// 	{
// 		msg << (INT32)pHero->GetFlyUp();
// 	}
// 	else
// 	{
// 		msg << 0;
// 	}

// 	pCreator->SendMessage(msg);
// }
// //-----------------------------------------------------------------------------
// // 塔战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnTowerFinish(Player* pCreator)
// {
// 	if (!VALID(pCreator))
// 		return;

// 	list<tagLootData> listNormal;
// 	list<tagLootData> listRandom;
// 	list<tagLootData> listBoss;

// 	if (m_eWinner == ESC_Attack)
// 	{
// 		const tagTowerEntry* pTowerEntry = sResMgr.GetTowerEntry(m_stMisc.tower.dwTowerID);
// 		if(!VALID(pTowerEntry))
// 			return;

// 		pCreator->GetTowerController().FloorPass((INT32)m_stMisc.tower.dwTowerID);

// 		// 固定掉落
// 		if(VALID(pTowerEntry->dwNormatlLootID))
// 		{
// 			list<tagLootData> listLoot;
// 			sResMgr.GenLootData(pTowerEntry->dwNormatlLootID, listLoot);

// 			// 获得奖励
// 			while (!listLoot.empty())
// 			{
// 				tagLootData& stData = listLoot.front();
// 				pCreator->GainLoot(stData, ELCID_Tower_Reward);
// 				listNormal.push_back(stData);
// 				listLoot.pop_front();
// 			}
// 		}

// 		// boss掉落
// 		if(VALID(pTowerEntry->dwBossLootID))
// 		{
// 			list<tagLootData> listLoot;
// 			sResMgr.GenLootData(pTowerEntry->dwBossLootID, listLoot);

// 			// 获得奖励
// 			while (!listLoot.empty())
// 			{
// 				tagLootData& stData = listLoot.front();
// 				pCreator->GainLoot(stData, ELCID_Stage_Reward);
// 				listRandom.push_back(stData);
// 				listLoot.pop_front();
// 			}
// 		}

// 		// 称号掉落
// 		if(VALID(pTowerEntry->dwTitleLootID))
// 		{
// 			list<tagLootData> listLoot;
// 			sResMgr.GenLootData(pTowerEntry->dwTitleLootID, listLoot);

// 			// 获得奖励
// 			while (!listLoot.empty())
// 			{
// 				tagLootData& stData = listLoot.front();
// 				pCreator->GainLoot(stData, ELCID_Stage_Reward);
// 				listRandom.push_back(stData);
// 				listLoot.pop_front();
// 			}
// 		}

// 		sLogMgr.LogFunctionPlaying(pCreator->GetID(), ELCID_Tower_Reward, 2, m_stMisc.tower.dwTowerID);
// 	}

// 	evtTowerBattleEnd evt;
// 	evt.nFloor = m_stMisc.tower.dwTowerID;
// 	pCreator->SendEvent(evt);

// 	CreateProtoMsg(msg, MS_TowerFinish, );
// 	msg << (bool)(m_eWinner == ESC_Attack);
// 	msg << (UINT32)m_stMisc.tower.dwTowerID;
// 	msg << 1;		// todo达成条件

// 	// 普通奖励
// 	msg << (INT32)listNormal.size();
// 	for(list<tagLootData>::const_iterator cit = listNormal.begin(); cit != listNormal.end(); ++cit)
// 	{
// 		CreateProtoMsg(loot_normal, LootData, );
// 		loot_normal << (INT32)cit->eType;
// 		loot_normal << (UINT32)cit->dwTypeMisc;
// 		loot_normal << (INT32)cit->nNum;
// 		msg << loot_normal;
// 	}

// 	// 随机奖励
// 	msg << (INT32)listRandom.size();
// 	for(list<tagLootData>::const_iterator cit = listRandom.begin(); cit != listRandom.end(); ++cit)
// 	{
// 		CreateProtoMsg(loot_random, LootData, );
// 		loot_random << (INT32)cit->eType;
// 		loot_random << (UINT32)cit->dwTypeMisc;
// 		loot_random << (INT32)cit->nNum;
// 		msg << loot_random;
// 	}

// 	pCreator->SendMessage(msg);
// }

// //-----------------------------------------------------------------------------
// // 战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnScriptFinish()
// {
// 	if( VALID(m_pScript) )
// 	{
// 		DOUBLE nAttackDmg = m_MuitlGroup[ESC_Attack].GetTotalDmg();
// 		DOUBLE nDefenceDmg = m_MuitlGroup[ESC_Defence].GetTotalDmg();
// 		m_pScript->OnCombatFinish(m_n64Creator, nAttackDmg, nDefenceDmg, (INT32)m_eWinner, m_stMisc);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 设置场景怪物
// //-----------------------------------------------------------------------------
// VOID Scene::SetGroupEntry(DWORD dwGroupID)
// {
// 	m_pEntityGroupEntry = sEntityGroupEntry(dwGroupID);
// }

// //-----------------------------------------------------------------------------
// // 设置对阵双方
// //-----------------------------------------------------------------------------
// VOID Scene::SetTeamMember(ECamp eCamp, INT64 n64ID)
// {
// 	if( IS_PLAYER(n64ID) )
// 	{
// 		SetPlayerID(eCamp, n64ID);
// 	}
// 	else
// 	{
// 		SetGroupEntry(n64ID);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 场景伤害百分比加成
// //-----------------------------------------------------------------------------
// FLOAT Scene::GetSceneDmgMod()
// {
// 	if( !VALID(m_MuitlGroup[ESC_Attack].GetPlayerID()) )
// 		return 0.0f;

// 	if( !VALID(m_MuitlGroup[ESC_Defence].GetPlayerID()) )
// 		return 0.0f;

// 	if( m_nCurRound <= 5 )
// 		return 0.0f;

// 	return ((FLOAT)m_nCurRound - 5.0f) * 2000.0f;
// }

// //-----------------------------------------------------------------------------
// // 等级压制
// //-----------------------------------------------------------------------------
// FLOAT Scene::GetLevelSuppress(EntityHero* pCaster, EntityHero* pTarget)
// {
// 	if( (GetStateFlag() & ESST_Suppress) == 0 )
// 		return 0.0f;

// 	if( pCaster->GetCamp() == pTarget->GetCamp() )
// 		return 0.0f;

// 	if( pCaster->GetLevel() > pTarget->GetLevel() )
// 	{
// 		INT32 nLevelDiff = pCaster->GetLevel() - pTarget->GetLevel();
// 		nLevelDiff = Min(nLevelDiff, 30);
// 		const tagEntityLevelSuppress* pEntry = sResMgr.GetLevelSuppress(nLevelDiff);
// 		if( !VALID(pEntry) )
// 			return 0.0f;

// 		if( pCaster->GetCamp() == ESC_Attack )
// 		{
// 			return pEntry->fPlayerDmgDoneInc;
// 		}
// 		else
// 		{
// 			return pEntry->fDmgDoneInc;
// 		}
// 	}

// 	if( pCaster->GetLevel() < pTarget->GetLevel() )
// 	{
// 		INT32 nLevelDiff = pTarget->GetLevel() - pCaster->GetLevel();
// 		nLevelDiff = Min(nLevelDiff, 30);
// 		const tagEntityLevelSuppress* pEntry = sResMgr.GetLevelSuppress(nLevelDiff);
// 		if( !VALID(pEntry) )
// 			return 0.0f;

// 		if( pCaster->GetCamp() == ESC_Attack )
// 		{
// 			return -(pEntry->fDmgTakenDec);
// 		}
// 		else
// 		{
// 			return -(pEntry->fPlayerDmgTakenDec);
// 		}
// 	}

// 	return 0.0f;
// }

// //-----------------------------------------------------------------------------
// // 矿产争夺结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnMineFinish( Player* pCreator )
// {
// 	if (!VALID(pCreator))
// 		return;

// 	const tagPlayerLevelUpInfo* pLevelUpInfo = sPlayerLevelUpInfo(pCreator->GetHeroContainer().GetMaster()->GetLevel());

// 	INT nMicaBooty = 0;

// 	// 固定奖励
// 	if (m_eWinner == ESC_Attack)
// 	{
// 		if (VALID(pLevelUpInfo))
// 		{
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, pLevelUpInfo->nPvPGold, ELCID_MineBattleFinish, INVALID);
// 			pCreator->GetCurrencyMgr().Inc(EMT_Honour, pLevelUpInfo->nPvPHonour, ELCID_MineBattleFinish, INVALID);

// 			for(INT n = 0; n < X_Max_Summon_Num; ++n)
// 			{
// 				HeroData* pHeroData = pCreator->GetHeroContainer().GetHeroGroup(n);
// 				if(VALID(pHeroData))
// 					pCreator->GainHeroExp(pHeroData, (INT32)(pCreator->GetLevelExpFactor() * (FLOAT)pLevelUpInfo->nPvPExp));
// 			}
// 		}

// 		// 扣除耐力
// 		pCreator->GetAttController().ModAttValueWithLog(EPA_Stamina, -2, ELCID_MineBattleFinish);

// 		// 占领成功
// 		sMineMgr.OnOccupyMine(m_stMisc.mine.dwMineID, m_n64Creator, m_stMisc.mine.bPVP ? m_stMisc.mine.n64PlayerID : INVALID, &nMicaBooty);

// 		// 游戏事件
// 		evtOccupyMine evt;
// 		evt.dwMineID = m_stMisc.mine.dwMineID;
// 		evt.bPVP = m_stMisc.mine.bPVP;
// 		pCreator->SendEvent(evt);
// 	}
// 	else
// 	{
// 		// 占领失败
// 		sMineMgr.OnAttackFailed(pCreator);
// 	}

// 	CreateProtoMsg(msg, MS_MineFinish, );

// 	CreateProtoMsg(result, PVPResult, );
// 	result << (bool)(m_eWinner == ESC_Attack);
// 	result << (UINT32)m_dwServerID;
// 	result << m_n64RecordID;

// 	DWORD dwEntityID[X_Max_Summon_Num];
// 	INT nEntityNum = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 	result << string(pCreator->GetPlayerInfo()->szPlayerName);
// 	result << (INT32)pCreator->GetLevel();
// 	result << (INT32)pCreator->GetPlayerScore();
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	DWORD dwLevel[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	DWORD dwStar[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	DWORD dwQuality[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);
// 	if (m_stMisc.mine.bPVP)
// 	{
// 		const tagPlayerInfo* pPreyInfo = sPlayerMgr.GetPlayerInfoByGUID(m_stMisc.mine.n64PlayerID);
// 		result << (VALID(pPreyInfo) ? string(pPreyInfo->szPlayerName) : string());
// 		result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nLevel) : INT32(1));
// 		result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nPlayerScore) : INT32(0));
// 		nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 		InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 		InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 		InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 		InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);
// 	}
// 	else
// 	{
// 		result << string();
// 		result << (INT32)1;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 	}

// 	msg << result;

// 	// Booty
// 	msg << nMicaBooty;

// 	pCreator->SendMessage(msg);
// }

// //-------------------------------------------- ---------------------------------
// // 劫镖结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnSeaTradeFinish(Player* pCreator)
// {
// 	if (!VALID(pCreator))
// 		return;

// 	FLOAT fMoneyReward = 0.0f;
// 	FLOAT fHonourReward = 0.0f;
// 	FLOAT fRate = 1.0f;

// 	if ( !VALID(pCreator) ) return;

// 	// 固定奖励
// 	if (m_eWinner == ESC_Attack)
// 	{
// 		if( pCreator->GetGrapSeaTrderNum() < MAX_GRAP_SEA_TRADE_NUM )
// 		{
// 			// 增加被劫次数
// 			tagSeaTrader* pTrader = sSeaTradeMgr.GetSeaTrader(m_stMisc.sea_trade.n64TraderID);
// 			if( VALID(pTrader) )
// 			{
// 				pTrader->n8GrapedNum++;
// 				if( pTrader->n8GrapedNum > MAX_GRAPED_NUM )
// 					pTrader->n8GrapedNum = MAX_GRAPED_NUM;

// 				DataSaveSet(SeaTrader, m_stMisc.sea_trade.n64TraderID, n8GrapedNum, pTrader->n8GrapedNum);
// 			}

// 			// 增加劫镖次数
// 			pCreator->IncGrapSeaTradeNum();

// 			if( sSeaTradeMgr.IsSeaHot() )
// 			{
// 				fRate = 1.5f;
// 			}

// 			const tagSeaTradeEntry* pEntry = sSeaTradeEntry(m_stMisc.sea_trade.n16Level);
// 			fMoneyReward = pEntry->fMoneyReward * TRADE_REWARD_RATE[m_stMisc.sea_trade.n16GoodsType] * 0.2f * fRate;
// 			fHonourReward = (FLOAT)pEntry->fHonourReward * TRADE_REWARD_RATE[m_stMisc.sea_trade.n16GoodsType] * 0.2f * fRate;
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, fMoneyReward, ELCID_Grap_Sea_Trade, INVALID);
// 			pCreator->GetCurrencyMgr().Inc(EMT_Honour, fHonourReward, ELCID_Grap_Sea_Trade, INVALID);
// 		}
// 	}
// 	else
// 	{
// 		//cd
// 		pCreator->SetGrapCd(IncTime(UCLOCK->CurrentClock(), X_GRAP_CD));
// 	}

// 	evtGrapSeaTrade evt;
// 	evt.n64Winner = (m_eWinner == ESC_Attack) ? pCreator->GetID() : m_MuitlGroup[ESC_Defence].GetPlayerID();
// 	evt.n64Attacker = pCreator->GetID();
// 	pCreator->SendEvent(evt);

// 	CreateProtoMsg(msg, MS_SeaTradeFinish, );
// 	msg << (INT32)fMoneyReward;
// 	msg << (INT32)fHonourReward;

// 	CreateProtoMsg(result, PVPResult, );
// 	result << (bool)(m_eWinner == ESC_Attack);
// 	result << (UINT32)m_dwServerID;
// 	result << m_n64RecordID;

// 	DWORD dwEntityID[X_Max_Summon_Num];
// 	INT nEntityNum = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 	result << string(pCreator->GetPlayerInfo()->szPlayerName);
// 	result << (INT32)pCreator->GetLevel();
// 	result << (INT32)pCreator->GetPlayerScore();
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	DWORD dwLevel[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	DWORD dwStar[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	DWORD dwQuality[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);

// 	const tagPlayerInfo* pPreyInfo = sPlayerMgr.GetPlayerInfoByGUID(m_MuitlGroup[ESC_Defence].GetPlayerID());
// 	result << (VALID(pPreyInfo) ? string(pPreyInfo->szPlayerName) : string());
// 	result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nLevel) : INT32(1));
// 	result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nPlayerScore) : INT32(0));
// 	nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);

// 	msg << result;

// 	pCreator->SendMessage(msg);
// }

// //-----------------------------------------------------------------------------
// // 魔王悬赏结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnDevilOfferFinish( Player* pCreator )
// {
// 	if (!VALID(pCreator))
// 		return;

// 	INT nMasterIndex = m_MuitlGroup[ESC_Defence].GetMasterIndex();
// 	EntityHero* pMaster = m_MuitlGroup[ESC_Defence].GetEntityHero(nMasterIndex);
// 	INT nHP = INVALID, nDamage = 0;
// 	BOOL bWin = (m_eWinner == ESC_Attack);
// 	if (VALID(pMaster))
// 	{
// 		nHP = pMaster->GetAttController().GetAttValue(EHA_CurHP);
// 		nDamage = pMaster->GetTotalDmg() - pMaster->GetTotalHeal();
// 	}
// 	else
// 	{
// 		nDamage = m_MuitlGroup[ESC_Defence].GetTotalDmg() - m_MuitlGroup[ESC_Defence].GetTotalHeal();
// 	}

// 	INT nCurHP = sFriendMgr.OnDevilCombatFinish(pCreator, m_stMisc.devil_offer.n64PlayerID, m_stMisc.devil_offer.nCacheIndex, nDamage, nHP, bWin);
// 	if (nCurHP < 0)
// 		return;

// 	INT nPos = sFriendMgr.GetDevilOfferRankPos(pCreator->GetID(), m_stMisc.devil_offer.n64PlayerID, m_stMisc.devil_offer.nCacheIndex);

// 	CreateProtoMsg(msg, MS_DevilFinish, );

// 	msg << (bool)!!bWin;
// 	msg << (INT32)nCurHP;
// 	msg << (INT32)nDamage;
// 	msg << (INT32)nPos;

// 	pCreator->SendMessage(msg);
// }

// //-----------------------------------------------------------------------------
// // 荣耀之战结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnBHonourFinish(Player* pCreator)
// {
// 	if (!VALID(pCreator))
// 		return;

// 	INT64 n64SrcID = m_stMisc.battle_honour.n64SrcID;
// 	INT64 n64DesID = m_stMisc.battle_honour.n64DesID;
// 	tagBHonourData* pSrcData = sBattleHonour.GetBHonourData(n64SrcID);
// 	tagBHonourData* pDesData = sBattleHonour.GetBHonourData(n64DesID);
// 	if( !VALID(pSrcData) || !VALID(pDesData) )
// 		return;

// 	INT32 nSrcRewardWin = pDesData->nWinNum;
// 	INT32 nDstRewardWin = pSrcData->nWinNum;

// 	const tagBattleHonourEntry* pBHEntry = sResMgr.GetBattleHonourEntry(pCreator->GetLevel());
// 	if( !VALID(pBHEntry) )
// 		return;

// 	INT32 nSrcAddScore = pBHEntry->nBHonourScoreBase;;
// 	INT32 nSrcAddHonour = pBHEntry->nBHonourHonourBase;
// 	INT32 nDstAddScore = pBHEntry->nBHonourScoreBase;
// 	INT32 nDstAddHonour = pBHEntry->nBHonourHonourBase;

// 	if (m_eWinner == ESC_Attack)
// 	{
// 		nDstAddScore = pBHEntry->nBHonourScoreBase / 2;
// 		nDstAddHonour = pBHEntry->nBHonourHonourBase / 2 ;

// 		pSrcData->nWinNum += 1;
// 		pDesData->nWinNum = -1;

// 		nSrcRewardWin = Max(pSrcData->nWinNum, nSrcRewardWin);
// 		nDstRewardWin = 0;

// 		// 跨服玩家数据
// 		if( pDesData->IsCrossWorldPlayer() )
// 		{
// 			// 重置当前血量
// 			for( INT32 i = 0 ; i < X_Hero_Max_Group; ++i )
// 			{
// 				EntityHero* pEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(i);
// 				if( VALID(pEntity) )
// 				{
// 					pDesData->stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = pDesData->stRecord.stHeroRecord[i].nAtt[EHA_MaxHP];
// 				}
// 			}
// 		}

// 			// 记录当前血量
// 		for( INT32 i = 0 ; i < X_Hero_Max_Group; ++i )
// 		{
// 			EntityHero* pEntity = m_MuitlGroup[ESC_Attack].GetEntityHero(i);
// 			if( VALID(pEntity) )
// 			{
// 				if( pEntity->IsDead() )
// 				{
// 					pSrcData->stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 				}
// 				else
// 				{
// 					pSrcData->stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = pEntity->GetAttController().GetAttValue(EHA_CurHP);
// 				}
// 			}
// 			else
// 			{
// 				pSrcData->stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 			}
// 		}
// 	}
// 	else
// 	{
// 		nSrcAddScore = pBHEntry->nBHonourScoreBase / 2;
// 		nSrcAddHonour = pBHEntry->nBHonourHonourBase / 2;

// 		pDesData->nWinNum += 1;
// 		pSrcData->nWinNum = -1;

// 		nDstRewardWin = Max(pDesData->nWinNum, nDstRewardWin);
// 		nSrcRewardWin = 0;

// 		// 记录当前血量
// 		for( INT32 i = 0 ; i < X_Hero_Max_Group; ++i )
// 		{
// 			EntityHero* pEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(i);
// 			if( VALID(pEntity) )
// 			{
// 				if( pEntity->IsDead() )
// 				{
// 					pDesData->stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 				}
// 				else
// 				{
// 					pDesData->stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = pEntity->GetAttController().GetAttValue(EHA_CurHP);
// 				}
// 			}
// 			else
// 			{
// 				pDesData->stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 			}
// 		}
// 	}

// 	pSrcData->nBattleNum += 1;
// 	pDesData->nBattleNum += 1;
// 	INT32 nAttackNum =  m_MuitlGroup[ESC_Attack].GetAttackNum() + m_MuitlGroup[ESC_Defence].GetAttackNum();
// 	sBattleHonour.Reback2WaitList(pSrcData, nAttackNum, n64DesID);
// 	sBattleHonour.Reback2WaitList(pDesData, nAttackNum, n64SrcID);

// 	// 积分奖励
// 	if( nSrcRewardWin > 0 )
// 	{
// 		nSrcAddScore += nSrcRewardWin * pBHEntry->nBHonourScoreRate;
// 		nSrcAddHonour += nSrcRewardWin * pBHEntry->nBHonourHonourRate;
// 	}

// 	if( nDstRewardWin > 0 )
// 	{
// 		nDstAddScore += nDstRewardWin * pBHEntry->nBHonourScoreRate;
// 		nDstAddHonour += nDstRewardWin * pBHEntry->nBHonourHonourRate;
// 	}

// 	if( (pSrcData->nRewardHonour + nSrcAddHonour) >= pBHEntry->nBHonourMaxHonour )
// 	{
// 		nSrcAddHonour = pBHEntry->nBHonourMaxHonour - pSrcData->nRewardHonour;
// 		pSrcData->nRewardHonour = pBHEntry->nBHonourMaxHonour;
// 	}
// 	else
// 	{
// 		pSrcData->nRewardHonour += nSrcAddHonour;
// 	}

// 	if( (pDesData->nRewardHonour + nDstAddHonour) >= pBHEntry->nBHonourMaxHonour )
// 	{
// 		nDstAddHonour = pBHEntry->nBHonourMaxHonour - pDesData->nRewardHonour;
// 		pDesData->nRewardHonour = pBHEntry->nBHonourMaxHonour;
// 	}
// 	else
// 	{
// 		pDesData->nRewardHonour += nDstAddHonour;
// 	}

// 	if( VALID(pSrcData->GetGuildID()) )
// 	{
// 		pSrcData->nScore += nSrcAddScore;
// 		pSrcData->nContribute += pBHEntry->nBHonourContribute;
// 	}

// 	if( VALID(pDesData->GetGuildID()) )
// 	{
// 		pDesData->nScore += nDstAddScore;
// 		pDesData->nContribute += pBHEntry->nBHonourContribute;
// 	}

// 	INT32 nSrcGuildScore = 0;
// 	INT32 nDstGuildScore = 0;

// 	Guild* pSrcGuild = NULL;
// 	Guild* pDstGuild = NULL;

// 	// 荣誉和帮贡奖励
// 	Player* pSrcPlayer = sPlayerMgr.GetPlayerByGUID(n64SrcID);
// 	if( VALID(pSrcPlayer) )
// 	{
// 		pSrcPlayer->ModContribution(pBHEntry->nBHonourContribute, ELCID_Battle_Honour, INVALID);
// 		pSrcPlayer->ModBattleHonourScore(nSrcAddScore);
// 		pSrcPlayer->GetCurrencyMgr().Inc(EMT_Honour, nSrcAddHonour, ELCID_Battle_Honour, INVALID);

// 		sLogMgr.LogFunctionPlaying(pSrcPlayer->GetID(), ELCID_Battle_Honour, 2);
// 		sLogMgr.LogActivity(n64SrcID, X_BHonour_Activity_ID);
// 		evtBattleHonour evt;
// 		pSrcPlayer->SendEvent(evt);

// 		pSrcGuild = sGuildMgr.GetGuild(pSrcPlayer->GetGuildID());
// 		if( VALID(pSrcGuild) )
// 		{
// 			nSrcGuildScore = pSrcGuild->AddBHonourScore(nSrcAddScore);
// 		}
// 	}

// 	Player* pDstPlayer = sPlayerMgr.GetPlayerByGUID(n64DesID);
// 	if( VALID(pDstPlayer) )
// 	{
// 		pDstPlayer->ModContribution(pBHEntry->nBHonourContribute, ELCID_Battle_Honour, INVALID);
// 		pDstPlayer->ModBattleHonourScore(nDstAddScore);
// 		pDstPlayer->GetCurrencyMgr().Inc(EMT_Honour, nDstAddHonour, ELCID_Battle_Honour, INVALID);

// 		sLogMgr.LogFunctionPlaying(pDstPlayer->GetID(), ELCID_Battle_Honour, 2);
// 		sLogMgr.LogActivity(n64DesID, X_BHonour_Activity_ID);
// 		evtBattleHonour evt;
// 		pDstPlayer->SendEvent(evt);

// 		pDstGuild = sGuildMgr.GetGuild(pDstPlayer->GetGuildID());
// 		if( VALID(pDstGuild) )
// 		{
// 			nDstGuildScore = pDstGuild->AddBHonourScore(nDstAddScore);
// 		}
// 	}

// 	// 同步帮会战斗录像
// 	tagRecordInfo* pRecordInfo = sRecordMgr.GetBHonourPlayerRecordInfo(n64SrcID, PACK_RECORD_ID(m_dwID, sServer.GetWorldID()));
// 	if( VALID(pRecordInfo) )
// 	{
// 		CreateProtoMsg(msg, MS_SynGuildBHonourRecordInfo, );
// 		sRecordMgr.FillRecordInfo(msg, pRecordInfo);

// 		if( VALID(pSrcPlayer) )
// 		{
// 			if( VALID(pSrcGuild) )
// 			{
// 				pSrcGuild->SendGuildMessage(msg);
// 			}
// 			else
// 			{
// 				pSrcPlayer->SendMessage(msg);
// 			}
// 		}

// 		if( VALID(pDstPlayer) )
// 		{
// 			if( VALID(pDstGuild) )
// 			{
// 				if( VALID(pSrcGuild) && pSrcGuild->GetID() != pDstGuild->GetID() )
// 				{
// 					pDstGuild->SendGuildMessage(msg);
// 				}
// 				else if( !VALID(pSrcGuild))
// 				{
// 					pDstGuild->SendGuildMessage(msg);
// 				}
// 			}
// 			else
// 			{
// 				pDstPlayer->SendMessage(msg);
// 			}
// 		}
// 	}

// 	// 发送事件
// 	EVENT_CREATE(pSrcEvent, BHonourFinish);
// 	pSrcEvent->n64PlayerID	= pSrcData->GetPlayerID();
// 	pSrcEvent->n64GuildID	= pSrcData->GetGuildID();
// 	pSrcEvent->nWinNum		= pSrcData->nWinNum;
// 	pSrcEvent->nAddScore	= nSrcAddScore;
// 	pSrcEvent->nGuildScore	= nSrcGuildScore;
// 	EVENT_DISPATCH(pSrcEvent);

// 	// 不是跨服玩家数据
// 	if( !pDesData->IsCrossWorldPlayer() )
// 	{
// 		EVENT_CREATE(pDstEvent, BHonourFinish);
// 		pDstEvent->n64PlayerID	= pDesData->GetPlayerID();
// 		pDstEvent->n64GuildID	= pDesData->GetGuildID();
// 		pDstEvent->nWinNum		= pDesData->nWinNum;
// 		pDstEvent->nAddScore	= nDstAddScore;
// 		pDstEvent->nGuildScore	= nDstGuildScore;
// 		EVENT_DISPATCH(pDstEvent);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 城堡争夺战结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnCastleBattleFinish( Player* pCreator )
// {
// 	// 通知Battle
// 	sCastleMgr.OnCombatFinish((m_eWinner == ESC_Attack), m_n64Creator, m_stMisc.castle_battle.n64GuildID, m_stMisc.castle_battle.nCastleIndex, m_stMisc.castle_battle.nGrade, m_n64RecordID);

// 	// 通知客户端
// 	CreateProtoMsg(msg, MS_CastleFinish, );

// 	CreateProtoMsg(result, PVPResult, );
// 	result << (bool)(m_eWinner == ESC_Attack);
// 	result << (UINT32)m_dwServerID;
// 	result << m_n64RecordID;

// 	DWORD dwEntityID[X_Max_Summon_Num];
// 	INT nEntityNum = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 	result << string(m_MuitlGroup[ESC_Attack].GetPlayerName());
// 	result << (INT32)m_MuitlGroup[ESC_Attack].GetPlayerLevel();
// 	result << (INT32)m_MuitlGroup[ESC_Attack].GetPlayerScore();
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	DWORD dwLevel[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	DWORD dwStar[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	DWORD dwQuality[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);

// 	if (!VALID(m_stMisc.castle_battle.dwRobotID))
// 	{
// 		const tagPlayerInfo* pPreyInfo = sPlayerMgr.GetPlayerInfoByGUID(m_MuitlGroup[ESC_Defence].GetPlayerID());
// 		result << (VALID(pPreyInfo) ? string(pPreyInfo->szPlayerName) : string());
// 		result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nLevel) : INT32(1));
// 		result << (VALID(pPreyInfo) ? INT32(pPreyInfo->nPlayerScore) : INT32(0));
// 		nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 		InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 		InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 		InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 		m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 		InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);
// 	}
// 	else
// 	{
// 		result << string();
// 		result << (INT32)1;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 		result << (INT32)0;
// 	}

// 	msg << result;

// 	sBattleSession.SendOutlandMessage(m_dwServerID, m_n64Creator, msg);
// }

// //-----------------------------------------------------------------------------
// // 单人队伍关卡战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnTeamRoomSingleFinish(tagTeamSceneData* pData)
// {
// 	const tagRoomEntry* pEntry= sResMgr.GetRoomEntry(m_stMisc.room.dwID);
// 	if(!VALID(pEntry))
// 		return;

// 	Player* pLeader = NULL;

// 	for( INT32 i = 0; i < MAX_TEAM_MEM; ++i )
// 	{
// 		if( !VALID(pData->stEnterInfo.n64Attack[i]) )
// 			continue;

// 		Player* pPlayer = sPlayerMgr.GetPlayerByGUID(pData->stEnterInfo.n64Attack[i]);
// 		if (!VALID(pPlayer))
// 			continue;

// 		// 双倍消耗
// 		const FlagArray<ETSF_End>& flagTeamStage = pPlayer->GetTeamStageFlag();
// 		INT32 nCostFactor = 1;
// 		if (flagTeamStage.IsSet(ETSF_Double))
// 			nCostFactor = 2;

// 		//倍率
// 		DWORD dRate = INVALID;

// 		// 奖励
// 		list<tagLootData>::iterator it;
// 		list<tagLootData> listLoot;
// 		listLoot.clear();
// 		if( m_eWinner == ESC_Attack )
// 		{
// 			for (INT n = 0; n < nCostFactor; ++n)
// 			{
// 				// 固定掉落
// 				sResMgr.GenLootGroupData(pEntry->dwNormatlLootID, listLoot);
// 			}

// 			// 开启buff增加金币
// 			if (pPlayer->GetTeamStageBuffTime() > UCLOCK->CurrentClock())
// 				listLoot.push_back(tagLootData(ELT_Currency, EMT_Gold, 500000 * nCostFactor));

// 			//组队副本琼玉奖励翻倍
// 			if( VALID(sScriptMgr.GetScriptPlayer()) )
// 			{
// 				dRate = sScriptMgr.GetScriptPlayer()->OnGetAvyRewardRate(m_stMisc.room.dwID,39);
// 				if( VALID(dRate) && dRate > 0)
// 				{
// 					it = listLoot.begin();
// 					for ( it; it != listLoot.end(); ++it )
// 					{
// 						if ((*it).eType == ELT_Currency && (*it).dwTypeMisc == EMT_JuniorDemon)
// 						{
// 							(*it).nNum = (*it).nNum * dRate;
// 						}
// 					}
// 				}
// 			}
// 		}

// 		// 非雇佣玩家
// 		if( !(pData->stEnterInfo.stMisc.room.dwHireMask & (1 << i)) )
// 		{
// 			// 通知客户端
// 			CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 			msg << (bool)(m_eWinner == ESC_Attack);

// 			if( m_eWinner == ESC_Attack )
// 			{
// 				// 掉落奖励
// 				it = listLoot.begin();
// 				while( it != listLoot.end() )
// 				{
// 					pPlayer->GainLoot(*it, ELCID_TeamStage_Reward, m_stMisc.room.dwID);
// 					++it;
// 				}
// 				tagLootData stReward;
// 				if (pPlayer->OnStageLootAvy(ESLAT_TeamStage, stReward))
// 				{
// 					listLoot.push_back(stReward);
// 				}
// 				sLogMgr.LogFunctionPlaying(pPlayer->GetID(), ELCID_TeamStage_Reward, 2, pData->stEnterInfo.stMisc.room.dwID, INVALID);
// 			}
// 			msg << (INT32)listLoot.size();
// 			it = listLoot.begin();
// 			while( it != listLoot.end() )
// 			{
// 				CreateProtoMsg(loot, LootData, );
// 				loot << (INT32)(*it).eType;
// 				loot << (UINT32)(*it).dwTypeMisc;
// 				loot << (INT32)((*it).nNum);

// 				msg << loot;

// 				++it;
// 			}

// 			// 扣除耐力
// 			pPlayer->GetAttController().ModAttValueWithLog(EPA_Stamina, -sConstParam->nRoomStageStaminaSingle * nCostFactor, ELCID_TeamStage_Reward);
// 			sLogMgr.LogActivity(pPlayer->GetID(), X_TeamStage_Activity_ID);
// 			pPlayer->OnActivityAvy(EAAT_TeamStage, 1);
// 			pPlayer->OnScoreAvy(ESAT_TeamStage, 1);
// 			//精彩活动Log
// 			if( VALID(dRate) && dRate > 0)
// 			{
// 				sLogMgr.LogActivity(pPlayer->GetID(), 1039);
// 			}

// 			// 胜利通过该关卡
// 			if( m_eWinner == ESC_Attack )
// 				pPlayer->SetTeamRoomMask(pEntry->nIndex);

// 			if( pPlayer->GetID() == m_stMisc.room.n64LeaderID )
// 				pLeader = pPlayer;

// 			Room* pRoom = sRoomMgr.GetRoom(m_stMisc.room.dwRoomID);
// 			// 队员改为未准备状态
// 			if( pPlayer != pLeader && VALID(pRoom))
// 			{
// 				pRoom->Ready(pPlayer, FALSE);
// 			}

// 			// 发送消息
// 			pPlayer->SendMessage(msg);

// 			if( !(pData->stEnterInfo.stMisc.room.dwHireMask & (1 << i)) )
// 			{
// 				// 发送事件
// 				evtTeamRoom evt;
// 				evt.nWinner = m_eWinner;
// 				evt.nChapterID = m_stMisc.room.dwID;
// 				evt.bSingle = TRUE;
// 				pPlayer->SendEvent(evt);
// 			}
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 单人女娲遗迹
// //-----------------------------------------------------------------------------
// VOID Scene::OnRemainsSingleFinish(tagTeamSceneData* pData)
// {
// 	const tagRoomEntry* pEntry= sResMgr.GetRoomEntry(m_stMisc.remains.dwID);
// 	if(!VALID(pEntry))
// 		return;

// 	const tagRemainsEntry* pRemainsEntry = sResMgr.GetRemainsEntry(m_stMisc.remains.nRemainsFloor);
// 	if(!VALID(pRemainsEntry))
// 		return;

// 	Player* pLeader = NULL;

// 	for( INT32 i = 0; i < MAX_TEAM_MEM; ++i )
// 	{
// 		if( !VALID(pData->stEnterInfo.n64Attack[i]) )
// 			continue;

// 		// 奖励
// 		list<tagLootData>::iterator it;
// 		list<tagLootData> listLoot;
// 		listLoot.clear();
// 		if( m_eWinner == ESC_Attack )
// 		{
// 			// 固定掉落
// 			sResMgr.GenLootGroupData(pEntry->dwNormatlLootID, listLoot);
// 		}

// 		// 通知客户端
// 		CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 		msg << (bool)(m_eWinner == ESC_Attack);

// 		Player* pPlayer = sPlayerMgr.GetPlayerByGUID(pData->stEnterInfo.n64Attack[i]);
// 		if( VALID(pPlayer) )
// 		{
// 			if( m_eWinner == ESC_Attack )
// 			{
// 				if(m_stMisc.remains.nRemainsFloor == (pPlayer->GetPlayerInfo()->nRemainsFloor + 1))
// 				{
// 					sResMgr.GenLootData(pRemainsEntry->dwLootID, listLoot);
// 				}

// 				// 掉落奖励
// 				msg << (INT32)listLoot.size();
// 				it = listLoot.begin();
// 				while( it != listLoot.end() )
// 				{
// 					CreateProtoMsg(loot, LootData, );
// 					loot << (INT32)(*it).eType;
// 					loot << (UINT32)(*it).dwTypeMisc;
// 					loot << (INT32)((*it).nNum);

// 					msg << loot;

// 					++it;
// 				}

// 				it = listLoot.begin();
// 				while( it != listLoot.end() )
// 				{
// 					pPlayer->GainLoot(*it, ELCID_Remains_Reward, m_stMisc.room.dwID);
// 					++it;
// 				}
// 				sLogMgr.LogFunctionPlaying(pPlayer->GetID(), ELCID_Remains_Reward, 2, pData->stEnterInfo.stMisc.room.dwID, INVALID);
// 			}
// 			else
// 			{
// 				msg << 0;
// 			}

// 			// 改为为准备状态
// 			Room* pRoom = sRoomMgr.GetRoom(m_stMisc.remains.dwRoomID);
// 			if( pPlayer != pLeader && VALID(pRoom))
// 			{
// 				pRoom->Ready(pPlayer, FALSE);
// 			}

// 			// 发送消息
// 			pPlayer->SendMessage(msg);
// 		}
// 	}
// }

// VOID Scene::OnTeamRoomFinish(tagCrossSceneData* pData)
// {
// 	const tagRoomEntry* pEntry= sResMgr.GetRoomEntry(m_stMisc.room.dwRoomID);
// 	if(!VALID(pEntry))
// 		return;

// 	TMap<INT64, INT32> mapMarriageMember;
// 	for( INT32 i = 0; i < MAX_TEAM_MEM; ++i )
// 	{
// 		if( VALID(pData->stTeam[ESC_Attack].stGroupRecord[i].n64PlayerID) )
// 		{
// 			Player* pPlayer = sPlayerMgr.GetPlayerByGUID(pData->stTeam[ESC_Attack].stGroupRecord[i].n64PlayerID);
// 			if (VALID(pPlayer) && VALID(pPlayer->GetMarriageID()))
// 				mapMarriageMember.ModifyValue(pPlayer->GetMarriageID(), 1);
// 		}
// 	}

// 	// 获取夫妻奖励值
// 	TMap<INT64, INT32>::Iterator itMapMarriage = mapMarriageMember.Begin();
// 	INT32 nCount = 0;
// 	INT64 nKey = INVALID;
// 	while(mapMarriageMember.PeekNext(itMapMarriage, nKey, nCount))
// 	{
// 		if (nCount == 2)
// 		{
// 			Marriage* pMarriage = sMarriageMgr.GetMarriage(nKey);
// 			if( VALID(pMarriage) )
// 			{
// 				pMarriage->AddConjugalLove(sConstParam->nRoomStageCoupleReward);
// 			}
// 		}
// 	}

// 	for( INT32 i = 0; i < MAX_TEAM_MEM; ++i )
// 	{
// 		if( !VALID(pData->stTeam[ESC_Attack].stGroupRecord[i].n64PlayerID) )
// 			continue;

// 		Player* pPlayer = sPlayerMgr.GetPlayerByGUID(pData->stTeam[ESC_Attack].stGroupRecord[i].n64PlayerID);
// 		if (!VALID(pPlayer))
// 			continue;

// 		//倍率
// 		DWORD dRate = INVALID;

// 		// 双倍消耗
// 		const FlagArray<ETSF_End>& flagTeamStage = pPlayer->GetTeamStageFlag();
// 		INT32 nCostFactor = 1;
// 		if (flagTeamStage.IsSet(ETSF_Double))
// 			nCostFactor = 2;

// 		// 奖励
// 		list<tagLootData>::iterator it;
// 		list<tagLootData> listLoot;
// 		listLoot.clear();

// 		if( m_eWinner == ESC_Attack )
// 		{
// 			// 固定掉落
// 			for (INT n = 0; n < nCostFactor; ++n)
// 			{
// 				sResMgr.GenLootGroupData(pEntry->dwNormatlLootID, listLoot);
// 			}

// 			// 开启buff增加金币
// 			if (pPlayer->GetTeamStageBuffTime() > UCLOCK->CurrentClock())
// 				listLoot.push_back(tagLootData(ELT_Currency, EMT_Gold, 500000 * nCostFactor));

// 			//组队副本琼玉奖励翻倍
// 			if( VALID(sScriptMgr.GetScriptPlayer()) )
// 			{
// 				dRate = sScriptMgr.GetScriptPlayer()->OnGetAvyRewardRate(m_stMisc.room.dwRoomID,39);
// 				if( VALID(dRate) && dRate > 0)
// 				{
// 					it = listLoot.begin();
// 					for ( it; it != listLoot.end(); ++it )
// 					{
// 						if ((*it).eType == ELT_Currency && (*it).dwTypeMisc == EMT_JuniorDemon)
// 						{
// 							(*it).nNum = (*it).nNum * dRate;
// 						}
// 					}
// 				}
// 			}
// 		}

// 		// 非雇佣玩家
// 		if( !(m_stMisc.room.dwHireMask & (1 << i)) )
// 		{
// 			// 通知客户端
// 			CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 			msg << (bool)(m_eWinner == ESC_Attack);

// 			if( m_eWinner == ESC_Attack )
// 			{
// 				// 掉落奖励
// 				it = listLoot.begin();
// 				while( it != listLoot.end() )
// 				{
// 					pPlayer->GainLoot(*it, ELCID_TeamStage_Reward);
// 					++it;
// 				}
// 				tagLootData stReward;
// 				if (pPlayer->OnStageLootAvy(ESLAT_TeamStage, stReward))
// 				{
// 					listLoot.push_back(stReward);
// 				}
// 				sLogMgr.LogFunctionPlaying(pPlayer->GetID(), ELCID_TeamStage_Reward, 2, m_stMisc.room.dwRoomID, INVALID);

// 				// 本服好友增加亲密度
// 				BOOL bFriendAchieve = FALSE;
// 				for(INT32 m = 0; m < MAX_TEAM_MEM; ++m)
// 				{
// 					INT64 n64FriendID = pData->stTeam[ESC_Attack].stGroupRecord[m].n64PlayerID;
// 					if(pPlayer->GetID() == n64FriendID)
// 						continue;

// 					//if(SERVER_INDEX_PLAYER(pPlayer->GetID()) != SERVER_INDEX_PLAYER(n64FriendID))
// 					//	continue;

// 					if(!sFriendMgr.IsFriend(pPlayer->GetID(), n64FriendID))
// 						continue;

// 					sFriendMgr.AddFriendVal(pPlayer, n64FriendID, 1);
// 					bFriendAchieve = TRUE;
// 				}

// 				if(bFriendAchieve)
// 				{
// 					evtTeamRoomFriendPass evt;
// 					pPlayer->SendEvent(evt);
// 				}
// 			}
// 			msg << (INT32)listLoot.size();
// 			it = listLoot.begin();
// 			while( it != listLoot.end() )
// 			{
// 				CreateProtoMsg(loot, LootData, );
// 				loot << (INT32)(*it).eType;
// 				loot << (UINT32)(*it).dwTypeMisc;
// 				loot << (INT32)((*it).nNum);

// 				msg << loot;

// 				++it;
// 			}
// 			//精彩活动Log
// 			if( VALID(dRate) && dRate > 0)
// 			{
// 				sLogMgr.LogActivity(pPlayer->GetID(), 1039);
// 			}

// 			// 扣除耐力
// 			pPlayer->GetAttController().ModAttValueWithLog(EPA_Stamina, -sConstParam->nRoomStageStamina * nCostFactor, ELCID_TeamStage_Reward);

// 			// 设置CD
// 			pPlayer->SetRoomStageCD(IncTime(UCLOCK->CurrentClock(), sConstParam->dwRoomStageCD));

// 			// 发送消息
// 			pPlayer->SendMessage(msg);

// 			if( !(m_stMisc.room.dwHireMask & (1 << i)) )
// 			{
// 				// 发送事件
// 				evtTeamRoom evt;
// 				evt.nWinner = m_eWinner;
// 				evt.nChapterID = m_stMisc.room.dwRoomID;
// 				evt.bSingle = FALSE;
// 				pPlayer->SendEvent(evt);
// 			}
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 组队boss
// //-----------------------------------------------------------------------------
// VOID Scene::OnRaidBossFinish(tagTeamSceneData* pData)
// {
// 	Room* pRoom = sRoomMgr.GetRoom(m_stMisc.raid_boss.dwRoomID);
// 	if(!VALID(pRoom))
// 		return;

// 	const tagRaidBossEntry* pRaidEntry = sResMgr.GetRaidBossEntry(m_stMisc.raid_boss.nBossLevel);
// 	if(!VALID(pRaidEntry))
// 		return;

// 	// 奖励
// 	list<tagLootData>::iterator it;
// 	list<tagLootData> listLoot;

// 	if( m_eWinner == ESC_Attack )
// 	{
// 		// boss数量改变
// 		//sRoomMgr.ModRaidBossNum(-1);
// 	}

// 	CreateProtoMsg(msgNotReward, MS_TeamRoomFinish, );
// 	msgNotReward << (bool)(m_eWinner == ESC_Attack);
// 	msgNotReward << (INT32)0;

// 	Player* pLeader = sPlayerMgr.GetPlayerByGUID(pData->stEnterInfo.stMisc.raid_boss.n64LeaderID);
// 	//
// 	// 判断队伍成员是否来自同一军团
// 	INT64 n64GuildID = INVALID;
// 	BOOL bSameGuild = TRUE;
// 	for (INT j = 0; j < MAX_TEAM_MEM; j++)
// 	{
// 		if( VALID(pData->stEnterInfo.n64Attack[j]) )
// 		{
// 			tagPlayerInfo *pInfo = sPlayerMgr.GetPlayerInfoByGUID(pData->stEnterInfo.n64Attack[j]);
// 			if ( !VALID(pInfo) || !VALID(pInfo->pGuildMem) || !VALID(pInfo->pGuildMem->n64GuildID) )
// 			{
// 				bSameGuild = FALSE;
// 				break;
// 			}

// 			if ( !VALID(n64GuildID) )
// 			{
// 				n64GuildID = pInfo->pGuildMem->n64GuildID;
// 			}
// 			else if ( n64GuildID != pInfo->pGuildMem->n64GuildID )
// 			{
// 				bSameGuild = FALSE;
// 				break;
// 			}
// 		}
// 		else
// 		{
// 			bSameGuild = FALSE;
// 			break;
// 		}
// 	}
// 	//
// 	for( INT32 i = 0; i < MAX_TEAM_MEM; ++i )
// 	{
// 		if( VALID(pData->stEnterInfo.n64Attack[i]) )
// 		{
// 			Player* pPlayer = sPlayerMgr.GetPlayerByGUID(pData->stEnterInfo.n64Attack[i]);
// 			if( VALID(pPlayer) )
// 			{
// 				if( pPlayer->GetID() == m_stMisc.raid_boss.n64LeaderID )
// 					pLeader = pPlayer;

// 				// 队员改为未准备状态
// 				if( pPlayer != pLeader )
// 				{
// 					pRoom->Ready(pPlayer, FALSE);
// 				}

// 				// 不管胜利失败都要重置cd
// 				//pPlayer->SetRaidBossBattleTime(IncTime(UCLOCK->CurrentClock(), sConstParam->nRaidBossCD));

// 				if( m_eWinner == ESC_Attack )
// 				{
// 					if( pPlayer->GetRaidBossValidNum() > 0 )
// 					{
// 						// 固定掉落
// 						listLoot.clear();
// 						sResMgr.GenLootData(pRaidEntry->dwNormatlLootID, listLoot);
// 						//3个同军团，多奖励一个琼玉
// 						if(bSameGuild)
// 						{
// 							tagLootData qiongyu;
// 							qiongyu.eType = ELT_Currency;
// 							qiongyu.nNum = 1;
// 							qiongyu.dwTypeMisc = EMT_JuniorDemon;

// 							listLoot.push_back(qiongyu);
// 						}
// 						// 掉落奖励
// 						it = listLoot.begin();
// 						while( it != listLoot.end() )
// 						{
// 							pPlayer->GainLoot(*it, ELCID_RaidBoss_Reward);
// 							++it;
// 						}

// 						// 发送消息
// 						CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 						msg << (bool)(m_eWinner == ESC_Attack);
// 						msg << (INT32)listLoot.size();
// 						it = listLoot.begin();
// 						while( it != listLoot.end() )
// 						{
// 							CreateProtoMsg(loot, LootData, );
// 							loot << (INT32)(*it).eType;
// 							loot << (UINT32)(*it).dwTypeMisc;
// 							loot << (INT32)(*it).nNum;
// 							msg << loot;
// 							++it;
// 						}
// 						pPlayer->SendMessage(msg);

// 						pPlayer->ModRaidBossValidNum(-1);
// 						sLogMgr.LogFunctionPlaying(pPlayer->GetID(), ELCID_RaidBoss_Reward, 2, pData->stEnterInfo.stMisc.raid_boss.nBossLevel, INVALID);
// 						sLogMgr.LogActivity(pPlayer->GetID(), X_RaidBoss_Activity_ID);
// 					}
// 					else
// 					{
// 						// 发送事件
// 						EVENT_CREATE(pEvent, GoodPerson);
// 						pEvent->n64PlayerID	= pPlayer->GetID();
// 						pEvent->nValueAdd = 1;
// 						EVENT_DISPATCH(pEvent);

// 						// 发送消息
// 						pPlayer->SendMessage(msgNotReward);

// 						// 增加侠义值(好心值)
// 						pPlayer->GetCurrencyMgr().Inc(EMT_Courage, sConstParam->nBossTeamCourage, ELCID_TeamBoss_Courage, INVALID);
// 					}

// 					// 队长没有次数后自动解散队伍
// 					if(pPlayer->GetRaidBossValidNum() <= 0 && VALID(pLeader) && (pPlayer->GetID() == pLeader->GetID()))
// 					{
// 						sRoomMgr.LeaveRoom(pLeader);
// 					}
// 				}

// 				else
// 				{
// 					// 发送消息
// 					CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 					msg << (bool)(m_eWinner == ESC_Attack);
// 					msg << 0;
// 					pPlayer->SendMessage(msg);
// 				}

// 				// 发送事件
// 				evtRaidBoss evt;
// 				evt.bWin = (m_eWinner == ESC_Attack);
// 				evt.nLevel = m_stMisc.raid_boss.nBossLevel;
// 				pPlayer->SendEvent(evt);
// 			}
// 		}
// 	}

// 	pRoom->SyncRoomData();

// 	if( VALID(pLeader) && (pLeader->GetRaidBossValidNum() <= 0) )
// 	{
// 		pRoom->LeaveRoom(pLeader);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 3v3结束
// //-----------------------------------------------------------------------------
// VOID Scene::On3v3Finish()
// {
// 	s3v3Mgr.OnBattleFinish(m_eWinner, m_stMisc.c3v3.n64TeamLightID, m_stMisc.c3v3.n64TeamDarkID);
// }

// //-----------------------------------------------------------------------------
// // 跨服3v3结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnCross3v3Finish(tagCrossSceneData *pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	const fxDescriptor* pDescriptor = fxDescriptorDatabase::Inst()->GetMessageTypeByName("Cross3v3Member");
// 	if ( !VALID(pDescriptor) )	return;

// 	// 填充队员信息
// 	fxMessage msgMem[MAX_TEAM_MEM*2];
// 	INT32 nIdx = 0;
// 	for (INT i = ESC_Attack; i < ESC_End; i++ )
// 	{
// 		for (INT j = 0; j < MAX_TEAM_MEM; j++ )
// 		{
// 			msgMem[nIdx].SetMessage(pDescriptor);
// 			msgMem[nIdx] << string(pData->stTeam[i].stGroupRecord[j].szName);
// 			msgMem[nIdx] << pData->stTeam[i].stGroupRecord[j].nLevel;
// 			msgMem[nIdx] << pData->stTeam[i].stGroupRecord[j].nPlayerScore;
// 			msgMem[nIdx] << (INT32)(pData->stTeam[i].stGroupRecord[j].n16HeadProtrait);
// 			msgMem[nIdx] << (bool)(0 == strncmp(pData->stTeam[i].stGroupRecord[j].szName, pData->stTeam[i].szLeaderName, SHORT_STRING));

// 			++nIdx;
// 		}
// 	}

// 	for (INT i = ESC_Attack; i < ESC_End; i++ )
// 	{
// 		if ( pData->stTeam[i].dwWorldID != sServer.GetWorldID() )
// 			continue;

// 		// 判断队伍成员是否来自同一军团
// 		INT64 n64GuildID = INVALID;
// 		BOOL bSameGuild = TRUE;
// 		for (INT j = 0; j < MAX_TEAM_MEM; j++)
// 		{
// 			tagPlayerInfo *pInfo = sPlayerMgr.GetPlayerInfoByGUID(pData->stTeam[i].stGroupRecord[j].n64PlayerID);
// 			if ( !VALID(pInfo) || !VALID(pInfo->pGuildMem) || !VALID(pInfo->pGuildMem->n64GuildID) )
// 			{
// 				bSameGuild = FALSE;
// 				break;
// 			}

// 			if ( !VALID(n64GuildID) )
// 			{
// 				n64GuildID = pInfo->pGuildMem->n64GuildID;
// 			}
// 			else if ( n64GuildID != pInfo->pGuildMem->n64GuildID )
// 			{
// 				bSameGuild = FALSE;
// 				break;
// 			}
// 		}

// 		BOOL bWin = (INT)m_eWinner == i;

// 		Room* pRoom = sRoomMgr.GetRoom( pData->stTeam[i].n64TeamID);

// 		for (INT j = 0; j < MAX_TEAM_MEM; j++ )
// 		{
// 			Player *pPlayer = sPlayerMgr.GetPlayerByGUID( pData->stTeam[i].stGroupRecord[j].n64PlayerID);
// 			if ( VALID(pPlayer) )
// 			{
// 				const tagPlayerLevelUpInfo* pRewardEntry = sPlayerLevelUpInfo( pPlayer->GetLevel() );
// 				if ( VALID(pRewardEntry) )
// 				{
// 					CreateProtoMsg(msg, MS_Cross3v3BattleEnd, );
// 					msg << (INT32)m_eWinner;
// 					msg << 2;

// 					// 掉落
// 					for (INT m = 0; m < 2; m++ )
// 					{
// 						CreateProtoMsg(loot, LootData, );

// 						if ( pPlayer->GetCro3v3TotalTimes() < X_Cro3v3_JoinTimes )
// 						{
// 							loot << (INT32)ELT_Currency;
// 							loot << (UINT32)pRewardEntry->n3v3Type[m];
// 							loot << (INT32)pRewardEntry->n3v3Num[m][bWin];
// 							msg << loot;

// 							pPlayer->GetCurrencyMgr().Inc(pRewardEntry->n3v3Type[m], pRewardEntry->n3v3Num[m][bWin], ELCID_3v3_Battle, INVALID);
// 						}
// 						else
// 						{
// 							loot << (INT32)INVALID;
// 							loot << (UINT32)INVALID;
// 							loot << (INT32)0;
// 							msg << loot;
// 						}

// 					}

// 					// 队员信息
// 					msg << (MAX_TEAM_MEM*2);
// 					for ( INT k = 0; k < MAX_TEAM_MEM*2; k++ )
// 					{
// 						msg << msgMem[k];
// 					}

// 					pPlayer->SendMessage(msg);
// 				}

// 				pPlayer->ModifyCro3v3Info(bWin, bSameGuild);
// 				sLogMgr.LogActivity(pPlayer->GetID(), X_Cross3v3_Activity_ID);
// 				// 重置队员状态
// 				if ( VALID(pRoom ) )	pRoom->Ready(pPlayer, FALSE);
// 			}
// 		}

// 		if ( VALID(pRoom) )
// 		{
// 			pRoom->SetInCro3v3(FALSE);
// 			pRoom->SyncRoomData();
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 群雄争霸海选赛
// //-----------------------------------------------------------------------------
// VOID Scene::OnKingSeaFinish(tagCrossSceneData *pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	Player *pCreator = sPlayerMgr.GetPlayerByGUID(pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID);
// 	if ( !VALID(pCreator) )
// 	{
// 		ErrLog("king_battle cannot find palyer on finish %lld!\n", pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID);
// 		return;
// 	}

// 	// 扣除耐力
// 	if ( ESC_Attack == m_eWinner )
// 	{
// 		pCreator->GetAttController().ModAttValueWithLog(EPA_Stamina, -X_KingSea_Stamina, ELCID_KingSea_Finish);
// 	}

// 	const tagPlayerLevelUpInfo* pLevelUpInfo = sPlayerLevelUpInfo(pCreator->GetHeroContainer().GetMaster()->GetLevel());
// 	// 固定奖励
// 	if (VALID(pLevelUpInfo))
// 	{
// 		float fRatio = (ESC_Attack == m_eWinner) ? 1.0f : 0.5f;
// 		fRatio = fRatio * X_KingSea_Stamina / 2.0f;
// 		if ( ESC_Attack == m_eWinner )
// 		{
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, pLevelUpInfo->nPvPGold*fRatio, ELCID_KingSea_Finish, INVALID);
// 			pCreator->GetCurrencyMgr().Inc(EMT_Honour, 50, ELCID_KingSea_Finish, INVALID);

// 			for(INT n = 0; n < X_Max_Summon_Num; ++n)
// 			{
// 				HeroData* pHeroData = pCreator->GetHeroContainer().GetHeroGroup(n);
// 				if(VALID(pHeroData))
// 					pCreator->GainHeroExp(pHeroData, (INT32)(pCreator->GetLevelExpFactor() * (FLOAT)pLevelUpInfo->nPvPExp));
// 			}
// 		}

// 	}

// 	// 金子奖励
// 	if ( ESC_Attack == m_eWinner && pCreator->GetKingBattleController().GetDiamond() < X_KingSea_Diamond_Limit )
// 	{
// 		if ( VALID(pLevelUpInfo) )
// 		{
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, pLevelUpInfo->nKingSeaGold, ELCID_KingSea_Finish, INVALID);
// 			pCreator->GetKingBattleController().IncDiamond( 1 );
// 		}
// 	}

// 	pCreator->GetKingBattleController().OnKingSeaFinish(ESC_Attack == m_eWinner);

// 	CreateProtoMsg(msg, MS_KingSeaFinish, );

// 	CreateProtoMsg(result, PVPResult, );
// 	result << (bool)(m_eWinner == ESC_Attack);
// 	result << (UINT32)m_dwServerID;
// 	result << m_n64RecordID;

// 	DWORD dwEntityID[X_Max_Summon_Num];
// 	result << string(m_MuitlGroup[ESC_Attack].GetPlayerName());
// 	result << m_MuitlGroup[ESC_Attack].GetPlayerLevel();
// 	result << m_MuitlGroup[ESC_Attack].GetPlayerScore();
// 	INT32 nEntityNum = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	DWORD dwLevel[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	DWORD dwStar[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	DWORD dwQuality[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);

// 	result << string(m_MuitlGroup[ESC_Defence].GetPlayerName());
// 	result << m_MuitlGroup[ESC_Defence].GetPlayerLevel();
// 	result << m_MuitlGroup[ESC_Defence].GetPlayerScore();
// 	nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);

// 	msg << result;

// 	tagKingSeaInfo &stKingSeaInfo = pCreator->GetKingBattleController().GetKingSeaInfo();
// 	msg << stKingSeaInfo.nWinTimes;
// 	msg << (UINT32)stKingSeaInfo.dwCD;
// 	msg << stKingSeaInfo.nDiamond;
// 	msg << (bool)stKingSeaInfo.bFailed;

// 	pCreator->SendMessage(msg);

// 	sLogMgr.LogFunctionPlaying(pCreator->GetID(), ELCID_KingSea_Finish, 2);
// 	sLogMgr.LogActivity(pCreator->GetID(), X_KingSea_ActivityID);
// }

// //-----------------------------------------------------------------------------
// // // 群雄争霸排位赛
// //-----------------------------------------------------------------------------
// VOID Scene::OnKingRankFinish(tagCrossSceneData *pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	Player *pCreator = sPlayerMgr.GetPlayerByGUID(pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID);
// 	if ( !VALID(pCreator) )
// 	{
// 		return;
// 	}

// 	// 扣除耐力
// 	if ( ESC_Attack == m_eWinner )
// 	{
// 		pCreator->GetAttController().ModAttValueWithLog(EPA_Stamina, -X_KingSea_Stamina, ELCID_KingRank_Finish);
// 	}

// 	const tagPlayerLevelUpInfo* pLevelUpInfo = sPlayerLevelUpInfo(pCreator->GetHeroContainer().GetMaster()->GetLevel());

// 	// 固定奖励
// 	if (VALID(pLevelUpInfo))
// 	{
// 		float fRatio = (ESC_Attack == m_eWinner) ? 1.0f : 0.5f;
// 		fRatio = fRatio * X_KingSea_Stamina / 2.0f;
// 		if ( ESC_Attack == m_eWinner )
// 		{
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, pLevelUpInfo->nPvPGold*fRatio, ELCID_KingRank_Finish, INVALID);
// 			pCreator->GetCurrencyMgr().Inc(EMT_Honour, 50, ELCID_KingRank_Finish, INVALID);

// 			for(INT n = 0; n < X_Max_Summon_Num; ++n)
// 			{
// 				HeroData* pHeroData = pCreator->GetHeroContainer().GetHeroGroup(n);
// 				if(VALID(pHeroData))
// 					pCreator->GainHeroExp(pHeroData, (INT32)(pCreator->GetLevelExpFactor() * (FLOAT)pLevelUpInfo->nPvPExp));
// 			}
// 		}

// 	}

// 	// 金子奖励
// 	if ( ESC_Attack == m_eWinner && pCreator->GetKingBattleController().GetDiamond() < X_KingSea_Diamond_Limit )
// 	{
// 		if ( VALID(pLevelUpInfo) )
// 		{
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, pLevelUpInfo->nKingSeaGold, ELCID_KingSea_Finish, INVALID);
// 			pCreator->GetKingBattleController().IncDiamond( 1 );
// 		}
// 	}

// 	pCreator->GetKingBattleController().OnKingRankFinish(ESC_Attack == m_eWinner, pData->stTeam[ESC_Defence].stGroupRecord[0].n64PlayerID, &pData->stTeam[ESC_Attack].stGroupRecord[0]);

// 	CreateProtoMsg(msg, MS_KingRankFinish, );

// 	CreateProtoMsg(result, PVPResult, );
// 	result << (bool)(m_eWinner == ESC_Attack);
// 	result << (UINT32)m_dwServerID;
// 	result << m_n64RecordID;

// 	DWORD dwEntityID[X_Max_Summon_Num];
// 	result << string(m_MuitlGroup[ESC_Attack].GetPlayerName());
// 	result << m_MuitlGroup[ESC_Attack].GetPlayerLevel();
// 	result << m_MuitlGroup[ESC_Attack].GetPlayerScore();
// 	INT32 nEntityNum = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	DWORD dwLevel[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	DWORD dwStar[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	DWORD dwQuality[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);

// 	result << string(m_MuitlGroup[ESC_Defence].GetPlayerName());
// 	result << m_MuitlGroup[ESC_Defence].GetPlayerLevel();
// 	result << m_MuitlGroup[ESC_Defence].GetPlayerScore();
// 	nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);

// 	msg << result;

// 	tagKingSeaInfo &stKingSea = pCreator->GetKingBattleController().GetKingSeaInfo();
// 	msg << (UINT32)stKingSea.dwCD;
// 	msg << stKingSea.nDiamond;

// 	pCreator->SendMessage(msg);

// 	sLogMgr.LogFunctionPlaying(pCreator->GetID(), ELCID_KingRank_Finish, 2);
// }

// //-----------------------------------------------------------------------------
// // // 跨服王者之战
// //-----------------------------------------------------------------------------
// VOID Scene::OnKingCupFinish(tagCrossSceneData *pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	// 通知battle战斗结果
// 	tagMWB_KingCupCombat stSend;
// 	stSend.nQueueIdx = m_stMisc.kingcup.nQueueIdx;
// 	stSend.nWinner = m_eWinner;
// 	stSend.nCompeterIdx[ESC_Attack] = m_stMisc.kingcup.nCompeterIdx[ESC_Attack];
// 	stSend.nCompeterIdx[ESC_Defence] = m_stMisc.kingcup.nCompeterIdx[ESC_Defence];
// 	stSend.nLoop = m_stMisc.kingcup.nLoop;
// 	tagMWD_SaveSceneRecord *pRecord = sRecordMgr.GetRecordInCache( m_n64RecordID );
// 	if ( VALID(pRecord) )
// 	{
// 		stSend.stRecord = pRecord->stRecord;
// 	}

// 	// 发送奖励
// 	if ( sServer.GetWorldID() == pData->stTeam[ESC_Attack].dwWorldID )
// 	{
// 		sMailMgr.SendKingCupBattleReward( pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID, m_eWinner == ESC_Attack, m_stMisc.kingcup.nLoop);
// 	}

// 	if ( sServer.GetWorldID() == pData->stTeam[ESC_Defence].dwWorldID )
// 	{
// 		sMailMgr.SendKingCupBattleReward( pData->stTeam[ESC_Defence].stGroupRecord[0].n64PlayerID, m_eWinner == ESC_Defence, m_stMisc.kingcup.nLoop);
// 	}

// 	sBattleSession.SendMessage(&stSend, stSend.dwSize);

// }

// //-----------------------------------------------------------------------------
// /// 公平竞技场
// //-----------------------------------------------------------------------------
// VOID Scene::OnJustArenaFinish(tagCrossSceneData *pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	INT32 nAttackNum = m_MuitlGroup[ESC_Attack].GetAttackNum() + m_MuitlGroup[ESC_Defence].GetAttackNum();

// 	// 发送奖励
// 	if ( sServer.GetWorldID() == pData->stTeam[ESC_Attack].dwWorldID )
// 	{
// 		sJustArenaMgr.OnCombatFinish( pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID, m_eWinner == ESC_Attack, m_n64RecordID, pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID, nAttackNum);
// 	}

// 	if ( sServer.GetWorldID() == pData->stTeam[ESC_Defence].dwWorldID )
// 	{
// 		sJustArenaMgr.OnCombatFinish( pData->stTeam[ESC_Defence].stGroupRecord[0].n64PlayerID, m_eWinner == ESC_Defence, m_n64RecordID, pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID, nAttackNum);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 跨服帮会对冲战结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnGuildBumpFinish(tagCrossSceneData* pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	tagGroupRecord* pAttackRecord = &(pData->stTeam[ESC_Attack].stGroupRecord[0]);
// 	tagGroupRecord* pDefenceRecord = &(pData->stTeam[ESC_Defence].stGroupRecord[0]);

// 	// 获胜方当前血量
// 	INT32 nCurHP[X_Hero_Max_Group];
// 	for( INT32 i = 0 ; i < X_Hero_Max_Group; ++i )
// 	{
// 		EntityHero* pEntity = m_MuitlGroup[m_eWinner].GetEntityHero(i);
// 		if( VALID(pEntity) )
// 		{
// 			if( pEntity->IsDead() )
// 			{
// 				nCurHP[i] = 0;
// 			}
// 			else
// 			{
// 				nCurHP[i] = pEntity->GetAttController().GetAttValue(EHA_CurHP);
// 			}
// 		}
// 		else
// 		{
// 			nCurHP[i] = 0;
// 		}
// 	}

// 	// 通知battle战斗结果
// 	if( sServer.GetWorldID() == pData->stTeam[ESC_Attack].dwWorldID )
// 	{
// 		tagMWB_GuildBumpCombat stSend;
// 		stSend.n64GuildID = m_stMisc.bump.nAttackGuildID;
// 		stSend.stRecordInfo.stSrcPlayer.n64PlayerID = pAttackRecord->n64PlayerID;
// 		stSend.stRecordInfo.stSrcPlayer.nLevel = pAttackRecord->nLevel;
// 		stSend.stRecordInfo.stSrcPlayer.nProtrait = pAttackRecord->n16HeadProtrait;
// 		stSend.stRecordInfo.stSrcPlayer.nPlayerScore = pAttackRecord->nPlayerScore;
// 		memcpy(stSend.stRecordInfo.stSrcPlayer.szPlayerName, pAttackRecord->szName, sizeof(pAttackRecord->szName) );
// 		stSend.stRecordInfo.stDstPlayer.n64PlayerID = pDefenceRecord->n64PlayerID;
// 		stSend.stRecordInfo.stDstPlayer.nLevel = pDefenceRecord->nLevel;
// 		stSend.stRecordInfo.stDstPlayer.nProtrait = pDefenceRecord->n16HeadProtrait;
// 		stSend.stRecordInfo.stDstPlayer.nPlayerScore = pDefenceRecord->nPlayerScore;
// 		memcpy(stSend.stRecordInfo.stDstPlayer.szPlayerName, pDefenceRecord->szName, sizeof(pDefenceRecord->szName) );
// 		stSend.stRecordInfo.n64RecordID = m_n64RecordID;
// 		stSend.stRecordInfo.bWin = (m_eWinner == ESC_Attack);
// 		memcpy(stSend.nCurHP, nCurHP, sizeof(stSend.nCurHP) );
// 		stSend.nRoadIndex = m_stMisc.bump.nRoadIndex;
// 		sBattleSession.SendMessage(&stSend, stSend.dwSize);
// 		sLogMgr.LogActivity(pAttackRecord->n64PlayerID, X_BUMP_AVY);
// 	}

// 	if( sServer.GetWorldID() == pData->stTeam[ESC_Defence].dwWorldID )
// 	{
// 		tagMWB_GuildBumpCombat stSend;
// 		stSend.n64GuildID = m_stMisc.bump.nDefenceGuildID;
// 		stSend.stRecordInfo.stSrcPlayer.n64PlayerID = pDefenceRecord->n64PlayerID;
// 		stSend.stRecordInfo.stSrcPlayer.nLevel = pDefenceRecord->nLevel;
// 		stSend.stRecordInfo.stSrcPlayer.nProtrait = pDefenceRecord->n16HeadProtrait;
// 		stSend.stRecordInfo.stDstPlayer.nPlayerScore = pDefenceRecord->nPlayerScore;
// 		memcpy(stSend.stRecordInfo.stSrcPlayer.szPlayerName, pDefenceRecord->szName, sizeof(pDefenceRecord->szName) );
// 		stSend.stRecordInfo.stDstPlayer.n64PlayerID = pAttackRecord->n64PlayerID;
// 		stSend.stRecordInfo.stDstPlayer.nLevel = pAttackRecord->nLevel;
// 		stSend.stRecordInfo.stDstPlayer.nProtrait = pAttackRecord->n16HeadProtrait;
// 		stSend.stRecordInfo.stSrcPlayer.nPlayerScore = pAttackRecord->nPlayerScore;
// 		memcpy(stSend.stRecordInfo.stDstPlayer.szPlayerName, pAttackRecord->szName, sizeof(pAttackRecord->szName) );
// 		stSend.stRecordInfo.n64RecordID = m_n64RecordID;
// 		stSend.stRecordInfo.bWin = (m_eWinner == ESC_Defence);
// 		memcpy(stSend.nCurHP, nCurHP, sizeof(stSend.nCurHP) );
// 		stSend.nRoadIndex = m_stMisc.bump.nRoadIndex;
// 		sBattleSession.SendMessage(&stSend, stSend.dwSize);
// 		sLogMgr.LogActivity(pDefenceRecord->n64PlayerID, X_BUMP_AVY);
// 	}
// }

// //-----------------------------------------------------------------------------
// // 杯赛战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnCupFinish()
// {
// 	sCupMgr.OnCupBattleFinish(m_eWinner, m_stMisc.cup.dwID, m_stMisc.cup.nLiveNum, m_n64RecordID);
// }

// //-----------------------------------------------------------------------------
// // 帮会Boss战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnGuildBossFinish(Player* pCreator)
// {
// 	if (!VALID(pCreator))
// 		return;

// 	// 记录当前血量
// 	INT32 nCurHP[X_Hero_Max_Group] = {0};
// 	for( INT32 i = 0 ; i < X_Hero_Max_Group; ++i )
// 	{
// 		EntityHero* pEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(i);
// 		if( VALID(pEntity) )
// 		{
// 			if( pEntity->IsDead() )
// 			{
// 				nCurHP[i] = 0;
// 			}
// 			else
// 			{
// 				nCurHP[i] = pEntity->GetAttController().GetAttValue(EHA_CurHP);
// 			}
// 		}
// 		else
// 		{
// 			nCurHP[i] = 0;
// 		}
// 	}

// 	const tagBossCaseEntry* pEntry = sGuildBoss.OnGuildBossFinish(pCreator, m_stMisc.guild_boss.n64GuildID, m_stMisc.guild_boss.nCase, m_eWinner, nCurHP);
// 	if( !VALID(pEntry) )
// 		return;

// 	// 增加个人贡献
// 	pCreator->ModGuildBossDevote((m_MuitlGroup[ESC_Defence].GetTotalDmg() / 10000) + sConstParam->nBossDevoteAttack );

// 	// 奖励
// 	list<tagLootData> listLoot;

// 	// 固定掉落
// 	sResMgr.GenLootData(pEntry->dwLootID, listLoot);

// 	// 发送消息
// 	CreateProtoMsg(msg, MS_GuildBossFinish, );
// 	msg << (bool)(m_eWinner == ESC_Attack);
// 	msg << (INT32)listLoot.size();
// 	list<tagLootData>::iterator it = listLoot.begin();
// 	while( it != listLoot.end() )
// 	{
// 		CreateProtoMsg(loot, LootData, );
// 		loot << (INT32)(*it).eType;
// 		loot << (UINT32)(*it).dwTypeMisc;
// 		loot << (INT32)(*it).nNum;

// 		msg << loot;

// 		pCreator->GainLoot(*it, ELCID_Guild_Boss_Attack);

// 		++it;
// 	}

// 	pCreator->SendMessage(msg);
// 	sLogMgr.LogActivity(pCreator->GetID(), X_GuildBoss_Activity_ID);
// }

// //-----------------------------------------------------------------------------
// // 夫妻塔
// //-----------------------------------------------------------------------------
// VOID Scene::OnCoupleTowerFinish(tagTeamSceneData* pData)
// {
// 	if(!VALID(pData))
// 		return;

// 	Room* pRoom = sRoomMgr.GetRoom(m_stMisc.couple_tower.dwRoomID);
// 	if(!VALID(pRoom))
// 		return;

// 	const tagCPTowerEntry* pCPTowerEntry = sResMgr.GetCPTowerEntry(m_stMisc.couple_tower.nFloor);
// 	if(!VALID(pCPTowerEntry))
// 		return;

// 	Marriage* pMarriage = sMarriageMgr.GetMarriage(pData->stEnterInfo.stMisc.couple_tower.n64MarriageID);
// 	if(!VALID(pMarriage))
// 		return;

// 	// 奖励
// 	list<tagLootData>::iterator it;
// 	list<tagLootData> listLoot;
// 	INT32 nFloor = pData->stEnterInfo.stMisc.couple_tower.nFloor;

// 	CreateProtoMsg(msgNotReward, MS_TeamRoomFinish, );
// 	msgNotReward << (bool)(m_eWinner == ESC_Attack);
// 	msgNotReward << (INT32)0;

// 	Player* pLeader = sPlayerMgr.GetPlayerByGUID(pData->stEnterInfo.stMisc.couple_tower.n64LeaderID);
// 	Player* pSpouse = sPlayerMgr.GetPlayerByGUID(pMarriage->GetSpouseID(pData->stEnterInfo.stMisc.couple_tower.n64LeaderID));

// 	// 固定掉落
// 	listLoot.clear();
// 	sResMgr.GenLootData(pCPTowerEntry->dwLootID, listLoot);

// 	if(VALID(pSpouse))
// 	{
// 		pRoom->Ready(pSpouse, FALSE);

// 		if( m_eWinner == ESC_Attack && nFloor > pSpouse->GetCPTowerFloor() )
// 		{
// 			// 掉落奖励
// 			it = listLoot.begin();
// 			while( it != listLoot.end() )
// 			{
// 				pSpouse->GainLoot(*it, ELCID_Couple_Tower_Reward);
// 				++it;
// 			}

// 			// 发送消息
// 			CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 			msg << (bool)(m_eWinner == ESC_Attack);
// 			msg << (INT32)listLoot.size();
// 			it = listLoot.begin();
// 			while( it != listLoot.end() )
// 			{
// 				CreateProtoMsg(loot, LootData, );
// 				loot << (INT32)(*it).eType;
// 				loot << (UINT32)(*it).dwTypeMisc;
// 				loot << (INT32)(*it).nNum;
// 				msg << loot;
// 				++it;
// 			}
// 			pSpouse->SendMessage(msg);

// 			pSpouse->SetCPTowerFloor(nFloor);
// 		}
// 		else
// 		{
// 			// 发送消息
// 			CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 			msg << (bool)(m_eWinner == ESC_Attack);
// 			msg << 0;
// 			pSpouse->SendMessage(msg);
// 		}
// 	}

// 	if(VALID(pLeader))
// 	{
// 		if( m_eWinner == ESC_Attack && nFloor > pLeader->GetCPTowerFloor() )
// 		{
// 			// 掉落奖励
// 			it = listLoot.begin();
// 			while( it != listLoot.end() )
// 			{
// 				pLeader->GainLoot(*it, ELCID_Couple_Tower_Reward);
// 				++it;
// 			}

// 			// 发送消息
// 			CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 			msg << (bool)(m_eWinner == ESC_Attack);
// 			msg << (INT32)listLoot.size();
// 			it = listLoot.begin();
// 			while( it != listLoot.end() )
// 			{
// 				CreateProtoMsg(loot, LootData, );
// 				loot << (INT32)(*it).eType;
// 				loot << (UINT32)(*it).dwTypeMisc;
// 				loot << (INT32)(*it).nNum;
// 				msg << loot;
// 				++it;
// 			}
// 			pLeader->SendMessage(msg);

// 			pLeader->SetCPTowerFloor(nFloor);
// 		}
// 		else
// 		{
// 			// 发送消息
// 			CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 			msg << (bool)(m_eWinner == ESC_Attack);
// 			msg << 0;
// 			pLeader->SendMessage(msg);
// 		}
// 	}

// 	pRoom->SyncRoomData();
// }

// //-----------------------------------------------------------------------------
// // 跨服帮会车轮战结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnGuildWheelFinish(tagCrossSceneData* pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	tagGroupRecord* pAttackRecord = &(pData->stTeam[ESC_Attack].stGroupRecord[0]);
// 	tagGroupRecord* pDefenceRecord = &(pData->stTeam[ESC_Defence].stGroupRecord[0]);

// 	// 获胜方当前血量
// 	INT32 nCurHP[X_Hero_Max_Group];
// 	for( INT32 i = 0 ; i < X_Hero_Max_Group; ++i )
// 	{
// 		EntityHero* pEntity = m_MuitlGroup[m_eWinner].GetEntityHero(i);
// 		if( VALID(pEntity) )
// 		{
// 			if( pEntity->IsDead() )
// 			{
// 				nCurHP[i] = 0;
// 			}
// 			else
// 			{
// 				nCurHP[i] = pEntity->GetAttController().GetAttValue(EHA_CurHP);
// 			}
// 		}
// 		else
// 		{
// 			nCurHP[i] = 0;
// 		}
// 	}

// 	if(sServer.GetWorldID() == pData->stTeam[ESC_Attack].dwWorldID )
// 	{
// 		tagMWB_WheelCombat stSend;
// 		stSend.n64GuildID = m_stMisc.guild_wheel.nAttackGuildID;
// 		stSend.stRecordInfo.stSrcPlayer.n64PlayerID = pAttackRecord->n64PlayerID;
// 		stSend.stRecordInfo.stSrcPlayer.nLevel = pAttackRecord->nLevel;
// 		stSend.stRecordInfo.stSrcPlayer.nProtrait = pAttackRecord->n16HeadProtrait;
// 		stSend.stRecordInfo.stSrcPlayer.nPlayerScore = pAttackRecord->nPlayerScore;
// 		memcpy(stSend.stRecordInfo.stSrcPlayer.szPlayerName, pAttackRecord->szName, sizeof(pAttackRecord->szName) );
// 		memcpy(stSend.stRecordInfo.stSrcPlayer.szGuildName, pAttackRecord->szGuildName, sizeof(pAttackRecord->szGuildName) );
// 		stSend.stRecordInfo.stDstPlayer.n64PlayerID = pDefenceRecord->n64PlayerID;
// 		stSend.stRecordInfo.stDstPlayer.nLevel = pDefenceRecord->nLevel;
// 		stSend.stRecordInfo.stDstPlayer.nProtrait = pDefenceRecord->n16HeadProtrait;
// 		stSend.stRecordInfo.stDstPlayer.nPlayerScore = pDefenceRecord->nPlayerScore;
// 		memcpy(stSend.stRecordInfo.stDstPlayer.szPlayerName, pDefenceRecord->szName, sizeof(pDefenceRecord->szName) );
// 		memcpy(stSend.stRecordInfo.stDstPlayer.szGuildName, pDefenceRecord->szGuildName, sizeof(pDefenceRecord->szGuildName) );
// 		stSend.stRecordInfo.n64RecordID = m_n64RecordID;
// 		stSend.stRecordInfo.bWin = (m_eWinner == ESC_Attack);
// 		memcpy(stSend.nCurHP, nCurHP, sizeof(stSend.nCurHP) );

// 		stSend.nAttackIndex =  pData->stMisc.guild_wheel.n16AttackIndex;
// 		stSend.nDefenceIndex = pData->stMisc.guild_wheel.n16DefenceIndex;

// 		if ( m_eWinner == ESC_Attack )
// 		{
// 			stSend.stRecordInfo.stSrcPlayer.nKillCount = pData->stMisc.guild_wheel.n16AttackKill + 1;
// 			stSend.stRecordInfo.stDstPlayer.nKillCount = pData->stMisc.guild_wheel.n16DefenceKill;
// 		}
// 		else
// 		{
// 			stSend.stRecordInfo.stSrcPlayer.nKillCount = pData->stMisc.guild_wheel.n16AttackKill;
// 			stSend.stRecordInfo.stDstPlayer.nKillCount = pData->stMisc.guild_wheel.n16DefenceKill+1;
// 		}

// 		// 通知battle战斗结果
// 		sBattleSession.SendMessage(&stSend, stSend.dwSize);
// 		sLogMgr.LogActivity(pAttackRecord->n64PlayerID, X_WHEEL_WAR_AVY);
// 	}
// 	if (sServer.GetWorldID() == pData->stTeam[ESC_Defence].dwWorldID )
// 	{
// 		tagMWB_WheelCombat stSend;
// 		stSend.n64GuildID = m_stMisc.guild_wheel.nDefenceGuildID;
// 		stSend.stRecordInfo.stSrcPlayer.n64PlayerID = pDefenceRecord->n64PlayerID;
// 		stSend.stRecordInfo.stSrcPlayer.nLevel = pDefenceRecord->nLevel;
// 		stSend.stRecordInfo.stSrcPlayer.nProtrait = pDefenceRecord->n16HeadProtrait;
// 		stSend.stRecordInfo.stSrcPlayer.nPlayerScore = pDefenceRecord->nPlayerScore;
// 		memcpy(stSend.stRecordInfo.stSrcPlayer.szPlayerName, pDefenceRecord->szName, sizeof(pDefenceRecord->szName) );
// 		memcpy(stSend.stRecordInfo.stSrcPlayer.szGuildName, pDefenceRecord->szGuildName, sizeof(pDefenceRecord->szGuildName) );

// 		stSend.stRecordInfo.stDstPlayer.n64PlayerID = pAttackRecord->n64PlayerID;
// 		stSend.stRecordInfo.stDstPlayer.nLevel = pAttackRecord->nLevel;
// 		stSend.stRecordInfo.stDstPlayer.nProtrait = pAttackRecord->n16HeadProtrait;
// 		stSend.stRecordInfo.stDstPlayer.nPlayerScore = pAttackRecord->nPlayerScore;
// 		memcpy(stSend.stRecordInfo.stDstPlayer.szPlayerName, pAttackRecord->szName, sizeof(pAttackRecord->szName) );
// 		memcpy(stSend.stRecordInfo.stDstPlayer.szGuildName, pAttackRecord->szGuildName, sizeof(pAttackRecord->szGuildName) );
// 		stSend.stRecordInfo.n64RecordID = m_n64RecordID;
// 		stSend.stRecordInfo.bWin = (m_eWinner == ESC_Defence);
// 		memcpy(stSend.nCurHP, nCurHP, sizeof(stSend.nCurHP) );
// 		stSend.nAttackIndex =  pData->stMisc.guild_wheel.n16AttackIndex;
// 		stSend.nDefenceIndex = pData->stMisc.guild_wheel.n16DefenceIndex;

// 		if ( m_eWinner == ESC_Attack )
// 		{
// 			stSend.stRecordInfo.stDstPlayer.nKillCount = pData->stMisc.guild_wheel.n16AttackKill + 1;
// 			stSend.stRecordInfo.stSrcPlayer.nKillCount = pData->stMisc.guild_wheel.n16DefenceKill;
// 		}
// 		else
// 		{
// 			stSend.stRecordInfo.stDstPlayer.nKillCount = pData->stMisc.guild_wheel.n16AttackKill;
// 			stSend.stRecordInfo.stSrcPlayer.nKillCount = pData->stMisc.guild_wheel.n16DefenceKill + 1;
// 		}

// 		// 通知battle战斗结果
// 		sBattleSession.SendMessage(&stSend, stSend.dwSize);
// 		sLogMgr.LogActivity(pDefenceRecord->n64PlayerID, X_WHEEL_WAR_AVY);
// 	}

// }

// //-----------------------------------------------------------------------------
// // 牢笼战场结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnIslandBattleFinish(tagCrossSceneData* pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	if (!m_stMisc.island_battle.bNeedRet)
// 		return;

// 	INT nModGrade = 0;
// 	bool bBoom = false;

// 	tagMWB_IslandCombat stSend;

// 	stSend.nIslandID = m_stMisc.island_battle.nIslandID;
// 	stSend.byHorizontal = m_stMisc.island_battle.byHorizontal;
// 	stSend.byVertical = m_stMisc.island_battle.byVertical;
// 	if (ESC_Attack == m_eWinner)
// 	{
// 		stSend.n64LoserID = m_MuitlGroup[ESC_Defence].GetPlayerID();
// 		nModGrade = m_stMisc.island_battle.n16AttackDead;
// 		bBoom = m_stMisc.island_battle.bAttackBoom;
// 	}
// 	else
// 	{
// 		stSend.n64LoserID = m_MuitlGroup[ESC_Attack].GetPlayerID();
// 		nModGrade = m_stMisc.island_battle.n16DefenceDead;
// 		bBoom = m_stMisc.island_battle.bDefenceBoom;
// 	}

// 	m_MuitlGroup[m_eWinner].Save2DB(&stSend.stRecord);

// 	// 获胜方当前血量
// 	for( INT32 i = 0 ; i < X_Max_Summon_Num; ++i )
// 	{
// 		EntityHero* pEntity = m_MuitlGroup[m_eWinner].GetEntityHero(i);
// 		if( VALID(pEntity) && !pEntity->IsDead() )
// 		{
// 			stSend.stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = pEntity->GetAttController().GetAttValue(EHA_CurHP);
// 		}
// 		else
// 		{
// 			stSend.stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 		}
// 	}

// 	// 移除越战越勇属性
// 	for (INT n = 0; n < X_Max_Summon_Num; n++)
// 	{
// 		if (VALID(stSend.stRecord.stHeroRecord[n].dwEntityID))
// 		{
// 			if (nModGrade > 0)
// 			{
// 				stSend.stRecord.stHeroRecord[n].nAttModPct[EHA_AttackPower] -= (nModGrade * X_Island_ComboDead_AtkMod);
// 				stSend.stRecord.stHeroRecord[n].nAttMod[EHA_DmgDec] -= (nModGrade * X_Island_ComboDead_DefMod);
// 			}
// 			if (bBoom)
// 			{
// 				stSend.stRecord.stHeroRecord[n].nAttModPct[EHA_DefenceMelee] -= X_Island_Boom_Mod;
// 				stSend.stRecord.stHeroRecord[n].nAttModPct[EHA_DefenceMagic] -= X_Island_Boom_Mod;
// 			}
// 		}
// 	}
// 	stSend.stRecord.CalAttMod(EHA_AttackPower);
// 	stSend.stRecord.CalAttMod(EHA_DmgDec);
// 	stSend.stRecord.CalAttMod(EHA_DefenceMelee);
// 	stSend.stRecord.CalAttMod(EHA_DefenceMagic);

// 	// 通知battle战斗结果
// 	sBattleSession.SendMessage(&stSend, stSend.dwSize);
// }

// //-----------------------------------------------------------------------------
// // 牧野之战结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnMakinoBattleFinish(tagCrossSceneData* pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	if (!m_stMisc.makino_battle.bNeedRet)
// 		return;

// 	INT nModGrade = 0;
// 	bool bFlag = false;

// 	tagMWB_MakinoCombat stSend;

// 	stSend.nMakinoID = m_stMisc.makino_battle.nMakinoID;
// 	stSend.byHorizontal = m_stMisc.makino_battle.byHorizontal;
// 	stSend.byVertical = m_stMisc.makino_battle.byVertical;
// 	if (ESC_Attack == m_eWinner)
// 	{
// 		stSend.n64LoserID = m_MuitlGroup[ESC_Defence].GetPlayerID();
// 		nModGrade = m_stMisc.makino_battle.n16AttackDead;
// 		bFlag = m_stMisc.makino_battle.bAttackFlag;
// 	}
// 	else
// 	{
// 		stSend.n64LoserID = m_MuitlGroup[ESC_Attack].GetPlayerID();
// 		nModGrade = m_stMisc.makino_battle.n16DefenceDead;
// 		bFlag = m_stMisc.makino_battle.bDefenceFlag;
// 	}

// 	m_MuitlGroup[m_eWinner].Save2DB(&stSend.stRecord);

// 	// 获胜方当前血量
// 	for( INT32 i = 0 ; i < X_Max_Summon_Num; ++i )
// 	{
// 		EntityHero* pEntity = m_MuitlGroup[m_eWinner].GetEntityHero(i);
// 		if( VALID(pEntity) && !pEntity->IsDead() )
// 		{
// 			stSend.stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = pEntity->GetAttController().GetAttValue(EHA_CurHP);
// 		}
// 		else
// 		{
// 			stSend.stRecord.stHeroRecord[i].nAtt[EHA_CurHP] = 0;
// 		}
// 	}

// 	// 移除越战越勇属性
// 	for (INT n = 0; n < X_Max_Summon_Num; n++)
// 	{
// 		if (VALID(stSend.stRecord.stHeroRecord[n].dwEntityID))
// 		{
// 			if (nModGrade > 0)
// 			{
// 				stSend.stRecord.stHeroRecord[n].nAttModPct[EHA_AttackPower] -= (nModGrade * X_Makino_ComboDead_AtkMod);
// 				stSend.stRecord.stHeroRecord[n].nAttMod[EHA_DmgDec] -= (nModGrade * X_Makino_ComboDead_DefMod);
// 			}
// 			if (bFlag)
// 			{
// 				stSend.stRecord.stHeroRecord[n].nAttModPct[EHA_DefenceMelee] -= X_Makino_Flag_Mod;
// 				stSend.stRecord.stHeroRecord[n].nAttModPct[EHA_DefenceMagic] -= X_Makino_Flag_Mod;
// 			}
// 		}
// 	}
// 	stSend.stRecord.CalAttMod(EHA_AttackPower);
// 	stSend.stRecord.CalAttMod(EHA_DmgDec);
// 	stSend.stRecord.CalAttMod(EHA_DefenceMelee);
// 	stSend.stRecord.CalAttMod(EHA_DefenceMagic);

// 	// 通知battle战斗结果
// 	sBattleSession.SendMessage(&stSend, stSend.dwSize);
// }

// //-----------------------------------------------------------------------------
// // 女娲遗迹结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnRemainsBattleFinish(tagCrossSceneData* pData)
// {
// 	const tagRoomEntry* pEntry= sResMgr.GetRoomEntry(m_stMisc.remains.dwRoomID);
// 	if(!VALID(pEntry))
// 		return;

// 	const tagRemainsEntry* pRemainsEntry = sResMgr.GetRemainsEntry(m_stMisc.remains.nRemainsFloor);
// 	if(!VALID(pRemainsEntry))
// 		return;

// 	for( INT32 i = 0; i < MAX_TEAM_MEM; ++i )
// 	{
// 		if( !VALID(pData->stTeam[ESC_Attack].stGroupRecord[i].n64PlayerID) )
// 			continue;

// 		// 奖励
// 		list<tagLootData>::iterator it;
// 		list<tagLootData> listLoot;
// 		listLoot.clear();

// 		if( m_eWinner == ESC_Attack )
// 		{
// 			// 固定掉落
// 			sResMgr.GenLootGroupData(pEntry->dwNormatlLootID, listLoot);
// 		}

// 		// 通知客户端
// 		CreateProtoMsg(msg, MS_TeamRoomFinish, );
// 		msg << (bool)(m_eWinner == ESC_Attack);

// 		Player* pPlayer = sPlayerMgr.GetPlayerByGUID(pData->stTeam[ESC_Attack].stGroupRecord[i].n64PlayerID);
// 		if( VALID(pPlayer) )
// 		{
// 			if( m_eWinner == ESC_Attack )
// 			{
// 				if((pPlayer->GetPlayerInfo()->nRemainsFloor + 1) == pData->stMisc.remains.nRemainsFloor)
// 				{
// 					sResMgr.GenLootData(pRemainsEntry->dwLootID, listLoot);
// 				}

// 				msg << (INT32)listLoot.size();
// 				it = listLoot.begin();
// 				while( it != listLoot.end() )
// 				{
// 					CreateProtoMsg(loot, LootData, );
// 					loot << (INT32)(*it).eType;
// 					loot << (UINT32)(*it).dwTypeMisc;
// 					loot << (INT32)((*it).nNum);

// 					msg << loot;

// 					++it;
// 				}

// 				// 掉落奖励
// 				it = listLoot.begin();
// 				while( it != listLoot.end() )
// 				{
// 					pPlayer->GainLoot(*it, ELCID_Remains_Reward);
// 					++it;
// 				}

// 				sLogMgr.LogFunctionPlaying(pPlayer->GetID(), ELCID_Remains_Reward, 2, m_stMisc.remains.dwRoomID, INVALID);

// 				if((pPlayer->GetPlayerInfo()->nRemainsFloor + 1) == m_stMisc.remains.nRemainsFloor)
// 					pPlayer->SetRemainsFloor(m_stMisc.remains.nRemainsFloor);

// 				evtRemainsPass evt;
// 				evt.nFloor = m_stMisc.remains.nRemainsFloor;
// 				pPlayer->SendEvent(evt);
// 			}
// 			else
// 			{
// 				msg << 0;
// 			}

// 			// 发送消息
// 			pPlayer->SendMessage(msg);
// 		}
// 	}
// }

// //-----------------------------------------------------------------------------
// // 山海战斗结束
// //-----------------------------------------------------------------------------
// VOID Scene::OnSquareBattleFinish(Player* pPlayer, tagSquareSceneData* pData, std::list<tagLootData>& listLoot)
// {
// 	if (!VALID(pPlayer) || !VALID(pData))
// 		return;

// 	tagSquareBeast* pAttackBeast = pPlayer->GetBeastContainer().GetBeast(pData->stMisc.square.dwAttackBeastTypeID);
// 	if (!VALID(pAttackBeast))
// 		return;

// 	EntityHero* pAttackEntity = m_MuitlGroup[ESC_Attack].GetEntityHero(1);
// 	EntityHero* pDefenceEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(1);
// 	if (!VALID(pAttackEntity) || !VALID(pDefenceEntity))
// 		return;

// 	// 抢夺资源
// 	if (pData->stSquareGrab.nGrabType == 1 && m_eWinner == ESC_Attack)
// 	{
// 		OnSquareGrabRes(pPlayer, pAttackBeast, pData);
// 	}

// 	// 抢夺城堡
// 	if (pData->stSquareGrab.nGrabType == 2 && m_eWinner == ESC_Attack)
// 	{
// 		OnSquareGrabCastle(pPlayer, pAttackBeast, pData);
// 	}

// 	// 通知battle战斗结果
// 	tagMWB_SquareCombatFinish stSend;
// 	memcpy(&stSend.stData, pData, sizeof(tagSquareSceneData));
// 	sBattleSession.SendMessage(&stSend, stSend.dwSize);

// 	// 通知客户端战斗结果
// 	/*CreateProtoMsg(msg, MS_SquareBattleFinish, );
// 	msg << (bool)(m_eWinner == ESC_Attack);

// 	// 进攻方战报
// 	msg << (UINT32)pAttackEntity->GetEntry()->dwTypeID;
// 	msg << string(pPlayer->GetPlayerName());
// 	msg << pAttackEntity->GetLevel();
// 	msg << m_MuitlGroup[ESC_Attack].GetPlayerScore();

// 	// 防守方战报
// 	msg << (UINT32)pDefenceEntity->GetEntry()->dwTypeID;
// 	msg << string(pData->stDefenceBeast.szName);
// 	msg << pDefenceEntity->GetLevel();
// 	msg << m_MuitlGroup[ESC_Defence].GetPlayerScore();

// 	msg << (INT32)listLoot.size();
// 	for (auto it = listLoot.begin(); it != listLoot.end(); ++it)
// 	{
// 		CreateProtoMsg(loot, LootData, );
// 		loot << (INT32)(*it).eType;
// 		loot << (UINT32)(*it).dwTypeMisc;
// 		loot << (INT32)((*it).nNum);
// 		msg << loot;
// 	}
// 	pPlayer->SendMessage(msg);*/

// 	// 发送山海通知
// 	SendOwnerSquareNotice(pPlayer, pData, listLoot);
// 	//　如果是pvp，添加跨服通知
// 	if (VALID(m_stMisc.square.n64Defencer))
// 	{
// 		SendEnemySquareNotice(pPlayer, pData);
// 	}
// }

// VOID Scene::SendOwnerSquareNotice(Player* pPlayer, tagSquareSceneData* pData, std::list<tagLootData>& listLoot)
// {
// 	if (!VALID(pPlayer) || !VALID(pData))
// 		return;
// 	tagSquareBeast* pAttackBeast = pPlayer->GetBeastContainer().GetBeast(pData->stMisc.square.dwAttackBeastTypeID);
// 	if (!VALID(pAttackBeast))
// 		return;
// 	// 生成山海通知
// 	tagSquareNotice stNotice;
// 	stNotice.n64OwnerID = pPlayer->GetID();
// 	stNotice.n64EnemyPlayerID = m_stMisc.square.n64Defencer;
// 	stNotice.dwErrorCode = E_Success;
// 	stNotice.dwTime = UCLOCK->CurrentClock();
// 	stNotice.dwTargetID = pData->dwMonsterTypeID;
// 	stNotice.nX = SQUARE_HIG_ID(pData->nPos);
// 	stNotice.nY = SQUARE_LOW_ID(pData->nPos);
// 	strcpy(stNotice.szEnemyName, pData->szEnemyName);
// 	strcpy(stNotice.szGuildName, pData->szGuildName);
// 	stNotice.dwWinner = m_eWinner;
// 	// 抢夺资源
// 	for (INT i = 0; i < ESMCT_RecordEnd; ++i)
// 	{
// 		stNotice.nGrabRes[i] = pData->stSquareGrab.nGrabNum[i];
// 	}
// 	switch(pData->nTargetType)
// 	{
// 	case ESMUT_Res:
// 		{
// 			stNotice.eType = ESNT_ResAttack;
// 			stNotice.dwTargetID = pData->dwResTypeID;
// 			for (INT n = 0; n != pData->stMisc.square.nDefencerBeastNum; ++n)
// 			{
// 				stNotice.stEnemyBeast[n].dwBeastTypeID = pData->stDefenceBeast.stBeastRecord[n].dwTypeID;
// 				stNotice.stEnemyBeast[n].nBeastLevel = pData->stDefenceBeast.stBeastRecord[n].nLevel;
// 			}
// 		}
// 		break;
// 	case ESMUT_Monster:
// 		{
// 			stNotice.eType = ESNT_Monster;
// 		}
// 		break;
// 	case ESMUT_Castle:
// 		{
// 			stNotice.eType = ESNT_Battle;
// 			for (INT n = 0; n != pData->stMisc.square.nDefencerBeastNum; ++n)
// 			{
// 				stNotice.stEnemyBeast[n].dwBeastTypeID = pData->stDefenceBeast.stBeastRecord[n+1].dwTypeID;
// 				stNotice.stEnemyBeast[n].nBeastLevel = pData->stDefenceBeast.stBeastRecord[n+1].nLevel;
// 			}
// 		}
// 		break;
// 	default:
// 		break;
// 	}
// 	memcpy(stNotice.n64RecordID, pData->n64RecordID, sizeof(stNotice.n64RecordID));
// 	INT nIndex = 0;
// 	for (auto it = listLoot.begin(); it != listLoot.end(); ++it)
// 	{
// 		stNotice.stAttachment[nIndex].eType = (*it).eType;
// 		stNotice.stAttachment[nIndex].dwTypeMisc = (*it).dwTypeMisc;
// 		stNotice.stAttachment[nIndex].nNum = (*it).nNum;
// 		nIndex++;
// 	}
// 	stNotice.stOwnerBeast[0].dwBeastTypeID = pAttackBeast->dwTypeID;
// 	stNotice.stOwnerBeast[0].nBeastLevel = pAttackBeast->nLevel;

// 	pPlayer->GetSquareController().AddSquareNotice(stNotice);
// }

// VOID Scene::SendEnemySquareNotice(Player* pPlayer, tagSquareSceneData* pData)
// {
// 	if (!VALID(pPlayer) || !VALID(pData))
// 		return;
// 	tagSquareBeast* pAttackBeast = pPlayer->GetBeastContainer().GetBeast(pData->stMisc.square.dwAttackBeastTypeID);
// 	if (!VALID(pAttackBeast))
// 		return;
// 	tagSquareNotice stDefendNotice;
// 	stDefendNotice.n64OwnerID = m_stMisc.square.n64Defencer;
// 	stDefendNotice.n64EnemyPlayerID = pPlayer->GetID();
// 	INT64 n64GuildID = pPlayer->GetGuildID();
// 	Guild *pGuild = sGuildMgr.GetGuild(n64GuildID);
// 	if (VALID(pGuild))
// 	{
// 		strcpy(stDefendNotice.szGuildName, pGuild->GetName());
// 	}
// 	strcpy(stDefendNotice.szEnemyName, pPlayer->GetPlayerName());
// 	stDefendNotice.dwErrorCode = E_Success;
// 	stDefendNotice.dwTime = UCLOCK->CurrentClock();
// 	stDefendNotice.nX = SQUARE_HIG_ID(pData->nPos);
// 	stDefendNotice.nY = SQUARE_LOW_ID(pData->nPos);
// 	stDefendNotice.dwWinner = !m_eWinner;
// 	stDefendNotice.dwTargetID = pData->stMisc.square.dwAttackBeastTypeID;
// 	// 被抢夺资源
// 	for (INT i = 0; i < ESMCT_RecordEnd; ++i)
// 	{
// 		stDefendNotice.nGrabRes[i] = pData->stSquareGrab.nGrabNum[i];
// 	}
// 	switch(pData->nTargetType)
// 	{
// 	case ESMUT_Res:
// 		{
// 			stDefendNotice.eType = ESNT_ResDefend;
// 			stDefendNotice.dwTargetID = pData->dwResTypeID;
// 			for (INT n = 0; n != pData->stMisc.square.nDefencerBeastNum; ++n)
// 			{
// 				stDefendNotice.stOwnerBeast[n].dwBeastTypeID = pData->stDefenceBeast.stBeastRecord[n].dwTypeID;
// 				stDefendNotice.stOwnerBeast[n].nBeastLevel = pData->stDefenceBeast.stBeastRecord[n].nLevel;
// 			}
// 		}
// 		break;
// 	case ESMUT_Castle:
// 		{
// 			stDefendNotice.eType = ESNT_Defend;
// 			for (INT n = 0; n != pData->stMisc.square.nDefencerBeastNum; ++n)
// 			{
// 				stDefendNotice.stOwnerBeast[n].dwBeastTypeID = pData->stDefenceBeast.stBeastRecord[n+1].dwTypeID;
// 				stDefendNotice.stOwnerBeast[n].nBeastLevel = pData->stDefenceBeast.stBeastRecord[n+1].nLevel;
// 			}
// 		}
// 		break;
// 	default:
// 		break;
// 	}
// 	memcpy(stDefendNotice.n64RecordID, pData->n64RecordID, sizeof(stDefendNotice.n64RecordID));
// 	stDefendNotice.stEnemyBeast[0].dwBeastTypeID = pAttackBeast->dwTypeID;
// 	stDefendNotice.stEnemyBeast[0].nBeastLevel = pAttackBeast->nLevel;

// 	tagMWB_AddSquareNotice stSendNotice;
// 	stSendNotice.n64PlayerID = m_stMisc.square.n64Defencer;
// 	memcpy(&stSendNotice.stNotice, &stDefendNotice, sizeof(stSendNotice.stNotice));
// 	sBattleSession.SendMessage(&stSendNotice, stSendNotice.dwSize);
// }

// //-----------------------------------------------------------------------------
// // 抢夺资源点
// //-----------------------------------------------------------------------------
// VOID Scene::OnSquareGrabRes(Player* pPlayer, tagSquareBeast* pAttackBeast, tagSquareSceneData* pData)
// {
// 	tagBeastScienceRecord* pScienceRecord = &(pData->stDefenceBeast.stBeastRecord[0].stScienceRecord);
// 	INT32 nTmpValue = 0;

// 	// 基础剩余
// 	INT32 nBasePercent = 5000;

// 	// 防守科技加层
// 	INT32 nGuardPercent = pScienceRecord->nScienceValue[EBST_1] + pScienceRecord->nResGuardValue;

// 	for (INT i = 0; i < ESMCT_RecordEnd; ++i)
// 	{
// 		// 剩余资源 = 资源总量 * （基础剩余 + ESST_28 + ESST_20） / 10000
// 		INT32 nTotalRes = pData->stSquareGrab.nResNum[i];
// 		INT32 nLeftNum = nTotalRes * (nBasePercent + nGuardPercent) / 10000;
// 		pData->stSquareGrab.nLeftNum[i] = nLeftNum;

// 		// 抢夺资源 = 资源总量 * （基础抢夺 + 对应异兽科技） / 10000
// 		INT32 nGrabPercent = nBasePercent - nGuardPercent;
// 		for (INT n = 0; n < SquareBeast_Passive_Skill_Num; ++n)
// 		{
// 			if (i == ESMCT_Gold)
// 				nGrabPercent += pPlayer->GetSquareController().GetScienceValueByID(pAttackBeast->dwPassiveSkillID[n], ESST_31, nTmpValue);

// 			if (i == ESMCT_Wood)
// 				nGrabPercent += pPlayer->GetSquareController().GetScienceValueByID(pAttackBeast->dwPassiveSkillID[n], ESST_29, nTmpValue);

// 			if (i == ESMCT_Stone)
// 				nGrabPercent += pPlayer->GetSquareController().GetScienceValueByID(pAttackBeast->dwPassiveSkillID[n], ESST_30, nTmpValue);
// 		}

// 		pData->stSquareGrab.nGrabNum[i] = nTotalRes * nGrabPercent / 10000;
// 	}
// }

// //-----------------------------------------------------------------------------
// // 抢夺城堡 （只用计算抢夺资源数）
// //-----------------------------------------------------------------------------
// VOID Scene::OnSquareGrabCastle(Player* pPlayer, tagSquareBeast* pAttackBeast, tagSquareSceneData* pData)
// {
// 	tagCastleScienceRecord* pCastleRecord = &(pData->stCastleRecord);

// 	for (INT i = 0; i < ESMCT_RecordEnd; ++i)
// 	{
// 		INT32 nMinute = CalcTimeDiff(UCLOCK->CurrentClock(), pCastleRecord->dwHarvestTime[i])/60;
// 		INT32 nProduce = min(nMinute * pCastleRecord->nIncPerMin[i], pCastleRecord->nMaxBuildProduce[i]);
// 		pData->stSquareGrab.nGrabNum[i] = pCastleRecord->nStorage[i] + nProduce;
// 	}
// }

// //-----------------------------------------------------------------------------
// // 保存录像
// //-----------------------------------------------------------------------------
// VOID Scene::Save2DB(BOOL bFinishSave, DWORD dwMisc)
// {
// 	if( m_bRecord )
// 		return;

// 	INT32 nRecordType = X_RECORD_TYPE_MAPPING[m_pEntry->eType];

// 	if( nRecordType == ESRT_Null )
// 		return;

// 	if( !VALID(m_MuitlGroup[ESC_Attack].GetPlayerID()) )
// 		return;

// 	if( !VALID(m_MuitlGroup[ESC_Defence].GetPlayerID()) && nRecordType != ESRT_Square )
// 	{
// 		if( nRecordType == ESRT_Grab || nRecordType == ESRT_Arena )
// 		{
// 			sFantasyMgr.SaveCombatData(this, NULL);
// 		}

// 		return;
// 	}

// 	INT64 nRecordID = PACK_RECORD_ID(m_dwID, sServer.GetWorldID());
// 	tagMWD_SaveSceneRecord* pMsg = sRecordMgr.CreateRecord(nRecordID);
// 	tagSceneRecord* pRecord = &pMsg->stRecord;
// 	pRecord->dwSceneID = m_pEntry->dwID;
// 	pRecord->nType = nRecordType;
// 	pRecord->dwSeed = m_dwSeed;
// 	pRecord->dwRecordTime = UCLOCK->CurrentClock();
// 	pRecord->nWinner = m_eWinner;
// 	pRecord->stMisc = m_stMisc;
// 	pRecord->nSize = sizeof(tagSceneRecord);
// 	m_MuitlGroup[ESC_Attack].Save2DB(&(pRecord->stGroupRecord[ESC_Attack]));
// 	m_MuitlGroup[ESC_Defence].Save2DB(&(pRecord->stGroupRecord[ESC_Defence]));
// 	m_n64RecordID = pRecord->n64RecordID;

// 	switch( nRecordType )
// 	{
// 	case ESRT_Grab:
// 		{
// 			if( !bFinishSave )
// 			{
// 				sRecordMgr.ReleaseRecord(pMsg);
// 				return;
// 			}

// 			Player* pCreator = sPlayerMgr.GetPlayerByGUID(m_n64Creator);
// 			pRecord->dwRecodeMisc[0] = pCreator->GetGrabController().GetTargetID();
// 			pRecord->dwRecodeMisc[1] = dwMisc;

// 			// 保存进攻方战斗数据
// 			sFantasyMgr.SaveCombatData(this, &(pRecord->stGroupRecord[ESC_Attack]));
// 		}
// 		break;

// 	case ESRT_Arena:
// 		{
// 			const tagArenaPlayer* pArenaAttacker = sArenaMgr.GetArenaPlayer(m_n64Creator);
// 			const tagArenaPlayer* pArenaDefencer = sArenaMgr.GetArenaPlayer(m_stMisc.arena.n64DstPlayerID);
// 			if (pArenaAttacker->stData.nCurRank <= pArenaDefencer->stData.nCurRank)
// 			{
// 				sRecordMgr.ReleaseRecord(pMsg);
// 				return;
// 			}

// 			pRecord->dwRecodeMisc[0] = pArenaAttacker->stData.nCurRank;
// 			pRecord->dwRecodeMisc[1] = pArenaDefencer->stData.nCurRank;

// 			// 保存进攻方战斗数据
// 			sFantasyMgr.SaveCombatData(this, &(pRecord->stGroupRecord[ESC_Attack]));
// 		}
// 		break;

// 	case ESRT_Mine:
// 		{
// 			pRecord->dwRecodeMisc[0] = m_stMisc.mine.dwMineID;
// 		}
// 		break;

// 	case ESRT_BHonour:
// 		{
// 			tagBHonourData* pSrcData = sBattleHonour.GetBHonourData(pRecord->stGroupRecord[ESC_Attack].n64PlayerID);
// 			tagBHonourData* pDesData = sBattleHonour.GetBHonourData(pRecord->stGroupRecord[ESC_Defence].n64PlayerID);
// 			if( VALID(pSrcData) && VALID(pDesData) )
// 			{
// 				INT32 nSrcRewardWin = pDesData->nWinNum;
// 				INT32 nDstRewardWin = pSrcData->nWinNum;

// 				if( pSrcData->nWinNum == -1 )
// 				{
// 					pSrcData->stRecord =  pRecord->stGroupRecord[ESC_Attack];
// 				}

// 				if( pDesData->nWinNum == -1 )
// 				{
// 					pDesData->stRecord =  pRecord->stGroupRecord[ESC_Defence];
// 				}

// 				if( m_eWinner == ESC_Attack )
// 				{
// 					nDstRewardWin = 0;
// 					nSrcRewardWin = Max(pSrcData->nWinNum+1, nSrcRewardWin);
// 					pRecord->dwRecodeMisc[0] = pSrcData->nWinNum+1;
// 					pRecord->dwRecodeMisc[1] = pDesData->nWinNum;
// 				}
// 				else
// 				{
// 					nSrcRewardWin = 0;
// 					nDstRewardWin = Max(pDesData->nWinNum+1, nDstRewardWin);
// 					pRecord->dwRecodeMisc[1] = pDesData->nWinNum+1;
// 					pRecord->dwRecodeMisc[0] = pSrcData->nWinNum;
// 				}

// 				pRecord->dwRecodeMisc[2] = nSrcRewardWin;
// 				pRecord->dwRecodeMisc[3] = nDstRewardWin;
// 				pRecord->dwRecodeMisc[4] = pSrcData->nRewardHonour;
// 				pRecord->dwRecodeMisc[5] = pDesData->nRewardHonour;
// 				pRecord->dwRecodeMisc[6] = pSrcData->GetLevel();
// 				pRecord->dwRecodeMisc[7] = pDesData->GetLevel();
// 			}
// 		}
// 		break;

// 	case ESRT_Castle:
// 		pRecord->dwRecodeMisc[0] = m_stMisc.castle_battle.nGrade;
// 		break;

// 	case ESRT_3v3:
// 		{
// 			s3v3Mgr.SaveRecord(m_stMisc.c3v3.n64TeamLightID, m_stMisc.c3v3.n64TeamDarkID, pRecord->n64RecordID);
// 		}
// 		break;
// 	case ESRT_Cross3v3:
// 		{
// 			/*tagCrossSceneData* pData = sSceneMgr.GetCrossTeamData(m_dwCrossSerial);
// 			if ( VALID(pData) )
// 			{
// 				INT64 n64AttackTeamID = INVALID, n64DefenceTeamID = INVALID;
// 				if ( pData->stTeam[ESC_Attack].dwWorldID == sServer.GetWorldID() )
// 				{
// 					n64AttackTeamID = pData->stTeam[ESC_Attack].n64TeamID;
// 				}

// 				if ( pData->stTeam[ESC_Defence].dwWorldID == sServer.GetWorldID() )
// 				{
// 					n64DefenceTeamID = pData->stTeam[ESC_Defence].n64TeamID;
// 				}

// 				sCross3v3TeamMgr.SaveRecord(n64AttackTeamID, n64DefenceTeamID, pRecord->n64RecordID);
// 			}
// 			*/

// 		}
// 		break;
// 	case ESRT_KingRank:
// 		{
// 			pRecord->dwRecodeMisc[0] = m_stMisc.kingrank.nAttRank;
// 			pRecord->dwRecodeMisc[1] = m_stMisc.kingrank.nDefRank;

// 			// 守方不是本服玩家，发送录像到守方所在服务器
// 			tagCrossSceneData* pData = sSceneMgr.GetCrossTeamData(m_dwCrossSerial);
// 			if ( VALID( pData) && ( pData->stTeam[ESC_Defence].dwWorldID != sServer.GetWorldID() ) )
// 			{
// 				tagMWB_SyncKingRankRecord stSend;
// 				stSend.dwWorldID = pData->stTeam[ESC_Defence].dwWorldID;
// 				stSend.stRecord	= *pRecord;

// 				sBattleSession.SendMessage(&stSend, stSend.dwSize);

// 			}

// 		}
// 		break;

// 	case ESRT_IslandBattle:
// 	case ESRT_MakinoBattle:
// 		{
// 			pRecord->dwRecodeMisc[0] = UCLOCK->PastPeriod(ETP_Day, 1);
// 		}
// 		break;

// 	case ESRT_Square:
// 		{
// 			// 守方不是本服玩家，发送录像到守方所在服务器
// 			tagSquareSceneData* pData = sSceneMgr.GetSquareSceneData(m_dwSquareSerial);
// 			for (INT i = 0; i != SquareNoticeRecord_Num; ++i)
// 			{
// 				if (!VALID(pData->n64RecordID[i]))
// 				{
// 					pData->n64RecordID[i] = nRecordID;
// 					break;
// 				}
// 			}

// 			if (VALID(pData) && (pData->stDefenceBeast.dwWorldID != sServer.GetWorldID()))
// 			{
// 				tagMWB_SyncSquareCombatRecord stSend;
// 				stSend.dwWorldID = pData->stDefenceBeast.dwWorldID;
// 				stSend.stRecord = *pRecord;

// 				sBattleSession.SendMessage(&stSend, stSend.dwSize);
// 			}

// 			// 保存异兽战力
// 			EntityHero* pAttackEntity = m_MuitlGroup[ESC_Attack].GetEntityHero(1);
// 			EntityHero* pDefenceEntity = m_MuitlGroup[ESC_Defence].GetEntityHero(1);
// 			if (VALID(pAttackEntity) && VALID(pDefenceEntity))
// 			{
// 				INT32 nAttackScore = (pAttackEntity->GetAttController().GetAttValue(EHA_AttackPower) * 2 + pAttackEntity->GetAttController().GetAttValue(EHA_MaxHP) / 4);
// 				INT32 nDefenceScore = (pDefenceEntity->GetAttController().GetAttValue(EHA_AttackPower) * 2 + pDefenceEntity->GetAttController().GetAttValue(EHA_MaxHP) / 4);
// 				m_MuitlGroup[ESC_Attack].SetPlayerScore(nAttackScore);
// 				m_MuitlGroup[ESC_Attack].SetPlayerScore(nDefenceScore);
// 			}
// 		}

// 	default:
// 		break;
// 	}

// 	sRecordMgr.Save2DB(pMsg);
// }

// //-------------------------------------------------------------------------
// // 发送场景消息
// //-------------------------------------------------------------------------
// VOID Scene::SendSceneMessage(Player* pPlayer, fxMessage& msg)
// {
// 	if( VALID(m_dwTeamSerial) || IsNeedBroadCast() )
// 	{
// 		if( m_bDelete )
// 		{
// 			SendTeamMessageByID(msg);
// 		}
// 		else
// 		{
// 			for( INT32 i = 0; i < m_nMaxTeamSyn; ++i )
// 			{
// 				m_pTeamSyn[i]->SendMessage(msg);
// 			}
// 		}
// 	}
// 	// 本服玩家
// 	else if (!VALID(m_dwServerID) || (m_dwServerID == sServer.GetWorldID()))
// 	{
// 		if( !VALID(pPlayer) )
// 		{
// 			pPlayer  = sPlayerMgr.GetPlayerByGUID(m_n64Creator);
// 		}

// 		if( VALID(pPlayer) && pPlayer->IsOnline() )
// 		{
// 			pPlayer->SendMessage(msg);
// 		}
// 	}
// 	// 跨服玩家
// 	else
// 	{
// 		sBattleSession.SendOutlandMessage(m_dwServerID, m_n64Creator, msg);
// 	}
// }

// VOID Scene::SendTeamMessageByID(fxMessage& msg)
// {
// 	for( INT32 i = 0; i < m_nMaxTeamSyn; ++i )
// 	{
// 		Player* pPlayer  = sPlayerMgr.GetPlayerByGUID(m_n64TeamSynID[i]);
// 		if( VALID(pPlayer) )
// 		{
// 			pPlayer->SendMessage(msg);
// 		}
// 	}
// }

// //-------------------------------------------- ---------------------------------
// // 劫镖结束 跨服
// //-----------------------------------------------------------------------------
// VOID Scene::OnSeaTradeCrossFinish(tagCrossSceneData* pData)
// {
// 	Player *pCreator = sPlayerMgr.GetPlayerByGUID(pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID);
// 	if(!VALID(pCreator)) return;

// 	if( pCreator->GetGrapSeaTrderNum() >= MAX_GRAP_SEA_TRADE_NUM ) return;
// 	// 扣除耐力
// 	//pCreator->GetAttController().ModAttValueWithLog(EPA_Stamina, -2, ELCID_Grap_Sea_TradeCross);

// 	//tagLootData stLootData;
// 	FLOAT fMoneyReward = 0.0f;
// 	FLOAT fHonourReward = 0.0f;
// 	FLOAT fRate = 1.0f;

// 	// 固定奖励
// 	if (m_eWinner == ESC_Attack)
// 	{
// 		if( pCreator->GetGrapSeaTrderNum() < MAX_GRAP_SEA_TRADE_NUM )
// 		{
// 			// 增加劫镖次数
// 			pCreator->IncGrapSeaTradeNum();
// 			// 随机掉落
// 			/*stLootData = sResMgr.GetRandomLootData();
// 			m_listLoot.push_back(stLootData);
// */
// 			if( sSeaTradeMgr.IsSeaHot() )
// 			{
// 				fRate = 1.5f;
// 			}
// 			const tagSeaTradeEntry* pEntry = sSeaTradeEntry(m_stMisc.sea_trade_cross.n16Level);
// 			fMoneyReward = pEntry->fMoneyRewardCross * TRADE_REWARD_RATE[m_stMisc.sea_trade_cross.n16GoodsType] * 0.2f * fRate;
// 			fHonourReward = (FLOAT)pEntry->fHonourRewardCross * TRADE_REWARD_RATE[m_stMisc.sea_trade_cross.n16GoodsType] * 0.2f * fRate;
// 			pCreator->GetCurrencyMgr().Inc(EMT_Gold, fMoneyReward, ELCID_Grap_Sea_TradeCross, INVALID);
// 			pCreator->GetCurrencyMgr().Inc(EMT_Honour, fHonourReward, ELCID_Grap_Sea_TradeCross, INVALID);
// 		}

// 		//send battle
// 		tagMWB_GrapSeaCrossFinish stSendBattle;
// 		stSendBattle.dwWorldID   = sServer.GetWorldID();
// 		stSendBattle.n64TargetID = m_stMisc.sea_trade_cross.n64TraderID;
// 		sBattleSession.SendMessage(&stSendBattle, stSendBattle.dwSize);
// 	}
// 	else
// 	{
// 		pCreator->SetGrapCd(IncTime(UCLOCK->CurrentClock(), X_GRAP_CD));
// 	}

// 	evtGrapSeaTrade evt;
// 	evt.n64Winner = (m_eWinner == ESC_Attack) ? pCreator->GetID() : m_MuitlGroup[ESC_Defence].GetPlayerID();
// 	evt.n64Attacker = pCreator->GetID();
// 	pCreator->SendEvent(evt);

// 	CreateProtoMsg(msg, MS_SeaCrossTradeFinish, );
// 	msg << (INT32)fMoneyReward;
// 	msg << (INT32)fHonourReward;

// 	CreateProtoMsg(result, PVPResult, );
// 	result << (bool)(m_eWinner == ESC_Attack);
// 	result << (UINT32)m_dwServerID;
// 	result << m_n64RecordID;

// 	/*CreateProtoMsg(loot_data, LootData, );
// 	loot_data << (INT32)stLootData.eType;
// 	loot_data << (UINT32)stLootData.dwTypeMisc;
// 	loot_data << (INT32)stLootData.nNum;

// 	result << loot_data;*/

// 	DWORD dwEntityID[X_Max_Summon_Num];
// 	INT nEntityNum = m_MuitlGroup[ESC_Attack].ExportEntityID(dwEntityID);
// 	result << string(pCreator->GetPlayerInfo()->szPlayerName);
// 	result << (INT32)pCreator->GetLevel();
// 	result << (INT32)pCreator->GetPlayerScore();
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	//添加英雄的等级，星级，战斗结束界面新需求
// 	DWORD dwLevel[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	DWORD dwStar[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	DWORD dwQuality[X_Max_Summon_Num];
// 	m_MuitlGroup[ESC_Attack].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);

// 	result << string(m_MuitlGroup[ESC_Defence].GetPlayerName());
// 	result << m_MuitlGroup[ESC_Defence].GetPlayerLevel();
// 	result << m_MuitlGroup[ESC_Defence].GetPlayerScore();
// 	nEntityNum = m_MuitlGroup[ESC_Defence].ExportEntityID(dwEntityID);
// 	InsertArrayProtoMsg(result, dwEntityID, nEntityNum, UINT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityLevel(dwLevel);
// 	InsertArrayProtoMsg(result, dwLevel, nEntityNum, INT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityStar(dwStar);
// 	InsertArrayProtoMsg(result, dwStar, nEntityNum, INT32);
// 	m_MuitlGroup[ESC_Defence].ExportEntityQuality(dwQuality);
// 	InsertArrayProtoMsg(result, dwQuality, nEntityNum, INT32);

// 	msg << result;

// 	pCreator->SendMessage(msg);
// }

// //-----------------------------------------------------------------------------
// /// 好友pk
// //-----------------------------------------------------------------------------
// VOID Scene::OnFriendPkFinish(tagCrossSceneData *pData)
// {
// 	if ( !VALID(pData) )
// 		return;

// 	if ( sServer.GetWorldID() == pData->stTeam[ESC_Attack].dwWorldID )
// 	{
// 		Player *pPlayer = sPlayerMgr.GetPlayerByGUID( pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID );
// 		// 发送战斗录像
// 		if ( VALID(pPlayer ) )
// 		{
// 			CreateProtoMsg(msg, MS_FriendPkCombatFinish, );
// 			msg << m_n64RecordID;
// 			msg << pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID;

// 			pPlayer->SendMessage(msg);
// 		}
// 	}

// 	if ( sServer.GetWorldID() == pData->stTeam[ESC_Defence].dwWorldID )
// 	{
// 		Player *pPlayer = sPlayerMgr.GetPlayerByGUID( pData->stTeam[ESC_Defence].stGroupRecord[0].n64PlayerID );
// 		// 发送战斗录像
// 		if ( VALID(pPlayer ) )
// 		{
// 			CreateProtoMsg(msg, MS_FriendPkCombatFinish, );
// 			msg << m_n64RecordID;
// 			msg << pData->stTeam[ESC_Attack].stGroupRecord[0].n64PlayerID;

// 			pPlayer->SendMessage(msg);
// 		}
// 	}
// }
