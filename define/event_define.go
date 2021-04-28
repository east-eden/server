package define

const (
	Event_Type_Begin         int32 = iota
	Event_Type_Sign          int32 = iota - 1 // 0 签到
	Event_Type_PlayerLevelup                  // 1 角色升级
	Event_Type_HeroLevelup                    // 2 英雄升级
	Event_Type_End
)
