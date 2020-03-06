package define

const (
	Item_TypeItem    = iota // 普通物品
	Item_TypeEquip          // 装备
	Item_TypePresent        // 礼包
	Item_TypeRune           // 御魂
)

const (
	Item_TypeEx_RuneNotDefined = 0 // 未鉴定御魂
	Item_TypeEx_RuneDefined    = 1 // 已鉴定御魂
)

const (
	Item_Effect_Null       = -1 // 无效果
	Item_Effect_Loot       = 0  // 掉落
	Item_Effect_RuneDefine = 1  // 鉴定御魂

	Item_Effect_End = 2
)

// item entry
type ItemEntry struct {
	ID          int32   `json:"_id"`
	Name        string  `json:"Name"`
	Desc        string  `json:"Desc"`
	Type        int32   `json:"Type"`
	SubType     int32   `json:"SubType"`
	Quality     int32   `json:"Quality"`
	MaxStack    int32   `json:"MaxStack"`
	EffectType  int32   `json:"EffectType"`
	EffectValue []int32 `json:"EffectValue"`

	EquipEnchantID int32 `json:"EquipEnchantID"`
}
