package item

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
)

type DefaultItem struct {
	ID      int64 `gorm:"type:bigint(20);primary_key;column:id;default:0;not null"`
	OwnerID int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:0;not null"`
	TypeID  int32 `gorm:"type:int(10);column:type_id;default:0;not null"`
	entry   *define.ItemEntry
}

func defaultNewItem(id int64) Item {
	return &DefaultItem{
		ID: id,
	}
}

func defaultMigrate(ds *db.Datastore) {
	ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(DefaultItem{})
}

func (i *DefaultItem) TableName() string {
	return "item"
}

func (h *DefaultItem) GetID() int64 {
	return h.ID
}

func (h *DefaultItem) GetOwnerID() int64 {
	return h.OwnerID
}

func (h *DefaultItem) GetTypeID() int32 {
	return h.TypeID
}

func (h *DefaultItem) Entry() *define.ItemEntry {
	return h.entry
}

func (h *DefaultItem) SetOwnerID(id int64) {
	h.OwnerID = id
}

func (h *DefaultItem) SetTypeID(id int32) {
	h.TypeID = id
}

func (h *DefaultItem) SetEntry(e *define.ItemEntry) {
	h.entry = e
}
