package item

import (
	"context"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	"github.com/yokaiio/yokai_server/game/att"
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
	"go.mongodb.org/mongo-driver/bson"
)

type Item interface {
	Entry() *define.ItemEntry
	EquipEnchantEntry() *define.EquipEnchantEntry
	GetID() int64
	GetOwnerID() int64
	GetTypeID() int32
	GetNum() int32
	GetEquipObj() int64
	GetAttManager() *att.AttManager

	SetOwnerID(int64)
	SetTypeID(int32)
	SetNum(int32)
	SetEquipObj(int64)
	SetEntry(*define.ItemEntry)
	SetEquipEnchantEntry(*define.EquipEnchantEntry)
	SetAttManager(*att.AttManager)
}

func NewItem(id int64) Item {
	return defaultNewItem(id)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}

func LoadAll(ds *db.Datastore, ownerID int64, tableName string) interface{} {
	list := make([]*ItemV1, 0)

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
