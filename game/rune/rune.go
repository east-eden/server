package rune

import (
	"context"
	"log"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Rune struct {
	ID      int64 `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerID int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	TypeID  int32 `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`

	EquipObj   int64             `gorm:"type:bigint(20);column:equip_obj;default:-1;not null" bson:"equip_obj"`
	entry      *define.RuneEntry `gorm:"-" bson:"-"`
	attManager *att.AttManager   `gorm:"-" bson:"-"`
}

func NewRune(id int64) *Rune {
	return &Rune{
		ID:       id,
		EquipObj: -1,
	}
}

func Migrate(ds *db.Datastore) {
	coll := ds.Database().Collection("Rune")

	// creck index
	idx := coll.Indexes()

	opts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := idx.List(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	indexExist := false
	for cursor.Next(context.Background()) {
		var result bson.M
		cursor.Decode(&result)
		if result["name"] == "owner_id" {
			indexExist = true
			break
		}
	}

	// create index
	if !indexExist {
		_, err := coll.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bsonx.Doc{{"owner_id", bsonx.Int32(1)}},
				Options: options.Index().SetName("owner_id"),
			},
		)

		if err != nil {
			logger.Warn("collection Rune create index owner_id failed:", err)
		}
	}
}

func LoadAll(ds *db.Datastore, ownerID int64, tableName string) []*Rune {
	list := make([]*Rune, 0)

	ctx, _ := context.WithTimeout(context.Background(), define.DatastoreTimeout)
	cur, err := ds.Database().Collection(tableName).Find(ctx, bson.D{{"owner_id", ownerID}})
	defer cur.Close(ctx)
	if err != nil {
		logger.Warn("rune load all error:", err)
		return list
	}

	for cur.Next(ctx) {
		var r Rune
		if err := cur.Decode(&r); err != nil {
			logger.Warn("rune decode failed:", err)
			continue
		}

		list = append(list, &r)
	}

	return list
}

func (r *Rune) GetID() int64 {
	return r.ID
}

func (r *Rune) GetOwnerID() int64 {
	return r.OwnerID
}

func (r *Rune) GetTypeID() int32 {
	return r.TypeID
}

func (r *Rune) GetEquipObj() int64 {
	return r.EquipObj
}

func (r *Rune) GetAttManager() *att.AttManager {
	return r.attManager
}

func (r *Rune) Entry() *define.RuneEntry {
	return r.entry
}

func (r *Rune) SetOwnerID(id int64) {
	r.OwnerID = id
}

func (r *Rune) SetTypeID(id int32) {
	r.TypeID = id
}

func (r *Rune) SetEquipObj(id int64) {
	r.EquipObj = id
}

func (r *Rune) SetEntry(e *define.RuneEntry) {
	r.entry = e
}

func (r *Rune) SetAttManager(m *att.AttManager) {
	r.attManager = m
}
