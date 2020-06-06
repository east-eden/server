package player

import (
	"context"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/blade"
	"github.com/yokaiio/yokai_server/store"
	"github.com/yokaiio/yokai_server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BladeManager struct {
	owner    *Player
	mapBlade map[int64]blade.Blade

	coll *mongo.Collection
	sync.RWMutex
}

func NewBladeManager(owner *Player) *BladeManager {
	m := &BladeManager{
		owner:    owner,
		mapBlade: make(map[int64]blade.Blade, 0),
	}

	return m
}

func (m *BladeManager) TableName() string {
	return "blade"
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

func (m *BladeManager) LoadAll() {
	bladeList, err := m.owner.store.LoadArrayFromCacheAndDB(store.StoreType_Blade, "owner_id", m.owner.GetID(), blade.GetBladePool())
	if err != nil {
		logger.Error("load blade manager failed:", err)
	}

	for _, i := range bladeList {
		err := m.initLoadedBlade(i.(blade.Blade))
		if err != nil {
			logger.Error("load blade failed:", err)
		}
	}
}

func (m *BladeManager) createEntryBlade(entry *define.BladeEntry) blade.Blade {
	if entry == nil {
		logger.Error("createEntryBlade with nil BladeEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Blade)
	if err != nil {
		logger.Error(err)
		return nil
	}

	b := blade.NewBlade(
		blade.Id(id),
		blade.OwnerId(m.owner.GetID()),
		blade.Entry(entry),
		blade.TypeId(entry.ID),
	)

	b.GetAttManager().SetBaseAttId(entry.AttID)
	m.mapBlade[b.Options().Id] = b
	b.GetAttManager().CalcAtt()

	return b
}

func (m *BladeManager) initLoadedBlade(b blade.Blade) error {
	entry := entries.GetBladeEntry(b.Options().TypeId)

	if b.Options().Entry == nil {
		return fmt.Errorf("blade<%d> entry invalid", b.Options().TypeId)
	}

	b.Options().Entry = entry
	b.GetAttManager().SetBaseAttId(entry.AttID)

	m.mapBlade[b.Options().Id] = b
	b.CalcAtt()
	return nil
}

func (m *BladeManager) GetBlade(id int64) blade.Blade {
	return m.mapBlade[id]
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
	bladeEntry := entries.GetBladeEntry(typeId)
	blade := m.createEntryBlade(bladeEntry)
	if blade == nil {
		return nil
	}

	m.owner.store.SaveObjectToCacheAndDB(store.StoreType_Blade, blade)
	return blade
}

func (m *BladeManager) DelBlade(id int64) {
	_, ok := m.mapBlade[id]
	if !ok {
		return
	}

	delete(m.mapBlade, id)

	filter := bson.D{{"_id", id}}
	m.coll.DeleteOne(context.Background(), filter)
}

func (m *BladeManager) BladeAddExp(id int64, exp int64) {
	blade, ok := m.mapBlade[id]

	if ok {
		blade.Options().Exp += exp

		filter := bson.D{{"_id", blade.Options().Id}}
		update := bson.D{{"$set",
			bson.D{
				{"exp", blade.Options().Exp},
			},
		}}
		m.coll.UpdateOne(context.Background(), filter, update)
	}
}

func (m *BladeManager) BladeAddLevel(id int64, level int32) {
	blade, ok := m.mapBlade[id]

	if ok {
		blade.Options().Level += level

		filter := bson.D{{"_id", blade.Options().Id}}
		update := bson.D{{"$set",
			bson.D{
				{"level", blade.Options().Level},
			},
		}}
		m.coll.UpdateOne(context.Background(), filter, update)
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
