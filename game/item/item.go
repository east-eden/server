package item

type Item interface {
	Init() error
	ID() int64
}

var (
	DefaultItem defaultItem = newDefaultItem()
)

func NewItem() Item {
	return newDefaultItem()
}
