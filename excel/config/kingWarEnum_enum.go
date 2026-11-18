package config

// 皇城争锋-拳皇/KingWarEnum
type EMKingWarMatchType int

const (
	EMKingWarMatchType_BEGIN EMKingWarMatchType = 1 //枚举开始
	EMKingWarMatchType_Week1 EMKingWarMatchType = 1 //周赛1
	EMKingWarMatchType_Week2 EMKingWarMatchType = 2 //周赛2
	EMKingWarMatchType_Week3 EMKingWarMatchType = 3 //周赛3
	EMKingWarMatchType_Final EMKingWarMatchType = 4 //决赛
	EMKingWarMatchType_END   EMKingWarMatchType = 5 //枚举结束
)

// 皇城争锋-拳皇/KingWarEnum
type KingWarGroupType int

const (
	KingWarGroupType_BEGIN   KingWarGroupType = 1 //枚举开始
	KingWarGroupType_King    KingWarGroupType = 1 //王者
	KingWarGroupType_Diamond KingWarGroupType = 2 //钻石
	KingWarGroupType_Bench   KingWarGroupType = 3 //预备
	KingWarGroupType_END     KingWarGroupType = 4 //枚举结束
)

// 皇城争锋-拳皇/KingWarEnum
type KingWarProgress int

const (
	KingWarProgress_BEGIN           KingWarProgress = 1  //枚举开始
	KingWarProgress_Week1Sign       KingWarProgress = 1  //周赛1报名
	KingWarProgress_Week1Notice     KingWarProgress = 2  //周赛1公示
	KingWarProgress_Week1Battle     KingWarProgress = 3  //周赛1战斗
	KingWarProgress_Week1Exhibution KingWarProgress = 4  //周赛1展示
	KingWarProgress_Week1Empty      KingWarProgress = 5  //周赛1空窗
	KingWarProgress_Week2Sign       KingWarProgress = 6  //周赛2报名
	KingWarProgress_Week2Notice     KingWarProgress = 7  //周赛2公示
	KingWarProgress_Week2Battle     KingWarProgress = 8  //周赛2战斗
	KingWarProgress_Week2Exhibution KingWarProgress = 9  //周赛2展示
	KingWarProgress_Week2Empty      KingWarProgress = 10 //周赛2空窗
	KingWarProgress_Week3Sign       KingWarProgress = 11 //周赛3报名
	KingWarProgress_Week3Notice     KingWarProgress = 12 //周赛3公示
	KingWarProgress_Week3Battle     KingWarProgress = 13 //周赛3战斗
	KingWarProgress_Week3Exhibution KingWarProgress = 14 //周赛3展示
	KingWarProgress_Week3Empty      KingWarProgress = 15 //周赛3空窗
	KingWarProgress_FinalNotice     KingWarProgress = 16 //决赛公示
	KingWarProgress_MercenaryStart  KingWarProgress = 17 //佣兵开始
	KingWarProgress_MercenaryEnd    KingWarProgress = 18 //佣兵结束
	KingWarProgress_FinalBattle     KingWarProgress = 19 //决赛战斗
	KingWarProgress_FinalExhibution KingWarProgress = 20 //决赛展示
	KingWarProgress_FinalEmpty      KingWarProgress = 21 //决赛空窗
	KingWarProgress_Close           KingWarProgress = 22 //活动结束
	KingWarProgress_END             KingWarProgress = 23 //枚举结束
)

// 皇城争锋-拳皇/KingWarEnum
type KingWarBattleState int

const (
	KingWarBattleState_BEGIN    KingWarBattleState = 1 //枚举开始
	KingWarBattleState_Safe     KingWarBattleState = 1 //安全期
	KingWarBattleState_CoreCity KingWarBattleState = 2 //核心城
	KingWarBattleState_Free     KingWarBattleState = 3 //混战
	KingWarBattleState_END      KingWarBattleState = 4 //枚举结束
)

// 皇城争锋-拳皇/KingWarEnum
type EMKingWarMapGridType int

const (
	EMKingWarMapGridType_BEGIN       EMKingWarMapGridType = 0  //枚举开始
	EMKingWarMapGridType_Empty       EMKingWarMapGridType = 0  //空地
	EMKingWarMapGridType_Base1       EMKingWarMapGridType = 1  //大本营1
	EMKingWarMapGridType_Base2       EMKingWarMapGridType = 2  //大本营2
	EMKingWarMapGridType_Base3       EMKingWarMapGridType = 3  //大本营3
	EMKingWarMapGridType_Base4       EMKingWarMapGridType = 4  //大本营4
	EMKingWarMapGridType_Base5       EMKingWarMapGridType = 5  //大本营5
	EMKingWarMapGridType_Base6       EMKingWarMapGridType = 6  //大本营6
	EMKingWarMapGridType_Battery     EMKingWarMapGridType = 7  //炮台
	EMKingWarMapGridType_Training    EMKingWarMapGridType = 8  //训练场
	EMKingWarMapGridType_Fortress    EMKingWarMapGridType = 9  //城寨
	EMKingWarMapGridType_Hall        EMKingWarMapGridType = 10 //王殿
	EMKingWarMapGridType_PalaceRound EMKingWarMapGridType = 11 //皇居外围
	EMKingWarMapGridType_Palace      EMKingWarMapGridType = 12 //皇居
	EMKingWarMapGridType_Grass       EMKingWarMapGridType = 13 //草地
	EMKingWarMapGridType_END         EMKingWarMapGridType = 14 //枚举结束
)

// 皇城争锋-拳皇/KingWarEnum
type EMKingWarCampType int

const (
	EMKingWarCampType_BEGIN EMKingWarCampType = 0 //枚举开始
	EMKingWarCampType_Empty EMKingWarCampType = 0 //空
	EMKingWarCampType_A     EMKingWarCampType = 1 //阵营1
	EMKingWarCampType_B     EMKingWarCampType = 2 //阵营2
	EMKingWarCampType_C     EMKingWarCampType = 3 //阵营3
	EMKingWarCampType_D     EMKingWarCampType = 4 //阵营4
	EMKingWarCampType_E     EMKingWarCampType = 5 //阵营5
	EMKingWarCampType_F     EMKingWarCampType = 6 //阵营6
	EMKingWarCampType_END   EMKingWarCampType = 7 //枚举结束
)

// 皇城争锋-拳皇/KingWarEnum
type EMKingWarMapItemType int

const (
	EMKingWarMapItemType_BEGIN         EMKingWarMapItemType = 1  //枚举开始
	EMKingWarMapItemType_ItemSize      EMKingWarMapItemType = 1  //地块大小
	EMKingWarMapItemType_Background    EMKingWarMapItemType = 2  //地块显示
	EMKingWarMapItemType_Path          EMKingWarMapItemType = 3  //路径显示
	EMKingWarMapItemType_NormalProtect EMKingWarMapItemType = 4  //普通保护动画
	EMKingWarMapItemType_CoreProtect   EMKingWarMapItemType = 5  //核心城保护动画
	EMKingWarMapItemType_NormalBattle  EMKingWarMapItemType = 6  //普通战斗动画
	EMKingWarMapItemType_CoreBattle    EMKingWarMapItemType = 7  //核心城战斗动画
	EMKingWarMapItemType_NormalOccupy  EMKingWarMapItemType = 8  //普通占领动画
	EMKingWarMapItemType_CoreOccupy    EMKingWarMapItemType = 9  //核心城占领动画
	EMKingWarMapItemType_NormalSelect  EMKingWarMapItemType = 10 //普通选中
	EMKingWarMapItemType_CoreSelect    EMKingWarMapItemType = 11 //核心选中
	EMKingWarMapItemType_CorePath      EMKingWarMapItemType = 12 //核心城路径显示
	EMKingWarMapItemType_END           EMKingWarMapItemType = 13 //枚举结束
)

func (e EMKingWarMatchType) String() string {

	if e == EMKingWarMatchType_Week1 {
		return "Week1"
	}
	if e == EMKingWarMatchType_Week2 {
		return "Week2"
	}
	if e == EMKingWarMatchType_Week3 {
		return "Week3"
	}
	if e == EMKingWarMatchType_Final {
		return "Final"
	}

	return ""

}

func GetEMKingWarMatchType(tp string) (EMKingWarMatchType, bool) {
	switch tp {
	case "Week1":
		return EMKingWarMatchType_Week1, true
	case "Week2":
		return EMKingWarMatchType_Week2, true
	case "Week3":
		return EMKingWarMatchType_Week3, true
	case "Final":
		return EMKingWarMatchType_Final, true

	}
	return -1, false

}

func GetEMKingWarMatchType_ByAlias(tp string) (EMKingWarMatchType, bool) {
	switch tp {
	case "周赛1":
		return EMKingWarMatchType_Week1, true
	case "周赛2":
		return EMKingWarMatchType_Week2, true
	case "周赛3":
		return EMKingWarMatchType_Week3, true
	case "决赛":
		return EMKingWarMatchType_Final, true

	}
	return -1, false

}

func (e KingWarGroupType) String() string {

	if e == KingWarGroupType_King {
		return "King"
	}
	if e == KingWarGroupType_Diamond {
		return "Diamond"
	}
	if e == KingWarGroupType_Bench {
		return "Bench"
	}

	return ""

}

func GetKingWarGroupType(tp string) (KingWarGroupType, bool) {
	switch tp {
	case "King":
		return KingWarGroupType_King, true
	case "Diamond":
		return KingWarGroupType_Diamond, true
	case "Bench":
		return KingWarGroupType_Bench, true

	}
	return -1, false

}

func GetKingWarGroupType_ByAlias(tp string) (KingWarGroupType, bool) {
	switch tp {
	case "王者":
		return KingWarGroupType_King, true
	case "钻石":
		return KingWarGroupType_Diamond, true
	case "预备":
		return KingWarGroupType_Bench, true

	}
	return -1, false

}

func (e KingWarProgress) String() string {

	if e == KingWarProgress_Week1Sign {
		return "Week1Sign"
	}
	if e == KingWarProgress_Week1Notice {
		return "Week1Notice"
	}
	if e == KingWarProgress_Week1Battle {
		return "Week1Battle"
	}
	if e == KingWarProgress_Week1Exhibution {
		return "Week1Exhibution"
	}
	if e == KingWarProgress_Week1Empty {
		return "Week1Empty"
	}
	if e == KingWarProgress_Week2Sign {
		return "Week2Sign"
	}
	if e == KingWarProgress_Week2Notice {
		return "Week2Notice"
	}
	if e == KingWarProgress_Week2Battle {
		return "Week2Battle"
	}
	if e == KingWarProgress_Week2Exhibution {
		return "Week2Exhibution"
	}
	if e == KingWarProgress_Week2Empty {
		return "Week2Empty"
	}
	if e == KingWarProgress_Week3Sign {
		return "Week3Sign"
	}
	if e == KingWarProgress_Week3Notice {
		return "Week3Notice"
	}
	if e == KingWarProgress_Week3Battle {
		return "Week3Battle"
	}
	if e == KingWarProgress_Week3Exhibution {
		return "Week3Exhibution"
	}
	if e == KingWarProgress_Week3Empty {
		return "Week3Empty"
	}
	if e == KingWarProgress_FinalNotice {
		return "FinalNotice"
	}
	if e == KingWarProgress_MercenaryStart {
		return "MercenaryStart"
	}
	if e == KingWarProgress_MercenaryEnd {
		return "MercenaryEnd"
	}
	if e == KingWarProgress_FinalBattle {
		return "FinalBattle"
	}
	if e == KingWarProgress_FinalExhibution {
		return "FinalExhibution"
	}
	if e == KingWarProgress_FinalEmpty {
		return "FinalEmpty"
	}
	if e == KingWarProgress_Close {
		return "Close"
	}

	return ""

}

func GetKingWarProgress(tp string) (KingWarProgress, bool) {
	switch tp {
	case "Week1Sign":
		return KingWarProgress_Week1Sign, true
	case "Week1Notice":
		return KingWarProgress_Week1Notice, true
	case "Week1Battle":
		return KingWarProgress_Week1Battle, true
	case "Week1Exhibution":
		return KingWarProgress_Week1Exhibution, true
	case "Week1Empty":
		return KingWarProgress_Week1Empty, true
	case "Week2Sign":
		return KingWarProgress_Week2Sign, true
	case "Week2Notice":
		return KingWarProgress_Week2Notice, true
	case "Week2Battle":
		return KingWarProgress_Week2Battle, true
	case "Week2Exhibution":
		return KingWarProgress_Week2Exhibution, true
	case "Week2Empty":
		return KingWarProgress_Week2Empty, true
	case "Week3Sign":
		return KingWarProgress_Week3Sign, true
	case "Week3Notice":
		return KingWarProgress_Week3Notice, true
	case "Week3Battle":
		return KingWarProgress_Week3Battle, true
	case "Week3Exhibution":
		return KingWarProgress_Week3Exhibution, true
	case "Week3Empty":
		return KingWarProgress_Week3Empty, true
	case "FinalNotice":
		return KingWarProgress_FinalNotice, true
	case "MercenaryStart":
		return KingWarProgress_MercenaryStart, true
	case "MercenaryEnd":
		return KingWarProgress_MercenaryEnd, true
	case "FinalBattle":
		return KingWarProgress_FinalBattle, true
	case "FinalExhibution":
		return KingWarProgress_FinalExhibution, true
	case "FinalEmpty":
		return KingWarProgress_FinalEmpty, true
	case "Close":
		return KingWarProgress_Close, true

	}
	return -1, false

}

func GetKingWarProgress_ByAlias(tp string) (KingWarProgress, bool) {
	switch tp {
	case "周赛1报名":
		return KingWarProgress_Week1Sign, true
	case "周赛1公示":
		return KingWarProgress_Week1Notice, true
	case "周赛1战斗":
		return KingWarProgress_Week1Battle, true
	case "周赛1展示":
		return KingWarProgress_Week1Exhibution, true
	case "周赛1空窗":
		return KingWarProgress_Week1Empty, true
	case "周赛2报名":
		return KingWarProgress_Week2Sign, true
	case "周赛2公示":
		return KingWarProgress_Week2Notice, true
	case "周赛2战斗":
		return KingWarProgress_Week2Battle, true
	case "周赛2展示":
		return KingWarProgress_Week2Exhibution, true
	case "周赛2空窗":
		return KingWarProgress_Week2Empty, true
	case "周赛3报名":
		return KingWarProgress_Week3Sign, true
	case "周赛3公示":
		return KingWarProgress_Week3Notice, true
	case "周赛3战斗":
		return KingWarProgress_Week3Battle, true
	case "周赛3展示":
		return KingWarProgress_Week3Exhibution, true
	case "周赛3空窗":
		return KingWarProgress_Week3Empty, true
	case "决赛公示":
		return KingWarProgress_FinalNotice, true
	case "佣兵开始":
		return KingWarProgress_MercenaryStart, true
	case "佣兵结束":
		return KingWarProgress_MercenaryEnd, true
	case "决赛战斗":
		return KingWarProgress_FinalBattle, true
	case "决赛展示":
		return KingWarProgress_FinalExhibution, true
	case "决赛空窗":
		return KingWarProgress_FinalEmpty, true
	case "活动结束":
		return KingWarProgress_Close, true

	}
	return -1, false

}

func (e KingWarBattleState) String() string {

	if e == KingWarBattleState_Safe {
		return "Safe"
	}
	if e == KingWarBattleState_CoreCity {
		return "CoreCity"
	}
	if e == KingWarBattleState_Free {
		return "Free"
	}

	return ""

}

func GetKingWarBattleState(tp string) (KingWarBattleState, bool) {
	switch tp {
	case "Safe":
		return KingWarBattleState_Safe, true
	case "CoreCity":
		return KingWarBattleState_CoreCity, true
	case "Free":
		return KingWarBattleState_Free, true

	}
	return -1, false

}

func GetKingWarBattleState_ByAlias(tp string) (KingWarBattleState, bool) {
	switch tp {
	case "安全期":
		return KingWarBattleState_Safe, true
	case "核心城":
		return KingWarBattleState_CoreCity, true
	case "混战":
		return KingWarBattleState_Free, true

	}
	return -1, false

}

func (e EMKingWarMapGridType) String() string {

	if e == EMKingWarMapGridType_Empty {
		return "Empty"
	}
	if e == EMKingWarMapGridType_Base1 {
		return "Base1"
	}
	if e == EMKingWarMapGridType_Base2 {
		return "Base2"
	}
	if e == EMKingWarMapGridType_Base3 {
		return "Base3"
	}
	if e == EMKingWarMapGridType_Base4 {
		return "Base4"
	}
	if e == EMKingWarMapGridType_Base5 {
		return "Base5"
	}
	if e == EMKingWarMapGridType_Base6 {
		return "Base6"
	}
	if e == EMKingWarMapGridType_Battery {
		return "Battery"
	}
	if e == EMKingWarMapGridType_Training {
		return "Training"
	}
	if e == EMKingWarMapGridType_Fortress {
		return "Fortress"
	}
	if e == EMKingWarMapGridType_Hall {
		return "Hall"
	}
	if e == EMKingWarMapGridType_PalaceRound {
		return "PalaceRound"
	}
	if e == EMKingWarMapGridType_Palace {
		return "Palace"
	}
	if e == EMKingWarMapGridType_Grass {
		return "Grass"
	}

	return ""

}

func GetEMKingWarMapGridType(tp string) (EMKingWarMapGridType, bool) {
	switch tp {
	case "Empty":
		return EMKingWarMapGridType_Empty, true
	case "Base1":
		return EMKingWarMapGridType_Base1, true
	case "Base2":
		return EMKingWarMapGridType_Base2, true
	case "Base3":
		return EMKingWarMapGridType_Base3, true
	case "Base4":
		return EMKingWarMapGridType_Base4, true
	case "Base5":
		return EMKingWarMapGridType_Base5, true
	case "Base6":
		return EMKingWarMapGridType_Base6, true
	case "Battery":
		return EMKingWarMapGridType_Battery, true
	case "Training":
		return EMKingWarMapGridType_Training, true
	case "Fortress":
		return EMKingWarMapGridType_Fortress, true
	case "Hall":
		return EMKingWarMapGridType_Hall, true
	case "PalaceRound":
		return EMKingWarMapGridType_PalaceRound, true
	case "Palace":
		return EMKingWarMapGridType_Palace, true
	case "Grass":
		return EMKingWarMapGridType_Grass, true

	}
	return -1, false

}

func GetEMKingWarMapGridType_ByAlias(tp string) (EMKingWarMapGridType, bool) {
	switch tp {
	case "空地":
		return EMKingWarMapGridType_Empty, true
	case "大本营1":
		return EMKingWarMapGridType_Base1, true
	case "大本营2":
		return EMKingWarMapGridType_Base2, true
	case "大本营3":
		return EMKingWarMapGridType_Base3, true
	case "大本营4":
		return EMKingWarMapGridType_Base4, true
	case "大本营5":
		return EMKingWarMapGridType_Base5, true
	case "大本营6":
		return EMKingWarMapGridType_Base6, true
	case "炮台":
		return EMKingWarMapGridType_Battery, true
	case "训练场":
		return EMKingWarMapGridType_Training, true
	case "城寨":
		return EMKingWarMapGridType_Fortress, true
	case "王殿":
		return EMKingWarMapGridType_Hall, true
	case "皇居外围":
		return EMKingWarMapGridType_PalaceRound, true
	case "皇居":
		return EMKingWarMapGridType_Palace, true
	case "草地":
		return EMKingWarMapGridType_Grass, true

	}
	return -1, false

}

func (e EMKingWarCampType) String() string {

	if e == EMKingWarCampType_Empty {
		return "Empty"
	}
	if e == EMKingWarCampType_A {
		return "A"
	}
	if e == EMKingWarCampType_B {
		return "B"
	}
	if e == EMKingWarCampType_C {
		return "C"
	}
	if e == EMKingWarCampType_D {
		return "D"
	}
	if e == EMKingWarCampType_E {
		return "E"
	}
	if e == EMKingWarCampType_F {
		return "F"
	}

	return ""

}

func GetEMKingWarCampType(tp string) (EMKingWarCampType, bool) {
	switch tp {
	case "Empty":
		return EMKingWarCampType_Empty, true
	case "A":
		return EMKingWarCampType_A, true
	case "B":
		return EMKingWarCampType_B, true
	case "C":
		return EMKingWarCampType_C, true
	case "D":
		return EMKingWarCampType_D, true
	case "E":
		return EMKingWarCampType_E, true
	case "F":
		return EMKingWarCampType_F, true

	}
	return -1, false

}

func GetEMKingWarCampType_ByAlias(tp string) (EMKingWarCampType, bool) {
	switch tp {
	case "空":
		return EMKingWarCampType_Empty, true
	case "阵营1":
		return EMKingWarCampType_A, true
	case "阵营2":
		return EMKingWarCampType_B, true
	case "阵营3":
		return EMKingWarCampType_C, true
	case "阵营4":
		return EMKingWarCampType_D, true
	case "阵营5":
		return EMKingWarCampType_E, true
	case "阵营6":
		return EMKingWarCampType_F, true

	}
	return -1, false

}

func (e EMKingWarMapItemType) String() string {

	if e == EMKingWarMapItemType_ItemSize {
		return "ItemSize"
	}
	if e == EMKingWarMapItemType_Background {
		return "Background"
	}
	if e == EMKingWarMapItemType_Path {
		return "Path"
	}
	if e == EMKingWarMapItemType_NormalProtect {
		return "NormalProtect"
	}
	if e == EMKingWarMapItemType_CoreProtect {
		return "CoreProtect"
	}
	if e == EMKingWarMapItemType_NormalBattle {
		return "NormalBattle"
	}
	if e == EMKingWarMapItemType_CoreBattle {
		return "CoreBattle"
	}
	if e == EMKingWarMapItemType_NormalOccupy {
		return "NormalOccupy"
	}
	if e == EMKingWarMapItemType_CoreOccupy {
		return "CoreOccupy"
	}
	if e == EMKingWarMapItemType_NormalSelect {
		return "NormalSelect"
	}
	if e == EMKingWarMapItemType_CoreSelect {
		return "CoreSelect"
	}
	if e == EMKingWarMapItemType_CorePath {
		return "CorePath"
	}

	return ""

}

func GetEMKingWarMapItemType(tp string) (EMKingWarMapItemType, bool) {
	switch tp {
	case "ItemSize":
		return EMKingWarMapItemType_ItemSize, true
	case "Background":
		return EMKingWarMapItemType_Background, true
	case "Path":
		return EMKingWarMapItemType_Path, true
	case "NormalProtect":
		return EMKingWarMapItemType_NormalProtect, true
	case "CoreProtect":
		return EMKingWarMapItemType_CoreProtect, true
	case "NormalBattle":
		return EMKingWarMapItemType_NormalBattle, true
	case "CoreBattle":
		return EMKingWarMapItemType_CoreBattle, true
	case "NormalOccupy":
		return EMKingWarMapItemType_NormalOccupy, true
	case "CoreOccupy":
		return EMKingWarMapItemType_CoreOccupy, true
	case "NormalSelect":
		return EMKingWarMapItemType_NormalSelect, true
	case "CoreSelect":
		return EMKingWarMapItemType_CoreSelect, true
	case "CorePath":
		return EMKingWarMapItemType_CorePath, true

	}
	return -1, false

}

func GetEMKingWarMapItemType_ByAlias(tp string) (EMKingWarMapItemType, bool) {
	switch tp {
	case "地块大小":
		return EMKingWarMapItemType_ItemSize, true
	case "地块显示":
		return EMKingWarMapItemType_Background, true
	case "路径显示":
		return EMKingWarMapItemType_Path, true
	case "普通保护动画":
		return EMKingWarMapItemType_NormalProtect, true
	case "核心城保护动画":
		return EMKingWarMapItemType_CoreProtect, true
	case "普通战斗动画":
		return EMKingWarMapItemType_NormalBattle, true
	case "核心城战斗动画":
		return EMKingWarMapItemType_CoreBattle, true
	case "普通占领动画":
		return EMKingWarMapItemType_NormalOccupy, true
	case "核心城占领动画":
		return EMKingWarMapItemType_CoreOccupy, true
	case "普通选中":
		return EMKingWarMapItemType_NormalSelect, true
	case "核心选中":
		return EMKingWarMapItemType_CoreSelect, true
	case "核心城路径显示":
		return EMKingWarMapItemType_CorePath, true

	}
	return -1, false

}
