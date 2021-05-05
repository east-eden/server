package scene

import (
	"bitbucket.org/funplus/server/define"
)

const (
	Camp_Max_Unit   = 50  // 每个阵营最多20个单位
	Camp_Max_Spell  = 10  // 每个阵营所属技能最多10个
	Camp_Max_Energy = 100 // 阵营符文能量最大值
)

type SceneCamp struct {
	scene        *Scene
	actionIdx    int    // 当前行动unit索引
	camp         int32  // 阵营
	aliveUnitNum int32  // 存活的单位数
	playerId     int64  // 所属玩家id
	playerLevel  int32  // 玩家等级
	playerScore  int64  // 玩家战力
	playerName   string // 玩家名字
	serverName   string // 服务器名字
	guildName    string // 工会名字
	guildId      int64  // 工会id
	portrait     int32  // 玩家头像id
	// INT32					m_nMasterIndex;							// 主角索引

	// 阵营所属技能
	energy  int32 // 能量
	spellCd []int // 技能cd

	// 所有单位

}

func NewSceneCamp(scene *Scene, camp int32) *SceneCamp {
	return &SceneCamp{
		scene:     scene,
		actionIdx: 0,
		camp:      camp,

		spellCd: make([]int, 0, Camp_Max_Spell),
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

func (c *SceneCamp) IsLoopEnd() bool {
	return c.actionIdx >= Camp_Max_Unit
}

func (c *SceneCamp) ResetLoopIndex() {
	c.actionIdx = 0
}

func (c *SceneCamp) IsValid() bool {
	return c.aliveUnitNum != 0
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

//-----------------------------------------------------------------------------
// 更新
//-----------------------------------------------------------------------------
func (c *SceneCamp) Update() {
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
