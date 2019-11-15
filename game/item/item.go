package item

import "github.com/yokaiio/yokai_server/game/define"

type Item interface {
	ID() int64
	Entry() *define.ItemEntry
}

func NewItem(id int64, typeID int32) Item {
	return newDefaultItem(id, typeID)
}
