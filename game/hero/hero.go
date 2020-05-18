package hero

import (
	"context"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/rune"
	"go.mongodb.org/mongo-driver/bson"
)

type Hero interface {
	define.PluginObj

	Entry() *define.HeroEntry
	GetOwnerID() int64
	GetOwnerType() int32
	GetTypeID() int32
	GetExp() int64
	GetEquipBar() *item.EquipBar
	GetAttManager() *att.AttManager
	GetRuneBox() *rune.RuneBox

	SetOwnerID(int64)
	SetOwnerType(int32)
	SetTypeID(int32)
	SetExp(int64)
	SetLevel(int32)
	SetEntry(*define.HeroEntry)
	SetAttManager(*att.AttManager)
	SetRuneBox(*rune.RuneBox)
	SetEquipBar(*item.EquipBar)

	AddExp(int64) int64
	AddLevel(int32) int32
	BeforeDelete()
	CalcAtt()
}

func NewHero(id int64) Hero {
	return defaultNewHero(id)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}

func LoadAll(ds *db.Datastore, ownerID int64, tableName string) interface{} {
	list := make([]*HeroV1, 0)

	if ds == nil {
		return list
	}

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
