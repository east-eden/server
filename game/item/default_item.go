package item

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/define"
	"github.com/yokaiio/yokai_server/game/global"
)

type DefaultItem struct {
	ID      int64 `gorm:"type:bigint(20);primary_key;column:id;default:0;not null"`
	OwnerID int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:0;not null"`
	TypeID  int32 `gorm:"type:int(10);column:type_id;default:0;not null"`
	entry   *define.ItemEntry
}

func defaultNewItem(id int64, ownerID int64, typeID int32) Item {
	return &DefaultItem{
		ID:      id,
		OwnerID: ownerID,
		TypeID:  typeID,
		entry:   global.GetItemEntry(typeID),
	}
}

func defaultMigrate(ds *db.Datastore) {
	ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(DefaultItem{})
}

func (i *DefaultItem) TableName() string {
	return "item"
}

func (i *DefaultItem) GetID() int64 {
	return i.ID
}

func (i *DefaultItem) GetTypeID() int32 {
	return i.TypeID
}

func (i *DefaultItem) Entry() *define.ItemEntry {
	return i.entry
}
