package talent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"github.com/yokaiio/yokai_server/internal/global"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TalentManager struct {
	Owner      define.PluginObj `gorm:"-" bson:"-"`
	OwnerID    int64            `gorm:"type:bigint(20);primary_key;column:owner_id;index:owner_id;default:-1;not null" bson:"_id"`
	OwnerType  int32            `gorm:"type:int(10);primary_key;column:owner_type;index:owner_type;default:-1;not null" bson:"owner_type"`
	TalentJson string           `gorm:"type:varchar(5120);column:talent_json" bson:"-"`
	Talents    []*Talent        `json:"talents" bson:"talents"`

	ds           *db.Datastore `bson:"-"`
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

	// init talents
	//m.initTalents()

	return m
}

func (m *TalentManager) TableName() string {
	return "talent"
}

func Migrate(ds *db.Datastore) {
	//ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(TalentManager{})
}

func (m *TalentManager) initTalents() {
	//for n := 0; n < define.Talent_End; n++ {
	//m.Talents = append(m.Talents, &Talent{
	//ID:      int32(n),
	//Value:   0,
	//MaxHold: 100000000,
	//entry:   global.GetTalentEntry(int32(n)),
	//})
	//}
}

func (m *TalentManager) LoadFromDB() {
	res := m.ds.Database().Collection(m.TableName()).FindOne(context.Background(), bson.M{"_id": m.OwnerID})
	if res.Err() == mongo.ErrNoDocuments {
		m.ds.Database().Collection(m.TableName()).InsertOne(context.Background(), m)
	} else {
		res.Decode(m)
	}

	// unmarshal json to talent value
	if len(m.TalentJson) > 0 {
		err := json.Unmarshal([]byte(m.TalentJson), &m.Talents)
		if err != nil {
			logger.Error("unmarshal talent json failed:", err)
		}
	}

	// init entry
	for _, v := range m.Talents {
		v.entry = global.GetTalentEntry(int32(v.ID))
	}
}

func (m *TalentManager) Save() error {
	data, err := json.Marshal(m.Talents)
	if err != nil {
		return fmt.Errorf("json marshal failed:%s", err.Error())
	}

	m.TalentJson = string(data)

	filter := bson.D{{"_id", m.OwnerID}}
	update := bson.D{{"$set", m}}
	op := options.Update().SetUpsert(true)
	m.ds.Database().Collection(m.TableName()).UpdateOne(context.Background(), filter, update, op)
	return nil
}

func (m *TalentManager) AddTalent(id int32) error {
	t := &Talent{ID: id, entry: global.GetTalentEntry(int32(id))}

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

	filter := bson.D{{"_id", m.OwnerID}}
	update := bson.D{{"$set", m}}
	op := options.Update().SetUpsert(true)
	m.ds.Database().Collection(m.TableName()).UpdateOne(context.Background(), filter, update, op)
	return nil
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
