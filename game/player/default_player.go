package player

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/hero"
	"github.com/yokaiio/yokai_server/game/item"
)

type defaultPlayer struct {
	DS          *db.Datastore
	itemManager *item.ItemManager
	heroManager *hero.HeroManager

	ID    int64  `gorm:"type:bigint(20);primary_key;column:id;default:0;not null"`
	Name  string `gorm:"type:varchar(32);column:name;not null"`
	Exp   int64  `gorm:"type:bigint(20);column:exp;default:0;not null"`
	Level int32  `gorm:"type:int(10);column:level;default:1;not null"`
}

func newDefaultPlayer(id int64, name string, ds *db.Datastore) Player {
	return &defaultPlayer{
		DS:          ds,
		ID:          id,
		Name:        name,
		Exp:         0,
		Level:       1,
		itemManager: item.NewItemManager(),
		heroManager: hero.NewHeroManager(),
	}
}

func defaultMigrate(ds *db.Datastore) {
	ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(defaultPlayer{})
	item.Migrate(ds)
	hero.Migrate(ds)
}

func (p *defaultPlayer) TableName() string {
	return "player"
}

func (p *defaultPlayer) HeroManager() *hero.HeroManager {
	return p.heroManager
}

func (p *defaultPlayer) ItemManager() *item.ItemManager {
	return p.itemManager
}

func (p *defaultPlayer) ChangeExp(add int64) {
	p.Exp += add

	p.DS.ORM().Model(p).Updates(defaultPlayer{
		Exp: p.Exp,
	})
}

func (p *defaultPlayer) ChangeLevel(add int32) {
	p.Level += add

	p.DS.ORM().Model(p).Updates(defaultPlayer{
		Level: p.Level,
	})
}
