package item

import (
	"context"
	"log"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type ItemV1 struct {
	Opts       *Options        `bson:"inline"`
	attManager *att.AttManager `gorm:"-" bson:"-"`
}

func newPoolItemV1() interface{} {
	h := &ItemV1{
		Opts: DefaultOptions(),
	}

	h.attManager = att.NewAttManager(-1)

	return h
}

func migrateV1(ds *db.Datastore) {
	coll := ds.Database().Collection("item")

	// check index
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
			logger.Warn("collection item create index owner_id failed:", err)
		}
	}
}

func (i *ItemV1) Options() *Options {
	return i.Opts
}

func (i *ItemV1) GetID() int64 {
	return i.Opts.Id
}

func (i *ItemV1) GetOwnerID() int64 {
	return i.Opts.OwnerId
}

func (i *ItemV1) GetTypeID() int32 {
	return i.Opts.TypeId
}

func (i *ItemV1) GetAttManager() *att.AttManager {
	return i.attManager
}

func (i *ItemV1) Entry() *define.ItemEntry {
	return i.Opts.Entry
}

func (i *ItemV1) EquipEnchantEntry() *define.EquipEnchantEntry {
	return i.Opts.EquipEnchantEntry
}

func (i *ItemV1) GetEquipObj() int64 {
	return i.Opts.EquipObj
}

func (i *ItemV1) SetEquipObj(obj int64) {
	i.Opts.EquipObj = obj
}

func (i *ItemV1) CalcAtt() {
	i.attManager.Reset()
}
