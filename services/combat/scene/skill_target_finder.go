package scene

import (
	"container/list"

	"bitbucket.org/funplus/server/define"
)

type SkillTargetSelector func(*Skill)

var (
	launchTypeFn map[int32]SkillTargetSelector
	targetTypeFn map[int32]SkillTargetSelector
	rangeTypeFn  map[int32]SkillTargetSelector
	scopeTypeFn  map[int32]SkillTargetSelector
)

func init() {
	launchTypeFn = make(map[int32]SkillTargetSelector)
	targetTypeFn = make(map[int32]SkillTargetSelector)
	rangeTypeFn = make(map[int32]SkillTargetSelector)
	scopeTypeFn = make(map[int32]SkillTargetSelector)

	// 注册发起类型处理器
	registerLaunchType()

	// 注册目标类型处理器
	registerTargetType()

	// 注册技能范围类型处理器
	registerRangeType()

	// 注册技能作用对象处理器
	registerScopeType()
}

func registerLaunchType() {
	// 以技能caster为发起
	launchTypeFn[define.SkillLaunchType_Self] = func(s *Skill) {
		s.opts.TargetPosition = s.opts.Caster.GetPosition()
	}

	// 以技能target为发起
	launchTypeFn[define.SkillLaunchType_Target] = func(s *Skill) {
		s.opts.TargetPosition = s.opts.Target.GetPosition()
	}
}

func registerTargetType() {
	// 自身周围
	targetTypeFn[define.SkillTargetType_SelfRound] = func(s *Skill) {
		s.opts.TargetPosition = s.opts.Caster.GetPosition()
	}

	// 选定空间
	targetTypeFn[define.SkillTargetType_SelectRound] = func(s *Skill) {
	}

	// 友军单体
	targetTypeFn[define.SkillTargetType_FriendlySingle] = func(s *Skill) {
		s.listTargets.Init()
		s.listTargets.PushBack(s.opts.Target)
	}

	// 敌军单体
	targetTypeFn[define.SkillTargetType_EnemySingle] = func(s *Skill) {
		s.listTargets.Init()
		s.listTargets.PushBack(s.opts.Target)
	}
}

func registerRangeType() {
	// 圆形
	rangeTypeFn[define.SkillRangeType_Circle] = func(s *Skill) {
		skillStartPosition := s.opts.TargetPosition
		radius := int32(s.opts.Entry.TargetLength.IntPart())

		var next *list.Element
		for e := s.listTargets.Front(); e != nil; e = next {
			next = e.Next()
			targetPosition := e.Value.(*SceneEntity).GetPosition()

			// 判断距离是否在技能范围内
			if !IsInDistance(skillStartPosition, targetPosition, radius) {
				s.listTargets.Remove(e)
			}
		}
	}

	// 矩形
	rangeTypeFn[define.SkillRangeType_Rectangle] = func(s *Skill) {

	}

	// 扇形
	rangeTypeFn[define.SkillRangeType_Fan] = func(s *Skill) {

	}
}

func registerScopeType() {
	// 选定目标
	scopeTypeFn[define.SkillScopeType_SelectTarget] = func(s *Skill) {
		var next *list.Element
		for e := s.listTargets.Front(); e != nil; e = next {
			next = e.Next()
			if e.Value.(*SceneEntity) != s.opts.Target {
				s.listTargets.Remove(e)
			}
		}
	}

	// 除目标以外的友军
	scopeTypeFn[define.SkillScopeType_FriendlyTroops] = func(s *Skill) {
		var next *list.Element
		for e := s.listTargets.Front(); e != nil; e = next {
			next = e.Next()
			target := e.Value.(*SceneEntity)

			// 目标
			if target == s.opts.Target {
				s.listTargets.Remove(e)
				continue
			}

			// 非友军
			if target.GetCamp().camp != s.opts.Caster.GetCamp().camp {
				s.listTargets.Remove(e)
				continue
			}
		}
	}

	// 所有友军
	scopeTypeFn[define.SkillScopeType_AllFriendlyTroops] = func(s *Skill) {
		var next *list.Element
		for e := s.listTargets.Front(); e != nil; e = next {
			next = e.Next()
			target := e.Value.(*SceneEntity)

			// 非友军
			if target.GetCamp().camp != s.opts.Caster.GetCamp().camp {
				s.listTargets.Remove(e)
				continue
			}
		}
	}

	// 除目标以外敌军
	scopeTypeFn[define.SkillScopeType_EnemyTroops] = func(s *Skill) {
		var next *list.Element
		for e := s.listTargets.Front(); e != nil; e = next {
			next = e.Next()
			target := e.Value.(*SceneEntity)

			// 目标
			if target == s.opts.Target {
				s.listTargets.Remove(e)
				continue
			}

			// 非敌军
			if target.GetCamp().camp == s.opts.Caster.GetCamp().camp {
				s.listTargets.Remove(e)
				continue
			}
		}
	}

	// 所有敌军
	scopeTypeFn[define.SkillScopeType_AllEnemyTroops] = func(s *Skill) {
		var next *list.Element
		for e := s.listTargets.Front(); e != nil; e = next {
			next = e.Next()
			target := e.Value.(*SceneEntity)

			// 非敌军
			if target.GetCamp().camp == s.opts.Caster.GetCamp().camp {
				s.listTargets.Remove(e)
				continue
			}
		}
	}
}

// 筛选目标列表
func (s *Skill) selectTargets() {
	// 混乱状态特殊处理
	// if s.opts.SpellType == define.SpellType_Melee {
	// 	if s.opts.Caster == nil {
	// 		return
	// 	}

	// 	if s.opts.Caster.HasState(define.HeroState_Chaos) {
	// 		s.findTargetChaos()
	// 		return
	// 	}

	// 	if s.opts.Caster.HasState(define.HeroState_Taunt) {
	// 		s.findTargetEnemySingle()
	// 		return
	// 	}
	// }

	// 通过发起类型获取技能开始时坐标
	launchTypeFn[s.opts.Entry.SkillLaunch](s)

	// 通过目标类型筛选目标
	targetTypeFn[s.opts.Entry.TargetType](s)

	// 通过技能范围筛选目标
	rangeTypeFn[s.opts.Entry.RangeType](s)

	// 通过作用对象筛选目标
	scopeTypeFn[s.opts.Entry.Scope](s)
}

// 选定目标
func (s *Skill) findTargetSelectTarget() {
	s.listTargets.PushBack(s.opts.Target)
}

// 除目标外友军
func (s *Skill) findTargetFriendlyTroops() {

}

//-------------------------------------------------------------------------------
// 自己
//-------------------------------------------------------------------------------
func (s *Skill) findTargetSelf() {
	/*if (!VALID(m_pCaster))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	if( IsTargetValid(m_pCaster) )
		m_listTarget.PushBack(m_pCaster);
	*/
}

//-------------------------------------------------------------------------------
// 目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemySingle() {
	/*if (!VALID(m_pCaster) || !VALID(m_pTarget))
		return;

	Scene* pCasterScene = m_pCaster->GetScene();
	Scene* pTargetScene = m_pTarget->GetScene();
	if (!VALID(pCasterScene) || !VALID(pTargetScene) || (pCasterScene != pTargetScene))
		return;

	if (IsTargetValid(m_pTarget))
		m_listTarget.PushBack(m_pTarget);
	*/
}

//-------------------------------------------------------------------------------
// 敌方后排单体目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemySingleBack() {
	/*if (!VALID(m_pCaster))
		return;

	// 反击后排单体还是攻击当前指定目标
	if(ERMT_BeatBack == m_eSpellType)
	{
		if ( VALID(m_pTarget) && IsTargetValid(m_pTarget))
		{
			m_listTarget.PushBack(m_pTarget);
			return;
		}
	}

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());
	EntityHero* pEntity = static_cast<EntityGroup*>(m_pCaster->GetFather())->FindTargetByPriority(m_pCaster->GetLocation() & 0xF, &group, FALSE);

	if (IsTargetValid(pEntity))
		m_listTarget.PushBack(pEntity);
	*/
}

//-------------------------------------------------------------------------------
// 友方血量最少目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetFrienHPMin() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetCamp());

	TList<EntityHero*> listTemp;

	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (!IsTargetValid(pEntity))
			continue;

		listTemp.PushBack(pEntity);
	}

	listTemp.GetList().sort(ThreatSortMinHP);

	INT nTargetNum = Min(listTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		m_listTarget.PushBack(listTemp.PopFront());
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方血量最多目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemyHPMax() {
	/*if (!VALID(m_pCaster))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	TList<EntityHero*> listTemp;

	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (!IsTargetValid(pEntity))
			continue;

		listTemp.PushBack(pEntity);
	}

	listTemp.GetList().sort(ThreatSortMaxHP);

	INT nTargetNum = Min(listTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		m_listTarget.PushBack(listTemp.PopFront());
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方怒气最多目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemyRageMax() {
	/*if (!VALID(m_pCaster))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	TList<EntityHero*> listTemp;

	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (!IsTargetValid(pEntity))
			continue;

		listTemp.PushBack(pEntity);
	}

	listTemp.GetList().sort(ThreatSortMaxRage);

	INT nTargetNum = Min(listTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		m_listTarget.PushBack(listTemp.PopFront());
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方直线目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemyColumn() {
	/*if (!VALID(m_pCaster) || !VALID(m_pTarget))
		return;

	if (IsTargetValid(m_pTarget))
		m_listTarget.PushBack(m_pTarget);

	INT nIndex = m_pTarget->GetLocation() & 0xF;
	if (nIndex < 3)
	{
		Scene* pScene = m_pTarget->GetScene();
		if (!VALID(pScene))
			return;

		EntityHero* pEntity = pScene->GetGroup(m_pTarget->GetCamp()).GetEntityHero(nIndex + 3);
		if (!IsTargetValid(pEntity))
			return;

		m_listTarget.PushBack(pEntity);
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方横排目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemyFrontline(checkBack bool) {
	/*if (!VALID(m_pCaster))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	BOOL bHasTarget = FALSE;
	for (INT n = 0; n < 3; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);
		if (!VALID(pEntity))
			continue;

		if (IsTargetValid(pEntity))
		{
			bHasTarget = TRUE;
			m_listTarget.PushBack(pEntity);
		}
	}

	if (!bHasTarget && bCheckBack)
	{
		for (INT n = 3; n < 6; n++)
		{
			EntityHero* pEntity = group.GetEntityHero(n);
			if (VALID(pEntity))
			{
				FindTargetEnemySupporter(FALSE);
				break;
			}
		}
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方后排目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemySupporter(checkFront bool) {
	/*if (!VALID(m_pCaster))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	BOOL bHasTarget = FALSE;
	for (INT n = 3; n < 6; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);
		if (!VALID(pEntity))
			continue;

		if (IsTargetValid(pEntity))
		{
			bHasTarget = TRUE;
			m_listTarget.PushBack(pEntity);
		}
	}

	if (!bHasTarget && bCheckFront)
	{
		for (INT n = 0; n < 3; n++)
		{
			EntityHero* pEntity = group.GetEntityHero(n);
			if (VALID(pEntity))
			{
				FindTargetEnemyFrontline(FALSE);
				break;
			}
		}
	}
	*/
}

//-------------------------------------------------------------------------------
// 友方随机目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetFriendRandom() {
	/*if (!VALID(m_pCaster))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetCamp());

	SimpleVector<INT> vecTemp(6);
	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);
		if (!VALID(pEntity))
			continue;

		if (IsTargetValid(pEntity))
			vecTemp.PushBack(n);
	}

	if (vecTemp.Empty())
		return;

	INT nTargetNum = Min(vecTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		INT nRandIndex = m_pCaster->GetScene()->GetRandom().Rand(0, (vecTemp.Size()-n-1));
		m_listTarget.PushBack(group.GetEntityHero(vecTemp[nRandIndex]));
		vecTemp[nRandIndex] = vecTemp[(vecTemp.Size()-n-1)];
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方随机目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemyRandom() {
	/*if (!VALID(m_pCaster))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	SimpleVector<INT> vecTemp(6);
	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);
		if (!VALID(pEntity))
			continue;

		if (IsTargetValid(pEntity))
			vecTemp.PushBack(n);
	}

	if (vecTemp.Empty())
		return;

	INT nTargetNum = Min(vecTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		INT nRandIndex = m_pCaster->GetScene()->GetRandom().Rand(0,(vecTemp.Size()-n-1));
		m_listTarget.PushBack(group.GetEntityHero(vecTemp[nRandIndex]));
		vecTemp[nRandIndex] = vecTemp[(vecTemp.Size()-n-1)];
	}
	*/
}

//-------------------------------------------------------------------------------
// 友方全体目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetFriendAll() {
	/*if (!VALID(m_pCaster))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetCamp());

	SimpleVector<INT> vecTemp(6);
	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);
		if (!VALID(pEntity))
			continue;

		if (IsTargetValid(pEntity))
			m_listTarget.PushBack(pEntity);
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方全体目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemyAll() {
	/*if (!VALID(m_pCaster))
		return;

	Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	SimpleVector<INT> vecTemp(6);
	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);
		if (!VALID(pEntity))
			continue;

		if (IsTargetValid(pEntity))
			m_listTarget.PushBack(pEntity);
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方符文携带者
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemyRune() {
	//if (!VALID(m_pCaster))
	//	return;

	//Scene* pScene = m_pCaster->GetScene();
	//if (!VALID(pScene))
	//	return;

	//EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	//SimpleVector<INT> vecTemp(6);
	//for (INT n = 0; n < X_Max_Summon_Num; n++)
	//{
	//	EntityHero* pEntity = group.GetEntityHero(n);
	//	if (!VALID(pEntity))
	//		continue;

	//	if( !pEntity->IsCarryRune() )
	//		continue;

	//	if (IsTargetValid(pEntity))
	//		vecTemp.PushBack(n);
	//}

	//if (vecTemp.Empty())
	//	return;

	//INT nTargetNum = Min(vecTemp.Size(), m_pEntry->nTargetNum);
	//for (INT n = 0; n < nTargetNum; n++)
	//{
	//	INT nRandIndex = m_pCaster->GetScene()->GetRandom().Rand(0, nTargetNum - n - 1);
	//	m_listTarget.PushBack(group.GetEntityHero(vecTemp[nRandIndex]));
	//	vecTemp[nRandIndex] = vecTemp[nTargetNum - n - 1];
	//}
}

//-------------------------------------------------------------------------------
// 友方符文携带者
//-------------------------------------------------------------------------------
func (s *Skill) findTargetFriendRune() {
	//if (!VALID(m_pCaster))
	//	return;

	//Scene* pScene = m_pCaster->GetScene();
	//if (!VALID(pScene))
	//	return;

	//EntityGroup& group = pScene->GetGroup(m_pCaster->GetCamp());

	//SimpleVector<INT> vecTemp(6);
	//for (INT n = 0; n < X_Max_Summon_Num; n++)
	//{
	//	EntityHero* pEntity = group.GetEntityHero(n);
	//	if (!VALID(pEntity))
	//		continue;

	//	if( !pEntity->IsCarryRune() )
	//		continue;

	//	if (IsTargetValid(pEntity))
	//		vecTemp.PushBack(n);
	//}

	//if (vecTemp.Empty())
	//	return;

	//INT nTargetNum = Min(vecTemp.Size(), m_pEntry->nTargetNum);
	//for (INT n = 0; n < nTargetNum; n++)
	//{
	//	INT nRandIndex = m_pCaster->GetScene()->GetRandom().Rand(0, nTargetNum - n - 1);
	//	m_listTarget.PushBack(group.GetEntityHero(vecTemp[nRandIndex]));
	//	vecTemp[nRandIndex] = vecTemp[nTargetNum - n - 1];
	//}
}

//-------------------------------------------------------------------------------
// 下一个将要行动的敌人
//-------------------------------------------------------------------------------
func (s *Skill) findTargetNextAttack() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	INT16 n16LoopIndx = group.GetLoopIndex();

	for (INT n = n16LoopIndx; n < n16LoopIndx + X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n % X_Max_Summon_Num);

		if (!IsTargetValid(pEntity))
			continue;;

		m_listTarget.PushBack(pEntity);
		break;
	}
	*/
}

//-------------------------------------------------------------------------------
// 友方怒气最低
//-------------------------------------------------------------------------------
func (s *Skill) findTargetFriendRageMin() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetCamp());

	TList<EntityHero*> listTemp;

	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (!IsTargetValid(pEntity))
			continue;

		listTemp.PushBack(pEntity);
	}

	listTemp.GetList().sort(ThreatSortMinRage);

	INT nTargetNum = Min(listTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		m_listTarget.PushBack(listTemp.PopFront());
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌放前横排随机
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemyFrontLineRandom() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	SimpleVector<INT> vecTemp(3);
	for (INT n = 0; n < 3; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (IsTargetValid(pEntity))
		{
			vecTemp.PushBack(n);
		}
	}

	if( vecTemp.Size() == 0 )
	{
		for (INT n = 3; n < 6; n++)
		{
			EntityHero* pEntity = group.GetEntityHero(n);

			if (IsTargetValid(pEntity))
			{
				vecTemp.PushBack(n);
			}
		}
	}

	INT nTargetNum = Min(vecTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		INT nRandIndex = m_pCaster->GetScene()->GetRandom().Rand(0, (vecTemp.Size()-n-1));
		m_listTarget.PushBack(group.GetEntityHero(vecTemp[nRandIndex]));
		vecTemp[nRandIndex] = vecTemp[(vecTemp.Size()-n-1)];
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌放后横排随机
//-------------------------------------------------------------------------------
func (s *Skill) findTargetEnemyBackLineRandom() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	SimpleVector<INT> vecTemp(3);
	for (INT n = 3; n < 6; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (IsTargetValid(pEntity))
		{
			vecTemp.PushBack(n);
		}
	}

	if( vecTemp.Size() == 0 )
	{
		for (INT n = 0; n < 3; n++)
		{
			EntityHero* pEntity = group.GetEntityHero(n);

			if (IsTargetValid(pEntity))
			{
				vecTemp.PushBack(n);
			}
		}
	}

	INT nTargetNum = Min(vecTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		INT nRandIndex = m_pCaster->GetScene()->GetRandom().Rand(0, (vecTemp.Size()-n-1));
		m_listTarget.PushBack(group.GetEntityHero(vecTemp[nRandIndex]));
		vecTemp[nRandIndex] = vecTemp[(vecTemp.Size()-n-1)];
	}
	*/
}

//-------------------------------------------------------------------------------
// 友方前横排随机
//-------------------------------------------------------------------------------
func (s *Skill) findTargetFriendFrontLineRandom() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetCamp());

	SimpleVector<INT> vecTemp(3);
	for (INT n = 0; n < 3; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (IsTargetValid(pEntity))
		{
			vecTemp.PushBack(n);
		}
	}

	if( vecTemp.Size() == 0 )
	{
		for (INT n = 3; n < 6; n++)
		{
			EntityHero* pEntity = group.GetEntityHero(n);

			if (IsTargetValid(pEntity))
			{
				vecTemp.PushBack(n);
			}
		}
	}

	INT nTargetNum = Min(vecTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		INT nRandIndex = m_pCaster->GetScene()->GetRandom().Rand(0, (vecTemp.Size()-n-1));
		m_listTarget.PushBack(group.GetEntityHero(vecTemp[nRandIndex]));
		vecTemp[nRandIndex] = vecTemp[(vecTemp.Size()-n-1)];
	}
	*/
}

//-------------------------------------------------------------------------------
// 友方后横排随机
//-------------------------------------------------------------------------------
func (s *Skill) findTargetFriendBackLineRandom() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetCamp());

	SimpleVector<INT> vecTemp(3);
	for (INT n = 3; n < 6; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (IsTargetValid(pEntity))
		{
			vecTemp.PushBack(n);
		}
	}

	if( vecTemp.Size() == 0 )
	{
		for (INT n = 0; n < 3; n++)
		{
			EntityHero* pEntity = group.GetEntityHero(n);

			if (IsTargetValid(pEntity))
			{
				vecTemp.PushBack(n);
			}
		}
	}

	INT nTargetNum = Min(vecTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		INT nRandIndex = m_pCaster->GetScene()->GetRandom().Rand(0, (vecTemp.Size()-n-1));
		m_listTarget.PushBack(group.GetEntityHero(vecTemp[nRandIndex]));
		vecTemp[nRandIndex] = vecTemp[(vecTemp.Size()-n-1)];
	}
	*/
}

//-------------------------------------------------------------------------------
// 下一个将要行动的敌人所在横排
//-------------------------------------------------------------------------------
func (s *Skill) findTargetNextAttackRow() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	INT16 n16LoopIndx = group.GetLoopIndex();

	for (INT n = n16LoopIndx; n < n16LoopIndx + X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n % X_Max_Summon_Num);

		if (!IsTargetValid(pEntity))
			continue;

		INT nIndex = pEntity->GetLocation() & 0xF;

		if( nIndex < 3 )
		{
			for (INT n = 0; n < 3; n++)
			{
				EntityHero* pEntity = group.GetEntityHero(n);

				if (IsTargetValid(pEntity))
				{
					m_listTarget.PushBack(pEntity);
				}
			}
		}
		else
		{
			for (INT n = 3; n < 6; n++)
			{
				EntityHero* pEntity = group.GetEntityHero(n);

				if (IsTargetValid(pEntity))
				{
					m_listTarget.PushBack(pEntity);
				}
			}
		}

		break;
	}
	*/
}

//-------------------------------------------------------------------------------
// 下一个将要行动的敌人所在竖排
//-------------------------------------------------------------------------------
func (s *Skill) findTargetNextAttackConlumn() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	INT16 n16LoopIndx = group.GetLoopIndex();

	for (INT n = n16LoopIndx; n < n16LoopIndx + X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n % X_Max_Summon_Num);

		if (!IsTargetValid(pEntity))
			continue;

		m_listTarget.PushBack(pEntity);
		INT nIndex = pEntity->GetLocation() & 0xF;

		if( nIndex < 3 )
		{
			EntityHero* pEntity = group.GetEntityHero(nIndex + 3);

			if (IsTargetValid(pEntity))
			{
				m_listTarget.PushBack(pEntity);
			}
		}
		else
		{
			EntityHero* pEntity = group.GetEntityHero(nIndex - 3);

			if (IsTargetValid(pEntity))
			{
				m_listTarget.PushBack(pEntity);
			}
		}
		break;
	}
	*/
}

//-------------------------------------------------------------------------------
// 将要行动的敌人相邻目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetNextAttackBorder() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	INT16 n16LoopIndx = group.GetLoopIndex();
	INT16 nRealIndex = INVALID;

	for (INT n = n16LoopIndx; n < n16LoopIndx + X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n % X_Max_Summon_Num);

		if ( IsTargetValid(pEntity))
		{
			nRealIndex = n % X_Max_Summon_Num;
			break;
		}
	}

	if( VALID(nRealIndex) )
	{
		for( INT32 i = 0; i < X_Max_Summon_Num; ++i )
		{
			if( !VALID(XBorderTarget_Priority[nRealIndex][i] ) )
				break;

			EntityHero* pEntity = group.GetEntityHero(XBorderTarget_Priority[nRealIndex][i]);
			if (IsTargetValid(pEntity))
			{
				m_listTarget.PushBack(pEntity);
			}
		}
	}
	*/
}

//-------------------------------------------------------------------------------
// 将要行动的敌人周围所在目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetNextAttackExplode() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	INT16 n16LoopIndx = group.GetLoopIndex();
	INT16 nRealIndex = INVALID;

	for (INT n = n16LoopIndx; n < n16LoopIndx + X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n % X_Max_Summon_Num);

		if ( IsTargetValid(pEntity))
		{
			nRealIndex = n % X_Max_Summon_Num;
			break;
		}
	}

	if( VALID(nRealIndex) )
	{
		for( INT32 i = 0; i < X_Max_Summon_Num; ++i )
		{
			if( !VALID(XExplodeTarget_Priority[nRealIndex][i] ) )
				break;

			EntityHero* pEntity = group.GetEntityHero(XExplodeTarget_Priority[nRealIndex][i]);
			if (IsTargetValid(pEntity))
			{
				m_listTarget.PushBack(pEntity);
			}
		}
	}
	*/
}

//-------------------------------------------------------------------------------
// 友方攻击力最大目标
//-------------------------------------------------------------------------------
func (s *Skill) findCasterMaxAttack() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetCamp());

	TList<EntityHero*> listTemp;

	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (!IsTargetValid(pEntity))
			continue;

		listTemp.PushBack(pEntity);
	}

	listTemp.GetList().sort(ThreatSortMaxAttack);

	INT nTargetNum = Min(listTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		m_listTarget.PushBack(listTemp.PopFront());
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方攻击力最大目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetMaxAttack() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	TList<EntityHero*> listTemp;

	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (!IsTargetValid(pEntity))
			continue;

		listTemp.PushBack(pEntity);
	}

	listTemp.GetList().sort(ThreatSortMaxAttack);

	INT nTargetNum = Min(listTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		m_listTarget.PushBack(listTemp.PopFront());
	}
	*/
}

//-------------------------------------------------------------------------------
// 混乱状态选取目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetChaos() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	TList<EntityHero*> listTemp;
	EntityHero* pEntity = NULL;

	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		pEntity = group.GetEntityHero(n);

		if (!IsTargetValid(pEntity))
			continue;

		listTemp.PushBack(pEntity);
	}

	EntityGroup& group1  = pScene->GetGroup(m_pCaster->GetCamp());
	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		pEntity = group1.GetEntityHero(n);

		if (IsTargetValid(pEntity) && pEntity != m_pCaster )
		{
			listTemp.PushBack(pEntity);
		}
	}

	if( RandPeek(pScene, pEntity, listTemp) )
	{
		m_listTarget.PushBack(pEntity);
		m_pTarget = pEntity;
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方血量最少的目标
//-------------------------------------------------------------------------------
func (s *Skill) findEnemyHPMin() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetOtherCamp());

	TList<EntityHero*> listTemp;

	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (!IsTargetValid(pEntity))
			continue;

		listTemp.PushBack(pEntity);
	}

	listTemp.GetList().sort(ThreatSortMinHP);

	INT nTargetNum = Min(listTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		m_listTarget.PushBack(listTemp.PopFront());
	}
	*/
}

//-------------------------------------------------------------------------------
// 敌方血量最少的目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetFriendRageMax() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group = pScene->GetGroup(m_pCaster->GetCamp());

	TList<EntityHero*> listTemp;

	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		EntityHero* pEntity = group.GetEntityHero(n);

		if (!IsTargetValid(pEntity))
			continue;

		listTemp.PushBack(pEntity);
	}

	listTemp.GetList().sort(ThreatSortMaxRage);

	INT nTargetNum = Min(listTemp.Size(), m_pEntry->nTargetNum);
	for (INT n = 0; n < nTargetNum; n++)
	{
		m_listTarget.PushBack(listTemp.PopFront());
	}
	*/
}

//-------------------------------------------------------------------------------
// 魅惑状态选取目标
//-------------------------------------------------------------------------------
func (s *Skill) findTargetCharm() {
	/*Scene* pScene = m_pCaster->GetScene();
	if (!VALID(pScene))
		return;

	EntityGroup& group1 = pScene->GetGroup(m_pCaster->GetCamp());
	TList<EntityHero*> listTemp;
	EntityHero* pEntity = NULL;
	for (INT n = 0; n < X_Max_Summon_Num; n++)
	{
		pEntity = group1.GetEntityHero(n);

		if (IsTargetValid(pEntity) && pEntity != m_pCaster)
		{
			listTemp.PushBack(pEntity);
		}
	}

	if (RandPeek(pScene, pEntity, listTemp))
	{
		m_listTarget.PushBack(pEntity);
		m_pTarget = pEntity;
	}
	*/
}
