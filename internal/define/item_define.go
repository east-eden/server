package define

const (
	Item_TypeItem = iota
	Item_TypeEquip
)

// item entry
type ItemEntry struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	ItemType int32  `json:"item_type"`
	MaxStack int32  `json:"max_stack"`
	EquipPos int32  `json:"equip_pos"`
	Quality  int32  `json:"quality"`
	Price    int32  `json:"price"`
}
