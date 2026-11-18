package config

// 跨服配置-拳皇/serverGameplayType
type GameplayType int

const (
	GameplayType_BEGIN                     GameplayType = 1  //枚举开始
	GameplayType_ManorGrab                 GameplayType = 1  //领地掠夺
	GameplayType_Arena                     GameplayType = 2  //竞技场
	GameplayType_Bounty                    GameplayType = 3  //赏金战
	GameplayType_CrossArena                GameplayType = 4  //跨服竞技32强
	GameplayType_CampBattle                GameplayType = 5  //阵营战
	GameplayType_CityWar                   GameplayType = 6  //城市争夺战
	GameplayType_AllianceDuel              GameplayType = 7  //公会对决
	GameplayType_Alliance                  GameplayType = 8  //公会
	GameplayType_WorldBoss                 GameplayType = 9  //世界boss
	GameplayType_KofRoad                   GameplayType = 10 //拳皇之路
	GameplayType_ELFTower                  GameplayType = 11 //塔防精灵
	GameplayType_Fighting                  GameplayType = 12 //格斗圣殿
	GameplayType_DailyTopic                GameplayType = 13 //每日话题
	GameplayType_LeaderboardGameplay       GameplayType = 14 //新榜玩法
	GameplayType_LeaderboardArena          GameplayType = 15 //新榜竞技场
	GameplayType_LeaderboardActiveAlliance GameplayType = 16 //新榜活跃公会
	GameplayType_UpNurture                 GameplayType = 17 //Up养成
	GameplayType_SendShare                 GameplayType = 18 //跨服频道
	GameplayType_Strongest108              GameplayType = 19 //最强108
	GameplayType_LeaderboardRecruit        GameplayType = 20 //招募冲榜
	GameplayType_LeaderboardArtifact       GameplayType = 21 //神器冲榜
	GameplayType_LeaderboardMysticism      GameplayType = 22 //秘术冲榜
	GameplayType_LeaderboardBattlePower    GameplayType = 23 //战力冲榜
	GameplayType_Arcade                    GameplayType = 24 //街机厅
	GameplayType_FashionSend               GameplayType = 25 //时装赠送
	GameplayType_HeroSoulSend              GameplayType = 26 //八门赠送
	GameplayType_Team                      GameplayType = 27 //组队
	GameplayType_LeaderboardStage          GameplayType = 28 //赛事冲榜
	GameplayType_LeaderboardTower          GameplayType = 29 //吉斯塔冲榜
	GameplayType_AllianceTeamLeague        GameplayType = 30 //团赛
	GameplayType_PatrolShare               GameplayType = 31 //探秘分享
	GameplayType_RugalRecruit              GameplayType = 32 //卢卡尔招募
	GameplayType_ManorCheer                GameplayType = 33 //领地助威
	GameplayType_TeamArena                 GameplayType = 34 //组队竞技场
	GameplayType_ArcadeMidLevel            GameplayType = 35 //街机厅中级
	GameplayType_ArcadeHighLevel           GameplayType = 36 //街机厅高级
	GameplayType_ArcadeRank                GameplayType = 37 //街机厅排行
	GameplayType_DasheBackRecruit          GameplayType = 38 //大蛇返场招募
	GameplayType_ServerAllGift             GameplayType = 39 //全服赠送
	GameplayType_Crazy8Recruit             GameplayType = 40 //疯八招募
	GameplayType_RugalBackRecruit          GameplayType = 41 //卢卡尔返场招募
	GameplayType_ArcadeLevel8              GameplayType = 42 //街机厅8服
	GameplayType_ArcadeLevel12             GameplayType = 43 //街机厅12服
	GameplayType_ArcadeLevel24             GameplayType = 44 //街机厅24服
	GameplayType_ArcadeLevel48             GameplayType = 45 //街机厅2战区
	GameplayType_ArcadeLevel72             GameplayType = 46 //街机厅3战区
	GameplayType_KingWarWeek               GameplayType = 47 //皇城争锋周赛
	GameplayType_KingWarFinal              GameplayType = 48 //皇城争锋决赛
	GameplayType_KingWarMap                GameplayType = 49 //皇城争锋地图赛
	GameplayType_Crazy8BackRecruit         GameplayType = 50 //疯八返场招募
	GameplayType_CrazyLeonaRecruit         GameplayType = 51 //狂暴莉安娜招募
	GameplayType_AllianceWar               GameplayType = 52 //联盟攻城战
	GameplayType_CrazyLeonaBackRecruit     GameplayType = 53 //狂暴莉安娜返场招募
	GameplayType_END                       GameplayType = 54 //枚举结束
)

// 跨服配置-拳皇/serverGameplayType
type PlatformType int

const (
	PlatformType_BEGIN   PlatformType = 0 //枚举开始
	PlatformType_Unknown PlatformType = 0 //未知平台
	PlatformType_Wechat  PlatformType = 1 //微信
	PlatformType_Douyin  PlatformType = 2 //抖音
	PlatformType_Alipay  PlatformType = 3 //支付宝
	PlatformType_Meituan PlatformType = 4 //美团
	PlatformType_END     PlatformType = 5 //枚举结束
)

// 跨服配置-拳皇/serverGameplayType
type CrossType int

const (
	CrossType_BEGIN                   CrossType = 0  //枚举开始
	CrossType_CrossWechat             CrossType = 0  //微信跨服方案
	CrossType_CrossDouyin             CrossType = 1  //抖音跨服方案
	CrossType_CrossWechatGroupBefor4  CrossType = 2  //跨服加公会测试方案
	CrossType_CrossWechatGroupBefor5  CrossType = 3  //跨服对战激进方案
	CrossType_CrossWechatGroupBefor6  CrossType = 4  //微小3战区特殊处理方案
	CrossType_CrossWechatGroupBefor7  CrossType = 5  //微小4战区特殊处理方案
	CrossType_CrossWechatGroupBefor8  CrossType = 6  //跨服加公会测试方案4服
	CrossType_CrossWechatGroupBefor9  CrossType = 7  //微小11战区特殊处理方案
	CrossType_CrossWechatGroupBefor10 CrossType = 8  //抖小14战区特殊处理方案
	CrossType_CrossWechatGroupBefor11 CrossType = 9  //微小12战区特殊处理方案
	CrossType_CrossWechatGroupBefor12 CrossType = 10 //抖小15战区特殊处理方案
	CrossType_CrossWechatGroupBefor13 CrossType = 11 //微小13战区特殊处理方案
	CrossType_CrossDouyinGroupBefor16 CrossType = 12 //抖小16战区特殊处理方案
	CrossType_END                     CrossType = 13 //枚举结束
)

// 跨服配置-拳皇/serverGameplayType
type GameplayCrossDivideType int

const (
	GameplayCrossDivideType_BEGIN GameplayCrossDivideType = 0 //枚举开始
	GameplayCrossDivideType_Snake GameplayCrossDivideType = 0 //积分蛇形分布
	GameplayCrossDivideType_Order GameplayCrossDivideType = 1 //积分顺序分布
	GameplayCrossDivideType_END   GameplayCrossDivideType = 2 //枚举结束
)

func (e GameplayType) String() string {

	if e == GameplayType_ManorGrab {
		return "ManorGrab"
	}
	if e == GameplayType_Arena {
		return "Arena"
	}
	if e == GameplayType_Bounty {
		return "Bounty"
	}
	if e == GameplayType_CrossArena {
		return "CrossArena"
	}
	if e == GameplayType_CampBattle {
		return "CampBattle"
	}
	if e == GameplayType_CityWar {
		return "CityWar"
	}
	if e == GameplayType_AllianceDuel {
		return "AllianceDuel"
	}
	if e == GameplayType_Alliance {
		return "Alliance"
	}
	if e == GameplayType_WorldBoss {
		return "WorldBoss"
	}
	if e == GameplayType_KofRoad {
		return "KofRoad"
	}
	if e == GameplayType_ELFTower {
		return "ELFTower"
	}
	if e == GameplayType_Fighting {
		return "Fighting"
	}
	if e == GameplayType_DailyTopic {
		return "DailyTopic"
	}
	if e == GameplayType_LeaderboardGameplay {
		return "LeaderboardGameplay"
	}
	if e == GameplayType_LeaderboardArena {
		return "LeaderboardArena"
	}
	if e == GameplayType_LeaderboardActiveAlliance {
		return "LeaderboardActiveAlliance"
	}
	if e == GameplayType_UpNurture {
		return "UpNurture"
	}
	if e == GameplayType_SendShare {
		return "SendShare"
	}
	if e == GameplayType_Strongest108 {
		return "Strongest108"
	}
	if e == GameplayType_LeaderboardRecruit {
		return "LeaderboardRecruit"
	}
	if e == GameplayType_LeaderboardArtifact {
		return "LeaderboardArtifact"
	}
	if e == GameplayType_LeaderboardMysticism {
		return "LeaderboardMysticism"
	}
	if e == GameplayType_LeaderboardBattlePower {
		return "LeaderboardBattlePower"
	}
	if e == GameplayType_Arcade {
		return "Arcade"
	}
	if e == GameplayType_FashionSend {
		return "FashionSend"
	}
	if e == GameplayType_HeroSoulSend {
		return "HeroSoulSend"
	}
	if e == GameplayType_Team {
		return "Team"
	}
	if e == GameplayType_LeaderboardStage {
		return "LeaderboardStage"
	}
	if e == GameplayType_LeaderboardTower {
		return "LeaderboardTower"
	}
	if e == GameplayType_AllianceTeamLeague {
		return "AllianceTeamLeague"
	}
	if e == GameplayType_PatrolShare {
		return "PatrolShare"
	}
	if e == GameplayType_RugalRecruit {
		return "RugalRecruit"
	}
	if e == GameplayType_ManorCheer {
		return "ManorCheer"
	}
	if e == GameplayType_TeamArena {
		return "TeamArena"
	}
	if e == GameplayType_ArcadeMidLevel {
		return "ArcadeMidLevel"
	}
	if e == GameplayType_ArcadeHighLevel {
		return "ArcadeHighLevel"
	}
	if e == GameplayType_ArcadeRank {
		return "ArcadeRank"
	}
	if e == GameplayType_DasheBackRecruit {
		return "DasheBackRecruit"
	}
	if e == GameplayType_ServerAllGift {
		return "ServerAllGift"
	}
	if e == GameplayType_Crazy8Recruit {
		return "Crazy8Recruit"
	}
	if e == GameplayType_RugalBackRecruit {
		return "RugalBackRecruit"
	}
	if e == GameplayType_ArcadeLevel8 {
		return "ArcadeLevel8"
	}
	if e == GameplayType_ArcadeLevel12 {
		return "ArcadeLevel12"
	}
	if e == GameplayType_ArcadeLevel24 {
		return "ArcadeLevel24"
	}
	if e == GameplayType_ArcadeLevel48 {
		return "ArcadeLevel48"
	}
	if e == GameplayType_ArcadeLevel72 {
		return "ArcadeLevel72"
	}
	if e == GameplayType_KingWarWeek {
		return "KingWarWeek"
	}
	if e == GameplayType_KingWarFinal {
		return "KingWarFinal"
	}
	if e == GameplayType_KingWarMap {
		return "KingWarMap"
	}
	if e == GameplayType_Crazy8BackRecruit {
		return "Crazy8BackRecruit"
	}
	if e == GameplayType_CrazyLeonaRecruit {
		return "CrazyLeonaRecruit"
	}
	if e == GameplayType_AllianceWar {
		return "AllianceWar"
	}
	if e == GameplayType_CrazyLeonaBackRecruit {
		return "CrazyLeonaBackRecruit"
	}

	return ""

}

func GetGameplayType(tp string) (GameplayType, bool) {
	switch tp {
	case "ManorGrab":
		return GameplayType_ManorGrab, true
	case "Arena":
		return GameplayType_Arena, true
	case "Bounty":
		return GameplayType_Bounty, true
	case "CrossArena":
		return GameplayType_CrossArena, true
	case "CampBattle":
		return GameplayType_CampBattle, true
	case "CityWar":
		return GameplayType_CityWar, true
	case "AllianceDuel":
		return GameplayType_AllianceDuel, true
	case "Alliance":
		return GameplayType_Alliance, true
	case "WorldBoss":
		return GameplayType_WorldBoss, true
	case "KofRoad":
		return GameplayType_KofRoad, true
	case "ELFTower":
		return GameplayType_ELFTower, true
	case "Fighting":
		return GameplayType_Fighting, true
	case "DailyTopic":
		return GameplayType_DailyTopic, true
	case "LeaderboardGameplay":
		return GameplayType_LeaderboardGameplay, true
	case "LeaderboardArena":
		return GameplayType_LeaderboardArena, true
	case "LeaderboardActiveAlliance":
		return GameplayType_LeaderboardActiveAlliance, true
	case "UpNurture":
		return GameplayType_UpNurture, true
	case "SendShare":
		return GameplayType_SendShare, true
	case "Strongest108":
		return GameplayType_Strongest108, true
	case "LeaderboardRecruit":
		return GameplayType_LeaderboardRecruit, true
	case "LeaderboardArtifact":
		return GameplayType_LeaderboardArtifact, true
	case "LeaderboardMysticism":
		return GameplayType_LeaderboardMysticism, true
	case "LeaderboardBattlePower":
		return GameplayType_LeaderboardBattlePower, true
	case "Arcade":
		return GameplayType_Arcade, true
	case "FashionSend":
		return GameplayType_FashionSend, true
	case "HeroSoulSend":
		return GameplayType_HeroSoulSend, true
	case "Team":
		return GameplayType_Team, true
	case "LeaderboardStage":
		return GameplayType_LeaderboardStage, true
	case "LeaderboardTower":
		return GameplayType_LeaderboardTower, true
	case "AllianceTeamLeague":
		return GameplayType_AllianceTeamLeague, true
	case "PatrolShare":
		return GameplayType_PatrolShare, true
	case "RugalRecruit":
		return GameplayType_RugalRecruit, true
	case "ManorCheer":
		return GameplayType_ManorCheer, true
	case "TeamArena":
		return GameplayType_TeamArena, true
	case "ArcadeMidLevel":
		return GameplayType_ArcadeMidLevel, true
	case "ArcadeHighLevel":
		return GameplayType_ArcadeHighLevel, true
	case "ArcadeRank":
		return GameplayType_ArcadeRank, true
	case "DasheBackRecruit":
		return GameplayType_DasheBackRecruit, true
	case "ServerAllGift":
		return GameplayType_ServerAllGift, true
	case "Crazy8Recruit":
		return GameplayType_Crazy8Recruit, true
	case "RugalBackRecruit":
		return GameplayType_RugalBackRecruit, true
	case "ArcadeLevel8":
		return GameplayType_ArcadeLevel8, true
	case "ArcadeLevel12":
		return GameplayType_ArcadeLevel12, true
	case "ArcadeLevel24":
		return GameplayType_ArcadeLevel24, true
	case "ArcadeLevel48":
		return GameplayType_ArcadeLevel48, true
	case "ArcadeLevel72":
		return GameplayType_ArcadeLevel72, true
	case "KingWarWeek":
		return GameplayType_KingWarWeek, true
	case "KingWarFinal":
		return GameplayType_KingWarFinal, true
	case "KingWarMap":
		return GameplayType_KingWarMap, true
	case "Crazy8BackRecruit":
		return GameplayType_Crazy8BackRecruit, true
	case "CrazyLeonaRecruit":
		return GameplayType_CrazyLeonaRecruit, true
	case "AllianceWar":
		return GameplayType_AllianceWar, true
	case "CrazyLeonaBackRecruit":
		return GameplayType_CrazyLeonaBackRecruit, true

	}
	return -1, false

}

func GetGameplayType_ByAlias(tp string) (GameplayType, bool) {
	switch tp {
	case "领地掠夺":
		return GameplayType_ManorGrab, true
	case "竞技场":
		return GameplayType_Arena, true
	case "赏金战":
		return GameplayType_Bounty, true
	case "跨服竞技32强":
		return GameplayType_CrossArena, true
	case "阵营战":
		return GameplayType_CampBattle, true
	case "城市争夺战":
		return GameplayType_CityWar, true
	case "公会对决":
		return GameplayType_AllianceDuel, true
	case "公会":
		return GameplayType_Alliance, true
	case "世界boss":
		return GameplayType_WorldBoss, true
	case "拳皇之路":
		return GameplayType_KofRoad, true
	case "塔防精灵":
		return GameplayType_ELFTower, true
	case "格斗圣殿":
		return GameplayType_Fighting, true
	case "每日话题":
		return GameplayType_DailyTopic, true
	case "新榜玩法":
		return GameplayType_LeaderboardGameplay, true
	case "新榜竞技场":
		return GameplayType_LeaderboardArena, true
	case "新榜活跃公会":
		return GameplayType_LeaderboardActiveAlliance, true
	case "Up养成":
		return GameplayType_UpNurture, true
	case "跨服频道":
		return GameplayType_SendShare, true
	case "最强108":
		return GameplayType_Strongest108, true
	case "招募冲榜":
		return GameplayType_LeaderboardRecruit, true
	case "神器冲榜":
		return GameplayType_LeaderboardArtifact, true
	case "秘术冲榜":
		return GameplayType_LeaderboardMysticism, true
	case "战力冲榜":
		return GameplayType_LeaderboardBattlePower, true
	case "街机厅":
		return GameplayType_Arcade, true
	case "时装赠送":
		return GameplayType_FashionSend, true
	case "八门赠送":
		return GameplayType_HeroSoulSend, true
	case "组队":
		return GameplayType_Team, true
	case "赛事冲榜":
		return GameplayType_LeaderboardStage, true
	case "吉斯塔冲榜":
		return GameplayType_LeaderboardTower, true
	case "团赛":
		return GameplayType_AllianceTeamLeague, true
	case "探秘分享":
		return GameplayType_PatrolShare, true
	case "卢卡尔招募":
		return GameplayType_RugalRecruit, true
	case "领地助威":
		return GameplayType_ManorCheer, true
	case "组队竞技场":
		return GameplayType_TeamArena, true
	case "街机厅中级":
		return GameplayType_ArcadeMidLevel, true
	case "街机厅高级":
		return GameplayType_ArcadeHighLevel, true
	case "街机厅排行":
		return GameplayType_ArcadeRank, true
	case "大蛇返场招募":
		return GameplayType_DasheBackRecruit, true
	case "全服赠送":
		return GameplayType_ServerAllGift, true
	case "疯八招募":
		return GameplayType_Crazy8Recruit, true
	case "卢卡尔返场招募":
		return GameplayType_RugalBackRecruit, true
	case "街机厅8服":
		return GameplayType_ArcadeLevel8, true
	case "街机厅12服":
		return GameplayType_ArcadeLevel12, true
	case "街机厅24服":
		return GameplayType_ArcadeLevel24, true
	case "街机厅2战区":
		return GameplayType_ArcadeLevel48, true
	case "街机厅3战区":
		return GameplayType_ArcadeLevel72, true
	case "皇城争锋周赛":
		return GameplayType_KingWarWeek, true
	case "皇城争锋决赛":
		return GameplayType_KingWarFinal, true
	case "皇城争锋地图赛":
		return GameplayType_KingWarMap, true
	case "疯八返场招募":
		return GameplayType_Crazy8BackRecruit, true
	case "狂暴莉安娜招募":
		return GameplayType_CrazyLeonaRecruit, true
	case "联盟攻城战":
		return GameplayType_AllianceWar, true
	case "狂暴莉安娜返场招募":
		return GameplayType_CrazyLeonaBackRecruit, true

	}
	return -1, false

}

func (e PlatformType) String() string {

	if e == PlatformType_Unknown {
		return "Unknown"
	}
	if e == PlatformType_Wechat {
		return "Wechat"
	}
	if e == PlatformType_Douyin {
		return "Douyin"
	}
	if e == PlatformType_Alipay {
		return "Alipay"
	}
	if e == PlatformType_Meituan {
		return "Meituan"
	}

	return ""

}

func GetPlatformType(tp string) (PlatformType, bool) {
	switch tp {
	case "Unknown":
		return PlatformType_Unknown, true
	case "Wechat":
		return PlatformType_Wechat, true
	case "Douyin":
		return PlatformType_Douyin, true
	case "Alipay":
		return PlatformType_Alipay, true
	case "Meituan":
		return PlatformType_Meituan, true

	}
	return -1, false

}

func GetPlatformType_ByAlias(tp string) (PlatformType, bool) {
	switch tp {
	case "未知平台":
		return PlatformType_Unknown, true
	case "微信":
		return PlatformType_Wechat, true
	case "抖音":
		return PlatformType_Douyin, true
	case "支付宝":
		return PlatformType_Alipay, true
	case "美团":
		return PlatformType_Meituan, true

	}
	return -1, false

}

func (e CrossType) String() string {

	if e == CrossType_CrossWechat {
		return "CrossWechat"
	}
	if e == CrossType_CrossDouyin {
		return "CrossDouyin"
	}
	if e == CrossType_CrossWechatGroupBefor4 {
		return "CrossWechatGroupBefor4"
	}
	if e == CrossType_CrossWechatGroupBefor5 {
		return "CrossWechatGroupBefor5"
	}
	if e == CrossType_CrossWechatGroupBefor6 {
		return "CrossWechatGroupBefor6"
	}
	if e == CrossType_CrossWechatGroupBefor7 {
		return "CrossWechatGroupBefor7"
	}
	if e == CrossType_CrossWechatGroupBefor8 {
		return "CrossWechatGroupBefor8"
	}
	if e == CrossType_CrossWechatGroupBefor9 {
		return "CrossWechatGroupBefor9"
	}
	if e == CrossType_CrossWechatGroupBefor10 {
		return "CrossWechatGroupBefor10"
	}
	if e == CrossType_CrossWechatGroupBefor11 {
		return "CrossWechatGroupBefor11"
	}
	if e == CrossType_CrossWechatGroupBefor12 {
		return "CrossWechatGroupBefor12"
	}
	if e == CrossType_CrossWechatGroupBefor13 {
		return "CrossWechatGroupBefor13"
	}
	if e == CrossType_CrossDouyinGroupBefor16 {
		return "CrossDouyinGroupBefor16"
	}

	return ""

}

func GetCrossType(tp string) (CrossType, bool) {
	switch tp {
	case "CrossWechat":
		return CrossType_CrossWechat, true
	case "CrossDouyin":
		return CrossType_CrossDouyin, true
	case "CrossWechatGroupBefor4":
		return CrossType_CrossWechatGroupBefor4, true
	case "CrossWechatGroupBefor5":
		return CrossType_CrossWechatGroupBefor5, true
	case "CrossWechatGroupBefor6":
		return CrossType_CrossWechatGroupBefor6, true
	case "CrossWechatGroupBefor7":
		return CrossType_CrossWechatGroupBefor7, true
	case "CrossWechatGroupBefor8":
		return CrossType_CrossWechatGroupBefor8, true
	case "CrossWechatGroupBefor9":
		return CrossType_CrossWechatGroupBefor9, true
	case "CrossWechatGroupBefor10":
		return CrossType_CrossWechatGroupBefor10, true
	case "CrossWechatGroupBefor11":
		return CrossType_CrossWechatGroupBefor11, true
	case "CrossWechatGroupBefor12":
		return CrossType_CrossWechatGroupBefor12, true
	case "CrossWechatGroupBefor13":
		return CrossType_CrossWechatGroupBefor13, true
	case "CrossDouyinGroupBefor16":
		return CrossType_CrossDouyinGroupBefor16, true

	}
	return -1, false

}

func GetCrossType_ByAlias(tp string) (CrossType, bool) {
	switch tp {
	case "微信跨服方案":
		return CrossType_CrossWechat, true
	case "抖音跨服方案":
		return CrossType_CrossDouyin, true
	case "跨服加公会测试方案":
		return CrossType_CrossWechatGroupBefor4, true
	case "跨服对战激进方案":
		return CrossType_CrossWechatGroupBefor5, true
	case "微小3战区特殊处理方案":
		return CrossType_CrossWechatGroupBefor6, true
	case "微小4战区特殊处理方案":
		return CrossType_CrossWechatGroupBefor7, true
	case "跨服加公会测试方案4服":
		return CrossType_CrossWechatGroupBefor8, true
	case "微小11战区特殊处理方案":
		return CrossType_CrossWechatGroupBefor9, true
	case "抖小14战区特殊处理方案":
		return CrossType_CrossWechatGroupBefor10, true
	case "微小12战区特殊处理方案":
		return CrossType_CrossWechatGroupBefor11, true
	case "抖小15战区特殊处理方案":
		return CrossType_CrossWechatGroupBefor12, true
	case "微小13战区特殊处理方案":
		return CrossType_CrossWechatGroupBefor13, true
	case "抖小16战区特殊处理方案":
		return CrossType_CrossDouyinGroupBefor16, true

	}
	return -1, false

}

func (e GameplayCrossDivideType) String() string {

	if e == GameplayCrossDivideType_Snake {
		return "Snake"
	}
	if e == GameplayCrossDivideType_Order {
		return "Order"
	}

	return ""

}

func GetGameplayCrossDivideType(tp string) (GameplayCrossDivideType, bool) {
	switch tp {
	case "Snake":
		return GameplayCrossDivideType_Snake, true
	case "Order":
		return GameplayCrossDivideType_Order, true

	}
	return -1, false

}

func GetGameplayCrossDivideType_ByAlias(tp string) (GameplayCrossDivideType, bool) {
	switch tp {
	case "积分蛇形分布":
		return GameplayCrossDivideType_Snake, true
	case "积分顺序分布":
		return GameplayCrossDivideType_Order, true

	}
	return -1, false

}
