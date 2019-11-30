package item

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/game/define"
	"github.com/yokaiio/yokai_server/game/global"
)

type defaultItem struct {
	ID     int64 `gorm:"type:bigint(20);primary_key;column:id;default:0;not null"`
	TypeID int32 `gorm:"type:int(10);column:type_id;default:0;not null"`
	entry  *define.ItemEntry
}

func defaultNewItem(id int64, typeID int32) Item {
	return &defaultItem{
		ID:     id,
		TypeID: typeID,
		entry:  global.GetItemEntry(typeID),
	}
}

func defaultMigrate(ds *db.Datastore) {
	ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(defaultItem{})
}

func (i *defaultItem) TableName() string {
	return "item"
}

func (i *defaultItem) GetID() int64 {
	return i.ID
}

func (i *defaultItem) GetTypeID() int32 {
	return i.TypeID
}

func (i *defaultItem) Entry() *define.ItemEntry {
	return i.entry
}
