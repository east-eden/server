package item

import "github.com/yokaiio/yokai_server/game/define"

type Item interface {
	Init() error
	ID() int64
	Entry() *define.ItemEntry
}

var (
	DefaultItem Item = newDefaultItem()
)

func NewItem() Item {
	return newDefaultItem()
}
