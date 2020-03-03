package item

import (
	"context"
	"log"
	"time"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/mgo.v2/bson"
)

type DefaultItem struct {
	ID      int64 `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerID int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	TypeID  int32 `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`
	Num     int32 `gorm:"type:int(10);column:num;default:0;not null" bson:"num"`

	EquipObj          int64                     `gorm:"type:bigint(20);column:equip_obj;default:-1;not null" bson:"equip_obj"`
	entry             *define.ItemEntry         `gorm:"-" bson:"-"`
	equipEnchantEntry *define.EquipEnchantEntry `gorm:"-" bson:"-"`
}

func defaultNewItem(id int64) Item {
	return &DefaultItem{
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

func (h *DefaultItem) GetID() int64 {
	return h.ID
}

func (h *DefaultItem) GetOwnerID() int64 {
	return h.OwnerID
}

func (h *DefaultItem) GetTypeID() int32 {
	return h.TypeID
}

func (h *DefaultItem) GetNum() int32 {
	return h.Num
}

func (h *DefaultItem) GetEquipObj() int64 {
	return h.EquipObj
}

func (h *DefaultItem) Entry() *define.ItemEntry {
	return h.entry
}

func (h *DefaultItem) EquipEnchantEntry() *define.EquipEnchantEntry {
	return h.equipEnchantEntry
}

func (h *DefaultItem) SetOwnerID(id int64) {
	h.OwnerID = id
}

func (h *DefaultItem) SetTypeID(id int32) {
	h.TypeID = id
}

func (h *DefaultItem) SetNum(num int32) {
	h.Num = num
}

func (h *DefaultItem) SetEquipObj(id int64) {
	h.EquipObj = id
}

func (h *DefaultItem) SetEntry(e *define.ItemEntry) {
	h.entry = e
}

func (h *DefaultItem) SetEquipEnchantEntry(e *define.EquipEnchantEntry) {
	h.equipEnchantEntry = e
}
