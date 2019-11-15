package item

import (
	"github.com/yokaiio/yokai_server/game/define"
	"github.com/yokaiio/yokai_server/game/global"
)

type defaultItem struct {
	id     int64
	typeID int32
	entry  *define.ItemEntry
}

func defaultNewItem(id int64, typeID int32) Item {
	return &defaultItem{
		id:     id,
		typeID: typeID,
		entry:  global.GetItemEntry(typeID),
	}
}

func (i *defaultItem) ID() int64 {
	return i.id
}

func (i *defaultItem) Entry() *define.ItemEntry {
	return i.entry
}
