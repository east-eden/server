package item

import (
	"time"

	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
)

type ItemV1 struct {
	Options    `bson:"inline" json:",inline"`
	attManager *att.AttManager `json:"-" bson:"-"`
}

func newPoolItemV1() interface{} {
	h := &ItemV1{
		Options: DefaultOptions(),
	}

	h.attManager = att.NewAttManager(-1)

	return h
}

// StoreObjector interface
func (i *ItemV1) AfterLoad() {

}

func (i *ItemV1) GetExpire() *time.Timer {
	return nil
}

func (i *ItemV1) GetObjID() interface{} {
	return i.Options.Id
}

func (i *ItemV1) TableName() string {
	return "item"
}

func (i *ItemV1) GetOptions() *Options {
	return &i.Options
}

func (i *ItemV1) GetID() int64 {
	return i.Options.Id
}

func (i *ItemV1) GetOwnerID() int64 {
	return i.Options.OwnerId
}

func (i *ItemV1) GetTypeID() int32 {
	return i.Options.TypeId
}

func (i *ItemV1) GetAttManager() *att.AttManager {
	return i.attManager
}

func (i *ItemV1) Entry() *define.ItemEntry {
	return i.Options.Entry
}

func (i *ItemV1) EquipEnchantEntry() *define.EquipEnchantEntry {
	return i.Options.EquipEnchantEntry
}

func (i *ItemV1) GetEquipObj() int64 {
	return i.Options.EquipObj
}

func (i *ItemV1) SetEquipObj(obj int64) {
	i.Options.EquipObj = obj
}

func (i *ItemV1) CalcAtt() {
	i.attManager.Reset()
}
