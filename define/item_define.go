package define

type ItemType int

const (
	Item_TypeItem    ItemType = iota // 普通物品
	Item_TypeEquip                   // 装备
	Item_TypePresent                 // 礼包
	Item_TypeCrystal                 // 晶石
)

const (
	Item_Effect_Null          = -1 // 无效果
	Item_Effect_Loot          = 0  // 掉落
	Item_Effect_CrystalDefine = 1  // 鉴定晶石

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
	Item_Quality_End                    = 5
)

type ContainerType int

const (
	Container_Null     ContainerType = -1
	Container_Begin    ContainerType = 0
	Container_Material ContainerType = 0 // 材料与消耗品
	Container_Equip    ContainerType = 1 // 装备背包
	Container_Crystal  ContainerType = 2 // 晶石背包
	Container_End      ContainerType = 3
)

// 装备位置
type EquipPosType int

const (
	Equip_Pos_Begin   EquipPosType = 0
	Equip_Pos_Weapon  EquipPosType = 0 // 武器
	Equip_Pos_Clothes EquipPosType = 1 // 衣服
	Equip_Pos_Shoe    EquipPosType = 2 // 鞋子
	Equip_Pos_Jewel   EquipPosType = 3 // 饰品
	Equip_Pos_End     EquipPosType = 4
)
