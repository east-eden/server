package config

// LD关卡表-拳皇/LDLevelType
type LDLevelType int

const (
	LDLevelType_BEGIN               LDLevelType = 1   //枚举开始
	LDLevelType_Guide               LDLevelType = 1   //引导
	LDLevelType_Normal              LDLevelType = 2   //普通
	LDLevelType_Mutation            LDLevelType = 3   //异变
	LDLevelType_PVP                 LDLevelType = 4   //pvp
	LDLevelType_PvpFriend           LDLevelType = 5   //友谊赛
	LDLevelType_Athena              LDLevelType = 6   //雅呆娜
	LDLevelType_ContributeChallenge LDLevelType = 7   //贡献挑战
	LDLevelType_HuntZombie          LDLevelType = 8   //狩猎僵尸
	LDLevelType_Test                LDLevelType = 99  //测试
	LDLevelType_DPS                 LDLevelType = 100 //秒伤测试
	LDLevelType_TD                  LDLevelType = 101 //塔防
	LDLevelType_WorkShop            LDLevelType = 102 //创意工坊
	LDLevelType_WorkShop_Publish    LDLevelType = 103 //创意工坊尝试发布前试玩
	LDLevelType_END                 LDLevelType = 104 //枚举结束
)

// LD关卡表-拳皇/LDLevelType
type LDLevelConfigsId int

const (
	LDLevelConfigsId_BEGIN               LDLevelConfigsId = 0   //枚举开始
	LDLevelConfigsId_Teach               LDLevelConfigsId = 0   //新手教学关
	LDLevelConfigsId_Normal              LDLevelConfigsId = 1   //正式1关
	LDLevelConfigsId_MutationDifficulty1 LDLevelConfigsId = 2   //异变难度1
	LDLevelConfigsId_MutationPhysicalUp  LDLevelConfigsId = 3   //异变物抗
	LDLevelConfigsId_MutationMagicUp     LDLevelConfigsId = 4   //异变魔抗
	LDLevelConfigsId_MutationFaint       LDLevelConfigsId = 5   //异变免疫眩晕
	LDLevelConfigsId_MutationSlowSpeed   LDLevelConfigsId = 6   //异变免疫减速
	LDLevelConfigsId_MutationHardest     LDLevelConfigsId = 7   //超高难度
	LDLevelConfigsId_MutationLevel1      LDLevelConfigsId = 8   //异变1关
	LDLevelConfigsId_MutationBlodUp      LDLevelConfigsId = 9   //异变血量提升
	LDLevelConfigsId_MutationRandomTask  LDLevelConfigsId = 10  //异变随机任务
	LDLevelConfigsId_WorkShopTemp        LDLevelConfigsId = 11  //创意工坊
	LDLevelConfigsId_Pvp                 LDLevelConfigsId = 100 //pvp
	LDLevelConfigsId_PvpFriendship       LDLevelConfigsId = 101 //友谊赛
	LDLevelConfigsId_Athena              LDLevelConfigsId = 200 //雅呆娜
	LDLevelConfigsId_HuntZombie          LDLevelConfigsId = 210 //狩猎僵尸
	LDLevelConfigsId_END                 LDLevelConfigsId = 211 //枚举结束
)

// LD关卡表-拳皇/LDLevelType
type LDLevelGeneratorType int

const (
	LDLevelGeneratorType_BEGIN LDLevelGeneratorType = 1 //枚举开始
	LDLevelGeneratorType_Main  LDLevelGeneratorType = 1 //主线
	LDLevelGeneratorType_Boss  LDLevelGeneratorType = 2 //boss
	LDLevelGeneratorType_Hunt  LDLevelGeneratorType = 3 //狩猎
	LDLevelGeneratorType_END   LDLevelGeneratorType = 4 //枚举结束
)

// LD关卡表-拳皇/LDLevelType
type LDRoleBattleType int

const (
	LDRoleBattleType_BEGIN LDRoleBattleType = 1 //枚举开始
	LDRoleBattleType_Front LDRoleBattleType = 1 //前排
	LDRoleBattleType_Back  LDRoleBattleType = 2 //后排
	LDRoleBattleType_END   LDRoleBattleType = 3 //枚举结束
)

// LD关卡表-拳皇/LDLevelType
type LDBossWarningBgType int

const (
	LDBossWarningBgType_BEGIN  LDBossWarningBgType = 0 //枚举开始
	LDBossWarningBgType_Red    LDBossWarningBgType = 0 //红
	LDBossWarningBgType_Yellow LDBossWarningBgType = 1 //黄
	LDBossWarningBgType_Blue   LDBossWarningBgType = 2 //蓝
	LDBossWarningBgType_END    LDBossWarningBgType = 3 //枚举结束
)

// LD关卡表-拳皇/LDLevelType
type LDMonsterType int

const (
	LDMonsterType_BEGIN LDMonsterType = 1 //枚举开始
	LDMonsterType_One   LDMonsterType = 1 //类型1
	LDMonsterType_Two   LDMonsterType = 2 //类型2
	LDMonsterType_Tree  LDMonsterType = 3 //类型3
	LDMonsterType_Four  LDMonsterType = 4 //类型4
	LDMonsterType_Five  LDMonsterType = 5 //类型5
	LDMonsterType_END   LDMonsterType = 6 //枚举结束
)

// LD关卡表-拳皇/LDLevelType
type LDLevelEnchantType int

const (
	LDLevelEnchantType_BEGIN           LDLevelEnchantType = 0 //枚举开始
	LDLevelEnchantType_GrayBlueEnchant LDLevelEnchantType = 0 //灰色、蓝色品质强化
	LDLevelEnchantType_PurpleEnchant   LDLevelEnchantType = 1 //紫色品质强化
	LDLevelEnchantType_LegendEnchant   LDLevelEnchantType = 2 //传说品质强化
	LDLevelEnchantType_SummonEnchant   LDLevelEnchantType = 3 //召唤概率强化
	LDLevelEnchantType_END             LDLevelEnchantType = 4 //枚举结束
)

// LD关卡表-拳皇/LDLevelType
type LevelGeneratorType int

const (
	LevelGeneratorType_BEGIN LevelGeneratorType = 1 //枚举开始
	LevelGeneratorType_Main  LevelGeneratorType = 1 //主线
	LevelGeneratorType_Boss  LevelGeneratorType = 2 //boss
	LevelGeneratorType_Hunt  LevelGeneratorType = 3 //狩猎
	LevelGeneratorType_END   LevelGeneratorType = 4 //枚举结束
)

func (e LDLevelType) String() string {

	if e == LDLevelType_Guide {
		return "Guide"
	}
	if e == LDLevelType_Normal {
		return "Normal"
	}
	if e == LDLevelType_Mutation {
		return "Mutation"
	}
	if e == LDLevelType_PVP {
		return "PVP"
	}
	if e == LDLevelType_PvpFriend {
		return "PvpFriend"
	}
	if e == LDLevelType_Athena {
		return "Athena"
	}
	if e == LDLevelType_ContributeChallenge {
		return "ContributeChallenge"
	}
	if e == LDLevelType_HuntZombie {
		return "HuntZombie"
	}
	if e == LDLevelType_Test {
		return "Test"
	}
	if e == LDLevelType_DPS {
		return "DPS"
	}
	if e == LDLevelType_TD {
		return "TD"
	}
	if e == LDLevelType_WorkShop {
		return "WorkShop"
	}
	if e == LDLevelType_WorkShop_Publish {
		return "WorkShop_Publish"
	}

	return ""

}

func GetLDLevelType(tp string) (LDLevelType, bool) {
	switch tp {
	case "Guide":
		return LDLevelType_Guide, true
	case "Normal":
		return LDLevelType_Normal, true
	case "Mutation":
		return LDLevelType_Mutation, true
	case "PVP":
		return LDLevelType_PVP, true
	case "PvpFriend":
		return LDLevelType_PvpFriend, true
	case "Athena":
		return LDLevelType_Athena, true
	case "ContributeChallenge":
		return LDLevelType_ContributeChallenge, true
	case "HuntZombie":
		return LDLevelType_HuntZombie, true
	case "Test":
		return LDLevelType_Test, true
	case "DPS":
		return LDLevelType_DPS, true
	case "TD":
		return LDLevelType_TD, true
	case "WorkShop":
		return LDLevelType_WorkShop, true
	case "WorkShop_Publish":
		return LDLevelType_WorkShop_Publish, true

	}
	return -1, false

}

func GetLDLevelType_ByAlias(tp string) (LDLevelType, bool) {
	switch tp {
	case "引导":
		return LDLevelType_Guide, true
	case "普通":
		return LDLevelType_Normal, true
	case "异变":
		return LDLevelType_Mutation, true
	case "pvp":
		return LDLevelType_PVP, true
	case "友谊赛":
		return LDLevelType_PvpFriend, true
	case "雅呆娜":
		return LDLevelType_Athena, true
	case "贡献挑战":
		return LDLevelType_ContributeChallenge, true
	case "狩猎僵尸":
		return LDLevelType_HuntZombie, true
	case "测试":
		return LDLevelType_Test, true
	case "秒伤测试":
		return LDLevelType_DPS, true
	case "塔防":
		return LDLevelType_TD, true
	case "创意工坊":
		return LDLevelType_WorkShop, true
	case "创意工坊尝试发布前试玩":
		return LDLevelType_WorkShop_Publish, true

	}
	return -1, false

}

func (e LDLevelConfigsId) String() string {

	if e == LDLevelConfigsId_Teach {
		return "Teach"
	}
	if e == LDLevelConfigsId_Normal {
		return "Normal"
	}
	if e == LDLevelConfigsId_MutationDifficulty1 {
		return "MutationDifficulty1"
	}
	if e == LDLevelConfigsId_MutationPhysicalUp {
		return "MutationPhysicalUp"
	}
	if e == LDLevelConfigsId_MutationMagicUp {
		return "MutationMagicUp"
	}
	if e == LDLevelConfigsId_MutationFaint {
		return "MutationFaint"
	}
	if e == LDLevelConfigsId_MutationSlowSpeed {
		return "MutationSlowSpeed"
	}
	if e == LDLevelConfigsId_MutationHardest {
		return "MutationHardest"
	}
	if e == LDLevelConfigsId_MutationLevel1 {
		return "MutationLevel1"
	}
	if e == LDLevelConfigsId_MutationBlodUp {
		return "MutationBlodUp"
	}
	if e == LDLevelConfigsId_MutationRandomTask {
		return "MutationRandomTask"
	}
	if e == LDLevelConfigsId_WorkShopTemp {
		return "WorkShopTemp"
	}
	if e == LDLevelConfigsId_Pvp {
		return "Pvp"
	}
	if e == LDLevelConfigsId_PvpFriendship {
		return "PvpFriendship"
	}
	if e == LDLevelConfigsId_Athena {
		return "Athena"
	}
	if e == LDLevelConfigsId_HuntZombie {
		return "HuntZombie"
	}

	return ""

}

func GetLDLevelConfigsId(tp string) (LDLevelConfigsId, bool) {
	switch tp {
	case "Teach":
		return LDLevelConfigsId_Teach, true
	case "Normal":
		return LDLevelConfigsId_Normal, true
	case "MutationDifficulty1":
		return LDLevelConfigsId_MutationDifficulty1, true
	case "MutationPhysicalUp":
		return LDLevelConfigsId_MutationPhysicalUp, true
	case "MutationMagicUp":
		return LDLevelConfigsId_MutationMagicUp, true
	case "MutationFaint":
		return LDLevelConfigsId_MutationFaint, true
	case "MutationSlowSpeed":
		return LDLevelConfigsId_MutationSlowSpeed, true
	case "MutationHardest":
		return LDLevelConfigsId_MutationHardest, true
	case "MutationLevel1":
		return LDLevelConfigsId_MutationLevel1, true
	case "MutationBlodUp":
		return LDLevelConfigsId_MutationBlodUp, true
	case "MutationRandomTask":
		return LDLevelConfigsId_MutationRandomTask, true
	case "WorkShopTemp":
		return LDLevelConfigsId_WorkShopTemp, true
	case "Pvp":
		return LDLevelConfigsId_Pvp, true
	case "PvpFriendship":
		return LDLevelConfigsId_PvpFriendship, true
	case "Athena":
		return LDLevelConfigsId_Athena, true
	case "HuntZombie":
		return LDLevelConfigsId_HuntZombie, true

	}
	return -1, false

}

func GetLDLevelConfigsId_ByAlias(tp string) (LDLevelConfigsId, bool) {
	switch tp {
	case "新手教学关":
		return LDLevelConfigsId_Teach, true
	case "正式1关":
		return LDLevelConfigsId_Normal, true
	case "异变难度1":
		return LDLevelConfigsId_MutationDifficulty1, true
	case "异变物抗":
		return LDLevelConfigsId_MutationPhysicalUp, true
	case "异变魔抗":
		return LDLevelConfigsId_MutationMagicUp, true
	case "异变免疫眩晕":
		return LDLevelConfigsId_MutationFaint, true
	case "异变免疫减速":
		return LDLevelConfigsId_MutationSlowSpeed, true
	case "超高难度":
		return LDLevelConfigsId_MutationHardest, true
	case "异变1关":
		return LDLevelConfigsId_MutationLevel1, true
	case "异变血量提升":
		return LDLevelConfigsId_MutationBlodUp, true
	case "异变随机任务":
		return LDLevelConfigsId_MutationRandomTask, true
	case "创意工坊":
		return LDLevelConfigsId_WorkShopTemp, true
	case "pvp":
		return LDLevelConfigsId_Pvp, true
	case "友谊赛":
		return LDLevelConfigsId_PvpFriendship, true
	case "雅呆娜":
		return LDLevelConfigsId_Athena, true
	case "狩猎僵尸":
		return LDLevelConfigsId_HuntZombie, true

	}
	return -1, false

}

func (e LDLevelGeneratorType) String() string {

	if e == LDLevelGeneratorType_Main {
		return "Main"
	}
	if e == LDLevelGeneratorType_Boss {
		return "Boss"
	}
	if e == LDLevelGeneratorType_Hunt {
		return "Hunt"
	}

	return ""

}

func GetLDLevelGeneratorType(tp string) (LDLevelGeneratorType, bool) {
	switch tp {
	case "Main":
		return LDLevelGeneratorType_Main, true
	case "Boss":
		return LDLevelGeneratorType_Boss, true
	case "Hunt":
		return LDLevelGeneratorType_Hunt, true

	}
	return -1, false

}

func GetLDLevelGeneratorType_ByAlias(tp string) (LDLevelGeneratorType, bool) {
	switch tp {
	case "主线":
		return LDLevelGeneratorType_Main, true
	case "boss":
		return LDLevelGeneratorType_Boss, true
	case "狩猎":
		return LDLevelGeneratorType_Hunt, true

	}
	return -1, false

}

func (e LDRoleBattleType) String() string {

	if e == LDRoleBattleType_Front {
		return "Front"
	}
	if e == LDRoleBattleType_Back {
		return "Back"
	}

	return ""

}

func GetLDRoleBattleType(tp string) (LDRoleBattleType, bool) {
	switch tp {
	case "Front":
		return LDRoleBattleType_Front, true
	case "Back":
		return LDRoleBattleType_Back, true

	}
	return -1, false

}

func GetLDRoleBattleType_ByAlias(tp string) (LDRoleBattleType, bool) {
	switch tp {
	case "前排":
		return LDRoleBattleType_Front, true
	case "后排":
		return LDRoleBattleType_Back, true

	}
	return -1, false

}

func (e LDBossWarningBgType) String() string {

	if e == LDBossWarningBgType_Red {
		return "Red"
	}
	if e == LDBossWarningBgType_Yellow {
		return "Yellow"
	}
	if e == LDBossWarningBgType_Blue {
		return "Blue"
	}

	return ""

}

func GetLDBossWarningBgType(tp string) (LDBossWarningBgType, bool) {
	switch tp {
	case "Red":
		return LDBossWarningBgType_Red, true
	case "Yellow":
		return LDBossWarningBgType_Yellow, true
	case "Blue":
		return LDBossWarningBgType_Blue, true

	}
	return -1, false

}

func GetLDBossWarningBgType_ByAlias(tp string) (LDBossWarningBgType, bool) {
	switch tp {
	case "红":
		return LDBossWarningBgType_Red, true
	case "黄":
		return LDBossWarningBgType_Yellow, true
	case "蓝":
		return LDBossWarningBgType_Blue, true

	}
	return -1, false

}

func (e LDMonsterType) String() string {

	if e == LDMonsterType_One {
		return "One"
	}
	if e == LDMonsterType_Two {
		return "Two"
	}
	if e == LDMonsterType_Tree {
		return "Tree"
	}
	if e == LDMonsterType_Four {
		return "Four"
	}
	if e == LDMonsterType_Five {
		return "Five"
	}

	return ""

}

func GetLDMonsterType(tp string) (LDMonsterType, bool) {
	switch tp {
	case "One":
		return LDMonsterType_One, true
	case "Two":
		return LDMonsterType_Two, true
	case "Tree":
		return LDMonsterType_Tree, true
	case "Four":
		return LDMonsterType_Four, true
	case "Five":
		return LDMonsterType_Five, true

	}
	return -1, false

}

func GetLDMonsterType_ByAlias(tp string) (LDMonsterType, bool) {
	switch tp {
	case "类型1":
		return LDMonsterType_One, true
	case "类型2":
		return LDMonsterType_Two, true
	case "类型3":
		return LDMonsterType_Tree, true
	case "类型4":
		return LDMonsterType_Four, true
	case "类型5":
		return LDMonsterType_Five, true

	}
	return -1, false

}

func (e LDLevelEnchantType) String() string {

	if e == LDLevelEnchantType_GrayBlueEnchant {
		return "GrayBlueEnchant"
	}
	if e == LDLevelEnchantType_PurpleEnchant {
		return "PurpleEnchant"
	}
	if e == LDLevelEnchantType_LegendEnchant {
		return "LegendEnchant"
	}
	if e == LDLevelEnchantType_SummonEnchant {
		return "SummonEnchant"
	}

	return ""

}

func GetLDLevelEnchantType(tp string) (LDLevelEnchantType, bool) {
	switch tp {
	case "GrayBlueEnchant":
		return LDLevelEnchantType_GrayBlueEnchant, true
	case "PurpleEnchant":
		return LDLevelEnchantType_PurpleEnchant, true
	case "LegendEnchant":
		return LDLevelEnchantType_LegendEnchant, true
	case "SummonEnchant":
		return LDLevelEnchantType_SummonEnchant, true

	}
	return -1, false

}

func GetLDLevelEnchantType_ByAlias(tp string) (LDLevelEnchantType, bool) {
	switch tp {
	case "灰色、蓝色品质强化":
		return LDLevelEnchantType_GrayBlueEnchant, true
	case "紫色品质强化":
		return LDLevelEnchantType_PurpleEnchant, true
	case "传说品质强化":
		return LDLevelEnchantType_LegendEnchant, true
	case "召唤概率强化":
		return LDLevelEnchantType_SummonEnchant, true

	}
	return -1, false

}

func (e LevelGeneratorType) String() string {

	if e == LevelGeneratorType_Main {
		return "Main"
	}
	if e == LevelGeneratorType_Boss {
		return "Boss"
	}
	if e == LevelGeneratorType_Hunt {
		return "Hunt"
	}

	return ""

}

func GetLevelGeneratorType(tp string) (LevelGeneratorType, bool) {
	switch tp {
	case "Main":
		return LevelGeneratorType_Main, true
	case "Boss":
		return LevelGeneratorType_Boss, true
	case "Hunt":
		return LevelGeneratorType_Hunt, true

	}
	return -1, false

}

func GetLevelGeneratorType_ByAlias(tp string) (LevelGeneratorType, bool) {
	switch tp {
	case "主线":
		return LevelGeneratorType_Main, true
	case "boss":
		return LevelGeneratorType_Boss, true
	case "狩猎":
		return LevelGeneratorType_Hunt, true

	}
	return -1, false

}
