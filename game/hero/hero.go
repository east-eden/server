package hero

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
)

type Hero interface {
	define.PluginObj

	Entry() *define.HeroEntry
	GetOwnerID() int64
	GetTypeID() int32
	GetExp() int64
	GetEquips() [define.Hero_MaxEquip]int64
	GetEquip(int32) int64

	SetOwnerID(int64)
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

func LoadAll(ds *db.Datastore, ownerID int64) interface{} {
	list := make([]*DefaultHero, 0)
	ds.ORM().Where("owner_id = ?", ownerID).Find(&list)
	return list
}
