package hero

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/define"
)

type Hero interface {
	GetID() int64
	Entry() *define.HeroEntry
}

func NewHero(id int64, typeID int32) Hero {
	return defaultNewHero(id, typeID)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}
