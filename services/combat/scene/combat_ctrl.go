package scene

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"

	log "github.com/rs/zerolog/log"
	"github.com/willf/bitset"
	"github.com/east-eden/server/define"
	"github.com/east-eden/server/entries"
)

type AuraTrigger struct {
	Aura     *Aura
	EffIndex int32
}

type CombatCtrl struct {
	mapSpells map[uint64]*Spell // 技能列表
	owner     SceneUnit         // 拥有者
	idGen     uint64            // id generator

	auraPool                sync.Pool                               // aura 池
	arrayAura               [define.Combat_MaxAura]*Aura            // 当前aura列表
	listDelAura             *list.List                              // 待删除aura列表 List<*Aura>
	listSpellResultTrigger  *list.List                              // 技能作用结果触发器 List<*AuraTrigger>
	listServentStateTrigger [define.StateChangeMode_End]*list.List  // 技能作用结果触发器 List<*AuraTrigger>
	listDmgModTrigger       [define.Combat_DmgModTypeNum]*list.List // DmgMod触发器 List<*AuraTrigger>
	listBehaviourTrigger    [define.BehaviourType_End]*list.List    // 行为触发器 List<*AuraTrigger>
	listAuraStateTrigger    [define.StateChangeMode_End]*list.List  // Aura状态改变触发器 List<*AuraTrigger>

	auraStateBitSet *bitset.BitSet
}

func NewCombatCtrl(owner SceneUnit) *CombatCtrl {
	c := &CombatCtrl{
		mapSpells:              make(map[uint64]*Spell, define.Combat_MaxSpell),
		listDelAura:            list.New(),
		listSpellResultTrigger: list.New(),
		auraStateBitSet:        bitset.New(define.AuraFlagNum),
	}

	c.auraPool.New = func() interface{} { return c.createAura() }

	for k := range c.listServentStateTrigger {
		c.listServentStateTrigger[k] = list.New()
	}

	for k := range c.listDmgModTrigger {
		c.listDmgModTrigger[k] = list.New()
	}

	for k := range c.listBehaviourTrigger {
		c.listBehaviourTrigger[k] = list.New()
	}

	for k := range c.listAuraStateTrigger {
		c.listAuraStateTrigger[k] = list.New()
	}

	c.owner = owner
	c.idGen = 0
	return c
}

//-------------------------------------------------------------------------------
// 创建与销毁Aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) createAura() *Aura {
	return &Aura{
		opts: DefaultAuraOptions(),
	}
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

	s := NewSpell(
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
	// 更新删除buff
	for n := 0; n < define.Combat_MaxAura; n++ {
		if c.arrayAura[n] != nil {
			c.arrayAura[n].RoundEnd()
		}
	}

	if c.listDelAura.Len() > 0 {
		for e := c.listDelAura.Front(); e != nil; e = e.Next() {
			c.deleteAura(e.Value.(*Aura))
		}

		c.listDelAura.Init()
	}
}

//-------------------------------------------------------------------------------
// 技能结果触发
//-------------------------------------------------------------------------------
func (c *CombatCtrl) TriggerBySpellResult(isCaster bool, target SceneUnit, dmgInfo *CalcDamageInfo) {
	if dmgInfo.ProcEx&uint32(define.AuraEventEx_Internal_Cant_Trigger) != 0 {
		return
	}

	for trigger := c.listSpellResultTrigger.Front(); trigger != nil; trigger = trigger.Next() {
		auraTrigger := trigger.Value.(*AuraTrigger)

		if auraTrigger.Aura == nil || auraTrigger.Aura.IsRemoved() || auraTrigger.Aura.IsHangup() {
			continue
		}

		triggerEntry := entries.GetAuraTriggerEntry(auraTrigger.Aura.Opts().Entry.TriggerId[auraTrigger.EffIndex])
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
		if uint32(define.AuraEventEx_Trigger_Always) != triggerEntry.TriggerMisc2 {
			if isCaster {
				if triggerEntry.TriggerMisc1 != 0 && (dmgInfo.ProcCaster&triggerEntry.TriggerMisc1 == 0) {
					continue
				}

				if triggerEntry.TriggerMisc2 != 0 && (dmgInfo.ProcEx&triggerEntry.TriggerMisc2 == 0) {
					continue
				}
			} else {
				if triggerEntry.TriggerMisc1 != 0 && (dmgInfo.ProcTarget&triggerEntry.TriggerMisc1 == 0) {
					continue
				}

				if triggerEntry.TriggerMisc2 != 0 && (dmgInfo.ProcEx&triggerEntry.TriggerMisc2 == 0) {
					continue
				}
			}

			// 验证触发条件
			if !c.checkTriggerCondition(triggerEntry, target) {
				continue
			}
		}

		// 计算触发几率
		if c.owner.GetScene().Rand(1, 10000) > int(triggerEntry.EventProp) {
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
	if add {
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
		if triggerEntry.TriggerMisc2&(1<<state) == 0 {
			continue
		}

		// 计算触发几率
		if c.owner.GetScene().Rand(1, 10000) > int(triggerEntry.EventProp) {
			continue
		}

		// 验证触发条件
		if !c.checkTriggerCondition(triggerEntry, nil) {
			continue
		}

		// 作用效果
		auraTrigger.Aura.CalAuraEffect(define.AuraEffectStep_Effect, auraTrigger.EffIndex, nil, nil)
	}

	c.removeAuraByState(state)

}

//-------------------------------------------------------------------------------
// 行为触发
//-------------------------------------------------------------------------------
func (c *CombatCtrl) TriggerByBehaviour(behaviour define.EBehaviourType,
	target SceneUnit,
	procCaster, procEx uint32,
	spellType define.ESpellType) (triggerCount int32) {

	if behaviour < 0 || behaviour >= define.BehaviourType_End {
		return
	}

	listTrigger := c.listBehaviourTrigger[behaviour]
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

		// 检查响应行为
		if triggerEntry.FamilyMask&(1<<behaviour) == 0 {
			continue
		}

		if triggerEntry.SpellTypeMask != 0 {
			if triggerEntry.SpellTypeMask&(1<<spellType) == 0 {
				continue
			}
		}

		// 计算触发几率
		if c.owner.GetScene().Rand(1, 10000) > int(triggerEntry.EventProp) {
			continue
		}

		if procCaster != 0 && triggerEntry.TriggerMisc1 != 0 {
			if procCaster&triggerEntry.TriggerMisc1 == 0 {
				continue
			}
		}

		if procEx != 0 && triggerEntry.TriggerMisc2 != 0 {
			if procEx&triggerEntry.TriggerMisc2 == 0 {
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

	return
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
		if triggerEntry.TriggerMisc2&(1<<state) == 0 {
			continue
		}

		// 计算触发几率
		if c.owner.GetScene().Rand(1, 10000) > int(triggerEntry.EventProp) {
			continue
		}

		// 验证触发条件
		if !c.checkTriggerCondition(triggerEntry, nil) {
			continue
		}

		// 作用效果
		auraTrigger.Aura.CalAuraEffect(define.AuraEffectStep_Effect, auraTrigger.EffIndex, nil, nil)
	}
}

//-------------------------------------------------------------------------------
// 计算效果参数
//-------------------------------------------------------------------------------
func (c *CombatCtrl) CalSpellPoint(spellBase *define.SpellBase, points []int32, multiple []float32, level uint32) {
	scene := c.owner.GetScene()
	if scene == nil {
		return
	}

	var basePoint int32 = 0
	for i := 0; i < define.SpellEffectNum; i++ {
		basePoint = spellBase.BasePoints[i]
		points[i] = basePoint + int32(level)*spellBase.LevelPoints[i]
		multiple[i] = float32(spellBase.Multiple[i])
	}
}

//-------------------------------------------------------------------------------
// 根据目标等级计算效果
//-------------------------------------------------------------------------------
func (c *CombatCtrl) CalDecByTargetPoint(spellBase *define.SpellBase, points []int32, multiple []float32, level uint32) {
	basePoint := 0.0
	for i := 0; i < define.SpellEffectNum; i++ {
		basePoint = float64(spellBase.BasePoints[i])

		if c.owner.GetLevel() > level {
			basePoint -= float64(int32(c.owner.GetLevel()-level) * spellBase.Multiple[i] / 10000.0 * spellBase.BasePoints[i])
			points[i] = int32(basePoint)
			multiple[i] = 10000.0

			if basePoint < float64(spellBase.LevelPoints[i]) {
				points[i] = spellBase.LevelPoints[i]
			} else {
				points[i] = int32(basePoint)
			}
		}
	}
}

func (c *CombatCtrl) ClearAllAura() {
	for k, aura := range c.arrayAura {
		if aura != nil {
			c.auraPool.Put(aura)
			c.arrayAura[k] = nil
		}
	}

	c.listDelAura.Init()
	c.listSpellResultTrigger.Init()

	for n := 0; define.EStateChangeMode(n) < define.StateChangeMode_End; n++ {
		c.listServentStateTrigger[n].Init()
		c.listAuraStateTrigger[n].Init()
	}

	for n := 0; n < define.Combat_DmgModTypeNum; n++ {
		c.listDmgModTrigger[n].Init()
	}

	for n := 0; define.EBehaviourType(n) < define.BehaviourType_End; n++ {
		c.listBehaviourTrigger[n].Init()
	}
}

//-------------------------------------------------------------------------------
// 添加Aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) AddAura(auraId uint32,
	caster SceneUnit,
	amount int32,
	level uint32,
	spellType define.ESpellType,
	ragePctMod float32,
	wrapTime int32) define.EAuraAddResult {

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

	tempAura.Init(
		WithAuraEntry(auraEntry),
		WithAuraOwner(c.owner),
		WithAuraCaster(caster),
		WithAuraAmount(amount),
		WithAuraLevel(level),
		WithAuraSpellType(spellType),
		WithAuraRagePctMod(ragePctMod),
		WithAuraCurWrapTimes(uint8(wrapTime)),
	)

	// 检查效果
	if define.AuraAddResult_Success != tempAura.CalAuraEffect(define.AuraEffectStep_Check, -1, nil, nil) {
		c.auraPool.Put(tempAura)
		return define.AuraAddResult_Immunity
	}

	// 取得可用空位
	aura, wrapResult := c.generateAura(tempAura)
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
	if aura == nil || (mode&define.AuraRemoveMode_Removed == 0) || aura.IsRemoved() {
		return false
	}

	auraEntry := aura.Opts().Entry
	if auraEntry == nil {
		return false
	}

	slotIndex := aura.opts.SlotIndex

	// 是否计算过Apply效果
	calcApplyEffect := (aura.GetRemoveMode()&define.AuraRemoveMode_Running != 0) && !aura.IsHangup()

	// 设置标志位
	aura.AddRemoveMode(mode)

	// 转移到废弃列表
	c.arrayAura[slotIndex] = nil
	c.listDelAura.PushBack(aura)

	// 发送同步消息
	if define.AuraRemoveMode_Destroy&mode != 0 {
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
func (c *CombatCtrl) generateAura(aura *Aura) (newAura *Aura, wrapResult define.EAuraWrapResult) {
	newAura = nil
	wrapResult = define.AuraWrapResult_Add

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
		if wrapResult < result {
			wrapResult = result
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

		aura.opts.SlotIndex = int8(validSlot)
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

		aura.opts.SlotIndex = int8(validSlot)

		// 删除
		for e := listReplace.Front(); e != nil; e = e.Next() {
			index := e.Value.(int)
			c.RemoveAura(c.arrayAura[index], define.AuraRemoveMode_Replace)
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
	if aura == nil || (aura.GetRemoveMode()&define.AuraRemoveMode_Registered != 0) {
		return
	}

	c.registerAuraTrigger(aura)

	// Add Aura State
	c.AddAuraState(aura.Opts().Entry.AuraState)

	// Mod Aura Mode
	aura.AddRemoveMode(define.AuraRemoveMode_Registered)
}

func (c *CombatCtrl) UnregisterAura(aura *Aura) {
	if aura == nil || (aura.GetRemoveMode()&define.AuraRemoveMode_Registered == 0) {
		return
	}

	c.unRegisterAuraTrigger(aura)

	// Dec Aura State
	c.DecAuraState(aura.Opts().Entry.AuraState)

	aura.DecRemoveMode(define.AuraRemoveMode_Registered)
}

//-------------------------------------------------------------------------------
// 注册触发器
//-------------------------------------------------------------------------------
func (c *CombatCtrl) registerAuraTrigger(aura *Aura) {
	if aura == nil {
		return
	}

	auraEntry := aura.Opts().Entry
	if auraEntry == nil {
		return
	}

	var index int32
	for index = 0; index < define.SpellEffectNum; index++ {
		if auraEntry.TriggerId[index] == 0 {
			continue
		}

		auraTriggerEntry := entries.GetAuraTriggerEntry(auraEntry.TriggerId[index])
		if auraTriggerEntry == nil {
			continue
		}

		if c.isDmgModEffect(auraEntry.Effects[index]) {
			c.registerDmgMod(aura, index)
			continue
		}

		auraTrigger := &AuraTrigger{
			Aura:     aura,
			EffIndex: index,
		}

		switch auraTriggerEntry.TriggerType {
		case define.AuraTrigger_SpellResult:
			if auraTriggerEntry.AddHead {
				c.listSpellResultTrigger.PushFront(auraTrigger)
			} else {
				c.listSpellResultTrigger.PushBack(auraTrigger)
			}

		case define.AuraTrigger_State:
			if auraTriggerEntry.TriggerMisc1 >= uint32(define.StateChangeMode_Begin) && auraTriggerEntry.TriggerMisc1 < uint32(define.StateChangeMode_End) {
				c.listServentStateTrigger[auraTriggerEntry.TriggerMisc1].PushBack(auraTrigger)
			}

		case define.AuraTrigger_Behaviour:
			for n := define.BehaviourType_BeforeMelee; n < define.BehaviourType_End; n++ {
				if (1<<n)&auraTriggerEntry.FamilyMask != 0 {
					c.listBehaviourTrigger[n].PushBack(auraTrigger)
				}
			}

		case define.AuraTrigger_AuraState:
			if auraTriggerEntry.TriggerMisc1 >= uint32(define.StateChangeMode_Begin) && auraTriggerEntry.TriggerMisc1 < uint32(define.StateChangeMode_End) {
				c.listAuraStateTrigger[auraTriggerEntry.TriggerMisc1].PushBack(auraTrigger)
			}

		default:
			continue
		}
	}
}

func (c *CombatCtrl) unRegisterAuraTrigger(aura *Aura) {
	if aura == nil {
		return
	}

	auraEntry := aura.Opts().Entry
	if auraEntry == nil {
		return
	}

	var index int32
	for index = 0; index < define.SpellEffectNum; index++ {
		if auraEntry.TriggerId[index] == 0 {
			continue
		}

		auraTriggerEntry := entries.GetAuraTriggerEntry(auraEntry.TriggerId[index])
		if auraTriggerEntry == nil {
			continue
		}

		if c.isDmgModEffect(auraEntry.Effects[index]) {
			c.unRegisterDmgMod(aura, index)
			continue
		}

		switch auraTriggerEntry.TriggerType {
		case define.AuraTrigger_SpellResult:
			var next *list.Element
			for e := c.listSpellResultTrigger.Front(); e != nil; e = next {
				next = e.Next()
				auraTrigger := e.Value.(*AuraTrigger)
				if auraTrigger.Aura == aura && auraTrigger.EffIndex == index {
					c.listSpellResultTrigger.Remove(e)
					break
				}
			}

		case define.AuraTrigger_State:
			listTrigger := c.listServentStateTrigger[auraTriggerEntry.TriggerMisc1]
			var next *list.Element
			for e := listTrigger.Front(); e != nil; e = next {
				next = e.Next()
				auraTrigger := e.Value.(*AuraTrigger)
				if auraTrigger.Aura == aura && auraTrigger.EffIndex == index {
					listTrigger.Remove(e)
					break
				}
			}

		case define.AuraTrigger_Behaviour:
			for n := define.BehaviourType_Begin; n < define.BehaviourType_End; n++ {
				if (1<<n)&auraTriggerEntry.TriggerMisc1 != 0 {
					listTrigger := c.listBehaviourTrigger[n]
					var next *list.Element
					for e := listTrigger.Front(); e != nil; e = next {
						next = e.Next()
						auraTrigger := e.Value.(*AuraTrigger)
						if auraTrigger.Aura == aura && auraTrigger.EffIndex == index {
							listTrigger.Remove(e)
							break
						}
					}
				}
			}

		case define.AuraTrigger_AuraState:
			listTrigger := c.listAuraStateTrigger[auraTriggerEntry.TriggerMisc1]
			var next *list.Element
			for e := listTrigger.Front(); e != nil; e = next {
				next = e.Next()
				auraTrigger := e.Value.(*AuraTrigger)
				if auraTrigger.Aura == aura && auraTrigger.EffIndex == index {
					listTrigger.Remove(e)
					break
				}
			}

		default:
			continue
		}
	}
}

//-------------------------------------------------------------------------------
// 计算aura效果
//-------------------------------------------------------------------------------
func (c *CombatCtrl) CalAuraEffect(curRound int32) {

	for n := 0; n < define.Combat_MaxAura; n++ {
		if c.arrayAura[n] != nil && !c.owner.HasState(define.HeroState_Freeze) {
			if c.arrayAura[n].Opts().Entry.RoundUpdateMask&(1<<curRound) == 0 {
				continue
			}

			c.arrayAura[n].CalAuraEffect(define.AuraEffectStep_Effect, -1, nil, nil)
		}
	}
}

//-------------------------------------------------------------------------------
// 清除aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) deleteAura(aura *Aura) {
	if aura == nil {
		return
	}

	removeMode := aura.GetRemoveMode()
	if aura.IsApplyLockValid() {
		if removeMode&define.AuraRemoveMode_Destroy == 0 {
			aura.CalAuraEffect(define.AuraEffectStep_Remove, -1, removeMode, nil)
		}
	}

	// 反注册
	c.UnregisterAura(aura)

	// 释放内存
	c.auraPool.Put(aura)
}

//-------------------------------------------------------------------------------
// 伤害触发器
//-------------------------------------------------------------------------------
func (c *CombatCtrl) registerDmgMod(aura *Aura, index int32) {
	if aura == nil || !(index >= 0 && index < define.SpellEffectNum) {
		return
	}

	auraEntry := aura.Opts().Entry
	if auraEntry == nil {
		return
	}

	for n := 0; n < define.Combat_DmgModTypeNum; n++ {
		if auraEntry.Effects[index] == define.Combat_DmgModType[n] {
			auraTrigger := &AuraTrigger{
				Aura:     aura,
				EffIndex: index,
			}

			c.listDmgModTrigger[n].PushBack(auraTrigger)
			return
		}
	}
}

func (c *CombatCtrl) unRegisterDmgMod(aura *Aura, index int32) {
	if aura == nil || !(index >= 0 && index < define.SpellEffectNum) {
		return
	}

	auraEntry := aura.Opts().Entry
	if auraEntry == nil {
		return
	}

	for n := 0; n < define.Combat_DmgModTypeNum; n++ {
		if auraEntry.Effects[index] == define.Combat_DmgModType[n] {
			listTrigger := c.listDmgModTrigger[n]
			var next *list.Element
			for e := listTrigger.Front(); e != nil; e = next {
				next = e.Next()
				auraTrigger := e.Value.(*AuraTrigger)
				if auraTrigger.Aura == aura && auraTrigger.EffIndex == index {
					listTrigger.Remove(e)
					return
				}
			}
		}
	}
}

func (c *CombatCtrl) isDmgModEffect(tp define.EAuraEffectType) bool {
	switch tp {
	case define.AuraEffectType_DmgMod:
		fallthrough
	case define.AuraEffectType_DmgFix:
		fallthrough
	case define.AuraEffectType_AbsorbAllDmg:
		return true
	}

	return false
}

func (c *CombatCtrl) TriggerByDmgMod(caster bool, target SceneUnit, dmgInfo *CalcDamageInfo) {
	for n := 0; n < define.Combat_DmgModTypeNum; n++ {
		listTrigger := c.listDmgModTrigger[n]
		for e := listTrigger.Front(); e != nil; e = e.Next() {
			auraTrigger := e.Value.(*AuraTrigger)

			// 是否已经废弃
			if auraTrigger.Aura == nil || auraTrigger.Aura.IsRemoved() || auraTrigger.Aura.IsHangup() {
				continue
			}

			triggerEntry := entries.GetAuraTriggerEntry(auraTrigger.Aura.Opts().Entry.TriggerId[auraTrigger.EffIndex])
			if triggerEntry == nil {
				continue
			}

			// 检查技能
			var spellBase *define.SpellBase = nil
			if spellEntry := entries.GetSpellEntry(dmgInfo.SpellId); spellEntry != nil {
				spellBase = &spellEntry.SpellBase
			} else if auraEntry := entries.GetAuraEntry(dmgInfo.SpellId); auraEntry != nil {
				spellBase = &auraEntry.SpellBase
			} else {
				continue
			}

			if triggerEntry.SpellId != 0 && triggerEntry.SpellId != dmgInfo.SpellId {
				continue
			}

			if triggerEntry.FamilyMask > 0 && (triggerEntry.FamilyMask&spellBase.FamilyMask == 0) {
				continue
			}

			if triggerEntry.FamilyRace > 0 && (triggerEntry.FamilyRace != spellBase.FamilyRace) {
				continue
			}

			if (triggerEntry.DmgInfoType > 0) && (triggerEntry.DmgInfoType&(1<<dmgInfo.Type) == 0) {
				continue
			}

			if triggerEntry.SchoolType != define.SchoolType_Null && (triggerEntry.SchoolType == dmgInfo.SchoolType) {
				continue
			}

			// 检查响应事件
			if uint32(define.AuraEventEx_Trigger_Always) != triggerEntry.TriggerMisc2 {
				if caster {
					if triggerEntry.TriggerMisc1 > 0 && (dmgInfo.ProcCaster&triggerEntry.TriggerMisc1 == 0) {
						continue
					}

					if triggerEntry.TriggerMisc2 > 0 && (dmgInfo.ProcEx&triggerEntry.TriggerMisc2 == 0) {
						continue
					}
				} else {
					if triggerEntry.TriggerMisc1 > 0 && (dmgInfo.ProcTarget&triggerEntry.TriggerMisc1 == 0) {
						continue
					}

					if triggerEntry.TriggerMisc2 > 0 && (dmgInfo.ProcEx&triggerEntry.TriggerMisc2 == 0) {
						continue
					}
				}

				// 验证触发条件
				if !c.checkTriggerCondition(triggerEntry, target) {
					continue
				}
			}

			// 计算触发几率
			if c.owner.GetScene().Rand(1, 10000) > int(triggerEntry.EventProp) {
				continue
			}

			// 作用效果
			auraTrigger.Aura.CalAuraEffect(define.AuraEffectStep_Effect, auraTrigger.EffIndex, dmgInfo, target)
		}
	}
}

//-------------------------------------------------------------------------------
// 触发器条件检查
//-------------------------------------------------------------------------------
func (c *CombatCtrl) checkTriggerCondition(auraTriggerEntry *define.AuraTriggerEntry, target SceneUnit) bool {
	if auraTriggerEntry == nil {
		return false
	}

	switch auraTriggerEntry.ConditionType {
	case define.AuraEventCondition_HPLowerFlat:
		if c.owner.Opts().AttManager.GetAttValue(define.Att_CurHP) < int64(auraTriggerEntry.ConditionMisc1) {
			return true
		}

	case define.AuraEventCondition_HPLowerPct:
		if c.owner.Opts().AttManager.GetAttValue(define.Att_CurHP)/c.owner.Opts().AttManager.GetAttValue(define.Att_MaxHP)*10000.0 < int64(auraTriggerEntry.ConditionMisc1) {
			return true
		}

	case define.AuraEventCondition_HPHigherFlat:
		if c.owner.Opts().AttManager.GetAttValue(define.Att_CurHP) >= int64(auraTriggerEntry.ConditionMisc1) {
			return true
		}

	case define.AuraEventCondition_HPHigherPct:
		if c.owner.Opts().AttManager.GetAttValue(define.Att_CurHP)/c.owner.Opts().AttManager.GetAttValue(define.Att_MaxHP)*10000.0 >= int64(auraTriggerEntry.ConditionMisc1) {
			return true
		}

	case define.AuraEventCondition_AnyUnitState:
		if c.owner.HasStateAny(uint32(auraTriggerEntry.ConditionMisc1)) {
			return true
		}

	case define.AuraEventCondition_AllUnitState:
		return true
		/*if c.owner.hasstateall(auratriggerentry.conditionmisc1) {
			return true
		}*/

	case define.AuraEventCondition_AuraState:
		if c.HasAuraState(auraTriggerEntry.ConditionMisc1) {
			return true
		}

	case define.AuraEventCondition_TargetClass:
		return true
		/*if target != nil && (1<<target.Opts().Entry.Class&auraTriggerEntry.ConditionMisc1 != 0) {
			return true
		}*/

	case define.AuraEventCondition_StrongTarget:
		if target != nil && target.Opts().AttManager.GetAttValue(define.Att_CurHP) > c.owner.Opts().AttManager.GetAttValue(define.Att_CurHP) {
			return true
		}

	case define.AuraEventCondition_TargetAuraState:
		if target != nil && target.CombatCtrl().HasAuraState(auraTriggerEntry.ConditionMisc1) {
			return true
		}

	case define.AuraEventCondition_TargetAllUnitState:
		return true
		/*if target != nil && target.HasStateAll(auraTriggerEntry.ConditionMisc1) {
			return true
		}*/

	case define.AuraEventCondition_TargetAnyUnitState:
		if target != nil && target.HasStateAny(uint32(auraTriggerEntry.ConditionMisc1)) {
			return true
		}

	case define.AuraEventCondition_None:
		return true

	default:
		return true
	}

	return false
}

//-------------------------------------------------------------------------------
// 按状态删除aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) removeAuraByState(state define.EHeroState) {
	for index := 0; index < define.Combat_MaxAura; index++ {
		if c.arrayAura[index] == nil || (c.arrayAura[index].GetRemoveMode()&define.AuraRemoveMode_Running == 0) {
			continue
		}

		auraEntry := c.arrayAura[index].Opts().Entry
		if auraEntry == nil {
			continue
		}

		// Add State:		1-Conflict	0-Ignore
		// Remove State:	1-Ignore	0-Dependence
		if auraEntry.OwnerStateCheckBitSet.Test(uint(state)) && auraEntry.OwnerStateLimitBitSet.Test(uint(state)) != c.owner.HasState(state) {
			c.RemoveAura(c.arrayAura[index], define.AuraRemoveMode_Dispel)
		}
	}
}

//-------------------------------------------------------------------------------
// 按施放者和ID删除aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) removeAuraByCaster(caster SceneUnit) {
	if caster == nil {
		return
	}

	for index := 0; index < define.Combat_MaxAura; index++ {
		aura := c.arrayAura[index]
		if aura == nil || (aura.GetRemoveMode()&define.AuraRemoveMode_Running) == 0 {
			continue
		}

		auraEntry := aura.Opts().Entry
		if auraEntry == nil {
			continue
		}

		if auraEntry.DependCaster && aura.opts.Caster == caster {
			c.RemoveAura(aura, define.AuraRemoveMode_Interrupt)
		}
	}
}

//-------------------------------------------------------------------------------
// 是否有指定TypeID的aura
//-------------------------------------------------------------------------------
func (c *CombatCtrl) GetAuraByIDCaster(auraId uint32, caster SceneUnit) *Aura {
	for index := 0; index < define.Combat_MaxAura; index++ {
		aura := c.arrayAura[index]
		if aura == nil {
			continue
		}

		if caster != nil {
			if aura.opts.Caster != caster {
				continue
			}
		}

		if aura.opts.Entry.ID == auraId {
			return aura
		}
	}

	return nil
}

//-------------------------------------------------------------------------------
// 驱散
//-------------------------------------------------------------------------------
func (c *CombatCtrl) DispelAura(dispelType uint32, num uint32) bool {
	// 是否驱散成功
	dispel := false

	for index := 0; index < define.Combat_MaxAura; index++ {
		aura := c.arrayAura[index]
		if aura == nil {
			continue
		}

		if aura.Opts().Entry.DispelFlags&dispelType != 0 {
			aura.Disperse()
			dispel = true

			if num--; num <= 0 {
				break
			}
		}
	}

	return dispel
}

//-------------------------------------------------------------------------------
// 强化作用时间
//-------------------------------------------------------------------------------
func (c *CombatCtrl) ModAuraDuration(modType uint32, modDuration uint32) {
	for index := 0; index < define.Combat_MaxAura; index++ {
		aura := c.arrayAura[index]
		if aura == nil {
			continue
		}

		if aura.Opts().Entry.DurationFlags&modType != 0 {
			aura.ModDuration(modDuration)
		}
	}
}

//-------------------------------------------------------------------------------
// 增益和减益buff数量
//-------------------------------------------------------------------------------
func (c *CombatCtrl) GetPositiveAndNegativeNum(posNum int32, negNum int32) (newPosNum int32, newNegNum int32) {
	newPosNum, newNegNum = posNum, negNum
	for index := 0; index < define.Combat_MaxAura; index++ {
		aura := c.arrayAura[index]
		if aura == nil {
			continue
		}

		// 被动buff不计数
		if aura.opts.Entry.Passive {
			continue
		}

		// 不可见buff不计数
		if !aura.Opts().Entry.HaveVisual {
			continue
		}

		if aura.Opts().Entry.EffectPriority > 0 {
			newPosNum++
		}

		if aura.Opts().Entry.EffectPriority < 0 {
			newNegNum++
		}
	}

	return
}

//-------------------------------------------------------------------------------
// aura state
//-------------------------------------------------------------------------------
func (c *CombatCtrl) AddAuraState(auraState int32) {
	if auraState < 0 || auraState >= define.AuraFlagNum {
		return
	}

	newState := !c.auraStateBitSet.Test(uint(auraState))
	c.auraStateBitSet.Set(uint(auraState))

	if newState {
		c.TriggerByAuraState(auraState, true)
	}
}

func (c *CombatCtrl) DecAuraState(auraState int32) {
	if auraState < 0 || auraState >= define.AuraFlagNum {
		return
	}

	if !c.auraStateBitSet.Test(uint(auraState)) {
		return
	}

	c.auraStateBitSet.Clear(uint(auraState))

	if !c.auraStateBitSet.Test(uint(auraState)) {
		c.TriggerByAuraState(auraState, false)
	}
}

func (c *CombatCtrl) HasAuraState(auraState int32) bool {
	return c.auraStateBitSet.Test(uint(auraState))
}

func (c *CombatCtrl) HasAuraStateAny(auraStateMask uint32) bool {
	compare := bitset.From([]uint64{uint64(auraStateMask)})
	return c.auraStateBitSet.Intersection(compare).Any()
}
