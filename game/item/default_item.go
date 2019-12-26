package item

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
)

type DefaultItem struct {
	ID      int64 `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerID int64 `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	TypeID  int32 `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`
	Num     int32 `gorm:"type:int(10);column:num;default:0;not null" bson:"num"`

	EquipObj int64             `gorm:"type:bigint(20);column:equip_obj;default:-1;not null" bson:"equip_obj"`
	entry    *define.ItemEntry `gorm:"-" bson:"-"`
}

func defaultNewItem(id int64) Item {
	return &DefaultItem{
		ID:       id,
		Num:      1,
		EquipObj: -1,
	}
}

func defaultMigrate(ds *db.Datastore) {
	//ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(DefaultItem{})
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

func (h *DefaultItem) GetNum() int32 {
	return h.Num
}

func (h *DefaultItem) GetEquipObj() int64 {
	return h.EquipObj
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

func (h *DefaultItem) SetNum(num int32) {
	h.Num = num
}

func (h *DefaultItem) SetEquipObj(id int64) {
	h.EquipObj = id
}

func (h *DefaultItem) SetEntry(e *define.ItemEntry) {
	h.entry = e
}
