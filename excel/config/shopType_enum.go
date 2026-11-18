package config

// 商店商品表-拳皇/ShopType
type ShopId int

const (
	ShopId_BEGIN                 ShopId = 1     //枚举开始
	ShopId_CoinShop              ShopId = 1     //金币商店
	ShopId_DiamondShop           ShopId = 2     //钻石商店
	ShopId_CommonShop            ShopId = 3     //通用商店
	ShopId_GiftCard              ShopId = 4     //权益卡
	ShopId_LevelBattlePass       ShopId = 5     //通行证
	ShopId_BlackMarket           ShopId = 6     //黑市
	ShopId_DailySuperPackage     ShopId = 7     //日特惠
	ShopId_WeekSuperPackage      ShopId = 8     //周特惠
	ShopId_MonthSuperPackage     ShopId = 9     //月特惠
	ShopId_AccumulatedRecharge   ShopId = 10    //累计充值
	ShopId_SprintWelfare         ShopId = 11    //冲刺福利
	ShopId_AllianceShop          ShopId = 12    //联盟商店
	ShopId_CorpsRecruit          ShopId = 13    //佣兵招募
	ShopId_TerritoryResources    ShopId = 14    //领地资源
	ShopId_CoachRecruit          ShopId = 15    //教官招募
	ShopId_CoachUpgrade          ShopId = 16    //教官成长
	ShopId_RelicsUpgrade         ShopId = 17    //圣物成长
	ShopId_CorpLeader            ShopId = 18    //团长等级成长
	ShopId_BoxTimeTicket         ShopId = 19    //加速券
	ShopId_VipDailyShop          ShopId = 20    //Vip特权礼包商店
	ShopId_VipShop               ShopId = 21    //Vip商店
	ShopId_FirstRecharge         ShopId = 22    //首充
	ShopId_MallDaily             ShopId = 30    //每日精选
	ShopId_GuildShop             ShopId = 31    //公会商店
	ShopId_TDShop                ShopId = 32    //TD商店
	ShopId_ExpeditionShop        ShopId = 33    //远征商店
	ShopId_JiJinShop             ShopId = 34    //基金商城
	ShopId_CrossSportsShop       ShopId = 35    //跨服竞技商城
	ShopId_CrossArenaCheer       ShopId = 36    //跨服竞技助威礼包
	ShopId_CrossArenaCheer1      ShopId = 39    //跨服竞技备战礼包
	ShopId_MonthCardShop         ShopId = 43    //月卡商店
	ShopId_WeekCircleRecruit     ShopId = 44    //周循环招募商店
	ShopId_WeeklyCycleArtifact   ShopId = 45    //周循环神器商店
	ShopId_DayOptionalPackage    ShopId = 46    //每日自选礼包商店
	ShopId_WeeklyOptionalPackage ShopId = 47    //每周自选礼包
	ShopId_MonthOptionalPackage  ShopId = 4701  //每月自选礼包
	ShopId_HeroOptionalPackage   ShopId = 48    //每周英雄自选礼包
	ShopId_FishBaitShop          ShopId = 49    //钓鱼鱼饵商店
	ShopId_FishGalaShop          ShopId = 50    //鱼获宴会商店
	ShopId_WorkshopShop          ShopId = 51    //创意工坊商店
	ShopId_NewbieHeroRecruit     ShopId = 52    //新手招募商店
	ShopId_WeekCircleAlliTech    ShopId = 53    //周循环公会科技商店
	ShopId_SoulCompose           ShopId = 54    //八门合成
	ShopId_StarDiamond           ShopId = 56    //星钻商城
	ShopId_WorldBossShop         ShopId = 66    //世界boss商店
	ShopId_FishDailyShop         ShopId = 67    //钓鱼日常礼包
	ShopId_FishCardShop          ShopId = 68    //钓鱼权益卡商店
	ShopId_SpTDShop              ShopId = 71    //南镇商店
	ShopId_NewbieGiftPack        ShopId = 72    //新手活动.新手礼包
	ShopId_SpringExchange        ShopId = 75    //繁花物语_兑换商店
	ShopId_SpringDaily           ShopId = 76    //繁花物语_每日特惠
	ShopId_SpringShop            ShopId = 77    //繁花物语_活动特惠
	ShopId_SummerExchange        ShopId = 78    //盛夏庆典_兑换商店
	ShopId_SummerDaily           ShopId = 79    //盛夏庆典_每日特惠
	ShopId_SummerShop            ShopId = 80    //盛夏庆典_活动特惠
	ShopId_FlipCardSign          ShopId = 86    //拳皇翻翻乐-签到有礼
	ShopId_FlipCardGift          ShopId = 86001 //拳皇翻翻乐每日礼包
	ShopId_CampBattleShop        ShopId = 90    //阵营战商店
	ShopId_BecomeRichShop        ShopId = 92    //一夜暴富商店
	ShopId_FirstRechargeNew      ShopId = 93    //新版首充
	ShopId_Strongest108          ShopId = 95    //最强108商店
	ShopId_TeamLeague            ShopId = 96    //公会团赛礼包
	ShopId_Trade                 ShopId = 98    //贸易商店
	ShopId_KingWar               ShopId = 99    //皇城争锋商店
	ShopId_ManorShop             ShopId = 1006  //领地商店
	ShopId_ManorMeteorShop       ShopId = 1007  //领地陨石商店
	ShopId_RankShop              ShopId = 1028  //冲榜商店
	ShopId_END                   ShopId = 86002 //枚举结束
)

// 商店商品表-拳皇/ShopType
type ShopGoodsIs int

const (
	ShopGoodsIs_BEGIN                ShopGoodsIs = 1        //枚举开始
	ShopGoodsIs_FirstRechargeGoodsId ShopGoodsIs = 40100901 //首充商品id
	ShopGoodsIs_END                  ShopGoodsIs = 40100902 //枚举结束
)

// 商店商品表-拳皇/ShopType
type PackageId int

const (
	PackageId_BEGIN                 PackageId = 1        //枚举开始
	PackageId_DailyFree             PackageId = 40100300 //每日礼包免费
	PackageId_DailyAll              PackageId = 40100305 //每日礼包一件打包
	PackageId_MonthCard             PackageId = 40100400 //月卡
	PackageId_ForeverCard           PackageId = 40100401 //终身卡
	PackageId_WeekHeroCard          PackageId = 40100403 //英雄特惠周卡
	PackageId_WeekFishCard          PackageId = 40100404 //钓鱼特惠周卡
	PackageId_WeekTreasureCard      PackageId = 40100405 //宝物特惠周卡
	PackageId_EquityCard            PackageId = 40100406 //畅玩卡
	PackageId_HightPassCard         PackageId = 40100402 //高级通行证
	PackageId_FirstRecharge         PackageId = 40100901 //首充礼包
	PackageId_FirstRecharge1        PackageId = 40100902 //首充第二档
	PackageId_FirstRecharge2        PackageId = 40100903 //首充第三档
	PackageId_FirstRecharge3        PackageId = 40100904 //首充第四档
	PackageId_FirstRechargeNew      PackageId = 93100001 //新首充礼包
	PackageId_FirstRechargeNew1     PackageId = 93100002 //新首充第二档
	PackageId_FirstRechargeNew2     PackageId = 93100003 //新首充第三档
	PackageId_FirstRechargeNew3     PackageId = 93100004 //新首充第四档
	PackageId_GrowUp1               PackageId = 40101001 //成长基金1
	PackageId_GrowUp2               PackageId = 40101002 //成长基金2
	PackageId_Challenge1            PackageId = 40101003 //挑战基金1
	PackageId_Challenge2            PackageId = 40101004 //挑战基金2
	PackageId_Honor1                PackageId = 40101005 //荣誉基金1
	PackageId_Honor2                PackageId = 40101006 //荣誉基金2
	PackageId_ArtifactPass          PackageId = 40101007 //神器基金
	PackageId_CollectPass           PackageId = 40101008 //收藏基金
	PackageId_NestPass              PackageId = 40101009 //音巢基金
	PackageId_RecruitPass           PackageId = 40101010 //招募基金
	PackageId_HeroStarPass          PackageId = 40101011 //格斗家等级基金
	PackageId_SecretArtsPass        PackageId = 40101012 //秘术基金
	PackageId_HeroSoulPass          PackageId = 40101013 //八门基金
	PackageId_HeroSelectPack        PackageId = 48100304 //英雄自选礼包打包id
	PackageId_HeroSelectFreeID      PackageId = 48100300 //英雄自选礼包免费id
	PackageId_CrossArenaTicket      PackageId = 40100104 //跨服竞技场门票
	PackageId_AutoSelectPack        PackageId = 47000001 //每周自选礼包免费id
	PackageId_RecruitTicket         PackageId = 40100103 //高级招募券
	PackageId_FishCardID            PackageId = 40100407 //钓鱼权益卡id
	PackageId_UpRecruitTicket       PackageId = 40100110 //Up招募券
	PackageId_WishRecruitTicket     PackageId = 40100111 //心愿招募券
	PackageId_FlipCardSign          PackageId = 8600001  //拳皇翻翻乐-签到有礼
	PackageId_FlipCardBoxTicket     PackageId = 40100113 //拳皇翻翻乐翻卡券
	PackageId_OneBuy                PackageId = 40100408 //一元购买
	PackageId_RecruitActSign        PackageId = 92000000 //招募活动签到有礼
	PackageId_RecruitPrivilageCard  PackageId = 40100114 //招募冲刺特权卡
	PackageId_ArtifactPrivilageCard PackageId = 40100115 //神器冲刺特权卡
	PackageId_MiShuPrivilageCard    PackageId = 40100116 //秘术冲刺特权卡
	PackageId_SoulCard              PackageId = 40100409 //战魂特权卡
	PackageId_END                   PackageId = 93100005 //枚举结束
)

// 商店商品表-拳皇/ShopType
type EMShopType int

const (
	EMShopType_BEGIN    EMShopType = 1 //枚举开始
	EMShopType_PlayShop EMShopType = 1 //玩法商城
	EMShopType_END      EMShopType = 2 //枚举结束
)

func (e ShopId) String() string {

	if e == ShopId_CoinShop {
		return "CoinShop"
	}
	if e == ShopId_DiamondShop {
		return "DiamondShop"
	}
	if e == ShopId_CommonShop {
		return "CommonShop"
	}
	if e == ShopId_GiftCard {
		return "GiftCard"
	}
	if e == ShopId_LevelBattlePass {
		return "LevelBattlePass"
	}
	if e == ShopId_BlackMarket {
		return "BlackMarket"
	}
	if e == ShopId_DailySuperPackage {
		return "DailySuperPackage"
	}
	if e == ShopId_WeekSuperPackage {
		return "WeekSuperPackage"
	}
	if e == ShopId_MonthSuperPackage {
		return "MonthSuperPackage"
	}
	if e == ShopId_AccumulatedRecharge {
		return "AccumulatedRecharge"
	}
	if e == ShopId_SprintWelfare {
		return "SprintWelfare"
	}
	if e == ShopId_AllianceShop {
		return "AllianceShop"
	}
	if e == ShopId_CorpsRecruit {
		return "CorpsRecruit"
	}
	if e == ShopId_TerritoryResources {
		return "TerritoryResources"
	}
	if e == ShopId_CoachRecruit {
		return "CoachRecruit"
	}
	if e == ShopId_CoachUpgrade {
		return "CoachUpgrade"
	}
	if e == ShopId_RelicsUpgrade {
		return "RelicsUpgrade"
	}
	if e == ShopId_CorpLeader {
		return "CorpLeader"
	}
	if e == ShopId_BoxTimeTicket {
		return "BoxTimeTicket"
	}
	if e == ShopId_VipDailyShop {
		return "VipDailyShop"
	}
	if e == ShopId_VipShop {
		return "VipShop"
	}
	if e == ShopId_FirstRecharge {
		return "FirstRecharge"
	}
	if e == ShopId_MallDaily {
		return "MallDaily"
	}
	if e == ShopId_GuildShop {
		return "GuildShop"
	}
	if e == ShopId_TDShop {
		return "TDShop"
	}
	if e == ShopId_ExpeditionShop {
		return "ExpeditionShop"
	}
	if e == ShopId_JiJinShop {
		return "JiJinShop"
	}
	if e == ShopId_CrossSportsShop {
		return "CrossSportsShop"
	}
	if e == ShopId_CrossArenaCheer {
		return "CrossArenaCheer"
	}
	if e == ShopId_CrossArenaCheer1 {
		return "CrossArenaCheer1"
	}
	if e == ShopId_MonthCardShop {
		return "MonthCardShop"
	}
	if e == ShopId_WeekCircleRecruit {
		return "WeekCircleRecruit"
	}
	if e == ShopId_WeeklyCycleArtifact {
		return "WeeklyCycleArtifact"
	}
	if e == ShopId_DayOptionalPackage {
		return "DayOptionalPackage"
	}
	if e == ShopId_WeeklyOptionalPackage {
		return "WeeklyOptionalPackage"
	}
	if e == ShopId_MonthOptionalPackage {
		return "MonthOptionalPackage"
	}
	if e == ShopId_HeroOptionalPackage {
		return "HeroOptionalPackage"
	}
	if e == ShopId_FishBaitShop {
		return "FishBaitShop"
	}
	if e == ShopId_FishGalaShop {
		return "FishGalaShop"
	}
	if e == ShopId_WorkshopShop {
		return "WorkshopShop"
	}
	if e == ShopId_NewbieHeroRecruit {
		return "NewbieHeroRecruit"
	}
	if e == ShopId_WeekCircleAlliTech {
		return "WeekCircleAlliTech"
	}
	if e == ShopId_SoulCompose {
		return "SoulCompose"
	}
	if e == ShopId_StarDiamond {
		return "StarDiamond"
	}
	if e == ShopId_WorldBossShop {
		return "WorldBossShop"
	}
	if e == ShopId_FishDailyShop {
		return "FishDailyShop"
	}
	if e == ShopId_FishCardShop {
		return "FishCardShop"
	}
	if e == ShopId_SpTDShop {
		return "SpTDShop"
	}
	if e == ShopId_NewbieGiftPack {
		return "NewbieGiftPack"
	}
	if e == ShopId_SpringExchange {
		return "SpringExchange"
	}
	if e == ShopId_SpringDaily {
		return "SpringDaily"
	}
	if e == ShopId_SpringShop {
		return "SpringShop"
	}
	if e == ShopId_SummerExchange {
		return "SummerExchange"
	}
	if e == ShopId_SummerDaily {
		return "SummerDaily"
	}
	if e == ShopId_SummerShop {
		return "SummerShop"
	}
	if e == ShopId_FlipCardSign {
		return "FlipCardSign"
	}
	if e == ShopId_FlipCardGift {
		return "FlipCardGift"
	}
	if e == ShopId_CampBattleShop {
		return "CampBattleShop"
	}
	if e == ShopId_BecomeRichShop {
		return "BecomeRichShop"
	}
	if e == ShopId_FirstRechargeNew {
		return "FirstRechargeNew"
	}
	if e == ShopId_Strongest108 {
		return "Strongest108"
	}
	if e == ShopId_TeamLeague {
		return "TeamLeague"
	}
	if e == ShopId_Trade {
		return "Trade"
	}
	if e == ShopId_KingWar {
		return "KingWar"
	}
	if e == ShopId_ManorShop {
		return "ManorShop"
	}
	if e == ShopId_ManorMeteorShop {
		return "ManorMeteorShop"
	}
	if e == ShopId_RankShop {
		return "RankShop"
	}

	return ""

}

func GetShopId(tp string) (ShopId, bool) {
	switch tp {
	case "CoinShop":
		return ShopId_CoinShop, true
	case "DiamondShop":
		return ShopId_DiamondShop, true
	case "CommonShop":
		return ShopId_CommonShop, true
	case "GiftCard":
		return ShopId_GiftCard, true
	case "LevelBattlePass":
		return ShopId_LevelBattlePass, true
	case "BlackMarket":
		return ShopId_BlackMarket, true
	case "DailySuperPackage":
		return ShopId_DailySuperPackage, true
	case "WeekSuperPackage":
		return ShopId_WeekSuperPackage, true
	case "MonthSuperPackage":
		return ShopId_MonthSuperPackage, true
	case "AccumulatedRecharge":
		return ShopId_AccumulatedRecharge, true
	case "SprintWelfare":
		return ShopId_SprintWelfare, true
	case "AllianceShop":
		return ShopId_AllianceShop, true
	case "CorpsRecruit":
		return ShopId_CorpsRecruit, true
	case "TerritoryResources":
		return ShopId_TerritoryResources, true
	case "CoachRecruit":
		return ShopId_CoachRecruit, true
	case "CoachUpgrade":
		return ShopId_CoachUpgrade, true
	case "RelicsUpgrade":
		return ShopId_RelicsUpgrade, true
	case "CorpLeader":
		return ShopId_CorpLeader, true
	case "BoxTimeTicket":
		return ShopId_BoxTimeTicket, true
	case "VipDailyShop":
		return ShopId_VipDailyShop, true
	case "VipShop":
		return ShopId_VipShop, true
	case "FirstRecharge":
		return ShopId_FirstRecharge, true
	case "MallDaily":
		return ShopId_MallDaily, true
	case "GuildShop":
		return ShopId_GuildShop, true
	case "TDShop":
		return ShopId_TDShop, true
	case "ExpeditionShop":
		return ShopId_ExpeditionShop, true
	case "JiJinShop":
		return ShopId_JiJinShop, true
	case "CrossSportsShop":
		return ShopId_CrossSportsShop, true
	case "CrossArenaCheer":
		return ShopId_CrossArenaCheer, true
	case "CrossArenaCheer1":
		return ShopId_CrossArenaCheer1, true
	case "MonthCardShop":
		return ShopId_MonthCardShop, true
	case "WeekCircleRecruit":
		return ShopId_WeekCircleRecruit, true
	case "WeeklyCycleArtifact":
		return ShopId_WeeklyCycleArtifact, true
	case "DayOptionalPackage":
		return ShopId_DayOptionalPackage, true
	case "WeeklyOptionalPackage":
		return ShopId_WeeklyOptionalPackage, true
	case "MonthOptionalPackage":
		return ShopId_MonthOptionalPackage, true
	case "HeroOptionalPackage":
		return ShopId_HeroOptionalPackage, true
	case "FishBaitShop":
		return ShopId_FishBaitShop, true
	case "FishGalaShop":
		return ShopId_FishGalaShop, true
	case "WorkshopShop":
		return ShopId_WorkshopShop, true
	case "NewbieHeroRecruit":
		return ShopId_NewbieHeroRecruit, true
	case "WeekCircleAlliTech":
		return ShopId_WeekCircleAlliTech, true
	case "SoulCompose":
		return ShopId_SoulCompose, true
	case "StarDiamond":
		return ShopId_StarDiamond, true
	case "WorldBossShop":
		return ShopId_WorldBossShop, true
	case "FishDailyShop":
		return ShopId_FishDailyShop, true
	case "FishCardShop":
		return ShopId_FishCardShop, true
	case "SpTDShop":
		return ShopId_SpTDShop, true
	case "NewbieGiftPack":
		return ShopId_NewbieGiftPack, true
	case "SpringExchange":
		return ShopId_SpringExchange, true
	case "SpringDaily":
		return ShopId_SpringDaily, true
	case "SpringShop":
		return ShopId_SpringShop, true
	case "SummerExchange":
		return ShopId_SummerExchange, true
	case "SummerDaily":
		return ShopId_SummerDaily, true
	case "SummerShop":
		return ShopId_SummerShop, true
	case "FlipCardSign":
		return ShopId_FlipCardSign, true
	case "FlipCardGift":
		return ShopId_FlipCardGift, true
	case "CampBattleShop":
		return ShopId_CampBattleShop, true
	case "BecomeRichShop":
		return ShopId_BecomeRichShop, true
	case "FirstRechargeNew":
		return ShopId_FirstRechargeNew, true
	case "Strongest108":
		return ShopId_Strongest108, true
	case "TeamLeague":
		return ShopId_TeamLeague, true
	case "Trade":
		return ShopId_Trade, true
	case "KingWar":
		return ShopId_KingWar, true
	case "ManorShop":
		return ShopId_ManorShop, true
	case "ManorMeteorShop":
		return ShopId_ManorMeteorShop, true
	case "RankShop":
		return ShopId_RankShop, true

	}
	return -1, false

}

func GetShopId_ByAlias(tp string) (ShopId, bool) {
	switch tp {
	case "金币商店":
		return ShopId_CoinShop, true
	case "钻石商店":
		return ShopId_DiamondShop, true
	case "通用商店":
		return ShopId_CommonShop, true
	case "权益卡":
		return ShopId_GiftCard, true
	case "通行证":
		return ShopId_LevelBattlePass, true
	case "黑市":
		return ShopId_BlackMarket, true
	case "日特惠":
		return ShopId_DailySuperPackage, true
	case "周特惠":
		return ShopId_WeekSuperPackage, true
	case "月特惠":
		return ShopId_MonthSuperPackage, true
	case "累计充值":
		return ShopId_AccumulatedRecharge, true
	case "冲刺福利":
		return ShopId_SprintWelfare, true
	case "联盟商店":
		return ShopId_AllianceShop, true
	case "佣兵招募":
		return ShopId_CorpsRecruit, true
	case "领地资源":
		return ShopId_TerritoryResources, true
	case "教官招募":
		return ShopId_CoachRecruit, true
	case "教官成长":
		return ShopId_CoachUpgrade, true
	case "圣物成长":
		return ShopId_RelicsUpgrade, true
	case "团长等级成长":
		return ShopId_CorpLeader, true
	case "加速券":
		return ShopId_BoxTimeTicket, true
	case "Vip特权礼包商店":
		return ShopId_VipDailyShop, true
	case "Vip商店":
		return ShopId_VipShop, true
	case "首充":
		return ShopId_FirstRecharge, true
	case "每日精选":
		return ShopId_MallDaily, true
	case "公会商店":
		return ShopId_GuildShop, true
	case "TD商店":
		return ShopId_TDShop, true
	case "远征商店":
		return ShopId_ExpeditionShop, true
	case "基金商城":
		return ShopId_JiJinShop, true
	case "跨服竞技商城":
		return ShopId_CrossSportsShop, true
	case "跨服竞技助威礼包":
		return ShopId_CrossArenaCheer, true
	case "跨服竞技备战礼包":
		return ShopId_CrossArenaCheer1, true
	case "月卡商店":
		return ShopId_MonthCardShop, true
	case "周循环招募商店":
		return ShopId_WeekCircleRecruit, true
	case "周循环神器商店":
		return ShopId_WeeklyCycleArtifact, true
	case "每日自选礼包商店":
		return ShopId_DayOptionalPackage, true
	case "每周自选礼包":
		return ShopId_WeeklyOptionalPackage, true
	case "每月自选礼包":
		return ShopId_MonthOptionalPackage, true
	case "每周英雄自选礼包":
		return ShopId_HeroOptionalPackage, true
	case "钓鱼鱼饵商店":
		return ShopId_FishBaitShop, true
	case "鱼获宴会商店":
		return ShopId_FishGalaShop, true
	case "创意工坊商店":
		return ShopId_WorkshopShop, true
	case "新手招募商店":
		return ShopId_NewbieHeroRecruit, true
	case "周循环公会科技商店":
		return ShopId_WeekCircleAlliTech, true
	case "八门合成":
		return ShopId_SoulCompose, true
	case "星钻商城":
		return ShopId_StarDiamond, true
	case "世界boss商店":
		return ShopId_WorldBossShop, true
	case "钓鱼日常礼包":
		return ShopId_FishDailyShop, true
	case "钓鱼权益卡商店":
		return ShopId_FishCardShop, true
	case "南镇商店":
		return ShopId_SpTDShop, true
	case "新手活动.新手礼包":
		return ShopId_NewbieGiftPack, true
	case "繁花物语_兑换商店":
		return ShopId_SpringExchange, true
	case "繁花物语_每日特惠":
		return ShopId_SpringDaily, true
	case "繁花物语_活动特惠":
		return ShopId_SpringShop, true
	case "盛夏庆典_兑换商店":
		return ShopId_SummerExchange, true
	case "盛夏庆典_每日特惠":
		return ShopId_SummerDaily, true
	case "盛夏庆典_活动特惠":
		return ShopId_SummerShop, true
	case "拳皇翻翻乐-签到有礼":
		return ShopId_FlipCardSign, true
	case "拳皇翻翻乐每日礼包":
		return ShopId_FlipCardGift, true
	case "阵营战商店":
		return ShopId_CampBattleShop, true
	case "一夜暴富商店":
		return ShopId_BecomeRichShop, true
	case "新版首充":
		return ShopId_FirstRechargeNew, true
	case "最强108商店":
		return ShopId_Strongest108, true
	case "公会团赛礼包":
		return ShopId_TeamLeague, true
	case "贸易商店":
		return ShopId_Trade, true
	case "皇城争锋商店":
		return ShopId_KingWar, true
	case "领地商店":
		return ShopId_ManorShop, true
	case "领地陨石商店":
		return ShopId_ManorMeteorShop, true
	case "冲榜商店":
		return ShopId_RankShop, true

	}
	return -1, false

}

func (e ShopGoodsIs) String() string {

	if e == ShopGoodsIs_FirstRechargeGoodsId {
		return "FirstRechargeGoodsId"
	}

	return ""

}

func GetShopGoodsIs(tp string) (ShopGoodsIs, bool) {
	switch tp {
	case "FirstRechargeGoodsId":
		return ShopGoodsIs_FirstRechargeGoodsId, true

	}
	return -1, false

}

func GetShopGoodsIs_ByAlias(tp string) (ShopGoodsIs, bool) {
	switch tp {
	case "首充商品id":
		return ShopGoodsIs_FirstRechargeGoodsId, true

	}
	return -1, false

}

func (e PackageId) String() string {

	if e == PackageId_DailyFree {
		return "DailyFree"
	}
	if e == PackageId_DailyAll {
		return "DailyAll"
	}
	if e == PackageId_MonthCard {
		return "MonthCard"
	}
	if e == PackageId_ForeverCard {
		return "ForeverCard"
	}
	if e == PackageId_WeekHeroCard {
		return "WeekHeroCard"
	}
	if e == PackageId_WeekFishCard {
		return "WeekFishCard"
	}
	if e == PackageId_WeekTreasureCard {
		return "WeekTreasureCard"
	}
	if e == PackageId_EquityCard {
		return "EquityCard"
	}
	if e == PackageId_HightPassCard {
		return "HightPassCard"
	}
	if e == PackageId_FirstRecharge {
		return "FirstRecharge"
	}
	if e == PackageId_FirstRecharge1 {
		return "FirstRecharge1"
	}
	if e == PackageId_FirstRecharge2 {
		return "FirstRecharge2"
	}
	if e == PackageId_FirstRecharge3 {
		return "FirstRecharge3"
	}
	if e == PackageId_FirstRechargeNew {
		return "FirstRechargeNew"
	}
	if e == PackageId_FirstRechargeNew1 {
		return "FirstRechargeNew1"
	}
	if e == PackageId_FirstRechargeNew2 {
		return "FirstRechargeNew2"
	}
	if e == PackageId_FirstRechargeNew3 {
		return "FirstRechargeNew3"
	}
	if e == PackageId_GrowUp1 {
		return "GrowUp1"
	}
	if e == PackageId_GrowUp2 {
		return "GrowUp2"
	}
	if e == PackageId_Challenge1 {
		return "Challenge1"
	}
	if e == PackageId_Challenge2 {
		return "Challenge2"
	}
	if e == PackageId_Honor1 {
		return "Honor1"
	}
	if e == PackageId_Honor2 {
		return "Honor2"
	}
	if e == PackageId_ArtifactPass {
		return "ArtifactPass"
	}
	if e == PackageId_CollectPass {
		return "CollectPass"
	}
	if e == PackageId_NestPass {
		return "NestPass"
	}
	if e == PackageId_RecruitPass {
		return "RecruitPass"
	}
	if e == PackageId_HeroStarPass {
		return "HeroStarPass"
	}
	if e == PackageId_SecretArtsPass {
		return "SecretArtsPass"
	}
	if e == PackageId_HeroSoulPass {
		return "HeroSoulPass"
	}
	if e == PackageId_HeroSelectPack {
		return "HeroSelectPack"
	}
	if e == PackageId_HeroSelectFreeID {
		return "HeroSelectFreeID"
	}
	if e == PackageId_CrossArenaTicket {
		return "CrossArenaTicket"
	}
	if e == PackageId_AutoSelectPack {
		return "AutoSelectPack"
	}
	if e == PackageId_RecruitTicket {
		return "RecruitTicket"
	}
	if e == PackageId_FishCardID {
		return "FishCardID"
	}
	if e == PackageId_UpRecruitTicket {
		return "UpRecruitTicket"
	}
	if e == PackageId_WishRecruitTicket {
		return "WishRecruitTicket"
	}
	if e == PackageId_FlipCardSign {
		return "FlipCardSign"
	}
	if e == PackageId_FlipCardBoxTicket {
		return "FlipCardBoxTicket"
	}
	if e == PackageId_OneBuy {
		return "OneBuy"
	}
	if e == PackageId_RecruitActSign {
		return "RecruitActSign"
	}
	if e == PackageId_RecruitPrivilageCard {
		return "RecruitPrivilageCard"
	}
	if e == PackageId_ArtifactPrivilageCard {
		return "ArtifactPrivilageCard"
	}
	if e == PackageId_MiShuPrivilageCard {
		return "MiShuPrivilageCard"
	}
	if e == PackageId_SoulCard {
		return "SoulCard"
	}

	return ""

}

func GetPackageId(tp string) (PackageId, bool) {
	switch tp {
	case "DailyFree":
		return PackageId_DailyFree, true
	case "DailyAll":
		return PackageId_DailyAll, true
	case "MonthCard":
		return PackageId_MonthCard, true
	case "ForeverCard":
		return PackageId_ForeverCard, true
	case "WeekHeroCard":
		return PackageId_WeekHeroCard, true
	case "WeekFishCard":
		return PackageId_WeekFishCard, true
	case "WeekTreasureCard":
		return PackageId_WeekTreasureCard, true
	case "EquityCard":
		return PackageId_EquityCard, true
	case "HightPassCard":
		return PackageId_HightPassCard, true
	case "FirstRecharge":
		return PackageId_FirstRecharge, true
	case "FirstRecharge1":
		return PackageId_FirstRecharge1, true
	case "FirstRecharge2":
		return PackageId_FirstRecharge2, true
	case "FirstRecharge3":
		return PackageId_FirstRecharge3, true
	case "FirstRechargeNew":
		return PackageId_FirstRechargeNew, true
	case "FirstRechargeNew1":
		return PackageId_FirstRechargeNew1, true
	case "FirstRechargeNew2":
		return PackageId_FirstRechargeNew2, true
	case "FirstRechargeNew3":
		return PackageId_FirstRechargeNew3, true
	case "GrowUp1":
		return PackageId_GrowUp1, true
	case "GrowUp2":
		return PackageId_GrowUp2, true
	case "Challenge1":
		return PackageId_Challenge1, true
	case "Challenge2":
		return PackageId_Challenge2, true
	case "Honor1":
		return PackageId_Honor1, true
	case "Honor2":
		return PackageId_Honor2, true
	case "ArtifactPass":
		return PackageId_ArtifactPass, true
	case "CollectPass":
		return PackageId_CollectPass, true
	case "NestPass":
		return PackageId_NestPass, true
	case "RecruitPass":
		return PackageId_RecruitPass, true
	case "HeroStarPass":
		return PackageId_HeroStarPass, true
	case "SecretArtsPass":
		return PackageId_SecretArtsPass, true
	case "HeroSoulPass":
		return PackageId_HeroSoulPass, true
	case "HeroSelectPack":
		return PackageId_HeroSelectPack, true
	case "HeroSelectFreeID":
		return PackageId_HeroSelectFreeID, true
	case "CrossArenaTicket":
		return PackageId_CrossArenaTicket, true
	case "AutoSelectPack":
		return PackageId_AutoSelectPack, true
	case "RecruitTicket":
		return PackageId_RecruitTicket, true
	case "FishCardID":
		return PackageId_FishCardID, true
	case "UpRecruitTicket":
		return PackageId_UpRecruitTicket, true
	case "WishRecruitTicket":
		return PackageId_WishRecruitTicket, true
	case "FlipCardSign":
		return PackageId_FlipCardSign, true
	case "FlipCardBoxTicket":
		return PackageId_FlipCardBoxTicket, true
	case "OneBuy":
		return PackageId_OneBuy, true
	case "RecruitActSign":
		return PackageId_RecruitActSign, true
	case "RecruitPrivilageCard":
		return PackageId_RecruitPrivilageCard, true
	case "ArtifactPrivilageCard":
		return PackageId_ArtifactPrivilageCard, true
	case "MiShuPrivilageCard":
		return PackageId_MiShuPrivilageCard, true
	case "SoulCard":
		return PackageId_SoulCard, true

	}
	return -1, false

}

func GetPackageId_ByAlias(tp string) (PackageId, bool) {
	switch tp {
	case "每日礼包免费":
		return PackageId_DailyFree, true
	case "每日礼包一件打包":
		return PackageId_DailyAll, true
	case "月卡":
		return PackageId_MonthCard, true
	case "终身卡":
		return PackageId_ForeverCard, true
	case "英雄特惠周卡":
		return PackageId_WeekHeroCard, true
	case "钓鱼特惠周卡":
		return PackageId_WeekFishCard, true
	case "宝物特惠周卡":
		return PackageId_WeekTreasureCard, true
	case "畅玩卡":
		return PackageId_EquityCard, true
	case "高级通行证":
		return PackageId_HightPassCard, true
	case "首充礼包":
		return PackageId_FirstRecharge, true
	case "首充第二档":
		return PackageId_FirstRecharge1, true
	case "首充第三档":
		return PackageId_FirstRecharge2, true
	case "首充第四档":
		return PackageId_FirstRecharge3, true
	case "新首充礼包":
		return PackageId_FirstRechargeNew, true
	case "新首充第二档":
		return PackageId_FirstRechargeNew1, true
	case "新首充第三档":
		return PackageId_FirstRechargeNew2, true
	case "新首充第四档":
		return PackageId_FirstRechargeNew3, true
	case "成长基金1":
		return PackageId_GrowUp1, true
	case "成长基金2":
		return PackageId_GrowUp2, true
	case "挑战基金1":
		return PackageId_Challenge1, true
	case "挑战基金2":
		return PackageId_Challenge2, true
	case "荣誉基金1":
		return PackageId_Honor1, true
	case "荣誉基金2":
		return PackageId_Honor2, true
	case "神器基金":
		return PackageId_ArtifactPass, true
	case "收藏基金":
		return PackageId_CollectPass, true
	case "音巢基金":
		return PackageId_NestPass, true
	case "招募基金":
		return PackageId_RecruitPass, true
	case "格斗家等级基金":
		return PackageId_HeroStarPass, true
	case "秘术基金":
		return PackageId_SecretArtsPass, true
	case "八门基金":
		return PackageId_HeroSoulPass, true
	case "英雄自选礼包打包id":
		return PackageId_HeroSelectPack, true
	case "英雄自选礼包免费id":
		return PackageId_HeroSelectFreeID, true
	case "跨服竞技场门票":
		return PackageId_CrossArenaTicket, true
	case "每周自选礼包免费id":
		return PackageId_AutoSelectPack, true
	case "高级招募券":
		return PackageId_RecruitTicket, true
	case "钓鱼权益卡id":
		return PackageId_FishCardID, true
	case "Up招募券":
		return PackageId_UpRecruitTicket, true
	case "心愿招募券":
		return PackageId_WishRecruitTicket, true
	case "拳皇翻翻乐-签到有礼":
		return PackageId_FlipCardSign, true
	case "拳皇翻翻乐翻卡券":
		return PackageId_FlipCardBoxTicket, true
	case "一元购买":
		return PackageId_OneBuy, true
	case "招募活动签到有礼":
		return PackageId_RecruitActSign, true
	case "招募冲刺特权卡":
		return PackageId_RecruitPrivilageCard, true
	case "神器冲刺特权卡":
		return PackageId_ArtifactPrivilageCard, true
	case "秘术冲刺特权卡":
		return PackageId_MiShuPrivilageCard, true
	case "战魂特权卡":
		return PackageId_SoulCard, true

	}
	return -1, false

}

func (e EMShopType) String() string {

	if e == EMShopType_PlayShop {
		return "PlayShop"
	}

	return ""

}

func GetEMShopType(tp string) (EMShopType, bool) {
	switch tp {
	case "PlayShop":
		return EMShopType_PlayShop, true

	}
	return -1, false

}

func GetEMShopType_ByAlias(tp string) (EMShopType, bool) {
	switch tp {
	case "玩法商城":
		return EMShopType_PlayShop, true

	}
	return -1, false

}
