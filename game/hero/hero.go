package hero

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
)

type Hero interface {
	Entry() *define.HeroEntry
	GetID() int64
	GetOwnerID() int64
	GetTypeID() int32

	SetOwnerID(int64)
	SetTypeID(int32)
	SetEntry(*define.HeroEntry)
}

func NewHero(id int64) Hero {
	return defaultNewHero(id)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}

func LoadAll(ds *db.Datastore, ownerID int64) interface{} {
	list := make([]*DefaultHero, 0)
	ds.ORM().Where("owner_id = ?", ownerID).Find(&list)
	return list
}
