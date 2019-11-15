package item

import "github.com/yokaiio/yokai_server/game/define"

type defaultItem struct {
	id    int64
	entry *define.ItemEntry
}

func newDefaultItem() Item {
	return &defaultItem{}
}

func (i *defaultItem) Init() error {
	i.id = 1
	return nil
}

func (i *defaultItem) ID() int64 {
	return i.id
}

func (i *defaultItem) Entry() *define.ItemEntry {
	return i.entry
}
