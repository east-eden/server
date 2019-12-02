package hero

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/define"
	"github.com/yokaiio/yokai_server/game/global"
)

type DefaultHero struct {
	ID      int64 `gorm:"type:bigint(20);primary_key;column:id;default:0;not null"`
	OwnerID int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:0;not null"`
	TypeID  int32 `gorm:"type:int(10);column:type_id;default:0;not null"`
	entry   *define.HeroEntry
}

func defaultNewHero(id int64, ownerID int64, typeID int32) Hero {
	return &DefaultHero{
		ID:      id,
		OwnerID: ownerID,
		TypeID:  typeID,
		entry:   global.GetHeroEntry(typeID),
	}
}

func defaultMigrate(ds *db.Datastore) {
	ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(DefaultHero{})
}

func (h *DefaultHero) TableName() string {
	return "hero"
}

func (h *DefaultHero) GetID() int64 {
	return h.ID
}

func (h *DefaultHero) Entry() *define.HeroEntry {
	return h.entry
}
