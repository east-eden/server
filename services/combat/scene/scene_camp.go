package scene

import (
	"container/list"
	"fmt"
	"sync/atomic"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	pbGlobal "bitbucket.org/funplus/server/proto/global"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

const (
	Camp_Max_Unit   = 50  // 每个阵营最多20个单位
	Camp_Max_Spell  = 10  // 每个阵营所属技能最多10个
	Camp_Max_Energy = 100 // 阵营符文能量最大值
)

type SceneCamp struct {
	scene        *Scene
	entityIdGen  int64
	entityMap    *treemap.Map // 战斗unit列表
	actionIdx    int          // 当前行动unit索引
	camp         int32        // 阵营
	aliveUnitNum int32        // 存活的单位数
	playerId     int64        // 所属玩家id
	playerLevel  int32        // 玩家等级
	playerScore  int64        // 玩家战力
	playerName   string       // 玩家名字
	serverName   string       // 服务器名字
	guildName    string       // 工会名字
	guildId      int64        // 工会id
	portrait     int32        // 玩家头像id
	// INT32					m_nMasterIndex;							// 主角索引

	// 阵营所属技能
	energy  int32 // 能量
	spellCd []int // 技能cd

	// 所有单位

	spellList *list.List // 场景内技能列表
}

func NewSceneCamp(scene *Scene, camp int32) *SceneCamp {
	return &SceneCamp{
		scene:     scene,
		entityMap: treemap.NewWith(utils.Int64Comparator),
		actionIdx: 0,
		camp:      camp,

		spellList: list.New(),
		spellCd:   make([]int, 0, Camp_Max_Spell),
	}
}

// 获取对方阵营
func (c *SceneCamp) GetOtherCamp() int32 {
	if c.camp == define.Scene_Camp_Attack {
		return define.Scene_Camp_Defence
	} else {
		return define.Scene_Camp_Attack
	}
}

// 获取战斗单位
func (c *SceneCamp) GetUnit(id int64) (*SceneEntity, bool) {
	val, ok := c.entityMap.Get(id)
	if ok {
		return val.(*SceneEntity), ok
	}

	return nil, ok
}

func (c *SceneCamp) GetUnitsLen() int {
	return c.entityMap.Size()
}

// 寻找单位
func (c *SceneCamp) FindUnitByHead() (*SceneEntity, bool) {
	if c.entityMap.Size() == 0 {
		return nil, false
	}

	return c.entityMap.Values()[0].(*SceneEntity), true
}

func (c *SceneCamp) IsLoopEnd() bool {
	return c.actionIdx >= Camp_Max_Unit
}

func (c *SceneCamp) ResetLoopIndex() {
	c.actionIdx = 0
}

func (c *SceneCamp) IsValid() bool {
	return c.aliveUnitNum != 0
}

// 战斗单位死亡
func (c *SceneCamp) OnUnitDead(u *SceneEntity) {
	c.aliveUnitNum--
	c.scene.OnUnitDead(u)
}

// 战斗单位消亡
func (c *SceneCamp) OnUnitDisappear(u *SceneEntity) {

}

func (c *SceneCamp) addSpell(opts ...SpellOption) {
	spell := NewSkill()
	spell.Init(opts...)
	c.spellList.PushBack(spell)
}

func (s *SceneCamp) AddEntityByPB(unitInfo *pbGlobal.EntityInfo) error {
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

func (s *SceneCamp) AddEntityByOptions(opts ...EntityOption) error {
	id := atomic.AddInt64(&s.entityIdGen, 1)
	e, err := NewSceneEntity(id, opts...)
	if err != nil {
		return err
	}

	s.entityMap.Put(id, e)
	return nil
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

//-----------------------------------------------------------------------------
// 更新
//-----------------------------------------------------------------------------
func (c *SceneCamp) Update() {
	c.updateUnits()
	c.updateSpells()
}

// 更新阵营内技能
func (c *SceneCamp) updateSpells() {
	var next *list.Element
	for e := c.spellList.Front(); e != nil; e = next {
		next = e.Next()

		s := e.Value.(*Skill)
		s.Update()

		// 删除已作用玩的技能
		if s.IsCompleted() {
			c.spellList.Remove(e)
		}
	}
}

// 更新阵营内单位
func (c *SceneCamp) updateUnits() {
	it := c.entityMap.Iterator()
	for it.Next() {
		it.Value().(*SceneEntity).Update()
	}
}

//-----------------------------------------------------------------------------
// 清空所有单位
//-----------------------------------------------------------------------------
func (c *SceneCamp) ClearUnit() {
	c.entityMap.Clear()
}

//-----------------------------------------------------------------------------
// 释放符文技能
//-----------------------------------------------------------------------------
func (c *SceneCamp) CastCampSpell() {
	// RuneSet::iterator it = m_setRune.begin();

	// INT32 nRuneIndex = (*it) / 10000;

	// // 判断能量是否足够
	// if( VALID(m_pRuneSpellEntry[nRuneIndex]) && m_nEnergy > m_pRuneSpellEntry[nRuneIndex]->nEnergyCost)
	// {
	// 	ModeAttEnergy(-(m_pRuneSpellEntry[nRuneIndex]->nEnergyCost));

	// 	// 释放技能
	// 	if( VALID(m_ArrayHero[m_nMasterIndex]) )
	// 	{
	// 		EntityGroup& group = GetScene()->GetGroup(GetOtherCamp());
	// 		EntityHero* pTarget = FindTargetByPriority(m_nMasterIndex, &group, FALSE);

	// 		m_ArrayHero[m_nMasterIndex]->CastRuneSpell(m_pRuneSpellEntry[nRuneIndex], pTarget, m_n8RuneLevel[nRuneIndex]);
	// 	}

	// 	m_n8RuneWeight[nRuneIndex]+= m_pRuneSpellEntry[nRuneIndex]->nRuneCD;
	// 	m_n8RuneCD[nRuneIndex] += m_pRuneSpellEntry[nRuneIndex]->nRuneCD;

	// 	m_setRune.erase(it);
	// }
}

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
// 攻击
//-----------------------------------------------------------------------------
func (c *SceneCamp) Attack(dst *SceneCamp) {
	// compile comment
	// EntityGroup* pTarget = static_cast<EntityGroup*>(pEntity);
	// BOOL bBreak = FALSE;
	// for( INT32 i = m_n16LoopIndex; i < X_Max_Summon_Num; ++i )
	// {
	// 	++m_n16LoopIndex;

	// 	if( VALID(m_ArrayHero[i]) && m_ArrayHero[i]->IsValid() )
	// 	{
	// 		EntityHero* pHero = FindTargetByPriority(i, pTarget, TRUE);

	// 		if( VALID(pHero) )
	// 		{
	// 			m_ArrayHero[i]->Attack(pHero);
	// 			m_ArrayHero[i]->GetCombatController().CalAuraEffect(GetScene()->GetCurRound());

	// 			// 风怒状态
	// 			if( m_ArrayHero[i]->HasState(EHS_Anger) )
	// 			{
	// 				EntityHero* pHero = FindTargetByPriority(i, pTarget, TRUE);
	// 				if( VALID(pHero) )
	// 				{
	// 					m_ArrayHero[i]->Attack(pHero);
	// 				}
	// 			}

	// 			AddAttackNum();
	// 			bBreak = TRUE;
	// 		}
	// 	}

	// 	if( bBreak )
	// 		break;
	// }
}

//-----------------------------------------------------------------------------
// 改变阵营内符文能量
//-----------------------------------------------------------------------------
func (c *SceneCamp) ModAttEnergy(mod int32) {
	c.energy += mod
	if c.energy > Camp_Max_Energy {
		c.energy = Camp_Max_Energy
	}
}

//-----------------------------------------------------------------------------
// 战斗开始时触发
//-----------------------------------------------------------------------------
func (c *SceneCamp) TriggerByStartBehaviour() {
	it := c.entityMap.Iterator()
	for it.Next() {
		u := it.Value().(*SceneEntity)
		u.opts.CombatCtrl.TriggerByBehaviour(define.BehaviourType_Start, u, 0, 0, define.SpellType_Null)
	}
}

//-----------------------------------------------------------------------------
// 计算帮会和符文产生的伤害改变属性
//-----------------------------------------------------------------------------
func (c *SceneCamp) CalDmgModAtt() {
	// m_nDmgModAtt[EDM_RaceDoneKindom] += pPlayer->GetScienceSkillValue(ESCS_DoneKindom);
	// m_nDmgModAtt[EDM_RaceTakenKindom] -= pPlayer->GetScienceSkillValue(ESCS_TakenKindom);
	// m_nDmgModAtt[EDM_RaceDoneHell] += pPlayer->GetScienceSkillValue(ESCS_DoneHell);
	// m_nDmgModAtt[EDM_RaceTakenHell] -= pPlayer->GetScienceSkillValue(ESCS_TakenHell);
	// m_nDmgModAtt[EDM_RaceDoneForest] += pPlayer->GetScienceSkillValue(ESCS_DoneForest);
	// m_nDmgModAtt[EDM_RaceTakenForest] -= pPlayer->GetScienceSkillValue(ESCS_TakenForest);
	// m_nDmgModAtt[EDM_RaceDoneWild] += pPlayer->GetScienceSkillValue(ESCS_DoneWild);
	// m_nDmgModAtt[EDM_RaceTakenWild] -= pPlayer->GetScienceSkillValue(ESCS_TakenWild);
	// m_nDmgModAtt[EDM_RaceDoneOther] += pPlayer->GetScienceSkillValue(ESCS_DoneForest);
	// m_nDmgModAtt[EDM_RaceTakenOther] -= pPlayer->GetScienceSkillValue(ESCS_TakenForest);
}

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
