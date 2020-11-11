package scene

import (
	"fmt"
	"sync/atomic"

	log "github.com/rs/zerolog/log"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
)

type AuraTrigger struct {
	Aura *Aura
	EffIndex int32
}

type CombatCtrl struct {
	mapSpells map[uint64]*Spell // 技能列表
	owner     SceneUnit         // 拥有者
	idGen     uint64            // id generator

	auraPool sync.Pool // aura 池
	arrayAura [define.Combat_MaxAura]*Aura // 当前aura列表
	listDelAura *list.List  // 待删除aura列表 List<*Aura>
	listSpellResultTrigger *list.List 		// 技能作用结果触发器 List<*AuraTrigger>
	listServentStateTrigger [define.StateChangeMode_End]*list.List 		// 技能作用结果触发器 List<*AuraTrigger>
	listDmgModTrigger [3]*list.List 		// DmgMod触发器 List<*AuraTrigger>
	listBehaviourTrigger [define.BehaviourType_End]*list.List 		// 行为触发器 List<*AuraTrigger>
	listAuraStateTrigger [define.StateChangeMode_End]*list.List 		// Aura状态改变触发器 List<*AuraTrigger>
}

func NewCombatCtrl(owner SceneUnit) *CombatCtrl {
	c := &CombatCtrl{
		mapSpells: make(map[uint64]*Spell, define.Combat_MaxSpell),
		arrayAura: make([]*Aura, define.Combat_MaxAura, define.Combat_MaxAura),
		listDelAura: list.New(),
		listSpellResultTrigger: list.New(),
	}

	c.auraPool.New = c.createAura(-1)

	for k, _ := range c.listServentStateTrigger {
		c.listServentStateTrigger[k] = list.New()
	}

	for k, _ := range c.listDmgModTrigger {
		c.listDmgModTrigger[k] = list.New()
	}

	for k, _ := range c.listBehaviourTrigger {
		c.listBehaviourTrigger[k] = list.New()
	}

	for k, _ := range c.listAuraStateTrigger {
		c.listAuraStateTrigger[k] = list.New()
	}

	c.owner = owner
	idGen = 0
	return c
}

//-------------------------------------------------------------------------------
// 创建与销毁Aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) createAura(index int32) *Aura {
	aura := &Aura{}
	aura.Reset(c.owner, index)
	return aura
}

//-------------------------------------------------------------------------------
// 施放技能
//-------------------------------------------------------------------------------
func (c *CombatCtrl) CastSpell(spellId uint32, caster, target SceneUnit, triggered bool) error {
	if len(c.mapSpells) >= define.Combat_MaxSpell {
		err := fmt.Errorf("spell list length >= <%d>", len(c.mapSpells))
		log.Warn().Err(err).Uint32("spell_id", spellId).Send()
		return err
	}

	entry := entries.GetSpellEntry(spellId)
	if entry == nil {
		err := fmt.Errorf("get spell entry failed")
		log.Warn().Err(err).Uint32("spell_id", spellId).Send()
		return err
	}

	s := NewSpell(spellId,
		WithSpellEntry(entry),
		WithSpellCaster(caster),
		WithSpellTarget(target),
		WithSpellTriggered(triggered),
	)

	s.Cast()

	c.mapSpells[atomic.AddUint64(&c.idGen, 1)] = s

	return nil
}

func (c *CombatCtrl) Update() {
	for spellGuid, s := range c.mapSpells {
		log.Info().Uint32("spell_id", s.GetOptions().Id).Msg("spell update...")
		delete(c.mapSpells, spellGuid)
	}

}

//-------------------------------------------------------------------------------
// 技能结果触发
//-------------------------------------------------------------------------------
func (c *CombatCtrl) TriggerBySpellResult(isCaster bool, target SceneUnit, dmgInfo *define.tagCalcDamageInfo) {
	if dmgInfo.ProcEx & define.AuraEventEx_Internal_Cant_Trigger != 0 {
		return
	}

	for trigger := c.listSpellResultTrigger.Front(); trigger != nil; trigger = trigger.Next() {
		auraTrigger := trigger.Value.(*AuraTrigger)

		if auraTrigger.Aura == nil || auraTrigger.Aura.IsRemoved() || auraTrigger.Aura.IsHangup() {
			continue
		}

		triggerEntry := entries.GetAuraTriggerEntry(auraTrigger.Aura.Opts().Entry.TriggerID[auraTrigger.EffIndex])
		if triggerEntry == nil {
			continue
		}

		// 检查技能
		spellEntry := entries.GetSpellEntry(dmgInfo.SpellId)
		if spellEntry == nil {
			continue
		}

		// todo check
		//if (VALID(pTriggerEntry->dwSpellID) && (SpellType(pTriggerEntry->dwSpellID) != SpellType(DmgInfo.dwSpellID)))
			//continue;

		//if (VALID(pTriggerEntry->dwFamilyMask) && !(pTriggerEntry->dwFamilyMask & pSpellEntry->dwFamilyMask))
			//continue;

		//if (VALID(pTriggerEntry->eFamilyRace) && (pTriggerEntry->eFamilyRace != pSpellEntry->eFamilyRace))
			//continue;

		//if( VALID(pTriggerEntry->dwDmgInfoType) && !(pTriggerEntry->dwDmgInfoType & (1 << DmgInfo.eType)) )
			//continue;

		//if( VALID(pTriggerEntry->eSchool) && pTriggerEntry->eSchool != DmgInfo.eSchool )
			//continue;

		// 检查响应事件
		if define.AuraEventEx_Trigger_Always != triggerEntry.TriggerMisc2 {
			if isCaster {
				if triggerEntry.TriggerMisc1 != 0 && (dmgInfo.ProcCaster & triggerEntry.TriggerMisc1 == 0) {
					continue
				}

				if triggerEntry.TriggerMisc2 != 0 && (dmgInfo.ProcEx & triggerEntry.TriggerMisc2 == 0) {
					continue
				}
			} else {
				if triggerEntry.TriggerMisc1 != 0 && (dmgInfo.dwProcTarget & triggerEntry.TriggerMisc1 == 0) {
					continue
				}

				if triggerEntry.TriggerMisc2 != 0 && (dmgInfo.dwProcEx & triggerEntry.TriggerMisc2 == 0) {
					continue
				}
			}

			// 验证触发条件
			if !c.checkTriggerCondition(triggerEntry, target) {
				continue
			}
		}

		// 计算触发几率
		if c.owner.GetScene().GetRandom().Rand(1, 10000) > triggerEntry.EventProp {
			continue
		}

		// 作用效果
		auraTrigger.Aura.CalAuraEffect(define.AuraEffectStep_Effect, auraTrigger.EffIndex, dmgInfo, target)
	}
}

//-------------------------------------------------------------------------------
// 状态改变触发
//-------------------------------------------------------------------------------
func (c *CombatCtrl) TriggerByServentState(state define.EHeroState, add bool) {
	var listTrigger *list.List
	if bAdd {
		listTrigger = c.listServentStateTrigger[define.StateChangeMode_Add]
	} else {
		listTrigger = c.listServentStateTrigger[define.StateChangeMode_Remove]
	}

	for trigger := listTrigger.Front(); trigger != nil; trigger = trigger.Next() {
		auraTrigger := trigger.Value.(*AuraTrigger)

		// 是否已经废弃
		if auraTrigger.Aura == nil || auraTrigger.Aura.IsRemoved() || auraTrigger.Aura.IsHangup() {
			continue
		}

		triggerEntry := entries.GetAuraTriggerEntry(auraTrigger.Aura.Opts().Entry.TriggerId[auraTrigger.EffIndex])
		if triggerEntry == nil {
			continue
		}

		// 检查响应状态
		if triggerEntry.TriggerMisc2 & (1 << state) == 0 {
			continue
		}

		// 计算触发几率
		if c.owner.GetScene().GetRandom().Rand(1, 10000) > triggerEntry.EventProp {
			continue
		}

		// 验证触发条件
		if !c.checkTriggerCondition(triggerEntry) {
			continue
		}

		// 作用效果
		auraTrigger.Aura.CalAuraEffect(define.AuraEffectStep_Effect, auraTrigger.EffIndex)
	}

	c.removeAuraByState(state)

}

//-------------------------------------------------------------------------------
// 行为触发
//-------------------------------------------------------------------------------
func (c *CombatCtrl) TriggerByBehaviour(behaviour define.EBehaviourType, 
target SceneUnit, 
ProcCaster , ProcEx uint32, 
spellType define.ESpellType) (triggerCount int32) {

	if behaviour < 0 || behaviour >= define.BehaviourType_End {
		return
	}

	listTrigger := c.listBehaviourTrigger[behaviour]
	for trigger := listTrigger.Front(); trigger != nil; trigger = trigger.Next() {
		auraTrigger := trigger.Value.(*AuraTrigger)

		// 是否已经废弃
		if auraTrigger.Aura == nil || auraTrigger.Aura.IsRemoved() || auraTrigger.IsHangup() {
			continue
		}

		triggerEntry := entries.GetAuraTriggerEntry(auraTrigger.Aura.Opts().Entry.TriggerId[auraTrigger.EffIndex])
		if triggerEntry == nil {
			continue
		}

		// 检查响应行为
		if triggerEntry.FamilyMask & (1 << behaviour) == 0 {
			continue
		}

		if triggerEntry.SpellTypeMask != 0 {
			if triggerEntry.SpellTypeMask & (1 << spellType) == 0 {
				continue
			}
		}

		// 计算触发几率
		if c.owner.GetScene().GetRandom().Rand(1, 10000) > triggerEntry.EventProp {
			continue
		}

		if procCaster != 0 && triggerEntry.TriggerMisc1 != 0 {
			if procCaster & triggerEntry.TriggerMisc1 == 0 {
				continue
			}
		}

		if procEx != 0 && triggerEntry.TriggerMisc2 != 0 {
			if procEx & triggerEntry.TriggerMisc2 == 0 {
				continue
			}
		}

		// 验证触发条件
		if !c.checkTriggerCondition(triggerEntry, target) {
			continue
		}

		// 作用效果
		result := auraTrigger.Aura.CalAuraEffect(define.AuraEffectStep_Effect, auraTrigger.EffIndex, nil, target)
		if define.AuraAddResult_Success == result {
			triggerCount++
		}
	}

}

//-------------------------------------------------------------------------------
// aura state改变触发
//-------------------------------------------------------------------------------
func (c *CombatCtrl) TriggerByAuraState(state int32, add bool) {
	var listTrigger *list.List
	if add {
		listTrigger = c.listAuraStateTrigger[define.StateChangeMode_Add]
	} else {
		listTrigger = c.listAuraStateTrigger[define.StateChangeMode_Remove]
	}

	for trigger := listTrigger.Front(); trigger != nil; trigger = trigger.Next() {
		auraTrigger := trigger.Value.(*AuraTrigger)
		
		// 是否已经废弃
		if auraTrigger.Aura == nil || auraTrigger.Aura.IsRemoved() || auraTrigger.Aura.IsHangup() {
			continue
		}

		triggerEntry := entries.GetAuraTriggerEntry(auraTrigger.Aura.Opts().Entry.TriggerId[auraTrigger.EffIndex])
		if triggerEntry == nil {
			continue
		}

		// 检查响应状态
		if triggerEntry.TriggerMisc2 & (1 << state) == 0 {
			continue
		}

		// 计算触发几率
		if c.owner.GetScene().GetRandom().Rand(1, 10000) > triggerEntry.EventProp {
			continue
		}

		// 验证触发条件
		if !c.checkTriggerCondition(triggerEntry) {
			continue
		}

		// 作用效果
		auraTrigger.Aura.CalAuraEffect(define.AuraEffectStep_Effect, auraTrigger.EffIndex)
	}
}

//-------------------------------------------------------------------------------
// 计算效果参数
//-------------------------------------------------------------------------------
func (c *CombatCtrl) CalSpellPoint(spellEntry *define.SpellEntry, points []int32, multiple []float32, level int32) {
	scene := c.owner.GetScene()
	if scene == nil {
		return
	}

	var basePoint int32 = 0
	var randPoint int32 = 0
	for i := 0; i < define.SpellEffectNum; i++ {
		basePoint = spellEntry.BasePoints[i]
		points[i] = basePoint + level * spellEntry.LevelPoints[i]
		multiple[i] = spellEntry.Multiple[i]
	}
}

//-------------------------------------------------------------------------------
// 根据目标等级计算效果
//-------------------------------------------------------------------------------
func (c *CombatCtrl) CalDecByTargetPoint(spellEntry *define.SpellEntry, points []int32, multiple []float32, level int32) {
	basePoint := 0.0
	for i := 0; i < define.SpellEffectNum; i++ {
		basePoint = spellEntry.BasePoints[i]

		if c.owner.GetLevel() > level {
			basePoint -= (c.owner.GetLevel() - level) * spellEntry.Multiple[i] / 10000.0 * spellEntry.BasePoints[i]
			points[i] = basePoint
			multiple[i] = 10000.0

			if basePoint < spellEntry.LevelPoints[i] {
				points[i] = spellEntry.LevelPoints[i]
			} else {
				points[i] = basePoint
			}
		}
	}
}



func (c *CombatCtrl) ClearAllAura() {
	for k, aura := range c.arrayAura {
		if aura != nil  {
			c.auraPool.Put(aura)
			c.arrayAura[k] = nil
		}
	}

	c.listDelAura.Init()
	c.listSpellResultTrigger.Init()

	for n := 0; n < define.StateChangeMode_End; n++ {
		c.listServentStateTrigger[n].Init()
		c.listAuraStateTrigger[n].Init()
	}

	for  n := 0; n < 3; n++ {
		c.listDmgModTrigger[n].Init()
	}

	for n := 0; n < define.BehaviourType_End; n++ {
		c.listBehaviourTrigger[n].Init()
	}
}

//-------------------------------------------------------------------------------
// 添加Aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) AddAura(auraId uint32, 
caster SceneUnit, 
amount int32, 
level int32, 
spellType define.ESpellType, 
ragePctMod float32, 
wrapTime int32) uint32 {

	auraEntry := entries.GetAuraEntry(auraId)
	if auraEntry == nil {
		return define.AuraAddResult_Null
	}

	// 检查免疫状态
	if c.owner.HasImmunityAny(define.ImmunityType_Mechanic, auraEntry.MechanicFlags) {
		return define.AuraAddResult_Immunity
	}

	// 生成Aura
	tempAura := c.auraPool.Get().(*Aura)
	if tempAura == nil {
		return define.AuraAddResult_Null
	}

	if !tempAura.Init(auraEntry, c.owner, caster, amount, level, spellType, ragePctMod, wrapTime) {
		c.auraPool.Put(tempAura)
		return define.AuraAddResult_Null
	}

	// 检查效果
	if define.AuraAddResult_Success != tempAura.CalAuraEffect(define.AuraEffectStep_Check) {
		c.auraPool.Put(tempAura)
		return define.AuraAddResult_Immunity
	}

	// 取得可用空位
	wrapResult define.EAuraWrapResult
	aura := c.generateAura(tempAura, wrapResult)
	if aura == nil {
		c.auraPool.Put(tempAura)
		return define.AuraAddResult_Full
	}

	if wrapResult == define.AuraWrapResult_Invalid {
		c.auraPool.Put(tempAura)
		return define.AuraAddResult_Inferior
	}
	
	switch wrapResult {
	case define.AuraWrapResult_Replace:
		fallthrough
	case define.AuraWrapResult_Add:
		fallthrough
	case define.AuraWrapResult_Wrap:
		aura.CalcApplyEffect(true, true)
	default:
		c.auraPool.Put(tempAura)
		return define.AuraAddResult_Inferior
	}

	return define.AuraAddResult_Success
}

//-------------------------------------------------------------------------------
// 移除Aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) RemoveAura(aura *Aura, mode define.EAuraRemoveMode) bool {
	if aura == nil || (mode & define.AuraRemoveMode_Removed == 0) || aura.IsRemoved() {
		return false
	}

	auraEntry := aura.Opts().Entry
	if auraEntry == nil {
		return false
	}

	slotIndex := aura.GetSlotIndex()

	// 是否计算过Apply效果
	calcApplyEffect := (aura.GetRemoveMode() & define.AuraRemoveMode_Running != 0) && !aura.IsHangup()

	// 设置标志位
	aura.AddRemoveMode(mode)

	// 转移到废弃列表
	c.arrayAura[slotIndex] = nil
	c.listDelAura.PushBack(aura)

	// 发送同步消息
	if define.AuraRemoveMode_Destroy & mode != 0 {
		return true
	}

	aura.SyncAuraToClient(define.AuraSyncStep_Remove)

	if !aura.IsApplyLocked() {
		if calcApplyEffect {
			aura.CalAuraEffect(define.AuraEffectStep_Remove, -1, &mode, nil)
		}

		aura.InvalidApplyLock()
	}

	return true
}

//-------------------------------------------------------------------------------
// 配置Aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) generateAura(aura *Aura, wrapResult *define.EAuraWrapResult) (newAura *Aura, newWrapResult define.EAuraWrapResult) {
	newAura = nil
	newWrapResult = define.AuraWrapResult_Add

	if aura == nil {
		return
	}

	auraEntry := aura.Opts().Entry
	if auraEntry == nil {
		return
	}

	validSlot := -1
	listReplace := list.New()
	for index := 0; index < define.Combat_MaxAura; index++ {
		if c.arrayAura[index] == nil {
			if validSlot == -1 {
				validSlot = index
			}
			continue
		}

		result := c.arrayAura[index].CheckWrapResult(aura)
		if newWrapResult < result {
			newWrapResult = result
		}

		switch result {
			case define.AuraWrapResult_Add:
				continue

			case define.AuraWrapResult_Invalid:
				return

			case define.AuraWrapResult_Wrap:
				fallthrough
			case define.AuraWrapResult_Replace:
				listReplace.PushBack(index)
		}
	}

	switch wrapResult {
	case define.AuraWrapResult_Add:
		if validSlot == -1 {
			return
		}

		aura.SetSlotIndex(validSlot)
		c.arrayAura[validSlot] = aura

	case define.AuraWrapResult_Wrap:
		fallthrough
	case define.AuraWrapResult_Replace:
		e := listReplace.Front()
		if e == nil {
			return
		}

		// 取得slot
		replaceSlot := e.Value.(int)

		// 检查是否是合理位置
		if replaceSlot < 0 || replaceSlot >= define.Combat_MaxAura {
			validSlot = replaceSlot
		}

		// 没位置了
		if validSlot == -1 {
			return
		}

		aura.SetSlotIndex(validSlot)

		// 删除
		for e := listReplace.Front(); e != nil; e = e.Next() {
			index := e.Value.(int)
			c.removeAura(c.arrayAura[index], define.AuraRemoveMode_Replace)
		}

		// 放入
		c.arrayAura[validSlot] = aura
	default:
		return
	}
	
	newAura = c.arrayAura[validSlot]
	return
}

//-------------------------------------------------------------------------------
// 注册Aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) RegisterAura(aura *Aura) {
	if aura == nil || (aura.GetRemoveMode() & define.AuraRemoveMode_Registered != 0) {
		return
	}

	c.registerAuraTrigger(aura)

	// Add Aura State
	c.addAuraState(aura.Opts().Entry.AuraState)

	// Mod Aura Mode
	aura.AddRemoveMode(define.AuraRemoveMode_Registered)
}

func (c *CombatCtrl) UnregisterAura(aura *Aura) {
	if aura == nil || (aura.GetRemoveMode() & define.AuraRemoveMode_Registered == 0) {
		return
	}

	c.unRegisterAuraTrigger(aura)

	// Dec Aura State
	c.decAuraState(aura.Opts().Entry.AuraState)

	aura.DecRemoveMode(define.AuraRemoveMode_Registered)
}

//-------------------------------------------------------------------------------
// 注册触发器
//-------------------------------------------------------------------------------
VOID CombatController::RegisterAuraTrigger( Aura* pAura )
{
	if (!VALID(pAura))
		return;

	const tagAuraEntry* pAuraEntry = pAura->GetAuraEntry();
	if (!VALID(pAuraEntry))
		return;

	for (INT index=0; index<EFFECT_NUM; index++)
	{
		if (!VALID(pAuraEntry->dwTriggerID[index]))
			continue;

		const tagAuraTriggerEntry* pTriggerEntry = sAuraTriggerEntry(pAuraEntry->dwTriggerID[index]);
		if (!VALID(pTriggerEntry))
			continue;

		if (IsDmgModEffect(pAuraEntry->eEffect[index]))
		{
			RegisterDmgMod(pAura, index);
			continue;
		}

		tagAuraTrigger sTrigger;
		sTrigger.pAura		=	pAura;
		sTrigger.nEffIndex	=	index;

		switch (pTriggerEntry->eTriggerType)
		{
		case EATT_SpellResult:
			{
				if( pTriggerEntry->bAddHead )
				{
					m_listSpellResultTrigger.push_front(sTrigger);
				}
				else
				{
					m_listSpellResultTrigger.push_back(sTrigger);
				}
			}
			break;

		case EATT_State:
			if (MIsBetween(pTriggerEntry->dwTriggerMisc1, 0, ESCM_End))			m_listServentStateTrigger[pTriggerEntry->dwTriggerMisc1].push_back(sTrigger);
			break;

		case EATT_Behaviour:
			for (INT n = 0; n < EBT_End; n++)
			{
				if ((1 << n) & pTriggerEntry->dwFamilyMask)					m_listBehaviourTrigger[n].push_back(sTrigger);
			}
			break;

		case EATT_AuraState:
			if (MIsBetween(pTriggerEntry->dwTriggerMisc1, 0, ESCM_End))			m_listAuraStateTrigger[pTriggerEntry->dwTriggerMisc1].push_back(sTrigger);
			break;

		default:
			continue;
		}
	}
}

VOID CombatController::UnRegisterAuraTrigger( Aura* pAura )
{
	if (!VALID(pAura))
		return;

	const tagAuraEntry* pAuraEntry = pAura->GetAuraEntry();
	if (!VALID(pAuraEntry))
		return;

	for (INT index=0; index<EFFECT_NUM; index++)
	{
		if (!VALID(pAuraEntry->dwTriggerID[index]))
			continue;

		const tagAuraTriggerEntry* pTriggerEntry = sAuraTriggerEntry(pAuraEntry->dwTriggerID[index]);
		if (!VALID(pTriggerEntry))
			continue;

		if (IsDmgModEffect(pAuraEntry->eEffect[index]))
		{
			UnregisterDmgMod(pAura, index);
			continue;
		}

		switch (pTriggerEntry->eTriggerType)
		{
		case EATT_SpellResult:
			{
				AuraTriggerList::iterator iter = m_listSpellResultTrigger.begin();
				for (; iter != m_listSpellResultTrigger.end(); ++iter)
				{
					if ((iter->pAura == pAura) && (iter->nEffIndex == index))
					{
						m_listSpellResultTrigger.erase(iter);
						break;
					}
				}
			}
			break;

		case EATT_State:
			{
				AuraTriggerList& listTrigger = m_listServentStateTrigger[pTriggerEntry->dwTriggerMisc1];
				AuraTriggerList::iterator iter = listTrigger.begin();
				for (; iter != listTrigger.end(); ++iter)
				{
					if ((iter->pAura == pAura) && (iter->nEffIndex == index))
					{
						listTrigger.erase(iter);
						break;
					}
				}
			}
			break;

		case EATT_Behaviour:
			for (INT n = 0; n < EBT_End; n++)
			{
				if ((1 << n) & pTriggerEntry->dwTriggerMisc1)
				{
					AuraTriggerList& listTrigger = m_listBehaviourTrigger[n];
					AuraTriggerList::iterator iter = listTrigger.begin();
					for (; iter != listTrigger.end(); ++iter)
					{
						if ((iter->pAura == pAura) && (iter->nEffIndex == index))
						{
							listTrigger.erase(iter);
							break;
						}
					}
				}
			}
			break;

		case EATT_AuraState:
			{
				AuraTriggerList& listTrigger = m_listAuraStateTrigger[pTriggerEntry->dwTriggerMisc1];
				AuraTriggerList::iterator iter = listTrigger.begin();
				for (; iter != listTrigger.end(); ++iter)
				{
					if ((iter->pAura == pAura) && (iter->nEffIndex == index))
					{
						listTrigger.erase(iter);
						break;
					}
				}
			}
			break;

		default:
			continue;
		}
	}
}

//-------------------------------------------------------------------------------
// 计算aura效果
//-------------------------------------------------------------------------------
VOID CombatController::CalAuraEffect(INT32 nCurRound)
{
	if( m_AuraNum == 0 )
		return;

	for (INT n = 0; n < m_Aura.Size(); n++)
	{
		if (VALID(m_Aura[n]) && !m_pOwner->HasState(EHS_Freeze))
		{
			if( !(m_Aura[n]->GetAuraEntry()->dwRoundUpdateMask & (1 << nCurRound)) )
				continue;

			m_Aura[n]->CalAuraEffect(EAES_Effect);
		}

	}
}

//-------------------------------------------------------------------------------
// 清除aura
//-------------------------------------------------------------------------------
VOID CombatController::DeleteAura( Aura* pAura )
{
	if (!VALID(pAura))
		return;

	EAuraRemoveMode eRemoveMode = pAura->GetRemoveMode();
	if (pAura->IsApplyLockValid())
	{
		if (!(EARM_Destroy & eRemoveMode))
			pAura->CalAuraEffect(EAES_Remove, INVALID, &eRemoveMode, NULL);
	}

	// 反注册
	UnregisterAura(pAura);

	// 释放内存
	ReleaseAura(pAura);
}

//-------------------------------------------------------------------------------
// 伤害触发器
//-------------------------------------------------------------------------------
VOID CombatController::RegisterDmgMod( Aura* pAura, INT nIndex )
{
	if (!VALID(pAura) || !MIsBetween(nIndex, 0, EFFECT_NUM))	return;

	const tagAuraEntry* pAuraEntry = pAura->GetAuraEntry();
	ASSERT(VALID(pAuraEntry));

	for (INT n=0; n<X_DmgModTypeNum; n++)
	{
		if (pAuraEntry->eEffect[nIndex] == AuraDmgModType[n])
		{
			tagAuraTrigger stTrigger;
			stTrigger.pAura		= pAura;
			stTrigger.nEffIndex	= nIndex;

			m_listDmgModTrigger[n].push_back(stTrigger);

			return;
		}
	}
}

VOID CombatController::UnregisterDmgMod( Aura* pAura, INT nIndex )
{
	if (!VALID(pAura) || !MIsBetween(nIndex, 0, EFFECT_NUM))	return;

	const tagAuraEntry* pAuraEntry = pAura->GetAuraEntry();
	ASSERT(VALID(pAuraEntry));

	for (INT n=0; n<X_DmgModTypeNum; n++)
	{
		if (pAuraEntry->eEffect[nIndex] == AuraDmgModType[n])
		{
			AuraTriggerList::iterator iter = m_listDmgModTrigger[n].begin();
			for (; iter != m_listDmgModTrigger[n].end(); ++iter)
			{
				if ((iter->pAura == pAura) && (iter->nEffIndex == nIndex))
				{
					m_listDmgModTrigger[n].erase(iter);
					return;
				}
			}
		}
	}
}

BOOL CombatController::IsDmgModEffect( EAuraEffectType eType )
{
	switch (eType)
	{
	case EAFT_DmgMod:
		//case ESAET_NewDmg:
	case EAFT_DmgFix:
	case EAFT_AbsorbAllDmg:
		return TRUE;
	}

	return FALSE;
}

VOID CombatController::TriggerByDmgMod( BOOL bCaster, EntityHero* pTarget, tagCalcDamageInfo& DmgInfo )
{
	for (INT n=0; n<X_DmgModTypeNum; n++)
	{
		AuraTriggerList::iterator iter = m_listDmgModTrigger[n].begin();
		for (; iter != m_listDmgModTrigger[n].end(); ++iter)
		{
			tagAuraTrigger sTrigger = *iter;

			// 是否已经废弃
			if (!VALID(sTrigger.pAura) || sTrigger.pAura->IsRemoved() || sTrigger.pAura->IsHangup())
				continue;

			const tagAuraTriggerEntry* pTriggerEntry = sAuraTriggerEntry(sTrigger.pAura->GetAuraEntry()->dwTriggerID[sTrigger.nEffIndex]);
			if (!VALID(pTriggerEntry))
				continue;

			// 检查技能
			const tagSpellBase* pSpellEntry = IS_SPELL(DmgInfo.dwSpellID) ? sSpellEntry(DmgInfo.dwSpellID) : (const tagSpellBase*)sAuraEntry(DmgInfo.dwSpellID);
			if (!VALID(pSpellEntry))
				continue;

			if (VALID(pTriggerEntry->dwSpellID) && (pTriggerEntry->dwSpellID != DmgInfo.dwSpellID))
				continue;

			if (VALID(pTriggerEntry->dwFamilyMask) && !(pTriggerEntry->dwFamilyMask & pSpellEntry->dwFamilyMask))
				continue;

			if (VALID(pTriggerEntry->eFamilyRace) && (pTriggerEntry->eFamilyRace != pSpellEntry->eFamilyRace))
				continue;

			if( VALID(pTriggerEntry->dwDmgInfoType) && !(pTriggerEntry->dwDmgInfoType & (1 << DmgInfo.eType)) )
				continue;

			if( VALID(pTriggerEntry->eSchool) && !(pTriggerEntry->eSchool != DmgInfo.eSchool) )
				continue;

			// 检查响应事件
			if (EAEE_Trigger_Always != pTriggerEntry->dwTriggerMisc2)
			{
				if (bCaster)
				{
					if (VALID(pTriggerEntry->dwTriggerMisc1) && !(DmgInfo.dwProcCaster & pTriggerEntry->dwTriggerMisc1))
						continue;

					if (VALID(pTriggerEntry->dwTriggerMisc2) && !(DmgInfo.dwProcEx & pTriggerEntry->dwTriggerMisc2))
						continue;
				}
				else
				{
					if (VALID(pTriggerEntry->dwTriggerMisc1) && !(DmgInfo.dwProcTarget & pTriggerEntry->dwTriggerMisc1))
						continue;

					if (VALID(pTriggerEntry->dwTriggerMisc2) && !(DmgInfo.dwProcEx & pTriggerEntry->dwTriggerMisc2))
						continue;
				}

				// 验证触发条件
				if (!CheckTriggerCondition(pTriggerEntry, pTarget))
					continue;
			}

			// 计算触发几率
			if( m_pOwner->GetScene()->GetRandom().Rand(1, 10000) > pTriggerEntry->nEventProp )
				continue;

			// 作用效果
			sTrigger.pAura->CalAuraEffect(EAES_Effect, sTrigger.nEffIndex, &DmgInfo, pTarget);
		}
	}
}

//-------------------------------------------------------------------------------
// 触发器条件检查
//-------------------------------------------------------------------------------
BOOL CombatController::CheckTriggerCondition( const tagAuraTriggerEntry* pTrigger, EntityHero* pTarget/* = NULL*/)
{
	if (!VALID(pTrigger))
		return FALSE;

	switch (pTrigger->eConditionType)
	{
	case EAEC_HPLowerFlat:
		if (m_pOwner->GetAttController().GetAttValue(EHA_CurHP) < pTrigger->nConditionMisc1)
			return TRUE;
		break;

	case EAEC_HPLowerPct:
		if ((FLOAT)m_pOwner->GetAttController().GetAttValue(EHA_CurHP) / (FLOAT)m_pOwner->GetAttController().GetAttValue(EHA_MaxHP) * 10000.0f
			< pTrigger->nConditionMisc1)
			return TRUE;
		break;

	case EAEC_HPHigherFlat:
		if (m_pOwner->GetAttController().GetAttValue(EHA_CurHP) >= pTrigger->nConditionMisc1)
			return TRUE;
		break;

	case EAEC_HPHigherPct:
		if ((FLOAT)m_pOwner->GetAttController().GetAttValue(EHA_CurHP) / (FLOAT)m_pOwner->GetAttController().GetAttValue(EHA_MaxHP) * 10000.0f
			>= pTrigger->nConditionMisc1)
			return TRUE;
		break;

	case EAEC_AnyUnitState:
		{
			DWORD dwStateMask = pTrigger->nConditionMisc1;
			if (m_pOwner->HasStateAny(dwStateMask))
				return TRUE;
		}
		break;

	case EAEC_AllUnitState:
		{
			DWORD dwStateMask = pTrigger->nConditionMisc1;
			if (m_pOwner->HasStateAll(dwStateMask))
				return TRUE;
		}
		break;

	case EAEC_AuraState:
		if (HasAuraState(pTrigger->nConditionMisc1))
			return TRUE;
		break;

	case EAEC_TargetClass:
		if (VALID(pTarget) && (1 << pTarget->GetEntry()->eClass) & pTrigger->nConditionMisc1)
			return TRUE;
		break;

	case EAEC_StrongTarget:
		if (VALID(pTarget) && pTarget->GetAttController().GetAttValue(EHA_CurHP) > m_pOwner->GetAttController().GetAttValue(EHA_CurHP))
			return TRUE;
		break;

	case EAEC_TargetAuraState:
		{
			if (VALID(pTarget) && pTarget->GetCombatController().HasAuraState(pTrigger->nConditionMisc1))
				return TRUE;
		}
		break;

	case EAEC_TargetAllUnitState:
		{
			DWORD dwStateMask = pTrigger->nConditionMisc1;
			if (VALID(pTarget) && pTarget->HasStateAll(dwStateMask))
				return TRUE;
		}
		break;

	case EAEC_TargetAnyUnitState:
		{
			DWORD dwStateMask = pTrigger->nConditionMisc1;
			if (VALID(pTarget) && pTarget->HasStateAny(dwStateMask))
				return TRUE;
		}
		break;

	case EAEC_None:
	default:
		return TRUE;
	}

	return FALSE;
}

//-------------------------------------------------------------------------------
// 按状态删除aura
//-------------------------------------------------------------------------------
VOID CombatController::RemoveAuraByState( EHeroState eState )
{
	for (INT index=0; index<m_Aura.Size(); index++)
	{
		if (!VALID(m_Aura[index]) || !(m_Aura[index]->GetRemoveMode() & EARM_Running))
			continue;

		const tagAuraEntry* pAuraEntry = m_Aura[index]->GetAuraEntry();
		if (!VALID(pAuraEntry))
			continue;

		// Add State:		1-Conflict	0-Ignore
		// Remove State:	1-Ignore	0-Dependence
		if (pAuraEntry->OwnerStateCheckFlag.IsSet(eState) && (pAuraEntry->OwnerStateLimit.IsSet(eState) != m_pOwner->HasState(eState)) )
		{
			RemoveAura(m_Aura[index], EARM_Dispel);
		}
	}
}

//-------------------------------------------------------------------------------
// 按施放者和ID删除aura
//-------------------------------------------------------------------------------
VOID CombatController::RemoveAuraByCaster( EntityHero* pCaster )
{
	if (!VALID(pCaster))
		return;

	for (INT index=0; index<m_Aura.Size(); index++)
	{
		if (!VALID(m_Aura[index]) || !(m_Aura[index]->GetRemoveMode() & EARM_Running))
			continue;

		const tagAuraEntry* pAuraEntry = m_Aura[index]->GetAuraEntry();
		if (!VALID(pAuraEntry))
			continue;

		if (pAuraEntry->bDependCaster && m_Aura[index]->GetCaster() == pCaster)
		{
			RemoveAura(m_Aura[index], EARM_Interrupt);
		}
	}
}

//-------------------------------------------------------------------------------
// 是否有指定TypeID的aura
//-------------------------------------------------------------------------------
Aura* CombatController::GetAuraByIDCaster( DWORD dwAuraID, EntityHero* pCaster)
{
	for (INT index=0; index<m_Aura.Size(); index++)
	{
		if (!VALID(m_Aura[index]) )
			continue;

		if( VALID(pCaster) )
		{
			if( m_Aura[index]->GetCaster() != pCaster )
				continue;
		}

		if (m_Aura[index]->GetAuraID() == dwAuraID )
		{
			return m_Aura[index];
		}
	}

	return NULL;
}

//-------------------------------------------------------------------------------
// 驱散
//-------------------------------------------------------------------------------
BOOL CombatController::DispelAura(DWORD dwDispelType, DWORD dwNum)
{
	BOOL bDisperse = FALSE;			// 是否驱散成功

	for (INT index=0; index<m_Aura.Size(); index++)
	{
		if (!VALID(m_Aura[index]))
			continue;

		if ( 0 != (m_Aura[index]->GetAuraEntry()->dwDispelFlags & dwDispelType) )
		{
			m_Aura[index]->Disperse();
			bDisperse = TRUE;

			if( --dwNum <= 0 )
				break;
		}
	}

	return bDisperse;
}

//-------------------------------------------------------------------------------
// 强化作用时间
//-------------------------------------------------------------------------------
VOID CombatController::ModAuraDuration(DWORD dwModType, DWORD dwModDuration)
{
	for (INT index=0; index<m_Aura.Size(); index++)
	{
		if (!VALID(m_Aura[index]))
			continue;

		if ( 0 != (m_Aura[index]->GetAuraEntry()->dwDurationFlags & dwModType) )
		{
			m_Aura[index]->ModDuration(dwModDuration);
		}
	}
}

//-------------------------------------------------------------------------------
// 增益和减益buff数量
//-------------------------------------------------------------------------------
VOID CombatController::GetPositiveAndNegativeNum(INT32 &nPosNum, INT32 &nNegNum)
{
	for (INT index=0; index<m_Aura.Size(); index++)
	{
		if (!VALID(m_Aura[index]))
			continue;

		// 被动buff不计数
		if (m_Aura[index]->IsPassive())
			continue;

		// 不可见buff不计数
		if (!m_Aura[index]->GetAuraEntry()->bHaveVisual)
			continue;

		if ( m_Aura[index]->GetAuraEntry()->nEffectPriority > 0 )
		{
			nPosNum++;
		}

		if ( m_Aura[index]->GetAuraEntry()->nEffectPriority < 0 )
		{
			nNegNum++;
		}
	}
}

//-------------------------------------------------------------------------------
// aura state
//-------------------------------------------------------------------------------
VOID CombatController::AddAuraState( INT32 nAuraState )
{
	if (!VALID(nAuraState))
		return;

	ASSERT(MIsBetween(nAuraState, 0, X_AuraFlagNum));

	BOOL bNewState = !m_AuraState.IsSet(nAuraState);

	m_AuraState.Set(nAuraState);

	if (bNewState)
	{
		TriggerByAuraState(nAuraState, TRUE);
	}
}

VOID CombatController::DecAuraState( INT32 nAuraState )
{
	if (!VALID(nAuraState))
		return;

	ASSERT(MIsBetween(nAuraState, 0, X_AuraFlagNum));

	if (!m_AuraState.IsSet(nAuraState))
		return;

	m_AuraState.Unset(nAuraState);

	if (!m_AuraState.IsSet(nAuraState))
	{
		TriggerByAuraState(nAuraState, FALSE);
	}
}

BOOL CombatController::HasAuraState( INT32 nAuraState )
{
	if (!VALID(nAuraState))
		return TRUE;

	return m_AuraState.IsSet(nAuraState);
}

BOOL CombatController::HasAuraStateAny( DWORD dwAuraStateMask )
{
	if (!VALID(dwAuraStateMask))
		return TRUE;

	return m_AuraState.IsSetAny(dwAuraStateMask);
}



