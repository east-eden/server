package config

// 关卡表-拳皇/LevelType
type LevelType int

const (
	LevelType_BEGIN       LevelType = 1  //枚举开始
	LevelType_Normal      LevelType = 1  //普通
	LevelType_Hard        LevelType = 2  //困难
	LevelType_HardPlus    LevelType = 3  //噩梦
	LevelType_Test        LevelType = 4  //测试
	LevelType_ACT         LevelType = 5  //副玩法
	LevelType_Purgatory   LevelType = 6  //炼狱
	LevelType_Destroy     LevelType = 7  //毁灭
	LevelType_TD          LevelType = 8  //塔防
	LevelType_Chaos       LevelType = 9  //混沌
	LevelType_UpTest      LevelType = 10 //招募试玩
	LevelType_FeedNewBee  LevelType = 11 //Feed新手
	LevelType_Vanity      LevelType = 12 //虚空
	LevelType_QuickBattle LevelType = 13 //快速战斗
	LevelType_KofStar     LevelType = 14 //拳皇大明星
	LevelType_SpTDBoss    LevelType = 15 //南镇Boss战
	LevelType_Doom        LevelType = 16 //终焉
	LevelType_SnowMan     LevelType = 17 //雪人BOSS
	LevelType_Silence     LevelType = 18 //永寂
	LevelType_END         LevelType = 19 //枚举结束
)

// 关卡表-拳皇/LevelType
type RoleBattleType int

const (
	RoleBattleType_BEGIN RoleBattleType = 1 //枚举开始
	RoleBattleType_Front RoleBattleType = 1 //前排
	RoleBattleType_Back  RoleBattleType = 2 //后排
	RoleBattleType_END   RoleBattleType = 3 //枚举结束
)

// 关卡表-拳皇/LevelType
type MonsterType int

const (
	MonsterType_BEGIN MonsterType = 1 //枚举开始
	MonsterType_One   MonsterType = 1 //类型1
	MonsterType_Two   MonsterType = 2 //类型2
	MonsterType_Tree  MonsterType = 3 //类型3
	MonsterType_Four  MonsterType = 4 //类型4
	MonsterType_Five  MonsterType = 5 //类型5
	MonsterType_END   MonsterType = 6 //枚举结束
)

// 关卡表-拳皇/LevelType
type BossWarningBgType int

const (
	BossWarningBgType_BEGIN  BossWarningBgType = 0 //枚举开始
	BossWarningBgType_Red    BossWarningBgType = 0 //红
	BossWarningBgType_Yellow BossWarningBgType = 1 //黄
	BossWarningBgType_Blue   BossWarningBgType = 2 //蓝
	BossWarningBgType_END    BossWarningBgType = 3 //枚举结束
)

// 关卡表-拳皇/LevelType
type WeekdayType int

const (
	WeekdayType_BEGIN    WeekdayType = 0 //枚举开始
	WeekdayType_Sunday   WeekdayType = 0 //周日
	WeekdayType_Monday   WeekdayType = 1 //周一
	WeekdayType_Tuesday  WeekdayType = 2 //周二
	WeekdayType_Wenseday WeekdayType = 3 //周三
	WeekdayType_Thursday WeekdayType = 4 //周四
	WeekdayType_Friday   WeekdayType = 5 //周五
	WeekdayType_Saturday WeekdayType = 6 //周六
	WeekdayType_END      WeekdayType = 7 //枚举结束
)

// 关卡表-拳皇/LevelType
type LevelMemberCamp int

const (
	LevelMemberCamp_BEGIN        LevelMemberCamp = 1  //枚举开始
	LevelMemberCamp_Self         LevelMemberCamp = 1  //己方英雄
	LevelMemberCamp_Others       LevelMemberCamp = 2  //对方英雄
	LevelMemberCamp_All          LevelMemberCamp = 3  //全部
	LevelMemberCamp_SelfMonster  LevelMemberCamp = 4  //己方怪物
	LevelMemberCamp_EnemyMonster LevelMemberCamp = 5  //对方怪物
	LevelMemberCamp_SelfWall     LevelMemberCamp = 6  //己方城墙
	LevelMemberCamp_EnemyWall    LevelMemberCamp = 7  //对方城墙
	LevelMemberCamp_EnemyAll     LevelMemberCamp = 8  //对方全体
	LevelMemberCamp_SelfAll      LevelMemberCamp = 9  //己方全体
	LevelMemberCamp_END          LevelMemberCamp = 10 //枚举结束
)

// 关卡表-拳皇/LevelType
type LevelImmueType int

const (
	LevelImmueType_BEGIN      LevelImmueType = 1 //枚举开始
	LevelImmueType_Stun       LevelImmueType = 1 //眩晕
	LevelImmueType_Decelerate LevelImmueType = 2 //减速
	LevelImmueType_END        LevelImmueType = 3 //枚举结束
)

// 关卡表-拳皇/LevelType
type LevelGridOwnerType int

const (
	LevelGridOwnerType_BEGIN   LevelGridOwnerType = 1 //枚举开始
	LevelGridOwnerType_Blue    LevelGridOwnerType = 1 //蓝方
	LevelGridOwnerType_Red     LevelGridOwnerType = 2 //红方
	LevelGridOwnerType_Neutral LevelGridOwnerType = 3 //中立
	LevelGridOwnerType_END     LevelGridOwnerType = 4 //枚举结束
)

func (e LevelType) String() string {

	if e == LevelType_Normal {
		return "Normal"
	}
	if e == LevelType_Hard {
		return "Hard"
	}
	if e == LevelType_HardPlus {
		return "HardPlus"
	}
	if e == LevelType_Test {
		return "Test"
	}
	if e == LevelType_ACT {
		return "ACT"
	}
	if e == LevelType_Purgatory {
		return "Purgatory"
	}
	if e == LevelType_Destroy {
		return "Destroy"
	}
	if e == LevelType_TD {
		return "TD"
	}
	if e == LevelType_Chaos {
		return "Chaos"
	}
	if e == LevelType_UpTest {
		return "UpTest"
	}
	if e == LevelType_FeedNewBee {
		return "FeedNewBee"
	}
	if e == LevelType_Vanity {
		return "Vanity"
	}
	if e == LevelType_QuickBattle {
		return "QuickBattle"
	}
	if e == LevelType_KofStar {
		return "KofStar"
	}
	if e == LevelType_SpTDBoss {
		return "SpTDBoss"
	}
	if e == LevelType_Doom {
		return "Doom"
	}
	if e == LevelType_SnowMan {
		return "SnowMan"
	}
	if e == LevelType_Silence {
		return "Silence"
	}

	return ""

}

func GetLevelType(tp string) (LevelType, bool) {
	switch tp {
	case "Normal":
		return LevelType_Normal, true
	case "Hard":
		return LevelType_Hard, true
	case "HardPlus":
		return LevelType_HardPlus, true
	case "Test":
		return LevelType_Test, true
	case "ACT":
		return LevelType_ACT, true
	case "Purgatory":
		return LevelType_Purgatory, true
	case "Destroy":
		return LevelType_Destroy, true
	case "TD":
		return LevelType_TD, true
	case "Chaos":
		return LevelType_Chaos, true
	case "UpTest":
		return LevelType_UpTest, true
	case "FeedNewBee":
		return LevelType_FeedNewBee, true
	case "Vanity":
		return LevelType_Vanity, true
	case "QuickBattle":
		return LevelType_QuickBattle, true
	case "KofStar":
		return LevelType_KofStar, true
	case "SpTDBoss":
		return LevelType_SpTDBoss, true
	case "Doom":
		return LevelType_Doom, true
	case "SnowMan":
		return LevelType_SnowMan, true
	case "Silence":
		return LevelType_Silence, true

	}
	return -1, false

}

func GetLevelType_ByAlias(tp string) (LevelType, bool) {
	switch tp {
	case "普通":
		return LevelType_Normal, true
	case "困难":
		return LevelType_Hard, true
	case "噩梦":
		return LevelType_HardPlus, true
	case "测试":
		return LevelType_Test, true
	case "副玩法":
		return LevelType_ACT, true
	case "炼狱":
		return LevelType_Purgatory, true
	case "毁灭":
		return LevelType_Destroy, true
	case "塔防":
		return LevelType_TD, true
	case "混沌":
		return LevelType_Chaos, true
	case "招募试玩":
		return LevelType_UpTest, true
	case "Feed新手":
		return LevelType_FeedNewBee, true
	case "虚空":
		return LevelType_Vanity, true
	case "快速战斗":
		return LevelType_QuickBattle, true
	case "拳皇大明星":
		return LevelType_KofStar, true
	case "南镇Boss战":
		return LevelType_SpTDBoss, true
	case "终焉":
		return LevelType_Doom, true
	case "雪人BOSS":
		return LevelType_SnowMan, true
	case "永寂":
		return LevelType_Silence, true

	}
	return -1, false

}

func (e RoleBattleType) String() string {

	if e == RoleBattleType_Front {
		return "Front"
	}
	if e == RoleBattleType_Back {
		return "Back"
	}

	return ""

}

func GetRoleBattleType(tp string) (RoleBattleType, bool) {
	switch tp {
	case "Front":
		return RoleBattleType_Front, true
	case "Back":
		return RoleBattleType_Back, true

	}
	return -1, false

}

func GetRoleBattleType_ByAlias(tp string) (RoleBattleType, bool) {
	switch tp {
	case "前排":
		return RoleBattleType_Front, true
	case "后排":
		return RoleBattleType_Back, true

	}
	return -1, false

}

func (e MonsterType) String() string {

	if e == MonsterType_One {
		return "One"
	}
	if e == MonsterType_Two {
		return "Two"
	}
	if e == MonsterType_Tree {
		return "Tree"
	}
	if e == MonsterType_Four {
		return "Four"
	}
	if e == MonsterType_Five {
		return "Five"
	}

	return ""

}

func GetMonsterType(tp string) (MonsterType, bool) {
	switch tp {
	case "One":
		return MonsterType_One, true
	case "Two":
		return MonsterType_Two, true
	case "Tree":
		return MonsterType_Tree, true
	case "Four":
		return MonsterType_Four, true
	case "Five":
		return MonsterType_Five, true

	}
	return -1, false

}

func GetMonsterType_ByAlias(tp string) (MonsterType, bool) {
	switch tp {
	case "类型1":
		return MonsterType_One, true
	case "类型2":
		return MonsterType_Two, true
	case "类型3":
		return MonsterType_Tree, true
	case "类型4":
		return MonsterType_Four, true
	case "类型5":
		return MonsterType_Five, true

	}
	return -1, false

}

func (e BossWarningBgType) String() string {

	if e == BossWarningBgType_Red {
		return "Red"
	}
	if e == BossWarningBgType_Yellow {
		return "Yellow"
	}
	if e == BossWarningBgType_Blue {
		return "Blue"
	}

	return ""

}

func GetBossWarningBgType(tp string) (BossWarningBgType, bool) {
	switch tp {
	case "Red":
		return BossWarningBgType_Red, true
	case "Yellow":
		return BossWarningBgType_Yellow, true
	case "Blue":
		return BossWarningBgType_Blue, true

	}
	return -1, false

}

func GetBossWarningBgType_ByAlias(tp string) (BossWarningBgType, bool) {
	switch tp {
	case "红":
		return BossWarningBgType_Red, true
	case "黄":
		return BossWarningBgType_Yellow, true
	case "蓝":
		return BossWarningBgType_Blue, true

	}
	return -1, false

}

func (e WeekdayType) String() string {

	if e == WeekdayType_Sunday {
		return "Sunday"
	}
	if e == WeekdayType_Monday {
		return "Monday"
	}
	if e == WeekdayType_Tuesday {
		return "Tuesday"
	}
	if e == WeekdayType_Wenseday {
		return "Wenseday"
	}
	if e == WeekdayType_Thursday {
		return "Thursday"
	}
	if e == WeekdayType_Friday {
		return "Friday"
	}
	if e == WeekdayType_Saturday {
		return "Saturday"
	}

	return ""

}

func GetWeekdayType(tp string) (WeekdayType, bool) {
	switch tp {
	case "Sunday":
		return WeekdayType_Sunday, true
	case "Monday":
		return WeekdayType_Monday, true
	case "Tuesday":
		return WeekdayType_Tuesday, true
	case "Wenseday":
		return WeekdayType_Wenseday, true
	case "Thursday":
		return WeekdayType_Thursday, true
	case "Friday":
		return WeekdayType_Friday, true
	case "Saturday":
		return WeekdayType_Saturday, true

	}
	return -1, false

}

func GetWeekdayType_ByAlias(tp string) (WeekdayType, bool) {
	switch tp {
	case "周日":
		return WeekdayType_Sunday, true
	case "周一":
		return WeekdayType_Monday, true
	case "周二":
		return WeekdayType_Tuesday, true
	case "周三":
		return WeekdayType_Wenseday, true
	case "周四":
		return WeekdayType_Thursday, true
	case "周五":
		return WeekdayType_Friday, true
	case "周六":
		return WeekdayType_Saturday, true

	}
	return -1, false

}

func (e LevelMemberCamp) String() string {

	if e == LevelMemberCamp_Self {
		return "Self"
	}
	if e == LevelMemberCamp_Others {
		return "Others"
	}
	if e == LevelMemberCamp_All {
		return "All"
	}
	if e == LevelMemberCamp_SelfMonster {
		return "SelfMonster"
	}
	if e == LevelMemberCamp_EnemyMonster {
		return "EnemyMonster"
	}
	if e == LevelMemberCamp_SelfWall {
		return "SelfWall"
	}
	if e == LevelMemberCamp_EnemyWall {
		return "EnemyWall"
	}
	if e == LevelMemberCamp_EnemyAll {
		return "EnemyAll"
	}
	if e == LevelMemberCamp_SelfAll {
		return "SelfAll"
	}

	return ""

}

func GetLevelMemberCamp(tp string) (LevelMemberCamp, bool) {
	switch tp {
	case "Self":
		return LevelMemberCamp_Self, true
	case "Others":
		return LevelMemberCamp_Others, true
	case "All":
		return LevelMemberCamp_All, true
	case "SelfMonster":
		return LevelMemberCamp_SelfMonster, true
	case "EnemyMonster":
		return LevelMemberCamp_EnemyMonster, true
	case "SelfWall":
		return LevelMemberCamp_SelfWall, true
	case "EnemyWall":
		return LevelMemberCamp_EnemyWall, true
	case "EnemyAll":
		return LevelMemberCamp_EnemyAll, true
	case "SelfAll":
		return LevelMemberCamp_SelfAll, true

	}
	return -1, false

}

func GetLevelMemberCamp_ByAlias(tp string) (LevelMemberCamp, bool) {
	switch tp {
	case "己方英雄":
		return LevelMemberCamp_Self, true
	case "对方英雄":
		return LevelMemberCamp_Others, true
	case "全部":
		return LevelMemberCamp_All, true
	case "己方怪物":
		return LevelMemberCamp_SelfMonster, true
	case "对方怪物":
		return LevelMemberCamp_EnemyMonster, true
	case "己方城墙":
		return LevelMemberCamp_SelfWall, true
	case "对方城墙":
		return LevelMemberCamp_EnemyWall, true
	case "对方全体":
		return LevelMemberCamp_EnemyAll, true
	case "己方全体":
		return LevelMemberCamp_SelfAll, true

	}
	return -1, false

}

func (e LevelImmueType) String() string {

	if e == LevelImmueType_Stun {
		return "Stun"
	}
	if e == LevelImmueType_Decelerate {
		return "Decelerate"
	}

	return ""

}

func GetLevelImmueType(tp string) (LevelImmueType, bool) {
	switch tp {
	case "Stun":
		return LevelImmueType_Stun, true
	case "Decelerate":
		return LevelImmueType_Decelerate, true

	}
	return -1, false

}

func GetLevelImmueType_ByAlias(tp string) (LevelImmueType, bool) {
	switch tp {
	case "眩晕":
		return LevelImmueType_Stun, true
	case "减速":
		return LevelImmueType_Decelerate, true

	}
	return -1, false

}

func (e LevelGridOwnerType) String() string {

	if e == LevelGridOwnerType_Blue {
		return "Blue"
	}
	if e == LevelGridOwnerType_Red {
		return "Red"
	}
	if e == LevelGridOwnerType_Neutral {
		return "Neutral"
	}

	return ""

}

func GetLevelGridOwnerType(tp string) (LevelGridOwnerType, bool) {
	switch tp {
	case "Blue":
		return LevelGridOwnerType_Blue, true
	case "Red":
		return LevelGridOwnerType_Red, true
	case "Neutral":
		return LevelGridOwnerType_Neutral, true

	}
	return -1, false

}

func GetLevelGridOwnerType_ByAlias(tp string) (LevelGridOwnerType, bool) {
	switch tp {
	case "蓝方":
		return LevelGridOwnerType_Blue, true
	case "红方":
		return LevelGridOwnerType_Red, true
	case "中立":
		return LevelGridOwnerType_Neutral, true

	}
	return -1, false

}
