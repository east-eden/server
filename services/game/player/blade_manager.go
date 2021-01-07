package player

import (
	"errors"
	"fmt"
	"sync"

	"github.com/east-eden/server/define"
	"github.com/east-eden/server/excel/auto"
	"github.com/east-eden/server/services/game/blade"
	"github.com/east-eden/server/services/game/talent"
	"github.com/east-eden/server/store"
	"github.com/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

type BladeManager struct {
	owner    *Player
	mapBlade map[int64]blade.Blade

	sync.RWMutex
}

func NewBladeManager(owner *Player) *BladeManager {
	m := &BladeManager{
		owner:    owner,
		mapBlade: make(map[int64]blade.Blade, 0),
	}

	return m
}

// interface of cost_loot
func (m *BladeManager) GetCostLootType() int32 {
	return define.CostLoot_Blade
}

func (m *BladeManager) CanCost(typeMisc int32, num int32) error {
	return nil
}

func (m *BladeManager) DoCost(typeMisc int32, num int32) error {
	return nil
}

func (m *BladeManager) CanGain(typeMisc int32, num int32) error {
	return nil
}

func (m *BladeManager) GainLoot(typeMisc int32, num int32) error {
	return nil
}

func (m *BladeManager) LoadAll() error {
	bladeList, err := store.GetStore().LoadArray(define.StoreType_Blade, m.owner.GetID(), blade.GetBladePool())
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("BladeManager LoadAll: %w", err)
	}

	for _, i := range bladeList {
		err := m.initLoadedBlade(i.(blade.Blade))
		if err != nil {
			return fmt.Errorf("BladeManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *BladeManager) createEntryBlade(entry *auto.BladeEntry) blade.Blade {
	if entry == nil {
		log.Error().Msg("createEntryBlade with nil BladeEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Blade)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}

	b := blade.NewBlade(
		blade.Id(id),
		blade.OwnerId(m.owner.GetID()),
		blade.Entry(entry),
		blade.TypeId(entry.Id),
	)

	// blade's talent
	tm := talent.NewTalentManager(b)
	b.SetTalentManager(tm)

	// b.GetAttManager().SetBaseAttId(entry.AttId)
	m.mapBlade[b.GetOptions().Id] = b
	b.GetAttManager().CalcAtt()

	return b
}

func (m *BladeManager) initLoadedBlade(b blade.Blade) error {
	entry, _ := auto.GetBladeEntry(b.GetOptions().TypeId)

	if b.GetOptions().Entry == nil {
		return fmt.Errorf("blade<%d> entry invalid", b.GetOptions().TypeId)
	}

	b.GetOptions().Entry = entry
	// b.GetAttManager().SetBaseAttId(entry.AttID)

	m.mapBlade[b.GetOptions().Id] = b
	b.CalcAtt()
	return nil
}

func (m *BladeManager) GetBlade(id int64) (blade.Blade, error) {
	if b, ok := m.mapBlade[id]; ok {
		return b, nil
	} else {
		return nil, fmt.Errorf("invalid blade id<%d>", id)
	}
}

func (m *BladeManager) GetBladeNums() int {
	return len(m.mapBlade)
}

func (m *BladeManager) GetBladeList() []blade.Blade {
	list := make([]blade.Blade, 0)

	m.RLock()
	for _, v := range m.mapBlade {
		list = append(list, v)
	}
	m.RUnlock()

	return list
}

func (m *BladeManager) AddBlade(typeId int32) blade.Blade {
	bladeEntry, ok := auto.GetBladeEntry(typeId)
	if !ok {
		return nil
	}

	blade := m.createEntryBlade(bladeEntry)
	if blade == nil {
		return nil
	}

	store.GetStore().SaveObject(define.StoreType_Blade, blade)
	return blade
}

func (m *BladeManager) DelBlade(id int64) {
	b, ok := m.mapBlade[id]
	if !ok {
		return
	}

	delete(m.mapBlade, id)
	store.GetStore().DeleteObject(define.StoreType_Blade, b)
	blade.ReleasePoolBlade(b)
}

func (m *BladeManager) BladeAddExp(id int64, exp int64) {
	b, ok := m.mapBlade[id]

	if ok {
		b.GetOptions().Exp += exp

		fields := map[string]interface{}{
			"exp": b.GetOptions().Exp,
		}
		store.GetStore().SaveFields(define.StoreType_Blade, b, fields)
	}
}

func (m *BladeManager) BladeAddLevel(id int64, level int32) {
	b, ok := m.mapBlade[id]

	if ok {
		b.GetOptions().Level += level

		fields := map[string]interface{}{
			"level": b.GetOptions().Level,
		}
		store.GetStore().SaveFields(define.StoreType_Blade, b, fields)
	}
}

func (m *BladeManager) PutonEquip(bladeID int64, equipID int64) error {
	/*blade, ok := m.mapBlade[bladeID]*/
	//if !ok {
	//return fmt.Errorf("invalid bladeid")
	/*}*/

	return nil
}

func (m *BladeManager) TakeoffEquip(bladeID int64) error {
	/*blade, ok := m.mapBlade[bladeID]*/
	//if !ok {
	//return fmt.Errorf("invalid blade_id")
	/*}*/

	return nil
}
