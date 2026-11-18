package config

// 功能模块配置-拳皇/SystemEnum
type SysIconGroup int

const (
	SysIconGroup_BEGIN                  SysIconGroup = 0  //枚举开始
	SysIconGroup_None                   SysIconGroup = 0  //无
	SysIconGroup_HomeBottom             SysIconGroup = 1  //主界面底部
	SysIconGroup_Adventure              SysIconGroup = 2  //挑战
	SysIconGroup_HomeLeft               SysIconGroup = 3  //主界面左边
	SysIconGroup_HomeRight              SysIconGroup = 4  //主界面右边
	SysIconGroup_NewbieActivity         SysIconGroup = 5  //新手活动
	SysIconGroup_HomeCentral            SysIconGroup = 6  //主界面中部
	SysIconGroup_Test                   SysIconGroup = 7  //测试
	SysIconGroup_HomeRightTop           SysIconGroup = 8  //主界面右上边
	SysIconGroup_HomeActiveBannerTop    SysIconGroup = 9  //主界面上边banner活动
	SysIconGroup_HomeActiveBannerBottom SysIconGroup = 10 //主界面下边banner活动
	SysIconGroup_Activity               SysIconGroup = 11 //活动
	SysIconGroup_WelfareActivity        SysIconGroup = 12 //精彩福利
	SysIconGroup_WeeklyCycle            SysIconGroup = 13 //周循环
	SysIconGroup_XinShouZhaoMu          SysIconGroup = 14 //新手招募
	SysIconGroup_HomeRightBottom        SysIconGroup = 15 //主界面右下边
	SysIconGroup_HomeRightTop2          SysIconGroup = 16 //主界面右上边第二列
	SysIconGroup_BattleQuickBar         SysIconGroup = 17 //战斗快捷栏
	SysIconGroup_HomeTop                SysIconGroup = 18 //主界面顶部
	SysIconGroup_AllianceChallenge      SysIconGroup = 19 //公会挑战
	SysIconGroup_UpRecruit              SysIconGroup = 20 //Up招募
	SysIconGroup_WishRecruit            SysIconGroup = 21 //自选招募
	SysIconGroup_VeteranActivity        SysIconGroup = 22 //精英活动
	SysIconGroup_VeteranHeroRecruit     SysIconGroup = 23 //精英招募
	SysIconGroup_HomeRightTop3          SysIconGroup = 24 //主界面右上边第三列
	SysIconGroup_LeaderBoardRush        SysIconGroup = 25 //冲榜排行榜
	SysIconGroup_BannerPreview          SysIconGroup = 26 //玩法活动预览
	SysIconGroup_FestivalPartyActivity  SysIconGroup = 27 //双旦节日活动
	SysIconGroup_Recruit                SysIconGroup = 28 //招募界面
	SysIconGroup_AlipayAct              SysIconGroup = 29 //支付宝活动
	SysIconGroup_Arena                  SysIconGroup = 30 //竞技场界面
	SysIconGroup_KingWarGift            SysIconGroup = 31 //皇城备战
	SysIconGroup_END                    SysIconGroup = 32 //枚举结束
)

// 功能模块配置-拳皇/SystemEnum
type SysIconShowRule int

const (
	SysIconShowRule_BEGIN             SysIconShowRule = 1 //枚举开始
	SysIconShowRule_Always            SysIconShowRule = 1 //始终显示
	SysIconShowRule_Unlocked          SysIconShowRule = 2 //解锁显示
	SysIconShowRule_CanEnter          SysIconShowRule = 3 //可进入显示
	SysIconShowRule_None              SysIconShowRule = 4 //不显示
	SysIconShowRule_PreviewOrCanEnter SysIconShowRule = 5 //预热或可进入显示
	SysIconShowRule_OnlyPreview       SysIconShowRule = 6 //仅预热显示
	SysIconShowRule_ShowWithActivity  SysIconShowRule = 7 //跟随活动显示
	SysIconShowRule_END               SysIconShowRule = 8 //枚举结束
)

// 功能模块配置-拳皇/SystemEnum
type SysOpenRuleExType int

const (
	SysOpenRuleExType_BEGIN           SysOpenRuleExType = 0 //枚举开始
	SysOpenRuleExType_None            SysOpenRuleExType = 0 //无
	SysOpenRuleExType_CloseOnIOSAudit SysOpenRuleExType = 1 //ios审核服关闭
	SysOpenRuleExType_OpenAlways      SysOpenRuleExType = 2 //始终开启
	SysOpenRuleExType_CloseAlways     SysOpenRuleExType = 3 //始终关闭
	SysOpenRuleExType_Server          SysOpenRuleExType = 4 //指定服务器开启
	SysOpenRuleExType_CloseOnServer   SysOpenRuleExType = 5 //指定服务器关闭
	SysOpenRuleExType_TargetPlatform  SysOpenRuleExType = 6 //指定平台
	SysOpenRuleExType_END             SysOpenRuleExType = 7 //枚举结束
)

// 功能模块配置-拳皇/SystemEnum
type ViewTag int

const (
	ViewTag_BEGIN ViewTag = 0 //枚举开始
	ViewTag_None  ViewTag = 0 //无
	ViewTag_END   ViewTag = 1 //枚举结束
)

// 功能模块配置-拳皇/SystemEnum
type WeatherTimeType int

const (
	WeatherTimeType_BEGIN   WeatherTimeType = -1 //枚举开始
	WeatherTimeType_None    WeatherTimeType = -1 //关闭
	WeatherTimeType_Morning WeatherTimeType = 0  //清晨
	WeatherTimeType_Noon    WeatherTimeType = 1  //中午
	WeatherTimeType_Dusk    WeatherTimeType = 2  //黄昏
	WeatherTimeType_Night   WeatherTimeType = 3  //夜晚
	WeatherTimeType_Day     WeatherTimeType = 4  //白天
	WeatherTimeType_Sunny   WeatherTimeType = 5  //晴
	WeatherTimeType_Rainy   WeatherTimeType = 6  //雨
	WeatherTimeType_Snow    WeatherTimeType = 7  //雪
	WeatherTimeType_Sand    WeatherTimeType = 8  //沙尘暴
	WeatherTimeType_END     WeatherTimeType = 9  //枚举结束
)

// 功能模块配置-拳皇/SystemEnum
type PlayHelpGroup int

const (
	PlayHelpGroup_BEGIN   PlayHelpGroup = 1 //枚举开始
	PlayHelpGroup_KingWar PlayHelpGroup = 1 //皇城争锋
	PlayHelpGroup_END     PlayHelpGroup = 2 //枚举结束
)

func (e SysIconGroup) String() string {

	if e == SysIconGroup_None {
		return "None"
	}
	if e == SysIconGroup_HomeBottom {
		return "HomeBottom"
	}
	if e == SysIconGroup_Adventure {
		return "Adventure"
	}
	if e == SysIconGroup_HomeLeft {
		return "HomeLeft"
	}
	if e == SysIconGroup_HomeRight {
		return "HomeRight"
	}
	if e == SysIconGroup_NewbieActivity {
		return "NewbieActivity"
	}
	if e == SysIconGroup_HomeCentral {
		return "HomeCentral"
	}
	if e == SysIconGroup_Test {
		return "Test"
	}
	if e == SysIconGroup_HomeRightTop {
		return "HomeRightTop"
	}
	if e == SysIconGroup_HomeActiveBannerTop {
		return "HomeActiveBannerTop"
	}
	if e == SysIconGroup_HomeActiveBannerBottom {
		return "HomeActiveBannerBottom"
	}
	if e == SysIconGroup_Activity {
		return "Activity"
	}
	if e == SysIconGroup_WelfareActivity {
		return "WelfareActivity"
	}
	if e == SysIconGroup_WeeklyCycle {
		return "WeeklyCycle"
	}
	if e == SysIconGroup_XinShouZhaoMu {
		return "XinShouZhaoMu"
	}
	if e == SysIconGroup_HomeRightBottom {
		return "HomeRightBottom"
	}
	if e == SysIconGroup_HomeRightTop2 {
		return "HomeRightTop2"
	}
	if e == SysIconGroup_BattleQuickBar {
		return "BattleQuickBar"
	}
	if e == SysIconGroup_HomeTop {
		return "HomeTop"
	}
	if e == SysIconGroup_AllianceChallenge {
		return "AllianceChallenge"
	}
	if e == SysIconGroup_UpRecruit {
		return "UpRecruit"
	}
	if e == SysIconGroup_WishRecruit {
		return "WishRecruit"
	}
	if e == SysIconGroup_VeteranActivity {
		return "VeteranActivity"
	}
	if e == SysIconGroup_VeteranHeroRecruit {
		return "VeteranHeroRecruit"
	}
	if e == SysIconGroup_HomeRightTop3 {
		return "HomeRightTop3"
	}
	if e == SysIconGroup_LeaderBoardRush {
		return "LeaderBoardRush"
	}
	if e == SysIconGroup_BannerPreview {
		return "BannerPreview"
	}
	if e == SysIconGroup_FestivalPartyActivity {
		return "FestivalPartyActivity"
	}
	if e == SysIconGroup_Recruit {
		return "Recruit"
	}
	if e == SysIconGroup_AlipayAct {
		return "AlipayAct"
	}
	if e == SysIconGroup_Arena {
		return "Arena"
	}
	if e == SysIconGroup_KingWarGift {
		return "KingWarGift"
	}

	return ""

}

func GetSysIconGroup(tp string) (SysIconGroup, bool) {
	switch tp {
	case "None":
		return SysIconGroup_None, true
	case "HomeBottom":
		return SysIconGroup_HomeBottom, true
	case "Adventure":
		return SysIconGroup_Adventure, true
	case "HomeLeft":
		return SysIconGroup_HomeLeft, true
	case "HomeRight":
		return SysIconGroup_HomeRight, true
	case "NewbieActivity":
		return SysIconGroup_NewbieActivity, true
	case "HomeCentral":
		return SysIconGroup_HomeCentral, true
	case "Test":
		return SysIconGroup_Test, true
	case "HomeRightTop":
		return SysIconGroup_HomeRightTop, true
	case "HomeActiveBannerTop":
		return SysIconGroup_HomeActiveBannerTop, true
	case "HomeActiveBannerBottom":
		return SysIconGroup_HomeActiveBannerBottom, true
	case "Activity":
		return SysIconGroup_Activity, true
	case "WelfareActivity":
		return SysIconGroup_WelfareActivity, true
	case "WeeklyCycle":
		return SysIconGroup_WeeklyCycle, true
	case "XinShouZhaoMu":
		return SysIconGroup_XinShouZhaoMu, true
	case "HomeRightBottom":
		return SysIconGroup_HomeRightBottom, true
	case "HomeRightTop2":
		return SysIconGroup_HomeRightTop2, true
	case "BattleQuickBar":
		return SysIconGroup_BattleQuickBar, true
	case "HomeTop":
		return SysIconGroup_HomeTop, true
	case "AllianceChallenge":
		return SysIconGroup_AllianceChallenge, true
	case "UpRecruit":
		return SysIconGroup_UpRecruit, true
	case "WishRecruit":
		return SysIconGroup_WishRecruit, true
	case "VeteranActivity":
		return SysIconGroup_VeteranActivity, true
	case "VeteranHeroRecruit":
		return SysIconGroup_VeteranHeroRecruit, true
	case "HomeRightTop3":
		return SysIconGroup_HomeRightTop3, true
	case "LeaderBoardRush":
		return SysIconGroup_LeaderBoardRush, true
	case "BannerPreview":
		return SysIconGroup_BannerPreview, true
	case "FestivalPartyActivity":
		return SysIconGroup_FestivalPartyActivity, true
	case "Recruit":
		return SysIconGroup_Recruit, true
	case "AlipayAct":
		return SysIconGroup_AlipayAct, true
	case "Arena":
		return SysIconGroup_Arena, true
	case "KingWarGift":
		return SysIconGroup_KingWarGift, true

	}
	return -1, false

}

func GetSysIconGroup_ByAlias(tp string) (SysIconGroup, bool) {
	switch tp {
	case "无":
		return SysIconGroup_None, true
	case "主界面底部":
		return SysIconGroup_HomeBottom, true
	case "挑战":
		return SysIconGroup_Adventure, true
	case "主界面左边":
		return SysIconGroup_HomeLeft, true
	case "主界面右边":
		return SysIconGroup_HomeRight, true
	case "新手活动":
		return SysIconGroup_NewbieActivity, true
	case "主界面中部":
		return SysIconGroup_HomeCentral, true
	case "测试":
		return SysIconGroup_Test, true
	case "主界面右上边":
		return SysIconGroup_HomeRightTop, true
	case "主界面上边banner活动":
		return SysIconGroup_HomeActiveBannerTop, true
	case "主界面下边banner活动":
		return SysIconGroup_HomeActiveBannerBottom, true
	case "活动":
		return SysIconGroup_Activity, true
	case "精彩福利":
		return SysIconGroup_WelfareActivity, true
	case "周循环":
		return SysIconGroup_WeeklyCycle, true
	case "新手招募":
		return SysIconGroup_XinShouZhaoMu, true
	case "主界面右下边":
		return SysIconGroup_HomeRightBottom, true
	case "主界面右上边第二列":
		return SysIconGroup_HomeRightTop2, true
	case "战斗快捷栏":
		return SysIconGroup_BattleQuickBar, true
	case "主界面顶部":
		return SysIconGroup_HomeTop, true
	case "公会挑战":
		return SysIconGroup_AllianceChallenge, true
	case "Up招募":
		return SysIconGroup_UpRecruit, true
	case "自选招募":
		return SysIconGroup_WishRecruit, true
	case "精英活动":
		return SysIconGroup_VeteranActivity, true
	case "精英招募":
		return SysIconGroup_VeteranHeroRecruit, true
	case "主界面右上边第三列":
		return SysIconGroup_HomeRightTop3, true
	case "冲榜排行榜":
		return SysIconGroup_LeaderBoardRush, true
	case "玩法活动预览":
		return SysIconGroup_BannerPreview, true
	case "双旦节日活动":
		return SysIconGroup_FestivalPartyActivity, true
	case "招募界面":
		return SysIconGroup_Recruit, true
	case "支付宝活动":
		return SysIconGroup_AlipayAct, true
	case "竞技场界面":
		return SysIconGroup_Arena, true
	case "皇城备战":
		return SysIconGroup_KingWarGift, true

	}
	return -1, false

}

func (e SysIconShowRule) String() string {

	if e == SysIconShowRule_Always {
		return "Always"
	}
	if e == SysIconShowRule_Unlocked {
		return "Unlocked"
	}
	if e == SysIconShowRule_CanEnter {
		return "CanEnter"
	}
	if e == SysIconShowRule_None {
		return "None"
	}
	if e == SysIconShowRule_PreviewOrCanEnter {
		return "PreviewOrCanEnter"
	}
	if e == SysIconShowRule_OnlyPreview {
		return "OnlyPreview"
	}
	if e == SysIconShowRule_ShowWithActivity {
		return "ShowWithActivity"
	}

	return ""

}

func GetSysIconShowRule(tp string) (SysIconShowRule, bool) {
	switch tp {
	case "Always":
		return SysIconShowRule_Always, true
	case "Unlocked":
		return SysIconShowRule_Unlocked, true
	case "CanEnter":
		return SysIconShowRule_CanEnter, true
	case "None":
		return SysIconShowRule_None, true
	case "PreviewOrCanEnter":
		return SysIconShowRule_PreviewOrCanEnter, true
	case "OnlyPreview":
		return SysIconShowRule_OnlyPreview, true
	case "ShowWithActivity":
		return SysIconShowRule_ShowWithActivity, true

	}
	return -1, false

}

func GetSysIconShowRule_ByAlias(tp string) (SysIconShowRule, bool) {
	switch tp {
	case "始终显示":
		return SysIconShowRule_Always, true
	case "解锁显示":
		return SysIconShowRule_Unlocked, true
	case "可进入显示":
		return SysIconShowRule_CanEnter, true
	case "不显示":
		return SysIconShowRule_None, true
	case "预热或可进入显示":
		return SysIconShowRule_PreviewOrCanEnter, true
	case "仅预热显示":
		return SysIconShowRule_OnlyPreview, true
	case "跟随活动显示":
		return SysIconShowRule_ShowWithActivity, true

	}
	return -1, false

}

func (e SysOpenRuleExType) String() string {

	if e == SysOpenRuleExType_None {
		return "None"
	}
	if e == SysOpenRuleExType_CloseOnIOSAudit {
		return "CloseOnIOSAudit"
	}
	if e == SysOpenRuleExType_OpenAlways {
		return "OpenAlways"
	}
	if e == SysOpenRuleExType_CloseAlways {
		return "CloseAlways"
	}
	if e == SysOpenRuleExType_Server {
		return "Server"
	}
	if e == SysOpenRuleExType_CloseOnServer {
		return "CloseOnServer"
	}
	if e == SysOpenRuleExType_TargetPlatform {
		return "TargetPlatform"
	}

	return ""

}

func GetSysOpenRuleExType(tp string) (SysOpenRuleExType, bool) {
	switch tp {
	case "None":
		return SysOpenRuleExType_None, true
	case "CloseOnIOSAudit":
		return SysOpenRuleExType_CloseOnIOSAudit, true
	case "OpenAlways":
		return SysOpenRuleExType_OpenAlways, true
	case "CloseAlways":
		return SysOpenRuleExType_CloseAlways, true
	case "Server":
		return SysOpenRuleExType_Server, true
	case "CloseOnServer":
		return SysOpenRuleExType_CloseOnServer, true
	case "TargetPlatform":
		return SysOpenRuleExType_TargetPlatform, true

	}
	return -1, false

}

func GetSysOpenRuleExType_ByAlias(tp string) (SysOpenRuleExType, bool) {
	switch tp {
	case "无":
		return SysOpenRuleExType_None, true
	case "ios审核服关闭":
		return SysOpenRuleExType_CloseOnIOSAudit, true
	case "始终开启":
		return SysOpenRuleExType_OpenAlways, true
	case "始终关闭":
		return SysOpenRuleExType_CloseAlways, true
	case "指定服务器开启":
		return SysOpenRuleExType_Server, true
	case "指定服务器关闭":
		return SysOpenRuleExType_CloseOnServer, true
	case "指定平台":
		return SysOpenRuleExType_TargetPlatform, true

	}
	return -1, false

}

func (e ViewTag) String() string {

	if e == ViewTag_None {
		return "None"
	}

	return ""

}

func GetViewTag(tp string) (ViewTag, bool) {
	switch tp {
	case "None":
		return ViewTag_None, true

	}
	return -1, false

}

func GetViewTag_ByAlias(tp string) (ViewTag, bool) {
	switch tp {
	case "无":
		return ViewTag_None, true

	}
	return -1, false

}

func (e WeatherTimeType) String() string {

	if e == WeatherTimeType_None {
		return "None"
	}
	if e == WeatherTimeType_Morning {
		return "Morning"
	}
	if e == WeatherTimeType_Noon {
		return "Noon"
	}
	if e == WeatherTimeType_Dusk {
		return "Dusk"
	}
	if e == WeatherTimeType_Night {
		return "Night"
	}
	if e == WeatherTimeType_Day {
		return "Day"
	}
	if e == WeatherTimeType_Sunny {
		return "Sunny"
	}
	if e == WeatherTimeType_Rainy {
		return "Rainy"
	}
	if e == WeatherTimeType_Snow {
		return "Snow"
	}
	if e == WeatherTimeType_Sand {
		return "Sand"
	}

	return ""

}

func GetWeatherTimeType(tp string) (WeatherTimeType, bool) {
	switch tp {
	case "None":
		return WeatherTimeType_None, true
	case "Morning":
		return WeatherTimeType_Morning, true
	case "Noon":
		return WeatherTimeType_Noon, true
	case "Dusk":
		return WeatherTimeType_Dusk, true
	case "Night":
		return WeatherTimeType_Night, true
	case "Day":
		return WeatherTimeType_Day, true
	case "Sunny":
		return WeatherTimeType_Sunny, true
	case "Rainy":
		return WeatherTimeType_Rainy, true
	case "Snow":
		return WeatherTimeType_Snow, true
	case "Sand":
		return WeatherTimeType_Sand, true

	}
	return -1, false

}

func GetWeatherTimeType_ByAlias(tp string) (WeatherTimeType, bool) {
	switch tp {
	case "关闭":
		return WeatherTimeType_None, true
	case "清晨":
		return WeatherTimeType_Morning, true
	case "中午":
		return WeatherTimeType_Noon, true
	case "黄昏":
		return WeatherTimeType_Dusk, true
	case "夜晚":
		return WeatherTimeType_Night, true
	case "白天":
		return WeatherTimeType_Day, true
	case "晴":
		return WeatherTimeType_Sunny, true
	case "雨":
		return WeatherTimeType_Rainy, true
	case "雪":
		return WeatherTimeType_Snow, true
	case "沙尘暴":
		return WeatherTimeType_Sand, true

	}
	return -1, false

}

func (e PlayHelpGroup) String() string {

	if e == PlayHelpGroup_KingWar {
		return "KingWar"
	}

	return ""

}

func GetPlayHelpGroup(tp string) (PlayHelpGroup, bool) {
	switch tp {
	case "KingWar":
		return PlayHelpGroup_KingWar, true

	}
	return -1, false

}

func GetPlayHelpGroup_ByAlias(tp string) (PlayHelpGroup, bool) {
	switch tp {
	case "皇城争锋":
		return PlayHelpGroup_KingWar, true

	}
	return -1, false

}
