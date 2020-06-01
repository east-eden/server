package item

import (
	"context"
	"sync"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/store"
	"go.mongodb.org/mongo-driver/bson"
)

// item create pool
var itemPool = &sync.Pool{New: newPoolItemV1}

func NewPoolItem() Item {
	return itemPool.Get().(Item)
}

func ReleasePoolItem(x interface{}) {
	itemPool.Put(x)
}

type Item interface {
	Options() *Options
	Entry() *define.ItemEntry
	EquipEnchantEntry() *define.EquipEnchantEntry
	GetAttManager() *att.AttManager

	GetEquipObj() int64
	SetEquipObj(int64)

	CalcAtt()
}

func NewItem(opts ...Option) Item {
	i := NewPoolItem()

	for _, o := range opts {
		o(i.Options())
	}

	return i
}

func Migrate(ds *store.Datastore) {
	migrateV1(ds)
}

func LoadAll(ds *store.Datastore, ownerID int64, tableName string) interface{} {
	list := make([]*ItemV1, 0)

	if ds == nil {
		return list
	}

	ctx, _ := context.WithTimeout(context.Background(), define.DatastoreTimeout)
	cur, err := ds.Database().Collection(tableName).Find(ctx, bson.D{{"owner_id", ownerID}})
	defer cur.Close(ctx)
	if err != nil {
		logger.Warn("item load all error:", err)
		return list
	}

	for cur.Next(ctx) {
		var i ItemV1
		if err := cur.Decode(&i); err != nil {
			logger.Warn("item decode failed:", err)
			continue
		}

		list = append(list, &i)
	}

	return list
}
