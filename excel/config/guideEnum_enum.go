package config

// 引导表-拳皇/guideEnum
type GuideFramePos int

const (
	GuideFramePos_BEGIN GuideFramePos = 1 //枚举开始
	GuideFramePos_Up    GuideFramePos = 1 //上方
	GuideFramePos_Down  GuideFramePos = 2 //下方
	GuideFramePos_END   GuideFramePos = 3 //枚举结束
)

// 引导表-拳皇/guideEnum
type GuideTalkDirection int

const (
	GuideTalkDirection_BEGIN GuideTalkDirection = 1 //枚举开始
	GuideTalkDirection_Left  GuideTalkDirection = 1 //左
	GuideTalkDirection_Right GuideTalkDirection = 2 //右
	GuideTalkDirection_END   GuideTalkDirection = 3 //枚举结束
)

// 引导表-拳皇/guideEnum
type GuideHandAction int

const (
	GuideHandAction_BEGIN     GuideHandAction = 1 //枚举开始
	GuideHandAction_Click     GuideHandAction = 1 //点击
	GuideHandAction_Slide     GuideHandAction = 2 //滑动
	GuideHandAction_Circle    GuideHandAction = 3 //画圈
	GuideHandAction_Idle      GuideHandAction = 4 //放置
	GuideHandAction_LongPress GuideHandAction = 5 //长按
	GuideHandAction_END       GuideHandAction = 6 //枚举结束
)

// 引导表-拳皇/guideEnum
type GuideCustomAction int

const (
	GuideCustomAction_BEGIN               GuideCustomAction = 1  //枚举开始
	GuideCustomAction_FirstStartFight     GuideCustomAction = 1  //开始战斗
	GuideCustomAction_RunNext             GuideCustomAction = 2  //执行下一条
	GuideCustomAction_ShowHeroRange       GuideCustomAction = 3  //显示英雄范围
	GuideCustomAction_ChangeWave          GuideCustomAction = 4  //切换波次
	GuideCustomAction_Unlock              GuideCustomAction = 5  //解锁功能
	GuideCustomAction_AddItem             GuideCustomAction = 6  //发放货币
	GuideCustomAction_CloseBottomView     GuideCustomAction = 7  //关闭底部弹窗
	GuideCustomAction_HeroFollowBoss      GuideCustomAction = 8  //英雄追Boss
	GuideCustomAction_CheckHeroPos        GuideCustomAction = 9  //检查英雄位置
	GuideCustomAction_StopCheckHeroPos    GuideCustomAction = 10 //停止检查英雄位置
	GuideCustomAction_AssignSummonResults GuideCustomAction = 11 //指定召唤结果
	GuideCustomAction_StopSpawnEnemy      GuideCustomAction = 12 //停止刷怪
	GuideCustomAction_ShakeMap            GuideCustomAction = 13 //地图震动
	GuideCustomAction_HeroAtkMulti        GuideCustomAction = 14 //英雄攻击乘
	GuideCustomAction_AssignNextTroop     GuideCustomAction = 15 //指定下次阵容
	GuideCustomAction_GetNewHero          GuideCustomAction = 16 //获得新英雄
	GuideCustomAction_BattlePause         GuideCustomAction = 17 //战斗暂停
	GuideCustomAction_BattleResume        GuideCustomAction = 18 //战斗继续
	GuideCustomAction_BeginReplaceHero    GuideCustomAction = 19 //开始替换引导
	GuideCustomAction_EndReplaceHero      GuideCustomAction = 20 //结束替换引导
	GuideCustomAction_ScrollToEntry       GuideCustomAction = 21 //滑动到节点
	GuideCustomAction_END                 GuideCustomAction = 22 //枚举结束
)

// 引导表-拳皇/guideEnum
type GuideCondiType int

const (
	GuideCondiType_BEGIN                    GuideCondiType = 1  //枚举开始
	GuideCondiType_UIFocus                  GuideCondiType = 1  //界面置顶
	GuideCondiType_FunctionUnlock           GuideCondiType = 2  //功能解锁
	GuideCondiType_TaskFinish               GuideCondiType = 3  //完成任务
	GuideCondiType_ClickAny                 GuideCondiType = 4  //点击任意位置
	GuideCondiType_ClickTarget              GuideCondiType = 5  //点击按钮
	GuideCondiType_GroupFinish              GuideCondiType = 6  //完成引导组
	GuideCondiType_StepFinish               GuideCondiType = 7  //完成引导
	GuideCondiType_SelfCD                   GuideCondiType = 8  //间隔时间
	GuideCondiType_GlobalCD                 GuideCondiType = 9  //全局间隔时间
	GuideCondiType_Manual                   GuideCondiType = 10 //主动触发
	GuideCondiType_IdleTime                 GuideCondiType = 11 //无操作时间大于
	GuideCondiType_TaskIdLess               GuideCondiType = 12 //任务不大于
	GuideCondiType_DailyTimes               GuideCondiType = 13 //每日次数
	GuideCondiType_HeroLevelUp              GuideCondiType = 14 //英雄等级
	GuideCondiType_EquipUnlock              GuideCondiType = 15 //装备解锁
	GuideCondiType_PassMainLevel            GuideCondiType = 16 //到达关卡
	GuideCondiType_Wait                     GuideCondiType = 17 //延迟
	GuideCondiType_LevelGreate              GuideCondiType = 18 //关卡大于等于
	GuideCondiType_GuideBattleStart         GuideCondiType = 19 //开始引导关
	GuideCondiType_Coin                     GuideCondiType = 20 //局内金币
	GuideCondiType_CanSummon                GuideCondiType = 21 //可召唤
	GuideCondiType_CanMerge                 GuideCondiType = 22 //可合成
	GuideCondiType_WaveEqual                GuideCondiType = 23 //波数等于
	GuideCondiType_GameTime                 GuideCondiType = 24 //关卡时间大于
	GuideCondiType_EnemyCnt                 GuideCondiType = 25 //敌人数量
	GuideCondiType_EnemyClear               GuideCondiType = 26 //清空敌人
	GuideCondiType_GameEnd                  GuideCondiType = 27 //关卡结束
	GuideCondiType_KillBoss                 GuideCondiType = 28 //击杀Boss
	GuideCondiType_MoveFinished             GuideCondiType = 29 //移动完成
	GuideCondiType_UnlockFinished           GuideCondiType = 30 //解锁动画完成
	GuideCondiType_ShakeMapFinished         GuideCondiType = 31 //地图震动完成
	GuideCondiType_GambleSuccess            GuideCondiType = 32 //抽卡完成
	GuideCondiType_BattleTimes              GuideCondiType = 33 //完成挑战次数
	GuideCondiType_HighestWave              GuideCondiType = 34 //最高波数
	GuideCondiType_HighestWaveLess          GuideCondiType = 35 //最高波数小于
	GuideCondiType_RoleCreateVersion        GuideCondiType = 36 //创角版本号
	GuideCondiType_RoleCreateVersionLess    GuideCondiType = 37 //创角版本号小于
	GuideCondiType_MaxSp                    GuideCondiType = 38 //能放大招
	GuideCondiType_NoChapterUpdate          GuideCondiType = 39 //没有章节动画
	GuideCondiType_FirstGetRedHero          GuideCondiType = 40 //有一个红英雄
	GuideCondiType_RoleCreateVersionGreater GuideCondiType = 41 //创角版本号大于
	GuideCondiType_PlayerLevelGreater       GuideCondiType = 42 //战队等级大于等于
	GuideCondiType_CanSetOffHero            GuideCondiType = 43 //可上阵英雄
	GuideCondiType_PlayerLevelLess          GuideCondiType = 44 //战队等级小于
	GuideCondiType_DestinyBind              GuideCondiType = 45 //绑定宿命英雄
	GuideCondiType_END                      GuideCondiType = 46 //枚举结束
)

// 引导表-拳皇/guideEnum
type GuideViewKey int

const (
	GuideViewKey_BEGIN                  GuideViewKey = 0  //枚举开始
	GuideViewKey_None                   GuideViewKey = 0  //无
	GuideViewKey_BtnHomeTask            GuideViewKey = 1  //任务奖励
	GuideViewKey_BtnHomeTaskGuide       GuideViewKey = 2  //任务引导
	GuideViewKey_Recuit                 GuideViewKey = 3  //招募
	GuideViewKey_BtnRecuitFree          GuideViewKey = 4  //免费招募
	GuideViewKey_BtnStartFight          GuideViewKey = 5  //开始战斗
	GuideViewKey_BtnMainLevelOk         GuideViewKey = 6  //主线胜利
	GuideViewKey_BtnNewHero             GuideViewKey = 7  //新英雄确认
	GuideViewKey_Afk                    GuideViewKey = 8  //离线挂机
	GuideViewKey_BtnAfkReward           GuideViewKey = 9  //离线全部奖励
	GuideViewKey_Hero                   GuideViewKey = 10 //英雄入口
	GuideViewKey_BtnHeroLevel           GuideViewKey = 11 //英雄升级
	GuideViewKey_BtnHeroClose           GuideViewKey = 12 //英雄关闭
	GuideViewKey_BtnMainLevel           GuideViewKey = 13 //主线入口
	GuideViewKey_ClickHero              GuideViewKey = 14 //上阵英雄
	GuideViewKey_BtnHeroEquip           GuideViewKey = 15 //英雄装备
	GuideViewKey_BtnEquipLevelUp        GuideViewKey = 16 //装备强化
	GuideViewKey_Adventure              GuideViewKey = 17 //挑战入口
	GuideViewKey_ArtsTrial              GuideViewKey = 18 //试炼入口
	GuideViewKey_TowerChallenge         GuideViewKey = 19 //爬塔入口
	GuideViewKey_Arena                  GuideViewKey = 20 //竞技场入口
	GuideViewKey_BtnArenaChallenge      GuideViewKey = 21 //竞技场挑战入口
	GuideViewKey_BtnArenaItemChallenge  GuideViewKey = 22 //竞技场挑战对手
	GuideViewKey_BtnSummon              GuideViewKey = 23 //召唤按钮
	GuideViewKey_GridCanMerge           GuideViewKey = 24 //可合成格子
	GuideViewKey_BtnMerge               GuideViewKey = 25 //合成按钮
	GuideViewKey_MoveHeroStartGrid      GuideViewKey = 26 //移动英雄起始格子
	GuideViewKey_MoveHeroEndGrid        GuideViewKey = 27 //移动英雄目标格子
	GuideViewKey_ItemStone              GuideViewKey = 28 //幸运石图标
	GuideViewKey_EntryGamble            GuideViewKey = 29 //抽卡入口
	GuideViewKey_ItemGamble1            GuideViewKey = 30 //抽卡选项
	GuideViewKey_EntryUpgrade           GuideViewKey = 31 //升级入口
	GuideViewKey_ItemUpgade2            GuideViewKey = 32 //紫色升级按钮
	GuideViewKey_ItemUpgade4            GuideViewKey = 33 //概率升级按钮
	GuideViewKey_BtnEpic                GuideViewKey = 34 //神话按钮
	GuideViewKey_NodeEpicNeeds          GuideViewKey = 35 //神话需求列表
	GuideViewKey_BtnMergeEpic           GuideViewKey = 36 //神话合成按钮
	GuideViewKey_BtnAward               GuideViewKey = 37 //结算返回按钮
	GuideViewKey_BtnPvp                 GuideViewKey = 38 //开始挑战按钮
	GuideViewKey_EntryTroop             GuideViewKey = 39 //今日阵容按钮
	GuideViewKey_BtnRefreshTroop        GuideViewKey = 40 //阵容刷新按钮
	GuideViewKey_BtnCloseTroop          GuideViewKey = 41 //阵容返回按钮
	GuideViewKey_BtnAdventureReward     GuideViewKey = 42 //挑战奖励按钮
	GuideViewKey_BtnGet                 GuideViewKey = 43 //领取挑战奖励
	GuideViewKey_BtnMatch               GuideViewKey = 44 //随机组队按钮
	GuideViewKey_BattleFundCard         GuideViewKey = 45 //局内基金卡
	GuideViewKey_ItemChooseBossAward    GuideViewKey = 46 //Boss奖励选项1
	GuideViewKey_BtnChooseBossAward     GuideViewKey = 47 //选择Boss奖励按钮
	GuideViewKey_Expedition             GuideViewKey = 48 //远征入口
	GuideViewKey_RoyalBoss              GuideViewKey = 49 //大蛇模拟战入口
	GuideViewKey_BoxGift                GuideViewKey = 50 //收藏入口
	GuideViewKey_BoxGiftBoxBtn          GuideViewKey = 51 //收藏箱子
	GuideViewKey_BoxGiftBoxReawrd       GuideViewKey = 52 //收藏领取按钮
	GuideViewKey_MaxSP                  GuideViewKey = 53 //大招英雄
	GuideViewKey_BtnArtsTrialFightType1 GuideViewKey = 54 //格斗家试炼挑战按钮
	GuideViewKey_BtnArtsTrialFightItem1 GuideViewKey = 55 //格斗家试炼挑战按钮2
	GuideViewKey_BtnHandUltimate        GuideViewKey = 56 //自动战斗
	GuideViewKey_HomeMenu               GuideViewKey = 57 //菜单
	GuideViewKey_DailyMission           GuideViewKey = 58 //每日任务
	GuideViewKey_Manor                  GuideViewKey = 59 //领地入口
	GuideViewKey_ManorRes               GuideViewKey = 60 //领地资源
	GuideViewKey_ManorWork              GuideViewKey = 61 //一键上车
	GuideViewKey_ManorApply             GuideViewKey = 62 //领地确定
	GuideViewKey_ManorSearch            GuideViewKey = 63 //领地搜索
	GuideViewKey_BoxGiftUpgrade         GuideViewKey = 64 //收藏升级
	GuideViewKey_TabAstrology           GuideViewKey = 65 //神器祈福
	GuideViewKey_BtnReborn              GuideViewKey = 66 //重生按钮
	GuideViewKey_BtnReplaceOk           GuideViewKey = 67 //英雄替换确定
	GuideViewKey_NoRedHero              GuideViewKey = 68 //非红英雄
	GuideViewKey_RedHero                GuideViewKey = 69 //红英雄
	GuideViewKey_BtnReplaceChoose       GuideViewKey = 70 //选择替换英雄
	GuideViewKey_BtnChooseOk            GuideViewKey = 71 //英雄选择确定
	GuideViewKey_BoxGiftPage2           GuideViewKey = 72 //收藏页签2
	GuideViewKey_WorldBoss              GuideViewKey = 73 //世界boss入口
	GuideViewKey_Fish                   GuideViewKey = 74 //钓鱼入口
	GuideViewKey_SpTD                   GuideViewKey = 75 //南镇风云入口
	GuideViewKey_SpRefreshPower         GuideViewKey = 76 //南镇刷新能力
	GuideViewKey_SpRefresh              GuideViewKey = 77 //南镇刷新
	GuideViewKey_SpLevelUp              GuideViewKey = 78 //南镇升级
	GuideViewKey_SpHero                 GuideViewKey = 79 //南镇英雄
	GuideViewKey_SpNext                 GuideViewKey = 80 //南镇下一关
	GuideViewKey_SpTeam                 GuideViewKey = 81 //南镇阵容编辑
	GuideViewKey_Destiny                GuideViewKey = 82 //宿命入口
	GuideViewKey_DestinyGearTab         GuideViewKey = 83 //宿命齿轮入口
	GuideViewKey_DestinyGearDefault     GuideViewKey = 84 //宿命齿轮默认英雄
	GuideViewKey_DestinyGearGoBtn       GuideViewKey = 85 //宿命齿轮前往按钮
	GuideViewKey_TeamEntry              GuideViewKey = 86 //队伍试炼入口
	GuideViewKey_END                    GuideViewKey = 87 //枚举结束
)

func (e GuideFramePos) String() string {

	if e == GuideFramePos_Up {
		return "Up"
	}
	if e == GuideFramePos_Down {
		return "Down"
	}

	return ""

}

func GetGuideFramePos(tp string) (GuideFramePos, bool) {
	switch tp {
	case "Up":
		return GuideFramePos_Up, true
	case "Down":
		return GuideFramePos_Down, true

	}
	return -1, false

}

func GetGuideFramePos_ByAlias(tp string) (GuideFramePos, bool) {
	switch tp {
	case "上方":
		return GuideFramePos_Up, true
	case "下方":
		return GuideFramePos_Down, true

	}
	return -1, false

}

func (e GuideTalkDirection) String() string {

	if e == GuideTalkDirection_Left {
		return "Left"
	}
	if e == GuideTalkDirection_Right {
		return "Right"
	}

	return ""

}

func GetGuideTalkDirection(tp string) (GuideTalkDirection, bool) {
	switch tp {
	case "Left":
		return GuideTalkDirection_Left, true
	case "Right":
		return GuideTalkDirection_Right, true

	}
	return -1, false

}

func GetGuideTalkDirection_ByAlias(tp string) (GuideTalkDirection, bool) {
	switch tp {
	case "左":
		return GuideTalkDirection_Left, true
	case "右":
		return GuideTalkDirection_Right, true

	}
	return -1, false

}

func (e GuideHandAction) String() string {

	if e == GuideHandAction_Click {
		return "Click"
	}
	if e == GuideHandAction_Slide {
		return "Slide"
	}
	if e == GuideHandAction_Circle {
		return "Circle"
	}
	if e == GuideHandAction_Idle {
		return "Idle"
	}
	if e == GuideHandAction_LongPress {
		return "LongPress"
	}

	return ""

}

func GetGuideHandAction(tp string) (GuideHandAction, bool) {
	switch tp {
	case "Click":
		return GuideHandAction_Click, true
	case "Slide":
		return GuideHandAction_Slide, true
	case "Circle":
		return GuideHandAction_Circle, true
	case "Idle":
		return GuideHandAction_Idle, true
	case "LongPress":
		return GuideHandAction_LongPress, true

	}
	return -1, false

}

func GetGuideHandAction_ByAlias(tp string) (GuideHandAction, bool) {
	switch tp {
	case "点击":
		return GuideHandAction_Click, true
	case "滑动":
		return GuideHandAction_Slide, true
	case "画圈":
		return GuideHandAction_Circle, true
	case "放置":
		return GuideHandAction_Idle, true
	case "长按":
		return GuideHandAction_LongPress, true

	}
	return -1, false

}

func (e GuideCustomAction) String() string {

	if e == GuideCustomAction_FirstStartFight {
		return "FirstStartFight"
	}
	if e == GuideCustomAction_RunNext {
		return "RunNext"
	}
	if e == GuideCustomAction_ShowHeroRange {
		return "ShowHeroRange"
	}
	if e == GuideCustomAction_ChangeWave {
		return "ChangeWave"
	}
	if e == GuideCustomAction_Unlock {
		return "Unlock"
	}
	if e == GuideCustomAction_AddItem {
		return "AddItem"
	}
	if e == GuideCustomAction_CloseBottomView {
		return "CloseBottomView"
	}
	if e == GuideCustomAction_HeroFollowBoss {
		return "HeroFollowBoss"
	}
	if e == GuideCustomAction_CheckHeroPos {
		return "CheckHeroPos"
	}
	if e == GuideCustomAction_StopCheckHeroPos {
		return "StopCheckHeroPos"
	}
	if e == GuideCustomAction_AssignSummonResults {
		return "AssignSummonResults"
	}
	if e == GuideCustomAction_StopSpawnEnemy {
		return "StopSpawnEnemy"
	}
	if e == GuideCustomAction_ShakeMap {
		return "ShakeMap"
	}
	if e == GuideCustomAction_HeroAtkMulti {
		return "HeroAtkMulti"
	}
	if e == GuideCustomAction_AssignNextTroop {
		return "AssignNextTroop"
	}
	if e == GuideCustomAction_GetNewHero {
		return "GetNewHero"
	}
	if e == GuideCustomAction_BattlePause {
		return "BattlePause"
	}
	if e == GuideCustomAction_BattleResume {
		return "BattleResume"
	}
	if e == GuideCustomAction_BeginReplaceHero {
		return "BeginReplaceHero"
	}
	if e == GuideCustomAction_EndReplaceHero {
		return "EndReplaceHero"
	}
	if e == GuideCustomAction_ScrollToEntry {
		return "ScrollToEntry"
	}

	return ""

}

func GetGuideCustomAction(tp string) (GuideCustomAction, bool) {
	switch tp {
	case "FirstStartFight":
		return GuideCustomAction_FirstStartFight, true
	case "RunNext":
		return GuideCustomAction_RunNext, true
	case "ShowHeroRange":
		return GuideCustomAction_ShowHeroRange, true
	case "ChangeWave":
		return GuideCustomAction_ChangeWave, true
	case "Unlock":
		return GuideCustomAction_Unlock, true
	case "AddItem":
		return GuideCustomAction_AddItem, true
	case "CloseBottomView":
		return GuideCustomAction_CloseBottomView, true
	case "HeroFollowBoss":
		return GuideCustomAction_HeroFollowBoss, true
	case "CheckHeroPos":
		return GuideCustomAction_CheckHeroPos, true
	case "StopCheckHeroPos":
		return GuideCustomAction_StopCheckHeroPos, true
	case "AssignSummonResults":
		return GuideCustomAction_AssignSummonResults, true
	case "StopSpawnEnemy":
		return GuideCustomAction_StopSpawnEnemy, true
	case "ShakeMap":
		return GuideCustomAction_ShakeMap, true
	case "HeroAtkMulti":
		return GuideCustomAction_HeroAtkMulti, true
	case "AssignNextTroop":
		return GuideCustomAction_AssignNextTroop, true
	case "GetNewHero":
		return GuideCustomAction_GetNewHero, true
	case "BattlePause":
		return GuideCustomAction_BattlePause, true
	case "BattleResume":
		return GuideCustomAction_BattleResume, true
	case "BeginReplaceHero":
		return GuideCustomAction_BeginReplaceHero, true
	case "EndReplaceHero":
		return GuideCustomAction_EndReplaceHero, true
	case "ScrollToEntry":
		return GuideCustomAction_ScrollToEntry, true

	}
	return -1, false

}

func GetGuideCustomAction_ByAlias(tp string) (GuideCustomAction, bool) {
	switch tp {
	case "开始战斗":
		return GuideCustomAction_FirstStartFight, true
	case "执行下一条":
		return GuideCustomAction_RunNext, true
	case "显示英雄范围":
		return GuideCustomAction_ShowHeroRange, true
	case "切换波次":
		return GuideCustomAction_ChangeWave, true
	case "解锁功能":
		return GuideCustomAction_Unlock, true
	case "发放货币":
		return GuideCustomAction_AddItem, true
	case "关闭底部弹窗":
		return GuideCustomAction_CloseBottomView, true
	case "英雄追Boss":
		return GuideCustomAction_HeroFollowBoss, true
	case "检查英雄位置":
		return GuideCustomAction_CheckHeroPos, true
	case "停止检查英雄位置":
		return GuideCustomAction_StopCheckHeroPos, true
	case "指定召唤结果":
		return GuideCustomAction_AssignSummonResults, true
	case "停止刷怪":
		return GuideCustomAction_StopSpawnEnemy, true
	case "地图震动":
		return GuideCustomAction_ShakeMap, true
	case "英雄攻击乘":
		return GuideCustomAction_HeroAtkMulti, true
	case "指定下次阵容":
		return GuideCustomAction_AssignNextTroop, true
	case "获得新英雄":
		return GuideCustomAction_GetNewHero, true
	case "战斗暂停":
		return GuideCustomAction_BattlePause, true
	case "战斗继续":
		return GuideCustomAction_BattleResume, true
	case "开始替换引导":
		return GuideCustomAction_BeginReplaceHero, true
	case "结束替换引导":
		return GuideCustomAction_EndReplaceHero, true
	case "滑动到节点":
		return GuideCustomAction_ScrollToEntry, true

	}
	return -1, false

}

func (e GuideCondiType) String() string {

	if e == GuideCondiType_UIFocus {
		return "UIFocus"
	}
	if e == GuideCondiType_FunctionUnlock {
		return "FunctionUnlock"
	}
	if e == GuideCondiType_TaskFinish {
		return "TaskFinish"
	}
	if e == GuideCondiType_ClickAny {
		return "ClickAny"
	}
	if e == GuideCondiType_ClickTarget {
		return "ClickTarget"
	}
	if e == GuideCondiType_GroupFinish {
		return "GroupFinish"
	}
	if e == GuideCondiType_StepFinish {
		return "StepFinish"
	}
	if e == GuideCondiType_SelfCD {
		return "SelfCD"
	}
	if e == GuideCondiType_GlobalCD {
		return "GlobalCD"
	}
	if e == GuideCondiType_Manual {
		return "Manual"
	}
	if e == GuideCondiType_IdleTime {
		return "IdleTime"
	}
	if e == GuideCondiType_TaskIdLess {
		return "TaskIdLess"
	}
	if e == GuideCondiType_DailyTimes {
		return "DailyTimes"
	}
	if e == GuideCondiType_HeroLevelUp {
		return "HeroLevelUp"
	}
	if e == GuideCondiType_EquipUnlock {
		return "EquipUnlock"
	}
	if e == GuideCondiType_PassMainLevel {
		return "PassMainLevel"
	}
	if e == GuideCondiType_Wait {
		return "Wait"
	}
	if e == GuideCondiType_LevelGreate {
		return "LevelGreate"
	}
	if e == GuideCondiType_GuideBattleStart {
		return "GuideBattleStart"
	}
	if e == GuideCondiType_Coin {
		return "Coin"
	}
	if e == GuideCondiType_CanSummon {
		return "CanSummon"
	}
	if e == GuideCondiType_CanMerge {
		return "CanMerge"
	}
	if e == GuideCondiType_WaveEqual {
		return "WaveEqual"
	}
	if e == GuideCondiType_GameTime {
		return "GameTime"
	}
	if e == GuideCondiType_EnemyCnt {
		return "EnemyCnt"
	}
	if e == GuideCondiType_EnemyClear {
		return "EnemyClear"
	}
	if e == GuideCondiType_GameEnd {
		return "GameEnd"
	}
	if e == GuideCondiType_KillBoss {
		return "KillBoss"
	}
	if e == GuideCondiType_MoveFinished {
		return "MoveFinished"
	}
	if e == GuideCondiType_UnlockFinished {
		return "UnlockFinished"
	}
	if e == GuideCondiType_ShakeMapFinished {
		return "ShakeMapFinished"
	}
	if e == GuideCondiType_GambleSuccess {
		return "GambleSuccess"
	}
	if e == GuideCondiType_BattleTimes {
		return "BattleTimes"
	}
	if e == GuideCondiType_HighestWave {
		return "HighestWave"
	}
	if e == GuideCondiType_HighestWaveLess {
		return "HighestWaveLess"
	}
	if e == GuideCondiType_RoleCreateVersion {
		return "RoleCreateVersion"
	}
	if e == GuideCondiType_RoleCreateVersionLess {
		return "RoleCreateVersionLess"
	}
	if e == GuideCondiType_MaxSp {
		return "MaxSp"
	}
	if e == GuideCondiType_NoChapterUpdate {
		return "NoChapterUpdate"
	}
	if e == GuideCondiType_FirstGetRedHero {
		return "FirstGetRedHero"
	}
	if e == GuideCondiType_RoleCreateVersionGreater {
		return "RoleCreateVersionGreater"
	}
	if e == GuideCondiType_PlayerLevelGreater {
		return "PlayerLevelGreater"
	}
	if e == GuideCondiType_CanSetOffHero {
		return "CanSetOffHero"
	}
	if e == GuideCondiType_PlayerLevelLess {
		return "PlayerLevelLess"
	}
	if e == GuideCondiType_DestinyBind {
		return "DestinyBind"
	}

	return ""

}

func GetGuideCondiType(tp string) (GuideCondiType, bool) {
	switch tp {
	case "UIFocus":
		return GuideCondiType_UIFocus, true
	case "FunctionUnlock":
		return GuideCondiType_FunctionUnlock, true
	case "TaskFinish":
		return GuideCondiType_TaskFinish, true
	case "ClickAny":
		return GuideCondiType_ClickAny, true
	case "ClickTarget":
		return GuideCondiType_ClickTarget, true
	case "GroupFinish":
		return GuideCondiType_GroupFinish, true
	case "StepFinish":
		return GuideCondiType_StepFinish, true
	case "SelfCD":
		return GuideCondiType_SelfCD, true
	case "GlobalCD":
		return GuideCondiType_GlobalCD, true
	case "Manual":
		return GuideCondiType_Manual, true
	case "IdleTime":
		return GuideCondiType_IdleTime, true
	case "TaskIdLess":
		return GuideCondiType_TaskIdLess, true
	case "DailyTimes":
		return GuideCondiType_DailyTimes, true
	case "HeroLevelUp":
		return GuideCondiType_HeroLevelUp, true
	case "EquipUnlock":
		return GuideCondiType_EquipUnlock, true
	case "PassMainLevel":
		return GuideCondiType_PassMainLevel, true
	case "Wait":
		return GuideCondiType_Wait, true
	case "LevelGreate":
		return GuideCondiType_LevelGreate, true
	case "GuideBattleStart":
		return GuideCondiType_GuideBattleStart, true
	case "Coin":
		return GuideCondiType_Coin, true
	case "CanSummon":
		return GuideCondiType_CanSummon, true
	case "CanMerge":
		return GuideCondiType_CanMerge, true
	case "WaveEqual":
		return GuideCondiType_WaveEqual, true
	case "GameTime":
		return GuideCondiType_GameTime, true
	case "EnemyCnt":
		return GuideCondiType_EnemyCnt, true
	case "EnemyClear":
		return GuideCondiType_EnemyClear, true
	case "GameEnd":
		return GuideCondiType_GameEnd, true
	case "KillBoss":
		return GuideCondiType_KillBoss, true
	case "MoveFinished":
		return GuideCondiType_MoveFinished, true
	case "UnlockFinished":
		return GuideCondiType_UnlockFinished, true
	case "ShakeMapFinished":
		return GuideCondiType_ShakeMapFinished, true
	case "GambleSuccess":
		return GuideCondiType_GambleSuccess, true
	case "BattleTimes":
		return GuideCondiType_BattleTimes, true
	case "HighestWave":
		return GuideCondiType_HighestWave, true
	case "HighestWaveLess":
		return GuideCondiType_HighestWaveLess, true
	case "RoleCreateVersion":
		return GuideCondiType_RoleCreateVersion, true
	case "RoleCreateVersionLess":
		return GuideCondiType_RoleCreateVersionLess, true
	case "MaxSp":
		return GuideCondiType_MaxSp, true
	case "NoChapterUpdate":
		return GuideCondiType_NoChapterUpdate, true
	case "FirstGetRedHero":
		return GuideCondiType_FirstGetRedHero, true
	case "RoleCreateVersionGreater":
		return GuideCondiType_RoleCreateVersionGreater, true
	case "PlayerLevelGreater":
		return GuideCondiType_PlayerLevelGreater, true
	case "CanSetOffHero":
		return GuideCondiType_CanSetOffHero, true
	case "PlayerLevelLess":
		return GuideCondiType_PlayerLevelLess, true
	case "DestinyBind":
		return GuideCondiType_DestinyBind, true

	}
	return -1, false

}

func GetGuideCondiType_ByAlias(tp string) (GuideCondiType, bool) {
	switch tp {
	case "界面置顶":
		return GuideCondiType_UIFocus, true
	case "功能解锁":
		return GuideCondiType_FunctionUnlock, true
	case "完成任务":
		return GuideCondiType_TaskFinish, true
	case "点击任意位置":
		return GuideCondiType_ClickAny, true
	case "点击按钮":
		return GuideCondiType_ClickTarget, true
	case "完成引导组":
		return GuideCondiType_GroupFinish, true
	case "完成引导":
		return GuideCondiType_StepFinish, true
	case "间隔时间":
		return GuideCondiType_SelfCD, true
	case "全局间隔时间":
		return GuideCondiType_GlobalCD, true
	case "主动触发":
		return GuideCondiType_Manual, true
	case "无操作时间大于":
		return GuideCondiType_IdleTime, true
	case "任务不大于":
		return GuideCondiType_TaskIdLess, true
	case "每日次数":
		return GuideCondiType_DailyTimes, true
	case "英雄等级":
		return GuideCondiType_HeroLevelUp, true
	case "装备解锁":
		return GuideCondiType_EquipUnlock, true
	case "到达关卡":
		return GuideCondiType_PassMainLevel, true
	case "延迟":
		return GuideCondiType_Wait, true
	case "关卡大于等于":
		return GuideCondiType_LevelGreate, true
	case "开始引导关":
		return GuideCondiType_GuideBattleStart, true
	case "局内金币":
		return GuideCondiType_Coin, true
	case "可召唤":
		return GuideCondiType_CanSummon, true
	case "可合成":
		return GuideCondiType_CanMerge, true
	case "波数等于":
		return GuideCondiType_WaveEqual, true
	case "关卡时间大于":
		return GuideCondiType_GameTime, true
	case "敌人数量":
		return GuideCondiType_EnemyCnt, true
	case "清空敌人":
		return GuideCondiType_EnemyClear, true
	case "关卡结束":
		return GuideCondiType_GameEnd, true
	case "击杀Boss":
		return GuideCondiType_KillBoss, true
	case "移动完成":
		return GuideCondiType_MoveFinished, true
	case "解锁动画完成":
		return GuideCondiType_UnlockFinished, true
	case "地图震动完成":
		return GuideCondiType_ShakeMapFinished, true
	case "抽卡完成":
		return GuideCondiType_GambleSuccess, true
	case "完成挑战次数":
		return GuideCondiType_BattleTimes, true
	case "最高波数":
		return GuideCondiType_HighestWave, true
	case "最高波数小于":
		return GuideCondiType_HighestWaveLess, true
	case "创角版本号":
		return GuideCondiType_RoleCreateVersion, true
	case "创角版本号小于":
		return GuideCondiType_RoleCreateVersionLess, true
	case "能放大招":
		return GuideCondiType_MaxSp, true
	case "没有章节动画":
		return GuideCondiType_NoChapterUpdate, true
	case "有一个红英雄":
		return GuideCondiType_FirstGetRedHero, true
	case "创角版本号大于":
		return GuideCondiType_RoleCreateVersionGreater, true
	case "战队等级大于等于":
		return GuideCondiType_PlayerLevelGreater, true
	case "可上阵英雄":
		return GuideCondiType_CanSetOffHero, true
	case "战队等级小于":
		return GuideCondiType_PlayerLevelLess, true
	case "绑定宿命英雄":
		return GuideCondiType_DestinyBind, true

	}
	return -1, false

}

func (e GuideViewKey) String() string {

	if e == GuideViewKey_None {
		return "None"
	}
	if e == GuideViewKey_BtnHomeTask {
		return "BtnHomeTask"
	}
	if e == GuideViewKey_BtnHomeTaskGuide {
		return "BtnHomeTaskGuide"
	}
	if e == GuideViewKey_Recuit {
		return "Recuit"
	}
	if e == GuideViewKey_BtnRecuitFree {
		return "BtnRecuitFree"
	}
	if e == GuideViewKey_BtnStartFight {
		return "BtnStartFight"
	}
	if e == GuideViewKey_BtnMainLevelOk {
		return "BtnMainLevelOk"
	}
	if e == GuideViewKey_BtnNewHero {
		return "BtnNewHero"
	}
	if e == GuideViewKey_Afk {
		return "Afk"
	}
	if e == GuideViewKey_BtnAfkReward {
		return "BtnAfkReward"
	}
	if e == GuideViewKey_Hero {
		return "Hero"
	}
	if e == GuideViewKey_BtnHeroLevel {
		return "BtnHeroLevel"
	}
	if e == GuideViewKey_BtnHeroClose {
		return "BtnHeroClose"
	}
	if e == GuideViewKey_BtnMainLevel {
		return "BtnMainLevel"
	}
	if e == GuideViewKey_ClickHero {
		return "ClickHero"
	}
	if e == GuideViewKey_BtnHeroEquip {
		return "BtnHeroEquip"
	}
	if e == GuideViewKey_BtnEquipLevelUp {
		return "BtnEquipLevelUp"
	}
	if e == GuideViewKey_Adventure {
		return "Adventure"
	}
	if e == GuideViewKey_ArtsTrial {
		return "ArtsTrial"
	}
	if e == GuideViewKey_TowerChallenge {
		return "TowerChallenge"
	}
	if e == GuideViewKey_Arena {
		return "Arena"
	}
	if e == GuideViewKey_BtnArenaChallenge {
		return "BtnArenaChallenge"
	}
	if e == GuideViewKey_BtnArenaItemChallenge {
		return "BtnArenaItemChallenge"
	}
	if e == GuideViewKey_BtnSummon {
		return "BtnSummon"
	}
	if e == GuideViewKey_GridCanMerge {
		return "GridCanMerge"
	}
	if e == GuideViewKey_BtnMerge {
		return "BtnMerge"
	}
	if e == GuideViewKey_MoveHeroStartGrid {
		return "MoveHeroStartGrid"
	}
	if e == GuideViewKey_MoveHeroEndGrid {
		return "MoveHeroEndGrid"
	}
	if e == GuideViewKey_ItemStone {
		return "ItemStone"
	}
	if e == GuideViewKey_EntryGamble {
		return "EntryGamble"
	}
	if e == GuideViewKey_ItemGamble1 {
		return "ItemGamble1"
	}
	if e == GuideViewKey_EntryUpgrade {
		return "EntryUpgrade"
	}
	if e == GuideViewKey_ItemUpgade2 {
		return "ItemUpgade2"
	}
	if e == GuideViewKey_ItemUpgade4 {
		return "ItemUpgade4"
	}
	if e == GuideViewKey_BtnEpic {
		return "BtnEpic"
	}
	if e == GuideViewKey_NodeEpicNeeds {
		return "NodeEpicNeeds"
	}
	if e == GuideViewKey_BtnMergeEpic {
		return "BtnMergeEpic"
	}
	if e == GuideViewKey_BtnAward {
		return "BtnAward"
	}
	if e == GuideViewKey_BtnPvp {
		return "BtnPvp"
	}
	if e == GuideViewKey_EntryTroop {
		return "EntryTroop"
	}
	if e == GuideViewKey_BtnRefreshTroop {
		return "BtnRefreshTroop"
	}
	if e == GuideViewKey_BtnCloseTroop {
		return "BtnCloseTroop"
	}
	if e == GuideViewKey_BtnAdventureReward {
		return "BtnAdventureReward"
	}
	if e == GuideViewKey_BtnGet {
		return "BtnGet"
	}
	if e == GuideViewKey_BtnMatch {
		return "BtnMatch"
	}
	if e == GuideViewKey_BattleFundCard {
		return "BattleFundCard"
	}
	if e == GuideViewKey_ItemChooseBossAward {
		return "ItemChooseBossAward"
	}
	if e == GuideViewKey_BtnChooseBossAward {
		return "BtnChooseBossAward"
	}
	if e == GuideViewKey_Expedition {
		return "Expedition"
	}
	if e == GuideViewKey_RoyalBoss {
		return "RoyalBoss"
	}
	if e == GuideViewKey_BoxGift {
		return "BoxGift"
	}
	if e == GuideViewKey_BoxGiftBoxBtn {
		return "BoxGiftBoxBtn"
	}
	if e == GuideViewKey_BoxGiftBoxReawrd {
		return "BoxGiftBoxReawrd"
	}
	if e == GuideViewKey_MaxSP {
		return "MaxSP"
	}
	if e == GuideViewKey_BtnArtsTrialFightType1 {
		return "BtnArtsTrialFightType1"
	}
	if e == GuideViewKey_BtnArtsTrialFightItem1 {
		return "BtnArtsTrialFightItem1"
	}
	if e == GuideViewKey_BtnHandUltimate {
		return "BtnHandUltimate"
	}
	if e == GuideViewKey_HomeMenu {
		return "HomeMenu"
	}
	if e == GuideViewKey_DailyMission {
		return "DailyMission"
	}
	if e == GuideViewKey_Manor {
		return "Manor"
	}
	if e == GuideViewKey_ManorRes {
		return "ManorRes"
	}
	if e == GuideViewKey_ManorWork {
		return "ManorWork"
	}
	if e == GuideViewKey_ManorApply {
		return "ManorApply"
	}
	if e == GuideViewKey_ManorSearch {
		return "ManorSearch"
	}
	if e == GuideViewKey_BoxGiftUpgrade {
		return "BoxGiftUpgrade"
	}
	if e == GuideViewKey_TabAstrology {
		return "TabAstrology"
	}
	if e == GuideViewKey_BtnReborn {
		return "BtnReborn"
	}
	if e == GuideViewKey_BtnReplaceOk {
		return "BtnReplaceOk"
	}
	if e == GuideViewKey_NoRedHero {
		return "NoRedHero"
	}
	if e == GuideViewKey_RedHero {
		return "RedHero"
	}
	if e == GuideViewKey_BtnReplaceChoose {
		return "BtnReplaceChoose"
	}
	if e == GuideViewKey_BtnChooseOk {
		return "BtnChooseOk"
	}
	if e == GuideViewKey_BoxGiftPage2 {
		return "BoxGiftPage2"
	}
	if e == GuideViewKey_WorldBoss {
		return "WorldBoss"
	}
	if e == GuideViewKey_Fish {
		return "Fish"
	}
	if e == GuideViewKey_SpTD {
		return "SpTD"
	}
	if e == GuideViewKey_SpRefreshPower {
		return "SpRefreshPower"
	}
	if e == GuideViewKey_SpRefresh {
		return "SpRefresh"
	}
	if e == GuideViewKey_SpLevelUp {
		return "SpLevelUp"
	}
	if e == GuideViewKey_SpHero {
		return "SpHero"
	}
	if e == GuideViewKey_SpNext {
		return "SpNext"
	}
	if e == GuideViewKey_SpTeam {
		return "SpTeam"
	}
	if e == GuideViewKey_Destiny {
		return "Destiny"
	}
	if e == GuideViewKey_DestinyGearTab {
		return "DestinyGearTab"
	}
	if e == GuideViewKey_DestinyGearDefault {
		return "DestinyGearDefault"
	}
	if e == GuideViewKey_DestinyGearGoBtn {
		return "DestinyGearGoBtn"
	}
	if e == GuideViewKey_TeamEntry {
		return "TeamEntry"
	}

	return ""

}

func GetGuideViewKey(tp string) (GuideViewKey, bool) {
	switch tp {
	case "None":
		return GuideViewKey_None, true
	case "BtnHomeTask":
		return GuideViewKey_BtnHomeTask, true
	case "BtnHomeTaskGuide":
		return GuideViewKey_BtnHomeTaskGuide, true
	case "Recuit":
		return GuideViewKey_Recuit, true
	case "BtnRecuitFree":
		return GuideViewKey_BtnRecuitFree, true
	case "BtnStartFight":
		return GuideViewKey_BtnStartFight, true
	case "BtnMainLevelOk":
		return GuideViewKey_BtnMainLevelOk, true
	case "BtnNewHero":
		return GuideViewKey_BtnNewHero, true
	case "Afk":
		return GuideViewKey_Afk, true
	case "BtnAfkReward":
		return GuideViewKey_BtnAfkReward, true
	case "Hero":
		return GuideViewKey_Hero, true
	case "BtnHeroLevel":
		return GuideViewKey_BtnHeroLevel, true
	case "BtnHeroClose":
		return GuideViewKey_BtnHeroClose, true
	case "BtnMainLevel":
		return GuideViewKey_BtnMainLevel, true
	case "ClickHero":
		return GuideViewKey_ClickHero, true
	case "BtnHeroEquip":
		return GuideViewKey_BtnHeroEquip, true
	case "BtnEquipLevelUp":
		return GuideViewKey_BtnEquipLevelUp, true
	case "Adventure":
		return GuideViewKey_Adventure, true
	case "ArtsTrial":
		return GuideViewKey_ArtsTrial, true
	case "TowerChallenge":
		return GuideViewKey_TowerChallenge, true
	case "Arena":
		return GuideViewKey_Arena, true
	case "BtnArenaChallenge":
		return GuideViewKey_BtnArenaChallenge, true
	case "BtnArenaItemChallenge":
		return GuideViewKey_BtnArenaItemChallenge, true
	case "BtnSummon":
		return GuideViewKey_BtnSummon, true
	case "GridCanMerge":
		return GuideViewKey_GridCanMerge, true
	case "BtnMerge":
		return GuideViewKey_BtnMerge, true
	case "MoveHeroStartGrid":
		return GuideViewKey_MoveHeroStartGrid, true
	case "MoveHeroEndGrid":
		return GuideViewKey_MoveHeroEndGrid, true
	case "ItemStone":
		return GuideViewKey_ItemStone, true
	case "EntryGamble":
		return GuideViewKey_EntryGamble, true
	case "ItemGamble1":
		return GuideViewKey_ItemGamble1, true
	case "EntryUpgrade":
		return GuideViewKey_EntryUpgrade, true
	case "ItemUpgade2":
		return GuideViewKey_ItemUpgade2, true
	case "ItemUpgade4":
		return GuideViewKey_ItemUpgade4, true
	case "BtnEpic":
		return GuideViewKey_BtnEpic, true
	case "NodeEpicNeeds":
		return GuideViewKey_NodeEpicNeeds, true
	case "BtnMergeEpic":
		return GuideViewKey_BtnMergeEpic, true
	case "BtnAward":
		return GuideViewKey_BtnAward, true
	case "BtnPvp":
		return GuideViewKey_BtnPvp, true
	case "EntryTroop":
		return GuideViewKey_EntryTroop, true
	case "BtnRefreshTroop":
		return GuideViewKey_BtnRefreshTroop, true
	case "BtnCloseTroop":
		return GuideViewKey_BtnCloseTroop, true
	case "BtnAdventureReward":
		return GuideViewKey_BtnAdventureReward, true
	case "BtnGet":
		return GuideViewKey_BtnGet, true
	case "BtnMatch":
		return GuideViewKey_BtnMatch, true
	case "BattleFundCard":
		return GuideViewKey_BattleFundCard, true
	case "ItemChooseBossAward":
		return GuideViewKey_ItemChooseBossAward, true
	case "BtnChooseBossAward":
		return GuideViewKey_BtnChooseBossAward, true
	case "Expedition":
		return GuideViewKey_Expedition, true
	case "RoyalBoss":
		return GuideViewKey_RoyalBoss, true
	case "BoxGift":
		return GuideViewKey_BoxGift, true
	case "BoxGiftBoxBtn":
		return GuideViewKey_BoxGiftBoxBtn, true
	case "BoxGiftBoxReawrd":
		return GuideViewKey_BoxGiftBoxReawrd, true
	case "MaxSP":
		return GuideViewKey_MaxSP, true
	case "BtnArtsTrialFightType1":
		return GuideViewKey_BtnArtsTrialFightType1, true
	case "BtnArtsTrialFightItem1":
		return GuideViewKey_BtnArtsTrialFightItem1, true
	case "BtnHandUltimate":
		return GuideViewKey_BtnHandUltimate, true
	case "HomeMenu":
		return GuideViewKey_HomeMenu, true
	case "DailyMission":
		return GuideViewKey_DailyMission, true
	case "Manor":
		return GuideViewKey_Manor, true
	case "ManorRes":
		return GuideViewKey_ManorRes, true
	case "ManorWork":
		return GuideViewKey_ManorWork, true
	case "ManorApply":
		return GuideViewKey_ManorApply, true
	case "ManorSearch":
		return GuideViewKey_ManorSearch, true
	case "BoxGiftUpgrade":
		return GuideViewKey_BoxGiftUpgrade, true
	case "TabAstrology":
		return GuideViewKey_TabAstrology, true
	case "BtnReborn":
		return GuideViewKey_BtnReborn, true
	case "BtnReplaceOk":
		return GuideViewKey_BtnReplaceOk, true
	case "NoRedHero":
		return GuideViewKey_NoRedHero, true
	case "RedHero":
		return GuideViewKey_RedHero, true
	case "BtnReplaceChoose":
		return GuideViewKey_BtnReplaceChoose, true
	case "BtnChooseOk":
		return GuideViewKey_BtnChooseOk, true
	case "BoxGiftPage2":
		return GuideViewKey_BoxGiftPage2, true
	case "WorldBoss":
		return GuideViewKey_WorldBoss, true
	case "Fish":
		return GuideViewKey_Fish, true
	case "SpTD":
		return GuideViewKey_SpTD, true
	case "SpRefreshPower":
		return GuideViewKey_SpRefreshPower, true
	case "SpRefresh":
		return GuideViewKey_SpRefresh, true
	case "SpLevelUp":
		return GuideViewKey_SpLevelUp, true
	case "SpHero":
		return GuideViewKey_SpHero, true
	case "SpNext":
		return GuideViewKey_SpNext, true
	case "SpTeam":
		return GuideViewKey_SpTeam, true
	case "Destiny":
		return GuideViewKey_Destiny, true
	case "DestinyGearTab":
		return GuideViewKey_DestinyGearTab, true
	case "DestinyGearDefault":
		return GuideViewKey_DestinyGearDefault, true
	case "DestinyGearGoBtn":
		return GuideViewKey_DestinyGearGoBtn, true
	case "TeamEntry":
		return GuideViewKey_TeamEntry, true

	}
	return -1, false

}

func GetGuideViewKey_ByAlias(tp string) (GuideViewKey, bool) {
	switch tp {
	case "无":
		return GuideViewKey_None, true
	case "任务奖励":
		return GuideViewKey_BtnHomeTask, true
	case "任务引导":
		return GuideViewKey_BtnHomeTaskGuide, true
	case "招募":
		return GuideViewKey_Recuit, true
	case "免费招募":
		return GuideViewKey_BtnRecuitFree, true
	case "开始战斗":
		return GuideViewKey_BtnStartFight, true
	case "主线胜利":
		return GuideViewKey_BtnMainLevelOk, true
	case "新英雄确认":
		return GuideViewKey_BtnNewHero, true
	case "离线挂机":
		return GuideViewKey_Afk, true
	case "离线全部奖励":
		return GuideViewKey_BtnAfkReward, true
	case "英雄入口":
		return GuideViewKey_Hero, true
	case "英雄升级":
		return GuideViewKey_BtnHeroLevel, true
	case "英雄关闭":
		return GuideViewKey_BtnHeroClose, true
	case "主线入口":
		return GuideViewKey_BtnMainLevel, true
	case "上阵英雄":
		return GuideViewKey_ClickHero, true
	case "英雄装备":
		return GuideViewKey_BtnHeroEquip, true
	case "装备强化":
		return GuideViewKey_BtnEquipLevelUp, true
	case "挑战入口":
		return GuideViewKey_Adventure, true
	case "试炼入口":
		return GuideViewKey_ArtsTrial, true
	case "爬塔入口":
		return GuideViewKey_TowerChallenge, true
	case "竞技场入口":
		return GuideViewKey_Arena, true
	case "竞技场挑战入口":
		return GuideViewKey_BtnArenaChallenge, true
	case "竞技场挑战对手":
		return GuideViewKey_BtnArenaItemChallenge, true
	case "召唤按钮":
		return GuideViewKey_BtnSummon, true
	case "可合成格子":
		return GuideViewKey_GridCanMerge, true
	case "合成按钮":
		return GuideViewKey_BtnMerge, true
	case "移动英雄起始格子":
		return GuideViewKey_MoveHeroStartGrid, true
	case "移动英雄目标格子":
		return GuideViewKey_MoveHeroEndGrid, true
	case "幸运石图标":
		return GuideViewKey_ItemStone, true
	case "抽卡入口":
		return GuideViewKey_EntryGamble, true
	case "抽卡选项":
		return GuideViewKey_ItemGamble1, true
	case "升级入口":
		return GuideViewKey_EntryUpgrade, true
	case "紫色升级按钮":
		return GuideViewKey_ItemUpgade2, true
	case "概率升级按钮":
		return GuideViewKey_ItemUpgade4, true
	case "神话按钮":
		return GuideViewKey_BtnEpic, true
	case "神话需求列表":
		return GuideViewKey_NodeEpicNeeds, true
	case "神话合成按钮":
		return GuideViewKey_BtnMergeEpic, true
	case "结算返回按钮":
		return GuideViewKey_BtnAward, true
	case "开始挑战按钮":
		return GuideViewKey_BtnPvp, true
	case "今日阵容按钮":
		return GuideViewKey_EntryTroop, true
	case "阵容刷新按钮":
		return GuideViewKey_BtnRefreshTroop, true
	case "阵容返回按钮":
		return GuideViewKey_BtnCloseTroop, true
	case "挑战奖励按钮":
		return GuideViewKey_BtnAdventureReward, true
	case "领取挑战奖励":
		return GuideViewKey_BtnGet, true
	case "随机组队按钮":
		return GuideViewKey_BtnMatch, true
	case "局内基金卡":
		return GuideViewKey_BattleFundCard, true
	case "Boss奖励选项1":
		return GuideViewKey_ItemChooseBossAward, true
	case "选择Boss奖励按钮":
		return GuideViewKey_BtnChooseBossAward, true
	case "远征入口":
		return GuideViewKey_Expedition, true
	case "大蛇模拟战入口":
		return GuideViewKey_RoyalBoss, true
	case "收藏入口":
		return GuideViewKey_BoxGift, true
	case "收藏箱子":
		return GuideViewKey_BoxGiftBoxBtn, true
	case "收藏领取按钮":
		return GuideViewKey_BoxGiftBoxReawrd, true
	case "大招英雄":
		return GuideViewKey_MaxSP, true
	case "格斗家试炼挑战按钮":
		return GuideViewKey_BtnArtsTrialFightType1, true
	case "格斗家试炼挑战按钮2":
		return GuideViewKey_BtnArtsTrialFightItem1, true
	case "自动战斗":
		return GuideViewKey_BtnHandUltimate, true
	case "菜单":
		return GuideViewKey_HomeMenu, true
	case "每日任务":
		return GuideViewKey_DailyMission, true
	case "领地入口":
		return GuideViewKey_Manor, true
	case "领地资源":
		return GuideViewKey_ManorRes, true
	case "一键上车":
		return GuideViewKey_ManorWork, true
	case "领地确定":
		return GuideViewKey_ManorApply, true
	case "领地搜索":
		return GuideViewKey_ManorSearch, true
	case "收藏升级":
		return GuideViewKey_BoxGiftUpgrade, true
	case "神器祈福":
		return GuideViewKey_TabAstrology, true
	case "重生按钮":
		return GuideViewKey_BtnReborn, true
	case "英雄替换确定":
		return GuideViewKey_BtnReplaceOk, true
	case "非红英雄":
		return GuideViewKey_NoRedHero, true
	case "红英雄":
		return GuideViewKey_RedHero, true
	case "选择替换英雄":
		return GuideViewKey_BtnReplaceChoose, true
	case "英雄选择确定":
		return GuideViewKey_BtnChooseOk, true
	case "收藏页签2":
		return GuideViewKey_BoxGiftPage2, true
	case "世界boss入口":
		return GuideViewKey_WorldBoss, true
	case "钓鱼入口":
		return GuideViewKey_Fish, true
	case "南镇风云入口":
		return GuideViewKey_SpTD, true
	case "南镇刷新能力":
		return GuideViewKey_SpRefreshPower, true
	case "南镇刷新":
		return GuideViewKey_SpRefresh, true
	case "南镇升级":
		return GuideViewKey_SpLevelUp, true
	case "南镇英雄":
		return GuideViewKey_SpHero, true
	case "南镇下一关":
		return GuideViewKey_SpNext, true
	case "南镇阵容编辑":
		return GuideViewKey_SpTeam, true
	case "宿命入口":
		return GuideViewKey_Destiny, true
	case "宿命齿轮入口":
		return GuideViewKey_DestinyGearTab, true
	case "宿命齿轮默认英雄":
		return GuideViewKey_DestinyGearDefault, true
	case "宿命齿轮前往按钮":
		return GuideViewKey_DestinyGearGoBtn, true
	case "队伍试炼入口":
		return GuideViewKey_TeamEntry, true

	}
	return -1, false

}
