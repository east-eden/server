package hero

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
)

type DefaultHero struct {
	ID      int64 `gorm:"type:bigint(20);primary_key;column:id;default:0;not null"`
	OwnerID int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:0;not null"`
	TypeID  int32 `gorm:"type:int(10);column:type_id;default:0;not null"`
	entry   *define.HeroEntry
}

func defaultNewHero(id int64) Hero {
	return &DefaultHero{
		ID: id,
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

func (h *DefaultHero) GetOwnerID() int64 {
	return h.OwnerID
}

func (h *DefaultHero) GetTypeID() int32 {
	return h.TypeID
}

func (h *DefaultHero) Entry() *define.HeroEntry {
	return h.entry
}

func (h *DefaultHero) SetOwnerID(id int64) {
	h.OwnerID = id
}

func (h *DefaultHero) SetTypeID(id int32) {
	h.TypeID = id
}

func (h *DefaultHero) SetEntry(e *define.HeroEntry) {
	h.entry = e
}
