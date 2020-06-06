package item

import (
	"time"

	"github.com/yokaiio/yokai_server/define"
	"github.com/yokaiio/yokai_server/game/att"
)

type ItemV1 struct {
	Opts       *Options        `bson:"inline"`
	attManager *att.AttManager `gorm:"-" bson:"-"`
}

func newPoolItemV1() interface{} {
	h := &ItemV1{
		Opts: DefaultOptions(),
	}

	h.attManager = att.NewAttManager(-1)

	return h
}

func (i *ItemV1) Options() *Options {
	return i.Opts
}

func (i *ItemV1) GetID() int64 {
	return i.Opts.Id
}

func (i *ItemV1) GetOwnerID() int64 {
	return i.Opts.OwnerId
}

func (i *ItemV1) GetTypeID() int32 {
	return i.Opts.TypeId
}

func (i *ItemV1) GetAttManager() *att.AttManager {
	return i.attManager
}

func (i *ItemV1) Entry() *define.ItemEntry {
	return i.Opts.Entry
}

func (i *ItemV1) EquipEnchantEntry() *define.EquipEnchantEntry {
	return i.Opts.EquipEnchantEntry
}

func (i *ItemV1) GetEquipObj() int64 {
	return i.Opts.EquipObj
}

func (i *ItemV1) SetEquipObj(obj int64) {
	i.Opts.EquipObj = obj
}

func (i *ItemV1) AfterLoad() {

}

func (i *ItemV1) GetExpire() *time.Timer {
	return nil
}

func (i *ItemV1) GetObjID() interface{} {
	return i.Opts.Id
}

func (i *ItemV1) TableName() string {
	return "item"
}

func (i *ItemV1) CalcAtt() {
	i.attManager.Reset()
}
