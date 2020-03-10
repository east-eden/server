package hero

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

type HeroV1 struct {
	ID        int64                       `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerID   int64                       `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	OwnerType int32                       `gorm:"type:int(10);column:owner_type;index:owner_type;default:-1;not null" bson:"owner_type"`
	TypeID    int32                       `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`
	Exp       int64                       `gorm:"type:bigint(20);column:exp;default:0;not null" bson:"exp"`
	Level     int32                       `gorm:"type:int(10);column:level;default:1;not null" bson:"level"`
	Equips    [define.Hero_MaxEquip]int64 `gorm:"-" bson:"-"`
	entry     *define.HeroEntry           `gorm:"-" bson:"-"`
}

func defaultNewHero(id int64) Hero {
	return &HeroV1{
		ID:        id,
		OwnerID:   -1,
		OwnerType: -1,
		TypeID:    -1,
		Exp:       0,
		Level:     1,
		Equips:    [define.Hero_MaxEquip]int64{-1, -1, -1, -1},
	}
}

func defaultMigrate(ds *db.Datastore) {
	coll := ds.Database().Collection("hero")

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
			logger.Warn("collection hero create index owner_id failed:", err)
		}
	}

}

func (h *HeroV1) GetType() int32 {
	return define.Plugin_Hero
}

func (h *HeroV1) GetID() int64 {
	return h.ID
}

func (h *HeroV1) GetLevel() int32 {
	return h.Level
}

func (h *HeroV1) GetOwnerID() int64 {
	return h.OwnerID
}

func (h *HeroV1) GetOwnerType() int32 {
	return h.OwnerType
}

func (h *HeroV1) GetTypeID() int32 {
	return h.TypeID
}

func (h *HeroV1) GetExp() int64 {
	return h.Exp
}

func (h *HeroV1) GetEquips() [define.Hero_MaxEquip]int64 {
	return h.Equips
}

func (h *HeroV1) GetEquip(pos int32) int64 {
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return -1
	}

	return h.Equips[pos]
}

func (h *HeroV1) Entry() *define.HeroEntry {
	return h.entry
}

func (h *HeroV1) SetOwnerID(id int64) {
	h.OwnerID = id
}

func (h *HeroV1) SetOwnerType(tp int32) {
	h.OwnerType = tp
}

func (h *HeroV1) SetTypeID(id int32) {
	h.TypeID = id
}

func (h *HeroV1) SetExp(exp int64) {
	h.Exp = exp
}

func (h *HeroV1) SetLevel(level int32) {
	h.Level = level
}

func (h *HeroV1) SetEntry(e *define.HeroEntry) {
	h.entry = e
}

func (h *HeroV1) AddExp(exp int64) int64 {
	h.Exp += exp
	return h.Exp
}

func (h *HeroV1) AddLevel(level int32) int32 {
	h.Level += level
	return h.Level
}

func (h *HeroV1) BeforeDelete() {

}

func (h *HeroV1) SetEquip(equipID int64, pos int32) {
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return
	}

	h.Equips[pos] = equipID
}

func (h *HeroV1) UnsetEquip(pos int32) {
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return
	}

	h.Equips[pos] = -1
}
