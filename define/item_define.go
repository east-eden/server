package define

type ItemType int

const (
	Item_TypeItem    ItemType = iota // 普通物品
	Item_TypeEquip                   // 装备
	Item_TypePresent                 // 礼包
	Item_TypeRune                    // 御魂(未鉴定)
)

const (
	Item_Effect_Null       = -1 // 无效果
	Item_Effect_Loot       = 0  // 掉落
	Item_Effect_RuneDefine = 1  // 鉴定御魂

	Item_Effect_End = 2
)

type ItemQualityType int

const (
	Item_Quality_Begin  ItemQualityType = 0
	Item_Quality_White  ItemQualityType = 0 // 白
	Item_Quality_Green  ItemQualityType = 1 // 绿
	Item_Quality_Blue   ItemQualityType = 2 // 蓝
	Item_Quality_Purple ItemQualityType = 3 // 紫
	Item_Quality_Orange ItemQualityType = 4 // 橙
	Item_Quality_End
)

type ContainerType int

const (
	Container_Null     ContainerType = -1
	Container_Begin    ContainerType = 0
	Container_Material ContainerType = 0 // 材料与消耗品
	Container_Equip    ContainerType = 1 // 装备背包
	Container_Crystal  ContainerType = 2 // 晶石背包
	Container_End
)
