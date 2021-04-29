package define

// 任务类型
const (
	Quest_Type_Begin    int32 = iota
	Quest_Type_Main     int32 = iota - 1 // 0 主线任务
	Quest_Type_Daily                     // 1 日常任务
	Quest_Type_Weekly                    // 2 周常任务
	Quest_Type_Activity                  // 3 活动任务
	Quest_Type_End
)

// 任务刷新方式
const (
	Quest_Refresh_Type_Begin  int32 = iota
	Quest_Refresh_Type_None   int32 = iota - 1 // 0 不刷新
	Quest_Refresh_Type_Daily                   // 1 跨天刷新
	Quest_Refresh_Type_Weekly                  // 2 跨周刷新
	Quest_Refresh_Type_End
)

// 任务状态
const (
	Quest_State_Type_Begin     int32 = iota
	Quest_State_Type_Opened    int32 = iota - 1 // 0 任务开启
	Quest_State_Type_Completed                  // 1 任务完成
	Quest_State_Type_Rewarded                   // 2 任务已领取奖励
	Quest_State_Type_End
)

// 任务目标类型
const (
	QuestObj_Type_Begin          int32 = iota
	QuestObj_Type_StagePass      int32 = iota - 1 // 0 通关关卡
	QuestObj_Type_ChapterPass                     // 1 通关章节
	QuestObj_Type_PlayerLevel                     // 2 队伍等级达到XX级
	QuestObj_Type_HeroLevelTimes                  // 3 英雄升级XX次
	QuestObj_Type_GainHero                        // 4 获得英雄
	QuestObj_Type_GainHeroNum                     // 5 获得英雄XX个
	QuestObj_Type_End
)

// 任务拥有者类型
const (
	QuestOwner_Type_Begin      int32 = iota
	QuestOwner_Type_Player     int32 = iota - 1 // 0 玩家任务
	QuestOwner_Type_Collection                  // 1 收集品任务
	QuestOwner_End
)
