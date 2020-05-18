package talent

import (
	"context"
	"fmt"
	"sync"

	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/entries"
	"github.com/yokaiio/yokai_server/game/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Talent struct {
	ID    int32               `json:"id" bson:"talent_id"`
	entry *define.TalentEntry `bson:"-"`
}

type TalentManager struct {
	Owner     define.PluginObj `gorm:"-" bson:"-"`
	OwnerID   int64            `gorm:"type:bigint(20);primary_key;column:owner_id;index:owner_id;default:-1;not null" bson:"_id"`
	OwnerType int32            `gorm:"type:int(10);primary_key;column:owner_type;index:owner_type;default:-1;not null" bson:"owner_type"`
	Talents   []*Talent        `json:"talents" bson:"talents"`

	ds           *db.Datastore     `bson:"-"`
	coll         *mongo.Collection `bson:"-"`
	sync.RWMutex `bson:"-"`
}

func NewTalentManager(owner define.PluginObj, ds *db.Datastore) *TalentManager {
	m := &TalentManager{
		Owner:     owner,
		OwnerID:   owner.GetID(),
		OwnerType: owner.GetType(),
		ds:        ds,
		Talents:   make([]*Talent, 0),
	}

	if ds != nil {
		m.coll = ds.Database().Collection(m.TableName())
	}

	return m
}

func (m *TalentManager) TableName() string {
	return "talent"
}

func (m *TalentManager) LoadFromDB() {
	res := m.coll.FindOne(context.Background(), bson.D{{"_id", m.OwnerID}})
	if res.Err() == mongo.ErrNoDocuments {
		m.coll.InsertOne(context.Background(), m)
	} else {
		res.Decode(m)
	}

	// init entry
	for _, v := range m.Talents {
		v.entry = entries.GetTalentEntry(int32(v.ID))
	}
}

func (m *TalentManager) save() error {
	filter := bson.D{{"_id", m.OwnerID}}
	update := bson.D{{"$set", m}}
	op := options.FindOneAndUpdate().SetUpsert(true)
	res := m.coll.FindOneAndUpdate(context.Background(), filter, update, op)
	return res.Err()
}

func (m *TalentManager) AddTalent(id int32) error {
	t := &Talent{ID: id, entry: entries.GetTalentEntry(int32(id))}

	if t.entry == nil {
		return fmt.Errorf("add not exist talent entry:%d", id)
	}

	if m.Owner.GetLevel() < t.entry.LevelLimit {
		return fmt.Errorf("level limit:%d", t.entry.LevelLimit)
	}

	// check group_id
	for _, v := range m.Talents {
		if v.ID == t.ID {
			return fmt.Errorf("add existed talent:%d", id)
		}

		// check group_id
		if t.entry.GroupID != -1 && t.entry.GroupID == v.entry.GroupID {
			return fmt.Errorf("talent group_id conflict:%d", t.entry.GroupID)
		}
	}

	m.Talents = append(m.Talents, t)

	return m.save()
}

func (m *TalentManager) GetTalent(id int32) *Talent {

	for _, v := range m.Talents {
		if v.ID == id {
			return v
		}
	}

	return nil
}

func (m *TalentManager) GetTalentList() []*Talent {

	return m.Talents
}
