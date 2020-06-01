package hero

import (
	"context"
	"sync"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/rune"
	"github.com/yokaiio/yokai_server/game/store"
	"go.mongodb.org/mongo-driver/bson"
)

// hero create pool
var heroPool = &sync.Pool{New: newPoolHeroV1}

func NewPoolHero() Hero {
	return heroPool.Get().(Hero)
}

func ReleasePoolHero(x interface{}) {
	heroPool.Put(x)
}

type Hero interface {
	define.PluginObj

	Options() *Options
	GetEquipBar() *item.EquipBar
	GetAttManager() *att.AttManager
	GetRuneBox() *rune.RuneBox

	AddExp(int64) int64
	AddLevel(int32) int32
	BeforeDelete()
	CalcAtt()
}

func NewHero(opts ...Option) Hero {
	h := NewPoolHero()

	for _, o := range opts {
		o(h.Options())
	}

	return h
}

func Migrate(ds *store.Datastore) {
	migrateV1(ds)
}

func LoadAll(ds *store.Datastore, ownerID int64, tableName string) interface{} {
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
