package define

const (
	Item_TypeItem    = iota // 普通物品
	Item_TypeEquip          // 装备
	Item_TypePresent        // 礼包
	Item_TypeRune           // 御魂(未鉴定)
)

const (
	Item_Effect_Null       = -1 // 无效果
	Item_Effect_Loot       = 0  // 掉落
	Item_Effect_RuneDefine = 1  // 鉴定御魂

	Item_Effect_End = 2
)
