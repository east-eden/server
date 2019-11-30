package player

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
)

type Player interface {
	TableName() string

	HeroManager() *hero.HeroManager
	ItemManager() *item.ItemManager

	ChangeExp(int64)
	ChangeLevel(int32)
}

func NewPlayer(id int64, name string, db *db.Datastore) Player {
	return newDefaultPlayer(id, name, db)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}
