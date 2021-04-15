package scene

import (
	"errors"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/utils"
	log "github.com/rs/zerolog/log"
	"github.com/willf/bitset"
)

const (
	Unit_Energy_OnBeDamaged = 2 // 受伤害增加能量
	Unit_Init_AuraNum       = 3 // 初始化buff数量
)

var (
	ErrEntityInvalidHeroEntry = errors.New("invalid hero entry")
)

type SceneEntity struct {
	id            int64
	level         uint32
	posX          int16                                       // x坐标
	posZ          int16                                       // y坐标
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

	opts *EntityOptions
}

func NewSceneEntity(id int64, opts ...EntityOption) (*SceneEntity, error) {
	e := &SceneEntity{
		opts: DefaultUnitOptions(),
	}

	for _, o := range opts {
		o(e.opts)
	}

	if e.opts.Entry == nil {
		return nil, ErrEntityInvalidHeroEntry
	}

	e.opts.AttManager.SetBaseAttId(e.opts.Entry.AttId)
	e.opts.AttManager.CalcAtt()

	return e, nil
}

func (s *SceneEntity) Guid() int64 {
	return s.id
}

func (s *SceneEntity) GetLevel() uint32 {
	return s.level
}

func (s *SceneEntity) GetScene() *Scene {
	return s.opts.Scene
}

func (s *SceneEntity) GetCamp() int32 {
	return 0
}

func (s *SceneEntity) ActionCtrl() *ActionCtrl {
	return s.opts.ActionCtrl
}

func (s *SceneEntity) CombatCtrl() *CombatCtrl {
	return s.opts.CombatCtrl
}

func (s *SceneEntity) MoveCtrl() *MoveCtrl {
	return s.opts.MoveCtrl
}

func (s *SceneEntity) Opts() *EntityOptions {
	return s.opts
}

func (s *SceneEntity) Update() {
	log.Info().
		Int64("id", s.id).
		Int32("type_id", s.opts.TypeId).
		Int32("pos_x", s.opts.PosX).
		Int32("pos_z", s.opts.PosZ).
		Msg("creature start UpdateSpell")

	s.ActionCtrl().Update()
	s.CombatCtrl().Update()
	s.MoveCtrl().Update()
}

func (s *SceneEntity) HasState(e define.EHeroState) bool {
	return s.opts.State.Test(uint(e))
}

func (s *SceneEntity) HasStateAny(flag uint32) bool {
	compare := utils.FromCountableBitset([]uint64{uint64(flag)}, []int16{})
	return s.opts.State.Intersection(compare).Any()
}

func (s *SceneEntity) GetState64() uint64 {
	return s.opts.State.Bytes()[0]
}

func (s *SceneEntity) HasImmunityAny(tp define.EImmunityType, flag uint32) bool {
	compare := bitset.From([]uint64{uint64(flag)})
	return s.opts.Immunity[tp].Intersection(compare).Any()
}

//-----------------------------------------------------------------------------
// 初始化
//-----------------------------------------------------------------------------
func (s *SceneEntity) InitByScene(scene *Scene, camp *SceneCamp, posX, posY int16, v2 define.Vector2) {
	s.posX = posX
	s.posZ = posY
	s.scene = scene
	s.camp = camp

	s.initSpell()
	s.initAura()
}

//-----------------------------------------------------------------------------
// 进攻
//-----------------------------------------------------------------------------
func (s *SceneEntity) Attack(target *SceneEntity) {
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
func (s *SceneEntity) BeatBack(target *SceneEntity) {
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
func (s *SceneEntity) OnDead(caster *SceneEntity, spellId int32) {
	if s.HasState(define.HeroState_Dead) {
		return
	}

	s.camp.OnUnitDead(s)

	// 清空当前值
	s.opts.AttManager.SetBaseAttValue(define.Att_CurHP, 0)

	// 设置为死亡状态
	s.AddState(define.HeroState_Dead, 1)
}

//-----------------------------------------------------------------------------
// 造成伤害
//-----------------------------------------------------------------------------
func (s *SceneEntity) OnDamage(target *SceneEntity, dmgInfo *CalcDamageInfo) {
	s.opts.CombatCtrl.TriggerByDmgMod(true, target, dmgInfo)
}

//-----------------------------------------------------------------------------
// 改变符文能量
//-----------------------------------------------------------------------------
func (s *SceneEntity) ModAttEnergy(mod int32) {
	s.camp.ModAttEnergy(mod)
}

//-----------------------------------------------------------------------------
// 承受伤害
//-----------------------------------------------------------------------------
func (s *SceneEntity) OnBeDamaged(caster *SceneEntity, dmgInfo *CalcDamageInfo) {
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
func (s *SceneEntity) DoneDamage(caster *SceneEntity, dmgInfo *CalcDamageInfo) {
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
			if int64(s.opts.AttManager.GetBaseAttValue(define.Att_CurHP)) <= dmgInfo.Damage {
				dmgInfo.Damage = int64(s.opts.AttManager.GetBaseAttValue(define.Att_CurHP) - 1)
				s.opts.AttManager.SetBaseAttValue(define.Att_CurHP, 1)

				// 伤害统计
				s.totalDmgRecv += dmgInfo.Damage
				caster.totalDmgDone += dmgInfo.Damage

				dmgInfo.ProcTarget |= (1 << define.AuraEvent_Taken_Any_Damage)
				dmgInfo.ProcEx |= (1 << define.AuraEventEx_UnDead)
			} else {
				// 伤害统计
				s.totalDmgRecv += dmgInfo.Damage
				caster.totalDmgDone += dmgInfo.Damage

				s.opts.AttManager.ModBaseAttValue(define.Att_CurHP, int32(-dmgInfo.Damage))
			}
		} else {
			// 伤害统计
			s.totalDmgRecv += dmgInfo.Damage
			caster.totalDmgDone += dmgInfo.Damage

			s.opts.AttManager.ModBaseAttValue(define.Att_CurHP, int32(-dmgInfo.Damage))

			if s.opts.AttManager.GetBaseAttValue(define.Att_CurHP) <= 0 {
				// 刚刚死亡
				s.OnDead(caster, dmgInfo.SpellId)
			}
		}

		// 治疗
	case define.DmgInfo_Heal:
		s.opts.AttManager.ModBaseAttValue(define.Att_CurHP, int32(dmgInfo.Damage))

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
func (s *SceneEntity) AddToState(state define.EHeroState) {
	s.opts.CombatCtrl.TriggerByServentState(state, true)
}

//-----------------------------------------------------------------------------
// 脱离状态
//-----------------------------------------------------------------------------
func (s *SceneEntity) EscFromState(state define.EHeroState) {
	s.opts.CombatCtrl.TriggerByServentState(state, false)
}

//-----------------------------------------------------------------------------
// 免疫
//-----------------------------------------------------------------------------
func (s *SceneEntity) AddToImmunity(immunityType define.EImmunityType, immunity int) {
	switch immunityType {
	case define.ImmunityType_Mechanic:
		// 删除指定类型的Aura
		//s.opts.CombatCtrl.RemoveAura(immunity)
	}
}

//-----------------------------------------------------------------------------
// 初始化伤害加成
//-----------------------------------------------------------------------------
func (s *SceneEntity) InitDmgModAtt() {
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
func (s *SceneEntity) InitAttribute(heroInfo *define.HeroInfo) {
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

	s.opts.AttManager.SetBaseAttId(int32(heroEntry.AttId))
	s.opts.AttManager.CalcAtt()
	s.opts.AttManager.SetBaseAttValue(define.Att_CurHP, s.opts.AttManager.GetBaseAttValue(define.Att_MaxHP))
}

//-----------------------------------------------------------------------------
// 技能初始化
//-----------------------------------------------------------------------------
func (s *SceneEntity) initSpell() {
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
func (s *SceneEntity) initAura() {
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
func (s *SceneEntity) SetNormalSpell(spellId uint32) {
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
func (s *SceneEntity) AddState(state define.EHeroState, count int16) {
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

func (s *SceneEntity) DecState(state define.EHeroState, count int16) {
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
func (s *SceneEntity) Save2DB(pRecord interface{}) {
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
func (s *SceneEntity) InitHeroDmgModAtt(info *define.HeroInfo, pos int32) {
	// ZeroMemory(m_nHeroDmgModAtt,sizeof(m_nHeroDmgModAtt));
}
