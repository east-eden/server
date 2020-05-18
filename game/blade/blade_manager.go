package blade

import (
	"context"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BladeManager struct {
	Owner    define.PluginObj
	mapBlade map[int64]*Blade

	ds   *db.Datastore
	coll *mongo.Collection
	sync.RWMutex
	wg utils.WaitGroupWrapper
}

func NewBladeManager(obj define.PluginObj, ds *db.Datastore) *BladeManager {
	m := &BladeManager{
		Owner:    obj,
		ds:       ds,
		mapBlade: make(map[int64]*Blade, 0),
	}

	if ds != nil {
		m.coll = ds.Database().Collection(m.TableName())
	}

	return m
}

func Migrate(ds *db.Datastore) {
	//ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(Blade{})
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

func (m *BladeManager) LoadFromDB() {
	ctx, _ := context.WithTimeout(context.Background(), define.DatastoreTimeout)
	cur, err := m.coll.Find(ctx, bson.D{{"owner_id", m.Owner.GetID()}})
	defer cur.Close(ctx)

	if err != nil {
		logger.Warn("blade_manager load from db error:", err)
		return
	}

	for cur.Next(ctx) {
		var b Blade
		if err := cur.Decode(&b); err != nil {
			logger.Warn("blade_manager decode failed:", err)
			continue
		}

		m.newDBBlade(&b)
	}

	if err := cur.Err(); err != nil {
		logger.Fatal(err)
	}

	m.wg.Wait()
}

func (m *BladeManager) newEntryBlade(entry *define.BladeEntry) *Blade {
	if entry == nil {
		logger.Error("newEntryBlade with nil BladeEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Blade)
	if err != nil {
		logger.Error(err)
		return nil
	}

	blade := newBlade(id, m.Owner, m.ds)
	blade.OwnerID = m.Owner.GetID()
	blade.TypeID = entry.ID
	blade.Entry = entry

	m.mapBlade[blade.GetID()] = blade

	return blade
}

func (m *BladeManager) newDBBlade(b *Blade) *Blade {
	blade := newBlade(b.GetID(), m.Owner, m.ds)
	blade.OwnerID = m.Owner.GetID()
	blade.TypeID = b.TypeID
	blade.Entry = entries.GetBladeEntry(b.TypeID)

	m.mapBlade[blade.GetID()] = blade

	// load from db
	m.wg.Wrap(blade.LoadFromDB)

	return blade
}

func (m *BladeManager) GetBlade(id int64) *Blade {
	return m.mapBlade[id]
}

func (m *BladeManager) GetBladeNums() int {
	return len(m.mapBlade)
}

func (m *BladeManager) GetBladeList() []*Blade {
	list := make([]*Blade, 0)

	m.RLock()
	for _, v := range m.mapBlade {
		list = append(list, v)
	}
	m.RUnlock()

	return list
}

func (m *BladeManager) AddBlade(typeID int32) *Blade {
	bladeEntry := entries.GetBladeEntry(typeID)
	blade := m.newEntryBlade(bladeEntry)
	if blade == nil {
		return nil
	}

	filter := bson.D{{"_id", blade.GetID()}}
	update := bson.D{{"$set", blade}}
	op := options.Update().SetUpsert(true)
	m.coll.UpdateOne(context.Background(), filter, update, op)
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
		blade.Exp += exp

		filter := bson.D{{"_id", blade.GetID()}}
		update := bson.D{{"$set",
			bson.D{
				{"exp", blade.Exp},
			},
		}}
		m.coll.UpdateOne(context.Background(), filter, update)
	}
}

func (m *BladeManager) BladeAddLevel(id int64, level int32) {
	blade, ok := m.mapBlade[id]

	if ok {
		blade.Level += level

		filter := bson.D{{"_id", blade.GetID()}}
		update := bson.D{{"$set",
			bson.D{
				{"level", blade.Level},
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
