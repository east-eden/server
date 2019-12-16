package item

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
)

type Item interface {
	Entry() *define.ItemEntry
	GetID() int64
	GetOwnerID() int64
	GetTypeID() int32
	GetNum() int32
	GetEquipObj() int64

	SetOwnerID(int64)
	SetTypeID(int32)
	SetNum(int32)
	SetEquipObj(int64)
	SetEntry(*define.ItemEntry)
}

func NewItem(id int64) Item {
	return defaultNewItem(id)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}

func LoadAll(ds *db.Datastore, ownerID int64) interface{} {
	list := make([]*DefaultItem, 0)
	ds.ORM().Where("owner_id = ?", ownerID).Find(&list)
	return list
}
