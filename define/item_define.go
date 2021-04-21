package define

// 物品主类型
const (
	Item_TypeItem    int32 = iota // 0 普通物品
	Item_TypeEquip                // 1 装备
	Item_TypeCrystal              // 2 晶石
	Item_TypePresent              // 3 礼包
)

// 物品子类型
const (
	// 物品子类型
	Item_SubType_Item_Begin      int32 = 0
	Item_SubType_Item_Normal     int32 = 0 // 0 普通物品
	Item_SubType_Item_EquipExp   int32 = 1 // 装备经验
	Item_SubType_Item_CrystalExp int32 = 2 // 晶石经验
	Item_SubType_Item_HeroExp    int32 = 3 // 英雄经验
	Item_SubType_Item_End        int32 = 4

	// 装备子类型
	Item_SubType_Equip_Begin        int32 = 0
	Item_SubType_Equip_OneHandSword int32 = 0 // 0 单手剑
	Item_SubType_Equip_TwoHandSword int32 = 1 // 1 双手剑
	Item_SubType_Equip_Polearm      int32 = 2 // 2 长柄
	Item_SubType_Equip_Pistol       int32 = 3 // 3 手枪
	Item_SubType_Equip_Staff        int32 = 4 // 4 法杖
	Item_SubType_Equip_Summon       int32 = 5 // 5 召唤物
	Item_SubType_Equip_End          int32 = 6

	// 晶石子类型
	Item_SubType_Crystal_Begin   int32 = 0
	Item_SubType_Crystal_Earth   int32 = 0 // 0 地
	Item_SubType_Crystal_Water   int32 = 1 // 1 水
	Item_SubType_Crystal_Fire    int32 = 2 // 2 火
	Item_SubType_Crystal_Wind    int32 = 3 // 3 风
	Item_SubType_Crystal_Thunder int32 = 4 // 4 雷
	Item_SubType_Crystal_Time    int32 = 5 // 5 时
	Item_SubType_Crystal_Space   int32 = 6 // 6 空
	Item_SubType_Crystal_Steel   int32 = 7 // 7 钢
	Item_SubType_Crystal_Destroy int32 = 8 // 8 灭
	Item_SubType_Crystal_End     int32 = 9
)

const (
	Item_Effect_Null int32 = -1   // 无效果
	Item_Effect_Loot int32 = iota // 掉落

	Item_Effect_End
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

const (
	Equip_Max_Promote_Times = 5 // 装备最多突破5次
	Equip_Max_Starup_Times  = 6 // 装备升星最多6次

	Crystal_Subtype_Num = Item_SubType_Crystal_End - Item_SubType_Crystal_Begin // 晶石子类型个数
)
