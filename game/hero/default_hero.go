package hero

import (
	"github.com/yokaiio/yokai_server/game/db"
	"github.com/yokaiio/yokai_server/internal/define"
)

type DefaultHero struct {
	ID        int64                       `gorm:"type:bigint(20);primary_key;column:id;default:-1;not null" bson:"_id"`
	OwnerID   int64                       `gorm:"type:bigint(20);column:owner_id;index:owner_id;default:-1;not null" bson:"owner_id"`
	OwnerType int32                       `gorm:"type:int(10);column:owner_type;index:owner_type;default:-1;not null" bson:"owner_type"`
	TypeID    int32                       `gorm:"type:int(10);column:type_id;default:-1;not null" bson:"type_id"`
	Exp       int64                       `gorm:"type:bigint(20);column:exp;default:0;not null" bson:"exp"`
	Level     int32                       `gorm:"type:int(10);column:level;default:1;not null" bson:"level"`
	Equips    [define.Hero_MaxEquip]int64 `gorm:"-" bson:"-"`
	entry     *define.HeroEntry           `gorm:"-" bson:"-"`
}

func defaultNewHero(id int64) Hero {
	return &DefaultHero{
		ID:        id,
		OwnerID:   -1,
		OwnerType: -1,
		TypeID:    -1,
		Exp:       0,
		Level:     1,
		Equips:    [define.Hero_MaxEquip]int64{-1, -1, -1, -1},
	}
}

func defaultMigrate(ds *db.Datastore) {
	//ds.ORM().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(DefaultHero{})
}

func (h *DefaultHero) GetType() int32 {
	return define.Plugin_Hero
}

func (h *DefaultHero) GetID() int64 {
	return h.ID
}

func (h *DefaultHero) GetLevel() int32 {
	return h.Level
}

func (h *DefaultHero) GetOwnerID() int64 {
	return h.OwnerID
}

func (h *DefaultHero) GetOwnerType() int32 {
	return h.OwnerType
}

func (h *DefaultHero) GetTypeID() int32 {
	return h.TypeID
}

func (h *DefaultHero) GetExp() int64 {
	return h.Exp
}

func (h *DefaultHero) GetEquips() [define.Hero_MaxEquip]int64 {
	return h.Equips
}

func (h *DefaultHero) GetEquip(pos int32) int64 {
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return -1
	}

	return h.Equips[pos]
}

func (h *DefaultHero) Entry() *define.HeroEntry {
	return h.entry
}

func (h *DefaultHero) SetOwnerID(id int64) {
	h.OwnerID = id
}

func (h *DefaultHero) SetOwnerType(tp int32) {
	h.OwnerType = tp
}

func (h *DefaultHero) SetTypeID(id int32) {
	h.TypeID = id
}

func (h *DefaultHero) SetExp(exp int64) {
	h.Exp = exp
}

func (h *DefaultHero) SetLevel(level int32) {
	h.Level = level
}

func (h *DefaultHero) SetEntry(e *define.HeroEntry) {
	h.entry = e
}

func (h *DefaultHero) AddExp(exp int64) int64 {
	h.Exp += exp
	return h.Exp
}

func (h *DefaultHero) AddLevel(level int32) int32 {
	h.Level += level
	return h.Level
}

func (h *DefaultHero) BeforeDelete() {

}

func (h *DefaultHero) SetEquip(equipID int64, pos int32) {
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return
	}

	h.Equips[pos] = equipID
}

func (h *DefaultHero) UnsetEquip(pos int32) {
	if pos < 0 || pos >= define.Hero_MaxEquip {
		return
	}

	h.Equips[pos] = -1
}
