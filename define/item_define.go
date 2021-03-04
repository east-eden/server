package define

const (
	Item_TypeItem    int32 = iota // 普通物品
	Item_TypeEquip                // 装备
	Item_TypePresent              // 礼包
	Item_TypeCrystal              // 晶石
)

const (
	Item_Effect_Null          = -1 // 无效果
	Item_Effect_Loot          = 0  // 掉落
	Item_Effect_CrystalDefine = 1  // 鉴定晶石

	Item_Effect_End = 2
)

const (
	Item_Quality_Begin  int32 = iota
	Item_Quality_White  int32 = iota - 1 // 白
	Item_Quality_Green                   // 绿
	Item_Quality_Blue                    // 蓝
	Item_Quality_Purple                  // 紫
	Item_Quality_Orange                  // 橙
	Item_Quality_Red                     // 红
	Item_Quality_End
)

const (
	Container_Null     int32 = -1
	Container_Begin    int32 = 0
	Container_Material int32 = 0 // 材料与消耗品
	Container_Equip    int32 = 1 // 装备背包
	Container_Crystal  int32 = 2 // 晶石背包
	Container_End      int32 = 3
)

// 装备位置
const (
	Equip_Pos_Begin   int32 = 0
	Equip_Pos_Weapon  int32 = 0 // 武器
	Equip_Pos_Clothes int32 = 1 // 衣服
	Equip_Pos_Shoe    int32 = 2 // 鞋子
	Equip_Pos_Jewel   int32 = 3 // 饰品
	Equip_Pos_End     int32 = 4
)
