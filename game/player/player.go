package player

import (
	"context"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/yokaiio/yokai_server/game/blade"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/token"
	"github.com/yokaiio/yokai_server/internal/define"
	"gopkg.in/mgo.v2/bson"
)

type Player interface {
	define.PluginObj

	LoadFromDB()
	AfterLoad()
	Save()

	GetClientID() int64
	GetName() string
	GetExp() int64

	SetClientID(int64)
	SetName(string)
	SetExp(int64)
	SetLevel(int32)

	HeroManager() *hero.HeroManager
	ItemManager() *item.ItemManager
	TokenManager() *token.TokenManager
	BladeManager() *blade.BladeManager

	ChangeExp(int64)
	ChangeLevel(int32)
}

// player proto
func NewPlayer(id int64, name string, db *db.Datastore) Player {
	return newDefaultPlayer(id, name, db)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}

func LoadAll(ds *db.Datastore, tableName string) interface{} {
	list := make([]*DefaultPlayer, 0)

	ctx, _ := context.WithTimeout(context.Background(), define.DatastoreTimeout)
	cur, err := ds.Database().Collection(tableName).Find(ctx, bson.D{})
	defer cur.Close(ctx)

	if err != nil {
		logger.Warn("player load all failed:", err)
		return list
	}

	for cur.Next(ctx) {
		var p DefaultPlayer
		if err := cur.Decode(&p); err != nil {
			logger.Warn("player decode failed:", err)
			continue
		}

		list = append(list, &p)
	}

	return list
}
