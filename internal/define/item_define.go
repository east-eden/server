package define

const (
	Item_TypeItem    = iota // 普通物品
	Item_TypeEquip          // 装备
	Item_TypePresent        // 礼包
	Item_TypeVerb           // 御魂
)

const (
	Item_TypeEx_VerbNotDefined = 0 // 未鉴定御魂
	Item_TypeEx_VerbDefined    = 1 // 已鉴定御魂
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
}
