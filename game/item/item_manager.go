package item

type ItemManager struct {
	mapItem map[int64]Item
}

func NewItemManager() *ItemManager {
	return &ItemManager{
		mapItem: make(map[int64]Item, 0),
	}
}
