package player

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
)

type Player interface {
	TableName() string
	LoadFromDB()

	GetID() int64
	GetName() string
	GetExp() int64
	GetLevel() int32

	SetName(string)
	SetExp(int64)
	SetLevel(int32)

	HeroManager() *hero.HeroManager
	ItemManager() *item.ItemManager

	ChangeExp(int64)
	ChangeLevel(int32)
}

// player proto
func NewPlayer(id int64, db *db.Datastore) Player {
	return newDefaultPlayer(id, db)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}

func LoadAll(ds *db.Datastore) interface{} {
	list := make([]*DefaultPlayer, 0)
	ds.ORM().Find(&list)
	return list
}
