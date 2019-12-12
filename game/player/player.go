package player

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
	"github.com/yokaiio/yokai_server/game/talent"
	"github.com/yokaiio/yokai_server/game/token"
)

type Player interface {
	TableName() string
	LoadFromDB()
	AfterLoad()
	Save()

	GetClientID() int64
	GetID() int64
	GetName() string
	GetExp() int64
	GetLevel() int32

	SetClientID(int64)
	SetName(string)
	SetExp(int64)
	SetLevel(int32)

	HeroManager() *hero.HeroManager
	ItemManager() *item.ItemManager
	TokenManager() *token.TokenManager
	TalentManager() *talent.TalentManager

	ChangeExp(int64)
	ChangeLevel(int32)
}

// player proto
func NewPlayer(id int64, name string, db *db.Datastore) Player {
	return newDefaultPlayer(id, name, db)
}

func Migrate(ds *db.Datastore) {
	defaultMigrate(ds)
}

func LoadAll(ds *db.Datastore) interface{} {
	list := make([]*DefaultPlayer, 0)
	ds.ORM().Find(&list)
	return list
}
