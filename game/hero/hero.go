package hero

import (
	"context"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"go.mongodb.org/mongo-driver/bson"
)

type Hero interface {
	define.PluginObj

	Entry() *define.HeroEntry
	GetOwnerID() int64
	GetOwnerType() int32
	GetTypeID() int32
	GetExp() int64
	GetEquips() [define.Hero_MaxEquip]int64
	GetEquip(int32) int64

	SetOwnerID(int64)
	SetOwnerType(int32)
	SetTypeID(int32)
	SetExp(int64)
	SetLevel(int32)
	SetEntry(*define.HeroEntry)

	AddExp(int64) int64
	AddLevel(int32) int32
	BeforeDelete()
	SetEquip(int64, int32)
	UnsetEquip(int32)
}

func NewHero(id int64) Hero {
	return defaultNewHero(id)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}

func LoadAll(ds *db.Datastore, ownerID int64, tableName string) interface{} {
	{

	}

	list := make([]*HeroV1, 0)

	ctx, _ := context.WithTimeout(context.Background(), define.DatastoreTimeout)
	cur, err := ds.Database().Collection(tableName).Find(ctx, bson.D{{"owner_id", ownerID}})
	defer cur.Close(ctx)

	if err != nil {
		logger.Warn("hero loadall failed:", err)
		return list
	}

	for cur.Next(ctx) {
		var h HeroV1
		if err := cur.Decode(&h); err != nil {
			logger.Warn("hero decode failed:", err)
			continue
		}

		list = append(list, &h)
	}

	return list
}
