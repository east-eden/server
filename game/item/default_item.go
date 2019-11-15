package item

type defaultItem struct {
	id int64
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
