package player

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/rune"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbGame "github.com/yokaiio/yokai_server/proto/game"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RuneManager struct {
	owner   *Player
	mapRune map[int64]*rune.Rune

	ds   *db.Datastore
	coll *mongo.Collection
	sync.RWMutex
}

func NewRuneManager(owner *Player, ds *db.Datastore) *RuneManager {
	m := &RuneManager{
		owner:   owner,
		mapRune: make(map[int64]*rune.Rune, 0),
		ds:      ds,
	}

	m.coll = ds.Database().Collection(m.TableName())
	return m
}

func (m *RuneManager) TableName() string {
	return "rune"
}

func (m *RuneManager) save(r *rune.Rune) {
	go func() {
		filter := bson.D{{"_id", r.ID}}
		update := bson.D{{"$set", r}}
		op := options.Update().SetUpsert(true)
		m.coll.UpdateOne(context.Background(), filter, update, op)
	}()
}

func (m *RuneManager) delete(id int64) {
	go func() {
		m.coll.DeleteOne(context.Background(), bson.D{{"_id", id}})
	}()
}

func (m *RuneManager) createRune(typeID int32) *rune.Rune {
	runeEntry := global.GetRuneEntry(typeID)
	r := m.createEntryRune(runeEntry)
	if r == nil {
		logger.Warning("new rune failed when createRune:", typeID)
		return nil
	}

	m.mapRune[r.GetID()] = r
	m.save(r)

	return r
}

func (m *RuneManager) delRune(id int64) {
	r, ok := m.mapRune[id]
	if !ok {
		return
	}

	r.SetEquipObj(-1)
	delete(m.mapRune, id)
	m.delete(id)
}

func (m *RuneManager) createEntryRune(entry *define.RuneEntry) *rune.Rune {
	if entry == nil {
		logger.Error("createEntryRune with nil RuneEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Rune)
	if err != nil {
		logger.Error(err)
		return nil
	}

	r := rune.NewRune(id)
	r.SetOwnerID(m.owner.GetID())
	r.SetTypeID(entry.ID)
	r.SetEntry(entry)

	m.mapRune[r.GetID()] = r

	return i
}

func (m *RuneManager) createDBRune(r *rune.Rune) *rune.Rune {
	newRune := rune.NewRune(r.GetID())
	newRune.SetOwnerID(r.GetOwnerID())
	newRune.SetTypeID(r.GetTypeID())
	newRune.SetEquipObj(i.GetEquipObj())

	entry := global.GetRuneEntry(i.GetTypeID())
	newRune.SetEntry(entry)

	m.mapRune[newRune.GetID()] = newRune

	return newRune
}

// interface of cost_loot
func (m *RuneManager) GetCostLootType() int32 {
	return define.CostLoot_Rune
}

func (m *RuneManager) CanCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("rune manager check item<%d> cost failed, wrong number<%d>", typeMisc, num)
	}

	var fixNum int32 = 0
	for _, v := range m.mapRune {
		if v.GetTypeID() == typeMisc && v.GetEquipObj() != -1 {
			fixNum += 1
		}
	}

	if fixNum >= num {
		return nil
	}

	return fmt.Errorf("not enough rune<%d>, num<%d>", typeMisc, num)
}

func (m *RuneManager) DoCost(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("rune manager cost item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	return m.CostRuneByTypeID(typeMisc, num)
}

func (m *RuneManager) CanGain(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("rune manager check gain item<%d> failed, wrong number<%d>", typeMisc, num)
	}

	// todo bag max item

	return nil
}

func (m *RuneManager) GainLoot(typeMisc int32, num int32) error {
	if num <= 0 {
		return fmt.Errorf("rune manager gain rune<%d> failed, wrong number<%d>", typeMisc, num)
	}

	return m.AddRuneByTypeID(typeMisc, num)
}

func (m *RuneManager) LoadFromDB() {
	l := rune.LoadAll(m.ds, m.owner.GetID(), m.TableName())
	sliceItem := make([]*rune.Rune, 0)

	listRune := reflect.ValueOf(l)
	if listRune.Kind() != reflect.Slice {
		logger.Error("load rune returns non-slice type")
		return
	}

	for n := 0; n < listRune.Len(); n++ {
		p := listRune.Index(n)
		sliceRune = append(sliceRune, p.Interface().(*rune.Rune))
	}

	for _, v := range sliceRune {
		m.createDBRune(v)
	}
}

func (m *RuneManager) GetRune(id int64) *rune.Rune {
	return m.mapRune[id]
}

func (m *RuneManager) GetRuneNums() int {
	return len(m.mapRune)
}

func (m *RuneManager) GetRuneList() []*rune.Rune {
	list := make([]*rune.Rune, 0)

	for _, v := range m.mapRune {
		list = append(list, v)
	}

	return list
}

func (m *RuneManager) AddRuneByTypeID(typeID int32, num int32) error {
	if num <= 0 {
		return nil
	}

	for n := 0; n < num; n++ {
		if r := m.createRune(typeID); r == nil {
			break
		}

		m.SendRuneAdd(i)
	}

	return nil
}

func (m *RuneManager) DeleteRune(id int64) error {
	if i := m.GetRune(id); i == nil {
		return fmt.Errorf("cannot find rune<%d> while DeleteRune", id)
	}

	m.delRune(id)
	m.SendRuneDelete(id)

	return nil
}

func (m *RuneManager) CostRuneByTypeID(typeID int32, num int32) error {
	if num < 0 {
		return fmt.Errorf("dec rune error, invalid number:%d", num)
	}

	decNum := num
	for _, v := range m.mapRune {
		if decNum <= 0 {
			break
		}

		if v.Entry().ID == typeID {
			decNum--
			delID := v.GetID()
			m.delRune(delID)
			m.SendRuneDelete(delID)
		}
	}

	if decNum > 0 {
		logger.WithFields(logger.Fields{
			"need_dec":   num,
			"actual_dec": num - decNum,
		}).Warning("CostRuneByTypeID warning")
	}

	return nil
}

func (m *RuneManager) CostRuneByID(id int64) error {
	r := m.GetRune(id)
	if r == nil {
		return fmt.Errorf("cannot find rune by id:%d", id)
	}

	m.delRune(id)
	m.SendRuneDelete(id)

	return nil
}

func (m *RuneManager) SetRuneEquiped(id int64, objID int64) {
	r, ok := m.mapRune[id]
	if !ok {
		return
	}

	r.SetEquipObj(objID)
	m.save(r)
	m.SendRuneUpdate(r)
}

func (m *RuneManager) SetRuneUnEquiped(id int64) {
	i, ok := m.mapRune[id]
	if !ok {
		return
	}

	i.SetEquipObj(-1)
	m.save(i)
	m.SendRuneUpdate(i)
}

func (m *RuneManager) SendRuneAdd(r *rune.Rune) {
	msg := &pbGame.M2C_RuneAdd{
		Rune: &pbGame.Rune{
			Id:     r.GetID(),
			TypeId: r.GetTypeID(),
		},
	}

	m.owner.SendProtoMessage(msg)
}

func (m *RuneManager) SendRuneDelete(id int64) {
	msg := &pbGame.M2C_DelRune{
		RuneId: id,
	}

	m.owner.SendProtoMessage(msg)
}

func (m *RuneManager) SendRuneUpdate(r *rune.Rune) {
	msg := &pbGame.M2C_RuneUpdate{
		Rune: &pbGame.Rune{
			Id:     r.GetID(),
			TypeId: r.GetTypeID(),
		},
	}

	m.owner.SendProtoMessage(msg)
}
