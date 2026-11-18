package config

// 玩家表-拳皇/playerEvent
type PlayerEventType int

const (
	PlayerEventType_BEGIN                    PlayerEventType = 1   //枚举开始
	PlayerEventType_LevelUp                  PlayerEventType = 1   //玩家升级
	PlayerEventType_StagePass                PlayerEventType = 2   //关卡通过
	PlayerEventType_ArenaChallenge           PlayerEventType = 3   //竞技场挑战
	PlayerEventType_GainItem                 PlayerEventType = 4   //获得道具
	PlayerEventType_WorldBossChallenge       PlayerEventType = 5   //世界boss挑战
	PlayerEventType_MainQuestComplete        PlayerEventType = 6   //完成主线任务
	PlayerEventType_CostDiamond              PlayerEventType = 7   //消耗钻石
	PlayerEventType_WatchAds                 PlayerEventType = 8   //看广告
	PlayerEventType_CostItem                 PlayerEventType = 9   //消耗道具
	PlayerEventType_Recharge                 PlayerEventType = 10  //充值
	PlayerEventType_ShopBuy                  PlayerEventType = 11  //商店购买
	PlayerEventType_CorpsPromotionLevelup    PlayerEventType = 12  //佣兵团晋升
	PlayerEventType_HeroGain                 PlayerEventType = 13  //获得英雄
	PlayerEventType_HeroStar                 PlayerEventType = 14  //英雄升星
	PlayerEventType_HeroLineup               PlayerEventType = 15  //英雄上阵
	PlayerEventType_HeroLinedown             PlayerEventType = 16  //英雄下阵
	PlayerEventType_HeroDraw                 PlayerEventType = 17  //英雄招募
	PlayerEventType_TopFivePowerChange       PlayerEventType = 18  //前五战力变动
	PlayerEventType_HeroLevelUp              PlayerEventType = 19  //英雄升级
	PlayerEventType_TowerFloorPass           PlayerEventType = 20  //爬塔通关
	PlayerEventType_ArenaWin                 PlayerEventType = 21  //竞技场胜利
	PlayerEventType_HeroEquipLevelUp         PlayerEventType = 22  //英雄装备升级
	PlayerEventType_MainChallenge            PlayerEventType = 23  //主线关卡挑战
	PlayerEventType_QuikAfkReward            PlayerEventType = 24  //领取快速离线奖励
	PlayerEventType_CommonAfkReward          PlayerEventType = 25  //领取离线奖励
	PlayerEventType_CommonHeroRecruit        PlayerEventType = 26  //普通英雄招募
	PlayerEventType_SeniorHeroRecruit        PlayerEventType = 27  //高级高级招募
	PlayerEventType_TrialChallenge           PlayerEventType = 28  //试炼挑战
	PlayerEventType_TowerChallenge           PlayerEventType = 29  //爬塔挑战
	PlayerEventType_EightGateCombine         PlayerEventType = 30  //八门合成
	PlayerEventType_ArtifactDraw             PlayerEventType = 31  //神器召唤
	PlayerEventType_HeroSoulRefining         PlayerEventType = 32  //战魂洗练
	PlayerEventType_ArtifactActivate         PlayerEventType = 33  //神器激活
	PlayerEventType_ArtifactBreak            PlayerEventType = 34  //神器突破
	PlayerEventType_ArtifactUpStar           PlayerEventType = 35  //神器升星
	PlayerEventType_GodheadUpLv              PlayerEventType = 36  //神格升级
	PlayerEventType_GodheadReset             PlayerEventType = 37  //神格重置
	PlayerEventType_StatueUpLv               PlayerEventType = 38  //雕像升级
	PlayerEventType_StatusBlessUpLv          PlayerEventType = 39  //雕像祝福升级
	PlayerEventType_FameUpLv                 PlayerEventType = 40  //名人堂（翅膀）升级
	PlayerEventType_HeroReset                PlayerEventType = 41  //英雄重生
	PlayerEventType_ShopCostCurrency         PlayerEventType = 42  //商店消耗货币
	PlayerEventType_LeaveAlliance            PlayerEventType = 43  //离开工会
	PlayerEventType_MonthCardPoint           PlayerEventType = 44  //月卡积分改变
	PlayerEventType_PrivilegeCardActive      PlayerEventType = 45  //权益卡生效
	PlayerEventType_PrivilegeCardDeactive    PlayerEventType = 46  //权益卡失效
	PlayerEventType_ManorTechLevelup         PlayerEventType = 47  //领地科技升级
	PlayerEventType_FishSuccess              PlayerEventType = 48  //钓鱼成功
	PlayerEventType_FishSell                 PlayerEventType = 49  //出售鱼获
	PlayerEventType_FishEquipStrength        PlayerEventType = 50  //强化渔具
	PlayerEventType_FishOrderCommit          PlayerEventType = 51  //提交钓鱼订单
	PlayerEventType_FishOrderAward           PlayerEventType = 52  //完成钓鱼订单
	PlayerEventType_VipRightsChange          PlayerEventType = 53  //贵族权益改变
	PlayerEventType_AllianceBossChallenge    PlayerEventType = 54  //公会boss挑战
	PlayerEventType_FarwarriorWin            PlayerEventType = 55  //远征胜利
	PlayerEventType_FarwarriorPass           PlayerEventType = 56  //远征通关
	PlayerEventType_BoxFavorUpLv             PlayerEventType = 57  //好感度升级
	PlayerEventType_RechargeDay              PlayerEventType = 58  //充值天
	PlayerEventType_LDStagePass              PlayerEventType = 59  //LD关卡通关
	PlayerEventType_LDSummonHero             PlayerEventType = 60  //LD召唤英雄
	PlayerEventType_LDGambleSuccess          PlayerEventType = 61  //LD赌博成功
	PlayerEventType_LDGambleFail             PlayerEventType = 62  //LD赌博失败
	PlayerEventType_LDHuntSuccess            PlayerEventType = 63  //LD狩猎成功
	PlayerEventType_LDRoomMonsterAdd         PlayerEventType = 64  //LD房间怪物增加
	PlayerEventType_LDHeroEnchantSuccess     PlayerEventType = 65  //LD英雄强化成功
	PlayerEventType_LDHeroEnchantFail        PlayerEventType = 66  //LD英雄强化失败
	PlayerEventType_LDHeroKillMonster        PlayerEventType = 67  //LD英雄击杀怪物
	PlayerEventType_LDCostItem               PlayerEventType = 68  //LD消耗道具
	PlayerEventType_LDMergeLegendHero        PlayerEventType = 69  //LD合成神话英雄
	PlayerEventType_LDHeroCardLevelup        PlayerEventType = 70  //LD英雄卡牌培养升级
	PlayerEventType_LDHeroSell               PlayerEventType = 71  //LD英雄出售
	PlayerEventType_AllianceBuild            PlayerEventType = 72  //联盟建设
	PlayerEventType_AllianceTechGacha        PlayerEventType = 73  //联盟科技抽卡
	PlayerEventType_HeroPromotion            PlayerEventType = 74  //英雄晋升
	PlayerEventType_DayChange                PlayerEventType = 75  //跨天
	PlayerEventType_AddAvatar                PlayerEventType = 76  //获得头像
	PlayerEventType_AddAvatarFrame           PlayerEventType = 77  //获得头像框
	PlayerEventType_AddTitle                 PlayerEventType = 78  //获得称号
	PlayerEventType_AddChatBubble            PlayerEventType = 79  //获得聊天气泡
	PlayerEventType_GiftBoxLevelup           PlayerEventType = 80  //格斗家礼物箱子升级
	PlayerEventType_NestChallengeProgress    PlayerEventType = 81  //音巢挑战进度改变
	PlayerEventType_AllianceTechGachaChip    PlayerEventType = 82  //联盟科技抽到碎片
	PlayerEventType_AllianceTechUpgrade      PlayerEventType = 83  //联盟科技升级
	PlayerEventType_FetterActive             PlayerEventType = 84  //羁绊激活
	PlayerEventType_ManorLevelup             PlayerEventType = 85  //领地升级
	PlayerEventType_HeroAwaken               PlayerEventType = 86  //英雄觉醒
	PlayerEventType_AllianceTechWear         PlayerEventType = 87  //联盟科技穿戴
	PlayerEventType_BoxGiftOpenBox           PlayerEventType = 88  //格斗家礼物开箱
	PlayerEventType_VipExpChange             PlayerEventType = 89  //Vip经验改变
	PlayerEventType_ManorGotRes              PlayerEventType = 90  //领地获得资源
	PlayerEventType_EquipHeroSoul            PlayerEventType = 91  //装备战魂
	PlayerEventType_AllianceBargain          PlayerEventType = 92  //公会砍价
	PlayerEventType_QuestGainAward           PlayerEventType = 93  //任务领奖
	PlayerEventType_ManorPlunder             PlayerEventType = 94  //领地掠夺
	PlayerEventType_GetPurchaseAward         PlayerEventType = 95  //领取充值道具
	PlayerEventType_TotalPowerChange         PlayerEventType = 96  //总战力变动
	PlayerEventType_HeroReplace              PlayerEventType = 97  //英雄互换
	PlayerEventType_KofRoadLvUp              PlayerEventType = 98  //拳皇路晋升
	PlayerEventType_JoinCityWar              PlayerEventType = 99  //参与城池战
	PlayerEventType_CityWarOccupyCity        PlayerEventType = 100 //占领城池战城池
	PlayerEventType_CrossArenaRank           PlayerEventType = 101 //跨服32强排名
	PlayerEventType_BountyPersonRank         PlayerEventType = 102 //赏金战个人排名
	PlayerEventType_BountyAllianceRank       PlayerEventType = 103 //赏金战公会排名
	PlayerEventType_RecruitLevelup           PlayerEventType = 104 //招募升级
	PlayerEventType_FightingRank             PlayerEventType = 105 //格斗圣殿升排名
	PlayerEventType_FameForge                PlayerEventType = 106 //名人堂铸造
	PlayerEventType_GainCampBattleToken      PlayerEventType = 107 //获得阵营战商店代币
	PlayerEventType_NewbiePlayerBindSuccess  PlayerEventType = 108 //邀请玩家绑定成功
	PlayerEventType_NewbiePlayerLevelUp      PlayerEventType = 109 //邀请玩家升级
	PlayerEventType_ContinuousScoreChange    PlayerEventType = 110 //连续充值活动积分变化
	PlayerEventType_BindDestinyHero          PlayerEventType = 111 //绑定宿命英雄
	PlayerEventType_GainPrivilegeCard        PlayerEventType = 112 //激活特权卡
	PlayerEventType_OpenNewSystem            PlayerEventType = 113 //开启新系统
	PlayerEventType_BoxGiftReplace           PlayerEventType = 114 //格斗家礼物替换
	PlayerEventType_ElfBattleResult          PlayerEventType = 115 //南镇战斗结果
	PlayerEventType_WorldBossDamage          PlayerEventType = 116 //世界boss伤害
	PlayerEventType_ChannelChatTimes         PlayerEventType = 117 //频道发言
	PlayerEventType_AnyOneOfShopBuy          PlayerEventType = 118 //多个商店购买
	PlayerEventType_HeroChipCost             PlayerEventType = 119 //英雄碎片消耗
	PlayerEventType_HistoryBattlePowerChange PlayerEventType = 120 //历史最高战力变化
	PlayerEventType_ActPuzzleActive          PlayerEventType = 121 //拼图点亮格子
	PlayerEventType_WishHeroRecruit          PlayerEventType = 122 //自选招募
	PlayerEventType_UpHeroRecruit            PlayerEventType = 123 //Up招募
	PlayerEventType_TeamStagePass            PlayerEventType = 124 //组队关卡通关
	PlayerEventType_HeroDepromotion          PlayerEventType = 125 //英雄降品
	PlayerEventType_RealRecharge             PlayerEventType = 126 //真实充值
	PlayerEventType_ActivityRelicChallenge   PlayerEventType = 127 //活动遗迹关卡挑战
	PlayerEventType_GameClubFollow           PlayerEventType = 128 //关注游戏圈
	PlayerEventType_GameClubBrowsingSeconds  PlayerEventType = 129 //游戏圈浏览时长
	PlayerEventType_GameClubLikeCount        PlayerEventType = 130 //游戏圈内点赞帖子
	PlayerEventType_GameClubCommentCount     PlayerEventType = 131 //评论官方帖子
	PlayerEventType_GameClubPost             PlayerEventType = 132 //游戏圈发帖
	PlayerEventType_TradeJoin                PlayerEventType = 133 //参与贸易
	PlayerEventType_TradeHoldAsset           PlayerEventType = 134 //贸易持有资产
	PlayerEventType_TradeEarn                PlayerEventType = 135 //贸易赚取
	PlayerEventType_TradeRank                PlayerEventType = 136 //贸易名次
	PlayerEventType_AvyTaskGotAward          PlayerEventType = 137 //活动任务领奖
	PlayerEventType_OpenFashionBox           PlayerEventType = 138 //开启时装
	PlayerEventType_KingWarPlayerScoreChange PlayerEventType = 139 //皇城个人积分改变
	PlayerEventType_KingWarPlayerCostEnergy  PlayerEventType = 140 //皇城累计消耗精力
	PlayerEventType_AllianceDissolve         PlayerEventType = 141 //联盟解散
	PlayerEventType_AllianceJoin             PlayerEventType = 142 //加入联盟
	PlayerEventType_AniActIslandDraw         PlayerEventType = 143 //周年海岛探险
	PlayerEventType_PetMerge                 PlayerEventType = 144 //宠物合成
	PlayerEventType_END                      PlayerEventType = 145 //枚举结束
)

// 玩家表-拳皇/playerEvent
type AvatarFrameTimeType int

const (
	AvatarFrameTimeType_BEGIN      AvatarFrameTimeType = 1 //枚举开始
	AvatarFrameTimeType_Season     AvatarFrameTimeType = 1 //赛季
	AvatarFrameTimeType_ConfigTime AvatarFrameTimeType = 2 //配置时间
	AvatarFrameTimeType_Forever    AvatarFrameTimeType = 3 //永久
	AvatarFrameTimeType_END        AvatarFrameTimeType = 4 //枚举结束
)

// 玩家表-拳皇/playerEvent
type AvatarGroupType int

const (
	AvatarGroupType_BEGIN   AvatarGroupType = 1 //枚举开始
	AvatarGroupType_Fashion AvatarGroupType = 1 //装扮头像
	AvatarGroupType_Hero    AvatarGroupType = 2 //格斗家头像
	AvatarGroupType_END     AvatarGroupType = 3 //枚举结束
)

// 玩家表-拳皇/playerEvent
type AvatarFrameGroupType int

const (
	AvatarFrameGroupType_BEGIN   AvatarFrameGroupType = 1 //枚举开始
	AvatarFrameGroupType_Fashion AvatarFrameGroupType = 1 //装扮头像框
	AvatarFrameGroupType_Hero    AvatarFrameGroupType = 2 //头像框
	AvatarFrameGroupType_END     AvatarFrameGroupType = 3 //枚举结束
)

func (e PlayerEventType) String() string {

	if e == PlayerEventType_LevelUp {
		return "LevelUp"
	}
	if e == PlayerEventType_StagePass {
		return "StagePass"
	}
	if e == PlayerEventType_ArenaChallenge {
		return "ArenaChallenge"
	}
	if e == PlayerEventType_GainItem {
		return "GainItem"
	}
	if e == PlayerEventType_WorldBossChallenge {
		return "WorldBossChallenge"
	}
	if e == PlayerEventType_MainQuestComplete {
		return "MainQuestComplete"
	}
	if e == PlayerEventType_CostDiamond {
		return "CostDiamond"
	}
	if e == PlayerEventType_WatchAds {
		return "WatchAds"
	}
	if e == PlayerEventType_CostItem {
		return "CostItem"
	}
	if e == PlayerEventType_Recharge {
		return "Recharge"
	}
	if e == PlayerEventType_ShopBuy {
		return "ShopBuy"
	}
	if e == PlayerEventType_CorpsPromotionLevelup {
		return "CorpsPromotionLevelup"
	}
	if e == PlayerEventType_HeroGain {
		return "HeroGain"
	}
	if e == PlayerEventType_HeroStar {
		return "HeroStar"
	}
	if e == PlayerEventType_HeroLineup {
		return "HeroLineup"
	}
	if e == PlayerEventType_HeroLinedown {
		return "HeroLinedown"
	}
	if e == PlayerEventType_HeroDraw {
		return "HeroDraw"
	}
	if e == PlayerEventType_TopFivePowerChange {
		return "TopFivePowerChange"
	}
	if e == PlayerEventType_HeroLevelUp {
		return "HeroLevelUp"
	}
	if e == PlayerEventType_TowerFloorPass {
		return "TowerFloorPass"
	}
	if e == PlayerEventType_ArenaWin {
		return "ArenaWin"
	}
	if e == PlayerEventType_HeroEquipLevelUp {
		return "HeroEquipLevelUp"
	}
	if e == PlayerEventType_MainChallenge {
		return "MainChallenge"
	}
	if e == PlayerEventType_QuikAfkReward {
		return "QuikAfkReward"
	}
	if e == PlayerEventType_CommonAfkReward {
		return "CommonAfkReward"
	}
	if e == PlayerEventType_CommonHeroRecruit {
		return "CommonHeroRecruit"
	}
	if e == PlayerEventType_SeniorHeroRecruit {
		return "SeniorHeroRecruit"
	}
	if e == PlayerEventType_TrialChallenge {
		return "TrialChallenge"
	}
	if e == PlayerEventType_TowerChallenge {
		return "TowerChallenge"
	}
	if e == PlayerEventType_EightGateCombine {
		return "EightGateCombine"
	}
	if e == PlayerEventType_ArtifactDraw {
		return "ArtifactDraw"
	}
	if e == PlayerEventType_HeroSoulRefining {
		return "HeroSoulRefining"
	}
	if e == PlayerEventType_ArtifactActivate {
		return "ArtifactActivate"
	}
	if e == PlayerEventType_ArtifactBreak {
		return "ArtifactBreak"
	}
	if e == PlayerEventType_ArtifactUpStar {
		return "ArtifactUpStar"
	}
	if e == PlayerEventType_GodheadUpLv {
		return "GodheadUpLv"
	}
	if e == PlayerEventType_GodheadReset {
		return "GodheadReset"
	}
	if e == PlayerEventType_StatueUpLv {
		return "StatueUpLv"
	}
	if e == PlayerEventType_StatusBlessUpLv {
		return "StatusBlessUpLv"
	}
	if e == PlayerEventType_FameUpLv {
		return "FameUpLv"
	}
	if e == PlayerEventType_HeroReset {
		return "HeroReset"
	}
	if e == PlayerEventType_ShopCostCurrency {
		return "ShopCostCurrency"
	}
	if e == PlayerEventType_LeaveAlliance {
		return "LeaveAlliance"
	}
	if e == PlayerEventType_MonthCardPoint {
		return "MonthCardPoint"
	}
	if e == PlayerEventType_PrivilegeCardActive {
		return "PrivilegeCardActive"
	}
	if e == PlayerEventType_PrivilegeCardDeactive {
		return "PrivilegeCardDeactive"
	}
	if e == PlayerEventType_ManorTechLevelup {
		return "ManorTechLevelup"
	}
	if e == PlayerEventType_FishSuccess {
		return "FishSuccess"
	}
	if e == PlayerEventType_FishSell {
		return "FishSell"
	}
	if e == PlayerEventType_FishEquipStrength {
		return "FishEquipStrength"
	}
	if e == PlayerEventType_FishOrderCommit {
		return "FishOrderCommit"
	}
	if e == PlayerEventType_FishOrderAward {
		return "FishOrderAward"
	}
	if e == PlayerEventType_VipRightsChange {
		return "VipRightsChange"
	}
	if e == PlayerEventType_AllianceBossChallenge {
		return "AllianceBossChallenge"
	}
	if e == PlayerEventType_FarwarriorWin {
		return "FarwarriorWin"
	}
	if e == PlayerEventType_FarwarriorPass {
		return "FarwarriorPass"
	}
	if e == PlayerEventType_BoxFavorUpLv {
		return "BoxFavorUpLv"
	}
	if e == PlayerEventType_RechargeDay {
		return "RechargeDay"
	}
	if e == PlayerEventType_LDStagePass {
		return "LDStagePass"
	}
	if e == PlayerEventType_LDSummonHero {
		return "LDSummonHero"
	}
	if e == PlayerEventType_LDGambleSuccess {
		return "LDGambleSuccess"
	}
	if e == PlayerEventType_LDGambleFail {
		return "LDGambleFail"
	}
	if e == PlayerEventType_LDHuntSuccess {
		return "LDHuntSuccess"
	}
	if e == PlayerEventType_LDRoomMonsterAdd {
		return "LDRoomMonsterAdd"
	}
	if e == PlayerEventType_LDHeroEnchantSuccess {
		return "LDHeroEnchantSuccess"
	}
	if e == PlayerEventType_LDHeroEnchantFail {
		return "LDHeroEnchantFail"
	}
	if e == PlayerEventType_LDHeroKillMonster {
		return "LDHeroKillMonster"
	}
	if e == PlayerEventType_LDCostItem {
		return "LDCostItem"
	}
	if e == PlayerEventType_LDMergeLegendHero {
		return "LDMergeLegendHero"
	}
	if e == PlayerEventType_LDHeroCardLevelup {
		return "LDHeroCardLevelup"
	}
	if e == PlayerEventType_LDHeroSell {
		return "LDHeroSell"
	}
	if e == PlayerEventType_AllianceBuild {
		return "AllianceBuild"
	}
	if e == PlayerEventType_AllianceTechGacha {
		return "AllianceTechGacha"
	}
	if e == PlayerEventType_HeroPromotion {
		return "HeroPromotion"
	}
	if e == PlayerEventType_DayChange {
		return "DayChange"
	}
	if e == PlayerEventType_AddAvatar {
		return "AddAvatar"
	}
	if e == PlayerEventType_AddAvatarFrame {
		return "AddAvatarFrame"
	}
	if e == PlayerEventType_AddTitle {
		return "AddTitle"
	}
	if e == PlayerEventType_AddChatBubble {
		return "AddChatBubble"
	}
	if e == PlayerEventType_GiftBoxLevelup {
		return "GiftBoxLevelup"
	}
	if e == PlayerEventType_NestChallengeProgress {
		return "NestChallengeProgress"
	}
	if e == PlayerEventType_AllianceTechGachaChip {
		return "AllianceTechGachaChip"
	}
	if e == PlayerEventType_AllianceTechUpgrade {
		return "AllianceTechUpgrade"
	}
	if e == PlayerEventType_FetterActive {
		return "FetterActive"
	}
	if e == PlayerEventType_ManorLevelup {
		return "ManorLevelup"
	}
	if e == PlayerEventType_HeroAwaken {
		return "HeroAwaken"
	}
	if e == PlayerEventType_AllianceTechWear {
		return "AllianceTechWear"
	}
	if e == PlayerEventType_BoxGiftOpenBox {
		return "BoxGiftOpenBox"
	}
	if e == PlayerEventType_VipExpChange {
		return "VipExpChange"
	}
	if e == PlayerEventType_ManorGotRes {
		return "ManorGotRes"
	}
	if e == PlayerEventType_EquipHeroSoul {
		return "EquipHeroSoul"
	}
	if e == PlayerEventType_AllianceBargain {
		return "AllianceBargain"
	}
	if e == PlayerEventType_QuestGainAward {
		return "QuestGainAward"
	}
	if e == PlayerEventType_ManorPlunder {
		return "ManorPlunder"
	}
	if e == PlayerEventType_GetPurchaseAward {
		return "GetPurchaseAward"
	}
	if e == PlayerEventType_TotalPowerChange {
		return "TotalPowerChange"
	}
	if e == PlayerEventType_HeroReplace {
		return "HeroReplace"
	}
	if e == PlayerEventType_KofRoadLvUp {
		return "KofRoadLvUp"
	}
	if e == PlayerEventType_JoinCityWar {
		return "JoinCityWar"
	}
	if e == PlayerEventType_CityWarOccupyCity {
		return "CityWarOccupyCity"
	}
	if e == PlayerEventType_CrossArenaRank {
		return "CrossArenaRank"
	}
	if e == PlayerEventType_BountyPersonRank {
		return "BountyPersonRank"
	}
	if e == PlayerEventType_BountyAllianceRank {
		return "BountyAllianceRank"
	}
	if e == PlayerEventType_RecruitLevelup {
		return "RecruitLevelup"
	}
	if e == PlayerEventType_FightingRank {
		return "FightingRank"
	}
	if e == PlayerEventType_FameForge {
		return "FameForge"
	}
	if e == PlayerEventType_GainCampBattleToken {
		return "GainCampBattleToken"
	}
	if e == PlayerEventType_NewbiePlayerBindSuccess {
		return "NewbiePlayerBindSuccess"
	}
	if e == PlayerEventType_NewbiePlayerLevelUp {
		return "NewbiePlayerLevelUp"
	}
	if e == PlayerEventType_ContinuousScoreChange {
		return "ContinuousScoreChange"
	}
	if e == PlayerEventType_BindDestinyHero {
		return "BindDestinyHero"
	}
	if e == PlayerEventType_GainPrivilegeCard {
		return "GainPrivilegeCard"
	}
	if e == PlayerEventType_OpenNewSystem {
		return "OpenNewSystem"
	}
	if e == PlayerEventType_BoxGiftReplace {
		return "BoxGiftReplace"
	}
	if e == PlayerEventType_ElfBattleResult {
		return "ElfBattleResult"
	}
	if e == PlayerEventType_WorldBossDamage {
		return "WorldBossDamage"
	}
	if e == PlayerEventType_ChannelChatTimes {
		return "ChannelChatTimes"
	}
	if e == PlayerEventType_AnyOneOfShopBuy {
		return "AnyOneOfShopBuy"
	}
	if e == PlayerEventType_HeroChipCost {
		return "HeroChipCost"
	}
	if e == PlayerEventType_HistoryBattlePowerChange {
		return "HistoryBattlePowerChange"
	}
	if e == PlayerEventType_ActPuzzleActive {
		return "ActPuzzleActive"
	}
	if e == PlayerEventType_WishHeroRecruit {
		return "WishHeroRecruit"
	}
	if e == PlayerEventType_UpHeroRecruit {
		return "UpHeroRecruit"
	}
	if e == PlayerEventType_TeamStagePass {
		return "TeamStagePass"
	}
	if e == PlayerEventType_HeroDepromotion {
		return "HeroDepromotion"
	}
	if e == PlayerEventType_RealRecharge {
		return "RealRecharge"
	}
	if e == PlayerEventType_ActivityRelicChallenge {
		return "ActivityRelicChallenge"
	}
	if e == PlayerEventType_GameClubFollow {
		return "GameClubFollow"
	}
	if e == PlayerEventType_GameClubBrowsingSeconds {
		return "GameClubBrowsingSeconds"
	}
	if e == PlayerEventType_GameClubLikeCount {
		return "GameClubLikeCount"
	}
	if e == PlayerEventType_GameClubCommentCount {
		return "GameClubCommentCount"
	}
	if e == PlayerEventType_GameClubPost {
		return "GameClubPost"
	}
	if e == PlayerEventType_TradeJoin {
		return "TradeJoin"
	}
	if e == PlayerEventType_TradeHoldAsset {
		return "TradeHoldAsset"
	}
	if e == PlayerEventType_TradeEarn {
		return "TradeEarn"
	}
	if e == PlayerEventType_TradeRank {
		return "TradeRank"
	}
	if e == PlayerEventType_AvyTaskGotAward {
		return "AvyTaskGotAward"
	}
	if e == PlayerEventType_OpenFashionBox {
		return "OpenFashionBox"
	}
	if e == PlayerEventType_KingWarPlayerScoreChange {
		return "KingWarPlayerScoreChange"
	}
	if e == PlayerEventType_KingWarPlayerCostEnergy {
		return "KingWarPlayerCostEnergy"
	}
	if e == PlayerEventType_AllianceDissolve {
		return "AllianceDissolve"
	}
	if e == PlayerEventType_AllianceJoin {
		return "AllianceJoin"
	}
	if e == PlayerEventType_AniActIslandDraw {
		return "AniActIslandDraw"
	}
	if e == PlayerEventType_PetMerge {
		return "PetMerge"
	}

	return ""

}

func GetPlayerEventType(tp string) (PlayerEventType, bool) {
	switch tp {
	case "LevelUp":
		return PlayerEventType_LevelUp, true
	case "StagePass":
		return PlayerEventType_StagePass, true
	case "ArenaChallenge":
		return PlayerEventType_ArenaChallenge, true
	case "GainItem":
		return PlayerEventType_GainItem, true
	case "WorldBossChallenge":
		return PlayerEventType_WorldBossChallenge, true
	case "MainQuestComplete":
		return PlayerEventType_MainQuestComplete, true
	case "CostDiamond":
		return PlayerEventType_CostDiamond, true
	case "WatchAds":
		return PlayerEventType_WatchAds, true
	case "CostItem":
		return PlayerEventType_CostItem, true
	case "Recharge":
		return PlayerEventType_Recharge, true
	case "ShopBuy":
		return PlayerEventType_ShopBuy, true
	case "CorpsPromotionLevelup":
		return PlayerEventType_CorpsPromotionLevelup, true
	case "HeroGain":
		return PlayerEventType_HeroGain, true
	case "HeroStar":
		return PlayerEventType_HeroStar, true
	case "HeroLineup":
		return PlayerEventType_HeroLineup, true
	case "HeroLinedown":
		return PlayerEventType_HeroLinedown, true
	case "HeroDraw":
		return PlayerEventType_HeroDraw, true
	case "TopFivePowerChange":
		return PlayerEventType_TopFivePowerChange, true
	case "HeroLevelUp":
		return PlayerEventType_HeroLevelUp, true
	case "TowerFloorPass":
		return PlayerEventType_TowerFloorPass, true
	case "ArenaWin":
		return PlayerEventType_ArenaWin, true
	case "HeroEquipLevelUp":
		return PlayerEventType_HeroEquipLevelUp, true
	case "MainChallenge":
		return PlayerEventType_MainChallenge, true
	case "QuikAfkReward":
		return PlayerEventType_QuikAfkReward, true
	case "CommonAfkReward":
		return PlayerEventType_CommonAfkReward, true
	case "CommonHeroRecruit":
		return PlayerEventType_CommonHeroRecruit, true
	case "SeniorHeroRecruit":
		return PlayerEventType_SeniorHeroRecruit, true
	case "TrialChallenge":
		return PlayerEventType_TrialChallenge, true
	case "TowerChallenge":
		return PlayerEventType_TowerChallenge, true
	case "EightGateCombine":
		return PlayerEventType_EightGateCombine, true
	case "ArtifactDraw":
		return PlayerEventType_ArtifactDraw, true
	case "HeroSoulRefining":
		return PlayerEventType_HeroSoulRefining, true
	case "ArtifactActivate":
		return PlayerEventType_ArtifactActivate, true
	case "ArtifactBreak":
		return PlayerEventType_ArtifactBreak, true
	case "ArtifactUpStar":
		return PlayerEventType_ArtifactUpStar, true
	case "GodheadUpLv":
		return PlayerEventType_GodheadUpLv, true
	case "GodheadReset":
		return PlayerEventType_GodheadReset, true
	case "StatueUpLv":
		return PlayerEventType_StatueUpLv, true
	case "StatusBlessUpLv":
		return PlayerEventType_StatusBlessUpLv, true
	case "FameUpLv":
		return PlayerEventType_FameUpLv, true
	case "HeroReset":
		return PlayerEventType_HeroReset, true
	case "ShopCostCurrency":
		return PlayerEventType_ShopCostCurrency, true
	case "LeaveAlliance":
		return PlayerEventType_LeaveAlliance, true
	case "MonthCardPoint":
		return PlayerEventType_MonthCardPoint, true
	case "PrivilegeCardActive":
		return PlayerEventType_PrivilegeCardActive, true
	case "PrivilegeCardDeactive":
		return PlayerEventType_PrivilegeCardDeactive, true
	case "ManorTechLevelup":
		return PlayerEventType_ManorTechLevelup, true
	case "FishSuccess":
		return PlayerEventType_FishSuccess, true
	case "FishSell":
		return PlayerEventType_FishSell, true
	case "FishEquipStrength":
		return PlayerEventType_FishEquipStrength, true
	case "FishOrderCommit":
		return PlayerEventType_FishOrderCommit, true
	case "FishOrderAward":
		return PlayerEventType_FishOrderAward, true
	case "VipRightsChange":
		return PlayerEventType_VipRightsChange, true
	case "AllianceBossChallenge":
		return PlayerEventType_AllianceBossChallenge, true
	case "FarwarriorWin":
		return PlayerEventType_FarwarriorWin, true
	case "FarwarriorPass":
		return PlayerEventType_FarwarriorPass, true
	case "BoxFavorUpLv":
		return PlayerEventType_BoxFavorUpLv, true
	case "RechargeDay":
		return PlayerEventType_RechargeDay, true
	case "LDStagePass":
		return PlayerEventType_LDStagePass, true
	case "LDSummonHero":
		return PlayerEventType_LDSummonHero, true
	case "LDGambleSuccess":
		return PlayerEventType_LDGambleSuccess, true
	case "LDGambleFail":
		return PlayerEventType_LDGambleFail, true
	case "LDHuntSuccess":
		return PlayerEventType_LDHuntSuccess, true
	case "LDRoomMonsterAdd":
		return PlayerEventType_LDRoomMonsterAdd, true
	case "LDHeroEnchantSuccess":
		return PlayerEventType_LDHeroEnchantSuccess, true
	case "LDHeroEnchantFail":
		return PlayerEventType_LDHeroEnchantFail, true
	case "LDHeroKillMonster":
		return PlayerEventType_LDHeroKillMonster, true
	case "LDCostItem":
		return PlayerEventType_LDCostItem, true
	case "LDMergeLegendHero":
		return PlayerEventType_LDMergeLegendHero, true
	case "LDHeroCardLevelup":
		return PlayerEventType_LDHeroCardLevelup, true
	case "LDHeroSell":
		return PlayerEventType_LDHeroSell, true
	case "AllianceBuild":
		return PlayerEventType_AllianceBuild, true
	case "AllianceTechGacha":
		return PlayerEventType_AllianceTechGacha, true
	case "HeroPromotion":
		return PlayerEventType_HeroPromotion, true
	case "DayChange":
		return PlayerEventType_DayChange, true
	case "AddAvatar":
		return PlayerEventType_AddAvatar, true
	case "AddAvatarFrame":
		return PlayerEventType_AddAvatarFrame, true
	case "AddTitle":
		return PlayerEventType_AddTitle, true
	case "AddChatBubble":
		return PlayerEventType_AddChatBubble, true
	case "GiftBoxLevelup":
		return PlayerEventType_GiftBoxLevelup, true
	case "NestChallengeProgress":
		return PlayerEventType_NestChallengeProgress, true
	case "AllianceTechGachaChip":
		return PlayerEventType_AllianceTechGachaChip, true
	case "AllianceTechUpgrade":
		return PlayerEventType_AllianceTechUpgrade, true
	case "FetterActive":
		return PlayerEventType_FetterActive, true
	case "ManorLevelup":
		return PlayerEventType_ManorLevelup, true
	case "HeroAwaken":
		return PlayerEventType_HeroAwaken, true
	case "AllianceTechWear":
		return PlayerEventType_AllianceTechWear, true
	case "BoxGiftOpenBox":
		return PlayerEventType_BoxGiftOpenBox, true
	case "VipExpChange":
		return PlayerEventType_VipExpChange, true
	case "ManorGotRes":
		return PlayerEventType_ManorGotRes, true
	case "EquipHeroSoul":
		return PlayerEventType_EquipHeroSoul, true
	case "AllianceBargain":
		return PlayerEventType_AllianceBargain, true
	case "QuestGainAward":
		return PlayerEventType_QuestGainAward, true
	case "ManorPlunder":
		return PlayerEventType_ManorPlunder, true
	case "GetPurchaseAward":
		return PlayerEventType_GetPurchaseAward, true
	case "TotalPowerChange":
		return PlayerEventType_TotalPowerChange, true
	case "HeroReplace":
		return PlayerEventType_HeroReplace, true
	case "KofRoadLvUp":
		return PlayerEventType_KofRoadLvUp, true
	case "JoinCityWar":
		return PlayerEventType_JoinCityWar, true
	case "CityWarOccupyCity":
		return PlayerEventType_CityWarOccupyCity, true
	case "CrossArenaRank":
		return PlayerEventType_CrossArenaRank, true
	case "BountyPersonRank":
		return PlayerEventType_BountyPersonRank, true
	case "BountyAllianceRank":
		return PlayerEventType_BountyAllianceRank, true
	case "RecruitLevelup":
		return PlayerEventType_RecruitLevelup, true
	case "FightingRank":
		return PlayerEventType_FightingRank, true
	case "FameForge":
		return PlayerEventType_FameForge, true
	case "GainCampBattleToken":
		return PlayerEventType_GainCampBattleToken, true
	case "NewbiePlayerBindSuccess":
		return PlayerEventType_NewbiePlayerBindSuccess, true
	case "NewbiePlayerLevelUp":
		return PlayerEventType_NewbiePlayerLevelUp, true
	case "ContinuousScoreChange":
		return PlayerEventType_ContinuousScoreChange, true
	case "BindDestinyHero":
		return PlayerEventType_BindDestinyHero, true
	case "GainPrivilegeCard":
		return PlayerEventType_GainPrivilegeCard, true
	case "OpenNewSystem":
		return PlayerEventType_OpenNewSystem, true
	case "BoxGiftReplace":
		return PlayerEventType_BoxGiftReplace, true
	case "ElfBattleResult":
		return PlayerEventType_ElfBattleResult, true
	case "WorldBossDamage":
		return PlayerEventType_WorldBossDamage, true
	case "ChannelChatTimes":
		return PlayerEventType_ChannelChatTimes, true
	case "AnyOneOfShopBuy":
		return PlayerEventType_AnyOneOfShopBuy, true
	case "HeroChipCost":
		return PlayerEventType_HeroChipCost, true
	case "HistoryBattlePowerChange":
		return PlayerEventType_HistoryBattlePowerChange, true
	case "ActPuzzleActive":
		return PlayerEventType_ActPuzzleActive, true
	case "WishHeroRecruit":
		return PlayerEventType_WishHeroRecruit, true
	case "UpHeroRecruit":
		return PlayerEventType_UpHeroRecruit, true
	case "TeamStagePass":
		return PlayerEventType_TeamStagePass, true
	case "HeroDepromotion":
		return PlayerEventType_HeroDepromotion, true
	case "RealRecharge":
		return PlayerEventType_RealRecharge, true
	case "ActivityRelicChallenge":
		return PlayerEventType_ActivityRelicChallenge, true
	case "GameClubFollow":
		return PlayerEventType_GameClubFollow, true
	case "GameClubBrowsingSeconds":
		return PlayerEventType_GameClubBrowsingSeconds, true
	case "GameClubLikeCount":
		return PlayerEventType_GameClubLikeCount, true
	case "GameClubCommentCount":
		return PlayerEventType_GameClubCommentCount, true
	case "GameClubPost":
		return PlayerEventType_GameClubPost, true
	case "TradeJoin":
		return PlayerEventType_TradeJoin, true
	case "TradeHoldAsset":
		return PlayerEventType_TradeHoldAsset, true
	case "TradeEarn":
		return PlayerEventType_TradeEarn, true
	case "TradeRank":
		return PlayerEventType_TradeRank, true
	case "AvyTaskGotAward":
		return PlayerEventType_AvyTaskGotAward, true
	case "OpenFashionBox":
		return PlayerEventType_OpenFashionBox, true
	case "KingWarPlayerScoreChange":
		return PlayerEventType_KingWarPlayerScoreChange, true
	case "KingWarPlayerCostEnergy":
		return PlayerEventType_KingWarPlayerCostEnergy, true
	case "AllianceDissolve":
		return PlayerEventType_AllianceDissolve, true
	case "AllianceJoin":
		return PlayerEventType_AllianceJoin, true
	case "AniActIslandDraw":
		return PlayerEventType_AniActIslandDraw, true
	case "PetMerge":
		return PlayerEventType_PetMerge, true

	}
	return -1, false

}

func GetPlayerEventType_ByAlias(tp string) (PlayerEventType, bool) {
	switch tp {
	case "玩家升级":
		return PlayerEventType_LevelUp, true
	case "关卡通过":
		return PlayerEventType_StagePass, true
	case "竞技场挑战":
		return PlayerEventType_ArenaChallenge, true
	case "获得道具":
		return PlayerEventType_GainItem, true
	case "世界boss挑战":
		return PlayerEventType_WorldBossChallenge, true
	case "完成主线任务":
		return PlayerEventType_MainQuestComplete, true
	case "消耗钻石":
		return PlayerEventType_CostDiamond, true
	case "看广告":
		return PlayerEventType_WatchAds, true
	case "消耗道具":
		return PlayerEventType_CostItem, true
	case "充值":
		return PlayerEventType_Recharge, true
	case "商店购买":
		return PlayerEventType_ShopBuy, true
	case "佣兵团晋升":
		return PlayerEventType_CorpsPromotionLevelup, true
	case "获得英雄":
		return PlayerEventType_HeroGain, true
	case "英雄升星":
		return PlayerEventType_HeroStar, true
	case "英雄上阵":
		return PlayerEventType_HeroLineup, true
	case "英雄下阵":
		return PlayerEventType_HeroLinedown, true
	case "英雄招募":
		return PlayerEventType_HeroDraw, true
	case "前五战力变动":
		return PlayerEventType_TopFivePowerChange, true
	case "英雄升级":
		return PlayerEventType_HeroLevelUp, true
	case "爬塔通关":
		return PlayerEventType_TowerFloorPass, true
	case "竞技场胜利":
		return PlayerEventType_ArenaWin, true
	case "英雄装备升级":
		return PlayerEventType_HeroEquipLevelUp, true
	case "主线关卡挑战":
		return PlayerEventType_MainChallenge, true
	case "领取快速离线奖励":
		return PlayerEventType_QuikAfkReward, true
	case "领取离线奖励":
		return PlayerEventType_CommonAfkReward, true
	case "普通英雄招募":
		return PlayerEventType_CommonHeroRecruit, true
	case "高级高级招募":
		return PlayerEventType_SeniorHeroRecruit, true
	case "试炼挑战":
		return PlayerEventType_TrialChallenge, true
	case "爬塔挑战":
		return PlayerEventType_TowerChallenge, true
	case "八门合成":
		return PlayerEventType_EightGateCombine, true
	case "神器召唤":
		return PlayerEventType_ArtifactDraw, true
	case "战魂洗练":
		return PlayerEventType_HeroSoulRefining, true
	case "神器激活":
		return PlayerEventType_ArtifactActivate, true
	case "神器突破":
		return PlayerEventType_ArtifactBreak, true
	case "神器升星":
		return PlayerEventType_ArtifactUpStar, true
	case "神格升级":
		return PlayerEventType_GodheadUpLv, true
	case "神格重置":
		return PlayerEventType_GodheadReset, true
	case "雕像升级":
		return PlayerEventType_StatueUpLv, true
	case "雕像祝福升级":
		return PlayerEventType_StatusBlessUpLv, true
	case "名人堂（翅膀）升级":
		return PlayerEventType_FameUpLv, true
	case "英雄重生":
		return PlayerEventType_HeroReset, true
	case "商店消耗货币":
		return PlayerEventType_ShopCostCurrency, true
	case "离开工会":
		return PlayerEventType_LeaveAlliance, true
	case "月卡积分改变":
		return PlayerEventType_MonthCardPoint, true
	case "权益卡生效":
		return PlayerEventType_PrivilegeCardActive, true
	case "权益卡失效":
		return PlayerEventType_PrivilegeCardDeactive, true
	case "领地科技升级":
		return PlayerEventType_ManorTechLevelup, true
	case "钓鱼成功":
		return PlayerEventType_FishSuccess, true
	case "出售鱼获":
		return PlayerEventType_FishSell, true
	case "强化渔具":
		return PlayerEventType_FishEquipStrength, true
	case "提交钓鱼订单":
		return PlayerEventType_FishOrderCommit, true
	case "完成钓鱼订单":
		return PlayerEventType_FishOrderAward, true
	case "贵族权益改变":
		return PlayerEventType_VipRightsChange, true
	case "公会boss挑战":
		return PlayerEventType_AllianceBossChallenge, true
	case "远征胜利":
		return PlayerEventType_FarwarriorWin, true
	case "远征通关":
		return PlayerEventType_FarwarriorPass, true
	case "好感度升级":
		return PlayerEventType_BoxFavorUpLv, true
	case "充值天":
		return PlayerEventType_RechargeDay, true
	case "LD关卡通关":
		return PlayerEventType_LDStagePass, true
	case "LD召唤英雄":
		return PlayerEventType_LDSummonHero, true
	case "LD赌博成功":
		return PlayerEventType_LDGambleSuccess, true
	case "LD赌博失败":
		return PlayerEventType_LDGambleFail, true
	case "LD狩猎成功":
		return PlayerEventType_LDHuntSuccess, true
	case "LD房间怪物增加":
		return PlayerEventType_LDRoomMonsterAdd, true
	case "LD英雄强化成功":
		return PlayerEventType_LDHeroEnchantSuccess, true
	case "LD英雄强化失败":
		return PlayerEventType_LDHeroEnchantFail, true
	case "LD英雄击杀怪物":
		return PlayerEventType_LDHeroKillMonster, true
	case "LD消耗道具":
		return PlayerEventType_LDCostItem, true
	case "LD合成神话英雄":
		return PlayerEventType_LDMergeLegendHero, true
	case "LD英雄卡牌培养升级":
		return PlayerEventType_LDHeroCardLevelup, true
	case "LD英雄出售":
		return PlayerEventType_LDHeroSell, true
	case "联盟建设":
		return PlayerEventType_AllianceBuild, true
	case "联盟科技抽卡":
		return PlayerEventType_AllianceTechGacha, true
	case "英雄晋升":
		return PlayerEventType_HeroPromotion, true
	case "跨天":
		return PlayerEventType_DayChange, true
	case "获得头像":
		return PlayerEventType_AddAvatar, true
	case "获得头像框":
		return PlayerEventType_AddAvatarFrame, true
	case "获得称号":
		return PlayerEventType_AddTitle, true
	case "获得聊天气泡":
		return PlayerEventType_AddChatBubble, true
	case "格斗家礼物箱子升级":
		return PlayerEventType_GiftBoxLevelup, true
	case "音巢挑战进度改变":
		return PlayerEventType_NestChallengeProgress, true
	case "联盟科技抽到碎片":
		return PlayerEventType_AllianceTechGachaChip, true
	case "联盟科技升级":
		return PlayerEventType_AllianceTechUpgrade, true
	case "羁绊激活":
		return PlayerEventType_FetterActive, true
	case "领地升级":
		return PlayerEventType_ManorLevelup, true
	case "英雄觉醒":
		return PlayerEventType_HeroAwaken, true
	case "联盟科技穿戴":
		return PlayerEventType_AllianceTechWear, true
	case "格斗家礼物开箱":
		return PlayerEventType_BoxGiftOpenBox, true
	case "Vip经验改变":
		return PlayerEventType_VipExpChange, true
	case "领地获得资源":
		return PlayerEventType_ManorGotRes, true
	case "装备战魂":
		return PlayerEventType_EquipHeroSoul, true
	case "公会砍价":
		return PlayerEventType_AllianceBargain, true
	case "任务领奖":
		return PlayerEventType_QuestGainAward, true
	case "领地掠夺":
		return PlayerEventType_ManorPlunder, true
	case "领取充值道具":
		return PlayerEventType_GetPurchaseAward, true
	case "总战力变动":
		return PlayerEventType_TotalPowerChange, true
	case "英雄互换":
		return PlayerEventType_HeroReplace, true
	case "拳皇路晋升":
		return PlayerEventType_KofRoadLvUp, true
	case "参与城池战":
		return PlayerEventType_JoinCityWar, true
	case "占领城池战城池":
		return PlayerEventType_CityWarOccupyCity, true
	case "跨服32强排名":
		return PlayerEventType_CrossArenaRank, true
	case "赏金战个人排名":
		return PlayerEventType_BountyPersonRank, true
	case "赏金战公会排名":
		return PlayerEventType_BountyAllianceRank, true
	case "招募升级":
		return PlayerEventType_RecruitLevelup, true
	case "格斗圣殿升排名":
		return PlayerEventType_FightingRank, true
	case "名人堂铸造":
		return PlayerEventType_FameForge, true
	case "获得阵营战商店代币":
		return PlayerEventType_GainCampBattleToken, true
	case "邀请玩家绑定成功":
		return PlayerEventType_NewbiePlayerBindSuccess, true
	case "邀请玩家升级":
		return PlayerEventType_NewbiePlayerLevelUp, true
	case "连续充值活动积分变化":
		return PlayerEventType_ContinuousScoreChange, true
	case "绑定宿命英雄":
		return PlayerEventType_BindDestinyHero, true
	case "激活特权卡":
		return PlayerEventType_GainPrivilegeCard, true
	case "开启新系统":
		return PlayerEventType_OpenNewSystem, true
	case "格斗家礼物替换":
		return PlayerEventType_BoxGiftReplace, true
	case "南镇战斗结果":
		return PlayerEventType_ElfBattleResult, true
	case "世界boss伤害":
		return PlayerEventType_WorldBossDamage, true
	case "频道发言":
		return PlayerEventType_ChannelChatTimes, true
	case "多个商店购买":
		return PlayerEventType_AnyOneOfShopBuy, true
	case "英雄碎片消耗":
		return PlayerEventType_HeroChipCost, true
	case "历史最高战力变化":
		return PlayerEventType_HistoryBattlePowerChange, true
	case "拼图点亮格子":
		return PlayerEventType_ActPuzzleActive, true
	case "自选招募":
		return PlayerEventType_WishHeroRecruit, true
	case "Up招募":
		return PlayerEventType_UpHeroRecruit, true
	case "组队关卡通关":
		return PlayerEventType_TeamStagePass, true
	case "英雄降品":
		return PlayerEventType_HeroDepromotion, true
	case "真实充值":
		return PlayerEventType_RealRecharge, true
	case "活动遗迹关卡挑战":
		return PlayerEventType_ActivityRelicChallenge, true
	case "关注游戏圈":
		return PlayerEventType_GameClubFollow, true
	case "游戏圈浏览时长":
		return PlayerEventType_GameClubBrowsingSeconds, true
	case "游戏圈内点赞帖子":
		return PlayerEventType_GameClubLikeCount, true
	case "评论官方帖子":
		return PlayerEventType_GameClubCommentCount, true
	case "游戏圈发帖":
		return PlayerEventType_GameClubPost, true
	case "参与贸易":
		return PlayerEventType_TradeJoin, true
	case "贸易持有资产":
		return PlayerEventType_TradeHoldAsset, true
	case "贸易赚取":
		return PlayerEventType_TradeEarn, true
	case "贸易名次":
		return PlayerEventType_TradeRank, true
	case "活动任务领奖":
		return PlayerEventType_AvyTaskGotAward, true
	case "开启时装":
		return PlayerEventType_OpenFashionBox, true
	case "皇城个人积分改变":
		return PlayerEventType_KingWarPlayerScoreChange, true
	case "皇城累计消耗精力":
		return PlayerEventType_KingWarPlayerCostEnergy, true
	case "联盟解散":
		return PlayerEventType_AllianceDissolve, true
	case "加入联盟":
		return PlayerEventType_AllianceJoin, true
	case "周年海岛探险":
		return PlayerEventType_AniActIslandDraw, true
	case "宠物合成":
		return PlayerEventType_PetMerge, true

	}
	return -1, false

}

func (e AvatarFrameTimeType) String() string {

	if e == AvatarFrameTimeType_Season {
		return "Season"
	}
	if e == AvatarFrameTimeType_ConfigTime {
		return "ConfigTime"
	}
	if e == AvatarFrameTimeType_Forever {
		return "Forever"
	}

	return ""

}

func GetAvatarFrameTimeType(tp string) (AvatarFrameTimeType, bool) {
	switch tp {
	case "Season":
		return AvatarFrameTimeType_Season, true
	case "ConfigTime":
		return AvatarFrameTimeType_ConfigTime, true
	case "Forever":
		return AvatarFrameTimeType_Forever, true

	}
	return -1, false

}

func GetAvatarFrameTimeType_ByAlias(tp string) (AvatarFrameTimeType, bool) {
	switch tp {
	case "赛季":
		return AvatarFrameTimeType_Season, true
	case "配置时间":
		return AvatarFrameTimeType_ConfigTime, true
	case "永久":
		return AvatarFrameTimeType_Forever, true

	}
	return -1, false

}

func (e AvatarGroupType) String() string {

	if e == AvatarGroupType_Fashion {
		return "Fashion"
	}
	if e == AvatarGroupType_Hero {
		return "Hero"
	}

	return ""

}

func GetAvatarGroupType(tp string) (AvatarGroupType, bool) {
	switch tp {
	case "Fashion":
		return AvatarGroupType_Fashion, true
	case "Hero":
		return AvatarGroupType_Hero, true

	}
	return -1, false

}

func GetAvatarGroupType_ByAlias(tp string) (AvatarGroupType, bool) {
	switch tp {
	case "装扮头像":
		return AvatarGroupType_Fashion, true
	case "格斗家头像":
		return AvatarGroupType_Hero, true

	}
	return -1, false

}

func (e AvatarFrameGroupType) String() string {

	if e == AvatarFrameGroupType_Fashion {
		return "Fashion"
	}
	if e == AvatarFrameGroupType_Hero {
		return "Hero"
	}

	return ""

}

func GetAvatarFrameGroupType(tp string) (AvatarFrameGroupType, bool) {
	switch tp {
	case "Fashion":
		return AvatarFrameGroupType_Fashion, true
	case "Hero":
		return AvatarFrameGroupType_Hero, true

	}
	return -1, false

}

func GetAvatarFrameGroupType_ByAlias(tp string) (AvatarFrameGroupType, bool) {
	switch tp {
	case "装扮头像框":
		return AvatarFrameGroupType_Fashion, true
	case "头像框":
		return AvatarFrameGroupType_Hero, true

	}
	return -1, false

}
