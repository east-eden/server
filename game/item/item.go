package item

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/define"
)

type Item interface {
	GetID() int64
	GetTypeID() int32
	Entry() *define.ItemEntry
}

func NewItem(id int64, typeID int32) Item {
	return defaultNewItem(id, typeID)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}
