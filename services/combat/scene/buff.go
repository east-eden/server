package scene

import (
	"math"

	"bitbucket.org/funplus/server/define"
)

type Buff struct {
	opts *BuffOptions

	EffectTimes uint8 // 剩余作用次数
	CurDuration uint8 // 当前剩余回合数

	AddEffLock int16                  // 追加时效果锁
	removeMode define.EAuraRemoveMode // 删除标志
	SpellType  define.ESpellType

	Multiple     [define.SpellEffectNum]float32
	CurPoint     [define.SpellEffectNum]int32
	EffMisc      [define.SpellEffectNum]int32  // 效果变量
	TriggerCount [define.SpellEffectNum]uint32 // 每回合已触发次数
	UpdateTime   int32
	TriggerCd    [define.SpellEffectNum]int32 // 触发CD
}

//-------------------------------------------------------------------------------
// 初始化
//-------------------------------------------------------------------------------
func (b *Buff) Init(opts ...AuraOption) {
	b.opts = DefaultBuffOptions()

	for _, o := range opts {
		o(b.opts)
	}

	b.removeMode = define.AuraRemoveMode_Null
	b.AddEffLock = 0
	b.EffectTimes = uint8(b.opts.Entry.EffectTimes)
	b.CurDuration = uint8(b.opts.Entry.Duration)
	b.UpdateTime = 0

	for n := 0; n < define.SpellEffectNum; n++ {
		b.EffMisc[n] = 0
		b.CurPoint[n] = 0
		b.Multiple[n] = float32(b.opts.Entry.Multiple[n])
		b.TriggerCount[n] = b.opts.Entry.TriggerCount[n]
		b.TriggerCd[n] = 0
	}

	if b.opts.Entry.DecByTarget {
		b.opts.Owner.CombatCtrl().CalDecByTargetPoint(&b.opts.Entry.SpellBase, b.CurPoint[:], b.Multiple[:], b.opts.Level)
	} else {
		if b.opts.Caster != nil {
			b.opts.Caster.CombatCtrl().CalSpellPoint(&b.opts.Entry.SpellBase, b.CurPoint[:], b.Multiple[:], b.opts.Level)
		} else {
			b.opts.Owner.CombatCtrl().CalSpellPoint(&b.opts.Entry.SpellBase, b.CurPoint[:], b.Multiple[:], b.opts.Level)
		}
	}

	for n := 0; n < define.SpellEffectNum; n++ {
		b.CurPoint[n] *= int32(b.opts.CurWrapTimes)
	}
}

func (a *Buff) Opts() *BuffOptions {
	return a.opts
}

func (a *Buff) lockApply() {
	if a.AddEffLock >= 0 {
		a.AddEffLock++
	}
}

func (a *Buff) unlockApply() {
	if a.AddEffLock > 0 {
		a.AddEffLock--
	}
}

func (a *Buff) IsApplyLocked() bool {
	return a.AddEffLock > 0
}

func (a *Buff) InvalidApplyLock() {
	a.AddEffLock = -1
}

func (a *Buff) IsApplyLockValid() bool {
	return a.AddEffLock != -1
}

func (a *Buff) IsHangup() bool {
	return a.removeMode&define.AuraRemoveMode_Hangup != 0
}

func (a *Buff) IsRemoved() bool {
	return a.removeMode&define.AuraRemoveMode_Removed != 0
}

func (a *Buff) GetRemoveMode() define.EAuraRemoveMode {
	return a.removeMode
}

func (a *Buff) AddRemoveMode(mode define.EAuraRemoveMode) {
	a.removeMode |= mode
}

func (a *Buff) DecRemoveMode(mode define.EAuraRemoveMode) {
	a.removeMode &= ^mode
}

func (a *Buff) isNoCd(index int32) bool {
	if a.opts.Entry.TriggerCd[index] == 0 {
		return true
	}

	return a.TriggerCd[index] < 0
}

//-------------------------------------------------------------------------------
// 计算效果
//-------------------------------------------------------------------------------
func (a *Buff) CalcApplyEffect(register bool, sync bool) {
	a.lockApply()
	a.CalAuraEffect(define.AuraEffectStep_Apply, -1, nil, nil)
	a.unlockApply()

	a.AddRemoveMode(define.AuraRemoveMode_Running)

	if !a.IsRemoved() {
		if register {
			a.opts.Owner.CombatCtrl().RegisterAura(a)
		}

		if sync {
			a.SyncAuraToClient(define.AuraSyncStep_Add)
		}
	}
}

func (a *Buff) CalAuraEffect(step define.EAuraEffectStep, effIndex int32, param1 interface{}, param2 interface{}) define.EAuraAddResult {
	result := define.AuraAddResult_Null

	if effIndex >= 0 {
		if a.TriggerCount[effIndex] != 0 && a.isNoCd(effIndex) {
			eff := a.opts.Entry.Effects[effIndex]
			result = auraEffectsHandlers[eff](a, step, effIndex, param1, param2)
			a.TriggerCount[effIndex]--
			a.TriggerCd[effIndex] = a.opts.Entry.TriggerCd[effIndex]

			// 作用次数更新
			if step == define.AuraEffectStep_Effect && result == define.AuraAddResult_Success {
				a.consume()
			}
		}
	} else {
		for index := 0; index < define.SpellEffectNum; index++ {
			if step == define.AuraEffectStep_Effect && a.opts.Entry.TriggerId[index] > 0 {
				continue
			}

			removeMode := *param1.(*define.EAuraRemoveMode)
			if step == define.AuraEffectStep_Remove && (removeMode&define.EAuraRemoveMode(a.opts.Entry.RemoveEffect[index]) == 0) {
				continue
			}

			eff := a.opts.Entry.Effects[index]
			if eff == define.AuraEffectType_Null {
				continue
			}

			result = auraEffectsHandlers[eff](a, step, effIndex, param1, param2)
		}
	}

	return result
}

//-------------------------------------------------------------------------------
// 同步客户端
//-------------------------------------------------------------------------------
func (a *Buff) SyncAuraToClient(step define.EAuraSyncStep) {
	// 有动作或者特效就九宫格群发(如果是玩家，再发给小队)
	if a.opts.Entry.HaveVisual {
		scene := a.opts.Owner.GetScene()
		if scene == nil {
			return
		}

		switch step {
		case define.AuraSyncStep_Add:
			if !scene.IsOnlyRecord() {
				// todo send
				//CreateSceneProtoMsg(msg, MS_AddAura,);
				//*msg << (UINT32)(VALID(m_pCaster) ? m_pCaster->GetLocation() : INVALID);
				//*msg << (UINT32)m_pOwner->GetLocation();
				//*msg << (UINT32)GetAuraID();
				//pScene->AddMsgList(msg);
			}

		case define.AuraSyncStep_Remove:
			if !scene.IsOnlyRecord() {
				//CreateSceneProtoMsg(msg, MS_RemoveAura,);
				//*msg << (UINT32)(VALID(m_pCaster) ? m_pCaster->GetLocation() : INVALID);
				//*msg << (UINT32)m_pOwner->GetLocation();
				//*msg << (UINT32)GetAuraID();
				//pScene->AddMsgList(msg);
			}
		}
	}
}

//-------------------------------------------------------------------------------
// 叠加结果
//-------------------------------------------------------------------------------
func (a *Buff) CheckWrapResult(aura *Buff) define.EAuraWrapResult {
	if aura == nil {
		return define.AuraWrapResult_Invalid
	}

	// 相同Aura && 相同Caster
	if a.opts.Entry.ID == aura.opts.Entry.ID {
		if a.opts.Caster == aura.opts.Caster {
			return define.AuraWrapResult_Wrap
		} else {
			if a.opts.Entry.MultiWrap {
				return define.AuraWrapResult_Wrap
			} else {
				return define.AuraWrapResult_Add
			}
		}
	}

	// 同类Aura
	if define.SpellType(a.opts.Entry.ID) == define.SpellType(aura.opts.Entry.ID) {

		// 同一个人
		if a.opts.Caster == aura.opts.Caster || a.opts.Entry.MultiWrap {
			if aura.MorePowerfulThan(a) {
				return define.AuraWrapResult_Replace
			} else {
				return define.AuraWrapResult_Invalid
			}
		} else {
			return define.AuraWrapResult_Add
		}
	}

	// 不同类Aura
	return define.AuraWrapResult_Add
}

//-------------------------------------------------------------------------------
// 等级比较
//-------------------------------------------------------------------------------
func (a *Buff) MorePowerfulThan(aura *Buff) bool {

	if define.SpellType(a.opts.Entry.ID) == define.SpellType(aura.opts.Entry.ID) {
		return define.SpellLevel(a.opts.Entry.ID) >= define.SpellLevel(aura.opts.Entry.ID)
	} else {
		return math.Abs(float64(a.opts.Entry.EffectPriority)) >= math.Abs(float64(aura.opts.Entry.EffectPriority))
	}
}

//-------------------------------------------------------------------------------
// 消耗作用次数
//-------------------------------------------------------------------------------
func (a *Buff) consume() {
	if a.opts.Entry.AuraCastType != define.BuffCasting_Times {
		return
	}

	a.EffectTimes--
	if a.EffectTimes <= 0 {
		a.opts.Owner.CombatCtrl().RemoveAura(a, define.AuraRemoveMode_Consume)
	}
}

//-------------------------------------------------------------------------------
// 回合结束
//-------------------------------------------------------------------------------
func (a *Buff) RoundEnd() {
	a.UpdateTime++
	a.CurDuration--

	if a.CurDuration <= 0 {
		if a.opts.Entry.Duration > 0 {
			a.opts.Owner.CombatCtrl().RemoveAura(a, define.AuraRemoveMode_Default)
			return
		}
	}

	for index := 0; index < define.SpellEffectNum; index++ {
		if a.opts.Entry.TriggerId[index] > 0 {
			if a.UpdateTime%6 == 0 {
				// 重置每回合已触发次数
				a.TriggerCount[index] = a.opts.Entry.TriggerCount[index]
			}
		}

		if a.opts.Entry.TriggerCd[index] != 0 {
			a.TriggerCd[index]--
		}
	}
}

//-------------------------------------------------------------------------------
// 驱散
//-------------------------------------------------------------------------------
func (a *Buff) Disperse() {
	if a.IsApplyLocked() {
		return
	}

	a.opts.Owner.CombatCtrl().RemoveAura(a, define.AuraRemoveMode_Dispel)
}

//-------------------------------------------------------------------------------
// 强化Aura作用时间
//-------------------------------------------------------------------------------
func (a *Buff) ModDuration(modDuration uint32) {
	a.CurDuration += uint8(modDuration)
}

//-------------------------------------------------------------------------------
// 计算伤害
//-------------------------------------------------------------------------------
func (a *Buff) CalDamage(baseDamage int64, damageInfo *CalcDamageInfo, target *SceneEntity) {
	if a.opts.SpellType == define.SpellType_Rune {
		damageInfo.Damage = baseDamage
		return
	}

	casterAttManager := a.opts.Caster.Opts().AttManager
	targetAttManager := target.Opts().AttManager
	baseDamage += int64(casterAttManager.GetAttValue(define.Att_DmgInc)) - int64(targetAttManager.GetAttValue(define.Att_DmgDec))

	if a.opts.SpellType == define.SpellType_Rage {
		dmgMod := float64(a.opts.RagePctMod) * float64(baseDamage)
		baseDamage += int64(dmgMod)
	}

	// todo
	damageInfo.Damage = baseDamage

	//FLOAT fPctDmgMod = (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PctDmgInc) - pTarget->GetAttController().GetAttValue(EHA_PctDmgDec));

	//if( m_pCaster->GetScene()->GetStateFlag() & ESSF_PVP )
	//fPctDmgMod += (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PVPPctDmgInc) - pTarget->GetAttController().GetAttValue(EHA_PVPPctDmgDec));

	//// 伤害类型加成
	//if (DamageInfo.eSchool == EIS_Physics)
	//{
	//fPctDmgMod = (FLOAT)(m_pCaster->GetDmgModAtt(EDM_DamageDonePctPhysics) + pTarget->GetDmgModAtt(EDM_DamageTakenPctPhysics));
	//}
	//else if(DamageInfo.eSchool == EIS_Magic)
	//{
	//fPctDmgMod = (FLOAT)(m_pCaster->GetDmgModAtt(EDM_DamageDonePctMagic) + pTarget->GetDmgModAtt(EDM_DamageTakenPctMagic));
	//}

	//// 伤害种族加成
	//fPctDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_RaceDoneKindom + pTarget->GetEntry()->eRace));
	//fPctDmgMod += (FLOAT)(pTarget->GetDmgModAtt(EDM_RaceTakenKindom + m_pCaster->GetEntry()->eRace));

	//fPctDmgMod += m_pCaster->GetScene()->GetSceneDmgMod();
	//fPctDmgMod += m_pCaster->GetScene()->GetLevelSuppress(m_pCaster, pTarget);
	//// 英雄势力伤害加成
	//fPctDmgMod += (FLOAT)m_pCaster->GetDmgModAtt(pTarget->GetDmgModeByForce());
	//// 判断百分比下限
	//if( fPctDmgMod < -7000.0f )
	//fPctDmgMod = -7000.0f;

	//nBaseDamage += ((fPctDmgMod / 10000.0f) * (FLOAT)nBaseDamage);

	//if (nBaseDamage < 1)
	//{
	//DamageInfo.nDamage = 1;
	//}

	//if (DamageInfo.dwProcEx & EAEE_Critical_Hit)
	//{
	//INT nCrit = m_pCaster->GetAttController().GetAttValue(EHA_CritInc) - pTarget->GetAttController().GetAttValue(EHA_CritDec);
	//nBaseDamage *= (Max(10000, 17500 + nCrit) /10000.0f);
	//}

	//INT32 nMinDmg = (INT32)((FLOAT)m_pCaster->GetAttController().GetAttValue(EHA_AttackPower) * 0.05f);
	//DamageInfo.nDamage = Max(nMinDmg, nBaseDamage);
}

//-------------------------------------------------------------------------------
// 计算治疗
//-------------------------------------------------------------------------------
func (a *Buff) CalHeal(baseHeal int32, damageInfo *CalcDamageInfo, target *SceneEntity) {
	// 重伤状态无法加血
	if target.HasState(define.HeroState_Injury) {
		damageInfo.Damage = 0
		return
	}

	// todo
	// 中毒状态加血效果减75%
	//FLOAT fHealPct		= pTarget->HasState(EHS_Poison) ? 0.25f : 1.0f;

	//if( m_eSpellType == ERMT_Rune || m_eSpellType == ERMT_Pet)
	//{
	//DamageInfo.nDamage = nBaseHeal * fHealPct;
	//return;
	//}

	//if( m_eSpellType == ERMT_Rage )
	//{
	//FLOAT fDmgMod = m_fRagePctMod * (FLOAT)(nBaseHeal);
	//nBaseHeal += fDmgMod;
	//}

	//FLOAT fDmgMod = (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PctDmgInc));

	//if( m_pCaster->GetScene()->GetStateFlag() & ESSF_PVP )
	//fDmgMod += (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_PVPPctDmgInc));

	//fDmgMod += (FLOAT)(m_pCaster->GetDmgModAtt(EDM_DamageDonePctHeal) + pTarget->GetDmgModAtt(EDM_DamageTakenPcttHeal));
	////	fDmgMod += (FLOAT)(m_pCaster->GetAttController().GetAttValue(EHA_HealPctIncDone) + pTarget->GetAttController().GetAttValue(EHA_HealPctIncTaken));
	//fDmgMod = fDmgMod / 10000.0f;
	//fDmgMod = fDmgMod * (FLOAT)nBaseHeal;
	//nBaseHeal += fDmgMod;

	//if (DamageInfo.dwProcEx & EAEE_Critical_Hit)
	//{
	//INT nCrit = m_pCaster->GetAttController().GetAttValue(EHA_CritInc);
	//nBaseHeal *= (Max(10000, 17500 + nCrit) /10000.0f);
	//}

	//DamageInfo.nDamage = Max(0, nBaseHeal) * fHealPct;
}
