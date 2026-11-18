package config

// 阵营战-拳皇/CampBattleMapEnum
type EMCampBattleMapGridType int

const (
	EMCampBattleMapGridType_BEGIN          EMCampBattleMapGridType = 0  //枚举开始
	EMCampBattleMapGridType_Empty          EMCampBattleMapGridType = 0  //空地
	EMCampBattleMapGridType_Block          EMCampBattleMapGridType = 1  //障碍
	EMCampBattleMapGridType_Hollow         EMCampBattleMapGridType = 2  //镂空
	EMCampBattleMapGridType_RadarEagle     EMCampBattleMapGridType = 3  //雷达鹰
	EMCampBattleMapGridType_RadarWolf      EMCampBattleMapGridType = 4  //雷达狼
	EMCampBattleMapGridType_Center         EMCampBattleMapGridType = 5  //中心
	EMCampBattleMapGridType_BaseEagle      EMCampBattleMapGridType = 6  //大本营鹰
	EMCampBattleMapGridType_BaseWolf       EMCampBattleMapGridType = 7  //大本营狼
	EMCampBattleMapGridType_Building1Eagle EMCampBattleMapGridType = 8  //建筑物1级鹰
	EMCampBattleMapGridType_Building2Eagle EMCampBattleMapGridType = 9  //建筑物2级鹰
	EMCampBattleMapGridType_Building3Eagle EMCampBattleMapGridType = 10 //建筑物3级鹰
	EMCampBattleMapGridType_Building4Eagle EMCampBattleMapGridType = 11 //建筑物4级鹰
	EMCampBattleMapGridType_Building1Wolf  EMCampBattleMapGridType = 12 //建筑物1级狼
	EMCampBattleMapGridType_Building2Wolf  EMCampBattleMapGridType = 13 //建筑物2级狼
	EMCampBattleMapGridType_Building3Wolf  EMCampBattleMapGridType = 14 //建筑物3级狼
	EMCampBattleMapGridType_Building4Wolf  EMCampBattleMapGridType = 15 //建筑物4级狼
	EMCampBattleMapGridType_Tree           EMCampBattleMapGridType = 16 //树
	EMCampBattleMapGridType_Stone          EMCampBattleMapGridType = 17 //石头
	EMCampBattleMapGridType_Ruins          EMCampBattleMapGridType = 18 //废墟
	EMCampBattleMapGridType_END            EMCampBattleMapGridType = 19 //枚举结束
)

// 阵营战-拳皇/CampBattleMapEnum
type EMCampBattleBuildingType int

const (
	EMCampBattleBuildingType_BEGIN  EMCampBattleBuildingType = 0 //枚举开始
	EMCampBattleBuildingType_Empty  EMCampBattleBuildingType = 0 //空地
	EMCampBattleBuildingType_Block  EMCampBattleBuildingType = 1 //障碍
	EMCampBattleBuildingType_Hollow EMCampBattleBuildingType = 2 //镂空
	EMCampBattleBuildingType_Radar  EMCampBattleBuildingType = 3 //雷达
	EMCampBattleBuildingType_Center EMCampBattleBuildingType = 4 //中心
	EMCampBattleBuildingType_Base   EMCampBattleBuildingType = 5 //大本营
	EMCampBattleBuildingType_Normal EMCampBattleBuildingType = 6 //普通建筑
	EMCampBattleBuildingType_END    EMCampBattleBuildingType = 7 //枚举结束
)

// 阵营战-拳皇/CampBattleMapEnum
type EMCampBattleType int

const (
	EMCampBattleType_BEGIN EMCampBattleType = 0 //枚举开始
	EMCampBattleType_Empty EMCampBattleType = 0 //无
	EMCampBattleType_Eagle EMCampBattleType = 1 //鹰
	EMCampBattleType_Wolf  EMCampBattleType = 2 //狼
	EMCampBattleType_END   EMCampBattleType = 3 //枚举结束
)

// 阵营战-拳皇/CampBattleMapEnum
type EMCampBattleSeizeType int

const (
	EMCampBattleSeizeType_BEGIN   EMCampBattleSeizeType = 0 //枚举开始
	EMCampBattleSeizeType_None    EMCampBattleSeizeType = 0 //无法占领
	EMCampBattleSeizeType_Achieve EMCampBattleSeizeType = 1 //到达时长
	EMCampBattleSeizeType_History EMCampBattleSeizeType = 2 //历史时长
	EMCampBattleSeizeType_END     EMCampBattleSeizeType = 3 //枚举结束
)

func (e EMCampBattleMapGridType) String() string {

	if e == EMCampBattleMapGridType_Empty {
		return "Empty"
	}
	if e == EMCampBattleMapGridType_Block {
		return "Block"
	}
	if e == EMCampBattleMapGridType_Hollow {
		return "Hollow"
	}
	if e == EMCampBattleMapGridType_RadarEagle {
		return "RadarEagle"
	}
	if e == EMCampBattleMapGridType_RadarWolf {
		return "RadarWolf"
	}
	if e == EMCampBattleMapGridType_Center {
		return "Center"
	}
	if e == EMCampBattleMapGridType_BaseEagle {
		return "BaseEagle"
	}
	if e == EMCampBattleMapGridType_BaseWolf {
		return "BaseWolf"
	}
	if e == EMCampBattleMapGridType_Building1Eagle {
		return "Building1Eagle"
	}
	if e == EMCampBattleMapGridType_Building2Eagle {
		return "Building2Eagle"
	}
	if e == EMCampBattleMapGridType_Building3Eagle {
		return "Building3Eagle"
	}
	if e == EMCampBattleMapGridType_Building4Eagle {
		return "Building4Eagle"
	}
	if e == EMCampBattleMapGridType_Building1Wolf {
		return "Building1Wolf"
	}
	if e == EMCampBattleMapGridType_Building2Wolf {
		return "Building2Wolf"
	}
	if e == EMCampBattleMapGridType_Building3Wolf {
		return "Building3Wolf"
	}
	if e == EMCampBattleMapGridType_Building4Wolf {
		return "Building4Wolf"
	}
	if e == EMCampBattleMapGridType_Tree {
		return "Tree"
	}
	if e == EMCampBattleMapGridType_Stone {
		return "Stone"
	}
	if e == EMCampBattleMapGridType_Ruins {
		return "Ruins"
	}

	return ""

}

func GetEMCampBattleMapGridType(tp string) (EMCampBattleMapGridType, bool) {
	switch tp {
	case "Empty":
		return EMCampBattleMapGridType_Empty, true
	case "Block":
		return EMCampBattleMapGridType_Block, true
	case "Hollow":
		return EMCampBattleMapGridType_Hollow, true
	case "RadarEagle":
		return EMCampBattleMapGridType_RadarEagle, true
	case "RadarWolf":
		return EMCampBattleMapGridType_RadarWolf, true
	case "Center":
		return EMCampBattleMapGridType_Center, true
	case "BaseEagle":
		return EMCampBattleMapGridType_BaseEagle, true
	case "BaseWolf":
		return EMCampBattleMapGridType_BaseWolf, true
	case "Building1Eagle":
		return EMCampBattleMapGridType_Building1Eagle, true
	case "Building2Eagle":
		return EMCampBattleMapGridType_Building2Eagle, true
	case "Building3Eagle":
		return EMCampBattleMapGridType_Building3Eagle, true
	case "Building4Eagle":
		return EMCampBattleMapGridType_Building4Eagle, true
	case "Building1Wolf":
		return EMCampBattleMapGridType_Building1Wolf, true
	case "Building2Wolf":
		return EMCampBattleMapGridType_Building2Wolf, true
	case "Building3Wolf":
		return EMCampBattleMapGridType_Building3Wolf, true
	case "Building4Wolf":
		return EMCampBattleMapGridType_Building4Wolf, true
	case "Tree":
		return EMCampBattleMapGridType_Tree, true
	case "Stone":
		return EMCampBattleMapGridType_Stone, true
	case "Ruins":
		return EMCampBattleMapGridType_Ruins, true

	}
	return -1, false

}

func GetEMCampBattleMapGridType_ByAlias(tp string) (EMCampBattleMapGridType, bool) {
	switch tp {
	case "空地":
		return EMCampBattleMapGridType_Empty, true
	case "障碍":
		return EMCampBattleMapGridType_Block, true
	case "镂空":
		return EMCampBattleMapGridType_Hollow, true
	case "雷达鹰":
		return EMCampBattleMapGridType_RadarEagle, true
	case "雷达狼":
		return EMCampBattleMapGridType_RadarWolf, true
	case "中心":
		return EMCampBattleMapGridType_Center, true
	case "大本营鹰":
		return EMCampBattleMapGridType_BaseEagle, true
	case "大本营狼":
		return EMCampBattleMapGridType_BaseWolf, true
	case "建筑物1级鹰":
		return EMCampBattleMapGridType_Building1Eagle, true
	case "建筑物2级鹰":
		return EMCampBattleMapGridType_Building2Eagle, true
	case "建筑物3级鹰":
		return EMCampBattleMapGridType_Building3Eagle, true
	case "建筑物4级鹰":
		return EMCampBattleMapGridType_Building4Eagle, true
	case "建筑物1级狼":
		return EMCampBattleMapGridType_Building1Wolf, true
	case "建筑物2级狼":
		return EMCampBattleMapGridType_Building2Wolf, true
	case "建筑物3级狼":
		return EMCampBattleMapGridType_Building3Wolf, true
	case "建筑物4级狼":
		return EMCampBattleMapGridType_Building4Wolf, true
	case "树":
		return EMCampBattleMapGridType_Tree, true
	case "石头":
		return EMCampBattleMapGridType_Stone, true
	case "废墟":
		return EMCampBattleMapGridType_Ruins, true

	}
	return -1, false

}

func (e EMCampBattleBuildingType) String() string {

	if e == EMCampBattleBuildingType_Empty {
		return "Empty"
	}
	if e == EMCampBattleBuildingType_Block {
		return "Block"
	}
	if e == EMCampBattleBuildingType_Hollow {
		return "Hollow"
	}
	if e == EMCampBattleBuildingType_Radar {
		return "Radar"
	}
	if e == EMCampBattleBuildingType_Center {
		return "Center"
	}
	if e == EMCampBattleBuildingType_Base {
		return "Base"
	}
	if e == EMCampBattleBuildingType_Normal {
		return "Normal"
	}

	return ""

}

func GetEMCampBattleBuildingType(tp string) (EMCampBattleBuildingType, bool) {
	switch tp {
	case "Empty":
		return EMCampBattleBuildingType_Empty, true
	case "Block":
		return EMCampBattleBuildingType_Block, true
	case "Hollow":
		return EMCampBattleBuildingType_Hollow, true
	case "Radar":
		return EMCampBattleBuildingType_Radar, true
	case "Center":
		return EMCampBattleBuildingType_Center, true
	case "Base":
		return EMCampBattleBuildingType_Base, true
	case "Normal":
		return EMCampBattleBuildingType_Normal, true

	}
	return -1, false

}

func GetEMCampBattleBuildingType_ByAlias(tp string) (EMCampBattleBuildingType, bool) {
	switch tp {
	case "空地":
		return EMCampBattleBuildingType_Empty, true
	case "障碍":
		return EMCampBattleBuildingType_Block, true
	case "镂空":
		return EMCampBattleBuildingType_Hollow, true
	case "雷达":
		return EMCampBattleBuildingType_Radar, true
	case "中心":
		return EMCampBattleBuildingType_Center, true
	case "大本营":
		return EMCampBattleBuildingType_Base, true
	case "普通建筑":
		return EMCampBattleBuildingType_Normal, true

	}
	return -1, false

}

func (e EMCampBattleType) String() string {

	if e == EMCampBattleType_Empty {
		return "Empty"
	}
	if e == EMCampBattleType_Eagle {
		return "Eagle"
	}
	if e == EMCampBattleType_Wolf {
		return "Wolf"
	}

	return ""

}

func GetEMCampBattleType(tp string) (EMCampBattleType, bool) {
	switch tp {
	case "Empty":
		return EMCampBattleType_Empty, true
	case "Eagle":
		return EMCampBattleType_Eagle, true
	case "Wolf":
		return EMCampBattleType_Wolf, true

	}
	return -1, false

}

func GetEMCampBattleType_ByAlias(tp string) (EMCampBattleType, bool) {
	switch tp {
	case "无":
		return EMCampBattleType_Empty, true
	case "鹰":
		return EMCampBattleType_Eagle, true
	case "狼":
		return EMCampBattleType_Wolf, true

	}
	return -1, false

}

func (e EMCampBattleSeizeType) String() string {

	if e == EMCampBattleSeizeType_None {
		return "None"
	}
	if e == EMCampBattleSeizeType_Achieve {
		return "Achieve"
	}
	if e == EMCampBattleSeizeType_History {
		return "History"
	}

	return ""

}

func GetEMCampBattleSeizeType(tp string) (EMCampBattleSeizeType, bool) {
	switch tp {
	case "None":
		return EMCampBattleSeizeType_None, true
	case "Achieve":
		return EMCampBattleSeizeType_Achieve, true
	case "History":
		return EMCampBattleSeizeType_History, true

	}
	return -1, false

}

func GetEMCampBattleSeizeType_ByAlias(tp string) (EMCampBattleSeizeType, bool) {
	switch tp {
	case "无法占领":
		return EMCampBattleSeizeType_None, true
	case "到达时长":
		return EMCampBattleSeizeType_Achieve, true
	case "历史时长":
		return EMCampBattleSeizeType_History, true

	}
	return -1, false

}
