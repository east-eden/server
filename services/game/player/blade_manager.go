package player

import (
	"errors"
	"fmt"

	"bitbucket.org/east-eden/server/define"
	"bitbucket.org/east-eden/server/excel/auto"
	"bitbucket.org/east-eden/server/services/game/blade"
	"bitbucket.org/east-eden/server/services/game/talent"
	"bitbucket.org/east-eden/server/store"
	"bitbucket.org/east-eden/server/utils"
	log "github.com/rs/zerolog/log"
)

type BladeManager struct {
	owner    *Player                `bson:"-" json:"-"`
	BladeMap map[int64]*blade.Blade `bson:"blade_map" json:"blade_map"`
}

func NewBladeManager(owner *Player) *BladeManager {
	m := &BladeManager{
		owner:    owner,
		BladeMap: make(map[int64]*blade.Blade),
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
	err := store.GetStore().LoadObject(define.StoreType_Blade, m.owner.ID, m)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("BladeManager LoadAll: %w", err)
	}

	for _, b := range m.BladeMap {
		err := m.initLoadedBlade(b)
		if err != nil {
			return fmt.Errorf("BladeManager LoadAll: %w", err)
		}
	}

	return nil
}

func (m *BladeManager) createEntryBlade(entry *auto.BladeEntry) *blade.Blade {
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
	m.BladeMap[b.GetOptions().Id] = b
	b.GetAttManager().CalcAtt()

	return b
}

func (m *BladeManager) initLoadedBlade(b *blade.Blade) error {
	entry, _ := auto.GetBladeEntry(b.GetOptions().TypeId)

	if b.GetOptions().Entry == nil {
		return fmt.Errorf("blade<%d> entry invalid", b.GetOptions().TypeId)
	}

	b.GetOptions().Entry = entry

	m.BladeMap[b.GetOptions().Id] = b
	b.CalcAtt()
	return nil
}

func (m *BladeManager) GetBlade(id int64) (*blade.Blade, error) {
	if b, ok := m.BladeMap[id]; ok {
		return b, nil
	} else {
		return nil, fmt.Errorf("invalid blade id<%d>", id)
	}
}

func (m *BladeManager) GetBladeNums() int {
	return len(m.BladeMap)
}

func (m *BladeManager) GetBladeList() []*blade.Blade {
	list := make([]*blade.Blade, 0)

	for _, v := range m.BladeMap {
		list = append(list, v)
	}

	return list
}

func (m *BladeManager) AddBlade(typeId int32) *blade.Blade {
	bladeEntry, ok := auto.GetBladeEntry(typeId)
	if !ok {
		return nil
	}

	blade := m.createEntryBlade(bladeEntry)
	if blade == nil {
		return nil
	}

	fields := map[string]interface{}{}
	fields[fmt.Sprintf("blade_map.id_%d", blade.Id)] = blade
	err := store.GetStore().SaveFields(define.StoreType_Blade, m.owner.ID, fields)
	if pass := utils.ErrCheck(err, "AddBlade SaveObject failed", typeId); !pass {
		m.delBlade(blade)
		return nil
	}

	return blade
}

func (m *BladeManager) delBlade(b *blade.Blade) {
	delete(m.BladeMap, b.Id)
	blade.ReleasePoolBlade(b)
}

func (m *BladeManager) DelBlade(id int64) {
	b, ok := m.BladeMap[id]
	if !ok {
		return
	}

	fieldsName := []string{fmt.Sprintf("blade_map.id_%d", id)}
	err := store.GetStore().DeleteFields(define.StoreType_Blade, m.owner.ID, fieldsName)
	utils.ErrPrint(err, "DelBlade DeleteObject failed", id)
	m.delBlade(b)
}

func (m *BladeManager) BladeAddExp(id int64, exp int64) {
	b, ok := m.BladeMap[id]

	if ok {
		b.GetOptions().Exp += exp

		fields := map[string]interface{}{}
		fields[fmt.Sprintf("blade_map.id_%d.exp", id)] = b.GetOptions().Exp
		err := store.GetStore().SaveFields(define.StoreType_Blade, m.owner.ID, fields)
		utils.ErrPrint(err, "BladeAddExp SaveFields failed", id, exp)
	}
}

func (m *BladeManager) BladeAddLevel(id int64, level int32) {
	b, ok := m.BladeMap[id]

	if ok {
		b.GetOptions().Level += level

		fields := map[string]interface{}{}
		fields[fmt.Sprintf("blade_map.id_%d.level", id)] = b.GetOptions().Level
		err := store.GetStore().SaveFields(define.StoreType_Blade, m.owner.ID, fields)
		utils.ErrPrint(err, "BladeAddLevel SaveFields failed", id, level)
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
