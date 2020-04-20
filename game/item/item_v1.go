package item

import (
	"context"
	"log"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/mgo.v2/bson"
)

type ItemV1 struct {
	ID      int64 `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerID int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	TypeID  int32 `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`
	Num     int32 `gorm:"type:int(10);column:num;default:0;not null" bson:"num"`

	EquipObj          int64                     `gorm:"type:bigint(20);column:equip_obj;default:-1;not null" bson:"equip_obj"`
	entry             *define.ItemEntry         `gorm:"-" bson:"-"`
	equipEnchantEntry *define.EquipEnchantEntry `gorm:"-" bson:"-"`
	attManager        *att.AttManager           `gorm:"-" bson:"-"`
}

func defaultNewItem(id int64) Item {
	return &ItemV1{
		ID:       id,
		Num:      1,
		EquipObj: -1,
	}
}

func defaultMigrate(ds *db.Datastore) {
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

func (h *ItemV1) GetID() int64 {
	return h.ID
}

func (h *ItemV1) GetOwnerID() int64 {
	return h.OwnerID
}

func (h *ItemV1) GetTypeID() int32 {
	return h.TypeID
}

func (h *ItemV1) GetNum() int32 {
	return h.Num
}

func (h *ItemV1) GetEquipObj() int64 {
	return h.EquipObj
}

func (h *ItemV1) GetAttManager() *att.AttManager {
	return h.attManager
}

func (h *ItemV1) Entry() *define.ItemEntry {
	return h.entry
}

func (h *ItemV1) EquipEnchantEntry() *define.EquipEnchantEntry {
	return h.equipEnchantEntry
}

func (h *ItemV1) SetOwnerID(id int64) {
	h.OwnerID = id
}

func (h *ItemV1) SetTypeID(id int32) {
	h.TypeID = id
}

func (h *ItemV1) SetNum(num int32) {
	h.Num = num
}

func (h *ItemV1) SetEquipObj(id int64) {
	h.EquipObj = id
}

func (h *ItemV1) SetEntry(e *define.ItemEntry) {
	h.entry = e
}

func (h *ItemV1) SetEquipEnchantEntry(e *define.EquipEnchantEntry) {
	h.equipEnchantEntry = e
}

func (h *ItemV1) SetAttManager(m *att.AttManager) {
	h.attManager = m
}
