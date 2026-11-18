package config

// 玩家表-拳皇/settleType
type SettleType int

const (
	SettleType_BEGIN                          SettleType = 1  //枚举开始
	SettleType_LeaderboardWeeklyArena         SettleType = 1  //竞技冲榜
	SettleType_LeaderboardWeeklyAlliance      SettleType = 2  //公会冲榜
	SettleType_LeaderboardBattlePower         SettleType = 3  //战力冲榜
	SettleType_LeaderboardStage               SettleType = 4  //赛事冲榜
	SettleType_LeaderboardTower               SettleType = 5  //吉斯塔冲榜
	SettleType_LeaderboardRecruit             SettleType = 6  //招募冲榜
	SettleType_LeaderboardArtifact            SettleType = 7  //神器冲榜
	SettleType_LeaderboardMysticism           SettleType = 8  //秘术冲榜
	SettleType_BountyPersonAward              SettleType = 9  //赏金战个人奖励结算
	SettleType_BountyAllianceAward            SettleType = 10 //赏金战公会奖励结算
	SettleType_BountyAllianceAwardMail        SettleType = 11 //赏金战公会奖励邮件
	SettleType_BountyPersonAwardMail          SettleType = 12 //赏金战个人奖励邮件
	SettleType_BountyChallenge                SettleType = 13 //赏金战boss挑战
	SettleType_RugalRecruit                   SettleType = 14 //卢卡尔招募
	SettleType_TeamLeagueScoreAward           SettleType = 15 //团赛积分奖励
	SettleType_TeamLeagueRankAward            SettleType = 16 //团赛排行奖励
	SettleType_UpRecruitRank                  SettleType = 17 //Up招募排行结算
	SettleType_WorldBossRank                  SettleType = 18 //世界boss排行结算
	SettleType_DasheBackRecruit               SettleType = 19 //大蛇返场招募
	SettleType_Crazy8Recruit                  SettleType = 20 //疯八招募
	SettleType_RugalBackRecruit               SettleType = 21 //卢卡尔返场招募
	SettleType_KingWarMapAllianceScore        SettleType = 22 //皇城争锋公会积分排名奖励结算
	SettleType_KingWarMapPlayerScore          SettleType = 23 //皇城争锋个人积分排名奖励结算
	SettleType_KingWarMapMercenaryScore       SettleType = 24 //皇城争锋佣兵积分排名奖励结算
	SettleType_KingWarGuessMail               SettleType = 25 //皇城争锋竞猜邮件结算
	SettleType_KingWarWeekQualificationMail   SettleType = 26 //皇城争锋周赛入围邮件
	SettleType_KingWarWeekUnqualificationMail SettleType = 27 //皇城争锋周赛未入围邮件
	SettleType_KingWarFinalQualificationMail  SettleType = 28 //皇城争锋决赛入围邮件
	SettleType_KingWarMercenaryAssignMail     SettleType = 29 //皇城争锋佣兵分配奖励邮件
	SettleType_Crazy8BackRecruit              SettleType = 30 //疯八返场招募
	SettleType_CrazyLeonaRecruit              SettleType = 31 //疯莉安娜招募
	SettleType_CrazyLeonaBackRecruit          SettleType = 32 //狂暴莉安娜返场招募
	SettleType_END                            SettleType = 33 //枚举结束
)

func (e SettleType) String() string {

	if e == SettleType_LeaderboardWeeklyArena {
		return "LeaderboardWeeklyArena"
	}
	if e == SettleType_LeaderboardWeeklyAlliance {
		return "LeaderboardWeeklyAlliance"
	}
	if e == SettleType_LeaderboardBattlePower {
		return "LeaderboardBattlePower"
	}
	if e == SettleType_LeaderboardStage {
		return "LeaderboardStage"
	}
	if e == SettleType_LeaderboardTower {
		return "LeaderboardTower"
	}
	if e == SettleType_LeaderboardRecruit {
		return "LeaderboardRecruit"
	}
	if e == SettleType_LeaderboardArtifact {
		return "LeaderboardArtifact"
	}
	if e == SettleType_LeaderboardMysticism {
		return "LeaderboardMysticism"
	}
	if e == SettleType_BountyPersonAward {
		return "BountyPersonAward"
	}
	if e == SettleType_BountyAllianceAward {
		return "BountyAllianceAward"
	}
	if e == SettleType_BountyAllianceAwardMail {
		return "BountyAllianceAwardMail"
	}
	if e == SettleType_BountyPersonAwardMail {
		return "BountyPersonAwardMail"
	}
	if e == SettleType_BountyChallenge {
		return "BountyChallenge"
	}
	if e == SettleType_RugalRecruit {
		return "RugalRecruit"
	}
	if e == SettleType_TeamLeagueScoreAward {
		return "TeamLeagueScoreAward"
	}
	if e == SettleType_TeamLeagueRankAward {
		return "TeamLeagueRankAward"
	}
	if e == SettleType_UpRecruitRank {
		return "UpRecruitRank"
	}
	if e == SettleType_WorldBossRank {
		return "WorldBossRank"
	}
	if e == SettleType_DasheBackRecruit {
		return "DasheBackRecruit"
	}
	if e == SettleType_Crazy8Recruit {
		return "Crazy8Recruit"
	}
	if e == SettleType_RugalBackRecruit {
		return "RugalBackRecruit"
	}
	if e == SettleType_KingWarMapAllianceScore {
		return "KingWarMapAllianceScore"
	}
	if e == SettleType_KingWarMapPlayerScore {
		return "KingWarMapPlayerScore"
	}
	if e == SettleType_KingWarMapMercenaryScore {
		return "KingWarMapMercenaryScore"
	}
	if e == SettleType_KingWarGuessMail {
		return "KingWarGuessMail"
	}
	if e == SettleType_KingWarWeekQualificationMail {
		return "KingWarWeekQualificationMail"
	}
	if e == SettleType_KingWarWeekUnqualificationMail {
		return "KingWarWeekUnqualificationMail"
	}
	if e == SettleType_KingWarFinalQualificationMail {
		return "KingWarFinalQualificationMail"
	}
	if e == SettleType_KingWarMercenaryAssignMail {
		return "KingWarMercenaryAssignMail"
	}
	if e == SettleType_Crazy8BackRecruit {
		return "Crazy8BackRecruit"
	}
	if e == SettleType_CrazyLeonaRecruit {
		return "CrazyLeonaRecruit"
	}
	if e == SettleType_CrazyLeonaBackRecruit {
		return "CrazyLeonaBackRecruit"
	}

	return ""

}

func GetSettleType(tp string) (SettleType, bool) {
	switch tp {
	case "LeaderboardWeeklyArena":
		return SettleType_LeaderboardWeeklyArena, true
	case "LeaderboardWeeklyAlliance":
		return SettleType_LeaderboardWeeklyAlliance, true
	case "LeaderboardBattlePower":
		return SettleType_LeaderboardBattlePower, true
	case "LeaderboardStage":
		return SettleType_LeaderboardStage, true
	case "LeaderboardTower":
		return SettleType_LeaderboardTower, true
	case "LeaderboardRecruit":
		return SettleType_LeaderboardRecruit, true
	case "LeaderboardArtifact":
		return SettleType_LeaderboardArtifact, true
	case "LeaderboardMysticism":
		return SettleType_LeaderboardMysticism, true
	case "BountyPersonAward":
		return SettleType_BountyPersonAward, true
	case "BountyAllianceAward":
		return SettleType_BountyAllianceAward, true
	case "BountyAllianceAwardMail":
		return SettleType_BountyAllianceAwardMail, true
	case "BountyPersonAwardMail":
		return SettleType_BountyPersonAwardMail, true
	case "BountyChallenge":
		return SettleType_BountyChallenge, true
	case "RugalRecruit":
		return SettleType_RugalRecruit, true
	case "TeamLeagueScoreAward":
		return SettleType_TeamLeagueScoreAward, true
	case "TeamLeagueRankAward":
		return SettleType_TeamLeagueRankAward, true
	case "UpRecruitRank":
		return SettleType_UpRecruitRank, true
	case "WorldBossRank":
		return SettleType_WorldBossRank, true
	case "DasheBackRecruit":
		return SettleType_DasheBackRecruit, true
	case "Crazy8Recruit":
		return SettleType_Crazy8Recruit, true
	case "RugalBackRecruit":
		return SettleType_RugalBackRecruit, true
	case "KingWarMapAllianceScore":
		return SettleType_KingWarMapAllianceScore, true
	case "KingWarMapPlayerScore":
		return SettleType_KingWarMapPlayerScore, true
	case "KingWarMapMercenaryScore":
		return SettleType_KingWarMapMercenaryScore, true
	case "KingWarGuessMail":
		return SettleType_KingWarGuessMail, true
	case "KingWarWeekQualificationMail":
		return SettleType_KingWarWeekQualificationMail, true
	case "KingWarWeekUnqualificationMail":
		return SettleType_KingWarWeekUnqualificationMail, true
	case "KingWarFinalQualificationMail":
		return SettleType_KingWarFinalQualificationMail, true
	case "KingWarMercenaryAssignMail":
		return SettleType_KingWarMercenaryAssignMail, true
	case "Crazy8BackRecruit":
		return SettleType_Crazy8BackRecruit, true
	case "CrazyLeonaRecruit":
		return SettleType_CrazyLeonaRecruit, true
	case "CrazyLeonaBackRecruit":
		return SettleType_CrazyLeonaBackRecruit, true

	}
	return -1, false

}

func GetSettleType_ByAlias(tp string) (SettleType, bool) {
	switch tp {
	case "竞技冲榜":
		return SettleType_LeaderboardWeeklyArena, true
	case "公会冲榜":
		return SettleType_LeaderboardWeeklyAlliance, true
	case "战力冲榜":
		return SettleType_LeaderboardBattlePower, true
	case "赛事冲榜":
		return SettleType_LeaderboardStage, true
	case "吉斯塔冲榜":
		return SettleType_LeaderboardTower, true
	case "招募冲榜":
		return SettleType_LeaderboardRecruit, true
	case "神器冲榜":
		return SettleType_LeaderboardArtifact, true
	case "秘术冲榜":
		return SettleType_LeaderboardMysticism, true
	case "赏金战个人奖励结算":
		return SettleType_BountyPersonAward, true
	case "赏金战公会奖励结算":
		return SettleType_BountyAllianceAward, true
	case "赏金战公会奖励邮件":
		return SettleType_BountyAllianceAwardMail, true
	case "赏金战个人奖励邮件":
		return SettleType_BountyPersonAwardMail, true
	case "赏金战boss挑战":
		return SettleType_BountyChallenge, true
	case "卢卡尔招募":
		return SettleType_RugalRecruit, true
	case "团赛积分奖励":
		return SettleType_TeamLeagueScoreAward, true
	case "团赛排行奖励":
		return SettleType_TeamLeagueRankAward, true
	case "Up招募排行结算":
		return SettleType_UpRecruitRank, true
	case "世界boss排行结算":
		return SettleType_WorldBossRank, true
	case "大蛇返场招募":
		return SettleType_DasheBackRecruit, true
	case "疯八招募":
		return SettleType_Crazy8Recruit, true
	case "卢卡尔返场招募":
		return SettleType_RugalBackRecruit, true
	case "皇城争锋公会积分排名奖励结算":
		return SettleType_KingWarMapAllianceScore, true
	case "皇城争锋个人积分排名奖励结算":
		return SettleType_KingWarMapPlayerScore, true
	case "皇城争锋佣兵积分排名奖励结算":
		return SettleType_KingWarMapMercenaryScore, true
	case "皇城争锋竞猜邮件结算":
		return SettleType_KingWarGuessMail, true
	case "皇城争锋周赛入围邮件":
		return SettleType_KingWarWeekQualificationMail, true
	case "皇城争锋周赛未入围邮件":
		return SettleType_KingWarWeekUnqualificationMail, true
	case "皇城争锋决赛入围邮件":
		return SettleType_KingWarFinalQualificationMail, true
	case "皇城争锋佣兵分配奖励邮件":
		return SettleType_KingWarMercenaryAssignMail, true
	case "疯八返场招募":
		return SettleType_Crazy8BackRecruit, true
	case "疯莉安娜招募":
		return SettleType_CrazyLeonaRecruit, true
	case "狂暴莉安娜返场招募":
		return SettleType_CrazyLeonaBackRecruit, true

	}
	return -1, false

}
