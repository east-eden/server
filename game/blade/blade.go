package blade

import (
	"context"
	"log"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/talent"
	"github.com/yokaiio/yokai_server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Blade struct {
	ID        int64              `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerID   int64              `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	OwnerType int32              `gorm:"type:int(10);column:owner_type;index:owner_type;default:-1;not null" bson:"owner_type"`
	TypeID    int32              `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`
	Exp       int64              `gorm:"type:bigint(20);column:exp;default:0;not null" bson:"exp"`
	Level     int32              `gorm:"type:int(10);column:level;default:1;not null" bson:"level"`
	Entry     *define.BladeEntry `gorm:"-" bson:"-"`

	talentManager *talent.TalentManager  `bson:"-"`
	wg            utils.WaitGroupWrapper `bson:"-"`
}

func newBlade(id int64, owner define.PluginObj, ds *db.Datastore) *Blade {
	b := &Blade{
		ID:        id,
		OwnerID:   owner.GetID(),
		OwnerType: owner.GetType(),
		TypeID:    -1,
		Exp:       0,
		Level:     1,
	}

	b.talentManager = talent.NewTalentManager(b, ds)
	return b
}

func defaultMigrate(ds *db.Datastore) {
	coll := ds.Database().Collection("blade")

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
			logger.Warn("collection blade create index owner_id failed:", err)
		}
	}
}

func (b *Blade) GetType() int32 {
	return define.Plugin_Blade
}

func (b *Blade) GetID() int64 {
	return b.ID
}

func (b *Blade) GetLevel() int32 {
	return b.Level
}

func (b *Blade) LoadFromDB() {
	b.wg.Wrap(b.talentManager.LoadFromDB)
	b.wg.Wait()
}

func (b *Blade) TalentManager() *talent.TalentManager {
	return b.talentManager
}
