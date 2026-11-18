package config

// 任务表-拳皇/ConditionType
type SubConditionOperationType int

const (
	SubConditionOperationType_BEGIN SubConditionOperationType = 0 //枚举开始
	SubConditionOperationType_Equal SubConditionOperationType = 0 //等于
	SubConditionOperationType_Reach SubConditionOperationType = 1 //达到
	SubConditionOperationType_Less  SubConditionOperationType = 2 //小于
	SubConditionOperationType_Lte   SubConditionOperationType = 3 //小于等于
	SubConditionOperationType_END   SubConditionOperationType = 4 //枚举结束
)

// 任务表-拳皇/ConditionType
type CounterType int

const (
	CounterType_BEGIN         CounterType = 0 //枚举开始
	CounterType_Trigger       CounterType = 0 //触发计数
	CounterType_Accumulate    CounterType = 1 //累积计数
	CounterType_RepeatTrigger CounterType = 2 //重复触发
	CounterType_END           CounterType = 3 //枚举结束
)

// 任务表-拳皇/ConditionType
type SubConditionRelation int

const (
	SubConditionRelation_BEGIN SubConditionRelation = 0 //枚举开始
	SubConditionRelation_And   SubConditionRelation = 0 //满足全部子条件
	SubConditionRelation_Or    SubConditionRelation = 1 //满足一个子条件
	SubConditionRelation_END   SubConditionRelation = 2 //枚举结束
)

// 任务表-拳皇/ConditionType
type SubConditionType int

const (
	SubConditionType_BEGIN                           SubConditionType = 0   //枚举开始
	SubConditionType_Empty                           SubConditionType = 0   //空
	SubConditionType_PlayerLevel                     SubConditionType = 1   //玩家等级
	SubConditionType_ChapterStage                    SubConditionType = 2   //章节关卡
	SubConditionType_GainItemNum                     SubConditionType = 3   //获得指定数量道具
	SubConditionType_ArenaChallengeTimes             SubConditionType = 4   //竞技场挑战次数
	SubConditionType_WorldBossChallengeTimes         SubConditionType = 5   //世界boss挑战次数
	SubConditionType_MainQuestComplete               SubConditionType = 6   //完成主线任务
	SubConditionType_LevelId                         SubConditionType = 7   //到达关卡
	SubConditionType_RechargeCount                   SubConditionType = 8   //累计充值
	SubConditionType_GoodsNotEnoughClick             SubConditionType = 9   //道具不足时点击
	SubConditionType_GoodsNotEnough                  SubConditionType = 10  //道具不足时
	SubConditionType_CorpsPromotionLevel             SubConditionType = 11  //佣兵团晋升等级
	SubConditionType_DrawHeroTimes                   SubConditionType = 12  //英雄招募次数
	SubConditionType_HeroStarTimes                   SubConditionType = 13  //英雄升星次数
	SubConditionType_EquipStrengthTimes              SubConditionType = 14  //强化任意装备次数
	SubConditionType_ChapterChallengeTimes           SubConditionType = 15  //挑战主线关卡次数
	SubConditionType_HeroLevelUpTimes                SubConditionType = 16  //升级任意英雄次数
	SubConditionType_QuikAfkTimes                    SubConditionType = 17  //领取快速离线奖励次数
	SubConditionType_CommonAfkTimes                  SubConditionType = 18  //领取离线奖励次数
	SubConditionType_CommonHeroRecruit               SubConditionType = 19  //普通英雄招募次数
	SubConditionType_SeniorHeroRecruit               SubConditionType = 20  //高级英雄招募次数
	SubConditionType_TrialChallengeTimes             SubConditionType = 21  //挑战试炼次数
	SubConditionType_DiamondShopBuy                  SubConditionType = 22  //钻石商店购买道具次数
	SubConditionType_TowerChallengeTimes             SubConditionType = 23  //爬塔挑战次数
	SubConditionType_EightGateCombineTimes           SubConditionType = 24  //合成任意八门次数
	SubConditionType_ArtifactDrawTimes               SubConditionType = 25  //神器召唤次数
	SubConditionType_ArenaScore                      SubConditionType = 26  //竞技场积分
	SubConditionType_BattleFailed                    SubConditionType = 27  //战斗失败
	SubConditionType_AnyHeroStar                     SubConditionType = 28  //任意英雄星级
	SubConditionType_ArenaWinTimes                   SubConditionType = 29  //竞技场胜利次数
	SubConditionType_ShopBuyCnt                      SubConditionType = 30  //商店购买次数
	SubConditionType_ShopItemBuyCnt                  SubConditionType = 31  //商店商品购买次数
	SubConditionType_EquipLevelNum                   SubConditionType = 32  //装备等级x件
	SubConditionType_HeroLevelNum                    SubConditionType = 33  //英雄等级x个
	SubConditionType_ShopCurrencyCost                SubConditionType = 34  //商店货币花费
	SubConditionType_TowerFloorPass                  SubConditionType = 35  //通过塔某层
	SubConditionType_FishingNum                      SubConditionType = 37  //钓鱼X条
	SubConditionType_FishingCustom                   SubConditionType = 38  //钓到X条指定Y鱼
	SubConditionType_FishingType                     SubConditionType = 39  //钓到X条指定Y类型鱼
	SubConditionType_FishingItem                     SubConditionType = 40  //钓到X次指定物品Y
	SubConditionType_SellFish                        SubConditionType = 41  //出售鱼获得的银币X
	SubConditionType_SellRubbish                     SubConditionType = 42  //出售杂物获得的银币X
	SubConditionType_BuyBait                         SubConditionType = 43  //购买饵料次数X
	SubConditionType_FishingLength                   SubConditionType = 44  //钓到X条Y长度以上的鱼
	SubConditionType_FishingQuality                  SubConditionType = 45  //钓到X条Y品质以上的鱼
	SubConditionType_UpgradeFishingEquip             SubConditionType = 46  //强化渔具X次
	SubConditionType_FishOrderCnt                    SubConditionType = 47  //完成X次Y星级以上订单
	SubConditionType_FishingRubbishQuality           SubConditionType = 48  //钓到X次Y品质以上杂物
	SubConditionType_FishingCommitCustom             SubConditionType = 49  //钓到并提交订单X条指定Y鱼
	SubConditionType_FishingCommitType               SubConditionType = 50  //钓到并提交订单X条指定Y类型鱼
	SubConditionType_AllianceBossChallenge           SubConditionType = 51  //进行公会挑战N次
	SubConditionType_FarwarriorWin                   SubConditionType = 52  //远征战胜N波敌人
	SubConditionType_FarwarriorPass                  SubConditionType = 53  //远征通关N次
	SubConditionType_FarmUpLv                        SubConditionType = 54  //名人堂升级至X阶段Y级
	SubConditionType_FarmUpLvCnt                     SubConditionType = 55  //名人堂升级X次
	SubConditionType_BoxFavorUpLv                    SubConditionType = 56  //好感度达到X级
	SubConditionType_CostItem                        SubConditionType = 57  //消耗X道具Y个
	SubConditionType_RechargeDay                     SubConditionType = 58  //累计充值天数
	SubConditionType_SummonHero                      SubConditionType = 59  //召唤英雄次数
	SubConditionType_GambleFail                      SubConditionType = 60  //赌博失败次数
	SubConditionType_GambleSuccess                   SubConditionType = 61  //赌博成功次数
	SubConditionType_HoldItemNum                     SubConditionType = 62  //持有道具数量
	SubConditionType_EverOwnRarityHero               SubConditionType = 63  //拥有过某品质全部英雄
	SubConditionType_GetRarityHero                   SubConditionType = 64  //获得某品质英雄数量
	SubConditionType_HuntSuccessBeforeWave           SubConditionType = 65  //在x回合前狩猎成功y次
	SubConditionType_MonsterNum                      SubConditionType = 66  //场上怪物数量累积达到x
	SubConditionType_SummonHeroRarity                SubConditionType = 67  //召唤品质x英雄达到y个
	SubConditionType_HeroNum                         SubConditionType = 68  //场上存在英雄数量
	SubConditionType_HeroEverEnchantLevel            SubConditionType = 69  //英雄曾强化到x级
	SubConditionType_HeroEnchantFailTimes            SubConditionType = 70  //英雄强化失败x次
	SubConditionType_HeroEnchantSuccTimes            SubConditionType = 71  //英雄强化成功次数
	SubConditionType_HeroKillMonsterNum              SubConditionType = 72  //英雄击杀怪物数量
	SubConditionType_HeroEnchantLevelNums            SubConditionType = 73  //场上存在强化x级的英雄y个
	SubConditionType_SummonHeroTimes                 SubConditionType = 74  //召唤某英雄x次
	SubConditionType_PassLevelIdTimes                SubConditionType = 75  //通关某关卡ID多少次
	SubConditionType_RecruitHeroTimes                SubConditionType = 76  //招募某英雄x次
	SubConditionType_HeroSellTimes                   SubConditionType = 77  //出售英雄次数
	SubConditionType_TrailLevelIdPass                SubConditionType = 78  //试炼通关关卡
	SubConditionType_AllianceBuildTimes              SubConditionType = 79  //进行公会建设N次
	SubConditionType_AllianceTechGacha               SubConditionType = 80  //进行公会科技抽奖N次
	SubConditionType_SoulCompose                     SubConditionType = 81  //八门合成
	SubConditionType_ManorLevelupTimes               SubConditionType = 82  //领地升级x次
	SubConditionType_ManorLevelAchieve               SubConditionType = 83  //领地达到x级
	SubConditionType_BoxGiftOpenBox                  SubConditionType = 84  //格斗家礼物开箱次数
	SubConditionType_BoxGiftBoxLv                    SubConditionType = 85  //格斗家礼物箱子等级
	SubConditionType_TotalVipExp                     SubConditionType = 86  //累积vip经验
	SubConditionType_ManorGotRes                     SubConditionType = 87  //领地成功拉车X次
	SubConditionType_EquipHeroSoulTotal              SubConditionType = 88  //穿戴X个本源
	SubConditionType_AllianceBargainTimes            SubConditionType = 89  //公会砍价X次
	SubConditionType_QuestGainTypeAward              SubConditionType = 90  //领取某类型任务奖励
	SubConditionType_OwnHero                         SubConditionType = 91  //拥有某个英雄
	SubConditionType_ManorPlunder                    SubConditionType = 92  //领地掠夺X次
	SubConditionType_ArtifactDrawGearTimes           SubConditionType = 93  //神器档位召唤次数
	SubConditionType_CostHeroChip                    SubConditionType = 94  //消耗英雄碎片
	SubConditionType_AllianceTechObtain              SubConditionType = 95  //公会科技抽到碎片品质
	SubConditionType_SpecificHeroStar                SubConditionType = 96  //特定英雄星级
	SubConditionType_ArtifactDrawCost                SubConditionType = 97  //神器抽取消耗
	SubConditionType_FollowAllianceOccupyLevelCities SubConditionType = 98  //跟随公会占领X个X级城市
	SubConditionType_PhaseJoinCityWar                SubConditionType = 99  //阶段首次城市争夺参战
	SubConditionType_CreatePlayerDays                SubConditionType = 100 //创建角色满x天
	SubConditionType_CrossArenaRankCnt               SubConditionType = 101 //跨服32强排名次数
	SubConditionType_BountyPersonRankCnt             SubConditionType = 102 //赏金个人排名次数
	SubConditionType_BountyAllianceRankCnt           SubConditionType = 103 //赏金公会排名次数
	SubConditionType_LoginDays                       SubConditionType = 104 //登录天数
	SubConditionType_FightingRank                    SubConditionType = 105 //格斗圣殿排名
	SubConditionType_FameForgeTimes                  SubConditionType = 106 //名人堂锻造次数
	SubConditionType_ActivityOpenDays                SubConditionType = 107 //活动开启x天
	SubConditionType_GainCampBattleToken             SubConditionType = 108 //获得阵营战商店代币数量
	SubConditionType_BindDestinyHero                 SubConditionType = 109 //绑定宿命格斗家
	SubConditionType_InvitePlayerLevel               SubConditionType = 110 //邀请x个玩家升到x级
	SubConditionType_ContinuousScore                 SubConditionType = 111 //连续充值活动积分
	SubConditionType_MyselfJoinOtherFlipCard         SubConditionType = 112 //自己参加别人翻牌聚会
	SubConditionType_OtherJoinMyselfrFlipCard        SubConditionType = 113 //别人参加自己翻牌聚会
	SubConditionType_BindDestinyHeroStar             SubConditionType = 114 //绑定宿命英雄达到x星
	SubConditionType_ActivatePrivilegeCard           SubConditionType = 115 //激活权益卡
	SubConditionType_SystemOpened                    SubConditionType = 116 //开启某个系统
	SubConditionType_BindAndHaveDestinyHero          SubConditionType = 117 //绑定且拥有宿命英雄
	SubConditionType_TDPassTimes                     SubConditionType = 118 //呆呆通关次数
	SubConditionType_TopFivePower                    SubConditionType = 119 //前五战力
	SubConditionType_HaveRarityHeros                 SubConditionType = 120 //拥有X个X品质格斗家
	SubConditionType_HaveRarityArtifacts             SubConditionType = 121 //拥有X个X品质神器
	SubConditionType_HaveStarHeros                   SubConditionType = 122 //拥有X个X星级格斗家
	SubConditionType_HaveRarityGift                  SubConditionType = 123 //拥有X个X品质收藏手办
	SubConditionType_ElfMaxKill                      SubConditionType = 124 //南镇击杀数
	SubConditionType_WorldBossDamage                 SubConditionType = 125 //世界boss伤害
	SubConditionType_ActPuzzleActive                 SubConditionType = 126 //拼图格子点亮
	SubConditionType_ChannelChatTimes                SubConditionType = 127 //在X频道发言N次
	SubConditionType_AnyOneOfShopBuyCnt              SubConditionType = 128 //在多个商店购买N次
	SubConditionType_ServerOpenDays                  SubConditionType = 129 //开服天数
	SubConditionType_UnlimitedDrawLoginDays          SubConditionType = 130 //随心抽登录天数
	SubConditionType_HasLifeCard                     SubConditionType = 131 //拥有终身卡
	SubConditionType_WishHeroRecruit                 SubConditionType = 132 //自选招募次数
	SubConditionType_UpHeroRecruit                   SubConditionType = 133 //Up招募次数
	SubConditionType_PassAsTeamLeader                SubConditionType = 134 //作为队长通关x次
	SubConditionType_PassTeamStage                   SubConditionType = 135 //组队关卡通关x次
	SubConditionType_RoyalBossChallengeDamage        SubConditionType = 136 //大蛇boss挑战伤害
	SubConditionType_WorldBossModeDamage             SubConditionType = 138 //世界boss模式伤害
	SubConditionType_WorldBossModeAttack             SubConditionType = 139 //世界boss模式挑战
	SubConditionType_AllianceBossDamage              SubConditionType = 140 //公会boss伤害
	SubConditionType_FightingWinCnt                  SubConditionType = 141 //格斗盛典胜利场次
	SubConditionType_FightingChallengeCnt            SubConditionType = 142 //格斗盛典挑战次数
	SubConditionType_TDPassXTimes                    SubConditionType = 143 //呆呆通关多少波x次
	SubConditionType_TeamStage                       SubConditionType = 144 //组队关卡
	SubConditionType_GameClubFollow                  SubConditionType = 145 //关注游戏圈
	SubConditionType_GameClubBrowsingSeconds         SubConditionType = 146 //游戏圈浏览时长
	SubConditionType_GameClubLikeCount               SubConditionType = 147 //游戏圈内点赞帖子
	SubConditionType_GameClubCommentCount            SubConditionType = 148 //评论官方帖子
	SubConditionType_GameClubPost                    SubConditionType = 149 //游戏圈发帖数
	SubConditionType_TradeJoinTimes                  SubConditionType = 150 //参与X次贸易
	SubConditionType_TradeEarn                       SubConditionType = 151 //贸易赚取X贸易币
	SubConditionType_TradeAsset                      SubConditionType = 152 //贸易单局持有X贸易币
	SubConditionType_TradeRank                       SubConditionType = 153 //贸易获得前X名Y次
	SubConditionType_AnyShopItemBuyCnt               SubConditionType = 154 //商店任意商品购买次数
	SubConditionType_AvyTaskGotAward                 SubConditionType = 155 //活动任务领奖
	SubConditionType_OpenFashionBox                  SubConditionType = 156 //开启时装次数
	SubConditionType_KingWarPlayerScoreChange        SubConditionType = 157 //皇城个人积分改变
	SubConditionType_KingWarPlayerCostEnergy         SubConditionType = 158 //皇城累计消耗精力
	SubConditionType_AniActIslandDraw                SubConditionType = 159 //周年海岛探险
	SubConditionType_PetMerge                        SubConditionType = 160 //宠物合成
	SubConditionType_END                             SubConditionType = 161 //枚举结束
)

func (e SubConditionOperationType) String() string {

	if e == SubConditionOperationType_Equal {
		return "Equal"
	}
	if e == SubConditionOperationType_Reach {
		return "Reach"
	}
	if e == SubConditionOperationType_Less {
		return "Less"
	}
	if e == SubConditionOperationType_Lte {
		return "Lte"
	}

	return ""

}

func GetSubConditionOperationType(tp string) (SubConditionOperationType, bool) {
	switch tp {
	case "Equal":
		return SubConditionOperationType_Equal, true
	case "Reach":
		return SubConditionOperationType_Reach, true
	case "Less":
		return SubConditionOperationType_Less, true
	case "Lte":
		return SubConditionOperationType_Lte, true

	}
	return -1, false

}

func GetSubConditionOperationType_ByAlias(tp string) (SubConditionOperationType, bool) {
	switch tp {
	case "等于":
		return SubConditionOperationType_Equal, true
	case "达到":
		return SubConditionOperationType_Reach, true
	case "小于":
		return SubConditionOperationType_Less, true
	case "小于等于":
		return SubConditionOperationType_Lte, true

	}
	return -1, false

}

func (e CounterType) String() string {

	if e == CounterType_Trigger {
		return "Trigger"
	}
	if e == CounterType_Accumulate {
		return "Accumulate"
	}
	if e == CounterType_RepeatTrigger {
		return "RepeatTrigger"
	}

	return ""

}

func GetCounterType(tp string) (CounterType, bool) {
	switch tp {
	case "Trigger":
		return CounterType_Trigger, true
	case "Accumulate":
		return CounterType_Accumulate, true
	case "RepeatTrigger":
		return CounterType_RepeatTrigger, true

	}
	return -1, false

}

func GetCounterType_ByAlias(tp string) (CounterType, bool) {
	switch tp {
	case "触发计数":
		return CounterType_Trigger, true
	case "累积计数":
		return CounterType_Accumulate, true
	case "重复触发":
		return CounterType_RepeatTrigger, true

	}
	return -1, false

}

func (e SubConditionRelation) String() string {

	if e == SubConditionRelation_And {
		return "And"
	}
	if e == SubConditionRelation_Or {
		return "Or"
	}

	return ""

}

func GetSubConditionRelation(tp string) (SubConditionRelation, bool) {
	switch tp {
	case "And":
		return SubConditionRelation_And, true
	case "Or":
		return SubConditionRelation_Or, true

	}
	return -1, false

}

func GetSubConditionRelation_ByAlias(tp string) (SubConditionRelation, bool) {
	switch tp {
	case "满足全部子条件":
		return SubConditionRelation_And, true
	case "满足一个子条件":
		return SubConditionRelation_Or, true

	}
	return -1, false

}

func (e SubConditionType) String() string {

	if e == SubConditionType_Empty {
		return "Empty"
	}
	if e == SubConditionType_PlayerLevel {
		return "PlayerLevel"
	}
	if e == SubConditionType_ChapterStage {
		return "ChapterStage"
	}
	if e == SubConditionType_GainItemNum {
		return "GainItemNum"
	}
	if e == SubConditionType_ArenaChallengeTimes {
		return "ArenaChallengeTimes"
	}
	if e == SubConditionType_WorldBossChallengeTimes {
		return "WorldBossChallengeTimes"
	}
	if e == SubConditionType_MainQuestComplete {
		return "MainQuestComplete"
	}
	if e == SubConditionType_LevelId {
		return "LevelId"
	}
	if e == SubConditionType_RechargeCount {
		return "RechargeCount"
	}
	if e == SubConditionType_GoodsNotEnoughClick {
		return "GoodsNotEnoughClick"
	}
	if e == SubConditionType_GoodsNotEnough {
		return "GoodsNotEnough"
	}
	if e == SubConditionType_CorpsPromotionLevel {
		return "CorpsPromotionLevel"
	}
	if e == SubConditionType_DrawHeroTimes {
		return "DrawHeroTimes"
	}
	if e == SubConditionType_HeroStarTimes {
		return "HeroStarTimes"
	}
	if e == SubConditionType_EquipStrengthTimes {
		return "EquipStrengthTimes"
	}
	if e == SubConditionType_ChapterChallengeTimes {
		return "ChapterChallengeTimes"
	}
	if e == SubConditionType_HeroLevelUpTimes {
		return "HeroLevelUpTimes"
	}
	if e == SubConditionType_QuikAfkTimes {
		return "QuikAfkTimes"
	}
	if e == SubConditionType_CommonAfkTimes {
		return "CommonAfkTimes"
	}
	if e == SubConditionType_CommonHeroRecruit {
		return "CommonHeroRecruit"
	}
	if e == SubConditionType_SeniorHeroRecruit {
		return "SeniorHeroRecruit"
	}
	if e == SubConditionType_TrialChallengeTimes {
		return "TrialChallengeTimes"
	}
	if e == SubConditionType_DiamondShopBuy {
		return "DiamondShopBuy"
	}
	if e == SubConditionType_TowerChallengeTimes {
		return "TowerChallengeTimes"
	}
	if e == SubConditionType_EightGateCombineTimes {
		return "EightGateCombineTimes"
	}
	if e == SubConditionType_ArtifactDrawTimes {
		return "ArtifactDrawTimes"
	}
	if e == SubConditionType_ArenaScore {
		return "ArenaScore"
	}
	if e == SubConditionType_BattleFailed {
		return "BattleFailed"
	}
	if e == SubConditionType_AnyHeroStar {
		return "AnyHeroStar"
	}
	if e == SubConditionType_ArenaWinTimes {
		return "ArenaWinTimes"
	}
	if e == SubConditionType_ShopBuyCnt {
		return "ShopBuyCnt"
	}
	if e == SubConditionType_ShopItemBuyCnt {
		return "ShopItemBuyCnt"
	}
	if e == SubConditionType_EquipLevelNum {
		return "EquipLevelNum"
	}
	if e == SubConditionType_HeroLevelNum {
		return "HeroLevelNum"
	}
	if e == SubConditionType_ShopCurrencyCost {
		return "ShopCurrencyCost"
	}
	if e == SubConditionType_TowerFloorPass {
		return "TowerFloorPass"
	}
	if e == SubConditionType_FishingNum {
		return "FishingNum"
	}
	if e == SubConditionType_FishingCustom {
		return "FishingCustom"
	}
	if e == SubConditionType_FishingType {
		return "FishingType"
	}
	if e == SubConditionType_FishingItem {
		return "FishingItem"
	}
	if e == SubConditionType_SellFish {
		return "SellFish"
	}
	if e == SubConditionType_SellRubbish {
		return "SellRubbish"
	}
	if e == SubConditionType_BuyBait {
		return "BuyBait"
	}
	if e == SubConditionType_FishingLength {
		return "FishingLength"
	}
	if e == SubConditionType_FishingQuality {
		return "FishingQuality"
	}
	if e == SubConditionType_UpgradeFishingEquip {
		return "UpgradeFishingEquip"
	}
	if e == SubConditionType_FishOrderCnt {
		return "FishOrderCnt"
	}
	if e == SubConditionType_FishingRubbishQuality {
		return "FishingRubbishQuality"
	}
	if e == SubConditionType_FishingCommitCustom {
		return "FishingCommitCustom"
	}
	if e == SubConditionType_FishingCommitType {
		return "FishingCommitType"
	}
	if e == SubConditionType_AllianceBossChallenge {
		return "AllianceBossChallenge"
	}
	if e == SubConditionType_FarwarriorWin {
		return "FarwarriorWin"
	}
	if e == SubConditionType_FarwarriorPass {
		return "FarwarriorPass"
	}
	if e == SubConditionType_FarmUpLv {
		return "FarmUpLv"
	}
	if e == SubConditionType_FarmUpLvCnt {
		return "FarmUpLvCnt"
	}
	if e == SubConditionType_BoxFavorUpLv {
		return "BoxFavorUpLv"
	}
	if e == SubConditionType_CostItem {
		return "CostItem"
	}
	if e == SubConditionType_RechargeDay {
		return "RechargeDay"
	}
	if e == SubConditionType_SummonHero {
		return "SummonHero"
	}
	if e == SubConditionType_GambleFail {
		return "GambleFail"
	}
	if e == SubConditionType_GambleSuccess {
		return "GambleSuccess"
	}
	if e == SubConditionType_HoldItemNum {
		return "HoldItemNum"
	}
	if e == SubConditionType_EverOwnRarityHero {
		return "EverOwnRarityHero"
	}
	if e == SubConditionType_GetRarityHero {
		return "GetRarityHero"
	}
	if e == SubConditionType_HuntSuccessBeforeWave {
		return "HuntSuccessBeforeWave"
	}
	if e == SubConditionType_MonsterNum {
		return "MonsterNum"
	}
	if e == SubConditionType_SummonHeroRarity {
		return "SummonHeroRarity"
	}
	if e == SubConditionType_HeroNum {
		return "HeroNum"
	}
	if e == SubConditionType_HeroEverEnchantLevel {
		return "HeroEverEnchantLevel"
	}
	if e == SubConditionType_HeroEnchantFailTimes {
		return "HeroEnchantFailTimes"
	}
	if e == SubConditionType_HeroEnchantSuccTimes {
		return "HeroEnchantSuccTimes"
	}
	if e == SubConditionType_HeroKillMonsterNum {
		return "HeroKillMonsterNum"
	}
	if e == SubConditionType_HeroEnchantLevelNums {
		return "HeroEnchantLevelNums"
	}
	if e == SubConditionType_SummonHeroTimes {
		return "SummonHeroTimes"
	}
	if e == SubConditionType_PassLevelIdTimes {
		return "PassLevelIdTimes"
	}
	if e == SubConditionType_RecruitHeroTimes {
		return "RecruitHeroTimes"
	}
	if e == SubConditionType_HeroSellTimes {
		return "HeroSellTimes"
	}
	if e == SubConditionType_TrailLevelIdPass {
		return "TrailLevelIdPass"
	}
	if e == SubConditionType_AllianceBuildTimes {
		return "AllianceBuildTimes"
	}
	if e == SubConditionType_AllianceTechGacha {
		return "AllianceTechGacha"
	}
	if e == SubConditionType_SoulCompose {
		return "SoulCompose"
	}
	if e == SubConditionType_ManorLevelupTimes {
		return "ManorLevelupTimes"
	}
	if e == SubConditionType_ManorLevelAchieve {
		return "ManorLevelAchieve"
	}
	if e == SubConditionType_BoxGiftOpenBox {
		return "BoxGiftOpenBox"
	}
	if e == SubConditionType_BoxGiftBoxLv {
		return "BoxGiftBoxLv"
	}
	if e == SubConditionType_TotalVipExp {
		return "TotalVipExp"
	}
	if e == SubConditionType_ManorGotRes {
		return "ManorGotRes"
	}
	if e == SubConditionType_EquipHeroSoulTotal {
		return "EquipHeroSoulTotal"
	}
	if e == SubConditionType_AllianceBargainTimes {
		return "AllianceBargainTimes"
	}
	if e == SubConditionType_QuestGainTypeAward {
		return "QuestGainTypeAward"
	}
	if e == SubConditionType_OwnHero {
		return "OwnHero"
	}
	if e == SubConditionType_ManorPlunder {
		return "ManorPlunder"
	}
	if e == SubConditionType_ArtifactDrawGearTimes {
		return "ArtifactDrawGearTimes"
	}
	if e == SubConditionType_CostHeroChip {
		return "CostHeroChip"
	}
	if e == SubConditionType_AllianceTechObtain {
		return "AllianceTechObtain"
	}
	if e == SubConditionType_SpecificHeroStar {
		return "SpecificHeroStar"
	}
	if e == SubConditionType_ArtifactDrawCost {
		return "ArtifactDrawCost"
	}
	if e == SubConditionType_FollowAllianceOccupyLevelCities {
		return "FollowAllianceOccupyLevelCities"
	}
	if e == SubConditionType_PhaseJoinCityWar {
		return "PhaseJoinCityWar"
	}
	if e == SubConditionType_CreatePlayerDays {
		return "CreatePlayerDays"
	}
	if e == SubConditionType_CrossArenaRankCnt {
		return "CrossArenaRankCnt"
	}
	if e == SubConditionType_BountyPersonRankCnt {
		return "BountyPersonRankCnt"
	}
	if e == SubConditionType_BountyAllianceRankCnt {
		return "BountyAllianceRankCnt"
	}
	if e == SubConditionType_LoginDays {
		return "LoginDays"
	}
	if e == SubConditionType_FightingRank {
		return "FightingRank"
	}
	if e == SubConditionType_FameForgeTimes {
		return "FameForgeTimes"
	}
	if e == SubConditionType_ActivityOpenDays {
		return "ActivityOpenDays"
	}
	if e == SubConditionType_GainCampBattleToken {
		return "GainCampBattleToken"
	}
	if e == SubConditionType_BindDestinyHero {
		return "BindDestinyHero"
	}
	if e == SubConditionType_InvitePlayerLevel {
		return "InvitePlayerLevel"
	}
	if e == SubConditionType_ContinuousScore {
		return "ContinuousScore"
	}
	if e == SubConditionType_MyselfJoinOtherFlipCard {
		return "MyselfJoinOtherFlipCard"
	}
	if e == SubConditionType_OtherJoinMyselfrFlipCard {
		return "OtherJoinMyselfrFlipCard"
	}
	if e == SubConditionType_BindDestinyHeroStar {
		return "BindDestinyHeroStar"
	}
	if e == SubConditionType_ActivatePrivilegeCard {
		return "ActivatePrivilegeCard"
	}
	if e == SubConditionType_SystemOpened {
		return "SystemOpened"
	}
	if e == SubConditionType_BindAndHaveDestinyHero {
		return "BindAndHaveDestinyHero"
	}
	if e == SubConditionType_TDPassTimes {
		return "TDPassTimes"
	}
	if e == SubConditionType_TopFivePower {
		return "TopFivePower"
	}
	if e == SubConditionType_HaveRarityHeros {
		return "HaveRarityHeros"
	}
	if e == SubConditionType_HaveRarityArtifacts {
		return "HaveRarityArtifacts"
	}
	if e == SubConditionType_HaveStarHeros {
		return "HaveStarHeros"
	}
	if e == SubConditionType_HaveRarityGift {
		return "HaveRarityGift"
	}
	if e == SubConditionType_ElfMaxKill {
		return "ElfMaxKill"
	}
	if e == SubConditionType_WorldBossDamage {
		return "WorldBossDamage"
	}
	if e == SubConditionType_ActPuzzleActive {
		return "ActPuzzleActive"
	}
	if e == SubConditionType_ChannelChatTimes {
		return "ChannelChatTimes"
	}
	if e == SubConditionType_AnyOneOfShopBuyCnt {
		return "AnyOneOfShopBuyCnt"
	}
	if e == SubConditionType_ServerOpenDays {
		return "ServerOpenDays"
	}
	if e == SubConditionType_UnlimitedDrawLoginDays {
		return "UnlimitedDrawLoginDays"
	}
	if e == SubConditionType_HasLifeCard {
		return "HasLifeCard"
	}
	if e == SubConditionType_WishHeroRecruit {
		return "WishHeroRecruit"
	}
	if e == SubConditionType_UpHeroRecruit {
		return "UpHeroRecruit"
	}
	if e == SubConditionType_PassAsTeamLeader {
		return "PassAsTeamLeader"
	}
	if e == SubConditionType_PassTeamStage {
		return "PassTeamStage"
	}
	if e == SubConditionType_RoyalBossChallengeDamage {
		return "RoyalBossChallengeDamage"
	}
	if e == SubConditionType_WorldBossModeDamage {
		return "WorldBossModeDamage"
	}
	if e == SubConditionType_WorldBossModeAttack {
		return "WorldBossModeAttack"
	}
	if e == SubConditionType_AllianceBossDamage {
		return "AllianceBossDamage"
	}
	if e == SubConditionType_FightingWinCnt {
		return "FightingWinCnt"
	}
	if e == SubConditionType_FightingChallengeCnt {
		return "FightingChallengeCnt"
	}
	if e == SubConditionType_TDPassXTimes {
		return "TDPassXTimes"
	}
	if e == SubConditionType_TeamStage {
		return "TeamStage"
	}
	if e == SubConditionType_GameClubFollow {
		return "GameClubFollow"
	}
	if e == SubConditionType_GameClubBrowsingSeconds {
		return "GameClubBrowsingSeconds"
	}
	if e == SubConditionType_GameClubLikeCount {
		return "GameClubLikeCount"
	}
	if e == SubConditionType_GameClubCommentCount {
		return "GameClubCommentCount"
	}
	if e == SubConditionType_GameClubPost {
		return "GameClubPost"
	}
	if e == SubConditionType_TradeJoinTimes {
		return "TradeJoinTimes"
	}
	if e == SubConditionType_TradeEarn {
		return "TradeEarn"
	}
	if e == SubConditionType_TradeAsset {
		return "TradeAsset"
	}
	if e == SubConditionType_TradeRank {
		return "TradeRank"
	}
	if e == SubConditionType_AnyShopItemBuyCnt {
		return "AnyShopItemBuyCnt"
	}
	if e == SubConditionType_AvyTaskGotAward {
		return "AvyTaskGotAward"
	}
	if e == SubConditionType_OpenFashionBox {
		return "OpenFashionBox"
	}
	if e == SubConditionType_KingWarPlayerScoreChange {
		return "KingWarPlayerScoreChange"
	}
	if e == SubConditionType_KingWarPlayerCostEnergy {
		return "KingWarPlayerCostEnergy"
	}
	if e == SubConditionType_AniActIslandDraw {
		return "AniActIslandDraw"
	}
	if e == SubConditionType_PetMerge {
		return "PetMerge"
	}

	return ""

}

func GetSubConditionType(tp string) (SubConditionType, bool) {
	switch tp {
	case "Empty":
		return SubConditionType_Empty, true
	case "PlayerLevel":
		return SubConditionType_PlayerLevel, true
	case "ChapterStage":
		return SubConditionType_ChapterStage, true
	case "GainItemNum":
		return SubConditionType_GainItemNum, true
	case "ArenaChallengeTimes":
		return SubConditionType_ArenaChallengeTimes, true
	case "WorldBossChallengeTimes":
		return SubConditionType_WorldBossChallengeTimes, true
	case "MainQuestComplete":
		return SubConditionType_MainQuestComplete, true
	case "LevelId":
		return SubConditionType_LevelId, true
	case "RechargeCount":
		return SubConditionType_RechargeCount, true
	case "GoodsNotEnoughClick":
		return SubConditionType_GoodsNotEnoughClick, true
	case "GoodsNotEnough":
		return SubConditionType_GoodsNotEnough, true
	case "CorpsPromotionLevel":
		return SubConditionType_CorpsPromotionLevel, true
	case "DrawHeroTimes":
		return SubConditionType_DrawHeroTimes, true
	case "HeroStarTimes":
		return SubConditionType_HeroStarTimes, true
	case "EquipStrengthTimes":
		return SubConditionType_EquipStrengthTimes, true
	case "ChapterChallengeTimes":
		return SubConditionType_ChapterChallengeTimes, true
	case "HeroLevelUpTimes":
		return SubConditionType_HeroLevelUpTimes, true
	case "QuikAfkTimes":
		return SubConditionType_QuikAfkTimes, true
	case "CommonAfkTimes":
		return SubConditionType_CommonAfkTimes, true
	case "CommonHeroRecruit":
		return SubConditionType_CommonHeroRecruit, true
	case "SeniorHeroRecruit":
		return SubConditionType_SeniorHeroRecruit, true
	case "TrialChallengeTimes":
		return SubConditionType_TrialChallengeTimes, true
	case "DiamondShopBuy":
		return SubConditionType_DiamondShopBuy, true
	case "TowerChallengeTimes":
		return SubConditionType_TowerChallengeTimes, true
	case "EightGateCombineTimes":
		return SubConditionType_EightGateCombineTimes, true
	case "ArtifactDrawTimes":
		return SubConditionType_ArtifactDrawTimes, true
	case "ArenaScore":
		return SubConditionType_ArenaScore, true
	case "BattleFailed":
		return SubConditionType_BattleFailed, true
	case "AnyHeroStar":
		return SubConditionType_AnyHeroStar, true
	case "ArenaWinTimes":
		return SubConditionType_ArenaWinTimes, true
	case "ShopBuyCnt":
		return SubConditionType_ShopBuyCnt, true
	case "ShopItemBuyCnt":
		return SubConditionType_ShopItemBuyCnt, true
	case "EquipLevelNum":
		return SubConditionType_EquipLevelNum, true
	case "HeroLevelNum":
		return SubConditionType_HeroLevelNum, true
	case "ShopCurrencyCost":
		return SubConditionType_ShopCurrencyCost, true
	case "TowerFloorPass":
		return SubConditionType_TowerFloorPass, true
	case "FishingNum":
		return SubConditionType_FishingNum, true
	case "FishingCustom":
		return SubConditionType_FishingCustom, true
	case "FishingType":
		return SubConditionType_FishingType, true
	case "FishingItem":
		return SubConditionType_FishingItem, true
	case "SellFish":
		return SubConditionType_SellFish, true
	case "SellRubbish":
		return SubConditionType_SellRubbish, true
	case "BuyBait":
		return SubConditionType_BuyBait, true
	case "FishingLength":
		return SubConditionType_FishingLength, true
	case "FishingQuality":
		return SubConditionType_FishingQuality, true
	case "UpgradeFishingEquip":
		return SubConditionType_UpgradeFishingEquip, true
	case "FishOrderCnt":
		return SubConditionType_FishOrderCnt, true
	case "FishingRubbishQuality":
		return SubConditionType_FishingRubbishQuality, true
	case "FishingCommitCustom":
		return SubConditionType_FishingCommitCustom, true
	case "FishingCommitType":
		return SubConditionType_FishingCommitType, true
	case "AllianceBossChallenge":
		return SubConditionType_AllianceBossChallenge, true
	case "FarwarriorWin":
		return SubConditionType_FarwarriorWin, true
	case "FarwarriorPass":
		return SubConditionType_FarwarriorPass, true
	case "FarmUpLv":
		return SubConditionType_FarmUpLv, true
	case "FarmUpLvCnt":
		return SubConditionType_FarmUpLvCnt, true
	case "BoxFavorUpLv":
		return SubConditionType_BoxFavorUpLv, true
	case "CostItem":
		return SubConditionType_CostItem, true
	case "RechargeDay":
		return SubConditionType_RechargeDay, true
	case "SummonHero":
		return SubConditionType_SummonHero, true
	case "GambleFail":
		return SubConditionType_GambleFail, true
	case "GambleSuccess":
		return SubConditionType_GambleSuccess, true
	case "HoldItemNum":
		return SubConditionType_HoldItemNum, true
	case "EverOwnRarityHero":
		return SubConditionType_EverOwnRarityHero, true
	case "GetRarityHero":
		return SubConditionType_GetRarityHero, true
	case "HuntSuccessBeforeWave":
		return SubConditionType_HuntSuccessBeforeWave, true
	case "MonsterNum":
		return SubConditionType_MonsterNum, true
	case "SummonHeroRarity":
		return SubConditionType_SummonHeroRarity, true
	case "HeroNum":
		return SubConditionType_HeroNum, true
	case "HeroEverEnchantLevel":
		return SubConditionType_HeroEverEnchantLevel, true
	case "HeroEnchantFailTimes":
		return SubConditionType_HeroEnchantFailTimes, true
	case "HeroEnchantSuccTimes":
		return SubConditionType_HeroEnchantSuccTimes, true
	case "HeroKillMonsterNum":
		return SubConditionType_HeroKillMonsterNum, true
	case "HeroEnchantLevelNums":
		return SubConditionType_HeroEnchantLevelNums, true
	case "SummonHeroTimes":
		return SubConditionType_SummonHeroTimes, true
	case "PassLevelIdTimes":
		return SubConditionType_PassLevelIdTimes, true
	case "RecruitHeroTimes":
		return SubConditionType_RecruitHeroTimes, true
	case "HeroSellTimes":
		return SubConditionType_HeroSellTimes, true
	case "TrailLevelIdPass":
		return SubConditionType_TrailLevelIdPass, true
	case "AllianceBuildTimes":
		return SubConditionType_AllianceBuildTimes, true
	case "AllianceTechGacha":
		return SubConditionType_AllianceTechGacha, true
	case "SoulCompose":
		return SubConditionType_SoulCompose, true
	case "ManorLevelupTimes":
		return SubConditionType_ManorLevelupTimes, true
	case "ManorLevelAchieve":
		return SubConditionType_ManorLevelAchieve, true
	case "BoxGiftOpenBox":
		return SubConditionType_BoxGiftOpenBox, true
	case "BoxGiftBoxLv":
		return SubConditionType_BoxGiftBoxLv, true
	case "TotalVipExp":
		return SubConditionType_TotalVipExp, true
	case "ManorGotRes":
		return SubConditionType_ManorGotRes, true
	case "EquipHeroSoulTotal":
		return SubConditionType_EquipHeroSoulTotal, true
	case "AllianceBargainTimes":
		return SubConditionType_AllianceBargainTimes, true
	case "QuestGainTypeAward":
		return SubConditionType_QuestGainTypeAward, true
	case "OwnHero":
		return SubConditionType_OwnHero, true
	case "ManorPlunder":
		return SubConditionType_ManorPlunder, true
	case "ArtifactDrawGearTimes":
		return SubConditionType_ArtifactDrawGearTimes, true
	case "CostHeroChip":
		return SubConditionType_CostHeroChip, true
	case "AllianceTechObtain":
		return SubConditionType_AllianceTechObtain, true
	case "SpecificHeroStar":
		return SubConditionType_SpecificHeroStar, true
	case "ArtifactDrawCost":
		return SubConditionType_ArtifactDrawCost, true
	case "FollowAllianceOccupyLevelCities":
		return SubConditionType_FollowAllianceOccupyLevelCities, true
	case "PhaseJoinCityWar":
		return SubConditionType_PhaseJoinCityWar, true
	case "CreatePlayerDays":
		return SubConditionType_CreatePlayerDays, true
	case "CrossArenaRankCnt":
		return SubConditionType_CrossArenaRankCnt, true
	case "BountyPersonRankCnt":
		return SubConditionType_BountyPersonRankCnt, true
	case "BountyAllianceRankCnt":
		return SubConditionType_BountyAllianceRankCnt, true
	case "LoginDays":
		return SubConditionType_LoginDays, true
	case "FightingRank":
		return SubConditionType_FightingRank, true
	case "FameForgeTimes":
		return SubConditionType_FameForgeTimes, true
	case "ActivityOpenDays":
		return SubConditionType_ActivityOpenDays, true
	case "GainCampBattleToken":
		return SubConditionType_GainCampBattleToken, true
	case "BindDestinyHero":
		return SubConditionType_BindDestinyHero, true
	case "InvitePlayerLevel":
		return SubConditionType_InvitePlayerLevel, true
	case "ContinuousScore":
		return SubConditionType_ContinuousScore, true
	case "MyselfJoinOtherFlipCard":
		return SubConditionType_MyselfJoinOtherFlipCard, true
	case "OtherJoinMyselfrFlipCard":
		return SubConditionType_OtherJoinMyselfrFlipCard, true
	case "BindDestinyHeroStar":
		return SubConditionType_BindDestinyHeroStar, true
	case "ActivatePrivilegeCard":
		return SubConditionType_ActivatePrivilegeCard, true
	case "SystemOpened":
		return SubConditionType_SystemOpened, true
	case "BindAndHaveDestinyHero":
		return SubConditionType_BindAndHaveDestinyHero, true
	case "TDPassTimes":
		return SubConditionType_TDPassTimes, true
	case "TopFivePower":
		return SubConditionType_TopFivePower, true
	case "HaveRarityHeros":
		return SubConditionType_HaveRarityHeros, true
	case "HaveRarityArtifacts":
		return SubConditionType_HaveRarityArtifacts, true
	case "HaveStarHeros":
		return SubConditionType_HaveStarHeros, true
	case "HaveRarityGift":
		return SubConditionType_HaveRarityGift, true
	case "ElfMaxKill":
		return SubConditionType_ElfMaxKill, true
	case "WorldBossDamage":
		return SubConditionType_WorldBossDamage, true
	case "ActPuzzleActive":
		return SubConditionType_ActPuzzleActive, true
	case "ChannelChatTimes":
		return SubConditionType_ChannelChatTimes, true
	case "AnyOneOfShopBuyCnt":
		return SubConditionType_AnyOneOfShopBuyCnt, true
	case "ServerOpenDays":
		return SubConditionType_ServerOpenDays, true
	case "UnlimitedDrawLoginDays":
		return SubConditionType_UnlimitedDrawLoginDays, true
	case "HasLifeCard":
		return SubConditionType_HasLifeCard, true
	case "WishHeroRecruit":
		return SubConditionType_WishHeroRecruit, true
	case "UpHeroRecruit":
		return SubConditionType_UpHeroRecruit, true
	case "PassAsTeamLeader":
		return SubConditionType_PassAsTeamLeader, true
	case "PassTeamStage":
		return SubConditionType_PassTeamStage, true
	case "RoyalBossChallengeDamage":
		return SubConditionType_RoyalBossChallengeDamage, true
	case "WorldBossModeDamage":
		return SubConditionType_WorldBossModeDamage, true
	case "WorldBossModeAttack":
		return SubConditionType_WorldBossModeAttack, true
	case "AllianceBossDamage":
		return SubConditionType_AllianceBossDamage, true
	case "FightingWinCnt":
		return SubConditionType_FightingWinCnt, true
	case "FightingChallengeCnt":
		return SubConditionType_FightingChallengeCnt, true
	case "TDPassXTimes":
		return SubConditionType_TDPassXTimes, true
	case "TeamStage":
		return SubConditionType_TeamStage, true
	case "GameClubFollow":
		return SubConditionType_GameClubFollow, true
	case "GameClubBrowsingSeconds":
		return SubConditionType_GameClubBrowsingSeconds, true
	case "GameClubLikeCount":
		return SubConditionType_GameClubLikeCount, true
	case "GameClubCommentCount":
		return SubConditionType_GameClubCommentCount, true
	case "GameClubPost":
		return SubConditionType_GameClubPost, true
	case "TradeJoinTimes":
		return SubConditionType_TradeJoinTimes, true
	case "TradeEarn":
		return SubConditionType_TradeEarn, true
	case "TradeAsset":
		return SubConditionType_TradeAsset, true
	case "TradeRank":
		return SubConditionType_TradeRank, true
	case "AnyShopItemBuyCnt":
		return SubConditionType_AnyShopItemBuyCnt, true
	case "AvyTaskGotAward":
		return SubConditionType_AvyTaskGotAward, true
	case "OpenFashionBox":
		return SubConditionType_OpenFashionBox, true
	case "KingWarPlayerScoreChange":
		return SubConditionType_KingWarPlayerScoreChange, true
	case "KingWarPlayerCostEnergy":
		return SubConditionType_KingWarPlayerCostEnergy, true
	case "AniActIslandDraw":
		return SubConditionType_AniActIslandDraw, true
	case "PetMerge":
		return SubConditionType_PetMerge, true

	}
	return -1, false

}

func GetSubConditionType_ByAlias(tp string) (SubConditionType, bool) {
	switch tp {
	case "空":
		return SubConditionType_Empty, true
	case "玩家等级":
		return SubConditionType_PlayerLevel, true
	case "章节关卡":
		return SubConditionType_ChapterStage, true
	case "获得指定数量道具":
		return SubConditionType_GainItemNum, true
	case "竞技场挑战次数":
		return SubConditionType_ArenaChallengeTimes, true
	case "世界boss挑战次数":
		return SubConditionType_WorldBossChallengeTimes, true
	case "完成主线任务":
		return SubConditionType_MainQuestComplete, true
	case "到达关卡":
		return SubConditionType_LevelId, true
	case "累计充值":
		return SubConditionType_RechargeCount, true
	case "道具不足时点击":
		return SubConditionType_GoodsNotEnoughClick, true
	case "道具不足时":
		return SubConditionType_GoodsNotEnough, true
	case "佣兵团晋升等级":
		return SubConditionType_CorpsPromotionLevel, true
	case "英雄招募次数":
		return SubConditionType_DrawHeroTimes, true
	case "英雄升星次数":
		return SubConditionType_HeroStarTimes, true
	case "强化任意装备次数":
		return SubConditionType_EquipStrengthTimes, true
	case "挑战主线关卡次数":
		return SubConditionType_ChapterChallengeTimes, true
	case "升级任意英雄次数":
		return SubConditionType_HeroLevelUpTimes, true
	case "领取快速离线奖励次数":
		return SubConditionType_QuikAfkTimes, true
	case "领取离线奖励次数":
		return SubConditionType_CommonAfkTimes, true
	case "普通英雄招募次数":
		return SubConditionType_CommonHeroRecruit, true
	case "高级英雄招募次数":
		return SubConditionType_SeniorHeroRecruit, true
	case "挑战试炼次数":
		return SubConditionType_TrialChallengeTimes, true
	case "钻石商店购买道具次数":
		return SubConditionType_DiamondShopBuy, true
	case "爬塔挑战次数":
		return SubConditionType_TowerChallengeTimes, true
	case "合成任意八门次数":
		return SubConditionType_EightGateCombineTimes, true
	case "神器召唤次数":
		return SubConditionType_ArtifactDrawTimes, true
	case "竞技场积分":
		return SubConditionType_ArenaScore, true
	case "战斗失败":
		return SubConditionType_BattleFailed, true
	case "任意英雄星级":
		return SubConditionType_AnyHeroStar, true
	case "竞技场胜利次数":
		return SubConditionType_ArenaWinTimes, true
	case "商店购买次数":
		return SubConditionType_ShopBuyCnt, true
	case "商店商品购买次数":
		return SubConditionType_ShopItemBuyCnt, true
	case "装备等级x件":
		return SubConditionType_EquipLevelNum, true
	case "英雄等级x个":
		return SubConditionType_HeroLevelNum, true
	case "商店货币花费":
		return SubConditionType_ShopCurrencyCost, true
	case "通过塔某层":
		return SubConditionType_TowerFloorPass, true
	case "钓鱼X条":
		return SubConditionType_FishingNum, true
	case "钓到X条指定Y鱼":
		return SubConditionType_FishingCustom, true
	case "钓到X条指定Y类型鱼":
		return SubConditionType_FishingType, true
	case "钓到X次指定物品Y":
		return SubConditionType_FishingItem, true
	case "出售鱼获得的银币X":
		return SubConditionType_SellFish, true
	case "出售杂物获得的银币X":
		return SubConditionType_SellRubbish, true
	case "购买饵料次数X":
		return SubConditionType_BuyBait, true
	case "钓到X条Y长度以上的鱼":
		return SubConditionType_FishingLength, true
	case "钓到X条Y品质以上的鱼":
		return SubConditionType_FishingQuality, true
	case "强化渔具X次":
		return SubConditionType_UpgradeFishingEquip, true
	case "完成X次Y星级以上订单":
		return SubConditionType_FishOrderCnt, true
	case "钓到X次Y品质以上杂物":
		return SubConditionType_FishingRubbishQuality, true
	case "钓到并提交订单X条指定Y鱼":
		return SubConditionType_FishingCommitCustom, true
	case "钓到并提交订单X条指定Y类型鱼":
		return SubConditionType_FishingCommitType, true
	case "进行公会挑战N次":
		return SubConditionType_AllianceBossChallenge, true
	case "远征战胜N波敌人":
		return SubConditionType_FarwarriorWin, true
	case "远征通关N次":
		return SubConditionType_FarwarriorPass, true
	case "名人堂升级至X阶段Y级":
		return SubConditionType_FarmUpLv, true
	case "名人堂升级X次":
		return SubConditionType_FarmUpLvCnt, true
	case "好感度达到X级":
		return SubConditionType_BoxFavorUpLv, true
	case "消耗X道具Y个":
		return SubConditionType_CostItem, true
	case "累计充值天数":
		return SubConditionType_RechargeDay, true
	case "召唤英雄次数":
		return SubConditionType_SummonHero, true
	case "赌博失败次数":
		return SubConditionType_GambleFail, true
	case "赌博成功次数":
		return SubConditionType_GambleSuccess, true
	case "持有道具数量":
		return SubConditionType_HoldItemNum, true
	case "拥有过某品质全部英雄":
		return SubConditionType_EverOwnRarityHero, true
	case "获得某品质英雄数量":
		return SubConditionType_GetRarityHero, true
	case "在x回合前狩猎成功y次":
		return SubConditionType_HuntSuccessBeforeWave, true
	case "场上怪物数量累积达到x":
		return SubConditionType_MonsterNum, true
	case "召唤品质x英雄达到y个":
		return SubConditionType_SummonHeroRarity, true
	case "场上存在英雄数量":
		return SubConditionType_HeroNum, true
	case "英雄曾强化到x级":
		return SubConditionType_HeroEverEnchantLevel, true
	case "英雄强化失败x次":
		return SubConditionType_HeroEnchantFailTimes, true
	case "英雄强化成功次数":
		return SubConditionType_HeroEnchantSuccTimes, true
	case "英雄击杀怪物数量":
		return SubConditionType_HeroKillMonsterNum, true
	case "场上存在强化x级的英雄y个":
		return SubConditionType_HeroEnchantLevelNums, true
	case "召唤某英雄x次":
		return SubConditionType_SummonHeroTimes, true
	case "通关某关卡ID多少次":
		return SubConditionType_PassLevelIdTimes, true
	case "招募某英雄x次":
		return SubConditionType_RecruitHeroTimes, true
	case "出售英雄次数":
		return SubConditionType_HeroSellTimes, true
	case "试炼通关关卡":
		return SubConditionType_TrailLevelIdPass, true
	case "进行公会建设N次":
		return SubConditionType_AllianceBuildTimes, true
	case "进行公会科技抽奖N次":
		return SubConditionType_AllianceTechGacha, true
	case "八门合成":
		return SubConditionType_SoulCompose, true
	case "领地升级x次":
		return SubConditionType_ManorLevelupTimes, true
	case "领地达到x级":
		return SubConditionType_ManorLevelAchieve, true
	case "格斗家礼物开箱次数":
		return SubConditionType_BoxGiftOpenBox, true
	case "格斗家礼物箱子等级":
		return SubConditionType_BoxGiftBoxLv, true
	case "累积vip经验":
		return SubConditionType_TotalVipExp, true
	case "领地成功拉车X次":
		return SubConditionType_ManorGotRes, true
	case "穿戴X个本源":
		return SubConditionType_EquipHeroSoulTotal, true
	case "公会砍价X次":
		return SubConditionType_AllianceBargainTimes, true
	case "领取某类型任务奖励":
		return SubConditionType_QuestGainTypeAward, true
	case "拥有某个英雄":
		return SubConditionType_OwnHero, true
	case "领地掠夺X次":
		return SubConditionType_ManorPlunder, true
	case "神器档位召唤次数":
		return SubConditionType_ArtifactDrawGearTimes, true
	case "消耗英雄碎片":
		return SubConditionType_CostHeroChip, true
	case "公会科技抽到碎片品质":
		return SubConditionType_AllianceTechObtain, true
	case "特定英雄星级":
		return SubConditionType_SpecificHeroStar, true
	case "神器抽取消耗":
		return SubConditionType_ArtifactDrawCost, true
	case "跟随公会占领X个X级城市":
		return SubConditionType_FollowAllianceOccupyLevelCities, true
	case "阶段首次城市争夺参战":
		return SubConditionType_PhaseJoinCityWar, true
	case "创建角色满x天":
		return SubConditionType_CreatePlayerDays, true
	case "跨服32强排名次数":
		return SubConditionType_CrossArenaRankCnt, true
	case "赏金个人排名次数":
		return SubConditionType_BountyPersonRankCnt, true
	case "赏金公会排名次数":
		return SubConditionType_BountyAllianceRankCnt, true
	case "登录天数":
		return SubConditionType_LoginDays, true
	case "格斗圣殿排名":
		return SubConditionType_FightingRank, true
	case "名人堂锻造次数":
		return SubConditionType_FameForgeTimes, true
	case "活动开启x天":
		return SubConditionType_ActivityOpenDays, true
	case "获得阵营战商店代币数量":
		return SubConditionType_GainCampBattleToken, true
	case "绑定宿命格斗家":
		return SubConditionType_BindDestinyHero, true
	case "邀请x个玩家升到x级":
		return SubConditionType_InvitePlayerLevel, true
	case "连续充值活动积分":
		return SubConditionType_ContinuousScore, true
	case "自己参加别人翻牌聚会":
		return SubConditionType_MyselfJoinOtherFlipCard, true
	case "别人参加自己翻牌聚会":
		return SubConditionType_OtherJoinMyselfrFlipCard, true
	case "绑定宿命英雄达到x星":
		return SubConditionType_BindDestinyHeroStar, true
	case "激活权益卡":
		return SubConditionType_ActivatePrivilegeCard, true
	case "开启某个系统":
		return SubConditionType_SystemOpened, true
	case "绑定且拥有宿命英雄":
		return SubConditionType_BindAndHaveDestinyHero, true
	case "呆呆通关次数":
		return SubConditionType_TDPassTimes, true
	case "前五战力":
		return SubConditionType_TopFivePower, true
	case "拥有X个X品质格斗家":
		return SubConditionType_HaveRarityHeros, true
	case "拥有X个X品质神器":
		return SubConditionType_HaveRarityArtifacts, true
	case "拥有X个X星级格斗家":
		return SubConditionType_HaveStarHeros, true
	case "拥有X个X品质收藏手办":
		return SubConditionType_HaveRarityGift, true
	case "南镇击杀数":
		return SubConditionType_ElfMaxKill, true
	case "世界boss伤害":
		return SubConditionType_WorldBossDamage, true
	case "拼图格子点亮":
		return SubConditionType_ActPuzzleActive, true
	case "在X频道发言N次":
		return SubConditionType_ChannelChatTimes, true
	case "在多个商店购买N次":
		return SubConditionType_AnyOneOfShopBuyCnt, true
	case "开服天数":
		return SubConditionType_ServerOpenDays, true
	case "随心抽登录天数":
		return SubConditionType_UnlimitedDrawLoginDays, true
	case "拥有终身卡":
		return SubConditionType_HasLifeCard, true
	case "自选招募次数":
		return SubConditionType_WishHeroRecruit, true
	case "Up招募次数":
		return SubConditionType_UpHeroRecruit, true
	case "作为队长通关x次":
		return SubConditionType_PassAsTeamLeader, true
	case "组队关卡通关x次":
		return SubConditionType_PassTeamStage, true
	case "大蛇boss挑战伤害":
		return SubConditionType_RoyalBossChallengeDamage, true
	case "世界boss模式伤害":
		return SubConditionType_WorldBossModeDamage, true
	case "世界boss模式挑战":
		return SubConditionType_WorldBossModeAttack, true
	case "公会boss伤害":
		return SubConditionType_AllianceBossDamage, true
	case "格斗盛典胜利场次":
		return SubConditionType_FightingWinCnt, true
	case "格斗盛典挑战次数":
		return SubConditionType_FightingChallengeCnt, true
	case "呆呆通关多少波x次":
		return SubConditionType_TDPassXTimes, true
	case "组队关卡":
		return SubConditionType_TeamStage, true
	case "关注游戏圈":
		return SubConditionType_GameClubFollow, true
	case "游戏圈浏览时长":
		return SubConditionType_GameClubBrowsingSeconds, true
	case "游戏圈内点赞帖子":
		return SubConditionType_GameClubLikeCount, true
	case "评论官方帖子":
		return SubConditionType_GameClubCommentCount, true
	case "游戏圈发帖数":
		return SubConditionType_GameClubPost, true
	case "参与X次贸易":
		return SubConditionType_TradeJoinTimes, true
	case "贸易赚取X贸易币":
		return SubConditionType_TradeEarn, true
	case "贸易单局持有X贸易币":
		return SubConditionType_TradeAsset, true
	case "贸易获得前X名Y次":
		return SubConditionType_TradeRank, true
	case "商店任意商品购买次数":
		return SubConditionType_AnyShopItemBuyCnt, true
	case "活动任务领奖":
		return SubConditionType_AvyTaskGotAward, true
	case "开启时装次数":
		return SubConditionType_OpenFashionBox, true
	case "皇城个人积分改变":
		return SubConditionType_KingWarPlayerScoreChange, true
	case "皇城累计消耗精力":
		return SubConditionType_KingWarPlayerCostEnergy, true
	case "周年海岛探险":
		return SubConditionType_AniActIslandDraw, true
	case "宠物合成":
		return SubConditionType_PetMerge, true

	}
	return -1, false

}
