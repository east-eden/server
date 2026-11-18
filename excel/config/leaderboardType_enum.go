package config

// 排行榜-拳皇/LeaderboardType
type LeaderboardType int

const (
	LeaderboardType_BEGIN  LeaderboardType = 0 //枚举开始
	LeaderboardType_Hero   LeaderboardType = 0 //英雄榜
	LeaderboardType_Rush   LeaderboardType = 1 //冲刺榜
	LeaderboardType_Weekly LeaderboardType = 2 //周榜
	LeaderboardType_END    LeaderboardType = 3 //枚举结束
)

// 排行榜-拳皇/LeaderboardType
type LeaderboardSubType int

const (
	LeaderboardSubType_BEGIN   LeaderboardSubType = 1 //枚举开始
	LeaderboardSubType_Grow    LeaderboardSubType = 1 //成长
	LeaderboardSubType_Collect LeaderboardSubType = 2 //收集
	LeaderboardSubType_END     LeaderboardSubType = 3 //枚举结束
)

// 排行榜-拳皇/LeaderboardType
type LeaderboardId int

const (
	LeaderboardId_BEGIN            LeaderboardId = 1  //枚举开始
	LeaderboardId_NormalBegin      LeaderboardId = 1  //普通榜开始
	LeaderboardId_Stage            LeaderboardId = 1  //主线榜
	LeaderboardId_Tower            LeaderboardId = 2  //吉斯塔榜
	LeaderboardId_TD               LeaderboardId = 3  //TD榜
	LeaderboardId_ElfTower         LeaderboardId = 4  //南镇风云榜
	LeaderboardId_RarityHero       LeaderboardId = 5  //拥有X个X品质格斗家
	LeaderboardId_RarityArtifact   LeaderboardId = 6  //拥有X个X品质神器
	LeaderboardId_StarHero         LeaderboardId = 7  //拥有X个X星级格斗家
	LeaderboardId_RarityCollection LeaderboardId = 8  //拥有X个X品质收藏手办
	LeaderboardId_RushStage        LeaderboardId = 9  //冒险
	LeaderboardId_RushPower        LeaderboardId = 10 //战力
	LeaderboardId_RushTower        LeaderboardId = 11 //吉斯塔
	LeaderboardId_NormalEnd        LeaderboardId = 12 //普通榜结束
	LeaderboardId_WeeklyArena      LeaderboardId = 12 //竞技冲榜
	LeaderboardId_WeeklyAlliance   LeaderboardId = 13 //工会冲榜
	LeaderboardId_Stage2           LeaderboardId = 14 //主线榜2
	LeaderboardId_Tower2           LeaderboardId = 15 //吉斯塔榜2
	LeaderboardId_Fishing          LeaderboardId = 16 //钓鱼
	LeaderboardId_Fishing2         LeaderboardId = 17 //钓鱼2
	LeaderboardId_Recruit          LeaderboardId = 18 //招募冲榜
	LeaderboardId_Artifact         LeaderboardId = 19 //神器冲榜
	LeaderboardId_Mysticism        LeaderboardId = 20 //秘术冲榜
	LeaderboardId_BattlePower      LeaderboardId = 21 //战力冲榜
	LeaderboardId_StageLeaderboard LeaderboardId = 22 //赛事冲榜
	LeaderboardId_TowerLeaderboard LeaderboardId = 23 //吉斯塔冲榜
	LeaderboardId_END              LeaderboardId = 24 //枚举结束
)

func (e LeaderboardType) String() string {

	if e == LeaderboardType_Hero {
		return "Hero"
	}
	if e == LeaderboardType_Rush {
		return "Rush"
	}
	if e == LeaderboardType_Weekly {
		return "Weekly"
	}

	return ""

}

func GetLeaderboardType(tp string) (LeaderboardType, bool) {
	switch tp {
	case "Hero":
		return LeaderboardType_Hero, true
	case "Rush":
		return LeaderboardType_Rush, true
	case "Weekly":
		return LeaderboardType_Weekly, true

	}
	return -1, false

}

func GetLeaderboardType_ByAlias(tp string) (LeaderboardType, bool) {
	switch tp {
	case "英雄榜":
		return LeaderboardType_Hero, true
	case "冲刺榜":
		return LeaderboardType_Rush, true
	case "周榜":
		return LeaderboardType_Weekly, true

	}
	return -1, false

}

func (e LeaderboardSubType) String() string {

	if e == LeaderboardSubType_Grow {
		return "Grow"
	}
	if e == LeaderboardSubType_Collect {
		return "Collect"
	}

	return ""

}

func GetLeaderboardSubType(tp string) (LeaderboardSubType, bool) {
	switch tp {
	case "Grow":
		return LeaderboardSubType_Grow, true
	case "Collect":
		return LeaderboardSubType_Collect, true

	}
	return -1, false

}

func GetLeaderboardSubType_ByAlias(tp string) (LeaderboardSubType, bool) {
	switch tp {
	case "成长":
		return LeaderboardSubType_Grow, true
	case "收集":
		return LeaderboardSubType_Collect, true

	}
	return -1, false

}

func (e LeaderboardId) String() string {

	if e == LeaderboardId_NormalBegin {
		return "NormalBegin"
	}
	if e == LeaderboardId_Stage {
		return "Stage"
	}
	if e == LeaderboardId_Tower {
		return "Tower"
	}
	if e == LeaderboardId_TD {
		return "TD"
	}
	if e == LeaderboardId_ElfTower {
		return "ElfTower"
	}
	if e == LeaderboardId_RarityHero {
		return "RarityHero"
	}
	if e == LeaderboardId_RarityArtifact {
		return "RarityArtifact"
	}
	if e == LeaderboardId_StarHero {
		return "StarHero"
	}
	if e == LeaderboardId_RarityCollection {
		return "RarityCollection"
	}
	if e == LeaderboardId_RushStage {
		return "RushStage"
	}
	if e == LeaderboardId_RushPower {
		return "RushPower"
	}
	if e == LeaderboardId_RushTower {
		return "RushTower"
	}
	if e == LeaderboardId_NormalEnd {
		return "NormalEnd"
	}
	if e == LeaderboardId_WeeklyArena {
		return "WeeklyArena"
	}
	if e == LeaderboardId_WeeklyAlliance {
		return "WeeklyAlliance"
	}
	if e == LeaderboardId_Stage2 {
		return "Stage2"
	}
	if e == LeaderboardId_Tower2 {
		return "Tower2"
	}
	if e == LeaderboardId_Fishing {
		return "Fishing"
	}
	if e == LeaderboardId_Fishing2 {
		return "Fishing2"
	}
	if e == LeaderboardId_Recruit {
		return "Recruit"
	}
	if e == LeaderboardId_Artifact {
		return "Artifact"
	}
	if e == LeaderboardId_Mysticism {
		return "Mysticism"
	}
	if e == LeaderboardId_BattlePower {
		return "BattlePower"
	}
	if e == LeaderboardId_StageLeaderboard {
		return "StageLeaderboard"
	}
	if e == LeaderboardId_TowerLeaderboard {
		return "TowerLeaderboard"
	}

	return ""

}

func GetLeaderboardId(tp string) (LeaderboardId, bool) {
	switch tp {
	case "NormalBegin":
		return LeaderboardId_NormalBegin, true
	case "Stage":
		return LeaderboardId_Stage, true
	case "Tower":
		return LeaderboardId_Tower, true
	case "TD":
		return LeaderboardId_TD, true
	case "ElfTower":
		return LeaderboardId_ElfTower, true
	case "RarityHero":
		return LeaderboardId_RarityHero, true
	case "RarityArtifact":
		return LeaderboardId_RarityArtifact, true
	case "StarHero":
		return LeaderboardId_StarHero, true
	case "RarityCollection":
		return LeaderboardId_RarityCollection, true
	case "RushStage":
		return LeaderboardId_RushStage, true
	case "RushPower":
		return LeaderboardId_RushPower, true
	case "RushTower":
		return LeaderboardId_RushTower, true
	case "NormalEnd":
		return LeaderboardId_NormalEnd, true
	case "WeeklyArena":
		return LeaderboardId_WeeklyArena, true
	case "WeeklyAlliance":
		return LeaderboardId_WeeklyAlliance, true
	case "Stage2":
		return LeaderboardId_Stage2, true
	case "Tower2":
		return LeaderboardId_Tower2, true
	case "Fishing":
		return LeaderboardId_Fishing, true
	case "Fishing2":
		return LeaderboardId_Fishing2, true
	case "Recruit":
		return LeaderboardId_Recruit, true
	case "Artifact":
		return LeaderboardId_Artifact, true
	case "Mysticism":
		return LeaderboardId_Mysticism, true
	case "BattlePower":
		return LeaderboardId_BattlePower, true
	case "StageLeaderboard":
		return LeaderboardId_StageLeaderboard, true
	case "TowerLeaderboard":
		return LeaderboardId_TowerLeaderboard, true

	}
	return -1, false

}

func GetLeaderboardId_ByAlias(tp string) (LeaderboardId, bool) {
	switch tp {
	case "普通榜开始":
		return LeaderboardId_NormalBegin, true
	case "主线榜":
		return LeaderboardId_Stage, true
	case "吉斯塔榜":
		return LeaderboardId_Tower, true
	case "TD榜":
		return LeaderboardId_TD, true
	case "南镇风云榜":
		return LeaderboardId_ElfTower, true
	case "拥有X个X品质格斗家":
		return LeaderboardId_RarityHero, true
	case "拥有X个X品质神器":
		return LeaderboardId_RarityArtifact, true
	case "拥有X个X星级格斗家":
		return LeaderboardId_StarHero, true
	case "拥有X个X品质收藏手办":
		return LeaderboardId_RarityCollection, true
	case "冒险":
		return LeaderboardId_RushStage, true
	case "战力":
		return LeaderboardId_RushPower, true
	case "吉斯塔":
		return LeaderboardId_RushTower, true
	case "普通榜结束":
		return LeaderboardId_NormalEnd, true
	case "竞技冲榜":
		return LeaderboardId_WeeklyArena, true
	case "工会冲榜":
		return LeaderboardId_WeeklyAlliance, true
	case "主线榜2":
		return LeaderboardId_Stage2, true
	case "吉斯塔榜2":
		return LeaderboardId_Tower2, true
	case "钓鱼":
		return LeaderboardId_Fishing, true
	case "钓鱼2":
		return LeaderboardId_Fishing2, true
	case "招募冲榜":
		return LeaderboardId_Recruit, true
	case "神器冲榜":
		return LeaderboardId_Artifact, true
	case "秘术冲榜":
		return LeaderboardId_Mysticism, true
	case "战力冲榜":
		return LeaderboardId_BattlePower, true
	case "赛事冲榜":
		return LeaderboardId_StageLeaderboard, true
	case "吉斯塔冲榜":
		return LeaderboardId_TowerLeaderboard, true

	}
	return -1, false

}
