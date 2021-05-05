package define

const (
	Event_Type_Begin         int32 = iota
	Event_Type_Null          int32 = iota - 1 // 0 无
	Event_Type_Sign                           // 签到
	Event_Type_PlayerLevelup                  // 角色升级
	Event_Type_HeroLevelup                    // 英雄升级
	Event_Type_HeroGain                       // 获得英雄
	Event_Type_StagePass                      // 关卡通关
	Event_Type_ChapterPass                    // 章节通关

	Event_Type_End
)
