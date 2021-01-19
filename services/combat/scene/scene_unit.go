package scene

import (
	"e.coding.net/mmstudio/blade/server/define"
	"e.coding.net/mmstudio/blade/server/excel/auto"
	"e.coding.net/mmstudio/blade/server/utils"
	log "github.com/rs/zerolog/log"
	"github.com/willf/bitset"
)

const (
	Unit_Energy_OnBeDamaged = 2 // 受伤害增加能量
	Unit_Init_AuraNum       = 3 // 初始化buff数量
)

type SceneUnit struct {
	id            int64
	level         uint32
	posX          int16                                       // x坐标
	posY          int16                                       // y坐标
	TauntId       int64                                       // 被嘲讽目标
	v2            define.Vector2                              // 朝向
	scene         *Scene                                      // 场景
	camp          *SceneCamp                                  // 场景阵营
	normalSpell   *define.SpellEntry                          // 普攻技能
	specialSpell  *define.SpellEntry                          // 特殊技能
	passiveSpells [define.Spell_PassiveNum]*define.SpellEntry // 被动技能列表

	// 伤害统计
	totalDmgRecv int64 // 总共受到的伤害
	totalDmgDone int64 // 总共造成的伤害
	totalHeal    int64 // 总共产生治疗
	attackNum    int   // 攻击次数

	opts *UnitOptions
}

func NewSceneUnit(id int64, opts ...UnitOption) *SceneUnit {
	u := &SceneUnit{
		opts: DefaultUnitOptions(),
	}

	for _, o := range opts {
		o(u.opts)
	}

	return u
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

func (s *SceneUnit) ActionCtrl() *ActionCtrl {
	return s.opts.ActionCtrl
}

func (s *SceneUnit) CombatCtrl() *CombatCtrl {
	return s.opts.CombatCtrl
}

func (s *SceneUnit) MoveCtrl() *MoveCtrl {
	return s.opts.MoveCtrl
}

func (s *SceneUnit) Opts() *UnitOptions {
	return s.opts
}

func (s *SceneUnit) Update() {
	log.Info().
		Int64("id", s.id).
		Int32("type_id", s.opts.TypeId).
		Int32("pos_x", s.opts.PosX).
		Int32("pos_y", s.opts.PosY).
		Msg("creature start UpdateSpell")

	s.ActionCtrl().Update()
	s.CombatCtrl().Update()
	s.MoveCtrl().Update()
}

func (s *SceneUnit) HasState(e define.EHeroState) bool {
	return s.opts.State.Test(uint(e))
}

func (s *SceneUnit) HasStateAny(flag uint32) bool {
	compare := utils.FromCountableBitset([]uint64{uint64(flag)}, []int16{})
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
	s.AddState(define.HeroState_Dead, 1)
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
	heroEntry, ok := auto.GetHeroEntry(heroInfo.TypeId)
	if !ok {
		log.Warn().Int32("type_id", heroInfo.TypeId).Msg("cannot find hero entry")
		return
	}

	s.opts.AttManager.SetBaseAttId(int32(heroEntry.AttID))
	s.opts.AttManager.CalcAtt()
	s.opts.AttManager.SetAttValue(define.Att_Plus_CurHP, s.opts.AttManager.GetAttValue(define.Att_Plus_MaxHP))
}

//-----------------------------------------------------------------------------
// 技能初始化
//-----------------------------------------------------------------------------
func (s *SceneUnit) initSpell() {
	// todo 设置初始技能
	// s.normalSpell = auto.GetSpellEntry(s.opts.Entry.NormalSpellId)
	// s.specialSpell = auto.GetSpellEntry(s.opts.Entry.SpecialSpellId)

	// 被动技能
	for n := 0; n < define.Spell_PassiveNum; n++ {
		passiveSpellEntry := s.passiveSpells[n]
		if passiveSpellEntry == nil {
			continue
		}

		err := s.opts.CombatCtrl.CastSpell(passiveSpellEntry, s, s, false)
		utils.ErrPrint(err, "InitSpell failed", passiveSpellEntry.ID, s.opts.TypeId)
	}
}

//-----------------------------------------------------------------------------
// 初始化被动技能
//-----------------------------------------------------------------------------
func (s *SceneUnit) initAura() {
	// 增加初始被动Aura
	for n := 0; n < Unit_Init_AuraNum; n++ {
		// todo
		// if s.opts.Entry.PassiveAuraId[n] == -1 {
		// 	continue
		// }

		// s.opts.CombatCtrl.AddAura(s.opts.Entry.PassiveAuraId[n], s, 0, 0, define.SpellType_Null, 0, 1)
	}
}

//-----------------------------------------------------------------------------
// 设置普通攻击
//-----------------------------------------------------------------------------
func (s *SceneUnit) SetNormalSpell(spellId uint32) {
	// todo
	// spellEntry, ok := auto.GetSpellEntry(spellId)
	// if !ok {
	// 	return
	// }

	// s.normalSpell = spellEntry
}

//-------------------------------------------------------------------------------
// 状态
//-------------------------------------------------------------------------------
func (s *SceneUnit) AddState(state define.EHeroState, count int16) {
	new := !s.HasState(state)

	s.opts.State.Set(uint(state), count)

	// todo 进入新状态处理
	if new {
		// Scene* pScene = GetScene();
		// if (VALID(pScene) && !pScene->IsOnlyRecord() )
		// {
		// 	CreateSceneProtoMsg(msg, MS_SetState,);
		// 	*msg << (UINT32)GetLocation();
		// 	*msg << (UINT32)eState;
		// 	pScene->AddMsgList(msg);
		// }

		// 追加状态处理
		s.AddToState(state)
	}
}

func (s *SceneUnit) DecState(state define.EHeroState, count int16) {
	if !s.HasState(state) {
		return
	}

	s.opts.State.Clear(uint(state), count)

	// todo 退出状态处理
	if !s.HasState(state) {
		// Scene* pScene = GetScene();
		// if (VALID(pScene) && !pScene->IsOnlyRecord() )
		// {
		// 	CreateSceneProtoMsg(msg, MS_UnsetState, );
		// 	*msg << (UINT32)GetLocation();
		// 	*msg << (UINT32)eState;
		// 	pScene->AddMsgList(msg);
		// }

		s.EscFromState(state)
	}
}

//-------------------------------------------------------------------------------
// todo 保存录像
//-------------------------------------------------------------------------------
func (s *SceneUnit) Save2DB(pRecord interface{}) {
	// pRecord->dwEntityID = m_pEntry->dwTypeID;
	// pRecord->nFashionID = m_nFashionID;
	// pRecord->dwMountTypeID = m_dwMountTypeID;
	// pRecord->nStateFlag = m_n16HeroState;
	// pRecord->nFlyUp = m_nFly_Up;
	// pRecord->nLevel = m_nLevel;
	// pRecord->nRageLevel = m_n16RageLevel;
	// pRecord->nStarLevel = m_nStar;
	// pRecord->nQuality = m_nQuality;
	// memcpy(pRecord->nAtt, m_AttRecord.ExportAtt(), sizeof(pRecord->nAtt));
	// memcpy(pRecord->nBaseAtt, m_AttRecord.ExportBaseAtt(), sizeof(pRecord->nBaseAtt));
	// memcpy(pRecord->nBaseAttModPct, m_AttRecord.ExportBaseAttModPct(), sizeof(pRecord->nBaseAttModPct));
	// memcpy(pRecord->nAttMod, m_AttRecord.ExportAttMod(), sizeof(pRecord->nAttMod));
	// memcpy(pRecord->nAttModPct, m_AttRecord.ExportAttModPct(), sizeof(pRecord->nAttModPct));
	// memcpy(pRecord->dwPassiveSpell, m_AttRecord.ExportPassiveSpell(), sizeof(pRecord->dwPassiveSpell));
}

//-----------------------------------------------------------------------------
// todo 初始化伤害加成
//-----------------------------------------------------------------------------
func (s *SceneUnit) InitHeroDmgModAtt(info *define.HeroInfo, pos int32) {
	// ZeroMemory(m_nHeroDmgModAtt,sizeof(m_nHeroDmgModAtt));
}
