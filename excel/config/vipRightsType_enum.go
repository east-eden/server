package config

// vip表-拳皇/VipRightsType
type VipRightsType int

const (
	VipRightsType_BEGIN                       VipRightsType = 1  //枚举开始
	VipRightsType_OfflineGoldAdd              VipRightsType = 1  //离线，金币道具奖励+%，只影响领取全部
	VipRightsType_OfflineHeroExpAdd           VipRightsType = 2  //离线，英雄经验道具掉落加成+%，只影响领取全部
	VipRightsType_OfflineTimeLimitAdd         VipRightsType = 3  //离线，上限时间+X小时
	VipRightsType_OfflineQuickGetTimesAdd     VipRightsType = 4  //离线，每日快速巡逻次数+X
	VipRightsType_ArenaTicketsLimitAdd        VipRightsType = 5  //竞技场，挑战券存储上限+X
	VipRightsType_LineupHeroAttackAdd         VipRightsType = 6  //全体上阵英雄攻击+%
	VipRightsType_LineupHeroDefenceAdd        VipRightsType = 7  //全体上阵英雄防御+%
	VipRightsType_LineupHeroHPAdd             VipRightsType = 8  //全体上阵英雄血量+%
	VipRightsType_ArtsTrialChallengeTimesAdd  VipRightsType = 9  //武道试炼，3个副本每日免费挑战次数+X
	VipRightsType_ArtifactDrawAdd             VipRightsType = 10 //神器召唤最后一档可以用钻石召唤
	VipRightsType_TeamExpAdd                  VipRightsType = 11 //战队经验获取加成+%
	VipRightsType_ArenaRefreshCdSub           VipRightsType = 12 //竞技场刷新cd减少（秒）
	VipRightsType_ArenaShieldDurationAdd      VipRightsType = 13 //终身卡竞技场保护时间增加（秒）
	VipRightsType_SkipAllAds                  VipRightsType = 14 //免除所有广告
	VipRightsType_FameExpAdd                  VipRightsType = 15 //荣誉每次升级经验额外+X
	VipRightsType_OfflineQuickGetFreeTimesAdd VipRightsType = 16 //离线，每日免费巡逻次数+X
	VipRightsType_OfflineQuickGetPayTimesAdd  VipRightsType = 17 //离线，每日付费巡逻次数+X
	VipRightsType_RecruitProbUp               VipRightsType = 18 //招募心愿单英雄起始几率提升
	VipRightsType_AutoNextLevel               VipRightsType = 19 //自动挑战下一关
	VipRightsType_LDRefreshFormation          VipRightsType = 20 //LD阵容视频刷新次数+X
	VipRightsType_LDSpecialCard               VipRightsType = 21 //LD特殊卡视频获取次数+X
	VipRightsType_LifeCardSpeedUpOpenBox      VipRightsType = 22 //开盒加速
	VipRightsType_FishCardDelegate            VipRightsType = 23 //钓鱼委托
	VipRightsType_FishRodPropAdd              VipRightsType = 24 //鱼竿属性+%
	VipRightsType_FishHookPropAdd             VipRightsType = 25 //鱼钩属性+%
	VipRightsType_FishWirePropAdd             VipRightsType = 26 //鱼线属性+X
	VipRightsType_FishFloatPropAdd            VipRightsType = 27 //鱼漂属性+%
	VipRightsType_FishSpeedAdd                VipRightsType = 28 //钓鱼速度+%
	VipRightsType_FishSellRubbishAdd          VipRightsType = 29 //出售杂物银币增加+%
	VipRightsType_HeroSoulUnlock              VipRightsType = 30 //八门提前解锁
	VipRightsType_ManorButler                 VipRightsType = 31 //领地管家
	VipRightsType_SpTDSweep                   VipRightsType = 32 //南镇扫荡
	VipRightsType_PatrolShare                 VipRightsType = 33 //探秘分享
	VipRightsType_Send                        VipRightsType = 34 //赠送
	VipRightsType_RainbowSoulComposeProbAdd   VipRightsType = 35 //永恒八门合成概率+%
	VipRightsType_SkipBattle                  VipRightsType = 39 //跳过战斗
	VipRightsType_ExtraBoxGiftLevel           VipRightsType = 40 //额外祝福等级
	VipRightsType_ManorCheer                  VipRightsType = 41 //领地助威
	VipRightsType_TerritoryOutput             VipRightsType = 42 //领地产出
	VipRightsType_TeamShare                   VipRightsType = 43 //组队分享
	VipRightsType_ChatGroupAnnounce           VipRightsType = 44 //交流群发布
	VipRightsType_TeamTicketAdd               VipRightsType = 45 //组队每日额外领取门票次数
	VipRightsType_RoyalBossAdd                VipRightsType = 46 //大蛇模拟战每日额外领奖次数
	VipRightsType_SendFashion                 VipRightsType = 47 //装扮赠送
	VipRightsType_MillionaireReduceCoin       VipRightsType = 48 //百万大富翁装备消耗减少
	VipRightsType_END                         VipRightsType = 49 //枚举结束
)

func (e VipRightsType) String() string {

	if e == VipRightsType_OfflineGoldAdd {
		return "OfflineGoldAdd"
	}
	if e == VipRightsType_OfflineHeroExpAdd {
		return "OfflineHeroExpAdd"
	}
	if e == VipRightsType_OfflineTimeLimitAdd {
		return "OfflineTimeLimitAdd"
	}
	if e == VipRightsType_OfflineQuickGetTimesAdd {
		return "OfflineQuickGetTimesAdd"
	}
	if e == VipRightsType_ArenaTicketsLimitAdd {
		return "ArenaTicketsLimitAdd"
	}
	if e == VipRightsType_LineupHeroAttackAdd {
		return "LineupHeroAttackAdd"
	}
	if e == VipRightsType_LineupHeroDefenceAdd {
		return "LineupHeroDefenceAdd"
	}
	if e == VipRightsType_LineupHeroHPAdd {
		return "LineupHeroHPAdd"
	}
	if e == VipRightsType_ArtsTrialChallengeTimesAdd {
		return "ArtsTrialChallengeTimesAdd"
	}
	if e == VipRightsType_ArtifactDrawAdd {
		return "ArtifactDrawAdd"
	}
	if e == VipRightsType_TeamExpAdd {
		return "TeamExpAdd"
	}
	if e == VipRightsType_ArenaRefreshCdSub {
		return "ArenaRefreshCdSub"
	}
	if e == VipRightsType_ArenaShieldDurationAdd {
		return "ArenaShieldDurationAdd"
	}
	if e == VipRightsType_SkipAllAds {
		return "SkipAllAds"
	}
	if e == VipRightsType_FameExpAdd {
		return "FameExpAdd"
	}
	if e == VipRightsType_OfflineQuickGetFreeTimesAdd {
		return "OfflineQuickGetFreeTimesAdd"
	}
	if e == VipRightsType_OfflineQuickGetPayTimesAdd {
		return "OfflineQuickGetPayTimesAdd"
	}
	if e == VipRightsType_RecruitProbUp {
		return "RecruitProbUp"
	}
	if e == VipRightsType_AutoNextLevel {
		return "AutoNextLevel"
	}
	if e == VipRightsType_LDRefreshFormation {
		return "LDRefreshFormation"
	}
	if e == VipRightsType_LDSpecialCard {
		return "LDSpecialCard"
	}
	if e == VipRightsType_LifeCardSpeedUpOpenBox {
		return "LifeCardSpeedUpOpenBox"
	}
	if e == VipRightsType_FishCardDelegate {
		return "FishCardDelegate"
	}
	if e == VipRightsType_FishRodPropAdd {
		return "FishRodPropAdd"
	}
	if e == VipRightsType_FishHookPropAdd {
		return "FishHookPropAdd"
	}
	if e == VipRightsType_FishWirePropAdd {
		return "FishWirePropAdd"
	}
	if e == VipRightsType_FishFloatPropAdd {
		return "FishFloatPropAdd"
	}
	if e == VipRightsType_FishSpeedAdd {
		return "FishSpeedAdd"
	}
	if e == VipRightsType_FishSellRubbishAdd {
		return "FishSellRubbishAdd"
	}
	if e == VipRightsType_HeroSoulUnlock {
		return "HeroSoulUnlock"
	}
	if e == VipRightsType_ManorButler {
		return "ManorButler"
	}
	if e == VipRightsType_SpTDSweep {
		return "SpTDSweep"
	}
	if e == VipRightsType_PatrolShare {
		return "PatrolShare"
	}
	if e == VipRightsType_Send {
		return "Send"
	}
	if e == VipRightsType_RainbowSoulComposeProbAdd {
		return "RainbowSoulComposeProbAdd"
	}
	if e == VipRightsType_SkipBattle {
		return "SkipBattle"
	}
	if e == VipRightsType_ExtraBoxGiftLevel {
		return "ExtraBoxGiftLevel"
	}
	if e == VipRightsType_ManorCheer {
		return "ManorCheer"
	}
	if e == VipRightsType_TerritoryOutput {
		return "TerritoryOutput"
	}
	if e == VipRightsType_TeamShare {
		return "TeamShare"
	}
	if e == VipRightsType_ChatGroupAnnounce {
		return "ChatGroupAnnounce"
	}
	if e == VipRightsType_TeamTicketAdd {
		return "TeamTicketAdd"
	}
	if e == VipRightsType_RoyalBossAdd {
		return "RoyalBossAdd"
	}
	if e == VipRightsType_SendFashion {
		return "SendFashion"
	}
	if e == VipRightsType_MillionaireReduceCoin {
		return "MillionaireReduceCoin"
	}

	return ""

}

func GetVipRightsType(tp string) (VipRightsType, bool) {
	switch tp {
	case "OfflineGoldAdd":
		return VipRightsType_OfflineGoldAdd, true
	case "OfflineHeroExpAdd":
		return VipRightsType_OfflineHeroExpAdd, true
	case "OfflineTimeLimitAdd":
		return VipRightsType_OfflineTimeLimitAdd, true
	case "OfflineQuickGetTimesAdd":
		return VipRightsType_OfflineQuickGetTimesAdd, true
	case "ArenaTicketsLimitAdd":
		return VipRightsType_ArenaTicketsLimitAdd, true
	case "LineupHeroAttackAdd":
		return VipRightsType_LineupHeroAttackAdd, true
	case "LineupHeroDefenceAdd":
		return VipRightsType_LineupHeroDefenceAdd, true
	case "LineupHeroHPAdd":
		return VipRightsType_LineupHeroHPAdd, true
	case "ArtsTrialChallengeTimesAdd":
		return VipRightsType_ArtsTrialChallengeTimesAdd, true
	case "ArtifactDrawAdd":
		return VipRightsType_ArtifactDrawAdd, true
	case "TeamExpAdd":
		return VipRightsType_TeamExpAdd, true
	case "ArenaRefreshCdSub":
		return VipRightsType_ArenaRefreshCdSub, true
	case "ArenaShieldDurationAdd":
		return VipRightsType_ArenaShieldDurationAdd, true
	case "SkipAllAds":
		return VipRightsType_SkipAllAds, true
	case "FameExpAdd":
		return VipRightsType_FameExpAdd, true
	case "OfflineQuickGetFreeTimesAdd":
		return VipRightsType_OfflineQuickGetFreeTimesAdd, true
	case "OfflineQuickGetPayTimesAdd":
		return VipRightsType_OfflineQuickGetPayTimesAdd, true
	case "RecruitProbUp":
		return VipRightsType_RecruitProbUp, true
	case "AutoNextLevel":
		return VipRightsType_AutoNextLevel, true
	case "LDRefreshFormation":
		return VipRightsType_LDRefreshFormation, true
	case "LDSpecialCard":
		return VipRightsType_LDSpecialCard, true
	case "LifeCardSpeedUpOpenBox":
		return VipRightsType_LifeCardSpeedUpOpenBox, true
	case "FishCardDelegate":
		return VipRightsType_FishCardDelegate, true
	case "FishRodPropAdd":
		return VipRightsType_FishRodPropAdd, true
	case "FishHookPropAdd":
		return VipRightsType_FishHookPropAdd, true
	case "FishWirePropAdd":
		return VipRightsType_FishWirePropAdd, true
	case "FishFloatPropAdd":
		return VipRightsType_FishFloatPropAdd, true
	case "FishSpeedAdd":
		return VipRightsType_FishSpeedAdd, true
	case "FishSellRubbishAdd":
		return VipRightsType_FishSellRubbishAdd, true
	case "HeroSoulUnlock":
		return VipRightsType_HeroSoulUnlock, true
	case "ManorButler":
		return VipRightsType_ManorButler, true
	case "SpTDSweep":
		return VipRightsType_SpTDSweep, true
	case "PatrolShare":
		return VipRightsType_PatrolShare, true
	case "Send":
		return VipRightsType_Send, true
	case "RainbowSoulComposeProbAdd":
		return VipRightsType_RainbowSoulComposeProbAdd, true
	case "SkipBattle":
		return VipRightsType_SkipBattle, true
	case "ExtraBoxGiftLevel":
		return VipRightsType_ExtraBoxGiftLevel, true
	case "ManorCheer":
		return VipRightsType_ManorCheer, true
	case "TerritoryOutput":
		return VipRightsType_TerritoryOutput, true
	case "TeamShare":
		return VipRightsType_TeamShare, true
	case "ChatGroupAnnounce":
		return VipRightsType_ChatGroupAnnounce, true
	case "TeamTicketAdd":
		return VipRightsType_TeamTicketAdd, true
	case "RoyalBossAdd":
		return VipRightsType_RoyalBossAdd, true
	case "SendFashion":
		return VipRightsType_SendFashion, true
	case "MillionaireReduceCoin":
		return VipRightsType_MillionaireReduceCoin, true

	}
	return -1, false

}

func GetVipRightsType_ByAlias(tp string) (VipRightsType, bool) {
	switch tp {
	case "离线，金币道具奖励+%，只影响领取全部":
		return VipRightsType_OfflineGoldAdd, true
	case "离线，英雄经验道具掉落加成+%，只影响领取全部":
		return VipRightsType_OfflineHeroExpAdd, true
	case "离线，上限时间+X小时":
		return VipRightsType_OfflineTimeLimitAdd, true
	case "离线，每日快速巡逻次数+X":
		return VipRightsType_OfflineQuickGetTimesAdd, true
	case "竞技场，挑战券存储上限+X":
		return VipRightsType_ArenaTicketsLimitAdd, true
	case "全体上阵英雄攻击+%":
		return VipRightsType_LineupHeroAttackAdd, true
	case "全体上阵英雄防御+%":
		return VipRightsType_LineupHeroDefenceAdd, true
	case "全体上阵英雄血量+%":
		return VipRightsType_LineupHeroHPAdd, true
	case "武道试炼，3个副本每日免费挑战次数+X":
		return VipRightsType_ArtsTrialChallengeTimesAdd, true
	case "神器召唤最后一档可以用钻石召唤":
		return VipRightsType_ArtifactDrawAdd, true
	case "战队经验获取加成+%":
		return VipRightsType_TeamExpAdd, true
	case "竞技场刷新cd减少（秒）":
		return VipRightsType_ArenaRefreshCdSub, true
	case "终身卡竞技场保护时间增加（秒）":
		return VipRightsType_ArenaShieldDurationAdd, true
	case "免除所有广告":
		return VipRightsType_SkipAllAds, true
	case "荣誉每次升级经验额外+X":
		return VipRightsType_FameExpAdd, true
	case "离线，每日免费巡逻次数+X":
		return VipRightsType_OfflineQuickGetFreeTimesAdd, true
	case "离线，每日付费巡逻次数+X":
		return VipRightsType_OfflineQuickGetPayTimesAdd, true
	case "招募心愿单英雄起始几率提升":
		return VipRightsType_RecruitProbUp, true
	case "自动挑战下一关":
		return VipRightsType_AutoNextLevel, true
	case "LD阵容视频刷新次数+X":
		return VipRightsType_LDRefreshFormation, true
	case "LD特殊卡视频获取次数+X":
		return VipRightsType_LDSpecialCard, true
	case "开盒加速":
		return VipRightsType_LifeCardSpeedUpOpenBox, true
	case "钓鱼委托":
		return VipRightsType_FishCardDelegate, true
	case "鱼竿属性+%":
		return VipRightsType_FishRodPropAdd, true
	case "鱼钩属性+%":
		return VipRightsType_FishHookPropAdd, true
	case "鱼线属性+X":
		return VipRightsType_FishWirePropAdd, true
	case "鱼漂属性+%":
		return VipRightsType_FishFloatPropAdd, true
	case "钓鱼速度+%":
		return VipRightsType_FishSpeedAdd, true
	case "出售杂物银币增加+%":
		return VipRightsType_FishSellRubbishAdd, true
	case "八门提前解锁":
		return VipRightsType_HeroSoulUnlock, true
	case "领地管家":
		return VipRightsType_ManorButler, true
	case "南镇扫荡":
		return VipRightsType_SpTDSweep, true
	case "探秘分享":
		return VipRightsType_PatrolShare, true
	case "赠送":
		return VipRightsType_Send, true
	case "永恒八门合成概率+%":
		return VipRightsType_RainbowSoulComposeProbAdd, true
	case "跳过战斗":
		return VipRightsType_SkipBattle, true
	case "额外祝福等级":
		return VipRightsType_ExtraBoxGiftLevel, true
	case "领地助威":
		return VipRightsType_ManorCheer, true
	case "领地产出":
		return VipRightsType_TerritoryOutput, true
	case "组队分享":
		return VipRightsType_TeamShare, true
	case "交流群发布":
		return VipRightsType_ChatGroupAnnounce, true
	case "组队每日额外领取门票次数":
		return VipRightsType_TeamTicketAdd, true
	case "大蛇模拟战每日额外领奖次数":
		return VipRightsType_RoyalBossAdd, true
	case "装扮赠送":
		return VipRightsType_SendFashion, true
	case "百万大富翁装备消耗减少":
		return VipRightsType_MillionaireReduceCoin, true

	}
	return -1, false

}
